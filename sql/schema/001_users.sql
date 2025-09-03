-- +goose Up

CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    closed_at TIMESTAMP NULL,
    name TEXT NOT NULL
);

-- +goose Down

    DROP TABLE users;