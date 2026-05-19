-- +goose Up
-- Module: runtime.session
-- Purpose: Create durable runtime session lifecycle records for authenticated player sessions.

CREATE TABLE runtime_sessions (
    session_id TEXT PRIMARY KEY,
    actor_kind TEXT NOT NULL,
    actor_id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    session_status TEXT NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    revocation_reason TEXT,
    access_token_record_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT runtime_sessions_player_fk
        FOREIGN KEY (player_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT runtime_sessions_access_token_fk
        FOREIGN KEY (access_token_record_id)
        REFERENCES authentication_access_tokens(token_record_id)
        ON DELETE RESTRICT,
    CONSTRAINT runtime_sessions_session_id_not_blank CHECK (length(btrim(session_id)) > 0),
    CONSTRAINT runtime_sessions_actor_kind_valid CHECK (actor_kind = 'player'),
    CONSTRAINT runtime_sessions_actor_id_not_blank CHECK (length(btrim(actor_id)) > 0),
    CONSTRAINT runtime_sessions_player_id_not_blank CHECK (length(btrim(player_id)) > 0),
    CONSTRAINT runtime_sessions_actor_id_matches_player CHECK (actor_id = player_id),
    CONSTRAINT runtime_sessions_status_valid CHECK (session_status IN ('active', 'expired', 'revoked')),
    CONSTRAINT runtime_sessions_expires_after_issued CHECK (expires_at > issued_at),
    CONSTRAINT runtime_sessions_last_seen_after_issued CHECK (last_seen_at >= issued_at),
    CONSTRAINT runtime_sessions_revoked_state_matches_time CHECK (
        revoked_at IS NULL OR session_status = 'revoked'
    ),
    CONSTRAINT runtime_sessions_revocation_reason_not_blank CHECK (
        revocation_reason IS NULL OR length(btrim(revocation_reason)) > 0
    ),
    CONSTRAINT runtime_sessions_revocation_reason_matches_time CHECK (
        revocation_reason IS NULL OR revoked_at IS NOT NULL
    ),
    CONSTRAINT runtime_sessions_revoked_after_issued CHECK (
        revoked_at IS NULL OR revoked_at >= issued_at
    ),
    CONSTRAINT runtime_sessions_updated_after_created CHECK (updated_at >= created_at)
);

COMMENT ON TABLE runtime_sessions IS
    'Durable runtime session lifecycle records; stores no raw secret material, verifier digest, or connection state.';

CREATE INDEX runtime_sessions_player_status_idx
    ON runtime_sessions (player_id, session_status);

CREATE INDEX runtime_sessions_status_expires_at_idx
    ON runtime_sessions (session_status, expires_at);

CREATE INDEX runtime_sessions_access_token_record_idx
    ON runtime_sessions (access_token_record_id)
    WHERE access_token_record_id IS NOT NULL;

CREATE INDEX runtime_sessions_last_seen_at_idx
    ON runtime_sessions (last_seen_at);

-- +goose Down
DROP INDEX IF EXISTS runtime_sessions_last_seen_at_idx;
DROP INDEX IF EXISTS runtime_sessions_access_token_record_idx;
DROP INDEX IF EXISTS runtime_sessions_status_expires_at_idx;
DROP INDEX IF EXISTS runtime_sessions_player_status_idx;
DROP TABLE IF EXISTS runtime_sessions;
