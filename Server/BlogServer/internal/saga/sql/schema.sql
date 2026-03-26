-- Schema
CREATE SCHEMA IF NOT EXISTS saga;

CREATE TABLE saga.sagas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    saga_type TEXT NOT NULL,
    status TEXT NOT NULL, -- pending, running, completed, failed, compensating

    current_step TEXT,
    payload JSONB, -- input data for the saga
    context JSONB, -- accumulated results between steps

    error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE saga.saga_steps (
    id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    saga_id UUID REFERENCES sagas(id),

    step_index INT NOT NULL,
    step_name TEXT NOT NULL,
    status TEXT NOT NULL, -- pending, completed, failed, compensating, compensated

    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 3,
    next_retry_at TIMESTAMPTZ,

    context JSONB,

    last_error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    compensated_at TIMESTAMPTZ,

    UNIQUE (saga_id, step_index)
);

CREATE TABLE saga.dead_letter_queue (
    id UUID PRIMARY KEY,

    saga_id UUID,
    step_index INT,

    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,

    error TEXT NOT NULL,

    failed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_saga_steps_retry 
ON saga.saga_steps (status, next_retry_at);

CREATE INDEX idx_sagas_status 
ON saga.sagas (status);