package friends

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRepositoryInterfaceIsStorageNeutral(t *testing.T) {
	var _ Repository = recordingRepository{}
}

func TestLifecycleStatusAndConflictVocabulariesAreClosedSets(t *testing.T) {
	for _, state := range []FriendRelationshipLifecycleState{
		FriendRelationshipLifecyclePending,
		FriendRelationshipLifecycleFriends,
		FriendRelationshipLifecycleRejected,
		FriendRelationshipLifecycleRemoved,
	} {
		if !state.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", state)
		}
	}
	if FriendRelationshipLifecycleState("muted").IsValid() {
		t.Fatal(`FriendRelationshipLifecycleState("muted").IsValid() = true, want false`)
	}

	for _, status := range []FriendRelationshipStatus{
		FriendRelationshipStatusPending,
		FriendRelationshipStatusFriends,
		FriendRelationshipStatusBlocked,
		FriendRelationshipStatusEnded,
	} {
		if !status.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", status)
		}
	}
	if FriendRelationshipStatus("all").IsValid() {
		t.Fatal(`FriendRelationshipStatus("all").IsValid() = true, want false`)
	}

	for _, class := range []FriendRelationshipConflictClass{
		FriendRelationshipConflictNotFound,
		FriendRelationshipConflictTargetPlayerNotFound,
		FriendRelationshipConflictSelfRelationshipForbidden,
		FriendRelationshipConflictDuplicatePendingRequest,
		FriendRelationshipConflictAlreadyFriends,
		FriendRelationshipConflictBlockedRelationship,
		FriendRelationshipConflictInvalidTransition,
		FriendRelationshipConflictVersionMismatch,
		FriendRelationshipConflictStaleVersion,
		FriendRelationshipConflictPairIdentity,
		FriendRelationshipConflictStorageUnavailable,
	} {
		if !class.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", class)
		}
	}
	if FriendRelationshipConflictClass("private_history_visible").IsValid() {
		t.Fatal(`FriendRelationshipConflictClass("private_history_visible").IsValid() = true, want false`)
	}
}

func TestNormalizeFriendRelationshipPairCanonicalizesAndRejectsSelf(t *testing.T) {
	pair, err := NormalizeFriendRelationshipPair(FriendRelationshipPair{
		PlayerLowID:  " player-z ",
		PlayerHighID: " player-a ",
	})
	if err != nil {
		t.Fatalf("NormalizeFriendRelationshipPair() error = %v, want nil", err)
	}
	if pair.PlayerLowID != "player-a" || pair.PlayerHighID != "player-z" {
		t.Fatalf("pair = %#v, want canonical unordered pair", pair)
	}

	if _, err := NormalizeFriendRelationshipPair(FriendRelationshipPair{PlayerLowID: "player-1", PlayerHighID: "player-1"}); err == nil {
		t.Fatal("NormalizeFriendRelationshipPair(self) error = nil, want rejection")
	}
}

func TestNormalizeFriendRelationshipRecordTrimsTimesAndPairMembers(t *testing.T) {
	createdAt := time.Date(2026, 5, 25, 8, 0, 0, 0, time.FixedZone("test", 8*60*60))
	rejectedAt := createdAt.Add(time.Minute)
	blockedAt := createdAt.Add(2 * time.Minute)

	record, err := NormalizeFriendRelationshipRecord(FriendRelationship{
		RelationshipID:      " relationship-1 ",
		Pair:                FriendRelationshipPair{PlayerLowID: " player-b ", PlayerHighID: " player-a "},
		LifecycleState:      FriendRelationshipLifecycleRejected,
		RequestedByPlayerID: " player-a ",
		RespondedByPlayerID: " player-b ",
		BlockState:          FriendRelationshipBlockState{BlockedByHighAt: &blockedAt},
		Version:             FriendRelationshipVersion(3),
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt.Add(3 * time.Minute),
		StateChangedAt:      createdAt.Add(time.Minute),
		RejectedAt:          &rejectedAt,
	})
	if err != nil {
		t.Fatalf("NormalizeFriendRelationshipRecord() error = %v, want nil", err)
	}

	if record.RelationshipID != FriendRelationshipID("relationship-1") ||
		record.Pair.PlayerLowID != "player-a" ||
		record.Pair.PlayerHighID != "player-b" ||
		record.RequestedByPlayerID != "player-a" ||
		record.RespondedByPlayerID != "player-b" ||
		record.LifecycleState != FriendRelationshipLifecycleRejected {
		t.Fatalf("record = %#v, want trimmed canonical relationship", record)
	}
	if record.CreatedAt.Location() != time.UTC ||
		record.UpdatedAt.Location() != time.UTC ||
		record.StateChangedAt.Location() != time.UTC ||
		record.RejectedAt == nil ||
		record.RejectedAt.Location() != time.UTC ||
		record.BlockState.BlockedByHighAt == nil ||
		record.BlockState.BlockedByHighAt.Location() != time.UTC {
		t.Fatalf("record times = %#v, want UTC-normalized times", record)
	}
}

func TestNormalizeFriendRelationshipRecordRejectsInvalidShape(t *testing.T) {
	valid := validFriendRelationshipRecord()
	tests := []struct {
		name   string
		mutate func(*FriendRelationship)
	}{
		{name: "relationship_id", mutate: func(r *FriendRelationship) { r.RelationshipID = " " }},
		{name: "pair", mutate: func(r *FriendRelationship) { r.Pair.PlayerHighID = r.Pair.PlayerLowID }},
		{name: "lifecycle_state", mutate: func(r *FriendRelationship) { r.LifecycleState = "archived" }},
		{name: "requested_by_non_member", mutate: func(r *FriendRelationship) { r.RequestedByPlayerID = "player-x" }},
		{name: "responded_by_non_member", mutate: func(r *FriendRelationship) { r.RespondedByPlayerID = "player-x" }},
		{name: "removed_by_non_member", mutate: func(r *FriendRelationship) { r.RemovedByPlayerID = "player-x" }},
		{name: "version", mutate: func(r *FriendRelationship) { r.Version = 0 }},
		{name: "created_at", mutate: func(r *FriendRelationship) { r.CreatedAt = time.Time{} }},
		{name: "updated_at", mutate: func(r *FriendRelationship) { r.UpdatedAt = time.Time{} }},
		{name: "state_changed_at", mutate: func(r *FriendRelationship) { r.StateChangedAt = time.Time{} }},
		{name: "updated_before_created", mutate: func(r *FriendRelationship) { r.UpdatedAt = r.CreatedAt.Add(-time.Second) }},
		{name: "state_changed_before_created", mutate: func(r *FriendRelationship) { r.StateChangedAt = r.CreatedAt.Add(-time.Second) }},
		{name: "blocked_before_created", mutate: func(r *FriendRelationship) { at := r.CreatedAt.Add(-time.Second); r.BlockState.BlockedByLowAt = &at }},
		{name: "rejected_before_created", mutate: func(r *FriendRelationship) { at := r.CreatedAt.Add(-time.Second); r.RejectedAt = &at }},
		{name: "removed_before_created", mutate: func(r *FriendRelationship) { at := r.CreatedAt.Add(-time.Second); r.RemovedAt = &at }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := valid
			tt.mutate(&record)
			if _, err := NormalizeFriendRelationshipRecord(record); err == nil {
				t.Fatal("NormalizeFriendRelationshipRecord() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeSendFriendRequestInputAndBlockInputsRejectSelf(t *testing.T) {
	sendInput, err := NormalizeSendFriendRequestInput(SendFriendRequestInput{
		RelationshipID: " relationship-1 ",
		Actor:          FriendRelationshipActor{PlayerID: " player-1 "},
		TargetPlayerID: " player-2 ",
	})
	if err != nil {
		t.Fatalf("NormalizeSendFriendRequestInput() error = %v, want nil", err)
	}
	if sendInput.RelationshipID != FriendRelationshipID("relationship-1") || sendInput.Actor.PlayerID != "player-1" || sendInput.TargetPlayerID != "player-2" {
		t.Fatalf("send input = %#v, want trimmed values", sendInput)
	}

	if _, err := NormalizeSendFriendRequestInput(SendFriendRequestInput{RelationshipID: "relationship-1", Actor: FriendRelationshipActor{PlayerID: "player-1"}, TargetPlayerID: "player-1"}); err == nil {
		t.Fatal("NormalizeSendFriendRequestInput(self) error = nil, want rejection")
	}
	if _, err := NormalizeBlockPlayerInput(BlockPlayerInput{Actor: FriendRelationshipActor{PlayerID: "player-1"}, TargetPlayerID: "player-1"}); err == nil {
		t.Fatal("NormalizeBlockPlayerInput(self) error = nil, want rejection")
	}
	if _, err := NormalizeUnblockPlayerInput(UnblockPlayerInput{Actor: FriendRelationshipActor{PlayerID: "player-1"}, TargetPlayerID: "player-1"}); err == nil {
		t.Fatal("NormalizeUnblockPlayerInput(self) error = nil, want rejection")
	}
}

func TestNormalizePairMutationInputsRequireActorInPairAndValidVersion(t *testing.T) {
	expected := FriendRelationshipVersion(4)
	acceptInput, err := NormalizeAcceptFriendRequestInput(AcceptFriendRequestInput{
		Actor:           FriendRelationshipActor{PlayerID: " player-b "},
		Pair:            FriendRelationshipPair{PlayerLowID: " player-b ", PlayerHighID: " player-a "},
		ExpectedVersion: &expected,
	})
	if err != nil {
		t.Fatalf("NormalizeAcceptFriendRequestInput() error = %v, want nil", err)
	}
	if acceptInput.Pair.PlayerLowID != "player-a" ||
		acceptInput.Pair.PlayerHighID != "player-b" ||
		acceptInput.Actor.PlayerID != "player-b" ||
		acceptInput.ExpectedVersion == nil ||
		*acceptInput.ExpectedVersion != expected {
		t.Fatalf("accept input = %#v, want canonical pair and expected version", acceptInput)
	}
	expected = 9
	if *acceptInput.ExpectedVersion != FriendRelationshipVersion(4) {
		t.Fatalf("accept expected version aliases caller pointer, got %d", *acceptInput.ExpectedVersion)
	}

	for name, err := range map[string]error{
		"accept_actor_not_pair_member": mustErrNormalizeAccept(AcceptFriendRequestInput{Actor: FriendRelationshipActor{PlayerID: "player-x"}, Pair: FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"}}),
		"reject_bad_version":           mustErrNormalizeReject(RejectFriendRequestInput{Actor: FriendRelationshipActor{PlayerID: "player-a"}, Pair: FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"}, ExpectedVersion: versionPtr(0)}),
		"remove_actor_not_pair_member": mustErrNormalizeRemove(RemoveFriendInput{Actor: FriendRelationshipActor{PlayerID: "player-x"}, Pair: FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"}}),
	} {
		if err == nil {
			t.Fatalf("%s error = nil, want rejection", name)
		}
	}
}

func TestNormalizeListFriendRelationshipsResultNormalizesRecordsAndCopiesSlice(t *testing.T) {
	first := validFriendRelationshipRecord()
	second := validFriendRelationshipRecord()
	second.RelationshipID = "relationship-2"
	second.Pair = FriendRelationshipPair{PlayerLowID: " player-c ", PlayerHighID: " player-b "}
	second.RequestedByPlayerID = " player-b "

	records := []FriendRelationship{first, second}
	result, err := NormalizeListFriendRelationshipsResult(ListFriendRelationshipsResult{
		Relationships: records,
		NextPairToken: " pair-token ",
	})
	if err != nil {
		t.Fatalf("NormalizeListFriendRelationshipsResult() error = %v, want nil", err)
	}
	if len(result.Relationships) != 2 || result.NextPairToken != "pair-token" {
		t.Fatalf("result = %#v, want normalized records and token", result)
	}
	if result.Relationships[1].Pair.PlayerLowID != "player-b" ||
		result.Relationships[1].Pair.PlayerHighID != "player-c" ||
		result.Relationships[1].RequestedByPlayerID != "player-b" {
		t.Fatalf("second relationship = %#v, want canonical normalized record", result.Relationships[1])
	}

	records[0].RelationshipID = "mutated-after-normalize"
	if result.Relationships[0].RelationshipID != FriendRelationshipID("relationship-1") {
		t.Fatalf("result relationship aliases caller slice, got %#v", result.Relationships[0])
	}

	invalid := validFriendRelationshipRecord()
	invalid.Version = 0
	if _, err := NormalizeListFriendRelationshipsResult(ListFriendRelationshipsResult{Relationships: []FriendRelationship{invalid}}); err == nil {
		t.Fatal("NormalizeListFriendRelationshipsResult(invalid record) error = nil, want rejection")
	}
}

func TestNormalizeListFriendRelationshipsInputDefaultsAndRejectsInvalidPagination(t *testing.T) {
	input, err := NormalizeListFriendRelationshipsInput(ListFriendRelationshipsInput{
		PlayerID:       " player-1 ",
		Status:         FriendRelationshipStatus(" friends "),
		AfterPairToken: " pair-token ",
	})
	if err != nil {
		t.Fatalf("NormalizeListFriendRelationshipsInput() error = %v, want nil", err)
	}
	if input.PlayerID != "player-1" || input.Status != FriendRelationshipStatusFriends || input.Limit != DefaultListFriendRelationshipsLimit || input.AfterPairToken != "pair-token" {
		t.Fatalf("list input = %#v, want trimmed status, default limit, and cursor", input)
	}

	valid := ListFriendRelationshipsInput{PlayerID: "player-1", Limit: DefaultListFriendRelationshipsLimit}
	tests := []struct {
		name   string
		mutate func(*ListFriendRelationshipsInput)
	}{
		{name: "missing_player_id", mutate: func(i *ListFriendRelationshipsInput) { i.PlayerID = " " }},
		{name: "invalid_status", mutate: func(i *ListFriendRelationshipsInput) { i.Status = "hidden" }},
		{name: "limit_negative", mutate: func(i *ListFriendRelationshipsInput) { i.Limit = -1 }},
		{name: "limit_too_large", mutate: func(i *ListFriendRelationshipsInput) { i.Limit = MaxListFriendRelationshipsLimit + 1 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := valid
			tt.mutate(&input)
			if _, err := NormalizeListFriendRelationshipsInput(input); err == nil {
				t.Fatal("NormalizeListFriendRelationshipsInput() error = nil, want rejection")
			}
		})
	}
}

func TestConflictAndRepositoryErrorsAreTypedAndRedacted(t *testing.T) {
	conflict := FriendRelationshipConflict{
		Class:          FriendRelationshipConflictVersionMismatch,
		Expected:       FriendRelationshipVersion(3),
		Actual:         FriendRelationshipVersion(4),
		Retryable:      true,
		RedactedReason: "version_mismatch",
	}
	if !conflict.Class.IsValid() {
		t.Fatalf("%q IsValid() = false, want true", conflict.Class)
	}
	if strings.Contains(conflict.Error(), "player-1") || strings.Contains(conflict.Error(), "relationship-1") {
		t.Fatalf("conflict error leaks social graph material: %q", conflict.Error())
	}

	repositoryErr := &FriendRelationshipRepositoryError{
		Kind:           ErrFriendRelationshipUnavailable,
		Conflict:       conflict,
		Operation:      "accept",
		RedactedReason: "storage_unavailable",
		Err:            errors.New("driver leaked player-1 relationship-1"),
	}
	if !errors.Is(repositoryErr, ErrFriendRelationshipUnavailable) {
		t.Fatalf("errors.Is(repositoryErr, ErrFriendRelationshipUnavailable) = false, want true")
	}
	if !errors.Is(repositoryErr, conflict) {
		t.Fatalf("errors.Is(repositoryErr, conflict) = false, want true")
	}
	if strings.Contains(repositoryErr.Error(), "player-1") ||
		strings.Contains(repositoryErr.Error(), "relationship-1") ||
		strings.Contains(repositoryErr.Error(), "driver leaked") {
		t.Fatalf("repository error leaks wrapped detail: %q", repositoryErr.Error())
	}
}

func TestFriendRelationshipHasNoSecretTransportProtocolOrDistributedFields(t *testing.T) {
	forbiddenFragments := []string{
		"Token",
		"Credential",
		"Verifier",
		"Digest",
		"Connection",
		"Subprotocol",
		"Remote",
		"Authorization",
		"Cookie",
		"Chat",
		"Group",
		"Party",
		"Match",
		"Pitaya",
		"Nakama",
	}
	relationshipType := reflect.TypeOf(FriendRelationship{})
	for i := 0; i < relationshipType.NumField(); i++ {
		fieldName := relationshipType.Field(i).Name
		for _, fragment := range forbiddenFragments {
			if strings.Contains(fieldName, fragment) {
				t.Fatalf("FriendRelationship field %q contains forbidden fragment %q", fieldName, fragment)
			}
		}
	}
}

func validFriendRelationshipRecord() FriendRelationship {
	now := time.Date(2026, 5, 25, 1, 2, 3, 0, time.UTC)
	return FriendRelationship{
		RelationshipID:      "relationship-1",
		Pair:                FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"},
		LifecycleState:      FriendRelationshipLifecyclePending,
		RequestedByPlayerID: "player-a",
		Version:             InitialFriendRelationshipVersion,
		CreatedAt:           now,
		UpdatedAt:           now,
		StateChangedAt:      now,
	}
}

func versionPtr(version FriendRelationshipVersion) *FriendRelationshipVersion {
	return &version
}

func mustErrNormalizeAccept(input AcceptFriendRequestInput) error {
	_, err := NormalizeAcceptFriendRequestInput(input)
	return err
}

func mustErrNormalizeReject(input RejectFriendRequestInput) error {
	_, err := NormalizeRejectFriendRequestInput(input)
	return err
}

func mustErrNormalizeRemove(input RemoveFriendInput) error {
	_, err := NormalizeRemoveFriendInput(input)
	return err
}

type recordingRepository struct{}

func (recordingRepository) CreateOrUpdateFriendRequest(context.Context, SendFriendRequestInput) (FriendRelationship, error) {
	return FriendRelationship{}, nil
}

func (recordingRepository) GetRelationshipByPair(context.Context, GetFriendRelationshipInput) (FriendRelationship, error) {
	return FriendRelationship{}, nil
}

func (recordingRepository) ListRelationshipsForPlayer(context.Context, ListFriendRelationshipsInput) (ListFriendRelationshipsResult, error) {
	return ListFriendRelationshipsResult{}, nil
}

func (recordingRepository) AcceptFriendRequest(context.Context, AcceptFriendRequestInput) (FriendRelationship, error) {
	return FriendRelationship{}, nil
}

func (recordingRepository) RejectFriendRequest(context.Context, RejectFriendRequestInput) (FriendRelationship, error) {
	return FriendRelationship{}, nil
}

func (recordingRepository) RemoveFriend(context.Context, RemoveFriendInput) (FriendRelationship, error) {
	return FriendRelationship{}, nil
}

func (recordingRepository) SetPlayerBlock(context.Context, BlockPlayerInput) (FriendRelationship, error) {
	return FriendRelationship{}, nil
}

func (recordingRepository) ClearPlayerBlock(context.Context, UnblockPlayerInput) (FriendRelationship, error) {
	return FriendRelationship{}, nil
}
