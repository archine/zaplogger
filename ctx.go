package zaplogger

import (
	"context"
	"go.uber.org/zap"
)

// WithContext returns a new logger with the given context.
// If the logger has no ApplyFields function, it returns the default logger.
// Otherwise, it returns a new logger with the fields returned by the ApplyFields function.
func WithContext(ctx context.Context) *zap.Logger {
	if defaultLogger.conf.ApplyFields == nil {
		return defaultLogger.Logger
	}
	return defaultLogger.With(defaultLogger.conf.ApplyFields(ctx)...)
}
