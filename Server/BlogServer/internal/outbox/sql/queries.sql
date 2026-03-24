-- name: InsertRecord :exec
INSERT INTO outbox.outbox_events
    (topic, payload)
VALUES ($1,$2);

-- name: UpdateProcessedAt :exec
UPDATE outbox.outbox_events
    SET processed_at = NOW()
WHERE id = ANY($1::bigint[]);

-- name: GetUnprocessedEvent :many
SELECT id, topic, payload, retries
FROM outbox.outbox_events
WHERE processed_at IS NULL AND retries < 3
ORDER BY created_at
LIMIT 50
FOR UPDATE SKIP LOCKED;

-- name: UpdateRetiresInBatch :exec
UPDATE outbox.outbox_events
SET retries = retries + 1
WHERE id = ANY($1::bigint[]);