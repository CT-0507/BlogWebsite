package domain

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
)

type Blog struct {
	BlogID   int64  `json:"id"`
	AuthorID string `json:"authorID"`
	Title    string `json:"title"`
	URLSlug  string `json:"urlSlug"`
	Content  string `json:"content"`
	Status   string `json:"status"`
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
	Status  string     `json:"status"`
	model.Audit
}

type AuthorData struct {
	AuthorID    string `json:"authorID"`
	Slug        string `json:"slug"`
	DisplayName string `json:"displayName"`
}
