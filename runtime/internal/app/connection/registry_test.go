package connection

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRegisterOpenConnectionStoresServerObservedActiveState(t *testing.T) {
	openedAt := time.Date(2026, 5, 18, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	registry := NewInMemoryRegistry(staticClock{now: openedAt})

	record, err := registry.RegisterOpenConnection(context.Background(), OpenConnection{
		ConnectionID:    " connection-1 ",
		ConnectionEpoch: 7,
	})
	if err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	if record.ConnectionID != "connection-1" ||
		record.ConnectionEpoch != 7 ||
		record.State != StateOpenUnbound ||
		!record.OpenedAt.Equal(openedAt.UTC()) ||
		!record.CreatedAt.Equal(openedAt.UTC()) ||
		!record.UpdatedAt.Equal(openedAt.UTC()) {
		t.Fatalf("registered record = %#v, want normalized open unbound state", record)
	}
	if record.ActorKind != "" ||
		record.PlayerID != "" ||
		record.RuntimeSessionID != "" ||
		record.AccessTokenRecordID != "" {
		t.Fatalf("registered record contains identity linkage before validation: %#v", record)
	}
}

func TestRegisterOpenConnectionRejectsDuplicateActiveSameIDAndEpoch(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() initial error = %v, want nil", err)
	}
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); !hasRegistryCode(err, ErrorCodeConnectionAlreadyOpen) {
		t.Fatalf("RegisterOpenConnection() duplicate error = %v, want %s", err, ErrorCodeConnectionAlreadyOpen)
	}
}

func TestRegisterOpenConnectionSupersedesEarlierActiveEpoch(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() initial error = %v, want nil", err)
	}
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 2)); err != nil {
		t.Fatalf("RegisterOpenConnection() new epoch error = %v, want nil", err)
	}
	oldRecord, ok := registry.FindConnectionByID(context.Background(), "connection-1", 1)
	if !ok {
		t.Fatal("FindConnectionByID(old epoch) ok = false, want true")
	}
	if oldRecord.State != StateSuperseded ||
		oldRecord.SupersededAt == nil ||
		!oldRecord.SupersededAt.Equal(fixedTime()) ||
		oldRecord.SupersededByEpoch != 2 {
		t.Fatalf("old epoch record = %#v, want superseded by new epoch", oldRecord)
	}
	newRecord, ok := registry.FindConnectionByID(context.Background(), "connection-1", 2)
	if !ok || newRecord.State != StateOpenUnbound {
		t.Fatalf("new epoch record = %#v, ok = %v, want open unbound", newRecord, ok)
	}
}

func TestRegisterOpenConnectionRejectsStaleEpochAfterNewerObservedEpoch(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 2)); err != nil {
		t.Fatalf("RegisterOpenConnection() initial error = %v, want nil", err)
	}
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); !hasRegistryCode(err, ErrorCodeConnectionEpochStale) {
		t.Fatalf("RegisterOpenConnection() stale epoch error = %v, want %s", err, ErrorCodeConnectionEpochStale)
	}
}

func TestRegisterOpenConnectionRequiresNewerEpochAfterTerminalState(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() initial error = %v, want nil", err)
	}
	if _, err := registry.MarkConnectionClosed(context.Background(), MarkClosed{
		ConnectionID:     "connection-1",
		ConnectionEpoch:  1,
		CloseReasonClass: "transport closed",
	}); err != nil {
		t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
	}

	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); !hasRegistryCode(err, ErrorCodeConnectionEpochStale) {
		t.Fatalf("RegisterOpenConnection() same epoch after closed error = %v, want %s", err, ErrorCodeConnectionEpochStale)
	}

	reopened, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 2))
	if err != nil {
		t.Fatalf("RegisterOpenConnection() newer epoch after closed error = %v, want nil", err)
	}
	if reopened.State != StateOpenUnbound || reopened.CloseReasonClass != "" || reopened.ClosedAt != nil || reopened.ConnectionEpoch != 2 {
		t.Fatalf("reopened record = %#v, want fresh active record at newer epoch", reopened)
	}
}

func TestBindConnectionIdentityRequiresValidatedPlayerIdentity(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}

	tests := []struct {
		name    string
		command BindIdentity
		code    ErrorCode
	}{
		{
			name: "missing player",
			command: BindIdentity{
				ConnectionID:    "connection-1",
				ConnectionEpoch: 1,
				ActorKind:       ActorKindPlayer,
			},
			code: ErrorCodeIdentityInvalid,
		},
		{
			name: "metadata only session and token record without player",
			command: BindIdentity{
				ConnectionID:        "connection-1",
				ConnectionEpoch:     1,
				ActorKind:           ActorKindPlayer,
				RuntimeSessionID:    "session-1",
				AccessTokenRecordID: "token-record-1",
			},
			code: ErrorCodeIdentityInvalid,
		},
		{
			name: "unsupported actor kind",
			command: BindIdentity{
				ConnectionID:    "connection-1",
				ConnectionEpoch: 1,
				ActorKind:       ActorKind("service"),
				PlayerID:        "player-1",
			},
			code: ErrorCodeUnsupportedActorKind,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := registry.BindConnectionIdentity(context.Background(), tt.command); !hasRegistryCode(err, tt.code) {
				t.Fatalf("BindConnectionIdentity() error = %v, want %s", err, tt.code)
			}
		})
	}
}

func TestBindConnectionIdentityStoresValidatedLinkageOnly(t *testing.T) {
	validatedAt := time.Date(2026, 5, 18, 9, 10, 11, 0, time.FixedZone("test", 8*60*60))
	registry := NewInMemoryRegistry(staticClock{now: validatedAt})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}

	record, err := registry.BindConnectionIdentity(context.Background(), BindIdentity{
		ConnectionID:        " connection-1 ",
		ConnectionEpoch:     1,
		ActorKind:           "",
		PlayerID:            " player-1 ",
		RuntimeSessionID:    " session-1 ",
		AccessTokenRecordID: " token-record-1 ",
	})
	if err != nil {
		t.Fatalf("BindConnectionIdentity() error = %v, want nil", err)
	}

	if record.State != StateBound ||
		record.ActorKind != ActorKindPlayer ||
		record.PlayerID != "player-1" ||
		record.RuntimeSessionID != "session-1" ||
		record.AccessTokenRecordID != "token-record-1" {
		t.Fatalf("bound record = %#v, want validated player/session/token linkage", record)
	}
	if record.BoundAt == nil ||
		record.LastSeenAt == nil ||
		!record.BoundAt.Equal(validatedAt.UTC()) ||
		!record.LastSeenAt.Equal(validatedAt.UTC()) {
		t.Fatalf("bound times = %#v, want validated_at UTC", record)
	}
}

func TestBindConnectionIdentityRejectsTerminalRecords(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() error = %v, want nil", err)
	}
	if _, err := registry.MarkConnectionClosed(context.Background(), MarkClosed{
		ConnectionID:     "connection-1",
		ConnectionEpoch:  1,
		CloseReasonClass: "transport_closed",
	}); err != nil {
		t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
	}
	if _, err := registry.BindConnectionIdentity(context.Background(), bindCommand("connection-1", 1, "player-1", "session-1", "token-1")); !hasRegistryCode(err, ErrorCodeConnectionNotActive) {
		t.Fatalf("BindConnectionIdentity() error = %v, want %s", err, ErrorCodeConnectionNotActive)
	}
}

func TestActiveListsReturnOnlyBoundActiveRecordsByTarget(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-1")
	registerAndBind(t, registry, "connection-2", 1, "player-1", "session-2", "token-2")
	registerAndBind(t, registry, "connection-3", 1, "player-2", "session-3", "token-3")
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-4", 1)); err != nil {
		t.Fatalf("RegisterOpenConnection() unbound error = %v, want nil", err)
	}
	if _, err := registry.MarkConnectionClosed(context.Background(), MarkClosed{
		ConnectionID:     "connection-2",
		ConnectionEpoch:  1,
		CloseReasonClass: "transport_closed",
	}); err != nil {
		t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
	}
	if _, err := registry.MarkConnectionInvalidated(context.Background(), Invalidation{
		ConnectionID:      "connection-3",
		ConnectionEpoch:   1,
		InvalidationClass: "token_revoked",
	}); err != nil {
		t.Fatalf("MarkConnectionInvalidated() error = %v, want nil", err)
	}

	byPlayer := registry.ListConnectionsByPlayerID(context.Background(), " player-1 ")
	assertConnectionIDs(t, byPlayer, []ConnectionID{"connection-1"})
	bySession := registry.ListConnectionsByRuntimeSessionID(context.Background(), " session-1 ")
	assertConnectionIDs(t, bySession, []ConnectionID{"connection-1"})
	byToken := registry.ListConnectionsByAccessTokenRecordID(context.Background(), " token-1 ")
	assertConnectionIDs(t, byToken, []ConnectionID{"connection-1"})

	if got := registry.ListConnectionsByPlayerID(context.Background(), "player-3"); len(got) != 0 {
		t.Fatalf("ListConnectionsByPlayerID(player-3) = %#v, want empty", got)
	}
	if got := registry.ListConnectionsByRuntimeSessionID(context.Background(), " "); len(got) != 0 {
		t.Fatalf("ListConnectionsByRuntimeSessionID(empty) = %#v, want empty", got)
	}
}

func TestMarkConnectionClosedAndInvalidatedAreTerminalAndPolicyNeutral(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-1")
	registerAndBind(t, registry, "connection-2", 1, "player-1", "session-2", "token-2")

	closed, err := registry.MarkConnectionClosed(context.Background(), MarkClosed{
		ConnectionID:     "connection-1",
		ConnectionEpoch:  1,
		ClosedAt:         fixedTime().Add(time.Minute),
		CloseReasonClass: " transport observed ",
	})
	if err != nil {
		t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
	}
	if closed.State != StateClosed || closed.ClosedAt == nil || closed.CloseReasonClass != "transport_observed" {
		t.Fatalf("closed record = %#v, want closed terminal record", closed)
	}

	invalidated, err := registry.MarkConnectionInvalidated(context.Background(), Invalidation{
		ConnectionID:      "connection-2",
		ConnectionEpoch:   1,
		InvalidatedAt:     fixedTime().Add(2 * time.Minute),
		InvalidationClass: " token revoked ",
	})
	if err != nil {
		t.Fatalf("MarkConnectionInvalidated() error = %v, want nil", err)
	}
	if invalidated.State != StateInvalidated ||
		invalidated.InvalidatedAt == nil ||
		invalidated.InvalidationClass != "token_revoked" ||
		invalidated.ClosedAt != nil ||
		invalidated.CloseReasonClass != "" {
		t.Fatalf("invalidated record = %#v, want invalidated without socket close side effect", invalidated)
	}

	if got := registry.ListConnectionsByPlayerID(context.Background(), "player-1"); len(got) != 0 {
		t.Fatalf("ListConnectionsByPlayerID(player-1) = %#v, want terminal records excluded", got)
	}
	if _, err := registry.MarkConnectionInvalidated(context.Background(), Invalidation{
		ConnectionID:      "connection-1",
		ConnectionEpoch:   1,
		InvalidationClass: "again",
	}); !hasRegistryCode(err, ErrorCodeConnectionNotActive) {
		t.Fatalf("MarkConnectionInvalidated() on closed error = %v, want %s", err, ErrorCodeConnectionNotActive)
	}
}

func TestReturnedRecordsAndSlicesAreCopies(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-1")

	found, ok := registry.FindConnectionByID(context.Background(), "connection-1", 1)
	if !ok {
		t.Fatal("FindConnectionByID() ok = false, want true")
	}
	if found.BoundAt == nil || found.LastSeenAt == nil {
		t.Fatalf("found record = %#v, want time pointers", found)
	}
	mutatedTime := fixedTime().Add(10 * time.Hour)
	*found.BoundAt = mutatedTime
	found.PlayerID = "mutated"

	foundAgain, ok := registry.FindConnectionByID(context.Background(), "connection-1", 1)
	if !ok {
		t.Fatal("FindConnectionByID() after mutation ok = false, want true")
	}
	if foundAgain.PlayerID != "player-1" || foundAgain.BoundAt.Equal(mutatedTime) {
		t.Fatalf("FindConnectionByID() returned aliased record = %#v", foundAgain)
	}

	list := registry.ListConnectionsByPlayerID(context.Background(), "player-1")
	if len(list) != 1 {
		t.Fatalf("ListConnectionsByPlayerID() len = %d, want 1", len(list))
	}
	*list[0].LastSeenAt = mutatedTime
	list[0].RuntimeSessionID = "mutated"
	listAgain := registry.ListConnectionsByPlayerID(context.Background(), "player-1")
	if len(listAgain) != 1 || listAgain[0].RuntimeSessionID != "session-1" || listAgain[0].LastSeenAt.Equal(mutatedTime) {
		t.Fatalf("ListConnectionsByPlayerID() returned aliased slice records = %#v", listAgain)
	}
}

func TestFindConnectionByIDReturnsTerminalRecordForLifecycleInspection(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-1")
	if _, err := registry.MarkConnectionClosed(context.Background(), MarkClosed{
		ConnectionID:     "connection-1",
		ConnectionEpoch:  1,
		CloseReasonClass: "transport_closed",
	}); err != nil {
		t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
	}

	record, ok := registry.FindConnectionByID(context.Background(), "connection-1", 1)
	if !ok {
		t.Fatal("FindConnectionByID() ok = false, want true")
	}
	if record.State != StateClosed || record.CloseReasonClass != "transport_closed" {
		t.Fatalf("FindConnectionByID() = %#v, want closed lifecycle record", record)
	}
}

func TestFindConnectionByIDReturnsSupersededRecordForLifecycleInspection(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-1")
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 2)); err != nil {
		t.Fatalf("RegisterOpenConnection() new epoch error = %v, want nil", err)
	}

	record, ok := registry.FindConnectionByID(context.Background(), "connection-1", 1)
	if !ok {
		t.Fatal("FindConnectionByID() ok = false, want true")
	}
	if record.State != StateSuperseded || record.SupersededAt == nil || record.SupersededByEpoch != 2 {
		t.Fatalf("FindConnectionByID() = %#v, want superseded lifecycle record", record)
	}
	if got := registry.ListConnectionsByPlayerID(context.Background(), "player-1"); len(got) != 0 {
		t.Fatalf("ListConnectionsByPlayerID(player-1) = %#v, want superseded record excluded", got)
	}
}

func TestPresenceForPlayerReportsOnlineFromActiveBoundConnections(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime().Add(5 * time.Minute)})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-2", "token-2")
	registerAndBind(t, registry, "connection-2", 1, "player-1", "session-1", "token-1")
	registerAndBind(t, registry, "connection-3", 1, "player-2", "session-3", "token-3")
	if _, err := registry.MarkConnectionClosed(context.Background(), MarkClosed{
		ConnectionID:     "connection-2",
		ConnectionEpoch:  1,
		ClosedAt:         fixedTime().Add(10 * time.Minute),
		CloseReasonClass: "transport_closed",
	}); err != nil {
		t.Fatalf("MarkConnectionClosed() error = %v, want nil", err)
	}

	presence := registry.PresenceForPlayer(context.Background(), " player-1 ")

	if presence.PlayerID != "player-1" ||
		presence.Status != PresenceStatusOnline ||
		presence.ConnectionCount != 1 ||
		len(presence.ActiveConnections) != 1 {
		t.Fatalf("PresenceForPlayer() = %#v, want one online active connection", presence)
	}
	connection := presence.ActiveConnections[0]
	if connection.ConnectionID != "connection-1" ||
		connection.ConnectionEpoch != 1 ||
		connection.RuntimeSessionID != "session-2" ||
		connection.LastSeenAt == nil ||
		!connection.LastSeenAt.Equal(fixedTime().Add(5*time.Minute)) {
		t.Fatalf("presence connection = %#v, want active bound connection metadata", connection)
	}
	if presence.LastSeenAt == nil || !presence.LastSeenAt.Equal(fixedTime().Add(5*time.Minute)) {
		t.Fatalf("LastSeenAt = %v, want registry binding time", presence.LastSeenAt)
	}
	if !reflect.DeepEqual(presence.RuntimeSessionIDs, []RuntimeSessionID{"session-2"}) ||
		!reflect.DeepEqual(presence.AccessTokenRecords, []AccessTokenRecordID{"token-2"}) {
		t.Fatalf("presence ids = %#v/%#v, want sorted active ids", presence.RuntimeSessionIDs, presence.AccessTokenRecords)
	}
	if !presence.ObservedAt.Equal(fixedTime().Add(5 * time.Minute)) {
		t.Fatalf("ObservedAt = %v, want registry clock", presence.ObservedAt)
	}
}

func TestPresenceForPlayerReportsOfflineForNoActiveBoundConnections(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-1")
	if _, err := registry.MarkConnectionInvalidated(context.Background(), Invalidation{
		ConnectionID:      "connection-1",
		ConnectionEpoch:   1,
		InvalidationClass: "token_revoked",
	}); err != nil {
		t.Fatalf("MarkConnectionInvalidated() error = %v, want nil", err)
	}

	presence := registry.PresenceForPlayer(context.Background(), "player-1")

	if presence.Status != PresenceStatusOffline ||
		presence.ConnectionCount != 0 ||
		len(presence.ActiveConnections) != 0 ||
		presence.LastSeenAt != nil {
		t.Fatalf("PresenceForPlayer() = %#v, want offline without active connections", presence)
	}
}

func TestPresenceForPlayerReturnsCopies(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})
	registerAndBind(t, registry, "connection-1", 1, "player-1", "session-1", "token-1")

	presence := registry.PresenceForPlayer(context.Background(), "player-1")
	if presence.LastSeenAt == nil || len(presence.ActiveConnections) != 1 || presence.ActiveConnections[0].LastSeenAt == nil {
		t.Fatalf("PresenceForPlayer() = %#v, want time pointers", presence)
	}
	mutatedTime := fixedTime().Add(time.Hour)
	*presence.LastSeenAt = mutatedTime
	*presence.ActiveConnections[0].LastSeenAt = mutatedTime
	presence.ActiveConnections[0].ConnectionID = "mutated"
	presence.RuntimeSessionIDs[0] = "mutated"
	presence.AccessTokenRecords[0] = "mutated"

	presenceAgain := registry.PresenceForPlayer(context.Background(), "player-1")
	if presenceAgain.ActiveConnections[0].ConnectionID != "connection-1" ||
		presenceAgain.LastSeenAt.Equal(mutatedTime) ||
		presenceAgain.ActiveConnections[0].LastSeenAt.Equal(mutatedTime) ||
		presenceAgain.RuntimeSessionIDs[0] != "session-1" ||
		presenceAgain.AccessTokenRecords[0] != "token-1" {
		t.Fatalf("PresenceForPlayer() returned aliased presence = %#v", presenceAgain)
	}
}

func TestRegistryRejectsInvalidConnectionShapeAndClock(t *testing.T) {
	registry := NewInMemoryRegistry(staticClock{now: time.Time{}})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 1)); !hasRegistryCode(err, ErrorCodeClockUnavailable) {
		t.Fatalf("RegisterOpenConnection() zero clock error = %v, want %s", err, ErrorCodeClockUnavailable)
	}

	registry = NewInMemoryRegistry(staticClock{now: fixedTime()})
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand(" ", 1)); !hasRegistryCode(err, ErrorCodeConnectionInvalid) {
		t.Fatalf("RegisterOpenConnection() missing id error = %v, want %s", err, ErrorCodeConnectionInvalid)
	}
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand("connection-1", 0)); !hasRegistryCode(err, ErrorCodeConnectionInvalid) {
		t.Fatalf("RegisterOpenConnection() missing epoch error = %v, want %s", err, ErrorCodeConnectionInvalid)
	}
	if _, err := registry.MarkConnectionInvalidated(context.Background(), Invalidation{
		ConnectionID:    "connection-1",
		ConnectionEpoch: 1,
	}); !hasRegistryCode(err, ErrorCodeInvalidationInvalid) {
		t.Fatalf("MarkConnectionInvalidated() missing class error = %v, want %s", err, ErrorCodeInvalidationInvalid)
	}
}

func TestRegistryHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	registry := NewInMemoryRegistry(staticClock{now: fixedTime()})

	if _, err := registry.RegisterOpenConnection(ctx, openCommand("connection-1", 1)); !errors.Is(err, context.Canceled) {
		t.Fatalf("RegisterOpenConnection() error = %v, want context.Canceled", err)
	}
	if _, ok := registry.FindConnectionByID(ctx, "connection-1", 1); ok {
		t.Fatal("FindConnectionByID() ok = true with canceled context, want false")
	}
	if got := registry.ListConnectionsByPlayerID(ctx, "player-1"); got != nil {
		t.Fatalf("ListConnectionsByPlayerID() = %#v with canceled context, want nil", got)
	}
}

func TestRegistryRecordDoesNotContainRawProofMaterialFields(t *testing.T) {
	recordType := reflect.TypeOf(Record{})
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
	} {
		if _, ok := recordType.FieldByName(forbidden); ok {
			t.Fatalf("Record contains forbidden raw proof or transport metadata field %s", forbidden)
		}
	}
}

func registerAndBind(t *testing.T, registry *InMemoryRegistry, connectionID string, epoch ConnectionEpoch, playerID string, sessionID string, tokenRecordID string) {
	t.Helper()
	if _, err := registry.RegisterOpenConnection(context.Background(), openCommand(connectionID, epoch)); err != nil {
		t.Fatalf("RegisterOpenConnection(%s, %d) error = %v, want nil", connectionID, epoch, err)
	}
	if _, err := registry.BindConnectionIdentity(context.Background(), bindCommand(connectionID, epoch, playerID, sessionID, tokenRecordID)); err != nil {
		t.Fatalf("BindConnectionIdentity(%s, %d) error = %v, want nil", connectionID, epoch, err)
	}
}

func openCommand(connectionID string, epoch ConnectionEpoch) OpenConnection {
	return OpenConnection{
		ConnectionID:    ConnectionID(connectionID),
		ConnectionEpoch: epoch,
	}
}

func bindCommand(connectionID string, epoch ConnectionEpoch, playerID string, sessionID string, tokenRecordID string) BindIdentity {
	return BindIdentity{
		ConnectionID:        ConnectionID(connectionID),
		ConnectionEpoch:     epoch,
		ActorKind:           ActorKindPlayer,
		PlayerID:            PlayerID(playerID),
		RuntimeSessionID:    RuntimeSessionID(sessionID),
		AccessTokenRecordID: AccessTokenRecordID(tokenRecordID),
	}
}

func assertConnectionIDs(t *testing.T, records []Record, expected []ConnectionID) {
	t.Helper()
	actual := make([]ConnectionID, 0, len(records))
	for _, record := range records {
		actual = append(actual, record.ConnectionID)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("connection ids = %#v, want %#v from records %#v", actual, expected, records)
	}
}

func hasRegistryCode(err error, code ErrorCode) bool {
	var registryErr *RegistryError
	return errors.As(err, &registryErr) && registryErr.Code == code
}

func fixedTime() time.Time {
	return time.Date(2026, 5, 18, 1, 2, 3, 0, time.UTC)
}

type staticClock struct {
	now time.Time
}

func (c staticClock) Now() time.Time {
	return c.now
}
