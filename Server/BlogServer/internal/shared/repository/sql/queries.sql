-- name: GetUserTokenVersionByID :one
SELECT u.token_version
FROM users.users u
WHERE u.user_id = $1;