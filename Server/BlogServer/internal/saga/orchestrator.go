package saga

import (
	"context"
	"errors"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/messaging"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/google/uuid"
)

type Orchestrator struct {
	registry   domain.Registry
	txManager  database.TxManager
	repo       domain.SagaRepository
	outboxRepo outbox.OutboxRepository
}

var OrderSaga = domain.SagaDefinition{
	Name: "order",
	Steps: []domain.Step{
		{
			Name:           "create_order",
			ActionType:     "CreateOrder",
			CompensateType: "CancelOrder",
			Next:           "reserve_inventory",
		},
		{
			Name:           "reserve_inventory",
			ActionType:     "ReserveInventory",
			CompensateType: "ReleaseInventory",
			Next:           "charge_payment",
		},
		{
			Name:           "charge_payment",
			ActionType:     "ChargePayment",
			CompensateType: "RefundPayment",
			Next:           "Complete",
		},
	},
}

func NewOrchestrator(registry domain.Registry, txManager database.TxManager, repo domain.SagaRepository) *Orchestrator {
	return &Orchestrator{
		registry:  registry,
		txManager: txManager,
		repo:      repo,
	}
}

// payload map[string]interface{}

func (o *Orchestrator) StartSaga(ctx context.Context, defName string, payload map[string]interface{}) error {
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		currentStep := o.registry.GetStepByIndex(defName, 0)

		if currentStep == nil {
			return errors.New("")
		}

		saga := &domain.Saga{
			Type:        defName,
			Status:      domain.SagaRunning,
			CurrentStep: 0,
			Context:     payload,
		}

		// 1. Save saga + steps
		sagaID, err := o.repo.CreateSaga(ctx, saga)
		if err != nil {
			return err
		}

		uuid := uuid.MustParse(sagaID)

		step := &domain.SagaStep{
			SagaID:    uuid,
			StepIndex: 0,
			Name:      currentStep.Name,
			Context:   payload,
			Status:    domain.StepStatus(domain.SagaRunning),
		}

		err = o.repo.CreateStep(ctx, step)
		if err != nil {
			return err
		}
		// 2. Insert first event into outbox
		event := messaging.OutboxEvent{
			SagaID:    &uuid,
			EventType: currentStep.ActionType,
			Payload:   payload,
		}

		err = o.outboxRepo.Insert(ctx, &event)
		if err != nil {
			return err
		}

		return nil
	})
}

func (o *Orchestrator) HandleEvent(ctx context.Context, e *messaging.OutboxEvent) {
	o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		instance, err := o.repo.GetSagaByID(ctx, *e.SagaID)
		if err != nil {
			return err
		}

		currentStep := instance.CurrentStep
		err = o.repo.UpdateStep(ctx, *e.SagaID, currentStep, domain.SagaCompensating)
		if err != nil {
			return err
		}

		err = o.repo.UpdateSagaCurrentStep(ctx, instance.ID)
		if err != nil {
			return err
		}

		nextStep := o.registry.GetNextStep(instance.Type, currentStep)
		if nextStep == nil {
			return errors.New("")
		}

		err = o.outboxRepo.Insert(ctx, &messaging.OutboxEvent{
			SagaID:    &instance.ID,
			EventType: nextStep.ActionType,
			Payload:   instance.Payload,
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

		currentStep, err := o.repo.GetSagaByIndex(ctx, *e.SagaID, instance.CurrentStep)
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
			o.startCompensation(ctx, instance)
		}

		err = o.repo.UpdateStepRetryCount(ctx, *e.SagaID, instance.CurrentStep, "retrying")
		if err != nil {
			return err
		}

		// re-enqueue event
		o.outboxRepo.Insert(ctx, e)
		return nil
	})
}

func (o *Orchestrator) startCompensation(ctx context.Context, instance *domain.Saga) error {
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		// 1. mark saga as compensating
		err := o.repo.UpdateSagaStatus(ctx, instance.ID, "compensating")
		if err != nil {
			return err
		}

		// 2. find last completed step
		step, err := o.repo.GetSagaByIndex(ctx, instance.ID, instance.CurrentStep-1)
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
			EventType: step.Name,
			Payload:   step.Context, // usually use response data
		}

		err = o.outboxRepo.Insert(ctx, &event)
		if err != nil {
			return err
		}
		return nil
	})
}

func (o *Orchestrator) HandleCompensationSuccess(ctx context.Context, evt messaging.OutboxEvent) error {
	return o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {

		instance, err := o.repo.GetSagaByID(ctx, *evt.SagaID)
		if err != nil {
			return err
		}

		// 1. mark step compensated
		err = o.repo.UpdateStep(ctx, *evt.SagaID, instance.CurrentStep, domain.SagaStatus(domain.StepCompensated))
		if err != nil {
			return err
		}

		// 2. find previous completed step
		prevStep, err := o.repo.GetSagaByIndex(ctx, *evt.SagaID, instance.CurrentStep-1)
		if err != nil {
			return err
		}

		if prevStep == nil {
			// done compensating
			o.repo.UpdateSagaStatus(ctx, *evt.SagaID, "failed")
		}

		prevStepDef := o.registry.GetStepByIndex(instance.Type, instance.CurrentStep-1)

		// 3. enqueue next compensation
		event := messaging.OutboxEvent{
			ID:        uuid.New(),
			SagaID:    evt.SagaID,
			EventType: prevStepDef.CompensateType,
			// Payload:   prevStep.
		}

		err = o.outboxRepo.Insert(ctx, &event)
		if err != nil {
			return err
		}

		return nil
	})
}

func (o *Orchestrator) HandleCompensationFailure(ctx context.Context, evt messaging.OutboxEvent) {
	instance, err := o.repo.GetSagaByID(ctx, *evt.SagaID)
	if err != nil {
		return
	}

	currentStep, err := o.repo.GetSagaByIndex(ctx, *evt.SagaID, instance.CurrentStep)
	if err != nil {
		return
	}

	o.txManager.WithVoidTx(ctx, func(ctx context.Context) error {
		currentStep.RetryCount++

		if currentStep.RetryCount >= 3 {
			// send to DLQ
			o.repo.InsertDLQ(ctx, currentStep)

			// saga stuck → manual intervention
			o.repo.UpdateSagaStatus(ctx, instance.ID, "failed")
			return nil
		}

		// retry
		return o.outboxRepo.Insert(ctx, &evt)
	})
}
