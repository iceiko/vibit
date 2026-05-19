package session

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	ModuleName                    = "runtime.session"
	MaxListActiveSessionsPageSize = 500
)

type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "active"
	SessionStatusExpired SessionStatus = "expired"
	SessionStatusRevoked SessionStatus = "revoked"
)

func (s SessionStatus) IsValid() bool {
	switch s {
	case SessionStatusActive, SessionStatusExpired, SessionStatusRevoked:
		return true
	default:
		return false
	}
}

type ActorKind string

const (
	ActorKindPlayer ActorKind = "player"
)

func (k ActorKind) IsValid() bool {
	return k == ActorKindPlayer
}

type RuntimeSession struct {
	SessionID           string
	ActorKind           ActorKind
	ActorID             string
	PlayerID            string
	SessionStatus       SessionStatus
	IssuedAt            time.Time
	ExpiresAt           time.Time
	LastSeenAt          time.Time
	RevokedAt           *time.Time
	RevocationReason    string
	AccessTokenRecordID string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CreateRuntimeSessionMutation struct {
	SessionID           string
	ActorKind           ActorKind
	ActorID             string
	PlayerID            string
	SessionStatus       SessionStatus
	IssuedAt            time.Time
	ExpiresAt           time.Time
	LastSeenAt          time.Time
	AccessTokenRecordID string
	RequestedBy         string
}

type GetRuntimeSessionQuery struct {
	SessionID string
}

type FindActiveSessionByIDQuery struct {
	SessionID  string
	ObservedAt time.Time
}

type UpdateRuntimeSessionLastSeenMutation struct {
	SessionID   string
	LastSeenAt  time.Time
	RequestedBy string
}

type MarkRuntimeSessionExpiredMutation struct {
	SessionID   string
	ExpiredAt   time.Time
	RequestedBy string
}

type RevokeRuntimeSessionMutation struct {
	SessionID        string
	RevokedAt        time.Time
	RevocationReason string
	RequestedBy      string
}

type ListActiveSessionsForPlayerQuery struct {
	PlayerID   string
	ObservedAt time.Time
	Limit      int
}

type Repository interface {
	CreateRuntimeSession(context.Context, CreateRuntimeSessionMutation) (RuntimeSession, error)
	GetRuntimeSession(context.Context, GetRuntimeSessionQuery) (RuntimeSession, error)
	FindActiveSessionByID(context.Context, FindActiveSessionByIDQuery) (RuntimeSession, error)
	UpdateRuntimeSessionLastSeen(context.Context, UpdateRuntimeSessionLastSeenMutation) (RuntimeSession, error)
	MarkRuntimeSessionExpired(context.Context, MarkRuntimeSessionExpiredMutation) (RuntimeSession, error)
	RevokeRuntimeSession(context.Context, RevokeRuntimeSessionMutation) (RuntimeSession, error)
	ListActiveSessionsForPlayer(context.Context, ListActiveSessionsForPlayerQuery) ([]RuntimeSession, error)
}

func NormalizeRuntimeSessionRecord(record RuntimeSession) (RuntimeSession, error) {
	var err error
	record.SessionID, err = normalizeRequired("session_id", record.SessionID)
	if err != nil {
		return RuntimeSession{}, err
	}
	record.ActorKind, err = normalizeActorKind(record.ActorKind)
	if err != nil {
		return RuntimeSession{}, err
	}
	record.ActorID, err = normalizeRequired("actor_id", record.ActorID)
	if err != nil {
		return RuntimeSession{}, err
	}
	record.PlayerID, err = normalizeRequired("player_id", record.PlayerID)
	if err != nil {
		return RuntimeSession{}, err
	}
	if record.ActorKind == ActorKindPlayer && record.ActorID != record.PlayerID {
		return RuntimeSession{}, errors.New("session: actor_id must match player_id for player sessions")
	}
	if !record.SessionStatus.IsValid() {
		return RuntimeSession{}, errors.New("session: session_status is invalid")
	}
	record.IssuedAt, err = normalizeRequiredTime("issued_at", record.IssuedAt)
	if err != nil {
		return RuntimeSession{}, err
	}
	record.ExpiresAt, err = normalizeRequiredTime("expires_at", record.ExpiresAt)
	if err != nil {
		return RuntimeSession{}, err
	}
	record.LastSeenAt, err = normalizeRequiredTime("last_seen_at", record.LastSeenAt)
	if err != nil {
		return RuntimeSession{}, err
	}
	if !record.ExpiresAt.After(record.IssuedAt) {
		return RuntimeSession{}, errors.New("session: expires_at must be after issued_at")
	}
	if record.LastSeenAt.Before(record.IssuedAt) {
		return RuntimeSession{}, errors.New("session: last_seen_at must not be before issued_at")
	}
	if record.RevokedAt != nil {
		revokedAt := record.RevokedAt.UTC()
		if revokedAt.Before(record.IssuedAt) {
			return RuntimeSession{}, errors.New("session: revoked_at must not be before issued_at")
		}
		if record.SessionStatus != SessionStatusRevoked {
			return RuntimeSession{}, errors.New("session: revoked_at requires revoked status")
		}
		record.RevokedAt = &revokedAt
	}
	record.RevocationReason = strings.TrimSpace(record.RevocationReason)
	if record.RevocationReason != "" && record.RevokedAt == nil {
		return RuntimeSession{}, errors.New("session: revocation_reason requires revoked_at")
	}
	if record.SessionStatus == SessionStatusRevoked && record.RevokedAt == nil {
		return RuntimeSession{}, errors.New("session: revoked sessions require revoked_at")
	}
	record.AccessTokenRecordID = strings.TrimSpace(record.AccessTokenRecordID)
	record.CreatedAt, err = normalizeRequiredTime("created_at", record.CreatedAt)
	if err != nil {
		return RuntimeSession{}, err
	}
	record.UpdatedAt, err = normalizeRequiredTime("updated_at", record.UpdatedAt)
	if err != nil {
		return RuntimeSession{}, err
	}
	if record.UpdatedAt.Before(record.CreatedAt) {
		return RuntimeSession{}, errors.New("session: updated_at must not be before created_at")
	}
	return record, nil
}

func NormalizeCreateRuntimeSessionMutation(mutation CreateRuntimeSessionMutation) (CreateRuntimeSessionMutation, error) {
	var err error
	mutation.SessionID, err = normalizeRequired("session_id", mutation.SessionID)
	if err != nil {
		return CreateRuntimeSessionMutation{}, err
	}
	mutation.ActorKind, err = normalizeActorKind(mutation.ActorKind)
	if err != nil {
		return CreateRuntimeSessionMutation{}, err
	}
	mutation.ActorID, err = normalizeRequired("actor_id", mutation.ActorID)
	if err != nil {
		return CreateRuntimeSessionMutation{}, err
	}
	mutation.PlayerID, err = normalizeRequired("player_id", mutation.PlayerID)
	if err != nil {
		return CreateRuntimeSessionMutation{}, err
	}
	if mutation.ActorKind == ActorKindPlayer && mutation.ActorID != mutation.PlayerID {
		return CreateRuntimeSessionMutation{}, errors.New("session: actor_id must match player_id for player sessions")
	}
	if mutation.SessionStatus == "" {
		mutation.SessionStatus = SessionStatusActive
	}
	if mutation.SessionStatus != SessionStatusActive {
		return CreateRuntimeSessionMutation{}, errors.New("session: created sessions must start active")
	}
	mutation.IssuedAt, err = normalizeRequiredTime("issued_at", mutation.IssuedAt)
	if err != nil {
		return CreateRuntimeSessionMutation{}, err
	}
	mutation.ExpiresAt, err = normalizeRequiredTime("expires_at", mutation.ExpiresAt)
	if err != nil {
		return CreateRuntimeSessionMutation{}, err
	}
	mutation.LastSeenAt, err = normalizeRequiredTime("last_seen_at", mutation.LastSeenAt)
	if err != nil {
		return CreateRuntimeSessionMutation{}, err
	}
	if !mutation.ExpiresAt.After(mutation.IssuedAt) {
		return CreateRuntimeSessionMutation{}, errors.New("session: expires_at must be after issued_at")
	}
	if mutation.LastSeenAt.Before(mutation.IssuedAt) {
		return CreateRuntimeSessionMutation{}, errors.New("session: last_seen_at must not be before issued_at")
	}
	mutation.AccessTokenRecordID = strings.TrimSpace(mutation.AccessTokenRecordID)
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return CreateRuntimeSessionMutation{}, err
	}
	return mutation, nil
}

func NormalizeGetRuntimeSessionQuery(query GetRuntimeSessionQuery) (GetRuntimeSessionQuery, error) {
	sessionID, err := normalizeRequired("session_id", query.SessionID)
	if err != nil {
		return GetRuntimeSessionQuery{}, err
	}
	query.SessionID = sessionID
	return query, nil
}

func NormalizeFindActiveSessionByIDQuery(query FindActiveSessionByIDQuery) (FindActiveSessionByIDQuery, error) {
	var err error
	query.SessionID, err = normalizeRequired("session_id", query.SessionID)
	if err != nil {
		return FindActiveSessionByIDQuery{}, err
	}
	query.ObservedAt, err = normalizeRequiredTime("observed_at", query.ObservedAt)
	if err != nil {
		return FindActiveSessionByIDQuery{}, err
	}
	return query, nil
}

func NormalizeUpdateRuntimeSessionLastSeenMutation(mutation UpdateRuntimeSessionLastSeenMutation) (UpdateRuntimeSessionLastSeenMutation, error) {
	var err error
	mutation.SessionID, err = normalizeRequired("session_id", mutation.SessionID)
	if err != nil {
		return UpdateRuntimeSessionLastSeenMutation{}, err
	}
	mutation.LastSeenAt, err = normalizeRequiredTime("last_seen_at", mutation.LastSeenAt)
	if err != nil {
		return UpdateRuntimeSessionLastSeenMutation{}, err
	}
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return UpdateRuntimeSessionLastSeenMutation{}, err
	}
	return mutation, nil
}

func NormalizeMarkRuntimeSessionExpiredMutation(mutation MarkRuntimeSessionExpiredMutation) (MarkRuntimeSessionExpiredMutation, error) {
	var err error
	mutation.SessionID, err = normalizeRequired("session_id", mutation.SessionID)
	if err != nil {
		return MarkRuntimeSessionExpiredMutation{}, err
	}
	mutation.ExpiredAt, err = normalizeRequiredTime("expired_at", mutation.ExpiredAt)
	if err != nil {
		return MarkRuntimeSessionExpiredMutation{}, err
	}
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return MarkRuntimeSessionExpiredMutation{}, err
	}
	return mutation, nil
}

func NormalizeRevokeRuntimeSessionMutation(mutation RevokeRuntimeSessionMutation) (RevokeRuntimeSessionMutation, error) {
	var err error
	mutation.SessionID, err = normalizeRequired("session_id", mutation.SessionID)
	if err != nil {
		return RevokeRuntimeSessionMutation{}, err
	}
	mutation.RevokedAt, err = normalizeRequiredTime("revoked_at", mutation.RevokedAt)
	if err != nil {
		return RevokeRuntimeSessionMutation{}, err
	}
	mutation.RevocationReason, err = normalizeRequired("revocation_reason", mutation.RevocationReason)
	if err != nil {
		return RevokeRuntimeSessionMutation{}, err
	}
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return RevokeRuntimeSessionMutation{}, err
	}
	return mutation, nil
}

func NormalizeListActiveSessionsForPlayerQuery(query ListActiveSessionsForPlayerQuery) (ListActiveSessionsForPlayerQuery, error) {
	var err error
	query.PlayerID, err = normalizeRequired("player_id", query.PlayerID)
	if err != nil {
		return ListActiveSessionsForPlayerQuery{}, err
	}
	query.ObservedAt, err = normalizeRequiredTime("observed_at", query.ObservedAt)
	if err != nil {
		return ListActiveSessionsForPlayerQuery{}, err
	}
	if query.Limit <= 0 {
		return ListActiveSessionsForPlayerQuery{}, errors.New("session: limit must be positive")
	}
	if query.Limit > MaxListActiveSessionsPageSize {
		return ListActiveSessionsForPlayerQuery{}, fmt.Errorf("session: limit must be at most %d", MaxListActiveSessionsPageSize)
	}
	return query, nil
}

func normalizeRequired(name string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("session: %s is required", name)
	}
	return value, nil
}

func normalizeActorKind(value ActorKind) (ActorKind, error) {
	normalized := ActorKind(strings.TrimSpace(string(value)))
	if !normalized.IsValid() {
		return "", errors.New("session: actor_kind must be player")
	}
	return normalized, nil
}

func normalizeRequiredTime(name string, value time.Time) (time.Time, error) {
	if value.IsZero() {
		return time.Time{}, fmt.Errorf("session: %s is required", name)
	}
	return value.UTC(), nil
}
