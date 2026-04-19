package outbox

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	outboxdb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox/db"
)

func MapToOutboxEvent(event *outboxdb.OutboxOutboxEvent) *messaging.OutboxEvent {

	var err *string
	if event.Error.Valid {
		err = &event.Error.String
	}

	return &messaging.OutboxEvent{
		ID:         event.ID,
		SagaID:     event.SagaID,
		EventType:  event.EventType,
		Payload:    event.Payload,
		Context:    &event.Context,
		RetryCount: event.RetryCount,
		Processed:  event.ProcessedAt.Valid,
		CreatedAt:  event.CreatedAt.Time,
		Error:      err,
	}
}
