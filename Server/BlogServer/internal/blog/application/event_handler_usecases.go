package application

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
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

func (u *EventHandlerUsecases) OnAuthorCreated(c context.Context, payload []byte) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var evt domain.AuthorCreatedEvent
		if err := json.Unmarshal(payload, &evt); err != nil {
			return err
		}

		return u.repo.CreateUserIDAuthorProfileIDCacheRecord(c, evt.UserID, evt.AuthorID)
	})
}

func (u *EventHandlerUsecases) OnAuthorDeleted(c context.Context, payload []byte) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var evt domain.AuthorDeletedEvent
		if err := json.Unmarshal(payload, &evt); err != nil {
			return err
		}
		return u.repo.UpdateBlogStatusForDeletedAuthor(c, evt.AuthorID)
	})
}

func (u *EventHandlerUsecases) OnAuthorHardDeleted(c context.Context, payload []byte) error {
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		var evt domain.AuthorDeletedEvent
		if err := json.Unmarshal(payload, &evt); err != nil {
			return err
		}
		err := u.repo.DeleteAuthorCache(c, evt.AuthorID)
		if err != nil {
			return err
		}
		return u.repo.DeleteAuthorHardDeletedBlogs(c, evt.AuthorID)
	})
}
