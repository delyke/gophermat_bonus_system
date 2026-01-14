-- +goose Up

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        CREATE TYPE order_status AS ENUM (
            'NEW',
            'PROCESSING',
            'INVALID',
            'PROCESSED'
        );
    END IF;
END $$;
-- +goose StatementEnd

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS orders(
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id VARCHAR(64) NOT NULL,
    order_status order_status NOT NULL,
    accrual NUMERIC(14,2) NULL,
    user_uuid UUID NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_uuid)
        REFERENCES users(uuid)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_orders_user_uuid ON orders(user_uuid);

-- +goose Down
DROP TABLE IF EXISTS orders;

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
       DROP TYPE order_status;
    END IF;
END $$;
-- +goose StatementEnd