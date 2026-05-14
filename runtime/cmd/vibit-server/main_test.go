package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iceiko/vibit/runtime/internal/app"
)

func TestWebSocketEndpointIsMounted(t *testing.T) {
	handler, err := newHTTPHandler()
	if err != nil {
		t.Fatalf("newHTTPHandler() error = %v, want nil", err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	response, err := http.Get(server.URL + websocketPath)
	if err != nil {
		t.Fatalf("GET /v1/ws error = %v, want response", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		t.Fatal("/v1/ws returned 404, want mounted WebSocket endpoint")
	}
}

func TestHTTPHandlerFromEnvDefaultsToMemoryStore(t *testing.T) {
	handler, cleanup, err := newHTTPHandlerFromEnv(context.Background(), func(string) (string, bool) {
		return "", false
	})
	if err != nil {
		t.Fatalf("newHTTPHandlerFromEnv() error = %v, want nil", err)
	}
	defer cleanup()
	if handler == nil {
		t.Fatal("newHTTPHandlerFromEnv() handler = nil, want handler")
	}
}

func TestHTTPHandlerFromEnvRejectsUnsupportedStore(t *testing.T) {
	_, _, err := newHTTPHandlerFromEnv(context.Background(), func(name string) (string, bool) {
		if name == envRuntimeStore {
			return "unknown", true
		}
		return "", false
	})
	if err == nil {
		t.Fatal("newHTTPHandlerFromEnv() error = nil, want unsupported store error")
	}
	if !strings.Contains(err.Error(), envRuntimeStore) {
		t.Fatalf("newHTTPHandlerFromEnv() error = %v, want %s context", err, envRuntimeStore)
	}
}

func TestHTTPHandlerFromEnvRequiresPostgresDSNForPostgresStore(t *testing.T) {
	_, _, err := newHTTPHandlerFromEnv(context.Background(), func(name string) (string, bool) {
		if name == envRuntimeStore {
			return runtimeStorePostgres, true
		}
		return "", false
	})
	if err == nil {
		t.Fatal("newHTTPHandlerFromEnv(postgres) error = nil, want missing DSN error")
	}
	if !strings.Contains(err.Error(), "DSN") {
		t.Fatalf("newHTTPHandlerFromEnv(postgres) error = %v, want DSN context", err)
	}
}

func TestHTTPHandlerWithDispatcherAppliesMetadataOnlySessionValidation(t *testing.T) {
	frameHandler := newProtocolFrameHandlerWithDispatcher(routeDispatcherFunc(func(_ context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
		if request.Identity.Status != app.IdentityValidationMetadataOnly {
			t.Fatalf("request Identity.Status = %q, want %q", request.Identity.Status, app.IdentityValidationMetadataOnly)
		}
		if request.Identity.PlayerID != "player-1" || request.Identity.SessionID != "session-1" {
			t.Fatalf("request Identity = %#v, want normalized player/session metadata", request.Identity)
		}
		if request.Identity.PlayerIDValidated || request.Identity.SessionValidated {
			t.Fatalf("request Identity validation flags = %#v, want metadata-only identity", request.Identity)
		}
		return app.ApplicationResult{
			RequestID: request.RequestID,
			Route:     request.Route,
			Target:    request.Target,
			Session:   request.Session,
			Identity:  request.Identity,
		}, nil
	}))

	result, err := frameHandler.Dispatcher.Dispatch(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     app.RouteKey{Kind: app.MessageKindQuery, Module: "inventory", Name: "GetInventory"},
		Session:   app.Session{SessionID: "session-1", PlayerID: "player-1"},
	})
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}
	if result.Identity.Status != app.IdentityValidationMetadataOnly {
		t.Fatalf("result Identity.Status = %q, want %q", result.Identity.Status, app.IdentityValidationMetadataOnly)
	}
}

type routeDispatcherFunc func(context.Context, app.RouteRequest) (app.ApplicationResult, error)

func (f routeDispatcherFunc) Dispatch(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	return f(ctx, request)
}
