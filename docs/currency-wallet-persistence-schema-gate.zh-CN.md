# Currency Wallet Persistence Schema Gate

状态：Accepted v0.1
最后更新：2026-06-07
范围：在 migration source、balance tables、wallet transaction behavior、runtime behavior、protocol routes、generated output、repositories、adapters 或更广义 economy features 之前，定义未来 currency wallet persistence schema 的 gate
依赖：`docs/currency-wallet-lifecycle-boundary-gate.md`、`docs/postgresql-persistence-boundary.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`
规范决策：`ADR-0201`

英文文件 `docs/currency-wallet-persistence-schema-gate.md` 是权威版本。本文是配套简体中文翻译。

本文定义 currency wallet persistence schema gate。它是 gate artifact。它不添加 SQL migration source，不创建 currency catalog、wallet、balance、transaction、ledger、idempotency、reward、purchase、inventory 或 audit tables，不实现 currency wallet runtime behavior，不添加 protocol routes，不添加 Protobuf source 或 generated output，不添加 dependencies，不添加 repository interfaces，不添加 PostgreSQL adapters，不接入 startup，不改变 authentication/session behavior，不发布 SDK 或 generated client libraries，不创建 hosted deployments 或 release artifacts，不添加 distributed runtime behavior，也不添加 direct Nakama/Pitaya API compatibility。

## 1. 核心规则

currency wallet persistence schema gate 记录如下：

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

## 2. 产品意图

`ADR-0200` 已定义未来 currency wallet lifecycle vocabulary：catalog reads、wallet reads、balance reads、grants、spends、balance change recording 和 transaction reads。下一步是 data-first：balance mutation correctness 需要在 SQL、repositories、adapters、protocol 或 runtime handlers 出现前先固定 durable persistence posture。

本 gate 让未来 migration 可检查：

- table candidates 已知；
- player wallet identity 明确；
- balance row identity 与 wallet identity 分离；
- transaction 或 ledger records 是未来 mutation facts 的必需姿态；
- idempotency 锚定在 durable transaction records 上；
- currency catalog ownership 明确并继续延后；
- indexes、uniqueness、timestamps 和 redaction 在实现前被规划；
- future repository 和 PostgreSQL adapter ownership candidates 被命名。

该 gate 保持保守。它准备下一项 migration-source-only slice，但不添加 migration 文件。

## 3. 选定存储与表

第一版 currency wallet persistence target 是 PostgreSQL：

```yaml
selected_first_currency_wallet_store: postgres
future_migration_source_candidate: runtime/migrations/postgres/000008_create_currency_wallets.sql
future_repository_boundary: separate_future_work_item
future_postgresql_adapter: separate_future_work_item
```

第一版 schema candidate 可定义三个 logical runtime tables：

```yaml
future_tables:
  - currency_wallets
  - currency_wallet_balances
  - currency_wallet_transactions
```

currency catalog table 继续延后：

```yaml
future_currency_catalog_logical_table: deferred
currency_code_source: future_currency_catalog_or_closed_catalog_configuration
catalog_table_added_by_this_gate: false
```

理由：

- PostgreSQL 是 vibit 第一项已接受的 authoritative durable store。
- wallets 与 balances 需要在 grant/spend behavior 存在前具备 transactionally consistent 的姿态。
- transaction table 让 idempotency、auditability 和 support workflows 在 runtime behavior 编写前可检查。
- 独立 catalog table 可能有价值，但第一项 schema gate 不授权 catalog management、rewards、pricing、purchases 或 live-operations behavior。

## 4. 未来 `currency_wallets` 表候选

未来第一项 migration 可定义一个 logical wallet identity table：

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

`wallet_id` 是 server-generated opaque record id。它不是 authentication proof，默认也不是 log-safe。

第一版 logical wallet identity 是：

```text
owner_kind + owner_id
```

第一版 owner posture 只支持 player-owned wallets：

```yaml
owner_kind_first_value: player
owner_id_source: validated_request_identity_player_id
owner_player_fk_candidate: player_accounts(player_id)
owner_kind_check_candidate: owner_kind = 'player'
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
```

global、group、guild、party、match、server shard 或 operations account 等 future owner kinds 需要后续 gate。

## 5. Wallet Lifecycle State Representation

第一版 lifecycle state candidate 跟随 lifecycle gate：

```yaml
lifecycle_state_column: lifecycle_state
lifecycle_state_type: TEXT
allowed_lifecycle_states:
  - active
  - suspended
  - closed
```

规则：

- `active` wallets 可读，并可在未来接受 authorized grants 和 spends。
- `suspended` 与 `closed` 的 mutation behavior 继续延后到 runtime behavior gate。
- `wallet_version` 应为正数并由 server 管理。
- 应约束 `updated_at >= created_at`。
- 应约束 `state_changed_at >= created_at`。
- `closed_at` 与 `suspended_at` 是 persistence facts，不是 authorization proof。

## 6. 未来 `currency_wallet_balances` 表候选

未来第一项 migration 可定义一个 logical balance table：

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

Balance posture：

```yaml
currency_code_type_candidate: TEXT
currency_code_max_length_candidate: 64
amount_type_candidate: BIGINT
amount_unit: integer_minor_unit
negative_balance_allowed_by_default: false
balance_version_type_candidate: BIGINT
```

规则：

- `wallet_id` references `currency_wallets(wallet_id)`。
- `currency_code` 在后续 catalog gate 选择 normalization rules 前保持 case-sensitive。
- `currency_code` 必须 non-blank 且 bounded。
- `balance_amount >= 0` 是第一版默认姿态。
- `balance_version` 必须为正数并由 server 管理。
- row identity 是 `(wallet_id, currency_code)`。
- detailed balances 默认不是 log-safe。

## 7. 未来 `currency_wallet_transactions` 表候选

未来第一项 migration 可定义一个 logical transaction 或 ledger fact table：

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

第一版 transaction kind vocabulary 是：

```yaml
transaction_kinds:
  - grant
  - spend
```

规则：

- `transaction_id` 由 server 生成，默认不是 log-safe。
- `amount_delta` 记录 signed balance movement。
- grant deltas 为正；spend deltas 为负。
- `balance_after` 记录结果余额，用于 supportability 和 auditability。
- `idempotency_key` 与 `idempotency_scope` 是 durable duplicate-prevention fields，默认不是 log-safe。
- `metadata_json` 不得存储 raw credentials、raw tokens、verifier keys、digests、transport metadata、DSNs、payment secrets 或完整 external provider payloads。
- 第一版 schema 应确保重复 idempotency keys 不能在选定 scope 内重复应用 mutation。

## 8. Idempotency 与唯一性

第一版 idempotency posture 是：

```yaml
idempotency_owner: currency_wallet_transactions
idempotency_scope_candidate:
  - wallet_id
  - idempotency_scope
  - idempotency_key
duplicate_application_allowed: false
conflicting_duplicate_payload_public_error: CURRENCY_DUPLICATE_TRANSACTION
```

第一版 uniqueness 与 lookup posture 是：

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

精确 SQL index names 延后到 migration-source slice。

## 9. Catalog And Economy Integration Deferrals

Currency catalog persistence 有意延后：

```yaml
currency_catalog_table_added_by_this_gate: false
reward_integration_added_by_this_gate: false
inventory_integration_added_by_this_gate: false
purchase_behavior_added_by_this_gate: false
paid_currency_behavior_added_by_this_gate: false
```

Future currency catalog work 必须决定：

- code normalization；
- display metadata；
- precision 和 minor-unit semantics；
- 每种 currency 的 allowed grant/spend posture；
- soft-delete 或 deprecation behavior；
- catalog rows 是 runtime mutable、config-owned 还是 migration-owned。

Rewards、purchases、inventory pricing、live-operations grants、refunds、reservations、settlement、transfers 和 paid-currency behavior 仍是独立 future gates。

## 10. Future Repository And Adapter Boundaries

Future repository 和 adapter ownership candidates：

```yaml
future_repository_owner_candidate: runtime/internal/modules/currency
future_repository_interface_candidate: runtime/internal/modules/currency.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_postgresql_adapter_source_candidate: runtime/internal/platform/persistence/postgres/currency_wallet_repository.go
future_postgresql_adapter_tests_candidate: runtime/internal/platform/persistence/postgres/currency_wallet_repository_test.go
```

规则：

- repository interface 不得在后续 repository boundary 或 implementation work item 授权前出现。
- PostgreSQL SQL execution 必须留在 PostgreSQL platform owner package。
- domain/application code 不得 import `pgx`。
- Future mutations 必须运行在 application unit-of-work boundary 内。
- Future repository methods 不得接受 client-supplied actor ids 作为 proof。
- Future adapter errors 必须映射为 redacted module/application errors。

## 11. Redaction

默认不是 log-safe：

- wallet ids；
- player ids；
- transaction ids；
- detailed balances；
- idempotency keys；
- reason codes，当它们泄露 private operations 时；
- external references；
- transaction metadata；
- public contract 批准前的 currency catalog internals。

禁止的 persistence fields：

- raw credentials；
- raw access tokens；
- verifier keys；
- lookup digests；
- verifier digests；
- WebSocket connection ids；
- WebSocket subprotocols；
- remote addresses；
- authorization headers；
- cookies；
- 带 credentials 的 DSNs；
- full payment provider payloads；
- direct Nakama 或 Pitaya API path compatibility markers。

## 12. Future Test Expectations

Future migration-source checks 必须覆盖：

- table names；
- required columns；
- owner identity checks；
- wallet lifecycle state checks；
- balance amount non-negativity；
- positive version checks；
- timestamp checks；
- transaction kind checks；
- idempotency uniqueness；
- wallet owner uniqueness；
- balance uniqueness；
- forbidden secret、digest、transport、payment、Nakama 或 Pitaya compatibility columns。

Future repository 和 adapter tests 必须在 migration source 与 repository/adapter gates 后定义。

Future runtime behavior tests 继续延后，且在 runtime behavior 实现前必须覆盖 validated identity、grant/spend transactionality、duplicate idempotency handling、insufficient balance、invalid currency code、suspended/closed wallet behavior、public error collapse 和 redacted logs。

## 13. 非授权

本 gate 不授权：

- SQL migration source creation；
- currency catalog tables；
- currency wallet tables；
- currency balance tables；
- transaction or ledger tables；
- 独立于 planned transaction posture 的 idempotency tables；
- event/audit tables；
- repository interfaces；
- PostgreSQL adapters；
- runtime wallet behavior；
- grant/spend execution；
- reward integration；
- inventory integration；
- purchase behavior；
- paid currency behavior；
- reservations、settlement、refunds 或 transfers；
- protocol routes；
- Protobuf sources；
- generated output；
- dependencies；
- startup wiring；
- authentication/session behavior changes；
- hosted deployments；
- SDK publication；
- release artifacts；
- distributed runtime behavior；
- direct Nakama/Pitaya API compatibility。

## 14. 下一步

下一项 bounded work item 是：

```text
W-0294 Add currency wallet migration source
```

该 follow-up 可在遵循本 gate 且保持 migration-source-only 的前提下添加 SQL migration source。Repository interfaces、adapters、runtime behavior、protocol routes、generated output、dependencies、startup wiring、rewards、purchases、inventory integration、hosted surfaces、SDKs、distributed runtime 和 direct compatibility 仍属于后续 work。
