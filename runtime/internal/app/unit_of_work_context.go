package app

import (
	"context"

	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

type unitOfWorkContextKey struct{}

func ContextWithUnitOfWork(ctx context.Context, unit tx.UnitOfWork) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if unit == nil {
		return ctx
	}
	return context.WithValue(ctx, unitOfWorkContextKey{}, unit)
}

func UnitOfWorkFromContext(ctx context.Context) (tx.UnitOfWork, bool) {
	if ctx == nil {
		return nil, false
	}
	unit, ok := ctx.Value(unitOfWorkContextKey{}).(tx.UnitOfWork)
	if !ok || unit == nil {
		return nil, false
	}
	return unit, true
}
