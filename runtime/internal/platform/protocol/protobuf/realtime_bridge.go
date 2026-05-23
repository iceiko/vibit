package protobuf

import (
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/app/realtime"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	realtimev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/realtime/v1"
	"google.golang.org/protobuf/proto"
)

const realtimeModuleName = "realtime"

type RealtimeDeliveryFrame struct {
	ConnectionID    string
	ConnectionEpoch uint64
	Payload         []byte
}

func BuildRealtimeDeliveryFrame(intent realtime.DeliveryIntent) (RealtimeDeliveryFrame, error) {
	envelope, err := BuildRealtimeDeliveryEnvelope(intent)
	if err != nil {
		return RealtimeDeliveryFrame{}, err
	}
	payload, err := proto.Marshal(envelope)
	if err != nil {
		return RealtimeDeliveryFrame{}, err
	}
	return RealtimeDeliveryFrame{
		ConnectionID:    strings.TrimSpace(string(intent.ConnectionID)),
		ConnectionEpoch: uint64(intent.ConnectionEpoch),
		Payload:         payload,
	}, nil
}

func BuildRealtimeDeliveryEnvelope(intent realtime.DeliveryIntent) (*protocolv1.Envelope, error) {
	payload, route, err := realtimeProtoPayloadAndRoute(intent)
	if err != nil {
		return nil, err
	}
	target := realtimeTarget(intent)
	session := app.Session{
		ConnectionID:    strings.TrimSpace(string(intent.ConnectionID)),
		SessionID:       strings.TrimSpace(string(intent.RuntimeSessionID)),
		PlayerID:        strings.TrimSpace(string(intent.PlayerID)),
		ConnectionEpoch: uint64(intent.ConnectionEpoch),
	}
	return BuildEnvelope(route, strings.TrimSpace(intent.MessageID), target, session, payload)
}

func realtimeProtoPayloadAndRoute(intent realtime.DeliveryIntent) (proto.Message, app.RouteKey, error) {
	if err := validateRealtimeDeliveryIntent(intent); err != nil {
		return nil, app.RouteKey{}, err
	}

	switch intent.IntentKind {
	case realtime.MessageIntentServerNotice:
		return &realtimev1.ServerNotice{
			MessageId:   strings.TrimSpace(intent.MessageID),
			PayloadType: strings.TrimSpace(intent.PayloadType),
			Payload:     append([]byte(nil), intent.PayloadBytes...),
			AcceptedAt:  formatRealtimeDeliveryTime(intent.AcceptedAt),
		}, app.RouteKey{Kind: app.MessageKindSystem, Module: realtimeModuleName, Name: "ServerNotice"}, nil
	case realtime.MessageIntentDomainEventPush:
		return &realtimev1.DomainEventPush{
			MessageId:         strings.TrimSpace(intent.MessageID),
			SourcePayloadType: strings.TrimSpace(intent.PayloadType),
			Payload:           append([]byte(nil), intent.PayloadBytes...),
			AcceptedAt:        formatRealtimeDeliveryTime(intent.AcceptedAt),
		}, app.RouteKey{Kind: app.MessageKindEvent, Module: realtimeModuleName, Name: "DomainEventPush"}, nil
	case realtime.MessageIntentPresenceSignal:
		return &realtimev1.PresenceSignal{
			MessageId:         strings.TrimSpace(intent.MessageID),
			SourcePayloadType: strings.TrimSpace(intent.PayloadType),
			Payload:           append([]byte(nil), intent.PayloadBytes...),
			AcceptedAt:        formatRealtimeDeliveryTime(intent.AcceptedAt),
		}, app.RouteKey{Kind: app.MessageKindEvent, Module: realtimeModuleName, Name: "PresenceSignal"}, nil
	default:
		return nil, app.RouteKey{}, &PayloadBridgeError{Message: "realtime delivery intent kind is unsupported"}
	}
}

func validateRealtimeDeliveryIntent(intent realtime.DeliveryIntent) error {
	if strings.TrimSpace(string(intent.ConnectionID)) == "" || intent.ConnectionEpoch == 0 {
		return &PayloadBridgeError{Message: "realtime delivery target is invalid"}
	}
	if strings.TrimSpace(intent.PayloadType) == "" || len(intent.PayloadBytes) == 0 {
		return &PayloadBridgeError{Message: "realtime delivery payload is invalid"}
	}
	if intent.IntentKind == realtime.MessageIntentStreamMessage {
		return &PayloadBridgeError{Message: "realtime stream delivery is not supported"}
	}
	return nil
}

func realtimeTarget(intent realtime.DeliveryIntent) app.Target {
	playerID := strings.TrimSpace(string(intent.PlayerID))
	if playerID != "" {
		return app.Target{Scope: app.TargetScopePlayer, ID: playerID}
	}
	return app.Target{Scope: app.TargetScopeSystem, ID: "realtime"}
}

func formatRealtimeDeliveryTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}
