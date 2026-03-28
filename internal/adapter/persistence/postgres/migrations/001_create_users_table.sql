-- +goose Up
CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(255) NULL,
    password   VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    avatar     VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users