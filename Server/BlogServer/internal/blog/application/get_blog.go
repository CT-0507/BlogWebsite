package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type GetBlogUseCases struct {
	txManager *database.TxManager
	repo      domain.BlogRepository
}

func NewGetBlogUseCases(txManager *database.TxManager, repo domain.BlogRepository) *GetBlogUseCases {
	return &GetBlogUseCases{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *GetBlogUseCases) GetBlog(ctx context.Context, id int64) (*domain.Blog, error) {
	return s.repo.FindByID(ctx, id)
}
