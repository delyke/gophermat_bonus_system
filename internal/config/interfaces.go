package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PostgresConfig interface {
	Host() string
	Port() int
	User() string
	Password() string
	Database() string
	DatabaseURI() string
	MigrationDirectory() string
	URI() string
}

type BonusHTTPConfig interface {
	RunAddress() string
	ReadTimeout() time.Duration
}

type AccrualConfig interface {
	SystemAddress() string
}
