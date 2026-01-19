package defaults

import "time"

// Logger defaults
var (
	LoggerLogLevel = "debug"
	LoggerAsJSON   = true
)

// Accrual defaults
var (
	AccrualSystemAddress = "http://localhost:8081"
	AccrualWorkersCount  = 4
)

// Bonus HTTP defaults
var (
	BonusRunAddress  = ":8080"
	BonusReadTimeout = "10s"
)

// Postgres defaults
var (
	PgHost             = "localhost"
	PgPort             = 5435
	PgUser             = "bonus-service-user"
	PgPassword         = "bonus-service-password"
	PgDatabase         = "bonus-service"
	PgDatabaseURI      = ""
	MigrationDirectory = "./migrations"
)

// JWT defaults
var (
	JWTSecret = "bonus-service-jwt-secret"
	JWTTTL    = "1d"
)

// App Defaults
var (
	// AppShutdownDuration - maximal time duration for graceful shutdown
	AppShutdownDuration = 30 * time.Second
)
