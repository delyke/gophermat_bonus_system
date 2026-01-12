package defaults

// Logger defaults
var (
	LoggerLogLevel = "debug"
	LoggerAsJson   = true
)

// Accrual flags
var (
	AccrualSystemAddress = "localhost:8081"
)

// Bonus HTTP
var (
	BonusRunAddress  = ":8080"
	BonusReadTimeout = "10s"
)

// Postgres flags
var (
	PgHost             = "localhost"
	PgPort             = 5432
	PgUser             = "bonus-service-user"
	PgPassword         = "bonus-service-password"
	PgDatabase         = "bonus-service"
	PgDatabaseURI      = ""
	MigrationDirectory = "./migrations"
)
