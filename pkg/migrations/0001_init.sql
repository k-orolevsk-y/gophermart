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
    accrual DOUBLE PRECISION DEFAULT NULL,
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

-- Триггер который перед вставкой заказа, проверяет есть ли он в базе
-- и выкидывает исключения в зависимости от user_id

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_order()
    RETURNS TRIGGER AS $$
DECLARE
    existing_user_id uuid;
BEGIN
    SELECT user_id INTO existing_user_id FROM orders WHERE id = NEW.id;

    IF existing_user_id IS NOT NULL THEN
        IF existing_user_id <> NEW.user_id THEN
            RAISE EXCEPTION 'order already created by other user';
        ELSE
            RAISE EXCEPTION 'order already created by this user';
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER before_insert_order
    BEFORE INSERT ON orders
    FOR EACH ROW
EXECUTE FUNCTION check_order();

-- +goose Down