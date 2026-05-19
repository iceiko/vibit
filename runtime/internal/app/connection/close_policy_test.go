package connection

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClosePolicyInvalidatesBoundConnectionByIDAndEpochWithoutSocketClose(t *testing.T) {
	now := fixedTime().Add(time.Hour)
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 7, "player-1", "session-1", "token-record-1")
	policy := NewClosePolicy(registry, WithClosePolicyClock(staticClock{now: now}))

	result, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:           TargetConnection(" connection-1 ", 7),
		ReasonClass:      CloseReasonLogoutPresentedToken,
		PublicVisibility: ClosePublicVisibilityGenericReauthRequired,
		Retryability:     CloseRetryabilityNotRetryable,
	})
	if err != nil {
		t.Fatalf("RequestClose() error = %v, want nil", err)
	}

	if result.TransportAction != CloseTransportActionMarkInvalidatedOnly ||
		result.ReasonClass != CloseReasonLogoutPresentedToken ||
		result.PublicVisibility != ClosePublicVisibilityGenericReauthRequired ||
		result.Retryability != CloseRetryabilityNotRetryable ||
		!result.RequestedAt.Equal(now) {
		t.Fatalf("result metadata = %#v, want redacted mark-invalidated close intent metadata", result)
	}
	if len(result.Intents) != 1 || len(result.Skipped) != 0 {
		t.Fatalf("result intents/skipped = %#v/%#v, want one intent and no skipped records", result.Intents, result.Skipped)
	}

	intent := result.Intents[0]
	if intent.ConnectionID != "connection-1" ||
		intent.ConnectionEpoch != 7 ||
		intent.TargetKind != CloseTargetConnectionID ||
		intent.ActorKind != ActorKindPlayer ||
		intent.PlayerID != "player-1" ||
		intent.RuntimeSessionID != "session-1" ||
		intent.AccessTokenRecordID != "token-record-1" ||
		intent.ReasonClass != CloseReasonLogoutPresentedToken ||
		intent.TransportAction != CloseTransportActionMarkInvalidatedOnly ||
		intent.Outcome != CloseOutcomeInvalidated ||
		!intent.RecordedAt.Equal(now) {
		t.Fatalf("intent = %#v, want invalidated bound connection intent", intent)
	}

	record, ok := registry.FindConnectionByID(context.Background(), "connection-1", 7)
	if !ok {
		t.Fatal("FindConnectionByID() ok = false, want true")
	}
	if record.State != StateInvalidated ||
		record.InvalidatedAt == nil ||
		!record.InvalidatedAt.Equal(now) ||
		record.InvalidationClass != string(CloseReasonLogoutPresentedToken) ||
		record.ClosedAt != nil ||
		record.CloseReasonClass != "" {
		t.Fatalf("record = %#v, want invalidated registry state without concrete socket close", record)
	}
}

func TestClosePolicyTargetsOnlyActiveBoundRegistryRecords(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-record-1")
	registerAndBind(t, registry, "connection-2", 1, "player-1", "session-2", "token-record-2")
	registerAndBind(t, registry, "connection-3", 1, "player-2", "session-3", "token-record-3")
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-4", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() unbound error = %v, want nil", err)
	}
	if _, err := registry.MarkConnectionClosed(context.Background(), MarkClosed{
		ConnectionID:     "connection-2",
		ConnectionEpoch:  1,
		CloseReasonClass: "transport_observed",
	}); err != nil {
		t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
	}

	policy := NewClosePolicy(registry, WithClosePolicyClock(staticClock{now: fixedTime().Add(time.Minute)}))
	result, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetPlayer(" player-1 "),
		ReasonClass: CloseReasonTokenRevoked,
	})
	if err != nil {
		t.Fatalf("RequestClose() error = %v, want nil", err)
	}

	assertCloseIntentIDs(t, result.Intents, []ConnectionID{"connection-1"})
	if len(result.Skipped) != 0 {
		t.Fatalf("Skipped = %#v, want none for list target that only returns active bound records", result.Skipped)
	}

	if record, _ := registry.FindConnectionByID(context.Background(), "connection-1", 1); record.State != StateInvalidated {
		t.Fatalf("connection-1 state = %s, want invalidated", record.State)
	}
	if record, _ := registry.FindConnectionByID(context.Background(), "connection-2", 1); record.State != StateClosed {
		t.Fatalf("connection-2 state = %s, want closed untouched", record.State)
	}
	if record, _ := registry.FindConnectionByID(context.Background(), "connection-4", 1); record.State != StateOpenUnbound {
		t.Fatalf("connection-4 state = %s, want unbound untouched", record.State)
	}
}

func TestClosePolicyTargetsRuntimeSessionAndTokenRecordThroughRegistryOnly(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-record-1")
	registerAndBind(t, registry, "connection-2", 1, "player-1", "session-2", "token-record-2")
	policy := NewClosePolicy(registry, WithClosePolicyClock(staticClock{now: fixedTime().Add(time.Minute)}))

	sessionResult, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetRuntimeSession(" session-1 "),
		ReasonClass: CloseReasonSessionRevoked,
	})
	if err != nil {
		t.Fatalf("RequestClose(session) error = %v, want nil", err)
	}
	assertCloseIntentIDs(t, sessionResult.Intents, []ConnectionID{"connection-1"})

	tokenResult, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetAccessTokenRecord(" token-record-2 "),
		ReasonClass: CloseReasonLogoutPresentedToken,
	})
	if err != nil {
		t.Fatalf("RequestClose(token) error = %v, want nil", err)
	}
	assertCloseIntentIDs(t, tokenResult.Intents, []ConnectionID{"connection-2"})

	missingResult, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetRuntimeSession("session-from-client-metadata-only"),
		ReasonClass: CloseReasonSessionRevoked,
	})
	if err != nil {
		t.Fatalf("RequestClose(missing session) error = %v, want nil", err)
	}
	if len(missingResult.Intents) != 0 || len(missingResult.Skipped) != 0 {
		t.Fatalf("missingResult = %#v, want no synthesized intents from metadata-only target", missingResult)
	}
}

func TestClosePolicySkipsConnectionTargetWhenRecordIsNotBoundActive(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	policy := NewClosePolicy(registry, WithClosePolicyClock(staticClock{now: fixedTime().Add(time.Minute)}))

	result, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetConnection("connection-1", 1),
		ReasonClass: CloseReasonAdministrativeAction,
	})
	if err != nil {
		t.Fatalf("RequestClose() error = %v, want nil", err)
	}
	if len(result.Intents) != 0 || len(result.Skipped) != 0 {
		t.Fatalf("result = %#v, want no intent for unbound connection target", result)
	}
	record, _ := registry.FindConnectionByID(context.Background(), "connection-1", 1)
	if record.State != StateOpenUnbound {
		t.Fatalf("record state = %s, want unbound untouched", record.State)
	}
}

func TestClosePolicyNormalizesDefaultsAndRejectsInvalidInputBeforeMutation(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-record-1")
	policy := NewClosePolicy(registry, WithClosePolicyClock(staticClock{now: fixedTime().Add(time.Minute)}))

	tests := []struct {
		name    string
		command CloseConnectionsCommand
		code    ClosePolicyErrorCode
	}{
		{
			name: "missing target",
			command: CloseConnectionsCommand{
				ReasonClass: CloseReasonTokenRevoked,
			},
			code: ClosePolicyErrorCodeTargetInvalid,
		},
		{
			name: "missing reason",
			command: CloseConnectionsCommand{
				Target: TargetPlayer("player-1"),
			},
			code: ClosePolicyErrorCodeReasonInvalid,
		},
		{
			name: "unknown reason",
			command: CloseConnectionsCommand{
				Target:      TargetPlayer("player-1"),
				ReasonClass: CloseReasonClass("raw token expired"),
			},
			code: ClosePolicyErrorCodeReasonInvalid,
		},
		{
			name: "invalid visibility",
			command: CloseConnectionsCommand{
				Target:           TargetPlayer("player-1"),
				ReasonClass:      CloseReasonTokenRevoked,
				PublicVisibility: ClosePublicVisibility("show raw reason"),
			},
			code: ClosePolicyErrorCodeReasonInvalid,
		},
		{
			name: "invalid retryability",
			command: CloseConnectionsCommand{
				Target:       TargetPlayer("player-1"),
				ReasonClass:  CloseReasonTokenRevoked,
				Retryability: CloseRetryability("maybe with token details"),
			},
			code: ClosePolicyErrorCodeReasonInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := policy.RequestClose(context.Background(), tt.command); !hasClosePolicyCode(err, tt.code) {
				t.Fatalf("RequestClose() error = %v, want %s", err, tt.code)
			}
			record, ok := registry.FindConnectionByID(context.Background(), "connection-1", 1)
			if !ok || record.State != StateBound {
				t.Fatalf("record = %#v, ok = %v, want original bound state after rejected command", record, ok)
			}
		})
	}

	result, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetPlayer("player-1"),
		ReasonClass: CloseReasonPolicyViolation,
	})
	if err != nil {
		t.Fatalf("RequestClose(defaults) error = %v, want nil", err)
	}
	if result.PublicVisibility != ClosePublicVisibilitySilent || result.Retryability != CloseRetryabilityUnknown {
		t.Fatalf("result defaults = %#v, want silent and unknown retryability", result)
	}
}

func TestClosePolicyErrorsAreRedactedAndContextAware(t *testing.T) {
	policy := NewClosePolicy(nil)
	if _, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetPlayer("player-1"),
		ReasonClass: CloseReasonTokenRevoked,
	}); !hasClosePolicyCode(err, ClosePolicyErrorCodeRegistryUnavailable) {
		t.Fatalf("RequestClose(nil registry) error = %v, want %s", err, ClosePolicyErrorCodeRegistryUnavailable)
	}

	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	policy = NewClosePolicy(registry, WithClosePolicyClock(staticClock{}))
	if _, err := policy.RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetPlayer("player-1"),
		ReasonClass: CloseReasonTokenRevoked,
	}); !hasClosePolicyCode(err, ClosePolicyErrorCodeClockUnavailable) {
		t.Fatalf("RequestClose(zero clock) error = %v, want %s", err, ClosePolicyErrorCodeClockUnavailable)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	policy = NewClosePolicy(registry, WithClosePolicyClock(staticClock{now: fixedTime()}))
	if _, err := policy.RequestClose(ctx, CloseConnectionsCommand{
		Target:      TargetPlayer("player-1"),
		ReasonClass: CloseReasonTokenRevoked,
	}); !errors.Is(err, context.Canceled) {
		t.Fatalf("RequestClose(canceled) error = %v, want context.Canceled", err)
	}

	err := closePolicyError(ClosePolicyErrorCodeTargetInvalid, errors.New("raw-token-secret"))
	if strings.Contains(err.Error(), "raw-token-secret") {
		t.Fatalf("ClosePolicyError.Error() = %q, leaked wrapped detail", err.Error())
	}
}

func TestClosePolicyResultCopiesSlices(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-record-1")
	result, err := NewClosePolicy(registry, WithClosePolicyClock(staticClock{now: fixedTime().Add(time.Minute)})).RequestClose(context.Background(), CloseConnectionsCommand{
		Target:      TargetPlayer("player-1"),
		ReasonClass: CloseReasonTokenRevoked,
	})
	if err != nil {
		t.Fatalf("RequestClose() error = %v, want nil", err)
	}

	copied := copyClosePolicyResult(result)
	copied.Intents[0].PlayerID = "mutated"
	if result.Intents[0].PlayerID != "player-1" {
		t.Fatalf("copyClosePolicyResult() aliased intents: result = %#v copied = %#v", result, copied)
	}
}

func TestCloseIntentDoesNotContainRawProofTransportOrCloseReasonFields(t *testing.T) {
	intentType := reflect.TypeOf(CloseIntent{})
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
		"CloseCode",
		"CloseReason",
		"ReasonText",
	} {
		if _, ok := intentType.FieldByName(forbidden); ok {
			t.Fatalf("CloseIntent contains forbidden proof, transport, or close reason field %s", forbidden)
		}
	}
}

func assertCloseIntentIDs(t *testing.T, intents []CloseIntent, expected []ConnectionID) {
	t.Helper()
	actual := make([]ConnectionID, 0, len(intents))
	for _, intent := range intents {
		actual = append(actual, intent.ConnectionID)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("close intent ids = %#v, want %#v from intents %#v", actual, expected, intents)
	}
}

func hasClosePolicyCode(err error, code ClosePolicyErrorCode) bool {
	var policyErr *ClosePolicyError
	return errors.As(err, &policyErr) && policyErr.Code == code
}
