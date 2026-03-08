-- name: InsertRecord :exec
INSERT INTO outbox.outbox_events
        (aggregate_type, aggregate_id, event_type, payload)
VALUES ($1,$2,$3,$4);

-- name: UpdateProcessedAt :exec
UPDATE outbox.outbox_events
    SET processed_at = NOW()
WHERE id = ANY($1::bigint[]);

-- name: GetUnprocessedEvent :many
SELECT id, event_type, payload
FROM outbox.outbox_events
WHERE processed_at IS NULL
ORDER BY created_at
LIMIT 50
FOR UPDATE SKIP LOCKED;