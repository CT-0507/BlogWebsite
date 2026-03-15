-- internal/blog/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS blogs;

CREATE TABLE blogs.blogs (
    blog_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES users.users(user_id),
    title TEXT NOT NULL,
    content TEXT,
    active VARCHAR(10) NOT NULL DEFAULT 'true',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users.users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users.users(user_id),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users.users(user_id)
);

CREATE TABLE blogs.tags (
    tag_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users.users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users.users(user_id),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users.users(user_id)
);

CREATE TABLE blogs.blog_tags (
    tag_id BIGINT NOT NULL REFERENCES blogs.tags(tag_id) ON DELETE CASCADE,
    blog_id BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users.users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users.users(user_id),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users.users(user_id),

    PRIMARY KEY (blog_id, tag_id)
);

CREATE INDEX ON blogs.blogs (title);
CREATE INDEX ON blogs.blog_tags (tag_id);