package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appconnection "github.com/iceiko/vibit/runtime/internal/app/connection"
	"github.com/iceiko/vibit/runtime/internal/platform/transport/ws"
)

func TestConnectionLifecycleRegistryObserverRegistersOpenAndClosedConnection(t *testing.T) {
	registry := appconnection.NewInMemoryRegistry(staticClock{now: fixedTime()})
	observer := connectionLifecycleRegistryObserver{registry: registry}

	openedAt := fixedTime().Add(time.Minute)
	if err := observer.ConnectionOpened(context.Background(), ws.ConnectionLifecycleEvent{
		ConnectionID:    " ws-1 ",
		ConnectionEpoch: 7,
		ObservedAt:      openedAt,
	}); err != nil {
		t.Fatalf("ConnectionOpened() error = %v, want nil", err)
	}
	record, ok := registry.FindConnectionByID(context.Background(), "ws-1", 7)
	if !ok || record.State != appconnection.StateOpenUnbound || !record.OpenedAt.Equal(openedAt.UTC()) {
		t.Fatalf("opened record = %#v, ok = %v, want open registry record", record, ok)
	}

	closedAt := fixedTime().Add(2 * time.Minute)
	if err := observer.ConnectionClosed(context.Background(), ws.ConnectionLifecycleEvent{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 7,
		ObservedAt:      closedAt,
	}); err != nil {
		t.Fatalf("ConnectionClosed() error = %v, want nil", err)
	}
	record, ok = registry.FindConnectionByID(context.Background(), "ws-1", 7)
	if !ok || record.State != appconnection.StateClosed || record.CloseReasonClass != "transport_closed" ||
		record.ClosedAt == nil || !record.ClosedAt.Equal(closedAt.UTC()) {
		t.Fatalf("closed record = %#v, ok = %v, want closed registry record", record, ok)
	}
}

func TestRegistryConnectionBinderBindsValidatedConnectionIntoRegistry(t *testing.T) {
	registry := appconnection.NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), appconnection.OpenConnection{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 1,
		OpenedAt:        fixedTime(),
	}); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	boundAt := fixedTime().Add(time.Minute)
	binder := registryConnectionBinder{
		registry: registry,
		binder: app.ConnectionBinder{
			Validator: staticRouteAccessTokenValidator{
				result: app.RouteAccessTokenValidationResult{
					Valid: true,
					Identity: app.ValidatedPlayerIdentity("player-1", app.Session{
						ConnectionID:    "ws-1",
						ConnectionEpoch: 1,
						SessionID:       "runtime-session-1",
						PlayerID:        "player-1",
					}),
				},
			},
			Clock: staticClock{now: boundAt},
		},
	}

	result, err := binder.BindConnection(context.Background(), app.ConnectionBindingRequest{
		AccessToken:     "redacted-token-text",
		Route:           app.BindConnectionRoute(),
		ConnectionID:    "ws-1",
		ConnectionEpoch: 1,
	})
	if err != nil {
		t.Fatalf("BindConnection() error = %v, want nil", err)
	}
	if !result.Bound {
		t.Fatalf("BindConnection() result = %#v, want bound", result)
	}
	record, ok := registry.FindConnectionByID(context.Background(), "ws-1", 1)
	if !ok ||
		record.State != appconnection.StateBound ||
		record.PlayerID != "player-1" ||
		record.RuntimeSessionID != "runtime-session-1" ||
		record.BoundAt == nil ||
		!record.BoundAt.Equal(boundAt.UTC()) {
		t.Fatalf("registry record = %#v, ok = %v, want bound identity linkage", record, ok)
	}

	presence := registry.PresenceForPlayer(context.Background(), "player-1")
	if presence.Status != appconnection.PresenceStatusOnline || presence.ConnectionCount != 1 {
		t.Fatalf("PresenceForPlayer() = %#v, want online after bind", presence)
	}
}

func TestRegistryConnectionBinderFailsClosedWhenRegistryBindFails(t *testing.T) {
	registry := appconnection.NewInMemoryRegistry(staticClock{now: fixedTime()})
	binder := registryConnectionBinder{
		registry: registry,
		binder: app.ConnectionBinder{
			Validator: staticRouteAccessTokenValidator{
				result: app.RouteAccessTokenValidationResult{
					Valid: true,
					Identity: app.ValidatedPlayerIdentity("player-1", app.Session{
						ConnectionID:    "missing",
						ConnectionEpoch: 1,
						SessionID:       "runtime-session-1",
						PlayerID:        "player-1",
					}),
				},
			},
			Clock: staticClock{now: fixedTime()},
		},
	}

	result, err := binder.BindConnection(context.Background(), app.ConnectionBindingRequest{
		AccessToken:     "redacted-token-text",
		Route:           app.BindConnectionRoute(),
		ConnectionID:    "missing",
		ConnectionEpoch: 1,
	})

	if err == nil {
		t.Fatal("BindConnection() error = nil, want fail-closed registry unavailable error")
	}
	var appErr *app.ApplicationError
	if !errors.As(err, &appErr) || appErr.Code != app.ErrorCodeConnectionBindingUnavailable {
		t.Fatalf("BindConnection() error = %v, want CONNECTION_BINDING_UNAVAILABLE", err)
	}
	if result.Bound || result.BindingStatus != app.ConnectionBindingStatusRejected {
		t.Fatalf("BindConnection() result = %#v, want rejected", result)
	}
}

func fixedTime() time.Time {
	return time.Date(2026, 5, 20, 1, 2, 3, 0, time.UTC)
}

type staticClock struct {
	now time.Time
}

func (c staticClock) Now() time.Time {
	return c.now
}
