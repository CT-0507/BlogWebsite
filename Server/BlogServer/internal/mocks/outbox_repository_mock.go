package mocks_test

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	outboxdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox/db"
	"github.com/stretchr/testify/mock"
)

type MockOutboxRepository struct {
	mock.Mock
}

func (m *MockOutboxRepository) Insert(ctx context.Context, topic string, payload []byte) error {
	args := m.Called(ctx, topic, payload)

	return args.Error(0)
}

func (m *MockOutboxRepository) UpdateProcessedAt(ctx context.Context, q *outboxdb.Queries, outboxID []int64) error {
	args := m.Called(ctx, q, outboxID)

	return args.Error(0)
}

func (m *MockOutboxRepository) GetUnprocessedEvent(ctx context.Context, q *outboxdb.Queries) ([]messaging.OutboxEvent, error) {
	args := m.Called(ctx, q)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]messaging.OutboxEvent), args.Error(1)
}
