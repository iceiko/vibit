-- +goose Up
-- Module: player
-- Purpose: Create player account lifecycle state and lifecycle event records.

CREATE TABLE player_accounts (
    player_id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    account_state TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    disabled_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT player_accounts_player_id_not_blank CHECK (length(btrim(player_id)) > 0),
    CONSTRAINT player_accounts_display_name_not_blank CHECK (length(btrim(display_name)) > 0),
    CONSTRAINT player_accounts_account_state_valid CHECK (account_state IN ('active', 'disabled', 'deleted')),
    CONSTRAINT player_accounts_disabled_at_matches_state CHECK (
        disabled_at IS NULL OR account_state IN ('disabled', 'deleted')
    ),
    CONSTRAINT player_accounts_deleted_at_matches_state CHECK (
        deleted_at IS NULL OR account_state = 'deleted'
    ),
    CONSTRAINT player_accounts_deleted_after_created CHECK (
        deleted_at IS NULL OR deleted_at >= created_at
    ),
    CONSTRAINT player_accounts_disabled_after_created CHECK (
        disabled_at IS NULL OR disabled_at >= created_at
    ),
    CONSTRAINT player_accounts_updated_after_created CHECK (updated_at >= created_at)
);

CREATE TABLE player_account_events (
    event_id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    player_id TEXT NOT NULL,
    requested_by TEXT NOT NULL,
    account_state TEXT NOT NULL,
    display_name TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT player_account_events_player_fk
        FOREIGN KEY (player_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT player_account_events_event_id_not_blank CHECK (length(btrim(event_id)) > 0),
    CONSTRAINT player_account_events_event_type_not_blank CHECK (length(btrim(event_type)) > 0),
    CONSTRAINT player_account_events_player_id_not_blank CHECK (length(btrim(player_id)) > 0),
    CONSTRAINT player_account_events_requested_by_not_blank CHECK (length(btrim(requested_by)) > 0),
    CONSTRAINT player_account_events_account_state_valid CHECK (account_state IN ('active', 'disabled', 'deleted')),
    CONSTRAINT player_account_events_display_name_not_blank CHECK (length(btrim(display_name)) > 0),
    CONSTRAINT player_account_events_metadata_object CHECK (jsonb_typeof(metadata) = 'object')
);

COMMENT ON TABLE player_account_events IS
    'Player account lifecycle event records; first required event type: PlayerAccountCreated.';

CREATE INDEX player_accounts_account_state_idx
    ON player_accounts (account_state);

CREATE INDEX player_account_events_player_occurred_at_idx
    ON player_account_events (player_id, occurred_at);

CREATE INDEX player_account_events_event_type_occurred_at_idx
    ON player_account_events (event_type, occurred_at);

-- +goose Down
DROP INDEX IF EXISTS player_account_events_event_type_occurred_at_idx;
DROP INDEX IF EXISTS player_account_events_player_occurred_at_idx;
DROP INDEX IF EXISTS player_accounts_account_state_idx;
DROP TABLE IF EXISTS player_account_events;
DROP TABLE IF EXISTS player_accounts;
