package main

import (
	"net/http"
	"net/http/httptest"
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
