package app

import (
	"context"
	"strings"
	"time"
)

type ConnectionBindingStatus string

const (
	ConnectionBindingStatusUnspecified ConnectionBindingStatus = "unspecified"
	ConnectionBindingStatusBound       ConnectionBindingStatus = "bound"
	ConnectionBindingStatusRejected    ConnectionBindingStatus = "rejected"
)

type ConnectionBindingClock interface {
	Now() time.Time
}

type ConnectionBindingRequest struct {
	AccessToken      string
	Route            RouteKey
	ConnectionID     string
	ConnectionEpoch  uint64
	ClientInstanceID string
}

type ConnectionBindingResult struct {
	BindingStatus    ConnectionBindingStatus
	Bound            bool
	Identity         RequestIdentity
	PublicErrorCode  ErrorCode
	ConnectionID     string
	ConnectionEpoch  uint64
	ClientInstanceID string
	BoundAt          time.Time
}

type ConnectionBinder struct {
	Validator RouteAccessTokenValidator
	Clock     ConnectionBindingClock
}

func NewConnectionBinder(validator RouteAccessTokenValidator) ConnectionBinder {
	return ConnectionBinder{Validator: validator}
}

func (b ConnectionBinder) BindConnection(ctx context.Context, request ConnectionBindingRequest) (ConnectionBindingResult, error) {
	route := normalizeRouteForAuthentication(request.Route)
	if route == (RouteKey{}) {
		route = BindConnectionRoute()
	}

	connectionID := strings.TrimSpace(request.ConnectionID)
	if connectionID == "" {
		return rejectedConnectionBindingResult(request, ErrorCodeConnectionBindingUnavailable),
			connectionBindingError(route, ErrorCodeConnectionBindingUnavailable)
	}

	accessToken := request.AccessToken
	if strings.TrimSpace(accessToken) == "" {
		return rejectedConnectionBindingResult(request, ErrorCodeConnectionBindingTokenMissing),
			connectionBindingError(route, ErrorCodeConnectionBindingTokenMissing)
	}
	if strings.TrimSpace(accessToken) != accessToken {
		return rejectedConnectionBindingResult(request, ErrorCodeConnectionBindingTokenMalformed),
			connectionBindingError(route, ErrorCodeConnectionBindingTokenMalformed)
	}

	validator := b.Validator
	if validator == nil {
		return rejectedConnectionBindingResult(request, ErrorCodeConnectionBindingUnavailable),
			connectionBindingError(route, ErrorCodeConnectionBindingUnavailable)
	}

	validation, err := validator.ValidateRouteAccessToken(ctx, RouteAccessTokenValidationRequest{
		AccessToken:     accessToken,
		Route:           route,
		ConnectionID:    connectionID,
		ConnectionEpoch: request.ConnectionEpoch,
	})
	if err != nil {
		code := connectionBindingErrorCode(validation.PublicErrorCode, true)
		return rejectedConnectionBindingResult(request, code), connectionBindingError(route, code)
	}
	if !validation.Valid {
		code := connectionBindingErrorCode(validation.PublicErrorCode, false)
		return rejectedConnectionBindingResult(request, code), connectionBindingError(route, code)
	}
	if !identitySatisfiesConnectionBinding(validation.Identity) {
		return rejectedConnectionBindingResult(request, ErrorCodeConnectionBindingTokenInvalid),
			connectionBindingError(route, ErrorCodeConnectionBindingTokenInvalid)
	}

	identity := validation.Identity
	identity.ConnectionID = connectionID
	identity.ConnectionEpoch = request.ConnectionEpoch
	identity.SessionValidated = false

	return ConnectionBindingResult{
		BindingStatus:    ConnectionBindingStatusBound,
		Bound:            true,
		Identity:         identity,
		ConnectionID:     connectionID,
		ConnectionEpoch:  request.ConnectionEpoch,
		ClientInstanceID: strings.TrimSpace(request.ClientInstanceID),
		BoundAt:          b.now(),
	}, nil
}

func identitySatisfiesConnectionBinding(identity RequestIdentity) bool {
	return identity.Status == IdentityValidationValidated &&
		identity.ActorKind == ActorKindPlayer &&
		strings.TrimSpace(identity.PlayerID) != "" &&
		identity.PlayerIDValidated
}

func rejectedConnectionBindingResult(request ConnectionBindingRequest, code ErrorCode) ConnectionBindingResult {
	return ConnectionBindingResult{
		BindingStatus:    ConnectionBindingStatusRejected,
		Bound:            false,
		PublicErrorCode:  code,
		ConnectionID:     strings.TrimSpace(request.ConnectionID),
		ConnectionEpoch:  request.ConnectionEpoch,
		ClientInstanceID: strings.TrimSpace(request.ClientInstanceID),
		Identity: RequestIdentity{
			Status: IdentityValidationUnknown,
		},
	}
}

func connectionBindingErrorCode(code ErrorCode, dependencyFailed bool) ErrorCode {
	switch code {
	case ErrorCodeAuthenticationTokenMissing, ErrorCodeConnectionBindingTokenMissing:
		return ErrorCodeConnectionBindingTokenMissing
	case ErrorCodeAuthenticationTokenMalformed, ErrorCodeConnectionBindingTokenMalformed:
		return ErrorCodeConnectionBindingTokenMalformed
	case ErrorCodeAuthenticationTokenUnavailable, ErrorCodeConnectionBindingUnavailable:
		return ErrorCodeConnectionBindingUnavailable
	case ErrorCodeAuthenticationTokenInvalid, ErrorCodeConnectionBindingTokenInvalid:
		return ErrorCodeConnectionBindingTokenInvalid
	default:
		if dependencyFailed {
			return ErrorCodeConnectionBindingUnavailable
		}
		return ErrorCodeConnectionBindingTokenInvalid
	}
}

func connectionBindingError(route RouteKey, code ErrorCode) *ApplicationError {
	return &ApplicationError{
		Code:    code,
		Message: connectionBindingMessage(code),
		Route:   route,
	}
}

func connectionBindingMessage(code ErrorCode) string {
	switch code {
	case ErrorCodeConnectionBindingTokenMissing:
		return "connection binding token proof is missing"
	case ErrorCodeConnectionBindingTokenMalformed:
		return "connection binding token proof is malformed"
	case ErrorCodeConnectionBindingUnavailable:
		return "connection binding validation is unavailable"
	case ErrorCodeConnectionBindingRequired:
		return "connection binding is required"
	default:
		return "connection binding token proof is invalid"
	}
}

func (b ConnectionBinder) now() time.Time {
	if b.Clock == nil {
		return time.Now().UTC()
	}
	value := b.Clock.Now()
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}
