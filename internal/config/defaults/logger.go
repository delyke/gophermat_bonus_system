package defaults

type loggerDefaultConfig struct {
	Level  string
	AsJSON bool
}

type loggerConfig struct {
	raw loggerDefaultConfig
}

func NewDefaultLoggerConfig() *loggerConfig {
	var raw loggerDefaultConfig
	raw.Level = LoggerLogLevel
	raw.AsJSON = LoggerAsJSON
	return &loggerConfig{raw: raw}
}

func (ld *loggerConfig) Level() string {
	return ld.raw.Level
}

func (ld *loggerConfig) AsJSON() bool {
	return ld.raw.AsJSON
}
