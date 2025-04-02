-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expired_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1
RETURNING *;