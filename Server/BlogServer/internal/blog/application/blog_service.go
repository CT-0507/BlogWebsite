package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
)

type BlogService interface {
	CreateWithOutBox(c context.Context, blog *domain.Blog, userID string) error
	OnBlogPosted(c context.Context, payload []byte) error
	GetAll(ctx context.Context) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsByAuthorID(ctx context.Context, authorID string) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsByNickname(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error)
	GetBlog(ctx context.Context, id int64) (*domain.BlogWithAuthorData, error)
	GetBlogByUrlSlug(ctx context.Context, slug string) (*domain.BlogWithAuthorData, error)
	DeleteBlog(ctx context.Context, id int64, userID string) (*int64, error)
	VerifyAuthorIDByUserID(c context.Context, userID string) (string, error)
	// Event handler
	OnAuthorCreated(c context.Context, payload []byte) error
	OnAuthorDeleted(c context.Context, payload []byte) error
	OnAuthorHardDeleted(c context.Context, payload []byte) error
}
