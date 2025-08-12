package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() *zap.Logger {
	env := os.Getenv("APP_ENV")
	var logger *zap.Logger
	var err error

	if env == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err = cfg.Build()
	}

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	return logger
}