package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.App.Port == "" {
		return fmt.Errorf("port is required")
	}
	if c.App.PublicAPIURL == "" {
		return fmt.Errorf("public api url is required")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("database url is required")
	}
	if c.Database.MaxConns < 1 {
		return fmt.Errorf("database max connections must be positive")
	}
	if c.Tunnel.MaxConnections < 1 {
		return fmt.Errorf("tunnel max connections must be positive")
	}
	return nil
}
