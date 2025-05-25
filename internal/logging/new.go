package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New initializes a zap logger with the specified log level.
// It configures the logger to use a production encoder with ISO8601 time format.
// The log level can be set to "debug", "info", "warn", or "error".
// The logger is configured to output to stderr and does not include caller or stacktrace information.
func New(logLevel string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

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

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(lv),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	return zap.Must(config.Build())
}
