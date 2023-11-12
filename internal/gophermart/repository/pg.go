package repository

import (
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type Pg struct {
	db postgres.PgSQL
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

func (p *Pg) Order() *pgCategoryOrders {
	return &pgCategoryOrders{db: p.db}
}

func (p *Pg) UserWithdraw() *pgCategoryUserWithdraw {
	return &pgCategoryUserWithdraw{db: p.db}
}

func (p *Pg) ParsePgError(err error) *pgconn.PgError {
	var pgError *pgconn.PgError
	if !errors.As(err, &pgError) {
		return &pgconn.PgError{}
	}

	return pgError
}

func (p *Pg) Close() error {
	return p.db.Close()
}
