package protobuf

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
	apppresence "github.com/iceiko/vibit/runtime/internal/app/presence"
	appstorage "github.com/iceiko/vibit/runtime/internal/app/storage"
	authenticationv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/authentication/v1"
	inventoryv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/inventory/v1"
	presencev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/presence/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	storagev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/storage/v1"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	storagemodule "github.com/iceiko/vibit/runtime/internal/modules/storage"
	"google.golang.org/protobuf/proto"
)

func TestFrameHandlerRejectsMissingWrapperOnProtectedRouteBeforeDomainDispatch(t *testing.T) {
	var dispatched bool
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
			dispatched = true
			return app.ApplicationResult{}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			t.Fatal("validator was called without wrapper")
			return app.RouteAccessTokenValidationResult{}, nil
		})),
	}
	requestPayload := mustMarshalEnvelope(t, inventory.GetInventoryRoute(), &inventoryv1.GetInventoryRequest{
		PlayerId:    "player-1",
		RequestedBy: "player-1",
	})

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{ConnectionID: "ws-1", Payload: requestPayload})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil error envelope", err)
	}
	if dispatched {
		t.Fatal("domain dispatcher was called, want route protection to stop dispatch")
	}
	envelope := mustUnmarshalSingleResponse(t, responses)
	assertErrorEnvelope(t, envelope, app.ErrorCodeAuthenticationTokenMissing)
}

func TestFrameHandlerRejectsMissingAccessTokenInWrapperBeforeDomainDispatch(t *testing.T) {
	var dispatched bool
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
			dispatched = true
			return app.ApplicationResult{}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			t.Fatal("validator was called with missing token")
			return app.RouteAccessTokenValidationResult{}, nil
		})),
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID: "ws-1",
		Payload: mustMarshalAuthenticatedEnvelope(t, inventory.GetInventoryRoute(), "", &inventoryv1.GetInventoryRequest{
			PlayerId:    "player-1",
			RequestedBy: "player-1",
		}),
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil error envelope", err)
	}
	if dispatched {
		t.Fatal("domain dispatcher was called, want route protection to stop dispatch")
	}
	envelope := mustUnmarshalSingleResponse(t, responses)
	assertErrorEnvelope(t, envelope, app.ErrorCodeAuthenticationTokenMissing)
}

func TestFrameHandlerMapsMalformedWrapperToMalformedTokenError(t *testing.T) {
	var dispatched bool
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
			dispatched = true
			return app.ApplicationResult{}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			t.Fatal("validator was called with malformed wrapper")
			return app.RouteAccessTokenValidationResult{}, nil
		})),
	}

	envelope := mustBuildEnvelope(
		t,
		inventory.GetInventoryRoute(),
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{SessionID: "session-1", PlayerID: "metadata-player"},
		&authenticationv1.AuthenticatedRequest{
			AccessToken:      "redacted-token-text",
			InnerPayloadType: "vibit.inventory.v1.GetInventoryRequest",
		},
	)
	encoded, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal(envelope) error = %v, want nil", err)
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{ConnectionID: "ws-1", Payload: encoded})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil error envelope", err)
	}
	if dispatched {
		t.Fatal("domain dispatcher was called, want malformed wrapper to stop dispatch")
	}
	response := mustUnmarshalSingleResponse(t, responses)
	assertErrorEnvelope(t, response, app.ErrorCodeAuthenticationTokenMalformed)
	assertNoFrameErrorSecretLeak(t, response, "redacted-token-text")
}

func TestFrameHandlerMapsInvalidAndUnavailableTokenFailures(t *testing.T) {
	tests := []struct {
		name string
		code app.ErrorCode
		err  error
	}{
		{name: "invalid", code: app.ErrorCodeAuthenticationTokenInvalid, err: errors.New("redacted invalid token")},
		{name: "unavailable", code: app.ErrorCodeAuthenticationTokenUnavailable, err: errors.New("redacted store unavailable")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var dispatched bool
			handler := FrameHandler{
				Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
					dispatched = true
					return app.ApplicationResult{}, nil
				}),
				RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
					return app.RouteAccessTokenValidationResult{PublicErrorCode: tc.code}, tc.err
				})),
			}
			responses, err := handler.HandleFrame(context.Background(), FrameRequest{
				ConnectionID: "ws-1",
				Payload: mustMarshalAuthenticatedEnvelope(t, inventory.GetInventoryRoute(), "redacted-token-text", &inventoryv1.GetInventoryRequest{
					PlayerId:    "player-1",
					RequestedBy: "player-1",
				}),
			})
			if err != nil {
				t.Fatalf("HandleFrame() error = %v, want nil error envelope", err)
			}
			if dispatched {
				t.Fatal("domain dispatcher was called, want validation failure to stop dispatch")
			}
			response := mustUnmarshalSingleResponse(t, responses)
			assertErrorEnvelope(t, response, tc.code)
			assertNoFrameErrorSecretLeak(t, response, "redacted-token-text")
		})
	}
}

func TestFrameHandlerRejectsMetadataOnlyIdentityFromValidator(t *testing.T) {
	var dispatched bool
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
			dispatched = true
			return app.ApplicationResult{}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			return app.RouteAccessTokenValidationResult{
				Valid:    true,
				Identity: app.MetadataOnlyIdentityFromSession(app.Session{PlayerID: "player-1"}),
			}, nil
		})),
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID: "ws-1",
		Payload: mustMarshalAuthenticatedEnvelope(t, inventory.GetInventoryRoute(), "redacted-token-text", &inventoryv1.GetInventoryRequest{
			PlayerId:    "player-1",
			RequestedBy: "player-1",
		}),
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil error envelope", err)
	}
	if dispatched {
		t.Fatal("domain dispatcher was called, want metadata-only identity to be rejected")
	}
	response := mustUnmarshalSingleResponse(t, responses)
	assertErrorEnvelope(t, response, app.ErrorCodeAuthenticationTokenInvalid)
}

func TestFrameHandlerValidTokenDecodesInnerPayloadAfterValidation(t *testing.T) {
	var validatorCalled bool
	var received app.RouteRequest
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(_ context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
			received = request
			return app.ApplicationResult{
				RequestID: request.RequestID,
				Route:     request.Route,
				Target:    request.Target,
				Session:   request.Session,
				Identity:  request.Identity,
				Payload: inventory.GetInventoryResponse{
					PlayerID: "player-1",
					Items:    []inventory.Item{{ItemID: "item-1", Quantity: 2}},
				},
			}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(_ context.Context, request app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			validatorCalled = true
			if request.Route != inventory.GetInventoryRoute() || request.ConnectionID != "ws-1" || request.ConnectionEpoch != 7 {
				t.Fatalf("validator request = %#v, want route and frame/session metadata", request)
			}
			identity := app.ValidatedPlayerIdentity("player-1", app.Session{
				ConnectionID:    request.ConnectionID,
				ConnectionEpoch: request.ConnectionEpoch,
				PlayerID:        "metadata-player",
			})
			identity.SessionValidated = false
			return app.RouteAccessTokenValidationResult{
				Valid:    true,
				Identity: identity,
			}, nil
		})),
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID: "ws-1",
		Payload: mustMarshalAuthenticatedEnvelope(t, inventory.GetInventoryRoute(), "redacted-token-text", &inventoryv1.GetInventoryRequest{
			PlayerId:    "player-1",
			RequestedBy: "player-1",
		}),
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil", err)
	}
	if !validatorCalled {
		t.Fatal("validator was not called")
	}
	if received.PayloadType != "vibit.inventory.v1.GetInventoryRequest" {
		t.Fatalf("received PayloadType = %q, want inner inventory payload type", received.PayloadType)
	}
	if _, ok := received.Payload.(inventory.GetInventoryRequest); !ok {
		t.Fatalf("received Payload = %T, want inventory.GetInventoryRequest", received.Payload)
	}
	if received.Identity.Status != app.IdentityValidationValidated ||
		received.Identity.ActorKind != app.ActorKindPlayer ||
		received.Identity.PlayerID != "player-1" ||
		!received.Identity.PlayerIDValidated ||
		received.Identity.SessionValidated {
		t.Fatalf("received Identity = %#v, want validated player identity with SessionValidated=false", received.Identity)
	}

	response := mustUnmarshalSingleResponse(t, responses)
	if response.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_QUERY {
		t.Fatalf("response kind = %v, want query", response.GetKind())
	}
	if response.GetPayloadType() != "vibit.inventory.v1.GetInventoryResponse" {
		t.Fatalf("response PayloadType = %q, want GetInventoryResponse", response.GetPayloadType())
	}
}

func TestFrameHandlerProtectedPresenceQueryUsesValidatedSelfIdentity(t *testing.T) {
	var validatorCalled bool
	var received app.RouteRequest
	route := apppresence.GetPlayerPresenceRoute()
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(_ context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
			received = request
			payload, ok := request.Payload.(apppresence.GetPlayerPresenceRequest)
			if !ok {
				t.Fatalf("request Payload = %T, want GetPlayerPresenceRequest", request.Payload)
			}
			if payload.PlayerID != "player-1" {
				t.Fatalf("presence request payload = %#v, want player-1", payload)
			}
			observedAt := requestLoopClock{}.Now()
			return app.ApplicationResult{
				RequestID: request.RequestID,
				Route:     request.Route,
				Target:    request.Target,
				Session:   request.Session,
				Identity:  request.Identity,
				Payload: apppresence.GetPlayerPresenceResult{
					PlayerID:        "player-1",
					Status:          apppresence.PresenceStatusOnline,
					ConnectionCount: 1,
					ObservedAt:      observedAt,
					ActiveConnections: []apppresence.PresenceConnection{{
						ConnectionID:     "ws-1",
						ConnectionEpoch:  7,
						RuntimeSessionID: "session-1",
						OpenedAt:         observedAt,
					}},
					RuntimeSessionIDs: []string{"session-1"},
				},
			}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(_ context.Context, request app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			validatorCalled = true
			if request.Route != route || request.ConnectionID != "ws-1" || request.ConnectionEpoch != 7 {
				t.Fatalf("validator request = %#v, want presence route and frame/session metadata", request)
			}
			identity := app.ValidatedPlayerIdentity("player-1", app.Session{
				ConnectionID:    request.ConnectionID,
				ConnectionEpoch: request.ConnectionEpoch,
				PlayerID:        "metadata-player",
			})
			identity.SessionValidated = false
			return app.RouteAccessTokenValidationResult{
				Valid:    true,
				Identity: identity,
			}, nil
		})),
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID: "ws-1",
		Payload: mustMarshalAuthenticatedEnvelope(t, route, "redacted-token-text", &presencev1.GetPlayerPresenceRequest{
			PlayerId: "player-1",
		}),
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil", err)
	}
	if !validatorCalled {
		t.Fatal("validator was not called")
	}
	if received.PayloadType != "vibit.presence.v1.GetPlayerPresenceRequest" {
		t.Fatalf("received PayloadType = %q, want presence request payload type", received.PayloadType)
	}
	if received.Identity.Status != app.IdentityValidationValidated ||
		received.Identity.PlayerID != "player-1" ||
		!received.Identity.PlayerIDValidated ||
		received.Identity.SessionValidated {
		t.Fatalf("received Identity = %#v, want validated player identity with SessionValidated=false", received.Identity)
	}

	response := mustUnmarshalSingleResponse(t, responses)
	if response.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_QUERY {
		t.Fatalf("response kind = %v, want query", response.GetKind())
	}
	if response.GetPayloadType() != "vibit.presence.v1.GetPlayerPresenceResponse" {
		t.Fatalf("response payload_type = %q, want presence response", response.GetPayloadType())
	}
	responsePayload, err := DecodePayload(response.GetPayloadType(), response.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload(presence response) error = %v, want nil", err)
	}
	presenceResponse, ok := responsePayload.(*presencev1.GetPlayerPresenceResponse)
	if !ok {
		t.Fatalf("response payload = %T, want GetPlayerPresenceResponse", responsePayload)
	}
	if presenceResponse.GetPlayerId() != "player-1" ||
		presenceResponse.GetPresenceStatus() != presencev1.PresenceStatus_PRESENCE_STATUS_ONLINE ||
		presenceResponse.GetConnectionCount() != 1 ||
		len(presenceResponse.GetActiveConnections()) != 1 ||
		presenceResponse.GetActiveConnections()[0].GetConnectionId() != "ws-1" {
		t.Fatalf("presence response = %#v, want online self presence", presenceResponse)
	}
}

func TestFrameHandlerProtectedStorageRouteRequiresAuthenticatedWrapper(t *testing.T) {
	var dispatched bool
	route := appstorage.GetOwnStorageObjectRoute()
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
			dispatched = true
			return app.ApplicationResult{}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			t.Fatal("validator was called without authenticated wrapper")
			return app.RouteAccessTokenValidationResult{}, nil
		})),
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID: "ws-1",
		Payload: mustMarshalEnvelope(t, route, &storagev1.GetOwnStorageObjectRequest{
			Collection: "progress",
			Key:        "tutorial",
		}),
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil error envelope", err)
	}
	if dispatched {
		t.Fatal("domain dispatcher was called, want route protection to stop dispatch")
	}
	response := mustUnmarshalSingleResponse(t, responses)
	assertErrorEnvelope(t, response, app.ErrorCodeAuthenticationTokenMissing)
}

func TestFrameHandlerProtectedStorageRouteUsesValidatedIdentity(t *testing.T) {
	var validatorCalled bool
	var received app.RouteRequest
	route := appstorage.PutOwnStorageObjectRoute()
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(_ context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
			received = request
			payload, ok := request.Payload.(appstorage.PutOwnStorageObjectRequest)
			if !ok {
				t.Fatalf("request Payload = %T, want PutOwnStorageObjectRequest", request.Payload)
			}
			if payload.Collection != "progress" || payload.Key != "tutorial" || string(payload.ValueJSON) != `{"level":4}` {
				t.Fatalf("storage request payload = %#v, want mapped put request", payload)
			}
			if payload.Identity.Status != "" {
				t.Fatalf("storage payload Identity = %#v, want route handler to add identity later", payload.Identity)
			}
			return app.ApplicationResult{
				RequestID: request.RequestID,
				Route:     request.Route,
				Target:    request.Target,
				Session:   request.Session,
				Identity:  request.Identity,
				Payload: appstorage.StorageObjectResult{
					Status: appstorage.StorageObjectOperationStatusStored,
					Object: storagemodule.StorageObject{
						Identity: storagemodule.StorageObjectIdentity{
							Collection: "progress",
							Key:        "tutorial",
						},
						Value: storagemodule.StorageObjectValue{
							JSON: []byte(`{"level":4}`),
						},
						Version: 4,
					},
					Version: 4,
				},
			}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(_ context.Context, request app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			validatorCalled = true
			if request.Route != route || request.ConnectionID != "ws-1" || request.ConnectionEpoch != 7 {
				t.Fatalf("validator request = %#v, want storage route and frame/session metadata", request)
			}
			identity := app.ValidatedPlayerIdentity("player-1", app.Session{
				ConnectionID:    request.ConnectionID,
				ConnectionEpoch: request.ConnectionEpoch,
				PlayerID:        "metadata-player",
			})
			identity.SessionValidated = false
			return app.RouteAccessTokenValidationResult{
				Valid:    true,
				Identity: identity,
			}, nil
		})),
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID: "ws-1",
		Payload: mustMarshalAuthenticatedEnvelope(t, route, "redacted-token-text", &storagev1.PutOwnStorageObjectRequest{
			Collection: "progress",
			Key:        "tutorial",
			ValueJson:  `{"level":4}`,
		}),
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil", err)
	}
	if !validatorCalled {
		t.Fatal("validator was not called")
	}
	if received.PayloadType != "vibit.storage.v1.PutOwnStorageObjectRequest" {
		t.Fatalf("received PayloadType = %q, want storage request payload type", received.PayloadType)
	}
	if received.Identity.Status != app.IdentityValidationValidated ||
		received.Identity.PlayerID != "player-1" ||
		!received.Identity.PlayerIDValidated ||
		received.Identity.SessionValidated {
		t.Fatalf("received Identity = %#v, want validated request-token identity", received.Identity)
	}

	response := mustUnmarshalSingleResponse(t, responses)
	if response.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_COMMAND {
		t.Fatalf("response kind = %v, want command", response.GetKind())
	}
	if response.GetPayloadType() != "vibit.storage.v1.PutOwnStorageObjectResponse" {
		t.Fatalf("response payload_type = %q, want storage put response", response.GetPayloadType())
	}
	responsePayload, err := DecodePayload(response.GetPayloadType(), response.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload(storage response) error = %v, want nil", err)
	}
	storageResponse, ok := responsePayload.(*storagev1.PutOwnStorageObjectResponse)
	if !ok {
		t.Fatalf("response payload = %T, want PutOwnStorageObjectResponse", responsePayload)
	}
	if storageResponse.GetVersion() != 4 ||
		storageResponse.GetObject().GetCollection() != "progress" ||
		storageResponse.GetObject().GetValueJson() != `{"level":4}` {
		t.Fatalf("storage response = %#v, want mapped storage response", storageResponse)
	}
}

func TestFrameHandlerLeavesPublicAuthenticationRouteExplicit(t *testing.T) {
	publicRoute := app.AuthenticateWithDeviceCredentialRoute()
	var received app.RouteRequest
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(_ context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
			received = request
			payload, ok := request.Payload.(appauth.DeviceCredentialAuthenticationRequest)
			if !ok {
				t.Fatalf("request Payload = %T, want DeviceCredentialAuthenticationRequest", request.Payload)
			}
			if payload.CredentialProof != "redacted-credential-proof" ||
				payload.RequestedPlayerID != "player-1" ||
				payload.ClientInstanceID != "client-1" ||
				payload.AccountCreationIntent != appauth.AccountCreationIntentReject {
				t.Fatalf("request payload = %#v, want login request", payload)
			}
			return app.ApplicationResult{
				RequestID: request.RequestID,
				Route:     request.Route,
				Target:    request.Target,
				Session:   request.Session,
				Identity:  request.Identity,
				Payload: appauth.AuthenticationResult{
					Status:         appauth.AuthenticationStatusAuthenticated,
					ActorKind:      app.ActorKindPlayer,
					PlayerID:       "player-1",
					AccessToken:    "redacted-access-token",
					TokenType:      appauth.TokenTypeOpaqueAccess,
					SessionID:      "runtime-session-1",
					SessionCreated: true,
					TokenRecordID:  "token-record-1",
				},
			}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			t.Fatal("validator was called for public authentication route")
			return app.RouteAccessTokenValidationResult{}, nil
		})),
	}
	payload := mustMarshalEnvelope(t, publicRoute, &authenticationv1.AuthenticateWithDeviceCredentialRequest{
		CredentialProof:       "redacted-credential-proof",
		RequestedPlayerId:     "player-1",
		ClientInstanceId:      "client-1",
		AccountCreationIntent: authenticationv1.AccountCreationIntent_ACCOUNT_CREATION_INTENT_AUTHENTICATE_EXISTING_ONLY,
	})

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{ConnectionID: "ws-1", Payload: payload})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil", err)
	}
	if received.Route != publicRoute {
		t.Fatalf("received route = %#v, want public auth route", received.Route)
	}
	if received.Identity.Status != app.IdentityValidationMetadataOnly {
		t.Fatalf("received Identity.Status = %q, want metadata_only", received.Identity.Status)
	}
	response := mustUnmarshalSingleResponse(t, responses)
	if response.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_COMMAND {
		t.Fatalf("response kind = %v, want command", response.GetKind())
	}
	if response.GetPayloadType() != "vibit.authentication.v1.AuthenticateWithDeviceCredentialResponse" {
		t.Fatalf("response payload_type = %q, want login response", response.GetPayloadType())
	}
	responsePayload, err := DecodePayload(response.GetPayloadType(), response.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload(login response) error = %v, want nil", err)
	}
	loginResponse, ok := responsePayload.(*authenticationv1.AuthenticateWithDeviceCredentialResponse)
	if !ok {
		t.Fatalf("response payload = %T, want AuthenticateWithDeviceCredentialResponse", responsePayload)
	}
	if loginResponse.GetAccessToken() != "redacted-access-token" ||
		loginResponse.GetTokenType() != string(appauth.TokenTypeOpaqueAccess) ||
		loginResponse.GetTokenRecordId() != "token-record-1" {
		t.Fatalf("login response = %#v, want mapped authentication result", loginResponse)
	}
	if response.GetSession().GetConnectionId() != "ws-1" ||
		response.GetSession().GetSessionId() != "runtime-session-1" ||
		response.GetSession().GetPlayerId() != "player-1" {
		t.Fatalf("response session = %#v, want login-created runtime session carrier", response.GetSession())
	}
}

func TestFrameHandlerLeavesLogoutRouteServiceValidatedAndUnwrapped(t *testing.T) {
	logoutRoute := app.LogoutAccessTokenRoute()
	var received app.RouteRequest
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(_ context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
			received = request
			payload, ok := request.Payload.(appauth.LogoutAccessTokenRequest)
			if !ok {
				t.Fatalf("request Payload = %T, want LogoutAccessTokenRequest", request.Payload)
			}
			if payload.AccessToken != "redacted-access-token" || payload.LogoutReason != "player_logout" {
				t.Fatalf("request payload = %#v, want logout request", payload)
			}
			return app.ApplicationResult{
				RequestID: request.RequestID,
				Route:     request.Route,
				Target:    request.Target,
				Session:   request.Session,
				Identity:  request.Identity,
				Payload: appauth.LogoutAccessTokenResult{
					Status:        appauth.LogoutStatusRevoked,
					Revoked:       true,
					LogoutScope:   appauth.LogoutScopeToken,
					TokenRecordID: "token-record-1",
				},
			}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			t.Fatal("validator was called for service-validated logout route")
			return app.RouteAccessTokenValidationResult{}, nil
		})),
	}
	payload := mustMarshalEnvelope(t, logoutRoute, &authenticationv1.LogoutAccessTokenRequest{
		AccessToken:  "redacted-access-token",
		LogoutReason: "player_logout",
	})

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{ConnectionID: "ws-1", Payload: payload})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil", err)
	}
	if received.Route != logoutRoute {
		t.Fatalf("received route = %#v, want logout route", received.Route)
	}
	if received.PayloadType != "vibit.authentication.v1.LogoutAccessTokenRequest" {
		t.Fatalf("received PayloadType = %q, want logout request payload type", received.PayloadType)
	}
	if received.Identity.Status != app.IdentityValidationMetadataOnly {
		t.Fatalf("received Identity.Status = %q, want metadata_only", received.Identity.Status)
	}
	response := mustUnmarshalSingleResponse(t, responses)
	if response.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_COMMAND {
		t.Fatalf("response kind = %v, want command", response.GetKind())
	}
	if response.GetPayloadType() != "vibit.authentication.v1.LogoutAccessTokenResponse" {
		t.Fatalf("response payload_type = %q, want logout response", response.GetPayloadType())
	}
	responsePayload, err := DecodePayload(response.GetPayloadType(), response.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload(logout response) error = %v, want nil", err)
	}
	logoutResponse, ok := responsePayload.(*authenticationv1.LogoutAccessTokenResponse)
	if !ok {
		t.Fatalf("response payload = %T, want LogoutAccessTokenResponse", responsePayload)
	}
	if !logoutResponse.GetRevoked() ||
		logoutResponse.GetLogoutStatus() != string(appauth.LogoutStatusRevoked) ||
		logoutResponse.GetLogoutScope() != "presented_access_token" ||
		logoutResponse.GetTokenRecordId() != "token-record-1" {
		t.Fatalf("logout response = %#v, want mapped logout result", logoutResponse)
	}
}

func TestFrameHandlerRejectsAuthenticatedWrapperOnLogoutRoute(t *testing.T) {
	var dispatched bool
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
			dispatched = true
			return app.ApplicationResult{}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			t.Fatal("validator was called for wrapped logout route")
			return app.RouteAccessTokenValidationResult{}, nil
		})),
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID: "ws-1",
		Payload: mustMarshalAuthenticatedEnvelope(t, app.LogoutAccessTokenRoute(), "redacted-wrapper-token", &authenticationv1.LogoutAccessTokenRequest{
			AccessToken: "redacted-inner-token",
		}),
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil error envelope", err)
	}
	if dispatched {
		t.Fatal("domain dispatcher was called, want wrapped logout to stop dispatch")
	}
	response := mustUnmarshalSingleResponse(t, responses)
	assertErrorEnvelope(t, response, app.ErrorCodeAuthenticationTokenMalformed)
	assertNoFrameErrorSecretLeak(t, response, "redacted-wrapper-token", "redacted-inner-token")
}

func TestAuthenticatedRequestPayloadTypeKeepsEnvelopeUnchanged(t *testing.T) {
	envelope := mustBuildEnvelope(
		t,
		inventory.GetInventoryRoute(),
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{SessionID: "session-1", PlayerID: "metadata-player"},
		&authenticationv1.AuthenticatedRequest{
			AccessToken:      "redacted-token-text",
			InnerPayloadType: "vibit.inventory.v1.GetInventoryRequest",
			InnerPayload:     []byte("redacted-inner-payload"),
		},
	)

	if envelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_QUERY ||
		envelope.GetModule() != inventory.ModuleName ||
		envelope.GetName() != inventory.QueryGetInventory ||
		envelope.GetPayloadType() != "vibit.authentication.v1.AuthenticatedRequest" {
		t.Fatalf("envelope = %#v, want original route with authenticated payload wrapper", envelope)
	}
	if envelope.GetSession().GetPlayerId() != "metadata-player" {
		t.Fatalf("session player_id = %q, want metadata-only value unchanged", envelope.GetSession().GetPlayerId())
	}
}

func mustMarshalAuthenticatedEnvelope(t *testing.T, route app.RouteKey, accessToken string, innerPayload proto.Message) []byte {
	t.Helper()

	innerBytes, err := proto.Marshal(innerPayload)
	if err != nil {
		t.Fatalf("proto.Marshal(innerPayload) error = %v, want nil", err)
	}
	envelope := mustBuildEnvelope(
		t,
		route,
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{SessionID: "session-1", PlayerID: "metadata-player", ConnectionEpoch: 7},
		&authenticationv1.AuthenticatedRequest{
			AccessToken:      accessToken,
			InnerPayloadType: PayloadType(innerPayload),
			InnerPayload:     innerBytes,
		},
	)
	encoded, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal(envelope) error = %v, want nil", err)
	}
	return encoded
}

func mustUnmarshalSingleResponse(t *testing.T, responses [][]byte) *protocolv1.Envelope {
	t.Helper()

	if len(responses) != 1 {
		t.Fatalf("responses len = %d, want 1", len(responses))
	}
	return mustUnmarshalFrameEnvelope(t, responses[0])
}

func assertErrorEnvelope(t *testing.T, envelope *protocolv1.Envelope, code app.ErrorCode) {
	t.Helper()

	if envelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_ERROR {
		t.Fatalf("envelope kind = %v, want error", envelope.GetKind())
	}
	if envelope.GetError().GetCode() != string(code) {
		t.Fatalf("error code = %q, want %s", envelope.GetError().GetCode(), code)
	}
	if envelope.GetPayloadType() != "" || len(envelope.GetPayload()) != 0 {
		t.Fatalf("error payload_type = %q payload len = %d, want no payload", envelope.GetPayloadType(), len(envelope.GetPayload()))
	}
}

func assertNoFrameErrorSecretLeak(t *testing.T, envelope *protocolv1.Envelope, secrets ...string) {
	t.Helper()

	errorText := envelope.GetError().GetCode() + " " + envelope.GetError().GetMessage()
	for _, secret := range secrets {
		if secret != "" && strings.Contains(errorText, secret) {
			t.Fatalf("error envelope %q leaks secret %q", errorText, secret)
		}
	}
}

type routeTokenValidatorFunc func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error)

func (f routeTokenValidatorFunc) ValidateRouteAccessToken(ctx context.Context, request app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
	if f == nil {
		return app.RouteAccessTokenValidationResult{}, errors.New("route token validator function is nil")
	}
	return f(ctx, request)
}
