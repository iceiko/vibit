package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
	appconnection "github.com/iceiko/vibit/runtime/internal/app/connection"
	appfriends "github.com/iceiko/vibit/runtime/internal/app/friends"
	apppresence "github.com/iceiko/vibit/runtime/internal/app/presence"
	appstorage "github.com/iceiko/vibit/runtime/internal/app/storage"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/iceiko/vibit/runtime/internal/platform/persistence/postgres"
	vibitprotobuf "github.com/iceiko/vibit/runtime/internal/platform/protocol/protobuf"
	"github.com/iceiko/vibit/runtime/internal/platform/transport/ws"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

const (
	websocketPath = "/v1/ws"
	healthPath    = "/healthz"
	readinessPath = "/readyz"
	versionPath   = "/version"
	configPath    = "/configz"
)

const (
	runtimeVersion             = "0.1.0-pre-alpha"
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
	return newHTTPHandlerWithRuntimeInfo(runtimeInfoFromEnv(nil, runtimeStoreMemory, false))
}

func newHTTPHandlerWithRuntimeInfo(info runtimeInfo) (http.Handler, error) {
	dispatcher, err := bootstrap.NewInMemoryInventoryDispatcher()
	if err != nil {
		return nil, err
	}
	return newHTTPHandlerWithDispatcher(dispatcher, nil, nil, nil, info), nil
}

func newHTTPHandlerFromEnv(ctx context.Context, lookup func(string) (string, bool)) (http.Handler, func(), error) {
	store := runtimeStoreFromEnv(lookup)
	switch store {
	case runtimeStoreMemory:
		handler, err := newHTTPHandlerWithRuntimeInfo(runtimeInfoFromEnv(lookup, runtimeStoreMemory, false))
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

	connectionRegistry := appconnection.NewInMemoryRegistry(nil)
	if err := registerPresenceRoutes(persistentDispatcher, connectionRegistry); err != nil {
		pool.Close()
		return nil, func() {}, err
	}

	storageService, err := appstorage.NewService(appstorage.ServiceDependencies{
		UnitOfWorkRunner:  transactionRunner,
		ObjectIDGenerator: randomStorageObjectIDGenerator{},
	})
	if err != nil {
		pool.Close()
		return nil, func() {}, err
	}
	if err := (bootstrap.StorageRouteHandlers{Service: storageService}).RegisterRoutes(persistentDispatcher); err != nil {
		pool.Close()
		return nil, func() {}, err
	}

	friendsService, err := appfriends.NewService(appfriends.ServiceDependencies{
		UnitOfWorkRunner:        transactionRunner,
		RelationshipIDGenerator: randomFriendRelationshipIDGenerator{},
	})
	if err != nil {
		pool.Close()
		return nil, func() {}, err
	}
	if err := (bootstrap.FriendsRouteHandlers{Service: friendsService}).RegisterRoutes(persistentDispatcher); err != nil {
		pool.Close()
		return nil, func() {}, err
	}

	transactionalDispatcher := app.TransactionalDispatcher{
		Dispatcher: persistentDispatcher,
		Runner:     transactionRunner,
		BypassRoutes: []app.RouteKey{
			app.AuthenticateWithDeviceCredentialRoute(),
			app.LogoutAccessTokenRoute(),
			appstorage.PutOwnStorageObjectRoute(),
			appstorage.DeleteOwnStorageObjectRoute(),
			appfriends.SendFriendRequestRoute(),
			appfriends.AcceptFriendRequestRoute(),
			appfriends.RejectFriendRequestRoute(),
			appfriends.RemoveFriendRoute(),
			appfriends.BlockPlayerRoute(),
			appfriends.UnblockPlayerRoute(),
		},
	}
	connectionBinder := registryConnectionBinder{
		binder:   authComposition.ConnectionBinder,
		registry: connectionRegistry,
	}
	return newHTTPHandlerWithDispatcher(
		transactionalDispatcher,
		authComposition.RouteProtector,
		connectionBinder,
		connectionLifecycleRegistryObserver{registry: connectionRegistry},
		runtimeInfoFromEnv(lookup, runtimeStorePostgres, true),
	), pool.Close, nil
}

func registerPresenceRoutes(dispatcher *app.Dispatcher, registry *appconnection.InMemoryRegistry) error {
	return (apppresence.Handlers{Registry: registry}).RegisterRoutes(dispatcher)
}

func newHTTPHandlerWithDispatcher(
	dispatcher vibitprotobuf.ApplicationDispatcher,
	routeProtector vibitprotobuf.RouteProtector,
	connectionBinder vibitprotobuf.ConnectionBinder,
	lifecycle ws.ConnectionLifecycleObserver,
	info runtimeInfo,
) http.Handler {
	protocolHandler := newProtocolFrameHandlerWithDispatcher(dispatcher, routeProtector, connectionBinder)
	mux := http.NewServeMux()
	registerStatusHandlers(mux, info)
	mux.Handle(websocketPath, &ws.Server{
		Handler:   websocketProtocolHandler{handler: protocolHandler},
		Lifecycle: lifecycle,
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

type runtimeInfo struct {
	Version                      string
	RuntimeStore                 string
	WebSocketPath                string
	LocalAlphaRequestLoopScript  string
	PostgresConfigured           bool
	AuthenticationConfigRequired bool
	SecretsRedacted              bool
}

func runtimeInfoFromEnv(lookup func(string) (string, bool), store string, authenticationConfigRequired bool) runtimeInfo {
	return runtimeInfo{
		Version:                      runtimeVersion,
		RuntimeStore:                 store,
		WebSocketPath:                websocketPath,
		LocalAlphaRequestLoopScript:  "examples/local-alpha-request-loop.sh",
		PostgresConfigured:           envPresent(lookup, postgres.EnvDatabaseURL),
		AuthenticationConfigRequired: authenticationConfigRequired,
		SecretsRedacted:              true,
	}
}

func envPresent(lookup func(string) (string, bool), name string) bool {
	if lookup == nil {
		return false
	}
	value, ok := lookup(name)
	return ok && strings.TrimSpace(value) != ""
}

func registerStatusHandlers(mux *http.ServeMux, info runtimeInfo) {
	mux.HandleFunc(healthPath, func(w http.ResponseWriter, r *http.Request) {
		writeStatusJSON(w, http.StatusOK, map[string]any{
			"status":  "ok",
			"service": "vibit-runtime",
		})
	})
	mux.HandleFunc(readinessPath, func(w http.ResponseWriter, r *http.Request) {
		writeStatusJSON(w, http.StatusOK, map[string]any{
			"status":         "ready",
			"runtime_store":  info.RuntimeStore,
			"websocket_path": info.WebSocketPath,
		})
	})
	mux.HandleFunc(versionPath, func(w http.ResponseWriter, r *http.Request) {
		writeStatusJSON(w, http.StatusOK, map[string]any{
			"version": info.Version,
			"status":  "pre_alpha",
		})
	})
	mux.HandleFunc(configPath, func(w http.ResponseWriter, r *http.Request) {
		writeStatusJSON(w, http.StatusOK, map[string]any{
			"runtime_store":                  info.RuntimeStore,
			"websocket_path":                 info.WebSocketPath,
			"local_alpha_request_loop":       info.LocalAlphaRequestLoopScript,
			"postgres_configured":            info.PostgresConfigured,
			"authentication_config_required": info.AuthenticationConfigRequired,
			"secrets_redacted":               info.SecretsRedacted,
		})
	})
}

func writeStatusJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
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
		UnitOfWorkRunner:              unitOfWorkRunner,
		VerifierKeySet:                cfg.VerifierKeySet,
		AccessTokenRandom:             rand.Reader,
		DeviceCredentialRandom:        rand.Reader,
		Clock:                         systemAuthClock{},
		TokenRecordIDGenerator:        randomTokenRecordIDGenerator{},
		SessionIDGenerator:            randomSessionIDGenerator{},
		PlayerIDGenerator:             randomPlayerIDGenerator{},
		PlayerAccountEventIDGenerator: randomPlayerAccountEventIDGenerator{},
		CredentialRecordIDGenerator:   randomCredentialRecordIDGenerator{},
		AccessTokenLifetime:           cfg.AccessTokenLifetime,
		TokenAudience:                 cfg.TokenAudience,
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

type randomPlayerIDGenerator struct{}

func (randomPlayerIDGenerator) GeneratePlayerID(context.Context) (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("authentication startup: generate player id: %w", err)
	}
	return "player-" + hex.EncodeToString(raw[:]), nil
}

type randomPlayerAccountEventIDGenerator struct{}

func (randomPlayerAccountEventIDGenerator) GeneratePlayerAccountEventID(context.Context) (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("authentication startup: generate player account event id: %w", err)
	}
	return "player-account-event-" + hex.EncodeToString(raw[:]), nil
}

type randomCredentialRecordIDGenerator struct{}

func (randomCredentialRecordIDGenerator) GenerateCredentialRecordID(context.Context) (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("authentication startup: generate credential record id: %w", err)
	}
	return "auth-credential-" + hex.EncodeToString(raw[:]), nil
}

type randomStorageObjectIDGenerator struct{}

func (randomStorageObjectIDGenerator) GenerateStorageObjectID(context.Context) (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("storage startup: generate storage object id: %w", err)
	}
	return "storage-object-" + hex.EncodeToString(raw[:]), nil
}

type randomFriendRelationshipIDGenerator struct{}

func (randomFriendRelationshipIDGenerator) GenerateFriendRelationshipID(context.Context) (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("friends startup: generate friend relationship id: %w", err)
	}
	return "friend-relationship-" + hex.EncodeToString(raw[:]), nil
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
