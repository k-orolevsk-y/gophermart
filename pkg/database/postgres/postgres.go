package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/pkg/migrations"
)

type PgSQL interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	io.Closer

	Begin() (*sql.Tx, error)
	Beginx() (*sqlx.Tx, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContextWithReturnID(ctx context.Context, query string, args ...interface{}) (interface{}, error)
}

type postgresDatabase struct {
	*sqlx.DB
}

func New(logger *zap.Logger) (PgSQL, error) {
	db, err := sqlx.Open("pgx/v5", config.Config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	if config.Config.MigrationsFlag {
		if err = migrations.ApplyMigrations(db, logger); err != nil {
			return nil, fmt.Errorf("migrations.ApplyMigrations: %w", err)
		}
	}

	return &postgresDatabase{db}, err
}

func (db *postgresDatabase) ExecContextWithReturnID(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	query = fmt.Sprintf("%s RETURNING id", query)

	var id interface{}
	row := db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&id)
	return id, err
}
