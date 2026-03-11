package utils

import (
	"errors"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
	UserID       string
	Username     string
	FirstName    string
	LastName     string
	Roles        []string
	TokenVersion int
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY = os.Getenv("SECRET_REFRESH_KEY")

func GenerateAllTokens(username, firstName, lastName, userID string, roles []string, tokenVer int) (string, string, error) {

	signedToken, err := GenerateToken(
		username, firstName, lastName, userID, roles, tokenVer, time.Now(), time.Now().Add(30*time.Minute), SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err := GenerateToken(
		username, firstName, lastName, userID, roles, tokenVer, time.Now(), time.Now().Add(2*time.Hour), SECRET_REFRESH_KEY)
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func GenerateToken(username, firstName, lastName, userID string, roles []string, tokenVer int, issuedAt, expiredAt time.Time, key string) (string, error) {

	claims := &SignedDetails{
		UserID:       userID,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		TokenVersion: tokenVer,
		Roles:        roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "BlogServer",
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {

	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil

}

func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {

	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return []byte(SECRET_REFRESH_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
