package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type GetBlogUseCase struct {
	txManager *database.TxManager
	repo      domain.BlogRepository
}

func NewGetBlogUseCase(txManager *database.TxManager, repo domain.BlogRepository) *GetBlogUseCase {
	return &GetBlogUseCase{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *GetBlogUseCase) GetBlog(ctx context.Context, id int64) (*domain.Blog, error) {
	return s.repo.FindByID(ctx, id)
}
