package app

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestConnectionBinderBindsValidatedPlayerIdentity(t *testing.T) {
	boundAt := time.Date(2026, 5, 17, 8, 30, 0, 123, time.UTC)
	var received RouteAccessTokenValidationRequest
	binder := ConnectionBinder{
		Validator: routeAccessTokenValidatorFunc(func(_ context.Context, request RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
			received = request
			identity := ValidatedPlayerIdentity("player-1", Session{
				ConnectionID:    "other-connection",
				ConnectionEpoch: 1,
				PlayerID:        "player-1",
			})
			identity.SessionValidated = true
			return RouteAccessTokenValidationResult{
				Valid:    true,
				Identity: identity,
			}, nil
		}),
		Clock: fixedConnectionBindingClock{value: boundAt},
	}

	result, err := binder.BindConnection(context.Background(), ConnectionBindingRequest{
		AccessToken:      "redacted-token-text",
		Route:            BindConnectionRoute(),
		ConnectionID:     " ws-1 ",
		ConnectionEpoch:  7,
		ClientInstanceID: " client-1 ",
	})
	if err != nil {
		t.Fatalf("BindConnection() error = %v, want nil", err)
	}
	if received.Route != BindConnectionRoute() || received.ConnectionID != "ws-1" || received.ConnectionEpoch != 7 {
		t.Fatalf("validator request = %#v, want bind route and server-observed connection metadata", received)
	}
	if !result.Bound || result.BindingStatus != ConnectionBindingStatusBound {
		t.Fatalf("result = %#v, want bound result", result)
	}
	if result.ConnectionID != "ws-1" || result.ConnectionEpoch != 7 || result.ClientInstanceID != "client-1" {
		t.Fatalf("result connection fields = %#v, want normalized request metadata", result)
	}
	if !result.BoundAt.Equal(boundAt) {
		t.Fatalf("BoundAt = %v, want %v", result.BoundAt, boundAt)
	}
	if result.Identity.Status != IdentityValidationValidated ||
		result.Identity.ActorKind != ActorKindPlayer ||
		result.Identity.PlayerID != "player-1" ||
		!result.Identity.PlayerIDValidated ||
		result.Identity.SessionValidated ||
		result.Identity.ConnectionID != "ws-1" ||
		result.Identity.ConnectionEpoch != 7 {
		t.Fatalf("Identity = %#v, want validated player identity bound to connection with SessionValidated=false", result.Identity)
	}
}

func TestConnectionBinderRejectsMissingMalformedInvalidAndUnavailableToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		validator   RouteAccessTokenValidator
		want        ErrorCode
		wantNoCall  bool
		serviceCode ErrorCode
		serviceErr  error
	}{
		{name: "missing", token: "", want: ErrorCodeConnectionBindingTokenMissing, wantNoCall: true},
		{name: "malformed whitespace", token: " redacted-token-text ", want: ErrorCodeConnectionBindingTokenMalformed, wantNoCall: true},
		{name: "validator malformed", token: "redacted-token-text", serviceCode: ErrorCodeAuthenticationTokenMalformed, serviceErr: errors.New("redacted malformed token"), want: ErrorCodeConnectionBindingTokenMalformed},
		{name: "validator invalid", token: "redacted-token-text", serviceCode: ErrorCodeAuthenticationTokenInvalid, serviceErr: errors.New("redacted invalid token"), want: ErrorCodeConnectionBindingTokenInvalid},
		{name: "validator unavailable", token: "redacted-token-text", serviceCode: ErrorCodeAuthenticationTokenUnavailable, serviceErr: errors.New("redacted store unavailable"), want: ErrorCodeConnectionBindingUnavailable},
		{name: "missing validator", token: "redacted-token-text", want: ErrorCodeConnectionBindingUnavailable, wantNoCall: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			validator := tc.validator
			if validator == nil && !tc.wantNoCall {
				validator = routeAccessTokenValidatorFunc(func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
					called = true
					return RouteAccessTokenValidationResult{PublicErrorCode: tc.serviceCode}, tc.serviceErr
				})
			}
			binder := NewConnectionBinder(validator)

			result, err := binder.BindConnection(context.Background(), ConnectionBindingRequest{
				AccessToken:      tc.token,
				Route:            BindConnectionRoute(),
				ConnectionID:     "connection-1",
				ConnectionEpoch:  9,
				ClientInstanceID: "client-1",
			})
			assertApplicationErrorCode(t, err, tc.want)
			if result.Bound || result.BindingStatus != ConnectionBindingStatusRejected {
				t.Fatalf("result = %#v, want rejected unbound result", result)
			}
			if result.Identity.Status != IdentityValidationUnknown {
				t.Fatalf("result Identity.Status = %q, want unknown", result.Identity.Status)
			}
			if tc.wantNoCall && called {
				t.Fatal("validator was called, want local rejection before validation")
			}
			assertConnectionBindingNoSecretLeak(t, err, "redacted-token-text")
		})
	}
}

func TestConnectionBinderRejectsMetadataOnlyIdentity(t *testing.T) {
	binder := NewConnectionBinder(routeAccessTokenValidatorFunc(func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
		return RouteAccessTokenValidationResult{
			Valid:    true,
			Identity: MetadataOnlyIdentityFromSession(Session{PlayerID: "player-1"}),
		}, nil
	}))

	result, err := binder.BindConnection(context.Background(), ConnectionBindingRequest{
		AccessToken:  "redacted-token-text",
		Route:        BindConnectionRoute(),
		ConnectionID: "connection-1",
	})
	assertApplicationErrorCode(t, err, ErrorCodeConnectionBindingTokenInvalid)
	if result.Bound {
		t.Fatalf("result = %#v, want failed bind to stay unbound", result)
	}
}

type fixedConnectionBindingClock struct {
	value time.Time
}

func (c fixedConnectionBindingClock) Now() time.Time {
	return c.value
}

func assertConnectionBindingNoSecretLeak(t *testing.T, err error, secrets ...string) {
	t.Helper()

	if err == nil {
		return
	}
	text := err.Error()
	for _, secret := range secrets {
		if secret != "" && strings.Contains(text, secret) {
			t.Fatalf("error %q leaks secret %q", text, secret)
		}
	}
}
