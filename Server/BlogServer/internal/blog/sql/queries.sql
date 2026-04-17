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
    deleted_at = NOW()
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