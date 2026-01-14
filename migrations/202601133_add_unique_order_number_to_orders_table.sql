-- +goose Up
CREATE UNIQUE INDEX idx_orders_number_unique
    ON orders(order_id);

-- +goose Down
DROP INDEX IF EXISTS idx_orders_number_unique;