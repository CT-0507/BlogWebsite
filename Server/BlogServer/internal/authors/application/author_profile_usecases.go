package application

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type AuthorProfileUsecases struct {
	txManager  database.TxManager
	repo       domain.AuthorProfileRepository
	outboxRepo OutboxRepository
}

func NewAuthorProfileUsecases(txManager database.TxManager, repo domain.AuthorProfileRepository, outboxRepo OutboxRepository) *AuthorProfileUsecases {
	return &AuthorProfileUsecases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}
