package bootstrap

import (
	"github.com/iceiko/vibit/runtime/internal/app"
	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
)

func NewInMemoryInventoryDispatcher() (*app.Dispatcher, error) {
	dispatcher := app.NewDispatcher()
	handlers := inventory.Handlers{
		Repository:       inventory.NewMemoryRepository(),
		PermissionPolicy: inventory.StaticPermissionPolicy{GrantAllowed: true, ReadAllowed: true},
		CapacityPolicy:   inventory.MaxUniqueItemsCapacityPolicy{MaxUniqueItems: 256},
		EventIDs:         &inventory.IncrementingEventIDGenerator{Prefix: "inventory-event"},
		Clock:            inventory.SystemClock{},
	}
	if err := handlers.RegisterRoutes(dispatcher); err != nil {
		return nil, err
	}
	return dispatcher, nil
}
