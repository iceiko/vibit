-- +goose Up
-- Module: runtime.authentication
-- Purpose: Create opaque access-token verifier records for the selected first token posture.

CREATE TABLE authentication_access_tokens (
    token_record_id TEXT PRIMARY KEY,
    token_kind TEXT NOT NULL,
    token_status TEXT NOT NULL,
    actor_kind TEXT NOT NULL,
    player_id TEXT NOT NULL,
    credential_record_id TEXT NOT NULL,
    token_lookup_digest BYTEA NOT NULL,
    token_verifier_digest BYTEA NOT NULL,
    verifier_algorithm TEXT NOT NULL,
    verifier_version INTEGER NOT NULL,
    verifier_key_id TEXT NOT NULL,
    audience TEXT NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    revoked_reason TEXT,
    replaced_by_token_record_id TEXT,
    last_validated_at TIMESTAMPTZ,
    last_failed_validation_at TIMESTAMPTZ,
    cleanup_after TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT authentication_access_tokens_player_fk
        FOREIGN KEY (player_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT authentication_access_tokens_credential_fk
        FOREIGN KEY (credential_record_id)
        REFERENCES authentication_device_credentials(credential_record_id)
        ON DELETE RESTRICT,
    CONSTRAINT authentication_access_tokens_replaced_by_fk
        FOREIGN KEY (replaced_by_token_record_id)
        REFERENCES authentication_access_tokens(token_record_id)
        ON DELETE RESTRICT,
    CONSTRAINT authentication_access_tokens_record_id_not_blank CHECK (length(btrim(token_record_id)) > 0),
    CONSTRAINT authentication_access_tokens_kind_valid CHECK (token_kind = 'access_token'),
    CONSTRAINT authentication_access_tokens_status_valid CHECK (token_status IN ('active', 'expired', 'revoked')),
    CONSTRAINT authentication_access_tokens_actor_kind_valid CHECK (actor_kind = 'player'),
    CONSTRAINT authentication_access_tokens_player_id_not_blank CHECK (length(btrim(player_id)) > 0),
    CONSTRAINT authentication_access_tokens_credential_id_not_blank CHECK (length(btrim(credential_record_id)) > 0),
    CONSTRAINT authentication_access_tokens_lookup_digest_present CHECK (octet_length(token_lookup_digest) > 0),
    CONSTRAINT authentication_access_tokens_verifier_digest_present CHECK (octet_length(token_verifier_digest) > 0),
    CONSTRAINT authentication_access_tokens_verifier_algorithm_not_blank CHECK (length(btrim(verifier_algorithm)) > 0),
    CONSTRAINT authentication_access_tokens_verifier_version_positive CHECK (verifier_version > 0),
    CONSTRAINT authentication_access_tokens_verifier_key_id_not_blank CHECK (length(btrim(verifier_key_id)) > 0),
    CONSTRAINT authentication_access_tokens_audience_not_blank CHECK (length(btrim(audience)) > 0),
    CONSTRAINT authentication_access_tokens_expires_after_issued CHECK (expires_at > issued_at),
    CONSTRAINT authentication_access_tokens_revoked_state_matches_time CHECK (
        revoked_at IS NULL OR token_status = 'revoked'
    ),
    CONSTRAINT authentication_access_tokens_revoked_reason_not_blank CHECK (
        revoked_reason IS NULL OR length(btrim(revoked_reason)) > 0
    ),
    CONSTRAINT authentication_access_tokens_revoked_reason_matches_time CHECK (
        revoked_reason IS NULL OR revoked_at IS NOT NULL
    ),
    CONSTRAINT authentication_access_tokens_replaced_by_matches_revoked_state CHECK (
        replaced_by_token_record_id IS NULL OR token_status = 'revoked'
    ),
    CONSTRAINT authentication_access_tokens_replaced_not_self CHECK (
        replaced_by_token_record_id IS NULL OR replaced_by_token_record_id <> token_record_id
    ),
    CONSTRAINT authentication_access_tokens_cleanup_after_terminal_state CHECK (
        cleanup_after IS NULL OR token_status IN ('expired', 'revoked')
    ),
    CONSTRAINT authentication_access_tokens_cleanup_after_issued CHECK (
        cleanup_after IS NULL OR cleanup_after >= issued_at
    ),
    CONSTRAINT authentication_access_tokens_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT authentication_access_tokens_last_validated_after_issued CHECK (
        last_validated_at IS NULL OR last_validated_at >= issued_at
    ),
    CONSTRAINT authentication_access_tokens_last_failed_after_issued CHECK (
        last_failed_validation_at IS NULL OR last_failed_validation_at >= issued_at
    ),
    CONSTRAINT authentication_access_tokens_revoked_after_issued CHECK (
        revoked_at IS NULL OR revoked_at >= issued_at
    )
);

COMMENT ON TABLE authentication_access_tokens IS
    'Opaque access-token verifier records owned by runtime.authentication; stores non-plaintext lookup and verifier digests only.';

CREATE UNIQUE INDEX authentication_access_tokens_lookup_digest_uq
    ON authentication_access_tokens (token_lookup_digest);

CREATE INDEX authentication_access_tokens_player_status_idx
    ON authentication_access_tokens (player_id, token_status);

CREATE INDEX authentication_access_tokens_credential_status_idx
    ON authentication_access_tokens (credential_record_id, token_status);

CREATE INDEX authentication_access_tokens_status_expires_at_idx
    ON authentication_access_tokens (token_status, expires_at);

CREATE INDEX authentication_access_tokens_cleanup_after_idx
    ON authentication_access_tokens (cleanup_after)
    WHERE cleanup_after IS NOT NULL;

CREATE INDEX authentication_access_tokens_replaced_by_idx
    ON authentication_access_tokens (replaced_by_token_record_id)
    WHERE replaced_by_token_record_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS authentication_access_tokens_replaced_by_idx;
DROP INDEX IF EXISTS authentication_access_tokens_cleanup_after_idx;
DROP INDEX IF EXISTS authentication_access_tokens_status_expires_at_idx;
DROP INDEX IF EXISTS authentication_access_tokens_credential_status_idx;
DROP INDEX IF EXISTS authentication_access_tokens_player_status_idx;
DROP INDEX IF EXISTS authentication_access_tokens_lookup_digest_uq;
DROP TABLE IF EXISTS authentication_access_tokens;
