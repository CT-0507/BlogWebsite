package domain

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"

type AuthorProfile struct {
	AuthorID    string
	UserID      string
	DisplayName string
	Bio         string
	Avatar      string
	Slug        string
	SocialLink  string
	Status      string
	Email       string
	model.Audit
}

// Error models
type ErrFailedToCreateAuthorProfile struct{}

func (e *ErrFailedToCreateAuthorProfile) Error() string {
	return "Failed to create author profile"
}

type ErrFailedToFollowAuthor struct{}

func (e *ErrFailedToFollowAuthor) Error() string {
	return "Failed to follow author"
}

type ErrFailedToUnfollowAuthor struct{}

func (e *ErrFailedToUnfollowAuthor) Error() string {
	return "Failed to follow author"
}

// Event models
type AuthorCreatedEvent struct {
	AuthorID string
}

func (e AuthorCreatedEvent) EventName() string {
	return "authorProfile.created"
}

type AuthorFollowedEvent struct {
	AuthorID string
	UserID   string
}

func (e AuthorFollowedEvent) EventName() string {
	return "authorFollower.created"
}

type AuthorUnfollowedEvent struct {
	AuthorID string
	UserID   string
}

func (e AuthorUnfollowedEvent) EventName() string {
	return "authorFollower.deleted"
}

type UpdateAuthorFollowCountEvent struct {
	AuthorID string
	UserID   string
}

func (e UpdateAuthorFollowCountEvent) EventName() string {
	return "authorFollower.deleted"
}
