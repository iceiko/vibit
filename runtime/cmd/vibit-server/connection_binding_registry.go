package main

import (
	"context"

	"github.com/iceiko/vibit/runtime/internal/app"
	appconnection "github.com/iceiko/vibit/runtime/internal/app/connection"
)

type registryConnectionBinder struct {
	binder   app.ConnectionBinder
	registry *appconnection.InMemoryRegistry
}

func (b registryConnectionBinder) BindConnection(ctx context.Context, request app.ConnectionBindingRequest) (app.ConnectionBindingResult, error) {
	result, err := b.binder.BindConnection(ctx, request)
	if err != nil || !result.Bound || b.registry == nil {
		return result, err
	}

	recordErr := b.bindRegistryIdentity(ctx, result)
	if recordErr != nil {
		return app.ConnectionBindingResult{
				BindingStatus:    app.ConnectionBindingStatusRejected,
				Bound:            false,
				PublicErrorCode:  app.ErrorCodeConnectionBindingUnavailable,
				ConnectionID:     result.ConnectionID,
				ConnectionEpoch:  result.ConnectionEpoch,
				ClientInstanceID: result.ClientInstanceID,
				Identity: app.RequestIdentity{
					Status: app.IdentityValidationUnknown,
				},
			}, &app.ApplicationError{
				Code:    app.ErrorCodeConnectionBindingUnavailable,
				Message: "connection binding registry is unavailable",
				Route:   app.BindConnectionRoute(),
			}
	}

	return result, nil
}

func (b registryConnectionBinder) bindRegistryIdentity(ctx context.Context, result app.ConnectionBindingResult) error {
	_, err := b.registry.BindConnectionIdentity(ctx, appconnection.BindIdentity{
		ConnectionID:     appconnection.ConnectionID(result.ConnectionID),
		ConnectionEpoch:  appconnection.ConnectionEpoch(result.ConnectionEpoch),
		ActorKind:        appconnection.ActorKind(result.Identity.ActorKind),
		PlayerID:         appconnection.PlayerID(result.Identity.PlayerID),
		RuntimeSessionID: appconnection.RuntimeSessionID(result.Identity.SessionID),
		ValidatedAt:      result.BoundAt,
	})
	return err
}
