package protobuf

import (
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
	authenticationv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/authentication/v1"
)

func TestRouteRequestWithAuthenticationPayloadMapsLoginRequest(t *testing.T) {
	request := app.RouteRequest{
		Route: app.AuthenticateWithDeviceCredentialRoute(),
		Payload: &authenticationv1.AuthenticateWithDeviceCredentialRequest{
			CredentialProof:       "redacted-credential-proof",
			RequestedPlayerId:     "player-1",
			ClientInstanceId:      "client-1",
			AccountCreationIntent: authenticationv1.AccountCreationIntent_ACCOUNT_CREATION_INTENT_ALLOW_CREATE,
		},
	}

	mapped, err := RouteRequestWithDomainPayload(request)
	if err != nil {
		t.Fatalf("RouteRequestWithDomainPayload() error = %v, want nil", err)
	}
	payload, ok := mapped.Payload.(appauth.DeviceCredentialAuthenticationRequest)
	if !ok {
		t.Fatalf("mapped Payload = %T, want DeviceCredentialAuthenticationRequest", mapped.Payload)
	}
	if payload.CredentialProof != "redacted-credential-proof" ||
		payload.RequestedPlayerID != "player-1" ||
		payload.ClientInstanceID != "client-1" ||
		payload.AccountCreationIntent != appauth.AccountCreationIntentCreate {
		t.Fatalf("payload = %#v, want mapped authentication request", payload)
	}
}

func TestRouteRequestWithAuthenticationPayloadMapsLogoutRequest(t *testing.T) {
	request := app.RouteRequest{
		Route: app.LogoutAccessTokenRoute(),
		Payload: &authenticationv1.LogoutAccessTokenRequest{
			AccessToken:  "redacted-access-token",
			LogoutReason: "player_logout",
		},
	}

	mapped, err := RouteRequestWithDomainPayload(request)
	if err != nil {
		t.Fatalf("RouteRequestWithDomainPayload() error = %v, want nil", err)
	}
	payload, ok := mapped.Payload.(appauth.LogoutAccessTokenRequest)
	if !ok {
		t.Fatalf("mapped Payload = %T, want LogoutAccessTokenRequest", mapped.Payload)
	}
	if payload.AccessToken != "redacted-access-token" || payload.LogoutReason != "player_logout" {
		t.Fatalf("payload = %#v, want mapped logout request", payload)
	}
}

func TestProtoPayloadFromAuthenticationResultMapsLoginResponse(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 1, 2, 3, 4, time.UTC)
	expiresAt := issuedAt.Add(time.Hour)
	payload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: app.AuthenticateWithDeviceCredentialRoute(),
		Payload: appauth.AuthenticationResult{
			Status:           appauth.AuthenticationStatusAuthenticated,
			ActorKind:        app.ActorKindPlayer,
			PlayerID:         "player-1",
			AccessToken:      "redacted-access-token",
			TokenType:        appauth.TokenTypeOpaqueAccess,
			IssuedAt:         issuedAt,
			ExpiresAt:        expiresAt,
			SessionID:        "runtime-session-1",
			SessionCreated:   true,
			SessionExpiresAt: expiresAt,
			TokenRecordID:    "token-record-1",
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult() error = %v, want nil", err)
	}
	response, ok := payload.(*authenticationv1.AuthenticateWithDeviceCredentialResponse)
	if !ok {
		t.Fatalf("payload = %T, want AuthenticateWithDeviceCredentialResponse", payload)
	}
	if response.GetAuthenticationStatus() != string(appauth.AuthenticationStatusAuthenticated) ||
		response.GetActorKind() != string(app.ActorKindPlayer) ||
		response.GetPlayerId() != "player-1" ||
		response.GetAccessToken() != "redacted-access-token" ||
		response.GetTokenType() != string(appauth.TokenTypeOpaqueAccess) ||
		response.GetTokenRecordId() != "token-record-1" {
		t.Fatalf("response = %#v, want mapped authentication response", response)
	}
	if response.GetIssuedAt() != issuedAt.UTC().Format(time.RFC3339Nano) ||
		response.GetExpiresAt() != expiresAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("response times = %q %q, want RFC3339Nano UTC", response.GetIssuedAt(), response.GetExpiresAt())
	}
}

func TestProtoPayloadFromAuthenticationResultMapsLogoutResponse(t *testing.T) {
	revokedAt := time.Date(2026, 5, 18, 1, 2, 3, 4, time.FixedZone("test", 8*60*60))
	payload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: app.LogoutAccessTokenRoute(),
		Payload: appauth.LogoutAccessTokenResult{
			Status:        appauth.LogoutStatusRevoked,
			Revoked:       true,
			LogoutScope:   appauth.LogoutScopeToken,
			TokenRecordID: "token-record-1",
			RevokedAt:     revokedAt,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult() error = %v, want nil", err)
	}
	response, ok := payload.(*authenticationv1.LogoutAccessTokenResponse)
	if !ok {
		t.Fatalf("payload = %T, want LogoutAccessTokenResponse", payload)
	}
	if response.GetLogoutStatus() != string(appauth.LogoutStatusRevoked) ||
		!response.GetRevoked() ||
		response.GetLogoutScope() != "presented_access_token" ||
		response.GetTokenRecordId() != "token-record-1" {
		t.Fatalf("response = %#v, want mapped logout response", response)
	}
	if response.GetRevokedAt() != revokedAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("RevokedAt = %q, want RFC3339Nano UTC", response.GetRevokedAt())
	}
}
