-- +goose Up
DROP INDEX IF EXISTS idx_users_by_password;

-- +goose Down
CREATE INDEX IF NOT EXISTS idx_users_by_password ON users(login, password);