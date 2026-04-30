package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string `env:"BOT_TOKEN"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("./.env"); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg := &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
	}

	if cfg.BotToken == "" {
		return nil, fmt.Errorf("missing required telegram bot variables")
	}

	return cfg, nil
}
