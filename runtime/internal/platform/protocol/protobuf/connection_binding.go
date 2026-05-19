package protobuf

import (
	"context"
	"errors"
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	authenticationv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/authentication/v1"
	"google.golang.org/protobuf/proto"
)

func (h FrameHandler) handleConnectionBinding(ctx context.Context, request app.RouteRequest) ([][]byte, error) {
	decodedPayload, err := DecodePayload(request.PayloadType, request.PayloadBytes)
	if err != nil {
		return marshalApplicationResult(applicationErrorResultForRequest(request, connectionBindingMalformedError(request.Route)))
	}
	payload, ok := decodedPayload.(*authenticationv1.BindConnectionRequest)
	if !ok || payload == nil {
		return marshalApplicationResult(applicationErrorResultForRequest(request, connectionBindingMalformedError(request.Route)))
	}

	binder := h.ConnectionBinder
	if binder == nil {
		return marshalApplicationResult(applicationErrorResultForRequest(request, connectionBindingUnavailableError(request.Route)))
	}

	result, err := binder.BindConnection(ctx, app.ConnectionBindingRequest{
		AccessToken:      payload.GetAccessToken(),
		Route:            request.Route,
		ConnectionID:     request.Session.ConnectionID,
		ConnectionEpoch:  request.Session.ConnectionEpoch,
		ClientInstanceID: payload.GetClientInstanceId(),
	})
	applicationResult := app.ApplicationResult{
		RequestID: request.RequestID,
		Route:     request.Route,
		Target:    request.Target,
		Session: app.Session{
			ConnectionID:    request.Session.ConnectionID,
			ConnectionEpoch: request.Session.ConnectionEpoch,
		},
		Identity: request.Identity,
	}
	if result.Bound {
		applicationResult.Identity = result.Identity
		applicationResult.Session.ConnectionID = result.ConnectionID
		applicationResult.Session.ConnectionEpoch = result.ConnectionEpoch
		applicationResult.Session.PlayerID = result.Identity.PlayerID
		applicationResult.Payload = result
		return marshalApplicationResult(applicationResult)
	}

	var applicationError *app.ApplicationError
	if errors.As(err, &applicationError) {
		applicationResult.Error = applicationError
		return marshalApplicationResult(applicationResult)
	}
	if err != nil {
		return nil, err
	}
	applicationResult.Error = connectionBindingInvalidError(request.Route)
	return marshalApplicationResult(applicationResult)
}

func protoPayloadFromConnectionBindingRoute(route app.RouteKey, payload any) (proto.Message, bool, error) {
	if route != app.BindConnectionRoute() {
		return nil, false, nil
	}

	result, ok := payload.(app.ConnectionBindingResult)
	if !ok {
		if pointerResult, pointerOK := payload.(*app.ConnectionBindingResult); pointerOK && pointerResult != nil {
			result = *pointerResult
			ok = true
		}
	}
	if !ok {
		return nil, true, payloadBridgeError(route, "payload must be app.ConnectionBindingResult")
	}

	return &authenticationv1.BindConnectionResponse{
		BindingStatus:    protoConnectionBindingStatus(result.BindingStatus),
		ActorKind:        string(result.Identity.ActorKind),
		PlayerId:         result.Identity.PlayerID,
		ConnectionId:     result.ConnectionID,
		ConnectionEpoch:  result.ConnectionEpoch,
		SessionValidated: result.Identity.SessionValidated,
		BoundAt:          formatConnectionBindingTime(result.BoundAt),
	}, true, nil
}

func protoConnectionBindingStatus(status app.ConnectionBindingStatus) authenticationv1.ConnectionBindingStatus {
	switch status {
	case app.ConnectionBindingStatusBound:
		return authenticationv1.ConnectionBindingStatus_CONNECTION_BINDING_STATUS_BOUND
	case app.ConnectionBindingStatusRejected:
		return authenticationv1.ConnectionBindingStatus_CONNECTION_BINDING_STATUS_REJECTED
	default:
		return authenticationv1.ConnectionBindingStatus_CONNECTION_BINDING_STATUS_UNSPECIFIED
	}
}

func formatConnectionBindingTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func connectionBindingMalformedError(route app.RouteKey) *app.ApplicationError {
	return &app.ApplicationError{
		Code:    app.ErrorCodeConnectionBindingTokenMalformed,
		Message: "connection binding payload or token proof is malformed",
		Route:   route,
	}
}

func connectionBindingUnavailableError(route app.RouteKey) *app.ApplicationError {
	return &app.ApplicationError{
		Code:    app.ErrorCodeConnectionBindingUnavailable,
		Message: "connection binding validation is unavailable",
		Route:   route,
	}
}

func connectionBindingInvalidError(route app.RouteKey) *app.ApplicationError {
	return &app.ApplicationError{
		Code:    app.ErrorCodeConnectionBindingTokenInvalid,
		Message: "connection binding token proof is invalid",
		Route:   route,
	}
}
