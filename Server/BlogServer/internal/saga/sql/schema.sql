-- Schema
CREATE SCHEMA IF NOT EXISTS saga;

CREATE TABLE sagas (
    id UUID PRIMARY KEY,
    saga_type TEXT NOT NULL,
    status TEXT NOT NULL, -- pending, running, completed, failed, compensating

    current_step TEXT,
    payload JSONB, -- input data for the saga
    context JSONB, -- accumulated results between steps

    error TEXT,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE saga_steps (
    id UUID PRIMARY KEY,
    saga_id UUID REFERENCES sagas(id),

    step_index INT NOT NULL,
    step_name TEXT NOT NULL,
    status TEXT NOT NULL, -- pending, completed, failed, compensating, compensated

    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 3,
    next_retry_at TIMESTAMP,

    request JSONB,
    response JSONB,

    last_error TEXT,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    UNIQUE (saga_id, step_index)
);