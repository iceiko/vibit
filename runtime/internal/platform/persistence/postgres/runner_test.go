package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/iceiko/vibit/runtime/internal/platform/tx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestRunnerCommitsSuccessfulUnitOfWork(t *testing.T) {
	beginner := &fakeBeginner{tx: &fakeTx{}}
	runner := NewRunner(beginner)
	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "request-1")

	var called bool
	err := runner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		called = true
		if runCtx.Value(contextKey("request")) != "request-1" {
			t.Fatalf("run context value = %#v, want request-1", runCtx.Value(contextKey("request")))
		}
		if unit.Context() != runCtx {
			t.Fatalf("unit.Context() = %#v, want run context", unit.Context())
		}
		postgresUnit, ok := unit.(UnitOfWork)
		if !ok {
			t.Fatalf("unit type = %T, want postgres.UnitOfWork", unit)
		}
		repository, err := postgresUnit.NewInventoryRepository()
		if err != nil {
			t.Fatalf("NewInventoryRepository() error = %v, want nil", err)
		}
		if repository == nil {
			t.Fatal("NewInventoryRepository() = nil, want repository")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithinUnitOfWork() error = %v, want nil", err)
	}
	if !called {
		t.Fatal("WithinUnitOfWork() did not call function")
	}
	if beginner.calls != 1 {
		t.Fatalf("begin calls = %d, want 1", beginner.calls)
	}
	if beginner.tx.commits != 1 {
		t.Fatalf("commits = %d, want 1", beginner.tx.commits)
	}
	if beginner.tx.rollbacks != 0 {
		t.Fatalf("rollbacks = %d, want 0", beginner.tx.rollbacks)
	}
}

func TestRunnerRollsBackFunctionError(t *testing.T) {
	sentinel := errors.New("handler failed")
	beginner := &fakeBeginner{tx: &fakeTx{}}
	runner := NewRunner(beginner)

	err := runner.WithinUnitOfWork(context.Background(), func(context.Context, tx.UnitOfWork) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("WithinUnitOfWork() error = %v, want sentinel", err)
	}
	if beginner.tx.commits != 0 {
		t.Fatalf("commits = %d, want 0", beginner.tx.commits)
	}
	if beginner.tx.rollbacks != 1 {
		t.Fatalf("rollbacks = %d, want 1", beginner.tx.rollbacks)
	}
}

func TestRunnerRollsBackCommitError(t *testing.T) {
	sentinel := errors.New("commit failed")
	beginner := &fakeBeginner{tx: &fakeTx{commitErr: sentinel}}
	runner := NewRunner(beginner)

	err := runner.WithinUnitOfWork(context.Background(), func(context.Context, tx.UnitOfWork) error {
		return nil
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("WithinUnitOfWork() error = %v, want sentinel", err)
	}
	if !strings.Contains(err.Error(), "commit") {
		t.Fatalf("WithinUnitOfWork() error = %v, want commit context", err)
	}
	if beginner.tx.commits != 1 {
		t.Fatalf("commits = %d, want 1", beginner.tx.commits)
	}
	if beginner.tx.rollbacks != 1 {
		t.Fatalf("rollbacks = %d, want 1", beginner.tx.rollbacks)
	}
}

func TestRunnerRejectsMissingDependencies(t *testing.T) {
	runner := NewRunner(nil)

	err := runner.WithinUnitOfWork(context.Background(), func(context.Context, tx.UnitOfWork) error {
		return nil
	})
	if err == nil {
		t.Fatal("WithinUnitOfWork() error = nil, want missing beginner error")
	}

	runner = NewRunner(&fakeBeginner{tx: &fakeTx{}})
	err = runner.WithinUnitOfWork(context.Background(), nil)
	if err == nil {
		t.Fatal("WithinUnitOfWork(nil) error = nil, want missing function error")
	}
}

func TestNewPoolRunnerRejectsNilPool(t *testing.T) {
	runner := NewPoolRunner(nil)

	err := runner.WithinUnitOfWork(context.Background(), func(context.Context, tx.UnitOfWork) error {
		return nil
	})
	if err == nil {
		t.Fatal("WithinUnitOfWork() error = nil, want missing beginner error")
	}
}

func TestRunnerPropagatesBeginError(t *testing.T) {
	sentinel := errors.New("begin failed")
	runner := NewRunner(&fakeBeginner{err: sentinel})

	err := runner.WithinUnitOfWork(context.Background(), func(context.Context, tx.UnitOfWork) error {
		return nil
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("WithinUnitOfWork() error = %v, want sentinel", err)
	}
}

func TestUnitOfWorkRejectsMissingExecutor(t *testing.T) {
	unit := UnitOfWork{}

	_, err := unit.Executor()
	if err == nil {
		t.Fatal("Executor() error = nil, want missing executor error")
	}

	_, err = unit.NewInventoryRepository()
	if err == nil {
		t.Fatal("NewInventoryRepository() error = nil, want missing executor error")
	}
}

type fakeBeginner struct {
	tx    *fakeTx
	err   error
	calls int
}

func (b *fakeBeginner) Begin(context.Context) (pgx.Tx, error) {
	b.calls += 1
	if b.err != nil {
		return nil, b.err
	}
	return b.tx, nil
}

type fakeTx struct {
	commits   int
	rollbacks int
	commitErr error
}

func (t *fakeTx) Begin(context.Context) (pgx.Tx, error) {
	return nil, errors.New("fake tx: nested begin not supported")
}

func (t *fakeTx) Commit(context.Context) error {
	t.commits += 1
	return t.commitErr
}

func (t *fakeTx) Rollback(context.Context) error {
	t.rollbacks += 1
	return nil
}

func (t *fakeTx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	return 0, errors.New("fake tx: CopyFrom not implemented")
}

func (t *fakeTx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults {
	return nil
}

func (t *fakeTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (t *fakeTx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	return nil, errors.New("fake tx: Prepare not implemented")
}

func (t *fakeTx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, errors.New("fake tx: Exec not implemented")
}

func (t *fakeTx) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errors.New("fake tx: Query not implemented")
}

func (t *fakeTx) QueryRow(context.Context, string, ...any) pgx.Row {
	return fakeRow{err: errors.New("fake tx: QueryRow not implemented")}
}

func (t *fakeTx) Conn() *pgx.Conn {
	return nil
}
