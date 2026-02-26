package user

import (
	"context"

	userdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/db"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(c context.Context, q *userdb.Queries, user *User) (*User, error)
	CountByUsername(c context.Context, q *userdb.Queries, username string) (int64, error)
	GetUserByUsername(c context.Context, q *userdb.Queries, username string) (*User, error)
	UpdateLastLogout(c context.Context, q *userdb.Queries, userID uuid.UUID) error
	// FindAll(c context.Context, q *userdb.Queries) ([]User, error)
	// FindByID(c context.Context, q *userdb.Queries, id uuid.UUID) (*User, error)
	// Update(user *User, q *userdb.Queries) error
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
		Role:      user.Role,
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

// func (r *userRepository) FindByID(c context.Context, q *userdb.Queries, id uuid.UUID) (*User, error) {

// 	newUser, err := q.GetUserByID()(c, userdb.CreateUserParams{
// 		Email:     user.Email,
// 		Password:  user.Password,
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return UserDTOToBlog(&newUser), nil
// }
