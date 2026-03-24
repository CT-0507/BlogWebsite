package domain

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"

type AuthorProfile struct {
	AuthorID    string `json:"authorID"`
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName"`
	Bio         string `json:"bio"`
	Avatar      string `json:"avatar"`
	Slug        string `json:"slug"`
	SocialLink  string `json:"socialLink"`
	Status      string `json:"status"`
	Email       string `json:"email"`
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

// Event models
type AuthorCreatedEvent struct {
	AuthorID    string
	UserID      string
	Slug        string
	DisplayName string
}

func (e AuthorCreatedEvent) EventName() string {
	return "authorIdentity.created"
}

type AuthorDeletedEvent struct {
	AuthorID string
}

func (e AuthorDeletedEvent) EventName() string {
	return "authorIdentity.deleted"
}

type AuthorHardDeletedEvent struct {
	AuthorID string
}

func (e AuthorHardDeletedEvent) EventName() string {
	return "authorIdentity.hardDeleted"
}

type AuthorFollowedEvent struct {
	AuthorID    string
	UserID      string
	IsIncrement bool
}

func (e AuthorFollowedEvent) EventName() string {
	return "authorFollower.created"
}

type AuthorUnfollowedEvent struct {
	AuthorID    string
	UserID      string
	IsIncrement bool
}

func (e AuthorUnfollowedEvent) EventName() string {
	return "authorFollower.deleted"
}

type FollowCountChangedEvent struct {
	Slug        string
	AuthorID    string
	UserID      string
	IsIncrement bool
}

func (e FollowCountChangedEvent) EventName() string {
	eventType := "Increased"
	if !e.IsIncrement {
		eventType = "Decreased"
	}
	return "authorFollower.followerCount" + eventType
}

type BlogCountChangedEvent struct {
	BlogID   int64
	AuthorID string
}

func (e BlogCountChangedEvent) EventName() string {
	return "authorIdentity.blogCountIncreased"
}
