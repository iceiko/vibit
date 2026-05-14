package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/modules/player"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrPlayerAccountNotFound   = errors.New("postgres player account: not found")
	ErrPlayerAccountConflict   = errors.New("postgres player account: conflict")
	ErrPlayerAccountConstraint = errors.New("postgres player account: constraint violation")
)

type PlayerAccountRepository struct {
	executor Executor
}

var _ player.Repository = (*PlayerAccountRepository)(nil)

func NewPlayerAccountRepositoryForUnitOfWork(executor Executor) *PlayerAccountRepository {
	return &PlayerAccountRepository{executor: executor}
}

func (r *PlayerAccountRepository) CreatePlayerAccount(ctx context.Context, mutation player.CreatePlayerAccountMutation) (player.Account, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return player.Account{}, err
	}

	normalized, err := player.NormalizeCreatePlayerAccountMutation(mutation)
	if err != nil {
		return player.Account{}, fmt.Errorf("postgres player account: normalize create mutation: %w", err)
	}

	account, err := scanPlayerAccountRow(executor.QueryRow(
		ctx,
		insertPlayerAccountSQL,
		normalized.PlayerID,
		normalized.DisplayName,
		string(normalized.AccountState),
		normalized.OccurredAt,
	))
	if err != nil {
		return player.Account{}, mapPlayerAccountPostgresError("insert player account", err)
	}

	if _, err := executor.Exec(
		ctx,
		insertPlayerAccountEventSQL,
		normalized.EventID,
		player.EventPlayerAccountCreated,
		normalized.OccurredAt,
		account.PlayerID,
		normalized.RequestedBy,
		string(account.AccountState),
		account.DisplayName,
	); err != nil {
		return player.Account{}, mapPlayerAccountPostgresError("insert PlayerAccountCreated event", err)
	}

	return account, nil
}

func (r *PlayerAccountRepository) GetPlayerAccount(ctx context.Context, playerID string) (player.Account, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return player.Account{}, err
	}

	playerID, err = player.NormalizePlayerID(playerID)
	if err != nil {
		return player.Account{}, fmt.Errorf("postgres player account: normalize player_id: %w", err)
	}

	account, err := scanPlayerAccountRow(executor.QueryRow(ctx, getPlayerAccountSQL, playerID))
	if err != nil {
		return player.Account{}, mapPlayerAccountPostgresError("get player account", err)
	}
	return account, nil
}

func (r *PlayerAccountRepository) requireExecutor() (Executor, error) {
	if r == nil || r.executor == nil {
		return nil, errors.New("postgres player account: unit-of-work executor is required")
	}
	return r.executor, nil
}

func scanPlayerAccountRow(row pgx.Row) (player.Account, error) {
	var account player.Account
	var state string
	var disabledAt pgtype.Timestamptz
	var deletedAt pgtype.Timestamptz

	if err := row.Scan(
		&account.PlayerID,
		&account.DisplayName,
		&state,
		&account.CreatedAt,
		&account.UpdatedAt,
		&disabledAt,
		&deletedAt,
	); err != nil {
		return player.Account{}, err
	}

	normalizedPlayerID, err := player.NormalizePlayerID(account.PlayerID)
	if err != nil {
		return player.Account{}, fmt.Errorf("%w: row player_id: %v", ErrPlayerAccountConstraint, err)
	}
	account.PlayerID = normalizedPlayerID
	if strings.TrimSpace(account.DisplayName) == "" {
		return player.Account{}, fmt.Errorf("%w: row display_name is required", ErrPlayerAccountConstraint)
	}

	account.AccountState = player.AccountState(state)
	if !account.AccountState.IsValid() {
		return player.Account{}, fmt.Errorf("%w: row account_state %q is invalid", ErrPlayerAccountConstraint, state)
	}
	if account.CreatedAt.IsZero() {
		return player.Account{}, fmt.Errorf("%w: row created_at is required", ErrPlayerAccountConstraint)
	}
	if account.UpdatedAt.IsZero() {
		return player.Account{}, fmt.Errorf("%w: row updated_at is required", ErrPlayerAccountConstraint)
	}

	account.CreatedAt = account.CreatedAt.UTC()
	account.UpdatedAt = account.UpdatedAt.UTC()
	account.DisabledAt = nullableTimestamptzUTC(disabledAt)
	account.DeletedAt = nullableTimestamptzUTC(deletedAt)
	return account, nil
}

func nullableTimestamptzUTC(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	timestamp := value.Time.UTC()
	return &timestamp
}

func mapPlayerAccountPostgresError(action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrPlayerAccountNotFound, action)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s: %s", ErrPlayerAccountConflict, action, postgresConstraintLabel(pgErr))
		case "23502", "23503", "23514":
			return fmt.Errorf("%w: %s: %s", ErrPlayerAccountConstraint, action, postgresConstraintLabel(pgErr))
		default:
			return fmt.Errorf("postgres player account: %s failed with PostgreSQL code %s", action, pgErr.Code)
		}
	}

	if errors.Is(err, ErrPlayerAccountConstraint) {
		return err
	}
	return fmt.Errorf("postgres player account: %s: %w", action, err)
}

func postgresConstraintLabel(err *pgconn.PgError) string {
	if err == nil {
		return "unknown_constraint"
	}
	if err.ConstraintName != "" {
		return err.ConstraintName
	}
	if err.Code != "" {
		return err.Code
	}
	return "unknown_constraint"
}

const insertPlayerAccountSQL = `
INSERT INTO player_accounts (
    player_id,
    display_name,
    account_state,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $4)
RETURNING player_id, display_name, account_state, created_at, updated_at, disabled_at, deleted_at`

const insertPlayerAccountEventSQL = `
INSERT INTO player_account_events (
    event_id,
    event_type,
    occurred_at,
    player_id,
    requested_by,
    account_state,
    display_name
)
VALUES ($1, $2, $3, $4, $5, $6, $7)`

const getPlayerAccountSQL = `
SELECT player_id, display_name, account_state, created_at, updated_at, disabled_at, deleted_at
FROM player_accounts
WHERE player_id = $1`
