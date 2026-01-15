package defaults

import "time"

type jwtDefaultConfig struct {
	Secret string
	TTL    string
}

type jwtConfig struct {
	raw jwtDefaultConfig
}

func NewJwtConfig() *jwtConfig {
	var raw jwtDefaultConfig
	raw.Secret = JWTSecret
	raw.TTL = JWTTTL
	return &jwtConfig{raw: raw}
}

func (j *jwtConfig) Secret() string {
	return j.raw.Secret
}

func (j *jwtConfig) TTL() time.Duration {
	ttl, err := time.ParseDuration(j.raw.TTL)
	if err != nil {
		return 24 * time.Hour
	}
	return ttl
}
