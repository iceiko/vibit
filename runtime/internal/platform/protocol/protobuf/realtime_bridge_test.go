package protobuf

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app/connection"
	"github.com/iceiko/vibit/runtime/internal/app/realtime"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	realtimev1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/realtime/v1"
	"google.golang.org/protobuf/proto"
)

func TestBuildRealtimeDeliveryEnvelopeMapsServerNoticeToSystemEnvelope(t *testing.T) {
	intent := realtimeBridgeIntent(realtime.MessageIntentServerNotice)
	intent.PayloadType = "runtime.realtime.ServerNotice"
	intent.PayloadBytes = []byte(`{"title":"maintenance"}`)

	envelope, err := BuildRealtimeDeliveryEnvelope(intent)
	if err != nil {
		t.Fatalf("BuildRealtimeDeliveryEnvelope() error = %v, want nil", err)
	}

	if envelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_SYSTEM ||
		envelope.GetModule() != "realtime" ||
		envelope.GetName() != "ServerNotice" ||
		envelope.GetRequestId() != "message-1" ||
		envelope.GetPayloadType() != "vibit.realtime.v1.ServerNotice" {
		t.Fatalf("envelope = %#v, want realtime ServerNotice system envelope", envelope)
	}
	if envelope.GetSession().GetConnectionId() != "ws-1" ||
		envelope.GetSession().GetConnectionEpoch() != 7 ||
		envelope.GetSession().GetPlayerId() != "player-1" ||
		envelope.GetSession().GetSessionId() != "session-1" {
		t.Fatalf("session = %#v, want server-observed connection/session metadata", envelope.GetSession())
	}
	if envelope.GetTarget().GetScope() != protocolv1.TargetScope_TARGET_SCOPE_PLAYER ||
		envelope.GetTarget().GetId() != "player-1" {
		t.Fatalf("target = %#v, want player target", envelope.GetTarget())
	}

	payload := mustRealtimePayload[*realtimev1.ServerNotice](t, envelope)
	if payload.GetMessageId() != "message-1" ||
		payload.GetPayloadType() != "runtime.realtime.ServerNotice" ||
		string(payload.GetPayload()) != `{"title":"maintenance"}` ||
		payload.GetAcceptedAt() != realtimeBridgeTime().Format(time.RFC3339Nano) {
		t.Fatalf("payload = %#v, want copied ServerNotice payload", payload)
	}

	intent.PayloadBytes[0] = '['
	if string(payload.GetPayload()) != `{"title":"maintenance"}` {
		t.Fatalf("payload aliased intent bytes, got %q", string(payload.GetPayload()))
	}
}

func TestBuildRealtimeDeliveryFrameMapsDomainEventAndPresenceSignal(t *testing.T) {
	tests := []struct {
		name            string
		kind            realtime.MessageIntentKind
		wantRouteName   string
		wantPayloadType string
		decode          func(t *testing.T, envelope *protocolv1.Envelope)
	}{
		{
			name:            "domain event",
			kind:            realtime.MessageIntentDomainEventPush,
			wantRouteName:   "DomainEventPush",
			wantPayloadType: "vibit.realtime.v1.DomainEventPush",
			decode: func(t *testing.T, envelope *protocolv1.Envelope) {
				t.Helper()
				payload := mustRealtimePayload[*realtimev1.DomainEventPush](t, envelope)
				if payload.GetSourcePayloadType() != "inventory.ItemGranted" || string(payload.GetPayload()) != `{"item":"sword"}` {
					t.Fatalf("payload = %#v, want domain event bytes", payload)
				}
			},
		},
		{
			name:            "presence signal",
			kind:            realtime.MessageIntentPresenceSignal,
			wantRouteName:   "PresenceSignal",
			wantPayloadType: "vibit.realtime.v1.PresenceSignal",
			decode: func(t *testing.T, envelope *protocolv1.Envelope) {
				t.Helper()
				payload := mustRealtimePayload[*realtimev1.PresenceSignal](t, envelope)
				if payload.GetSourcePayloadType() != "presence.StatusChanged" || string(payload.GetPayload()) != `{"status":"online"}` {
					t.Fatalf("payload = %#v, want presence signal bytes", payload)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			intent := realtimeBridgeIntent(tc.kind)
			if tc.kind == realtime.MessageIntentDomainEventPush {
				intent.PayloadType = "inventory.ItemGranted"
				intent.PayloadBytes = []byte(`{"item":"sword"}`)
			} else {
				intent.PayloadType = "presence.StatusChanged"
				intent.PayloadBytes = []byte(`{"status":"online"}`)
			}
			frame, err := BuildRealtimeDeliveryFrame(intent)
			if err != nil {
				t.Fatalf("BuildRealtimeDeliveryFrame() error = %v, want nil", err)
			}
			if frame.ConnectionID != "ws-1" || frame.ConnectionEpoch != 7 || len(frame.Payload) == 0 {
				t.Fatalf("frame = %#v, want encoded ws-1/7 payload", frame)
			}

			var envelope protocolv1.Envelope
			if err := proto.Unmarshal(frame.Payload, &envelope); err != nil {
				t.Fatalf("proto.Unmarshal(frame payload) error = %v, want nil", err)
			}
			if envelope.GetKind() != protocolv1.MessageKind_MESSAGE_KIND_EVENT ||
				envelope.GetModule() != "realtime" ||
				envelope.GetName() != tc.wantRouteName ||
				envelope.GetPayloadType() != tc.wantPayloadType {
				t.Fatalf("envelope = %#v, want realtime event %s", &envelope, tc.wantRouteName)
			}
			tc.decode(t, &envelope)
		})
	}
}

func TestBuildRealtimeDeliveryEnvelopeRejectsUnsupportedOrMalformedIntent(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*realtime.DeliveryIntent)
	}{
		{
			name: "stream remains future only",
			mutate: func(intent *realtime.DeliveryIntent) {
				intent.IntentKind = realtime.MessageIntentStreamMessage
			},
		},
		{
			name: "unknown kind",
			mutate: func(intent *realtime.DeliveryIntent) {
				intent.IntentKind = realtime.MessageIntentKind("chat_message")
			},
		},
		{
			name: "missing connection id",
			mutate: func(intent *realtime.DeliveryIntent) {
				intent.ConnectionID = " "
			},
		},
		{
			name: "missing epoch",
			mutate: func(intent *realtime.DeliveryIntent) {
				intent.ConnectionEpoch = 0
			},
		},
		{
			name: "missing payload type",
			mutate: func(intent *realtime.DeliveryIntent) {
				intent.PayloadType = " "
			},
		},
		{
			name: "missing payload bytes",
			mutate: func(intent *realtime.DeliveryIntent) {
				intent.PayloadBytes = nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			intent := realtimeBridgeIntent(realtime.MessageIntentServerNotice)
			tc.mutate(&intent)
			_, err := BuildRealtimeDeliveryEnvelope(intent)
			var bridgeErr *PayloadBridgeError
			if !errors.As(err, &bridgeErr) {
				t.Fatalf("BuildRealtimeDeliveryEnvelope() error = %T %v, want *PayloadBridgeError", err, err)
			}
			if strings.Contains(bridgeErr.Error(), "raw-access-token") || strings.Contains(bridgeErr.Error(), "verifier") {
				t.Fatalf("bridge error leaks secret detail: %q", bridgeErr.Error())
			}
		})
	}
}

func TestRealtimeProtoPayloadShapeOmitsCredentialAndDeliveryGuaranteeFields(t *testing.T) {
	for _, message := range []proto.Message{
		&realtimev1.ServerNotice{},
		&realtimev1.DomainEventPush{},
		&realtimev1.PresenceSignal{},
	} {
		descriptor := message.ProtoReflect().Descriptor()
		fields := descriptor.Fields()
		fieldNames := map[string]bool{}
		for i := 0; i < fields.Len(); i++ {
			fieldNames[string(fields.Get(i).Name())] = true
		}
		for _, forbidden := range []string{
			"access_token",
			"access_token_record_id",
			"credential",
			"lookup_digest",
			"verifier_digest",
			"verifier_key_id",
			"authorization",
			"cookie",
			"ack",
			"retry",
			"offset",
			"sequence",
			"stream_id",
			"room_id",
			"group_id",
		} {
			if fieldNames[forbidden] {
				t.Fatalf("%s has forbidden field %q", descriptor.FullName(), forbidden)
			}
		}
	}
}

func TestRealtimeDeliveryFrameShapeIsTransportNeutral(t *testing.T) {
	typ := reflect.TypeOf(RealtimeDeliveryFrame{})
	for _, forbidden := range []string{
		"PlayerID",
		"RuntimeSessionID",
		"AccessToken",
		"AccessTokenRecordID",
		"Credential",
		"PayloadType",
		"MessageKind",
		"WebSocket",
		"Conn",
		"RemoteAddr",
		"Authorization",
		"Cookie",
		"QueryString",
		"Subprotocol",
		"Nakama",
		"Pitaya",
	} {
		if _, ok := typ.FieldByName(forbidden); ok {
			t.Fatalf("RealtimeDeliveryFrame contains forbidden policy, proof, or compatibility field %s", forbidden)
		}
	}
}

func realtimeBridgeIntent(kind realtime.MessageIntentKind) realtime.DeliveryIntent {
	return realtime.DeliveryIntent{
		IntentKind:          kind,
		ConnectionID:        " ws-1 ",
		ConnectionEpoch:     7,
		ActorKind:           connection.ActorKindPlayer,
		PlayerID:            " player-1 ",
		RuntimeSessionID:    " session-1 ",
		AccessTokenRecordID: "token-record-1",
		MessageID:           " message-1 ",
		PayloadType:         "runtime.realtime.ServerNotice",
		PayloadBytes:        []byte(`{"ok":true}`),
		AcceptedAt:          realtimeBridgeTime(),
	}
}

func realtimeBridgeTime() time.Time {
	return time.Date(2026, 5, 23, 12, 0, 0, 123, time.FixedZone("test", 8*60*60)).UTC()
}

func mustRealtimePayload[T proto.Message](t *testing.T, envelope *protocolv1.Envelope) T {
	t.Helper()
	decoded, err := DecodePayload(envelope.GetPayloadType(), envelope.GetPayload())
	if err != nil {
		t.Fatalf("DecodePayload(%q) error = %v, want nil", envelope.GetPayloadType(), err)
	}
	payload, ok := decoded.(T)
	if !ok {
		var zero T
		t.Fatalf("decoded payload = %T, want %T", decoded, zero)
	}
	return payload
}
