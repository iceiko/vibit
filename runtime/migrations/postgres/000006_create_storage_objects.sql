-- +goose Up
-- Module: storage
-- Purpose: Create player-owned small JSON storage object records.

CREATE TABLE storage_objects (
    object_id TEXT PRIMARY KEY,
    owner_kind TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    collection TEXT NOT NULL,
    object_key TEXT NOT NULL,
    value_json JSONB NOT NULL,
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT storage_objects_owner_player_fk
        FOREIGN KEY (owner_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT storage_objects_object_id_not_blank CHECK (length(btrim(object_id)) > 0),
    CONSTRAINT storage_objects_owner_kind_valid CHECK (owner_kind = 'player'),
    CONSTRAINT storage_objects_owner_id_not_blank CHECK (length(btrim(owner_id)) > 0),
    CONSTRAINT storage_objects_collection_not_blank CHECK (length(btrim(collection)) > 0),
    CONSTRAINT storage_objects_collection_length CHECK (length(collection) <= 128),
    CONSTRAINT storage_objects_object_key_not_blank CHECK (length(btrim(object_key)) > 0),
    CONSTRAINT storage_objects_object_key_length CHECK (length(object_key) <= 256),
    CONSTRAINT storage_objects_value_json_object CHECK (jsonb_typeof(value_json) = 'object'),
    CONSTRAINT storage_objects_version_positive CHECK (version > 0),
    CONSTRAINT storage_objects_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT storage_objects_deleted_after_created CHECK (
        deleted_at IS NULL OR deleted_at >= created_at
    )
);

COMMENT ON TABLE storage_objects IS
    'Player-owned small JSON storage object state; excludes secret material, digest material, transport metadata, blob bytes, file paths, or S3 object references.';

CREATE UNIQUE INDEX storage_objects_active_identity_uq
    ON storage_objects (owner_kind, owner_id, collection, object_key)
    WHERE deleted_at IS NULL;

CREATE INDEX storage_objects_owner_collection_idx
    ON storage_objects (owner_kind, owner_id, collection);

CREATE INDEX storage_objects_updated_at_idx
    ON storage_objects (updated_at);

-- +goose Down
DROP INDEX IF EXISTS storage_objects_updated_at_idx;
DROP INDEX IF EXISTS storage_objects_owner_collection_idx;
DROP INDEX IF EXISTS storage_objects_active_identity_uq;
DROP TABLE IF EXISTS storage_objects;
