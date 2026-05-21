package protobuf

import (
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	apppresence "github.com/iceiko/vibit/runtime/internal/app/presence"
	presencev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/presence/v1"
	"google.golang.org/protobuf/proto"
)

func routeRequestWithPresencePayload(request app.RouteRequest) (app.RouteRequest, bool, error) {
	switch request.Route {
	case apppresence.GetPlayerPresenceRoute():
		payload, ok := request.Payload.(*presencev1.GetPlayerPresenceRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.presence.v1.GetPlayerPresenceRequest")
		}

		request.Payload = apppresence.GetPlayerPresenceRequest{
			PlayerID: payload.GetPlayerId(),
		}
		return request, true, nil

	default:
		return request, false, nil
	}
}

func protoPayloadFromPresenceRoute(route app.RouteKey, payload any) (proto.Message, bool, error) {
	switch route {
	case apppresence.GetPlayerPresenceRoute():
		result, ok := payload.(apppresence.GetPlayerPresenceResult)
		if !ok {
			if pointerResult, pointerOK := payload.(*apppresence.GetPlayerPresenceResult); pointerOK && pointerResult != nil {
				result = *pointerResult
				ok = true
			}
		}
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be presence.GetPlayerPresenceResult")
		}

		activeConnections := make([]*presencev1.PresenceConnection, 0, len(result.ActiveConnections))
		for _, active := range result.ActiveConnections {
			activeConnections = append(activeConnections, &presencev1.PresenceConnection{
				ConnectionId:     active.ConnectionID,
				ConnectionEpoch:  active.ConnectionEpoch,
				RuntimeSessionId: active.RuntimeSessionID,
				LastSeenAt:       formatPresenceOptionalTime(active.LastSeenAt),
				OpenedAt:         formatPresenceTime(active.OpenedAt),
				BoundAt:          formatPresenceOptionalTime(active.BoundAt),
			})
		}

		return &presencev1.GetPlayerPresenceResponse{
			PlayerId:          result.PlayerID,
			PresenceStatus:    protoPresenceStatus(result.Status),
			ConnectionCount:   int32(result.ConnectionCount),
			ActiveConnections: activeConnections,
			RuntimeSessionIds: append([]string(nil), result.RuntimeSessionIDs...),
			LastSeenAt:        formatPresenceOptionalTime(result.LastSeenAt),
			ObservedAt:        formatPresenceTime(result.ObservedAt),
		}, true, nil

	default:
		return nil, false, nil
	}
}

func protoPresenceStatus(status apppresence.PresenceStatus) presencev1.PresenceStatus {
	switch status {
	case apppresence.PresenceStatusOnline:
		return presencev1.PresenceStatus_PRESENCE_STATUS_ONLINE
	case apppresence.PresenceStatusOffline:
		return presencev1.PresenceStatus_PRESENCE_STATUS_OFFLINE
	default:
		return presencev1.PresenceStatus_PRESENCE_STATUS_UNSPECIFIED
	}
}

func formatPresenceOptionalTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return formatPresenceTime(*value)
}

func formatPresenceTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}
