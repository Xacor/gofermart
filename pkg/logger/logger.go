package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.Logger {
	lvl := zap.NewAtomicLevel()

	switch strings.ToLower(level) {
	case "error":
		lvl.SetLevel(zapcore.ErrorLevel)
	case "warn":
		lvl.SetLevel(zapcore.WarnLevel)
	case "info":
		lvl.SetLevel(zapcore.InfoLevel)
	case "debug":
		lvl.SetLevel(zapcore.DebugLevel)
	default:
		lvl.SetLevel(zapcore.InfoLevel)
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = lvl

	zl := zap.Must(cfg.Build())
	return zl
}
