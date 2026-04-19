package domain

import (
	"context"

	"github.com/google/uuid"
)

type SagaRepository interface {
	CreateSaga(ctx context.Context, saga *Saga) (*uuid.UUID, error)
	CreateStep(ctx context.Context, step *SagaStep) error
	GetSagaByID(ctx context.Context, sagaID uuid.UUID) (*Saga, error)
	GetStepByIndex(ctx context.Context, sagaID uuid.UUID, stepIndex int32) (*SagaStep, error)
	UpdateStepStatus(ctx context.Context, sagaID uuid.UUID, stepIndex int32, status StepStatus) error
	UpdateSagaCurrentStep(ctx context.Context, sagaID uuid.UUID, step int32) error
	UpdateStepRetryCount(ctx context.Context, sagaID uuid.UUID, stepIndex int32, status StepStatus, err string) error
	UpdateStepOutput(ctx context.Context, sagaID uuid.UUID, stepIndex int32, status StepStatus, output []byte) error
	UpdateSagaContext(ctx context.Context, sagaID uuid.UUID, context []byte) error
	UpdateSagaStatus(ctx context.Context, sagaID uuid.UUID, status SagaStatus) error
	InsertDLQ(ctx context.Context, step *SagaStep, err string) error
	GetLastCompletedStep(ctx context.Context, sagaID uuid.UUID) (*SagaStep, error)
	UpdateLastCompetedStepStatus(ctx context.Context, sagaID uuid.UUID, status StepStatus) error
	GetCompensatingStep(ctx context.Context, sagaID uuid.UUID) (*SagaStep, error)
}
