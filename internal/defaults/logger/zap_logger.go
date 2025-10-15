package logger

import (
	"fmt"
	"mockium/pkg/core"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(logLevel, logDir string) (core.Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

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

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return &ZapLogger{logger: logger}, nil
}

type ZapLogger struct {
	logger *zap.Logger
}

func (inst *ZapLogger) Debug(msg string, fields ...any) {
	inst.logger.Debug(msg, inst.toZapFields(fields...)...)
}

func (inst *ZapLogger) Info(msg string, fields ...any) {
	inst.logger.Info(msg, inst.toZapFields(fields...)...)
}

func (inst *ZapLogger) Warn(msg string, fields ...any) {
	inst.logger.Warn(msg, inst.toZapFields(fields...)...)
}

func (inst *ZapLogger) Error(msg string, err error) {
	inst.logger.Error(msg, zap.Error(err))
}

func (inst *ZapLogger) Fatal(msg string, fields ...any) {
	inst.logger.Fatal(msg, inst.toZapFields(fields...)...)
}

// toZapFields преобразует []any в []zap.Field
func (inst *ZapLogger) toZapFields(fields ...any) []zap.Field {
	var zFields []zap.Field

	// ⚠️ самый простой и безопасный вариант — пытаться парсить пары key,value
	// например: logger.Info("msg", "key1", val1, "key2", val2)
	if len(fields)%2 == 0 {
		for i := 0; i < len(fields); i += 2 {
			key, ok := fields[i].(string)
			if !ok {
				continue
			}
			zFields = append(zFields, zap.Any(key, fields[i+1]))
		}
	} else {
		// если что-то не так, пишем всё одним полем
		zFields = append(zFields, zap.Any("fields", fields))
	}

	return zFields
}
