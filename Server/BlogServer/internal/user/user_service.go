package user

import (
	"context"
	"errors"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/config"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	userdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(c context.Context, user *User) (string, error)
	CheckExistedUsername(c context.Context, username string) (int64, error)
	RegisterUser(c context.Context, user *User) (string, error)
	LoginUser(c context.Context, username string, password string) (*User, error)
	LogoutUser(c context.Context, userID uuid.UUID) error
	// GetAll(c context.Context) ([]User, error)
	GetUserByID(c context.Context, userID uuid.UUID) (*User, error)
	UpdateEmail(c context.Context, userID uuid.UUID, email string) error
	UpdatePassword(c context.Context, userID uuid.UUID, userPassword *UpdatePasswordServiceParams) error
	UpdateBasicInfo(c context.Context, userID uuid.UUID, user *User) error
	GetUserNotifications(c context.Context, userID uuid.UUID) ([]Notification, error)
	CreateNotification(c context.Context, content string, userID uuid.UUID, createdBy uuid.UUID) (*Notification, error)
	UpdateNotificationStatus(c context.Context, notID int64, status bool, updatedBy *uuid.UUID) error
	GetHashedString(str string) (string, error)
	// Delete(c context.Context, id int64) (*int64, error)
}

type userService struct {
	pool *pgxpool.Pool
	repo UserRepository
}

func NewUserService(pool *pgxpool.Pool, repo UserRepository) UserService {
	return &userService{
		pool: pool,
		repo: repo,
	}
}

func (s *userService) withTx(
	ctx context.Context,
	fn func(q *userdb.Queries) (any, error),
) (any, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	queries := userdb.New(tx)
	result, err := fn(queries)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *userService) withTxExec(
	ctx context.Context,
	fn func(q *userdb.Queries) error,
) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	queries := userdb.New(tx)

	if err := fn(queries); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *userService) Create(c context.Context, user *User) (string, error) {
	newUserId, err := s.withTx(c, func(q *userdb.Queries) (any, error) {
		newUser, err := s.repo.Create(c, q, user)
		return newUser.UserID, err
	})
	return newUserId.(string), err
}

func (s *userService) RegisterUser(c context.Context, user *User) (string, error) {
	newUserId, err := s.withTx(c, func(q *userdb.Queries) (any, error) {
		count, err := s.repo.CountByUsername(c, q, user.Username)
		if err != nil {
			return "", err
		}

		if count > 0 {
			return "", &ErrUsernameAlreadyTaken{}
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return "", &ErrFailedToHashString{}
		}
		user.Password = hashedPassword
		newUser, err := s.repo.Create(c, q, user)
		if err != nil {
			return "", err
		}
		return newUser.UserID, nil
	})
	if err != nil {
		return "", err
	}
	newUserIdString, ok := newUserId.(string)
	if !ok {
		return "", errors.New("Cannot convert")
	}
	return newUserIdString, err
}

func (s *userService) CheckExistedUsername(c context.Context, username string) (int64, error) {
	q := userdb.New(s.pool)
	return s.repo.CountByUsername(c, q, username)
}

func (s *userService) LoginUser(c context.Context, username string, password string) (*User, error) {

	q := userdb.New(s.pool)
	foundUser, err := s.repo.GetUserByUsername(c, q, username)
	if err != nil {
		return nil, &ErrNotFound{}
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if compareErr != nil {
		return nil, &ErrPasswordNotMatched{}
	}

	return foundUser, nil
}

func (s *userService) LogoutUser(c context.Context, userID uuid.UUID) error {
	return s.withTxExec(c, func(q *userdb.Queries) error {
		return s.repo.UpdateLastLogout(c, q, userID)
	})
}

func (s *userService) GetUserByID(c context.Context, userID uuid.UUID) (*User, error) {
	q := userdb.New(s.pool)
	return s.repo.GetUserByID(c, q, userID)
}

func (s *userService) UpdateEmail(c context.Context, userID uuid.UUID, email string) error {
	return s.withTxExec(c, func(q *userdb.Queries) error {
		return s.repo.UpdateEmail(c, q, userID, email, &userID)
	})
}

func (s *userService) UpdateBasicInfo(c context.Context, userID uuid.UUID, user *User) error {
	return s.withTxExec(c, func(q *userdb.Queries) error {
		return s.repo.UpdateData(c, q, userID, user, &userID)
	})
}

func (s *userService) UpdatePassword(c context.Context, userID uuid.UUID, userPassword *UpdatePasswordServiceParams) error {
	return s.withTxExec(c, func(q *userdb.Queries) error {

		foundUser, err := s.repo.GetUserByID(c, q, userID)
		if err != nil {
			return err
		}
		compareErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userPassword.CurrentPassword))
		if compareErr != nil {
			return &ErrPasswordNotMatched{}
		}

		hashedNewPassword, err := utils.HashPassword(userPassword.NewPassword)
		if err != nil {
			return &ErrFailedToHashString{}
		}

		return s.repo.UpdatePassword(c, q, userID, hashedNewPassword, &userID)
	})
}

func (s *userService) GetUserNotifications(c context.Context, userID uuid.UUID) ([]Notification, error) {
	q := userdb.New(s.pool)
	return s.repo.GetNotificationsByUserID(c, q, userID)
}

func (s *userService) UpdateNotificationStatus(c context.Context, notID int64, status bool, updatedBy *uuid.UUID) error {
	q := userdb.New(s.pool)
	return s.repo.UpdateNotificationByID(c, q, notID, status, updatedBy)
}

func (s *userService) CreateNotification(c context.Context, content string, userID uuid.UUID, createdBy uuid.UUID) (*Notification, error) {
	systemId := uuid.MustParse(config.SYSTEM_ID)
	not, err := s.withTx(c, func(q *userdb.Queries) (any, error) {
		return s.repo.CreateNotification(c, q, content, userID, systemId)
	})
	if err != nil {
		return nil, err
	}
	return not.(*Notification), nil
}

// func (s *userService) UpdateNotificationsStatusByIds(c context.Context, not *Notification, userID uuid.UUID, createdBy uuid.UUID) error {
// 	systemId := uuid.MustParse(config.SYSTEM_ID)
// 	return s.withTxExec(c, func(q *userdb.Queries) error {
// 		return s.repo.UpdateNotificationByID(c, q, not, userID, &userID)
// 	})
// }

func (s *userService) GetHashedString(str string) (string, error) {
	return utils.HashPassword(str)
}
