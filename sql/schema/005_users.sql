-- +goose Up
-- Add new column password to users table
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;