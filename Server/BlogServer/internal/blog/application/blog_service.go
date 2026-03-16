package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/google/uuid"
)

type BlogService interface {
	CreateWithOutBox(c context.Context, blog *domain.Blog) error
	OnBlogPosted(c context.Context, payload []byte) error
	GetAll(ctx context.Context) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsByAuthorID(ctx context.Context, authorID uuid.UUID) ([]domain.BlogWithAuthorData, error)
	ListAuthorBlogsByNickname(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error)
	GetBlog(ctx context.Context, id int64) (*domain.BlogWithAuthorData, error)
	GetBlogByUrlSlug(ctx context.Context, slug string) (*domain.BlogWithAuthorData, error)
	DeleteBlog(ctx context.Context, id int64, userID uuid.UUID) (*int64, error)
}
