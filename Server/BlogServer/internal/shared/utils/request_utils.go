package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPreferredLang(c *gin.Context) *string {

	lang := c.GetHeader("Accept-Language")

	primaryLang := strings.Split(lang, ",")[0]

	return &primaryLang
}

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}
	tokenString := authHeader[len("Bearer "):]

	if tokenString == "" {
		return "", errors.New("Bearer token is required")
	}
	// tokenString, err := c.Cookie("access_token")
	// if err != nil {

	// 	return "", err
	// }

	return tokenString, nil

}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userID")

	if !exists {
		return uuid.UUID{}, errors.New("userID does not exists in this context")
	}

	id, err := uuid.Parse(userID.(string))

	if err != nil {
		return uuid.UUID{}, errors.New("unable to retrieve userID")
	}

	return id, nil

}

func GetUserIDStringFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")

	if !exists {
		return "", errors.New("userID does not exists in this context")
	}

	return userID.(string), nil

}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")

	if !exists {
		return "", errors.New("role does not exists in this context")
	}

	memberRole, ok := role.(string)

	if !ok {
		return "", errors.New("unable to retrieve role")
	}

	return memberRole, nil

}
