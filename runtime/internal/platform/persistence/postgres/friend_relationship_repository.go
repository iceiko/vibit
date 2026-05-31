package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/modules/friends"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type FriendRelationshipRepository struct {
	executor Executor
}

var _ friends.Repository = (*FriendRelationshipRepository)(nil)

func NewFriendRelationshipRepositoryForUnitOfWork(executor Executor) *FriendRelationshipRepository {
	return &FriendRelationshipRepository{executor: executor}
}

func (r *FriendRelationshipRepository) CreateOrUpdateFriendRequest(ctx context.Context, input friends.SendFriendRequestInput) (friends.FriendRelationship, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return friends.FriendRelationship{}, err
	}

	normalized, err := friends.NormalizeSendFriendRequestInput(input)
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("create_or_update_request", classifyFriendRelationshipInputError(err))
	}
	pair, err := friends.NormalizeFriendRelationshipPair(friends.FriendRelationshipPair{
		PlayerLowID:  normalized.Actor.PlayerID,
		PlayerHighID: normalized.TargetPlayerID,
	})
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("create_or_update_request", classifyFriendRelationshipInputError(err))
	}

	record, err := scanFriendRelationshipRow(executor.QueryRow(
		ctx,
		upsertFriendRequestSQL,
		normalized.RelationshipID,
		pair.PlayerLowID,
		pair.PlayerHighID,
		string(friends.FriendRelationshipLifecyclePending),
		normalized.Actor.PlayerID,
		int64(friends.InitialFriendRelationshipVersion),
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return friends.FriendRelationship{}, r.mapNoRowsForRelationshipMutation(ctx, executor, "create_or_update_request", pair, nil)
		}
		return friends.FriendRelationship{}, mapFriendRelationshipPostgresError("create_or_update_request", err)
	}
	return record, nil
}

func (r *FriendRelationshipRepository) GetRelationshipByPair(ctx context.Context, input friends.GetFriendRelationshipInput) (friends.FriendRelationship, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return friends.FriendRelationship{}, err
	}

	normalized, err := friends.NormalizeGetFriendRelationshipInput(input)
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("get_by_pair", classifyFriendRelationshipInputError(err))
	}

	record, err := scanFriendRelationshipRow(executor.QueryRow(
		ctx,
		getFriendRelationshipByPairSQL,
		normalized.Pair.PlayerLowID,
		normalized.Pair.PlayerHighID,
	))
	if err != nil {
		return friends.FriendRelationship{}, mapFriendRelationshipPostgresError("get_by_pair", err)
	}
	return record, nil
}

func (r *FriendRelationshipRepository) ListRelationshipsForPlayer(ctx context.Context, input friends.ListFriendRelationshipsInput) (friends.ListFriendRelationshipsResult, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return friends.ListFriendRelationshipsResult{}, err
	}

	normalized, err := friends.NormalizeListFriendRelationshipsInput(input)
	if err != nil {
		return friends.ListFriendRelationshipsResult{}, friendRelationshipInvalidInputError("list_for_player", classifyFriendRelationshipInputError(err))
	}

	rows, err := executor.Query(
		ctx,
		listFriendRelationshipsForPlayerSQL,
		normalized.PlayerID,
		normalized.AfterPairToken,
		string(normalized.Status),
		int32(normalized.Limit+1),
	)
	if err != nil {
		return friends.ListFriendRelationshipsResult{}, mapFriendRelationshipPostgresError("list_for_player", err)
	}
	defer rows.Close()

	relationships := make([]friends.FriendRelationship, 0, normalized.Limit)
	nextPairToken := ""
	for rows.Next() {
		record, err := scanFriendRelationshipScanner(rows)
		if err != nil {
			return friends.ListFriendRelationshipsResult{}, mapFriendRelationshipPostgresError("list_for_player", err)
		}
		if len(relationships) >= normalized.Limit {
			nextPairToken = friendRelationshipPairToken(record.Pair)
			continue
		}
		relationships = append(relationships, record)
	}
	if err := rows.Err(); err != nil {
		return friends.ListFriendRelationshipsResult{}, mapFriendRelationshipPostgresError("list_for_player", err)
	}
	return friends.NormalizeListFriendRelationshipsResult(friends.ListFriendRelationshipsResult{
		Relationships: relationships,
		NextPairToken: nextPairToken,
	})
}

func (r *FriendRelationshipRepository) AcceptFriendRequest(ctx context.Context, input friends.AcceptFriendRequestInput) (friends.FriendRelationship, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return friends.FriendRelationship{}, err
	}

	normalized, err := friends.NormalizeAcceptFriendRequestInput(input)
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("accept_request", classifyFriendRelationshipInputError(err))
	}

	query := acceptFriendRequestSQL
	args := []any{
		normalized.Pair.PlayerLowID,
		normalized.Pair.PlayerHighID,
		normalized.Actor.PlayerID,
		string(friends.FriendRelationshipLifecycleFriends),
	}
	if normalized.ExpectedVersion != nil {
		query = acceptFriendRequestWithExpectedVersionSQL
		args = append(args, int64(*normalized.ExpectedVersion))
	}
	record, err := scanFriendRelationshipRow(executor.QueryRow(
		ctx,
		query,
		args...,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return friends.FriendRelationship{}, r.mapNoRowsForRelationshipMutation(ctx, executor, "accept_request", normalized.Pair, normalized.ExpectedVersion)
		}
		return friends.FriendRelationship{}, mapFriendRelationshipPostgresError("accept_request", err)
	}
	return record, nil
}

func (r *FriendRelationshipRepository) RejectFriendRequest(ctx context.Context, input friends.RejectFriendRequestInput) (friends.FriendRelationship, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return friends.FriendRelationship{}, err
	}

	normalized, err := friends.NormalizeRejectFriendRequestInput(input)
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("reject_request", classifyFriendRelationshipInputError(err))
	}

	query := rejectFriendRequestSQL
	args := []any{
		normalized.Pair.PlayerLowID,
		normalized.Pair.PlayerHighID,
		normalized.Actor.PlayerID,
		string(friends.FriendRelationshipLifecycleRejected),
	}
	if normalized.ExpectedVersion != nil {
		query = rejectFriendRequestWithExpectedVersionSQL
		args = append(args, int64(*normalized.ExpectedVersion))
	}
	record, err := scanFriendRelationshipRow(executor.QueryRow(
		ctx,
		query,
		args...,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return friends.FriendRelationship{}, r.mapNoRowsForRelationshipMutation(ctx, executor, "reject_request", normalized.Pair, normalized.ExpectedVersion)
		}
		return friends.FriendRelationship{}, mapFriendRelationshipPostgresError("reject_request", err)
	}
	return record, nil
}

func (r *FriendRelationshipRepository) RemoveFriend(ctx context.Context, input friends.RemoveFriendInput) (friends.FriendRelationship, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return friends.FriendRelationship{}, err
	}

	normalized, err := friends.NormalizeRemoveFriendInput(input)
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("remove_friend", classifyFriendRelationshipInputError(err))
	}

	query := removeFriendSQL
	args := []any{
		normalized.Pair.PlayerLowID,
		normalized.Pair.PlayerHighID,
		normalized.Actor.PlayerID,
		string(friends.FriendRelationshipLifecycleRemoved),
	}
	if normalized.ExpectedVersion != nil {
		query = removeFriendWithExpectedVersionSQL
		args = append(args, int64(*normalized.ExpectedVersion))
	}
	record, err := scanFriendRelationshipRow(executor.QueryRow(
		ctx,
		query,
		args...,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return friends.FriendRelationship{}, r.mapNoRowsForRelationshipMutation(ctx, executor, "remove_friend", normalized.Pair, normalized.ExpectedVersion)
		}
		return friends.FriendRelationship{}, mapFriendRelationshipPostgresError("remove_friend", err)
	}
	return record, nil
}

func (r *FriendRelationshipRepository) SetPlayerBlock(ctx context.Context, input friends.BlockPlayerInput) (friends.FriendRelationship, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return friends.FriendRelationship{}, err
	}

	normalized, err := friends.NormalizeBlockPlayerInput(input)
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("set_player_block", classifyFriendRelationshipInputError(err))
	}
	pair, err := friends.NormalizeFriendRelationshipPair(friends.FriendRelationshipPair{
		PlayerLowID:  normalized.Actor.PlayerID,
		PlayerHighID: normalized.TargetPlayerID,
	})
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("set_player_block", classifyFriendRelationshipInputError(err))
	}

	query := setPlayerHighBlockSQL
	if normalized.Actor.PlayerID == pair.PlayerLowID {
		query = setPlayerLowBlockSQL
	}
	record, err := scanFriendRelationshipRow(executor.QueryRow(
		ctx,
		query,
		pair.PlayerLowID,
		pair.PlayerHighID,
		normalized.Actor.PlayerID,
		expectedFriendRelationshipVersionArg(normalized.ExpectedVersion),
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return friends.FriendRelationship{}, r.mapNoRowsForRelationshipMutation(ctx, executor, "set_player_block", pair, normalized.ExpectedVersion)
		}
		return friends.FriendRelationship{}, mapFriendRelationshipPostgresError("set_player_block", err)
	}
	return record, nil
}

func (r *FriendRelationshipRepository) ClearPlayerBlock(ctx context.Context, input friends.UnblockPlayerInput) (friends.FriendRelationship, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return friends.FriendRelationship{}, err
	}

	normalized, err := friends.NormalizeUnblockPlayerInput(input)
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("clear_player_block", classifyFriendRelationshipInputError(err))
	}
	pair, err := friends.NormalizeFriendRelationshipPair(friends.FriendRelationshipPair{
		PlayerLowID:  normalized.Actor.PlayerID,
		PlayerHighID: normalized.TargetPlayerID,
	})
	if err != nil {
		return friends.FriendRelationship{}, friendRelationshipInvalidInputError("clear_player_block", classifyFriendRelationshipInputError(err))
	}

	query := clearPlayerHighBlockSQL
	if normalized.Actor.PlayerID == pair.PlayerLowID {
		query = clearPlayerLowBlockSQL
	}
	record, err := scanFriendRelationshipRow(executor.QueryRow(
		ctx,
		query,
		pair.PlayerLowID,
		pair.PlayerHighID,
		expectedFriendRelationshipVersionArg(normalized.ExpectedVersion),
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return friends.FriendRelationship{}, r.mapNoRowsForRelationshipMutation(ctx, executor, "clear_player_block", pair, normalized.ExpectedVersion)
		}
		return friends.FriendRelationship{}, mapFriendRelationshipPostgresError("clear_player_block", err)
	}
	return record, nil
}

func (r *FriendRelationshipRepository) requireExecutor() (Executor, error) {
	if r == nil || r.executor == nil {
		return nil, errors.New("postgres friend relationship: unit-of-work executor is required")
	}
	return r.executor, nil
}

func (r *FriendRelationshipRepository) mapNoRowsForRelationshipMutation(ctx context.Context, executor Executor, operation string, pair friends.FriendRelationshipPair, expectedVersion *friends.FriendRelationshipVersion) error {
	current, err := scanFriendRelationshipRow(executor.QueryRow(
		ctx,
		getFriendRelationshipByPairSQL,
		pair.PlayerLowID,
		pair.PlayerHighID,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictNotFound, false)
	}
	if err != nil {
		return mapFriendRelationshipPostgresError(operation, err)
	}
	if current.BlockState.BlockedByLowAt != nil || current.BlockState.BlockedByHighAt != nil {
		return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictBlockedRelationship, false)
	}
	if expectedVersion != nil && current.Version != *expectedVersion {
		return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictStaleVersion, true)
	}
	switch operation {
	case "create_or_update_request":
		switch current.LifecycleState {
		case friends.FriendRelationshipLifecyclePending:
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictDuplicatePendingRequest, false)
		case friends.FriendRelationshipLifecycleFriends:
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictAlreadyFriends, false)
		default:
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictInvalidTransition, false)
		}
	case "accept_request", "reject_request":
		if current.LifecycleState != friends.FriendRelationshipLifecyclePending {
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictInvalidTransition, false)
		}
	case "remove_friend":
		if current.LifecycleState == friends.FriendRelationshipLifecycleRemoved || current.LifecycleState == friends.FriendRelationshipLifecycleRejected {
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictInvalidTransition, false)
		}
	}
	return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictVersionMismatch, true)
}

func scanFriendRelationshipRow(row pgx.Row) (friends.FriendRelationship, error) {
	return scanFriendRelationshipScanner(row)
}

func scanFriendRelationshipScanner(row scanner) (friends.FriendRelationship, error) {
	var record friends.FriendRelationship
	var relationshipID string
	var lifecycleState string
	var requestedBy pgtype.Text
	var respondedBy pgtype.Text
	var removedBy pgtype.Text
	var version int64
	var rejectedAt pgtype.Timestamptz
	var removedAt pgtype.Timestamptz
	var blockedByLowAt pgtype.Timestamptz
	var blockedByHighAt pgtype.Timestamptz

	if err := row.Scan(
		&relationshipID,
		&record.Pair.PlayerLowID,
		&record.Pair.PlayerHighID,
		&lifecycleState,
		&requestedBy,
		&respondedBy,
		&removedBy,
		&version,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.StateChangedAt,
		&rejectedAt,
		&removedAt,
		&blockedByLowAt,
		&blockedByHighAt,
	); err != nil {
		return friends.FriendRelationship{}, err
	}

	record.RelationshipID = friends.FriendRelationshipID(strings.TrimSpace(relationshipID))
	record.LifecycleState = friends.FriendRelationshipLifecycleState(strings.TrimSpace(lifecycleState))
	record.RequestedByPlayerID = nullableTextValue(requestedBy)
	record.RespondedByPlayerID = nullableTextValue(respondedBy)
	record.RemovedByPlayerID = nullableTextValue(removedBy)
	record.Version = friends.FriendRelationshipVersion(version)
	record.RejectedAt = nullableTimestamptzUTC(rejectedAt)
	record.RemovedAt = nullableTimestamptzUTC(removedAt)
	record.BlockState = friends.FriendRelationshipBlockState{
		BlockedByLowAt:  nullableTimestamptzUTC(blockedByLowAt),
		BlockedByHighAt: nullableTimestamptzUTC(blockedByHighAt),
	}

	normalized, err := friends.NormalizeFriendRelationshipRecord(record)
	if err != nil {
		return friends.FriendRelationship{}, fmt.Errorf("%w: row shape", friends.ErrFriendRelationshipUnavailable)
	}
	return normalized, nil
}

func mapFriendRelationshipPostgresError(operation string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictNotFound, false)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictPairIdentity, false)
		case "23503":
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictTargetPlayerNotFound, false)
		case "23502", "23514":
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipConflict, friends.FriendRelationshipConflictInvalidTransition, false)
		default:
			return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipUnavailable, friends.FriendRelationshipConflictStorageUnavailable, true)
		}
	}

	if errors.Is(err, friends.ErrFriendRelationshipUnavailable) {
		return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipUnavailable, friends.FriendRelationshipConflictStorageUnavailable, true)
	}
	return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipUnavailable, friends.FriendRelationshipConflictStorageUnavailable, true)
}

func friendRelationshipInvalidInputError(operation string, class friends.FriendRelationshipConflictClass) error {
	return friendRelationshipRepositoryError(operation, friends.ErrFriendRelationshipInvalidInput, class, false)
}

func friendRelationshipRepositoryError(operation string, kind error, class friends.FriendRelationshipConflictClass, retryable bool) error {
	reason := string(class)
	if reason == "" && kind != nil {
		reason = kind.Error()
	}
	return &friends.FriendRelationshipRepositoryError{
		Kind: kind,
		Conflict: friends.FriendRelationshipConflict{
			Class:          class,
			Retryable:      retryable,
			RedactedReason: reason,
		},
		Operation:      operation,
		RedactedReason: reason,
	}
}

func classifyFriendRelationshipInputError(err error) friends.FriendRelationshipConflictClass {
	text := strings.ToLower(err.Error())
	switch {
	case strings.Contains(text, "self relationship"):
		return friends.FriendRelationshipConflictSelfRelationshipForbidden
	case strings.Contains(text, "expected_version"):
		return friends.FriendRelationshipConflictVersionMismatch
	case strings.Contains(text, "actor must be pair member"):
		return friends.FriendRelationshipConflictInvalidTransition
	case strings.Contains(text, "pair"):
		return friends.FriendRelationshipConflictPairIdentity
	default:
		return friends.FriendRelationshipConflictInvalidTransition
	}
}

func expectedFriendRelationshipVersionArg(expected *friends.FriendRelationshipVersion) any {
	if expected == nil {
		return nil
	}
	return int64(*expected)
}

func friendRelationshipPairToken(pair friends.FriendRelationshipPair) string {
	return pair.PlayerLowID + "|" + pair.PlayerHighID
}

const friendRelationshipColumns = `
relationship_id,
player_low_id,
player_high_id,
lifecycle_state,
requested_by_player_id,
responded_by_player_id,
removed_by_player_id,
relationship_version,
created_at,
updated_at,
state_changed_at,
rejected_at,
removed_at,
blocked_by_low_at,
blocked_by_high_at`

const upsertFriendRequestSQL = `
INSERT INTO friend_relationships (
    relationship_id,
    player_low_id,
    player_high_id,
    lifecycle_state,
    requested_by_player_id,
    relationship_version,
    created_at,
    updated_at,
    state_changed_at
)
VALUES ($1, $2, $3, $4, $5, $6, now(), now(), now())
ON CONFLICT (player_low_id, player_high_id) DO UPDATE
SET relationship_id = EXCLUDED.relationship_id,
    lifecycle_state = EXCLUDED.lifecycle_state,
    requested_by_player_id = EXCLUDED.requested_by_player_id,
    responded_by_player_id = NULL,
    removed_by_player_id = NULL,
    relationship_version = friend_relationships.relationship_version + 1,
    updated_at = now(),
    state_changed_at = now(),
    rejected_at = NULL,
    removed_at = NULL
WHERE friend_relationships.lifecycle_state IN ('rejected', 'removed')
  AND friend_relationships.blocked_by_low_at IS NULL
  AND friend_relationships.blocked_by_high_at IS NULL
RETURNING ` + friendRelationshipColumns

const getFriendRelationshipByPairSQL = `
SELECT ` + friendRelationshipColumns + `
FROM friend_relationships
WHERE player_low_id = $1
  AND player_high_id = $2`

const listFriendRelationshipsForPlayerSQL = `
SELECT ` + friendRelationshipColumns + `
FROM friend_relationships
WHERE (player_low_id = $1 OR player_high_id = $1)
  AND (player_low_id || '|' || player_high_id) > $2
  AND CASE $3
      WHEN 'pending' THEN lifecycle_state = 'pending' AND blocked_by_low_at IS NULL AND blocked_by_high_at IS NULL
      WHEN 'friends' THEN lifecycle_state = 'friends' AND blocked_by_low_at IS NULL AND blocked_by_high_at IS NULL
      WHEN 'blocked' THEN blocked_by_low_at IS NOT NULL OR blocked_by_high_at IS NOT NULL
      WHEN 'ended' THEN lifecycle_state IN ('rejected', 'removed') AND blocked_by_low_at IS NULL AND blocked_by_high_at IS NULL
      ELSE TRUE
  END
ORDER BY player_low_id, player_high_id
LIMIT $4`

const acceptFriendRequestSQL = `
UPDATE friend_relationships
SET lifecycle_state = $4,
    responded_by_player_id = $3,
    removed_by_player_id = NULL,
    relationship_version = relationship_version + 1,
    updated_at = now(),
    state_changed_at = now(),
    rejected_at = NULL,
    removed_at = NULL
WHERE player_low_id = $1
  AND player_high_id = $2
  AND lifecycle_state = 'pending'
  AND blocked_by_low_at IS NULL
  AND blocked_by_high_at IS NULL
RETURNING ` + friendRelationshipColumns

const acceptFriendRequestWithExpectedVersionSQL = `
UPDATE friend_relationships
SET lifecycle_state = $4,
    responded_by_player_id = $3,
    removed_by_player_id = NULL,
    relationship_version = relationship_version + 1,
    updated_at = now(),
    state_changed_at = now(),
    rejected_at = NULL,
    removed_at = NULL
WHERE player_low_id = $1
  AND player_high_id = $2
  AND lifecycle_state = 'pending'
  AND blocked_by_low_at IS NULL
  AND blocked_by_high_at IS NULL
  AND relationship_version = $5
RETURNING ` + friendRelationshipColumns

const rejectFriendRequestSQL = `
UPDATE friend_relationships
SET lifecycle_state = $4,
    responded_by_player_id = $3,
    removed_by_player_id = NULL,
    relationship_version = relationship_version + 1,
    updated_at = now(),
    state_changed_at = now(),
    rejected_at = now(),
    removed_at = NULL
WHERE player_low_id = $1
  AND player_high_id = $2
  AND lifecycle_state = 'pending'
  AND blocked_by_low_at IS NULL
  AND blocked_by_high_at IS NULL
RETURNING ` + friendRelationshipColumns

const rejectFriendRequestWithExpectedVersionSQL = `
UPDATE friend_relationships
SET lifecycle_state = $4,
    responded_by_player_id = $3,
    removed_by_player_id = NULL,
    relationship_version = relationship_version + 1,
    updated_at = now(),
    state_changed_at = now(),
    rejected_at = now(),
    removed_at = NULL
WHERE player_low_id = $1
  AND player_high_id = $2
  AND lifecycle_state = 'pending'
  AND blocked_by_low_at IS NULL
  AND blocked_by_high_at IS NULL
  AND relationship_version = $5
RETURNING ` + friendRelationshipColumns

const removeFriendSQL = `
UPDATE friend_relationships
SET lifecycle_state = $4,
    removed_by_player_id = $3,
    relationship_version = relationship_version + 1,
    updated_at = now(),
    state_changed_at = now(),
    removed_at = now()
WHERE player_low_id = $1
  AND player_high_id = $2
  AND lifecycle_state IN ('pending', 'friends')
  AND blocked_by_low_at IS NULL
  AND blocked_by_high_at IS NULL
RETURNING ` + friendRelationshipColumns

const removeFriendWithExpectedVersionSQL = `
UPDATE friend_relationships
SET lifecycle_state = $4,
    removed_by_player_id = $3,
    relationship_version = relationship_version + 1,
    updated_at = now(),
    state_changed_at = now(),
    removed_at = now()
WHERE player_low_id = $1
  AND player_high_id = $2
  AND lifecycle_state IN ('pending', 'friends')
  AND blocked_by_low_at IS NULL
  AND blocked_by_high_at IS NULL
  AND relationship_version = $5
RETURNING ` + friendRelationshipColumns

const setPlayerLowBlockSQL = `
UPDATE friend_relationships
SET blocked_by_low_at = now(),
    lifecycle_state = CASE
        WHEN lifecycle_state IN ('pending', 'friends') THEN 'removed'
        ELSE lifecycle_state
    END,
    removed_by_player_id = CASE
        WHEN lifecycle_state IN ('pending', 'friends') THEN $3
        ELSE removed_by_player_id
    END,
    removed_at = CASE
        WHEN lifecycle_state IN ('pending', 'friends') AND removed_at IS NULL THEN now()
        ELSE removed_at
    END,
    relationship_version = relationship_version + 1,
    updated_at = now(),
    state_changed_at = now()
WHERE player_low_id = $1
  AND player_high_id = $2
  AND ($4::bigint IS NULL OR relationship_version = $4)
RETURNING ` + friendRelationshipColumns

const setPlayerHighBlockSQL = `
UPDATE friend_relationships
SET blocked_by_high_at = now(),
    lifecycle_state = CASE
        WHEN lifecycle_state IN ('pending', 'friends') THEN 'removed'
        ELSE lifecycle_state
    END,
    removed_by_player_id = CASE
        WHEN lifecycle_state IN ('pending', 'friends') THEN $3
        ELSE removed_by_player_id
    END,
    removed_at = CASE
        WHEN lifecycle_state IN ('pending', 'friends') AND removed_at IS NULL THEN now()
        ELSE removed_at
    END,
    relationship_version = relationship_version + 1,
    updated_at = now(),
    state_changed_at = now()
WHERE player_low_id = $1
  AND player_high_id = $2
  AND ($4::bigint IS NULL OR relationship_version = $4)
RETURNING ` + friendRelationshipColumns

const clearPlayerLowBlockSQL = `
UPDATE friend_relationships
SET blocked_by_low_at = NULL,
    lifecycle_state = lifecycle_state,
    relationship_version = relationship_version + 1,
    updated_at = now()
WHERE player_low_id = $1
  AND player_high_id = $2
  AND ($3::bigint IS NULL OR relationship_version = $3)
RETURNING ` + friendRelationshipColumns

const clearPlayerHighBlockSQL = `
UPDATE friend_relationships
SET blocked_by_high_at = NULL,
    lifecycle_state = lifecycle_state,
    relationship_version = relationship_version + 1,
    updated_at = now()
WHERE player_low_id = $1
  AND player_high_id = $2
  AND ($3::bigint IS NULL OR relationship_version = $3)
RETURNING ` + friendRelationshipColumns
