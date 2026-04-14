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
			StepName:  nextStepDef.Name,
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

const MaxRetries = 3

func (o *Orchestrator) HandleFailure(ctx context.Context, e *messaging.OutboxEvent) error {
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		instance, err := o.repo.GetSagaByID(ctx, *e.SagaID)
		if err != nil {
			return err
		}

		if instance.CurrentStep == 0 {
			err = o.repo.UpdateSagaStatus(ctx, instance.ID, domain.SagaFailed)
			if err != nil {
				log.Println(err)
				return err
			}
			err = o.repo.UpdateStepStatus(ctx, instance.ID, instance.CurrentStep, domain.StepCompensated)
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		}

		currentStep, err := o.repo.GetStepByIndex(ctx, *e.SagaID, instance.CurrentStep)
		if err != nil {
			return err
		}
		currentStep.RetryCount++

		stepDef := o.registry.GetStepByIndex(instance.Type, instance.CurrentStep)

		if currentStep.RetryCount >= stepDef.MaxRetries {
			// DLQ insert
			// o.repo.InsertDLQ(ctx, evt, err)

			// // mark saga failed
			// o.repo.MarkSagaFailed(ctx, evt.SagaID)
			// return nil
			o.StartCompensation(ctx, instance)
		}

		err = o.repo.UpdateStepRetryCount(ctx, *e.SagaID, instance.CurrentStep, "retrying", *e.Error)
		if err != nil {
			return err
		}

		// re-enqueue event
		o.outboxRepo.Insert(ctx, e)
		return nil
	})
}

func (o *Orchestrator) StartCompensation(ctx context.Context, instance *domain.Saga) error {
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		// 1. mark saga as compensating
		err := o.repo.UpdateSagaStatus(ctx, instance.ID, "compensating")
		if err != nil {
			return err
		}

		// 2. find last completed step
		step, err := o.repo.GetStepByIndex(ctx, instance.ID, instance.CurrentStep-1)
		if err != nil {
			return err
		}

		if step == nil {
			// nothing to compensate
			o.repo.UpdateSagaStatus(ctx, instance.ID, "failed")
		}

		// 3. enqueue compensation event
		event := messaging.OutboxEvent{
			ID:        uuid.New(),
			SagaID:    &instance.ID,
			EventType: step.StepName,
			Payload:   *step.Output, // usually use response data
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
		err = o.repo.UpdateStepStatus(ctx, *evt.SagaID, instance.CurrentStep, domain.StepCompensated)
		if err != nil {
			return err
		}

		// 2. find previous completed step
		prevStep, err := o.repo.GetStepByIndex(ctx, *evt.SagaID, instance.CurrentStep-1)
		if err != nil {
			return err
		}

		if prevStep == nil {
			// done compensating
			o.repo.UpdateSagaStatus(ctx, *evt.SagaID, "failed")
			return nil
		}

		// Update saga current step
		err = o.repo.UpdateSagaCurrentStep(ctx, instance.ID, instance.CurrentStep-1)
		if err != nil {
			return err
		}

		prevStepDef := o.registry.GetStepByIndex(instance.Type, instance.CurrentStep-1)

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

	currentStep, err := o.repo.GetStepByIndex(ctx, *evt.SagaID, instance.CurrentStep)
	if err != nil {
		return err
	}

	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		currentStep.RetryCount++

		prevStepDef := o.registry.GetStepByIndex(instance.Type, instance.CurrentStep-1)

		if currentStep.RetryCount >= prevStepDef.MaxRetries {
			// send to DLQ
			o.repo.InsertDLQ(ctx, currentStep, *evt.Error)

			// saga stuck → manual intervention
			o.repo.UpdateSagaStatus(ctx, instance.ID, domain.SagaFailed)
			return nil
		}

		err := o.repo.UpdateStepRetryCount(ctx, currentStep.SagaID, currentStep.StepIndex, domain.StepPending, *evt.Error)
		if err != nil {
			return err
		}
		// retry
		return o.outboxRepo.Insert(ctx, evt)
	})

}
