package saga

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/event"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/outbox"
	"github.com/google/uuid"
)

type Orchestrator struct {
	sagas map[string]*Saga
	steps map[string][]*SagaStep
	mu    sync.Mutex
	bus   *event.Bus
}

var OrderSaga = SagaDefinition{
	Name: "order",
	Steps: []Step{
		{Name: "create_order", ActionType: "CreateOrder"},
		{Name: "reserve_inventory", ActionType: "ReserveInventory"},
		{Name: "charge_payment", ActionType: "ChargePayment"},
	},
}

func NewOrchestrator(bus *event.Bus) *Orchestrator {
	return &Orchestrator{
		sagas: make(map[string]*Saga),
		steps: make(map[string][]*SagaStep),
		bus:   bus,
	}
}

func (o *Orchestrator) StartSaga(def SagaDefinition, payload map[string]interface{}) string {
	o.mu.Lock()
	defer o.mu.Unlock()

	id := uuid.New().String()

	saga := &Saga{
		ID:          id,
		Type:        def.Name,
		Status:      SagaRunning,
		CurrentStep: 0,
		Context:     payload,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	o.sagas[id] = saga

	var steps []*SagaStep
	for i, s := range def.Steps {
		steps = append(steps, &SagaStep{
			SagaID:    id,
			StepIndex: i,
			Name:      s.Name,
			Status:    StepPending,
		})
	}
	o.steps[id] = steps

	o.executeStep(saga, def)

	return id
}

func (o *Orchestrator) executeStep(saga *Saga, def SagaDefinition) {
	stepDef := def.Steps[saga.CurrentStep]
	payload, _ := json.Marshal(saga.Context)

	event := outbox.OutboxEvent{
		SagaID:    saga.ID,
		StepIndex: saga.CurrentStep,
		Type:      stepDef.ActionType,
		Payload:   payload,
	}

	eventPayload, _ := json.Marshal(event)

	o.bus.Publish(nil, "sage_step", eventPayload)
}

func (o *Orchestrator) HandleSuccess(evt outbox.OutboxEvent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	saga := o.sagas[evt.SagaID]
	step := o.steps[evt.SagaID][evt.StepIndex]

	step.Status = StepCompleted

	saga.CurrentStep++

	if saga.CurrentStep >= len(o.steps[evt.SagaID]) {
		saga.Status = SagaCompleted
		return
	}

	o.executeStep(saga, OrderSaga)
}

const MaxRetries = 3

func (o *Orchestrator) HandleFailure(evt outbox.OutboxEvent, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	step := o.steps[evt.SagaID][evt.StepIndex]
	step.RetryCount++
	step.LastError = err.Error()

	if step.RetryCount >= MaxRetries {
		// Send to DLQ
		// o.bus.dlq <- DeadLetterEvent{
		// 	Event:    evt,
		// 	FailedAt: time.Now(),
		// 	Reason:   err.Error(),
		// }

		// o.sagas[evt.SagaID].Status = SagaFailed
		return
	}

	payload, _ := json.Marshal(evt)

	// Retry
	go func() {
		time.Sleep(1 * time.Second) // backoff
		o.bus.Publish(nil, "", payload)
	}()
}
