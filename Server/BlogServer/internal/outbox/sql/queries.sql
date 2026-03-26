-- name: InsertRecord :exec
INSERT INTO outbox.outbox_events
    (saga_id, event_type, payload)
VALUES ($1, $2, $3);

-- name: UpdateProcessedAt :exec
UPDATE outbox.outbox_events
    SET processed_at = NOW()
WHERE id = ANY($1::UUID[]);

-- name: GetUnprocessedEvent :many
SELECT *
FROM outbox.outbox_events
WHERE processed_at IS NULL AND retries < 3
ORDER BY created_at
LIMIT 50
FOR UPDATE SKIP LOCKED;

-- name: UpdateRetiresInBatch :exec
UPDATE outbox.outbox_events
SET retry_count = retry_count + 1
WHERE id = ANY($1::UUID[]);