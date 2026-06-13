-- +goose Up
-- Module: currency
-- Purpose: Create durable player-owned currency wallet, balance, and transaction records.

CREATE TABLE currency_wallets (
    wallet_id TEXT PRIMARY KEY,
    owner_kind TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    lifecycle_state TEXT NOT NULL,
    wallet_version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    state_changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    suspended_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    CONSTRAINT currency_wallets_owner_player_fk
        FOREIGN KEY (owner_id)
        REFERENCES player_accounts(player_id)
        ON DELETE RESTRICT,
    CONSTRAINT currency_wallets_wallet_id_not_blank CHECK (length(btrim(wallet_id)) > 0),
    CONSTRAINT currency_wallets_owner_kind_valid CHECK (owner_kind = 'player'),
    CONSTRAINT currency_wallets_owner_id_not_blank CHECK (length(btrim(owner_id)) > 0),
    CONSTRAINT currency_wallets_lifecycle_state_valid CHECK (
        lifecycle_state IN ('active', 'suspended', 'closed')
    ),
    CONSTRAINT currency_wallets_wallet_version_positive CHECK (wallet_version > 0),
    CONSTRAINT currency_wallets_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT currency_wallets_state_changed_after_created CHECK (state_changed_at >= created_at),
    CONSTRAINT currency_wallets_suspended_after_created CHECK (
        suspended_at IS NULL OR suspended_at >= created_at
    ),
    CONSTRAINT currency_wallets_closed_after_created CHECK (
        closed_at IS NULL OR closed_at >= created_at
    )
);

COMMENT ON TABLE currency_wallets IS
    'Player-owned currency wallet lifecycle rows; excludes secret material, digest material, transport metadata, payment provider payloads, distributed runtime markers, or external compatibility columns.';

CREATE UNIQUE INDEX currency_wallets_owner_uq
    ON currency_wallets (owner_kind, owner_id);

CREATE INDEX currency_wallets_updated_at_idx
    ON currency_wallets (updated_at);

CREATE TABLE currency_wallet_balances (
    wallet_id TEXT NOT NULL,
    currency_code TEXT NOT NULL,
    balance_amount BIGINT NOT NULL DEFAULT 0,
    balance_version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT currency_wallet_balances_pk
        PRIMARY KEY (wallet_id, currency_code),
    CONSTRAINT currency_wallet_balances_wallet_fk
        FOREIGN KEY (wallet_id)
        REFERENCES currency_wallets(wallet_id)
        ON DELETE RESTRICT,
    CONSTRAINT currency_wallet_balances_wallet_id_not_blank CHECK (length(btrim(wallet_id)) > 0),
    CONSTRAINT currency_wallet_balances_currency_code_not_blank CHECK (length(btrim(currency_code)) > 0),
    CONSTRAINT currency_wallet_balances_currency_code_length CHECK (length(currency_code) <= 64),
    CONSTRAINT currency_wallet_balances_amount_non_negative CHECK (balance_amount >= 0),
    CONSTRAINT currency_wallet_balances_version_positive CHECK (balance_version > 0),
    CONSTRAINT currency_wallet_balances_updated_after_created CHECK (updated_at >= created_at)
);

COMMENT ON TABLE currency_wallet_balances IS
    'Current currency wallet balances keyed by wallet and currency code; excludes catalog ownership, runtime mutation behavior, protocol shape, or external compatibility columns.';

CREATE INDEX currency_wallet_balances_wallet_currency_idx
    ON currency_wallet_balances (wallet_id, currency_code);

CREATE INDEX currency_wallet_balances_updated_at_idx
    ON currency_wallet_balances (updated_at);

CREATE TABLE currency_wallet_transactions (
    transaction_id TEXT PRIMARY KEY,
    wallet_id TEXT NOT NULL,
    currency_code TEXT NOT NULL,
    transaction_kind TEXT NOT NULL,
    amount_delta BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    idempotency_key TEXT NOT NULL,
    idempotency_scope TEXT NOT NULL,
    actor_kind TEXT NOT NULL,
    actor_id TEXT,
    reason_code TEXT,
    external_reference TEXT,
    metadata_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT currency_wallet_transactions_wallet_fk
        FOREIGN KEY (wallet_id)
        REFERENCES currency_wallets(wallet_id)
        ON DELETE RESTRICT,
    CONSTRAINT currency_wallet_transactions_transaction_id_not_blank CHECK (length(btrim(transaction_id)) > 0),
    CONSTRAINT currency_wallet_transactions_wallet_id_not_blank CHECK (length(btrim(wallet_id)) > 0),
    CONSTRAINT currency_wallet_transactions_currency_code_not_blank CHECK (length(btrim(currency_code)) > 0),
    CONSTRAINT currency_wallet_transactions_currency_code_length CHECK (length(currency_code) <= 64),
    CONSTRAINT currency_wallet_transactions_kind_valid CHECK (
        transaction_kind IN ('grant', 'spend')
    ),
    CONSTRAINT currency_wallet_transactions_delta_direction CHECK (
        (transaction_kind = 'grant' AND amount_delta > 0) OR
        (transaction_kind = 'spend' AND amount_delta < 0)
    ),
    CONSTRAINT currency_wallet_transactions_balance_after_non_negative CHECK (balance_after >= 0),
    CONSTRAINT currency_wallet_transactions_idempotency_key_not_blank CHECK (length(btrim(idempotency_key)) > 0),
    CONSTRAINT currency_wallet_transactions_idempotency_key_length CHECK (length(idempotency_key) <= 256),
    CONSTRAINT currency_wallet_transactions_idempotency_scope_not_blank CHECK (length(btrim(idempotency_scope)) > 0),
    CONSTRAINT currency_wallet_transactions_idempotency_scope_length CHECK (length(idempotency_scope) <= 128),
    CONSTRAINT currency_wallet_transactions_actor_kind_valid CHECK (
        actor_kind IN ('system', 'player', 'operation')
    ),
    CONSTRAINT currency_wallet_transactions_actor_id_not_blank CHECK (
        actor_id IS NULL OR length(btrim(actor_id)) > 0
    ),
    CONSTRAINT currency_wallet_transactions_reason_code_length CHECK (
        reason_code IS NULL OR length(reason_code) <= 128
    ),
    CONSTRAINT currency_wallet_transactions_external_reference_length CHECK (
        external_reference IS NULL OR length(external_reference) <= 256
    ),
    CONSTRAINT currency_wallet_transactions_metadata_json_object CHECK (
        metadata_json IS NULL OR jsonb_typeof(metadata_json) = 'object'
    )
);

COMMENT ON TABLE currency_wallet_transactions IS
    'Currency wallet transaction fact rows for future grants and spends; excludes secret material, digest material, transport metadata, full payment payloads, reward execution, inventory integration, protocol shape, or external compatibility columns.';

CREATE UNIQUE INDEX currency_wallet_transactions_idempotency_uq
    ON currency_wallet_transactions (wallet_id, idempotency_scope, idempotency_key);

CREATE INDEX currency_wallet_transactions_wallet_created_at_idx
    ON currency_wallet_transactions (wallet_id, created_at);

CREATE INDEX currency_wallet_transactions_wallet_currency_created_at_idx
    ON currency_wallet_transactions (wallet_id, currency_code, created_at);

-- +goose Down
DROP INDEX IF EXISTS currency_wallet_transactions_wallet_currency_created_at_idx;
DROP INDEX IF EXISTS currency_wallet_transactions_wallet_created_at_idx;
DROP INDEX IF EXISTS currency_wallet_transactions_idempotency_uq;
DROP TABLE IF EXISTS currency_wallet_transactions;
DROP INDEX IF EXISTS currency_wallet_balances_updated_at_idx;
DROP INDEX IF EXISTS currency_wallet_balances_wallet_currency_idx;
DROP TABLE IF EXISTS currency_wallet_balances;
DROP INDEX IF EXISTS currency_wallets_updated_at_idx;
DROP INDEX IF EXISTS currency_wallets_owner_uq;
DROP TABLE IF EXISTS currency_wallets;
