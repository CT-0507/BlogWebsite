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
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newUserId, err := h.service.RegisterUser(ctx, &User{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  hashedPassword,
		Roles:     []string{"admin"},
	})
	if err != nil {
		switch {
		case errors.Is(err, &ErrUsernameAlreadyTaken{}):
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
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

	foundUser, loginErr := h.service.LoginUser(ctx, user.Username)
	if loginErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": loginErr.Error()})
		return
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if compareErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
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

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

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
