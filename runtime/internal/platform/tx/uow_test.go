package tx

import (
	"context"
	"errors"
	"testing"
)

func TestNoopRunnerRunsFunctionWithUnitOfWorkContext(t *testing.T) {
	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "request-1")
	runner := NoopRunner{}

	var called bool
	err := runner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit UnitOfWork) error {
		called = true
		if runCtx.Value(contextKey("request")) != "request-1" {
			t.Fatalf("run context value = %#v, want request-1", runCtx.Value(contextKey("request")))
		}
		if unit.Context() != runCtx {
			t.Fatalf("unit.Context() = %#v, want run context", unit.Context())
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithinUnitOfWork() error = %v, want nil", err)
	}
	if !called {
		t.Fatal("WithinUnitOfWork() did not call function")
	}
}

func TestNoopRunnerPropagatesFunctionError(t *testing.T) {
	runner := NoopRunner{}
	sentinel := errors.New("failed")

	err := runner.WithinUnitOfWork(context.Background(), func(context.Context, UnitOfWork) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("WithinUnitOfWork() error = %v, want sentinel", err)
	}
}

func TestNoopRunnerRejectsNilFunction(t *testing.T) {
	runner := NoopRunner{}

	err := runner.WithinUnitOfWork(context.Background(), nil)
	if err == nil {
		t.Fatal("WithinUnitOfWork() error = nil, want error")
	}
}
