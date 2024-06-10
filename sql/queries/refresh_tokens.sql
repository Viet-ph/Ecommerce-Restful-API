-- name: SaveTokenToDB :one
INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, revoked)
VALUES ($1, $2, $3, $4, $5, FALSE)
RETURNING *;

-- name: GetValidTokenByUserId :one
SELECT * FROM refresh_tokens
WHERE user_id = $1
AND revoked = FALSE
And expires_at > NOW();

-- name: GetTokenDetail :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens SET revoked = TRUE
WHERE token = $1;

-- name: DeleteTokenFromDB :exec
DELETE FROM refresh_tokens WHERE token = $1;

