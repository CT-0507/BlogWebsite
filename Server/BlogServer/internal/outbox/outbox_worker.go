package outbox

import (
	"context"
	"log"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

const MAX_RETRIES = 3

type OutboxRepository interface {
	Insert(ctx context.Context, event *messaging.OutboxEvent) error
	UpdateProcessedAt(ctx context.Context, outboxIDs []uuid.UUID) error
	GetUnprocessedEvent(ctx context.Context) ([]messaging.OutboxEvent, error)
	UpdateRetries(ctx context.Context, outboxIDs []uuid.UUID) error
}

type OutboxWorker struct {
	txManager  database.TxManager
	publisher  messaging.EventPublisher
	outboxRepo OutboxRepository
}

func NewOutboxWorker(txManager database.TxManager, publisher messaging.EventPublisher, outboxRepo OutboxRepository) *OutboxWorker {
	return &OutboxWorker{
		txManager:  txManager,
		publisher:  publisher,
		outboxRepo: outboxRepo,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			w.processBatch(ctx)

		case <-ctx.Done():
			return
		}
	}
}

func (w *OutboxWorker) processBatch(ctx context.Context) {
	w.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		events, err := w.outboxRepo.GetUnprocessedEvent(ctx)
		if err != nil {
			return err
		}

		var ids []uuid.UUID
		var failedIds []uuid.UUID

		for _, evt := range events {

			if evt.RetryCount >= MAX_RETRIES {
				continue
			}

			errs := w.handleEvent(ctx, &evt)
			if errs != nil {
				log.Println(errs)
				failedIds = append(failedIds, evt.ID)
				continue
			}

			ids = append(ids, evt.ID)
		}

		if len(ids) > 0 {
			log.Println("Update proccessed")
			log.Println(ids[0])

			err := w.outboxRepo.UpdateProcessedAt(ctx, ids)
			if err != nil {
				return err
			}
		}

		if len(failedIds) > 0 {
			log.Println(failedIds[0])
			log.Println("Update retry")
			err := w.outboxRepo.UpdateRetries(ctx, failedIds)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (w *OutboxWorker) handleEvent(ctx context.Context, evt *messaging.OutboxEvent) []error {

	log.Print("Proccess topic: " + evt.EventType)

	// switch evt.EventType {

	// case "blog.created", "authorFollower.created", "authorFollower.deleted", "authorIdentity.created", "authorIdentity.deleted", "authorIdentity.hardDeleted":

	// 	return w.publisher.Publish(*w.context, evt)

	// case "notification.created":

	// 	return w.publisher.Publish(evt)

	// }

	// return nil
	return w.publisher.Publish(ctx, evt)
}
