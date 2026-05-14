package inventory

import (
	"context"
	"errors"
	"sync"

	"github.com/iceiko/vibit/runtime/internal/app"
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

	return r.getInventoryLocked(playerID), nil
}

func (r *MemoryRepository) LockInventoryForMutation(_ context.Context, playerID string) (MutationLock, error) {
	r.mu.Lock()
	return &memoryMutationLock{
		repository: r,
		playerID:   playerID,
	}, nil
}

func (r *MemoryRepository) getInventoryLocked(playerID string) []Item {
	playerItems := r.items[playerID]
	items := make([]Item, 0, len(playerItems))
	for itemID, quantity := range playerItems {
		items = append(items, Item{ItemID: itemID, Quantity: quantity})
	}
	return items
}

func (r *MemoryRepository) grantItemLocked(mutation GrantItemMutation) Item {
	if r.items[mutation.PlayerID] == nil {
		r.items[mutation.PlayerID] = make(map[string]int64)
	}
	r.items[mutation.PlayerID][mutation.ItemID] += mutation.Quantity
	return Item{
		ItemID:   mutation.ItemID,
		Quantity: r.items[mutation.PlayerID][mutation.ItemID],
	}
}

type memoryMutationLock struct {
	repository *MemoryRepository
	playerID   string
	released   bool
}

func (l *memoryMutationLock) GetInventory(_ context.Context, playerID string) ([]Item, error) {
	if err := l.ensureUsable(playerID); err != nil {
		return nil, err
	}
	return l.repository.getInventoryLocked(playerID), nil
}

func (l *memoryMutationLock) GrantItem(_ context.Context, mutation GrantItemMutation) (Item, error) {
	if err := l.ensureUsable(mutation.PlayerID); err != nil {
		return Item{}, err
	}
	return l.repository.grantItemLocked(mutation), nil
}

func (l *memoryMutationLock) Release() {
	if l == nil || l.released {
		return
	}
	l.released = true
	l.repository.mu.Unlock()
}

func (l *memoryMutationLock) ensureUsable(playerID string) error {
	if l == nil || l.repository == nil {
		return errors.New("inventory: mutation lock is not initialized")
	}
	if l.released {
		return errors.New("inventory: mutation lock was released")
	}
	if playerID != l.playerID {
		return errors.New("inventory: mutation lock player_id mismatch")
	}
	return nil
}

type StaticPermissionPolicy struct {
	GrantAllowed bool
	ReadAllowed  bool
}

func (p StaticPermissionPolicy) CanGrantItem(context.Context, PermissionContext) (bool, error) {
	return p.GrantAllowed, nil
}

func (p StaticPermissionPolicy) CanReadInventory(context.Context, PermissionContext) (bool, error) {
	return p.ReadAllowed, nil
}

type MetadataOnlyDenyPermissionPolicy struct {
	AllowValidatedPlayerSelfRead bool
}

func (p MetadataOnlyDenyPermissionPolicy) CanGrantItem(_ context.Context, ctx PermissionContext) (bool, error) {
	if !ctx.Identity.PlayerIDValidated || ctx.Identity.Status != app.IdentityValidationValidated {
		return false, nil
	}
	return false, nil
}

func (p MetadataOnlyDenyPermissionPolicy) CanReadInventory(_ context.Context, ctx PermissionContext) (bool, error) {
	if !p.AllowValidatedPlayerSelfRead {
		return false, nil
	}
	if ctx.Identity.Status != app.IdentityValidationValidated || !ctx.Identity.PlayerIDValidated {
		return false, nil
	}
	return ctx.Identity.PlayerID == ctx.PlayerID && ctx.PlayerID != "", nil
}
