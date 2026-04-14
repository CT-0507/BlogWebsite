package domain

import (
	"io"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
)

type AuthorProfile struct {
	AuthorID    string  `json:"authorID"`
	UserID      string  `json:"userID"`
	DisplayName string  `json:"displayName"`
	Bio         *string `json:"bio"`
	Avatar      *string `json:"avatar"`
	Slug        string  `json:"slug"`
	SocialLink  *string `json:"socialLink"`
	Status      string  `json:"status"`
	Email       *string `json:"email"`
	// Derived data
	FollowerCount int32 `json:"followerCount"`
	BlogCount     int32 `json:"blogCount"`
	model.Audit
}

// Error models
type ErrFailedToCreateAuthorProfile struct {
	Message string
}

func (e *ErrFailedToCreateAuthorProfile) Error() string {
	return e.Message
}

type ErrFailedToFollowAuthor struct {
	Message string
}

func (e *ErrFailedToFollowAuthor) Error() string {
	return "Failed to follow author"
}

type ErrFailedToUnfollowAuthor struct {
	Message string
}

func (e *ErrFailedToUnfollowAuthor) Error() string {
	return "Failed to follow author"
}

type ErrAuthorNotFound struct {
	Message string
}

func (e *ErrAuthorNotFound) Error() string {
	return e.Message
}

// Shared
type CreateUserFileStorageParams struct {
	File        io.Reader
	FileName    string
	ContentType string
}
