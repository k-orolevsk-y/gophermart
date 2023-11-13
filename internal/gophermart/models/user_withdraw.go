package models

import (
	"time"

	"github.com/google/uuid"
)

type UserWithdraw struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	OrderID     int64     `db:"order_id"`
	Sum         float64   `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}
