package protobuf

import (
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	apppresence "github.com/iceiko/vibit/runtime/internal/app/presence"
	presencev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/presence/v1"
)

func TestRouteRequestWithPresencePayloadMapsGetPlayerPresenceRequest(t *testing.T) {
	request := mustRouteRequestForPresencePayload(&presencev1.GetPlayerPresenceRequest{
		PlayerId: "player-1",
	})

	mapped, err := RouteRequestWithDomainPayload(request)
	if err != nil {
		t.Fatalf("RouteRequestWithDomainPayload() error = %v, want nil", err)
	}
	payload, ok := mapped.Payload.(apppresence.GetPlayerPresenceRequest)
	if !ok {
		t.Fatalf("mapped Payload = %T, want GetPlayerPresenceRequest", mapped.Payload)
	}
	if payload.PlayerID != "player-1" {
		t.Fatalf("payload.PlayerID = %q, want player-1", payload.PlayerID)
	}
}

func TestProtoPayloadFromPresenceResultMapsGetPlayerPresenceResponse(t *testing.T) {
	lastSeenAt := time.Date(2026, 5, 20, 12, 0, 0, 1, time.FixedZone("test", 8*60*60))
	openedAt := lastSeenAt.Add(-time.Minute)
	boundAt := lastSeenAt.Add(-30 * time.Second)
	observedAt := lastSeenAt.Add(time.Second)

	payload, err := ProtoPayloadFromApplicationResult(applicationResultForPresence(apppresence.GetPlayerPresenceResult{
		PlayerID:          "player-1",
		Status:            apppresence.PresenceStatusOnline,
		ConnectionCount:   1,
		LastSeenAt:        &lastSeenAt,
		RuntimeSessionIDs: []string{"session-1"},
		ObservedAt:        observedAt,
		ActiveConnections: []apppresence.PresenceConnection{{
			ConnectionID:     "connection-1",
			ConnectionEpoch:  7,
			RuntimeSessionID: "session-1",
			LastSeenAt:       &lastSeenAt,
			OpenedAt:         openedAt,
			BoundAt:          &boundAt,
		}},
	}))
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult() error = %v, want nil", err)
	}
	response, ok := payload.(*presencev1.GetPlayerPresenceResponse)
	if !ok {
		t.Fatalf("payload = %T, want GetPlayerPresenceResponse", payload)
	}
	if response.GetPlayerId() != "player-1" ||
		response.GetPresenceStatus() != presencev1.PresenceStatus_PRESENCE_STATUS_ONLINE ||
		response.GetConnectionCount() != 1 ||
		len(response.GetActiveConnections()) != 1 ||
		len(response.GetRuntimeSessionIds()) != 1 ||
		response.GetRuntimeSessionIds()[0] != "session-1" {
		t.Fatalf("response = %#v, want mapped online presence", response)
	}
	active := response.GetActiveConnections()[0]
	if active.GetConnectionId() != "connection-1" ||
		active.GetConnectionEpoch() != 7 ||
		active.GetRuntimeSessionId() != "session-1" ||
		active.GetLastSeenAt() != lastSeenAt.UTC().Format(time.RFC3339Nano) ||
		active.GetOpenedAt() != openedAt.UTC().Format(time.RFC3339Nano) ||
		active.GetBoundAt() != boundAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("active connection = %#v, want mapped active metadata", active)
	}
	if response.GetLastSeenAt() != lastSeenAt.UTC().Format(time.RFC3339Nano) ||
		response.GetObservedAt() != observedAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("response times = %q/%q, want UTC RFC3339Nano", response.GetLastSeenAt(), response.GetObservedAt())
	}
}

func TestProtoPayloadFromPresenceResultMapsOfflineStatus(t *testing.T) {
	observedAt := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	payload, err := ProtoPayloadFromApplicationResult(applicationResultForPresence(apppresence.GetPlayerPresenceResult{
		PlayerID:   "player-1",
		Status:     apppresence.PresenceStatusOffline,
		ObservedAt: observedAt,
		LastSeenAt: nil,
	}))
	if err != nil {
		t.Fatalf("ProtoPayloadFromApplicationResult() error = %v, want nil", err)
	}
	response, ok := payload.(*presencev1.GetPlayerPresenceResponse)
	if !ok {
		t.Fatalf("payload = %T, want GetPlayerPresenceResponse", payload)
	}
	if response.GetPresenceStatus() != presencev1.PresenceStatus_PRESENCE_STATUS_OFFLINE ||
		response.GetConnectionCount() != 0 ||
		response.GetLastSeenAt() != "" ||
		len(response.GetActiveConnections()) != 0 {
		t.Fatalf("response = %#v, want offline presence without active connections", response)
	}
}

func mustRouteRequestForPresencePayload(payload *presencev1.GetPlayerPresenceRequest) app.RouteRequest {
	return app.RouteRequest{
		Route:   apppresence.GetPlayerPresenceRoute(),
		Payload: payload,
	}
}

func applicationResultForPresence(payload apppresence.GetPlayerPresenceResult) app.ApplicationResult {
	return app.ApplicationResult{
		Route:   apppresence.GetPlayerPresenceRoute(),
		Payload: payload,
	}
}
