package repository

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
)

type BlogRepository interface {
	Create(c context.Context, blog *domain.Blog) (*domain.Blog, error)
	FindAll(c context.Context) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsByAuthorID(c context.Context, authorID string) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsBySlug(c context.Context, nickname string) ([]domain.BlogWithAuthorData, error)
	FindByID(c context.Context, id int64) (*domain.BlogWithAuthorData, error)
	FindByUrlSlug(c context.Context, slug string, userID *string) (*domain.BlogWithAuthorData, error)
	// Update(user *Blog) error
	Delete(c context.Context, id int64, userId string) (*int64, error)
	// Cache table
	CreateUserIDAuthorProfileIDCacheRecord(c context.Context, userID string, authorID string, slug string, displayName string) error
	VerifyAuthorIDByUserID(c context.Context, userID string) (string, error)
	// Author deleted event
	UpdateBlogStatusForDeletedAuthor(c context.Context, authorID string) error
	DeleteAuthorHardDeletedBlogs(c context.Context, authorID string) error
	DeleteAuthorCache(c context.Context, authorID string) error
	MarkAuthorCacheAsDeleted(c context.Context, authorID string) error
	RestoreBlog(c context.Context, blogID int64, PreviousStatus string) error
	GetAuthorProfileByUserID(c context.Context, userID string) (*domain.AuthorData, error)
	// Blog metrics
	UpdateBlogReactionCount(c context.Context, blogID int64, transition ReactionTransition) error
}
