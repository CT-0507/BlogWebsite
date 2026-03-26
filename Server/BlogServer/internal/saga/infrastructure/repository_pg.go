package infrastructure

import (
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

// func (r *SagaRepository) CreateSaga(ctx context.Context, saga *domain.Saga, steps []domain.SagaStep) error {

// 	db := utils.GetExecutor(ctx, r.pool)

// 	q := sagadb.New()

// 	payloadBytes, err := json.Marshal(saga.Payload)
// 	if err != nil {
// 		return err
// 	}

// 	contextBytes, err := json.Marshal(saga.Payload)
// 	if err != nil {
// 		return err
// 	}

// 	// Insert saga
// 	id, err := q.CreateSaga(ctx, sagadb.CreateSagaParams{})
// }
