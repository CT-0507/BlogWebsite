package user

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	userdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/db"
)

func UserDTOToUser(userDTO *userdb.UsersUser) *User {
	return &User{
		UserID:       userDTO.UserID.String(),
		Username:     userDTO.Username,
		Password:     userDTO.Password,
		Email:        userDTO.Email.String,
		FirstName:    userDTO.FirstName,
		LastName:     userDTO.LastName,
		Roles:        []string{userDTO.Role},
		TokenVersion: int(userDTO.TokenVersion.Int32),
		Active:       userDTO.Active.String,
		Audit: model.Audit{
			CreatedAt: userDTO.CreatedAt.Time,
			CreatedBy: utils.UUIDPtr(userDTO.CreatedBy),
			UpdatedAt: userDTO.UpdatedAt.Time,
			UpdatedBy: utils.UUIDPtr(userDTO.UpdatedBy),
		},
	}
}

func NotificationDTOToNotification(notDTO *userdb.UsersNotification) *Notification {
	return &Notification{
		NotificationID: notDTO.NotificationID,
		UserID:         notDTO.UserID.String(),
		Content:        notDTO.Content,
		IsRead:         notDTO.IsRead,

		Audit: model.Audit{
			CreatedAt: notDTO.CreatedAt.Time,
		},
	}
}

// func UserDTOToUserLogin(user *userdb.UsersUser) *User {
// 	return &User{
// 		UserID:    userDTO.UserID.String(),
// 		Username:  userDTO.Username,
// 		Email:     userDTO.Email.String,
// 		FirstName: userDTO.FirstName,
// 		LastName:  userDTO.LastName,
// 		Role:      userDTO.Role,
// 		Active:    userDTO.Active.Bool,
// 		Audit: model.Audit{
// 			CreatedAt: userDTO.CreatedAt.Time,
// 			CreatedBy: utils.UUIDPtr(userDTO.CreatedBy),
// 			UpdatedAt: userDTO.UpdatedAt.Time,
// 			UpdatedBy: utils.UUIDPtr(userDTO.UpdatedBy),
// 		},
// 	}
// }
