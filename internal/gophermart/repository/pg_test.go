//go:build integration

package repository

import (
	"context"
	"flag"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
)

var (
	databaseDSN = flag.String("d", "", "postgres database dsn")
)

func NewTestPG(t *testing.T) (*pg, error) {
	config.Config.DatabaseURI = *databaseDSN // nolint
	config.Config.MigrationsFlag = true

	db, err := NewPG(zaptest.NewLogger(t))
	if err != nil {
		return nil, err
	}

	return db.(*pg), nil
}

func CloseTestPG(db *pg) error {
	_, err := db.db.ExecContext(
		context.Background(),
		`
			DROP SCHEMA public CASCADE;
			CREATE SCHEMA public;
		`,
	)
	if err != nil {
		return err
	}

	return db.Close()
}
