package application

import (
	"context"

	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCases struct {
	txManager  database.TxManager
	repo       domain.UserRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewAuthUseCases(txManager database.TxManager, repo domain.UserRepository, outboxRepo outboxrepo.OutboxRepository) *AuthUseCases {
	return &AuthUseCases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (u *AuthUseCases) RegisterUser(c context.Context, user *domain.User) (*uuid.UUID, error) {

	var userID *uuid.UUID
	// return newUserIdString, err
	err := u.txManager.WithVoidTx(c, func(ctx context.Context) error {
		count, err := u.repo.CountByUsername(c, user.Username)
		if err != nil {
			return err
		}

		if count > 0 {
			return &domain.ErrUsernameAlreadyTaken{}
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return &domain.ErrFailedToHashString{}
		}
		user.Password = hashedPassword
		newUser, err := u.repo.Create(c, user)
		if err != nil {
			return err
		}
		userID = &newUser.UserID
		return nil
	})

	if err != nil {
		return nil, err
	}

	return userID, nil
}

func (u *AuthUseCases) CheckExistedUsername(c context.Context, username string) (int64, error) {
	return u.repo.CountByUsername(c, username)
}

func (u *AuthUseCases) LoginUser(c context.Context, username string, password string) (*domain.User, error) {

	foundUser, err := u.repo.GetUserByUsername(c, username)
	if err != nil {
		return nil, &domain.ErrNotFound{}
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if compareErr != nil {
		return nil, &domain.ErrPasswordNotMatched{}
	}

	return foundUser, nil
}

func (u *AuthUseCases) LogoutUser(c context.Context, userID uuid.UUID) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {
		return u.repo.UpdateLastLogout(c, userID)
	})
}

func (s *AuthUseCases) GetUserByID(c context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.repo.GetUserByID(c, userID)
}
