package flag

import (
	"fmt"

	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

type postgresFlagConfig struct {
	Host               *string
	Port               *int
	User               *string
	Password           *string
	Database           *string
	DatabaseURI        *string
	MigrationDirectory *string
}

type postgresConfig struct {
	raw postgresFlagConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresFlagConfig
	flagWasSet := false

	if WasSet("pghost") {
		raw.Host = host
		flagWasSet = true
	}

	if WasSet("pgport") {
		raw.Port = port
		flagWasSet = true
	}

	if WasSet("pguser") {
		raw.User = user
		flagWasSet = true
	}

	if WasSet("pgpassword") {
		raw.Password = password
		flagWasSet = true
	}

	if WasSet("pgdb") {
		raw.Database = database
		flagWasSet = true
	}

	if WasSet("d") {
		raw.DatabaseURI = databaseURI
		flagWasSet = true
	}

	if WasSet("mdir") {
		raw.MigrationDirectory = migrationDirectory
		flagWasSet = true
	}

	if !flagWasSet {
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
