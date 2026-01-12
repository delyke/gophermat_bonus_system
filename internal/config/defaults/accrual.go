package defaults

type accrualDefaultConfig struct {
	SystemAddress string
}

type accrualConfig struct {
	raw accrualDefaultConfig
}

func NewAccrualDefaultConfig() *accrualConfig {
	var raw accrualDefaultConfig
	raw.SystemAddress = AccrualSystemAddress
	return &accrualConfig{raw: raw}
}

func (ac *accrualConfig) SystemAddress() string {
	return ac.raw.SystemAddress
}
