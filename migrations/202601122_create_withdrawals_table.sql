-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";


CREATE TABLE IF NOT EXISTS withdrawals(
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    order_id VARCHAR(64) NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    processed_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_withdrawals_user
        FOREIGN KEY (user_uuid)
        REFERENCES users(uuid)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_withdrawals_user_uuid ON withdrawals(user_uuid);

-- +goose Down
DROP TABLE IF EXISTS withdrawals;
