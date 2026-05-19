package app

import (
	"context"
	"errors"
	"strings"
)

const (
	ErrorCodeAuthenticationTokenMissing      ErrorCode = "AUTHENTICATION_TOKEN_MISSING"
	ErrorCodeAuthenticationTokenMalformed    ErrorCode = "AUTHENTICATION_TOKEN_MALFORMED"
	ErrorCodeAuthenticationTokenInvalid      ErrorCode = "AUTHENTICATION_TOKEN_INVALID"
	ErrorCodeAuthenticationTokenUnavailable  ErrorCode = "AUTHENTICATION_TOKEN_STORE_UNAVAILABLE"
	ErrorCodeConnectionBindingTokenMissing   ErrorCode = "CONNECTION_BINDING_TOKEN_MISSING"
	ErrorCodeConnectionBindingTokenMalformed ErrorCode = "CONNECTION_BINDING_TOKEN_MALFORMED"
	ErrorCodeConnectionBindingTokenInvalid   ErrorCode = "CONNECTION_BINDING_TOKEN_INVALID"
	ErrorCodeConnectionBindingUnavailable    ErrorCode = "CONNECTION_BINDING_UNAVAILABLE"
	ErrorCodeConnectionBindingRequired       ErrorCode = "CONNECTION_BINDING_REQUIRED"
)

var errRouteAccessTokenValidatorRequired = errors.New("app: route access-token validator is required")

type RouteProtectionRequirement string

const (
	RouteProtectionPublic                   RouteProtectionRequirement = "public"
	RouteProtectionRequestTokenRequired     RouteProtectionRequirement = "request_token_required"
	RouteProtectionBoundConnectionRequired  RouteProtectionRequirement = "bound_connection_required"
	RouteProtectionSessionValidatedRequired RouteProtectionRequirement = "session_validated_required"
	RouteProtectionBoundSessionRequired     RouteProtectionRequirement = "bound_session_required"
)

type RouteProtectionRouteRequirement struct {
	Route       RouteKey
	Requirement RouteProtectionRequirement
}

type RouteProtectionPolicy struct {
	PublicRoutes      []RouteKey
	RouteRequirements []RouteProtectionRouteRequirement
}

func DefaultRouteProtectionPolicy() RouteProtectionPolicy {
	return RouteProtectionPolicy{
		PublicRoutes: []RouteKey{
			AuthenticateWithDeviceCredentialRoute(),
			LogoutAccessTokenRoute(),
		},
	}
}

func AuthenticateWithDeviceCredentialRoute() RouteKey {
	return RouteKey{
		Kind:   MessageKindCommand,
		Module: "runtime.authentication",
		Name:   "AuthenticateWithDeviceCredential",
	}
}

func LogoutAccessTokenRoute() RouteKey {
	return RouteKey{
		Kind:   MessageKindCommand,
		Module: "runtime.authentication",
		Name:   "LogoutAccessToken",
	}
}

func BindConnectionRoute() RouteKey {
	return RouteKey{
		Kind:   MessageKindSystem,
		Module: "runtime.authentication",
		Name:   "BindConnection",
	}
}

func (p RouteProtectionPolicy) IsPublic(route RouteKey) bool {
	return p.RequirementFor(route) == RouteProtectionPublic
}

func (p RouteProtectionPolicy) RequirementFor(route RouteKey) RouteProtectionRequirement {
	normalizedRoute := normalizeRouteForAuthentication(route)
	if normalizedRoute == (RouteKey{}) {
		return RouteProtectionRequestTokenRequired
	}

	for _, routeRequirement := range p.RouteRequirements {
		if normalizeRouteForAuthentication(routeRequirement.Route) == normalizedRoute {
			return normalizeRouteProtectionRequirement(routeRequirement.Requirement)
		}
	}

	publicRoutes := p.PublicRoutes
	if len(publicRoutes) == 0 {
		publicRoutes = DefaultRouteProtectionPolicy().PublicRoutes
	}
	for _, publicRoute := range publicRoutes {
		if normalizeRouteForAuthentication(publicRoute) == normalizedRoute {
			return RouteProtectionPublic
		}
	}
	return RouteProtectionRequestTokenRequired
}

type RouteAccessTokenValidationRequest struct {
	AccessToken     string
	Route           RouteKey
	ConnectionID    string
	ConnectionEpoch uint64
}

type RouteAccessTokenValidationResult struct {
	Valid           bool
	Identity        RequestIdentity
	PublicErrorCode ErrorCode
}

type RouteAccessTokenValidator interface {
	ValidateRouteAccessToken(context.Context, RouteAccessTokenValidationRequest) (RouteAccessTokenValidationResult, error)
}

type RouteProtectionRequest struct {
	Request                  RouteRequest
	AccessToken              string
	ProofCarrierPresent      bool
	BoundIdentity            RequestIdentity
	SessionValidatedIdentity RequestIdentity
}

type RouteProtectionResult struct {
	Allowed  bool
	Public   bool
	Identity RequestIdentity
}

type RouteProtector struct {
	Policy    RouteProtectionPolicy
	Validator RouteAccessTokenValidator
}

func NewRouteProtector(validator RouteAccessTokenValidator) RouteProtector {
	return RouteProtector{
		Policy:    DefaultRouteProtectionPolicy(),
		Validator: validator,
	}
}

func (p RouteProtector) ProtectRoute(ctx context.Context, request RouteProtectionRequest) (RouteProtectionResult, error) {
	route := normalizeRouteForAuthentication(request.Request.Route)
	request.Request.Route = route
	requirement := p.Policy.RequirementFor(route)
	if requirement == RouteProtectionPublic {
		identity := request.Request.Identity
		if identity.Status == "" {
			identity = MetadataOnlyIdentityFromSession(request.Request.Session)
		}
		return RouteProtectionResult{
			Allowed:  true,
			Public:   true,
			Identity: identity,
		}, nil
	}

	switch requirement {
	case RouteProtectionBoundConnectionRequired:
		return protectBoundConnectionRoute(route, request)
	case RouteProtectionSessionValidatedRequired:
		return protectSessionValidatedRoute(route, request)
	case RouteProtectionBoundSessionRequired:
		return protectBoundSessionRoute(route, request)
	}

	if !request.ProofCarrierPresent || strings.TrimSpace(request.AccessToken) == "" {
		return RouteProtectionResult{}, routeAuthenticationError(route, ErrorCodeAuthenticationTokenMissing)
	}

	validator := p.Validator
	if validator == nil {
		return RouteProtectionResult{}, routeAuthenticationError(route, ErrorCodeAuthenticationTokenUnavailable)
	}

	validation, err := validator.ValidateRouteAccessToken(ctx, RouteAccessTokenValidationRequest{
		AccessToken:     request.AccessToken,
		Route:           route,
		ConnectionID:    request.Request.Session.ConnectionID,
		ConnectionEpoch: request.Request.Session.ConnectionEpoch,
	})
	if err != nil {
		code := validation.PublicErrorCode
		if code == "" {
			code = ErrorCodeAuthenticationTokenUnavailable
		}
		return RouteProtectionResult{}, routeAuthenticationError(route, normalizeAuthenticationErrorCode(code))
	}
	if !validation.Valid {
		code := validation.PublicErrorCode
		if code == "" {
			code = ErrorCodeAuthenticationTokenInvalid
		}
		return RouteProtectionResult{}, routeAuthenticationError(route, normalizeAuthenticationErrorCode(code))
	}
	if !identitySatisfiesProtectedRoute(validation.Identity) {
		return RouteProtectionResult{}, routeAuthenticationError(route, ErrorCodeAuthenticationTokenInvalid)
	}

	identity := validation.Identity
	identity.SessionValidated = false
	return RouteProtectionResult{
		Allowed:  true,
		Identity: identity,
	}, nil
}

func protectBoundConnectionRoute(route RouteKey, request RouteProtectionRequest) (RouteProtectionResult, error) {
	if !identitySatisfiesBoundConnectionRoute(request.BoundIdentity, request.Request.Session) {
		return RouteProtectionResult{}, routeAuthenticationError(route, ErrorCodeConnectionBindingRequired)
	}

	identity := request.BoundIdentity
	identity.SessionValidated = false
	return RouteProtectionResult{
		Allowed:  true,
		Identity: identity,
	}, nil
}

func protectSessionValidatedRoute(route RouteKey, request RouteProtectionRequest) (RouteProtectionResult, error) {
	if !identitySatisfiesSessionValidatedRoute(request.SessionValidatedIdentity, request.Request.Session) {
		return RouteProtectionResult{}, routeAuthenticationError(route, ErrorCodeSessionInvalid)
	}

	return RouteProtectionResult{
		Allowed:  true,
		Identity: request.SessionValidatedIdentity,
	}, nil
}

func protectBoundSessionRoute(route RouteKey, request RouteProtectionRequest) (RouteProtectionResult, error) {
	if !identitySatisfiesBoundConnectionRoute(request.BoundIdentity, request.Request.Session) {
		return RouteProtectionResult{}, routeAuthenticationError(route, ErrorCodeConnectionBindingRequired)
	}
	if !identitySatisfiesSessionValidatedRoute(request.SessionValidatedIdentity, request.Request.Session) {
		return RouteProtectionResult{}, routeAuthenticationError(route, ErrorCodeSessionInvalid)
	}
	if !validatedIdentitiesAgree(request.BoundIdentity, request.SessionValidatedIdentity) {
		return RouteProtectionResult{}, routeAuthenticationError(route, ErrorCodeConnectionBindingRequired)
	}

	identity := request.SessionValidatedIdentity
	if strings.TrimSpace(identity.ConnectionID) == "" {
		identity.ConnectionID = strings.TrimSpace(request.BoundIdentity.ConnectionID)
	}
	if identity.ConnectionEpoch == 0 {
		identity.ConnectionEpoch = request.BoundIdentity.ConnectionEpoch
	}
	return RouteProtectionResult{
		Allowed:  true,
		Identity: identity,
	}, nil
}

func identitySatisfiesProtectedRoute(identity RequestIdentity) bool {
	return identity.Status == IdentityValidationValidated &&
		identity.ActorKind == ActorKindPlayer &&
		strings.TrimSpace(identity.PlayerID) != "" &&
		identity.PlayerIDValidated
}

func identitySatisfiesBoundConnectionRoute(identity RequestIdentity, session Session) bool {
	if !identitySatisfiesProtectedRoute(identity) {
		return false
	}
	connectionID := strings.TrimSpace(identity.ConnectionID)
	if connectionID == "" {
		return false
	}
	requestConnectionID := strings.TrimSpace(session.ConnectionID)
	if requestConnectionID != "" && connectionID != requestConnectionID {
		return false
	}
	if session.ConnectionEpoch != 0 && identity.ConnectionEpoch != session.ConnectionEpoch {
		return false
	}
	return true
}

func identitySatisfiesSessionValidatedRoute(identity RequestIdentity, session Session) bool {
	if !identitySatisfiesProtectedRoute(identity) || !identity.SessionValidated {
		return false
	}
	sessionID := strings.TrimSpace(identity.SessionID)
	if sessionID == "" {
		return false
	}
	requestSessionID := strings.TrimSpace(session.SessionID)
	if requestSessionID != "" && sessionID != requestSessionID {
		return false
	}
	return true
}

func validatedIdentitiesAgree(left RequestIdentity, right RequestIdentity) bool {
	leftPlayerID := strings.TrimSpace(left.PlayerID)
	rightPlayerID := strings.TrimSpace(right.PlayerID)
	leftActorID := strings.TrimSpace(left.ActorID)
	rightActorID := strings.TrimSpace(right.ActorID)
	if left.ActorKind != right.ActorKind ||
		left.ActorKind != ActorKindPlayer ||
		leftPlayerID == "" ||
		leftPlayerID != rightPlayerID ||
		leftActorID == "" ||
		leftActorID != rightActorID {
		return false
	}

	leftSessionID := strings.TrimSpace(left.SessionID)
	rightSessionID := strings.TrimSpace(right.SessionID)
	if leftSessionID != "" && rightSessionID != "" && leftSessionID != rightSessionID {
		return false
	}

	leftConnectionID := strings.TrimSpace(left.ConnectionID)
	rightConnectionID := strings.TrimSpace(right.ConnectionID)
	if leftConnectionID != "" && rightConnectionID != "" && leftConnectionID != rightConnectionID {
		return false
	}
	if left.ConnectionEpoch != 0 && right.ConnectionEpoch != 0 && left.ConnectionEpoch != right.ConnectionEpoch {
		return false
	}
	return true
}

func normalizeAuthenticationErrorCode(code ErrorCode) ErrorCode {
	switch code {
	case ErrorCodeAuthenticationTokenMissing,
		ErrorCodeAuthenticationTokenMalformed,
		ErrorCodeAuthenticationTokenInvalid,
		ErrorCodeAuthenticationTokenUnavailable:
		return code
	default:
		return ErrorCodeAuthenticationTokenInvalid
	}
}

func normalizeRouteProtectionRequirement(requirement RouteProtectionRequirement) RouteProtectionRequirement {
	switch requirement {
	case RouteProtectionPublic,
		RouteProtectionRequestTokenRequired,
		RouteProtectionBoundConnectionRequired,
		RouteProtectionSessionValidatedRequired,
		RouteProtectionBoundSessionRequired:
		return requirement
	default:
		return RouteProtectionRequestTokenRequired
	}
}

func routeAuthenticationError(route RouteKey, code ErrorCode) *ApplicationError {
	return &ApplicationError{
		Code:    code,
		Message: routeAuthenticationMessage(code),
		Route:   route,
	}
}

func routeAuthenticationMessage(code ErrorCode) string {
	switch code {
	case ErrorCodeAuthenticationTokenMissing:
		return "authentication token proof is missing"
	case ErrorCodeAuthenticationTokenMalformed:
		return "authentication token proof is malformed"
	case ErrorCodeAuthenticationTokenUnavailable:
		return "authentication token validation is unavailable"
	case ErrorCodeConnectionBindingRequired:
		return "connection binding is required"
	case ErrorCodeSessionInvalid:
		return "runtime session validation is required"
	default:
		return "authentication token proof is invalid"
	}
}

func normalizeRouteForAuthentication(route RouteKey) RouteKey {
	return RouteKey{
		Kind:   MessageKind(strings.TrimSpace(string(route.Kind))),
		Module: strings.TrimSpace(route.Module),
		Name:   strings.TrimSpace(route.Name),
	}
}
