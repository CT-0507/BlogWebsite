package application

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
)

type EventHandlerUsecases struct {
	txManager database.TxManager
	repo      domain.BlogRepository
}

func NewEventHandlerUsecases(txManager database.TxManager, repo domain.BlogRepository) *EventHandlerUsecases {
	return &EventHandlerUsecases{
		txManager: txManager,
		repo:      repo,
	}
}

func (u *EventHandlerUsecases) OnAuthorCreated(c context.Context, evt *messaging.OutboxEvent) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var event domain.AuthorCreatedEvent
		if err := json.Unmarshal(evt.Payload, &event); err != nil {
			return err
		}

		return u.repo.CreateUserIDAuthorProfileIDCacheRecord(c, event.UserID, event.AuthorID, event.Slug, event.DisplayName)
	})
}

func (u *EventHandlerUsecases) OnAuthorDeleted(c context.Context, evt *messaging.OutboxEvent) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var event domain.AuthorDeletedEvent
		if err := json.Unmarshal(evt.Payload, &event); err != nil {
			return err
		}
		return u.repo.UpdateBlogStatusForDeletedAuthor(c, event.AuthorID)
	})
}

func (u *EventHandlerUsecases) OnAuthorHardDeleted(c context.Context, evt *messaging.OutboxEvent) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var event domain.AuthorDeletedEvent
		if err := json.Unmarshal(evt.Payload, &event); err != nil {
			return err
		}
		err := u.repo.DeleteAuthorCache(c, event.AuthorID)
		if err != nil {
			return err
		}
		return u.repo.DeleteAuthorHardDeletedBlogs(c, event.AuthorID)
	})
}
