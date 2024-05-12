-- name: CreateUser :one
INSERT INTO users (id, username, email, password, location, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $2
WHERE id = $1;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1;