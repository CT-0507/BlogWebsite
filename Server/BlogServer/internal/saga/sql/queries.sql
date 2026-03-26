-- name: CreateSaga :one
INSERT INTO saga.sagas (
  saga_type, status, current_step, payload, context, error
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: CreateSagaSteps :copyfrom
INSERT INTO saga.saga_steps (
  saga_id, step_index, step_name, status, retry_count, max_retries, next_retry_at, context
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: UpdateStepStatus :exec
UPDATE saga.saga_steps
SET status = 'completed'
WHERE saga_id = $1 AND step_index = $2;

-- name: UpdateSaga :exec
UPDATE saga.sagas
SET current_step = current_step + 1
WHERE id = $1;

-- name: UpdateStepRetries :exec
UPDATE saga.saga_steps
SET 
  retry_count = retry_count + 1,
  last_error = $1,
  next_retry_at = NOW() + INTERVAL '5 seconds'
WHERE saga_id = $2 AND step_index = $3;

-- name: UpdateSagaStatus :exec
UPDATE saga.sagas
SET status = $2
WHERE id = $1;