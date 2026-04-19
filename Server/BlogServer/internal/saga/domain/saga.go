package domain

import (
	"time"

	"github.com/google/uuid"
)

type SagaStatus string

const (
	SagaRunning      SagaStatus = "running"
	SagaCompleted    SagaStatus = "completed"
	SagaFailed       SagaStatus = "failed"
	SagaCompensating SagaStatus = "compensating"
)

type StepStatus string

const (
	StepPending      StepStatus = "pending"
	StepCompleted    StepStatus = "completed"
	StepFailed       StepStatus = "failed"
	StepCompensating StepStatus = "compensating"
	StepCompensated  StepStatus = "compensated"
)

type Saga struct {
	ID          uuid.UUID
	Type        string
	Status      SagaStatus
	CurrentStep int32
	Context     []byte
	Error       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SagaStep struct {
	ID         uuid.UUID
	SagaID     uuid.UUID
	StepIndex  int32
	StepName   string
	EventID    uuid.UUID
	Input      []byte
	Output     *[]byte
	Status     StepStatus
	RetryCount int32
}

type Step struct {
	Name           string
	ActionType     string
	CompensateType string
	Next           string
	MaxRetries     int32
}

type SagaDefinition struct {
	Name  string
	Steps []Step
}
