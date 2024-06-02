package ginplus

import (
	"fmt"
	"github.com/archine/zaplogger"
)

// Logger Implement the gin-plus Logger interface.
// You can use this logger to replace the default gin-plus logger
type Logger struct{}

func (l *Logger) Info(msg string, args ...any) {
	zaplogger.Info(fmt.Sprintf(msg, args...))
}

func (l *Logger) Warn(msg string, args ...any) {
	zaplogger.Warn(fmt.Sprintf(msg, args...))
}

func (l *Logger) Debug(msg string, args ...any) {
	zaplogger.Debug(fmt.Sprintf(msg, args...))
}

func (l *Logger) Error(msg string, args ...any) {
	zaplogger.Error(fmt.Sprintf(msg, args...))
}

func (l *Logger) Fatal(msg string, args ...any) {
	zaplogger.Fatal(fmt.Sprintf(msg, args...))
}
