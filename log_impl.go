package zaplogger

import (
	"context"
)

// GinPlusLoggerImpl Implement the gin-plus Logger interface.
// You can use this logger to replace the default gin-plus logger
type GinPlusLoggerImpl struct{}

func (g *GinPlusLoggerImpl) Info(msg string) {
	Info(msg)
}

func (g *GinPlusLoggerImpl) Warn(msg string) {
	Warn(msg)
}

func (g *GinPlusLoggerImpl) Debug(msg string) {
	Debug(msg)
}

func (g *GinPlusLoggerImpl) Error(msg string) {
	Error(msg)
}

func (g *GinPlusLoggerImpl) ErrorWithCtx(ctx context.Context, msg string) {
	WithContext(ctx).Error(msg)
}
