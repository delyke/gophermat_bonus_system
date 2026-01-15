package flag

import "github.com/delyke/gophermat_bonus_system/internal/config/defaults"

type loggerEnvConfig struct {
	Level  *string
	AsJSON *bool
}

type loggerConfig struct {
	raw loggerEnvConfig
}

// NewLoggerConfig - парсит флаги и сохраняет их в приватную структуру, доступную только для чтения
func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	flagWasSet := false

	if WasSet("l") {
		raw.Level = logLevel
		flagWasSet = true
	}

	if WasSet("j") {
		raw.AsJSON = asJSON
		flagWasSet = true
	}

	if !flagWasSet {
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
