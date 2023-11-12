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

func (o Order) CheckValidID() bool {
	if o.ID == 0 {
		return false
	}

	checkSum := func(number int64) int64 {
		var controlSum int64

		for i := 0; number > 0; i++ {
			cur := number % 10

			if i%2 == 0 {
				cur *= 2
				if cur > 9 {
					cur = cur%10 + cur/10
				}
			}

			controlSum += cur
			number /= 10
		}
		return controlSum % 10
	}

	return (o.ID%10+checkSum(o.ID/10))%10 == 0
}
