package main

import (
	"context"

	appconnection "github.com/iceiko/vibit/runtime/internal/app/connection"
	"github.com/iceiko/vibit/runtime/internal/platform/transport/ws"
)

type connectionLifecycleRegistryObserver struct {
	registry *appconnection.InMemoryRegistry
}

func (o connectionLifecycleRegistryObserver) ConnectionOpened(ctx context.Context, event ws.ConnectionLifecycleEvent) error {
	if o.registry == nil {
		return nil
	}
	_, err := o.registry.RegisterOpenConnection(ctx, appconnection.OpenConnection{
		ConnectionID:    appconnection.ConnectionID(event.ConnectionID),
		ConnectionEpoch: appconnection.ConnectionEpoch(event.ConnectionEpoch),
		OpenedAt:        event.ObservedAt,
	})
	return err
}

func (o connectionLifecycleRegistryObserver) ConnectionClosed(ctx context.Context, event ws.ConnectionLifecycleEvent) error {
	if o.registry == nil {
		return nil
	}
	_, err := o.registry.MarkConnectionClosed(ctx, appconnection.MarkClosed{
		ConnectionID:     appconnection.ConnectionID(event.ConnectionID),
		ConnectionEpoch:  appconnection.ConnectionEpoch(event.ConnectionEpoch),
		ClosedAt:         event.ObservedAt,
		CloseReasonClass: "transport_closed",
	})
	return err
}
