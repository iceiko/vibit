package session

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRepositoryInterfaceIsStorageNeutral(t *testing.T) {
	var _ Repository = recordingRepository{}
}

func TestSessionStatusIsValid(t *testing.T) {
	for _, status := range []SessionStatus{SessionStatusActive, SessionStatusExpired, SessionStatusRevoked} {
		if !status.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", status)
		}
	}
	if SessionStatus("disabled").IsValid() {
		t.Fatal(`SessionStatus("disabled").IsValid() = true, want false`)
	}
}

func TestNormalizeCreateRuntimeSessionMutationAcceptsActivePlayerSession(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 10, 0, 0, 0, time.FixedZone("test", 8*60*60))
	mutation, err := NormalizeCreateRuntimeSessionMutation(CreateRuntimeSessionMutation{
		SessionID:           " session-1 ",
		ActorKind:           ActorKind(" player "),
		ActorID:             " player-1 ",
		PlayerID:            " player-1 ",
		SessionStatus:       "",
		IssuedAt:            issuedAt,
		ExpiresAt:           issuedAt.Add(time.Hour),
		LastSeenAt:          issuedAt.Add(time.Minute),
		AccessTokenRecordID: " token-1 ",
		RequestedBy:         " authentication_service ",
	})
	if err != nil {
		t.Fatalf("NormalizeCreateRuntimeSessionMutation() error = %v, want nil", err)
	}
	if mutation.SessionID != "session-1" ||
		mutation.ActorKind != ActorKindPlayer ||
		mutation.ActorID != "player-1" ||
		mutation.PlayerID != "player-1" ||
		mutation.SessionStatus != SessionStatusActive ||
		mutation.AccessTokenRecordID != "token-1" ||
		mutation.RequestedBy != "authentication_service" {
		t.Fatalf("normalized mutation = %#v, want trimmed active player session", mutation)
	}
	if mutation.IssuedAt.Location() != time.UTC ||
		mutation.ExpiresAt.Location() != time.UTC ||
		mutation.LastSeenAt.Location() != time.UTC {
		t.Fatalf("normalized mutation times are not UTC: %#v", mutation)
	}
}

func TestNormalizeCreateRuntimeSessionMutationRejectsMissingRequiredFields(t *testing.T) {
	valid := validCreateRuntimeSessionMutation()
	tests := []struct {
		name   string
		mutate func(*CreateRuntimeSessionMutation)
	}{
		{name: "session_id", mutate: func(m *CreateRuntimeSessionMutation) { m.SessionID = " " }},
		{name: "actor_kind", mutate: func(m *CreateRuntimeSessionMutation) { m.ActorKind = "" }},
		{name: "actor_id", mutate: func(m *CreateRuntimeSessionMutation) { m.ActorID = " " }},
		{name: "player_id", mutate: func(m *CreateRuntimeSessionMutation) { m.PlayerID = " " }},
		{name: "issued_at", mutate: func(m *CreateRuntimeSessionMutation) { m.IssuedAt = time.Time{} }},
		{name: "expires_at", mutate: func(m *CreateRuntimeSessionMutation) { m.ExpiresAt = time.Time{} }},
		{name: "last_seen_at", mutate: func(m *CreateRuntimeSessionMutation) { m.LastSeenAt = time.Time{} }},
		{name: "requested_by", mutate: func(m *CreateRuntimeSessionMutation) { m.RequestedBy = " " }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutation := valid
			tt.mutate(&mutation)
			if _, err := NormalizeCreateRuntimeSessionMutation(mutation); err == nil {
				t.Fatal("NormalizeCreateRuntimeSessionMutation() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeCreateRuntimeSessionMutationRejectsInvalidStatus(t *testing.T) {
	mutation := validCreateRuntimeSessionMutation()
	mutation.SessionStatus = SessionStatusExpired
	if _, err := NormalizeCreateRuntimeSessionMutation(mutation); err == nil {
		t.Fatal("NormalizeCreateRuntimeSessionMutation() error = nil, want expired status rejection")
	}

	mutation = validCreateRuntimeSessionMutation()
	mutation.ActorKind = ActorKind("service")
	if _, err := NormalizeCreateRuntimeSessionMutation(mutation); err == nil {
		t.Fatal("NormalizeCreateRuntimeSessionMutation() error = nil, want unsupported actor kind rejection")
	}

	mutation = validCreateRuntimeSessionMutation()
	mutation.ActorID = "player-2"
	if _, err := NormalizeCreateRuntimeSessionMutation(mutation); err == nil {
		t.Fatal("NormalizeCreateRuntimeSessionMutation() error = nil, want player actor mismatch rejection")
	}
}

func TestNormalizeCreateRuntimeSessionMutationNormalizesTimesToUTC(t *testing.T) {
	issuedAt := time.Date(2026, 5, 17, 18, 30, 0, 0, time.FixedZone("test", 8*60*60))
	mutation := validCreateRuntimeSessionMutation()
	mutation.IssuedAt = issuedAt
	mutation.ExpiresAt = issuedAt.Add(2 * time.Hour)
	mutation.LastSeenAt = issuedAt.Add(time.Minute)
	normalized, err := NormalizeCreateRuntimeSessionMutation(mutation)
	if err != nil {
		t.Fatalf("NormalizeCreateRuntimeSessionMutation() error = %v, want nil", err)
	}
	if normalized.IssuedAt.Location() != time.UTC ||
		normalized.ExpiresAt.Location() != time.UTC ||
		normalized.LastSeenAt.Location() != time.UTC {
		t.Fatalf("normalized times = %#v, want UTC", normalized)
	}
}

func TestNormalizeCreateRuntimeSessionMutationRejectsInvalidLifetime(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*CreateRuntimeSessionMutation)
	}{
		{name: "expires_at_not_after_issued_at", mutate: func(m *CreateRuntimeSessionMutation) { m.ExpiresAt = m.IssuedAt }},
		{name: "last_seen_before_issued_at", mutate: func(m *CreateRuntimeSessionMutation) { m.LastSeenAt = m.IssuedAt.Add(-time.Second) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutation := validCreateRuntimeSessionMutation()
			tt.mutate(&mutation)
			if _, err := NormalizeCreateRuntimeSessionMutation(mutation); err == nil {
				t.Fatal("NormalizeCreateRuntimeSessionMutation() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeRuntimeSessionRecord(t *testing.T) {
	createdAt := time.Date(2026, 5, 17, 1, 0, 0, 0, time.FixedZone("test", 8*60*60))
	revokedAt := createdAt.Add(time.Minute)
	record, err := NormalizeRuntimeSessionRecord(RuntimeSession{
		SessionID:           " session-1 ",
		ActorKind:           ActorKind(" player "),
		ActorID:             " player-1 ",
		PlayerID:            " player-1 ",
		SessionStatus:       SessionStatusRevoked,
		IssuedAt:            createdAt,
		ExpiresAt:           createdAt.Add(time.Hour),
		LastSeenAt:          createdAt.Add(time.Minute),
		RevokedAt:           &revokedAt,
		RevocationReason:    " logout ",
		AccessTokenRecordID: " token-1 ",
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt.Add(2 * time.Minute),
	})
	if err != nil {
		t.Fatalf("NormalizeRuntimeSessionRecord() error = %v, want nil", err)
	}
	if record.SessionID != "session-1" ||
		record.ActorKind != ActorKindPlayer ||
		record.ActorID != "player-1" ||
		record.PlayerID != "player-1" ||
		record.RevocationReason != "logout" ||
		record.AccessTokenRecordID != "token-1" {
		t.Fatalf("normalized record = %#v, want trimmed fields", record)
	}
	if record.RevokedAt == nil || record.RevokedAt.Location() != time.UTC ||
		record.IssuedAt.Location() != time.UTC ||
		record.CreatedAt.Location() != time.UTC ||
		record.UpdatedAt.Location() != time.UTC {
		t.Fatalf("normalized record times = %#v, want UTC", record)
	}
}

func TestNormalizeRuntimeSessionRecordRejectsInvalidShape(t *testing.T) {
	valid := validRuntimeSessionRecord()
	tests := []struct {
		name   string
		mutate func(*RuntimeSession)
	}{
		{name: "invalid_status", mutate: func(r *RuntimeSession) { r.SessionStatus = "disabled" }},
		{name: "actor_mismatch", mutate: func(r *RuntimeSession) { r.ActorID = "player-2" }},
		{name: "revoked_status_without_revoked_at", mutate: func(r *RuntimeSession) { r.SessionStatus = SessionStatusRevoked; r.RevokedAt = nil }},
		{name: "revoked_at_without_revoked_status", mutate: func(r *RuntimeSession) { at := r.IssuedAt; r.SessionStatus = SessionStatusActive; r.RevokedAt = &at }},
		{name: "revocation_reason_without_revoked_at", mutate: func(r *RuntimeSession) { r.RevocationReason = "logout" }},
		{name: "updated_before_created", mutate: func(r *RuntimeSession) { r.UpdatedAt = r.CreatedAt.Add(-time.Second) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := valid
			tt.mutate(&record)
			if _, err := NormalizeRuntimeSessionRecord(record); err == nil {
				t.Fatal("NormalizeRuntimeSessionRecord() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeGetAndFindRuntimeSessionQueries(t *testing.T) {
	getQuery, err := NormalizeGetRuntimeSessionQuery(GetRuntimeSessionQuery{SessionID: " session-1 "})
	if err != nil {
		t.Fatalf("NormalizeGetRuntimeSessionQuery() error = %v, want nil", err)
	}
	if getQuery.SessionID != "session-1" {
		t.Fatalf("get query = %#v, want trimmed session id", getQuery)
	}

	observedAt := time.Date(2026, 5, 17, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	findQuery, err := NormalizeFindActiveSessionByIDQuery(FindActiveSessionByIDQuery{
		SessionID:  " session-1 ",
		ObservedAt: observedAt,
	})
	if err != nil {
		t.Fatalf("NormalizeFindActiveSessionByIDQuery() error = %v, want nil", err)
	}
	if findQuery.SessionID != "session-1" || findQuery.ObservedAt.Location() != time.UTC {
		t.Fatalf("find query = %#v, want trimmed session id and UTC observed_at", findQuery)
	}

	if _, err := NormalizeGetRuntimeSessionQuery(GetRuntimeSessionQuery{SessionID: " "}); err == nil {
		t.Fatal("NormalizeGetRuntimeSessionQuery() error = nil, want rejection")
	}
	if _, err := NormalizeFindActiveSessionByIDQuery(FindActiveSessionByIDQuery{SessionID: "session-1"}); err == nil {
		t.Fatal("NormalizeFindActiveSessionByIDQuery() error = nil, want missing observed_at rejection")
	}
}

func TestNormalizeUpdateRuntimeSessionLastSeenMutation(t *testing.T) {
	lastSeenAt := time.Date(2026, 5, 17, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	mutation, err := NormalizeUpdateRuntimeSessionLastSeenMutation(UpdateRuntimeSessionLastSeenMutation{
		SessionID:   " session-1 ",
		LastSeenAt:  lastSeenAt,
		RequestedBy: " session_validation ",
	})
	if err != nil {
		t.Fatalf("NormalizeUpdateRuntimeSessionLastSeenMutation() error = %v, want nil", err)
	}
	if mutation.SessionID != "session-1" ||
		mutation.RequestedBy != "session_validation" ||
		mutation.LastSeenAt.Location() != time.UTC {
		t.Fatalf("normalized mutation = %#v, want trimmed UTC fields", mutation)
	}
	if _, err := NormalizeUpdateRuntimeSessionLastSeenMutation(UpdateRuntimeSessionLastSeenMutation{
		SessionID:   "session-1",
		RequestedBy: "session_validation",
	}); err == nil {
		t.Fatal("NormalizeUpdateRuntimeSessionLastSeenMutation() error = nil, want missing last_seen_at rejection")
	}
}

func TestNormalizeMarkRuntimeSessionExpiredMutation(t *testing.T) {
	expiredAt := time.Date(2026, 5, 17, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	mutation, err := NormalizeMarkRuntimeSessionExpiredMutation(MarkRuntimeSessionExpiredMutation{
		SessionID:   " session-1 ",
		ExpiredAt:   expiredAt,
		RequestedBy: " cleanup ",
	})
	if err != nil {
		t.Fatalf("NormalizeMarkRuntimeSessionExpiredMutation() error = %v, want nil", err)
	}
	if mutation.SessionID != "session-1" ||
		mutation.RequestedBy != "cleanup" ||
		mutation.ExpiredAt.Location() != time.UTC {
		t.Fatalf("normalized mutation = %#v, want trimmed UTC fields", mutation)
	}
	if _, err := NormalizeMarkRuntimeSessionExpiredMutation(MarkRuntimeSessionExpiredMutation{
		SessionID:   "session-1",
		RequestedBy: "cleanup",
	}); err == nil {
		t.Fatal("NormalizeMarkRuntimeSessionExpiredMutation() error = nil, want missing expired_at rejection")
	}
}

func TestNormalizeRevokeRuntimeSessionMutation(t *testing.T) {
	revokedAt := time.Date(2026, 5, 17, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	mutation, err := NormalizeRevokeRuntimeSessionMutation(RevokeRuntimeSessionMutation{
		SessionID:        " session-1 ",
		RevokedAt:        revokedAt,
		RevocationReason: " logout ",
		RequestedBy:      " authentication_service ",
	})
	if err != nil {
		t.Fatalf("NormalizeRevokeRuntimeSessionMutation() error = %v, want nil", err)
	}
	if mutation.SessionID != "session-1" ||
		mutation.RevocationReason != "logout" ||
		mutation.RequestedBy != "authentication_service" ||
		mutation.RevokedAt.Location() != time.UTC {
		t.Fatalf("normalized mutation = %#v, want trimmed UTC fields", mutation)
	}
	if _, err := NormalizeRevokeRuntimeSessionMutation(RevokeRuntimeSessionMutation{
		SessionID:   "session-1",
		RevokedAt:   revokedAt,
		RequestedBy: "authentication_service",
	}); err == nil {
		t.Fatal("NormalizeRevokeRuntimeSessionMutation() error = nil, want missing reason rejection")
	}
}

func TestNormalizeListActiveSessionsForPlayerQuery(t *testing.T) {
	observedAt := time.Date(2026, 5, 17, 1, 2, 3, 0, time.FixedZone("test", 8*60*60))
	query, err := NormalizeListActiveSessionsForPlayerQuery(ListActiveSessionsForPlayerQuery{
		PlayerID:   " player-1 ",
		ObservedAt: observedAt,
		Limit:      50,
	})
	if err != nil {
		t.Fatalf("NormalizeListActiveSessionsForPlayerQuery() error = %v, want nil", err)
	}
	if query.PlayerID != "player-1" ||
		query.ObservedAt.Location() != time.UTC ||
		query.Limit != 50 {
		t.Fatalf("normalized query = %#v, want trimmed player id, UTC observed_at, and preserved limit", query)
	}

	tests := []struct {
		name   string
		query  ListActiveSessionsForPlayerQuery
		wantOK bool
	}{
		{name: "missing_player", query: ListActiveSessionsForPlayerQuery{ObservedAt: observedAt, Limit: 1}},
		{name: "missing_observed_at", query: ListActiveSessionsForPlayerQuery{PlayerID: "player-1", Limit: 1}},
		{name: "zero_limit", query: ListActiveSessionsForPlayerQuery{PlayerID: "player-1", ObservedAt: observedAt}},
		{name: "above_max_limit", query: ListActiveSessionsForPlayerQuery{PlayerID: "player-1", ObservedAt: observedAt, Limit: MaxListActiveSessionsPageSize + 1}},
		{name: "max_limit", query: ListActiveSessionsForPlayerQuery{PlayerID: "player-1", ObservedAt: observedAt, Limit: MaxListActiveSessionsPageSize}, wantOK: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NormalizeListActiveSessionsForPlayerQuery(tt.query)
			if tt.wantOK && err != nil {
				t.Fatalf("NormalizeListActiveSessionsForPlayerQuery() error = %v, want nil", err)
			}
			if !tt.wantOK && err == nil {
				t.Fatal("NormalizeListActiveSessionsForPlayerQuery() error = nil, want rejection")
			}
		})
	}
}

func TestRuntimeSessionHasNoSecretMaterialFields(t *testing.T) {
	forbidden := []string{
		"AccessToken",
		"RawAccessToken",
		"RawCredential",
		"CredentialProof",
		"TokenLookupDigest",
		"TokenVerifierDigest",
		"CredentialLookupDigest",
		"CredentialVerifierDigest",
		"VerifierKey",
		"VerifierDigest",
		"WebSocket",
		"ConnectionState",
	}
	types := []any{
		RuntimeSession{},
		CreateRuntimeSessionMutation{},
		GetRuntimeSessionQuery{},
		FindActiveSessionByIDQuery{},
		UpdateRuntimeSessionLastSeenMutation{},
		MarkRuntimeSessionExpiredMutation{},
		RevokeRuntimeSessionMutation{},
		ListActiveSessionsForPlayerQuery{},
	}
	for _, value := range types {
		typ := reflect.TypeOf(value)
		for i := 0; i < typ.NumField(); i++ {
			fieldName := typ.Field(i).Name
			for _, forbiddenPart := range forbidden {
				if strings.Contains(fieldName, forbiddenPart) && fieldName != "AccessTokenRecordID" {
					t.Fatalf("%s.%s exposes forbidden session material marker %q", typ.Name(), fieldName, forbiddenPart)
				}
			}
		}
	}
}

func validCreateRuntimeSessionMutation() CreateRuntimeSessionMutation {
	issuedAt := time.Date(2026, 5, 17, 1, 2, 3, 0, time.UTC)
	return CreateRuntimeSessionMutation{
		SessionID:           "session-1",
		ActorKind:           ActorKindPlayer,
		ActorID:             "player-1",
		PlayerID:            "player-1",
		SessionStatus:       SessionStatusActive,
		IssuedAt:            issuedAt,
		ExpiresAt:           issuedAt.Add(time.Hour),
		LastSeenAt:          issuedAt,
		AccessTokenRecordID: "token-1",
		RequestedBy:         "authentication_service",
	}
}

func validRuntimeSessionRecord() RuntimeSession {
	createdAt := time.Date(2026, 5, 17, 1, 2, 3, 0, time.UTC)
	return RuntimeSession{
		SessionID:           "session-1",
		ActorKind:           ActorKindPlayer,
		ActorID:             "player-1",
		PlayerID:            "player-1",
		SessionStatus:       SessionStatusActive,
		IssuedAt:            createdAt,
		ExpiresAt:           createdAt.Add(time.Hour),
		LastSeenAt:          createdAt,
		AccessTokenRecordID: "token-1",
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt,
	}
}

type recordingRepository struct{}

func (recordingRepository) CreateRuntimeSession(context.Context, CreateRuntimeSessionMutation) (RuntimeSession, error) {
	return RuntimeSession{}, nil
}

func (recordingRepository) GetRuntimeSession(context.Context, GetRuntimeSessionQuery) (RuntimeSession, error) {
	return RuntimeSession{}, nil
}

func (recordingRepository) FindActiveSessionByID(context.Context, FindActiveSessionByIDQuery) (RuntimeSession, error) {
	return RuntimeSession{}, nil
}

func (recordingRepository) UpdateRuntimeSessionLastSeen(context.Context, UpdateRuntimeSessionLastSeenMutation) (RuntimeSession, error) {
	return RuntimeSession{}, nil
}

func (recordingRepository) MarkRuntimeSessionExpired(context.Context, MarkRuntimeSessionExpiredMutation) (RuntimeSession, error) {
	return RuntimeSession{}, nil
}

func (recordingRepository) RevokeRuntimeSession(context.Context, RevokeRuntimeSessionMutation) (RuntimeSession, error) {
	return RuntimeSession{}, nil
}

func (recordingRepository) ListActiveSessionsForPlayer(context.Context, ListActiveSessionsForPlayerQuery) ([]RuntimeSession, error) {
	return nil, nil
}
