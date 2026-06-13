# Currency Wallet Repository Boundary

Status: Accepted v0.1
Last updated: 2026-06-07
Scope: Gate-only boundary for the future storage-neutral currency wallet repository after the PostgreSQL `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` migration source
Depends on: `docs/currency-wallet-lifecycle-boundary-gate.md`, `docs/currency-wallet-persistence-schema-gate.md`, `decisions/ADR-0202-currency-wallet-migration-source.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0203`

The paired Simplified Chinese translation is `docs/currency-wallet-repository-boundary.zh-CN.md`. The English file is authoritative.

This document defines the currency wallet repository boundary. It is a gate artifact. It does not add Go repository interfaces, PostgreSQL adapter behavior, runtime wallet behavior, grant/spend execution, protocol routes, Protobuf source, generated output, dependencies, migrations, automatic startup migration behavior, reward integration, inventory integration, purchase behavior, paid-currency behavior, catalog tables, event/audit tables, reservations, settlement, refunds, transfers, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The currency wallet repository boundary record is:

```yaml
currency_wallet_repository_boundary: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0295
decision: ADR-0203
check_rule: runtime.currency_wallet_repository_boundary
source_migration_source_decision: ADR-0202
source_migration_source: runtime/migrations/postgres/000008_create_currency_wallets.sql
source_schema_gate_decision: ADR-0201
source_lifecycle_gate_decision: ADR-0200
future_repository_owner_candidate: runtime/internal/modules/currency
future_repository_interface_candidate: runtime/internal/modules/currency.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
currency_wallets_logical_table: currency_wallets
currency_wallet_balances_logical_table: currency_wallet_balances
currency_wallet_transactions_logical_table: currency_wallet_transactions
repository_boundary_gate_only: true
repository_interface_added: false
postgresql_adapter_added: false
runtime_behavior_added: false
grant_spend_behavior_added: false
authentication_session_behavior_changed: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
dependency_added: false
migration_added: false
currency_catalog_table_added: false
currency_wallet_events_table_added: false
reward_integration_added: false
inventory_integration_added: false
purchase_behavior_added: false
paid_currency_behavior_added: false
reservation_behavior_added: false
settlement_behavior_added: false
refund_behavior_added: false
transfer_behavior_added: false
hosted_deployment_added: false
sdk_added: false
release_artifact_added: false
distributed_runtime_added: false
pitaya_distributed_architecture_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_repository_interface_work_item: W-0296
future_repository_interface_direction: implement_currency_wallet_repository_interface
```

## 2. Purpose

`W-0294` added the PostgreSQL migration source for `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`. The next useful boundary is the storage-neutral repository vocabulary that future implementation code can use without exposing SQL details, transport details, protocol assumptions, or direct external economy compatibility.

This boundary prepares Nakama/Hiro-class economy capability coverage by recording:

- repository ownership;
- candidate wallet, balance, transaction, actor, idempotency, and pagination value types;
- repository command and query vocabulary;
- validated identity and actor handoff rules;
- transaction, idempotency, and optimistic version posture;
- conflict, redaction, and error posture;
- PostgreSQL adapter expectations tied to the W-0294 migration;
- stop conditions for future implementation work.

This is still not a runtime feature. No route, handler, adapter, repository interface, or protocol message can use currency wallets until later bounded work items explicitly authorize them.

## 3. Ownership

The future repository is currency module-owned:

```yaml
future_repository_owner_candidate: runtime/internal/modules/currency
future_repository_interface_candidate: runtime/internal/modules/currency.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
currency_wallet_table_owner: runtime.currency
application_layer_owns_request_identity: true
postgresql_adapter_owns_sql_mapping: true
websocket_transport_owns_currency_wallets: false
protocol_adapter_owns_currency_wallets: false
authentication_module_owns_currency_wallets: false
player_module_owns_currency_wallets: false
inventory_module_owns_currency_wallets: false
storage_module_owns_currency_wallets: false
friends_module_owns_currency_wallets: false
```

Rules:

- The future repository interface must be storage-neutral and module-facing.
- The interface must not mention PostgreSQL, pgx, SQL rows, goose migrations, prepared statements, connection pools, transaction runners, or database driver errors.
- The PostgreSQL adapter may later implement the interface under `runtime/internal/platform/persistence/postgres`, but only after a separate adapter gate.
- Application or handler code may later call currency wallet behavior through module/application boundaries, not through SQL or transport state.
- Authentication and session code provide validated request identity; they do not own wallet records.
- Player account storage owns player lifecycle state, not currency wallet balances.
- Inventory owns inventory item state, not wallet balances or currency transactions.
- WebSocket transport owns connection plumbing, not economy state.
- Protocol adapters own wire conversion, not repository behavior.

## 4. Candidate Value Types

A later implementation gate may rename or reduce these shapes, but the first repository interface implementation should start from this vocabulary:

```yaml
candidate_value_types:
  - CurrencyWallet
  - CurrencyWalletID
  - CurrencyWalletOwner
  - CurrencyWalletLifecycleState
  - CurrencyWalletVersion
  - CurrencyWalletBalance
  - CurrencyCode
  - CurrencyAmount
  - CurrencyBalanceVersion
  - CurrencyWalletTransaction
  - CurrencyWalletTransactionID
  - CurrencyWalletTransactionKind
  - CurrencyWalletActor
  - CurrencyWalletIdempotencyKey
  - CurrencyWalletIdempotencyScope
  - CreateCurrencyWalletInput
  - GetCurrencyWalletInput
  - GetCurrencyWalletForOwnerInput
  - ListCurrencyWalletBalancesInput
  - RecordCurrencyGrantInput
  - RecordCurrencySpendInput
  - ListCurrencyWalletTransactionsInput
  - CurrencyWalletConflict
  - CurrencyWalletRepositoryError
```

First-posture record vocabulary:

```yaml
currency_wallet_record:
  wallet_id: server_generated_opaque_id
  owner_kind: player
  owner_id: validated_player_id_handoff
  lifecycle_state: active_or_suspended_or_closed
  wallet_version: server_managed_bigint_revision
  created_at: server_timestamp
  updated_at: server_timestamp
  state_changed_at: server_timestamp
  suspended_at: nullable_server_timestamp
  closed_at: nullable_server_timestamp
currency_wallet_balance_record:
  wallet_id: wallet_identity
  currency_code: bounded_case_sensitive_currency_code
  balance_amount: non_negative_bigint_minor_unit
  balance_version: server_managed_bigint_revision
  created_at: server_timestamp
  updated_at: server_timestamp
currency_wallet_transaction_record:
  transaction_id: server_generated_opaque_id
  wallet_id: wallet_identity
  currency_code: bounded_case_sensitive_currency_code
  transaction_kind: grant_or_spend
  amount_delta: signed_bigint_minor_unit
  balance_after: non_negative_bigint_minor_unit
  idempotency_key: bounded_duplicate_prevention_key
  idempotency_scope: bounded_duplicate_prevention_scope
  actor_kind: system_or_player_or_operation
  actor_id: nullable_actor_identifier
  reason_code: nullable_bounded_text
  external_reference: nullable_bounded_text
  metadata_json: nullable_small_json_object
  created_at: server_timestamp
```

Rules:

- `owner_kind` must remain a closed vocabulary. The first allowed owner kind is `player`.
- `owner_id` is input identity handoff, not proof. Future runtime behavior must derive it from validated request identity before calling the repository.
- `currency_code` remains case-sensitive until a later catalog gate selects normalization.
- `balance_amount` and `balance_after` are non-negative integer minor units.
- `wallet_version` and `balance_version` are server-managed and must not be client-authoritative state.
- `transaction_id`, `wallet_id`, detailed balances, idempotency fields, reason codes, external references, and metadata are not log-safe by default.

## 5. Candidate Repository Capabilities

The first storage-neutral capability family is:

```yaml
candidate_repository_capabilities:
  - CreateCurrencyWallet
  - GetCurrencyWallet
  - GetCurrencyWalletForOwner
  - ListCurrencyWalletBalances
  - RecordCurrencyGrant
  - RecordCurrencySpend
  - ListCurrencyWalletTransactions
```

Capability rules:

- `CreateCurrencyWallet` may create a wallet row for an already validated owner identity.
- `GetCurrencyWallet` is a storage lookup by server wallet id. It must not authenticate users, validate access tokens, or create request identity.
- `GetCurrencyWalletForOwner` must be owner-scoped and must not become arbitrary wallet search or admin inspection without a later gate.
- `ListCurrencyWalletBalances` must be wallet-scoped and pagination-ready where future catalog size requires it.
- `RecordCurrencyGrant` must insert one transaction fact and update or create one balance row in the same application unit of work.
- `RecordCurrencySpend` must insert one transaction fact and update one balance row in the same application unit of work, preserving insufficient-balance conflict posture.
- `ListCurrencyWalletTransactions` must be wallet-scoped and pagination-ready. It must not expose private support metadata without later protocol and redaction gates.
- All methods must return typed module-owned records and errors, not raw SQL rows or database driver errors.

The future repository interface may choose shorter names, but it must preserve the semantic split between wallet creation, wallet lookup, balance listing, grant recording, spend recording, transaction listing, idempotency, and conflict handling.

## 6. Identity, Transaction, And Idempotency Handoff

The repository boundary prepares identity, transaction, and idempotency handling without implementing behavior:

```yaml
wallet_owner_identity: owner_kind_plus_owner_id
owner_kind_first_value: player
actor_identity_source: validated_request_identity_or_server_operation_before_repository_call
client_supplied_actor_id_as_proof_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
idempotency_owner: currency_wallet_transactions
idempotency_unique_identity:
  - wallet_id
  - idempotency_scope
  - idempotency_key
transaction_required_for_balance_mutation: true
balance_mutation_without_transaction_allowed: false
```

Rules:

- The repository may receive normalized owner and actor ids as data, but ids are not authentication proof.
- Future runtime behavior must validate actor permissions before calling grant or spend repository capabilities.
- Duplicate idempotency keys must not double-apply a balance mutation within the same wallet and idempotency scope.
- Conflicting duplicate payload behavior is deferred to future runtime and protocol gates, but must not silently create a second transaction.
- Currency wallet mutations must run inside the application unit-of-work boundary.
- Repository methods must not accept transport metadata, WebSocket connection identifiers, cookies, headers, tokens, or sessions as identity proof.

## 7. Version And Conflict Handoff

The repository boundary prepares optimistic concurrency without implementing behavior:

```yaml
wallet_version_storage: BIGINT
balance_version_storage: BIGINT
initial_create_version: 1
version_owner: server
client_authoritative_version_allowed: false
expected_version_handoff: future_behavior_or_interface_gate
conflict_public_shape: deferred_to_protocol_gate
```

Candidate conflict classes:

```yaml
candidate_conflict_classes:
  - wallet_not_found
  - wallet_already_exists
  - wallet_owner_mismatch
  - wallet_not_active
  - currency_balance_not_found
  - invalid_currency_code
  - invalid_currency_amount
  - insufficient_balance
  - duplicate_transaction
  - conflicting_duplicate_transaction
  - invalid_idempotency_key
  - invalid_transaction_kind
  - version_mismatch
  - stale_wallet_version
  - stale_balance_version
  - storage_unavailable
```

Rules:

- Repository methods may later distinguish internal typed conflicts, but public protocol error mapping remains deferred.
- Version equality is not authentication proof.
- A stale expected version must not be collapsed into a hidden successful write.
- Owner mismatch must not reveal another player's wallet existence when future public behavior gates require leakage collapse.
- The PostgreSQL adapter must map unique-index, affected-row, check-constraint, and foreign-key outcomes into typed repository conflicts without exposing driver error text.

## 8. Redaction And Logging

Currency wallet state is private economy data.

```yaml
private_wallet_state_log_safe: false
wallet_id_log_safe: false_by_default
player_id_log_safe: false_by_default
detailed_balance_log_safe: false
transaction_id_log_safe: false_by_default
idempotency_key_log_safe: false
idempotency_scope_log_safe: false
reason_code_log_safe: conditional_after_validation
external_reference_log_safe: false
metadata_json_log_safe: false
forbidden_repository_material:
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
  - authorization_header
  - cookie
  - payment_provider_secret
  - payment_provider_payload
  - nakama_api_path
  - pitaya_server_id
```

Rules:

- Repository errors must be redacted and typed.
- Raw values and storage driver errors must not be logged by default.
- Transaction metadata may contain private support context and must not be logged by default.
- The repository must not store or return authentication material, token material, verifier digests, transport metadata, payment secrets, full provider payloads, or direct external compatibility markers.

## 9. PostgreSQL Adapter Expectations

The future PostgreSQL adapter may later map the repository to:

```yaml
logical_tables:
  - currency_wallets
  - currency_wallet_balances
  - currency_wallet_transactions
wallet_owner_unique_index: currency_wallets_owner_uq
balance_primary_key: currency_wallet_balances_pk
transaction_idempotency_unique_index: currency_wallet_transactions_idempotency_uq
transaction_wallet_index: currency_wallet_transactions_wallet_created_at_idx
transaction_wallet_currency_index: currency_wallet_transactions_wallet_currency_created_at_idx
wallet_version_column: wallet_version
balance_version_column: balance_version
```

Adapter rules:

- SQL must stay in `runtime/internal/platform/persistence/postgres`.
- Domain/application code must not import pgx.
- The adapter must run inside the application unit-of-work boundary.
- Balance mutation and transaction insert must be atomic.
- Idempotency uniqueness must be mapped to typed duplicate transaction outcomes.
- Check constraints and foreign-key failures must be mapped to redacted typed conflicts.
- Driver error messages, SQL statements, DSNs, and raw metadata must not be exposed in public errors.

## 10. Future Test Expectations

The future repository interface implementation must add focused tests for:

- wallet owner normalization;
- currency code validation posture;
- non-negative amount and signed delta validation;
- idempotency key and scope validation;
- grant/spend input normalization;
- transaction list pagination inputs;
- typed conflict vocabulary;
- redacted repository errors;
- storage neutrality, including no PostgreSQL imports.

Future PostgreSQL adapter tests must cover:

- wallet creation and owner uniqueness;
- wallet lookup by id and owner;
- balance listing;
- grant transaction plus balance update atomicity;
- spend transaction plus balance update atomicity;
- insufficient balance conflict mapping;
- duplicate idempotency handling;
- lifecycle state filtering;
- version conflict handling;
- redaction of database errors.

Future runtime behavior tests remain deferred and must cover validated identity, permissions, suspended and closed wallet behavior, duplicate request behavior, public error mapping, and protocol redaction.

## 11. Stop Conditions

Stop and create a new gate before any of the following:

- adding Go repository interface implementation;
- adding PostgreSQL adapter code;
- executing SQL from Go;
- adding runtime wallet behavior;
- adding grant or spend command execution;
- adding reward, inventory, purchase, paid-currency, reservation, settlement, refund, or transfer behavior;
- adding currency catalog or event/audit tables;
- adding protocol routes;
- adding Protobuf source or generated output;
- changing authentication/session behavior;
- adding dependencies;
- wiring startup;
- publishing SDKs or hosted surfaces;
- adding distributed runtime behavior;
- adding direct Nakama/Pitaya API compatibility.

## 12. Next Step

The next bounded work item is:

```text
W-0296 Implement storage-neutral currency wallet repository interface
```

That follow-up may add the storage-neutral Go repository interface if it follows this boundary and remains interface-only. PostgreSQL adapters, runtime behavior, protocol routes, generated output, dependencies, startup wiring, rewards, purchases, inventory integration, hosted surfaces, SDKs, distributed runtime, and direct compatibility remain later work.
