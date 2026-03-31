package infrastructure

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"
	sagadb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/infrastructure/db"
)

func mapSagaToDomainSaga(sagadb *sagadb.SagaSaga) *domain.Saga {
	return &domain.Saga{
		ID:          sagadb.ID,
		Type:        sagadb.SagaType,
		Status:      domain.SagaStatus(sagadb.Status),
		CurrentStep: sagadb.CurrentStep,
		Context:     sagadb.Context,
	}
}

func mapStepToDomainStep(stepdb *sagadb.SagaSagaStep) *domain.SagaStep {
	return &domain.SagaStep{
		ID:         stepdb.ID,
		SagaID:     stepdb.SagaID,
		StepIndex:  stepdb.StepIndex,
		StepName:   stepdb.StepName,
		EventID:    stepdb.EventID,
		Input:      stepdb.Input,
		Output:     &stepdb.Output,
		Status:     domain.StepStatus(stepdb.Status),
		RetryCount: stepdb.RetryCount,
	}
}
