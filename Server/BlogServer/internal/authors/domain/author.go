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
