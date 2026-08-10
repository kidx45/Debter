package util

import (
	"fmt"
	"os"

	"github.com/kidx45/Debter/internal/domain"
	"github.com/spf13/viper"
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
		if DB_URL == "" || DB_DRIVER_NAME == "" || PORT == "" {
			return domain.Config{}, fmt.Errorf("Incomplete Env Configuration")
		}
		config = domain.Config{
			DB_URL:         DB_URL,
			DB_DRIVER_NAME: DB_DRIVER_NAME,
			PORT:           PORT,
		}
		return config, nil
	}

	if err := viper.Unmarshal(&config); err != nil {
		return domain.Config{}, err
	}
	return config, nil
}
