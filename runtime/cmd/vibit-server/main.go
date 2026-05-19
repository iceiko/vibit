package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
	"github.com/iceiko/vibit/runtime/internal/app/bootstrap"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/iceiko/vibit/runtime/internal/platform/persistence/postgres"
	vibitprotobuf "github.com/iceiko/vibit/runtime/internal/platform/protocol/protobuf"
	"github.com/iceiko/vibit/runtime/internal/platform/transport/ws"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

const websocketPath = "/v1/ws"

const (
	envRuntimeStore            = "VIBIT_RUNTIME_STORE"
	envAuthAccessTokenLifetime = "VIBIT_AUTH_ACCESS_TOKEN_TTL"
	envAuthTokenAudience       = "VIBIT_AUTH_TOKEN_AUDIENCE"

	runtimeStoreMemory   = "memory"
	runtimeStorePostgres = "postgres"

	defaultAuthAccessTokenLifetime = time.Hour
	defaultAuthTokenAudience       = "vibit_gameplay_runtime_requests"
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
	return newHTTPHandlerWithDispatcher(dispatcher, nil, nil), nil
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

	transactionRunner := postgres.NewPoolRunner(pool)
	authComposition, err := newAuthenticationStartupComposition(lookup, transactionRunner)
	if err != nil {
		pool.Close()
		return nil, func() {}, err
	}
	if err := (bootstrap.AuthenticationRouteHandlers{Service: authComposition.Service}).RegisterRoutes(persistentDispatcher); err != nil {
		pool.Close()
		return nil, func() {}, err
	}

	transactionalDispatcher := app.TransactionalDispatcher{
		Dispatcher: persistentDispatcher,
		Runner:     transactionRunner,
		BypassRoutes: []app.RouteKey{
			app.AuthenticateWithDeviceCredentialRoute(),
			app.LogoutAccessTokenRoute(),
		},
	}
	return newHTTPHandlerWithDispatcher(transactionalDispatcher, authComposition.RouteProtector, authComposition.ConnectionBinder), pool.Close, nil
}

func newHTTPHandlerWithDispatcher(dispatcher vibitprotobuf.ApplicationDispatcher, routeProtector vibitprotobuf.RouteProtector, connectionBinder vibitprotobuf.ConnectionBinder) http.Handler {
	protocolHandler := newProtocolFrameHandlerWithDispatcher(dispatcher, routeProtector, connectionBinder)
	mux := http.NewServeMux()
	mux.Handle(websocketPath, &ws.Server{
		Handler: websocketProtocolHandler{handler: protocolHandler},
	})
	return mux
}

func newProtocolFrameHandlerWithDispatcher(dispatcher vibitprotobuf.ApplicationDispatcher, routeProtector vibitprotobuf.RouteProtector, connectionBinder vibitprotobuf.ConnectionBinder) vibitprotobuf.FrameHandler {
	return vibitprotobuf.FrameHandler{
		Dispatcher: app.SessionValidatingDispatcher{
			Dispatcher: dispatcher,
			Validator:  app.MetadataOnlySessionValidator{},
		},
		RouteProtector:   routeProtector,
		ConnectionBinder: connectionBinder,
	}
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

type authenticationStartupComposition struct {
	RouteProtector   app.RouteProtector
	ConnectionBinder app.ConnectionBinder
	Service          appauth.Service
}

func newAuthenticationStartupComposition(lookup func(string) (string, bool), unitOfWorkRunner tx.Runner) (authenticationStartupComposition, error) {
	cfg, err := authenticationStartupConfigFromEnv(lookup)
	if err != nil {
		return authenticationStartupComposition{}, err
	}

	service, err := appauth.NewService(appauth.ServiceDependencies{
		UnitOfWorkRunner:       unitOfWorkRunner,
		VerifierKeySet:         cfg.VerifierKeySet,
		AccessTokenRandom:      rand.Reader,
		Clock:                  systemAuthClock{},
		TokenRecordIDGenerator: randomTokenRecordIDGenerator{},
		SessionIDGenerator:     randomSessionIDGenerator{},
		AccessTokenLifetime:    cfg.AccessTokenLifetime,
		TokenAudience:          cfg.TokenAudience,
	})
	if err != nil {
		return authenticationStartupComposition{}, err
	}

	return authenticationStartupComposition{
		RouteProtector:   app.NewRouteProtector(appauth.NewRouteAccessTokenValidator(service)),
		ConnectionBinder: app.NewConnectionBinder(appauth.NewRouteAccessTokenValidator(service)),
		Service:          service,
	}, nil
}

type authenticationStartupConfig struct {
	VerifierKeySet      appauth.VerifierKeySet
	AccessTokenLifetime time.Duration
	TokenAudience       string
}

func authenticationStartupConfigFromEnv(lookup func(string) (string, bool)) (authenticationStartupConfig, error) {
	if lookup == nil {
		return authenticationStartupConfig{}, errors.New("authentication startup: environment lookup is required")
	}

	verifierKeySet, err := appauth.LoadVerifierKeySetFromEnvironment(lookup)
	if err != nil {
		return authenticationStartupConfig{}, fmt.Errorf("authentication startup: verifier key configuration: %w", err)
	}

	accessTokenLifetime, err := authAccessTokenLifetimeFromEnv(lookup)
	if err != nil {
		return authenticationStartupConfig{}, err
	}

	return authenticationStartupConfig{
		VerifierKeySet:      verifierKeySet,
		AccessTokenLifetime: accessTokenLifetime,
		TokenAudience:       authTokenAudienceFromEnv(lookup),
	}, nil
}

func authAccessTokenLifetimeFromEnv(lookup func(string) (string, bool)) (time.Duration, error) {
	if lookup == nil {
		return 0, errors.New("authentication startup: environment lookup is required")
	}

	value, ok := lookup(envAuthAccessTokenLifetime)
	value = strings.TrimSpace(value)
	if !ok || value == "" {
		return defaultAuthAccessTokenLifetime, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("authentication startup: %s must be a Go duration: %w", envAuthAccessTokenLifetime, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("authentication startup: %s must be positive", envAuthAccessTokenLifetime)
	}
	return duration, nil
}

func authTokenAudienceFromEnv(lookup func(string) (string, bool)) string {
	if lookup == nil {
		return defaultAuthTokenAudience
	}
	value, ok := lookup(envAuthTokenAudience)
	value = strings.TrimSpace(value)
	if !ok || value == "" {
		return defaultAuthTokenAudience
	}
	return value
}

type systemAuthClock struct{}

func (systemAuthClock) Now() time.Time {
	return time.Now().UTC()
}

type randomTokenRecordIDGenerator struct{}

func (randomTokenRecordIDGenerator) GenerateTokenRecordID(context.Context) (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("authentication startup: generate token record id: %w", err)
	}
	return "auth-token-" + hex.EncodeToString(raw[:]), nil
}

type randomSessionIDGenerator struct{}

func (randomSessionIDGenerator) GenerateSessionID(context.Context) (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("authentication startup: generate runtime session id: %w", err)
	}
	return "runtime-session-" + hex.EncodeToString(raw[:]), nil
}

type websocketProtocolHandler struct {
	handler vibitprotobuf.FrameHandler
}

func (h websocketProtocolHandler) HandleFrame(ctx context.Context, frame ws.Frame) ([][]byte, error) {
	return h.handler.HandleFrame(ctx, vibitprotobuf.FrameRequest{
		ConnectionID:    frame.ConnectionID,
		ConnectionEpoch: frame.ConnectionEpoch,
		RemoteAddr:      frame.RemoteAddr,
		Payload:         frame.Payload,
	})
}
