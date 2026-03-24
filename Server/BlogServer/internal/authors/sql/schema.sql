-- internal/authors/sql/schema.sql
CREATE SCHEMA IF NOT EXISTS authors;

CREATE TABLE authors.authors (
    author_id TEXT PRIMARY KEY,
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

-- 2. author_profiles Table 

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

-- 3. author_followers Table
CREATE TABLE authors.author_followers (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    author_id       VARCHAR(50) NOT NULL,
    user_id         VARCHAR(50) NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_followers_author
        FOREIGN KEY (author_id)
        REFERENCES authors.authors(author_id)
        ON DELETE CASCADE,
    UNIQUE(author_id, user_id)
);

-- 4. author_featured_blogs

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

-- CREATE TABLE author_stats (
--     author_id       VARCHAR(50) PRIMARY KEY,

--     follower_count  INT DEFAULT 0,
--     blog_count      INT DEFAULT 0,

--     updated_at      TIMESTAMP NOT NULL

-- CONSTRAINT fk_author_stats
--         FOREIGN KEY (author_id)
--         REFERENCES authors.authors(author_id)
--         ON DELETE CASCADE,
-- );
-- authors

CREATE INDEX idx_authors_slug ON authors.authors(slug);
CREATE INDEX idx_authors_user_id ON authors.authors(user_id);
CREATE INDEX idx_authors_status ON authors.authors(status);

-- author_followers
CREATE INDEX idx_followers_author_id ON authors.author_followers(author_id);
CREATE INDEX idx_followers_user_id ON authors.author_followers(user_id);

-- author_featured_blogs
CREATE INDEX idx_featured_author_id ON authors.author_featured_blogs(author_id);