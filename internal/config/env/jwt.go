package env

import (
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

type jwtEnvConfig struct {
	Secret *string `env:"JWT_SECRET"`
	TTL    *string `env:"JWT_TTL"`
}

type jwtConfig struct {
	raw jwtEnvConfig
}

func NewJwtConfig() (*jwtConfig, error) {
	var raw jwtEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	if raw.Secret == nil && raw.TTL == nil {
		return nil, nil
	}
	return &jwtConfig{raw: raw}, nil
}

func (j *jwtConfig) Secret() string {
	if j.raw.Secret == nil {
		return defaults.JWTSecret
	}
	return *j.raw.Secret
}

func (j *jwtConfig) TTL() time.Duration {
	var ttl time.Duration
	var err error
	if j == nil || j.raw.TTL == nil {
		ttl, err = time.ParseDuration(defaults.JWTTTL)
	} else {
		ttl, err = time.ParseDuration(*j.raw.TTL)
	}
	if err != nil {
		return 24 * time.Hour
	}
	return ttl
}
