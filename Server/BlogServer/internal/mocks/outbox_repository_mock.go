package mocks

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockOutboxRepository struct {
	mock.Mock
}

func (m *MockOutboxRepository) Insert(ctx context.Context, event *messaging.OutboxEvent) error {
	args := m.Called(ctx, event)

	return args.Error(0)
}

func (m *MockOutboxRepository) UpdateProcessedAt(ctx context.Context, outboxIDs []uuid.UUID) error {
	args := m.Called(ctx, outboxIDs)

	return args.Error(0)
}

func (m *MockOutboxRepository) UpdateRetries(ctx context.Context, outboxIDs []uuid.UUID) error {
	args := m.Called(ctx, outboxIDs)

	return args.Error(0)
}

func (m *MockOutboxRepository) GetUnprocessedEvent(ctx context.Context) ([]messaging.OutboxEvent, error) {

	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]messaging.OutboxEvent), args.Error(1)
}
