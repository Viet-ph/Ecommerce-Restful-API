-- name: CreateUser :one
INSERT INTO users (id, username, email, password, location, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1) AS exists;