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

    created_at,
    created_by,
    updated_at,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), $9, NOW(), $10
);

-- name: ListAuthorProfies :many
SELECT *
FROM authors.authors a
WHERE a.status = $1
    AND ($2 = 'check_null' AND deleted_at IS NULL)
    OR ($2 = 'check_not_null' AND deleted_at IS NOT NULL);

-- name: FindAuthorProfileBySlug :one
SELECT *
FROM authors.authors a
WHERE a.slug = $1 
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
    updated_by = $2
WHERE author_id = $3;

-- name: UpdateAuthorProfileDeleteAt :exec
UPDATE authors.authors
    SET status = $1,
    updated_at = NOW(),
    updated_by = $2,
    deleted_at = NOW(),
    deleted_by = $2
WHERE author_id = $3;

-- name: DeleteAuthorProfile :exec
DELETE FROM authors.authors
WHERE author_id = $1;