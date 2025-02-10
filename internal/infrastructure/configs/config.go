package configs

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment    Environment `env:"ENVIRONMENT" envDefault:"development"`
	HTTPServerAddr string      `env:"HTTP_SERVER_ADDR" envDefault:"0.0.0.0:8080"`
	LogLevel       string      `env:"LOG_LEVEL" envDefault:"debug"`

	PostgresURI string `env:"POSTGRES_URI" envDefault:"postgresql://postgres:postgres@localhost:5432/chat?sslmode=disable"`

	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDb       int    `env:"REDIS_DB"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}
