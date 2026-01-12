package flag

import "github.com/delyke/gophermat_bonus_system/internal/config/defaults"

type accrualEnvConfig struct {
	SystemAddress *string
}

type accrualConfig struct {
	raw accrualEnvConfig
}

func NewAccrualConfig() (*accrualConfig, error) {
	var raw accrualEnvConfig
	flagWasSet := false

	if WasSet("r") {
		raw.SystemAddress = systemAddress
		flagWasSet = true
	}

	if !flagWasSet {
		return nil, nil
	}
	return &accrualConfig{raw: raw}, nil
}

func (ac *accrualConfig) SystemAddress() string {
	if ac == nil && ac.raw.SystemAddress == nil {
		return defaults.AccrualSystemAddress
	}
	return *ac.raw.SystemAddress
}
