package protobuf

import (
	"fmt"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	inventoryv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/inventory/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"google.golang.org/protobuf/proto"
)

type PayloadBridgeError struct {
	Message string
}

func (e *PayloadBridgeError) Error() string {
	return e.Message
}

func RouteRequestFromEnvelopeForDispatch(envelope *protocolv1.Envelope) (app.RouteRequest, error) {
	request, err := RouteRequestFromEnvelope(envelope)
	if err != nil {
		return app.RouteRequest{}, err
	}
	return RouteRequestWithDomainPayload(request)
}

func RouteRequestWithDomainPayload(request app.RouteRequest) (app.RouteRequest, error) {
	switch request.Route {
	case inventory.GrantItemRoute():
		payload, ok := request.Payload.(*inventoryv1.GrantItemRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, payloadBridgeError(request.Route, "payload must be vibit.inventory.v1.GrantItemRequest")
		}
		request.Payload = inventory.GrantItemRequest{
			PlayerID:    payload.GetPlayerId(),
			ItemID:      payload.GetItemId(),
			Quantity:    payload.GetQuantity(),
			Reason:      payload.GetReason(),
			RequestedBy: payload.GetRequestedBy(),
		}
		return request, nil

	case inventory.GetInventoryRoute():
		payload, ok := request.Payload.(*inventoryv1.GetInventoryRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, payloadBridgeError(request.Route, "payload must be vibit.inventory.v1.GetInventoryRequest")
		}
		request.Payload = inventory.GetInventoryRequest{
			PlayerID:    payload.GetPlayerId(),
			RequestedBy: payload.GetRequestedBy(),
		}
		return request, nil

	default:
		return request, nil
	}
}

func BuildEnvelopeFromApplicationResult(result app.ApplicationResult) (*protocolv1.Envelope, error) {
	payload, err := ProtoPayloadFromApplicationResult(result)
	if err != nil {
		return nil, err
	}
	return BuildEnvelope(result.Route, result.RequestID, result.Target, result.Session, payload)
}

func BuildEnvelopeFromApplicationEvent(event app.ApplicationEvent, requestID string, target app.Target, session app.Session) (*protocolv1.Envelope, error) {
	payload, err := ProtoPayloadFromApplicationEvent(event)
	if err != nil {
		return nil, err
	}
	return BuildEnvelope(event.Route, requestID, target, session, payload)
}

func ProtoPayloadFromApplicationResult(result app.ApplicationResult) (proto.Message, error) {
	if result.Error != nil {
		return nil, &PayloadBridgeError{Message: "application error result payload mapping is not implemented"}
	}
	return protoPayloadFromRouteAndPayload(result.Route, result.Payload)
}

func ProtoPayloadFromApplicationEvent(event app.ApplicationEvent) (proto.Message, error) {
	return protoPayloadFromRouteAndPayload(event.Route, event.Payload)
}

func protoPayloadFromRouteAndPayload(route app.RouteKey, payload any) (proto.Message, error) {
	switch route {
	case inventory.GrantItemRoute():
		response, ok := payload.(inventory.GrantItemResponse)
		if !ok {
			if pointerResponse, pointerOK := payload.(*inventory.GrantItemResponse); pointerOK && pointerResponse != nil {
				response = *pointerResponse
				ok = true
			}
		}
		if !ok {
			return nil, payloadBridgeError(route, "payload must be inventory.GrantItemResponse")
		}
		return &inventoryv1.GrantItemResponse{
			PlayerId:    response.PlayerID,
			ItemId:      response.ItemID,
			Quantity:    response.Quantity,
			NewQuantity: response.NewQuantity,
			Event:       response.Event,
		}, nil

	case inventory.GetInventoryRoute():
		response, ok := payload.(inventory.GetInventoryResponse)
		if !ok {
			if pointerResponse, pointerOK := payload.(*inventory.GetInventoryResponse); pointerOK && pointerResponse != nil {
				response = *pointerResponse
				ok = true
			}
		}
		if !ok {
			return nil, payloadBridgeError(route, "payload must be inventory.GetInventoryResponse")
		}
		items := make([]*inventoryv1.InventoryItem, 0, len(response.Items))
		for _, item := range response.Items {
			items = append(items, &inventoryv1.InventoryItem{
				ItemId:   item.ItemID,
				Quantity: item.Quantity,
			})
		}
		return &inventoryv1.GetInventoryResponse{
			PlayerId: response.PlayerID,
			Items:    items,
		}, nil

	case app.RouteKey{Kind: app.MessageKindEvent, Module: inventory.ModuleName, Name: inventory.EventItemGranted}:
		event, ok := payload.(inventory.ItemGrantedEvent)
		if !ok {
			if pointerEvent, pointerOK := payload.(*inventory.ItemGrantedEvent); pointerOK && pointerEvent != nil {
				event = *pointerEvent
				ok = true
			}
		}
		if !ok {
			return nil, payloadBridgeError(route, "payload must be inventory.ItemGrantedEvent")
		}
		return &inventoryv1.ItemGranted{
			EventId:     event.EventID,
			OccurredAt:  formatEventTime(event.OccurredAt),
			PlayerId:    event.PlayerID,
			ItemId:      event.ItemID,
			Quantity:    event.Quantity,
			NewQuantity: event.NewQuantity,
			Reason:      event.Reason,
		}, nil

	default:
		protoPayload, ok := payload.(proto.Message)
		if !ok || protoPayload == nil {
			return nil, payloadBridgeError(route, "payload has no protocol bridge")
		}
		return protoPayload, nil
	}
}

func formatEventTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func payloadBridgeError(route app.RouteKey, message string) error {
	renderedRoute := app.RenderRouteKey(route)
	if renderedRoute == "" {
		renderedRoute = "<invalid>"
	}
	return &PayloadBridgeError{Message: fmt.Sprintf("%s: %s", renderedRoute, message)}
}
