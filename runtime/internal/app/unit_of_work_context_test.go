package app

import (
	"context"
	"testing"

	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

func TestUnitOfWorkContextRoundTrip(t *testing.T) {
	unit := tx.NoopUnitOfWork{}
	ctx := ContextWithUnitOfWork(context.Background(), unit)

	got, ok := UnitOfWorkFromContext(ctx)
	if !ok {
		t.Fatal("UnitOfWorkFromContext() ok = false, want true")
	}
	if got.Context() == nil {
		t.Fatal("UnitOfWorkFromContext() returned unit with nil context")
	}
}

func TestContextWithUnitOfWorkIgnoresNilUnit(t *testing.T) {
	ctx := ContextWithUnitOfWork(context.Background(), nil)

	if _, ok := UnitOfWorkFromContext(ctx); ok {
		t.Fatal("UnitOfWorkFromContext() ok = true, want false for nil unit")
	}
}
