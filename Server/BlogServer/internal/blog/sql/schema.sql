-- internal/blog/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS blogs;

CREATE TABLE blogs.blogs (
    blog_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id TEXT NOT NULL,
    url_slug VARCHAR(400) NOT NULL UNIQUE,
    title TEXT NOT NULL,

    content_json JSONB NOT NULL,

    content_text TEXT NOT NULL,
    thumbnail_url TEXT,
    title_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('english', coalesce(title, ''))
    ) STORED,
    content_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('english', coalesce(content_text, ''))
    ) STORED,

    status VARCHAR(20) NOT NULL DEFAULT 'active',

    -- blog ranking
    like_count BIGINT NOT NULL DEFAULT 0,
    dislike_count  BIGINT NOT NULL DEFAULT 0,
    daily_access_count BIGINT NOT NULL DEFAULT 0,
    weekly_access_count BIGINT NOT NULL DEFAULT 0,
    access_count BIGINT NOT NULL DEFAULT 0,

    -- moderation
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    report_count BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT
);

CREATE INDEX idx_title_vector ON blogs.blogs USING GIN (title_vector);
CREATE INDEX idx_content_vector ON blogs.blogs USING GIN (content_vector);

CREATE TABLE blogs.tags (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    UNIQUE(id),
    UNIQUE(name)
);

CREATE TABLE blogs.blog_tags (
    tag_id BIGINT NOT NULL REFERENCES blogs.tags(id) ON DELETE CASCADE,
    blog_id BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    PRIMARY KEY (blog_id, tag_id)
);

CREATE TABLE blogs.idx_user_author_profile (

    user_id TEXT NOT NULL,
    author_id TEXT NOT NULL,
    avatar TEXT,
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
CREATE INDEX idx_blogs_author_status
ON blogs.blogs (author_id, status);

-- Comment

CREATE TABLE blogs.comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blog_id BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,

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

    like_count INT NOT NULL DEFAULT 0,
    dislike_count  INT NOT NULL DEFAULT 0,

    depth SMALLINT  NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    UNIQUE(blog_id, actor_id)
);

CREATE TABLE blogs.comment_reactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    comment_id  UUID NOT NULL REFERENCES blogs.comments(id) ON DELETE CASCADE,

    user_id     VARCHAR(64) NOT NULL,
    type        VARCHAR(20) NOT NULL,  -- like, dislike, etc.
    status      VARCHAR(20) NOT NULL DEFAULT 'active', -- active | hidden | deleted

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,

    UNIQUE(comment_id, user_id)
);

CREATE TABLE blogs.blog_reactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blog_id     BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,

    user_id     TEXT NOT NULL,
    type        VARCHAR(20) NOT NULL,  -- like, dislike, etc.
    status      VARCHAR(20) NOT NULL DEFAULT 'active', -- active | hidden | deleted

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,

    UNIQUE(blog_id, user_id)
);

CREATE INDEX idx_comments_blog_id 
ON blogs.comments(blog_id, created_at);
CREATE INDEX idx_comments_root 
ON blogs.comments(root_comment_id, created_at);
CREATE INDEX idx_comments_parent 
ON blogs.comments(parent_comment_id);
CREATE INDEX idx_comments_status 
ON blogs.comments(status);
CREATE INDEX idx_blogs_reaction 
ON blogs.blog_reactions(blog_id, type) WHERE status = 'active';
CREATE INDEX idx_comments_reaction 
ON blogs.comment_reactions(comment_id, type) WHERE status = 'active';

CREATE INDEX idx_blog_reactions_active
ON blogs.blog_reactions (
    blog_id,
    type,
    created_at
)
WHERE status = 'active';

-- Ranking table

CREATE TABLE blogs.blog_ranking (
    
    blog_id BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,

        -- rankings
    rank_all_time INT,
    rank_trending INT,

    -- scores
    score_all_time DOUBLE PRECISION,
    score_trending DOUBLE PRECISION,
    
    like_count INT NOT NULL  DEFAULT 0,
    dislike_count INT NOT NULL  DEFAULT 0,
    comment_count INT NOT NULL  DEFAULT 0,
    weekly_access_count INT NOT NULL  DEFAULT 0,
    daily_access_count  INT NOT NULL  DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    computed_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(blog_id)

);
CREATE TABLE blogs.blog_metrics (
    
    blog_id BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,

        -- rankings
    date DATE NOT NULL,
    views BIGINT NOT NULL DEFAULT 1,

    UNIQUE(blog_id, date)

);

-- CREATE TABLE blogs.blog_request_tracking (
    
--     blog_id BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,
--     request_id TEXT NOT NULL,
--     -- rankings

--     UNIQUE(blog_id, request_id)

-- );

CREATE TABLE blogs.reports (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    blog_id BIGINT NOT NULL REFERENCES blogs.blogs(blog_id) ON DELETE CASCADE,
    user_id     VARCHAR(64) NOT NULL,
    user_display_name TEXT NOT NULL,
    reason TEXT NOT NULL,
    

    -- audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    UNIQUE(blog_id, user_id)
);