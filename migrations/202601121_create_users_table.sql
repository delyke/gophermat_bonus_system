-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    uuid   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(64) NOT NULL,
    password VARCHAR(256) NOT NULL,
    balance NUMERIC(14,2) NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX idx_users_login_unique
    ON users(login);

-- +goose Down
DROP TABLE IF EXISTS users;