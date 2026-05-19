package app

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app/session"
)

func TestPersistentSessionValidatorValidatesActivePersistedPlayerSession(t *testing.T) {
	observedAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	repository := &runtimeSessionValidationRepository{
		record: validRuntimeSessionValidationRecord(observedAt),
	}
	validator := PersistentSessionValidator{
		Repository: repository,
		Clock:      fixedRuntimeSessionValidationClock{value: observedAt},
	}

	result, err := validator.ValidateSession(context.Background(), RouteRequest{
		Session: Session{
			SessionID:       " session-1 ",
			ConnectionID:    " request-connection ",
			ConnectionEpoch: 3,
		},
		Identity: RequestIdentity{
			Status:            IdentityValidationValidated,
			ActorKind:         ActorKindPlayer,
			ActorID:           "player-1",
			PlayerID:          "player-1",
			SessionID:         " session-1 ",
			ConnectionID:      " bound-connection ",
			ConnectionEpoch:   7,
			PlayerIDValidated: true,
		},
	})
	if err != nil {
		t.Fatalf("ValidateSession() error = %v, want nil", err)
	}
	if !result.Valid || result.Reason != runtimeSessionValidationValidReason {
		t.Fatalf("result = %#v, want valid session result", result)
	}
	if repository.findCalls != 1 {
		t.Fatalf("FindActiveSessionByID calls = %d, want 1", repository.findCalls)
	}
	if repository.findQuery.SessionID != "session-1" || !repository.findQuery.ObservedAt.Equal(observedAt) {
		t.Fatalf("find query = %#v, want normalized session id and fixed observed_at", repository.findQuery)
	}
	if result.Identity.Status != IdentityValidationValidated ||
		result.Identity.ActorKind != ActorKindPlayer ||
		result.Identity.ActorID != "player-1" ||
		result.Identity.PlayerID != "player-1" ||
		result.Identity.SessionID != "session-1" ||
		result.Identity.ConnectionID != "bound-connection" ||
		result.Identity.ConnectionEpoch != 7 ||
		!result.Identity.PlayerIDValidated ||
		!result.Identity.SessionValidated {
		t.Fatalf("Identity = %#v, want session-validated player identity preserving server-bound connection", result.Identity)
	}
	if repository.mutated {
		t.Fatal("repository mutation was called, want validation to avoid last_seen/create/revoke mutations")
	}
}

func TestPersistentSessionValidatorRejectsMissingMalformedOrConflictingSessionIDBeforeLookup(t *testing.T) {
	tests := []struct {
		name    string
		request RouteRequest
	}{
		{
			name: "missing session id",
			request: RouteRequest{
				Identity: validatedRuntimeSessionIdentity("player-1", ""),
			},
		},
		{
			name: "malformed whitespace session id",
			request: RouteRequest{
				Session:  Session{SessionID: " "},
				Identity: validatedRuntimeSessionIdentity("player-1", ""),
			},
		},
		{
			name: "conflicting session id",
			request: RouteRequest{
				Session:  Session{SessionID: "session-from-envelope"},
				Identity: validatedRuntimeSessionIdentity("player-1", "session-from-identity"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &runtimeSessionValidationRepository{}
			validator := NewPersistentSessionValidator(repository)

			result, err := validator.ValidateSession(context.Background(), tt.request)
			if err != nil {
				t.Fatalf("ValidateSession() error = %v, want nil", err)
			}
			if result.Valid || result.Reason != runtimeSessionValidationInvalidReason {
				t.Fatalf("result = %#v, want collapsed invalid session", result)
			}
			if result.Reason != "session validation failed" {
				t.Fatalf("Reason = %q, want stable public invalid-session reason", result.Reason)
			}
			if repository.findCalls != 0 {
				t.Fatalf("FindActiveSessionByID calls = %d, want 0", repository.findCalls)
			}
		})
	}
}

func TestPersistentSessionValidatorRejectsMetadataOnlyOrUnvalidatedIdentityBeforeLookup(t *testing.T) {
	tests := []struct {
		name     string
		identity RequestIdentity
	}{
		{name: "metadata only", identity: MetadataOnlyIdentityFromSession(Session{SessionID: "session-1", PlayerID: "player-1"})},
		{name: "unknown", identity: RequestIdentity{Status: IdentityValidationUnknown, ActorKind: ActorKindPlayer, PlayerID: "player-1", ActorID: "player-1", PlayerIDValidated: true}},
		{name: "missing player validation", identity: RequestIdentity{Status: IdentityValidationValidated, ActorKind: ActorKindPlayer, PlayerID: "player-1", ActorID: "player-1"}},
		{name: "actor mismatch", identity: RequestIdentity{Status: IdentityValidationValidated, ActorKind: ActorKindPlayer, PlayerID: "player-1", ActorID: "player-2", PlayerIDValidated: true}},
		{name: "service actor", identity: RequestIdentity{Status: IdentityValidationValidated, ActorKind: ActorKindService, PlayerID: "player-1", ActorID: "player-1", PlayerIDValidated: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &runtimeSessionValidationRepository{}
			validator := NewPersistentSessionValidator(repository)

			result, err := validator.ValidateSession(context.Background(), RouteRequest{
				Session:  Session{SessionID: "session-1"},
				Identity: tt.identity,
			})
			if err != nil {
				t.Fatalf("ValidateSession() error = %v, want nil", err)
			}
			if result.Valid || result.Identity.SessionValidated {
				t.Fatalf("result = %#v, want invalid unvalidated session", result)
			}
			if repository.findCalls != 0 {
				t.Fatalf("FindActiveSessionByID calls = %d, want 0", repository.findCalls)
			}
		})
	}
}

func TestPersistentSessionValidatorCollapsesLookupAndRepositoryFailures(t *testing.T) {
	observedAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		record session.RuntimeSession
		err    error
	}{
		{name: "not found", err: errors.New("session session-1 not found")},
		{name: "expired row unavailable", err: errors.New("session session-1 is expired")},
		{name: "revoked row unavailable", err: errors.New("session session-1 is revoked")},
		{name: "repository unavailable", err: errors.New("database contains secret session-1 unavailable")},
		{name: "malformed record", record: session.RuntimeSession{SessionID: "session-1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &runtimeSessionValidationRepository{record: tt.record, err: tt.err}
			validator := PersistentSessionValidator{
				Repository: repository,
				Clock:      fixedRuntimeSessionValidationClock{value: observedAt},
			}

			result, err := validator.ValidateSession(context.Background(), RouteRequest{
				Session:  Session{SessionID: "session-1"},
				Identity: validatedRuntimeSessionIdentity("player-1", "session-1"),
			})
			if err != nil {
				t.Fatalf("ValidateSession() error = %v, want nil redacted collapse", err)
			}
			if result.Valid || result.Identity.SessionValidated || result.Reason != runtimeSessionValidationInvalidReason {
				t.Fatalf("result = %#v, want collapsed invalid session", result)
			}
			if strings.Contains(result.Reason, "session-1") || strings.Contains(result.Reason, "database") {
				t.Fatalf("Reason = %q, want redacted public reason", result.Reason)
			}
			if repository.mutated {
				t.Fatal("repository mutation was called, want lookup-only validation")
			}
		})
	}
}

func TestPersistentSessionValidatorCollapsesRecordIdentityMismatches(t *testing.T) {
	observedAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		mutate func(*session.RuntimeSession)
	}{
		{name: "session id mismatch", mutate: func(record *session.RuntimeSession) { record.SessionID = "session-2" }},
		{name: "actor id mismatch", mutate: func(record *session.RuntimeSession) { record.ActorID = "player-2"; record.PlayerID = "player-2" }},
		{name: "player id mismatch", mutate: func(record *session.RuntimeSession) { record.ActorID = "player-2"; record.PlayerID = "player-2" }},
		{name: "inactive status", mutate: func(record *session.RuntimeSession) {
			revokedAt := observedAt.Add(-time.Minute)
			record.SessionStatus = session.SessionStatusRevoked
			record.RevokedAt = &revokedAt
			record.RevocationReason = "logout"
		}},
		{name: "expired at observed time", mutate: func(record *session.RuntimeSession) { record.ExpiresAt = observedAt }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := validRuntimeSessionValidationRecord(observedAt)
			tt.mutate(&record)
			repository := &runtimeSessionValidationRepository{record: record}
			validator := PersistentSessionValidator{
				Repository: repository,
				Clock:      fixedRuntimeSessionValidationClock{value: observedAt},
			}

			result, err := validator.ValidateSession(context.Background(), RouteRequest{
				Session:  Session{SessionID: "session-1"},
				Identity: validatedRuntimeSessionIdentity("player-1", "session-1"),
			})
			if err != nil {
				t.Fatalf("ValidateSession() error = %v, want nil", err)
			}
			if result.Valid || result.Identity.SessionValidated || result.Reason != runtimeSessionValidationInvalidReason {
				t.Fatalf("result = %#v, want collapsed invalid mismatch", result)
			}
		})
	}
}

func TestPersistentSessionValidatorRejectsNilRepository(t *testing.T) {
	validator := NewPersistentSessionValidator(nil)

	result, err := validator.ValidateSession(context.Background(), RouteRequest{
		Session:  Session{SessionID: "session-1"},
		Identity: validatedRuntimeSessionIdentity("player-1", "session-1"),
	})
	if err != nil {
		t.Fatalf("ValidateSession() error = %v, want nil", err)
	}
	if result.Valid || result.Identity.SessionValidated || result.Reason != runtimeSessionValidationInvalidReason {
		t.Fatalf("result = %#v, want collapsed invalid session", result)
	}
}

type fixedRuntimeSessionValidationClock struct {
	value time.Time
}

func (c fixedRuntimeSessionValidationClock) Now() time.Time {
	return c.value
}

type runtimeSessionValidationRepository struct {
	record    session.RuntimeSession
	err       error
	findCalls int
	findQuery session.FindActiveSessionByIDQuery
	mutated   bool
}

func (r *runtimeSessionValidationRepository) CreateRuntimeSession(context.Context, session.CreateRuntimeSessionMutation) (session.RuntimeSession, error) {
	r.mutated = true
	return session.RuntimeSession{}, errors.New("unexpected CreateRuntimeSession")
}

func (r *runtimeSessionValidationRepository) GetRuntimeSession(context.Context, session.GetRuntimeSessionQuery) (session.RuntimeSession, error) {
	return session.RuntimeSession{}, errors.New("unexpected GetRuntimeSession")
}

func (r *runtimeSessionValidationRepository) FindActiveSessionByID(_ context.Context, query session.FindActiveSessionByIDQuery) (session.RuntimeSession, error) {
	r.findCalls++
	r.findQuery = query
	if r.err != nil {
		return session.RuntimeSession{}, r.err
	}
	return r.record, nil
}

func (r *runtimeSessionValidationRepository) UpdateRuntimeSessionLastSeen(context.Context, session.UpdateRuntimeSessionLastSeenMutation) (session.RuntimeSession, error) {
	r.mutated = true
	return session.RuntimeSession{}, errors.New("unexpected UpdateRuntimeSessionLastSeen")
}

func (r *runtimeSessionValidationRepository) MarkRuntimeSessionExpired(context.Context, session.MarkRuntimeSessionExpiredMutation) (session.RuntimeSession, error) {
	r.mutated = true
	return session.RuntimeSession{}, errors.New("unexpected MarkRuntimeSessionExpired")
}

func (r *runtimeSessionValidationRepository) RevokeRuntimeSession(context.Context, session.RevokeRuntimeSessionMutation) (session.RuntimeSession, error) {
	r.mutated = true
	return session.RuntimeSession{}, errors.New("unexpected RevokeRuntimeSession")
}

func (r *runtimeSessionValidationRepository) ListActiveSessionsForPlayer(context.Context, session.ListActiveSessionsForPlayerQuery) ([]session.RuntimeSession, error) {
	return nil, errors.New("unexpected ListActiveSessionsForPlayer")
}

func validRuntimeSessionValidationRecord(observedAt time.Time) session.RuntimeSession {
	issuedAt := observedAt.Add(-time.Hour)
	return session.RuntimeSession{
		SessionID:           "session-1",
		ActorKind:           session.ActorKindPlayer,
		ActorID:             "player-1",
		PlayerID:            "player-1",
		SessionStatus:       session.SessionStatusActive,
		IssuedAt:            issuedAt,
		ExpiresAt:           observedAt.Add(time.Hour),
		LastSeenAt:          issuedAt.Add(time.Minute),
		AccessTokenRecordID: "token-record-1",
		CreatedAt:           issuedAt,
		UpdatedAt:           issuedAt.Add(time.Minute),
	}
}

func validatedRuntimeSessionIdentity(playerID string, sessionID string) RequestIdentity {
	return RequestIdentity{
		Status:            IdentityValidationValidated,
		ActorKind:         ActorKindPlayer,
		ActorID:           playerID,
		PlayerID:          playerID,
		SessionID:         sessionID,
		PlayerIDValidated: true,
	}
}
