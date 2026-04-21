-- name: GetBlog :one
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
WHERE b.blog_id = $1 AND b.deleted_at IS NULL;

-- name: GetBlogByUrlSlug :one
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
WHERE b.url_slug = $1 AND b.deleted_at IS NULL;

-- name: ListBlogsByAuthor :many
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
    id, blog_id, content, actor_type, actor_id, parent_comment_id, root_comment_id, depth
)
SELECT u1, $1, $2, $3, $4, $5,
    CASE 
        WHEN $6 = TRUE THEN u1 
        ELSE u2
    END AS root_comment_id, 
    $7
FROM vals
RETURNING id;

-- name: GetBlogRootComment :many
SELECT *
FROM blogs.comments
WHERE blog_id = $1 AND status <> 'hidden' AND depth = 0;

-- name: GetCommentsByRootComment :many
SELECT *
FROM blogs.comments
WHERE root_comment_id = $1  AND status <> 'hidden';

-- name: GetCommentsByParentComment :many
SELECT *
FROM blogs.comments
WHERE parent_comment_id = $1 AND status <> 'hidden';

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

-- name: CreateBlogReaction :exec
INSERT INTO blogs.blog_reactions (
    blog_id, user_id, type
) VALUES (
    $1, $2, $3
);

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