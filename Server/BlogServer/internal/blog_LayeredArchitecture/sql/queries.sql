-- name: GetBlog :one
SELECT 
    b.blog_id, 
    b.title,
    b.author_id,
    CONCAT(u.first_name, ' ', u.last_name) as author_name,
    u.email,
    b.content,
    b.active,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by 
FROM blogs.blogs b
JOIN users.users u ON u.user_id = b.author_id
WHERE b.blog_id = $1 AND b.deleted_at IS NULL;

-- name: ListAllBlogsBlogs :many
SELECT blog_id, title, content, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by FROM blogs.blogs;

-- name: ListBlogs :many
SELECT 
    b.blog_id,
    u.user_id as author_id,
    u.first_name || u.last_name as author_name,
    b.title, 
    b.content, 
    b.active,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by 
FROM blogs.blogs b
JOIN users.users u ON u.user_id = b.author_id
WHERE b.deleted_at IS NULL;

-- name: CreateBlog :one
INSERT INTO blogs.blogs(
    author_id,
    title,
    content,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5
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