package logger

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

// Logger обертка над zap.Logger с enrich поддержкой контекста
type Logger struct {
	zapLogger    *zap.Logger
	dynamicLevel zap.AtomicLevel
}

// New создает новый logger с заданными настройками
func New(levelStr string, asJSON bool) (*Logger, error) {
	dynamicLevel := zap.NewAtomicLevelAt(parseLevel(levelStr))

	encoderCfg := buildProductionEncoderConfig()

	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		dynamicLevel,
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	return &Logger{
		zapLogger:    zapLogger,
		dynamicLevel: dynamicLevel,
	}, nil
}

// NewNop создаёт no-op логгер.
func NewNop() *Logger {
	return &Logger{
		zapLogger: zap.NewNop(),
	}
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",                 // время
		LevelKey:       "level",                     // уровень логирования
		NameKey:        "logger",                    // имя логгера, если используется
		CallerKey:      "caller",                    // откуда вызван лог
		MessageKey:     "message",                   // текст сообщения
		StacktraceKey:  "stacktrace",                // стектрейс для ошибок
		LineEnding:     zapcore.DefaultLineEnding,   // перенос строки
		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // читаемый ISO 8601 формат
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // полный caller
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// SetLevel динамически меняет уровень логирования
func (l *Logger) SetLevel(levelStr string) {
	if l == nil || l.dynamicLevel == (zap.AtomicLevel{}) {
		return
	}
	l.dynamicLevel.SetLevel(parseLevel(levelStr))
}

// Sync сбрасывает буферы логгера
func (l *Logger) Sync() error {
	if l == nil || l.zapLogger == nil {
		return nil
	}

	return l.zapLogger.Sync()
}

// With создает новый enrich-aware логгер с дополнительными полями
func (l *Logger) With(fields ...zap.Field) *Logger {
	if l == nil {
		return nil
	}
	return &Logger{
		zapLogger:    l.zapLogger.With(fields...),
		dynamicLevel: l.dynamicLevel,
	}
}

// WithContext создает enrich-aware логгер с контекстом
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if l == nil {
		return NewNop()
	}
	return &Logger{
		zapLogger:    l.zapLogger.With(fieldsFromContext(ctx)...),
		dynamicLevel: l.dynamicLevel,
	}
}

// Debug enrich-aware debug log.
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

// Info enrich-aware info log.
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

// Warn enrich-aware warn log.
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

// Error enrich-aware error log.
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

// Fatal enrich-aware fatal log.
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

// parseLevel конвертирует строковый уровень в zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// fieldsFromContext вытаскивает enrich-поля из контекста
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}
