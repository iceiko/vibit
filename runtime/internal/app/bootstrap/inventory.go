package bootstrap

import (
	"context"
	"errors"

	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

type InventoryRepositoryProvider interface {
	ForCommand(context.Context, tx.UnitOfWork) (inventory.Repository, error)
	ForQuery(context.Context) (inventory.Repository, error)
}

type StaticInventoryRepositoryProvider struct {
	Repository inventory.Repository
}

func (p StaticInventoryRepositoryProvider) ForCommand(context.Context, tx.UnitOfWork) (inventory.Repository, error) {
	return p.repository()
}

func (p StaticInventoryRepositoryProvider) ForQuery(context.Context) (inventory.Repository, error) {
	return p.repository()
}

func (p StaticInventoryRepositoryProvider) repository() (inventory.Repository, error) {
	if p.Repository == nil {
		return nil, errors.New("inventory bootstrap: repository is required")
	}
	return p.Repository, nil
}

type InventoryOptions struct {
	Repositories     InventoryRepositoryProvider
	PermissionPolicy inventory.PermissionPolicy
	CapacityPolicy   inventory.CapacityPolicy
	EventIDs         inventory.EventIDGenerator
	Clock            inventory.Clock
}

func NewInMemoryInventoryDispatcher() (*app.Dispatcher, error) {
	repository := inventory.NewMemoryRepository()
	return NewInventoryDispatcher(InventoryOptions{
		Repositories:     StaticInventoryRepositoryProvider{Repository: repository},
		PermissionPolicy: inventory.StaticPermissionPolicy{GrantAllowed: true, ReadAllowed: true},
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 256},
		EventIDs:         &inventory.IncrementingEventIDGenerator{Prefix: "inventory-event"},
		Clock:            inventory.SystemClock{},
	})
}

func NewInventoryDispatcher(options InventoryOptions) (*app.Dispatcher, error) {
	dispatcher := app.NewDispatcher()
	handlers := InventoryRouteHandlers{
		Repositories:     options.Repositories,
		PermissionPolicy: options.PermissionPolicy,
		CapacityPolicy:   options.CapacityPolicy,
		EventIDs:         options.EventIDs,
		Clock:            options.Clock,
	}
	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		return nil, err
	}
	return dispatcher, nil
}

type InventoryRouteHandlers struct {
	Repositories     InventoryRepositoryProvider
	PermissionPolicy inventory.PermissionPolicy
	CapacityPolicy   inventory.CapacityPolicy
	EventIDs         inventory.EventIDGenerator
	Clock            inventory.Clock
}

func (h InventoryRouteHandlers) RegisterRoutes(dispatcher *app.Dispatcher) error {
	if dispatcher == nil {
		return errors.New("inventory bootstrap: dispatcher is nil")
	}
	if err := dispatcher.Register(inventory.GrantItemRoute(), app.HandlerFunc(h.HandleGrantItemRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(inventory.GetInventoryRoute(), app.HandlerFunc(h.HandleGetInventoryRoute)); err != nil {
		return err
	}
	return nil
}

func (h InventoryRouteHandlers) HandleGrantItemRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	repository, err := h.commandRepository(ctx)
	if err != nil {
		return resultForRequest(request), err
	}
	return h.inventoryHandlers(repository).HandleGrantItemRoute(ctx, request)
}

func (h InventoryRouteHandlers) HandleGetInventoryRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	repository, err := h.queryRepository(ctx)
	if err != nil {
		return resultForRequest(request), err
	}
	return h.inventoryHandlers(repository).HandleGetInventoryRoute(ctx, request)
}

func (h InventoryRouteHandlers) commandRepository(ctx context.Context) (inventory.Repository, error) {
	if h.Repositories == nil {
		return nil, errors.New("inventory bootstrap: repository provider is required")
	}

	unit, ok := app.UnitOfWorkFromContext(ctx)
	if !ok {
		unit = tx.NoopUnitOfWork{}
	}
	return h.Repositories.ForCommand(ctx, unit)
}

func (h InventoryRouteHandlers) queryRepository(ctx context.Context) (inventory.Repository, error) {
	if h.Repositories == nil {
		return nil, errors.New("inventory bootstrap: repository provider is required")
	}
	return h.Repositories.ForQuery(ctx)
}

func (h InventoryRouteHandlers) inventoryHandlers(repository inventory.Repository) inventory.Handlers {
	return inventory.Handlers{
		Repository:       repository,
		PermissionPolicy: h.PermissionPolicy,
		CapacityPolicy:   h.CapacityPolicy,
		EventIDs:         h.EventIDs,
		Clock:            h.Clock,
	}
}

func resultForRequest(request app.RouteRequest) app.ApplicationResult {
	return app.ApplicationResult{
		RequestID: request.RequestID,
		Route:     request.Route,
		Target:    request.Target,
		Session:   request.Session,
	}
}
