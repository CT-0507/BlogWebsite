-- internal/outbox/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS outbox;

CREATE TABLE outbox.outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    saga_id UUID,

    event_type TEXT NOT NULL,
    context JSONB,
    payload JSONB NOT NULL,
    error TEXT,

    retry_count  INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_unprocessed
ON outbox.outbox_events (created_at)
WHERE processed_at IS NULL;