package flag

import (
	"time"

	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

type jwtFlagConfig struct {
	Secret *string
	TTL    *string
}

type jwtConfig struct {
	raw jwtFlagConfig
}

func NewJwtFlagConfig() (*jwtConfig, error) {
	raw := &jwtFlagConfig{}

	cfg := newFlagConfig(
		raw,
		func() bool {
			if WasSet("jwt_secret") {
				raw.Secret = jwtSecret
				return true
			}
			return false
		},
		func() bool {
			if WasSet("jwt_ttl") {
				raw.TTL = jwtTTL
				return true
			}
			return false
		},
	)

	if cfg == nil {
		return nil, nil
	}

	return &jwtConfig{raw: *cfg}, nil
}

func (j *jwtConfig) Secret() string {
	if j == nil || j.raw.Secret == nil {
		return defaults.JWTSecret
	}
	return *j.raw.Secret
}

func (j *jwtConfig) TTL() time.Duration {
	if j == nil {
		return parseDurationOrDefault(nil, defaults.JWTTTL, 24*time.Hour)
	}
	return parseDurationOrDefault(j.raw.TTL, defaults.JWTTTL, 24*time.Hour)
}
