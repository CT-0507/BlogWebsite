package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type ListBlogsUseCase struct {
	txManager *database.TxManager
	repo      domain.BlogRepository
}

func NewListBlogsUseCase(txManager *database.TxManager, repo domain.BlogRepository) *ListBlogsUseCase {
	return &ListBlogsUseCase{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *ListBlogsUseCase) ListBlogs(ctx context.Context) ([]domain.BlogWithAuthorData, error) {
	return s.repo.FindAll(ctx)
}
