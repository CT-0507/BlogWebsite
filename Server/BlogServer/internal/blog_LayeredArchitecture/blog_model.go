package blog

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/google/uuid"
)

type Tag struct {
}

type Blog struct {
	BlogID   int64     `json:"id"`
	AuthorID uuid.UUID `json:"author_id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Active   string    `json:"active"`
	// Tags    []Tag  `json:"tags"`
	// Images  []string
	model.Audit
}

type BlogWithAuthorDTO struct {
	BlogID  int64     `json:"blog_id"`
	Author  AuthorDTO `json:"author"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Active  string    `json:"active"`
	model.Audit
}

type AuthorDTO struct {
	AuthorID string `json:"author_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// event
type BlogCreatedEvent struct {
	BlogID    int64
	BlogTitle string
}

func (e BlogCreatedEvent) EventName() string {
	return "blog.created"
}
