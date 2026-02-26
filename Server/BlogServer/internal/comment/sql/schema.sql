CREATE TABLE comments (
    comment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
    blog_id BIGINT REFERENCES blogs(blog_id),
    user_id BIGINT REFERENCES users(user_id),
    content TEXT,
    active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES users(user_id),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users(user_id)
);