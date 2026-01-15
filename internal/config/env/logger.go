package env

import (
	"github.com/caarlos0/env/v11"

	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

type loggerEnvConfig struct {
	Level  *string `env:"LOGGER_LEVEL"`
	AsJSON *bool   `env:"LOGGER_AS_JSON"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

// NewLoggerConfig - парсит переменные окружения и сохраняет их в приватную структуру, доступную только для чтения
func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	if raw.Level == nil && raw.AsJSON == nil {
		return nil, nil
	}
	return &loggerConfig{raw: raw}, nil
}

func (cfg *loggerConfig) Level() string {
	if cfg == nil || cfg.raw.Level == nil {
		return defaults.LoggerLogLevel
	}
	return *cfg.raw.Level
}

func (cfg *loggerConfig) AsJSON() bool {
	if cfg == nil || cfg.raw.AsJSON == nil {
		return defaults.LoggerAsJSON
	}
	return *cfg.raw.AsJSON
}
