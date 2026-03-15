package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
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
