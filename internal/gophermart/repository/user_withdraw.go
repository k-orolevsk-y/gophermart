package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	pgCategoryUserWithdraw struct {
		db *sqlx.DB
	}

	UserWithdraw struct {
	}
)

func (pgCUW *pgCategoryUserWithdraw) Create(withdraw UserWithdraw) error {
	// TODO
	panic("TODO")
}

func (pgCUW *pgCategoryUserWithdraw) GetAllByUserID(userID uuid.UUID) ([]UserWithdraw, error) {
	// TODO
	panic("TODO")
}
