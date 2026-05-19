package protobuf

import (
	"context"
	"errors"
	"testing"

	"github.com/iceiko/vibit/runtime/internal/app"
	inventoryv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/inventory/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
)

func TestFrameHandlerDispatchesGrantItemAndBuildsResponseFrame(t *testing.T) {
	fixture := newInventoryRequestLoopFixture(t, true, true)
	handler := FrameHandler{Dispatcher: fixture.Dispatcher}
	requestPayload := mustMarshalEnvelope(t, inventory.GrantItemRoute(), &inventoryv1.GrantItemRequest{
		PlayerId:    "player-1",
		ItemId:      "item-1",
		Quantity:    3,
		Reason:      "reward",
		RequestedBy: "admin-1",
	})

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{
		ConnectionID: "ws-1",
		RemoteAddr:   "127.0.0.1:1",
		Payload:      requestPayload,
	})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil", err)
	}
	if len(responses) != 1 {
		t.Fatalf("responses len = %d, want 1", len(responses))
	}

	responseEnvelope := mustUnmarshalFrameEnvelope(t, responses[0])
	if responseEnvelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_COMMAND {
		t.Fatalf("response kind = %v, want command", responseEnvelope.GetKind())
	}
	if responseEnvelope.GetRequestId() != "request-1" {
		t.Fatalf("request_id = %q, want request-1", responseEnvelope.GetRequestId())
	}
	if responseEnvelope.GetPayloadType() != "vibit.inventory.v1.GrantItemResponse" {
		t.Fatalf("payload_type = %q, want GrantItemResponse", responseEnvelope.GetPayloadType())
	}
	if responseEnvelope.GetSession().GetConnectionId() != "ws-1" {
		t.Fatalf("connection_id = %q, want frame connection id", responseEnvelope.GetSession().GetConnectionId())
	}

	payload, err := DecodePayload(responseEnvelope.GetPayloadType(), responseEnvelope.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload() error = %v, want nil", err)
	}
	response, ok := payload.(*inventoryv1.GrantItemResponse)
	if !ok {
		t.Fatalf("response payload = %T, want *inventoryv1.GrantItemResponse", payload)
	}
	if response.GetPlayerId() != "player-1" || response.GetItemId() != "item-1" || response.GetQuantity() != 3 || response.GetNewQuantity() != 3 {
		t.Fatalf("response = %#v, want mapped GrantItemResponse", response)
	}
}

func TestFrameHandlerBuildsApplicationErrorFrame(t *testing.T) {
	fixture := newInventoryRequestLoopFixture(t, false, true)
	handler := FrameHandler{Dispatcher: fixture.Dispatcher}
	requestPayload := mustMarshalEnvelope(t, inventory.GrantItemRoute(), &inventoryv1.GrantItemRequest{
		PlayerId:    "player-1",
		ItemId:      "item-1",
		Quantity:    3,
		Reason:      "reward",
		RequestedBy: "admin-1",
	})

	responses, err := handler.HandleFrame(context.Background(), FrameRequest{Payload: requestPayload})
	if err != nil {
		t.Fatalf("HandleFrame() error = %v, want nil application error frame", err)
	}
	if len(responses) != 1 {
		t.Fatalf("responses len = %d, want 1", len(responses))
	}

	responseEnvelope := mustUnmarshalFrameEnvelope(t, responses[0])
	if responseEnvelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_ERROR {
		t.Fatalf("response kind = %v, want error", responseEnvelope.GetKind())
	}
	if responseEnvelope.GetPayloadType() != "" || len(responseEnvelope.GetPayload()) != 0 {
		t.Fatalf("error payload_type = %q payload len = %d, want no payload", responseEnvelope.GetPayloadType(), len(responseEnvelope.GetPayload()))
	}
	if responseEnvelope.GetError().GetCode() != string(inventory.ErrorCodeInventoryPermission) {
		t.Fatalf("error code = %q, want %s", responseEnvelope.GetError().GetCode(), inventory.ErrorCodeInventoryPermission)
	}
	if responseEnvelope.GetError().GetRequestId() != "request-1" {
		t.Fatalf("error request_id = %q, want request-1", responseEnvelope.GetError().GetRequestId())
	}
}

func TestRequestWithFrameMetadataRefreshesMetadataOnlyIdentity(t *testing.T) {
	request := app.RouteRequest{
		Session: app.Session{
			SessionID: "session-1",
			PlayerID:  "player-1",
		},
		Identity: app.MetadataOnlyIdentityFromSession(app.Session{
			SessionID: "session-1",
			PlayerID:  "player-1",
		}),
	}

	request = requestWithFrameMetadata(request, FrameRequest{ConnectionID: " ws-1 ", ConnectionEpoch: 11})

	if request.Session.ConnectionID != "ws-1" {
		t.Fatalf("Session.ConnectionID = %q, want ws-1", request.Session.ConnectionID)
	}
	if request.Identity.Status != app.IdentityValidationMetadataOnly {
		t.Fatalf("Identity.Status = %q, want %q", request.Identity.Status, app.IdentityValidationMetadataOnly)
	}
	if request.Identity.ConnectionID != "ws-1" || request.Identity.PlayerID != "player-1" || request.Identity.SessionID != "session-1" {
		t.Fatalf("Identity = %#v, want refreshed metadata-only identity", request.Identity)
	}
	if request.Identity.PlayerIDValidated || request.Identity.SessionValidated {
		t.Fatalf("Identity validation flags = %#v, want metadata-only identity", request.Identity)
	}
	if request.Session.ConnectionEpoch != 11 || request.Identity.ConnectionEpoch != 11 {
		t.Fatalf("connection epoch session=%d identity=%d, want 11", request.Session.ConnectionEpoch, request.Identity.ConnectionEpoch)
	}
}

func TestFrameHandlerRejectsMalformedFramePayload(t *testing.T) {
	fixture := newInventoryRequestLoopFixture(t, true, true)
	handler := FrameHandler{Dispatcher: fixture.Dispatcher}

	_, err := handler.HandleFrame(context.Background(), FrameRequest{Payload: []byte("not protobuf")})
	if err == nil {
		t.Fatal("HandleFrame() error = nil, want protobuf decode error")
	}
}

func TestFrameHandlerPropagatesInternalDispatchError(t *testing.T) {
	sentinel := errors.New("internal dispatch failed")
	handler := FrameHandler{Dispatcher: requestLoopDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
		return app.ApplicationResult{}, sentinel
	})}
	requestPayload := mustMarshalEnvelope(t, inventory.GetInventoryRoute(), &inventoryv1.GetInventoryRequest{
		PlayerId:    "player-1",
		RequestedBy: "player-1",
	})

	_, err := handler.HandleFrame(context.Background(), FrameRequest{Payload: requestPayload})
	if !errors.Is(err, sentinel) {
		t.Fatalf("HandleFrame() error = %v, want sentinel", err)
	}
}
