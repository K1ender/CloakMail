package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV" env-required:"true"`

	Database struct {
		Host     string `env:"DATABASE_HOST" env-required:"true"`
		Port     int    `env:"DATABASE_PORT" env-required:"true"`
		User     string `env:"DATABASE_USER" env-required:"true"`
		Password string `env:"DATABASE_PASSWORD" env-required:"true"`
		Name     string `env:"DATABASE_NAME" env-required:"true"`
	}

	Redis struct {
		Host     string `env:"REDIS_HOST" env-required:"true"`
		Port     int    `env:"REDIS_PORT" env-required:"true"`
		User     string `env:"REDIS_USER" env-required:"true"`
		Password string `env:"REDIS_PASSWORD" env-required:"true"`
	}

	HTTP struct {
		Host string `env:"HTTP_HOST" env-required:"true"`
		Port int    `env:"HTTP_PORT" env-required:"true"`
	}
}

const (
	LocalEnv = "local"
	ProdEnv  = "prod"
)

func MustInit() *Config {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err == nil {
		return cfg
	}

	if err := cleanenv.ReadConfig(".env", cfg); err == nil {
		return cfg
	}

	panic("Failed to load ENV variables")
}
