package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type pgCategoryUserWithdraw struct {
	db postgres.PgSQL
}

var (
	ErrorInsufficientFunds = errors.New("error insufficient funds in account")
)

func (pgCUW *pgCategoryUserWithdraw) Create(ctx context.Context, withdraw *models.UserWithdraw, user *models.User) error {
	tx, err := pgCUW.db.Beginx()
	if err != nil {
		return err
	}

	id, err := tx.ExecContextWithReturnID(
		ctx,
		"INSERT INTO users_withdrawals (user_id, order_id, sum, processed_at) VALUES ($1, $2, $3, $4)",
		withdraw.UserID, withdraw.OrderID, withdraw.Sum, withdraw.ProcessedAt,
	)
	if err != nil {
		return err
	}

	withdraw.ID, _ = uuid.Parse(id.(string))

	if user.Balance < withdraw.Sum {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return ErrorInsufficientFunds
	}

	user.Balance -= withdraw.Sum
	if _, err = tx.ExecContext(ctx, "UPDATE users SET balance = $1 WHERE id = $2", user.Balance, user.ID); err != nil {
		return err
	}

	return tx.Commit()
}

func (pgCUW *pgCategoryUserWithdraw) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserWithdraw, error) {
	var userWithdrawals []models.UserWithdraw
	err := pgCUW.db.SelectContext(ctx, &userWithdrawals, "SELECT * FROM users_withdrawals WHERE user_id = $1", userID)

	return userWithdrawals, err
}

func (pgCUW *pgCategoryUserWithdraw) GetWithdrawnSumByUserID(ctx context.Context, userID uuid.UUID) (float64, error) {
	var sum float64
	err := pgCUW.db.GetContext(ctx, &sum, "SELECT COALESCE(SUM(sum), 0.0) FROM users_withdrawals WHERE user_id = $1", userID)

	return sum, err
}
