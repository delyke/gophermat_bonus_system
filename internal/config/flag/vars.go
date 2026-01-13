package flag

// Logger flags
var (
	logLevel = FS.String("l", "", "Log level")
	asJson   = FS.Bool("j", false, "Use JSON log format")
)

// Accrual flags
var (
	systemAddress = FS.String("r", "", "Accrual System address")
)

// Bonus HTTP flags
var (
	runAddress  = FS.String("a", "", "HTTP listen address")
	readTimeout = FS.String("rt", "", "HTTP read timeout")
)

// Postgres flags
var (
	host               = FS.String("pghost", "", "Postgres host")
	port               = FS.Int("pgport", 0, "Postgres port")
	user               = FS.String("pguser", "", "Postgres user")
	password           = FS.String("pgpassword", "", "Postgres password")
	database           = FS.String("pgdb", "", "Postgres database")
	databaseURI        = FS.String("d", "", "Postgres database URI")
	migrationDirectory = FS.String("mdir", "", "Directory where migrations")
)

// JWT flags
var (
	jwtSecret = FS.String("jwt_secret", "", "Secret used to sign JWT")
	jwtTTL    = FS.String("jwt_ttl", "", "TTL of JWT")
)
