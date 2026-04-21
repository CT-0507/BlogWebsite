-- internal/blog/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS blogs;

CREATE TABLE blogs.blogs (
    blog_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id TEXT NOT NULL,
    url_slug VARCHAR(400) NOT NULL UNIQUE,
    title TEXT NOT NULL,
    content TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    -- blog ranking
    like_count INT NOT NULL DEFAULT 0,
    dislike_count  INT NOT NULL DEFAULT 0,

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
    slug TEXT NOT NULL,
    display_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    UNIQUE(user_id),
    UNIQUE(user_id, author_id)
);

CREATE INDEX ON blogs.blogs (title);
CREATE INDEX ON blogs.blogs (url_slug);
CREATE INDEX ON blogs.blog_tags (tag_id);

-- Comment

CREATE TABLE blogs.comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blog_id BIGINT NOT NULL REFERENCES blogs.tags(tag_id) ON DELETE CASCADE,

    content TEXT NOT NULL,

    -- identity
    actor_type VARCHAR(20) NOT NULL, -- 'user' | 'author' | 'deleted'
    actor_id TEXT, -- external ID (nullable if deleted)

    -- snapshot for rendering (no joins needed)
    actor_display_name  VARCHAR(100) NOT NULL,
    actor_avatar_url    TEXT,
    -- actor_badge         VARCHAR(50),

    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, hidden, deleted
    -- threading (2 levels)
    parent_comment_id   UUID NULL,
    root_comment_id     UUID NOT NULL,

    depth SMALLINT  NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE blogs.comment_reactions (
    id          UUID PRIMARY KEY,
    comment_id  UUID NOT NULL,

    user_id     VARCHAR(64) NOT NULL,
    type        VARCHAR(20) NOT NULL,  -- like, dislike, etc.

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(comment_id, user_id, type)
);

CREATE TABLE blogs.blog_reactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blog_id     BIGINT NOT NULL REFERENCES blogs.tags(tag_id) ON DELETE CASCADE,

    user_id     TEXT NOT NULL,
    type        VARCHAR(20) NOT NULL,  -- like, dislike, etc.
    status VARCHAR(20) NOT NULL, -- active | hidden | deleted

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,

    UNIQUE(blog_id, user_id, type)
);

CREATE INDEX idx_comments_post_id 
ON comments(post_id, created_at);
CREATE INDEX idx_comments_root 
ON comments(root_comment_id, created_at);
CREATE INDEX idx_comments_parent 
ON comments(parent_comment_id);
CREATE INDEX idx_comments_status 
ON comments(status);
CREATE INDEX idx_comments_reaction 
ON blog_reactions(blog_id, type) WHERE status = 'active';