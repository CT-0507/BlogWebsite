-- internal/blog/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    email VARCHAR(50) UNIQUE,
    password TEXT NOT NULL,
    nickname VARCHAR(50) NOT NULL UNIQUE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    role VARCHAR(10) NOT NULL DEFAULT 'user',
    status VARCHAR DEFAULT 'active',
    points INTEGER NOT NULL DEFAULT 0,
    -- refresh_token_id TEXT NOT REFERENCES users.refresh_tokens(token_id),
    token_version INT DEFAULT 0,
    last_logout TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users.users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users.users(user_id),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users.users(user_id)
);

CREATE TABLE users.refresh_tokens (
    token_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users.users(id),
    refresh_token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ
);

CREATE TABLE users.notifications (
    notification_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID REFERENCES users.users(user_id),
    content JSONB NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users.users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users.users(user_id),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users.users(user_id)
);

CREATE INDEX ON users.users (user_id);