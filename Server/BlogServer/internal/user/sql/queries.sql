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
    token_version = token_version + 1,
    updated_at = NOW(),
    updated_by = user_id
WHERE user_id = $1;

-- name: UpdateUserData :one
UPDATE users.users
    SET first_name = $1, 
    last_name = $2,
    updated_at = NOW(),
    updated_by = $3
WHERE user_id = $4
RETURNING user_id;

-- name: UpdateUserPassword :one
UPDATE users.users
    SET password = $1,
    updated_at = NOW(),
    updated_by = $3
WHERE user_id = $2
RETURNING user_id;

-- name: UpdateUserEmail :one
UPDATE users.users
    SET email = $1,
    updated_at = NOW(),
    updated_by = $3
WHERE user_id = $2
RETURNING user_id;

-- name: DeleteUser :one
DELETE FROM users.users
WHERE user_id = $1
RETURNING user_id;

-- name: GetUserNotiticationsByID :many
SELECT n.*
FROM users.users u
JOIN users.notifications n ON n.user_id = u.user_id
WHERE n.deleted_at IS NULL;

-- name: CreateNotification :one
INSERT INTO users.notifications (
    user_id,
    content,
    created_by,
    updated_by
)
VALUES ($1, $2, $3, $3)
RETURNING *;

-- name: CreateNotifications :copyfrom
INSERT INTO users.notifications (
    user_id,
    content,
    created_by,
    updated_by
)
VALUES ($1, $2, $3, $3);

-- name: UpdateNotification :exec
UPDATE users.notifications 
    SET is_read = $2,
    updated_at = NOW(),
    updated_by = $3
WHERE notification_id = $1;

-- name: UpdateNotificationStatus :exec
UPDATE users.notifications
    SET is_read = $1,
    updated_at = NOW(),
    updated_by = $2
WHERE notification_id = ANY(sqlc.arg(ids)::int[]);

-- name: MarkUserAsDeleted :exec
UPDATE users.users
SET status = 'deleted',
    updated_at = NOW(),
    updated_by = $1,
    deleted_at = NOW(),
    deleted_by = $1
WHERE users.user_id = $2;

-- name: RestoreUserByID :exec
UPDATE users.users
SET status = $1,
    updated_at = NOW(),
    updated_by = $2,
    deleted_at = NULL,
    deleted_by = NULL
WHERE users.user_id = $3;