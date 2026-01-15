-- +goose Up
CREATE UNIQUE INDEX idx_users_login_unique
    ON users(login);

-- +goose Down
DROP INDEX IF EXISTS idx_users_login_unique;