package domain

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/google/uuid"
)

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

type BlogWithAuthorData struct {
	BlogID  int64      `json:"blog_id"`
	Author  AuthorData `json:"author"`
	Title   string     `json:"title"`
	Content string     `json:"content"`
	Active  string     `json:"active"`
	model.Audit
}

type AuthorData struct {
	AuthorID string `json:"author_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
