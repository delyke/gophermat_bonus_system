package env

import (
	"fmt"
	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host               *string `env:"POSTGRES_HOST"`
	Port               *int    `env:"EXTERNAL_POSTGRES_PORT"`
	User               *string `env:"POSTGRES_USER"`
	Password           *string `env:"POSTGRES_PASSWORD"`
	Database           *string `env:"POSTGRES_DB"`
	DatabaseURI        *string `env:"DATABASE_URI"`
	MigrationDirectory *string `env:"MIGRATION_DIRECTORY"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

// NewPostgresConfig - парсит переменные окружения, относящиеся к Postgres
func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	if raw.Host == nil &&
		raw.User == nil &&
		raw.Port == nil &&
		raw.Database == nil &&
		raw.DatabaseURI == nil &&
		raw.Password == nil &&
		raw.MigrationDirectory == nil {
		return nil, nil
	}
	return &postgresConfig{raw: raw}, nil
}

func (c *postgresConfig) Host() string {
	if c == nil || c.raw.Host == nil {
		return defaults.PgHost
	}
	return *c.raw.Host
}

func (c *postgresConfig) Port() int {
	if c == nil || c.raw.Port == nil {
		return defaults.PgPort
	}
	return *c.raw.Port
}

func (c *postgresConfig) User() string {
	if c == nil || c.raw.User == nil {
		return defaults.PgUser
	}
	return *c.raw.User
}

func (c *postgresConfig) Password() string {
	if c == nil || c.raw.Password == nil {
		return defaults.PgPassword
	}
	return *c.raw.Password
}

func (c *postgresConfig) Database() string {
	if c == nil || c.raw.Database == nil {
		return defaults.PgDatabase
	}
	return *c.raw.Database
}

func (c *postgresConfig) MigrationDirectory() string {
	if c == nil || c.raw.MigrationDirectory == nil {
		return defaults.MigrationDirectory
	}
	return *c.raw.MigrationDirectory
}

func (c *postgresConfig) DatabaseURI() string {
	if c == nil || c.raw.DatabaseURI == nil {
		return defaults.PgDatabaseURI
	}
	return *c.raw.DatabaseURI
}

func (c *postgresConfig) URI() string {
	if c.DatabaseURI() != "" {
		return c.DatabaseURI()
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User(),
		c.Password(),
		c.Host(),
		c.Port(),
		c.Database(),
	)
}
