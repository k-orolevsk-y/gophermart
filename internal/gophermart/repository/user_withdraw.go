package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type (
	pgCategoryUserWithdraw struct {
		db postgres.PgSQL
	}

	UserWithdraw struct {
		ID          uuid.UUID `db:"id"`
		UserID      uuid.UUID `db:"user_id"`
		OrderID     int64     `db:"order_id"`
		Sum         int       `db:"sum"`
		ProcessedAt time.Time `db:"processed_at"`
	}
)

func (pgCUW *pgCategoryUserWithdraw) Create(ctx context.Context, withdraw *UserWithdraw) error {
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

func (pgCUW *pgCategoryUserWithdraw) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]UserWithdraw, error) {
	var userWithdrawals []UserWithdraw
	err := pgCUW.db.SelectContext(ctx, &userWithdrawals, "SELECT * FROM users_withdrawals WHERE user_id = $1", userID)

	return userWithdrawals, err
}
