package app

import (
	"context"
	"errors"
	"testing"
)

func TestMetadataOnlySessionValidatorNormalizesIdentity(t *testing.T) {
	validator := MetadataOnlySessionValidator{}
	result, err := validator.ValidateSession(context.Background(), RouteRequest{
		Session: Session{
			ConnectionID:    " connection-1 ",
			SessionID:       " session-1 ",
			PlayerID:        " player-1 ",
			ConnectionEpoch: 7,
		},
	})
	if err != nil {
		t.Fatalf("ValidateSession() error = %v, want nil", err)
	}
	if !result.Valid {
		t.Fatalf("ValidateSession() Valid = false, want true")
	}
	if result.Identity.Status != IdentityValidationMetadataOnly {
		t.Fatalf("Identity.Status = %q, want %q", result.Identity.Status, IdentityValidationMetadataOnly)
	}
	if result.Identity.PlayerID != "player-1" || result.Identity.SessionID != "session-1" || result.Identity.ConnectionID != "connection-1" {
		t.Fatalf("Identity = %#v, want normalized metadata-only identity", result.Identity)
	}
	if result.Identity.PlayerIDValidated || result.Identity.SessionValidated {
		t.Fatalf("Identity validation flags = %#v, want metadata-only identity", result.Identity)
	}
}

func TestSessionValidatingDispatcherPassesMetadataOnlyIdentity(t *testing.T) {
	route := RouteKey{Kind: MessageKindQuery, Module: "inventory", Name: "GetInventory"}
	var received RouteRequest
	dispatcher := SessionValidatingDispatcher{
		Dispatcher: routeDispatcherFunc(func(_ context.Context, request RouteRequest) (ApplicationResult, error) {
			received = request
			return resultForRequest(request), nil
		}),
	}

	request := RouteRequest{
		RequestID: "request-1",
		Route:     route,
		Session:   Session{SessionID: "session-1", PlayerID: "player-1"},
	}
	result, err := dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}
	if received.Identity.Status != IdentityValidationMetadataOnly {
		t.Fatalf("received Identity.Status = %q, want %q", received.Identity.Status, IdentityValidationMetadataOnly)
	}
	if received.Identity.PlayerID != "player-1" || received.Identity.SessionValidated || received.Identity.PlayerIDValidated {
		t.Fatalf("received Identity = %#v, want metadata-only player-1", received.Identity)
	}
	if result.Identity != received.Identity {
		t.Fatalf("result Identity = %#v, want %#v", result.Identity, received.Identity)
	}
}

func TestSessionValidatingDispatcherUsesInjectedValidatedIdentity(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}
	validatedIdentity := ValidatedPlayerIdentity("player-2", Session{
		ConnectionID: "connection-1",
		SessionID:    "session-1",
		PlayerID:     "player-1",
	})
	var received RouteRequest
	dispatcher := SessionValidatingDispatcher{
		Validator: SessionValidatorFunc(func(context.Context, RouteRequest) (SessionValidationResult, error) {
			return SessionValidationResult{Identity: validatedIdentity, Valid: true, Reason: "test validator"}, nil
		}),
		Dispatcher: routeDispatcherFunc(func(_ context.Context, request RouteRequest) (ApplicationResult, error) {
			received = request
			return resultForRequest(request), nil
		}),
	}

	result, err := dispatcher.Dispatch(context.Background(), RouteRequest{
		RequestID: "request-1",
		Route:     route,
		Session:   Session{ConnectionID: "connection-1", SessionID: "session-1", PlayerID: "player-1"},
	})
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}
	if received.Identity != validatedIdentity {
		t.Fatalf("received Identity = %#v, want injected validated identity %#v", received.Identity, validatedIdentity)
	}
	if result.Identity != validatedIdentity {
		t.Fatalf("result Identity = %#v, want injected validated identity %#v", result.Identity, validatedIdentity)
	}
}

func TestSessionValidatingDispatcherStopsInvalidSession(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "inventory", Name: "GrantItem"}
	var called bool
	dispatcher := SessionValidatingDispatcher{
		Validator: SessionValidatorFunc(func(context.Context, RouteRequest) (SessionValidationResult, error) {
			return SessionValidationResult{Valid: false, Reason: "invalid test session"}, nil
		}),
		Dispatcher: routeDispatcherFunc(func(context.Context, RouteRequest) (ApplicationResult, error) {
			called = true
			return ApplicationResult{}, nil
		}),
	}

	result, err := dispatcher.Dispatch(context.Background(), RouteRequest{RequestID: "request-1", Route: route})
	assertApplicationErrorCode(t, err, ErrorCodeSessionInvalid)
	if called {
		t.Fatal("inner dispatcher was called, want validation to stop dispatch")
	}
	if result.Error == nil || result.Error.Code != ErrorCodeSessionInvalid {
		t.Fatalf("result Error = %#v, want %s", result.Error, ErrorCodeSessionInvalid)
	}
}

func TestSessionValidatingDispatcherRequiresDispatcher(t *testing.T) {
	dispatcher := SessionValidatingDispatcher{}

	_, err := dispatcher.Dispatch(context.Background(), RouteRequest{RequestID: "request-1"})
	if err == nil {
		t.Fatal("Dispatch() error = nil, want missing dispatcher error")
	}
}

func TestSessionValidatingDispatcherPropagatesValidatorError(t *testing.T) {
	sentinel := errors.New("validator failed")
	dispatcher := SessionValidatingDispatcher{
		Validator: SessionValidatorFunc(func(context.Context, RouteRequest) (SessionValidationResult, error) {
			return SessionValidationResult{}, sentinel
		}),
		Dispatcher: routeDispatcherFunc(func(context.Context, RouteRequest) (ApplicationResult, error) {
			return ApplicationResult{}, nil
		}),
	}

	_, err := dispatcher.Dispatch(context.Background(), RouteRequest{RequestID: "request-1"})
	if !errors.Is(err, sentinel) {
		t.Fatalf("Dispatch() error = %v, want sentinel", err)
	}
}
