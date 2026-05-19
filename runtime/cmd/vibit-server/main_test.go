package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
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
	}), nil, nil)

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

func TestProtocolFrameHandlerCanReceiveRouteProtector(t *testing.T) {
	validator := staticRouteAccessTokenValidator{
		result: app.RouteAccessTokenValidationResult{
			Valid: true,
			Identity: app.ValidatedPlayerIdentity("player-1", app.Session{
				ConnectionID: "connection-1",
				PlayerID:     "player-1",
			}),
		},
	}
	protector := app.NewRouteProtector(validator)

	frameHandler := newProtocolFrameHandlerWithDispatcher(routeDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
		t.Fatal("dispatcher should not be called by this test")
		return app.ApplicationResult{}, nil
	}), protector, nil)

	if frameHandler.RouteProtector == nil {
		t.Fatal("RouteProtector = nil, want injected route protector")
	}
}

func TestProtocolFrameHandlerCanReceiveConnectionBinder(t *testing.T) {
	binder := app.NewConnectionBinder(staticRouteAccessTokenValidator{
		result: app.RouteAccessTokenValidationResult{
			Valid: true,
			Identity: app.ValidatedPlayerIdentity("player-1", app.Session{
				ConnectionID: "connection-1",
				PlayerID:     "player-1",
			}),
		},
	})

	frameHandler := newProtocolFrameHandlerWithDispatcher(routeDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
		t.Fatal("dispatcher should not be called by this test")
		return app.ApplicationResult{}, nil
	}), nil, binder)

	if frameHandler.ConnectionBinder == nil {
		t.Fatal("ConnectionBinder = nil, want injected connection binder")
	}
}

func TestAuthenticationStartupConfigFromEnvUsesDefaultsAndRedactsKeyValues(t *testing.T) {
	cfg, err := authenticationStartupConfigFromEnv(mapLookup(authStartupEnv(nil)))
	if err != nil {
		t.Fatalf("authenticationStartupConfigFromEnv() error = %v, want nil", err)
	}
	if cfg.AccessTokenLifetime != defaultAuthAccessTokenLifetime {
		t.Fatalf("AccessTokenLifetime = %v, want %v", cfg.AccessTokenLifetime, defaultAuthAccessTokenLifetime)
	}
	if cfg.TokenAudience != defaultAuthTokenAudience {
		t.Fatalf("TokenAudience = %q, want %q", cfg.TokenAudience, defaultAuthTokenAudience)
	}
	if cfg.VerifierKeySet.KeySetID() != "test-key-set" {
		t.Fatalf("KeySetID = %q, want test-key-set", cfg.VerifierKeySet.KeySetID())
	}
}

func TestAuthenticationStartupConfigFromEnvReadsExplicitLifetimeAndAudience(t *testing.T) {
	cfg, err := authenticationStartupConfigFromEnv(mapLookup(authStartupEnv(map[string]string{
		envAuthAccessTokenLifetime: "30m",
		envAuthTokenAudience:       " gameplay-test ",
	})))
	if err != nil {
		t.Fatalf("authenticationStartupConfigFromEnv() error = %v, want nil", err)
	}
	if cfg.AccessTokenLifetime != 30*time.Minute {
		t.Fatalf("AccessTokenLifetime = %v, want 30m", cfg.AccessTokenLifetime)
	}
	if cfg.TokenAudience != "gameplay-test" {
		t.Fatalf("TokenAudience = %q, want gameplay-test", cfg.TokenAudience)
	}
}

func TestAuthenticationStartupConfigFromEnvRejectsMissingVerifierKeyConfiguration(t *testing.T) {
	_, err := authenticationStartupConfigFromEnv(mapLookup(map[string]string{}))
	if err == nil {
		t.Fatal("authenticationStartupConfigFromEnv() error = nil, want missing verifier key config error")
	}
	if !strings.Contains(err.Error(), "verifier key configuration") {
		t.Fatalf("authenticationStartupConfigFromEnv() error = %v, want verifier key configuration context", err)
	}
	secretValues := authStartupEnv(nil)
	for _, value := range secretValues {
		if value != "" && strings.Contains(err.Error(), value) {
			t.Fatalf("authenticationStartupConfigFromEnv() error leaked secret value %q: %v", value, err)
		}
	}
}

func TestAuthenticationStartupConfigFromEnvRejectsInvalidAccessTokenLifetime(t *testing.T) {
	_, err := authenticationStartupConfigFromEnv(mapLookup(authStartupEnv(map[string]string{
		envAuthAccessTokenLifetime: "0s",
	})))
	if err == nil {
		t.Fatal("authenticationStartupConfigFromEnv() error = nil, want invalid lifetime error")
	}
	if !strings.Contains(err.Error(), envAuthAccessTokenLifetime) {
		t.Fatalf("authenticationStartupConfigFromEnv() error = %v, want lifetime env name", err)
	}
}

func TestAuthenticationStartupCompositionRequiresUnitOfWorkRunner(t *testing.T) {
	_, err := newAuthenticationStartupComposition(mapLookup(authStartupEnv(nil)), nil)
	if err == nil {
		t.Fatal("newAuthenticationStartupComposition() error = nil, want missing unit-of-work runner error")
	}
}

func TestAuthenticationStartupCompositionExposesRouteProtectorAndService(t *testing.T) {
	composition, err := newAuthenticationStartupComposition(mapLookup(authStartupEnv(nil)), startupNoopAuthenticationRunner{})
	if err != nil {
		t.Fatalf("newAuthenticationStartupComposition() error = %v, want nil", err)
	}
	routeResult, err := composition.RouteProtector.ProtectRoute(context.Background(), app.RouteProtectionRequest{
		Request: app.RouteRequest{Route: app.AuthenticateWithDeviceCredentialRoute()},
	})
	if err != nil {
		t.Fatalf("RouteProtector.ProtectRoute(public login) error = %v, want nil", err)
	}
	if !routeResult.Allowed || !routeResult.Public {
		t.Fatalf("route protection result = %#v, want public login route allowed", routeResult)
	}
	logoutRouteResult, err := composition.RouteProtector.ProtectRoute(context.Background(), app.RouteProtectionRequest{
		Request: app.RouteRequest{Route: app.LogoutAccessTokenRoute()},
	})
	if err != nil {
		t.Fatalf("RouteProtector.ProtectRoute(logout) error = %v, want nil", err)
	}
	if !logoutRouteResult.Allowed || !logoutRouteResult.Public {
		t.Fatalf("logout route protection result = %#v, want service-validated logout route allowed", logoutRouteResult)
	}
	if composition.ConnectionBinder.Validator == nil {
		t.Fatal("ConnectionBinder.Validator = nil, want startup connection binder")
	}
	_, err = composition.Service.ValidateAccessToken(context.Background(), appauth.AccessTokenValidationRequest{})
	if err == nil {
		t.Fatal("Service.ValidateAccessToken(empty) error = nil, want configured service to reject missing token")
	}
}

func TestRandomTokenRecordIDGeneratorShape(t *testing.T) {
	id, err := (randomTokenRecordIDGenerator{}).GenerateTokenRecordID(context.Background())
	if err != nil {
		t.Fatalf("GenerateTokenRecordID() error = %v, want nil", err)
	}
	if !regexp.MustCompile(`^auth-token-[0-9a-f]{32}$`).MatchString(id) {
		t.Fatalf("GenerateTokenRecordID() = %q, want auth-token- plus 32 lowercase hex chars", id)
	}
}

func TestRandomSessionIDGeneratorShape(t *testing.T) {
	id, err := (randomSessionIDGenerator{}).GenerateSessionID(context.Background())
	if err != nil {
		t.Fatalf("GenerateSessionID() error = %v, want nil", err)
	}
	if !regexp.MustCompile(`^runtime-session-[0-9a-f]{64}$`).MatchString(id) {
		t.Fatalf("GenerateSessionID() = %q, want runtime-session- plus 64 lowercase hex chars", id)
	}
}

type routeDispatcherFunc func(context.Context, app.RouteRequest) (app.ApplicationResult, error)

func (f routeDispatcherFunc) Dispatch(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	return f(ctx, request)
}

type staticRouteAccessTokenValidator struct {
	result app.RouteAccessTokenValidationResult
	err    error
}

type startupNoopAuthenticationRunner struct{}

func (startupNoopAuthenticationRunner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, tx.UnitOfWork) error) error {
	if fn == nil {
		return nil
	}
	return fn(ctx, tx.NoopUnitOfWork{})
}

func (v staticRouteAccessTokenValidator) ValidateRouteAccessToken(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
	return v.result, v.err
}

func authStartupEnv(overrides map[string]string) map[string]string {
	values := map[string]string{
		appauth.EnvVerifierKeySetID:      "test-key-set",
		appauth.EnvCredentialLookupKey:   encodedKey(1),
		appauth.EnvCredentialVerifierKey: encodedKey(2),
		appauth.EnvTokenLookupKey:        encodedKey(3),
		appauth.EnvTokenVerifierKey:      encodedKey(4),
	}
	for key, value := range overrides {
		values[key] = value
	}
	return values
}

func encodedKey(seed byte) string {
	value := make([]byte, appauth.MinVerifierKeyBytes)
	for i := range value {
		value[i] = seed + byte(i)
	}
	return base64.RawURLEncoding.EncodeToString(value)
}

func mapLookup(values map[string]string) func(string) (string, bool) {
	return func(name string) (string, bool) {
		value, ok := values[name]
		return value, ok
	}
}
