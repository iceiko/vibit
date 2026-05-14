package app

import (
	"context"
	"errors"
	"strings"
)

type SessionValidator interface {
	ValidateSession(context.Context, RouteRequest) (SessionValidationResult, error)
}

type SessionValidatorFunc func(context.Context, RouteRequest) (SessionValidationResult, error)

func (f SessionValidatorFunc) ValidateSession(ctx context.Context, request RouteRequest) (SessionValidationResult, error) {
	if f == nil {
		return SessionValidationResult{}, errors.New("app: session validator function is nil")
	}
	return f(ctx, request)
}

type MetadataOnlySessionValidator struct{}

func (MetadataOnlySessionValidator) ValidateSession(_ context.Context, request RouteRequest) (SessionValidationResult, error) {
	identity := request.Identity
	if identity.Status == "" || identity.Status == IdentityValidationMetadataOnly {
		identity = MetadataOnlyIdentityFromSession(request.Session)
	}
	return SessionValidationResult{
		Identity: identity,
		Valid:    true,
		Reason:   "metadata_only",
	}, nil
}

type SessionValidatingDispatcher struct {
	Dispatcher RouteDispatcher
	Validator  SessionValidator
}

func (d SessionValidatingDispatcher) Dispatch(ctx context.Context, request RouteRequest) (ApplicationResult, error) {
	dispatcher := d.Dispatcher
	if dispatcher == nil {
		return resultForRequest(request), errors.New("app: session validating dispatcher requires dispatcher")
	}

	validator := d.Validator
	if validator == nil {
		validator = MetadataOnlySessionValidator{}
	}

	validation, err := validator.ValidateSession(ctx, request)
	if err != nil {
		return resultForRequest(request), err
	}
	if !validation.Valid {
		result := resultForRequest(request)
		appErr := &ApplicationError{
			Code:    ErrorCodeSessionInvalid,
			Message: sessionValidationMessage(validation.Reason),
			Route:   request.Route,
		}
		result.Error = appErr
		return result, appErr
	}

	if validation.Identity.Status == "" {
		validation.Identity = MetadataOnlyIdentityFromSession(request.Session)
	}
	request.Identity = validation.Identity
	return dispatcher.Dispatch(ctx, request)
}

func sessionValidationMessage(reason string) string {
	normalizedReason := strings.TrimSpace(reason)
	if normalizedReason == "" {
		return "session validation failed"
	}
	return normalizedReason
}
