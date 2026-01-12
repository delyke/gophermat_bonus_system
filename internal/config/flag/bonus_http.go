package flag

import (
	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
	"time"
)

type httpEnvConfig struct {
	RunAddress  *string
	ReadTimeout *string
}

type httpConfig struct {
	raw httpEnvConfig
}

func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	flagWasSet := false

	if WasSet("a") {
		raw.RunAddress = runAddress
		flagWasSet = true
	}

	if WasSet("rt") {
		raw.ReadTimeout = readTimeout
		flagWasSet = true
	}

	if !flagWasSet {
		return nil, nil
	}
	return &httpConfig{raw: raw}, nil
}

func (c *httpConfig) RunAddress() string {
	if c == nil || c.raw.RunAddress == nil {
		return defaults.BonusRunAddress
	}
	return *c.raw.RunAddress
}

func (c *httpConfig) ReadTimeout() time.Duration {
	var timeout time.Duration
	var err error
	if c == nil || c.raw.ReadTimeout == nil {
		timeout, err = time.ParseDuration(defaults.BonusReadTimeout)
	} else {
		timeout, err = time.ParseDuration(*c.raw.ReadTimeout)
	}
	if err != nil {
		return 10 * time.Second
	}
	return timeout
}
