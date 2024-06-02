package zaplogger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger wrap the zap logger
type Logger struct {
	*zap.Logger
	conf *Config
}

// defaultLogger default zap logger wrapper
var defaultLogger *Logger

func Init(conf *Config) error {
	if conf == nil {
		return fmt.Errorf("the loggers configuration is nil, please check it")
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
	if conf.ApplyFields == nil {
		conf.ApplyFields = func(ctx context.Context) []zap.Field {
			return nil
		}
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
			return zapcore.NewCore(enc, os.Stdout, level), nil
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
		// level color
		ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	if conf.PrintStacktrace {
		conf.Options = append(conf.Options, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		conf.Options = append(conf.Options, zap.AddCaller())
	}
	encoder, err := conf.ApplyEncoder(conf.Format, ec, conf.BasicConfig) // apply encoder
	if err != nil {
		return err
	}
	cores, err := conf.ApplyCores(encoder, level, conf.BasicConfig) // apply cores
	if err != nil {
		return err
	}
	l := zap.New(cores, conf.Options...)
	defaultLogger = &Logger{l, conf}
	return nil
}

// DefaultLogger Get the default logger wrapper
func DefaultLogger() *Logger {
	return defaultLogger
}

// Info Output info level log
func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

// Debug Debug Output debug level log
func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

// Warn Output warn level log
func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

// Error Output error level log
func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

// Fatal Output fatal level log
func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}

// Panic Output panic level log
func Panic(msg string, fields ...zap.Field) {
	defaultLogger.Panic(msg, fields...)
}
