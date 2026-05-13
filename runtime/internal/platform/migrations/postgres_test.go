package migrations

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
)

const fakeDriverName = "vibit_postgres_migrations_fake"

var (
	registerFakeDriver sync.Once
	errFakeDatabase    = errors.New("fake database: live PostgreSQL is not available")
)

func TestNewPostgresRunnerListsSQLSources(t *testing.T) {
	db := openFakeDB(t)
	defer db.Close()

	runner, err := NewPostgresRunner(db, validMigrationFS())
	if err != nil {
		t.Fatalf("NewPostgresRunner() error = %v, want nil", err)
	}

	sources, err := runner.Sources()
	if err != nil {
		t.Fatalf("Sources() error = %v, want nil", err)
	}
	if len(sources) != 1 {
		t.Fatalf("sources len = %d, want 1", len(sources))
	}
	if sources[0].Type != "sql" {
		t.Fatalf("source type = %q, want sql", sources[0].Type)
	}
	if sources[0].Path != "000001_create_inventory_state.sql" {
		t.Fatalf("source path = %q, want first migration path", sources[0].Path)
	}
	if sources[0].Version != 1 {
		t.Fatalf("source version = %d, want 1", sources[0].Version)
	}
}

func TestNewPostgresRunnerRejectsMissingOptions(t *testing.T) {
	db := openFakeDB(t)
	defer db.Close()

	_, err := NewPostgresRunner(nil, validMigrationFS())
	if err == nil {
		t.Fatal("NewPostgresRunner(nil, fs) error = nil, want database handle error")
	}
	if !strings.Contains(err.Error(), "database handle") {
		t.Fatalf("NewPostgresRunner(nil, fs) error = %v, want database handle context", err)
	}

	_, err = NewPostgresRunner(db, nil)
	if err == nil {
		t.Fatal("NewPostgresRunner(db, nil) error = nil, want source filesystem error")
	}
	if !strings.Contains(err.Error(), "source filesystem") {
		t.Fatalf("NewPostgresRunner(db, nil) error = %v, want source filesystem context", err)
	}
}

func TestNewPostgresRunnerRejectsEmptySources(t *testing.T) {
	db := openFakeDB(t)
	defer db.Close()

	_, err := NewPostgresRunner(db, fstest.MapFS{})
	if err == nil {
		t.Fatal("NewPostgresRunner(empty fs) error = nil, want no migration sources error")
	}
	if !errors.Is(err, ErrNoMigrationSources) {
		t.Fatalf("NewPostgresRunner(empty fs) error = %v, want ErrNoMigrationSources", err)
	}
}

func TestNewPostgresRunnerFromDirRequiresExplicitDirectory(t *testing.T) {
	db := openFakeDB(t)
	defer db.Close()

	_, err := NewPostgresRunnerFromDir(db, " ")
	if err == nil {
		t.Fatal("NewPostgresRunnerFromDir() error = nil, want source directory error")
	}
	if !strings.Contains(err.Error(), "source directory") {
		t.Fatalf("NewPostgresRunnerFromDir() error = %v, want source directory context", err)
	}
}

func TestPostgresRunnerRequiresInitializedRunner(t *testing.T) {
	var runner *PostgresRunner

	_, err := runner.Sources()
	if err == nil {
		t.Fatal("Sources() error = nil, want runner error")
	}

	_, err = runner.Status(context.Background())
	if err == nil {
		t.Fatal("Status() error = nil, want runner error")
	}

	_, err = runner.Apply(context.Background())
	if err == nil {
		t.Fatal("Apply() error = nil, want runner error")
	}
}

func TestPostgresRunnerStatusAndApplyRequireLiveDatabase(t *testing.T) {
	db := openFakeDB(t)
	defer db.Close()

	runner, err := NewPostgresRunner(db, validMigrationFS())
	if err != nil {
		t.Fatalf("NewPostgresRunner() error = %v, want nil", err)
	}

	_, err = runner.Status(context.Background())
	if !errors.Is(err, errFakeDatabase) {
		t.Fatalf("Status() error = %v, want fake database error", err)
	}
	if !strings.Contains(err.Error(), "status") {
		t.Fatalf("Status() error = %v, want status context", err)
	}

	_, err = runner.Apply(context.Background())
	if !errors.Is(err, errFakeDatabase) {
		t.Fatalf("Apply() error = %v, want fake database error", err)
	}
	if !strings.Contains(err.Error(), "apply") {
		t.Fatalf("Apply() error = %v, want apply context", err)
	}
}

func validMigrationFS() fstest.MapFS {
	return fstest.MapFS{
		"000001_create_inventory_state.sql": {
			Data: []byte("-- +goose Up\nSELECT 1;\n-- +goose Down\nSELECT 1;\n"),
		},
	}
}

func openFakeDB(t *testing.T) *sql.DB {
	t.Helper()

	registerFakeDriver.Do(func() {
		sql.Register(fakeDriverName, fakeDriver{})
	})
	db, err := sql.Open(fakeDriverName, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v, want nil", err)
	}
	return db
}

type fakeDriver struct{}

func (fakeDriver) Open(string) (driver.Conn, error) {
	return fakeConn{}, nil
}

type fakeConn struct{}

func (fakeConn) Prepare(string) (driver.Stmt, error) {
	return nil, errFakeDatabase
}

func (fakeConn) Close() error {
	return nil
}

func (fakeConn) Begin() (driver.Tx, error) {
	return nil, errFakeDatabase
}

func (fakeConn) Ping(context.Context) error {
	return errFakeDatabase
}

func (fakeConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return nil, errFakeDatabase
}

func (fakeConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return nil, errFakeDatabase
}

type fakeRows struct{}

func (fakeRows) Columns() []string {
	return nil
}

func (fakeRows) Close() error {
	return nil
}

func (fakeRows) Next([]driver.Value) error {
	return io.EOF
}
