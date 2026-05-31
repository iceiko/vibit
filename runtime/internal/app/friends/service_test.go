package friends

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	friendsmodule "github.com/iceiko/vibit/runtime/internal/modules/friends"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

func TestNewServiceRequiresDependencies(t *testing.T) {
	_, err := NewService(ServiceDependencies{})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable)

	var runner *recordingUnitOfWorkRunner
	_, err = NewService(ServiceDependencies{
		UnitOfWorkRunner:        runner,
		RelationshipIDGenerator: staticRelationshipIDGenerator{value: "relationship-1"},
	})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable)

	_, err = NewService(ServiceDependencies{
		UnitOfWorkRunner: &recordingUnitOfWorkRunner{},
	})
	assertServiceError(t, err, OperationNewService, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable)
}

func TestSendFriendRequestRejectsMetadataOnlyIdentityBeforeUnitOfWork(t *testing.T) {
	repository := &fakeFriendRelationshipRepository{}
	runner := &recordingUnitOfWorkRunner{unit: &fakeFriendsUnitOfWork{repository: repository}}
	generator := &recordingRelationshipIDGenerator{value: "relationship-1"}
	service := mustNewService(t, runner, generator)

	result, err := service.SendFriendRequest(context.Background(), SendFriendRequestRequest{
		Identity: app.MetadataOnlyIdentityFromSession(app.Session{
			PlayerID:  "player-a",
			SessionID: "metadata-session",
		}),
		TargetPlayerID: "player-b",
	})

	assertServiceError(t, err, OperationSendFriendRequest, FailureClassUnauthenticated, PublicErrorFriendshipUnauthenticated)
	if result.Status != FriendRelationshipOperationStatusRejected ||
		result.PublicErrorCode != PublicErrorFriendshipUnauthenticated ||
		result.FailureClass != FailureClassUnauthenticated {
		t.Fatalf("result = %#v, want unauthenticated rejection", result)
	}
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if generator.calls != 0 {
		t.Fatalf("relationship id generator calls = %d, want 0", generator.calls)
	}
	if repository.totalCalls() != 0 {
		t.Fatalf("repository calls = %d, want 0", repository.totalCalls())
	}
}

func TestSendFriendRequestRejectsSelfTargetBeforeGenerationAndUnitOfWork(t *testing.T) {
	repository := &fakeFriendRelationshipRepository{}
	runner := &recordingUnitOfWorkRunner{unit: &fakeFriendsUnitOfWork{repository: repository}}
	generator := &recordingRelationshipIDGenerator{value: "relationship-1"}
	service := mustNewService(t, runner, generator)

	result, err := service.SendFriendRequest(context.Background(), SendFriendRequestRequest{
		Identity:       validatedPlayerIdentity("player-a"),
		TargetPlayerID: " player-a ",
	})

	assertServiceError(t, err, OperationSendFriendRequest, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest)
	if result.Status != FriendRelationshipOperationStatusRejected ||
		result.PublicErrorCode != PublicErrorFriendshipInvalidRequest {
		t.Fatalf("result = %#v, want invalid request rejection", result)
	}
	if runner.calls != 0 {
		t.Fatalf("unit-of-work calls = %d, want 0", runner.calls)
	}
	if generator.calls != 0 {
		t.Fatalf("relationship id generator calls = %d, want 0", generator.calls)
	}
	if repository.totalCalls() != 0 {
		t.Fatalf("repository calls = %d, want 0", repository.totalCalls())
	}
}

func TestSendFriendRequestUsesValidatedActorAndGeneratedRelationshipID(t *testing.T) {
	events := []string{}
	repository := &fakeFriendRelationshipRepository{
		events:       &events,
		createResult: validFriendRelationshipRecord(withRelationshipID("relationship-1")),
	}
	runner := &recordingUnitOfWorkRunner{
		unit:   &fakeFriendsUnitOfWork{repository: repository, events: &events},
		events: &events,
	}
	generator := staticRelationshipIDGenerator{value: " relationship-1 ", events: &events}
	service := mustNewService(t, runner, generator)

	result, err := service.SendFriendRequest(context.Background(), SendFriendRequestRequest{
		Identity:       validatedPlayerIdentity("player-a"),
		TargetPlayerID: " player-b ",
	})
	if err != nil {
		t.Fatalf("SendFriendRequest() error = %v, want nil", err)
	}
	if result.Status != FriendRelationshipOperationStatusSent ||
		result.PublicStatus != FriendRelationshipPublicStatusOutgoingRequestPending ||
		result.Relationship.RelationshipID != "relationship-1" ||
		result.Version != friendsmodule.InitialFriendRelationshipVersion {
		t.Fatalf("result = %#v, want sent outgoing pending relationship", result)
	}
	if repository.createCalls != 1 {
		t.Fatalf("CreateOrUpdateFriendRequest calls = %d, want 1", repository.createCalls)
	}
	if repository.lastCreateInput.RelationshipID != "relationship-1" ||
		repository.lastCreateInput.Actor.PlayerID != "player-a" ||
		repository.lastCreateInput.TargetPlayerID != "player-b" {
		t.Fatalf("CreateOrUpdateFriendRequest input = %#v, want server-derived actor and generated id", repository.lastCreateInput)
	}
	assertEvents(t, events, []string{"generate-relationship-id", "begin", "new-friend-relationship-repository", "create-or-update-friend-request", "commit"})
}

func TestAcceptFriendRequestRequiresIncomingPendingBeforeMutation(t *testing.T) {
	repository := &fakeFriendRelationshipRepository{
		getResult: validFriendRelationshipRecord(withRelationshipRequestedBy("player-a")),
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeFriendsUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticRelationshipIDGenerator{value: "relationship-1"})
	expected := friendsmodule.FriendRelationshipVersion(2)

	result, err := service.AcceptFriendRequest(context.Background(), AcceptFriendRequestRequest{
		Identity:        validatedPlayerIdentity("player-a"),
		TargetPlayerID:  "player-b",
		ExpectedVersion: &expected,
	})

	assertServiceError(t, err, OperationAcceptFriendRequest, FailureClassInvalidTransition, PublicErrorFriendshipInvalidTransition)
	if result.Status != FriendRelationshipOperationStatusRejected ||
		result.PublicErrorCode != PublicErrorFriendshipInvalidTransition {
		t.Fatalf("result = %#v, want invalid transition rejection", result)
	}
	if repository.getCalls != 1 {
		t.Fatalf("GetRelationshipByPair calls = %d, want 1", repository.getCalls)
	}
	if repository.acceptCalls != 0 {
		t.Fatalf("AcceptFriendRequest calls = %d, want 0 for outgoing pending", repository.acceptCalls)
	}
}

func TestAcceptAndRejectFriendRequestUseIncomingPendingRelationship(t *testing.T) {
	events := []string{}
	accepted := validFriendRelationshipRecord(
		withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleFriends),
		withRelationshipRequestedBy("player-b"),
		withRelationshipRespondedBy("player-a"),
		withRelationshipVersion(2),
	)
	rejected := validFriendRelationshipRecord(
		withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRejected),
		withRelationshipRequestedBy("player-b"),
		withRelationshipRespondedBy("player-a"),
		withRelationshipRejectedAt(fixedFriendRelationshipTime().Add(time.Minute)),
		withRelationshipVersion(3),
	)
	repository := &fakeFriendRelationshipRepository{
		events:       &events,
		getResult:    validFriendRelationshipRecord(withRelationshipRequestedBy("player-b")),
		acceptResult: accepted,
		rejectResult: rejected,
	}
	runner := &recordingUnitOfWorkRunner{
		unit:   &fakeFriendsUnitOfWork{repository: repository, events: &events},
		events: &events,
	}
	service := mustNewService(t, runner, staticRelationshipIDGenerator{value: "relationship-1"})
	expectedAccept := friendsmodule.FriendRelationshipVersion(1)

	acceptResult, err := service.AcceptFriendRequest(context.Background(), AcceptFriendRequestRequest{
		Identity:        validatedPlayerIdentity("player-a"),
		TargetPlayerID:  " player-b ",
		ExpectedVersion: &expectedAccept,
	})
	if err != nil {
		t.Fatalf("AcceptFriendRequest() error = %v, want nil", err)
	}
	if acceptResult.Status != FriendRelationshipOperationStatusAccepted ||
		acceptResult.PublicStatus != FriendRelationshipPublicStatusFriends ||
		acceptResult.Version != 2 {
		t.Fatalf("accept result = %#v, want accepted friends relationship", acceptResult)
	}
	if repository.lastAcceptInput.Actor.PlayerID != "player-a" ||
		repository.lastAcceptInput.Pair.PlayerLowID != "player-a" ||
		repository.lastAcceptInput.Pair.PlayerHighID != "player-b" ||
		repository.lastAcceptInput.ExpectedVersion == nil ||
		*repository.lastAcceptInput.ExpectedVersion != expectedAccept {
		t.Fatalf("AcceptFriendRequest input = %#v, want actor pair and expected version", repository.lastAcceptInput)
	}

	expectedReject := friendsmodule.FriendRelationshipVersion(2)
	rejectResult, err := service.RejectFriendRequest(context.Background(), RejectFriendRequestRequest{
		Identity:        validatedPlayerIdentity("player-a"),
		TargetPlayerID:  "player-b",
		ExpectedVersion: &expectedReject,
	})
	if err != nil {
		t.Fatalf("RejectFriendRequest() error = %v, want nil", err)
	}
	if rejectResult.Status != FriendRelationshipOperationStatusRequestRejected ||
		rejectResult.PublicStatus != FriendRelationshipPublicStatusRejected ||
		rejectResult.Version != 3 {
		t.Fatalf("reject result = %#v, want rejected relationship", rejectResult)
	}
	if repository.lastRejectInput.ExpectedVersion == nil || *repository.lastRejectInput.ExpectedVersion != expectedReject {
		t.Fatalf("RejectFriendRequest expected version = %#v, want %d", repository.lastRejectInput.ExpectedVersion, expectedReject)
	}
	if repository.getCalls != 2 || repository.acceptCalls != 1 || repository.rejectCalls != 1 {
		t.Fatalf("get/accept/reject calls = %d/%d/%d, want 2/1/1", repository.getCalls, repository.acceptCalls, repository.rejectCalls)
	}
}

func TestRemoveBlockAndUnblockUseExpectedVersionAndActorRelativeStatus(t *testing.T) {
	repository := &fakeFriendRelationshipRepository{
		removeResult: validFriendRelationshipRecord(
			withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRemoved),
			withRelationshipRemovedBy("player-a"),
			withRelationshipRemovedAt(fixedFriendRelationshipTime().Add(time.Minute)),
			withRelationshipVersion(2),
		),
		blockResult: validFriendRelationshipRecord(
			withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRemoved),
			withRelationshipRemovedBy("player-a"),
			withRelationshipRemovedAt(fixedFriendRelationshipTime().Add(time.Minute)),
			withRelationshipBlockedByLowAt(fixedFriendRelationshipTime().Add(time.Minute)),
			withRelationshipVersion(3),
		),
		unblockResult: validFriendRelationshipRecord(
			withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRemoved),
			withRelationshipRemovedBy("player-a"),
			withRelationshipRemovedAt(fixedFriendRelationshipTime().Add(time.Minute)),
			withRelationshipVersion(4),
		),
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeFriendsUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticRelationshipIDGenerator{value: "relationship-1"})
	removeExpected := friendsmodule.FriendRelationshipVersion(1)

	removeResult, err := service.RemoveFriend(context.Background(), RemoveFriendRequest{
		Identity:        validatedPlayerIdentity("player-a"),
		TargetPlayerID:  "player-b",
		ExpectedVersion: &removeExpected,
	})
	if err != nil {
		t.Fatalf("RemoveFriend() error = %v, want nil", err)
	}
	if removeResult.Status != FriendRelationshipOperationStatusRemoved ||
		removeResult.PublicStatus != FriendRelationshipPublicStatusRemoved ||
		removeResult.Version != 2 {
		t.Fatalf("remove result = %#v, want removed status", removeResult)
	}
	if repository.lastRemoveInput.ExpectedVersion == nil || *repository.lastRemoveInput.ExpectedVersion != removeExpected {
		t.Fatalf("RemoveFriend expected version = %#v, want %d", repository.lastRemoveInput.ExpectedVersion, removeExpected)
	}

	blockExpected := friendsmodule.FriendRelationshipVersion(2)
	blockResult, err := service.BlockPlayer(context.Background(), BlockPlayerRequest{
		Identity:        validatedPlayerIdentity("player-a"),
		TargetPlayerID:  "player-b",
		ExpectedVersion: &blockExpected,
	})
	if err != nil {
		t.Fatalf("BlockPlayer() error = %v, want nil", err)
	}
	if blockResult.Status != FriendRelationshipOperationStatusBlocked ||
		blockResult.PublicStatus != FriendRelationshipPublicStatusBlockedByActor ||
		blockResult.Version != 3 {
		t.Fatalf("block result = %#v, want actor block status", blockResult)
	}

	unblockExpected := friendsmodule.FriendRelationshipVersion(3)
	unblockResult, err := service.UnblockPlayer(context.Background(), UnblockPlayerRequest{
		Identity:        validatedPlayerIdentity("player-a"),
		TargetPlayerID:  "player-b",
		ExpectedVersion: &unblockExpected,
	})
	if err != nil {
		t.Fatalf("UnblockPlayer() error = %v, want nil", err)
	}
	if unblockResult.Status != FriendRelationshipOperationStatusUnblocked ||
		unblockResult.PublicStatus != FriendRelationshipPublicStatusRemoved ||
		unblockResult.Version != 4 {
		t.Fatalf("unblock result = %#v, want removed status without restored friendship", unblockResult)
	}
}

func TestListFriendRelationshipsIsActorScopedStatusFilteredAndCopiesResults(t *testing.T) {
	source := []friendsmodule.FriendRelationship{
		validFriendRelationshipRecord(withRelationshipRequestedBy("player-a")),
		validFriendRelationshipRecord(
			withRelationshipID("relationship-2"),
			withRelationshipPair("player-a", "player-c"),
			withRelationshipRequestedBy("player-c"),
		),
	}
	repository := &fakeFriendRelationshipRepository{
		listResult: friendsmodule.ListFriendRelationshipsResult{
			Relationships: source,
			NextPairToken: "player-a|player-d",
		},
	}
	runner := &recordingUnitOfWorkRunner{unit: &fakeFriendsUnitOfWork{repository: repository}}
	service := mustNewService(t, runner, staticRelationshipIDGenerator{value: "relationship-1"})

	result, err := service.ListFriendRelationships(context.Background(), ListFriendRelationshipsRequest{
		Identity:       validatedPlayerIdentity("player-a"),
		Status:         friendsmodule.FriendRelationshipStatusPending,
		Limit:          0,
		AfterPairToken: " ",
	})
	if err != nil {
		t.Fatalf("ListFriendRelationships() error = %v, want nil", err)
	}
	if result.Status != FriendRelationshipOperationStatusListed ||
		len(result.Relationships) != 2 ||
		result.Relationships[0].PublicStatus != FriendRelationshipPublicStatusOutgoingRequestPending ||
		result.Relationships[1].PublicStatus != FriendRelationshipPublicStatusIncomingRequestPending ||
		result.NextPairToken != "player-a|player-d" {
		t.Fatalf("result = %#v, want actor-relative listed relationships", result)
	}
	if repository.lastListInput.PlayerID != "player-a" ||
		repository.lastListInput.Status != friendsmodule.FriendRelationshipStatusPending ||
		repository.lastListInput.Limit != friendsmodule.DefaultListFriendRelationshipsLimit {
		t.Fatalf("ListRelationshipsForPlayer input = %#v, want actor-scoped defaulted input", repository.lastListInput)
	}
	source[0].RequestedByPlayerID = "player-b"
	if result.Relationships[0].Relationship.RequestedByPlayerID != "player-a" {
		t.Fatalf("listed relationship was not copied: %#v", result.Relationships[0])
	}
}

func TestGetFriendRelationshipStatusComputesActorRelativeStatuses(t *testing.T) {
	blockedAt := fixedFriendRelationshipTime().Add(time.Minute)
	tests := []struct {
		name       string
		record     friendsmodule.FriendRelationship
		getErr     error
		wantStatus FriendRelationshipPublicStatus
	}{
		{
			name:       "none when not found",
			getErr:     friendRelationshipRepositoryError(friendsmodule.FriendRelationshipConflictNotFound, "relationship player-a player-b not found"),
			wantStatus: FriendRelationshipPublicStatusNone,
		},
		{
			name:       "outgoing pending",
			record:     validFriendRelationshipRecord(withRelationshipRequestedBy("player-a")),
			wantStatus: FriendRelationshipPublicStatusOutgoingRequestPending,
		},
		{
			name:       "incoming pending",
			record:     validFriendRelationshipRecord(withRelationshipRequestedBy("player-b")),
			wantStatus: FriendRelationshipPublicStatusIncomingRequestPending,
		},
		{
			name: "friends",
			record: validFriendRelationshipRecord(
				withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleFriends),
				withRelationshipRequestedBy("player-b"),
				withRelationshipRespondedBy("player-a"),
				withRelationshipVersion(2),
			),
			wantStatus: FriendRelationshipPublicStatusFriends,
		},
		{
			name: "blocked by actor",
			record: validFriendRelationshipRecord(
				withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRemoved),
				withRelationshipRemovedBy("player-a"),
				withRelationshipRemovedAt(blockedAt),
				withRelationshipBlockedByLowAt(blockedAt),
				withRelationshipVersion(3),
			),
			wantStatus: FriendRelationshipPublicStatusBlockedByActor,
		},
		{
			name: "blocked actor",
			record: validFriendRelationshipRecord(
				withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRemoved),
				withRelationshipRemovedBy("player-b"),
				withRelationshipRemovedAt(blockedAt),
				withRelationshipBlockedByHighAt(blockedAt),
				withRelationshipVersion(3),
			),
			wantStatus: FriendRelationshipPublicStatusBlockedActor,
		},
		{
			name: "mutual blocked",
			record: validFriendRelationshipRecord(
				withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRemoved),
				withRelationshipRemovedBy("player-a"),
				withRelationshipRemovedAt(blockedAt),
				withRelationshipBlockedByLowAt(blockedAt),
				withRelationshipBlockedByHighAt(blockedAt),
				withRelationshipVersion(4),
			),
			wantStatus: FriendRelationshipPublicStatusMutualBlocked,
		},
		{
			name: "removed",
			record: validFriendRelationshipRecord(
				withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRemoved),
				withRelationshipRemovedBy("player-b"),
				withRelationshipRemovedAt(blockedAt),
				withRelationshipVersion(2),
			),
			wantStatus: FriendRelationshipPublicStatusRemoved,
		},
		{
			name: "rejected",
			record: validFriendRelationshipRecord(
				withRelationshipLifecycle(friendsmodule.FriendRelationshipLifecycleRejected),
				withRelationshipRequestedBy("player-a"),
				withRelationshipRespondedBy("player-b"),
				withRelationshipRejectedAt(blockedAt),
				withRelationshipVersion(2),
			),
			wantStatus: FriendRelationshipPublicStatusRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeFriendRelationshipRepository{
				getResult: tt.record,
				getErr:    tt.getErr,
			}
			runner := &recordingUnitOfWorkRunner{unit: &fakeFriendsUnitOfWork{repository: repository}}
			service := mustNewService(t, runner, staticRelationshipIDGenerator{value: "relationship-1"})

			result, err := service.GetFriendRelationshipStatus(context.Background(), GetFriendRelationshipStatusRequest{
				Identity:       validatedPlayerIdentity("player-a"),
				TargetPlayerID: "player-b",
			})
			if err != nil {
				t.Fatalf("GetFriendRelationshipStatus() error = %v, want nil", err)
			}
			if result.Status != FriendRelationshipOperationStatusFound ||
				result.PublicStatus != tt.wantStatus {
				t.Fatalf("result = %#v, want public status %s", result, tt.wantStatus)
			}
		})
	}
}

func TestRepositoryConflictsMapToRedactedPublicFailures(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantClass  FailureClass
		wantPublic PublicErrorCode
	}{
		{
			name:       "duplicate pending",
			err:        friendRelationshipRepositoryError(friendsmodule.FriendRelationshipConflictDuplicatePendingRequest, "duplicate player-a player-b"),
			wantClass:  FailureClassDuplicateRequest,
			wantPublic: PublicErrorFriendshipDuplicateRequest,
		},
		{
			name:       "target not found",
			err:        friendRelationshipRepositoryError(friendsmodule.FriendRelationshipConflictTargetPlayerNotFound, "target player-b missing"),
			wantClass:  FailureClassTargetNotFound,
			wantPublic: PublicErrorFriendshipTargetNotFound,
		},
		{
			name:       "blocked relationship",
			err:        friendRelationshipRepositoryError(friendsmodule.FriendRelationshipConflictBlockedRelationship, "blocked_by_low_at player-a player-b"),
			wantClass:  FailureClassBlockedRelationship,
			wantPublic: PublicErrorFriendshipBlockedRelationship,
		},
		{
			name:       "stale version",
			err:        friendRelationshipRepositoryError(friendsmodule.FriendRelationshipConflictStaleVersion, "expected=1 actual=3 player-a"),
			wantClass:  FailureClassVersionMismatch,
			wantPublic: PublicErrorFriendshipVersionMismatch,
		},
		{
			name:       "storage unavailable",
			err:        friendRelationshipRepositoryError(friendsmodule.FriendRelationshipConflictStorageUnavailable, "SELECT * FROM friend_relationships WHERE player_low_id='player-a'"),
			wantClass:  FailureClassDependencyUnavailable,
			wantPublic: PublicErrorFriendshipUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeFriendRelationshipRepository{createErr: tt.err}
			runner := &recordingUnitOfWorkRunner{unit: &fakeFriendsUnitOfWork{repository: repository}}
			service := mustNewService(t, runner, staticRelationshipIDGenerator{value: "relationship-1"})

			result, err := service.SendFriendRequest(context.Background(), SendFriendRequestRequest{
				Identity:       validatedPlayerIdentity("player-a"),
				TargetPlayerID: "player-b",
			})

			assertServiceError(t, err, OperationSendFriendRequest, tt.wantClass, tt.wantPublic)
			assertNoLeak(t, err, "player-a")
			assertNoLeak(t, err, "player-b")
			assertNoLeak(t, err, "SELECT *")
			if result.Status != FriendRelationshipOperationStatusRejected ||
				result.PublicErrorCode != tt.wantPublic ||
				result.FailureClass != tt.wantClass {
				t.Fatalf("result = %#v, want %s/%s rejection", result, tt.wantClass, tt.wantPublic)
			}
		})
	}
}

func TestServiceRequiresUnitOfWorkFriendRelationshipRepositoryCapability(t *testing.T) {
	runner := &recordingUnitOfWorkRunner{unit: tx.NoopUnitOfWork{}}
	service := mustNewService(t, runner, staticRelationshipIDGenerator{value: "relationship-1"})

	result, err := service.ListFriendRelationships(context.Background(), ListFriendRelationshipsRequest{
		Identity: validatedPlayerIdentity("player-a"),
		Limit:    10,
	})

	assertServiceError(t, err, OperationListFriendRelationships, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable)
	if result.Status != FriendRelationshipOperationStatusRejected {
		t.Fatalf("result = %#v, want rejected", result)
	}
}

func TestServiceKeepsProtocolTransportAndPersistenceDependenciesOut(t *testing.T) {
	source, err := os.ReadFile("service.go")
	if err != nil {
		t.Fatalf("ReadFile(service.go) error = %v, want nil", err)
	}
	for _, forbidden := range []string{
		"/internal/platform/persistence/postgres",
		"/internal/platform/transport/ws",
		"/internal/platform/protocol/protobuf",
		"/internal/generated/proto",
		"github.com/jackc/" + "pgx",
		"github.com/coder/" + "websocket",
	} {
		if strings.Contains(string(source), forbidden) {
			t.Fatalf("service.go contains forbidden dependency %q", forbidden)
		}
	}
}

func mustNewService(t *testing.T, runner UnitOfWorkRunner, generator RelationshipIDGenerator) Service {
	t.Helper()
	service, err := NewService(ServiceDependencies{
		UnitOfWorkRunner:        runner,
		RelationshipIDGenerator: generator,
	})
	if err != nil {
		t.Fatalf("NewService() error = %v, want nil", err)
	}
	return service
}

func validatedPlayerIdentity(playerID string) app.RequestIdentity {
	return app.ValidatedPlayerIdentity(playerID, app.Session{
		ConnectionID:    "connection-1",
		SessionID:       "session-1",
		PlayerID:        playerID,
		ConnectionEpoch: 7,
	})
}

func validFriendRelationshipRecord(options ...friendRelationshipRecordOption) friendsmodule.FriendRelationship {
	createdAt := fixedFriendRelationshipTime()
	record := friendsmodule.FriendRelationship{
		RelationshipID:      "relationship-1",
		Pair:                friendsmodule.FriendRelationshipPair{PlayerLowID: "player-a", PlayerHighID: "player-b"},
		LifecycleState:      friendsmodule.FriendRelationshipLifecyclePending,
		RequestedByPlayerID: "player-a",
		Version:             friendsmodule.InitialFriendRelationshipVersion,
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt.Add(time.Second),
		StateChangedAt:      createdAt.Add(time.Second),
	}
	for _, option := range options {
		option(&record)
	}
	normalized, err := friendsmodule.NormalizeFriendRelationshipRecord(record)
	if err != nil {
		panic(err)
	}
	return normalized
}

func fixedFriendRelationshipTime() time.Time {
	return time.Date(2026, 5, 26, 12, 0, 0, 0, time.UTC)
}

type friendRelationshipRecordOption func(*friendsmodule.FriendRelationship)

func withRelationshipID(id friendsmodule.FriendRelationshipID) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		record.RelationshipID = id
	}
}

func withRelationshipPair(low string, high string) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		record.Pair = friendsmodule.FriendRelationshipPair{PlayerLowID: low, PlayerHighID: high}
	}
}

func withRelationshipLifecycle(state friendsmodule.FriendRelationshipLifecycleState) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		record.LifecycleState = state
	}
}

func withRelationshipRequestedBy(playerID string) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		record.RequestedByPlayerID = playerID
	}
}

func withRelationshipRespondedBy(playerID string) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		record.RespondedByPlayerID = playerID
	}
}

func withRelationshipRemovedBy(playerID string) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		record.RemovedByPlayerID = playerID
	}
}

func withRelationshipVersion(version friendsmodule.FriendRelationshipVersion) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		record.Version = version
	}
}

func withRelationshipRejectedAt(rejectedAt time.Time) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		value := rejectedAt.UTC()
		record.RejectedAt = &value
	}
}

func withRelationshipRemovedAt(removedAt time.Time) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		value := removedAt.UTC()
		record.RemovedAt = &value
	}
}

func withRelationshipBlockedByLowAt(blockedAt time.Time) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		value := blockedAt.UTC()
		record.BlockState.BlockedByLowAt = &value
	}
}

func withRelationshipBlockedByHighAt(blockedAt time.Time) friendRelationshipRecordOption {
	return func(record *friendsmodule.FriendRelationship) {
		value := blockedAt.UTC()
		record.BlockState.BlockedByHighAt = &value
	}
}

func friendRelationshipVersionPtr(version friendsmodule.FriendRelationshipVersion) *friendsmodule.FriendRelationshipVersion {
	return &version
}

func friendRelationshipRepositoryError(class friendsmodule.FriendRelationshipConflictClass, rawDetail string) error {
	kind := friendsmodule.ErrFriendRelationshipConflict
	if class == friendsmodule.FriendRelationshipConflictStorageUnavailable {
		kind = friendsmodule.ErrFriendRelationshipUnavailable
	}
	return &friendsmodule.FriendRelationshipRepositoryError{
		Kind: kind,
		Conflict: friendsmodule.FriendRelationshipConflict{
			Class:          class,
			Retryable:      class == friendsmodule.FriendRelationshipConflictStorageUnavailable,
			RedactedReason: string(class),
		},
		Operation:      "test",
		RedactedReason: string(class),
		Err:            errors.New(rawDetail),
	}
}

func assertServiceError(t *testing.T, err error, operation Operation, class FailureClass, publicCode PublicErrorCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want service error %s/%s", class, publicCode)
	}
	var serviceErr *ServiceError
	if !errors.As(err, &serviceErr) {
		t.Fatalf("error = %T %v, want *ServiceError", err, err)
	}
	if serviceErr.Operation != operation ||
		serviceErr.Class != class ||
		serviceErr.PublicCode != publicCode {
		t.Fatalf("service error = %#v, want operation=%s class=%s public=%s", serviceErr, operation, class, publicCode)
	}
}

func assertNoLeak(t *testing.T, err error, forbidden string) {
	t.Helper()
	if err == nil || forbidden == "" {
		return
	}
	if strings.Contains(err.Error(), forbidden) {
		t.Fatalf("error %q leaks forbidden detail %q", err.Error(), forbidden)
	}
	if strings.Contains(err.Error(), "friend_relationships") ||
		strings.Contains(err.Error(), "blocked_by_") ||
		strings.Contains(err.Error(), "expected=") {
		t.Fatalf("error %q leaks private relationship detail", err.Error())
	}
}

func assertEvents(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("events = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("events = %#v, want %#v", got, want)
		}
	}
}

type recordingUnitOfWorkRunner struct {
	unit      tx.UnitOfWork
	events    *[]string
	calls     int
	commitErr error
}

func (r *recordingUnitOfWorkRunner) WithinUnitOfWork(ctx context.Context, fn func(context.Context, tx.UnitOfWork) error) error {
	r.calls += 1
	if ctx == nil {
		ctx = context.Background()
	}
	unit := r.unit
	if unit == nil {
		unit = tx.NoopUnitOfWork{}
	}
	appendEvent(r.events, "begin")
	if err := fn(ctx, unit); err != nil {
		appendEvent(r.events, "rollback")
		return err
	}
	if r.commitErr != nil {
		appendEvent(r.events, "rollback")
		return r.commitErr
	}
	appendEvent(r.events, "commit")
	return nil
}

type fakeFriendsUnitOfWork struct {
	ctx             context.Context
	repository      friendsmodule.Repository
	repositoryErr   error
	events          *[]string
	repositoryCalls int
}

func (u *fakeFriendsUnitOfWork) Context() context.Context {
	if u.ctx == nil {
		return context.Background()
	}
	return u.ctx
}

func (u *fakeFriendsUnitOfWork) NewFriendRelationshipRepository() (friendsmodule.Repository, error) {
	u.repositoryCalls += 1
	appendEvent(u.events, "new-friend-relationship-repository")
	if u.repositoryErr != nil {
		return nil, u.repositoryErr
	}
	return u.repository, nil
}

type fakeFriendRelationshipRepository struct {
	events *[]string

	createResult  friendsmodule.FriendRelationship
	createErr     error
	getResult     friendsmodule.FriendRelationship
	getErr        error
	listResult    friendsmodule.ListFriendRelationshipsResult
	listErr       error
	acceptResult  friendsmodule.FriendRelationship
	acceptErr     error
	rejectResult  friendsmodule.FriendRelationship
	rejectErr     error
	removeResult  friendsmodule.FriendRelationship
	removeErr     error
	blockResult   friendsmodule.FriendRelationship
	blockErr      error
	unblockResult friendsmodule.FriendRelationship
	unblockErr    error

	createCalls  int
	getCalls     int
	listCalls    int
	acceptCalls  int
	rejectCalls  int
	removeCalls  int
	blockCalls   int
	unblockCalls int

	lastCreateInput  friendsmodule.SendFriendRequestInput
	lastGetInput     friendsmodule.GetFriendRelationshipInput
	lastListInput    friendsmodule.ListFriendRelationshipsInput
	lastAcceptInput  friendsmodule.AcceptFriendRequestInput
	lastRejectInput  friendsmodule.RejectFriendRequestInput
	lastRemoveInput  friendsmodule.RemoveFriendInput
	lastBlockInput   friendsmodule.BlockPlayerInput
	lastUnblockInput friendsmodule.UnblockPlayerInput
}

func (r *fakeFriendRelationshipRepository) CreateOrUpdateFriendRequest(_ context.Context, input friendsmodule.SendFriendRequestInput) (friendsmodule.FriendRelationship, error) {
	r.createCalls += 1
	r.lastCreateInput = input
	appendEvent(r.events, "create-or-update-friend-request")
	if r.createErr != nil {
		return friendsmodule.FriendRelationship{}, r.createErr
	}
	return r.createResult, nil
}

func (r *fakeFriendRelationshipRepository) GetRelationshipByPair(_ context.Context, input friendsmodule.GetFriendRelationshipInput) (friendsmodule.FriendRelationship, error) {
	r.getCalls += 1
	r.lastGetInput = input
	appendEvent(r.events, "get-friend-relationship")
	if r.getErr != nil {
		return friendsmodule.FriendRelationship{}, r.getErr
	}
	return r.getResult, nil
}

func (r *fakeFriendRelationshipRepository) ListRelationshipsForPlayer(_ context.Context, input friendsmodule.ListFriendRelationshipsInput) (friendsmodule.ListFriendRelationshipsResult, error) {
	r.listCalls += 1
	r.lastListInput = input
	appendEvent(r.events, "list-friend-relationships")
	if r.listErr != nil {
		return friendsmodule.ListFriendRelationshipsResult{}, r.listErr
	}
	return r.listResult, nil
}

func (r *fakeFriendRelationshipRepository) AcceptFriendRequest(_ context.Context, input friendsmodule.AcceptFriendRequestInput) (friendsmodule.FriendRelationship, error) {
	r.acceptCalls += 1
	r.lastAcceptInput = input
	appendEvent(r.events, "accept-friend-request")
	if r.acceptErr != nil {
		return friendsmodule.FriendRelationship{}, r.acceptErr
	}
	return r.acceptResult, nil
}

func (r *fakeFriendRelationshipRepository) RejectFriendRequest(_ context.Context, input friendsmodule.RejectFriendRequestInput) (friendsmodule.FriendRelationship, error) {
	r.rejectCalls += 1
	r.lastRejectInput = input
	appendEvent(r.events, "reject-friend-request")
	if r.rejectErr != nil {
		return friendsmodule.FriendRelationship{}, r.rejectErr
	}
	return r.rejectResult, nil
}

func (r *fakeFriendRelationshipRepository) RemoveFriend(_ context.Context, input friendsmodule.RemoveFriendInput) (friendsmodule.FriendRelationship, error) {
	r.removeCalls += 1
	r.lastRemoveInput = input
	appendEvent(r.events, "remove-friend")
	if r.removeErr != nil {
		return friendsmodule.FriendRelationship{}, r.removeErr
	}
	return r.removeResult, nil
}

func (r *fakeFriendRelationshipRepository) SetPlayerBlock(_ context.Context, input friendsmodule.BlockPlayerInput) (friendsmodule.FriendRelationship, error) {
	r.blockCalls += 1
	r.lastBlockInput = input
	appendEvent(r.events, "set-player-block")
	if r.blockErr != nil {
		return friendsmodule.FriendRelationship{}, r.blockErr
	}
	return r.blockResult, nil
}

func (r *fakeFriendRelationshipRepository) ClearPlayerBlock(_ context.Context, input friendsmodule.UnblockPlayerInput) (friendsmodule.FriendRelationship, error) {
	r.unblockCalls += 1
	r.lastUnblockInput = input
	appendEvent(r.events, "clear-player-block")
	if r.unblockErr != nil {
		return friendsmodule.FriendRelationship{}, r.unblockErr
	}
	return r.unblockResult, nil
}

func (r *fakeFriendRelationshipRepository) totalCalls() int {
	return r.createCalls + r.getCalls + r.listCalls + r.acceptCalls + r.rejectCalls + r.removeCalls + r.blockCalls + r.unblockCalls
}

type staticRelationshipIDGenerator struct {
	value  string
	events *[]string
}

func (g staticRelationshipIDGenerator) GenerateFriendRelationshipID(context.Context) (string, error) {
	appendEvent(g.events, "generate-relationship-id")
	return g.value, nil
}

type recordingRelationshipIDGenerator struct {
	value string
	err   error
	calls int
}

func (g *recordingRelationshipIDGenerator) GenerateFriendRelationshipID(context.Context) (string, error) {
	g.calls += 1
	if g.err != nil {
		return "", g.err
	}
	return g.value, nil
}

func appendEvent(events *[]string, event string) {
	if events == nil {
		return
	}
	*events = append(*events, event)
}
