package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	RedisURL string
	Debug    bool
	AppName  string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		Port:     os.Getenv("PORT"),
		RedisURL: os.Getenv("REDIS_URL"),
		AppName:  os.Getenv("APP_NAME"),
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		return nil, fmt.Errorf("invalid DEBUG value: %w", err)
	}

	cfg.Debug = debug

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	if cfg.RedisURL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}

	if cfg.AppName == "" {
		return fmt.Errorf("APP_NAME is required")
	}

	return nil
}
