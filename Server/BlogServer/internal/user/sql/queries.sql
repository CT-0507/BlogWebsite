-- name: GetUserByID :one
SELECT 
    *
FROM users.users 
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: GetDeletedUserByID :one
SELECT 
    user_id, 
    username, 
    password, 
    first_name, 
    last_name, 
    created_at, 
    created_by, 
    updated_at, 
    updated_by,
    deleted_at,
    deleted_by
FROM users.users 
WHERE user_id = $1;

-- name: ListUsers :many
SELECT 
    user_id, 
    username, 
    password, 
    first_name, 
    last_name, 
    created_at, 
    created_by, 
    updated_at, 
    updated_by 
FROM users.users
WHERE deleted_at IS NULL;

-- name: ListWithDeleteUserUsers :many
SELECT user_id, username, password, first_name, last_name, created_at, created_by, updated_at, updated_by FROM users.users;

-- name: CountUserWithEmail :one
SELECT COUNT(*)
FROM users.users u
WHERE u.username = $1;

-- name: CreateUser :one
INSERT INTO users.users(
    username, 
    password, 
    first_name, 
    last_name,
    active,
    role
) VALUES (
    $1, $2, $3, $4, 'normal', $5
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users.users u
WHERE
    u.username = $1
    AND u.active <> 'banned' 
    AND u.deleted_at IS NULL;

-- name: UpdateLastLogout :exec
UPDATE users.users
    SET last_logout = NOW(),
    token_version = token_version + 1
    WHERE user_id = $1;

-- name: UpdateUserData :one
UPDATE users.users
    SET username = $1,
    first_name = $2, 
    last_name = $3
WHERE user_id = $4
RETURNING user_id;

-- name: UpdateUserPassword :one
UPDATE users.users
    SET password = $1
WHERE user_id = $2
RETURNING user_id;

-- name: UpdateUserEmail :one
UPDATE users.users
    SET email = $1
WHERE user_id = $2
RETURNING user_id;

-- name: DeleteUser :one
DELETE FROM users.users
WHERE user_id = $1
RETURNING user_id;