package config

import (
	"log/slog"

	"github.com/caarlos0/env/v6"
	"github.com/diianpro/tingerDog/internal/storage/postgres"
)

type Config struct {
	HTTPPort int `config:"HTTP_PORT"`

	Postgres postgres.Config
}

// New initialize Config structure
func New() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		slog.Error("New config: %v", err)
	}
	return &cfg, nil
}
