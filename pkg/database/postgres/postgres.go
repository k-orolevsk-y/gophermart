package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/pkg/migrations"
)

func New(logger *zap.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx/v5", config.Config.DatabaseURI)
	if err != nil {
		return nil, err
	}

	if config.Config.MigrationsFlag {
		if err = migrations.ApplyMigrations(db, logger); err != nil {
			return nil, fmt.Errorf("migrations.ApplyMigrations: %w", err)
		}
	}

	return db, err
}
