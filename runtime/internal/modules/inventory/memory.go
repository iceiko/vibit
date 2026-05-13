package inventory

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu    sync.Mutex
	items map[string]map[string]int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: make(map[string]map[string]int64),
	}
}

func (r *MemoryRepository) GetInventory(_ context.Context, playerID string) ([]Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	playerItems := r.items[playerID]
	items := make([]Item, 0, len(playerItems))
	for itemID, quantity := range playerItems {
		items = append(items, Item{ItemID: itemID, Quantity: quantity})
	}
	return items, nil
}

func (r *MemoryRepository) GrantItem(_ context.Context, mutation GrantItemMutation) (Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.items[mutation.PlayerID] == nil {
		r.items[mutation.PlayerID] = make(map[string]int64)
	}
	r.items[mutation.PlayerID][mutation.ItemID] += mutation.Quantity
	return Item{
		ItemID:   mutation.ItemID,
		Quantity: r.items[mutation.PlayerID][mutation.ItemID],
	}, nil
}

type StaticPermissionPolicy struct {
	GrantAllowed bool
	ReadAllowed  bool
}

func (p StaticPermissionPolicy) CanGrantItem(context.Context, string, string) (bool, error) {
	return p.GrantAllowed, nil
}

func (p StaticPermissionPolicy) CanReadInventory(context.Context, string, string) (bool, error) {
	return p.ReadAllowed, nil
}
