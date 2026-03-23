# Repository

##

✅ Avoid joins with Post/User modules

✅ Duplicate only what’s necessary

✅ Stay consistent via events later

##

### **🧱 1. authors Table (Core Aggregate)**

```sql
CREATE TABLE authors.authors (
    author_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL, -- external reference (NO FK)
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
    created_by TEXT REFERENCES NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT NOT NULL
);
```

- user_id → reference only (no join)

- follower_count, post_count → duplicated for fast reads

- status → avoids hard delete

###

### **🧱 2. author_profiles Table (Optional Split)**

```sql
CREATE TABLE authors.author_profiles (
    author_id      TEXT PRIMARY KEY,

    twitter_url    TEXT,
    github_url     TEXT,
    website_url    TEXT,

    created_at     TIMESTAMPTZ NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_author_profiles_author
        FOREIGN KEY (author_id)
        REFERENCES authors.authors(author_id)
        ON DELETE CASCADE
);
```

- can also be merged into authors

### **🧱 3. author_followers Table**

```sql
CREATE TABLE authors.authors.author_followers (
    id              VARCHAR(50) PRIMARY KEY,

    author_id       VARCHAR(50) NOT NULL,
    user_id         VARCHAR(50) NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_followers_author
        FOREIGN KEY (author_id)
        REFERENCES authors(id)
        ON DELETE CASCADE,
    UNIQUE(author_id, user_id)
);
```

🚀 With Denormalization

ALTER TABLE author_followers ADD COLUMN user_display_name VARCHAR(150);
ALTER TABLE author_followers ADD COLUMN user_avatar_url TEXT;

👉 Optional but useful for UI

### **🧱 3. author_featured_blogs**

```sql
CREATE TABLE authors.author_featured_blogs (
    id              int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    author_id       TEXT NOT NULL,
    blog_id         TEXT NOT NULL,  -- external (NO FK)

    position        INT NOT NULL, -- ordering

    created_at      TIMESTAMP NOT NULL,

    CONSTRAINT fk_featured_author
        FOREIGN KEY (author_id)
        REFERENCES authors.authors(author_id)
        ON DELETE CASCADE,

    UNIQUE(author_id, blog_id),
    UNIQUE(author_id, position)
);
```

We ONLY store:

```
author_id + post_id
```

No post data → no coupling

### **5. author_stats Table (Optional Optimization)**

```sql
CREATE TABLE author_stats (
    author_id       VARCHAR(50) PRIMARY KEY,

    follower_count  INT DEFAULT 0,
    blog_count      INT DEFAULT 0,

    updated_at      TIMESTAMP NOT NULL
);
```

### **6. Indexing Strategy**

authors

```sql
CREATE INDEX idx_authors_slug ON authors(slug);
CREATE INDEX idx_authors_user_id ON authors(user_id);
CREATE INDEX idx_authors_status ON authors(status);
```

author_followers

```sql
CREATE INDEX idx_followers_author_id ON author_followers(author_id);
CREATE INDEX idx_followers_user_id ON author_followers(user_id);
```

author_featured_posts

```sql
CREATE INDEX idx_featured_author_id ON author_featured_posts(author_id);
```

### **🔄 7. Event-Driven Sync (Future)**

Eventual consistency:

Events Author Module SHOULD CONSUME
From Post Module

PostCreated
→ increment post_count

PostDeleted
→ decrement post_count

PostUpdated
→ update duplicated fields (title, thumbnail)

From User Module

UserUpdated
→ update cached follower info (optional)

Events Author Module SHOULD EMIT

AuthorCreated

AuthorUpdated

AuthorDeleted

AuthorFollowed

AuthorUnfollowed

AuthorFeaturedPostsUpdated
