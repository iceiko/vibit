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
)

func TestRouteRequestFromEnvelopeForDispatchMapsGrantItemPayload(t *testing.T) {
	envelope, err := BuildEnvelope(
		inventory.GrantItemRoute(),
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{SessionID: "session-1", PlayerID: "player-1"},
		&inventoryv1.GrantItemRequest{
			PlayerId:    "player-1",
			ItemId:      "item-1",
			Quantity:    3,
			Reason:      "reward",
			RequestedBy: "admin-1",
		},
	)
	if err != nil {
		t.Fatalf("BuildEnvelope() error = %v, want nil", err)
	}

	request, err := RouteRequestFromEnvelopeForDispatch(envelope)
	if err != nil {
		t.Fatalf("RouteRequestFromEnvelopeForDispatch() error = %v, want nil", err)
	}

	payload, ok := request.Payload.(inventory.GrantItemRequest)
	if !ok {
		t.Fatalf("Payload = %T, want inventory.GrantItemRequest", request.Payload)
	}
	if payload.PlayerID != "player-1" || payload.ItemID != "item-1" || payload.Quantity != 3 || payload.Reason != "reward" || payload.RequestedBy != "admin-1" {
		t.Fatalf("Payload = %#v, want mapped GrantItem request", payload)
	}
	if request.PayloadType != "vibit.inventory.v1.GrantItemRequest" {
		t.Fatalf("PayloadType = %q, want generated payload type", request.PayloadType)
	}
	if len(request.PayloadBytes) == 0 {
		t.Fatal("PayloadBytes is empty, want original encoded payload")
	}
}

func TestRouteRequestFromEnvelopeForDispatchMapsGetInventoryPayload(t *testing.T) {
	envelope, err := BuildEnvelope(
		inventory.GetInventoryRoute(),
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{SessionID: "session-1", PlayerID: "player-1"},
		&inventoryv1.GetInventoryRequest{
			PlayerId:    "player-1",
			RequestedBy: "player-1",
		},
	)
	if err != nil {
		t.Fatalf("BuildEnvelope() error = %v, want nil", err)
	}

	request, err := RouteRequestFromEnvelopeForDispatch(envelope)
	if err != nil {
		t.Fatalf("RouteRequestFromEnvelopeForDispatch() error = %v, want nil", err)
	}

	payload, ok := request.Payload.(inventory.GetInventoryRequest)
	if !ok {
		t.Fatalf("Payload = %T, want inventory.GetInventoryRequest", request.Payload)
	}
	if payload.PlayerID != "player-1" || payload.RequestedBy != "player-1" {
		t.Fatalf("Payload = %#v, want mapped GetInventory request", payload)
	}
}

func TestRouteRequestWithDomainPayloadRejectsWrongInventoryPayload(t *testing.T) {
	_, err := RouteRequestWithDomainPayload(app.RouteRequest{
		Route:   inventory.GrantItemRoute(),
		Payload: &inventoryv1.GetInventoryRequest{},
	})
	if err == nil {
		t.Fatal("RouteRequestWithDomainPayload() error = nil, want bridge error")
	}

	var bridgeErr *PayloadBridgeError
	if !errors.As(err, &bridgeErr) {
		t.Fatalf("error = %v, want *PayloadBridgeError", err)
	}
}

func TestProtoPayloadFromApplicationResultMapsInventoryResponses(t *testing.T) {
	grantPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: inventory.GrantItemRoute(),
		Payload: inventory.GrantItemResponse{
			PlayerID:    "player-1",
			ItemID:      "item-1",
			Quantity:    2,
			NewQuantity: 5,
			Event:       inventory.EventItemGranted,
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult() grant error = %v, want nil", err)
	}
	grantResponse, ok := grantPayload.(*inventoryv1.GrantItemResponse)
	if !ok {
		t.Fatalf("grant payload = %T, want *inventoryv1.GrantItemResponse", grantPayload)
	}
	if grantResponse.GetPlayerId() != "player-1" || grantResponse.GetItemId() != "item-1" || grantResponse.GetQuantity() != 2 || grantResponse.GetNewQuantity() != 5 || grantResponse.GetEvent() != inventory.EventItemGranted {
		t.Fatalf("grant response = %#v, want mapped Protobuf response", grantResponse)
	}

	queryPayload, err := ProtoPayloadFromApplicationResult(app.ApplicationResult{
		Route: inventory.GetInventoryRoute(),
		Payload: inventory.GetInventoryResponse{
			PlayerID: "player-1",
			Items: []inventory.Item{
				{ItemID: "item-1", Quantity: 1},
				{ItemID: "item-2", Quantity: 3},
			},
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult() query error = %v, want nil", err)
	}
	queryResponse, ok := queryPayload.(*inventoryv1.GetInventoryResponse)
	if !ok {
		t.Fatalf("query payload = %T, want *inventoryv1.GetInventoryResponse", queryPayload)
	}
	if queryResponse.GetPlayerId() != "player-1" || len(queryResponse.GetItems()) != 2 {
		t.Fatalf("query response = %#v, want mapped inventory response", queryResponse)
	}
	if queryResponse.GetItems()[0].GetItemId() != "item-1" || queryResponse.GetItems()[1].GetQuantity() != 3 {
		t.Fatalf("query items = %#v, want mapped items", queryResponse.GetItems())
	}
}

func TestProtoPayloadFromApplicationEventMapsItemGranted(t *testing.T) {
	occurredAt := time.Date(2026, 5, 13, 9, 45, 30, 123, time.UTC)
	payload, err := ProtoPayloadFromApplicationEvent(app.ApplicationEvent{
		Route: app.RouteKey{Kind: app.MessageKindEvent, Module: inventory.ModuleName, Name: inventory.EventItemGranted},
		Payload: inventory.ItemGrantedEvent{
			EventID:     "event-1",
			OccurredAt:  occurredAt,
			PlayerID:    "player-1",
			ItemID:      "item-1",
			Quantity:    2,
			NewQuantity: 5,
			Reason:      "reward",
		},
	})
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationEvent() error = %v, want nil", err)
	}

	event, ok := payload.(*inventoryv1.ItemGranted)
	if !ok {
		t.Fatalf("payload = %T, want *inventoryv1.ItemGranted", payload)
	}
	if event.GetEventId() != "event-1" || event.GetOccurredAt() != "2026-05-13T09:45:30.000000123Z" || event.GetPlayerId() != "player-1" || event.GetItemId() != "item-1" || event.GetQuantity() != 2 || event.GetNewQuantity() != 5 || event.GetReason() != "reward" {
		t.Fatalf("event = %#v, want mapped ItemGranted event", event)
	}
}

func TestInventoryProtobufDomainBridgeDispatchesAndBuildsResponseEnvelope(t *testing.T) {
	fixture := newInventoryRequestLoopFixture(t, true, true)

	envelope := mustBuildEnvelope(
		t,
		inventory.GrantItemRoute(),
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{ConnectionID: "connection-1", SessionID: "session-1", PlayerID: "player-1", ConnectionEpoch: 7},
		&inventoryv1.GrantItemRequest{
			PlayerId:    "player-1",
			ItemId:      "item-1",
			Quantity:    2,
			Reason:      "reward",
			RequestedBy: "admin-1",
		},
	)

	request := mustRouteRequestFromEnvelope(t, envelope)

	result, err := fixture.Dispatcher.Dispatch(context.Background(), request)
	if err != nil {
		t.Fatalf("Dispatch() error = %v, want nil", err)
	}

	responseEnvelope, err := BuildEnvelopeFromApplicationResult(result)
	if err != nil {
		t.Fatalf("BuildEnvelopeFromApplicationResult() error = %v, want nil", err)
	}
	if responseEnvelope.GetRequestId() != "request-1" || responseEnvelope.GetPayloadType() != "vibit.inventory.v1.GrantItemResponse" {
		t.Fatalf("response envelope = %#v, want correlated GrantItemResponse envelope", responseEnvelope)
	}

	responsePayload, err := DecodePayload(responseEnvelope.GetPayloadType(), responseEnvelope.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload() response error = %v, want nil", err)
	}
	response, ok := responsePayload.(*inventoryv1.GrantItemResponse)
	if !ok {
		t.Fatalf("response payload = %T, want *inventoryv1.GrantItemResponse", responsePayload)
	}
	if response.GetPlayerId() != "player-1" || response.GetItemId() != "item-1" || response.GetQuantity() != 2 || response.GetNewQuantity() != 2 {
		t.Fatalf("response = %#v, want mapped grant response", response)
	}

	if len(result.Events) != 1 {
		t.Fatalf("result events len = %d, want 1", len(result.Events))
	}
	eventEnvelope, err := BuildEnvelopeFromApplicationEvent(result.Events[0], result.RequestID, result.Target, result.Session)
	if err != nil {
		t.Fatalf("BuildEnvelopeFromApplicationEvent() error = %v, want nil", err)
	}
	if eventEnvelope.GetPayloadType() != "vibit.inventory.v1.ItemGranted" {
		t.Fatalf("event payload type = %q, want ItemGranted", eventEnvelope.GetPayloadType())
	}
}

func TestInventoryProtobufDomainBridgeBuildsPermissionErrorEnvelope(t *testing.T) {
	fixture := newInventoryRequestLoopFixture(t, false, true)

	envelope := mustBuildEnvelope(
		t,
		inventory.GrantItemRoute(),
		"request-1",
		app.Target{Scope: app.TargetScopePlayer, ID: "player-1"},
		app.Session{ConnectionID: "connection-1", SessionID: "session-1", PlayerID: "player-1"},
		&inventoryv1.GrantItemRequest{
			PlayerId:    "player-1",
			ItemId:      "item-1",
			Quantity:    2,
			Reason:      "reward",
			RequestedBy: "admin-1",
		},
	)

	request := mustRouteRequestFromEnvelope(t, envelope)

	result, err := fixture.Dispatcher.Dispatch(context.Background(), request)
	if err == nil {
		t.Fatal("Dispatch() error = nil, want permission error")
	}
	if result.Error == nil {
		t.Fatal("result error = nil, want application error")
	}

	errorEnvelope, err := BuildEnvelopeFromApplicationResult(result)
	if err != nil {
		t.Fatalf("BuildEnvelopeFromApplicationResult() error = %v, want nil", err)
	}
	if errorEnvelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_ERROR {
		t.Fatalf("error kind = %v, want MESSAGE_KIND_ERROR", errorEnvelope.GetKind())
	}
	if errorEnvelope.GetError().GetCode() != string(inventory.ErrorCodeInventoryPermission) {
		t.Fatalf("error code = %q, want %s", errorEnvelope.GetError().GetCode(), inventory.ErrorCodeInventoryPermission)
	}
	if errorEnvelope.GetError().GetRequestId() != "request-1" {
		t.Fatalf("error request_id = %q, want request-1", errorEnvelope.GetError().GetRequestId())
	}
	if errorEnvelope.GetPayloadType() != "" || len(errorEnvelope.GetPayload()) != 0 {
		t.Fatalf("error payload_type = %q payload len = %d, want no success payload", errorEnvelope.GetPayloadType(), len(errorEnvelope.GetPayload()))
	}
}
