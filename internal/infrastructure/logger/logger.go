// Package logger provides a zap-based structured logger initialiser.
// In production mode it uses the JSON production preset;
// in development mode it uses the human-readable development preset.
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a configured *zap.Logger.
// logLevel accepts "debug", "info", "warn", "error".
// isProd controls whether to use the production (JSON) or development (console) encoder.
func New(logLevel string, isProd bool) (*zap.Logger, error) {
	level, err := parseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	var cfg zap.Config
	if isProd {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// Do not print stack traces for Warn level (like 404 client errors)
		cfg.DisableStacktrace = true
	}

	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("logger: failed to build: %w", err)
	}

	return logger, nil
}

// parseLevel converts a string log level to a zapcore.Level.
func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("logger: unknown log level %q, defaulting to info", level)
	}
}
