package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type pgCategoryUserWithdraw struct {
	db postgres.PgSQL
}

func (pgCUW *pgCategoryUserWithdraw) Create(ctx context.Context, withdraw *models.UserWithdraw) error {
	id, err := pgCUW.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO users_withdrawals (user_id, order_id, sum, processed_at) VALUES ($1, $2, $3, $4)",
		withdraw.UserID, withdraw.OrderID, withdraw.Sum, withdraw.ProcessedAt,
	)
	if err != nil {
		return err
	}

	withdraw.ID, _ = uuid.Parse(id.(string))
	return err
}

func (pgCUW *pgCategoryUserWithdraw) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserWithdraw, error) {
	var userWithdrawals []models.UserWithdraw
	err := pgCUW.db.SelectContext(ctx, &userWithdrawals, "SELECT * FROM users_withdrawals WHERE user_id = $1", userID)

	return userWithdrawals, err
}
