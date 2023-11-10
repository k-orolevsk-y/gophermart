package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type pgCategoryOrders struct {
	db postgres.PgSQL
}

func (pgCO *pgCategoryOrders) Create(ctx context.Context, order *models.Order) error {
	_, err := pgCO.db.ExecContext(
		ctx,
		"INSERT INTO orders (id, user_id, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5)",
		order.ID, order.UserID, order.Status, order.Accrual, order.UploadedAt,
	)

	return err
}

func (pgCO *pgCategoryOrders) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := pgCO.db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id = $1", userID)

	return orders, err
}
