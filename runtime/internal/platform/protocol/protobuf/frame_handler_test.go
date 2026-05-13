package protobuf

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	inventoryv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/inventory/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"google.golang.org/protobuf/proto"
)

func TestFrameHandlerDispatchesGrantItemAndBuildsResponseFrame(t *testing.T) {
	dispatcher := newFrameHandlerInventoryDispatcher(t, true, true)
	handler := FrameHandler{Dispatcher: dispatcher}
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
	dispatcher := newFrameHandlerInventoryDispatcher(t, false, true)
	handler := FrameHandler{Dispatcher: dispatcher}
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

func TestFrameHandlerRejectsMalformedFramePayload(t *testing.T) {
	handler := FrameHandler{Dispatcher: newFrameHandlerInventoryDispatcher(t, true, true)}

	_, err := handler.HandleFrame(context.Background(), FrameRequest{Payload: []byte("not protobuf")})
	if err == nil {
		t.Fatal("HandleFrame() error = nil, want protobuf decode error")
	}
}

func TestFrameHandlerPropagatesInternalDispatchError(t *testing.T) {
	sentinel := errors.New("internal dispatch failed")
	handler := FrameHandler{Dispatcher: frameHandlerDispatcherFunc(func(context.Context, app.RouteRequest) (app.ApplicationResult, error) {
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

func newFrameHandlerInventoryDispatcher(t *testing.T, grantAllowed bool, readAllowed bool) *app.Dispatcher {
	t.Helper()

	handlers := inventory.Handlers{
		Repository:       newFrameHandlerMemoryRepository(),
		PermissionPolicy: frameHandlerPermissionPolicy{grantAllowed: grantAllowed, readAllowed: readAllowed},
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 16},
		EventIDs:         &frameHandlerEventIDs{},
		Clock:            frameHandlerClock{},
	}
	dispatcher := app.NewDispatcher()
	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		t.Fatalf("RegisterRoutes() error = %v, want nil", err)
	}
	return dispatcher
}

func mustMarshalEnvelope(t *testing.T, route app.RouteKey, payload proto.Message) []byte {
	t.Helper()

	envelope, err := BuildEnvelope(
		route,
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{SessionID: "session-1", PlayerID: "player-1", ConnectionEpoch: 7},
		payload,
	)
	if err != nil {
		t.Fatalf("BuildEnvelope() error = %v, want nil", err)
	}
	encoded, err := proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("proto.Marshal(envelope) error = %v, want nil", err)
	}
	return encoded
}

func mustUnmarshalFrameEnvelope(t *testing.T, payload []byte) *protocolv1.Envelope {
	t.Helper()

	var envelope protocolv1.Envelope
	if err := proto.Unmarshal(payload, &envelope); err != nil {
		t.Fatalf("proto.Unmarshal(response) error = %v, want nil", err)
	}
	return &envelope
}

type frameHandlerDispatcherFunc func(context.Context, app.RouteRequest) (app.ApplicationResult, error)

func (f frameHandlerDispatcherFunc) Dispatch(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	return f(ctx, request)
}

type frameHandlerPermissionPolicy struct {
	grantAllowed bool
	readAllowed  bool
}

func (p frameHandlerPermissionPolicy) CanGrantItem(context.Context, string, string) (bool, error) {
	return p.grantAllowed, nil
}

func (p frameHandlerPermissionPolicy) CanReadInventory(context.Context, string, string) (bool, error) {
	return p.readAllowed, nil
}

type frameHandlerMemoryRepository struct {
	items map[string]map[string]int64
}

func newFrameHandlerMemoryRepository() *frameHandlerMemoryRepository {
	return &frameHandlerMemoryRepository{
		items: make(map[string]map[string]int64),
	}
}

func (r *frameHandlerMemoryRepository) GetInventory(_ context.Context, playerID string) ([]inventory.Item, error) {
	playerItems := r.items[playerID]
	items := make([]inventory.Item, 0, len(playerItems))
	for itemID, quantity := range playerItems {
		items = append(items, inventory.Item{ItemID: itemID, Quantity: quantity})
	}
	return items, nil
}

func (r *frameHandlerMemoryRepository) GrantItem(_ context.Context, mutation inventory.GrantItemMutation) (inventory.Item, error) {
	if r.items[mutation.PlayerID] == nil {
		r.items[mutation.PlayerID] = make(map[string]int64)
	}
	r.items[mutation.PlayerID][mutation.ItemID] += mutation.Quantity
	return inventory.Item{
		ItemID:   mutation.ItemID,
		Quantity: r.items[mutation.PlayerID][mutation.ItemID],
	}, nil
}

type frameHandlerEventIDs struct {
	next int
}

func (g *frameHandlerEventIDs) NewEventID() string {
	g.next += 1
	return "frame-event-1"
}

type frameHandlerClock struct{}

func (frameHandlerClock) Now() time.Time {
	return time.Date(2026, 5, 13, 12, 0, 0, 0, time.UTC)
}
