package authentication

import (
	"context"
	"testing"
	"time"
)

func TestRepositoryInterfaceIsStorageNeutral(t *testing.T) {
	var _ Repository = recordingRepository{}
}

func TestNormalizeStoreCredentialMutation(t *testing.T) {
	occurredAt := time.Date(2026, 5, 14, 10, 11, 12, 0, time.FixedZone("test", 8*60*60))
	lookupDigest := []byte{1, 2, 3}
	verifierDigest := []byte{4, 5, 6}

	mutation, err := NormalizeStoreCredentialMutation(StoreCredentialMutation{
		CredentialRecordID:       " credential-1 ",
		PlayerID:                 " player-1 ",
		CredentialKind:           " device_credential_login ",
		CredentialLookupDigest:   lookupDigest,
		CredentialVerifierDigest: verifierDigest,
		VerifierAlgorithm:        " hmac-sha256 ",
		VerifierVersion:          1,
		VerifierKeyID:            " key-1 ",
		ClientInstanceIDDigest:   []byte{7, 8},
		OccurredAt:               occurredAt,
		RequestedBy:              " maintainer ",
	})
	if err != nil {
		t.Fatalf("NormalizeStoreCredentialMutation() error = %v, want nil", err)
	}

	if mutation.CredentialRecordID != "credential-1" ||
		mutation.PlayerID != "player-1" ||
		mutation.CredentialKind != "device_credential_login" ||
		mutation.VerifierAlgorithm != "hmac-sha256" ||
		mutation.VerifierKeyID != "key-1" ||
		mutation.RequestedBy != "maintainer" {
		t.Fatalf("normalized mutation = %#v, want trimmed fields", mutation)
	}
	if mutation.OccurredAt.Location() != time.UTC {
		t.Fatalf("OccurredAt location = %s, want UTC", mutation.OccurredAt.Location())
	}
	lookupDigest[0] = 99
	verifierDigest[0] = 99
	if mutation.CredentialLookupDigest[0] == 99 || mutation.CredentialVerifierDigest[0] == 99 {
		t.Fatal("NormalizeStoreCredentialMutation() reused caller digest backing array")
	}
}

func TestNormalizeStoreCredentialMutationRejectsInvalidShape(t *testing.T) {
	valid := validStoreCredentialMutation()
	tests := []struct {
		name   string
		mutate func(*StoreCredentialMutation)
	}{
		{name: "credential_record_id", mutate: func(m *StoreCredentialMutation) { m.CredentialRecordID = " " }},
		{name: "player_id", mutate: func(m *StoreCredentialMutation) { m.PlayerID = " " }},
		{name: "credential_kind", mutate: func(m *StoreCredentialMutation) { m.CredentialKind = "email_password_login" }},
		{name: "lookup_digest", mutate: func(m *StoreCredentialMutation) { m.CredentialLookupDigest = nil }},
		{name: "verifier_digest", mutate: func(m *StoreCredentialMutation) { m.CredentialVerifierDigest = nil }},
		{name: "verifier_algorithm", mutate: func(m *StoreCredentialMutation) { m.VerifierAlgorithm = " " }},
		{name: "verifier_version", mutate: func(m *StoreCredentialMutation) { m.VerifierVersion = 0 }},
		{name: "verifier_key_id", mutate: func(m *StoreCredentialMutation) { m.VerifierKeyID = " " }},
		{name: "occurred_at", mutate: func(m *StoreCredentialMutation) { m.OccurredAt = time.Time{} }},
		{name: "requested_by", mutate: func(m *StoreCredentialMutation) { m.RequestedBy = " " }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutation := valid
			tt.mutate(&mutation)
			if _, err := NormalizeStoreCredentialMutation(mutation); err == nil {
				t.Fatal("NormalizeStoreCredentialMutation() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeStoreTokenMutation(t *testing.T) {
	issuedAt := time.Date(2026, 5, 14, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	expiresAt := issuedAt.Add(time.Hour)
	lookupDigest := []byte{1, 2, 3}

	mutation, err := NormalizeStoreTokenMutation(StoreTokenMutation{
		TokenRecordID:       " token-1 ",
		PlayerID:            " player-1 ",
		CredentialRecordID:  " credential-1 ",
		TokenKind:           " access_token ",
		ActorKind:           " player ",
		TokenLookupDigest:   lookupDigest,
		TokenVerifierDigest: []byte{4, 5, 6},
		VerifierAlgorithm:   " hmac-sha256 ",
		VerifierVersion:     1,
		VerifierKeyID:       " key-1 ",
		Audience:            " gameplay ",
		IssuedAt:            issuedAt,
		ExpiresAt:           expiresAt,
		RequestedBy:         " maintainer ",
	})
	if err != nil {
		t.Fatalf("NormalizeStoreTokenMutation() error = %v, want nil", err)
	}

	if mutation.TokenRecordID != "token-1" ||
		mutation.PlayerID != "player-1" ||
		mutation.CredentialRecordID != "credential-1" ||
		mutation.TokenKind != "access_token" ||
		mutation.ActorKind != "player" ||
		mutation.Audience != "gameplay" ||
		mutation.RequestedBy != "maintainer" {
		t.Fatalf("normalized mutation = %#v, want trimmed fields", mutation)
	}
	if mutation.IssuedAt.Location() != time.UTC || mutation.ExpiresAt.Location() != time.UTC {
		t.Fatalf("token times were not normalized to UTC: %#v", mutation)
	}
	lookupDigest[0] = 99
	if mutation.TokenLookupDigest[0] == 99 {
		t.Fatal("NormalizeStoreTokenMutation() reused caller digest backing array")
	}
}

func TestNormalizeStoreTokenMutationRejectsInvalidShape(t *testing.T) {
	valid := validStoreTokenMutation()
	tests := []struct {
		name   string
		mutate func(*StoreTokenMutation)
	}{
		{name: "token_record_id", mutate: func(m *StoreTokenMutation) { m.TokenRecordID = " " }},
		{name: "player_id", mutate: func(m *StoreTokenMutation) { m.PlayerID = " " }},
		{name: "credential_record_id", mutate: func(m *StoreTokenMutation) { m.CredentialRecordID = " " }},
		{name: "token_kind", mutate: func(m *StoreTokenMutation) { m.TokenKind = "refresh_token" }},
		{name: "actor_kind", mutate: func(m *StoreTokenMutation) { m.ActorKind = "service" }},
		{name: "lookup_digest", mutate: func(m *StoreTokenMutation) { m.TokenLookupDigest = nil }},
		{name: "verifier_digest", mutate: func(m *StoreTokenMutation) { m.TokenVerifierDigest = nil }},
		{name: "verifier_algorithm", mutate: func(m *StoreTokenMutation) { m.VerifierAlgorithm = " " }},
		{name: "verifier_version", mutate: func(m *StoreTokenMutation) { m.VerifierVersion = 0 }},
		{name: "verifier_key_id", mutate: func(m *StoreTokenMutation) { m.VerifierKeyID = " " }},
		{name: "audience", mutate: func(m *StoreTokenMutation) { m.Audience = " " }},
		{name: "issued_at", mutate: func(m *StoreTokenMutation) { m.IssuedAt = time.Time{} }},
		{name: "expires_at", mutate: func(m *StoreTokenMutation) { m.ExpiresAt = time.Time{} }},
		{name: "expires_before_issued", mutate: func(m *StoreTokenMutation) { m.ExpiresAt = m.IssuedAt }},
		{name: "requested_by", mutate: func(m *StoreTokenMutation) { m.RequestedBy = " " }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutation := valid
			tt.mutate(&mutation)
			if _, err := NormalizeStoreTokenMutation(mutation); err == nil {
				t.Fatal("NormalizeStoreTokenMutation() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeRevokeMutationsAndCleanupQuery(t *testing.T) {
	revokedAt := time.Date(2026, 5, 14, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	cleanupAfter := revokedAt.Add(7 * 24 * time.Hour)

	credentialMutation, err := NormalizeRevokeCredentialMutation(RevokeCredentialMutation{
		CredentialRecordID: " credential-1 ",
		RevokedAt:          revokedAt,
		RevokedReason:      " rotated ",
		RequestedBy:        " maintainer ",
	})
	if err != nil {
		t.Fatalf("NormalizeRevokeCredentialMutation() error = %v, want nil", err)
	}
	if credentialMutation.CredentialRecordID != "credential-1" ||
		credentialMutation.RevokedReason != "rotated" ||
		credentialMutation.RequestedBy != "maintainer" ||
		credentialMutation.RevokedAt.Location() != time.UTC {
		t.Fatalf("normalized credential revoke mutation = %#v, want trimmed UTC fields", credentialMutation)
	}

	tokenMutation, err := NormalizeRevokeTokenMutation(RevokeTokenMutation{
		TokenRecordID: " token-1 ",
		RevokedAt:     revokedAt,
		RevokedReason: " logout ",
		CleanupAfter:  &cleanupAfter,
		RequestedBy:   " maintainer ",
	})
	if err != nil {
		t.Fatalf("NormalizeRevokeTokenMutation() error = %v, want nil", err)
	}
	if tokenMutation.TokenRecordID != "token-1" ||
		tokenMutation.RevokedReason != "logout" ||
		tokenMutation.RequestedBy != "maintainer" ||
		tokenMutation.RevokedAt.Location() != time.UTC ||
		tokenMutation.CleanupAfter.Location() != time.UTC {
		t.Fatalf("normalized token revoke mutation = %#v, want trimmed UTC fields", tokenMutation)
	}

	query, err := NormalizeTokenCleanupQuery(TokenCleanupQuery{Now: revokedAt, Limit: 50})
	if err != nil {
		t.Fatalf("NormalizeTokenCleanupQuery() error = %v, want nil", err)
	}
	if query.Now.Location() != time.UTC || query.Limit != 50 {
		t.Fatalf("normalized cleanup query = %#v, want UTC time and preserved limit", query)
	}
}

func TestStatusesAreClosedSets(t *testing.T) {
	for _, status := range []CredentialStatus{CredentialStatusActive, CredentialStatusDisabled, CredentialStatusRevoked, CredentialStatusReplaced} {
		if !status.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", status)
		}
	}
	if CredentialStatus("pending").IsValid() {
		t.Fatal(`CredentialStatus("pending").IsValid() = true, want false`)
	}

	for _, status := range []TokenStatus{TokenStatusActive, TokenStatusExpired, TokenStatusRevoked} {
		if !status.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", status)
		}
	}
	if TokenStatus("rotated").IsValid() {
		t.Fatal(`TokenStatus("rotated").IsValid() = true, want false`)
	}
}

func validStoreCredentialMutation() StoreCredentialMutation {
	return StoreCredentialMutation{
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

func validStoreTokenMutation() StoreTokenMutation {
	issuedAt := time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC)
	return StoreTokenMutation{
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

type recordingRepository struct{}

func (recordingRepository) StoreCredential(context.Context, StoreCredentialMutation) (CredentialRecord, error) {
	return CredentialRecord{}, nil
}

func (recordingRepository) FindCredentialByLookupDigest(context.Context, []byte) (CredentialRecord, error) {
	return CredentialRecord{}, nil
}

func (recordingRepository) StoreToken(context.Context, StoreTokenMutation) (TokenRecord, error) {
	return TokenRecord{}, nil
}

func (recordingRepository) FindTokenByLookupDigest(context.Context, []byte) (TokenRecord, error) {
	return TokenRecord{}, nil
}

func (recordingRepository) RevokeCredential(context.Context, RevokeCredentialMutation) error {
	return nil
}

func (recordingRepository) RevokeToken(context.Context, RevokeTokenMutation) error {
	return nil
}

func (recordingRepository) ListTokensEligibleForCleanup(context.Context, TokenCleanupQuery) ([]TokenRecord, error) {
	return nil, nil
}
