-- name: GetBlog :one
SELECT 
    b.blog_id, 
    b.title,
    b.url_slug,
    b.author_id,
    b.content,
    b.like_count,
    b.dislike_count,
    b.status,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE b.blog_id = $1 AND b.deleted_at IS NULL;

-- name: GetBlogWithUserReaction :one
SELECT 
    b.*,
    i.slug,
    i.display_name,
    r.type AS reaction_type
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
LEFT JOIN blogs.blog_reactions r
    ON r.blog_id = b.blog_id
    AND r.user_id = $1
WHERE b.url_slug = $2 AND b.deleted_at IS NULL;

-- name: GetBlogByUrlSlug :one
SELECT 
    b.*,
    i.slug,
    i.display_name,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE b.url_slug = $1 AND b.deleted_at IS NULL;

-- name: ListBlogsByAuthor :many
SELECT
    b.blog_id, 
    b.title,
    b.url_slug,
    b.author_id,
    b.content,
    b.like_count,
    b.dislike_count,
    b.status,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE i.author_id = $1 AND b.deleted_at IS NULL AND b.status = $2;

-- name: ListBlogsByAuthorSlug :many
SELECT
    b.blog_id, 
    b.title,
    b.url_slug,
    b.author_id,
    b.content,
    b.status,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE i.slug = $1 AND b.deleted_at IS NULL AND b.status = $2;

-- name: ListAllBlogs :many
SELECT blog_id, title, content, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by FROM blogs.blogs;

-- name: ListBlogs :many
SELECT 
    b.blog_id,
    b.author_id,
    b.title, 
    b.url_slug,
    b.content,
    b.like_count,
    b.dislike_count, 
    b.status,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE b.deleted_at IS NULL;

-- name: CreateBlog :one
INSERT INTO blogs.blogs(
    author_id,
    title,
    url_slug,
    content,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateBlog :one
UPDATE blogs.blogs
    SET title = $1,
    content = $2
WHERE blog_id = $3
RETURNING blog_id;

-- name: HardDeleteBlog :one
DELETE FROM blogs.blogs
WHERE blog_id = $1
RETURNING blog_id;

-- name: DeleteBlog :one
UPDATE blogs.blogs
    SET deleted_by = $1,
    deleted_at = NOW(),
    status = 'deleted'
WHERE blog_id = $2
RETURNING blog_id;

-- name: CreateUserAuthorProfileIDCacheRecord :exec
INSERT INTO blogs.idx_user_author_profile (
    user_id,
    author_id,
    slug,
    display_name
) VALUES (
    $1, $2, $3, $4
);

-- name: VerifyAuthorIDByUserID :one
SELECT author_id
FROM blogs.idx_user_author_profile
WHERE user_id = $1;

-- name: UpdateBlogStatusForDeletedAuthor :exec
UPDATE blogs.blogs
SET status = 'author_deleted',
deleted_at = NOW()
WHERE blogs.author_id = $1;

-- name: DeleteAuthorHardDeletedBlogs :exec
DELETE FROM blogs.blogs
WHERE author_id = $1;

-- name: DeleteAuthorCache :exec
DELETE FROM blogs.idx_user_author_profile
WHERE author_id = $1;

-- name: MarkAuthorCacheAsDeleted :exec
UPDATE blogs.idx_user_author_profile
SET status = 'deleted', deleted_at = NOW()
WHERE author_id = $1;

-- name: RestoreBlog :exec
UPDATE blogs.blogs
SET status = $1,
deleted_at = null,
deleted_by = null
WHERE blog_id = $2;


-- comments

-- name: CreateComment :one
WITH vals AS (
    SELECT gen_random_uuid() AS u1,
            gen_random_uuid() AS u2
)
INSERT INTO blogs.comments (
    id, blog_id, content, actor_type, actor_id, actor_display_name, parent_comment_id, root_comment_id, depth
)
SELECT u1, $1, $2, $3, $4, $5, $6,
    CASE 
        WHEN $8 = 0 THEN u1 
        ELSE $7
    END AS root_comment_id, 
    $8
FROM vals
RETURNING *;

-- name: GetBlogRootComment :many
WITH child_counts AS (
    SELECT parent_comment_id, COUNT(*) AS cnt
    FROM blogs.comments
    GROUP BY parent_comment_id
)
SELECT
    p.*,
    COALESCE(cc.cnt, 0) AS reply_count
FROM blogs.comments p
LEFT JOIN child_counts cc
    ON cc.parent_comment_id = p.id
WHERE
    p.blog_id = $1
    AND p.status <> 'hidden'
    AND p.depth = 0;

-- name: GetBlogRootCommentWithUserReaction :many
WITH child_counts AS (
    SELECT parent_comment_id, COUNT(*) AS cnt
    FROM blogs.comments
    GROUP BY parent_comment_id
)
SELECT
    p.*,
    COALESCE(cc.cnt, 0) AS reply_count,
    r.type AS reaction_type
FROM blogs.comments p
LEFT JOIN child_counts cc
    ON cc.parent_comment_id = p.id
LEFT JOIN blogs.comment_reactions r
    ON r.comment_id = p.id
    AND r.user_id = $2
WHERE
    p.blog_id = $1
    AND p.status <> 'hidden'
    AND p.depth = 0;

-- name: GetBlogRootCommentCount :one
SELECT COUNT(*) AS total
FROM blogs.comments c
WHERE c.blog_id = $1
    AND c.status <> 'hidden';

-- name: GetCommentsByRootComment :many
SELECT *
FROM blogs.comments
WHERE root_comment_id = $1  AND status <> 'hidden';

-- name: GetCommentsByParentComment :many
SELECT p.*,
    COUNT(c.id) AS child_comment_count
FROM blogs.comments p
LEFT JOIN blogs.comments c
    ON c.parent_comment_id = p.id   AND c.status <> 'hidden'
WHERE p.parent_comment_id = $1 AND p.status <> 'hidden'
GROUP BY p.id;

-- name: GetCommentsByParentCommentUserWithReaction :many
WITH child_counts AS (
    SELECT parent_comment_id, COUNT(*) AS cnt
    FROM blogs.comments
    WHERE status <> 'hidden'
    GROUP BY parent_comment_id
)
SELECT 
    p.*,
    COALESCE(cc.cnt, 0) AS child_comment_count,
    r.type AS reaction_type
FROM blogs.comments p
LEFT JOIN child_counts cc
    ON cc.parent_comment_id = p.id
LEFT JOIN blogs.comment_reactions r
    ON r.comment_id = p.id
    AND r.user_id = $2
WHERE 
    p.parent_comment_id = $1 
    AND p.status <> 'hidden';
-- name: GetCommentByID :one
SELECT *
FROM blogs.comments
WHERE id = $1;

-- name: HideComment :one
UPDATE blogs.comments
SET status = 'hide', updated_at = NOW()
WHERE id = $1
RETURNING COUNT(*);

-- name: DeleteComment :one
UPDATE blogs.comments
SET status = 'delete', updated_at = NOW(), deleted_at = NOW()
WHERE id = $1
RETURNING COUNT(*);

-- name: UpdateBlogReactionType :exec
UPDATE blogs.blog_reactions
    SET type = $1
WHERE id = $2;

-- name: UpdateBlogReactionCount :exec
UPDATE blogs.blogs
    SET like_count = like_count + $1,
    dislike_count = dislike_count + $2
WHERE blog_id = $3;

-- name: UpsertBlogReaction :one
WITH existing AS (
    SELECT type
    FROM blogs.blog_reactions e
    WHERE e.blog_id = $1 AND e.user_id = $2
    FOR UPDATE
),
upsert AS (
    INSERT INTO blogs.blog_reactions (blog_id, user_id, type)
    VALUES ($1, $2, $3)
    ON CONFLICT (blog_id, user_id)
    DO UPDATE SET type = EXCLUDED.type
    RETURNING type AS new_type
)
SELECT
    COALESCE((SELECT type FROM existing), 'none')::VARCHAR(20) AS old_type,
    (SELECT new_type FROM upsert)::VARCHAR(20) AS new_type;

-- name: UpsertCommentReaction :one
WITH existing AS (
    SELECT type
    FROM blogs.comment_reactions e
    WHERE e.comment_id = $1 AND e.user_id = $2
    FOR UPDATE
),
upsert AS (
    INSERT INTO blogs.comment_reactions (comment_id, user_id, type)
    VALUES ($1, $2, $3)
    ON CONFLICT (comment_id, user_id)
    DO UPDATE SET type = EXCLUDED.type
    RETURNING type AS new_type
)
SELECT
    COALESCE((SELECT type FROM existing), 'none')::VARCHAR(20) AS old_type,
    COALESCE((SELECT new_type FROM upsert), 'none')::VARCHAR(20) AS new_type;

-- name: UpdateCommentReactionCount :exec
UPDATE blogs.comments
    SET like_count = like_count + $1,
    dislike_count = dislike_count + $2
WHERE id = $3;

-- name: SyncBlogLikeAndDislike :exec
UPDATE blogs.blogs b
SET
    like_count = COALESCE(x.like_count, 0),
    dislike_count = COALESCE(x.dislike_count, 0)
FROM (
    SELECT
        blog_id,
        COUNT(*) FILTER (WHERE type = 'like' AND status = 'active') AS like_count,
        COUNT(*) FILTER (WHERE type = 'dislike' AND status = 'active') AS dislike_count
    FROM blogs.blog_reactions
    GROUP BY blog_id
) x
WHERE x.blog_id = b.blog_id;

-- name: GetAuthorCacheByUserID :one
SELECT *
FROM blogs.idx_user_author_profile
WHERE user_id = $1;