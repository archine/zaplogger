package zaplogger

import (
	"fmt"
)

// GinPlusLoggerImpl Implement the gin-plus Logger interface.
// You can use this logger to replace the default gin-plus logger
type GinPlusLoggerImpl struct{}

func (g *GinPlusLoggerImpl) Info(msg string, args ...any) {
	Info(fmt.Sprintf(msg, args...))
}

func (g *GinPlusLoggerImpl) Warn(msg string, args ...any) {
	Warn(fmt.Sprintf(msg, args...))
}

func (g *GinPlusLoggerImpl) Debug(msg string, args ...any) {
	Debug(fmt.Sprintf(msg, args...))
}

func (g *GinPlusLoggerImpl) Error(msg string, args ...any) {
	Error(fmt.Sprintf(msg, args...))
}

func (g *GinPlusLoggerImpl) Fatal(msg string, args ...any) {
	Fatal(fmt.Sprintf(msg, args...))
}
