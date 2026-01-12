package defaults

import "fmt"

type postgresDefaultConfig struct {
	Host               string
	Port               int
	User               string
	Password           string
	Database           string
	DatabaseURI        string
	MigrationDirectory string
}

type postgresConfig struct {
	raw postgresDefaultConfig
}

func NewPostgresDefaultConfig() *postgresConfig {
	var raw postgresDefaultConfig
	raw.Host = PgHost
	raw.Port = PgPort
	raw.User = PgUser
	raw.Password = PgPassword
	raw.Database = PgDatabase
	raw.DatabaseURI = PgDatabaseURI
	raw.MigrationDirectory = MigrationDirectory
	return &postgresConfig{raw: raw}
}

func (cd *postgresConfig) Host() string {
	return cd.raw.Host
}
func (cd *postgresConfig) Port() int {
	return cd.raw.Port
}
func (cd *postgresConfig) User() string {
	return cd.raw.User
}
func (cd *postgresConfig) Password() string {
	return cd.raw.Password
}
func (cd *postgresConfig) Database() string {
	return cd.raw.Database
}
func (cd *postgresConfig) DatabaseURI() string {
	return cd.raw.DatabaseURI
}
func (cd *postgresConfig) MigrationDirectory() string {
	return cd.raw.MigrationDirectory
}
func (cd *postgresConfig) URI() string {
	if cd.DatabaseURI() != "" {
		return cd.DatabaseURI()
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cd.User(),
		cd.Password(),
		cd.Host(),
		cd.Port(),
		cd.Database(),
	)
}
