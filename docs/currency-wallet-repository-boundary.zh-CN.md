# Currency Wallet Repository Boundary

Status: Accepted v0.1
Last updated: 2026-06-07
Scope: 在 PostgreSQL `currency_wallets`、`currency_wallet_balances` 和 `currency_wallet_transactions` migration source 之后，为未来 storage-neutral currency wallet repository 定义 gate-only boundary
Depends on: `docs/currency-wallet-lifecycle-boundary-gate.md`, `docs/currency-wallet-persistence-schema-gate.md`, `decisions/ADR-0202-currency-wallet-migration-source.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0203`

配套英文原文是 `docs/currency-wallet-repository-boundary.md`。英文文件是权威版本。

本文定义 currency wallet repository boundary。它是 gate artifact。它不添加 Go repository interfaces、PostgreSQL adapter behavior、runtime wallet behavior、grant/spend execution、protocol routes、Protobuf source、generated output、dependencies、migrations、automatic startup migration behavior、reward integration、inventory integration、purchase behavior、paid-currency behavior、catalog tables、event/audit tables、reservations、settlement、refunds、transfers、hosted deployments、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Currency wallet repository boundary record 是：

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

`W-0294` 添加了 `currency_wallets`、`currency_wallet_balances` 和 `currency_wallet_transactions` 的 PostgreSQL migration source。下一个有用边界是 storage-neutral repository vocabulary，让未来实现代码可以在不暴露 SQL 细节、transport 细节、protocol 假设或 direct external economy compatibility 的情况下使用。

本 boundary 为 Nakama/Hiro-class economy capability coverage 记录：

- repository ownership；
- candidate wallet、balance、transaction、actor、idempotency 和 pagination value types；
- repository command/query vocabulary；
- validated identity 和 actor handoff rules；
- transaction、idempotency 和 optimistic version posture；
- conflict、redaction 和 error posture；
- 与 W-0294 migration 绑定的 PostgreSQL adapter expectations；
- future implementation work 的 stop conditions。

这仍不是 runtime feature。后续 bounded work item 明确授权前，route、handler、adapter、repository interface 或 protocol message 都不能使用 currency wallets。

## 3. Ownership

未来 repository 由 currency module 拥有：

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

规则：

- Future repository interface 必须是 storage-neutral 且面向 module。
- Interface 不得提及 PostgreSQL、pgx、SQL rows、goose migrations、prepared statements、connection pools、transaction runners 或 database driver errors。
- PostgreSQL adapter 后续可以在 `runtime/internal/platform/persistence/postgres` 下实现该 interface，但必须等 separate adapter gate 授权。
- Application 或 handler code 后续必须通过 module/application boundaries 调用 currency wallet behavior，而不是通过 SQL 或 transport state。
- Authentication 和 session code 提供 validated request identity；它们不拥有 wallet records。
- Player account storage 拥有 player lifecycle state，不拥有 currency wallet balances。
- Inventory 拥有 inventory item state，不拥有 wallet balances 或 currency transactions。
- WebSocket transport 拥有 connection plumbing，不拥有 economy state。
- Protocol adapters 拥有 wire conversion，不拥有 repository behavior。

## 4. Candidate Value Types

后续 implementation gate 可以重命名或缩减这些 shape，但第一版 repository interface implementation 应从以下 vocabulary 开始：

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

First-posture record vocabulary：

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

规则：

- `owner_kind` 必须保持 closed vocabulary。第一版允许值是 `player`。
- `owner_id` 是 input identity handoff，不是 proof。未来 runtime behavior 必须先从 validated request identity 派生它，再调用 repository。
- `currency_code` 在后续 catalog gate 选择 normalization 前保持 case-sensitive。
- `balance_amount` 和 `balance_after` 是 non-negative integer minor units。
- `wallet_version` 和 `balance_version` 由 server 管理，不能是 client-authoritative state。
- `transaction_id`、`wallet_id`、detailed balances、idempotency fields、reason codes、external references 和 metadata 默认不是 log-safe。

## 5. Candidate Repository Capabilities

第一版 storage-neutral capability family 是：

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

Capability rules：

- `CreateCurrencyWallet` 可以为已经验证的 owner identity 创建 wallet row。
- `GetCurrencyWallet` 是按 server wallet id 的 storage lookup。它不得 authenticate users、validate access tokens 或 create request identity。
- `GetCurrencyWalletForOwner` 必须 owner-scoped；没有后续 gate，不得变成 arbitrary wallet search 或 admin inspection。
- `ListCurrencyWalletBalances` 必须 wallet-scoped，并在未来 catalog size 需要时保持 pagination-ready。
- `RecordCurrencyGrant` 必须在同一个 application unit of work 中插入一个 transaction fact 并 update 或 create 一个 balance row。
- `RecordCurrencySpend` 必须在同一个 application unit of work 中插入一个 transaction fact 并 update 一个 balance row，同时保留 insufficient-balance conflict posture。
- `ListCurrencyWalletTransactions` 必须 wallet-scoped 且 pagination-ready。没有后续 protocol/redaction gate，不得暴露 private support metadata。
- 所有方法必须返回 module-owned typed records 和 errors，而不是 raw SQL rows 或 database driver errors。

未来 repository interface 可以选择更短命名，但必须保留 wallet creation、wallet lookup、balance listing、grant recording、spend recording、transaction listing、idempotency 和 conflict handling 的语义拆分。

## 6. Identity, Transaction, And Idempotency Handoff

Repository boundary 只准备 identity、transaction 和 idempotency handling，不实现行为：

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

规则：

- Repository 可以接收 normalized owner 和 actor ids 作为 data，但 ids 不是 authentication proof。
- 未来 runtime behavior 必须在调用 grant 或 spend repository capability 前验证 actor permissions。
- 同一 wallet 和 idempotency scope 下的 duplicate idempotency keys 不能 double-apply balance mutation。
- Conflicting duplicate payload behavior 留给未来 runtime/protocol gate，但不能静默创建第二个 transaction。
- Currency wallet mutations 必须在 application unit-of-work boundary 内运行。
- Repository methods 不得接受 transport metadata、WebSocket connection identifiers、cookies、headers、tokens 或 sessions 作为 identity proof。

## 7. Version And Conflict Handoff

Repository boundary 只准备 optimistic concurrency，不实现行为：

```yaml
wallet_version_storage: BIGINT
balance_version_storage: BIGINT
initial_create_version: 1
version_owner: server
client_authoritative_version_allowed: false
expected_version_handoff: future_behavior_or_interface_gate
conflict_public_shape: deferred_to_protocol_gate
```

Candidate conflict classes：

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

规则：

- Repository methods 后续可以区分 internal typed conflicts，但 public protocol error mapping 仍然 deferred。
- Version equality 不是 authentication proof。
- Stale expected version 不能被折叠成隐藏的成功写入。
- 当未来 public behavior gates 要求泄漏折叠时，owner mismatch 不得泄漏另一个玩家 wallet 的存在。
- PostgreSQL adapter 必须把 unique-index、affected-row、check-constraint 和 foreign-key outcomes 映射为 typed repository conflicts，不能暴露 driver error text。

## 8. Redaction And Logging

Currency wallet state 是 private economy data。

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

规则：

- Repository errors 必须 redacted 且 typed。
- Raw values 和 storage driver errors 默认不得记录日志。
- Transaction metadata 可能包含 private support context，默认不得记录日志。
- Repository 不得存储或返回 authentication material、token material、verifier digests、transport metadata、payment secrets、full provider payloads 或 direct external compatibility markers。

## 9. PostgreSQL Adapter Expectations

未来 PostgreSQL adapter 后续可以把 repository 映射到：

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

Adapter rules：

- SQL 必须留在 `runtime/internal/platform/persistence/postgres`。
- Domain/application code 不得 import pgx。
- Adapter 必须在 application unit-of-work boundary 内运行。
- Balance mutation 和 transaction insert 必须 atomic。
- Idempotency uniqueness 必须映射为 typed duplicate transaction outcomes。
- Check constraints 和 foreign-key failures 必须映射为 redacted typed conflicts。
- Driver error messages、SQL statements、DSNs 和 raw metadata 不得暴露在 public errors 中。

## 10. Future Test Expectations

未来 repository interface implementation 必须添加 focused tests，覆盖：

- wallet owner normalization；
- currency code validation posture；
- non-negative amount and signed delta validation；
- idempotency key and scope validation；
- grant/spend input normalization；
- transaction list pagination inputs；
- typed conflict vocabulary；
- redacted repository errors；
- storage neutrality, including no PostgreSQL imports。

未来 PostgreSQL adapter tests 必须覆盖：

- wallet creation and owner uniqueness；
- wallet lookup by id and owner；
- balance listing；
- grant transaction plus balance update atomicity；
- spend transaction plus balance update atomicity；
- insufficient balance conflict mapping；
- duplicate idempotency handling；
- lifecycle state filtering；
- version conflict handling；
- redaction of database errors。

未来 runtime behavior tests 仍然 deferred，必须覆盖 validated identity、permissions、suspended and closed wallet behavior、duplicate request behavior、public error mapping 和 protocol redaction。

## 11. Stop Conditions

在做以下任一事项前必须停止并创建新 gate：

- 添加 Go repository interface implementation；
- 添加 PostgreSQL adapter code；
- 从 Go 执行 SQL；
- 添加 runtime wallet behavior；
- 添加 grant 或 spend command execution；
- 添加 reward、inventory、purchase、paid-currency、reservation、settlement、refund 或 transfer behavior；
- 添加 currency catalog 或 event/audit tables；
- 添加 protocol routes；
- 添加 Protobuf source 或 generated output；
- 改变 authentication/session behavior；
- 添加 dependencies；
- wiring startup；
- 发布 SDKs 或 hosted surfaces；
- 添加 distributed runtime behavior；
- 添加 direct Nakama/Pitaya API compatibility。

## 12. Next Step

下一个 bounded work item 是：

```text
W-0296 Implement storage-neutral currency wallet repository interface
```

该 follow-up 可以在遵循本 boundary 且保持 interface-only 时添加 storage-neutral Go repository interface。PostgreSQL adapters、runtime behavior、protocol routes、generated output、dependencies、startup wiring、rewards、purchases、inventory integration、hosted surfaces、SDKs、distributed runtime 和 direct compatibility 仍是后续工作。
