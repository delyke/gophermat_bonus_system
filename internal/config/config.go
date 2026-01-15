package config

import (
	"os"

	"github.com/joho/godotenv"

	defaults "github.com/delyke/gophermat_bonus_system/internal/config/defaults"
	"github.com/delyke/gophermat_bonus_system/internal/config/env"
	flagCfg "github.com/delyke/gophermat_bonus_system/internal/config/flag"
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
}

// Partial — частичная конфигурация: поля nil означают “источник не задавал секцию”.
type partial struct {
	Logger   *LoggerConfig
	Postgres *PostgresConfig
	HTTP     *BonusHTTPConfig
	Accrual  *AccrualConfig
	JWT      *JWTConfig
}

// Load загружает конфигурацию, применяя приоритет:
// defaults -> env -> flags
//
// paths - опциональные пути к .env (если не переданы, будет попытка загрузить ".env").
func Load(paths ...string) error {
	*appConfig = defaultConfig()
	// 0) Парсит флаги один раз
	if err := flagCfg.Parse(); err != nil {
		return err
	}

	// 1) env слой
	envPartial, err := loadEnvPartial(paths...)
	if err != nil {
		return err
	}
	appConfig.apply(envPartial)

	// 2) flags слой (перекрывает env)
	flagPartial, err := loadFlagPartial()
	if err != nil {
		return err
	}
	appConfig.apply(flagPartial)
	return nil
}

// Get возвращает текущую конфигурацию.
func Get() *config {
	return appConfig
}

// apply накладывает Partial поверх текущей конфигурации.
func (c *config) apply(p partial) {
	if p.Logger != nil {
		c.Logger = *p.Logger
	}
	if p.Postgres != nil {
		c.Postgres = *p.Postgres
	}
	if p.Accrual != nil {
		c.Accrual = *p.Accrual
	}
	if p.HTTP != nil {
		c.HTTP = *p.HTTP
	}
	if p.JWT != nil {
		c.JWT = *p.JWT
	}
}

// defaultConfig - место для дефолтов.
func defaultConfig() config {
	var c config
	c.HTTP = defaults.NewHTTPDefaultConfig()
	c.Logger = defaults.NewDefaultLoggerConfig()
	c.Postgres = defaults.NewPostgresDefaultConfig()
	c.Accrual = defaults.NewAccrualDefaultConfig()
	c.JWT = defaults.NewJwtConfig()
	return c
}

func loadFlagPartial() (partial, error) {
	var out partial

	loggerCfg, err := flagCfg.NewLoggerConfig()
	if err != nil {
		return partial{}, err
	}
	if loggerCfg != nil {
		var v LoggerConfig = loggerCfg
		out.Logger = &v
	}

	postgresCfg, err := flagCfg.NewPostgresConfig()
	if err != nil {
		return partial{}, err
	}
	if postgresCfg != nil {
		var v PostgresConfig = postgresCfg
		out.Postgres = &v
	}

	accrualCfg, err := flagCfg.NewAccrualConfig()
	if err != nil {
		return partial{}, err
	}
	if accrualCfg != nil {
		var v AccrualConfig = accrualCfg
		out.Accrual = &v
	}

	httpCfg, err := flagCfg.NewHTTPConfig()
	if err != nil {
		return partial{}, err
	}
	if httpCfg != nil {
		var v BonusHTTPConfig = httpCfg
		out.HTTP = &v
	}

	jwtCfg, err := flagCfg.NewJwtFlagConfig()
	if err != nil {
		return partial{}, err
	}
	if jwtCfg != nil {
		var v JWTConfig = jwtCfg
		out.JWT = &v
	}

	return out, nil
}

func loadEnvPartial(paths ...string) (partial, error) {
	// Пытаемся загрузить .env. Отсутствие файла не считаем ошибкой.
	if len(paths) > 0 {
		if err := godotenv.Load(paths...); err != nil && !os.IsNotExist(err) {
			return partial{}, err
		}
	} else {
		// Если пути не переданы - попробуем стандартный ".env"
		if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
			return partial{}, err
		}
	}

	var out partial

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return partial{}, err
	}
	if loggerCfg != nil {
		var v LoggerConfig = loggerCfg
		out.Logger = &v
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return partial{}, err
	}
	if postgresCfg != nil {
		var v PostgresConfig = postgresCfg
		out.Postgres = &v
	}

	accrualCfg, err := env.NewAccrualConfig()
	if err != nil {
		return partial{}, err
	}
	if accrualCfg != nil {
		var v AccrualConfig = accrualCfg
		out.Accrual = &v
	}

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return partial{}, err
	}
	if httpCfg != nil {
		var v BonusHTTPConfig = httpCfg
		out.HTTP = &v
	}

	jwtCfg, err := env.NewJwtConfig()
	if err != nil {
		return partial{}, err
	}
	if jwtCfg != nil {
		var v JWTConfig = jwtCfg
		out.JWT = &v
	}

	return out, nil
}
