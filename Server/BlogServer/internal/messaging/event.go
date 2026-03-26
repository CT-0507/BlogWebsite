package messaging

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID         uuid.UUID
	SagaID     *uuid.UUID
	EventType  string
	RetryCount int
	Payload    any
	Processed  bool
	CreatedAt  time.Time
}

type EventPublisher interface {
	Publish(event OutboxEvent) error
}
