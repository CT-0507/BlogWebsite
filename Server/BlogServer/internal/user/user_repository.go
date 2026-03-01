package user

import (
	"context"

	userdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	Create(c context.Context, q *userdb.Queries, user *User) (*User, error)
	CountByUsername(c context.Context, q *userdb.Queries, username string) (int64, error)
	GetUserByUsername(c context.Context, q *userdb.Queries, username string) (*User, error)
	UpdateLastLogout(c context.Context, q *userdb.Queries, userID uuid.UUID) error
	GetUserByID(c context.Context, q *userdb.Queries, userID uuid.UUID) (*User, error)
	UpdateEmail(c context.Context, q *userdb.Queries, userID uuid.UUID, email string) error
	UpdateData(c context.Context, q *userdb.Queries, userID uuid.UUID, user *User) error
	UpdatePassword(c context.Context, q *userdb.Queries, userID uuid.UUID, hashedPassword string) error
	// Delete(c context.Context, q *userdb.Queries, id int64) (*int64, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(c context.Context, q *userdb.Queries, user *User) (*User, error) {

	newUser, err := q.CreateUser(c, userdb.CreateUserParams{
		Username:  user.Username,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Roles[0],
	})
	if err != nil {
		return nil, err
	}

	return UserDTOToUser(&newUser), nil
}

func (r *userRepository) CountByUsername(c context.Context, q *userdb.Queries, username string) (int64, error) {
	return q.CountUserWithEmail(c, username)
}

func (r *userRepository) GetUserByUsername(c context.Context, q *userdb.Queries, username string) (*User, error) {
	user, err := q.GetUserByUsername(c, username)
	if err != nil {
		return nil, err
	}
	return UserDTOToUser(&user), nil
}

func (r *userRepository) UpdateLastLogout(c context.Context, q *userdb.Queries, userID uuid.UUID) error {
	return q.UpdateLastLogout(c, userID)
}

func (r *userRepository) GetUserByID(c context.Context, q *userdb.Queries, userID uuid.UUID) (*User, error) {
	user, err := q.GetUserByID(c, userID)
	if err != nil {
		return nil, err
	}
	return UserDTOToUser(&user), nil
}

func (r *userRepository) UpdateEmail(c context.Context, q *userdb.Queries, userID uuid.UUID, email string) error {
	_, err := q.UpdateUserEmail(c, userdb.UpdateUserEmailParams{
		UserID: userID,
		Email: pgtype.Text{
			String: email,
			Valid:  true,
		},
	})
	return err
}

func (r *userRepository) UpdateData(c context.Context, q *userdb.Queries, userID uuid.UUID, user *User) error {
	_, err := q.UpdateUserData(c, userdb.UpdateUserDataParams{
		UserID:    userID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
	return err
}

func (r *userRepository) UpdatePassword(c context.Context, q *userdb.Queries, userID uuid.UUID, hashedPassword string) error {
	_, err := q.UpdateUserPassword(c, userdb.UpdateUserPasswordParams{
		UserID:   userID,
		Password: hashedPassword,
	})
	return err
}
