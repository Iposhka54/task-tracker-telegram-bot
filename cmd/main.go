package main

import (
	"context"
	"fmt"
	"os"
	"task-tracker-telegram-bot/internal/app"
	"task-tracker-telegram-bot/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := newLogger(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("config loaded",
		zap.String("log_level", cfg.LogLevel),
	)

	application, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal("Error while initialize application",
			zap.Error(err))
	}

	ctx := context.Background()
	if err := application.Run(ctx); err != nil {
		logger.Error("application error",
			zap.Error(err),
		)
	}
}

func newLogger(level string) (*zap.Logger, error) {
	zapLevel := zap.NewAtomicLevel()
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid LOG_LEVEL %q: %w", level, err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zapLevel
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build()
}
