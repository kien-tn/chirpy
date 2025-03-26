-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;