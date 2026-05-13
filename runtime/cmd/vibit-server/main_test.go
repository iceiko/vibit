package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
