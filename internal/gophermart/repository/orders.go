package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type pgCategoryOrders struct {
	db postgres.PgSQL
}

const (
	ErrorOrderByThisUser  = "order already created by this user"
	ErrorOrderByOtherUser = "order already created by other user"
)

func (pgCO *pgCategoryOrders) Create(ctx context.Context, order *models.Order) error {
	order.UploadedAt = time.Now()

	_, err := pgCO.db.ExecContext(
		ctx,
		"INSERT INTO orders (id, user_id, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5)",
		order.ID, order.UserID, order.Status, order.Accrual, order.UploadedAt,
	)

	return err
}

func (pgCO *pgCategoryOrders) Edit(ctx context.Context, order *models.Order) error {
	_, err := pgCO.db.ExecContext(
		ctx,
		"UPDATE orders SET status = $1, accrual = $2 WHERE id = $3",
		order.Status, order.Accrual, order.ID,
	)

	return err
}

func (pgCO *pgCategoryOrders) GetAccrualSumByUserID(ctx context.Context, userID uuid.UUID) (float64, error) {
	var sum float64
	err := pgCO.db.GetContext(ctx, &sum, "SELECT COALESCE(SUM(accrual), 0.0) FROM orders WHERE user_id = $1", userID)

	return sum, err
}

func (pgCO *pgCategoryOrders) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	err := pgCO.db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id = $1", userID)

	return orders, err
}
