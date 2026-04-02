package domain

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
)

type Orchestrator interface {
	StartSaga(ctx context.Context, evt *messaging.OutboxEvent) error
	HandleEvent(ctx context.Context, evt *messaging.OutboxEvent) error
	HandleFailure(ctx context.Context, e *messaging.OutboxEvent) error
	StartCompensation(ctx context.Context, instance *Saga) error
	HandleCompensationSuccess(ctx context.Context, evt messaging.OutboxEvent) error
	HandleCompensationFailure(ctx context.Context, evt messaging.OutboxEvent)
}
