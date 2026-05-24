-- +goose Up
-- Module: friends
-- Purpose: Create durable current-state friend relationship records.

CREATE TABLE friend_relationships (
    relationship_id TEXT PRIMARY KEY,
    player_low_id TEXT NOT NULL,
    player_high_id TEXT NOT NULL,
    lifecycle_state TEXT NOT NULL,
    requested_by_player_id TEXT,
    responded_by_player_id TEXT,
    removed_by_player_id TEXT,
    relationship_version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    state_changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    rejected_at TIMESTAMPTZ,
    removed_at TIMESTAMPTZ,
    blocked_by_low_at TIMESTAMPTZ,
    blocked_by_high_at TIMESTAMPTZ,
    CONSTRAINT friend_relationships_low_player_fk
        FOREIGN KEY (player_low_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT friend_relationships_high_player_fk
        FOREIGN KEY (player_high_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT friend_relationships_requested_by_fk
        FOREIGN KEY (requested_by_player_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT friend_relationships_responded_by_fk
        FOREIGN KEY (responded_by_player_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT friend_relationships_removed_by_fk
        FOREIGN KEY (removed_by_player_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT friend_relationships_relationship_id_not_blank CHECK (length(btrim(relationship_id)) > 0),
    CONSTRAINT friend_relationships_player_low_id_not_blank CHECK (length(btrim(player_low_id)) > 0),
    CONSTRAINT friend_relationships_player_high_id_not_blank CHECK (length(btrim(player_high_id)) > 0),
    CONSTRAINT friend_relationships_canonical_pair_order CHECK (player_low_id < player_high_id),
    CONSTRAINT friend_relationships_lifecycle_state_valid CHECK (
        lifecycle_state IN ('pending', 'friends', 'rejected', 'removed')
    ),
    CONSTRAINT friend_relationships_requested_by_pair_member CHECK (
        requested_by_player_id IS NULL OR
        requested_by_player_id = player_low_id OR
        requested_by_player_id = player_high_id
    ),
    CONSTRAINT friend_relationships_responded_by_pair_member CHECK (
        responded_by_player_id IS NULL OR
        responded_by_player_id = player_low_id OR
        responded_by_player_id = player_high_id
    ),
    CONSTRAINT friend_relationships_removed_by_pair_member CHECK (
        removed_by_player_id IS NULL OR
        removed_by_player_id = player_low_id OR
        removed_by_player_id = player_high_id
    ),
    CONSTRAINT friend_relationships_relationship_version_positive CHECK (relationship_version > 0),
    CONSTRAINT friend_relationships_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT friend_relationships_state_changed_after_created CHECK (state_changed_at >= created_at),
    CONSTRAINT friend_relationships_rejected_after_created CHECK (
        rejected_at IS NULL OR rejected_at >= created_at
    ),
    CONSTRAINT friend_relationships_removed_after_created CHECK (
        removed_at IS NULL OR removed_at >= created_at
    ),
    CONSTRAINT friend_relationships_blocked_by_low_after_created CHECK (
        blocked_by_low_at IS NULL OR blocked_by_low_at >= created_at
    ),
    CONSTRAINT friend_relationships_blocked_by_high_after_created CHECK (
        blocked_by_high_at IS NULL OR blocked_by_high_at >= created_at
    )
);

COMMENT ON TABLE friend_relationships IS
    'Current-state friendship relationship rows for canonical unordered player pairs; excludes secret material, digest material, transport metadata, chat, group, party, match, distributed-runtime, or external compatibility columns.';

CREATE UNIQUE INDEX friend_relationships_pair_uq
    ON friend_relationships (player_low_id, player_high_id);

CREATE INDEX friend_relationships_player_low_state_idx
    ON friend_relationships (player_low_id, lifecycle_state);

CREATE INDEX friend_relationships_player_high_state_idx
    ON friend_relationships (player_high_id, lifecycle_state);

CREATE INDEX friend_relationships_updated_at_idx
    ON friend_relationships (updated_at);

-- +goose Down
DROP INDEX IF EXISTS friend_relationships_updated_at_idx;
DROP INDEX IF EXISTS friend_relationships_player_high_state_idx;
DROP INDEX IF EXISTS friend_relationships_player_low_state_idx;
DROP INDEX IF EXISTS friend_relationships_pair_uq;
DROP TABLE IF EXISTS friend_relationships;
