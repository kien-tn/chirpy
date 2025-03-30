-- +goose Up
-- Add new column password to users table
ALTER TABLE users
ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_password;