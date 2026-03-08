package outbox

import (
	"context"
	"strconv"

	outboxdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox/db"
	"github.com/jackc/pgx/v5"
)

type OutboxRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, aggregateType string, blogId int64, eventName string, payload []byte) error
	UpdateProcessedAt(ctx context.Context, q *outboxdb.Queries, outboxID []int64) error
	GetUnprocessedEvent(ctx context.Context, q *outboxdb.Queries) ([]outboxdb.GetUnprocessedEventRow, error)
}

type outboxRepository struct{}

func New() OutboxRepository {
	return &outboxRepository{}
}

func (r *outboxRepository) Insert(ctx context.Context, tx pgx.Tx, aggregateType string, blogId int64, eventName string, payload []byte) error {
	q := outboxdb.New(tx)
	return q.InsertRecord(ctx, outboxdb.InsertRecordParams{
		AggregateType: aggregateType,
		AggregateID:   strconv.FormatInt(blogId, 10),
		EventType:     eventName,
		Payload:       payload,
	})
}

func (r *outboxRepository) UpdateProcessedAt(ctx context.Context, q *outboxdb.Queries, outboxID []int64) error {
	return q.UpdateProcessedAt(ctx, outboxID)
}

func (r *outboxRepository) GetUnprocessedEvent(ctx context.Context, q *outboxdb.Queries) ([]outboxdb.GetUnprocessedEventRow, error) {
	return q.GetUnprocessedEvent(ctx)
}
