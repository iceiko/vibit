package presence

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/app/connection"
)

const (
	ModuleName = "runtime.presence"

	QueryGetPlayerPresence = "GetPlayerPresence"

	ErrorCodePresenceQueryForbidden   app.ErrorCode = "PRESENCE_QUERY_FORBIDDEN"
	ErrorCodePresenceQueryUnavailable app.ErrorCode = "PRESENCE_QUERY_UNAVAILABLE"
)

type GetPlayerPresenceRequest struct {
	PlayerID string
}

type GetPlayerPresenceResult struct {
	PlayerID          string
	Status            PresenceStatus
	ActiveConnections []PresenceConnection
	ConnectionCount   int
	LastSeenAt        *time.Time
	RuntimeSessionIDs []string
	ObservedAt        time.Time
}

type PresenceStatus string

const (
	PresenceStatusOffline PresenceStatus = "offline"
	PresenceStatusOnline  PresenceStatus = "online"
)

type PresenceConnection struct {
	ConnectionID     string
	ConnectionEpoch  uint64
	RuntimeSessionID string
	LastSeenAt       *time.Time
	OpenedAt         time.Time
	BoundAt          *time.Time
}

type Registry interface {
	PresenceForPlayer(context.Context, connection.PlayerID) connection.PlayerPresence
}

type Handlers struct {
	Registry Registry
}

func GetPlayerPresenceRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: ModuleName, Name: QueryGetPlayerPresence}
}

func (h Handlers) RegisterRoutes(dispatcher *app.Dispatcher) error {
	if dispatcher == nil {
		return errors.New("presence: dispatcher is nil")
	}
	return dispatcher.Register(GetPlayerPresenceRoute(), app.HandlerFunc(h.HandleGetPlayerPresenceRoute))
}

func (h Handlers) HandleGetPlayerPresenceRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	payload, ok := request.Payload.(GetPlayerPresenceRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*GetPlayerPresenceRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		return applicationErrorResult(request, ErrorCodePresenceQueryForbidden, "presence query payload is malformed")
	}

	result, err := h.GetPlayerPresence(ctx, payload, request.Identity)
	if err != nil {
		var appErr *app.ApplicationError
		if errors.As(err, &appErr) {
			appResult := baseResult(request)
			appResult.Error = appErr
			return appResult, appErr
		}
		return baseResult(request), err
	}

	appResult := baseResult(request)
	appResult.PayloadType = "presence.GetPlayerPresenceResult"
	appResult.Payload = result
	return appResult, nil
}

func (h Handlers) GetPlayerPresence(ctx context.Context, request GetPlayerPresenceRequest, identity app.RequestIdentity) (GetPlayerPresenceResult, error) {
	playerID, err := playerIDForSelfQuery(request, identity)
	if err != nil {
		return GetPlayerPresenceResult{}, err
	}
	if h.Registry == nil {
		return GetPlayerPresenceResult{}, &app.ApplicationError{
			Code:    ErrorCodePresenceQueryUnavailable,
			Message: "presence query is unavailable",
			Route:   GetPlayerPresenceRoute(),
		}
	}

	snapshot := h.Registry.PresenceForPlayer(ctx, connection.PlayerID(playerID))
	return resultFromSnapshot(snapshot), nil
}

func playerIDForSelfQuery(request GetPlayerPresenceRequest, identity app.RequestIdentity) (string, error) {
	if identity.Status != app.IdentityValidationValidated ||
		identity.ActorKind != app.ActorKindPlayer ||
		!identity.PlayerIDValidated ||
		strings.TrimSpace(identity.PlayerID) == "" {
		return "", &app.ApplicationError{
			Code:    ErrorCodePresenceQueryForbidden,
			Message: "validated player identity is required",
			Route:   GetPlayerPresenceRoute(),
		}
	}

	authenticatedPlayerID := strings.TrimSpace(identity.PlayerID)
	requestedPlayerID := strings.TrimSpace(request.PlayerID)
	if requestedPlayerID == "" {
		return authenticatedPlayerID, nil
	}
	if requestedPlayerID != authenticatedPlayerID {
		return "", &app.ApplicationError{
			Code:    ErrorCodePresenceQueryForbidden,
			Message: "presence query may only target the authenticated player",
			Route:   GetPlayerPresenceRoute(),
		}
	}
	return authenticatedPlayerID, nil
}

func resultFromSnapshot(snapshot connection.PlayerPresence) GetPlayerPresenceResult {
	result := GetPlayerPresenceResult{
		PlayerID:          string(snapshot.PlayerID),
		Status:            presenceStatusFromConnection(snapshot.Status),
		ConnectionCount:   snapshot.ConnectionCount,
		LastSeenAt:        copyTime(snapshot.LastSeenAt),
		RuntimeSessionIDs: make([]string, 0, len(snapshot.RuntimeSessionIDs)),
		ObservedAt:        snapshot.ObservedAt,
		ActiveConnections: make([]PresenceConnection, 0, len(snapshot.ActiveConnections)),
	}
	for _, sessionID := range snapshot.RuntimeSessionIDs {
		if normalized := strings.TrimSpace(string(sessionID)); normalized != "" {
			result.RuntimeSessionIDs = append(result.RuntimeSessionIDs, normalized)
		}
	}
	for _, active := range snapshot.ActiveConnections {
		result.ActiveConnections = append(result.ActiveConnections, PresenceConnection{
			ConnectionID:     string(active.ConnectionID),
			ConnectionEpoch:  uint64(active.ConnectionEpoch),
			RuntimeSessionID: string(active.RuntimeSessionID),
			LastSeenAt:       copyTime(active.LastSeenAt),
			OpenedAt:         active.OpenedAt,
			BoundAt:          copyTime(active.BoundAt),
		})
	}
	return result
}

func presenceStatusFromConnection(status connection.PresenceStatus) PresenceStatus {
	switch status {
	case connection.PresenceStatusOnline:
		return PresenceStatusOnline
	default:
		return PresenceStatusOffline
	}
}

func applicationErrorResult(request app.RouteRequest, code app.ErrorCode, message string) (app.ApplicationResult, error) {
	result := baseResult(request)
	appErr := &app.ApplicationError{
		Code:    code,
		Message: message,
		Route:   request.Route,
	}
	result.Error = appErr
	return result, appErr
}

func baseResult(request app.RouteRequest) app.ApplicationResult {
	return app.ApplicationResult{
		RequestID: request.RequestID,
		Route:     request.Route,
		Target:    request.Target,
		Session:   request.Session,
		Identity:  request.Identity,
	}
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
