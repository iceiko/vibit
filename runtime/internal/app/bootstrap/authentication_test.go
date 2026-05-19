package bootstrap

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
)

func TestAuthenticationRouteHandlersRegisterPublicLoginRoute(t *testing.T) {
	dispatcher := app.NewDispatcher()
	service := &recordingDeviceCredentialAuthenticationService{}
	handlers := AuthenticationRouteHandlers{Service: service}

	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("RegisterRoutes() error = %v, want nil", err)
	}

	result, err := dispatcher.Dispatch(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     app.AuthenticateWithDeviceCredentialRoute(),
		Payload: appauth.DeviceCredentialAuthenticationRequest{
			CredentialProof:       "redacted-credential-proof",
			RequestedPlayerID:     "player-1",
			ClientInstanceID:      "client-1",
			AccountCreationIntent: appauth.AccountCreationIntentReject,
		},
	})
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}
	if !service.called {
		t.Fatal("service was not called")
	}
	if service.request.CredentialProof != "redacted-credential-proof" ||
		service.request.RequestedPlayerID != "player-1" ||
		service.request.ClientInstanceID != "client-1" ||
		service.request.AccountCreationIntent != appauth.AccountCreationIntentReject {
		t.Fatalf("service request = %#v, want mapped login payload", service.request)
	}
	payload, ok := result.Payload.(appauth.AuthenticationResult)
	if !ok {
		t.Fatalf("result Payload = %T, want authentication.AuthenticationResult", result.Payload)
	}
	if payload.AccessToken != "redacted-access-token" || payload.TokenRecordID != "token-record-1" {
		t.Fatalf("result payload = %#v, want service result", payload)
	}
}

func TestAuthenticationRouteHandlersRegisterLogoutRoute(t *testing.T) {
	dispatcher := app.NewDispatcher()
	service := &recordingDeviceCredentialAuthenticationService{}
	handlers := AuthenticationRouteHandlers{Service: service}

	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("RegisterRoutes() error = %v, want nil", err)
	}

	result, err := dispatcher.Dispatch(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     app.LogoutAccessTokenRoute(),
		Payload: appauth.LogoutAccessTokenRequest{
			AccessToken:  "redacted-access-token",
			LogoutReason: "player_logout",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}
	if !service.logoutCalled {
		t.Fatal("logout service was not called")
	}
	if service.logoutRequest.AccessToken != "redacted-access-token" || service.logoutRequest.LogoutReason != "player_logout" {
		t.Fatalf("logout request = %#v, want mapped logout payload", service.logoutRequest)
	}
	payload, ok := result.Payload.(appauth.LogoutAccessTokenResult)
	if !ok {
		t.Fatalf("result Payload = %T, want authentication.LogoutAccessTokenResult", result.Payload)
	}
	if !payload.Revoked || payload.TokenRecordID != "token-record-1" || payload.LogoutScope != appauth.LogoutScopeToken {
		t.Fatalf("result payload = %#v, want service logout result", payload)
	}
}

func TestAuthenticationRouteHandlersMapServiceErrorsWithoutSecretLeak(t *testing.T) {
	secret := "redacted-credential-proof"
	service := &recordingDeviceCredentialAuthenticationService{
		result: appauth.AuthenticationResult{
			Status:          appauth.AuthenticationStatusRejected,
			PublicErrorCode: appauth.PublicErrorAuthenticationCredentialInvalid,
		},
		err: &appauth.ServiceError{
			Operation:  appauth.OperationAuthenticateWithDeviceCredential,
			Class:      appauth.FailureClassLookupMiss,
			PublicCode: appauth.PublicErrorAuthenticationCredentialInvalid,
			Err:        errors.New("internal lookup miss"),
		},
	}
	handlers := AuthenticationRouteHandlers{Service: service}

	result, err := handlers.HandleAuthenticateWithDeviceCredentialRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     app.AuthenticateWithDeviceCredentialRoute(),
		Payload:   appauth.DeviceCredentialAuthenticationRequest{CredentialProof: secret},
	})
	if err == nil {
		t.Fatal("HandleAuthenticateWithDeviceCredentialRoute() error = nil, want public error")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appauth.PublicErrorAuthenticationCredentialInvalid) {
		t.Fatalf("result error = %#v, want invalid credential code", result.Error)
	}
	errorText := result.Error.Error()
	if strings.Contains(errorText, secret) || strings.Contains(errorText, "lookup miss") {
		t.Fatalf("error text %q leaked secret or internal detail", errorText)
	}
	if result.Payload != nil {
		t.Fatalf("result Payload = %#v, want nil on failed login", result.Payload)
	}
}

func TestAuthenticationRouteHandlersMapLogoutServiceErrorsWithoutSecretLeak(t *testing.T) {
	secret := "redacted-access-token"
	service := &recordingDeviceCredentialAuthenticationService{
		logoutResult: appauth.LogoutAccessTokenResult{
			Status:          appauth.LogoutStatusRejected,
			PublicErrorCode: appauth.PublicErrorAuthenticationTokenInvalid,
		},
		logoutErr: &appauth.ServiceError{
			Operation:  appauth.OperationLogoutAccessToken,
			Class:      appauth.FailureClassVerifierMismatch,
			PublicCode: appauth.PublicErrorAuthenticationTokenInvalid,
			Err:        errors.New("internal verifier mismatch"),
		},
	}
	handlers := AuthenticationRouteHandlers{Service: service}

	result, err := handlers.HandleLogoutAccessTokenRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     app.LogoutAccessTokenRoute(),
		Payload:   appauth.LogoutAccessTokenRequest{AccessToken: secret},
	})
	if err == nil {
		t.Fatal("HandleLogoutAccessTokenRoute() error = nil, want public error")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appauth.PublicErrorAuthenticationTokenInvalid) {
		t.Fatalf("result error = %#v, want invalid token code", result.Error)
	}
	errorText := result.Error.Error()
	if strings.Contains(errorText, secret) || strings.Contains(errorText, "verifier mismatch") {
		t.Fatalf("error text %q leaked secret or internal detail", errorText)
	}
	if result.Payload != nil {
		t.Fatalf("result Payload = %#v, want nil on failed logout", result.Payload)
	}
}

func TestAuthenticationRouteHandlersRejectMalformedPayload(t *testing.T) {
	service := &recordingDeviceCredentialAuthenticationService{}
	handlers := AuthenticationRouteHandlers{Service: service}

	result, err := handlers.HandleAuthenticateWithDeviceCredentialRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     app.AuthenticateWithDeviceCredentialRoute(),
		Payload:   "not a login request",
	})
	if err == nil {
		t.Fatal("HandleAuthenticateWithDeviceCredentialRoute() error = nil, want malformed proof error")
	}
	if service.called {
		t.Fatal("service was called for malformed payload")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appauth.PublicErrorAuthenticationProofMalformed) {
		t.Fatalf("result error = %#v, want malformed proof", result.Error)
	}
}

func TestAuthenticationRouteHandlersRejectMalformedLogoutPayload(t *testing.T) {
	service := &recordingDeviceCredentialAuthenticationService{}
	handlers := AuthenticationRouteHandlers{Service: service}

	result, err := handlers.HandleLogoutAccessTokenRoute(context.Background(), app.RouteRequest{
		RequestID: "request-1",
		Route:     app.LogoutAccessTokenRoute(),
		Payload:   "not a logout request",
	})
	if err == nil {
		t.Fatal("HandleLogoutAccessTokenRoute() error = nil, want malformed token error")
	}
	if service.logoutCalled {
		t.Fatal("service was called for malformed logout payload")
	}
	if result.Error == nil || result.Error.Code != app.ErrorCode(appauth.PublicErrorAuthenticationTokenMalformed) {
		t.Fatalf("result error = %#v, want malformed token", result.Error)
	}
}

type recordingDeviceCredentialAuthenticationService struct {
	called        bool
	request       appauth.DeviceCredentialAuthenticationRequest
	result        appauth.AuthenticationResult
	err           error
	logoutCalled  bool
	logoutRequest appauth.LogoutAccessTokenRequest
	logoutResult  appauth.LogoutAccessTokenResult
	logoutErr     error
}

func (s *recordingDeviceCredentialAuthenticationService) AuthenticateWithDeviceCredential(_ context.Context, request appauth.DeviceCredentialAuthenticationRequest) (appauth.AuthenticationResult, error) {
	s.called = true
	s.request = request
	if s.err != nil {
		return s.result, s.err
	}
	if s.result.Status != "" {
		return s.result, nil
	}
	return appauth.AuthenticationResult{
		Status:           appauth.AuthenticationStatusAuthenticated,
		ActorKind:        app.ActorKindPlayer,
		PlayerID:         "player-1",
		AccessToken:      "redacted-access-token",
		TokenType:        appauth.TokenTypeOpaqueAccess,
		IssuedAt:         time.Date(2026, 5, 17, 1, 2, 3, 0, time.UTC),
		ExpiresAt:        time.Date(2026, 5, 17, 2, 2, 3, 0, time.UTC),
		SessionID:        "runtime-session-1",
		SessionCreated:   true,
		SessionExpiresAt: time.Date(2026, 5, 17, 2, 2, 3, 0, time.UTC),
		TokenRecordID:    "token-record-1",
	}, nil
}

func (s *recordingDeviceCredentialAuthenticationService) LogoutAccessToken(_ context.Context, request appauth.LogoutAccessTokenRequest) (appauth.LogoutAccessTokenResult, error) {
	s.logoutCalled = true
	s.logoutRequest = request
	if s.logoutErr != nil {
		return s.logoutResult, s.logoutErr
	}
	if s.logoutResult.Status != "" {
		return s.logoutResult, nil
	}
	return appauth.LogoutAccessTokenResult{
		Status:        appauth.LogoutStatusRevoked,
		Revoked:       true,
		LogoutScope:   appauth.LogoutScopeToken,
		TokenRecordID: "token-record-1",
		RevokedAt:     time.Date(2026, 5, 18, 1, 2, 3, 0, time.UTC),
	}, nil
}
