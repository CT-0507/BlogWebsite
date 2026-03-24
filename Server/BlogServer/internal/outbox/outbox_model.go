package outbox

import "time"

type OutboxEvent struct {
	ID        int64
	SagaID    string
	StepIndex int
	Type      string
	EventName string
	Payload   []byte
	Processed bool
	CreatedAt time.Time
}
