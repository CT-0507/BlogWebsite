package domain

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/google/uuid"
)

type Orchestrator interface {
	StartSaga(ctx context.Context, def SagaDefinition, payload map[string]interface{}) error
	HandleSuccess(ctx context.Context, evt messaging.OutboxEvent) error
	HandleFailure(ctx context.Context, evt messaging.OutboxEvent, step SagaStep, err error) error
	startCompensation(ctx context.Context, sagaID uuid.UUID) error
	HandleCompensationSuccess(ctx context.Context, evt messaging.OutboxEvent) error
	HandleCompensationFailure(ctx context.Context, step SagaStep, evt messaging.OutboxEvent, err error)
}
