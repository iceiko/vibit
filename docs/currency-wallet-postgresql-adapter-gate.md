# Currency Wallet PostgreSQL Adapter Gate

Status: Accepted v0.1
Last updated: 2026-06-07
Scope: Gate-only boundary for the future PostgreSQL adapter implementing `runtime/internal/modules/currency.Repository`
Depends on: `runtime/internal/modules/currency/repository.go`, `docs/currency-wallet-repository-boundary.md`, `runtime/migrations/postgres/000008_create_currency_wallets.sql`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0205`

The paired Simplified Chinese translation is `docs/currency-wallet-postgresql-adapter-gate.zh-CN.md`. The English file is authoritative.

This document defines the currency wallet PostgreSQL adapter gate. It is a gate artifact. It does not add PostgreSQL adapter implementation, SQL execution behavior, unit-of-work factory wiring, runtime wallet behavior, grant/spend execution above the adapter, protocol routes, Protobuf source, generated output, dependencies, migrations, automatic startup migration behavior, authentication/session behavior changes, reward integration, inventory integration, purchase behavior, paid-currency behavior, catalog tables, event/audit tables, reservations, settlement, refunds, transfers, hosted deployments, SDK publication, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The currency wallet PostgreSQL adapter gate record is:

```yaml
currency_wallet_postgresql_adapter_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0297
decision: ADR-0205
check_rule: runtime.currency_wallet_postgresql_adapter_gate
source_repository_interface_decision: ADR-0204
repository_interface: runtime/internal/modules/currency.Repository
repository_interface_source: runtime/internal/modules/currency/repository.go
repository_tests: runtime/internal/modules/currency/repository_test.go
source_migration_source_decision: ADR-0202
source_migration_source: runtime/migrations/postgres/000008_create_currency_wallets.sql
currency_wallets_logical_table: currency_wallets
currency_wallet_balances_logical_table: currency_wallet_balances
currency_wallet_transactions_logical_table: currency_wallet_transactions
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_adapter_source_candidate: runtime/internal/platform/persistence/postgres/currency_wallet_repository.go
future_adapter_tests_candidate: runtime/internal/platform/persistence/postgres/currency_wallet_repository_test.go
future_constructor_candidate: NewCurrencyWalletRepositoryForUnitOfWork
unit_of_work_handoff_required: true
transaction_owner: caller_supplied_unit_of_work
sql_mapping_owner: postgresql_adapter
adapter_gate_only: true
postgresql_adapter_added: false
sql_execution_added: false
unit_of_work_factory_added: false
runtime_behavior_added: false
wallet_transaction_behavior_added: false
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
generated_client_library_added: false
release_artifact_added: false
distributed_runtime_added: false
pitaya_distributed_architecture_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_adapter_implementation_work_item: W-0298
future_adapter_implementation_direction: implement_currency_wallet_postgresql_adapter
```

## 2. Purpose

`W-0296` implemented the storage-neutral `runtime/internal/modules/currency.Repository` interface. The next useful boundary is the platform adapter gate that will later map that interface to the existing PostgreSQL `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` tables.

This gate records the future implementation shape before any adapter SQL is written:

- adapter ownership;
- constructor and executor handoff expectations;
- transaction and unit-of-work boundaries;
- SQL mapping posture for wallet creation, wallet lookup, owner lookup, balance listing, grant recording, spend recording, and transaction listing;
- idempotency, uniqueness, version, affected-row, insufficient-balance, and driver-error mapping;
- redaction posture for wallet, transaction, actor, metadata, SQL, DSN, and driver details;
- focused test expectations;
- stop conditions that keep runtime, protocol, generated output, integrations, hosted surfaces, distributed runtime, and direct compatibility out of the adapter gate.

This is not an implementation. The future adapter source path is named only so agents can verify later work against the accepted boundary.

## 3. Ownership

The future adapter owner is:

```yaml
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
repository_interface_owner: runtime/internal/modules/currency
sql_mapping_owner: runtime/internal/platform/persistence/postgres
transaction_owner: caller_supplied_unit_of_work
application_layer_owns_request_identity: true
currency_module_owns_repository_vocabulary: true
player_module_owns_player_lifecycle: true
inventory_module_owns_items_not_wallet_balances: true
websocket_transport_owns_currency_wallets: false
protocol_adapter_owns_currency_wallets: false
authentication_module_owns_currency_wallets: false
storage_module_owns_currency_wallets: false
friends_module_owns_currency_wallets: false
```

Rules:

- The adapter may later implement `currency.Repository` under the PostgreSQL platform package.
- The adapter must not move SQL into `runtime/internal/modules/currency`.
- The adapter must not own request authentication, player identity validation, route policy, protocol parsing, WebSocket state, reward execution, inventory grants, purchase validation, catalog ownership, settlement, refund, transfer, or distributed topology.
- The adapter must receive already-normalized repository input or call currency module normalizers before SQL mapping.
- The adapter must return currency module value types and typed repository errors, not driver-specific errors.
- The adapter may check player-account foreign-key outcomes as storage conflicts, but it must not become the player account lifecycle owner.

## 4. Future Constructor And Executor Handoff

The first adapter implementation should follow existing PostgreSQL adapter patterns:

```yaml
future_constructor_candidate: NewCurrencyWalletRepositoryForUnitOfWork
future_repository_interface: runtime/internal/modules/currency.Repository
executor_source: caller_supplied
transaction_control_sql_allowed: false
unit_of_work_handoff_required: true
connection_pool_owned_by_adapter: false
context_required: true
```

Rules:

- The constructor should accept an existing executor or query interface supplied by a unit-of-work boundary rather than owning a pool directly.
- The adapter must not issue `BEGIN`, `COMMIT`, or `ROLLBACK`; transaction ownership remains with the unit-of-work runner.
- The adapter must not create automatic startup migrations.
- The adapter must not add a new dependency if existing PostgreSQL platform dependencies already cover the implementation.
- Any required dependency change must be a separate dependency-adoption decision.

## 5. SQL Mapping Posture

The future adapter may map repository methods to the currency wallet tables:

```yaml
wallet_table: currency_wallets
balance_table: currency_wallet_balances
transaction_table: currency_wallet_transactions
wallet_primary_key_column: wallet_id
wallet_owner_columns:
  - owner_kind
  - owner_id
wallet_lifecycle_column: lifecycle_state
wallet_version_column: wallet_version
balance_primary_key_columns:
  - wallet_id
  - currency_code
balance_amount_column: balance_amount
balance_version_column: balance_version
transaction_primary_key_column: transaction_id
transaction_idempotency_columns:
  - wallet_id
  - idempotency_scope
  - idempotency_key
transaction_actor_columns:
  - actor_kind
  - actor_id
metadata_column: metadata_json
currency_wallets_owner_uq: currency_wallets_owner_uq
currency_wallet_balances_pk: currency_wallet_balances_pk
currency_wallet_transactions_idempotency_uq: currency_wallet_transactions_idempotency_uq
```

Method posture:

- `CreateCurrencyWallet` should insert a wallet row for an already validated owner. It must preserve repository validation, start at positive wallet version, map duplicate owner rows to typed conflicts, and return a normalized wallet record.
- `GetCurrencyWallet` should select by `wallet_id` and return a normalized `CurrencyWallet`.
- `GetCurrencyWalletForOwner` should select by `(owner_kind, owner_id)` and must not become arbitrary wallet search or admin inspection.
- `ListCurrencyWalletBalances` should be wallet-scoped, ordered deterministically by `currency_code`, bounded by repository limits, and pagination-ready.
- `RecordCurrencyGrant` should insert one transaction fact and update or create one balance row in the caller's unit of work. It must record positive grant deltas, maintain non-negative `balance_after`, handle idempotency uniqueness, and return the normalized transaction.
- `RecordCurrencySpend` should insert one transaction fact and update one balance row in the caller's unit of work. It must record negative spend deltas, reject insufficient balance without leaking private balance details, handle idempotency uniqueness, and return the normalized transaction only after the balance update succeeds.
- `ListCurrencyWalletTransactions` should be wallet-scoped, optionally currency-filtered, ordered deterministically by creation time and transaction id, bounded by repository limits, and pagination-ready.

Rules:

- SQL text must remain inside the PostgreSQL adapter package.
- Wallet ids, owner ids, transaction ids, actor ids, idempotency keys/scopes, reason codes, external references, metadata JSON, balances, and transaction history are not log-safe by default.
- Driver-specific constraint names may be used internally for error mapping, but public module errors must remain currency-module neutral.
- Affected-row counts must be checked for update operations.
- Balance mutations must not split transaction insertion and balance update across separate application units of work.
- The adapter must not hard-delete wallet or transaction history unless a later retention decision explicitly authorizes it.

## 6. Transaction And Unit-Of-Work Boundary

The future adapter participates in existing transaction handoff:

```yaml
unit_of_work_handoff_required: true
adapter_starts_transactions: false
adapter_commits_transactions: false
adapter_rolls_back_transactions: false
adapter_safe_for_existing_runner: true
balance_mutation_requires_caller_unit_of_work: true
```

Rules:

- Application services or runtime composition may later obtain the adapter through an explicit unit-of-work boundary.
- This gate does not add that factory or composition.
- Adapter methods must use the caller's context.
- The adapter must not hide transaction failures by returning successful wallet or transaction results.
- The adapter must not perform route policy, session validation, access-token validation, WebSocket close behavior, reward execution, inventory mutation, purchase validation, settlement, or refund behavior.

## 7. Error Mapping

The future adapter must collapse PostgreSQL details into currency module errors:

```yaml
repository_error_owner: runtime/internal/modules/currency
driver_error_public_leakage_allowed: false
constraint_name_public_leakage_allowed: false
private_wallet_state_public_leakage_allowed: false
conflict_classes:
  - wallet_not_found
  - owner_not_found
  - duplicate_wallet_owner
  - lifecycle_state_conflict
  - unsupported_owner_kind
  - invalid_currency_code
  - insufficient_balance
  - duplicate_transaction
  - conflicting_duplicate_transaction
  - version_mismatch
  - stale_wallet_version
  - stale_balance_version
  - storage_unavailable
```

Mapping expectations:

- `currency_wallets_owner_uq` conflicts should map to duplicate-wallet-owner style conflicts.
- `currency_wallet_transactions_idempotency_uq` conflicts should map to duplicate transaction or conflicting duplicate transaction according to the later implementation's payload comparison posture.
- Foreign-key failures for wallet owner or transaction wallet references should map to owner-not-found, wallet-not-found, invalid input, or storage conflict without exposing database constraint names.
- No affected row with expected wallet or balance version should map to missing wallet, missing balance, version mismatch, stale version, lifecycle conflict, or insufficient balance without leaking hidden wallet details.
- Malformed input should be rejected before SQL execution when possible.
- Unknown driver or executor failures should map to storage-unavailable style errors with redacted reasons.
- Raw SQL, DSNs, credentials, token material, verifier digests, wallet ids, owner ids, transaction ids, actor ids, idempotency material, balances, metadata, and driver stack details must not appear in public error strings.

## 8. Test Expectations

The later implementation slice should add focused PostgreSQL adapter tests. Fake-executor or query-capture tests are acceptable before live database verification is available.

Required test families for the implementation slice:

```yaml
future_tests:
  - constructor_requires_executor
  - create_currency_wallet_maps_insert_and_duplicate_owner_conflict
  - get_currency_wallet_selects_by_wallet_id
  - get_currency_wallet_for_owner_selects_by_owner_identity
  - list_currency_wallet_balances_is_wallet_scoped_ordered_and_bounded
  - record_currency_grant_inserts_transaction_and_updates_balance
  - record_currency_spend_checks_balance_and_insufficient_balance_conflict
  - idempotency_unique_conflicts_are_mapped
  - wallet_and_balance_version_checks_are_mapped
  - rows_scan_through_currency_normalizers
  - driver_errors_are_redacted
  - transaction_control_sql_is_absent
  - default_tests_do_not_require_live_postgresql
```

Rules:

- Tests must not require protocol routes.
- Tests must not require WebSocket transport.
- Tests must not print raw credentials, tokens, verifier keys, digests, DSNs, query strings, authorization headers, cookies, player ids, wallet ids, transaction ids, idempotency material, balances, or private wallet state.
- Live PostgreSQL verification may be added later when an implementation slice authorizes it; if unavailable, it must be explicitly recorded.

## 9. Relationship To Runtime, Protocol, And Integrations

This gate does not change runtime or protocol behavior:

```yaml
runtime_currency_handlers_added: false
runtime_currency_service_added: false
startup_wiring_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
reward_integration_added: false
inventory_integration_added: false
purchase_behavior_added: false
```

Rules:

- Future runtime wallet behavior needs a later runtime behavior gate.
- Future Protobuf messages and routes need a later protocol gate.
- Future reward, inventory, purchase, payment, catalog, transfer, refund, or settlement behavior needs separate bounded work.
- This gate must not change authentication/session behavior. Metadata-only `player_id` and `session_id` remain insufficient proof.

## 10. Reference Alignment

Nakama and Hiro guide the product capability need for durable currency and wallet state. vibit adapts that need through an explicit storage-neutral repository and PostgreSQL adapter boundary. Pitaya remains a future distributed architecture reference.

Rejected in this gate:

- direct Nakama API compatibility;
- direct Pitaya API compatibility;
- direct Hiro API compatibility;
- runtime economy handlers;
- catalog, reward, inventory, purchase, payment, settlement, refund, reservation, and transfer behavior;
- distributed runtime behavior.

## 11. Next Work Item

This gate opens:

```yaml
future_adapter_implementation_work_item: W-0298
future_adapter_implementation_direction: implement_currency_wallet_postgresql_adapter
```

`W-0298` may implement the currency wallet PostgreSQL adapter only within the source, test, constructor, unit-of-work, SQL mapping, redaction, and deferral posture recorded here.
