package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigFile("../../.env")
		_ = viper.ReadInConfig() // Ignore error if file not found, fallback to env vars
	}

	token := viper.GetString("TELEGRAM_TOKEN")
	if token == "" {
		return nil, errors.New("TELEGRAM_TOKEN is not set")
	}

	return &Config{
		TelegramToken: token,
	}, nil
}
