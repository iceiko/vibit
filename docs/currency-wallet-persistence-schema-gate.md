# Currency Wallet Persistence Schema Gate

Status: Accepted v0.1
Last updated: 2026-06-07
Scope: Gate for future currency wallet persistence schema before migration source, balance tables, wallet transaction behavior, runtime behavior, protocol routes, generated output, repositories, adapters, or broader economy features
Depends on: `docs/currency-wallet-lifecycle-boundary-gate.md`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0201`

The paired Simplified Chinese translation is `docs/currency-wallet-persistence-schema-gate.zh-CN.md`. The English file is authoritative.

This document defines the currency wallet persistence schema gate. It is a gate artifact. It does not add SQL migration source, create currency catalog, wallet, balance, transaction, ledger, idempotency, reward, purchase, inventory, or audit tables, implement currency wallet runtime behavior, add protocol routes, add Protobuf source or generated output, add dependencies, add repository interfaces, add PostgreSQL adapters, wire startup, change authentication/session behavior, publish SDKs or generated client libraries, create hosted deployments or release artifacts, add distributed runtime behavior, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The currency wallet persistence schema gate record is:

```yaml
currency_wallet_persistence_schema_gate: defined
completed_work_item: W-0293
decision: ADR-0201
check_rule: runtime.currency_wallet_persistence_schema_gate
source_lifecycle_gate_decision: ADR-0200
source_lifecycle_gate_standard: docs/currency-wallet-lifecycle-boundary-gate.md
gate_standard: docs/currency-wallet-persistence-schema-gate.md
gate_standard_translation: docs/currency-wallet-persistence-schema-gate.zh-CN.md
selected_capability_family: economy_inventory_rewards_currencies_and_progression
primary_product_reference: Nakama
secondary_product_reference: Hiro
pitaya_reference_status: deferred_future_architecture_reference
selected_first_currency_wallet_store: postgres
future_currency_wallets_logical_table: currency_wallets
future_currency_wallet_balances_logical_table: currency_wallet_balances
future_currency_wallet_transactions_logical_table: currency_wallet_transactions
future_currency_catalog_logical_table: deferred
future_currency_wallet_events_logical_table: deferred
future_migration_source_candidate: runtime/migrations/postgres/000008_create_currency_wallets.sql
future_repository_owner_candidate: runtime/internal/modules/currency
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
wallet_identity_posture_recorded: true
balance_identity_posture_recorded: true
transaction_ledger_posture_recorded: true
idempotency_posture_recorded: true
catalog_posture_recorded: true
index_uniqueness_posture_recorded: true
timestamp_posture_recorded: true
redaction_posture_recorded: true
future_repository_adapter_boundaries_recorded: true
future_migration_source_candidate_recorded: true
schema_gate_only: true
migration_source_added: false
currency_wallets_table_added: false
currency_wallet_balances_table_added: false
currency_wallet_transactions_table_added: false
currency_catalog_table_added: false
currency_wallet_events_table_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_migration_work_item: W-0294
future_migration_direction: add_currency_wallet_migration_source
```

## 2. Product Intent

`ADR-0200` defined the future currency wallet lifecycle vocabulary: catalog reads, wallet reads, balance reads, grants, spends, balance change recording, and transaction reads. The next step is data-first: balance mutation correctness depends on a stable durable persistence posture before SQL, repositories, adapters, protocol, or runtime handlers exist.

This gate makes the future migration inspectable:

- the table candidates are known;
- player wallet identity is explicit;
- balance row identity is separate from wallet identity;
- transaction or ledger records are required for future mutation facts;
- idempotency is anchored in durable transaction records;
- currency catalog ownership remains explicit and deferred;
- indexes, uniqueness, timestamps, and redaction are planned before implementation;
- future repository and PostgreSQL adapter ownership candidates are named.

The gate keeps the work conservative. It prepares the next migration-source-only slice but does not add the migration file.

## 3. Selected Store And Tables

The first currency wallet persistence target is PostgreSQL:

```yaml
selected_first_currency_wallet_store: postgres
future_migration_source_candidate: runtime/migrations/postgres/000008_create_currency_wallets.sql
future_repository_boundary: separate_future_work_item
future_postgresql_adapter: separate_future_work_item
```

The first schema candidate may define three logical runtime tables:

```yaml
future_tables:
  - currency_wallets
  - currency_wallet_balances
  - currency_wallet_transactions
```

The currency catalog table remains deferred:

```yaml
future_currency_catalog_logical_table: deferred
currency_code_source: future_currency_catalog_or_closed_catalog_configuration
catalog_table_added_by_this_gate: false
```

Rationale:

- PostgreSQL is vibit's first accepted authoritative durable store.
- Wallets and balances need transactional consistency before grant/spend behavior exists.
- A transaction table makes idempotency, auditability, and support workflows inspectable before runtime behavior is written.
- A separate catalog table can be useful, but the first schema gate should not authorize catalog management, rewards, pricing, purchases, or live-operations behavior.

## 4. Future `currency_wallets` Table Candidate

The future first migration may define one logical wallet identity table:

```yaml
currency_wallets:
  primary_key_candidate:
    - wallet_id
  required_columns:
    - wallet_id
    - owner_kind
    - owner_id
    - lifecycle_state
    - wallet_version
    - created_at
    - updated_at
    - state_changed_at
  nullable_columns:
    - suspended_at
    - closed_at
  forbidden_columns:
    - raw_access_token
    - raw_credential
    - credential_lookup_digest
    - credential_verifier_digest
    - token_lookup_digest
    - token_verifier_digest
    - verifier_key
    - websocket_connection_id
    - websocket_subprotocol
    - remote_address
    - nakama_api_path
    - pitaya_server_id
```

`wallet_id` is a server-generated opaque record id. It is not authentication proof and is not log-safe by default.

The first logical wallet identity is:

```text
owner_kind + owner_id
```

The first owner posture is player-owned wallets only:

```yaml
owner_kind_first_value: player
owner_id_source: validated_request_identity_player_id
owner_player_fk_candidate: player_accounts(player_id)
owner_kind_check_candidate: owner_kind = 'player'
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
```

Future owner kinds such as global, group, guild, party, match, server shard, or operations account require later gates.

## 5. Wallet Lifecycle State Representation

The first lifecycle state candidate follows the lifecycle gate:

```yaml
lifecycle_state_column: lifecycle_state
lifecycle_state_type: TEXT
allowed_lifecycle_states:
  - active
  - suspended
  - closed
```

Rules:

- `active` wallets may be read and may later accept authorized grants and spends.
- `suspended` and `closed` mutation behavior remains deferred to a runtime behavior gate.
- `wallet_version` should be positive and server-managed.
- `updated_at >= created_at` should be enforced.
- `state_changed_at >= created_at` should be enforced.
- `closed_at` and `suspended_at` are persistence facts, not authorization proof.

## 6. Future `currency_wallet_balances` Table Candidate

The future first migration may define one logical balance table:

```yaml
currency_wallet_balances:
  primary_key_candidate:
    - wallet_id
    - currency_code
  required_columns:
    - wallet_id
    - currency_code
    - balance_amount
    - balance_version
    - created_at
    - updated_at
  nullable_columns: []
```

Balance posture:

```yaml
currency_code_type_candidate: TEXT
currency_code_max_length_candidate: 64
amount_type_candidate: BIGINT
amount_unit: integer_minor_unit
negative_balance_allowed_by_default: false
balance_version_type_candidate: BIGINT
```

Rules:

- `wallet_id` references `currency_wallets(wallet_id)`.
- `currency_code` is case-sensitive until a later catalog gate chooses normalization rules.
- `currency_code` must be non-blank and bounded.
- `balance_amount >= 0` is the first default posture.
- `balance_version` must be positive and server-managed.
- The row identity is `(wallet_id, currency_code)`.
- Detailed balances are not log-safe by default.

## 7. Future `currency_wallet_transactions` Table Candidate

The future first migration may define one logical transaction or ledger fact table:

```yaml
currency_wallet_transactions:
  primary_key_candidate:
    - transaction_id
  required_columns:
    - transaction_id
    - wallet_id
    - currency_code
    - transaction_kind
    - amount_delta
    - balance_after
    - idempotency_key
    - idempotency_scope
    - actor_kind
    - created_at
  nullable_columns:
    - actor_id
    - reason_code
    - external_reference
    - metadata_json
```

The first transaction kind vocabulary is:

```yaml
transaction_kinds:
  - grant
  - spend
```

Rules:

- `transaction_id` is server-generated and not log-safe by default.
- `amount_delta` records signed balance movement.
- Grant deltas are positive; spend deltas are negative.
- `balance_after` records the resulting balance for supportability and auditability.
- `idempotency_key` and `idempotency_scope` are durable duplicate-prevention fields and are not log-safe by default.
- `metadata_json` must not store raw credentials, raw tokens, verifier keys, digests, transport metadata, DSNs, payment secrets, or full external provider payloads.
- The first schema should enforce that duplicate idempotency keys cannot double-apply a mutation within the chosen scope.

## 8. Idempotency And Uniqueness

The first idempotency posture is:

```yaml
idempotency_owner: currency_wallet_transactions
idempotency_scope_candidate:
  - wallet_id
  - idempotency_scope
  - idempotency_key
duplicate_application_allowed: false
conflicting_duplicate_payload_public_error: CURRENCY_DUPLICATE_TRANSACTION
```

The first uniqueness and lookup posture is:

```yaml
unique_wallet_owner_candidate:
  - owner_kind
  - owner_id
unique_balance_candidate:
  - wallet_id
  - currency_code
unique_idempotency_candidate:
  - wallet_id
  - idempotency_scope
  - idempotency_key
lookup_indexes:
  - owner_kind_owner_id
  - wallet_id_currency_code
  - wallet_id_created_at
  - wallet_id_currency_code_created_at
```

The exact SQL index names are deferred to the migration-source slice.

## 9. Catalog And Economy Integration Deferrals

Currency catalog persistence is intentionally deferred:

```yaml
currency_catalog_table_added_by_this_gate: false
reward_integration_added_by_this_gate: false
inventory_integration_added_by_this_gate: false
purchase_behavior_added_by_this_gate: false
paid_currency_behavior_added_by_this_gate: false
```

Future currency catalog work must decide:

- code normalization;
- display metadata;
- precision and minor-unit semantics;
- allowed grant/spend posture per currency;
- soft-delete or deprecation behavior;
- whether catalog rows are runtime mutable, config-owned, or migration-owned.

Rewards, purchases, inventory pricing, live-operations grants, refunds, reservations, settlement, transfers, and paid-currency behavior remain separate future gates.

## 10. Future Repository And Adapter Boundaries

Future repository and adapter ownership candidates are:

```yaml
future_repository_owner_candidate: runtime/internal/modules/currency
future_repository_interface_candidate: runtime/internal/modules/currency.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_postgresql_adapter_source_candidate: runtime/internal/platform/persistence/postgres/currency_wallet_repository.go
future_postgresql_adapter_tests_candidate: runtime/internal/platform/persistence/postgres/currency_wallet_repository_test.go
```

Rules:

- The repository interface must not exist until a later repository boundary or implementation work item authorizes it.
- PostgreSQL SQL execution must stay in the PostgreSQL platform owner package.
- Domain/application code must not import `pgx`.
- Future mutations must run inside the application unit-of-work boundary.
- Future repository methods must not accept client-supplied actor ids as proof.
- Future adapter errors must be mapped to redacted module/application errors.

## 11. Redaction

Not log-safe by default:

- wallet ids;
- player ids;
- transaction ids;
- detailed balances;
- idempotency keys;
- reason codes when they reveal private operations;
- external references;
- transaction metadata;
- currency catalog internals before public contract approval.

Forbidden persistence fields:

- raw credentials;
- raw access tokens;
- verifier keys;
- lookup digests;
- verifier digests;
- WebSocket connection ids;
- WebSocket subprotocols;
- remote addresses;
- authorization headers;
- cookies;
- DSNs with credentials;
- full payment provider payloads;
- direct Nakama or Pitaya API path compatibility markers.

## 12. Future Test Expectations

Future migration-source checks must cover:

- table names;
- required columns;
- owner identity checks;
- wallet lifecycle state checks;
- balance amount non-negativity;
- positive version checks;
- timestamp checks;
- transaction kind checks;
- idempotency uniqueness;
- wallet owner uniqueness;
- balance uniqueness;
- forbidden secret, digest, transport, payment, Nakama, or Pitaya compatibility columns.

Future repository and adapter tests must be defined after migration source and repository/adapter gates.

Future runtime behavior tests remain deferred and must cover validated identity, grant/spend transactionality, duplicate idempotency handling, insufficient balance, invalid currency code, suspended/closed wallet behavior, public error collapse, and redacted logs before runtime behavior is implemented.

## 13. Non-Authorization

This gate does not authorize:

- SQL migration source creation;
- currency catalog tables;
- currency wallet tables;
- currency balance tables;
- transaction or ledger tables;
- idempotency tables separate from planned transaction posture;
- event/audit tables;
- repository interfaces;
- PostgreSQL adapters;
- runtime wallet behavior;
- grant/spend execution;
- reward integration;
- inventory integration;
- purchase behavior;
- paid currency behavior;
- reservations, settlement, refunds, or transfers;
- protocol routes;
- Protobuf sources;
- generated output;
- dependencies;
- startup wiring;
- authentication/session behavior changes;
- hosted deployments;
- SDK publication;
- release artifacts;
- distributed runtime behavior;
- direct Nakama/Pitaya API compatibility.

## 14. Next Step

The next bounded work item is:

```text
W-0294 Add currency wallet migration source
```

That follow-up may add the SQL migration source if it follows this gate and remains migration-source-only. Repository interfaces, adapters, runtime behavior, protocol routes, generated output, dependencies, startup wiring, rewards, purchases, inventory integration, hosted surfaces, SDKs, distributed runtime, and direct compatibility remain later work.
