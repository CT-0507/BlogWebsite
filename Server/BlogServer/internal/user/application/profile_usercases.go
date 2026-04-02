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

type ProfileUseCases struct {
	txManager  database.TxManager
	repo       domain.UserRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewProfileUseCases(txManager database.TxManager, repo domain.UserRepository, outboxRepo outboxrepo.OutboxRepository) *ProfileUseCases {
	return &ProfileUseCases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (u *ProfileUseCases) UpdateEmail(c context.Context, userID uuid.UUID, email string) error {
	return u.repo.UpdateEmail(c, userID, email, &userID)
}

func (u *ProfileUseCases) UpdateBasicInfo(c context.Context, userID uuid.UUID, user *domain.User) error {
	return u.repo.UpdateData(c, userID, user, &userID)
}

func (u *ProfileUseCases) UpdatePassword(c context.Context, userID uuid.UUID, currentPassword string, newPassword string) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		foundUser, err := u.repo.GetUserByID(c, userID)
		if err != nil {
			return err
		}
		compareErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(currentPassword))
		if compareErr != nil {
			return &domain.ErrPasswordNotMatched{}
		}

		hashedNewPassword, err := utils.HashPassword(newPassword)
		if err != nil {
			return &domain.ErrFailedToHashString{}
		}

		return u.repo.UpdatePassword(c, userID, hashedNewPassword, &userID)
	})
}
