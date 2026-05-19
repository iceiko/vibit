package ws

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestCloseHandoffRequestsConcreteSocketCloseByConnectionIDAndEpoch(t *testing.T) {
	server := &Server{}
	conn := &fakeCloseHandoffConn{}
	server.registerCloseHandoffSocket(connectionMetadata{id: " ws-1 ", epoch: 7}, conn)

	requestedAt := time.Date(2026, 5, 19, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	result := server.RequestClose(context.Background(), CloseHandoffRequest{
		ConnectionID:    " ws-1 ",
		ConnectionEpoch: 7,
		RequestedAt:     requestedAt,
	})

	if result.Outcome != CloseHandoffOutcomeCloseRequested {
		t.Fatalf("RequestClose() outcome = %s, want %s", result.Outcome, CloseHandoffOutcomeCloseRequested)
	}
	if result.ConnectionID != "ws-1" || result.ConnectionEpoch != 7 {
		t.Fatalf("RequestClose() target = %q/%d, want ws-1/7", result.ConnectionID, result.ConnectionEpoch)
	}
	if result.ClosedAt == nil || !result.ClosedAt.Equal(requestedAt.UTC()) {
		t.Fatalf("RequestClose() ClosedAt = %v, want %v", result.ClosedAt, requestedAt.UTC())
	}
	if conn.closeNowCalls != 1 {
		t.Fatalf("CloseNow calls = %d, want 1", conn.closeNowCalls)
	}
}

func TestCloseHandoffRejectsStaleEpochWithoutClosingCurrentSocket(t *testing.T) {
	server := &Server{}
	current := &fakeCloseHandoffConn{}
	server.registerCloseHandoffSocket(connectionMetadata{id: "ws-1", epoch: 8}, current)

	result := server.RequestClose(context.Background(), CloseHandoffRequest{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 7,
	})

	if result.Outcome != CloseHandoffOutcomeEpochMismatch {
		t.Fatalf("RequestClose() outcome = %s, want %s", result.Outcome, CloseHandoffOutcomeEpochMismatch)
	}
	if current.closeNowCalls != 0 {
		t.Fatalf("stale epoch closed current socket %d times, want 0", current.closeNowCalls)
	}
}

func TestCloseHandoffReturnsRedactedMissingAndAlreadyClosedOutcomes(t *testing.T) {
	server := &Server{}

	missing := server.RequestClose(context.Background(), CloseHandoffRequest{
		ConnectionID:    "missing",
		ConnectionEpoch: 1,
	})
	if missing.Outcome != CloseHandoffOutcomeSocketNotFound {
		t.Fatalf("missing outcome = %s, want %s", missing.Outcome, CloseHandoffOutcomeSocketNotFound)
	}
	if missing.ClosedAt != nil {
		t.Fatalf("missing ClosedAt = %v, want nil", missing.ClosedAt)
	}

	conn := &fakeCloseHandoffConn{}
	server.registerCloseHandoffSocket(connectionMetadata{id: "ws-1", epoch: 1}, conn)
	first := server.RequestClose(context.Background(), CloseHandoffRequest{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 1,
	})
	if first.Outcome != CloseHandoffOutcomeCloseRequested {
		t.Fatalf("first close outcome = %s, want %s", first.Outcome, CloseHandoffOutcomeCloseRequested)
	}

	second := server.RequestClose(context.Background(), CloseHandoffRequest{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 1,
	})
	if second.Outcome != CloseHandoffOutcomeAlreadyClosed {
		t.Fatalf("second close outcome = %s, want %s", second.Outcome, CloseHandoffOutcomeAlreadyClosed)
	}
	if second.ClosedAt == nil {
		t.Fatal("second close ClosedAt = nil, want redacted close timestamp")
	}
	if conn.closeNowCalls != 1 {
		t.Fatalf("CloseNow calls = %d, want one concrete close attempt", conn.closeNowCalls)
	}
}

func TestCloseHandoffReturnsRedactedCloseFailedOutcome(t *testing.T) {
	server := &Server{}
	conn := &fakeCloseHandoffConn{err: errors.New("raw transport detail")}
	server.registerCloseHandoffSocket(connectionMetadata{id: "ws-1", epoch: 1}, conn)

	result := server.RequestClose(context.Background(), CloseHandoffRequest{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 1,
	})

	if result.Outcome != CloseHandoffOutcomeCloseFailed {
		t.Fatalf("RequestClose() outcome = %s, want %s", result.Outcome, CloseHandoffOutcomeCloseFailed)
	}
	if result.ClosedAt != nil {
		t.Fatalf("RequestClose() ClosedAt = %v, want nil after failed close", result.ClosedAt)
	}
	if conn.closeNowCalls != 1 {
		t.Fatalf("CloseNow calls = %d, want 1", conn.closeNowCalls)
	}
}

func TestCloseHandoffUnregistersOnlyMatchingEpoch(t *testing.T) {
	var table closeHandoffSocketTable
	oldConn := &fakeCloseHandoffConn{}
	currentConn := &fakeCloseHandoffConn{}
	table.register(connectionMetadata{id: "ws-1", epoch: 1}, oldConn)
	table.register(connectionMetadata{id: "ws-1", epoch: 2}, currentConn)

	table.unregister(connectionMetadata{id: "ws-1", epoch: 1})
	entry, outcome := table.find("ws-1", 2)
	if outcome != CloseHandoffOutcomeCloseRequested || entry == nil {
		t.Fatalf("find current after stale unregister outcome = %s, entry = %#v", outcome, entry)
	}

	table.unregister(connectionMetadata{id: "ws-1", epoch: 2})
	if _, outcome := table.find("ws-1", 2); outcome != CloseHandoffOutcomeSocketNotFound {
		t.Fatalf("find after matching unregister outcome = %s, want %s", outcome, CloseHandoffOutcomeSocketNotFound)
	}
}

func TestCloseHandoffRequestAndResultRemainCredentialNeutral(t *testing.T) {
	for _, typ := range []reflect.Type{
		reflect.TypeOf(CloseHandoffRequest{}),
		reflect.TypeOf(CloseHandoffResult{}),
	} {
		for _, forbidden := range []string{
			"PlayerID",
			"RuntimeSessionID",
			"AccessToken",
			"AccessTokenRecordID",
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
			if _, ok := typ.FieldByName(forbidden); ok {
				t.Fatalf("%s contains forbidden field %s", typ.Name(), forbidden)
			}
		}
	}
}

type fakeCloseHandoffConn struct {
	closeNowCalls int
	err           error
}

func (f *fakeCloseHandoffConn) CloseNow() error {
	f.closeNowCalls++
	return f.err
}
