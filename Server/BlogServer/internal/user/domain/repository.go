package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(c context.Context, user *User) (*User, error)
	CountByUsername(c context.Context, username string) (int64, error)
	GetUserByUsername(c context.Context, username string) (*User, error)
	UpdateLastLogout(c context.Context, userID uuid.UUID) error
	GetUserByID(c context.Context, userID uuid.UUID) (*User, error)
	UpdateEmail(c context.Context, userID uuid.UUID, email string, updatedBy *uuid.UUID) error
	UpdateData(c context.Context, userID uuid.UUID, user *User, updatedBy *uuid.UUID) error
	UpdatePassword(c context.Context, userID uuid.UUID, hashedPassword string, updatedBy *uuid.UUID) error
	GetNotificationsByUserID(c context.Context, userID uuid.UUID) ([]Notification, error)
	CreateNotification(c context.Context, content []byte, userID uuid.UUID, createdBy uuid.UUID) (*Notification, error)
	CreateNotifications(c context.Context, nots []Notification) error
	UpdateNotificationByID(c context.Context, notificationID int64, status bool, updatedBy *uuid.UUID) error
	MarkUserAsDeleted(c context.Context, userID uuid.UUID, updatedBy *uuid.UUID) error
}
