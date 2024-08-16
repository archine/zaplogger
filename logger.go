package zaplogger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger wraps the zap.Logger, providing additional configuration options.
type Logger struct {
	*zap.Logger
	conf *Config
}

// defaultLogger is the global instance of the Logger.
var defaultLogger *Logger

// Init initializes the Logger with the provided configuration.
// If the configuration is nil or incomplete, default values are applied.
func Init(conf *Config) error {
	if conf == nil {
		return fmt.Errorf("logger configuration is nil, please check it")
	}
	if conf.Level == "" {
		conf.Level = "debug"
	}
	level, err := zapcore.ParseLevel(conf.Level)
	if err != nil {
		return err
	}
	if conf.Format == "" {
		conf.Format = "console"
	}
	if conf.ApplyEncoder == nil {
		conf.ApplyEncoder = func(format string, ec zapcore.EncoderConfig, conf BasicConfig) (zapcore.Encoder, error) {
			if format == "json" {
				return zapcore.NewJSONEncoder(ec), nil
			}
			return zapcore.NewConsoleEncoder(ec), nil
		}
	}
	if conf.ApplyCores == nil {
		conf.ApplyCores = func(enc zapcore.Encoder, level zapcore.LevelEnabler, conf BasicConfig) (zapcore.Core, error) {
			return zapcore.NewCore(enc, os.Stderr, level), nil
		}
	}
	ec := zapcore.EncoderConfig{
		TimeKey:          "timestamp",
		LevelKey:         "level",
		CallerKey:        "caller",
		MessageKey:       "msg",
		FunctionKey:      zapcore.OmitKey,
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		LineEnding:       zapcore.DefaultLineEnding,
		StacktraceKey:    "stacktrace",
		ConsoleSeparator: conf.ConsoleSeparator,
	}
	if conf.LevelColor && conf.Format == "console" {
		// Apply level color coding when in console format
		ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	if conf.PrintStacktrace {
		conf.Options = append(conf.Options, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}
	encoder, err := conf.ApplyEncoder(conf.Format, ec, conf.BasicConfig) // Apply encoder configuration
	if err != nil {
		return err
	}
	cores, err := conf.ApplyCores(encoder, level, conf.BasicConfig) // Apply core configuration
	if err != nil {
		return err
	}
	l := zap.New(cores, conf.Options...)
	defaultLogger = &Logger{l, conf}
	return nil
}

// DefaultLogger returns the global instance of the Logger.
func DefaultLogger() *Logger {
	return defaultLogger
}

// Info logs an informational message.
func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

// Debug logs a debug-level message.
func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

// Warn logs a warning message.
func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

// Error logs an error message.
func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

// Fatal logs a fatal-level message and exits the application.
func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}

// Panic logs a panic-level message and panics.
func Panic(msg string, fields ...zap.Field) {
	defaultLogger.Panic(msg, fields...)
}
