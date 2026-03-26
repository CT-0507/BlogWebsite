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
	StepPending     StepStatus = "pending"
	StepCompleted   StepStatus = "completed"
	StepFailed      StepStatus = "failed"
	StepCompensated StepStatus = "compensated"
)

type Saga struct {
	ID          uuid.UUID
	Type        string
	Status      SagaStatus
	CurrentStep int32
	Payload     interface{}
	Context     map[string]interface{}
	Error       *string
	TotalSteps  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SagaStep struct {
	ID         uuid.UUID
	SagaID     uuid.UUID
	StepIndex  int
	Name       string
	Context    map[string]interface{}
	Status     StepStatus
	RetryCount int
}

type Step struct {
	Name           string
	ActionType     string
	CompensateType string
	Next           string
	MaxRetries     int
}

type SagaDefinition struct {
	Name  string
	Steps []Step
}
