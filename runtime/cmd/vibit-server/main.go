package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/iceiko/vibit/runtime/internal/app/bootstrap"
	vibitprotobuf "github.com/iceiko/vibit/runtime/internal/platform/protocol/protobuf"
	"github.com/iceiko/vibit/runtime/internal/platform/transport/ws"
)

const websocketPath = "/v1/ws"

func main() {
	addr := os.Getenv("VIBIT_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server, err := newServer(addr)
	if err != nil {
		log.Fatalf("build server: %v", err)
	}

	log.Printf("vibit runtime listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serve: %v", err)
	}
}

func newServer(addr string) (*http.Server, error) {
	handler, err := newHTTPHandler()
	if err != nil {
		return nil, err
	}
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}, nil
}

func newHTTPHandler() (http.Handler, error) {
	dispatcher, err := bootstrap.NewInMemoryInventoryDispatcher()
	if err != nil {
		return nil, err
	}

	protocolHandler := vibitprotobuf.FrameHandler{Dispatcher: dispatcher}
	mux := http.NewServeMux()
	mux.Handle(websocketPath, &ws.Server{
		Handler: websocketProtocolHandler{handler: protocolHandler},
	})
	return mux, nil
}

type websocketProtocolHandler struct {
	handler vibitprotobuf.FrameHandler
}

func (h websocketProtocolHandler) HandleFrame(ctx context.Context, frame ws.Frame) ([][]byte, error) {
	return h.handler.HandleFrame(ctx, vibitprotobuf.FrameRequest{
		ConnectionID: frame.ConnectionID,
		RemoteAddr:   frame.RemoteAddr,
		Payload:      frame.Payload,
	})
}
