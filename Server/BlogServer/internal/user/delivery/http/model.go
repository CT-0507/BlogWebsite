package http

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/user/domain"

type CreateUserRequest struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=20"`
	FirstName       string `json:"firstName" validate:"required,min=1,max=20"`
	LastName        string `json:"lastName" validate:"required,min=1,max=20"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type UserLoginResponse struct {
	UserID      string             `json:"userID"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Roles       []string           `json:"roles"`
	Email       string             `json:"email"`
	Active      string             `json:"active"`
	ProfileData domain.ProfileData `json:"profileData"`
	AccessToken string             `json:"accessToken"`
}

type UpdateUserEmailRequest struct {
	Email       string `json:"email" validate:"required,email,max=50"`
	ConfirmCode string `json:"confirmCode" validate:"required,max=6"`
}

type UpdateUserBasicInfoRequest struct {
	FirstName string `json:"firstName" validate:"required,min=1,max=20"`
	LastName  string `json:"lastName" validate:"required,min=1,max=20"`
}

type UpdatePasswordRequest struct {
	CurrentPassword    string `json:"currentPassword" validate:"required,min=8,max=20"`
	NewPassword        string `json:"newPassword" validate:"required,min=8,max=20,nefield=CurrentPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required,min=8,max=20,eqfield=NewPassword"`
}

type UpdatePasswordServiceParams struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type UpdateNotificationStatusRequest struct {
	NotId  int64 `json:"notificationID"`
	Status bool  `json:"status"`
}
