package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	pgCategoryOrders struct {
		db *sqlx.DB
	}

	Order struct {
	}
)

func (pgCO *pgCategoryOrders) Create(order Order) error {
	// TODO
	panic("TODO")
}

func (pgCO *pgCategoryOrders) GetAllByUserID(userID uuid.UUID) ([]Order, error) {
	// TODO
	panic("TODO")
}
