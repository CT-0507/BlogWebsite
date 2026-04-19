package infrastructure

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/domain"
	userdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/infrastructure/db"
)

func UserDTOToUser(userDTO *userdb.UsersUser) *domain.User {
	return &domain.User{
		UserID:       userDTO.UserID,
		Username:     userDTO.Username,
		Password:     userDTO.Password,
		Email:        userDTO.Email.String,
		FirstName:    userDTO.FirstName,
		LastName:     userDTO.LastName,
		Roles:        []string{userDTO.Role},
		TokenVersion: int(userDTO.TokenVersion.Int32),
		Status:       userDTO.Status.String,
		Audit: model.Audit{
			CreatedAt: userDTO.CreatedAt.Time,
			CreatedBy: utils.UUIDPtr(userDTO.CreatedBy),
			UpdatedAt: userDTO.UpdatedAt.Time,
			UpdatedBy: utils.UUIDPtr(userDTO.UpdatedBy),
		},
	}
}

func NotificationDTOToNotification(notDTO *userdb.UsersNotification) *domain.Notification {
	return &domain.Notification{
		NotificationID: notDTO.NotificationID,
		UserID:         "",
		Content:        notDTO.Content,
		IsRead:         notDTO.IsRead,

		Audit: model.Audit{
			CreatedAt: notDTO.CreatedAt.Time,
		},
	}
}
