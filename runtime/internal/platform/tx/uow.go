package tx

import (
	"context"
	"errors"
)

type UnitOfWork interface {
	Context() context.Context
}

type Runner interface {
	WithinUnitOfWork(context.Context, func(context.Context, UnitOfWork) error) error
}

type NoopUnitOfWork struct {
	ctx context.Context
}

func (u NoopUnitOfWork) Context() context.Context {
	if u.ctx == nil {
		return context.Background()
	}
	return u.ctx
}

type NoopRunner struct{}

func (NoopRunner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, UnitOfWork) error) error {
	if fn == nil {
		return errors.New("tx: unit-of-work function is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	unit := NoopUnitOfWork{ctx: ctx}
	return fn(unit.Context(), unit)
}
