package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	defaults "github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

// appConfig - глобальная конфигурация приложения
var appConfig = &config{}

// Config - итоговая конфигурация приложения (все секции обязательны после Load()).
type config struct {
	Logger   LoggerConfig
	Postgres PostgresConfig
	HTTP     BonusHTTPConfig
	Accrual  AccrualConfig
	JWT      JWTConfig
	App      AppConfig
}

// Load загружает конфигурацию, применяя приоритет:
// defaults -> env -> flags
//
// paths - опциональные пути к .env (если не переданы, будет попытка загрузить ".env").
func Load(paths ...string) error {
	// 1) env слой
	if err := loadEnvFile(paths...); err != nil {
		return err
	}

	v := viper.New()
	bindDefaults(v)
	if err := bindEnv(v); err != nil {
		return err
	}

	if err := bindFlags(v); err != nil {
		return err
	}

	*appConfig = buildConfig(v)
	return nil
}

// Get возвращает текущую конфигурацию.
func Get() *config {
	return appConfig
}

func loadEnvFile(paths ...string) error {
	if len(paths) > 0 {
		if err := godotenv.Load(paths...); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func bindDefaults(v *viper.Viper) {
	v.SetDefault("logger.level", defaults.LoggerLogLevel)
	v.SetDefault("logger.as_json", defaults.LoggerAsJSON)
	v.SetDefault("http.run_address", defaults.BonusRunAddress)
	v.SetDefault("http.read_timeout", defaults.BonusReadTimeout)
	v.SetDefault("postgres.host", defaults.PgHost)
	v.SetDefault("postgres.port", defaults.PgPort)
	v.SetDefault("postgres.user", defaults.PgUser)
	v.SetDefault("postgres.password", defaults.PgPassword)
	v.SetDefault("postgres.database", defaults.PgDatabase)
	v.SetDefault("postgres.database_uri", defaults.PgDatabaseURI)
	v.SetDefault("postgres.migration_directory", defaults.MigrationDirectory)
	v.SetDefault("accrual.system_address", defaults.AccrualSystemAddress)
	v.SetDefault("accrual.workers_count", defaults.AccrualWorkersCount)
	v.SetDefault("jwt.secret", defaults.JWTSecret)
	v.SetDefault("jwt.ttl", defaults.JWTTTL)
	v.SetDefault("app.shutdown_timeout", defaults.AppShutdownDuration)
}

func bindEnv(v *viper.Viper) error {
	if err := v.BindEnv("logger.level", "LOGGER_LEVEL"); err != nil {
		return err
	}
	if err := v.BindEnv("logger.as_json", "LOGGER_AS_JSON"); err != nil {
		return err
	}
	if err := v.BindEnv("http.run_address", "RUN_ADDRESS", "HTTP_RUN_ADDRESS"); err != nil {
		return err
	}
	if err := v.BindEnv("http.read_timeout", "HTTP_READ_TIMEOUT"); err != nil {
		return err
	}
	if err := v.BindEnv("postgres.host", "POSTGRES_HOST"); err != nil {
		return err
	}
	if err := v.BindEnv("postgres.port", "EXTERNAL_POSTGRES_PORT", "POSTGRES_PORT"); err != nil {
		return err
	}
	if err := v.BindEnv("postgres.user", "POSTGRES_USER"); err != nil {
		return err
	}
	if err := v.BindEnv("postgres.password", "POSTGRES_PASSWORD"); err != nil {
		return err
	}
	if err := v.BindEnv("postgres.database", "POSTGRES_DB"); err != nil {
		return err
	}
	if err := v.BindEnv("postgres.database_uri", "POSTGRES_URI", "DATABASE_URI"); err != nil {
		return err
	}
	if err := v.BindEnv("postgres.migration_directory", "MIGRATION_DIR", "MIGRATION_DIRECTORY"); err != nil {
		return err
	}
	if err := v.BindEnv("accrual.system_address", "ACCRUAL_SYSTEM_ADDRESS"); err != nil {
		return err
	}
	if err := v.BindEnv("accrual.workers_count", "ACCRUAL_WORKERS_COUNT"); err != nil {
		return err
	}
	if err := v.BindEnv("jwt.secret", "JWT_SECRET"); err != nil {
		return err
	}
	if err := v.BindEnv("jwt.ttl", "JWT_TTL"); err != nil {
		return err
	}

	if err := v.BindEnv("app.shutdown_timeout", "APP_SHUTDOWN_TIMEOUT"); err != nil {
		return err
	}
	return nil
}

func bindFlags(v *viper.Viper) error {
	fs := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	fs.String("l", "", "Log level")
	fs.Bool("j", false, "Use JSON log format")
	fs.StringP("accrual-url", "r", "", "Accrual system address")
	fs.Int("w", 0, "Accrual workers count")
	fs.StringP("address", "a", "", "HTTP listen address")
	fs.String("rt", "", "HTTP read timeout")
	fs.String("pghost", "", "Postgres host")
	fs.Int("pgport", 0, "Postgres port")
	fs.String("pguser", "", "Postgres user")
	fs.String("pgpassword", "", "Postgres password")
	fs.String("pgdb", "", "Postgres database")
	fs.StringP("database-uri", "d", "", "Postgres database URI")
	fs.String("mdir", "", "Directory where migrations")
	fs.String("jwt-secret", "", "Secret used to sign JWT")
	fs.String("jwt-ttl", "", "TTL of JWT")
	fs.String("shutdown-timeout", "", "Shutdown timeout")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	bind := func(key, name string) error {
		flag := fs.Lookup(name)
		if flag == nil {
			return fmt.Errorf("flag: %s not found", name)
		}
		return v.BindPFlag(key, flag)
	}

	if err := bind("logger.level", "l"); err != nil {
		return err
	}

	if err := bind("logger.as_json", "j"); err != nil {
		return err
	}

	if err := bind("accrual.system_address", "accrual-url"); err != nil {
		return err
	}

	if err := bind("accrual.workers_count", "w"); err != nil {
		return err
	}

	if err := bind("http.run_address", "address"); err != nil {
		return err
	}
	if err := bind("http.read_timeout", "rt"); err != nil {
		return err
	}
	if err := bind("postgres.host", "pghost"); err != nil {
		return err
	}
	if err := bind("postgres.port", "pgport"); err != nil {
		return err
	}
	if err := bind("postgres.user", "pguser"); err != nil {
		return err
	}
	if err := bind("postgres.password", "pgpassword"); err != nil {
		return err
	}
	if err := bind("postgres.database", "pgdb"); err != nil {
		return err
	}
	if err := bind("postgres.database_uri", "database-uri"); err != nil {
		return err
	}
	if err := bind("postgres.migration_directory", "mdir"); err != nil {
		return err
	}

	if err := bind("jwt.secret", "jwt-secret"); err != nil {
		return err
	}
	if err := bind("jwt.ttl", "jwt-ttl"); err != nil {
		return err
	}
	if err := bind("app.shutdown_timeout", "shutdown-timeout"); err != nil {
		return err
	}
	return nil
}

func buildConfig(v *viper.Viper) config {
	return config{
		Logger: &loggerConfig{
			level:  v.GetString("logger.level"),
			asJSON: v.GetBool("logger.as_json"),
		},
		HTTP: &httpConfig{
			runAddress:  v.GetString("http.run_address"),
			readTimeout: parseDuration(v.GetString("http.read_timeout"), 10*time.Second),
		},
		Postgres: &postgresConfig{
			host:               v.GetString("postgres.host"),
			port:               v.GetInt("postgres.port"),
			user:               v.GetString("postgres.user"),
			password:           v.GetString("postgres.password"),
			database:           v.GetString("postgres.database"),
			databaseURI:        v.GetString("postgres.database_uri"),
			migrationDirectory: v.GetString("postgres.migration_directory"),
		},
		Accrual: &accrualConfig{
			systemAddress: v.GetString("accrual.system_address"),
			workersCount:  v.GetInt("accrual.workers_count"),
		},
		JWT: &jwtConfig{
			secret: v.GetString("jwt.secret"),
			ttl:    parseDuration(v.GetString("jwt.ttl"), 24*time.Hour),
		},
		App: &aCfg{
			shutdownTimeout: parseDuration(v.GetString("app.shutdown_timeout"), 30*time.Second),
		},
	}
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return value
}

type aCfg struct {
	shutdownTimeout time.Duration
}

func (ac *aCfg) ShutdownTimeout() time.Duration {
	return ac.shutdownTimeout
}

type loggerConfig struct {
	level  string
	asJSON bool
}

func (lc *loggerConfig) Level() string {
	return lc.level
}

func (lc *loggerConfig) AsJSON() bool {
	return lc.asJSON
}

type httpConfig struct {
	runAddress  string
	readTimeout time.Duration
}

func (hc *httpConfig) RunAddress() string {
	return hc.runAddress
}

func (hc *httpConfig) ReadTimeout() time.Duration {
	return hc.readTimeout
}

type postgresConfig struct {
	host               string
	port               int
	user               string
	password           string
	database           string
	databaseURI        string
	migrationDirectory string
}

func (pc *postgresConfig) Host() string {
	return pc.host
}

func (pc *postgresConfig) Port() int {
	return pc.port
}

func (pc *postgresConfig) User() string {
	return pc.user
}

func (pc *postgresConfig) Password() string {
	return pc.password
}

func (pc *postgresConfig) Database() string {
	return pc.database
}

func (pc *postgresConfig) DatabaseURI() string {
	return pc.databaseURI
}

func (pc *postgresConfig) MigrationDirectory() string {
	return pc.migrationDirectory
}

func (pc *postgresConfig) URI() string {
	if pc.DatabaseURI() != "" {
		return pc.DatabaseURI()
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pc.User(), pc.Password(), pc.Host(), pc.Port(), pc.Database())
}

type accrualConfig struct {
	systemAddress string
	workersCount  int
}

func (ac *accrualConfig) SystemAddress() string {
	return ac.systemAddress
}

func (ac *accrualConfig) WorkersCount() int {
	return ac.workersCount
}

type jwtConfig struct {
	secret string
	ttl    time.Duration
}

func (jc *jwtConfig) Secret() string {
	return jc.secret
}

func (jc *jwtConfig) TTL() time.Duration {
	return jc.ttl
}
