package friends

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	ModuleName = "friends"

	DefaultListFriendRelationshipsLimit = 100
	MaxListFriendRelationshipsLimit     = 500

	InitialFriendRelationshipVersion FriendRelationshipVersion = 1
)

type FriendRelationshipID string

type FriendRelationshipLifecycleState string

const (
	FriendRelationshipLifecyclePending  FriendRelationshipLifecycleState = "pending"
	FriendRelationshipLifecycleFriends  FriendRelationshipLifecycleState = "friends"
	FriendRelationshipLifecycleRejected FriendRelationshipLifecycleState = "rejected"
	FriendRelationshipLifecycleRemoved  FriendRelationshipLifecycleState = "removed"
)

func (s FriendRelationshipLifecycleState) IsValid() bool {
	switch s {
	case FriendRelationshipLifecyclePending,
		FriendRelationshipLifecycleFriends,
		FriendRelationshipLifecycleRejected,
		FriendRelationshipLifecycleRemoved:
		return true
	default:
		return false
	}
}

type FriendRelationshipStatus string

const (
	FriendRelationshipStatusPending FriendRelationshipStatus = "pending"
	FriendRelationshipStatusFriends FriendRelationshipStatus = "friends"
	FriendRelationshipStatusBlocked FriendRelationshipStatus = "blocked"
	FriendRelationshipStatusEnded   FriendRelationshipStatus = "ended"
)

func (s FriendRelationshipStatus) IsValid() bool {
	switch s {
	case FriendRelationshipStatusPending,
		FriendRelationshipStatusFriends,
		FriendRelationshipStatusBlocked,
		FriendRelationshipStatusEnded:
		return true
	default:
		return false
	}
}

type FriendRelationshipVersion int64

func (v FriendRelationshipVersion) IsValid() bool {
	return v > 0
}

type FriendRelationshipPair struct {
	PlayerLowID  string
	PlayerHighID string
}

type FriendRelationshipActor struct {
	PlayerID string
}

type FriendRelationshipBlockState struct {
	BlockedByLowAt  *time.Time
	BlockedByHighAt *time.Time
}

type FriendRelationship struct {
	RelationshipID      FriendRelationshipID
	Pair                FriendRelationshipPair
	LifecycleState      FriendRelationshipLifecycleState
	RequestedByPlayerID string
	RespondedByPlayerID string
	RemovedByPlayerID   string
	BlockState          FriendRelationshipBlockState
	Version             FriendRelationshipVersion
	CreatedAt           time.Time
	UpdatedAt           time.Time
	StateChangedAt      time.Time
	RejectedAt          *time.Time
	RemovedAt           *time.Time
}

type SendFriendRequestInput struct {
	RelationshipID FriendRelationshipID
	Actor          FriendRelationshipActor
	TargetPlayerID string
}

type AcceptFriendRequestInput struct {
	Actor           FriendRelationshipActor
	Pair            FriendRelationshipPair
	ExpectedVersion *FriendRelationshipVersion
}

type RejectFriendRequestInput struct {
	Actor           FriendRelationshipActor
	Pair            FriendRelationshipPair
	ExpectedVersion *FriendRelationshipVersion
}

type RemoveFriendInput struct {
	Actor           FriendRelationshipActor
	Pair            FriendRelationshipPair
	ExpectedVersion *FriendRelationshipVersion
}

type BlockPlayerInput struct {
	Actor           FriendRelationshipActor
	TargetPlayerID  string
	ExpectedVersion *FriendRelationshipVersion
}

type UnblockPlayerInput struct {
	Actor           FriendRelationshipActor
	TargetPlayerID  string
	ExpectedVersion *FriendRelationshipVersion
}

type GetFriendRelationshipInput struct {
	Pair FriendRelationshipPair
}

type ListFriendRelationshipsInput struct {
	PlayerID       string
	Status         FriendRelationshipStatus
	Limit          int
	AfterPairToken string
}

type ListFriendRelationshipsResult struct {
	Relationships []FriendRelationship
	NextPairToken string
}

type Repository interface {
	CreateOrUpdateFriendRequest(context.Context, SendFriendRequestInput) (FriendRelationship, error)
	GetRelationshipByPair(context.Context, GetFriendRelationshipInput) (FriendRelationship, error)
	ListRelationshipsForPlayer(context.Context, ListFriendRelationshipsInput) (ListFriendRelationshipsResult, error)
	AcceptFriendRequest(context.Context, AcceptFriendRequestInput) (FriendRelationship, error)
	RejectFriendRequest(context.Context, RejectFriendRequestInput) (FriendRelationship, error)
	RemoveFriend(context.Context, RemoveFriendInput) (FriendRelationship, error)
	SetPlayerBlock(context.Context, BlockPlayerInput) (FriendRelationship, error)
	ClearPlayerBlock(context.Context, UnblockPlayerInput) (FriendRelationship, error)
}

type FriendRelationshipConflictClass string

const (
	FriendRelationshipConflictNotFound                  FriendRelationshipConflictClass = "relationship_not_found"
	FriendRelationshipConflictTargetPlayerNotFound      FriendRelationshipConflictClass = "target_player_not_found"
	FriendRelationshipConflictSelfRelationshipForbidden FriendRelationshipConflictClass = "self_relationship_forbidden"
	FriendRelationshipConflictDuplicatePendingRequest   FriendRelationshipConflictClass = "duplicate_pending_request"
	FriendRelationshipConflictAlreadyFriends            FriendRelationshipConflictClass = "already_friends"
	FriendRelationshipConflictBlockedRelationship       FriendRelationshipConflictClass = "blocked_relationship"
	FriendRelationshipConflictInvalidTransition         FriendRelationshipConflictClass = "invalid_transition"
	FriendRelationshipConflictVersionMismatch           FriendRelationshipConflictClass = "version_mismatch"
	FriendRelationshipConflictStaleVersion              FriendRelationshipConflictClass = "stale_relationship_version"
	FriendRelationshipConflictPairIdentity              FriendRelationshipConflictClass = "pair_identity_conflict"
	FriendRelationshipConflictStorageUnavailable        FriendRelationshipConflictClass = "storage_unavailable"
)

func (c FriendRelationshipConflictClass) IsValid() bool {
	switch c {
	case FriendRelationshipConflictNotFound,
		FriendRelationshipConflictTargetPlayerNotFound,
		FriendRelationshipConflictSelfRelationshipForbidden,
		FriendRelationshipConflictDuplicatePendingRequest,
		FriendRelationshipConflictAlreadyFriends,
		FriendRelationshipConflictBlockedRelationship,
		FriendRelationshipConflictInvalidTransition,
		FriendRelationshipConflictVersionMismatch,
		FriendRelationshipConflictStaleVersion,
		FriendRelationshipConflictPairIdentity,
		FriendRelationshipConflictStorageUnavailable:
		return true
	default:
		return false
	}
}

type FriendRelationshipConflict struct {
	Class          FriendRelationshipConflictClass
	Expected       FriendRelationshipVersion
	Actual         FriendRelationshipVersion
	Retryable      bool
	RedactedReason string
}

func (c FriendRelationshipConflict) Error() string {
	reason := strings.TrimSpace(c.RedactedReason)
	if reason == "" {
		reason = string(c.Class)
	}
	if reason == "" {
		return "friend relationship conflict"
	}
	return fmt.Sprintf("friend relationship conflict: %s", reason)
}

func (c FriendRelationshipConflict) Is(target error) bool {
	targetConflict, ok := target.(FriendRelationshipConflict)
	if !ok {
		return false
	}
	return c.Class == targetConflict.Class && c.Class != ""
}

var (
	ErrFriendRelationshipInvalidInput = errors.New("friend relationship repository: invalid input")
	ErrFriendRelationshipConflict     = errors.New("friend relationship repository: conflict")
	ErrFriendRelationshipUnavailable  = errors.New("friend relationship repository: storage unavailable")
)

type FriendRelationshipRepositoryError struct {
	Kind           error
	Conflict       FriendRelationshipConflict
	Operation      string
	RedactedReason string
	Err            error
}

func (e *FriendRelationshipRepositoryError) Error() string {
	if e == nil {
		return ""
	}
	reason := strings.TrimSpace(e.RedactedReason)
	if reason == "" && e.Conflict.Class != "" {
		reason = string(e.Conflict.Class)
	}
	if reason == "" && e.Kind != nil {
		reason = e.Kind.Error()
	}
	operation := strings.TrimSpace(e.Operation)
	if operation == "" {
		operation = "operation"
	}
	if reason == "" {
		return fmt.Sprintf("friend relationship repository %s failed", operation)
	}
	return fmt.Sprintf("friend relationship repository %s failed: %s", operation, reason)
}

func (e *FriendRelationshipRepositoryError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *FriendRelationshipRepositoryError) Is(target error) bool {
	if e == nil {
		return false
	}
	if errors.Is(e.Kind, target) || errors.Is(e.Err, target) {
		return true
	}
	if e.Conflict.Class != "" && errors.Is(e.Conflict, target) {
		return true
	}
	return target == ErrFriendRelationshipConflict && e.Conflict.Class != ""
}

func NormalizeFriendRelationshipRecord(record FriendRelationship) (FriendRelationship, error) {
	var err error
	record.RelationshipID, err = NormalizeFriendRelationshipID(record.RelationshipID)
	if err != nil {
		return FriendRelationship{}, err
	}
	record.Pair, err = NormalizeFriendRelationshipPair(record.Pair)
	if err != nil {
		return FriendRelationship{}, err
	}
	record.LifecycleState = FriendRelationshipLifecycleState(strings.TrimSpace(string(record.LifecycleState)))
	if !record.LifecycleState.IsValid() {
		return FriendRelationship{}, errors.New("friend relationship: lifecycle_state is invalid")
	}
	record.RequestedByPlayerID, err = normalizeOptionalPairMember("requested_by_player_id", record.RequestedByPlayerID, record.Pair)
	if err != nil {
		return FriendRelationship{}, err
	}
	record.RespondedByPlayerID, err = normalizeOptionalPairMember("responded_by_player_id", record.RespondedByPlayerID, record.Pair)
	if err != nil {
		return FriendRelationship{}, err
	}
	record.RemovedByPlayerID, err = normalizeOptionalPairMember("removed_by_player_id", record.RemovedByPlayerID, record.Pair)
	if err != nil {
		return FriendRelationship{}, err
	}
	if !record.Version.IsValid() {
		return FriendRelationship{}, errors.New("friend relationship: version must be positive")
	}
	record.CreatedAt, err = normalizeRequiredTime("created_at", record.CreatedAt)
	if err != nil {
		return FriendRelationship{}, err
	}
	record.UpdatedAt, err = normalizeRequiredTime("updated_at", record.UpdatedAt)
	if err != nil {
		return FriendRelationship{}, err
	}
	record.StateChangedAt, err = normalizeRequiredTime("state_changed_at", record.StateChangedAt)
	if err != nil {
		return FriendRelationship{}, err
	}
	if record.UpdatedAt.Before(record.CreatedAt) {
		return FriendRelationship{}, errors.New("friend relationship: updated_at must not be before created_at")
	}
	if record.StateChangedAt.Before(record.CreatedAt) {
		return FriendRelationship{}, errors.New("friend relationship: state_changed_at must not be before created_at")
	}
	record.BlockState, err = NormalizeFriendRelationshipBlockState(record.BlockState, record.CreatedAt)
	if err != nil {
		return FriendRelationship{}, err
	}
	if record.RejectedAt != nil {
		rejectedAt := record.RejectedAt.UTC()
		if rejectedAt.Before(record.CreatedAt) {
			return FriendRelationship{}, errors.New("friend relationship: rejected_at must not be before created_at")
		}
		record.RejectedAt = &rejectedAt
	}
	if record.RemovedAt != nil {
		removedAt := record.RemovedAt.UTC()
		if removedAt.Before(record.CreatedAt) {
			return FriendRelationship{}, errors.New("friend relationship: removed_at must not be before created_at")
		}
		record.RemovedAt = &removedAt
	}
	return record, nil
}

func NormalizeSendFriendRequestInput(input SendFriendRequestInput) (SendFriendRequestInput, error) {
	var err error
	input.RelationshipID, err = NormalizeFriendRelationshipID(input.RelationshipID)
	if err != nil {
		return SendFriendRequestInput{}, err
	}
	input.Actor, err = NormalizeFriendRelationshipActor(input.Actor)
	if err != nil {
		return SendFriendRequestInput{}, err
	}
	input.TargetPlayerID, err = normalizeRequired("target_player_id", input.TargetPlayerID)
	if err != nil {
		return SendFriendRequestInput{}, err
	}
	if input.Actor.PlayerID == input.TargetPlayerID {
		return SendFriendRequestInput{}, errors.New("friend relationship: self relationship is forbidden")
	}
	return input, nil
}

func NormalizeAcceptFriendRequestInput(input AcceptFriendRequestInput) (AcceptFriendRequestInput, error) {
	return normalizePairMutationInput(input.Actor, input.Pair, input.ExpectedVersion, func(actor FriendRelationshipActor, pair FriendRelationshipPair, expected *FriendRelationshipVersion) AcceptFriendRequestInput {
		return AcceptFriendRequestInput{Actor: actor, Pair: pair, ExpectedVersion: expected}
	})
}

func NormalizeRejectFriendRequestInput(input RejectFriendRequestInput) (RejectFriendRequestInput, error) {
	return normalizePairMutationInput(input.Actor, input.Pair, input.ExpectedVersion, func(actor FriendRelationshipActor, pair FriendRelationshipPair, expected *FriendRelationshipVersion) RejectFriendRequestInput {
		return RejectFriendRequestInput{Actor: actor, Pair: pair, ExpectedVersion: expected}
	})
}

func NormalizeRemoveFriendInput(input RemoveFriendInput) (RemoveFriendInput, error) {
	return normalizePairMutationInput(input.Actor, input.Pair, input.ExpectedVersion, func(actor FriendRelationshipActor, pair FriendRelationshipPair, expected *FriendRelationshipVersion) RemoveFriendInput {
		return RemoveFriendInput{Actor: actor, Pair: pair, ExpectedVersion: expected}
	})
}

func NormalizeBlockPlayerInput(input BlockPlayerInput) (BlockPlayerInput, error) {
	var err error
	input.Actor, err = NormalizeFriendRelationshipActor(input.Actor)
	if err != nil {
		return BlockPlayerInput{}, err
	}
	input.TargetPlayerID, err = normalizeRequired("target_player_id", input.TargetPlayerID)
	if err != nil {
		return BlockPlayerInput{}, err
	}
	if input.Actor.PlayerID == input.TargetPlayerID {
		return BlockPlayerInput{}, errors.New("friend relationship: self relationship is forbidden")
	}
	if input.ExpectedVersion != nil && !input.ExpectedVersion.IsValid() {
		return BlockPlayerInput{}, errors.New("friend relationship: expected_version must be positive")
	}
	if input.ExpectedVersion != nil {
		expected := *input.ExpectedVersion
		input.ExpectedVersion = &expected
	}
	return input, nil
}

func NormalizeUnblockPlayerInput(input UnblockPlayerInput) (UnblockPlayerInput, error) {
	var err error
	input.Actor, err = NormalizeFriendRelationshipActor(input.Actor)
	if err != nil {
		return UnblockPlayerInput{}, err
	}
	input.TargetPlayerID, err = normalizeRequired("target_player_id", input.TargetPlayerID)
	if err != nil {
		return UnblockPlayerInput{}, err
	}
	if input.Actor.PlayerID == input.TargetPlayerID {
		return UnblockPlayerInput{}, errors.New("friend relationship: self relationship is forbidden")
	}
	if input.ExpectedVersion != nil && !input.ExpectedVersion.IsValid() {
		return UnblockPlayerInput{}, errors.New("friend relationship: expected_version must be positive")
	}
	if input.ExpectedVersion != nil {
		expected := *input.ExpectedVersion
		input.ExpectedVersion = &expected
	}
	return input, nil
}

func NormalizeGetFriendRelationshipInput(input GetFriendRelationshipInput) (GetFriendRelationshipInput, error) {
	var err error
	input.Pair, err = NormalizeFriendRelationshipPair(input.Pair)
	if err != nil {
		return GetFriendRelationshipInput{}, err
	}
	return input, nil
}

func NormalizeListFriendRelationshipsInput(input ListFriendRelationshipsInput) (ListFriendRelationshipsInput, error) {
	var err error
	input.PlayerID, err = normalizeRequired("player_id", input.PlayerID)
	if err != nil {
		return ListFriendRelationshipsInput{}, err
	}
	if strings.TrimSpace(string(input.Status)) != "" {
		input.Status = FriendRelationshipStatus(strings.TrimSpace(string(input.Status)))
		if !input.Status.IsValid() {
			return ListFriendRelationshipsInput{}, errors.New("friend relationship: status is invalid")
		}
	}
	if input.Limit == 0 {
		input.Limit = DefaultListFriendRelationshipsLimit
	}
	if input.Limit < 0 || input.Limit > MaxListFriendRelationshipsLimit {
		return ListFriendRelationshipsInput{}, errors.New("friend relationship: list limit is invalid")
	}
	if strings.TrimSpace(input.AfterPairToken) != "" {
		input.AfterPairToken, err = normalizeRequired("after_pair_token", input.AfterPairToken)
		if err != nil {
			return ListFriendRelationshipsInput{}, err
		}
	}
	return input, nil
}

func NormalizeListFriendRelationshipsResult(result ListFriendRelationshipsResult) (ListFriendRelationshipsResult, error) {
	if len(result.Relationships) > 0 {
		relationships := make([]FriendRelationship, len(result.Relationships))
		for i, relationship := range result.Relationships {
			normalized, err := NormalizeFriendRelationshipRecord(relationship)
			if err != nil {
				return ListFriendRelationshipsResult{}, err
			}
			relationships[i] = normalized
		}
		result.Relationships = relationships
	}
	if strings.TrimSpace(result.NextPairToken) != "" {
		nextPairToken, err := normalizeRequired("next_pair_token", result.NextPairToken)
		if err != nil {
			return ListFriendRelationshipsResult{}, err
		}
		result.NextPairToken = nextPairToken
	}
	return result, nil
}

func NormalizeFriendRelationshipID(id FriendRelationshipID) (FriendRelationshipID, error) {
	value, err := normalizeRequired("relationship_id", string(id))
	if err != nil {
		return "", err
	}
	return FriendRelationshipID(value), nil
}

func NormalizeFriendRelationshipActor(actor FriendRelationshipActor) (FriendRelationshipActor, error) {
	var err error
	actor.PlayerID, err = normalizeRequired("actor_player_id", actor.PlayerID)
	if err != nil {
		return FriendRelationshipActor{}, err
	}
	return actor, nil
}

func NormalizeFriendRelationshipPair(pair FriendRelationshipPair) (FriendRelationshipPair, error) {
	low, err := normalizeRequired("player_low_id", pair.PlayerLowID)
	if err != nil {
		return FriendRelationshipPair{}, err
	}
	high, err := normalizeRequired("player_high_id", pair.PlayerHighID)
	if err != nil {
		return FriendRelationshipPair{}, err
	}
	if low == high {
		return FriendRelationshipPair{}, errors.New("friend relationship: self relationship is forbidden")
	}
	if high < low {
		low, high = high, low
	}
	return FriendRelationshipPair{PlayerLowID: low, PlayerHighID: high}, nil
}

func NormalizeFriendRelationshipBlockState(blockState FriendRelationshipBlockState, createdAt time.Time) (FriendRelationshipBlockState, error) {
	if blockState.BlockedByLowAt != nil {
		blockedAt := blockState.BlockedByLowAt.UTC()
		if !createdAt.IsZero() && blockedAt.Before(createdAt) {
			return FriendRelationshipBlockState{}, errors.New("friend relationship: blocked_by_low_at must not be before created_at")
		}
		blockState.BlockedByLowAt = &blockedAt
	}
	if blockState.BlockedByHighAt != nil {
		blockedAt := blockState.BlockedByHighAt.UTC()
		if !createdAt.IsZero() && blockedAt.Before(createdAt) {
			return FriendRelationshipBlockState{}, errors.New("friend relationship: blocked_by_high_at must not be before created_at")
		}
		blockState.BlockedByHighAt = &blockedAt
	}
	return blockState, nil
}

func normalizePairMutationInput[T any](actor FriendRelationshipActor, pair FriendRelationshipPair, expected *FriendRelationshipVersion, build func(FriendRelationshipActor, FriendRelationshipPair, *FriendRelationshipVersion) T) (T, error) {
	var zero T
	normalizedActor, err := NormalizeFriendRelationshipActor(actor)
	if err != nil {
		return zero, err
	}
	normalizedPair, err := NormalizeFriendRelationshipPair(pair)
	if err != nil {
		return zero, err
	}
	if !pairContainsActor(normalizedPair, normalizedActor) {
		return zero, errors.New("friend relationship: actor must be pair member")
	}
	if expected != nil && !expected.IsValid() {
		return zero, errors.New("friend relationship: expected_version must be positive")
	}
	if expected != nil {
		expectedValue := *expected
		expected = &expectedValue
	}
	return build(normalizedActor, normalizedPair, expected), nil
}

func normalizeOptionalPairMember(name string, value string, pair FriendRelationshipPair) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}
	normalized, err := normalizeRequired(name, value)
	if err != nil {
		return "", err
	}
	if normalized != pair.PlayerLowID && normalized != pair.PlayerHighID {
		return "", fmt.Errorf("friend relationship: %s must be a pair member", name)
	}
	return normalized, nil
}

func pairContainsActor(pair FriendRelationshipPair, actor FriendRelationshipActor) bool {
	return actor.PlayerID == pair.PlayerLowID || actor.PlayerID == pair.PlayerHighID
}

func normalizeRequired(name string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("friend relationship: %s is required", name)
	}
	return value, nil
}

func normalizeRequiredTime(name string, value time.Time) (time.Time, error) {
	if value.IsZero() {
		return time.Time{}, fmt.Errorf("friend relationship: %s is required", name)
	}
	return value.UTC(), nil
}
