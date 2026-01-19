package migrator

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
)

func TestGooseLoggerAdapterWrite(t *testing.T) {
	appLogger, err := logger.New("debug", true)
	if err != nil {
		panic(err)
	}
	adapter := &GooseLoggerAdapter{logger: appLogger}
	n, err := adapter.Write([]byte("test log"))

	require.NoError(t, err)
	require.Equal(t, len("test log"), n)
}

func TestNewMigrator(t *testing.T) {
	appLogger, err := logger.New("debug", true)
	if err != nil {
		panic(err)
	}
	m := NewMigrator(&sql.DB{}, "migrations", appLogger)
	require.Equal(t, "migrations", m.migrationsDir)
	require.NotNil(t, m.db)
}

func TestUpReturnsErrorForInvalidDB(t *testing.T) {
	appLogger, err := logger.New("debug", true)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("pgx", "postgres://invalid:invalid@127.0.0.1:1/dbname?sslmode=disable")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	m := NewMigrator(db, t.TempDir(), appLogger)
	err = m.Up(context.Background())

	require.Error(t, err)
}
