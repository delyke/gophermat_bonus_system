package logger

import (
	"context"

	"go.uber.org/zap"
)

// NoopLogger - пустой логгер, который используется как заглушка до инициализации логгера приложения
type NoopLogger struct{}

func (l *NoopLogger) Info(ctx context.Context, msg string, fields ...zap.Field)  {}
func (l *NoopLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {}
