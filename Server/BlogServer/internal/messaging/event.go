package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID         uuid.UUID
	SagaID     *uuid.UUID
	EventType  string
	RetryCount int32
	Payload    []byte
	Context    *[]byte
	Error      *string
	Processed  bool
	CreatedAt  time.Time
}

type EventPublisher interface {
	Publish(ctx context.Context, e *OutboxEvent) []error
}
