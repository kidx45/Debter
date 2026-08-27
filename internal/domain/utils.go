package domain

import "time"

type Config struct {
	DB_URL                 string        `mapstructure:"DB_URL"`
	PORT                   string        `mapstructure:"PORT"`
	DB_DRIVER_NAME         string        `mapstructure:"DB_DRIVER_NAME"`
	TOKEN_SYMMETRIC_KEY    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}
