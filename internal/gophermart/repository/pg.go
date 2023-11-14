package repository

import (
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type (
	Repository interface {
		User() RepositoryCategoryUser
		Order() RepositoryCategoryOrders
		UserWithdraw() RepositoryCategoryUserWithdraw

		ParsePgError(err error) *pgconn.PgError
		Close() error
	}

	pg struct {
		db postgres.PgSQL
	}
)

func NewPG(logger *zap.Logger) (Repository, error) {
	db, err := postgres.New(logger)
	if err != nil {
		return nil, err
	}

	return &pg{
		db: db,
	}, nil
}

func (p *pg) User() RepositoryCategoryUser {
	return &pgCategoryUser{db: p.db}
}

func (p *pg) Order() RepositoryCategoryOrders {
	return &pgCategoryOrders{db: p.db}
}

func (p *pg) UserWithdraw() RepositoryCategoryUserWithdraw {
	return &pgCategoryUserWithdraw{db: p.db}
}

func (p *pg) ParsePgError(err error) *pgconn.PgError {
	var pgError *pgconn.PgError
	if !errors.As(err, &pgError) {
		return &pgconn.PgError{}
	}

	return pgError
}

func (p *pg) Close() error {
	return p.db.Close()
}
