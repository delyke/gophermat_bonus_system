package env

import (
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

// httpEnvConfig - raw-структура для env.
type httpEnvConfig struct {
	RunAddress  *string `env:"RUN_ADDRESS"`
	ReadTimeout *string `env:"HTTP_READ_TIMEOUT"`
}

type httpConfig struct {
	raw httpEnvConfig
}

// NewHTTPConfig возвращает:
// (*httpConfig, nil) - если хотя-бы одно поле задано
// (nil, nil) - если env не содержит ни одного значения
// (nil, err) - если ошибка парсинга
func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	if raw.RunAddress == nil && raw.ReadTimeout == nil {
		return nil, nil
	}
	return &httpConfig{raw: raw}, nil
}

// RunAddress - безопасно возвращает адрес сервера, на котором будет слушаться приложение
func (h *httpConfig) RunAddress() string {
	if h == nil || h.raw.RunAddress == nil {
		return defaults.BonusRunAddress
	}
	return *h.raw.RunAddress
}

// ReadTimeout - безопасно возвращает таймаут
func (h *httpConfig) ReadTimeout() time.Duration {
	var timeout time.Duration
	var err error
	if h == nil || h.raw.ReadTimeout == nil {
		timeout, err = time.ParseDuration(defaults.BonusReadTimeout)
	} else {
		timeout, err = time.ParseDuration(*h.raw.ReadTimeout)
	}
	if err != nil {
		return 10 * time.Second
	}
	return timeout
}
