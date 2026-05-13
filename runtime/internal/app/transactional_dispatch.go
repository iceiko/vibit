package app

import (
	"context"
	"errors"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

type RouteDispatcher interface {
	Dispatch(context.Context, RouteRequest) (ApplicationResult, error)
}

type TransactionalDispatcher struct {
	Dispatcher RouteDispatcher
	Runner     tx.Runner
}

func (d TransactionalDispatcher) Dispatch(ctx context.Context, request RouteRequest) (ApplicationResult, error) {
	dispatcher := d.Dispatcher
	if dispatcher == nil {
		return resultForRequest(request), errors.New("app: transactional dispatcher requires dispatcher")
	}

	if MessageKind(strings.TrimSpace(string(request.Route.Kind))) != MessageKindCommand {
		return dispatcher.Dispatch(ctx, request)
	}

	runner := d.Runner
	if runner == nil {
		return resultForRequest(request), errors.New("app: transactional dispatcher requires unit-of-work runner for command routes")
	}

	var result ApplicationResult
	err := runner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		var dispatchErr error
		runCtx = ContextWithUnitOfWork(runCtx, unit)
		result, dispatchErr = dispatcher.Dispatch(runCtx, request)
		return dispatchErr
	})
	if result.RequestID == "" && result.Route == (RouteKey{}) {
		result = resultForRequest(request)
	}
	return result, err
}
