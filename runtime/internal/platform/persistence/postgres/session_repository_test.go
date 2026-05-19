package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app/session"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestSessionRepositoryCreateRuntimeSessionInsertsLifecycleRecord(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 10, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{values: activeSessionRowValues(issuedAt.UTC(), issuedAt.Add(time.Hour).UTC())},
		},
	}
	repository := NewSessionRepositoryForUnitOfWork(executor)

	record, err := repository.CreateRuntimeSession(context.Background(), session.CreateRuntimeSessionMutation{
		SessionID:           " session-1 ",
		ActorKind:           session.ActorKind(" player "),
		ActorID:             " player-1 ",
		PlayerID:            " player-1 ",
		IssuedAt:            issuedAt,
		ExpiresAt:           issuedAt.Add(time.Hour),
		LastSeenAt:          issuedAt.Add(time.Minute),
		AccessTokenRecordID: " token-1 ",
		RequestedBy:         " authentication_service ",
	})
	if err != nil {
		t.Fatalf("CreateRuntimeSession() error = %v, want nil", err)
	}
	if record.SessionID != "session-1" ||
		record.ActorKind != session.ActorKindPlayer ||
		record.PlayerID != "player-1" ||
		record.SessionStatus != session.SessionStatusActive ||
		record.AccessTokenRecordID != "token-1" {
		t.Fatalf("CreateRuntimeSession() record = %#v, want active runtime session", record)
	}
	if record.IssuedAt.Location() != time.UTC || record.ExpiresAt.Location() != time.UTC {
		t.Fatalf("CreateRuntimeSession() timestamps = %#v, want UTC", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "INSERT INTO runtime_sessions")
	assertSQLContains(t, call.sql, "RETURNING")
	assertArgs(t,
		call.args,
		"session-1",
		"player",
		"player-1",
		"player-1",
		string(session.SessionStatusActive),
		issuedAt.UTC(),
		issuedAt.Add(time.Hour).UTC(),
		issuedAt.Add(time.Minute).UTC(),
		pgtype.Text{String: "token-1", Valid: true},
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestSessionRepositoryGetRuntimeSessionSelectsBySessionID(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 1, 2, 3, 0, time.UTC)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{values: activeSessionRowValues(issuedAt, issuedAt.Add(time.Hour))},
		},
	}
	repository := NewSessionRepositoryForUnitOfWork(executor)

	record, err := repository.GetRuntimeSession(context.Background(), session.GetRuntimeSessionQuery{SessionID: " session-1 "})
	if err != nil {
		t.Fatalf("GetRuntimeSession() error = %v, want nil", err)
	}
	if record.SessionID != "session-1" {
		t.Fatalf("GetRuntimeSession() record = %#v, want session-1", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "FROM runtime_sessions")
	assertSQLContains(t, call.sql, "WHERE session_id = $1")
	assertArgs(t, call.args, "session-1")
}

func TestSessionRepositoryFindActiveSessionByIDFiltersStatusAndExpiration(t *testing.T) {
	observedAt := time.Date(2026, 5, 17, 2, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{values: activeSessionRowValues(observedAt.Add(-time.Minute).UTC(), observedAt.Add(time.Hour).UTC())},
		},
	}
	repository := NewSessionRepositoryForUnitOfWork(executor)

	if _, err := repository.FindActiveSessionByID(context.Background(), session.FindActiveSessionByIDQuery{
		SessionID:  " session-1 ",
		ObservedAt: observedAt,
	}); err != nil {
		t.Fatalf("FindActiveSessionByID() error = %v, want nil", err)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "expires_at > $2")
	assertSQLContains(t, call.sql, "session_status = $3")
	assertArgs(t, call.args, "session-1", observedAt.UTC(), string(session.SessionStatusActive))
}

func TestSessionRepositoryUpdateRuntimeSessionLastSeenReturnsUpdatedRow(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 1, 0, 0, 0, time.UTC)
	lastSeenAt := issuedAt.Add(15 * time.Minute)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{values: activeSessionRowValues(issuedAt, issuedAt.Add(time.Hour), withSessionLastSeen(lastSeenAt))},
		},
	}
	repository := NewSessionRepositoryForUnitOfWork(executor)

	record, err := repository.UpdateRuntimeSessionLastSeen(context.Background(), session.UpdateRuntimeSessionLastSeenMutation{
		SessionID:   " session-1 ",
		LastSeenAt:  lastSeenAt,
		RequestedBy: " session_validation ",
	})
	if err != nil {
		t.Fatalf("UpdateRuntimeSessionLastSeen() error = %v, want nil", err)
	}
	if !record.LastSeenAt.Equal(lastSeenAt.UTC()) {
		t.Fatalf("LastSeenAt = %s, want %s", record.LastSeenAt, lastSeenAt.UTC())
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "UPDATE runtime_sessions")
	assertSQLContains(t, call.sql, "last_seen_at = $2")
	assertSQLContains(t, call.sql, "updated_at = $2")
	assertArgs(t, call.args, "session-1", lastSeenAt.UTC())
}

func TestSessionRepositoryMarkRuntimeSessionExpiredTransitionsActiveRow(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 1, 0, 0, 0, time.UTC)
	expiredAt := issuedAt.Add(time.Hour)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{values: activeSessionRowValues(issuedAt, expiredAt.Add(time.Hour), withSessionStatus(session.SessionStatusExpired), withSessionUpdatedAt(expiredAt))},
		},
	}
	repository := NewSessionRepositoryForUnitOfWork(executor)

	record, err := repository.MarkRuntimeSessionExpired(context.Background(), session.MarkRuntimeSessionExpiredMutation{
		SessionID:   " session-1 ",
		ExpiredAt:   expiredAt,
		RequestedBy: " cleanup ",
	})
	if err != nil {
		t.Fatalf("MarkRuntimeSessionExpired() error = %v, want nil", err)
	}
	if record.SessionStatus != session.SessionStatusExpired {
		t.Fatalf("SessionStatus = %q, want expired", record.SessionStatus)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "session_status = $3")
	assertSQLContains(t, call.sql, "AND session_status = $4")
	assertArgs(t, call.args, "session-1", expiredAt.UTC(), string(session.SessionStatusExpired), string(session.SessionStatusActive))
}

func TestSessionRepositoryRevokeRuntimeSessionSetsRevokedFields(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 1, 0, 0, 0, time.UTC)
	revokedAt := issuedAt.Add(30 * time.Minute)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{values: activeSessionRowValues(
				issuedAt,
				issuedAt.Add(time.Hour),
				withSessionStatus(session.SessionStatusRevoked),
				withSessionRevokedAt(revokedAt),
				withSessionRevocationReason("logout"),
				withSessionUpdatedAt(revokedAt),
			)},
		},
	}
	repository := NewSessionRepositoryForUnitOfWork(executor)

	record, err := repository.RevokeRuntimeSession(context.Background(), session.RevokeRuntimeSessionMutation{
		SessionID:        " session-1 ",
		RevokedAt:        revokedAt,
		RevocationReason: " logout ",
		RequestedBy:      " authentication_service ",
	})
	if err != nil {
		t.Fatalf("RevokeRuntimeSession() error = %v, want nil", err)
	}
	if record.RevokedAt == nil || !record.RevokedAt.Equal(revokedAt.UTC()) || record.RevocationReason != "logout" {
		t.Fatalf("RevokeRuntimeSession() record = %#v, want revoked logout session", record)
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "revoked_at = $2")
	assertSQLContains(t, call.sql, "revocation_reason = $3")
	assertSQLContains(t, call.sql, "session_status = $4")
	assertArgs(t, call.args, "session-1", revokedAt.UTC(), "logout", string(session.SessionStatusRevoked))
}

func TestSessionRepositoryListActiveSessionsForPlayerUsesBoundedOrderedQuery(t *testing.T) {
	observedAt := time.Date(2026, 5, 17, 2, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{
		rowsResponses: []pgx.Rows{
			&sessionRows{
				values: [][]any{
					activeSessionRowValues(observedAt.Add(-time.Minute).UTC(), observedAt.Add(time.Hour).UTC()),
					activeSessionRowValues(observedAt.Add(-2*time.Minute).UTC(), observedAt.Add(2*time.Hour).UTC(), withSessionID("session-2")),
				},
			},
		},
	}
	repository := NewSessionRepositoryForUnitOfWork(executor)

	records, err := repository.ListActiveSessionsForPlayer(context.Background(), session.ListActiveSessionsForPlayerQuery{
		PlayerID:   " player-1 ",
		ObservedAt: observedAt,
		Limit:      2,
	})
	if err != nil {
		t.Fatalf("ListActiveSessionsForPlayer() error = %v, want nil", err)
	}
	if len(records) != 2 || records[0].SessionID != "session-1" || records[1].SessionID != "session-2" {
		t.Fatalf("ListActiveSessionsForPlayer() records = %#v, want two active sessions", records)
	}

	if len(executor.queries) != 1 {
		t.Fatalf("queries len = %d, want 1", len(executor.queries))
	}
	call := executor.queries[0]
	assertSQLContains(t, call.sql, "FROM runtime_sessions")
	assertSQLContains(t, call.sql, "WHERE player_id = $1")
	assertSQLContains(t, call.sql, "expires_at > $2")
	assertSQLContains(t, call.sql, "ORDER BY last_seen_at DESC, session_id")
	assertSQLContains(t, call.sql, "LIMIT $4")
	assertArgs(t, call.args, "player-1", observedAt.UTC(), string(session.SessionStatusActive), int32(2))
}

func TestSessionRepositoryMapsErrors(t *testing.T) {
	repository := NewSessionRepositoryForUnitOfWork(&recordingExecutor{})

	_, err := repository.GetRuntimeSession(context.Background(), session.GetRuntimeSessionQuery{SessionID: "session-1"})
	if !errors.Is(err, ErrSessionRecordNotFound) {
		t.Fatalf("GetRuntimeSession() error = %v, want ErrSessionRecordNotFound", err)
	}

	duplicate := &pgconn.PgError{Code: "23505", ConstraintName: "runtime_sessions_pkey"}
	repository = NewSessionRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{err: duplicate},
		},
	})
	_, err = repository.CreateRuntimeSession(context.Background(), validCreateRuntimeSessionMutation())
	if !errors.Is(err, ErrSessionRecordConflict) {
		t.Fatalf("CreateRuntimeSession() error = %v, want ErrSessionRecordConflict", err)
	}
	assertDoesNotLeakPgError(t, err)

	constraint := &pgconn.PgError{Code: "23503", ConstraintName: "runtime_sessions_player_fk"}
	repository = NewSessionRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{err: constraint},
		},
	})
	_, err = repository.CreateRuntimeSession(context.Background(), validCreateRuntimeSessionMutation())
	if !errors.Is(err, ErrSessionRecordConstraint) {
		t.Fatalf("CreateRuntimeSession() error = %v, want ErrSessionRecordConstraint", err)
	}
	assertDoesNotLeakPgError(t, err)

	repository = NewSessionRepositoryForUnitOfWork(&recordingExecutor{})
	_, err = repository.MarkRuntimeSessionExpired(context.Background(), session.MarkRuntimeSessionExpiredMutation{
		SessionID:   "session-1",
		ExpiredAt:   time.Date(2026, 5, 17, 1, 0, 0, 0, time.UTC),
		RequestedBy: "cleanup",
	})
	if !errors.Is(err, ErrSessionRecordStale) {
		t.Fatalf("MarkRuntimeSessionExpired() error = %v, want ErrSessionRecordStale", err)
	}
}

func TestSessionRepositoryRejectsInvalidRows(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 1, 0, 0, 0, time.UTC)
	values := activeSessionRowValues(issuedAt, issuedAt.Add(time.Hour))
	values[4] = "disabled"
	repository := NewSessionRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{values: values},
		},
	})

	_, err := repository.GetRuntimeSession(context.Background(), session.GetRuntimeSessionQuery{SessionID: "session-1"})
	if !errors.Is(err, ErrSessionRecordConstraint) {
		t.Fatalf("GetRuntimeSession() error = %v, want ErrSessionRecordConstraint", err)
	}
}

func TestSessionRepositoryRequiresUnitOfWorkExecutor(t *testing.T) {
	repository := NewSessionRepositoryForUnitOfWork(nil)

	_, err := repository.CreateRuntimeSession(context.Background(), validCreateRuntimeSessionMutation())
	if err == nil {
		t.Fatal("CreateRuntimeSession() error = nil, want executor error")
	}
	_, err = repository.ListActiveSessionsForPlayer(context.Background(), session.ListActiveSessionsForPlayerQuery{
		PlayerID:   "player-1",
		ObservedAt: time.Date(2026, 5, 17, 1, 0, 0, 0, time.UTC),
		Limit:      1,
	})
	if err == nil {
		t.Fatal("ListActiveSessionsForPlayer() error = nil, want executor error")
	}
}

func TestSessionRepositoryDefaultTestsDoNotRequireLivePostgreSQL(t *testing.T) {
	if os.Getenv("VIBIT_POSTGRES_TEST_DSN") != "" {
		t.Skip("live PostgreSQL environment is opt-in and not needed for this fake-executor test")
	}

	issuedAt := time.Date(2026, 5, 17, 1, 0, 0, 0, time.UTC)
	repository := NewSessionRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			sessionRow{values: activeSessionRowValues(issuedAt, issuedAt.Add(time.Hour))},
		},
	})

	if _, err := repository.GetRuntimeSession(context.Background(), session.GetRuntimeSessionQuery{SessionID: "session-1"}); err != nil {
		t.Fatalf("GetRuntimeSession() error = %v, want nil without live PostgreSQL", err)
	}
}

func TestPostgresUnitOfWorkCreatesSessionRepository(t *testing.T) {
	executor := &recordingExecutor{}
	unit := UnitOfWork{executor: executor}

	repository, err := unit.NewSessionRepository()
	if err != nil {
		t.Fatalf("NewSessionRepository() error = %v, want nil", err)
	}
	if repository == nil {
		t.Fatal("NewSessionRepository() = nil, want repository")
	}
}

func validCreateRuntimeSessionMutation() session.CreateRuntimeSessionMutation {
	issuedAt := time.Date(2026, 5, 17, 1, 0, 0, 0, time.UTC)
	return session.CreateRuntimeSessionMutation{
		SessionID:           "session-1",
		ActorKind:           session.ActorKindPlayer,
		ActorID:             "player-1",
		PlayerID:            "player-1",
		IssuedAt:            issuedAt,
		ExpiresAt:           issuedAt.Add(time.Hour),
		LastSeenAt:          issuedAt,
		AccessTokenRecordID: "token-1",
		RequestedBy:         "authentication_service",
	}
}

type sessionRowOption func([]any)

func withSessionID(sessionID string) sessionRowOption {
	return func(values []any) {
		values[0] = sessionID
	}
}

func withSessionStatus(status session.SessionStatus) sessionRowOption {
	return func(values []any) {
		values[4] = string(status)
	}
}

func withSessionLastSeen(lastSeenAt time.Time) sessionRowOption {
	return func(values []any) {
		values[7] = lastSeenAt
	}
}

func withSessionRevokedAt(revokedAt time.Time) sessionRowOption {
	return func(values []any) {
		values[8] = revokedAt
	}
}

func withSessionRevocationReason(reason string) sessionRowOption {
	return func(values []any) {
		values[9] = reason
	}
}

func withSessionUpdatedAt(updatedAt time.Time) sessionRowOption {
	return func(values []any) {
		values[12] = updatedAt
	}
}

func activeSessionRowValues(issuedAt time.Time, expiresAt time.Time, options ...sessionRowOption) []any {
	values := []any{
		"session-1",
		"player",
		"player-1",
		"player-1",
		"active",
		issuedAt,
		expiresAt,
		issuedAt,
		nil,
		nil,
		"token-1",
		issuedAt,
		issuedAt,
	}
	for _, option := range options {
		option(values)
	}
	return values
}

type sessionRow struct {
	values []any
	err    error
}

func (r sessionRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	assignSessionValues("session row", dest, r.values)
	return nil
}

type sessionRows struct {
	values [][]any
	index  int
	err    error
	closed bool
}

func (r *sessionRows) Close() {
	r.closed = true
}

func (r *sessionRows) Err() error {
	return r.err
}

func (r *sessionRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}

func (r *sessionRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *sessionRows) Next() bool {
	if r.index >= len(r.values) {
		r.closed = true
		return false
	}
	r.index += 1
	return true
}

func (r *sessionRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.values) {
		return errors.New("session rows: scan without current row")
	}
	assignSessionValues("session rows", dest, r.values[r.index-1])
	return nil
}

func (r *sessionRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.values) {
		return nil, errors.New("session rows: values without current row")
	}
	return append([]any(nil), r.values[r.index-1]...), nil
}

func (r *sessionRows) RawValues() [][]byte {
	return nil
}

func (r *sessionRows) Conn() *pgx.Conn {
	return nil
}

func assignSessionValues(label string, dest []any, values []any) {
	if len(dest) != len(values) {
		panic(label + ": destination count mismatch")
	}
	for i := range dest {
		assignSessionValue(label, dest[i], values[i])
	}
}

func assignSessionValue(label string, dest any, value any) {
	switch pointer := dest.(type) {
	case *string:
		switch typed := value.(type) {
		case string:
			*pointer = typed
		default:
			panic(label + ": unsupported string value")
		}
	case *time.Time:
		*pointer = value.(time.Time)
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
	default:
		panic(label + ": unsupported destination type")
	}
}
