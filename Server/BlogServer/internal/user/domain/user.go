package domain

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/google/uuid"
)

type ProfileData struct {
}

type Notification struct {
	NotificationID int64  `json:"notificationId"`
	UserID         string `json:"userId"`
	Content        []byte `json:"content"`
	IsRead         bool   `json:"isRead"`

	model.Audit
}

type User struct {
	UserID       uuid.UUID   `json:"user_id"`
	Username     string      `json:"username"`
	Password     string      `json:"password"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Roles        []string    `json:"roles"`
	Email        string      `json:"email"`
	Status       string      `json:"status"`
	ProfileData  ProfileData `json:"profileData"`
	Token        string      `json:"token"`
	TokenVersion int         `json:"token_version"`
	model.Audit
}

type ErrUsernameAlreadyTaken struct{}

func (e *ErrUsernameAlreadyTaken) Error() string {
	return "User already exists"
}

type ErrPasswordNotMatched struct{}

func (e *ErrPasswordNotMatched) Error() string {
	return "Current password does not matched"
}

type ErrFailedToHashString struct{}

func (e *ErrFailedToHashString) Error() string {
	return "Failed to hash string"
}

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "Target not found"
}
