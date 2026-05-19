package protobuf

import (
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
	authenticationv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/authentication/v1"
	"google.golang.org/protobuf/proto"
)

func routeRequestWithAuthenticationPayload(request app.RouteRequest) (app.RouteRequest, bool, error) {
	switch request.Route {
	case app.AuthenticateWithDeviceCredentialRoute():
		payload, ok := request.Payload.(*authenticationv1.AuthenticateWithDeviceCredentialRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.authentication.v1.AuthenticateWithDeviceCredentialRequest")
		}

		request.Payload = appauth.DeviceCredentialAuthenticationRequest{
			CredentialProof:       payload.GetCredentialProof(),
			RequestedPlayerID:     payload.GetRequestedPlayerId(),
			ClientInstanceID:      payload.GetClientInstanceId(),
			AccountCreationIntent: accountCreationIntentFromProto(payload.GetAccountCreationIntent()),
		}
		return request, true, nil

	case app.LogoutAccessTokenRoute():
		payload, ok := request.Payload.(*authenticationv1.LogoutAccessTokenRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.authentication.v1.LogoutAccessTokenRequest")
		}

		request.Payload = appauth.LogoutAccessTokenRequest{
			AccessToken:  payload.GetAccessToken(),
			LogoutReason: payload.GetLogoutReason(),
		}
		return request, true, nil

	default:
		return request, false, nil
	}
}

func protoPayloadFromAuthenticationRoute(route app.RouteKey, payload any) (proto.Message, bool, error) {
	switch route {
	case app.AuthenticateWithDeviceCredentialRoute():
		result, ok := payload.(appauth.AuthenticationResult)
		if !ok {
			if pointerResult, pointerOK := payload.(*appauth.AuthenticationResult); pointerOK && pointerResult != nil {
				result = *pointerResult
				ok = true
			}
		}
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be authentication.AuthenticationResult")
		}

		return &authenticationv1.AuthenticateWithDeviceCredentialResponse{
			AuthenticationStatus: string(result.Status),
			ActorKind:            string(result.ActorKind),
			PlayerId:             result.PlayerID,
			AccessToken:          result.AccessToken,
			TokenType:            string(result.TokenType),
			IssuedAt:             formatAuthenticationTime(result.IssuedAt),
			ExpiresAt:            formatAuthenticationTime(result.ExpiresAt),
			TokenRecordId:        result.TokenRecordID,
		}, true, nil

	case app.LogoutAccessTokenRoute():
		result, ok := payload.(appauth.LogoutAccessTokenResult)
		if !ok {
			if pointerResult, pointerOK := payload.(*appauth.LogoutAccessTokenResult); pointerOK && pointerResult != nil {
				result = *pointerResult
				ok = true
			}
		}
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be authentication.LogoutAccessTokenResult")
		}

		return &authenticationv1.LogoutAccessTokenResponse{
			LogoutStatus:  string(result.Status),
			Revoked:       result.Revoked,
			LogoutScope:   protocolLogoutScope(result.LogoutScope),
			RevokedAt:     formatAuthenticationTime(result.RevokedAt),
			TokenRecordId: result.TokenRecordID,
		}, true, nil

	default:
		return nil, false, nil
	}
}

func sessionCarrierFromAuthenticationPayload(route app.RouteKey, session app.Session, payload any) app.Session {
	if route != app.AuthenticateWithDeviceCredentialRoute() {
		return session
	}

	result, ok := payload.(appauth.AuthenticationResult)
	if !ok {
		if pointerResult, pointerOK := payload.(*appauth.AuthenticationResult); pointerOK && pointerResult != nil {
			result = *pointerResult
			ok = true
		}
	}
	if !ok || result.Status != appauth.AuthenticationStatusAuthenticated {
		return session
	}

	if sessionID := strings.TrimSpace(result.SessionID); sessionID != "" {
		session.SessionID = sessionID
	}
	if playerID := strings.TrimSpace(result.PlayerID); playerID != "" {
		session.PlayerID = playerID
	}
	return session
}

func accountCreationIntentFromProto(value authenticationv1.AccountCreationIntent) appauth.AccountCreationIntent {
	switch value {
	case authenticationv1.AccountCreationIntent_ACCOUNT_CREATION_INTENT_ALLOW_CREATE:
		return appauth.AccountCreationIntentCreate
	case authenticationv1.AccountCreationIntent_ACCOUNT_CREATION_INTENT_AUTHENTICATE_EXISTING_ONLY:
		return appauth.AccountCreationIntentReject
	default:
		return appauth.AccountCreationIntentUnspecified
	}
}

func formatAuthenticationTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func protocolLogoutScope(scope appauth.LogoutScope) string {
	switch scope {
	case appauth.LogoutScopeToken:
		return "presented_access_token"
	default:
		return ""
	}
}
