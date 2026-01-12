package defaults

type loggerDefaultConfig struct {
	Level  string
	AsJson bool
}

type loggerConfig struct {
	raw loggerDefaultConfig
}

func NewDefaultLoggerConfig() *loggerConfig {
	var raw loggerDefaultConfig
	raw.Level = LoggerLogLevel
	raw.AsJson = LoggerAsJson
	return &loggerConfig{raw: raw}
}

func (ld *loggerConfig) Level() string {
	return ld.raw.Level
}
func (ld *loggerConfig) AsJson() bool {
	return ld.raw.AsJson
}
