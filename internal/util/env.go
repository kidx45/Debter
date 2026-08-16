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
		DB_URL_TEST := os.Getenv("DB_URL_TEST")
		if DB_URL == "" || DB_DRIVER_NAME == "" || PORT == "" || DB_URL_TEST == "" {
			return domain.Config{}, fmt.Errorf("Incomplete Env Configuration")
		}
		config = domain.Config{
			DB_URL:         DB_URL,
			DB_DRIVER_NAME: DB_DRIVER_NAME,
			PORT:           PORT,
			DB_URL_TEST:    DB_URL_TEST,
		}
		return config, nil
	}

	if err := viper.Unmarshal(&config); err != nil {
		return domain.Config{}, err
	}
	return config, nil
}
