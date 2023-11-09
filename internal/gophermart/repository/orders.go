package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type (
	pgCategoryOrders struct {
		db postgres.PgSQL
	}

	Order struct {
		ID         int64     `db:"id"`
		UserID     uuid.UUID `db:"user_id"`
		Status     string    `db:"status"`
		Accrual    int       `db:"accrual"`
		UploadedAt time.Time `db:"uploaded_at"`
	}
)

func (pgCO *pgCategoryOrders) Create(ctx context.Context, order *Order) error {
	_, err := pgCO.db.ExecContext(
		ctx,
		"INSERT INTO orders (id, user_id, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5)",
		order.ID, order.UserID, order.Status, order.Accrual, order.UploadedAt,
	)

	return err
}

func (pgCO *pgCategoryOrders) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	var orders []Order
	err := pgCO.db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id = $1", userID)

	return orders, err
}
