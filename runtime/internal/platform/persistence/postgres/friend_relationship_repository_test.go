package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/modules/friends"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestFriendRelationshipRepositoryCreateOrUpdateInsertsPendingRequest(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 9, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			friendRelationshipRow{values: pendingFriendRelationshipRowValues(createdAt.UTC())},
		},
	}
	repository := NewFriendRelationshipRepositoryForUnitOfWork(executor)

	record, err := repository.CreateOrUpdateFriendRequest(context.Background(), friends.SendFriendRequestInput{
		RelationshipID: " relationship-1 ",
		Actor:          friends.FriendRelationshipActor{PlayerID: " player-a "},
		TargetPlayerID: " player-b ",
	})
	if err != nil {
		t.Fatalf("CreateOrUpdateFriendRequest() error = %v, want nil", err)
	}
	if record.RelationshipID != "relationship-1" ||
		record.Pair.PlayerLowID != "player-a" ||
		record.Pair.PlayerHighID != "player-b" ||
		record.LifecycleState != friends.FriendRelationshipLifecyclePending ||
		record.RequestedByPlayerID != "player-a" ||
		record.Version != friends.InitialFriendRelationshipVersion {
		t.Fatalf("CreateOrUpdateFriendRequest() record = %#v, want pending canonical relationship", record)
	}
	if record.CreatedAt.Location() != time.UTC || record.UpdatedAt.Location() != time.UTC {
		t.Fatalf("CreateOrUpdateFriendRequest() timestamps = %#v, want UTC", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "INSERT INTO friend_relationships")
	assertSQLContains(t, call.sql, "ON CONFLICT (player_low_id, player_high_id) DO UPDATE")
	assertSQLContains(t, call.sql, "relationship_version = friend_relationships.relationship_version + 1")
	assertSQLContains(t, call.sql, "RETURNING")
	assertArgs(t,
		call.args,
		friends.FriendRelationshipID("relationship-1"),
		"player-a",
		"player-b",
		string(friends.FriendRelationshipLifecyclePending),
		"player-a",
		int64(friends.InitialFriendRelationshipVersion),
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestFriendRelationshipRepositoryGetSelectsByCanonicalPair(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			friendRelationshipRow{values: pendingFriendRelationshipRowValues(createdAt)},
		},
	}
	repository := NewFriendRelationshipRepositoryForUnitOfWork(executor)

	record, err := repository.GetRelationshipByPair(context.Background(), friends.GetFriendRelationshipInput{
		Pair: friends.FriendRelationshipPair{PlayerLowID: " player-b ", PlayerHighID: " player-a "},
	})
	if err != nil {
		t.Fatalf("GetRelationshipByPair() error = %v, want nil", err)
	}
	if record.Pair.PlayerLowID != "player-a" || record.Pair.PlayerHighID != "player-b" {
		t.Fatalf("GetRelationshipByPair() pair = %#v, want canonical pair", record.Pair)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "FROM friend_relationships")
	assertSQLContains(t, call.sql, "player_low_id = $1")
	assertSQLContains(t, call.sql, "player_high_id = $2")
	assertArgs(t, call.args, "player-a", "player-b")
}

func TestFriendRelationshipRepositoryListIsPlayerScopedStatusFilteredAndOrdered(t *testing.T) {
	createdAt := time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC)
	executor := &recordingExecutor{
		rowsResponses: []pgx.Rows{
			&friendRelationshipRows{
				values: [][]any{
					pendingFriendRelationshipRowValues(createdAt, withFriendRelationshipPair("player-a", "player-b")),
					pendingFriendRelationshipRowValues(createdAt, withFriendRelationshipID("relationship-2"), withFriendRelationshipPair("player-a", "player-c")),
					pendingFriendRelationshipRowValues(createdAt, withFriendRelationshipID("relationship-3"), withFriendRelationshipPair("player-a", "player-d")),
				},
			},
		},
	}
	repository := NewFriendRelationshipRepositoryForUnitOfWork(executor)

	result, err := repository.ListRelationshipsForPlayer(context.Background(), friends.ListFriendRelationshipsInput{
		PlayerID:       " player-a ",
		Status:         friends.FriendRelationshipStatusPending,
		Limit:          2,
		AfterPairToken: " player-a|player-a ",
	})
	if err != nil {
		t.Fatalf("ListRelationshipsForPlayer() error = %v, want nil", err)
	}
	if len(result.Relationships) != 2 ||
		result.Relationships[0].Pair.PlayerHighID != "player-b" ||
		result.Relationships[1].Pair.PlayerHighID != "player-c" {
		t.Fatalf("ListRelationshipsForPlayer() relationships = %#v, want first two ordered rows", result.Relationships)
	}
	if result.NextPairToken != "player-a|player-d" {
		t.Fatalf("NextPairToken = %q, want overflow pair token", result.NextPairToken)
	}

	if len(executor.queries) != 1 {
		t.Fatalf("queries len = %d, want 1", len(executor.queries))
	}
	call := executor.queries[0]
	assertSQLContains(t, call.sql, "FROM friend_relationships")
	assertSQLContains(t, call.sql, "(player_low_id = $1 OR player_high_id = $1)")
	assertSQLContains(t, call.sql, "CASE $3")
	assertSQLContains(t, call.sql, "ORDER BY player_low_id, player_high_id")
	assertSQLContains(t, call.sql, "LIMIT $4")
	assertArgs(t, call.args, "player-a", "player-a|player-a", string(friends.FriendRelationshipStatusPending), int32(3))
}

func TestFriendRelationshipRepositoryLifecycleTransitionsUseExpectedVersion(t *testing.T) {
	changedAt := time.Date(2026, 5, 26, 2, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	expected := friends.FriendRelationshipVersion(3)

	tests := []struct {
		name       string
		call       func(friends.Repository) (friends.FriendRelationship, error)
		wantState  friends.FriendRelationshipLifecycleState
		wantSQL    []string
		wantArgs   []any
		rowOptions []friendRelationshipRowOption
	}{
		{
			name: "accept",
			call: func(repository friends.Repository) (friends.FriendRelationship, error) {
				return repository.AcceptFriendRequest(context.Background(), friends.AcceptFriendRequestInput{
					Actor:           friends.FriendRelationshipActor{PlayerID: " player-b "},
					Pair:            friends.FriendRelationshipPair{PlayerLowID: " player-b ", PlayerHighID: " player-a "},
					ExpectedVersion: &expected,
				})
			},
			wantState: friends.FriendRelationshipLifecycleFriends,
			wantSQL:   []string{"UPDATE friend_relationships", "lifecycle_state = $4", "responded_by_player_id = $3", "AND relationship_version = $5"},
			wantArgs:  []any{"player-a", "player-b", "player-b", string(friends.FriendRelationshipLifecycleFriends), int64(3)},
			rowOptions: []friendRelationshipRowOption{
				withFriendRelationshipLifecycle(friends.FriendRelationshipLifecycleFriends),
				withFriendRelationshipRespondedBy("player-b"),
				withFriendRelationshipVersion(4),
				withFriendRelationshipUpdatedAt(changedAt.UTC()),
				withFriendRelationshipStateChangedAt(changedAt.UTC()),
			},
		},
		{
			name: "reject",
			call: func(repository friends.Repository) (friends.FriendRelationship, error) {
				return repository.RejectFriendRequest(context.Background(), friends.RejectFriendRequestInput{
					Actor:           friends.FriendRelationshipActor{PlayerID: " player-b "},
					Pair:            friends.FriendRelationshipPair{PlayerLowID: " player-b ", PlayerHighID: " player-a "},
					ExpectedVersion: &expected,
				})
			},
			wantState: friends.FriendRelationshipLifecycleRejected,
			wantSQL:   []string{"UPDATE friend_relationships", "rejected_at = now()", "AND relationship_version = $5"},
			wantArgs:  []any{"player-a", "player-b", "player-b", string(friends.FriendRelationshipLifecycleRejected), int64(3)},
			rowOptions: []friendRelationshipRowOption{
				withFriendRelationshipLifecycle(friends.FriendRelationshipLifecycleRejected),
				withFriendRelationshipRespondedBy("player-b"),
				withFriendRelationshipRejectedAt(changedAt.UTC()),
				withFriendRelationshipVersion(4),
				withFriendRelationshipUpdatedAt(changedAt.UTC()),
				withFriendRelationshipStateChangedAt(changedAt.UTC()),
			},
		},
		{
			name: "remove",
			call: func(repository friends.Repository) (friends.FriendRelationship, error) {
				return repository.RemoveFriend(context.Background(), friends.RemoveFriendInput{
					Actor:           friends.FriendRelationshipActor{PlayerID: " player-a "},
					Pair:            friends.FriendRelationshipPair{PlayerLowID: " player-b ", PlayerHighID: " player-a "},
					ExpectedVersion: &expected,
				})
			},
			wantState: friends.FriendRelationshipLifecycleRemoved,
			wantSQL:   []string{"UPDATE friend_relationships", "removed_by_player_id = $3", "removed_at = now()", "blocked_by_low_at IS NULL", "blocked_by_high_at IS NULL", "AND relationship_version = $5"},
			wantArgs:  []any{"player-a", "player-b", "player-a", string(friends.FriendRelationshipLifecycleRemoved), int64(3)},
			rowOptions: []friendRelationshipRowOption{
				withFriendRelationshipLifecycle(friends.FriendRelationshipLifecycleRemoved),
				withFriendRelationshipRemovedBy("player-a"),
				withFriendRelationshipRemovedAt(changedAt.UTC()),
				withFriendRelationshipVersion(4),
				withFriendRelationshipUpdatedAt(changedAt.UTC()),
				withFriendRelationshipStateChangedAt(changedAt.UTC()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &recordingExecutor{
				rowResponses: []pgx.Row{
					friendRelationshipRow{values: pendingFriendRelationshipRowValues(changedAt.UTC(), tt.rowOptions...)},
				},
			}
			repository := NewFriendRelationshipRepositoryForUnitOfWork(executor)

			record, err := tt.call(repository)
			if err != nil {
				t.Fatalf("%s() error = %v, want nil", tt.name, err)
			}
			if record.LifecycleState != tt.wantState || record.Version != 4 {
				t.Fatalf("%s() record = %#v, want %s version 4", tt.name, record, tt.wantState)
			}
			if len(executor.queryRowCalls) != 1 {
				t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
			}
			call := executor.queryRowCalls[0]
			for _, want := range tt.wantSQL {
				assertSQLContains(t, call.sql, want)
			}
			assertArgs(t, call.args, tt.wantArgs...)
		})
	}
}

func TestFriendRelationshipRepositoryBlockAndUnblockUseActorSpecificColumns(t *testing.T) {
	changedAt := time.Date(2026, 5, 26, 3, 0, 0, 0, time.UTC)
	expected := friends.FriendRelationshipVersion(2)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			friendRelationshipRow{values: pendingFriendRelationshipRowValues(
				changedAt,
				withFriendRelationshipLifecycle(friends.FriendRelationshipLifecycleRemoved),
				withFriendRelationshipRemovedBy("player-a"),
				withFriendRelationshipRemovedAt(changedAt),
				withFriendRelationshipBlockedByLowAt(changedAt),
				withFriendRelationshipVersion(3),
			)},
			friendRelationshipRow{values: pendingFriendRelationshipRowValues(
				changedAt,
				withFriendRelationshipLifecycle(friends.FriendRelationshipLifecycleRemoved),
				withFriendRelationshipRemovedBy("player-a"),
				withFriendRelationshipRemovedAt(changedAt),
				withFriendRelationshipVersion(4),
			)},
		},
	}
	repository := NewFriendRelationshipRepositoryForUnitOfWork(executor)

	blocked, err := repository.SetPlayerBlock(context.Background(), friends.BlockPlayerInput{
		Actor:           friends.FriendRelationshipActor{PlayerID: " player-a "},
		TargetPlayerID:  " player-b ",
		ExpectedVersion: &expected,
	})
	if err != nil {
		t.Fatalf("SetPlayerBlock() error = %v, want nil", err)
	}
	if blocked.BlockState.BlockedByLowAt == nil || blocked.LifecycleState != friends.FriendRelationshipLifecycleRemoved {
		t.Fatalf("SetPlayerBlock() record = %#v, want low-side block and removed lifecycle", blocked)
	}

	unblocked, err := repository.ClearPlayerBlock(context.Background(), friends.UnblockPlayerInput{
		Actor:           friends.FriendRelationshipActor{PlayerID: " player-a "},
		TargetPlayerID:  " player-b ",
		ExpectedVersion: friendRelationshipVersionPtr(3),
	})
	if err != nil {
		t.Fatalf("ClearPlayerBlock() error = %v, want nil", err)
	}
	if unblocked.BlockState.BlockedByLowAt != nil || unblocked.LifecycleState != friends.FriendRelationshipLifecycleRemoved {
		t.Fatalf("ClearPlayerBlock() record = %#v, want low-side block cleared without restoring friendship", unblocked)
	}

	if len(executor.queryRowCalls) != 2 {
		t.Fatalf("query rows len = %d, want 2", len(executor.queryRowCalls))
	}
	blockCall := executor.queryRowCalls[0]
	assertSQLContains(t, blockCall.sql, "blocked_by_low_at = now()")
	assertSQLContains(t, blockCall.sql, "lifecycle_state = CASE")
	assertSQLContains(t, blockCall.sql, "removed_by_player_id = CASE")
	assertArgs(t, blockCall.args, "player-a", "player-b", "player-a", int64(2))

	unblockCall := executor.queryRowCalls[1]
	assertSQLContains(t, unblockCall.sql, "blocked_by_low_at = NULL")
	assertSQLContains(t, unblockCall.sql, "lifecycle_state")
	assertArgs(t, unblockCall.args, "player-a", "player-b", int64(3))
}

func TestFriendRelationshipRepositoryMapsConflictsAndRedactsDetails(t *testing.T) {
	repository := NewFriendRelationshipRepositoryForUnitOfWork(&recordingExecutor{})

	_, err := repository.GetRelationshipByPair(context.Background(), friends.GetFriendRelationshipInput{
		Pair: friends.FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"},
	})
	assertFriendRelationshipConflictClass(t, err, friends.FriendRelationshipConflictNotFound)

	duplicate := &pgconn.PgError{Code: "23505", ConstraintName: "friend_relationships_pair_uq", Detail: "player-a player-b"}
	repository = NewFriendRelationshipRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			friendRelationshipRow{err: duplicate},
		},
	})
	_, err = repository.CreateOrUpdateFriendRequest(context.Background(), validSendFriendRequestInput())
	assertFriendRelationshipConflictClass(t, err, friends.FriendRelationshipConflictPairIdentity)
	assertFriendRelationshipErrorRedacted(t, err)

	constraint := &pgconn.PgError{Code: "23503", ConstraintName: "friend_relationships_high_player_fk", Detail: "player-b"}
	repository = NewFriendRelationshipRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			friendRelationshipRow{err: constraint},
		},
	})
	_, err = repository.CreateOrUpdateFriendRequest(context.Background(), validSendFriendRequestInput())
	assertFriendRelationshipConflictClass(t, err, friends.FriendRelationshipConflictTargetPlayerNotFound)
	assertFriendRelationshipErrorRedacted(t, err)

	repository = NewFriendRelationshipRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			friendRelationshipRow{err: pgx.ErrNoRows},
			friendRelationshipRow{values: pendingFriendRelationshipRowValues(time.Date(2026, 5, 26, 1, 0, 0, 0, time.UTC), withFriendRelationshipVersion(5))},
		},
	})
	_, err = repository.AcceptFriendRequest(context.Background(), friends.AcceptFriendRequestInput{
		Actor:           friends.FriendRelationshipActor{PlayerID: "player-b"},
		Pair:            friends.FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"},
		ExpectedVersion: friendRelationshipVersionPtr(3),
	})
	assertFriendRelationshipConflictClass(t, err, friends.FriendRelationshipConflictStaleVersion)
	assertFriendRelationshipErrorRedacted(t, err)

	repository = NewFriendRelationshipRepositoryForUnitOfWork(&recordingExecutor{})
	_, err = repository.CreateOrUpdateFriendRequest(context.Background(), friends.SendFriendRequestInput{
		RelationshipID: "relationship-1",
		Actor:          friends.FriendRelationshipActor{PlayerID: "player-a"},
		TargetPlayerID: "player-a",
	})
	assertFriendRelationshipConflictClass(t, err, friends.FriendRelationshipConflictSelfRelationshipForbidden)
	assertFriendRelationshipErrorRedacted(t, err)
}

func TestFriendRelationshipRepositoryRejectsInvalidRows(t *testing.T) {
	values := pendingFriendRelationshipRowValues(time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC))
	values[3] = "hidden"
	repository := NewFriendRelationshipRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			friendRelationshipRow{values: values},
		},
	})

	_, err := repository.GetRelationshipByPair(context.Background(), friends.GetFriendRelationshipInput{
		Pair: friends.FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"},
	})
	assertFriendRelationshipConflictClass(t, err, friends.FriendRelationshipConflictStorageUnavailable)
	assertFriendRelationshipErrorRedacted(t, err)
}

func TestFriendRelationshipRepositoryRequiresUnitOfWorkExecutor(t *testing.T) {
	repository := NewFriendRelationshipRepositoryForUnitOfWork(nil)

	_, err := repository.CreateOrUpdateFriendRequest(context.Background(), validSendFriendRequestInput())
	if err == nil {
		t.Fatal("CreateOrUpdateFriendRequest() error = nil, want executor error")
	}

	_, err = repository.ListRelationshipsForPlayer(context.Background(), friends.ListFriendRelationshipsInput{
		PlayerID: "player-a",
		Limit:    1,
	})
	if err == nil {
		t.Fatal("ListRelationshipsForPlayer() error = nil, want executor error")
	}
}

func TestFriendRelationshipRepositoryDefaultTestsDoNotRequireLivePostgreSQL(t *testing.T) {
	if os.Getenv("VIBIT_POSTGRES_TEST_DSN") != "" {
		t.Skip("live PostgreSQL environment is opt-in and not needed for this fake-executor test")
	}

	repository := NewFriendRelationshipRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			friendRelationshipRow{values: pendingFriendRelationshipRowValues(time.Date(2026, 5, 26, 1, 2, 3, 0, time.UTC))},
		},
	})

	if _, err := repository.GetRelationshipByPair(context.Background(), friends.GetFriendRelationshipInput{
		Pair: friends.FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"},
	}); err != nil {
		t.Fatalf("GetRelationshipByPair() error = %v, want nil without live PostgreSQL", err)
	}
}

func TestPostgresUnitOfWorkCreatesFriendRelationshipRepository(t *testing.T) {
	executor := &recordingExecutor{}
	unit := UnitOfWork{executor: executor}

	repository, err := unit.NewFriendRelationshipRepository()
	if err != nil {
		t.Fatalf("NewFriendRelationshipRepository() error = %v, want nil", err)
	}
	if repository == nil {
		t.Fatal("NewFriendRelationshipRepository() = nil, want repository")
	}
}

func TestFriendRelationshipRepositoryImplementsFriendsRepository(t *testing.T) {
	var _ friends.Repository = (*FriendRelationshipRepository)(nil)
}

func validSendFriendRequestInput() friends.SendFriendRequestInput {
	return friends.SendFriendRequestInput{
		RelationshipID: "relationship-1",
		Actor:          friends.FriendRelationshipActor{PlayerID: "player-a"},
		TargetPlayerID: "player-b",
	}
}

func friendRelationshipVersionPtr(version friends.FriendRelationshipVersion) *friends.FriendRelationshipVersion {
	return &version
}

type friendRelationshipRowOption func([]any)

func withFriendRelationshipID(id friends.FriendRelationshipID) friendRelationshipRowOption {
	return func(values []any) {
		values[0] = id
	}
}

func withFriendRelationshipPair(low string, high string) friendRelationshipRowOption {
	return func(values []any) {
		values[1] = low
		values[2] = high
	}
}

func withFriendRelationshipLifecycle(state friends.FriendRelationshipLifecycleState) friendRelationshipRowOption {
	return func(values []any) {
		values[3] = string(state)
	}
}

func withFriendRelationshipRespondedBy(playerID string) friendRelationshipRowOption {
	return func(values []any) {
		values[5] = playerID
	}
}

func withFriendRelationshipRemovedBy(playerID string) friendRelationshipRowOption {
	return func(values []any) {
		values[6] = playerID
	}
}

func withFriendRelationshipVersion(version friends.FriendRelationshipVersion) friendRelationshipRowOption {
	return func(values []any) {
		values[7] = int64(version)
	}
}

func withFriendRelationshipUpdatedAt(updatedAt time.Time) friendRelationshipRowOption {
	return func(values []any) {
		values[9] = updatedAt
	}
}

func withFriendRelationshipStateChangedAt(stateChangedAt time.Time) friendRelationshipRowOption {
	return func(values []any) {
		values[10] = stateChangedAt
	}
}

func withFriendRelationshipRejectedAt(rejectedAt time.Time) friendRelationshipRowOption {
	return func(values []any) {
		values[11] = rejectedAt
	}
}

func withFriendRelationshipRemovedAt(removedAt time.Time) friendRelationshipRowOption {
	return func(values []any) {
		values[12] = removedAt
	}
}

func withFriendRelationshipBlockedByLowAt(blockedAt time.Time) friendRelationshipRowOption {
	return func(values []any) {
		values[13] = blockedAt
	}
}

func pendingFriendRelationshipRowValues(timestamp time.Time, options ...friendRelationshipRowOption) []any {
	values := []any{
		friends.FriendRelationshipID("relationship-1"),
		"player-a",
		"player-b",
		string(friends.FriendRelationshipLifecyclePending),
		"player-a",
		nil,
		nil,
		int64(friends.InitialFriendRelationshipVersion),
		timestamp,
		timestamp,
		timestamp,
		nil,
		nil,
		nil,
		nil,
	}
	for _, option := range options {
		option(values)
	}
	return values
}

type friendRelationshipRow struct {
	values []any
	err    error
}

func (r friendRelationshipRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	assignFriendRelationshipValues("friend relationship row", dest, r.values)
	return nil
}

type friendRelationshipRows struct {
	values [][]any
	index  int
	err    error
	closed bool
}

func (r *friendRelationshipRows) Close() {
	r.closed = true
}

func (r *friendRelationshipRows) Err() error {
	return r.err
}

func (r *friendRelationshipRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}

func (r *friendRelationshipRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *friendRelationshipRows) Next() bool {
	if r.index >= len(r.values) {
		r.closed = true
		return false
	}
	r.index += 1
	return true
}

func (r *friendRelationshipRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.values) {
		return errors.New("friend relationship rows: scan without current row")
	}
	assignFriendRelationshipValues("friend relationship rows", dest, r.values[r.index-1])
	return nil
}

func (r *friendRelationshipRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.values) {
		return nil, errors.New("friend relationship rows: values without current row")
	}
	return append([]any(nil), r.values[r.index-1]...), nil
}

func (r *friendRelationshipRows) RawValues() [][]byte {
	return nil
}

func (r *friendRelationshipRows) Conn() *pgx.Conn {
	return nil
}

func assignFriendRelationshipValues(label string, dest []any, values []any) {
	if len(dest) != len(values) {
		panic(label + ": destination count mismatch")
	}
	for i := range dest {
		assignFriendRelationshipValue(label, dest[i], values[i])
	}
}

func assignFriendRelationshipValue(label string, dest any, value any) {
	switch pointer := dest.(type) {
	case *friends.FriendRelationshipID:
		switch typed := value.(type) {
		case friends.FriendRelationshipID:
			*pointer = typed
		case string:
			*pointer = friends.FriendRelationshipID(typed)
		default:
			panic(label + ": unsupported relationship id value")
		}
	case *string:
		switch typed := value.(type) {
		case nil:
			*pointer = ""
		case friends.FriendRelationshipID:
			*pointer = string(typed)
		case string:
			*pointer = typed
		default:
			panic(label + ": unsupported string value")
		}
	case *int64:
		*pointer = value.(int64)
	case *time.Time:
		*pointer = value.(time.Time)
	case *pgtype.Text:
		switch text := value.(type) {
		case nil:
			*pointer = pgtype.Text{}
		case string:
			*pointer = pgtype.Text{String: text, Valid: true}
		case pgtype.Text:
			*pointer = text
		default:
			panic(label + ": unsupported text value")
		}
	case *pgtype.Timestamptz:
		switch timestamp := value.(type) {
		case nil:
			*pointer = pgtype.Timestamptz{}
		case time.Time:
			*pointer = pgtype.Timestamptz{Time: timestamp, Valid: true}
		case pgtype.Timestamptz:
			*pointer = timestamp
		default:
			panic(label + ": unsupported timestamptz value")
		}
	default:
		panic(label + ": unsupported destination type")
	}
}

func assertFriendRelationshipConflictClass(t *testing.T, err error, want friends.FriendRelationshipConflictClass) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want friends conflict class %q", want)
	}
	var repositoryErr *friends.FriendRelationshipRepositoryError
	if !errors.As(err, &repositoryErr) {
		t.Fatalf("error = %T %[1]v, want FriendRelationshipRepositoryError", err)
	}
	if !errors.Is(err, friends.ErrFriendRelationshipConflict) {
		t.Fatalf("errors.Is(err, ErrFriendRelationshipConflict) = false, want true for %v", err)
	}
	if repositoryErr.Conflict.Class != want {
		t.Fatalf("conflict class = %q, want %q", repositoryErr.Conflict.Class, want)
	}
}

func assertFriendRelationshipErrorRedacted(t *testing.T, err error) {
	t.Helper()
	text := err.Error()
	for _, forbidden := range []string{
		"player-a",
		"player-b",
		"relationship-1",
		"friend_relationships_pair_uq",
		"friend_relationships_high_player_fk",
		"SQL",
		"SELECT",
		"INSERT",
		"UPDATE",
		"driver",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("error %q leaks forbidden fragment %q", text, forbidden)
		}
	}
}
