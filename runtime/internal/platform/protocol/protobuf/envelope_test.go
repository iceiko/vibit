package protobuf

import (
	"testing"

	"github.com/iceiko/vibit/runtime/internal/app"
	inventoryv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/inventory/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"google.golang.org/protobuf/proto"
)

func TestBuildEnvelopeAndRouteRequestFromEnvelope(t *testing.T) {
	payload := &inventoryv1.GrantItemRequest{
		PlayerId:    "player-1",
		ItemId:      "item-1",
		Quantity:    3,
		Reason:      "test",
		RequestedBy: "agent",
	}

	envelope, err := BuildEnvelope(
		app.RouteKey{
			Kind:   app.MessageKindCommand,
			Module: " inventory ",
			Name:   " GrantItem ",
		},
		" request-1 ",
		app.Target{
			Scope: app.TargetScopePlayer,
			ID:    " player-1 ",
		},
		app.Session{
			ConnectionID:    " connection-1 ",
			SessionID:       " session-1 ",
			PlayerID:        " player-1 ",
			ConnectionEpoch: 7,
		},
		payload,
	)
	if err != nil {
		t.Fatalf("BuildEnvelope() error = %v", err)
	}

	if envelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_COMMAND {
		t.Fatalf("envelope kind = %v, want command", envelope.GetKind())
	}
	if envelope.GetPayloadType() != "vibit.inventory.v1.GrantItemRequest" {
		t.Fatalf("payload_type = %q, want generated message full name", envelope.GetPayloadType())
	}

	routeRequest, err := RouteRequestFromEnvelope(envelope)
	if err != nil {
		t.Fatalf("RouteRequestFromEnvelope() error = %v", err)
	}

	if routeRequest.RequestID != "request-1" {
		t.Fatalf("RequestID = %q, want %q", routeRequest.RequestID, "request-1")
	}
	if routeRequest.Route.Kind != app.MessageKindCommand {
		t.Fatalf("Route.Kind = %q, want command", routeRequest.Route.Kind)
	}
	if routeRequest.Route.Module != "inventory" || routeRequest.Route.Name != "GrantItem" {
		t.Fatalf("Route = %#v, want inventory GrantItem", routeRequest.Route)
	}
	if routeRequest.Target.Scope != app.TargetScopePlayer || routeRequest.Target.ID != "player-1" {
		t.Fatalf("Target = %#v, want player target", routeRequest.Target)
	}
	if routeRequest.Session.ConnectionID != "connection-1" || routeRequest.Session.SessionID != "session-1" || routeRequest.Session.PlayerID != "player-1" || routeRequest.Session.ConnectionEpoch != 7 {
		t.Fatalf("Session = %#v, want normalized session", routeRequest.Session)
	}
	if routeRequest.Identity.Status != app.IdentityValidationMetadataOnly {
		t.Fatalf("Identity.Status = %q, want %q", routeRequest.Identity.Status, app.IdentityValidationMetadataOnly)
	}
	if routeRequest.Identity.PlayerID != "player-1" || routeRequest.Identity.SessionID != "session-1" || routeRequest.Identity.ConnectionID != "connection-1" {
		t.Fatalf("Identity = %#v, want normalized metadata-only session identity", routeRequest.Identity)
	}
	if routeRequest.Identity.PlayerIDValidated || routeRequest.Identity.SessionValidated {
		t.Fatalf("Identity validation flags = %#v, want metadata-only identity", routeRequest.Identity)
	}

	var decoded inventoryv1.GrantItemRequest
	if err := proto.Unmarshal(routeRequest.PayloadBytes, &decoded); err != nil {
		t.Fatalf("proto.Unmarshal payload error = %v", err)
	}
	if decoded.GetPlayerId() != "player-1" || decoded.GetItemId() != "item-1" || decoded.GetQuantity() != 3 {
		t.Fatalf("decoded payload = %#v, want original grant request", &decoded)
	}

	decodedPayload, ok := routeRequest.Payload.(*inventoryv1.GrantItemRequest)
	if !ok {
		t.Fatalf("Payload = %T, want *inventoryv1.GrantItemRequest", routeRequest.Payload)
	}
	if decodedPayload.GetPlayerId() != "player-1" || decodedPayload.GetItemId() != "item-1" || decodedPayload.GetQuantity() != 3 {
		t.Fatalf("Payload = %#v, want decoded grant request", decodedPayload)
	}
}

func TestRouteRequestFromEnvelopeRejectsMissingPayloadType(t *testing.T) {
	_, err := RouteRequestFromEnvelope(&protocolv1.Envelope{
		Kind:    protocolv1.MessageKind_MESSAGE_KIND_COMMAND,
		Module:  "inventory",
		Name:    "GrantItem",
		Target:  &protocolv1.Target{Scope: protocolv1.TargetScope_TARGET_SCOPE_PLAYER},
		Session: &protocolv1.Session{},
	})
	if err == nil {
		t.Fatal("RouteRequestFromEnvelope() error = nil, want validation error")
	}
}

func TestRouteRequestFromEnvelopeRejectsUnknownPayloadType(t *testing.T) {
	_, err := RouteRequestFromEnvelope(&protocolv1.Envelope{
		Kind:        protocolv1.MessageKind_MESSAGE_KIND_COMMAND,
		Module:      "inventory",
		Name:        "GrantItem",
		Target:      &protocolv1.Target{Scope: protocolv1.TargetScope_TARGET_SCOPE_PLAYER},
		Session:     &protocolv1.Session{},
		PayloadType: "vibit.inventory.v1.DoesNotExist",
		Payload:     []byte{},
	})
	if err == nil {
		t.Fatal("RouteRequestFromEnvelope() error = nil, want validation error")
	}
}

func TestBuildEnvelopeRejectsEmptyRoute(t *testing.T) {
	_, err := BuildEnvelope(
		app.RouteKey{Kind: app.MessageKindCommand, Module: "", Name: "GrantItem"},
		"request-1",
		app.Target{},
		app.Session{},
		&inventoryv1.GrantItemRequest{},
	)
	if err == nil {
		t.Fatal("BuildEnvelope() error = nil, want validation error")
	}
}

func TestBuildEnvelopeRejectsUnsupportedKind(t *testing.T) {
	_, err := BuildEnvelope(
		app.RouteKey{Kind: "", Module: "inventory", Name: "GrantItem"},
		"request-1",
		app.Target{},
		app.Session{},
		&inventoryv1.GrantItemRequest{},
	)
	if err == nil {
		t.Fatal("BuildEnvelope() error = nil, want validation error")
	}
}

func TestGeneratedInventoryPayloadRoundTrip(t *testing.T) {
	source := &inventoryv1.GetInventoryRequest{
		PlayerId:    "player-1",
		RequestedBy: "agent",
	}

	encoded, err := proto.Marshal(source)
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}

	var decoded inventoryv1.GetInventoryRequest
	if err := proto.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("proto.Unmarshal() error = %v", err)
	}

	if decoded.GetPlayerId() != source.GetPlayerId() || decoded.GetRequestedBy() != source.GetRequestedBy() {
		t.Fatalf("decoded = %#v, want %#v", &decoded, source)
	}
}
