package protobuf

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	authenticationv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/authentication/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"google.golang.org/protobuf/proto"
)

func TestFrameHandlerBindsConnectionThroughSystemRoute(t *testing.T) {
	boundAt := time.Date(2026, 5, 17, 8, 45, 0, 0, time.UTC)
	var received app.ConnectionBindingRequest
	handler := FrameHandler{
		Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
			t.Fatal("dispatcher should not be called for BindConnection")
			return app.ApplicationResult{}, nil
		}),
		RouteProtector: app.NewRouteProtector(routeTokenValidatorFunc(func(context.Context, app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
			t.Fatal("route protector should not be called for BindConnection")
			return app.RouteAccessTokenValidationResult{}, nil
		})),
		ConnectionBinder: connectionBinderFunc(func(_ context.Context, request app.ConnectionBindingRequest) (app.ConnectionBindingResult, error) {
			received = request
			identity := app.ValidatedPlayerIdentity("player-1", app.Session{
				ConnectionID:    request.ConnectionID,
				ConnectionEpoch: request.ConnectionEpoch,
				PlayerID:        "player-1",
			})
			identity.SessionValidated = false
			return app.ConnectionBindingResult{
				BindingStatus:    app.ConnectionBindingStatusBound,
				Bound:            true,
				Identity:         identity,
				ConnectionID:     request.ConnectionID,
				ConnectionEpoch:  request.ConnectionEpoch,
				ClientInstanceID: request.ClientInstanceID,
				BoundAt:          boundAt,
			}, nil
		}),
	}

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID:    "ws-server-1",
		ConnectionEpoch: 3,
		Payload: mustMarshalEnvelope(t, app.BindConnectionRoute(), &authenticationv1.BindConnectionRequest{
			AccessToken:      "redacted-token-text",
			ClientInstanceId: "client-1",
		}),
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil", err)
	}
	if received.AccessToken != "redacted-token-text" ||
		received.Route != app.BindConnectionRoute() ||
		received.ConnectionID != "ws-server-1" ||
		received.ConnectionEpoch != 3 ||
		received.ClientInstanceID != "client-1" {
		t.Fatalf("binder request = %#v, want decoded payload and server-observed connection metadata", received)
	}

	envelope := mustUnmarshalSingleResponse(t, responses)
	if envelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_SYSTEM {
		t.Fatalf("response kind = %v, want system", envelope.GetKind())
	}
	if envelope.GetPayloadType() != "vibit.authentication.v1.BindConnectionResponse" {
		t.Fatalf("payload_type = %q, want BindConnectionResponse", envelope.GetPayloadType())
	}
	if envelope.GetSession().GetConnectionId() != "ws-server-1" ||
		envelope.GetSession().GetConnectionEpoch() != 3 ||
		envelope.GetSession().GetPlayerId() != "player-1" {
		t.Fatalf("response session = %#v, want bound server connection and player metadata", envelope.GetSession())
	}

	payload, err := DecodePayload(envelope.GetPayloadType(), envelope.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload() error = %v, want nil", err)
	}
	response, ok := payload.(*authenticationv1.BindConnectionResponse)
	if !ok {
		t.Fatalf("response payload = %T, want BindConnectionResponse", payload)
	}
	if response.GetBindingStatus() != authenticationv1.ConnectionBindingStatus_CONNECTION_BINDING_STATUS_BOUND ||
		response.GetActorKind() != string(app.ActorKindPlayer) ||
		response.GetPlayerId() != "player-1" ||
		response.GetConnectionId() != "ws-server-1" ||
		response.GetConnectionEpoch() != 3 ||
		response.GetSessionValidated() {
		t.Fatalf("BindConnectionResponse = %#v, want bound response without session validation", response)
	}
	if response.GetBoundAt() != boundAt.Format(time.RFC3339Nano) {
		t.Fatalf("BoundAt = %q, want %q", response.GetBoundAt(), boundAt.Format(time.RFC3339Nano))
	}
	if strings.Contains(response.String(), "redacted-token-text") {
		t.Fatalf("BindConnectionResponse leaks access token: %v", response)
	}
}

func TestFrameHandlerMapsConnectionBindingFailuresToErrorEnvelopes(t *testing.T) {
	tests := []struct {
		name    string
		binder  ConnectionBinder
		payload proto.Message
		want    app.ErrorCode
	}{
		{
			name:    "missing binder",
			binder:  nil,
			payload: &authenticationv1.BindConnectionRequest{AccessToken: "redacted-token-text"},
			want:    app.ErrorCodeConnectionBindingUnavailable,
		},
		{
			name: "malformed token",
			binder: connectionBinderFunc(func(context.Context, app.ConnectionBindingRequest) (app.ConnectionBindingResult, error) {
				return app.ConnectionBindingResult{BindingStatus: app.ConnectionBindingStatusRejected}, &app.ApplicationError{
					Code:    app.ErrorCodeConnectionBindingTokenMalformed,
					Message: "connection binding token proof is malformed",
					Route:   app.BindConnectionRoute(),
				}
			}),
			payload: &authenticationv1.BindConnectionRequest{AccessToken: " redacted-token-text "},
			want:    app.ErrorCodeConnectionBindingTokenMalformed,
		},
		{
			name: "invalid token",
			binder: connectionBinderFunc(func(context.Context, app.ConnectionBindingRequest) (app.ConnectionBindingResult, error) {
				return app.ConnectionBindingResult{BindingStatus: app.ConnectionBindingStatusRejected}, &app.ApplicationError{
					Code:    app.ErrorCodeConnectionBindingTokenInvalid,
					Message: "connection binding token proof is invalid",
					Route:   app.BindConnectionRoute(),
				}
			}),
			payload: &authenticationv1.BindConnectionRequest{AccessToken: "redacted-token-text"},
			want:    app.ErrorCodeConnectionBindingTokenInvalid,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := FrameHandler{
				Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
					t.Fatal("dispatcher should not be called for BindConnection failure")
					return app.ApplicationResult{}, nil
				}),
				ConnectionBinder: tc.binder,
			}
			responses, err := handler.HandleFrame(context.Background(), FrameRequest{
				ConnectionID: "ws-server-1",
				Payload:      mustMarshalEnvelopeWithPayloadType(t, app.BindConnectionRoute(), "vibit.authentication.v1.BindConnectionRequest", tc.payload),
			})
			if err != nil {
				t.Fatalf("HandleFrame() error = %v, want nil error envelope", err)
			}
			envelope := mustUnmarshalSingleResponse(t, responses)
			assertErrorEnvelope(t, envelope, tc.want)
			assertNoFrameErrorSecretLeak(t, envelope, "redacted-token-text")
		})
	}
}

func TestBindConnectionPayloadTypeKeepsEnvelopeUnchanged(t *testing.T) {
	envelope := mustBuildEnvelope(
		t,
		app.BindConnectionRoute(),
		"request-1",
		app.Target{Scope: app.TargetScopeSystem, ID: "runtime"},
		app.Session{SessionID: "session-1", PlayerID: "metadata-player", ConnectionEpoch: 7},
		&authenticationv1.BindConnectionRequest{
			AccessToken:      "redacted-token-text",
			ClientInstanceId: "client-1",
		},
	)

	if envelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_SYSTEM ||
		envelope.GetModule() != "runtime.authentication" ||
		envelope.GetName() != "BindConnection" ||
		envelope.GetPayloadType() != "vibit.authentication.v1.BindConnectionRequest" {
		t.Fatalf("envelope = %#v, want BindConnection system route with payload carrier", envelope)
	}
	if envelope.GetSession().GetPlayerId() != "metadata-player" {
		t.Fatalf("session player_id = %q, want metadata-only value unchanged", envelope.GetSession().GetPlayerId())
	}
}

type connectionBinderFunc func(context.Context, app.ConnectionBindingRequest) (app.ConnectionBindingResult, error)

func (f connectionBinderFunc) BindConnection(ctx context.Context, request app.ConnectionBindingRequest) (app.ConnectionBindingResult, error) {
	if f == nil {
		return app.ConnectionBindingResult{}, errors.New("connection binder function is nil")
	}
	return f(ctx, request)
}

func mustMarshalEnvelopeWithPayloadType(t *testing.T, route app.RouteKey, payloadType string, payload proto.Message) []byte {
	t.Helper()

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		t.Fatalf("proto.Marshal(payload) error = %v, want nil", err)
	}
	envelope := mustBuildEnvelope(
		t,
		route,
		"request-1",
		app.Target{Scope: app.TargetScopeSystem, ID: "runtime"},
		app.Session{SessionID: "session-1", PlayerID: "metadata-player"},
		&authenticationv1.BindConnectionRequest{},
	)
	envelope.PayloadType = payloadType
	envelope.Payload = payloadBytes
	encoded, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal(envelope) error = %v, want nil", err)
	}
	return encoded
}
