package app

import (
	"context"
	"errors"
	"testing"
)

func TestDispatcherDispatchesRegisteredCommand(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{
		Kind:   MessageKindCommand,
		Module: " inventory ",
		Name:   " GrantItem ",
	}
	request := RouteRequest{
		RequestID: "request-1",
		Route:     route,
		Target: Target{
			Scope: TargetScopePlayer,
			ID:    "player-1",
		},
		Session: Session{
			ConnectionID:    "connection-1",
			SessionID:       "session-1",
			PlayerID:        "player-1",
			ConnectionEpoch: 7,
		},
		PayloadType: "example.Payload",
		Payload:     "payload",
	}

	var received RouteRequest
	if err := dispatcher.Register(route, HandlerFunc(func(_ context.Context, req RouteRequest) (ApplicationResult, error) {
		received = req
		return ApplicationResult{
			PayloadType: "example.Response",
			Payload:     "ok",
		}, nil
	})); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	result, err := dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}

	wantRoute := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}
	if received.Route != wantRoute {
		t.Fatalf("handler received route = %#v, want %#v", received.Route, wantRoute)
	}
	if received.Payload != "payload" {
		t.Fatalf("handler received payload = %#v, want payload", received.Payload)
	}
	if received.Identity.Status != IdentityValidationMetadataOnly {
		t.Fatalf("handler received identity status = %q, want %q", received.Identity.Status, IdentityValidationMetadataOnly)
	}
	if received.Identity.PlayerID != "player-1" || received.Identity.SessionID != "session-1" || received.Identity.ConnectionID != "connection-1" {
		t.Fatalf("handler received identity = %#v, want normalized metadata-only session identity", received.Identity)
	}
	if received.Identity.PlayerIDValidated || received.Identity.SessionValidated {
		t.Fatalf("handler received identity validation flags = %#v, want metadata-only identity", received.Identity)
	}
	if result.RequestID != request.RequestID {
		t.Fatalf("result RequestID = %q, want %q", result.RequestID, request.RequestID)
	}
	if result.Route != wantRoute {
		t.Fatalf("result Route = %#v, want %#v", result.Route, wantRoute)
	}
	if result.Target != request.Target {
		t.Fatalf("result Target = %#v, want %#v", result.Target, request.Target)
	}
	if result.Session != request.Session {
		t.Fatalf("result Session = %#v, want %#v", result.Session, request.Session)
	}
	if result.Identity != received.Identity {
		t.Fatalf("result Identity = %#v, want %#v", result.Identity, received.Identity)
	}
	if result.Payload != "ok" {
		t.Fatalf("result Payload = %#v, want ok", result.Payload)
	}
}

func TestDispatcherDispatchesRegisteredQuery(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{Kind: MessageKindQuery, Module: "inventory", Name: "GetInventory"}

	if err := dispatcher.Register(route, HandlerFunc(func(_ context.Context, _ RouteRequest) (ApplicationResult, error) {
		return ApplicationResult{Payload: []string{"item-1"}}, nil
	})); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	result, err := dispatcher.Dispatch(context.Background(), RouteRequest{Route: route})
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}
	if len(result.Payload.([]string)) != 1 {
		t.Fatalf("result Payload = %#v, want one item", result.Payload)
	}
}

func TestDispatcherRejectsDuplicateRoute(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}
	handler := HandlerFunc(func(_ context.Context, _ RouteRequest) (ApplicationResult, error) {
		return ApplicationResult{}, nil
	})

	if err := dispatcher.Register(route, handler); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	err := dispatcher.Register(route, handler)
	assertApplicationErrorCode(t, err, ErrorCodeRouteAlreadyRegistered)
}

func TestDispatcherRejectsNilHandler(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}

	err := dispatcher.Register(route, nil)
	assertApplicationErrorCode(t, err, ErrorCodeNilHandler)
}

func TestDispatcherRejectsNilHandlerFunc(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}

	var handler HandlerFunc
	err := dispatcher.Register(route, handler)
	assertApplicationErrorCode(t, err, ErrorCodeNilHandler)
}

func TestDispatcherRejectsInvalidRoute(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory"}

	err := dispatcher.Register(route, HandlerFunc(func(_ context.Context, _ RouteRequest) (ApplicationResult, error) {
		return ApplicationResult{}, nil
	}))
	assertApplicationErrorCode(t, err, ErrorCodeInvalidRoute)
}

func TestDispatcherRejectsUnsupportedMessageKind(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{Kind: MessageKindEvent, Module: "inventory", Name: "ItemGranted"}

	err := dispatcher.Register(route, HandlerFunc(func(_ context.Context, _ RouteRequest) (ApplicationResult, error) {
		return ApplicationResult{}, nil
	}))
	assertApplicationErrorCode(t, err, ErrorCodeUnsupportedMessageKind)

	result, err := dispatcher.Dispatch(context.Background(), RouteRequest{Route: route})
	assertApplicationErrorCode(t, err, ErrorCodeUnsupportedMessageKind)
	if result.Error == nil || result.Error.Code != ErrorCodeUnsupportedMessageKind {
		t.Fatalf("result Error = %#v, want %s", result.Error, ErrorCodeUnsupportedMessageKind)
	}
}

func TestDispatcherRejectsUnknownRoute(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}

	result, err := dispatcher.Dispatch(context.Background(), RouteRequest{
		RequestID: "request-1",
		Route:     route,
	})
	assertApplicationErrorCode(t, err, ErrorCodeRouteNotFound)
	if result.RequestID != "request-1" {
		t.Fatalf("result RequestID = %q, want request-1", result.RequestID)
	}
	if result.Error == nil || result.Error.Code != ErrorCodeRouteNotFound {
		t.Fatalf("result Error = %#v, want %s", result.Error, ErrorCodeRouteNotFound)
	}
}

func TestDispatcherPropagatesHandlerErrorWithMetadata(t *testing.T) {
	dispatcher := NewDispatcher()
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}
	sentinel := errors.New("handler failed")

	if err := dispatcher.Register(route, HandlerFunc(func(_ context.Context, _ RouteRequest) (ApplicationResult, error) {
		return ApplicationResult{}, sentinel
	})); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	request := RouteRequest{
		RequestID: "request-1",
		Route:     route,
		Target:    Target{Scope: TargetScopePlayer, ID: "player-1"},
		Session:   Session{SessionID: "session-1"},
	}
	result, err := dispatcher.Dispatch(context.Background(), request)
	if !errors.Is(err, sentinel) {
		t.Fatalf("Dispatch() error = %v, want sentinel", err)
	}
	if result.RequestID != request.RequestID {
		t.Fatalf("result RequestID = %q, want %q", result.RequestID, request.RequestID)
	}
	if result.Route != route {
		t.Fatalf("result Route = %#v, want %#v", result.Route, route)
	}
	if result.Target != request.Target {
		t.Fatalf("result Target = %#v, want %#v", result.Target, request.Target)
	}
	if result.Session != request.Session {
		t.Fatalf("result Session = %#v, want %#v", result.Session, request.Session)
	}
}

func assertApplicationErrorCode(t *testing.T, err error, code ErrorCode) {
	t.Helper()

	var appErr *ApplicationError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *ApplicationError", err)
	}
	if appErr.Code != code {
		t.Fatalf("ApplicationError.Code = %s, want %s", appErr.Code, code)
	}
}
