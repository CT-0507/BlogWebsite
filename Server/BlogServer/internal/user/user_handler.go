package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/messages"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Description: create new blog
//   - @route POST /blogs
//   - @access Private
func (h *UserHandler) registerUser(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var user CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if isValidPassword := utils.IsValidPassword(user.Password); !isValidPassword {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "Password is not valid"})
		return
	}

	newUserId, err := h.service.RegisterUser(ctx, &User{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Roles:     []string{"admin"},
	})
	if err != nil {
		switch {
		case errors.Is(err, &ErrUsernameAlreadyTaken{}):
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		case errors.Is(err, &ErrFailedToHashString{}):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, newUserId)
}

func (h *UserHandler) loginUser(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var user UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if isValidPassword := utils.IsValidPassword(user.Password); !isValidPassword {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "Password is not valid"})
		return
	}

	foundUser, loginErr := h.service.LoginUser(ctx, user.Username, user.Password)
	if loginErr != nil {
		switch {
		case errors.Is(loginErr, &ErrNotFound{}):
		case errors.Is(loginErr, &ErrPasswordNotMatched{}):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password is invalid"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": loginErr.Error()})
		}
		return
	}

	token, refreshToken, err := utils.GenerateAllTokens(
		foundUser.Username, foundUser.FirstName, foundUser.LastName, foundUser.UserID, foundUser.Roles, foundUser.TokenVersion,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refresh_token", // name
		refreshToken,    // value
		60*60*24*30,     // maxAge (seconds) → 30 days
		"/",             // path
		"",              // domain ("" for current)
		true,            // secure (HTTPS only)
		true,            // httpOnly (no JS access)
	)

	c.JSON(http.StatusOK, &UserLoginResponse{
		UserID:      foundUser.UserID,
		FirstName:   foundUser.FirstName,
		LastName:    foundUser.LastName,
		Email:       foundUser.Email,
		Roles:       foundUser.Roles,
		Active:      foundUser.Active,
		AccessToken: token,
	})
}

func (h *UserHandler) logout(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refresh_token",
		"",
		-1, // delete
		"/",
		"",
		true,
		true,
	)

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unable to logout",
		})
		return
	}

	// Update last logout
	if h.service.LogoutUser(ctx, userID) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logout",
	})
}

func (h *UserHandler) getUserById(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unable to find user id",
		})
		return
	}

	foundUser, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusCreated, &UserLoginResponse{
		UserID:    foundUser.UserID,
		FirstName: foundUser.FirstName,
		LastName:  foundUser.LastName,
		Email:     foundUser.Email,
		Roles:     foundUser.Roles,
		Active:    foundUser.Active,
	})
}

func (h *UserHandler) RefreshTokenHandler(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		fmt.Println("error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve refresh token from cookie"})
		return
	}

	claim, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil || claim == nil {
		fmt.Println("error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	userID, err := uuid.Parse(claim.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating tokens"})
		return
	}

	foundUser, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	newToken, newRefreshToken, err := utils.GenerateAllTokens(foundUser.Username, foundUser.FirstName, foundUser.LastName, foundUser.UserID, foundUser.Roles, foundUser.TokenVersion)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "localhost", true, true) //expires in 1 week

	c.JSON(http.StatusOK, gin.H{"accessToken": newToken})
}

func (h *UserHandler) UpdateUserBasicInfo(c *gin.Context) {

	var userInfo UpdateUserBasicInfoRequest
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, userInfo); err != nil {
		c.JSON(http.StatusBadRequest, &err)
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unable to logout",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := h.service.UpdateBasicInfo(ctx, userID, &User{
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully change data",
		"firstName": userInfo.FirstName,
		"lastName":  userInfo.LastName})
}

func (h *UserHandler) UpdateUserEmail(c *gin.Context) {

	var userEmail UpdateUserEmailRequest
	if err := c.ShouldBindJSON(&userEmail); err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, userEmail); err != nil {
		c.JSON(http.StatusBadRequest, &err)
		return
	}

	if userEmail.ConfirmCode != "123456" {
		c.JSON(http.StatusForbidden, &gin.H{
			"message": "Confirm code is incorrect",
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unable to logout",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := h.service.UpdateEmail(ctx, userID, userEmail.Email); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &gin.H{"email": userEmail.Email})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {

	var userPassword UpdatePasswordRequest
	if err := c.ShouldBindJSON(&userPassword); err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(messages.ENGLISH, userPassword); err != nil {
		c.JSON(http.StatusBadRequest, &err)
		return
	}

	if isValidPassword := utils.IsValidPassword(userPassword.CurrentPassword); !isValidPassword {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "Current password is not valid"})
		return
	}

	if isValidPassword := utils.IsValidPassword(userPassword.NewPassword); !isValidPassword {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "New password is not valid"})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unable to logout",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := h.service.UpdatePassword(ctx, userID, &UpdatePasswordServiceParams{
		CurrentPassword: userPassword.CurrentPassword,
		NewPassword:     userPassword.NewPassword,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &gin.H{"message": "OK"})
}

func (h *UserHandler) GetChangeEmailCode(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{"code": "123456"})
}

func (h *UserHandler) GetNotifications(c *gin.Context) {

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{"message": "No user found"})
		return
	}

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	notifications, err := h.service.GetUserNotifications(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *UserHandler) UpdateNotification(c *gin.Context) {

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{"message": "No user found"})
		return
	}

	var requestJson UpdateNotificationStatusRequest
	if err := c.ShouldBindBodyWithJSON(&requestJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	if err := h.service.UpdateNotificationStatus(ctx, requestJson.NotId, requestJson.Status, &userID); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"message": "Success",
	})

}

func (h *UserHandler) GetHashedString(c *gin.Context) {

	str := c.Query("string")
	if str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string param is required"})
		return
	}

	if isValidPassword := utils.IsValidPassword(str); !isValidPassword {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "Password is not valid"})
		return
	}

	hashedString, err := h.service.GetHashedString(str)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"message": "Cannot hash"})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"hashed_string": hashedString,
	})
}
