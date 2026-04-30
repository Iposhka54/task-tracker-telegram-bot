package main

import (
	"context"
	"task-tracker-telegram-bot/internal/app"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	app, err := app.NewApp(logger)
	if err != nil {
		logger.Fatal("Error while initialize application",
			zap.Error(err))
	}

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		logger.Error("application error",
			zap.Error(err),
		)
	}
}
