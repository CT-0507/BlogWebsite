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

    CONSTRAINT comments_depth_check
        CHECK (depth IN (0, 1)),

    CONSTRAINT comments_parent_depth_check
        CHECK (
            (depth = 0 AND parent_comment_id IS NULL)
            OR
            (depth = 1 AND parent_comment_id IS NOT NULL)
        )
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

-- internal/authors/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS authors;

CREATE TABLE authors.authors (
    author_id VARCHAR(26) PRIMARY KEY CHECK (LENGTH(author_id) = 26),
    user_id TEXT NOT NULL UNIQUE, -- external reference (NO FK)
    display_name VARCHAR(100) NOT NULL,

    bio TEXT,

    avatar TEXT,
    slug VARCHAR(150) NOT NULL UNIQUE,
    social_link TEXT,

    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, suspended, deleted

    email TEXT,

    -- denormalized counters
    follower_count INT DEFAULT 0,
    blog_count INT DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT
);

CREATE TABLE authors.author_profiles (
    author_id      TEXT PRIMARY KEY,

    twitter_url    TEXT,
    github_url     TEXT,
    website_url    TEXT,

    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_author_profiles_author
        FOREIGN KEY (author_id)
        REFERENCES authors.authors(author_id)
        ON DELETE CASCADE
);

CREATE TABLE authors.author_followers (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    author_id       VARCHAR(50) NOT NULL,
    user_id         VARCHAR(50) NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    CONSTRAINT fk_followers_author
        FOREIGN KEY (author_id)
        REFERENCES authors.authors(author_id)
        ON DELETE CASCADE,
    UNIQUE(author_id, user_id)
);

CREATE TABLE authors.author_featured_blogs (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    author_id       TEXT NOT NULL,
    blog_id         TEXT NOT NULL,  -- external (NO FK)

    position        INT NOT NULL, -- ordering

    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_featured_author
        FOREIGN KEY (author_id)
        REFERENCES authors.authors(author_id)
        ON DELETE CASCADE,

    UNIQUE(author_id, blog_id),
    UNIQUE(author_id, position)
);

CREATE INDEX idx_authors_slug ON authors.authors(slug);
CREATE INDEX idx_authors_user_id ON authors.authors(user_id);
CREATE INDEX idx_authors_status ON authors.authors(status);

-- author_followers
CREATE INDEX idx_followers_author_id ON authors.author_followers(author_id);
CREATE INDEX idx_followers_user_id ON authors.author_followers(user_id);

-- author_featured_blogs
CREATE INDEX idx_featured_author_id ON authors.author_featured_blogs(author_id);

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

-- internal/blog/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    email VARCHAR(50) UNIQUE,
    password TEXT NOT NULL,
    nickname VARCHAR(50) UNIQUE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    avatar TEXT,
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

CREATE TABLE users.notifications (
    notification_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID REFERENCES users.users(user_id),
    content JSONB NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users.users(user_id) ON DELETE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users.users(user_id) ON DELETE CASCADE,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users.users(user_id) ON DELETE CASCADE
);

CREATE TABLE users.contacts (
    contact_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID REFERENCES users.users(user_id) NOT NULL,
    email VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON users.users (user_id);

insert into users.users (
    user_id,
    username,
    password,
    nickname,
    first_name,
    last_name,
    role
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'system',
    'password',
    'system',
    'system',
    'system',
    'system'
);

-- insert admin
insert into users.users (
    user_id,
    username,
    password,
    nickname,
    first_name,
    last_name,
    role
) VALUES (
    '00000000-0000-0000-0000-000000000002',
    'root',
    '$2a$10$MY335VAnFaJcj7E/UkyeNOTc7wWuk5uFQIbWGLk6a8vSW9C1FgTxm',
    'admin',
    'admin',
    'admin',
    'admin'
);

insert into users.users (
    user_id,
    username,
    password,
    nickname,
    first_name,
    last_name,
    role
) VALUES (
    '00000000-0000-0000-0000-000000000004',
    'user1',
    '$2a$10$MY335VAnFaJcj7E/UkyeNOTc7wWuk5uFQIbWGLk6a8vSW9C1FgTxm',
    'test-user1',
    'test',
    'user',
    'user'
);
insert into users.users (
    user_id,
    username,
    password,
    nickname,
    first_name,
    last_name,
    role
) VALUES (
    '00000000-0000-0000-0000-000000000005',
    'user2',
    '$2a$10$MY335VAnFaJcj7E/UkyeNOTc7wWuk5uFQIbWGLk6a8vSW9C1FgTxm',
    'test-user2',
    'test',
    'user',
    'user'
);
insert into users.users (
    user_id,
    username,
    password,
    nickname,
    first_name,
    last_name,
    role
) VALUES (
    '00000000-0000-0000-0000-000000000006',
    'user3',
    '$2a$10$MY335VAnFaJcj7E/UkyeNOTc7wWuk5uFQIbWGLk6a8vSW9C1FgTxm',
    'test-user3',
    'test',
    'user',
    'user'
);