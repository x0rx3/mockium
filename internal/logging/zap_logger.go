package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger initializes a zap logger with the specified log level.
// It configures the logger to use a production encoder with ISO8601 time format.
// The log level can be set to "debug", "info", "warn", or "error".
// The logger is configured to output to stderr and does not include caller or stacktrace information.
func NewZapLogger(logLevel, logDir string) (*zap.Logger, error) {
	var lv zapcore.Level
	switch logLevel {
	case "debug":
		lv = zap.DebugLevel
	case "info":
		lv = zap.InfoLevel
	case "warn":
		lv = zap.WarnLevel
	case "error":
		lv = zap.ErrorLevel
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(lv),
		Development: true,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "Time",
			LevelKey:    "Level",
			MessageKey:  "Message",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.RFC3339TimeEncoder,
		},
		OutputPaths:      []string{"stdout", logDir + "/app.log"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}
