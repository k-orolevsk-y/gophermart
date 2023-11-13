package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         int64     `db:"id" json:"number,string"`
	UserID     uuid.UUID `db:"user_id" json:"-"`
	Status     string    `db:"status" json:"status"`
	Accrual    *float64  `db:"accrual" json:"accrual,omitempty"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}
