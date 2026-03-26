package outbox

import (
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	outboxdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox/db"
)

func MapToOutboxEvent(event *outboxdb.OutboxOutboxEvent) *messaging.OutboxEvent {
	var eventUnmarshal any
	json.Unmarshal(event.Payload, &eventUnmarshal)
	return &messaging.OutboxEvent{
		SagaID:    event.SagaID,
		EventType: event.EventType,
		Payload:   eventUnmarshal,
		Processed: event.ProcessedAt.Valid,
		CreatedAt: event.CreatedAt.Time,
	}
}
