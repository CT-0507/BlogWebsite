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
	Roles        []string    `json:"roles"`
	Email        string      `json:"email"`
	Active       string      `json:"active"`
	ProfileData  ProfileData `json:"profileData"`
	Token        string      `json:"token"`
	TokenVersion int         `json:"token_version"`
	model.Audit
}

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
	UserID      string      `json:"userID"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	Roles       []string    `json:"roles"`
	Email       string      `json:"email"`
	Active      string      `json:"active"`
	ProfileData ProfileData `json:"profileData"`
	AccessToken string      `json:"accessToken"`
}

type ErrUsernameAlreadyTaken struct{}

func (e *ErrUsernameAlreadyTaken) Error() string {
	return "User already exists"
}
