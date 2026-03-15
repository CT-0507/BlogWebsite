package domain

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/google/uuid"
)

type Blog struct {
	BlogID   int64     `json:"id"`
	AuthorID uuid.UUID `json:"authorID"`
	Title    string    `json:"title"`
	URLSlug  string    `json:"urlSlug"`
	Content  string    `json:"content"`
	Active   string    `json:"active"`
	// Tags    []Tag  `json:"tags"`
	// Images  []string
	model.Audit
}

type BlogWithAuthorData struct {
	BlogID  int64      `json:"blogID"`
	Author  AuthorData `json:"author"`
	URLSlug string     `json:"urlSlug"`
	Title   string     `json:"title"`
	Content string     `json:"content"`
	Active  string     `json:"active"`
	model.Audit
}

type AuthorData struct {
	AuthorID string `json:"authorID"`
	Nickname string `json:"nickname"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}
