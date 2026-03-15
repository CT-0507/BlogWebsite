package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

type DeleteBlogUseCase struct {
	txManager *database.TxManager
	repo      domain.BlogRepository
}

func NewDeleteBlogUseCase(txManager *database.TxManager, repo domain.BlogRepository) *DeleteBlogUseCase {
	return &DeleteBlogUseCase{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *DeleteBlogUseCase) DeleteBlog(ctx context.Context, id int64, userID uuid.UUID) (*int64, error) {
	return s.repo.Delete(ctx, id, userID)
}
