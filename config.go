package zaplogger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BasicConfig struct {
	// Level log level, default debug.
	// 	Supports: error、info、trace、warn、panic、fetal、debug
	Level string `json:"level" yaml:"level" mapstructure:"level"`

	// LevelColor log level color.
	// 	Note: when the formatter is console, the log level is colored, default true.
	LevelColor bool `json:"level_color" yaml:"level_color" mapstructure:"level_color"`

	// Formatter log format, default console (supports: json、console)
	// 	json: output log in json format.
	Format string `json:"format" yaml:"format" mapstructure:"format"`

	// ConsoleSeparator console separator.
	// Note: when the formatter is console, the separator between the fields, default is "\t".
	ConsoleSeparator string `json:"console_separator" yaml:"console_separator" mapstructure:"console_separator"`

	// PrintStacktrace print stack trace. defalt false.
	PrintStacktrace bool `json:"print_stacktrace" yaml:"print_stacktrace" mapstructure:"print_stacktrace"`
}

// Config log configuration
type Config struct {
	BasicConfig

	// Options zap options.
	// Default zap.AddCaller and zap.AddStacktrace, you can add more options.
	Options []zap.Option

	// ApplyFields apply fields. default nil.
	// When you specify a context through the WithContext function, the function is fired.
	// The fields are added to the log on the return value.
	ApplyFields func(ctx context.Context) []zap.Field

	// ApplyEncoder apply encoder.
	// Default determined based on the format. if the formatter is console, the ConsoleEncoder is used, otherwise, the JsonEncoder is used.
	// If the ApplyEncoder function is specified, the function is fired, and the Encoder is created based on the return value.
	//
	// 	format: the format of the log output, json or console.
	// 	ec: an encoder configuration instance.
	ApplyEncoder func(format string, ec zapcore.EncoderConfig, conf BasicConfig) (zapcore.Encoder, error)

	// ApplyCores apply core.
	// Default create cores based on ApplyEncoder, os.Stdout and Level.
	// If the ApplyCores function is specified, the function is fired, and the Cores are created based on the return value.
	//
	// 	enc: an encoder instance created based on configuration or ApplyEncoder.
	// 	level: a log level instance. based on Level.
	ApplyCores func(enc zapcore.Encoder, level zapcore.LevelEnabler, conf BasicConfig) (zapcore.Core, error)
}
