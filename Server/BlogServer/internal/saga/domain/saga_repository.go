package domain

import (
	"context"

	"github.com/google/uuid"
)

type SagaRepository interface {
	CreateSaga(ctx context.Context, saga *Saga) (string, error)
	CreateStep(ctx context.Context, step *SagaStep) error
	GetSagaByID(ctx context.Context, sagaID uuid.UUID) (*Saga, error)
	GetSagaByIndex(ctx context.Context, sagaID uuid.UUID, stepIndex int32) (*SagaStep, error)
	UpdateStep(ctx context.Context, sagaID uuid.UUID, stepIndex int32, status SagaStatus) error
	UpdateStepRetryCount(ctx context.Context, sagaID uuid.UUID, stepIndex int32, status SagaStatus) error
	UpdateSagaCurrentStep(ctx context.Context, sagaID uuid.UUID) error
	UpdateSagaStatus(ctx context.Context, sagaID uuid.UUID, status SagaStatus) error
	InsertDLQ(ctx context.Context, step *SagaStep) error
}
