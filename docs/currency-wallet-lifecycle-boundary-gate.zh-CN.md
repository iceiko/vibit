# 货币钱包生命周期 Gate

状态：Accepted v0.1
最后更新：2026-06-07
范围：在 persistence、protocol、runtime behavior、rewards、purchases 或更大的 economy features 之前，为未来 currency wallet lifecycle 定义语义 gate
依赖：`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/agent-native-feature-request-test-workflow.md`
Canonical decision: `ADR-0200`

配套英文源文件是 `docs/currency-wallet-lifecycle-boundary-gate.md`。英文文件是权威版本。

本文定义 currency wallet lifecycle semantic gate。它是 gate artifact。它不添加 runtime currency wallet behavior、balance tables、wallet transaction behavior、reward integration、inventory integration、purchase behavior、grant/spend execution、audit/event tables、protocol routes、Protobuf source、generated output、migrations、repository interfaces、PostgreSQL adapters、dependencies、startup wiring、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts、distributed runtime behavior，或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

货币钱包生命周期 gate 记录如下：

```yaml
currency_wallet_lifecycle_boundary_gate: defined
completed_work_item: W-0292
decision: ADR-0200
check_rule: runtime.currency_wallet_lifecycle_boundary_gate
selection_decision: ADR-0199
gate_standard: docs/currency-wallet-lifecycle-boundary-gate.md
gate_standard_translation: docs/currency-wallet-lifecycle-boundary-gate.zh-CN.md
selected_capability_family: economy_inventory_rewards_currencies_and_progression
selected_module_candidate: currency
primary_product_reference: Nakama
secondary_product_reference: Hiro
pitaya_reference_status: deferred_future_architecture_reference
semantic_gate_only: true
future_persistence_schema_gate_work_item: W-0293
future_persistence_schema_gate_direction: define_currency_wallet_persistence_schema_gate
permission_required: validated_player_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Intent

Currency wallets 是 Nakama-class 和 Hiro-class game backends 的核心 economy primitive。它后续可以支撑 rewards、grants、spends、purchases、progression、live operations、inventory pricing、leaderboards、tournaments、quests 和 player support workflows。vibit 采用产品能力族，不采用任何外部 public API shape。

vibit 的姿态是：

- implementation 之前先定义 contract-first lifecycle semantics；
- balance mutations 必须 server-authoritative；
- player-owned wallet reads 和 player-initiated spends 需要 validated player identity；
- 未来 grants 和 administrative adjustments 需要 service-authoritative 姿态；
- 在任何 balance mutation 存在之前先定义 idempotency 和 transaction boundaries；
- detailed balances、ledger records、transaction ids、wallet ids 和 player ids 默认不是 log-safe；
- 除非后续 ADR 明确授权，否则拒绝 direct external API compatibility。

Pitaya 仍然只是 future distributed architecture reference。本 gate 不能引入 distributed routing、frontend/backend roles、RPC、cluster groups、service discovery 或 server-to-server messaging。

## 3. Future Semantic Scope

未来 currency wallet lifecycle 必须覆盖：

```yaml
semantic_scope:
  - currency_catalog_read
  - wallet_read
  - balance_read
  - grant_currency
  - spend_currency
  - balance_change_recording
  - transaction_read
```

生命周期是 player-wallet-oriented 并且 server-authoritative。未来 domain owner 是 currency/economy capability boundary，而不是 WebSocket transport、protocol adapters、authentication、storage objects、inventory、friends、rewards、purchases、leaderboards、matchmaking、match runtime、operations dashboards 或 distributed runtime。

Transfers、reservations、settlement、refunds、player-to-player exchange、store purchases、paid currency 和 live-operations grants 仍是未来决策。它们可能需要在基础 wallet lifecycle 和 persistence schema 稳定后单独定义 gate。

## 4. Future Contract Vocabulary

未来 command vocabulary 是：

```yaml
commands:
  - GrantCurrency
  - SpendCurrency
```

未来 query vocabulary 是：

```yaml
queries:
  - GetCurrencyWallet
  - ListCurrencyBalances
  - GetCurrencyTransaction
```

未来 event vocabulary 是：

```yaml
events:
  - CurrencyGranted
  - CurrencySpent
  - CurrencyBalanceChanged
  - CurrencyTransactionRecorded
```

未来 error vocabulary 是：

```yaml
errors:
  - CURRENCY_INVALID_CODE
  - CURRENCY_WALLET_NOT_FOUND
  - CURRENCY_INSUFFICIENT_BALANCE
  - CURRENCY_INVALID_AMOUNT
  - CURRENCY_DUPLICATE_TRANSACTION
  - CURRENCY_INVALID_TRANSITION
  - CURRENCY_METADATA_IDENTITY_NOT_AUTHENTICATED
```

这些 vocabulary 只是语义规划。它不创建 contract source files、generated shapes、protocol payloads、routes、repositories、migrations、adapters 或 runtime handlers。

## 5. Identity And Permissions

未来 player-owned wallet reads 和 player-initiated spends 要求：

```yaml
permission: validated_player_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
actor_identity_source: validated_request_identity
```

未来 grants、administrative adjustments、reward handoffs、purchase settlement、refunds 和 live-operations actions 必须在 implementation 之前由后续 service-authoritative permission model 定义。本 gate 除了拒绝 client-authoritative grants 外，不定义该模型。

规则：

- player actor 来自 validated request identity，不来自 client-supplied actor id。
- `player_id` 和 `session_id` metadata 不是 authentication proof。
- player-supplied target wallet ids 不是修改其他钱包的授权。
- Grants 永远不能 client-authoritative。
- 当 wallet existence、player existence、currency catalog configuration 或 private balance information 可能泄露时，public failures 必须 collapse 细节。
- wallet records、player ids、wallet ids、transaction ids、detailed balances、idempotency keys 和 ledger details 默认不是 log-safe。

## 6. Lifecycle State And Invariants

未来 public wallet lifecycle vocabulary 是：

```yaml
wallet_lifecycle_states:
  - active
  - suspended
  - closed
```

未来 currency balance posture 是：

```yaml
balance_posture:
  amount_unit: integer_minor_unit
  negative_balance_allowed_by_default: false
  currency_code_source: future_currency_catalog
  mutation_idempotency_required: true
  mutations_transactional: true
```

第一批 lifecycle invariants 是：

- Currency codes 必须在 mutation 前通过未来 currency catalog 验证。
- Grants 和 spends 的 amount 必须为 positive。
- Balance reads 不能要求 mutation authority。
- 除非后续 ADR 明确授权 overdraft behavior，否则 balances 绝不能变成 negative。
- 每个未来 mutation 都必须携带 idempotency key 或等价 duplicate-prevention key。
- 每个未来 mutation 都必须记录 durable transaction 或 ledger fact，之后才能认为 behavior 已实现。
- Grant 和 spend behavior 必须把 balance changes 与 transaction recording 放在同一个事务边界内。
- Spend 只适用于未来 active wallet 且 balance sufficient 的情况。
- Grant 只适用于后续 permission gate 定义的 server/service-authoritative paths。
- Suspended 或 closed wallet behavior 必须在 runtime implementation 前明确。
- Public query output 不能暴露超出未来 route contract 的 internal ledger details。

## 7. Future Persistence Gate

下一项 bounded work item 是：

```text
W-0293 Define currency wallet persistence schema gate
```

该 follow-up 应在 migration source 出现前定义 table candidates、currency catalog posture、wallet account identity、balance row identity、ledger 或 transaction record posture、idempotency key uniqueness、indexes、retention、redaction 和 repository/adapter boundaries。

本 gate 有意不决定：

- exact table names；
- 第一版 persistent model 是 balances plus ledger、ledger-derived balances，还是其他模式；
- currency catalog table ownership；
- idempotency key storage shape；
- transaction id format；
- reservation 或 settlement tables；
- reward、purchase、inventory 或 live-operations integration；
- repository interface shape；
- PostgreSQL adapter SQL；
- protocol routes 或 payloads。

## 8. Future Test Expectations

未来 behavior tests 必须在 implementation 前规划。

Positive tests：

- 通过 service-authoritative path grant currency；
- 为 validated player wallet spend currency；
- read one player-owned wallet；
- list player-owned currency balances；
- read a currency transaction by permitted owner；
- 通过 idempotency key 拒绝 duplicate mutations，且不能 double-apply balance changes。

Negative tests：

- invalid currency code；
- zero 或 negative mutation amount；
- insufficient balance；
- missing wallet；
- suspended 或 closed wallet transition；
- duplicate transaction with conflicting payload；
- metadata-only identity。

Permission and authentication tests：

- player-owned reads 和 spends 要求 validated player identity；
- client-supplied actor id 会被 ignored 或 rejected；
- metadata-only `player_id` 和 `session_id` 作为 proof 会被拒绝；
- client-authoritative grant attempts 在 mutation 前被拒绝。

Persistence and transaction tests：

- schema 和 repository tests 必须在 `W-0293` 之后、migration/adapter/runtime implementation 之前定义；
- mutation transitions 必须 transactional；
- balance changes 和 transaction records 必须在未来 unit-of-work boundary 内保持一致；
- duplicate idempotency keys 不能 double-apply mutations。

Failure and redaction tests：

- 在 privacy 要求 collapse 的地方，public errors 不泄露 private wallet、balance、catalog 或 transaction details；
- logs 不暴露 raw credentials、tokens、verifier keys、digests、transport metadata、wallet ids、player ids、detailed balances、idempotency keys 或 ledger internals。

Concurrency tests：

- simultaneous grant/spend/idempotency conflicts 必须在 runtime implementation 之前有明确 expected outcomes。

Integration and end-to-end tests：

- 推迟到 protocol routes 和 runtime handlers 获得授权后。

## 9. Non-Authorization

本 gate 不授权：

- runtime currency wallet behavior；
- balance tables；
- wallet transaction behavior；
- reward integration；
- inventory integration；
- purchase behavior；
- grant execution；
- spend execution；
- transfer execution；
- reservation 或 settlement behavior；
- audit/event tables；
- protocol routes；
- Protobuf source；
- generated output；
- migrations；
- repository interfaces；
- PostgreSQL adapters；
- startup wiring；
- authentication/session behavior changes；
- dependencies；
- SDK publication；
- hosted deployment；
- release artifacts；
- distributed runtime behavior；
- direct Nakama/Pitaya API compatibility。

任何未来 work 如果需要上述任一 surface，都必须创建单独 bounded work item，并通过自己的 repository check。
