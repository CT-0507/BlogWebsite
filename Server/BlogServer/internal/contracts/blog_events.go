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
