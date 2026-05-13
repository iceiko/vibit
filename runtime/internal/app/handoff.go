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

type RouteRequest struct {
	RequestID    string
	Route        RouteKey
	Target       Target
	Session      Session
	PayloadType  string
	Payload      any
	PayloadBytes []byte
}

func RenderRouteKey(route RouteKey) string {
	module := strings.TrimSpace(route.Module)
	name := strings.TrimSpace(route.Name)
	if module == "" || name == "" {
		return ""
	}
	return fmt.Sprintf("%s.%s", module, name)
}
