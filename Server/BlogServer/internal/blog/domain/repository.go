package domain

import (
	"context"

	"github.com/google/uuid"
)

type BlogRepository interface {
	Create(c context.Context, blog *Blog) (*Blog, error)
	FindAll(c context.Context) ([]BlogWithAuthorData, error)
	ListAuthorBlogsByAuthorID(c context.Context, authorID uuid.UUID) ([]BlogWithAuthorData, error)
	ListAuthorBlogsByNickname(c context.Context, nickname string) ([]BlogWithAuthorData, error)
	FindByID(c context.Context, id int64) (*BlogWithAuthorData, error)
	FindByUrlSlug(c context.Context, slug string) (*BlogWithAuthorData, error)
	// Update(user *Blog) error
	Delete(c context.Context, id int64, userId uuid.UUID) (*int64, error)
}
