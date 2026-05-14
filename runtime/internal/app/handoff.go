package app

import (
	"fmt"
	"strings"
)

type MessageKind string

const (
	MessageKindCommand   MessageKind = "command"
	MessageKindQuery     MessageKind = "query"
	MessageKindEvent     MessageKind = "event"
	MessageKindError     MessageKind = "error"
	MessageKindSystem    MessageKind = "system"
	MessageKindAck       MessageKind = "ack"
	MessageKindHeartbeat MessageKind = "heartbeat"
	MessageKindInput     MessageKind = "input"
	MessageKindState     MessageKind = "state"
)

type TargetScope string

const (
	TargetScopeGlobal TargetScope = "global"
	TargetScopePlayer TargetScope = "player"
	TargetScopeParty  TargetScope = "party"
	TargetScopeRoom   TargetScope = "room"
	TargetScopeMatch  TargetScope = "match"
	TargetScopeStream TargetScope = "stream"
	TargetScopeSystem TargetScope = "system"
)

type RouteKey struct {
	Kind   MessageKind
	Module string
	Name   string
}

type Target struct {
	Scope TargetScope
	ID    string
}

type Session struct {
	ConnectionID    string
	SessionID       string
	PlayerID        string
	ConnectionEpoch uint64
}

type IdentityValidationStatus string

const (
	IdentityValidationUnknown      IdentityValidationStatus = "unknown"
	IdentityValidationMetadataOnly IdentityValidationStatus = "metadata_only"
	IdentityValidationValidated    IdentityValidationStatus = "validated"
)

type ActorKind string

const (
	ActorKindUnknown ActorKind = "unknown"
	ActorKindPlayer  ActorKind = "player"
	ActorKindService ActorKind = "service"
	ActorKindAdmin   ActorKind = "admin"
)

type RequestIdentity struct {
	Status            IdentityValidationStatus
	ActorKind         ActorKind
	ActorID           string
	PlayerID          string
	SessionID         string
	ConnectionID      string
	ConnectionEpoch   uint64
	SessionValidated  bool
	PlayerIDValidated bool
}

type SessionValidationResult struct {
	Identity RequestIdentity
	Valid    bool
	Reason   string
}

type RouteRequest struct {
	RequestID    string
	Route        RouteKey
	Target       Target
	Session      Session
	Identity     RequestIdentity
	PayloadType  string
	Payload      any
	PayloadBytes []byte
}

func MetadataOnlyIdentityFromSession(session Session) RequestIdentity {
	playerID := strings.TrimSpace(session.PlayerID)
	sessionID := strings.TrimSpace(session.SessionID)
	connectionID := strings.TrimSpace(session.ConnectionID)

	actorKind := ActorKindUnknown
	actorID := ""
	if playerID != "" {
		actorKind = ActorKindPlayer
		actorID = playerID
	}

	return RequestIdentity{
		Status:          IdentityValidationMetadataOnly,
		ActorKind:       actorKind,
		ActorID:         actorID,
		PlayerID:        playerID,
		SessionID:       sessionID,
		ConnectionID:    connectionID,
		ConnectionEpoch: session.ConnectionEpoch,
	}
}

func ValidatedPlayerIdentity(playerID string, session Session) RequestIdentity {
	normalizedPlayerID := strings.TrimSpace(playerID)
	if normalizedPlayerID == "" {
		normalizedPlayerID = strings.TrimSpace(session.PlayerID)
	}

	identity := MetadataOnlyIdentityFromSession(Session{
		ConnectionID:    session.ConnectionID,
		SessionID:       session.SessionID,
		PlayerID:        normalizedPlayerID,
		ConnectionEpoch: session.ConnectionEpoch,
	})
	identity.Status = IdentityValidationValidated
	identity.ActorKind = ActorKindPlayer
	identity.ActorID = normalizedPlayerID
	identity.PlayerIDValidated = normalizedPlayerID != ""
	identity.SessionValidated = strings.TrimSpace(session.SessionID) != ""
	return identity
}

func RenderRouteKey(route RouteKey) string {
	module := strings.TrimSpace(route.Module)
	name := strings.TrimSpace(route.Name)
	if module == "" || name == "" {
		return ""
	}
	return fmt.Sprintf("%s.%s", module, name)
}
