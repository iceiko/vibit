package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	EnvDatabaseURL = "VIBIT_POSTGRES_DSN"
	EnvMaxConns    = "VIBIT_POSTGRES_MAX_CONNS"
	EnvMinConns    = "VIBIT_POSTGRES_MIN_CONNS"
)

type Config struct {
	DSN      string
	MaxConns int32
	MinConns int32
}

func ConfigFromEnv() (Config, error) {
	return ConfigFromEnvMap(os.LookupEnv)
}

func ConfigFromEnvMap(lookup func(string) (string, bool)) (Config, error) {
	if lookup == nil {
		return Config{}, errors.New("postgres config: environment lookup is required")
	}

	cfg := Config{}
	if value, ok := lookup(EnvDatabaseURL); ok {
		cfg.DSN = strings.TrimSpace(value)
	}
	if value, ok := lookup(EnvMaxConns); ok {
		parsed, err := parseOptionalPositiveInt32(EnvMaxConns, value)
		if err != nil {
			return Config{}, err
		}
		cfg.MaxConns = parsed
	}
	if value, ok := lookup(EnvMinConns); ok {
		parsed, err := parseOptionalPositiveInt32(EnvMinConns, value)
		if err != nil {
			return Config{}, err
		}
		cfg.MinConns = parsed
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.DSN) == "" {
		return errors.New("postgres config: DSN is required")
	}
	if c.MaxConns < 0 {
		return errors.New("postgres config: max connections must not be negative")
	}
	if c.MinConns < 0 {
		return errors.New("postgres config: min connections must not be negative")
	}
	if c.MaxConns > 0 && c.MinConns > c.MaxConns {
		return errors.New("postgres config: min connections must not exceed max connections")
	}
	return nil
}

func (c Config) PoolConfig() (*pgxpool.Config, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(strings.TrimSpace(c.DSN))
	if err != nil {
		return nil, fmt.Errorf("postgres config: parse DSN: %w", err)
	}
	if c.MaxConns > 0 {
		poolConfig.MaxConns = c.MaxConns
	}
	if c.MinConns > 0 {
		poolConfig.MinConns = c.MinConns
	}
	return poolConfig, nil
}

func OpenPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolConfig, err := cfg.PoolConfig()
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres config: open pool: %w", err)
	}
	return pool, nil
}

func parseOptionalPositiveInt32(name string, value string) (int32, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("postgres config: %s must be an integer", name)
	}
	if parsed < 0 {
		return 0, fmt.Errorf("postgres config: %s must not be negative", name)
	}
	return int32(parsed), nil
}
