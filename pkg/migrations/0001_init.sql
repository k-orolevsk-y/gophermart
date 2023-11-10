-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    login varchar(100) NOT NULL UNIQUE,
    password text NOT NULL,
    balance DOUBLE PRECISION NOT NULL,
    created_at timestamptz NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_login ON users (login);

CREATE TABLE IF NOT EXISTS orders (
    id bigint PRIMARY KEY  NOT NULL UNIQUE,
    user_id uuid NOT NULL,
    status varchar(28) NOT NULL DEFAULT 'NEW',
    accrual DOUBLE PRECISION NOT NULL DEFAULT 0,
    uploaded_at timestamptz NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_orders_userId ON orders (user_id);

CREATE TABLE IF NOT EXISTS users_withdrawals (
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    order_id bigint NOT NULL,
    sum int NOT NULL DEFAULT 0,
    processed_at timestamptz NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_withdrawals_userId_orderId ON users_withdrawals (user_id, order_id);

-- +goose Down