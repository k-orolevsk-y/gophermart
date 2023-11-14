package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type (
	RepositoryCategoryUserWithdraw interface {
		Create(ctx context.Context, withdraw *models.UserWithdraw) error
		GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserWithdraw, error)
		GetWithdrawnSumByUserID(ctx context.Context, userID uuid.UUID) (float64, error)
	}

	pgCategoryUserWithdraw struct {
		db postgres.PgSQL
	}
)

var (
	ErrorInsufficientFunds = errors.New("error insufficient funds in account")
)

func (pgCUW *pgCategoryUserWithdraw) Create(ctx context.Context, withdraw *models.UserWithdraw) error {
	tx, err := pgCUW.db.Beginx()
	if err != nil {
		return err
	}

	var balance float64
	if err = tx.GetContext(ctx, &balance, "SELECT COALESCE(SUM(accrual), 0.0) FROM orders WHERE user_id = $1", withdraw.UserID); err != nil {
		_ = tx.Rollback()
		return err
	}

	if balance < withdraw.Sum {
		_ = tx.Rollback()
		return ErrorInsufficientFunds
	}

	withdraw.ProcessedAt = time.Now()
	id, err := tx.ExecContextWithReturnID(
		ctx,
		"INSERT INTO users_withdrawals (user_id, order_id, sum, processed_at) VALUES ($1, $2, $3, $4)",
		withdraw.UserID, withdraw.OrderID, withdraw.Sum, withdraw.ProcessedAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	withdraw.ID, _ = uuid.Parse(id.(string))
	return tx.Commit()
}

func (pgCUW *pgCategoryUserWithdraw) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserWithdraw, error) {
	userWithdrawals := make([]models.UserWithdraw, 0)
	err := pgCUW.db.SelectContext(ctx, &userWithdrawals, "SELECT * FROM users_withdrawals WHERE user_id = $1", userID)

	return userWithdrawals, err
}

func (pgCUW *pgCategoryUserWithdraw) GetWithdrawnSumByUserID(ctx context.Context, userID uuid.UUID) (float64, error) {
	var sum float64
	err := pgCUW.db.GetContext(ctx, &sum, "SELECT COALESCE(SUM(sum), 0.0) FROM users_withdrawals WHERE user_id = $1", userID)

	return sum, err
}
