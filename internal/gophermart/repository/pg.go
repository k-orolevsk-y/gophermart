package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type Pg struct {
	db *sqlx.DB
}

func NewPG(logger *zap.Logger) (*Pg, error) {
	db, err := postgres.New(logger)
	if err != nil {
		return nil, err
	}

	return &Pg{
		db: db,
	}, nil
}

func (p *Pg) User() *pgCategoryUser {
	return &pgCategoryUser{db: p.db}
}
