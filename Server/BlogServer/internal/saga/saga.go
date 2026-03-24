package saga

import "time"

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
	ID          string
	Type        string
	Status      SagaStatus
	CurrentStep int
	Context     map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SagaStep struct {
	SagaID     string
	StepIndex  int
	Name       string
	Status     StepStatus
	RetryCount int
	LastError  string
}

type Step struct {
	Name       string
	ActionType string
	Compensate string
}

type SagaDefinition struct {
	Name  string
	Steps []Step
}
