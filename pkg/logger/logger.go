package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger creates a new Zap logger based on environment
func InitLogger(env string) (*zap.Logger, error) {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return cfg.Build()
}
