package flag

import (
	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
	"time"
)

type jwtFlagConfig struct {
	Secret *string
	TTL    *string
}

type jwtConfig struct {
	raw jwtFlagConfig
}

func NewJwtFlagConfig() (*jwtConfig, error) {
	var raw jwtFlagConfig
	flagWasSet := false

	if WasSet("jwt_secret") {
		raw.Secret = jwtSecret
		flagWasSet = true
	}

	if WasSet("jwt_ttl") {
		raw.TTL = jwtTTL
		flagWasSet = true
	}

	if !flagWasSet {
		return nil, nil
	}

	return &jwtConfig{raw: raw}, nil
}

func (j *jwtConfig) Secret() string {
	if j == nil || j.raw.Secret == nil {
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
