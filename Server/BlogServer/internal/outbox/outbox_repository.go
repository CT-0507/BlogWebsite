package outbox

import (
	"context"

	outboxdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepository interface {
	Insert(ctx context.Context, topic string, payload []byte) error
	UpdateProcessedAt(ctx context.Context, q *outboxdb.Queries, outboxID []int64) error
	GetUnprocessedEvent(ctx context.Context, q *outboxdb.Queries) ([]outboxdb.GetUnprocessedEventRow, error)
}

type outboxRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) OutboxRepository {
	return &outboxRepository{
		pool: pool,
	}
}

func (r *outboxRepository) Insert(ctx context.Context, topic string, payload []byte) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := outboxdb.New(db)
	return q.InsertRecord(ctx, outboxdb.InsertRecordParams{
		Topic:   topic,
		Payload: payload,
	})
}

func (r *outboxRepository) UpdateProcessedAt(ctx context.Context, q *outboxdb.Queries, outboxID []int64) error {
	return q.UpdateProcessedAt(ctx, outboxID)
}

func (r *outboxRepository) GetUnprocessedEvent(ctx context.Context, q *outboxdb.Queries) ([]outboxdb.GetUnprocessedEventRow, error) {
	return q.GetUnprocessedEvent(ctx)
}
