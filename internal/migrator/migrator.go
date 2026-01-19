package migrator

import (
	"context"
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
	logger        *logger.Logger
}

type GooseLoggerAdapter struct {
	logger *logger.Logger
}

func (g *GooseLoggerAdapter) Write(p []byte) (n int, err error) {
	if g.logger != nil {
		g.logger.Info(context.Background(), string(p))
	}

	return len(p), nil
}

func NewMigrator(db *sql.DB, migrationsDir string, appLogger *logger.Logger) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
		logger:        appLogger,
	}
}

func (m *Migrator) Up(ctx context.Context) error {
	if m.logger != nil {
		m.logger.Info(ctx, "Migrator Up")
	}
	goose.SetLogger(log.New(&GooseLoggerAdapter{logger: m.logger}, "", 0))
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}
	return nil
}
