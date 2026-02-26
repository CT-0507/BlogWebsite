package user

import (
	"context"
	"errors"

	userdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Create(c context.Context, user *User) (string, error)
	CheckExistedUsername(c context.Context, username string) (int64, error)
	RegisterUser(c context.Context, user *User) (string, error)
	LoginUser(c context.Context, username string) (*User, error)
	LogoutUser(c context.Context, userID uuid.UUID) error
	// GetAll(c context.Context) ([]User, error)
	// GetByID(c context.Context, id int64) (*User, error)
	// Update(user *User) error
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

func (s *userService) LoginUser(c context.Context, username string) (*User, error) {
	q := userdb.New(s.pool)
	return s.repo.GetUserByUsername(c, q, username)
}

func (s *userService) LogoutUser(c context.Context, userID uuid.UUID) error {
	return s.withTxExec(c, func(q *userdb.Queries) error {
		return s.repo.UpdateLastLogout(c, q, userID)
	})
}
