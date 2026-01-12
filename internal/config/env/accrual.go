package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

// accrualEnvConfig — raw-структура для env.
type accrualEnvConfig struct {
	SystemAddress *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// accrualConfig — публичный конфиг слоя env
type accrualConfig struct {
	raw accrualEnvConfig
}

// NewAccrualConfig возвращает:
//
//	(*accrualConfig, nil) - если хотя бы одно поле задано
//	(nil, nil) - если env не содержит ни одного значения
//	(nil, err) - если ошибка парсинга
func NewAccrualConfig() (*accrualConfig, error) {
	var raw accrualEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	// Если env ничего не задал — считаем конфиг отсутствующим
	if raw.SystemAddress == nil {
		return nil, nil
	}

	return &accrualConfig{raw: raw}, nil
}

// SystemAddress безопасно возвращает адрес Accural системы
func (ac *accrualConfig) SystemAddress() string {
	if ac == nil || ac.raw.SystemAddress == nil {
		return defaults.AccrualSystemAddress
	}
	return *ac.raw.SystemAddress
}
