// Copyright 2025 Mykhailo Bobrovskyi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configs

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
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

	Version string
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	fileVersion, err := os.ReadFile("VERSION")
	if err != nil {
		return nil, fmt.Errorf("error on reading VERSION file: %w", err)
	}

	c.Version = string(fileVersion)

	return c, nil
}
