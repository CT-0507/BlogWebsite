package application

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type UpdateAuthorProfileUsecases struct {
	txManager  *database.TxManager
	repo       domain.AuthorProfileRepository
	outboxRepo outbox.OutboxRepository
}

func NewUpdateAuthorProfileUsecases(txManager *database.TxManager, repo domain.AuthorProfileRepository, outboxRepo outbox.OutboxRepository) *UpdateAuthorProfileUsecases {
	return &UpdateAuthorProfileUsecases{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}
