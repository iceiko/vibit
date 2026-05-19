package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app/session"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrSessionRecordNotFound   = errors.New("postgres session: not found")
	ErrSessionRecordConflict   = errors.New("postgres session: conflict")
	ErrSessionRecordConstraint = errors.New("postgres session: constraint violation")
	ErrSessionRecordStale      = errors.New("postgres session: stale lifecycle state")
)

type SessionRepository struct {
	executor Executor
}

var _ session.Repository = (*SessionRepository)(nil)

func NewSessionRepositoryForUnitOfWork(executor Executor) *SessionRepository {
	return &SessionRepository{executor: executor}
}

func (r *SessionRepository) CreateRuntimeSession(ctx context.Context, mutation session.CreateRuntimeSessionMutation) (session.RuntimeSession, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return session.RuntimeSession{}, err
	}

	normalized, err := session.NormalizeCreateRuntimeSessionMutation(mutation)
	if err != nil {
		return session.RuntimeSession{}, fmt.Errorf("postgres session: normalize create runtime session mutation: %w", err)
	}

	record, err := scanSessionRow(executor.QueryRow(
		ctx,
		insertRuntimeSessionSQL,
		normalized.SessionID,
		string(normalized.ActorKind),
		normalized.ActorID,
		normalized.PlayerID,
		string(normalized.SessionStatus),
		normalized.IssuedAt,
		normalized.ExpiresAt,
		normalized.LastSeenAt,
		nullableText(normalized.AccessTokenRecordID),
	))
	if err != nil {
		return session.RuntimeSession{}, mapSessionPostgresError("create runtime session", err)
	}
	return record, nil
}

func (r *SessionRepository) GetRuntimeSession(ctx context.Context, query session.GetRuntimeSessionQuery) (session.RuntimeSession, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return session.RuntimeSession{}, err
	}

	normalized, err := session.NormalizeGetRuntimeSessionQuery(query)
	if err != nil {
		return session.RuntimeSession{}, fmt.Errorf("postgres session: normalize get runtime session query: %w", err)
	}

	record, err := scanSessionRow(executor.QueryRow(ctx, getRuntimeSessionSQL, normalized.SessionID))
	if err != nil {
		return session.RuntimeSession{}, mapSessionPostgresError("get runtime session", err)
	}
	return record, nil
}

func (r *SessionRepository) FindActiveSessionByID(ctx context.Context, query session.FindActiveSessionByIDQuery) (session.RuntimeSession, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return session.RuntimeSession{}, err
	}

	normalized, err := session.NormalizeFindActiveSessionByIDQuery(query)
	if err != nil {
		return session.RuntimeSession{}, fmt.Errorf("postgres session: normalize find active runtime session query: %w", err)
	}

	record, err := scanSessionRow(executor.QueryRow(
		ctx,
		findActiveRuntimeSessionSQL,
		normalized.SessionID,
		normalized.ObservedAt,
		string(session.SessionStatusActive),
	))
	if err != nil {
		return session.RuntimeSession{}, mapSessionPostgresError("find active runtime session", err)
	}
	return record, nil
}

func (r *SessionRepository) UpdateRuntimeSessionLastSeen(ctx context.Context, mutation session.UpdateRuntimeSessionLastSeenMutation) (session.RuntimeSession, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return session.RuntimeSession{}, err
	}

	normalized, err := session.NormalizeUpdateRuntimeSessionLastSeenMutation(mutation)
	if err != nil {
		return session.RuntimeSession{}, fmt.Errorf("postgres session: normalize update last seen mutation: %w", err)
	}

	record, err := scanSessionRow(executor.QueryRow(
		ctx,
		updateRuntimeSessionLastSeenSQL,
		normalized.SessionID,
		normalized.LastSeenAt,
	))
	if err != nil {
		return session.RuntimeSession{}, mapSessionPostgresError("update runtime session last_seen_at", err)
	}
	return record, nil
}

func (r *SessionRepository) MarkRuntimeSessionExpired(ctx context.Context, mutation session.MarkRuntimeSessionExpiredMutation) (session.RuntimeSession, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return session.RuntimeSession{}, err
	}

	normalized, err := session.NormalizeMarkRuntimeSessionExpiredMutation(mutation)
	if err != nil {
		return session.RuntimeSession{}, fmt.Errorf("postgres session: normalize mark expired mutation: %w", err)
	}

	record, err := scanSessionRow(executor.QueryRow(
		ctx,
		markRuntimeSessionExpiredSQL,
		normalized.SessionID,
		normalized.ExpiredAt,
		string(session.SessionStatusExpired),
		string(session.SessionStatusActive),
	))
	if err != nil {
		return session.RuntimeSession{}, mapSessionPostgresError("mark runtime session expired", err)
	}
	return record, nil
}

func (r *SessionRepository) RevokeRuntimeSession(ctx context.Context, mutation session.RevokeRuntimeSessionMutation) (session.RuntimeSession, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return session.RuntimeSession{}, err
	}

	normalized, err := session.NormalizeRevokeRuntimeSessionMutation(mutation)
	if err != nil {
		return session.RuntimeSession{}, fmt.Errorf("postgres session: normalize revoke runtime session mutation: %w", err)
	}

	record, err := scanSessionRow(executor.QueryRow(
		ctx,
		revokeRuntimeSessionSQL,
		normalized.SessionID,
		normalized.RevokedAt,
		normalized.RevocationReason,
		string(session.SessionStatusRevoked),
	))
	if err != nil {
		return session.RuntimeSession{}, mapSessionPostgresError("revoke runtime session", err)
	}
	return record, nil
}

func (r *SessionRepository) ListActiveSessionsForPlayer(ctx context.Context, query session.ListActiveSessionsForPlayerQuery) ([]session.RuntimeSession, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return nil, err
	}

	normalized, err := session.NormalizeListActiveSessionsForPlayerQuery(query)
	if err != nil {
		return nil, fmt.Errorf("postgres session: normalize list active runtime sessions query: %w", err)
	}

	rows, err := executor.Query(
		ctx,
		listActiveRuntimeSessionsForPlayerSQL,
		normalized.PlayerID,
		normalized.ObservedAt,
		string(session.SessionStatusActive),
		int32(normalized.Limit),
	)
	if err != nil {
		return nil, mapSessionPostgresError("list active runtime sessions for player", err)
	}
	defer rows.Close()

	records := []session.RuntimeSession{}
	for rows.Next() {
		record, err := scanSessionScanner(rows)
		if err != nil {
			return nil, mapSessionPostgresError("scan active runtime session", err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, mapSessionPostgresError("read active runtime sessions", err)
	}
	return records, nil
}

func (r *SessionRepository) requireExecutor() (Executor, error) {
	if r == nil || r.executor == nil {
		return nil, errors.New("postgres session: unit-of-work executor is required")
	}
	return r.executor, nil
}

func scanSessionRow(row pgx.Row) (session.RuntimeSession, error) {
	return scanSessionScanner(row)
}

func scanSessionScanner(row scanner) (session.RuntimeSession, error) {
	var record session.RuntimeSession
	var actorKind string
	var status string
	var revokedAt pgtype.Timestamptz
	var revocationReason pgtype.Text
	var accessTokenRecordID pgtype.Text

	if err := row.Scan(
		&record.SessionID,
		&actorKind,
		&record.ActorID,
		&record.PlayerID,
		&status,
		&record.IssuedAt,
		&record.ExpiresAt,
		&record.LastSeenAt,
		&revokedAt,
		&revocationReason,
		&accessTokenRecordID,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return session.RuntimeSession{}, err
	}

	record.ActorKind = session.ActorKind(strings.TrimSpace(actorKind))
	record.SessionStatus = session.SessionStatus(strings.TrimSpace(status))
	record.RevokedAt = nullableTimestamptzUTC(revokedAt)
	record.RevocationReason = nullableTextValue(revocationReason)
	record.AccessTokenRecordID = nullableTextValue(accessTokenRecordID)

	normalized, err := session.NormalizeRuntimeSessionRecord(record)
	if err != nil {
		return session.RuntimeSession{}, fmt.Errorf("%w: row shape: %v", ErrSessionRecordConstraint, err)
	}
	return normalized, nil
}

func mapSessionPostgresError(action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		if action == "mark runtime session expired" {
			return fmt.Errorf("%w: %s", ErrSessionRecordStale, action)
		}
		return fmt.Errorf("%w: %s", ErrSessionRecordNotFound, action)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s: %s", ErrSessionRecordConflict, action, postgresConstraintLabel(pgErr))
		case "23502", "23503", "23514":
			return fmt.Errorf("%w: %s: %s", ErrSessionRecordConstraint, action, postgresConstraintLabel(pgErr))
		default:
			return fmt.Errorf("postgres session: %s failed with PostgreSQL code %s", action, pgErr.Code)
		}
	}

	if errors.Is(err, ErrSessionRecordConstraint) || errors.Is(err, ErrSessionRecordStale) {
		return err
	}
	return fmt.Errorf("postgres session: %s: %w", action, err)
}

func nullableText(value string) pgtype.Text {
	if strings.TrimSpace(value) == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: strings.TrimSpace(value), Valid: true}
}

func nullableTextValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return strings.TrimSpace(value.String)
}

const runtimeSessionColumns = `
session_id,
actor_kind,
actor_id,
player_id,
session_status,
issued_at,
expires_at,
last_seen_at,
revoked_at,
revocation_reason,
access_token_record_id,
created_at,
updated_at`

const insertRuntimeSessionSQL = `
INSERT INTO runtime_sessions (
    session_id,
    actor_kind,
    actor_id,
    player_id,
    session_status,
    issued_at,
    expires_at,
    last_seen_at,
    access_token_record_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $6, $6)
RETURNING ` + runtimeSessionColumns

const getRuntimeSessionSQL = `
SELECT ` + runtimeSessionColumns + `
FROM runtime_sessions
WHERE session_id = $1`

const findActiveRuntimeSessionSQL = `
SELECT ` + runtimeSessionColumns + `
FROM runtime_sessions
WHERE session_id = $1
  AND expires_at > $2
  AND session_status = $3`

const updateRuntimeSessionLastSeenSQL = `
UPDATE runtime_sessions
SET last_seen_at = $2,
    updated_at = $2
WHERE session_id = $1
RETURNING ` + runtimeSessionColumns

const markRuntimeSessionExpiredSQL = `
UPDATE runtime_sessions
SET session_status = $3,
    updated_at = $2
WHERE session_id = $1
  AND session_status = $4
RETURNING ` + runtimeSessionColumns

const revokeRuntimeSessionSQL = `
UPDATE runtime_sessions
SET session_status = $4,
    revoked_at = $2,
    revocation_reason = $3,
    updated_at = $2
WHERE session_id = $1
RETURNING ` + runtimeSessionColumns

const listActiveRuntimeSessionsForPlayerSQL = `
SELECT ` + runtimeSessionColumns + `
FROM runtime_sessions
WHERE player_id = $1
  AND expires_at > $2
  AND session_status = $3
ORDER BY last_seen_at DESC, session_id
LIMIT $4`
