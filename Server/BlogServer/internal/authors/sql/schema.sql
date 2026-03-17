-- internal/authors/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS authors;

CREATE TABLE authors.authors (
    author_id TEXT PRIMARY KEY,
    user_id UUID REFERENCES users.users(user_id) ON DELETE CASCADE,
    display_name VARCHAR(100) NOT NULL,
    bio TEXT,
    avatar TEXT,
    slug TEXT NOT NULL UNIQUE,
    social_link TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    email TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users.users(user_id) ON DELETE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users.users(user_id) ON DELETE CASCADE,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users.users(user_id) ON DELETE CASCADE
);

CREATE INDEX ON authors.authors (slug);