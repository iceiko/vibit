package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/modules/authentication"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrAuthenticationRecordNotFound   = errors.New("postgres authentication: not found")
	ErrAuthenticationRecordConflict   = errors.New("postgres authentication: conflict")
	ErrAuthenticationRecordConstraint = errors.New("postgres authentication: constraint violation")
)

type AuthenticationRepository struct {
	executor Executor
}

var _ authentication.Repository = (*AuthenticationRepository)(nil)

func NewAuthenticationRepositoryForUnitOfWork(executor Executor) *AuthenticationRepository {
	return &AuthenticationRepository{executor: executor}
}

func (r *AuthenticationRepository) StoreCredential(ctx context.Context, mutation authentication.StoreCredentialMutation) (authentication.CredentialRecord, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return authentication.CredentialRecord{}, err
	}

	normalized, err := authentication.NormalizeStoreCredentialMutation(mutation)
	if err != nil {
		return authentication.CredentialRecord{}, fmt.Errorf("postgres authentication: normalize credential store mutation: %w", err)
	}

	record, err := scanAuthenticationCredentialRow(executor.QueryRow(
		ctx,
		insertAuthenticationCredentialSQL,
		normalized.CredentialRecordID,
		normalized.PlayerID,
		normalized.CredentialKind,
		string(authentication.CredentialStatusActive),
		normalized.CredentialLookupDigest,
		normalized.CredentialVerifierDigest,
		normalized.VerifierAlgorithm,
		int32(normalized.VerifierVersion),
		normalized.VerifierKeyID,
		nullableBytes(normalized.ClientInstanceIDDigest),
		normalized.OccurredAt,
	))
	if err != nil {
		return authentication.CredentialRecord{}, mapAuthenticationPostgresError("store credential record", err)
	}
	return record, nil
}

func (r *AuthenticationRepository) FindCredentialByLookupDigest(ctx context.Context, digest []byte) (authentication.CredentialRecord, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return authentication.CredentialRecord{}, err
	}

	normalizedDigest, err := authentication.NormalizeLookupDigest("credential_lookup_digest", digest)
	if err != nil {
		return authentication.CredentialRecord{}, fmt.Errorf("postgres authentication: normalize credential lookup digest: %w", err)
	}

	record, err := scanAuthenticationCredentialRow(executor.QueryRow(
		ctx,
		findAuthenticationCredentialByLookupDigestSQL,
		normalizedDigest,
	))
	if err != nil {
		return authentication.CredentialRecord{}, mapAuthenticationPostgresError("find credential record by lookup digest", err)
	}
	return record, nil
}

func (r *AuthenticationRepository) StoreToken(ctx context.Context, mutation authentication.StoreTokenMutation) (authentication.TokenRecord, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return authentication.TokenRecord{}, err
	}

	normalized, err := authentication.NormalizeStoreTokenMutation(mutation)
	if err != nil {
		return authentication.TokenRecord{}, fmt.Errorf("postgres authentication: normalize token store mutation: %w", err)
	}

	record, err := scanAuthenticationTokenRow(executor.QueryRow(
		ctx,
		insertAuthenticationTokenSQL,
		normalized.TokenRecordID,
		normalized.TokenKind,
		string(authentication.TokenStatusActive),
		normalized.ActorKind,
		normalized.PlayerID,
		normalized.CredentialRecordID,
		normalized.TokenLookupDigest,
		normalized.TokenVerifierDigest,
		normalized.VerifierAlgorithm,
		int32(normalized.VerifierVersion),
		normalized.VerifierKeyID,
		normalized.Audience,
		normalized.IssuedAt,
		normalized.ExpiresAt,
	))
	if err != nil {
		return authentication.TokenRecord{}, mapAuthenticationPostgresError("store token record", err)
	}
	return record, nil
}

func (r *AuthenticationRepository) FindTokenByLookupDigest(ctx context.Context, digest []byte) (authentication.TokenRecord, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return authentication.TokenRecord{}, err
	}

	normalizedDigest, err := authentication.NormalizeLookupDigest("token_lookup_digest", digest)
	if err != nil {
		return authentication.TokenRecord{}, fmt.Errorf("postgres authentication: normalize token lookup digest: %w", err)
	}

	record, err := scanAuthenticationTokenRow(executor.QueryRow(
		ctx,
		findAuthenticationTokenByLookupDigestSQL,
		normalizedDigest,
	))
	if err != nil {
		return authentication.TokenRecord{}, mapAuthenticationPostgresError("find token record by lookup digest", err)
	}
	return record, nil
}

func (r *AuthenticationRepository) RevokeCredential(ctx context.Context, mutation authentication.RevokeCredentialMutation) error {
	executor, err := r.requireExecutor()
	if err != nil {
		return err
	}

	normalized, err := authentication.NormalizeRevokeCredentialMutation(mutation)
	if err != nil {
		return fmt.Errorf("postgres authentication: normalize credential revoke mutation: %w", err)
	}

	commandTag, err := executor.Exec(
		ctx,
		revokeAuthenticationCredentialSQL,
		string(authentication.CredentialStatusRevoked),
		normalized.RevokedAt,
		normalized.RevokedReason,
		normalized.CredentialRecordID,
	)
	if err != nil {
		return mapAuthenticationPostgresError("revoke credential record", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: revoke credential record", ErrAuthenticationRecordNotFound)
	}
	return nil
}

func (r *AuthenticationRepository) RevokeToken(ctx context.Context, mutation authentication.RevokeTokenMutation) error {
	executor, err := r.requireExecutor()
	if err != nil {
		return err
	}

	normalized, err := authentication.NormalizeRevokeTokenMutation(mutation)
	if err != nil {
		return fmt.Errorf("postgres authentication: normalize token revoke mutation: %w", err)
	}

	commandTag, err := executor.Exec(
		ctx,
		revokeAuthenticationTokenSQL,
		string(authentication.TokenStatusRevoked),
		normalized.RevokedAt,
		normalized.RevokedReason,
		nullableTime(normalized.CleanupAfter),
		normalized.TokenRecordID,
	)
	if err != nil {
		return mapAuthenticationPostgresError("revoke token record", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: revoke token record", ErrAuthenticationRecordNotFound)
	}
	return nil
}

func (r *AuthenticationRepository) ListTokensEligibleForCleanup(ctx context.Context, query authentication.TokenCleanupQuery) ([]authentication.TokenRecord, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return nil, err
	}

	normalized, err := authentication.NormalizeTokenCleanupQuery(query)
	if err != nil {
		return nil, fmt.Errorf("postgres authentication: normalize token cleanup query: %w", err)
	}

	rows, err := executor.Query(ctx, listAuthenticationTokensEligibleForCleanupSQL, normalized.Now, int32(normalized.Limit))
	if err != nil {
		return nil, mapAuthenticationPostgresError("list token records eligible for cleanup", err)
	}
	defer rows.Close()

	records := []authentication.TokenRecord{}
	for rows.Next() {
		record, err := scanAuthenticationTokenScanner(rows)
		if err != nil {
			return nil, mapAuthenticationPostgresError("scan token cleanup record", err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, mapAuthenticationPostgresError("read token cleanup records", err)
	}
	return records, nil
}

func (r *AuthenticationRepository) requireExecutor() (Executor, error) {
	if r == nil || r.executor == nil {
		return nil, errors.New("postgres authentication: unit-of-work executor is required")
	}
	return r.executor, nil
}

type scanner interface {
	Scan(...any) error
}

func scanAuthenticationCredentialRow(row pgx.Row) (authentication.CredentialRecord, error) {
	return scanAuthenticationCredentialScanner(row)
}

func scanAuthenticationCredentialScanner(row scanner) (authentication.CredentialRecord, error) {
	var record authentication.CredentialRecord
	var status string
	var lastVerifiedAt pgtype.Timestamptz
	var disabledAt pgtype.Timestamptz
	var revokedAt pgtype.Timestamptz
	var replacedAt pgtype.Timestamptz
	var clientInstanceDigest []byte
	var replacedBy string

	if err := row.Scan(
		&record.CredentialRecordID,
		&record.PlayerID,
		&record.CredentialKind,
		&status,
		&record.CredentialLookupDigest,
		&record.CredentialVerifierDigest,
		&record.VerifierAlgorithm,
		&record.VerifierVersion,
		&record.VerifierKeyID,
		&clientInstanceDigest,
		&record.CreatedAt,
		&record.UpdatedAt,
		&lastVerifiedAt,
		&disabledAt,
		&record.DisabledReason,
		&revokedAt,
		&record.RevokedReason,
		&replacedAt,
		&replacedBy,
	); err != nil {
		return authentication.CredentialRecord{}, err
	}

	if err := validateAuthenticationText("credential_record_id", record.CredentialRecordID); err != nil {
		return authentication.CredentialRecord{}, err
	}
	if err := validateAuthenticationText("player_id", record.PlayerID); err != nil {
		return authentication.CredentialRecord{}, err
	}
	if record.CredentialKind != "device_credential_login" {
		return authentication.CredentialRecord{}, fmt.Errorf("%w: row credential_kind %q is invalid", ErrAuthenticationRecordConstraint, record.CredentialKind)
	}
	record.CredentialStatus = authentication.CredentialStatus(status)
	if !record.CredentialStatus.IsValid() {
		return authentication.CredentialRecord{}, fmt.Errorf("%w: row credential_status %q is invalid", ErrAuthenticationRecordConstraint, status)
	}
	if len(record.CredentialLookupDigest) == 0 {
		return authentication.CredentialRecord{}, fmt.Errorf("%w: row credential_lookup_digest is required", ErrAuthenticationRecordConstraint)
	}
	if len(record.CredentialVerifierDigest) == 0 {
		return authentication.CredentialRecord{}, fmt.Errorf("%w: row credential_verifier_digest is required", ErrAuthenticationRecordConstraint)
	}
	if err := validateAuthenticationText("verifier_algorithm", record.VerifierAlgorithm); err != nil {
		return authentication.CredentialRecord{}, err
	}
	if record.VerifierVersion <= 0 {
		return authentication.CredentialRecord{}, fmt.Errorf("%w: row verifier_version must be positive", ErrAuthenticationRecordConstraint)
	}
	if err := validateAuthenticationText("verifier_key_id", record.VerifierKeyID); err != nil {
		return authentication.CredentialRecord{}, err
	}
	if record.CreatedAt.IsZero() || record.UpdatedAt.IsZero() {
		return authentication.CredentialRecord{}, fmt.Errorf("%w: row credential timestamps are required", ErrAuthenticationRecordConstraint)
	}

	record.CredentialLookupDigest = cloneBytes(record.CredentialLookupDigest)
	record.CredentialVerifierDigest = cloneBytes(record.CredentialVerifierDigest)
	record.ClientInstanceIDDigest = cloneBytes(clientInstanceDigest)
	record.CreatedAt = record.CreatedAt.UTC()
	record.UpdatedAt = record.UpdatedAt.UTC()
	record.LastVerifiedAt = nullableTimestamptzUTC(lastVerifiedAt)
	record.DisabledAt = nullableTimestamptzUTC(disabledAt)
	record.RevokedAt = nullableTimestamptzUTC(revokedAt)
	record.ReplacedAt = nullableTimestamptzUTC(replacedAt)
	record.ReplacedByCredentialRecord = strings.TrimSpace(replacedBy)
	return record, nil
}

func scanAuthenticationTokenRow(row pgx.Row) (authentication.TokenRecord, error) {
	return scanAuthenticationTokenScanner(row)
}

func scanAuthenticationTokenScanner(row scanner) (authentication.TokenRecord, error) {
	var record authentication.TokenRecord
	var status string
	var revokedAt pgtype.Timestamptz
	var lastValidatedAt pgtype.Timestamptz
	var lastFailedValidationAt pgtype.Timestamptz
	var cleanupAfter pgtype.Timestamptz
	var replacedBy string

	if err := row.Scan(
		&record.TokenRecordID,
		&record.TokenKind,
		&status,
		&record.ActorKind,
		&record.PlayerID,
		&record.CredentialRecordID,
		&record.TokenLookupDigest,
		&record.TokenVerifierDigest,
		&record.VerifierAlgorithm,
		&record.VerifierVersion,
		&record.VerifierKeyID,
		&record.Audience,
		&record.IssuedAt,
		&record.ExpiresAt,
		&revokedAt,
		&record.RevokedReason,
		&replacedBy,
		&lastValidatedAt,
		&lastFailedValidationAt,
		&cleanupAfter,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return authentication.TokenRecord{}, err
	}

	if err := validateAuthenticationText("token_record_id", record.TokenRecordID); err != nil {
		return authentication.TokenRecord{}, err
	}
	if record.TokenKind != "access_token" {
		return authentication.TokenRecord{}, fmt.Errorf("%w: row token_kind %q is invalid", ErrAuthenticationRecordConstraint, record.TokenKind)
	}
	record.TokenStatus = authentication.TokenStatus(status)
	if !record.TokenStatus.IsValid() {
		return authentication.TokenRecord{}, fmt.Errorf("%w: row token_status %q is invalid", ErrAuthenticationRecordConstraint, status)
	}
	if record.ActorKind != "player" {
		return authentication.TokenRecord{}, fmt.Errorf("%w: row actor_kind %q is invalid", ErrAuthenticationRecordConstraint, record.ActorKind)
	}
	if err := validateAuthenticationText("player_id", record.PlayerID); err != nil {
		return authentication.TokenRecord{}, err
	}
	if err := validateAuthenticationText("credential_record_id", record.CredentialRecordID); err != nil {
		return authentication.TokenRecord{}, err
	}
	if len(record.TokenLookupDigest) == 0 {
		return authentication.TokenRecord{}, fmt.Errorf("%w: row token_lookup_digest is required", ErrAuthenticationRecordConstraint)
	}
	if len(record.TokenVerifierDigest) == 0 {
		return authentication.TokenRecord{}, fmt.Errorf("%w: row token_verifier_digest is required", ErrAuthenticationRecordConstraint)
	}
	if err := validateAuthenticationText("verifier_algorithm", record.VerifierAlgorithm); err != nil {
		return authentication.TokenRecord{}, err
	}
	if record.VerifierVersion <= 0 {
		return authentication.TokenRecord{}, fmt.Errorf("%w: row verifier_version must be positive", ErrAuthenticationRecordConstraint)
	}
	if err := validateAuthenticationText("verifier_key_id", record.VerifierKeyID); err != nil {
		return authentication.TokenRecord{}, err
	}
	if err := validateAuthenticationText("audience", record.Audience); err != nil {
		return authentication.TokenRecord{}, err
	}
	if record.IssuedAt.IsZero() || record.ExpiresAt.IsZero() || record.CreatedAt.IsZero() || record.UpdatedAt.IsZero() {
		return authentication.TokenRecord{}, fmt.Errorf("%w: row token timestamps are required", ErrAuthenticationRecordConstraint)
	}
	if !record.ExpiresAt.After(record.IssuedAt) {
		return authentication.TokenRecord{}, fmt.Errorf("%w: row expires_at must be after issued_at", ErrAuthenticationRecordConstraint)
	}

	record.TokenLookupDigest = cloneBytes(record.TokenLookupDigest)
	record.TokenVerifierDigest = cloneBytes(record.TokenVerifierDigest)
	record.IssuedAt = record.IssuedAt.UTC()
	record.ExpiresAt = record.ExpiresAt.UTC()
	record.RevokedAt = nullableTimestamptzUTC(revokedAt)
	record.ReplacedByTokenRecordID = strings.TrimSpace(replacedBy)
	record.LastValidatedAt = nullableTimestamptzUTC(lastValidatedAt)
	record.LastFailedValidationAt = nullableTimestamptzUTC(lastFailedValidationAt)
	record.CleanupAfter = nullableTimestamptzUTC(cleanupAfter)
	record.CreatedAt = record.CreatedAt.UTC()
	record.UpdatedAt = record.UpdatedAt.UTC()
	return record, nil
}

func validateAuthenticationText(name string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: row %s is required", ErrAuthenticationRecordConstraint, name)
	}
	return nil
}

func mapAuthenticationPostgresError(action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrAuthenticationRecordNotFound, action)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s: %s", ErrAuthenticationRecordConflict, action, postgresConstraintLabel(pgErr))
		case "23502", "23503", "23514":
			return fmt.Errorf("%w: %s: %s", ErrAuthenticationRecordConstraint, action, postgresConstraintLabel(pgErr))
		default:
			return fmt.Errorf("postgres authentication: %s failed with PostgreSQL code %s", action, pgErr.Code)
		}
	}

	if errors.Is(err, ErrAuthenticationRecordConstraint) {
		return err
	}
	return fmt.Errorf("postgres authentication: %s: %w", action, err)
}

func nullableBytes(value []byte) []byte {
	if len(value) == 0 {
		return nil
	}
	return cloneBytes(value)
}

func nullableTime(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: value.UTC(), Valid: true}
}

func cloneBytes(value []byte) []byte {
	if len(value) == 0 {
		return nil
	}
	cloned := make([]byte, len(value))
	copy(cloned, value)
	return cloned
}

const authenticationCredentialColumns = `
credential_record_id,
player_id,
credential_kind,
credential_status,
credential_lookup_digest,
credential_verifier_digest,
verifier_algorithm,
verifier_version,
verifier_key_id,
client_instance_id_digest,
created_at,
updated_at,
last_verified_at,
disabled_at,
disabled_reason,
revoked_at,
revoked_reason,
replaced_at,
replaced_by_credential_record_id`

const insertAuthenticationCredentialSQL = `
INSERT INTO authentication_device_credentials (
    credential_record_id,
    player_id,
    credential_kind,
    credential_status,
    credential_lookup_digest,
    credential_verifier_digest,
    verifier_algorithm,
    verifier_version,
    verifier_key_id,
    client_instance_id_digest,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
RETURNING ` + authenticationCredentialColumns

const findAuthenticationCredentialByLookupDigestSQL = `
SELECT ` + authenticationCredentialColumns + `
FROM authentication_device_credentials
WHERE credential_lookup_digest = $1`

const revokeAuthenticationCredentialSQL = `
UPDATE authentication_device_credentials
SET credential_status = $1,
    revoked_at = $2,
    revoked_reason = $3,
    updated_at = $2
WHERE credential_record_id = $4`

const authenticationTokenColumns = `
token_record_id,
token_kind,
token_status,
actor_kind,
player_id,
credential_record_id,
token_lookup_digest,
token_verifier_digest,
verifier_algorithm,
verifier_version,
verifier_key_id,
audience,
issued_at,
expires_at,
revoked_at,
revoked_reason,
replaced_by_token_record_id,
last_validated_at,
last_failed_validation_at,
cleanup_after,
created_at,
updated_at`

const insertAuthenticationTokenSQL = `
INSERT INTO authentication_access_tokens (
    token_record_id,
    token_kind,
    token_status,
    actor_kind,
    player_id,
    credential_record_id,
    token_lookup_digest,
    token_verifier_digest,
    verifier_algorithm,
    verifier_version,
    verifier_key_id,
    audience,
    issued_at,
    expires_at,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $13, $13)
RETURNING ` + authenticationTokenColumns

const findAuthenticationTokenByLookupDigestSQL = `
SELECT ` + authenticationTokenColumns + `
FROM authentication_access_tokens
WHERE token_lookup_digest = $1`

const revokeAuthenticationTokenSQL = `
UPDATE authentication_access_tokens
SET token_status = $1,
    revoked_at = $2,
    revoked_reason = $3,
    cleanup_after = $4,
    updated_at = $2
WHERE token_record_id = $5`

const listAuthenticationTokensEligibleForCleanupSQL = `
SELECT ` + authenticationTokenColumns + `
FROM authentication_access_tokens
WHERE cleanup_after IS NOT NULL
  AND cleanup_after <= $1
ORDER BY cleanup_after, token_record_id
LIMIT $2`
