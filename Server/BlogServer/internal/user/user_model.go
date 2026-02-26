package user

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"

type ProfileData struct {
}

type User struct {
	UserID       string      `json:"user_id"`
	Username     string      `json:"username"`
	Password     string      `json:"password"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Role         string      `json:"role"`
	Email        string      `json:"email"`
	Active       string      `json:"active"`
	ProfileData  ProfileData `json:"profileData"`
	Token        string      `json:"token"`
	TokenVersion int         `json:"token_version"`
	RefreshToken string      `json:"refresh_token"`
	model.Audit
}

type CreateUserRequest struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=20"`
	FirstName       string `json:"first_name" validate:"required,min=2,max=20"`
	LastName        string `json:"last_name" validate:"required,min=2,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type UserLoginResponse struct {
	UserID      string      `json:"user_id"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Role        string      `json:"role"`
	Email       string      `json:"email"`
	Active      string      `json:"active"`
	ProfileData ProfileData `json:"profileData"`
	Token       string      `json:"token"`
}

type ErrUsernameAlreadyTaken struct{}

func (e *ErrUsernameAlreadyTaken) Error() string {
	return "User already exists"
}
