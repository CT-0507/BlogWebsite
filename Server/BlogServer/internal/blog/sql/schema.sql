-- internal/blog/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS blogs;

CREATE TABLE blogs.blogs (
    blog_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id TEXT NOT NULL UNIQUE,
    url_slug VARCHAR(400) NOT NULL UNIQUE,
    title TEXT NOT NULL,
    content TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT
);

CREATE TABLE blogs.tags (
    tag_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT
);

CREATE TABLE blogs.blog_tags (
    tag_id BIGINT NOT NULL REFERENCES blogs.tags(tag_id) ON DELETE CASCADE,
    blog_id BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT,

    PRIMARY KEY (blog_id, tag_id)
);

CREATE TABLE blogs.idx_user_author_profile (
    user_id TEXT NOT NULL,
    author_id TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(user_id),
    UNIQUE(user_id, author_id)
);

CREATE INDEX ON blogs.blogs (title);
CREATE INDEX ON blogs.blogs (url_slug);
CREATE INDEX ON blogs.blog_tags (tag_id);