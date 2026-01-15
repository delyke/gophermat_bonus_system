package defaults

type accrualDefaultConfig struct {
	SystemAddress string
	WorkersCount  int
}

type accrualConfig struct {
	raw accrualDefaultConfig
}

func NewAccrualDefaultConfig() *accrualConfig {
	var raw accrualDefaultConfig
	raw.SystemAddress = AccrualSystemAddress
	raw.WorkersCount = AccrualWorkersCount
	return &accrualConfig{raw: raw}
}

func (ac *accrualConfig) SystemAddress() string {
	return ac.raw.SystemAddress
}

func (ac *accrualConfig) WorkersCount() int {
	return ac.raw.WorkersCount
}
