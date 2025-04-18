// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: tokens.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expired_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING token, created_at, updated_at, user_id, expired_at, revoked_at
`

type CreateRefreshTokenParams struct {
	Token     string
	UserID    uuid.UUID
	ExpiredAt time.Time
	RevokedAt sql.NullTime
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken,
		arg.Token,
		arg.UserID,
		arg.ExpiredAt,
		arg.RevokedAt,
	)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiredAt,
		&i.RevokedAt,
	)
	return i, err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT token, created_at, updated_at, user_id, expired_at, revoked_at FROM refresh_tokens WHERE token = $1
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiredAt,
		&i.RevokedAt,
	)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1
RETURNING token, created_at, updated_at, user_id, expired_at, revoked_at
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, revokeRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiredAt,
		&i.RevokedAt,
	)
	return i, err
}
