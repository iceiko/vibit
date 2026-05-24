package presence

import (
	"context"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/app/connection"
)

func TestGetPlayerPresenceReturnsSelfPresenceSnapshot(t *testing.T) {
	registry := connection.NewInMemoryRegistry(fixedClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), connection.OpenConnection{
		ConnectionID:    "connection-1",
		ConnectionEpoch: 7,
		OpenedAt:        fixedTime().Add(-time.Minute),
	}); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	if _, err := registry.BindConnectionIdentity(context.Background(), connection.BindIdentity{
		ConnectionID:        "connection-1",
		ConnectionEpoch:     7,
		ActorKind:           connection.ActorKindPlayer,
		PlayerID:            "player-1",
		RuntimeSessionID:    "session-1",
		AccessTokenRecordID: "token-record-1",
		ValidatedAt:         fixedTime(),
	}); err != nil {
		t.Fatalf("BindConnectionIdentity() error = %v, want nil", err)
	}

	result, err := (Handlers{Registry: registry}).GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{}, validatedIdentity("player-1"))
	if err != nil {
		t.Fatalf("GetPlayerPresence() error = %v, want nil", err)
	}

	if result.PlayerID != "player-1" ||
		result.Status != PresenceStatusOnline ||
		result.ConnectionCount != 1 ||
		len(result.ActiveConnections) != 1 ||
		len(result.RuntimeSessionIDs) != 1 ||
		result.RuntimeSessionIDs[0] != "session-1" {
		t.Fatalf("GetPlayerPresence() = %#v, want online self snapshot", result)
	}
	connection := result.ActiveConnections[0]
	if connection.ConnectionID != "connection-1" ||
		connection.ConnectionEpoch != 7 ||
		connection.RuntimeSessionID != "session-1" ||
		connection.LastSeenAt == nil ||
		!connection.LastSeenAt.Equal(fixedTime()) ||
		connection.BoundAt == nil ||
		!connection.BoundAt.Equal(fixedTime()) {
		t.Fatalf("active connection = %#v, want bounded connection metadata", connection)
	}
}

func TestGetPlayerPresenceAllowsExplicitSelfPlayerID(t *testing.T) {
	registry := connection.NewInMemoryRegistry(fixedClock{now: fixedTime()})
	result, err := (Handlers{Registry: registry}).GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{
		PlayerID: " player-1 ",
	}, validatedIdentity("player-1"))
	if err != nil {
		t.Fatalf("GetPlayerPresence() error = %v, want nil", err)
	}
	if result.PlayerID != "player-1" || result.Status != PresenceStatusOffline || result.ConnectionCount != 0 {
		t.Fatalf("GetPlayerPresence() = %#v, want offline self snapshot", result)
	}
}

func TestGetPlayerPresenceReportsOfflineAfterCloseAndInvalidation(t *testing.T) {
	tests := []struct {
		name     string
		terminal func(*testing.T, *connection.InMemoryRegistry)
	}{
		{
			name: "transport close",
			terminal: func(t *testing.T, registry *connection.InMemoryRegistry) {
				t.Helper()
				if _, err := registry.MarkConnectionClosed(context.Background(), connection.MarkClosed{
					ConnectionID:     "connection-1",
					ConnectionEpoch:  7,
					ClosedAt:         fixedTime().Add(time.Minute),
					CloseReasonClass: "transport_closed",
				}); err != nil {
					t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
				}
			},
		},
		{
			name: "policy invalidation",
			terminal: func(t *testing.T, registry *connection.InMemoryRegistry) {
				t.Helper()
				if _, err := registry.MarkConnectionInvalidated(context.Background(), connection.Invalidation{
					ConnectionID:      "connection-1",
					ConnectionEpoch:   7,
					InvalidatedAt:     fixedTime().Add(time.Minute),
					InvalidationClass: "token_revoked",
				}); err != nil {
					t.Fatalf("MarkConnectionInvalidated() error = %v, want nil", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := connection.NewInMemoryRegistry(fixedClock{now: fixedTime()})
			if _, err := registry.RegisterOpenConnection(context.Background(), connection.OpenConnection{
				ConnectionID:    "connection-1",
				ConnectionEpoch: 7,
				OpenedAt:        fixedTime().Add(-time.Minute),
			}); err != nil {
				t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
			}
			if _, err := registry.BindConnectionIdentity(context.Background(), connection.BindIdentity{
				ConnectionID:     "connection-1",
				ConnectionEpoch:  7,
				ActorKind:        connection.ActorKindPlayer,
				PlayerID:         "player-1",
				RuntimeSessionID: "session-1",
				ValidatedAt:      fixedTime(),
			}); err != nil {
				t.Fatalf("BindConnectionIdentity() error = %v, want nil", err)
			}
			handlers := Handlers{Registry: registry}

			online, err := handlers.GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{}, validatedIdentity("player-1"))
			if err != nil {
				t.Fatalf("GetPlayerPresence(online) error = %v, want nil", err)
			}
			if online.Status != PresenceStatusOnline || online.ConnectionCount != 1 {
				t.Fatalf("GetPlayerPresence(online) = %#v, want online with one active connection", online)
			}

			tt.terminal(t, registry)

			offline, err := handlers.GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{}, validatedIdentity("player-1"))
			if err != nil {
				t.Fatalf("GetPlayerPresence(offline) error = %v, want nil", err)
			}
			if offline.Status != PresenceStatusOffline ||
				offline.ConnectionCount != 0 ||
				len(offline.ActiveConnections) != 0 ||
				len(offline.RuntimeSessionIDs) != 0 ||
				offline.LastSeenAt != nil {
				t.Fatalf("GetPlayerPresence(offline) = %#v, want offline without active connection metadata", offline)
			}
		})
	}
}

func TestGetPlayerPresenceRejectsCrossPlayerQuery(t *testing.T) {
	_, err := (Handlers{Registry: connection.NewInMemoryRegistry(nil)}).GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{
		PlayerID: "player-2",
	}, validatedIdentity("player-1"))
	if err == nil {
		t.Fatal("GetPlayerPresence() error = nil, want forbidden")
	}
	assertApplicationError(t, err, ErrorCodePresenceQueryForbidden)
}

func TestGetPlayerPresenceRejectsMetadataOnlyIdentity(t *testing.T) {
	_, err := (Handlers{Registry: connection.NewInMemoryRegistry(nil)}).GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{}, app.MetadataOnlyIdentityFromSession(app.Session{PlayerID: "player-1"}))
	if err == nil {
		t.Fatal("GetPlayerPresence() error = nil, want forbidden")
	}
	assertApplicationError(t, err, ErrorCodePresenceQueryForbidden)
}

func TestGetPlayerPresenceFailsClosedWhenRegistryUnavailable(t *testing.T) {
	_, err := (Handlers{}).GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{}, validatedIdentity("player-1"))
	if err == nil {
		t.Fatal("GetPlayerPresence() error = nil, want unavailable")
	}
	assertApplicationError(t, err, ErrorCodePresenceQueryUnavailable)
}

func TestHandleGetPlayerPresenceRouteRejectsMalformedPayload(t *testing.T) {
	result, err := (Handlers{Registry: connection.NewInMemoryRegistry(nil)}).HandleGetPlayerPresenceRoute(context.Background(), app.RouteRequest{
		Route:    GetPlayerPresenceRoute(),
		Identity: validatedIdentity("player-1"),
		Payload:  "wrong",
	})
	if err == nil {
		t.Fatal("HandleGetPlayerPresenceRoute() error = nil, want malformed payload error")
	}
	if result.Error == nil || result.Error.Code != ErrorCodePresenceQueryForbidden {
		t.Fatalf("result error = %#v, want presence forbidden", result.Error)
	}
}

func TestGetPlayerPresenceReturnsCopies(t *testing.T) {
	registry := connection.NewInMemoryRegistry(fixedClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), connection.OpenConnection{
		ConnectionID:    "connection-1",
		ConnectionEpoch: 7,
		OpenedAt:        fixedTime(),
	}); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	if _, err := registry.BindConnectionIdentity(context.Background(), connection.BindIdentity{
		ConnectionID:     "connection-1",
		ConnectionEpoch:  7,
		ActorKind:        connection.ActorKindPlayer,
		PlayerID:         "player-1",
		RuntimeSessionID: "session-1",
		ValidatedAt:      fixedTime(),
	}); err != nil {
		t.Fatalf("BindConnectionIdentity() error = %v, want nil", err)
	}

	handlers := Handlers{Registry: registry}
	result, err := handlers.GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{}, validatedIdentity("player-1"))
	if err != nil {
		t.Fatalf("GetPlayerPresence() error = %v, want nil", err)
	}
	if result.LastSeenAt == nil || result.ActiveConnections[0].LastSeenAt == nil {
		t.Fatalf("GetPlayerPresence() = %#v, want time pointers", result)
	}
	mutated := fixedTime().Add(time.Hour)
	*result.LastSeenAt = mutated
	*result.ActiveConnections[0].LastSeenAt = mutated
	result.ActiveConnections[0].ConnectionID = "mutated"
	result.RuntimeSessionIDs[0] = "mutated"

	resultAgain, err := handlers.GetPlayerPresence(context.Background(), GetPlayerPresenceRequest{}, validatedIdentity("player-1"))
	if err != nil {
		t.Fatalf("GetPlayerPresence() again error = %v, want nil", err)
	}
	if resultAgain.ActiveConnections[0].ConnectionID != "connection-1" ||
		resultAgain.ActiveConnections[0].LastSeenAt.Equal(mutated) ||
		resultAgain.RuntimeSessionIDs[0] != "session-1" {
		t.Fatalf("GetPlayerPresence() returned aliased data = %#v", resultAgain)
	}
}

func validatedIdentity(playerID string) app.RequestIdentity {
	return app.RequestIdentity{
		Status:            app.IdentityValidationValidated,
		ActorKind:         app.ActorKindPlayer,
		ActorID:           playerID,
		PlayerID:          playerID,
		PlayerIDValidated: true,
	}
}

func assertApplicationError(t *testing.T, err error, code app.ErrorCode) {
	t.Helper()
	appErr, ok := err.(*app.ApplicationError)
	if !ok {
		t.Fatalf("error = %T %v, want *app.ApplicationError", err, err)
	}
	if appErr.Code != code {
		t.Fatalf("ApplicationError.Code = %q, want %q", appErr.Code, code)
	}
}

func fixedTime() time.Time {
	return time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}
