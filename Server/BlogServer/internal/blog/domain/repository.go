package domain

import (
	"context"
)

type BlogRepository interface {
	Create(c context.Context, blog *Blog) (*Blog, error)
	FindAll(c context.Context) ([]BlogWithAuthorData, error)
	ListAuthorBlogsByAuthorID(c context.Context, authorID string) ([]BlogWithAuthorData, error)
	ListAuthorBlogsBySlug(c context.Context, nickname string) ([]BlogWithAuthorData, error)
	FindByID(c context.Context, id int64) (*BlogWithAuthorData, error)
	FindByUrlSlug(c context.Context, slug string) (*BlogWithAuthorData, error)
	// Update(user *Blog) error
	Delete(c context.Context, id int64, userId string) (*int64, error)
	// Cache table
	CreateUserIDAuthorProfileIDCacheRecord(c context.Context, userID string, authorID string, slug string, displayName string) error
	VerifyAuthorIDByUserID(c context.Context, userID string) (string, error)
	// Author deleted event
	UpdateBlogStatusForDeletedAuthor(c context.Context, authorID string) error
	DeleteAuthorHardDeletedBlogs(c context.Context, authorID string) error
	DeleteAuthorCache(c context.Context, authorID string) error
}
