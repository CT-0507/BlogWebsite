package application

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts"
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/flows"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

type DeleteBlogUseCase struct {
	txManager  database.TxManager
	repo       domain.BlogRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewDeleteBlogUseCases(txManager database.TxManager, repo domain.BlogRepository, outboxRepo outboxrepo.OutboxRepository) *DeleteBlogUseCase {
	return &DeleteBlogUseCase{
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

// Soft delete blog
func (u *DeleteBlogUseCase) DeleteBlog(ctx context.Context, id int64, userID string) (*int64, error) {

	foundBlog, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if foundBlog == nil {
		return nil, errors.New("Blog not found")
	}

	eventPayload := &contracts.DeleteBlogKickstartPayload{
		BlogID:    foundBlog.BlogID,
		Status:    foundBlog.Status,
		AuthorID:  foundBlog.Author.AuthorID,
		DeletedBy: userID,
	}

	payload, err := json.Marshal(eventPayload)
	if err != nil {
		return nil, err
	}

	eventContext := &contracts.DeleteBlogKickstartContext{
		BlogID:    foundBlog.BlogID,
		Status:    foundBlog.Status,
		DeletedBy: userID,
	}

	context, err := json.Marshal(eventContext)
	if err != nil {
		return nil, err
	}

	sagaID := uuid.New()

	err = u.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
		SagaID:    &sagaID,
		EventType: flows.DeleteBlogSaga,
		Payload:   payload,
		Context:   &context,
	})

	return &foundBlog.BlogID, nil
}
