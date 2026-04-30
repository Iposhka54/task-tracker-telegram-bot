package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string `env:"BOT_TOKEN"`
	LogLevel string `env:"LOG_LEVEL"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("./.env"); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	cfg := &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
		LogLevel: strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))),
	}

	if cfg.BotToken == "" {
		return nil, fmt.Errorf("missing required telegram bot variables")
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return cfg, nil
}
