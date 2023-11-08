-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    Id uuid NOT NULL DEFAULT gen_random_uuid(),
    Login varchar(100) NOT NULL,
    Password text NOT NULL,
    Balance int NOT NULL,
    CreatedAt timestamptz NOT NULL,
    UpdatedAt timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    Id bigint NOT NULL,
    UserId uuid NOT NULL,
    Status varchar(28) NOT NULL DEFAULT 'NEW',
    Accrual int NOT NULL DEFAULT 0,
    UploadedAt timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS users_withdrawals (
    Id uuid NOT NULL DEFAULT gen_random_uuid(),
    UserId uuid NOT NULL,
    OrderId bigint NOT NULL,
    Sum int NOT NULL DEFAULT 0,
    ProcessedAt timestamptz NOT NULL
);

-- +goose Down