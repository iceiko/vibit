package authentication

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const ModuleName = "runtime.authentication"

type CredentialStatus string

const (
	CredentialStatusActive   CredentialStatus = "active"
	CredentialStatusDisabled CredentialStatus = "disabled"
	CredentialStatusRevoked  CredentialStatus = "revoked"
	CredentialStatusReplaced CredentialStatus = "replaced"
)

func (s CredentialStatus) IsValid() bool {
	switch s {
	case CredentialStatusActive, CredentialStatusDisabled, CredentialStatusRevoked, CredentialStatusReplaced:
		return true
	default:
		return false
	}
}

type TokenStatus string

const (
	TokenStatusActive  TokenStatus = "active"
	TokenStatusExpired TokenStatus = "expired"
	TokenStatusRevoked TokenStatus = "revoked"
)

func (s TokenStatus) IsValid() bool {
	switch s {
	case TokenStatusActive, TokenStatusExpired, TokenStatusRevoked:
		return true
	default:
		return false
	}
}

type CredentialRecord struct {
	CredentialRecordID         string
	PlayerID                   string
	CredentialKind             string
	CredentialStatus           CredentialStatus
	CredentialLookupDigest     []byte
	CredentialVerifierDigest   []byte
	VerifierAlgorithm          string
	VerifierVersion            int
	VerifierKeyID              string
	ClientInstanceIDDigest     []byte
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	LastVerifiedAt             *time.Time
	DisabledAt                 *time.Time
	DisabledReason             string
	RevokedAt                  *time.Time
	RevokedReason              string
	ReplacedAt                 *time.Time
	ReplacedByCredentialRecord string
}

type TokenRecord struct {
	TokenRecordID           string
	TokenKind               string
	TokenStatus             TokenStatus
	ActorKind               string
	PlayerID                string
	CredentialRecordID      string
	TokenLookupDigest       []byte
	TokenVerifierDigest     []byte
	VerifierAlgorithm       string
	VerifierVersion         int
	VerifierKeyID           string
	Audience                string
	IssuedAt                time.Time
	ExpiresAt               time.Time
	RevokedAt               *time.Time
	RevokedReason           string
	ReplacedByTokenRecordID string
	LastValidatedAt         *time.Time
	LastFailedValidationAt  *time.Time
	CleanupAfter            *time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type StoreCredentialMutation struct {
	CredentialRecordID       string
	PlayerID                 string
	CredentialKind           string
	CredentialLookupDigest   []byte
	CredentialVerifierDigest []byte
	VerifierAlgorithm        string
	VerifierVersion          int
	VerifierKeyID            string
	ClientInstanceIDDigest   []byte
	OccurredAt               time.Time
	RequestedBy              string
}

type StoreTokenMutation struct {
	TokenRecordID         string
	PlayerID              string
	CredentialRecordID    string
	TokenKind             string
	ActorKind             string
	TokenLookupDigest     []byte
	TokenVerifierDigest   []byte
	VerifierAlgorithm     string
	VerifierVersion       int
	VerifierKeyID         string
	Audience              string
	IssuedAt              time.Time
	ExpiresAt             time.Time
	ReplacesTokenRecordID string
	RequestedBy           string
}

type RevokeCredentialMutation struct {
	CredentialRecordID string
	RevokedAt          time.Time
	RevokedReason      string
	RequestedBy        string
}

type RevokeTokenMutation struct {
	TokenRecordID string
	RevokedAt     time.Time
	RevokedReason string
	CleanupAfter  *time.Time
	RequestedBy   string
}

type TokenCleanupQuery struct {
	Now   time.Time
	Limit int
}

type Repository interface {
	StoreCredential(context.Context, StoreCredentialMutation) (CredentialRecord, error)
	FindCredentialByLookupDigest(context.Context, []byte) (CredentialRecord, error)
	StoreToken(context.Context, StoreTokenMutation) (TokenRecord, error)
	FindTokenByLookupDigest(context.Context, []byte) (TokenRecord, error)
	RevokeCredential(context.Context, RevokeCredentialMutation) error
	RevokeToken(context.Context, RevokeTokenMutation) error
	ListTokensEligibleForCleanup(context.Context, TokenCleanupQuery) ([]TokenRecord, error)
}

func NormalizeStoreCredentialMutation(mutation StoreCredentialMutation) (StoreCredentialMutation, error) {
	var err error
	mutation.CredentialRecordID, err = normalizeRequired("credential_record_id", mutation.CredentialRecordID)
	if err != nil {
		return StoreCredentialMutation{}, err
	}
	mutation.PlayerID, err = normalizeRequired("player_id", mutation.PlayerID)
	if err != nil {
		return StoreCredentialMutation{}, err
	}
	mutation.CredentialKind, err = normalizeRequired("credential_kind", mutation.CredentialKind)
	if err != nil {
		return StoreCredentialMutation{}, err
	}
	if mutation.CredentialKind != "device_credential_login" {
		return StoreCredentialMutation{}, errors.New("authentication: credential_kind must be device_credential_login")
	}
	mutation.CredentialLookupDigest, err = normalizeRequiredDigest("credential_lookup_digest", mutation.CredentialLookupDigest)
	if err != nil {
		return StoreCredentialMutation{}, err
	}
	mutation.CredentialVerifierDigest, err = normalizeRequiredDigest("credential_verifier_digest", mutation.CredentialVerifierDigest)
	if err != nil {
		return StoreCredentialMutation{}, err
	}
	mutation.VerifierAlgorithm, err = normalizeRequired("verifier_algorithm", mutation.VerifierAlgorithm)
	if err != nil {
		return StoreCredentialMutation{}, err
	}
	if mutation.VerifierVersion <= 0 {
		return StoreCredentialMutation{}, errors.New("authentication: verifier_version must be positive")
	}
	mutation.VerifierKeyID, err = normalizeRequired("verifier_key_id", mutation.VerifierKeyID)
	if err != nil {
		return StoreCredentialMutation{}, err
	}
	if len(mutation.ClientInstanceIDDigest) > 0 {
		mutation.ClientInstanceIDDigest = cloneBytes(mutation.ClientInstanceIDDigest)
	}
	if mutation.OccurredAt.IsZero() {
		return StoreCredentialMutation{}, errors.New("authentication: occurred_at is required")
	}
	mutation.OccurredAt = mutation.OccurredAt.UTC()
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return StoreCredentialMutation{}, err
	}
	return mutation, nil
}

func NormalizeStoreTokenMutation(mutation StoreTokenMutation) (StoreTokenMutation, error) {
	var err error
	mutation.TokenRecordID, err = normalizeRequired("token_record_id", mutation.TokenRecordID)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	mutation.PlayerID, err = normalizeRequired("player_id", mutation.PlayerID)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	mutation.CredentialRecordID, err = normalizeRequired("credential_record_id", mutation.CredentialRecordID)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	mutation.TokenKind, err = normalizeRequired("token_kind", mutation.TokenKind)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	if mutation.TokenKind != "access_token" {
		return StoreTokenMutation{}, errors.New("authentication: token_kind must be access_token")
	}
	mutation.ActorKind, err = normalizeRequired("actor_kind", mutation.ActorKind)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	if mutation.ActorKind != "player" {
		return StoreTokenMutation{}, errors.New("authentication: actor_kind must be player")
	}
	mutation.TokenLookupDigest, err = normalizeRequiredDigest("token_lookup_digest", mutation.TokenLookupDigest)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	mutation.TokenVerifierDigest, err = normalizeRequiredDigest("token_verifier_digest", mutation.TokenVerifierDigest)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	mutation.VerifierAlgorithm, err = normalizeRequired("verifier_algorithm", mutation.VerifierAlgorithm)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	if mutation.VerifierVersion <= 0 {
		return StoreTokenMutation{}, errors.New("authentication: verifier_version must be positive")
	}
	mutation.VerifierKeyID, err = normalizeRequired("verifier_key_id", mutation.VerifierKeyID)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	mutation.Audience, err = normalizeRequired("audience", mutation.Audience)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	if mutation.IssuedAt.IsZero() {
		return StoreTokenMutation{}, errors.New("authentication: issued_at is required")
	}
	if mutation.ExpiresAt.IsZero() {
		return StoreTokenMutation{}, errors.New("authentication: expires_at is required")
	}
	mutation.IssuedAt = mutation.IssuedAt.UTC()
	mutation.ExpiresAt = mutation.ExpiresAt.UTC()
	if !mutation.ExpiresAt.After(mutation.IssuedAt) {
		return StoreTokenMutation{}, errors.New("authentication: expires_at must be after issued_at")
	}
	if mutation.ReplacesTokenRecordID != "" {
		mutation.ReplacesTokenRecordID, err = normalizeRequired("replaces_token_record_id", mutation.ReplacesTokenRecordID)
		if err != nil {
			return StoreTokenMutation{}, err
		}
	}
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return StoreTokenMutation{}, err
	}
	return mutation, nil
}

func NormalizeRevokeCredentialMutation(mutation RevokeCredentialMutation) (RevokeCredentialMutation, error) {
	var err error
	mutation.CredentialRecordID, err = normalizeRequired("credential_record_id", mutation.CredentialRecordID)
	if err != nil {
		return RevokeCredentialMutation{}, err
	}
	if mutation.RevokedAt.IsZero() {
		return RevokeCredentialMutation{}, errors.New("authentication: revoked_at is required")
	}
	mutation.RevokedAt = mutation.RevokedAt.UTC()
	mutation.RevokedReason, err = normalizeRequired("revoked_reason", mutation.RevokedReason)
	if err != nil {
		return RevokeCredentialMutation{}, err
	}
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return RevokeCredentialMutation{}, err
	}
	return mutation, nil
}

func NormalizeRevokeTokenMutation(mutation RevokeTokenMutation) (RevokeTokenMutation, error) {
	var err error
	mutation.TokenRecordID, err = normalizeRequired("token_record_id", mutation.TokenRecordID)
	if err != nil {
		return RevokeTokenMutation{}, err
	}
	if mutation.RevokedAt.IsZero() {
		return RevokeTokenMutation{}, errors.New("authentication: revoked_at is required")
	}
	mutation.RevokedAt = mutation.RevokedAt.UTC()
	mutation.RevokedReason, err = normalizeRequired("revoked_reason", mutation.RevokedReason)
	if err != nil {
		return RevokeTokenMutation{}, err
	}
	if mutation.CleanupAfter != nil {
		cleanupAfter := mutation.CleanupAfter.UTC()
		mutation.CleanupAfter = &cleanupAfter
	}
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return RevokeTokenMutation{}, err
	}
	return mutation, nil
}

func NormalizeTokenCleanupQuery(query TokenCleanupQuery) (TokenCleanupQuery, error) {
	if query.Now.IsZero() {
		return TokenCleanupQuery{}, errors.New("authentication: now is required")
	}
	query.Now = query.Now.UTC()
	if query.Limit <= 0 {
		return TokenCleanupQuery{}, errors.New("authentication: limit must be positive")
	}
	return query, nil
}

func NormalizeLookupDigest(name string, value []byte) ([]byte, error) {
	return normalizeRequiredDigest(name, value)
}

func normalizeRequired(name string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("authentication: %s is required", name)
	}
	return value, nil
}

func normalizeRequiredDigest(name string, value []byte) ([]byte, error) {
	if len(value) == 0 {
		return nil, fmt.Errorf("authentication: %s is required", name)
	}
	return cloneBytes(value), nil
}

func cloneBytes(value []byte) []byte {
	cloned := make([]byte, len(value))
	copy(cloned, value)
	return cloned
}
