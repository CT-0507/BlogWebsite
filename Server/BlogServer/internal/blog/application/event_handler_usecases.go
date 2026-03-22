package application

import (
	"context"
	"encoding/json"
	"log"

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
	log.Println("Inside on author created1")
	return u.txManager.WithVoidTx(c, func(ctx context.Context) error {

		log.Println("Inside on author created2")
		var evt domain.AuthorCreatedEvent
		if err := json.Unmarshal(payload, &evt); err != nil {
			return err
		}
		log.Println(evt)

		return u.repo.CreateUserIDAuthorProfileIDCacheRecord(c, evt.UserID, evt.AuthorID)
	})
}
