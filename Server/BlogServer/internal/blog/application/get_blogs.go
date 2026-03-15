package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

// List blog related use cases like list all, list with filter, etc
type ListBlogsUseCases struct {
	txManager *database.TxManager
	repo      domain.BlogRepository
}

func NewListBlogsUseCases(txManager *database.TxManager, repo domain.BlogRepository) *ListBlogsUseCases {
	return &ListBlogsUseCases{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *ListBlogsUseCases) ListBlogs(ctx context.Context) ([]domain.BlogWithAuthorData, error) {
	return s.repo.FindAll(ctx)
}

func (s *ListBlogsUseCases) ListAuthorBlogsByAuthorID(ctx context.Context, authorID uuid.UUID) ([]domain.BlogWithAuthorData, error) {
	return s.repo.ListAuthorBlogsByAuthorID(ctx, authorID)
}

func (s *ListBlogsUseCases) ListAuthorBlogsByNickname(ctx context.Context, nickname string) ([]domain.BlogWithAuthorData, error) {
	return s.repo.ListAuthorBlogsByNickname(ctx, nickname)
}
