package protobuf

import (
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
)

func BuildErrorEnvelopeFromApplicationResult(result app.ApplicationResult) (*protocolv1.Envelope, error) {
	if result.Error == nil {
		return nil, &PayloadBridgeError{Message: "application result has no error"}
	}
	return BuildErrorEnvelopeFromApplicationError(result.RequestID, result.Error, result.Target, result.Session)
}

func BuildErrorEnvelopeFromApplicationError(requestID string, applicationError *app.ApplicationError, target app.Target, session app.Session) (*protocolv1.Envelope, error) {
	if applicationError == nil {
		return nil, &PayloadBridgeError{Message: "application error is nil"}
	}

	route := app.RouteKey{
		Kind:   app.MessageKindError,
		Module: strings.TrimSpace(applicationError.Route.Module),
		Name:   strings.TrimSpace(applicationError.Route.Name),
	}
	if route.Module == "" {
		route.Module = "system"
	}
	if route.Name == "" {
		route.Name = "ApplicationError"
	}

	normalizedRequestID := strings.TrimSpace(requestID)
	return &protocolv1.Envelope{
		ProtocolVersion: 1,
		Kind:            protocolv1.MessageKind_MESSAGE_KIND_ERROR,
		RequestId:       normalizedRequestID,
		Module:          route.Module,
		Name:            route.Name,
		Target: &protocolv1.Target{
			Scope: protoTargetScope(target.Scope),
			Id:    strings.TrimSpace(target.ID),
		},
		Session: &protocolv1.Session{
			ConnectionId:    strings.TrimSpace(session.ConnectionID),
			SessionId:       strings.TrimSpace(session.SessionID),
			PlayerId:        strings.TrimSpace(session.PlayerID),
			ConnectionEpoch: session.ConnectionEpoch,
		},
		Error: &protocolv1.Error{
			Code:      strings.TrimSpace(string(applicationError.Code)),
			Message:   strings.TrimSpace(applicationError.Message),
			RequestId: normalizedRequestID,
			Retryable: false,
		},
	}, nil
}
