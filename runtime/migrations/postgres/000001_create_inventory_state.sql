-- +goose Up
-- Module: inventory
-- Purpose: Create authoritative inventory aggregate, quantity, and grant-record state.

CREATE TABLE inventory_accounts (
    player_id TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT inventory_accounts_player_id_not_blank CHECK (length(btrim(player_id)) > 0)
);

CREATE TABLE inventory_items (
    player_id TEXT NOT NULL,
    item_id TEXT NOT NULL,
    quantity BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (player_id, item_id),
    CONSTRAINT inventory_items_player_fk
        FOREIGN KEY (player_id)
        REFERENCES inventory_accounts(player_id)
        ON DELETE CASCADE,
    CONSTRAINT inventory_items_item_id_not_blank CHECK (length(btrim(item_id)) > 0),
    CONSTRAINT inventory_items_quantity_positive CHECK (quantity > 0)
);

CREATE TABLE inventory_item_grants (
    event_id TEXT PRIMARY KEY,
    occurred_at TIMESTAMPTZ NOT NULL,
    player_id TEXT NOT NULL,
    item_id TEXT NOT NULL,
    quantity BIGINT NOT NULL,
    new_quantity BIGINT NOT NULL,
    reason TEXT NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT inventory_item_grants_player_fk
        FOREIGN KEY (player_id)
        REFERENCES inventory_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT inventory_item_grants_player_id_not_blank CHECK (length(btrim(player_id)) > 0),
    CONSTRAINT inventory_item_grants_item_id_not_blank CHECK (length(btrim(item_id)) > 0),
    CONSTRAINT inventory_item_grants_event_id_not_blank CHECK (length(btrim(event_id)) > 0),
    CONSTRAINT inventory_item_grants_reason_not_blank CHECK (length(btrim(reason)) > 0),
    CONSTRAINT inventory_item_grants_quantity_positive CHECK (quantity > 0),
    CONSTRAINT inventory_item_grants_new_quantity_positive CHECK (new_quantity > 0),
    CONSTRAINT inventory_item_grants_new_quantity_covers_grant CHECK (new_quantity >= quantity)
);

CREATE INDEX inventory_item_grants_player_occurred_at_idx
    ON inventory_item_grants (player_id, occurred_at);

CREATE INDEX inventory_item_grants_item_id_idx
    ON inventory_item_grants (item_id);

-- +goose Down
DROP INDEX IF EXISTS inventory_item_grants_item_id_idx;
DROP INDEX IF EXISTS inventory_item_grants_player_occurred_at_idx;
DROP TABLE IF EXISTS inventory_item_grants;
DROP TABLE IF EXISTS inventory_items;
DROP TABLE IF EXISTS inventory_accounts;
