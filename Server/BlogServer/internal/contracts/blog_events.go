package contracts

import "github.com/google/uuid"

type BlogCreatedEvent struct {
	BlogID    int64
	BlogTitle string
	AuthorID  string
}

func (e BlogCreatedEvent) EventName() string {
	return "blog.created"
}

type BlogCreatedSagaContext struct {
	AuthorID string
	UserID   uuid.UUID
}

type BlogCreatedSagaPayload struct {
	AuthorID string
	UserID   uuid.UUID
	Status   string
	Title    string
	Content  string
	UrlSlug  string
}

type CreateBlogCompensationPayload struct {
	BlogID int64
}

type CreateBlogAuthorCacheSuccessContext struct {
	AuthorID string
	UserID   string
}

type DeleteBlogAuthorCacheContext struct {
	AuthorID string
}

type DeleteBlogKickstartPayload struct {
	BlogID    int64
	Status    string
	DeletedBy string
	AuthorID  string
}

type DeleteBlogKickstartContext struct {
	BlogID    int64
	Status    string
	DeletedBy string
}

type DeleteBlogContext struct {
	BlogID         int64
	PreviousStatus string
}

type DeleteBlogPayload struct {
	AuthorID string
}
