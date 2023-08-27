package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	Address        string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI    string `env:"DATABASE_URI"`
	LogLevel       string `env:"LOG_LEVEL"`
	SecretKey      string `env:"KEY"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	cfg.parseFlags()
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse envs: %w", err)
	}

	return cfg, nil
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.Address, "a", ":8181", "run address")
	flag.StringVar(&c.AccrualAddress, "r", "localhost:8080", "accrual system address")
	flag.StringVar(&c.DatabaseURI, "d", "", "database connection uri")
	flag.StringVar(&c.LogLevel, "l", "debug", "log level")
	flag.StringVar(&c.SecretKey, "k", "secret", "signing key")

	flag.Parse()
}
