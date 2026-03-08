package outbox

import "time"

type OutboxEvent struct {
	ID        int64
	EventName string
	Payload   []byte
	Processed bool
	CreatedAt time.Time
}
