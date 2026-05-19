package authentication

import (
	"context"

	"github.com/iceiko/vibit/runtime/internal/app"
)

type RouteAccessTokenValidator struct {
	Service Service
}

func NewRouteAccessTokenValidator(service Service) RouteAccessTokenValidator {
	return RouteAccessTokenValidator{Service: service}
}

func (v RouteAccessTokenValidator) ValidateRouteAccessToken(ctx context.Context, request app.RouteAccessTokenValidationRequest) (app.RouteAccessTokenValidationResult, error) {
	result, err := v.Service.ValidateAccessToken(ctx, AccessTokenValidationRequest{
		AccessToken:     request.AccessToken,
		Route:           request.Route,
		ConnectionID:    request.ConnectionID,
		ConnectionEpoch: request.ConnectionEpoch,
	})
	if result.Status == ValidationStatusValidated {
		return app.RouteAccessTokenValidationResult{
			Valid:    true,
			Identity: result.Identity,
		}, err
	}
	return app.RouteAccessTokenValidationResult{
		Valid:           false,
		Identity:        result.Identity,
		PublicErrorCode: app.ErrorCode(result.PublicErrorCode),
	}, err
}
