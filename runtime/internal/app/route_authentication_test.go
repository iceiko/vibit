package app

import (
	"context"
	"errors"
	"testing"
)

func TestRouteProtectorAllowsExplicitPublicAuthenticationRoute(t *testing.T) {
	route := AuthenticateWithDeviceCredentialRoute()
	protector := NewRouteProtector(nil)

	result, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{ConnectionID: "connection-1", PlayerID: "metadata-player"},
		},
	})
	if err != nil {
		t.Fatalf("ProtectRoute() error = %v, want nil", err)
	}
	if !result.Allowed || !result.Public {
		t.Fatalf("ProtectRoute() result = %#v, want allowed public route", result)
	}
	if result.Identity.Status != IdentityValidationMetadataOnly {
		t.Fatalf("Identity.Status = %q, want metadata_only", result.Identity.Status)
	}
}

func TestRouteProtectorAllowsServiceValidatedLogoutLifecycleRoute(t *testing.T) {
	route := LogoutAccessTokenRoute()
	protector := NewRouteProtector(routeAccessTokenValidatorFunc(func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
		t.Fatal("validator was called for service-validated logout route")
		return RouteAccessTokenValidationResult{}, nil
	}))

	result, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{ConnectionID: "connection-1", PlayerID: "metadata-player"},
		},
	})
	if err != nil {
		t.Fatalf("ProtectRoute() error = %v, want nil", err)
	}
	if !result.Allowed || !result.Public {
		t.Fatalf("ProtectRoute() result = %#v, want allowed service-validated lifecycle route", result)
	}
	if result.Identity.Status != IdentityValidationMetadataOnly {
		t.Fatalf("Identity.Status = %q, want metadata_only", result.Identity.Status)
	}
}

func TestRouteProtectorRejectsProtectedRouteWithoutWrapperOrToken(t *testing.T) {
	route := RouteKey{Kind: MessageKindQuery, Module: "inventory", Name: "GetInventory"}
	protector := NewRouteProtector(routeAccessTokenValidatorFunc(func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
		t.Fatal("validator was called without proof")
		return RouteAccessTokenValidationResult{}, nil
	}))

	result, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:    route,
			Identity: MetadataOnlyIdentityFromSession(Session{PlayerID: "player-1"}),
		},
	})
	assertApplicationErrorCode(t, err, ErrorCodeAuthenticationTokenMissing)
	if result.Allowed {
		t.Fatalf("ProtectRoute() result = %#v, want denied", result)
	}
}

func TestRouteProtectionPolicyClassifiesDefaultAndExplicitRoutes(t *testing.T) {
	inventoryRoute := RouteKey{Kind: MessageKindQuery, Module: "inventory", Name: "GetInventory"}
	boundRoute := RouteKey{Kind: MessageKindCommand, Module: "runtime.party", Name: "JoinBoundParty"}
	sessionRoute := RouteKey{Kind: MessageKindCommand, Module: "runtime.party", Name: "ClaimSessionReward"}

	policy := RouteProtectionPolicy{
		RouteRequirements: []RouteProtectionRouteRequirement{
			{Route: boundRoute, Requirement: RouteProtectionBoundConnectionRequired},
			{Route: sessionRoute, Requirement: RouteProtectionSessionValidatedRequired},
		},
	}

	if got := policy.RequirementFor(AuthenticateWithDeviceCredentialRoute()); got != RouteProtectionPublic {
		t.Fatalf("RequirementFor(authenticate) = %q, want public", got)
	}
	if got := policy.RequirementFor(LogoutAccessTokenRoute()); got != RouteProtectionPublic {
		t.Fatalf("RequirementFor(logout) = %q, want public service-validated route", got)
	}
	if got := policy.RequirementFor(inventoryRoute); got != RouteProtectionRequestTokenRequired {
		t.Fatalf("RequirementFor(inventory) = %q, want request_token_required", got)
	}
	if got := policy.RequirementFor(boundRoute); got != RouteProtectionBoundConnectionRequired {
		t.Fatalf("RequirementFor(bound) = %q, want bound_connection_required", got)
	}
	if got := policy.RequirementFor(sessionRoute); got != RouteProtectionSessionValidatedRequired {
		t.Fatalf("RequirementFor(session) = %q, want session_validated_required", got)
	}
}

func TestRouteProtectorMapsMalformedInvalidAndUnavailableFailures(t *testing.T) {
	route := RouteKey{Kind: MessageKindQuery, Module: "inventory", Name: "GetInventory"}
	tests := []struct {
		name string
		code ErrorCode
		err  error
		want ErrorCode
	}{
		{name: "malformed", code: ErrorCodeAuthenticationTokenMalformed, err: errors.New("redacted validation failure"), want: ErrorCodeAuthenticationTokenMalformed},
		{name: "invalid", code: ErrorCodeAuthenticationTokenInvalid, err: errors.New("redacted validation failure"), want: ErrorCodeAuthenticationTokenInvalid},
		{name: "unavailable", code: ErrorCodeAuthenticationTokenUnavailable, err: errors.New("redacted validation failure"), want: ErrorCodeAuthenticationTokenUnavailable},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			protector := NewRouteProtector(routeAccessTokenValidatorFunc(func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
				return RouteAccessTokenValidationResult{PublicErrorCode: tc.code}, tc.err
			}))

			_, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
				Request:             RouteRequest{Route: route},
				AccessToken:         "redacted-token-text",
				ProofCarrierPresent: true,
			})
			assertApplicationErrorCode(t, err, tc.want)
			assertRouteAuthNoSecretLeak(t, err, "redacted-token-text")
		})
	}
}

func TestRouteProtectorRejectsMetadataOnlyIdentityForProtectedRoute(t *testing.T) {
	route := RouteKey{Kind: MessageKindQuery, Module: "inventory", Name: "GetInventory"}
	protector := NewRouteProtector(routeAccessTokenValidatorFunc(func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
		return RouteAccessTokenValidationResult{
			Valid:    true,
			Identity: MetadataOnlyIdentityFromSession(Session{PlayerID: "player-1"}),
		}, nil
	}))

	_, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request:             RouteRequest{Route: route},
		AccessToken:         "redacted-token-text",
		ProofCarrierPresent: true,
	})
	assertApplicationErrorCode(t, err, ErrorCodeAuthenticationTokenInvalid)
}

func TestRouteProtectorRejectsMetadataOnlyIdentityForExplicitProtectedPolicyFamilies(t *testing.T) {
	tests := []struct {
		name        string
		requirement RouteProtectionRequirement
		request     RouteProtectionRequest
		want        ErrorCode
	}{
		{
			name:        "bound connection",
			requirement: RouteProtectionBoundConnectionRequired,
			request: RouteProtectionRequest{
				BoundIdentity: MetadataOnlyIdentityFromSession(Session{ConnectionID: "connection-1", PlayerID: "player-1"}),
			},
			want: ErrorCodeConnectionBindingRequired,
		},
		{
			name:        "session validated",
			requirement: RouteProtectionSessionValidatedRequired,
			request: RouteProtectionRequest{
				SessionValidatedIdentity: MetadataOnlyIdentityFromSession(Session{SessionID: "session-1", PlayerID: "player-1"}),
			},
			want: ErrorCodeSessionInvalid,
		},
		{
			name:        "bound session",
			requirement: RouteProtectionBoundSessionRequired,
			request: RouteProtectionRequest{
				BoundIdentity:            MetadataOnlyIdentityFromSession(Session{ConnectionID: "connection-1", PlayerID: "player-1"}),
				SessionValidatedIdentity: MetadataOnlyIdentityFromSession(Session{SessionID: "session-1", PlayerID: "player-1"}),
			},
			want: ErrorCodeConnectionBindingRequired,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := RouteKey{Kind: MessageKindCommand, Module: "runtime.test", Name: tc.name}
			protector := RouteProtector{
				Policy: RouteProtectionPolicy{
					RouteRequirements: []RouteProtectionRouteRequirement{
						{Route: route, Requirement: tc.requirement},
					},
				},
				Validator: routeAccessTokenValidatorFunc(func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
					t.Fatal("validator was called for explicit identity policy")
					return RouteAccessTokenValidationResult{}, nil
				}),
			}
			tc.request.Request.Route = route
			tc.request.Request.Session = Session{ConnectionID: "connection-1", SessionID: "session-1", PlayerID: "player-1", ConnectionEpoch: 3}

			result, err := protector.ProtectRoute(context.Background(), tc.request)
			assertApplicationErrorCode(t, err, tc.want)
			if result.Allowed {
				t.Fatalf("ProtectRoute() result = %#v, want denied", result)
			}
		})
	}
}

func TestRouteProtectorAllowsExplicitBoundConnectionRouteWithBoundIdentityOnly(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "runtime.party", Name: "JoinBoundParty"}
	protector := RouteProtector{
		Policy: RouteProtectionPolicy{
			RouteRequirements: []RouteProtectionRouteRequirement{
				{Route: route, Requirement: RouteProtectionBoundConnectionRequired},
			},
		},
		Validator: routeAccessTokenValidatorFunc(func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
			t.Fatal("validator was called for bound connection policy")
			return RouteAccessTokenValidationResult{}, nil
		}),
	}
	boundIdentity := ValidatedPlayerIdentity("player-1", Session{
		ConnectionID:    "connection-1",
		PlayerID:        "player-1",
		ConnectionEpoch: 7,
	})
	boundIdentity.SessionValidated = true

	result, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{ConnectionID: "connection-1", PlayerID: "player-1", ConnectionEpoch: 7},
		},
		BoundIdentity: boundIdentity,
	})
	if err != nil {
		t.Fatalf("ProtectRoute() error = %v, want nil", err)
	}
	if !result.Allowed || result.Public {
		t.Fatalf("ProtectRoute() result = %#v, want non-public allowed route", result)
	}
	if result.Identity.PlayerID != "player-1" || result.Identity.ConnectionID != "connection-1" || result.Identity.ConnectionEpoch != 7 || result.Identity.SessionValidated {
		t.Fatalf("Identity = %#v, want bound identity with SessionValidated=false", result.Identity)
	}
}

func TestRouteProtectorRejectsBoundConnectionRouteWithoutMatchingConnection(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "runtime.party", Name: "JoinBoundParty"}
	protector := RouteProtector{
		Policy: RouteProtectionPolicy{
			RouteRequirements: []RouteProtectionRouteRequirement{
				{Route: route, Requirement: RouteProtectionBoundConnectionRequired},
			},
		},
	}
	boundIdentity := ValidatedPlayerIdentity("player-1", Session{
		ConnectionID:    "connection-2",
		PlayerID:        "player-1",
		ConnectionEpoch: 7,
	})
	boundIdentity.SessionValidated = false

	_, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{ConnectionID: "connection-1", PlayerID: "player-1", ConnectionEpoch: 7},
		},
		BoundIdentity: boundIdentity,
	})
	assertApplicationErrorCode(t, err, ErrorCodeConnectionBindingRequired)
}

func TestRouteProtectorAllowsExplicitSessionValidatedRouteWithValidatedSessionIdentity(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "runtime.rewards", Name: "ClaimSessionReward"}
	protector := RouteProtector{
		Policy: RouteProtectionPolicy{
			RouteRequirements: []RouteProtectionRouteRequirement{
				{Route: route, Requirement: RouteProtectionSessionValidatedRequired},
			},
		},
	}
	sessionIdentity := ValidatedPlayerIdentity("player-1", Session{
		SessionID:       "session-1",
		ConnectionID:    "connection-1",
		PlayerID:        "player-1",
		ConnectionEpoch: 9,
	})
	sessionIdentity.SessionValidated = true

	result, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{SessionID: "session-1", ConnectionID: "connection-1", PlayerID: "player-1", ConnectionEpoch: 9},
		},
		SessionValidatedIdentity: sessionIdentity,
	})
	if err != nil {
		t.Fatalf("ProtectRoute() error = %v, want nil", err)
	}
	if !result.Allowed || !result.Identity.SessionValidated || result.Identity.SessionID != "session-1" {
		t.Fatalf("ProtectRoute() result = %#v, want session-validated identity", result)
	}
}

func TestRouteProtectorRejectsSessionValidatedRouteWithoutValidatedSession(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "runtime.rewards", Name: "ClaimSessionReward"}
	protector := RouteProtector{
		Policy: RouteProtectionPolicy{
			RouteRequirements: []RouteProtectionRouteRequirement{
				{Route: route, Requirement: RouteProtectionSessionValidatedRequired},
			},
		},
	}
	sessionIdentity := ValidatedPlayerIdentity("player-1", Session{
		SessionID: "session-1",
		PlayerID:  "player-1",
	})
	sessionIdentity.SessionValidated = false

	_, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{SessionID: "session-1", PlayerID: "player-1"},
		},
		SessionValidatedIdentity: sessionIdentity,
	})
	assertApplicationErrorCode(t, err, ErrorCodeSessionInvalid)
}

func TestRouteProtectorAllowsExplicitBoundSessionRouteWhenIdentitySourcesAgree(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "runtime.match", Name: "EnterBoundSession"}
	protector := RouteProtector{
		Policy: RouteProtectionPolicy{
			RouteRequirements: []RouteProtectionRouteRequirement{
				{Route: route, Requirement: RouteProtectionBoundSessionRequired},
			},
		},
	}
	boundIdentity := ValidatedPlayerIdentity("player-1", Session{
		ConnectionID:    "connection-1",
		PlayerID:        "player-1",
		ConnectionEpoch: 4,
	})
	boundIdentity.SessionValidated = false
	sessionIdentity := ValidatedPlayerIdentity("player-1", Session{
		SessionID:       "session-1",
		ConnectionID:    "connection-1",
		PlayerID:        "player-1",
		ConnectionEpoch: 4,
	})
	sessionIdentity.SessionValidated = true

	result, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{SessionID: "session-1", ConnectionID: "connection-1", PlayerID: "player-1", ConnectionEpoch: 4},
		},
		BoundIdentity:            boundIdentity,
		SessionValidatedIdentity: sessionIdentity,
	})
	if err != nil {
		t.Fatalf("ProtectRoute() error = %v, want nil", err)
	}
	if !result.Allowed || !result.Identity.SessionValidated || result.Identity.PlayerID != "player-1" || result.Identity.SessionID != "session-1" || result.Identity.ConnectionID != "connection-1" {
		t.Fatalf("ProtectRoute() result = %#v, want agreed bound session identity", result)
	}
}

func TestRouteProtectorRejectsExplicitBoundSessionRouteWhenIdentitySourcesMismatch(t *testing.T) {
	route := RouteKey{Kind: MessageKindCommand, Module: "runtime.match", Name: "EnterBoundSession"}
	protector := RouteProtector{
		Policy: RouteProtectionPolicy{
			RouteRequirements: []RouteProtectionRouteRequirement{
				{Route: route, Requirement: RouteProtectionBoundSessionRequired},
			},
		},
	}
	boundIdentity := ValidatedPlayerIdentity("player-1", Session{
		ConnectionID:    "connection-1",
		PlayerID:        "player-1",
		ConnectionEpoch: 4,
	})
	boundIdentity.SessionValidated = false
	sessionIdentity := ValidatedPlayerIdentity("player-2", Session{
		SessionID:       "session-1",
		ConnectionID:    "connection-1",
		PlayerID:        "player-2",
		ConnectionEpoch: 4,
	})
	sessionIdentity.SessionValidated = true

	_, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{SessionID: "session-1", ConnectionID: "connection-1", PlayerID: "player-1", ConnectionEpoch: 4},
		},
		BoundIdentity:            boundIdentity,
		SessionValidatedIdentity: sessionIdentity,
	})
	assertApplicationErrorCode(t, err, ErrorCodeConnectionBindingRequired)
}

func TestRouteProtectorPassesValidatedPlayerIdentityWithSessionUnvalidated(t *testing.T) {
	route := RouteKey{Kind: MessageKindQuery, Module: "inventory", Name: "GetInventory"}
	var received RouteAccessTokenValidationRequest
	validatedIdentity := ValidatedPlayerIdentity("player-1", Session{
		ConnectionID:    "connection-1",
		PlayerID:        "metadata-player",
		ConnectionEpoch: 9,
	})
	validatedIdentity.SessionValidated = true
	protector := NewRouteProtector(routeAccessTokenValidatorFunc(func(_ context.Context, request RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
		received = request
		return RouteAccessTokenValidationResult{
			Valid:    true,
			Identity: validatedIdentity,
		}, nil
	}))

	result, err := protector.ProtectRoute(context.Background(), RouteProtectionRequest{
		Request: RouteRequest{
			Route:   route,
			Session: Session{ConnectionID: "connection-1", PlayerID: "metadata-player", ConnectionEpoch: 9},
		},
		AccessToken:         "redacted-token-text",
		ProofCarrierPresent: true,
	})
	if err != nil {
		t.Fatalf("ProtectRoute() error = %v, want nil", err)
	}
	if !result.Allowed {
		t.Fatalf("ProtectRoute() result = %#v, want allowed", result)
	}
	if received.Route != route || received.ConnectionID != "connection-1" || received.ConnectionEpoch != 9 {
		t.Fatalf("validator request = %#v, want route and connection metadata", received)
	}
	if result.Identity.Status != IdentityValidationValidated ||
		result.Identity.ActorKind != ActorKindPlayer ||
		result.Identity.PlayerID != "player-1" ||
		!result.Identity.PlayerIDValidated ||
		result.Identity.SessionValidated {
		t.Fatalf("Identity = %#v, want validated player identity with SessionValidated=false", result.Identity)
	}
}

type routeAccessTokenValidatorFunc func(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error)

func (f routeAccessTokenValidatorFunc) ValidateRouteAccessToken(ctx context.Context, request RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error) {
	if f == nil {
		return RouteAccessTokenValidationResult{}, errors.New("route token validator function is nil")
	}
	return f(ctx, request)
}

func assertRouteAuthNoSecretLeak(t *testing.T, err error, secrets ...string) {
	t.Helper()

	if err == nil {
		return
	}
	text := err.Error()
	for _, secret := range secrets {
		if secret != "" && contains(text, secret) {
			t.Fatalf("error %q leaks secret %q", text, secret)
		}
	}
}

func contains(value string, needle string) bool {
	for i := 0; i+len(needle) <= len(value); i++ {
		if value[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
