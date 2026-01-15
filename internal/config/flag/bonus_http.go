package flag

import (
	"time"

	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

type httpEnvConfig struct {
	RunAddress  *string
	ReadTimeout *string
}

type httpConfig struct {
	raw httpEnvConfig
}

func NewHTTPConfig() (*httpConfig, error) {
	raw := &httpEnvConfig{}

	cfg := newFlagConfig(
		raw,
		func() bool {
			if WasSet("a") {
				raw.RunAddress = runAddress
				return true
			}
			return false
		},
		func() bool {
			if WasSet("rt") {
				raw.ReadTimeout = readTimeout
				return true
			}
			return false
		},
	)

	if cfg == nil {
		return nil, nil
	}
	return &httpConfig{raw: *cfg}, nil
}

func (c *httpConfig) RunAddress() string {
	if c == nil || c.raw.RunAddress == nil {
		return defaults.BonusRunAddress
	}
	return *c.raw.RunAddress
}

func (c *httpConfig) ReadTimeout() time.Duration {
	if c == nil {
		return parseDurationOrDefault(nil, defaults.BonusReadTimeout, 10*time.Second)
	}
	return parseDurationOrDefault(c.raw.ReadTimeout, defaults.BonusReadTimeout, 10*time.Second)
}
