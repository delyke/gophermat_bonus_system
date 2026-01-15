package defaults

import "time"

type httpDefaultConfig struct {
	RunAddress  string
	ReadTimeout string
}

type httpConfig struct {
	raw httpDefaultConfig
}

func NewHTTPDefaultConfig() *httpConfig {
	var raw httpDefaultConfig
	raw.RunAddress = BonusRunAddress
	raw.ReadTimeout = BonusReadTimeout
	return &httpConfig{raw: raw}
}

func (h *httpConfig) RunAddress() string {
	return h.raw.RunAddress
}

func (h *httpConfig) ReadTimeout() time.Duration {
	readTimeout, err := time.ParseDuration(h.raw.ReadTimeout)
	if err != nil {
		return 10 * time.Second
	}
	return readTimeout
}
