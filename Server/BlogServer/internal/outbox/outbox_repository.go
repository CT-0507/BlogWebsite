package outbox

import (
	"context"

	outboxdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
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

func (r *OutboxRepositoryImpl) Insert(ctx context.Context, topic string, payload []byte) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := outboxdb.New(db)
	return q.InsertRecord(ctx, outboxdb.InsertRecordParams{
		Topic:   topic,
		Payload: payload,
	})
}

func (r *OutboxRepositoryImpl) UpdateProcessedAt(ctx context.Context, q *outboxdb.Queries, outboxID []int64) error {
	return q.UpdateProcessedAt(ctx, outboxID)
}

func (r *OutboxRepositoryImpl) UpdateRetries(ctx context.Context, q *outboxdb.Queries, outboxID []int64) error {
	return q.UpdateRetiresInBatch(ctx, outboxID)
}

func (r *OutboxRepositoryImpl) GetUnprocessedEvent(ctx context.Context, q *outboxdb.Queries) ([]outboxdb.GetUnprocessedEventRow, error) {
	return q.GetUnprocessedEvent(ctx)
}
