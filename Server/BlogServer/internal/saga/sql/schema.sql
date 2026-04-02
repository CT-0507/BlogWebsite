-- Schema
CREATE SCHEMA IF NOT EXISTS saga;

CREATE TABLE saga.sagas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    saga_type TEXT NOT NULL,
    status TEXT NOT NULL, -- pending, running, completed, failed, compensating

    current_step INT NOT NULL DEFAULT 0,
    context JSONB, -- accumulated results between steps

    error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE saga.saga_steps (
    id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    saga_id UUID NOT NULL REFERENCES saga.sagas(id),

    step_index INT NOT NULL,
    step_name TEXT NOT NULL,
    status TEXT NOT NULL, -- pending, completed, failed, compensating, compensated

    -- dedup
    event_id UUID UNIQUE NOT NULL,

    retry_count INT NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMPTZ,

    input JSONB,
    output JSONB,

    last_error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    compensated_at TIMESTAMPTZ,

    UNIQUE (saga_id, step_index)
);

CREATE TABLE saga.dead_letter_queue (

    id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    saga_id UUID NOT NULL REFERENCES saga.sagas(id),

    step_index INT NOT NULL,
    step_name TEXT NOT NULL,
    status TEXT NOT NULL, -- pending, completed, failed, compensating, compensated

    -- dedup
    event_id UUID UNIQUE NOT NULL,

    retry_count INT NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMPTZ,

    input JSONB,
    output JSONB,

    last_error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    compensated_at TIMESTAMPTZ,


    error TEXT NOT NULL,

    failed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (saga_id, step_index)
);

CREATE INDEX idx_saga_steps_retry 
ON saga.saga_steps (status, next_retry_at);

CREATE INDEX idx_sagas_status 
ON saga.sagas (status);