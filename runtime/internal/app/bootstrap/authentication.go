package bootstrap

import (
	"context"
	"errors"

	"github.com/iceiko/vibit/runtime/internal/app"
	appauth "github.com/iceiko/vibit/runtime/internal/app/authentication"
)

type AuthenticationService interface {
	AuthenticateWithDeviceCredential(context.Context, appauth.DeviceCredentialAuthenticationRequest) (appauth.AuthenticationResult, error)
	LogoutAccessToken(context.Context, appauth.LogoutAccessTokenRequest) (appauth.LogoutAccessTokenResult, error)
}

type AuthenticationRouteHandlers struct {
	Service AuthenticationService
}

func (h AuthenticationRouteHandlers) RegisterRoutes(dispatcher *app.Dispatcher) error {
	if dispatcher == nil {
		return errors.New("authentication bootstrap: dispatcher is nil")
	}
	if err := dispatcher.Register(app.AuthenticateWithDeviceCredentialRoute(), app.HandlerFunc(h.HandleAuthenticateWithDeviceCredentialRoute)); err != nil {
		return err
	}
	return dispatcher.Register(app.LogoutAccessTokenRoute(), app.HandlerFunc(h.HandleLogoutAccessTokenRoute))
}

func (h AuthenticationRouteHandlers) HandleAuthenticateWithDeviceCredentialRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		result := resultForRequest(request)
		appErr := authenticationApplicationError(request.Route, appauth.PublicErrorAuthenticationNotImplemented)
		result.Error = appErr
		return result, appErr
	}

	payload, ok := request.Payload.(appauth.DeviceCredentialAuthenticationRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*appauth.DeviceCredentialAuthenticationRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		result := resultForRequest(request)
		appErr := authenticationApplicationError(request.Route, appauth.PublicErrorAuthenticationProofMalformed)
		result.Error = appErr
		return result, appErr
	}

	authenticationResult, err := service.AuthenticateWithDeviceCredential(ctx, payload)
	result := resultForRequest(request)
	if err != nil {
		appErr := authenticationApplicationError(request.Route, authenticationPublicErrorCode(authenticationResult, err))
		result.Error = appErr
		return result, appErr
	}

	result.PayloadType = "authentication.AuthenticationResult"
	result.Payload = authenticationResult
	return result, nil
}

func (h AuthenticationRouteHandlers) HandleLogoutAccessTokenRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		result := resultForRequest(request)
		appErr := authenticationApplicationError(request.Route, appauth.PublicErrorAuthenticationNotImplemented)
		result.Error = appErr
		return result, appErr
	}

	payload, ok := request.Payload.(appauth.LogoutAccessTokenRequest)
	if !ok {
		if pointerPayload, pointerOK := request.Payload.(*appauth.LogoutAccessTokenRequest); pointerOK && pointerPayload != nil {
			payload = *pointerPayload
			ok = true
		}
	}
	if !ok {
		result := resultForRequest(request)
		appErr := authenticationApplicationError(request.Route, appauth.PublicErrorAuthenticationTokenMalformed)
		result.Error = appErr
		return result, appErr
	}

	logoutResult, err := service.LogoutAccessToken(ctx, payload)
	result := resultForRequest(request)
	if err != nil {
		appErr := authenticationApplicationError(request.Route, logoutPublicErrorCode(logoutResult, err))
		result.Error = appErr
		return result, appErr
	}

	result.PayloadType = "authentication.LogoutAccessTokenResult"
	result.Payload = logoutResult
	return result, nil
}

func authenticationPublicErrorCode(result appauth.AuthenticationResult, err error) appauth.PublicErrorCode {
	if result.PublicErrorCode != "" {
		return result.PublicErrorCode
	}

	var serviceError *appauth.ServiceError
	if errors.As(err, &serviceError) && serviceError.PublicCode != "" {
		return serviceError.PublicCode
	}

	return appauth.PublicErrorAuthenticationCredentialInvalid
}

func logoutPublicErrorCode(result appauth.LogoutAccessTokenResult, err error) appauth.PublicErrorCode {
	if result.PublicErrorCode != "" {
		return result.PublicErrorCode
	}

	var serviceError *appauth.ServiceError
	if errors.As(err, &serviceError) && serviceError.PublicCode != "" {
		return serviceError.PublicCode
	}

	return appauth.PublicErrorAuthenticationTokenInvalid
}

func authenticationApplicationError(route app.RouteKey, publicCode appauth.PublicErrorCode) *app.ApplicationError {
	code := app.ErrorCode(publicCode)
	if code == "" {
		code = app.ErrorCode(appauth.PublicErrorAuthenticationCredentialInvalid)
	}
	return &app.ApplicationError{
		Code:    code,
		Message: authenticationErrorMessage(publicCode),
		Route:   route,
	}
}

func authenticationErrorMessage(code appauth.PublicErrorCode) string {
	switch code {
	case appauth.PublicErrorAuthenticationProofMissing:
		return "authentication proof is missing"
	case appauth.PublicErrorAuthenticationProofMalformed:
		return "authentication proof is malformed"
	case appauth.PublicErrorAuthenticationTokenMissing:
		return "authentication token proof is missing"
	case appauth.PublicErrorAuthenticationTokenMalformed:
		return "authentication token proof is malformed"
	case appauth.PublicErrorAuthenticationTokenInvalid:
		return "authentication token proof is invalid"
	case appauth.PublicErrorAuthenticationCredentialUnavailable:
		return "authentication credential store is unavailable"
	case appauth.PublicErrorAuthenticationTokenUnavailable:
		return "authentication token store is unavailable"
	case appauth.PublicErrorAuthenticationNotImplemented:
		return "authentication is not implemented"
	default:
		return "authentication credential is invalid"
	}
}
