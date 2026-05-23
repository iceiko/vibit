package realtime

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/app/connection"
)

func TestAcceptServerMessageResolvesPlayerCurrentConnections(t *testing.T) {
	registry := connection.NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-record-1")
	registerAndBind(t, registry, "connection-2", 1, "player-1", "session-2", "token-record-2")
	registerAndBind(t, registry, "connection-3", 1, "player-2", "session-3", "token-record-3")
	if _, err := registry.RegisterOpenConnection(context.Background(), connection.OpenConnection{
		ConnectionID:    "connection-4",
		ConnectionEpoch: 1,
	}); err != nil {
		t.Fatalf("RegisterOpenConnection(unbound) error = %v, want nil", err)
	}
	if _, err := registry.MarkConnectionClosed(context.Background(), connection.MarkClosed{
		ConnectionID:     "connection-2",
		ConnectionEpoch:  1,
		CloseReasonClass: "transport_closed",
	}); err != nil {
		t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
	}

	service := mustNewService(t, registry)
	payload := []byte(`{"title":"maintenance"}`)
	result, err := service.AcceptServerMessage(context.Background(), MessageIntent{
		Kind:           MessageIntentServerNotice,
		Target:         TargetPlayer(" player-1 "),
		SenderIdentity: serverIdentity(),
		PayloadType:    "runtime.realtime.ServerNotice",
		PayloadBytes:   payload,
		MessageID:      " notice-1 ",
	})
	if err != nil {
		t.Fatalf("AcceptServerMessage() error = %v, want nil", err)
	}

	if result.Outcome != DeliveryOutcomeAccepted ||
		result.IntentKind != MessageIntentServerNotice ||
		result.MessageID != "notice-1" ||
		result.PayloadType != "runtime.realtime.ServerNotice" ||
		len(result.Intents) != 1 {
		t.Fatalf("result = %#v, want one accepted intent", result)
	}
	intent := result.Intents[0]
	if intent.ConnectionID != "connection-1" ||
		intent.ConnectionEpoch != 1 ||
		intent.ActorKind != connection.ActorKindPlayer ||
		intent.PlayerID != "player-1" ||
		intent.RuntimeSessionID != "session-1" ||
		intent.AccessTokenRecordID != "token-record-1" ||
		string(intent.PayloadBytes) != `{"title":"maintenance"}` ||
		!intent.AcceptedAt.Equal(fixedTime()) {
		t.Fatalf("delivery intent = %#v, want active bound player connection", intent)
	}

	payload[0] = '['
	if string(result.Intents[0].PayloadBytes) != `{"title":"maintenance"}` {
		t.Fatalf("payload copy mutated to %q", string(result.Intents[0].PayloadBytes))
	}
}

func TestAcceptServerMessageTargetsSpecificBoundConnectionByIDAndEpoch(t *testing.T) {
	registry := connection.NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 7, "player-1", "session-1", "token-record-1")
	service := mustNewService(t, registry)

	result, err := service.AcceptServerMessage(context.Background(), MessageIntent{
		Kind:           MessageIntentDomainEventPush,
		Target:         TargetConnection(" connection-1 ", 7),
		SenderIdentity: serverIdentity(),
		PayloadType:    "inventory.ItemGranted",
		PayloadBytes:   []byte(`{"item":"sword"}`),
		MessageID:      "event-1",
	})
	if err != nil {
		t.Fatalf("AcceptServerMessage(connection) error = %v, want nil", err)
	}
	if result.Outcome != DeliveryOutcomeAccepted || len(result.Intents) != 1 {
		t.Fatalf("result = %#v, want one accepted connection intent", result)
	}
	if result.Intents[0].ConnectionID != "connection-1" || result.Intents[0].ConnectionEpoch != 7 {
		t.Fatalf("intent = %#v, want normalized connection target", result.Intents[0])
	}
}

func TestAcceptServerMessageRejectsMetadataOnlyOrPlayerSenderBeforeRecipientResolution(t *testing.T) {
	registry := &recordingRegistry{}
	service := mustNewService(t, registry)

	tests := []struct {
		name     string
		identity app.RequestIdentity
	}{
		{
			name: "metadata only player",
			identity: app.MetadataOnlyIdentityFromSession(app.Session{
				PlayerID: "player-1",
			}),
		},
		{
			name:     "validated player",
			identity: app.ValidatedPlayerIdentity("player-1", app.Session{PlayerID: "player-1"}),
		},
		{
			name: "empty service actor",
			identity: app.RequestIdentity{
				Status:    app.IdentityValidationValidated,
				ActorKind: app.ActorKindService,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.AcceptServerMessage(context.Background(), MessageIntent{
				Kind:           MessageIntentServerNotice,
				Target:         TargetPlayer("player-1"),
				SenderIdentity: tt.identity,
				PayloadType:    "runtime.realtime.ServerNotice",
				PayloadBytes:   []byte(`{}`),
			})
			assertServiceError(t, err, ErrorCodeRecipientForbidden, FailureClassForbidden)
			if result.Outcome != DeliveryOutcomeRecipientNotAuthorized ||
				result.FailureClass != FailureClassForbidden ||
				result.ErrorCode != ErrorCodeRecipientForbidden {
				t.Fatalf("result = %#v, want recipient-not-authorized", result)
			}
			if registry.findCalls != 0 || registry.listCalls != 0 {
				t.Fatalf("registry calls = find %d list %d, want none before identity rejection", registry.findCalls, registry.listCalls)
			}
		})
	}
}

func TestAcceptServerMessageReportsNoActiveRecipientWithoutSynthesizingTargets(t *testing.T) {
	registry := connection.NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), connection.OpenConnection{
		ConnectionID:    "connection-1",
		ConnectionEpoch: 1,
	}); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	service := mustNewService(t, registry)

	result, err := service.AcceptServerMessage(context.Background(), MessageIntent{
		Kind:           MessageIntentServerNotice,
		Target:         TargetConnection("connection-1", 1),
		SenderIdentity: serverIdentity(),
		PayloadType:    "runtime.realtime.ServerNotice",
		PayloadBytes:   []byte(`{}`),
	})
	assertServiceError(t, err, ErrorCodeNoActiveRecipient, FailureClassNoActiveRecipient)
	if result.Outcome != DeliveryOutcomeNoActiveRecipient || len(result.Intents) != 0 {
		t.Fatalf("result = %#v, want no active recipient and no delivery intents", result)
	}
}

func TestAcceptServerMessageFiltersUnboundPlayerCurrentConnections(t *testing.T) {
	registry := &recordingRegistry{
		listRecords: []connection.Record{
			{
				ConnectionID:    "connection-1",
				ConnectionEpoch: 1,
				State:           connection.StateOpenUnbound,
				PlayerID:        "player-1",
			},
			{
				ConnectionID:    "connection-2",
				ConnectionEpoch: 1,
				State:           connection.StateBound,
				ActorKind:       connection.ActorKindPlayer,
				PlayerID:        "player-1",
			},
		},
	}
	service := mustNewService(t, registry)

	result, err := service.AcceptServerMessage(context.Background(), MessageIntent{
		Kind:           MessageIntentServerNotice,
		Target:         TargetPlayer("player-1"),
		SenderIdentity: serverIdentity(),
		PayloadType:    "runtime.realtime.ServerNotice",
		PayloadBytes:   []byte(`{}`),
	})
	if err != nil {
		t.Fatalf("AcceptServerMessage() error = %v, want nil", err)
	}
	if len(result.Intents) != 1 || result.Intents[0].ConnectionID != "connection-2" {
		t.Fatalf("result intents = %#v, want only bound player connection", result.Intents)
	}
}

func TestAcceptServerMessageKeepsStreamSubscribersFutureOnly(t *testing.T) {
	registry := &recordingRegistry{}
	service := mustNewService(t, registry)

	result, err := service.AcceptServerMessage(context.Background(), MessageIntent{
		Kind:           MessageIntentStreamMessage,
		Target:         TargetStream("stream-1"),
		SenderIdentity: serverIdentity(),
		PayloadType:    "runtime.realtime.StreamMessage",
		PayloadBytes:   []byte(`{"text":"hello"}`),
	})
	assertServiceError(t, err, ErrorCodeRecipientForbidden, FailureClassForbidden)
	if result.Outcome != DeliveryOutcomeRecipientNotAuthorized ||
		result.Target.Kind != TargetStreamSubscribers ||
		len(result.Intents) != 0 {
		t.Fatalf("result = %#v, want stream target rejected without delivery", result)
	}
	if registry.findCalls != 0 || registry.listCalls != 0 {
		t.Fatalf("registry calls = find %d list %d, want none for deferred stream subscriptions", registry.findCalls, registry.listCalls)
	}
}

func TestAcceptServerMessageValidatesIntentBeforeRegistryResolution(t *testing.T) {
	registry := &recordingRegistry{}
	service := mustNewService(t, registry)

	tests := []struct {
		name   string
		intent MessageIntent
	}{
		{
			name: "unknown kind",
			intent: MessageIntent{
				Kind:           MessageIntentKind("chat_publish"),
				Target:         TargetPlayer("player-1"),
				SenderIdentity: serverIdentity(),
				PayloadType:    "runtime.realtime.ServerNotice",
				PayloadBytes:   []byte(`{}`),
			},
		},
		{
			name: "missing payload type",
			intent: MessageIntent{
				Kind:           MessageIntentServerNotice,
				Target:         TargetPlayer("player-1"),
				SenderIdentity: serverIdentity(),
				PayloadBytes:   []byte(`{}`),
			},
		},
		{
			name: "missing payload bytes",
			intent: MessageIntent{
				Kind:           MessageIntentServerNotice,
				Target:         TargetPlayer("player-1"),
				SenderIdentity: serverIdentity(),
				PayloadType:    "runtime.realtime.ServerNotice",
			},
		},
		{
			name: "missing target",
			intent: MessageIntent{
				Kind:           MessageIntentServerNotice,
				SenderIdentity: serverIdentity(),
				PayloadType:    "runtime.realtime.ServerNotice",
				PayloadBytes:   []byte(`{}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.AcceptServerMessage(context.Background(), tt.intent)
			assertServiceError(t, err, ErrorCodeIntentInvalid, FailureClassInvalidIntent)
			if result.Outcome != DeliveryOutcomePayloadInvalid || result.ErrorCode != ErrorCodeIntentInvalid {
				t.Fatalf("result = %#v, want payload-invalid intent rejection", result)
			}
			if registry.findCalls != 0 || registry.listCalls != 0 {
				t.Fatalf("registry calls = find %d list %d, want none before invalid intent rejection", registry.findCalls, registry.listCalls)
			}
		})
	}
}

func TestDeliveryResultsAreCopies(t *testing.T) {
	registry := connection.NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-record-1")
	service := mustNewService(t, registry)

	result, err := service.AcceptServerMessage(context.Background(), MessageIntent{
		Kind:           MessageIntentPresenceSignal,
		Target:         TargetPlayer("player-1"),
		SenderIdentity: serverIdentity(),
		PayloadType:    "presence.StatusChanged",
		PayloadBytes:   []byte(`{"status":"online"}`),
	})
	if err != nil {
		t.Fatalf("AcceptServerMessage() error = %v, want nil", err)
	}
	copied := copyDeliveryResult(result)
	copied.Intents[0].PlayerID = "mutated"
	copied.Intents[0].PayloadBytes[0] = '['
	if result.Intents[0].PlayerID != "player-1" || string(result.Intents[0].PayloadBytes) != `{"status":"online"}` {
		t.Fatalf("copyDeliveryResult aliased result = %#v copied = %#v", result, copied)
	}
}

func TestServiceErrorsAndDeliveryIntentShapeAreRedacted(t *testing.T) {
	err := serviceError(ErrorCodeDeliveryUnavailable, FailureClassDependencyUnavailable, errors.New("raw-access-token-secret"))
	if strings.Contains(err.Error(), "raw-access-token-secret") {
		t.Fatalf("ServiceError.Error() = %q, leaked wrapped detail", err.Error())
	}

	intentType := reflect.TypeOf(DeliveryIntent{})
	for _, forbidden := range []string{
		"AccessToken",
		"Credential",
		"RawToken",
		"RawCredential",
		"LookupDigest",
		"VerifierDigest",
		"VerifierKeyID",
		"Authorization",
		"Cookie",
		"QueryString",
		"Subprotocol",
		"RemoteAddr",
		"Socket",
		"Conn",
		"Writer",
		"CloseCode",
		"CloseReason",
		"Nakama",
		"Pitaya",
	} {
		if _, ok := intentType.FieldByName(forbidden); ok {
			t.Fatalf("DeliveryIntent contains forbidden proof, transport, or compatibility field %s", forbidden)
		}
	}
}

func TestNewServiceRequiresRegistry(t *testing.T) {
	_, err := NewService(ServiceDependencies{})
	assertServiceError(t, err, ErrorCodeDeliveryUnavailable, FailureClassDependencyUnavailable)
}

func registerAndBind(t *testing.T, registry *connection.InMemoryRegistry, connectionID string, epoch connection.ConnectionEpoch, playerID string, sessionID string, tokenRecordID string) {
	t.Helper()
	if _, err := registry.RegisterOpenConnection(context.Background(), connection.OpenConnection{
		ConnectionID:    connection.ConnectionID(connectionID),
		ConnectionEpoch: epoch,
	}); err != nil {
		t.Fatalf("RegisterOpenConnection(%s/%d) error = %v, want nil", connectionID, epoch, err)
	}
	if _, err := registry.BindConnectionIdentity(context.Background(), connection.BindIdentity{
		ConnectionID:        connection.ConnectionID(connectionID),
		ConnectionEpoch:     epoch,
		ActorKind:           connection.ActorKindPlayer,
		PlayerID:            connection.PlayerID(playerID),
		RuntimeSessionID:    connection.RuntimeSessionID(sessionID),
		AccessTokenRecordID: connection.AccessTokenRecordID(tokenRecordID),
	}); err != nil {
		t.Fatalf("BindConnectionIdentity(%s/%d) error = %v, want nil", connectionID, epoch, err)
	}
}

func mustNewService(t *testing.T, registry Registry) Service {
	t.Helper()
	service, err := NewService(ServiceDependencies{
		Registry: registry,
		Clock:    staticClock{now: fixedTime()},
	})
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}
	return service
}

func serverIdentity() app.RequestIdentity {
	return app.RequestIdentity{
		Status:    app.IdentityValidationValidated,
		ActorKind: app.ActorKindService,
		ActorID:   "runtime.realtime",
	}
}

func assertServiceError(t *testing.T, err error, code ErrorCode, class FailureClass) {
	t.Helper()
	var serviceErr *ServiceError
	if !errors.As(err, &serviceErr) {
		t.Fatalf("error = %T %v, want *ServiceError", err, err)
	}
	if serviceErr.Code != code || serviceErr.Class != class {
		t.Fatalf("ServiceError = %#v, want code %s class %s", serviceErr, code, class)
	}
}

type recordingRegistry struct {
	findCalls   int
	listCalls   int
	listRecords []connection.Record
}

func (r *recordingRegistry) FindConnectionByID(context.Context, connection.ConnectionID, connection.ConnectionEpoch) (connection.Record, bool) {
	r.findCalls++
	return connection.Record{}, false
}

func (r *recordingRegistry) ListConnectionsByPlayerID(context.Context, connection.PlayerID) []connection.Record {
	r.listCalls++
	return append([]connection.Record(nil), r.listRecords...)
}

type staticClock struct {
	now time.Time
}

func (c staticClock) Now() time.Time {
	return c.now
}

func fixedTime() time.Time {
	return time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC)
}
