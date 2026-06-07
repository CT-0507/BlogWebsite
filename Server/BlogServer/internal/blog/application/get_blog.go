package application

import (
	"context"
	"log"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/repository"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type GetBlogUseCases struct {
	txManager database.TxManager
	repo      repository.BlogRepository
}

func NewGetBlogUseCases(txManager database.TxManager, repo repository.BlogRepository) *GetBlogUseCases {
	return &GetBlogUseCases{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *GetBlogUseCases) GetBlog(ctx context.Context, id int64) (*domain.BlogWithAuthorData, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *GetBlogUseCases) GetBlogByUrlSlug(ctx context.Context, slug string, userID *string) (*domain.BlogWithAuthorData, error) {
	result, err := s.repo.FindByUrlSlug(ctx, slug, userID)

	if err != nil {
		return nil, err
	}

	if len(result.Tags) == 0 {
		result.Tags = []string{}
	}

	go func(id int64) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := s.repo.UpdateViewCount(ctx, id); err != nil {
			log.Println("daily metric failed:", err)
		}
	}(result.BlogID)

	return result, err
}
