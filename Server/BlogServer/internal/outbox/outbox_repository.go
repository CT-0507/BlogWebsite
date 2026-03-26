package outbox

import (
	"context"
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	outboxdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepositoryImpl struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *OutboxRepositoryImpl {
	return &OutboxRepositoryImpl{
		pool: pool,
	}
}

func (r *OutboxRepositoryImpl) Insert(ctx context.Context, event *messaging.OutboxEvent) error {

	db := utils.GetExecutor(ctx, r.pool)

	data, _ := json.Marshal(event.Payload)

	q := outboxdb.New(db)
	return q.InsertRecord(ctx, outboxdb.InsertRecordParams{
		SagaID:    event.SagaID,
		EventType: event.EventType,
		Payload:   data,
	})
}

func (r *OutboxRepositoryImpl) UpdateProcessedAt(ctx context.Context, outboxIDs []uuid.UUID) error {
	db := utils.GetExecutor(ctx, r.pool)
	q := outboxdb.New(db)
	return q.UpdateProcessedAt(ctx, outboxIDs)
}

func (r *OutboxRepositoryImpl) UpdateRetries(ctx context.Context, outboxIDs []uuid.UUID) error {
	db := utils.GetExecutor(ctx, r.pool)
	q := outboxdb.New(db)
	return q.UpdateRetiresInBatch(ctx, outboxIDs)
}

func (r *OutboxRepositoryImpl) GetUnprocessedEvent(ctx context.Context) ([]messaging.OutboxEvent, error) {
	db := utils.GetExecutor(ctx, r.pool)
	q := outboxdb.New(db)
	outboxEvents, err := q.GetUnprocessedEvent(ctx)
	if err != nil {
		return nil, err
	}

	var events []messaging.OutboxEvent
	for _, value := range outboxEvents {
		v := value
		events = append(events, *MapToOutboxEvent(&v))
	}

	return events, nil
}
