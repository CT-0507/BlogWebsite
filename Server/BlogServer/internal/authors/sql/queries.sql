-- name: CreateAuthorProfile :exec
INSERT INTO authors.authors (
    author_id,
    user_id,
    display_name,
    bio,
    avatar,
    slug,
    social_link,
    status,
    email,

    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, 
    $10, $10
);

-- name: ListAuthorProfies :many
SELECT *
FROM authors.authors a
WHERE a.status = $1
    AND ($2 = 'check_null' AND deleted_at IS NULL)
    OR ($2 = 'check_not_null' AND deleted_at IS NOT NULL);

-- name: GetAuthorProfileBySlug :one
SELECT *
FROM authors.authors a
WHERE a.slug = $1 
    AND a.status = $2 
    AND ($3 = 'check_null' AND deleted_at IS NULL)
    OR ($3 = 'check_not_null' AND deleted_at IS NOT NULL);

-- name: GetAuthorProfileByID :one
SELECT *
FROM authors.authors a
WHERE a.author_id = $1 
    AND a.status = $2 
    AND ($3 = 'check_null' AND deleted_at IS NULL)
    OR ($3 = 'check_not_null' AND deleted_at IS NOT NULL);

-- name: UpdateAuthorProfileBasic :exec
UPDATE authors.authors
    SET display_name = $1,
    bio = $2,
    avatar = $3,
    social_link = $4,
    updated_at = NOW(),
    updated_by = $5
WHERE author_id = $6;

-- name: UpdateAuthorProfileEmail :exec
UPDATE authors.authors
    SET email = $1,
    updated_at = NOW(),
    updated_by = $2
WHERE author_id = $3;

-- name: UpdateAuthorStatus :exec
UPDATE authors.authors
    SET status = $1,
    updated_at = NOW(),
    updated_by = $2,
    deleted_at = CASE
        WHEN $1 = 'deleted' THEN NOW()
        ELSE deleted_at
    END,
    deleted_by = CASE
        WHEN $1 = 'deleted' AND deleted_by IS NULL THEN $2
        ELSE deleted_by
    END
WHERE author_id = $3;

-- name: UpdateAuthorProfileDeleteAt :exec
UPDATE authors.authors
    SET status = $1,
    updated_at = NOW(),
    updated_by = $2,
    deleted_at = NOW(),
    deleted_by = $2
WHERE author_id = $3;

-- name: UpdateAuthorSlug :exec
UPDATE authors.authors
SET slug = $1,
updated_at = NOW(),
updated_by = $2
WHERE author_id = $3;

-- name: DeleteAuthorProfile :exec
DELETE FROM authors.authors
WHERE author_id = $1;

-- name: CreateAuthorFollower :exec
INSERT INTO authors.author_followers (
    author_id,
    user_id
) VALUES (
    $1, $2
);

-- name: DeleteAuthorFollower :exec
DELETE FROM authors.author_followers
WHERE author_id = $1 AND  user_id = $2;

-- name: GetAuthorFollowers :many
SELECT f.user_id
FROM authors.author_followers f
JOIN authors.authors a ON a.author_id = f.author_id
WHERE a.slug = $1
ORDER BY a.created_at;

-- name: GetAuthorFollowersByID :many
SELECT f.user_id
FROM authors.author_followers f
JOIN authors.authors a ON a.author_id = f.author_id
WHERE a.author_id = $1
ORDER BY a.created_at;

-- name: GetFollowedAuthors :many
SELECT author_id
FROM authors.author_followers
WHERE user_id = $1
ORDER BY created_at;

-- name: CreateAuthorFeatureBlogs :copyfrom
INSERT INTO authors.author_featured_blogs (
    author_id,
    blog_id,
    position
) VALUES (
    $1, $2, $3
);

-- name: GetAuthorFeatureBlogIDs :many
SELECT blog_id
FROM authors.author_featured_blogs f
JOIN authors.authors a ON author_id
WHERE a.slug = $1 AND a.status = $2;

-- name: UpdateAuthorBlogCount :exec
UPDATE authors.authors
SET blog_count = blog_count + $1
WHERE author_id = $2;

-- name: UpdateAuthorFollowerCount :exec
UPDATE authors.authors
SET follower_count = follower_count + $1
WHERE author_id = $2;
