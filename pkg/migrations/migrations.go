package migrations

import (
	"embed"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed *.sql
var sqlMigrationFiles embed.FS

func ApplyMigrations(db *sqlx.DB, logger *zap.Logger) error {
	goose.SetBaseFS(sqlMigrationFiles)
	goose.SetSequential(true)
	goose.SetLogger(&gooseLogger{logger.Sugar()})

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(db.DB, "."); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}

type gooseLogger struct {
	*zap.SugaredLogger
}

func (gl *gooseLogger) Printf(format string, v ...interface{}) {
	gl.Infof(strings.TrimSuffix(format, "\n"), v...)
}
