package domain

type Config struct {
	DB_URL                 string        `mapstructure:"DB_URL"`
	PORT                   string        `mapstructure:"PORT"`
	DB_DRIVER_NAME         string        `mapstructure:"DB_DRIVER_NAME"`
}

