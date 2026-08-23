package util

import (
	"fmt"
	"os"
	"time"

	"github.com/kidx45/Debter/internal/domain"
	"github.com/spf13/viper"
)

const (
	defaultAccessTokenDuration  = 15 * time.Minute
	defaultRefreshTokenDuration = 7 * 24 * time.Hour
)

func LoadEnv(path string) (domain.Config, error) {
	config := domain.Config{}
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error: %s", err)
		DB_URL := os.Getenv("DB_URL")
		DB_DRIVER_NAME := os.Getenv("DB_DRIVER_NAME")
		PORT := os.Getenv("PORT")
		TOKEN_SYMMETRIC_KEY := os.Getenv("TOKEN_SYMMETRIC_KEY")
		if DB_URL == "" || DB_DRIVER_NAME == "" || PORT == "" || TOKEN_SYMMETRIC_KEY == "" {
			return domain.Config{}, fmt.Errorf("Incomplete Env Configuration")
		}
		config = domain.Config{
			DB_URL:              DB_URL,
			DB_DRIVER_NAME:      DB_DRIVER_NAME,
			PORT:                PORT,
			TOKEN_SYMMETRIC_KEY: TOKEN_SYMMETRIC_KEY,
		}
	} else if err := viper.Unmarshal(&config); err != nil {
		return domain.Config{}, err
	}

	fmt.Print(len(config.TOKEN_SYMMETRIC_KEY))

	if config.ACCESS_TOKEN_DURATION <= 0 {
		config.ACCESS_TOKEN_DURATION = defaultAccessTokenDuration
	}
	if config.REFRESH_TOKEN_DURATION <= 0 {
		config.REFRESH_TOKEN_DURATION = defaultRefreshTokenDuration
	}

	return config, nil
}
