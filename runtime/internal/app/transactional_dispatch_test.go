package app

import (
	"context"
	"errors"
	"testing"

	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

func TestTransactionalDispatcherRunsCommandsInUnitOfWork(t *testing.T) {
	route := RouteKey{Kind: " command ", Module: "inventory", Name: "GrantItem"}
	request := RouteRequest{RequestID: "request-1", Route: route}
	runner := &recordingRunner{}
	inner := routeDispatcherFunc(func(ctx context.Context, req RouteRequest) (ApplicationResult, error) {
		if ctx.Value(recordingContextKey("uow")) != "active" {
			t.Fatalf("dispatcher context marker = %#v, want active", ctx.Value(recordingContextKey("uow")))
		}
		if _, ok := UnitOfWorkFromContext(ctx); !ok {
			t.Fatal("dispatcher context has no unit of work")
		}
		return resultForRequest(req), nil
	})
	dispatcher := TransactionalDispatcher{
		Dispatcher: inner,
		Runner:     runner,
	}

	result, err := dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}
	if result.RequestID != request.RequestID || result.Route != route {
		t.Fatalf("result = %#v, want request metadata", result)
	}
	if runner.calls != 1 {
		t.Fatalf("runner calls = %d, want 1", runner.calls)
	}
}

func TestTransactionalDispatcherDoesNotWrapQueries(t *testing.T) {
	route := RouteKey{Kind: MessageKindQuery, Module: "inventory", Name: "GetInventory"}
	request := RouteRequest{RequestID: "request-1", Route: route}
	runner := &recordingRunner{}
	var dispatched bool
	inner := routeDispatcherFunc(func(_ context.Context, req RouteRequest) (ApplicationResult, error) {
		dispatched = true
		return resultForRequest(req), nil
	})
	dispatcher := TransactionalDispatcher{
		Dispatcher: inner,
		Runner:     runner,
	}

	result, err := dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}
	if result.RequestID != request.RequestID || result.Route != route {
		t.Fatalf("result = %#v, want request metadata", result)
	}
	if !dispatched {
		t.Fatal("inner dispatcher was not called")
	}
	if runner.calls != 0 {
		t.Fatalf("runner calls = %d, want 0", runner.calls)
	}
}

func TestTransactionalDispatcherRequiresRunnerForCommands(t *testing.T) {
	dispatcher := TransactionalDispatcher{
		Dispatcher: routeDispatcherFunc(func(_ context.Context, req RouteRequest) (ApplicationResult, error) {
			return resultForRequest(req), nil
		}),
	}

	result, err := dispatcher.Dispatch(context.Background(), RouteRequest{
		RequestID: "request-1",
		Route:     RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"},
	})
	if err == nil {
		t.Fatal("Dispatch() error = nil, want missing runner error")
	}
	if result.RequestID != "request-1" {
		t.Fatalf("result RequestID = %q, want request-1", result.RequestID)
	}
}

func TestTransactionalDispatcherPropagatesCommandErrorWithResult(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}
	sentinel := errors.New("command failed")
	runner := &recordingRunner{}
	dispatcher := TransactionalDispatcher{
		Dispatcher: routeDispatcherFunc(func(_ context.Context, req RouteRequest) (ApplicationResult, error) {
			result := resultForRequest(req)
			result.Payload = "partial"
			return result, sentinel
		}),
		Runner: runner,
	}

	result, err := dispatcher.Dispatch(context.Background(), RouteRequest{
		RequestID: "request-1",
		Route:     route,
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("Dispatch() error = %v, want sentinel", err)
	}
	if result.RequestID != "request-1" || result.Route != route || result.Payload != "partial" {
		t.Fatalf("result = %#v, want handler result metadata and payload", result)
	}
	if runner.calls != 1 {
		t.Fatalf("runner calls = %d, want 1", runner.calls)
	}
}

func TestTransactionalDispatcherPropagatesRunnerErrorWithRequestResult(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}
	sentinel := errors.New("begin failed")
	dispatcher := TransactionalDispatcher{
		Dispatcher: routeDispatcherFunc(func(_ context.Context, req RouteRequest) (ApplicationResult, error) {
			return resultForRequest(req), nil
		}),
		Runner: failingRunner{err: sentinel},
	}

	result, err := dispatcher.Dispatch(context.Background(), RouteRequest{
		RequestID: "request-1",
		Route:     route,
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("Dispatch() error = %v, want sentinel", err)
	}
	if result.RequestID != "request-1" || result.Route != route {
		t.Fatalf("result = %#v, want request metadata", result)
	}
}

type routeDispatcherFunc func(context.Context, RouteRequest) (ApplicationResult, error)

func (f routeDispatcherFunc) Dispatch(ctx context.Context, request RouteRequest) (ApplicationResult, error) {
	return f(ctx, request)
}

type recordingContextKey string

type recordingRunner struct {
	calls int
}

func (r *recordingRunner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, tx.UnitOfWork) error) error {
	r.calls += 1
	runCtx := context.WithValue(ctx, recordingContextKey("uow"), "active")
	return fn(runCtx, tx.NoopUnitOfWork{})
}

type failingRunner struct {
	err error
}

func (r failingRunner) WithinUnitOfWork(context.Context, func(context.Context, tx.UnitOfWork) error) error {
	return r.err
}
