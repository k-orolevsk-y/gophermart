package models

import (
	"time"

	"github.com/google/uuid"
)

type UserWithdraw struct {
	ID          uuid.UUID `db:"id" json:"-"`
	UserID      uuid.UUID `db:"user_id" json:"-"`
	OrderID     int64     `db:"order_id" json:"order,string"`
	Sum         float64   `db:"sum" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
}
