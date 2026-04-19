-- name: CreateSaga :one
INSERT INTO saga.sagas (
  saga_type, status, current_step, context
) VALUES (
  $1, $2, $3, $4
)
RETURNING id;

-- name: CreateSagaStep :exec
INSERT INTO saga.saga_steps (
  saga_id, step_index, step_name, status, event_id, input
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: GetSagaByID :one
SELECT *
FROM saga.sagas
WHERE id = $1;

-- name: GetStepByIndexAndSagaID :one
SELECT *
FROM saga.saga_steps
WHERE saga_id = $1 AND step_index = $2;

-- name: UpdateStepStatusAndOutput :exec
UPDATE saga.saga_steps
SET status = $3, output = $4
WHERE saga_id = $1 AND step_index = $2;

-- name: UpdateStepStatus :exec
UPDATE saga.saga_steps
SET status = $3
WHERE saga_id = $1 AND step_index = $2;

-- name: UpdateSagaCurrentStep :exec
UPDATE saga.sagas
SET current_step = $2
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

-- name: UpdateSagaContextAndIncreaseStep :exec
UPDATE saga.sagas
SET context = context || $1::jsonb, current_step = current_step + 1
WHERE id = $2;

-- name: InsertDLQ :exec
INSERT INTO saga.dead_letter_queue (
  saga_id, step_index, step_name, status, event_id, input, output, last_error
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: GetLastCompletedStep :one
SELECT *
FROM saga.saga_steps
WHERE saga_id = $1 AND status = 'completed'
ORDER BY step_index DESC
LIMIT 1;

-- name: UpdateLastCompetedStepStatus :exec
UPDATE saga.saga_steps 
SET status = $1
WHERE id = (
  SELECT i.id
  FROM saga.saga_steps i
  WHERE i.saga_id = $2 AND i.status = 'completed'
  ORDER BY i.step_index DESC
  LIMIT 1
);

-- name: GetCompensatingStep :one
SELECT *
FROM saga.saga_steps
WHERE saga_id = $1 AND status = 'compensating'
ORDER BY step_index DESC
LIMIT 1;