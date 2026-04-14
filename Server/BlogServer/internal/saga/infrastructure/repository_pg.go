package infrastructure

import (
	"context"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"
	sagadb "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/infrastructure/db"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SagaRepository struct {
	pool *pgxpool.Pool
}

func NewSagaRepository(pool *pgxpool.Pool) *SagaRepository {
	return &SagaRepository{
		pool: pool,
	}
}

func (r *SagaRepository) CreateSaga(ctx context.Context, saga *domain.Saga) (*uuid.UUID, error) {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	// Insert saga
	id, err := q.CreateSaga(ctx, sagadb.CreateSagaParams{
		SagaType:    saga.Type,
		Status:      string(saga.Status),
		CurrentStep: saga.CurrentStep,
		Context:     saga.Context,
	})
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *SagaRepository) CreateStep(ctx context.Context, step *domain.SagaStep) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	return q.CreateSagaStep(ctx, sagadb.CreateSagaStepParams{
		SagaID:    step.SagaID,
		StepIndex: step.StepIndex,
		StepName:  step.StepName,
		Status:    string(step.Status),
		EventID:   step.EventID,
		Input:     step.Input,
	})
}

func (r *SagaRepository) GetSagaByID(ctx context.Context, sagaID uuid.UUID) (*domain.Saga, error) {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	saga, err := q.GetSagaByID(ctx, sagaID)
	if err != nil {
		return nil, err
	}

	return mapSagaToDomainSaga(&saga), nil
}

func (r *SagaRepository) GetStepByIndex(ctx context.Context, sagaID uuid.UUID, stepIndex int32) (*domain.SagaStep, error) {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	step, err := q.GetStepByIndexAndSagaID(ctx, sagadb.GetStepByIndexAndSagaIDParams{
		SagaID:    sagaID,
		StepIndex: stepIndex,
	})
	if err != nil {
		return nil, err
	}

	return mapStepToDomainStep(&step), nil
}

func (r *SagaRepository) UpdateStepStatus(ctx context.Context, sagaID uuid.UUID, stepIndex int32, status domain.StepStatus) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	return q.UpdateStepStatus(ctx, sagadb.UpdateStepStatusParams{
		SagaID:    sagaID,
		StepIndex: stepIndex,
		Status:    string(status),
	})
}

func (r *SagaRepository) UpdateSagaCurrentStep(ctx context.Context, sagaID uuid.UUID, step int32) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	return q.UpdateSagaCurrentStep(ctx, sagadb.UpdateSagaCurrentStepParams{
		ID:          sagaID,
		CurrentStep: step,
	})
}

func (r *SagaRepository) UpdateStepRetryCount(ctx context.Context, sagaID uuid.UUID, stepIndex int32, status domain.StepStatus, err string) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	return q.UpdateStepRetries(ctx, sagadb.UpdateStepRetriesParams{
		LastError: pgtype.Text{
			Valid:  true,
			String: err,
		},
		SagaID:    sagaID,
		StepIndex: stepIndex,
	})
}

func (r *SagaRepository) UpdateStepOutput(ctx context.Context, sagaID uuid.UUID, stepIndex int32, status domain.StepStatus, output []byte) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	return q.UpdateStepStatusAndOutput(ctx, sagadb.UpdateStepStatusAndOutputParams{
		SagaID:    sagaID,
		StepIndex: stepIndex,
		Status:    string(status),
		Output:    output,
	})
}

func (r *SagaRepository) UpdateSagaContext(ctx context.Context, sagaID uuid.UUID, context []byte) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	return q.UpdateSagaContextAndIncreaseStep(ctx, sagadb.UpdateSagaContextAndIncreaseStepParams{
		ID:      sagaID,
		Column1: context,
	})
}

func (r *SagaRepository) UpdateSagaStatus(ctx context.Context, sagaID uuid.UUID, status domain.SagaStatus) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	return q.UpdateSagaStatus(ctx, sagadb.UpdateSagaStatusParams{
		ID:     sagaID,
		Status: string(status),
	})
}

func (r *SagaRepository) InsertDLQ(ctx context.Context, step *domain.SagaStep, err string) error {

	db := utils.GetExecutor(ctx, r.pool)

	q := sagadb.New(db)

	return q.InsertDLQ(ctx, sagadb.InsertDLQParams{
		SagaID:    step.SagaID,
		StepIndex: step.StepIndex,
		StepName:  step.StepName,
		Status:    string(step.Status),
		EventID:   step.EventID,
		Input:     step.Input,
		Output:    *step.Output,
		LastError: pgtype.Text{
			Valid:  true,
			String: err,
		},
	})
}
