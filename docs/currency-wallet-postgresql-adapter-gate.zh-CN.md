# Currency Wallet PostgreSQL Adapter Gate

状态：Accepted v0.1
最后更新：2026-06-07
范围：未来实现 `runtime/internal/modules/currency.Repository` 的 PostgreSQL adapter 的 gate-only boundary
依赖：`runtime/internal/modules/currency/repository.go`、`docs/currency-wallet-repository-boundary.md`、`runtime/migrations/postgres/000008_create_currency_wallets.sql`、`docs/postgresql-persistence-boundary.md`、`docs/reference-game-server-alignment.md`
Canonical decision：`ADR-0205`

配对英文原文是 `docs/currency-wallet-postgresql-adapter-gate.md`。英文文件是权威版本。

本文定义 currency wallet PostgreSQL adapter gate。它只是 gate artifact，不添加 PostgreSQL adapter implementation、SQL execution behavior、unit-of-work factory wiring、runtime wallet behavior、adapter 之上的 grant/spend execution、protocol routes、Protobuf source、generated output、dependencies、migrations、automatic startup migration behavior、authentication/session behavior changes、reward integration、inventory integration、purchase behavior、paid-currency behavior、catalog tables、event/audit tables、reservations、settlement、refunds、transfers、hosted deployments、SDK publication、release artifacts、public announcements、paid promotion、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

currency wallet PostgreSQL adapter gate 记录是：

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

`W-0296` 已实现 storage-neutral `runtime/internal/modules/currency.Repository` interface。下一步有用边界是 platform adapter gate，未来它会把该 interface 映射到现有 PostgreSQL `currency_wallets`、`currency_wallet_balances` 和 `currency_wallet_transactions` tables。

该 gate 在写入任何 adapter SQL 前记录未来 implementation shape：

- adapter ownership；
- constructor 和 executor handoff expectations；
- transaction 和 unit-of-work boundaries；
- wallet creation、wallet lookup、owner lookup、balance listing、grant recording、spend recording、transaction listing 的 SQL mapping posture；
- idempotency、uniqueness、version、affected-row、insufficient-balance 和 driver-error mapping；
- wallet、transaction、actor、metadata、SQL、DSN 和 driver details 的 redaction posture；
- focused test expectations；
- 将 runtime、protocol、generated output、integrations、hosted surfaces、distributed runtime 和 direct compatibility 留在 adapter gate 范围外的 stop conditions。

这不是实现。未来 adapter source path 只是命名，方便 agent 用 accepted boundary 校验后续工作。

## 3. Ownership

未来 adapter owner 是：

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

规则：

- adapter 后续可以在 PostgreSQL platform package 下实现 `currency.Repository`。
- adapter 不得把 SQL 移到 `runtime/internal/modules/currency`。
- adapter 不拥有 request authentication、player identity validation、route policy、protocol parsing、WebSocket state、reward execution、inventory grants、purchase validation、catalog ownership、settlement、refund、transfer 或 distributed topology。
- adapter 必须接收 already-normalized repository input，或在 SQL mapping 前调用 currency module normalizers。
- adapter 必须返回 currency module value types 和 typed repository errors，而不是 driver-specific errors。
- adapter 可以把 player-account foreign-key outcomes 作为 storage conflicts 处理，但不得成为 player account lifecycle owner。

## 4. Future Constructor And Executor Handoff

第一版 adapter implementation 应沿用现有 PostgreSQL adapter patterns：

```yaml
future_constructor_candidate: NewCurrencyWalletRepositoryForUnitOfWork
future_repository_interface: runtime/internal/modules/currency.Repository
executor_source: caller_supplied
transaction_control_sql_allowed: false
unit_of_work_handoff_required: true
connection_pool_owned_by_adapter: false
context_required: true
```

规则：

- constructor 应接受 unit-of-work boundary 提供的 executor 或 query interface，而不是直接拥有 pool。
- adapter 不得发出 `BEGIN`、`COMMIT` 或 `ROLLBACK`；transaction ownership 保持在 unit-of-work runner。
- adapter 不得创建 automatic startup migrations。
- 如果现有 PostgreSQL platform dependencies 已覆盖 implementation，则 adapter 不得添加新 dependency。
- 任何必要 dependency change 都必须有单独 dependency-adoption decision。

## 5. SQL Mapping Posture

未来 adapter 可以把 repository methods 映射到 currency wallet tables：

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

方法姿态：

- `CreateCurrencyWallet` 应为 already validated owner 插入 wallet row，保留 repository validation，从正 wallet version 开始，把 duplicate owner rows 映射为 typed conflicts，并返回 normalized wallet record。
- `GetCurrencyWallet` 应按 `wallet_id` 查询并返回 normalized `CurrencyWallet`。
- `GetCurrencyWalletForOwner` 应按 `(owner_kind, owner_id)` 查询，且不得变成 arbitrary wallet search 或 admin inspection。
- `ListCurrencyWalletBalances` 应按 wallet scope 查询，按 `currency_code` deterministic ordering，遵守 repository limits，并保持 pagination-ready。
- `RecordCurrencyGrant` 应在 caller unit of work 中插入一条 transaction fact 并 update 或 create 一条 balance row。它必须记录 positive grant deltas，保持 non-negative `balance_after`，处理 idempotency uniqueness，并返回 normalized transaction。
- `RecordCurrencySpend` 应在 caller unit of work 中插入一条 transaction fact 并 update 一条 balance row。它必须记录 negative spend deltas，在不泄露 private balance details 的前提下拒绝 insufficient balance，处理 idempotency uniqueness，并只在 balance update 成功后返回 normalized transaction。
- `ListCurrencyWalletTransactions` 应按 wallet scope 查询，可选择 currency filter，按 creation time 和 transaction id deterministic ordering，遵守 repository limits，并保持 pagination-ready。

规则：

- SQL text 必须留在 PostgreSQL adapter package。
- Wallet ids、owner ids、transaction ids、actor ids、idempotency keys/scopes、reason codes、external references、metadata JSON、balances 和 transaction history 默认都不是 log-safe。
- Driver-specific constraint names 可以内部用于 error mapping，但 public module errors 必须保持 currency-module neutral。
- Update operations 必须检查 affected-row counts。
- Balance mutations 不得把 transaction insertion 和 balance update 拆到不同 application units of work。
- 除非后续 retention decision 明确授权，adapter 不得 hard-delete wallet 或 transaction history。

## 6. Transaction And Unit-Of-Work Boundary

未来 adapter 参与现有 transaction handoff：

```yaml
unit_of_work_handoff_required: true
adapter_starts_transactions: false
adapter_commits_transactions: false
adapter_rolls_back_transactions: false
adapter_safe_for_existing_runner: true
balance_mutation_requires_caller_unit_of_work: true
```

规则：

- Application services 或 runtime composition 后续可以通过 explicit unit-of-work boundary 获取 adapter。
- 本 gate 不添加该 factory 或 composition。
- Adapter methods 必须使用 caller context。
- Adapter 不得通过返回成功 wallet 或 transaction results 隐藏 transaction failures。
- Adapter 不得执行 route policy、session validation、access-token validation、WebSocket close behavior、reward execution、inventory mutation、purchase validation、settlement 或 refund behavior。

## 7. Error Mapping

未来 adapter 必须把 PostgreSQL details 折叠成 currency module errors：

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

Mapping expectations：

- `currency_wallets_owner_uq` conflicts 应映射为 duplicate-wallet-owner 类 conflicts。
- `currency_wallet_transactions_idempotency_uq` conflicts 应根据后续 implementation 的 payload comparison posture 映射为 duplicate transaction 或 conflicting duplicate transaction。
- Wallet owner 或 transaction wallet reference 的 foreign-key failures 应映射为 owner-not-found、wallet-not-found、invalid input 或 storage conflict，不暴露 database constraint names。
- 带 expected wallet 或 balance version 的 no affected row 应映射为 missing wallet、missing balance、version mismatch、stale version、lifecycle conflict 或 insufficient balance，不泄露 hidden wallet details。
- Malformed input 应尽可能在 SQL execution 前拒绝。
- Unknown driver 或 executor failures 应映射为 storage-unavailable 类 errors，并使用 redacted reasons。
- Raw SQL、DSNs、credentials、token material、verifier digests、wallet ids、owner ids、transaction ids、actor ids、idempotency material、balances、metadata 和 driver stack details 不得出现在 public error strings。

## 8. Test Expectations

后续 implementation slice 应添加 focused PostgreSQL adapter tests。没有 live database verification 前，fake-executor 或 query-capture tests 可以接受。

Implementation slice 的 required test families：

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

规则：

- Tests 不得要求 protocol routes。
- Tests 不得要求 WebSocket transport。
- Tests 不得打印 raw credentials、tokens、verifier keys、digests、DSNs、query strings、authorization headers、cookies、player ids、wallet ids、transaction ids、idempotency material、balances 或 private wallet state。
- Live PostgreSQL verification 后续可在 implementation slice 授权后添加；如果不可用，必须明确记录。

## 9. Relationship To Runtime, Protocol, And Integrations

本 gate 不改变 runtime 或 protocol behavior：

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

规则：

- 未来 runtime wallet behavior 需要后续 runtime behavior gate。
- 未来 Protobuf messages 和 routes 需要后续 protocol gate。
- 未来 reward、inventory、purchase、payment、catalog、transfer、refund 或 settlement behavior 需要单独 bounded work。
- 本 gate 不得改变 authentication/session behavior。Metadata-only `player_id` 和 `session_id` 仍然不是足够 proof。

## 10. Reference Alignment

Nakama 和 Hiro 指导 durable currency and wallet state 的 product capability need。vibit 通过 explicit storage-neutral repository 和 PostgreSQL adapter boundary 适配该能力。Pitaya 仍是 future distributed architecture reference。

本 gate 拒绝：

- direct Nakama API compatibility；
- direct Pitaya API compatibility；
- direct Hiro API compatibility；
- runtime economy handlers；
- catalog、reward、inventory、purchase、payment、settlement、refund、reservation 和 transfer behavior；
- distributed runtime behavior。

## 11. Next Work Item

本 gate 打开：

```yaml
future_adapter_implementation_work_item: W-0298
future_adapter_implementation_direction: implement_currency_wallet_postgresql_adapter
```

`W-0298` 只能在本文记录的 source、test、constructor、unit-of-work、SQL mapping、redaction 和 deferral posture 内实现 currency wallet PostgreSQL adapter。
