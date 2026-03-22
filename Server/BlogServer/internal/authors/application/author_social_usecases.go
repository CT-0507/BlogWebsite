package application

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type AuthorSocialUsecases struct {
	txManager  *database.TxManager
	repo       domain.AuthorProfileRepository
	outboxRepo outbox.OutboxRepository
}

func NewAuthorSocialUsecases(txManager *database.TxManager, repo domain.AuthorProfileRepository, outboxRepo outbox.OutboxRepository) *AuthorSocialUsecases {
	return &AuthorSocialUsecases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

func (u *AuthorSocialUsecases) SetFeatureBlogs(ctx context.Context, authorID string, blogIDs []string) error {
	return u.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		_, err := u.repo.CreateAuthorFeatureBlogs(ctx, authorID, blogIDs)
		return err
	})
}

func (u *AuthorSocialUsecases) GetFeatureBlogsByAuthorID(ctx context.Context, slug string) ([]string, error) {
	return u.repo.GetAuthorFeaturedBlogIDs(ctx, slug)
}
