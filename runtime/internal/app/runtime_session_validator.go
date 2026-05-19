package app

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app/session"
)

const (
	runtimeSessionValidationInvalidReason = "session validation failed"
	runtimeSessionValidationValidReason   = "session_validated"
)

type RuntimeSessionValidationClock interface {
	Now() time.Time
}

type PersistentSessionValidator struct {
	Repository session.Repository
	Clock      RuntimeSessionValidationClock
}

func NewPersistentSessionValidator(repository session.Repository) PersistentSessionValidator {
	return PersistentSessionValidator{Repository: repository}
}

func (v PersistentSessionValidator) ValidateSession(ctx context.Context, request RouteRequest) (SessionValidationResult, error) {
	repository := v.Repository
	if isNilSessionRepository(repository) {
		return invalidRuntimeSessionValidationResult(request), nil
	}

	identity := request.Identity
	if !identitySatisfiesRuntimeSessionValidation(identity) {
		return invalidRuntimeSessionValidationResult(request), nil
	}

	sessionID, ok := runtimeSessionIDFromRequest(request)
	if !ok {
		return invalidRuntimeSessionValidationResult(request), nil
	}

	observedAt := v.now()
	record, err := repository.FindActiveSessionByID(ctx, session.FindActiveSessionByIDQuery{
		SessionID:  sessionID,
		ObservedAt: observedAt,
	})
	if err != nil {
		return invalidRuntimeSessionValidationResult(request), nil
	}

	record, err = session.NormalizeRuntimeSessionRecord(record)
	if err != nil || !runtimeSessionRecordMatchesIdentity(record, sessionID, identity, observedAt) {
		return invalidRuntimeSessionValidationResult(request), nil
	}

	return SessionValidationResult{
		Identity: runtimeSessionValidatedIdentity(request, record),
		Valid:    true,
		Reason:   runtimeSessionValidationValidReason,
	}, nil
}

func identitySatisfiesRuntimeSessionValidation(identity RequestIdentity) bool {
	playerID := strings.TrimSpace(identity.PlayerID)
	actorID := strings.TrimSpace(identity.ActorID)
	return identity.Status == IdentityValidationValidated &&
		identity.ActorKind == ActorKindPlayer &&
		playerID != "" &&
		actorID == playerID &&
		identity.PlayerIDValidated
}

func runtimeSessionIDFromRequest(request RouteRequest) (string, bool) {
	identitySessionID := strings.TrimSpace(request.Identity.SessionID)
	requestSessionID := strings.TrimSpace(request.Session.SessionID)
	if identitySessionID != "" && requestSessionID != "" && identitySessionID != requestSessionID {
		return "", false
	}
	if identitySessionID != "" {
		return identitySessionID, true
	}
	if requestSessionID != "" {
		return requestSessionID, true
	}
	return "", false
}

func runtimeSessionRecordMatchesIdentity(record session.RuntimeSession, sessionID string, identity RequestIdentity, observedAt time.Time) bool {
	playerID := strings.TrimSpace(identity.PlayerID)
	actorID := strings.TrimSpace(identity.ActorID)
	return record.SessionID == sessionID &&
		record.ActorKind == session.ActorKindPlayer &&
		record.ActorID == actorID &&
		record.PlayerID == playerID &&
		record.SessionStatus == session.SessionStatusActive &&
		record.ExpiresAt.After(observedAt)
}

func runtimeSessionValidatedIdentity(request RouteRequest, record session.RuntimeSession) RequestIdentity {
	connectionID := strings.TrimSpace(request.Identity.ConnectionID)
	if connectionID == "" {
		connectionID = strings.TrimSpace(request.Session.ConnectionID)
	}

	connectionEpoch := request.Identity.ConnectionEpoch
	if connectionEpoch == 0 {
		connectionEpoch = request.Session.ConnectionEpoch
	}

	return RequestIdentity{
		Status:            IdentityValidationValidated,
		ActorKind:         ActorKindPlayer,
		ActorID:           record.ActorID,
		PlayerID:          record.PlayerID,
		SessionID:         record.SessionID,
		ConnectionID:      connectionID,
		ConnectionEpoch:   connectionEpoch,
		SessionValidated:  true,
		PlayerIDValidated: true,
	}
}

func invalidRuntimeSessionValidationResult(request RouteRequest) SessionValidationResult {
	identity := request.Identity
	if identity.Status == "" {
		identity = MetadataOnlyIdentityFromSession(request.Session)
	}
	identity.SessionValidated = false
	return SessionValidationResult{
		Identity: identity,
		Valid:    false,
		Reason:   runtimeSessionValidationInvalidReason,
	}
}

func (v PersistentSessionValidator) now() time.Time {
	if v.Clock == nil {
		return time.Now().UTC()
	}
	value := v.Clock.Now()
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}

func isNilSessionRepository(repository session.Repository) bool {
	if repository == nil {
		return true
	}
	reflected := reflect.ValueOf(repository)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}
