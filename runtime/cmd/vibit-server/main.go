package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/app/bootstrap"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/iceiko/vibit/runtime/internal/platform/persistence/postgres"
	vibitprotobuf "github.com/iceiko/vibit/runtime/internal/platform/protocol/protobuf"
	"github.com/iceiko/vibit/runtime/internal/platform/transport/ws"
)

const websocketPath = "/v1/ws"

const (
	envRuntimeStore = "VIBIT_RUNTIME_STORE"

	runtimeStoreMemory   = "memory"
	runtimeStorePostgres = "postgres"
)

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
	handler, cleanup, err := newHTTPHandlerFromEnv(context.Background(), os.LookupEnv)
	if err != nil {
		return nil, err
	}
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	server.RegisterOnShutdown(cleanup)
	return server, nil
}

func newHTTPHandler() (http.Handler, error) {
	dispatcher, err := bootstrap.NewInMemoryInventoryDispatcher()
	if err != nil {
		return nil, err
	}
	return newHTTPHandlerWithDispatcher(dispatcher), nil
}

func newHTTPHandlerFromEnv(ctx context.Context, lookup func(string) (string, bool)) (http.Handler, func(), error) {
	store := runtimeStoreFromEnv(lookup)
	switch store {
	case runtimeStoreMemory:
		handler, err := newHTTPHandler()
		return handler, func() {}, err
	case runtimeStorePostgres:
		return newPostgresHTTPHandler(ctx, lookup)
	default:
		return nil, func() {}, fmt.Errorf("unsupported %s %q", envRuntimeStore, store)
	}
}

func newPostgresHTTPHandler(ctx context.Context, lookup func(string) (string, bool)) (http.Handler, func(), error) {
	cfg, err := postgres.ConfigFromEnvMap(lookup)
	if err != nil {
		return nil, func() {}, err
	}

	pool, err := postgres.OpenPool(ctx, cfg)
	if err != nil {
		return nil, func() {}, err
	}

	queryRepository := postgres.NewInventoryRepositoryForUnitOfWork(pool)
	persistentDispatcher, err := bootstrap.NewInventoryDispatcher(bootstrap.InventoryOptions{
		Repositories: bootstrap.PostgresInventoryRepositoryProvider{
			QueryRepository: queryRepository,
		},
		PermissionPolicy: inventory.StaticPermissionPolicy{GrantAllowed: true, ReadAllowed: true},
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 256},
		EventIDs:         &inventory.IncrementingEventIDGenerator{Prefix: "inventory-event"},
		Clock:            inventory.SystemClock{},
	})
	if err != nil {
		pool.Close()
		return nil, func() {}, err
	}

	transactionalDispatcher := app.TransactionalDispatcher{
		Dispatcher: persistentDispatcher,
		Runner:     postgres.NewPoolRunner(pool),
	}
	return newHTTPHandlerWithDispatcher(transactionalDispatcher), pool.Close, nil
}

func newHTTPHandlerWithDispatcher(dispatcher vibitprotobuf.ApplicationDispatcher) http.Handler {
	protocolHandler := newProtocolFrameHandlerWithDispatcher(dispatcher)
	mux := http.NewServeMux()
	mux.Handle(websocketPath, &ws.Server{
		Handler: websocketProtocolHandler{handler: protocolHandler},
	})
	return mux
}

func newProtocolFrameHandlerWithDispatcher(dispatcher vibitprotobuf.ApplicationDispatcher) vibitprotobuf.FrameHandler {
	return vibitprotobuf.FrameHandler{Dispatcher: app.SessionValidatingDispatcher{
		Dispatcher: dispatcher,
		Validator:  app.MetadataOnlySessionValidator{},
	}}
}

func runtimeStoreFromEnv(lookup func(string) (string, bool)) string {
	if lookup == nil {
		return runtimeStoreMemory
	}
	if value, ok := lookup(envRuntimeStore); ok {
		value = strings.ToLower(strings.TrimSpace(value))
		if value != "" {
			return value
		}
	}
	return runtimeStoreMemory
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
