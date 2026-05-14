package postgres

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/modules/authentication"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestAuthenticationRepositoryStoreCredentialInsertsRecord(t *testing.T) {
	occurredAt := time.Date(2026, 5, 14, 9, 30, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	createdAt := occurredAt.Add(time.Second)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			authenticationCredentialRow{values: activeCredentialRowValues(createdAt)},
		},
	}
	repository := NewAuthenticationRepositoryForUnitOfWork(executor)

	record, err := repository.StoreCredential(context.Background(), authentication.StoreCredentialMutation{
		CredentialRecordID:       " credential-1 ",
		PlayerID:                 " player-1 ",
		CredentialKind:           " device_credential_login ",
		CredentialLookupDigest:   []byte{1, 2},
		CredentialVerifierDigest: []byte{3, 4},
		VerifierAlgorithm:        " hmac-sha256 ",
		VerifierVersion:          1,
		VerifierKeyID:            " key-1 ",
		ClientInstanceIDDigest:   []byte{5, 6},
		OccurredAt:               occurredAt,
		RequestedBy:              " maintainer ",
	})
	if err != nil {
		t.Fatalf("StoreCredential() error = %v, want nil", err)
	}

	if record.CredentialRecordID != "credential-1" ||
		record.PlayerID != "player-1" ||
		record.CredentialStatus != authentication.CredentialStatusActive {
		t.Fatalf("StoreCredential() record = %#v, want active credential", record)
	}
	if !record.CreatedAt.Equal(createdAt.UTC()) || record.CreatedAt.Location() != time.UTC {
		t.Fatalf("CreatedAt = %s, want UTC %s", record.CreatedAt, createdAt.UTC())
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "INSERT INTO authentication_device_credentials")
	assertSQLContains(t, call.sql, "RETURNING")
	assertArgs(t,
		call.args,
		"credential-1",
		"player-1",
		"device_credential_login",
		string(authentication.CredentialStatusActive),
		[]byte{1, 2},
		[]byte{3, 4},
		"hmac-sha256",
		int32(1),
		"key-1",
		[]byte{5, 6},
		occurredAt.UTC(),
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestAuthenticationRepositoryFindCredentialByLookupDigestSelectsRecord(t *testing.T) {
	timestamp := time.Date(2026, 5, 14, 1, 2, 3, 0, time.FixedZone("UTC+8", 8*60*60))
	lastVerifiedAt := timestamp.Add(time.Minute)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			authenticationCredentialRow{values: []any{
				"credential-1",
				"player-1",
				"device_credential_login",
				"active",
				[]byte{1, 2},
				[]byte{3, 4},
				"hmac-sha256",
				int32(1),
				"key-1",
				[]byte{5, 6},
				timestamp,
				timestamp,
				lastVerifiedAt,
				nil,
				"",
				nil,
				"",
				nil,
				"",
			}},
		},
	}
	repository := NewAuthenticationRepositoryForUnitOfWork(executor)

	record, err := repository.FindCredentialByLookupDigest(context.Background(), []byte{1, 2})
	if err != nil {
		t.Fatalf("FindCredentialByLookupDigest() error = %v, want nil", err)
	}

	if record.LastVerifiedAt == nil || !record.LastVerifiedAt.Equal(lastVerifiedAt.UTC()) || record.LastVerifiedAt.Location() != time.UTC {
		t.Fatalf("LastVerifiedAt = %#v, want UTC %s", record.LastVerifiedAt, lastVerifiedAt.UTC())
	}
	if !reflect.DeepEqual(record.ClientInstanceIDDigest, []byte{5, 6}) {
		t.Fatalf("ClientInstanceIDDigest = %#v, want [5 6]", record.ClientInstanceIDDigest)
	}
	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "FROM authentication_device_credentials")
	assertSQLContains(t, call.sql, "WHERE credential_lookup_digest = $1")
	assertArgs(t, call.args, []byte{1, 2})
}

func TestAuthenticationRepositoryStoreTokenInsertsRecord(t *testing.T) {
	issuedAt := time.Date(2026, 5, 14, 10, 30, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	expiresAt := issuedAt.Add(time.Hour)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			authenticationTokenRow{values: activeTokenRowValues(issuedAt.UTC(), expiresAt.UTC())},
		},
	}
	repository := NewAuthenticationRepositoryForUnitOfWork(executor)

	record, err := repository.StoreToken(context.Background(), authentication.StoreTokenMutation{
		TokenRecordID:       " token-1 ",
		PlayerID:            " player-1 ",
		CredentialRecordID:  " credential-1 ",
		TokenKind:           " access_token ",
		ActorKind:           " player ",
		TokenLookupDigest:   []byte{7, 8},
		TokenVerifierDigest: []byte{9, 10},
		VerifierAlgorithm:   " hmac-sha256 ",
		VerifierVersion:     1,
		VerifierKeyID:       " key-1 ",
		Audience:            " gameplay ",
		IssuedAt:            issuedAt,
		ExpiresAt:           expiresAt,
		RequestedBy:         " maintainer ",
	})
	if err != nil {
		t.Fatalf("StoreToken() error = %v, want nil", err)
	}

	if record.TokenRecordID != "token-1" ||
		record.TokenKind != "access_token" ||
		record.TokenStatus != authentication.TokenStatusActive ||
		record.PlayerID != "player-1" {
		t.Fatalf("StoreToken() record = %#v, want active token", record)
	}
	if !record.IssuedAt.Equal(issuedAt.UTC()) || record.IssuedAt.Location() != time.UTC {
		t.Fatalf("IssuedAt = %s, want UTC %s", record.IssuedAt, issuedAt.UTC())
	}

	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "INSERT INTO authentication_access_tokens")
	assertSQLContains(t, call.sql, "RETURNING")
	assertArgs(t,
		call.args,
		"token-1",
		"access_token",
		string(authentication.TokenStatusActive),
		"player",
		"player-1",
		"credential-1",
		[]byte{7, 8},
		[]byte{9, 10},
		"hmac-sha256",
		int32(1),
		"key-1",
		"gameplay",
		issuedAt.UTC(),
		expiresAt.UTC(),
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestAuthenticationRepositoryFindTokenByLookupDigestSelectsRecord(t *testing.T) {
	issuedAt := time.Date(2026, 5, 14, 1, 2, 3, 0, time.FixedZone("UTC+8", 8*60*60))
	expiresAt := issuedAt.Add(time.Hour)
	lastValidatedAt := issuedAt.Add(10 * time.Minute)
	executor := &recordingExecutor{
		rowResponses: []pgx.Row{
			authenticationTokenRow{values: []any{
				"token-1",
				"access_token",
				"active",
				"player",
				"player-1",
				"credential-1",
				[]byte{7, 8},
				[]byte{9, 10},
				"hmac-sha256",
				int32(1),
				"key-1",
				"gameplay",
				issuedAt,
				expiresAt,
				nil,
				"",
				"",
				lastValidatedAt,
				nil,
				nil,
				issuedAt,
				issuedAt,
			}},
		},
	}
	repository := NewAuthenticationRepositoryForUnitOfWork(executor)

	record, err := repository.FindTokenByLookupDigest(context.Background(), []byte{7, 8})
	if err != nil {
		t.Fatalf("FindTokenByLookupDigest() error = %v, want nil", err)
	}

	if record.LastValidatedAt == nil || !record.LastValidatedAt.Equal(lastValidatedAt.UTC()) || record.LastValidatedAt.Location() != time.UTC {
		t.Fatalf("LastValidatedAt = %#v, want UTC %s", record.LastValidatedAt, lastValidatedAt.UTC())
	}
	if !reflect.DeepEqual(record.TokenLookupDigest, []byte{7, 8}) {
		t.Fatalf("TokenLookupDigest = %#v, want [7 8]", record.TokenLookupDigest)
	}
	if len(executor.queryRowCalls) != 1 {
		t.Fatalf("query rows len = %d, want 1", len(executor.queryRowCalls))
	}
	call := executor.queryRowCalls[0]
	assertSQLContains(t, call.sql, "FROM authentication_access_tokens")
	assertSQLContains(t, call.sql, "WHERE token_lookup_digest = $1")
	assertArgs(t, call.args, []byte{7, 8})
}

func TestAuthenticationRepositoryRevokeCredentialUpdatesTerminalState(t *testing.T) {
	revokedAt := time.Date(2026, 5, 14, 2, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	executor := &recordingExecutor{}
	repository := NewAuthenticationRepositoryForUnitOfWork(executor)

	err := repository.RevokeCredential(context.Background(), authentication.RevokeCredentialMutation{
		CredentialRecordID: " credential-1 ",
		RevokedAt:          revokedAt,
		RevokedReason:      " rotated ",
		RequestedBy:        " maintainer ",
	})
	if err != nil {
		t.Fatalf("RevokeCredential() error = %v, want nil", err)
	}

	if len(executor.execs) != 1 {
		t.Fatalf("execs len = %d, want 1", len(executor.execs))
	}
	call := executor.execs[0]
	assertSQLContains(t, call.sql, "UPDATE authentication_device_credentials")
	assertSQLContains(t, call.sql, "credential_status = $1")
	assertSQLContains(t, call.sql, "WHERE credential_record_id = $4")
	assertArgs(t, call.args, string(authentication.CredentialStatusRevoked), revokedAt.UTC(), "rotated", "credential-1")
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestAuthenticationRepositoryRevokeTokenUpdatesTerminalStateAndCleanup(t *testing.T) {
	revokedAt := time.Date(2026, 5, 14, 2, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	cleanupAfter := revokedAt.Add(7 * 24 * time.Hour)
	executor := &recordingExecutor{}
	repository := NewAuthenticationRepositoryForUnitOfWork(executor)

	err := repository.RevokeToken(context.Background(), authentication.RevokeTokenMutation{
		TokenRecordID: " token-1 ",
		RevokedAt:     revokedAt,
		RevokedReason: " logout ",
		CleanupAfter:  &cleanupAfter,
		RequestedBy:   " maintainer ",
	})
	if err != nil {
		t.Fatalf("RevokeToken() error = %v, want nil", err)
	}

	if len(executor.execs) != 1 {
		t.Fatalf("execs len = %d, want 1", len(executor.execs))
	}
	call := executor.execs[0]
	assertSQLContains(t, call.sql, "UPDATE authentication_access_tokens")
	assertSQLContains(t, call.sql, "token_status = $1")
	assertSQLContains(t, call.sql, "cleanup_after = $4")
	assertSQLContains(t, call.sql, "WHERE token_record_id = $5")
	assertArgs(t,
		call.args,
		string(authentication.TokenStatusRevoked),
		revokedAt.UTC(),
		"logout",
		pgtype.Timestamptz{Time: cleanupAfter.UTC(), Valid: true},
		"token-1",
	)
	if hasTransactionControlSQL(executor.allSQL()) {
		t.Fatalf("repository SQL included transaction control: %#v", executor.allSQL())
	}
}

func TestAuthenticationRepositoryListTokensEligibleForCleanup(t *testing.T) {
	now := time.Date(2026, 5, 14, 10, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	issuedAt := now.Add(-2 * time.Hour)
	cleanupAfter := now.Add(-time.Minute)
	executor := &recordingExecutor{
		rowsResponses: []pgx.Rows{
			&authenticationTokenRows{
				values: [][]any{
					revokedTokenRowValues(issuedAt, issuedAt.Add(time.Hour), cleanupAfter),
				},
			},
		},
	}
	repository := NewAuthenticationRepositoryForUnitOfWork(executor)

	records, err := repository.ListTokensEligibleForCleanup(context.Background(), authentication.TokenCleanupQuery{
		Now:   now,
		Limit: 25,
	})
	if err != nil {
		t.Fatalf("ListTokensEligibleForCleanup() error = %v, want nil", err)
	}
	if len(records) != 1 || records[0].TokenRecordID != "token-1" || records[0].CleanupAfter == nil {
		t.Fatalf("ListTokensEligibleForCleanup() records = %#v, want one cleanup-eligible token", records)
	}
	if len(executor.queries) != 1 {
		t.Fatalf("queries len = %d, want 1", len(executor.queries))
	}
	call := executor.queries[0]
	assertSQLContains(t, call.sql, "FROM authentication_access_tokens")
	assertSQLContains(t, call.sql, "cleanup_after <= $1")
	assertSQLContains(t, call.sql, "ORDER BY cleanup_after, token_record_id")
	assertArgs(t, call.args, now.UTC(), int32(25))
}

func TestAuthenticationRepositoryMapsErrors(t *testing.T) {
	repository := NewAuthenticationRepositoryForUnitOfWork(&recordingExecutor{})

	_, err := repository.FindCredentialByLookupDigest(context.Background(), []byte{1})
	if !errors.Is(err, ErrAuthenticationRecordNotFound) {
		t.Fatalf("FindCredentialByLookupDigest() error = %v, want ErrAuthenticationRecordNotFound", err)
	}

	duplicate := &pgconn.PgError{Code: "23505", ConstraintName: "authentication_device_credentials_lookup_digest_uq"}
	repository = NewAuthenticationRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			authenticationCredentialRow{err: duplicate},
		},
	})
	_, err = repository.StoreCredential(context.Background(), validStoreCredentialMutation())
	if !errors.Is(err, ErrAuthenticationRecordConflict) {
		t.Fatalf("StoreCredential() error = %v, want ErrAuthenticationRecordConflict", err)
	}
	assertDoesNotLeakPgError(t, err)

	constraint := &pgconn.PgError{Code: "23503", ConstraintName: "authentication_access_tokens_player_fk"}
	repository = NewAuthenticationRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			authenticationTokenRow{err: constraint},
		},
	})
	_, err = repository.StoreToken(context.Background(), validStoreTokenMutation())
	if !errors.Is(err, ErrAuthenticationRecordConstraint) {
		t.Fatalf("StoreToken() error = %v, want ErrAuthenticationRecordConstraint", err)
	}
	assertDoesNotLeakPgError(t, err)

	repository = NewAuthenticationRepositoryForUnitOfWork(&recordingExecutor{execCommandTag: "UPDATE 0"})
	err = repository.RevokeToken(context.Background(), validRevokeTokenMutation())
	if !errors.Is(err, ErrAuthenticationRecordNotFound) {
		t.Fatalf("RevokeToken() error = %v, want ErrAuthenticationRecordNotFound", err)
	}
}

func TestAuthenticationRepositoryRejectsInvalidRows(t *testing.T) {
	values := activeTokenRowValues(time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC), time.Date(2026, 5, 14, 2, 2, 3, 0, time.UTC))
	values[2] = "pending"
	repository := NewAuthenticationRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			authenticationTokenRow{values: values},
		},
	})

	_, err := repository.FindTokenByLookupDigest(context.Background(), []byte{7, 8})
	if !errors.Is(err, ErrAuthenticationRecordConstraint) {
		t.Fatalf("FindTokenByLookupDigest() error = %v, want ErrAuthenticationRecordConstraint", err)
	}
}

func TestAuthenticationRepositoryRequiresUnitOfWorkExecutor(t *testing.T) {
	repository := NewAuthenticationRepositoryForUnitOfWork(nil)

	_, err := repository.StoreCredential(context.Background(), validStoreCredentialMutation())
	if err == nil {
		t.Fatal("StoreCredential() error = nil, want executor error")
	}
	_, err = repository.StoreToken(context.Background(), validStoreTokenMutation())
	if err == nil {
		t.Fatal("StoreToken() error = nil, want executor error")
	}
}

func TestAuthenticationRepositoryDefaultTestsDoNotRequireLivePostgreSQL(t *testing.T) {
	if os.Getenv("VIBIT_POSTGRES_TEST_DSN") != "" {
		t.Skip("live PostgreSQL environment is opt-in and not needed for this fake-executor test")
	}

	repository := NewAuthenticationRepositoryForUnitOfWork(&recordingExecutor{
		rowResponses: []pgx.Row{
			authenticationCredentialRow{values: activeCredentialRowValues(time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC))},
		},
	})

	if _, err := repository.FindCredentialByLookupDigest(context.Background(), []byte{1, 2}); err != nil {
		t.Fatalf("FindCredentialByLookupDigest() error = %v, want nil without live PostgreSQL", err)
	}
}

func TestPostgresUnitOfWorkCreatesAuthenticationRepository(t *testing.T) {
	executor := &recordingExecutor{}
	unit := UnitOfWork{executor: executor}

	repository, err := unit.NewAuthenticationRepository()
	if err != nil {
		t.Fatalf("NewAuthenticationRepository() error = %v, want nil", err)
	}
	if repository == nil {
		t.Fatal("NewAuthenticationRepository() = nil, want repository")
	}
}

func activeCredentialRowValues(timestamp time.Time) []any {
	return []any{
		"credential-1",
		"player-1",
		"device_credential_login",
		"active",
		[]byte{1, 2},
		[]byte{3, 4},
		"hmac-sha256",
		int32(1),
		"key-1",
		[]byte{5, 6},
		timestamp,
		timestamp,
		nil,
		nil,
		"",
		nil,
		"",
		nil,
		"",
	}
}

func activeTokenRowValues(issuedAt time.Time, expiresAt time.Time) []any {
	return []any{
		"token-1",
		"access_token",
		"active",
		"player",
		"player-1",
		"credential-1",
		[]byte{7, 8},
		[]byte{9, 10},
		"hmac-sha256",
		int32(1),
		"key-1",
		"gameplay",
		issuedAt,
		expiresAt,
		nil,
		"",
		"",
		nil,
		nil,
		nil,
		issuedAt,
		issuedAt,
	}
}

func revokedTokenRowValues(issuedAt time.Time, expiresAt time.Time, cleanupAfter time.Time) []any {
	values := activeTokenRowValues(issuedAt, expiresAt)
	values[2] = "revoked"
	values[14] = expiresAt.Add(-time.Minute)
	values[15] = "logout"
	values[19] = cleanupAfter
	return values
}

func validStoreCredentialMutation() authentication.StoreCredentialMutation {
	return authentication.StoreCredentialMutation{
		CredentialRecordID:       "credential-1",
		PlayerID:                 "player-1",
		CredentialKind:           "device_credential_login",
		CredentialLookupDigest:   []byte{1},
		CredentialVerifierDigest: []byte{2},
		VerifierAlgorithm:        "hmac-sha256",
		VerifierVersion:          1,
		VerifierKeyID:            "key-1",
		OccurredAt:               time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC),
		RequestedBy:              "maintainer",
	}
}

func validStoreTokenMutation() authentication.StoreTokenMutation {
	issuedAt := time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC)
	return authentication.StoreTokenMutation{
		TokenRecordID:       "token-1",
		PlayerID:            "player-1",
		CredentialRecordID:  "credential-1",
		TokenKind:           "access_token",
		ActorKind:           "player",
		TokenLookupDigest:   []byte{1},
		TokenVerifierDigest: []byte{2},
		VerifierAlgorithm:   "hmac-sha256",
		VerifierVersion:     1,
		VerifierKeyID:       "key-1",
		Audience:            "gameplay",
		IssuedAt:            issuedAt,
		ExpiresAt:           issuedAt.Add(time.Hour),
		RequestedBy:         "maintainer",
	}
}

func validRevokeTokenMutation() authentication.RevokeTokenMutation {
	return authentication.RevokeTokenMutation{
		TokenRecordID: "token-1",
		RevokedAt:     time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC),
		RevokedReason: "logout",
		RequestedBy:   "maintainer",
	}
}

type authenticationCredentialRow struct {
	values []any
	err    error
}

func (r authenticationCredentialRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	assignAuthenticationValues("authentication credential row", dest, r.values)
	return nil
}

type authenticationTokenRow struct {
	values []any
	err    error
}

func (r authenticationTokenRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	assignAuthenticationValues("authentication token row", dest, r.values)
	return nil
}

type authenticationTokenRows struct {
	values [][]any
	index  int
	err    error
	closed bool
}

func (r *authenticationTokenRows) Close() {
	r.closed = true
}

func (r *authenticationTokenRows) Err() error {
	return r.err
}

func (r *authenticationTokenRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}

func (r *authenticationTokenRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *authenticationTokenRows) Next() bool {
	if r.index >= len(r.values) {
		r.closed = true
		return false
	}
	r.index += 1
	return true
}

func (r *authenticationTokenRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.values) {
		return errors.New("authentication token rows: scan without current row")
	}
	assignAuthenticationValues("authentication token rows", dest, r.values[r.index-1])
	return nil
}

func (r *authenticationTokenRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.values) {
		return nil, errors.New("authentication token rows: values without current row")
	}
	return append([]any(nil), r.values[r.index-1]...), nil
}

func (r *authenticationTokenRows) RawValues() [][]byte {
	return nil
}

func (r *authenticationTokenRows) Conn() *pgx.Conn {
	return nil
}

func assignAuthenticationValues(label string, dest []any, values []any) {
	if len(dest) != len(values) {
		panic(label + ": destination count mismatch")
	}
	for i := range dest {
		assignAuthenticationValue(label, dest[i], values[i])
	}
}

func assignAuthenticationValue(label string, dest any, value any) {
	switch pointer := dest.(type) {
	case *string:
		switch typed := value.(type) {
		case nil:
			*pointer = ""
		case string:
			*pointer = typed
		default:
			panic(label + ": unsupported string value")
		}
	case *int:
		switch typed := value.(type) {
		case int:
			*pointer = typed
		case int32:
			*pointer = int(typed)
		default:
			panic(label + ": unsupported int value")
		}
	case *[]byte:
		switch typed := value.(type) {
		case nil:
			*pointer = nil
		case []byte:
			*pointer = append([]byte(nil), typed...)
		default:
			panic(label + ": unsupported bytes value")
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
	default:
		panic(label + ": unsupported destination type")
	}
}
