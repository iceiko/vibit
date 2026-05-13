package ws

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestServerBinaryFrameRoundTrip(t *testing.T) {
	received := make(chan Frame, 1)
	server := httptest.NewServer(&Server{
		AcceptOptions: &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		},
		Handler: FrameHandlerFunc(func(ctx context.Context, frame Frame) ([][]byte, error) {
			received <- frame
			return [][]byte{[]byte("ok:" + string(frame.Payload))}, nil
		}),
	})
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, websocketURL(server.URL), nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.CloseNow()

	if err := conn.Write(ctx, websocket.MessageBinary, []byte("ping")); err != nil {
		t.Fatalf("write binary frame: %v", err)
	}

	messageType, payload, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read response frame: %v", err)
	}
	if messageType != websocket.MessageBinary {
		t.Fatalf("response message type = %v, want %v", messageType, websocket.MessageBinary)
	}
	if string(payload) != "ok:ping" {
		t.Fatalf("response payload = %q, want %q", string(payload), "ok:ping")
	}

	select {
	case frame := <-received:
		if frame.ConnectionID == "" {
			t.Fatal("frame ConnectionID is empty")
		}
		if frame.RemoteAddr == "" {
			t.Fatal("frame RemoteAddr is empty")
		}
		if !bytes.Equal(frame.Payload, []byte("ping")) {
			t.Fatalf("frame payload = %q, want %q", string(frame.Payload), "ping")
		}
	case <-ctx.Done():
		t.Fatal("handler did not receive frame")
	}
}

func TestFrameCopiesPayload(t *testing.T) {
	payload := []byte("ping")
	frame := newFrame(connectionMetadata{
		id:         "ws-1",
		remoteAddr: "127.0.0.1:1",
	}, payload)

	payload[0] = 'x'

	if string(frame.Payload) != "ping" {
		t.Fatalf("frame payload = %q, want %q", string(frame.Payload), "ping")
	}
	if frame.ConnectionID != "ws-1" {
		t.Fatalf("frame ConnectionID = %q, want %q", frame.ConnectionID, "ws-1")
	}
	if frame.RemoteAddr != "127.0.0.1:1" {
		t.Fatalf("frame RemoteAddr = %q, want %q", frame.RemoteAddr, "127.0.0.1:1")
	}
}

func TestServerRejectsTextFrames(t *testing.T) {
	var called atomic.Bool
	server := httptest.NewServer(&Server{
		AcceptOptions: &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		},
		Handler: FrameHandlerFunc(func(ctx context.Context, frame Frame) ([][]byte, error) {
			called.Store(true)
			return nil, nil
		}),
	})
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, websocketURL(server.URL), nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.CloseNow()

	if err := conn.Write(ctx, websocket.MessageText, []byte("not-binary")); err != nil {
		t.Fatalf("write text frame: %v", err)
	}

	_, _, err = conn.Read(ctx)
	if err == nil {
		t.Fatal("read after text frame succeeded, want close error")
	}
	if status := websocket.CloseStatus(err); status != websocket.StatusUnsupportedData {
		t.Fatalf("close status = %v, want %v: %v", status, websocket.StatusUnsupportedData, err)
	}
	if called.Load() {
		t.Fatal("handler was called for a text frame")
	}
}

func TestServerRequiresFrameHandler(t *testing.T) {
	server := httptest.NewServer(&Server{
		AcceptOptions: &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		},
	})
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, resp, err := websocket.Dial(ctx, websocketURL(server.URL), nil)
	if conn != nil {
		defer conn.CloseNow()
	}
	if err == nil {
		t.Fatal("dial succeeded without a frame handler")
	}
	if resp == nil {
		t.Fatal("dial response is nil")
	}
	if resp.StatusCode != 500 {
		t.Fatalf("status code = %d, want 500", resp.StatusCode)
	}
}

func TestServerClosesWhenFrameHandlerFails(t *testing.T) {
	server := httptest.NewServer(&Server{
		AcceptOptions: &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		},
		Handler: FrameHandlerFunc(func(ctx context.Context, frame Frame) ([][]byte, error) {
			return nil, fmt.Errorf("handler failed")
		}),
	})
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, websocketURL(server.URL), nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.CloseNow()

	if err := conn.Write(ctx, websocket.MessageBinary, []byte("ping")); err != nil {
		t.Fatalf("write binary frame: %v", err)
	}

	_, _, err = conn.Read(ctx)
	if err == nil {
		t.Fatal("read after handler failure succeeded, want close error")
	}
	if status := websocket.CloseStatus(err); status != websocket.StatusInternalError {
		t.Fatalf("close status = %v, want %v: %v", status, websocket.StatusInternalError, err)
	}
}

func websocketURL(httpURL string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http")
}
