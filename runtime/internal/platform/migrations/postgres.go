package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/pressly/goose/v3"
)

var ErrNoMigrationSources = errors.New("postgres migrations: migration sources are required")

type MigrationSource struct {
	Type    string
	Path    string
	Version int64
}

type MigrationStatus struct {
	Source    MigrationSource
	State     string
	AppliedAt time.Time
}

type MigrationResult struct {
	Source    MigrationSource
	Direction string
	Duration  time.Duration
	Empty     bool
}

type PostgresRunner struct {
	provider *goose.Provider
}

func NewPostgresRunner(db *sql.DB, sourceFS fs.FS) (*PostgresRunner, error) {
	if db == nil {
		return nil, errors.New("postgres migrations: database handle is required")
	}
	if sourceFS == nil {
		return nil, errors.New("postgres migrations: source filesystem is required")
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		sourceFS,
		goose.WithDisableGlobalRegistry(true),
	)
	if err != nil {
		return nil, wrapProviderError("create provider", err)
	}
	return &PostgresRunner{provider: provider}, nil
}

func NewPostgresRunnerFromDir(db *sql.DB, sourceDir string) (*PostgresRunner, error) {
	sourceDir = strings.TrimSpace(sourceDir)
	if sourceDir == "" {
		return nil, errors.New("postgres migrations: source directory is required")
	}
	return NewPostgresRunner(db, os.DirFS(sourceDir))
}

func (r *PostgresRunner) Sources() ([]MigrationSource, error) {
	provider, err := r.requireProvider()
	if err != nil {
		return nil, err
	}
	sources := provider.ListSources()
	results := make([]MigrationSource, 0, len(sources))
	for _, source := range sources {
		results = append(results, convertSource(source))
	}
	return results, nil
}

func (r *PostgresRunner) Status(ctx context.Context) ([]MigrationStatus, error) {
	provider, err := r.requireProvider()
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	statuses, err := provider.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres migrations: status: %w", err)
	}
	results := make([]MigrationStatus, 0, len(statuses))
	for _, status := range statuses {
		if status == nil {
			continue
		}
		results = append(results, MigrationStatus{
			Source:    convertSource(status.Source),
			State:     string(status.State),
			AppliedAt: status.AppliedAt,
		})
	}
	return results, nil
}

func (r *PostgresRunner) Apply(ctx context.Context) ([]MigrationResult, error) {
	provider, err := r.requireProvider()
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	results, err := provider.Up(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres migrations: apply: %w", err)
	}
	converted := make([]MigrationResult, 0, len(results))
	for _, result := range results {
		if result == nil {
			continue
		}
		converted = append(converted, MigrationResult{
			Source:    convertSource(result.Source),
			Direction: result.Direction,
			Duration:  result.Duration,
			Empty:     result.Empty,
		})
	}
	return converted, nil
}

func (r *PostgresRunner) requireProvider() (*goose.Provider, error) {
	if r == nil || r.provider == nil {
		return nil, errors.New("postgres migrations: runner is required")
	}
	return r.provider, nil
}

func convertSource(source *goose.Source) MigrationSource {
	if source == nil {
		return MigrationSource{}
	}
	return MigrationSource{
		Type:    string(source.Type),
		Path:    source.Path,
		Version: source.Version,
	}
}

func wrapProviderError(operation string, err error) error {
	if errors.Is(err, goose.ErrNoMigrations) {
		return fmt.Errorf("%w: %v", ErrNoMigrationSources, err)
	}
	return fmt.Errorf("postgres migrations: %s: %w", operation, err)
}
