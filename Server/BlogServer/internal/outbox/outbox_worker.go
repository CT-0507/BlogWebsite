package outbox

import (
	"context"
	"log"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/event"
	outboxdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxWorker struct {
	db         *pgxpool.Pool
	bus        *event.Bus
	outboxRepo OutboxRepository
}

func NewOutboxWorker(db *pgxpool.Pool, bus *event.Bus, outboxRepo OutboxRepository) *OutboxWorker {
	return &OutboxWorker{
		db:         db,
		bus:        bus,
		outboxRepo: outboxRepo,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {

	log.Println("Worker started")

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

	tx, err := w.db.Begin(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback(ctx)
	q := outboxdb.New(tx)
	rows, err := w.outboxRepo.GetUnprocessedEvent(ctx, q)
	if err != nil {
		return
	}

	var ids []int64

	for _, row := range rows {

		err = w.handleEvent(row.Topic, row.Payload)
		if err != nil {
			continue
		}

		ids = append(ids, row.ID)
	}

	if len(ids) > 0 {
		err := w.outboxRepo.UpdateProcessedAt(ctx, q, ids)
		if err != nil {
			return
		}
	}

	tx.Commit(ctx)
}

// event
// type BlogCreatedEvent struct {
// 	BlogID int64
// 	Type   string
// }

func (w *OutboxWorker) handleEvent(topic string, payload []byte) error {

	log.Print("Proccess topic: ")
	log.Println(topic)

	switch topic {

	case "blog.created":
	case "authorFollower.created", "authorFollower.deleted", "authorIdentity.created":
		// var evt BlogCreatedEvent
		// json.Unmarshal(payload, &evt)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		log.Println("Here")

		return w.bus.Publish(ctx, topic, payload)

	case "notification.created":

		// var evt BlogCreatedEvent
		// json.Unmarshal(payload, &evt)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		return w.bus.Publish(ctx, topic, payload)

	}

	return nil
}
