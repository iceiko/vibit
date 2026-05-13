package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/modules/inventory"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Executor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type InventoryRepository struct {
	executor Executor
}

func NewInventoryRepositoryForUnitOfWork(executor Executor) *InventoryRepository {
	return &InventoryRepository{executor: executor}
}

func (r *InventoryRepository) GetInventory(ctx context.Context, playerID string) ([]inventory.Item, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return nil, err
	}
	playerID, err = normalizeRequired("player_id", playerID)
	if err != nil {
		return nil, err
	}

	rows, err := executor.Query(ctx, getInventorySQL, playerID)
	if err != nil {
		return nil, fmt.Errorf("postgres inventory: get inventory: %w", err)
	}
	defer rows.Close()

	items := []inventory.Item{}
	for rows.Next() {
		var item inventory.Item
		if err := rows.Scan(&item.ItemID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("postgres inventory: scan inventory item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres inventory: read inventory items: %w", err)
	}
	return items, nil
}

func (r *InventoryRepository) LockInventoryForMutation(ctx context.Context, playerID string) (inventory.MutationLock, error) {
	executor, err := r.requireExecutor()
	if err != nil {
		return nil, err
	}
	playerID, err = normalizeRequired("player_id", playerID)
	if err != nil {
		return nil, err
	}

	if _, err := executor.Exec(ctx, ensureInventoryAccountSQL, playerID); err != nil {
		return nil, fmt.Errorf("postgres inventory: ensure inventory account: %w", err)
	}

	var lockedPlayerID string
	if err := executor.QueryRow(ctx, lockInventoryAccountSQL, playerID).Scan(&lockedPlayerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("postgres inventory: inventory account was not created for %q", playerID)
		}
		return nil, fmt.Errorf("postgres inventory: lock inventory account: %w", err)
	}
	if lockedPlayerID != playerID {
		return nil, fmt.Errorf("postgres inventory: locked player_id mismatch")
	}

	return &inventoryMutationLock{
		repository: r,
		playerID:   playerID,
	}, nil
}

func (r *InventoryRepository) requireExecutor() (Executor, error) {
	if r == nil || r.executor == nil {
		return nil, errors.New("postgres inventory: unit-of-work executor is required")
	}
	return r.executor, nil
}

type inventoryMutationLock struct {
	repository *InventoryRepository
	playerID   string
	released   bool
}

func (l *inventoryMutationLock) GetInventory(ctx context.Context, playerID string) ([]inventory.Item, error) {
	if err := l.ensureUsable(playerID); err != nil {
		return nil, err
	}
	return l.repository.GetInventory(ctx, playerID)
}

func (l *inventoryMutationLock) GrantItem(ctx context.Context, mutation inventory.GrantItemMutation) (inventory.Item, error) {
	if err := l.ensureUsable(mutation.PlayerID); err != nil {
		return inventory.Item{}, err
	}

	normalized, err := normalizeGrantMutation(mutation)
	if err != nil {
		return inventory.Item{}, err
	}

	executor, err := l.repository.requireExecutor()
	if err != nil {
		return inventory.Item{}, err
	}

	var granted inventory.Item
	if err := executor.QueryRow(
		ctx,
		upsertInventoryItemSQL,
		normalized.PlayerID,
		normalized.ItemID,
		normalized.Quantity,
	).Scan(&granted.ItemID, &granted.Quantity); err != nil {
		return inventory.Item{}, fmt.Errorf("postgres inventory: grant item quantity: %w", err)
	}
	if granted.ItemID != normalized.ItemID {
		return inventory.Item{}, errors.New("postgres inventory: granted item_id mismatch")
	}
	if granted.Quantity < normalized.Quantity {
		return inventory.Item{}, errors.New("postgres inventory: grant returned invalid quantity")
	}

	if _, err := executor.Exec(
		ctx,
		recordInventoryItemGrantSQL,
		normalized.EventID,
		normalized.OccurredAt,
		normalized.PlayerID,
		normalized.ItemID,
		normalized.Quantity,
		granted.Quantity,
		normalized.Reason,
	); err != nil {
		return inventory.Item{}, fmt.Errorf("postgres inventory: record item grant: %w", err)
	}

	return granted, nil
}

func (l *inventoryMutationLock) Release() {
	if l == nil || l.released {
		return
	}
	// Transaction commit and rollback are owned by the application unit of work.
	l.released = true
}

func (l *inventoryMutationLock) ensureUsable(playerID string) error {
	if l == nil || l.repository == nil {
		return errors.New("postgres inventory: mutation lock is not initialized")
	}
	if l.released {
		return errors.New("postgres inventory: mutation lock was released")
	}
	playerID, err := normalizeRequired("player_id", playerID)
	if err != nil {
		return err
	}
	if playerID != l.playerID {
		return errors.New("postgres inventory: mutation lock player_id mismatch")
	}
	return nil
}

func normalizeGrantMutation(mutation inventory.GrantItemMutation) (inventory.GrantItemMutation, error) {
	var err error
	mutation.EventID, err = normalizeRequired("event_id", mutation.EventID)
	if err != nil {
		return inventory.GrantItemMutation{}, err
	}
	mutation.PlayerID, err = normalizeRequired("player_id", mutation.PlayerID)
	if err != nil {
		return inventory.GrantItemMutation{}, err
	}
	mutation.ItemID, err = normalizeRequired("item_id", mutation.ItemID)
	if err != nil {
		return inventory.GrantItemMutation{}, err
	}
	mutation.Reason, err = normalizeRequired("reason", mutation.Reason)
	if err != nil {
		return inventory.GrantItemMutation{}, err
	}
	if mutation.Quantity <= 0 {
		return inventory.GrantItemMutation{}, errors.New("postgres inventory: quantity must be positive")
	}
	if mutation.OccurredAt.IsZero() {
		return inventory.GrantItemMutation{}, errors.New("postgres inventory: occurred_at is required")
	}
	mutation.OccurredAt = mutation.OccurredAt.UTC()
	return mutation, nil
}

func normalizeRequired(name string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("postgres inventory: %s is required", name)
	}
	return value, nil
}

const getInventorySQL = `
SELECT item_id, quantity
FROM inventory_items
WHERE player_id = $1
ORDER BY item_id`

const ensureInventoryAccountSQL = `
INSERT INTO inventory_accounts (player_id)
VALUES ($1)
ON CONFLICT (player_id) DO NOTHING`

const lockInventoryAccountSQL = `
SELECT player_id
FROM inventory_accounts
WHERE player_id = $1
FOR UPDATE`

const upsertInventoryItemSQL = `
INSERT INTO inventory_items (player_id, item_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (player_id, item_id)
DO UPDATE SET
    quantity = inventory_items.quantity + EXCLUDED.quantity,
    updated_at = now()
RETURNING item_id, quantity`

const recordInventoryItemGrantSQL = `
INSERT INTO inventory_item_grants (
    event_id,
    occurred_at,
    player_id,
    item_id,
    quantity,
    new_quantity,
    reason
)
VALUES ($1, $2, $3, $4, $5, $6, $7)`
