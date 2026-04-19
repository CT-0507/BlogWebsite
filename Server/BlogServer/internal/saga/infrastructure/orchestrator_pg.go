package infrastructure

import (
	"context"
	"errors"
	"log"
	"time"

	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

type Orchestrator struct {
	registry   domain.Registry
	txManager  database.TxManager
	repo       domain.SagaRepository
	outboxRepo outboxrepo.OutboxRepository
}

func NewOrchestrator(registry domain.Registry, txManager database.TxManager, repo domain.SagaRepository, outboxRepo outboxrepo.OutboxRepository) *Orchestrator {
	return &Orchestrator{
		registry:   registry,
		txManager:  txManager,
		repo:       repo,
		outboxRepo: outboxRepo,
	}
}

// payload map[string]interface{}

func (o *Orchestrator) StartSaga(ctx context.Context, evt *messaging.OutboxEvent) error {
	log.Println(evt.EventType)
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		currentStep := o.registry.GetStepByIndex(evt.EventType, 0)

		if currentStep == nil {
			return errors.New("")
		}

		saga := &domain.Saga{
			Type:        evt.EventType,
			Status:      domain.SagaRunning,
			CurrentStep: 0,
			Context:     *evt.Context,
		}

		// 1. Save saga + steps
		sagaID, err := o.repo.CreateSaga(ctx, saga)
		if err != nil {
			return err
		}

		step := &domain.SagaStep{
			SagaID:    *sagaID,
			EventID:   evt.ID,
			StepIndex: 0,
			StepName:  currentStep.Name,
			Input:     evt.Payload,
			Status:    domain.StepStatus(domain.SagaRunning),
		}

		err = o.repo.CreateStep(ctx, step)
		if err != nil {
			return err
		}
		// 2. Insert first event into outbox
		event := &messaging.OutboxEvent{
			SagaID:    sagaID,
			EventType: currentStep.ActionType,
			Payload:   evt.Payload,
		}

		err = o.outboxRepo.Insert(ctx, event)
		if err != nil {
			return err
		}

		return nil
	})
}

func (o *Orchestrator) HandleEvent(ctx context.Context, evt *messaging.OutboxEvent) error {
	// Event payload for next step
	// Event context for updating saga context and step input
	ctxTimout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return o.txManager.WithVoidTx(ctxTimout, func(ctx context.Context) error {

		// Create next step first to short circuit in case event is proccessed
		instance, err := o.repo.GetSagaByID(ctx, *evt.SagaID)
		if err != nil {
			return err
		}
		currentStep := instance.CurrentStep

		// Update current Step status and output
		err = o.repo.UpdateStepOutput(ctx, *evt.SagaID, currentStep, domain.StepCompleted, *evt.Context)
		if err != nil {
			return err
		}

		nextStepDef := o.registry.GetNextStep(instance.Type, currentStep)
		if nextStepDef == nil {
			err = o.repo.UpdateSagaStatus(ctx, instance.ID, domain.SagaCompleted)
			if err != nil {
				return err
			}
			return nil
		}

		nextStep := &domain.SagaStep{
			SagaID:    *evt.SagaID,
			StepIndex: currentStep + 1,
			StepName:  nextStepDef.ActionType,
			EventID:   evt.ID,
			Input:     evt.Payload,
		}
		err = o.repo.CreateStep(ctx, nextStep)
		if err != nil {
			return err
		}

		// Update saga context and current step
		err = o.repo.UpdateSagaContext(ctx, instance.ID, *evt.Context)
		if err != nil {
			return err
		}

		// Insert event to outbox
		err = o.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    &instance.ID,
			EventType: nextStep.StepName,
			Payload:   nextStep.Input,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (o *Orchestrator) HandleFailure(ctx context.Context, e *messaging.OutboxEvent) error {
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		// 1. Get saga context
		instance, err := o.repo.GetSagaByID(ctx, *e.SagaID)
		if err != nil {
			return err
		}

		// 2. Get current step to check retry count
		currentStep, err := o.repo.GetStepByIndex(ctx, *e.SagaID, instance.CurrentStep)
		if err != nil {
			return err
		}

		stepDef := o.registry.GetStepByIndex(instance.Type, instance.CurrentStep)

		if currentStep.RetryCount >= stepDef.MaxRetries {
			err = o.StartCompensation(ctx, instance)
			if err != nil {
				return err
			}
			return nil
		}

		// 3. Update step status
		err = o.repo.UpdateStepRetryCount(ctx, *e.SagaID, instance.CurrentStep, "retrying", *e.Error)
		if err != nil {
			return err
		}

		// 4. re-enqueue event
		return o.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    e.SagaID,
			EventType: stepDef.ActionType,
			Payload:   currentStep.Input,
		})

	})
}

func (o *Orchestrator) StartCompensation(ctx context.Context, instance *domain.Saga) error {
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		if instance.CurrentStep == 0 {
			// nothing to compensate
			err := o.repo.UpdateSagaStatus(ctx, instance.ID, domain.SagaFailed)
			if err != nil {
				return err
			}
			err = o.repo.UpdateStepStatus(ctx, instance.ID, instance.CurrentStep, domain.StepFailed)
			if err != nil {
				return err
			}
			return nil
		}

		// 1. mark saga as compensating
		err := o.repo.UpdateSagaStatus(ctx, instance.ID, "compensating")
		if err != nil {
			return err
		}

		// 2. mark current step as failed
		err = o.repo.UpdateStepStatus(ctx, instance.ID, instance.CurrentStep, domain.StepFailed)
		if err != nil {
			return err
		}

		// 3. find last completed step
		step, err := o.repo.GetStepByIndex(ctx, instance.ID, instance.CurrentStep-1)
		if err != nil {
			return err
		}

		// 4. Update Step status as compensating
		err = o.repo.UpdateStepStatus(ctx, instance.ID, instance.CurrentStep-1, domain.StepCompensating)
		if err != nil {
			return err
		}

		compensationEventName := o.registry.GetStepByIndex(instance.Type, instance.CurrentStep).CompensateType

		// 3. enqueue compensation event
		event := messaging.OutboxEvent{
			ID:        uuid.New(),
			SagaID:    &instance.ID,
			EventType: compensationEventName,
			Payload:   *step.Output,
		}

		err = o.outboxRepo.Insert(ctx, &event)
		if err != nil {
			return err
		}
		return nil
	})
}

func (o *Orchestrator) HandleCompensationSuccess(ctx context.Context, evt *messaging.OutboxEvent) error {
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		instance, err := o.repo.GetSagaByID(ctx, *evt.SagaID)
		if err != nil {
			return err
		}

		// 1. mark step compensated
		err = o.repo.UpdateLastCompetedStepStatus(ctx, *evt.SagaID, domain.StepCompensated)
		if err != nil {
			return err
		}

		// 2. find previous completed step
		prevStep, err := o.repo.GetLastCompletedStep(ctx, *evt.SagaID)
		if err != nil {
			return err
		}

		if prevStep == nil {
			// done compensating
			err := o.repo.UpdateSagaStatus(ctx, *evt.SagaID, domain.SagaFailed)
			if err != nil {
				return err
			}
			return nil
		}

		prevStepDef := o.registry.GetStepByIndex(instance.Type, prevStep.StepIndex)

		// 3. enqueue next compensation
		event := messaging.OutboxEvent{
			ID:        uuid.New(),
			SagaID:    evt.SagaID,
			EventType: prevStepDef.CompensateType,
			Payload:   *prevStep.Output,
		}

		err = o.outboxRepo.Insert(ctx, &event)
		if err != nil {
			return err
		}

		return nil
	})
}

func (o *Orchestrator) HandleCompensationFailure(ctx context.Context, evt *messaging.OutboxEvent) error {

	instance, err := o.repo.GetSagaByID(ctx, *evt.SagaID)
	if err != nil {
		return err
	}

	compensatingStep, err := o.repo.GetCompensatingStep(ctx, *evt.SagaID)
	if err != nil {
		return err
	}

	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		prevStepDef := o.registry.GetStepByIndex(instance.Type, compensatingStep.StepIndex)

		if compensatingStep.RetryCount >= prevStepDef.MaxRetries {
			// send to DLQ
			err = o.repo.InsertDLQ(ctx, compensatingStep, *evt.Error)
			if err != nil {
				return err
			}

			// saga stuck → manual intervention
			err = o.repo.UpdateSagaStatus(ctx, instance.ID, domain.SagaFailed)
			if err != nil {
				return err
			}
			return nil
		}

		err := o.repo.UpdateStepRetryCount(ctx, compensatingStep.SagaID, compensatingStep.StepIndex, domain.StepCompensating, *evt.Error)
		if err != nil {
			return err
		}
		// retry
		return o.outboxRepo.Insert(ctx, evt)
	})

}
