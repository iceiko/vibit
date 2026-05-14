-- +goose Up
-- Module: runtime.authentication
-- Purpose: Create device credential verifier records for the selected first login posture.

CREATE TABLE authentication_device_credentials (
    credential_record_id TEXT PRIMARY KEY,
    player_id TEXT NOT NULL,
    credential_kind TEXT NOT NULL,
    credential_status TEXT NOT NULL,
    credential_lookup_digest BYTEA NOT NULL,
    credential_verifier_digest BYTEA NOT NULL,
    verifier_algorithm TEXT NOT NULL,
    verifier_version INTEGER NOT NULL,
    verifier_key_id TEXT NOT NULL,
    client_instance_id_digest BYTEA,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_verified_at TIMESTAMPTZ,
    disabled_at TIMESTAMPTZ,
    disabled_reason TEXT,
    revoked_at TIMESTAMPTZ,
    revoked_reason TEXT,
    replaced_at TIMESTAMPTZ,
    replaced_by_credential_record_id TEXT,
    CONSTRAINT authentication_device_credentials_player_fk
        FOREIGN KEY (player_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT authentication_device_credentials_replaced_by_fk
        FOREIGN KEY (replaced_by_credential_record_id)
        REFERENCES authentication_device_credentials(credential_record_id)
        ON DELETE RESTRICT,
    CONSTRAINT authentication_device_credentials_record_id_not_blank CHECK (length(btrim(credential_record_id)) > 0),
    CONSTRAINT authentication_device_credentials_player_id_not_blank CHECK (length(btrim(player_id)) > 0),
    CONSTRAINT authentication_device_credentials_kind_valid CHECK (credential_kind = 'device_credential_login'),
    CONSTRAINT authentication_device_credentials_status_valid CHECK (credential_status IN ('active', 'disabled', 'revoked', 'replaced')),
    CONSTRAINT authentication_device_credentials_lookup_digest_present CHECK (octet_length(credential_lookup_digest) > 0),
    CONSTRAINT authentication_device_credentials_verifier_digest_present CHECK (octet_length(credential_verifier_digest) > 0),
    CONSTRAINT authentication_device_credentials_verifier_algorithm_not_blank CHECK (length(btrim(verifier_algorithm)) > 0),
    CONSTRAINT authentication_device_credentials_verifier_version_positive CHECK (verifier_version > 0),
    CONSTRAINT authentication_device_credentials_verifier_key_id_not_blank CHECK (length(btrim(verifier_key_id)) > 0),
    CONSTRAINT authentication_device_credentials_client_instance_digest_present CHECK (
        client_instance_id_digest IS NULL OR octet_length(client_instance_id_digest) > 0
    ),
    CONSTRAINT authentication_device_credentials_disabled_reason_not_blank CHECK (
        disabled_reason IS NULL OR length(btrim(disabled_reason)) > 0
    ),
    CONSTRAINT authentication_device_credentials_revoked_reason_not_blank CHECK (
        revoked_reason IS NULL OR length(btrim(revoked_reason)) > 0
    ),
    CONSTRAINT authentication_device_credentials_disabled_state_matches_time CHECK (
        disabled_at IS NULL OR credential_status = 'disabled'
    ),
    CONSTRAINT authentication_device_credentials_disabled_reason_matches_time CHECK (
        disabled_reason IS NULL OR disabled_at IS NOT NULL
    ),
    CONSTRAINT authentication_device_credentials_revoked_state_matches_time CHECK (
        revoked_at IS NULL OR credential_status = 'revoked'
    ),
    CONSTRAINT authentication_device_credentials_revoked_reason_matches_time CHECK (
        revoked_reason IS NULL OR revoked_at IS NOT NULL
    ),
    CONSTRAINT authentication_device_credentials_replaced_state_matches_time CHECK (
        replaced_at IS NULL OR credential_status = 'replaced'
    ),
    CONSTRAINT authentication_device_credentials_replaced_by_matches_state CHECK (
        replaced_by_credential_record_id IS NULL OR credential_status = 'replaced'
    ),
    CONSTRAINT authentication_device_credentials_replaced_has_target CHECK (
        credential_status <> 'replaced' OR replaced_by_credential_record_id IS NOT NULL
    ),
    CONSTRAINT authentication_device_credentials_replaced_not_self CHECK (
        replaced_by_credential_record_id IS NULL OR replaced_by_credential_record_id <> credential_record_id
    ),
    CONSTRAINT authentication_device_credentials_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT authentication_device_credentials_last_verified_after_created CHECK (
        last_verified_at IS NULL OR last_verified_at >= created_at
    ),
    CONSTRAINT authentication_device_credentials_disabled_after_created CHECK (
        disabled_at IS NULL OR disabled_at >= created_at
    ),
    CONSTRAINT authentication_device_credentials_revoked_after_created CHECK (
        revoked_at IS NULL OR revoked_at >= created_at
    ),
    CONSTRAINT authentication_device_credentials_replaced_after_created CHECK (
        replaced_at IS NULL OR replaced_at >= created_at
    )
);

COMMENT ON TABLE authentication_device_credentials IS
    'Device credential verifier records owned by runtime.authentication; stores non-plaintext lookup and verifier digests only.';

CREATE UNIQUE INDEX authentication_device_credentials_lookup_digest_uq
    ON authentication_device_credentials (credential_lookup_digest);

CREATE UNIQUE INDEX authentication_device_credentials_one_active_per_player_uq
    ON authentication_device_credentials (player_id)
    WHERE credential_status = 'active';

CREATE INDEX authentication_device_credentials_player_status_idx
    ON authentication_device_credentials (player_id, credential_status);

CREATE INDEX authentication_device_credentials_status_updated_at_idx
    ON authentication_device_credentials (credential_status, updated_at);

CREATE INDEX authentication_device_credentials_replaced_by_idx
    ON authentication_device_credentials (replaced_by_credential_record_id)
    WHERE replaced_by_credential_record_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS authentication_device_credentials_replaced_by_idx;
DROP INDEX IF EXISTS authentication_device_credentials_status_updated_at_idx;
DROP INDEX IF EXISTS authentication_device_credentials_player_status_idx;
DROP INDEX IF EXISTS authentication_device_credentials_one_active_per_player_uq;
DROP INDEX IF EXISTS authentication_device_credentials_lookup_digest_uq;
DROP TABLE IF EXISTS authentication_device_credentials;
