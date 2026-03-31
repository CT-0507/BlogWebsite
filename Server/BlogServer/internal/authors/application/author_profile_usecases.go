package application

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type AuthorProfileUsecases struct {
	txManager  database.TxManager
	repo       domain.AuthorProfileRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewAuthorProfileUsecases(txManager database.TxManager, repo domain.AuthorProfileRepository, outboxRepo outboxrepo.OutboxRepository) *AuthorProfileUsecases {
	return &AuthorProfileUsecases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}
