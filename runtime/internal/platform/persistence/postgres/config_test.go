package postgres

import (
	"context"
	"strings"
	"testing"
)

func TestConfigFromEnvMapReadsExplicitPostgresSettings(t *testing.T) {
	lookup := mapLookup(map[string]string{
		EnvDatabaseURL: " postgres://user:pass@localhost:5432/vibit ",
		EnvMaxConns:    "8",
		EnvMinConns:    "2",
	})

	cfg, err := ConfigFromEnvMap(lookup)
	if err != nil {
		t.Fatalf("ConfigFromEnvMap() error = %v, want nil", err)
	}

	if cfg.DSN != "postgres://user:pass@localhost:5432/vibit" {
		t.Fatalf("DSN = %q, want trimmed DSN", cfg.DSN)
	}
	if cfg.MaxConns != 8 {
		t.Fatalf("MaxConns = %d, want 8", cfg.MaxConns)
	}
	if cfg.MinConns != 2 {
		t.Fatalf("MinConns = %d, want 2", cfg.MinConns)
	}
}

func TestConfigFromEnvMapRequiresDSN(t *testing.T) {
	_, err := ConfigFromEnvMap(mapLookup(map[string]string{}))
	if err == nil {
		t.Fatal("ConfigFromEnvMap() error = nil, want missing DSN error")
	}
	if !strings.Contains(err.Error(), "DSN") {
		t.Fatalf("ConfigFromEnvMap() error = %v, want DSN error", err)
	}
}

func TestConfigFromEnvMapRejectsInvalidPoolSize(t *testing.T) {
	_, err := ConfigFromEnvMap(mapLookup(map[string]string{
		EnvDatabaseURL: "postgres://localhost/vibit",
		EnvMaxConns:    "not-a-number",
	}))
	if err == nil {
		t.Fatal("ConfigFromEnvMap() error = nil, want invalid max connections error")
	}
	if !strings.Contains(err.Error(), EnvMaxConns) {
		t.Fatalf("ConfigFromEnvMap() error = %v, want max connections error", err)
	}
}

func TestConfigValidateRejectsMinGreaterThanMax(t *testing.T) {
	cfg := Config{
		DSN:      "postgres://localhost/vibit",
		MaxConns: 2,
		MinConns: 3,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want min greater than max error")
	}
	if !strings.Contains(err.Error(), "min connections") {
		t.Fatalf("Validate() error = %v, want min connections error", err)
	}
}

func TestConfigPoolConfigAppliesPoolSettingsWithoutOpeningConnection(t *testing.T) {
	cfg := Config{
		DSN:      "postgres://user:pass@localhost:5432/vibit?sslmode=disable",
		MaxConns: 9,
		MinConns: 1,
	}

	poolConfig, err := cfg.PoolConfig()
	if err != nil {
		t.Fatalf("PoolConfig() error = %v, want nil", err)
	}
	if poolConfig.MaxConns != 9 {
		t.Fatalf("MaxConns = %d, want 9", poolConfig.MaxConns)
	}
	if poolConfig.MinConns != 1 {
		t.Fatalf("MinConns = %d, want 1", poolConfig.MinConns)
	}
	if poolConfig.ConnConfig == nil || poolConfig.ConnConfig.Database != "vibit" {
		t.Fatalf("PoolConfig database = %#v, want vibit", poolConfig.ConnConfig)
	}
}

func TestOpenPoolRejectsInvalidConfigBeforeOpeningConnection(t *testing.T) {
	_, err := OpenPool(context.Background(), Config{})
	if err == nil {
		t.Fatal("OpenPool() error = nil, want missing DSN error")
	}
}

func mapLookup(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
