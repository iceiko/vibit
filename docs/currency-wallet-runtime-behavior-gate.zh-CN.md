# Currency Wallet Runtime Behavior Gate

Status: Accepted v0.1
Last updated: 2026-06-07
Scope: PostgreSQL adapter 之后 future application-owned currency wallet runtime behavior 的 gate-only boundary
Depends on: `docs/currency-wallet-lifecycle-boundary-gate.md`, `docs/currency-wallet-repository-boundary.md`, `docs/currency-wallet-postgresql-adapter-gate.md`, `runtime/internal/modules/currency/repository.go`, `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`, `docs/runtime-protocol-adapter.md`, `docs/bound-identity-route-policy-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0207`

This Simplified Chinese document translates `docs/currency-wallet-runtime-behavior-gate.md`. The English file is authoritative.

本文定义 currency wallet runtime behavior gate。它是 gate artifact。本文件不添加 runtime behavior implementation、runtime handlers、startup wiring、protocol routes、Protobuf source、generated output、repository interface changes、PostgreSQL adapter changes、migration changes、dependencies、authentication/session behavior changes、reward integration、inventory integration、purchase behavior、catalog tables、event/audit tables、payment behavior、reservation behavior、settlement behavior、refund behavior、transfer behavior、operations/admin behavior、SDK publication、generated client libraries、hosted deployments、release artifacts、public announcements、paid promotion、Pitaya-style distributed architecture，或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Currency wallet runtime behavior gate record 是：

```yaml
currency_wallet_runtime_behavior_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0299
decision: ADR-0207
check_rule: runtime.currency_wallet_runtime_behavior_gate
source_postgresql_adapter_implementation_decision: ADR-0206
source_postgresql_adapter: runtime/internal/platform/persistence/postgres/currency_wallet_repository.go
source_repository_interface_decision: ADR-0204
repository_interface: runtime/internal/modules/currency.Repository
repository_interface_source: runtime/internal/modules/currency/repository.go
future_runtime_owner_candidate: runtime/internal/app
future_currency_application_package_candidate: runtime/internal/app/currency
future_runtime_service_source_candidate: runtime/internal/app/currency/service.go
future_runtime_service_test_candidate: runtime/internal/app/currency/service_test.go
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
owner_identity_source: validated_request_identity_player_id
actor_identity_source: validated_request_identity_or_server_operation
first_owner_kind: player
first_actor_kinds:
  - player
  - system
route_policy_requirement: request_token_required
service_application_owner: runtime/internal/app
repository_handoff: unit_of_work_currency_wallet_repository_factory
unit_of_work_handoff_required: true
runtime_behavior_gate_only: true
runtime_behavior_added: false
runtime_handlers_added: false
startup_wiring_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
reward_integration_added: false
inventory_integration_added: false
purchase_behavior_added: false
currency_catalog_table_added: false
currency_wallet_events_table_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_runtime_behavior_implementation_work_item: W-0300
future_runtime_behavior_implementation_direction: implement_currency_wallet_runtime_behavior
```

## 2. Purpose

`W-0298` 已为 `runtime/internal/modules/currency.Repository` 实现 PostgreSQL adapter。下一步有价值的 boundary 不是 protocol route，也不是 reward integration，而是 runtime behavior gate：先定义 future application code 如何把 validated player request 或 server-authoritative operation 转换为 currency wallet repository operations。

这个 gate 在实现前记录：

- service 的 application ownership；
- wallet owner identity 来自 validated request identity；
- player-initiated 与 server-authoritative operation 的 actor derivation；
- wallet creation、balance reads、transaction reads、grants、spends 的 command/query posture；
- permission 与 route-policy posture；
- validation 与 conflict mapping expectations；
- unit-of-work 与 repository handoff；
- idempotency 与 redaction rules；
- fake-service test expectations；
- stop conditions，确保 protocol、generated output、authentication/session、reward、inventory、purchase、catalog、event/audit、payment、distributed runtime scope 不进入本 slice。

Nakama 提供 durable wallets、balances、grants、spends、transaction history 的 capability pressure。Pitaya 提醒 route/session context、handlers、persistence responsibilities 需要分离。vibit 通过 explicit application-owned behavior 和 checks 吸收这些参考，而不是复制 direct public API compatibility。

## 3. Ownership

Future runtime behavior 属于 application：

```yaml
future_runtime_owner_candidate: runtime/internal/app
future_currency_application_package_candidate: runtime/internal/app/currency
future_runtime_service_source_candidate: runtime/internal/app/currency/service.go
future_runtime_service_test_candidate: runtime/internal/app/currency/service_test.go
repository_interface_owner: runtime/internal/modules/currency
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
player_account_owner: runtime/internal/modules/player
```

Rules:

- Future service behavior 可以放在 `runtime/internal/app/currency`，或放在 implementation slice ratify 的等价 application-owned package。
- Service 只能通过 application 或 unit-of-work dependencies 调用 `runtime/internal/modules/currency.Repository`。
- 只有后续 implementation slice 明确授权 dependency handoff 时，service 才可通过现有 application-owned player repository capability 检查 player account existence 和 account state。
- Service 不得 import PostgreSQL adapter packages、SQL row types、migration packages、WebSocket transport packages、generated Protobuf packages、reward packages、inventory packages、purchase packages、payment provider SDKs、catalog packages、event/audit packages、SDK packages，或 distributed runtime packages。
- Currency module 继续拥有 storage-neutral wallet value types、normalizers、lifecycle vocabulary、transaction vocabulary、actor vocabulary、idempotency vocabulary 和 repository error vocabulary。
- PostgreSQL adapter 保持 persistence-only，不得派生 request identity、route policy、public protocol errors、reward decisions、inventory decisions、purchase decisions 或 payment decisions。
- Transport 与 protocol adapters 不拥有 currency wallet permission、economy policy 或 business behavior。

## 4. Request Identity And Owner Derivation

第一姿态是 player-wallet owned：

```yaml
first_owner_kind: player
owner_identity_source: validated_request_identity_player_id
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed_as_proof: false
client_supplied_wallet_id_allowed_for_lookup: false
server_wallet_lookup_by_owner_required: true
```

Rules:

- Future player-visible wallet operation 必须从 validated `app.RequestIdentity` 派生 `currency.CurrencyWalletOwner{Kind: player, ID: identity.PlayerID}`。
- `RequestIdentity.Status` 必须是 `validated`。
- `RequestIdentity.ActorKind` 必须是 `player`。
- `RequestIdentity.PlayerIDValidated` 必须为 true。
- `RequestIdentity.PlayerID` 必须非空，并且在同时出现 actor identity 时保持一致。
- Envelope/session metadata 中 metadata-only `player_id` 不能满足本 gate。
- 单独的 persisted `session_id` 不能成为 proof。
- 第一姿态下 client payload 不能选择其他 owner id。
- 第一姿态下 client payload 不应选择任意 wallet id；service 应尽可能通过 server-derived owner 解析 wallet。

本 gate 不改变 `RequestIdentity`、access-token validation、bound connection identity、durable runtime session validation 或 WebSocket handshake behavior。它只记录 future currency behavior 在访问 repository 前必须满足的 identity requirements。

## 5. Actor Derivation

Currency transactions 除 wallet ownership 外还需要 actor：

```yaml
actor_identity_source: validated_request_identity_or_server_operation
first_actor_kinds:
  - player
  - system
operation_actor_kind_reserved: true
client_supplied_actor_id_allowed_as_proof: false
```

Rules:

- Player-initiated wallet reads 和 spends 必须从 validated request identity 派生 actor identity。
- Server-authoritative grants 可以使用 service-owned `system` actor，但 implementation slice 必须明确暴露该 dependency 与 reason-code posture。
- `operation` actors 保留给后续 operations/admin behavior，本 gate 不启用。
- Client payload 不得提供 actor id 作为 proof。
- Actor ids、wallet ids、owner ids、idempotency keys、reason codes、external references、transaction ids 默认不 log-safe。

## 6. Future Runtime Behavior Shape

Future first implementation 可以暴露 application service，候选 operations：

```yaml
candidate_operations:
  - ensure_player_wallet
  - get_own_wallet
  - list_own_wallet_balances
  - grant_currency
  - spend_currency
  - list_own_wallet_transactions
```

Recommended first posture:

- `ensure_player_wallet` 通过 caller-owned unit of work 为 validated player 创建或查找 active wallet。
- `get_own_wallet` 读取 validated player 的 wallet。
- `list_own_wallet_balances` 使用 bounded pagination 列出 validated player's wallet balances。
- `grant_currency` 记录带 idempotency key、reason code 和 optional metadata 的 server-authoritative grant，不集成 reward、inventory、purchase、catalog 或 payment behavior。
- `spend_currency` 记录 player-authorized spend，输入 positive amount 和 idempotency key，并映射 insufficient-balance conflict；不集成 purchase、reservation、settlement、refund、transfer 或 paid-currency behavior。
- `list_own_wallet_transactions` 使用 bounded pagination 和 optional currency filter 列出 validated player's wallet transactions。

Rules:

- Runtime behavior 必须使用 server-derived owner 与 actor identity。
- Runtime behavior 应尽可能在 repository call 前验证 currency code、positive amount、idempotency fields、reason code、external reference、metadata JSON、expected wallet version、expected balance version、list limit 和 pagination cursor。
- Runtime behavior 不得在 first implementation 中暴露 cross-owner wallet reads、arbitrary wallet id lookup、admin wallet inspection、player-to-player transfer、paid-currency purchase、reward inventory integration、catalog lookup、exchange rates、reserves、settlement、refunds 或 payment provider behavior。
- Runtime behavior 不得添加 public protocol routes 或 generated output，除非后续 protocol gate 授权。

## 7. Candidate Application Service Shape

First implementation slice 应定义小型 application-owned service boundary。候选 input/output：

```yaml
candidate_request_fields:
  - request_identity
  - currency_code
  - amount
  - idempotency_key
  - idempotency_scope
  - reason_code
  - external_reference
  - metadata_json
  - expected_wallet_version
  - expected_balance_version
  - list_limit
  - after_currency_code
  - after_transaction_id
  - after_transaction_time
candidate_result_fields:
  - wallet
  - balance
  - balances
  - transaction
  - transactions
  - next_currency_code
  - next_transaction_id
  - next_transaction_time
  - public_error_code
```

Rules:

- Service 应接受 already-normalized application identity，而不是 raw tokens、cookies、headers、WebSocket subprotocol values 或 envelope proof carriers。
- Service 应在 repository handoff 前调用 currency module normalizers。
- Service 应避免在默认 errors/logs 中暴露 transaction metadata JSON、idempotency keys、wallet ids、owner ids、actor ids、reason codes、external references 和 platform errors。
- Service 应暴露 stable public error codes/classes，供后续 runtime handlers 映射。
- Gate slice 不得添加 route registration、Protobuf conversion、startup composition、reward composition、inventory composition、purchase composition 或 command/query dispatch wiring。

## 8. Validation Rules

Future runtime behavior 必须在 persistence 前执行 validation：

```yaml
validation_required:
  - validated_player_identity
  - server_derived_wallet_owner
  - active_wallet_required_for_mutation
  - currency_code_non_empty_length_bounded
  - amount_positive_for_grant_and_spend_requests
  - idempotency_key_non_empty_length_bounded
  - idempotency_scope_non_empty_length_bounded
  - reason_code_length_bounded
  - external_reference_length_bounded
  - metadata_json_top_level_object_when_present
  - expected_wallet_version_positive_when_present
  - expected_balance_version_positive_when_present
  - list_limit_bounded
  - pagination_cursor_bounded
```

Rules:

- Currency code 和 amount validation 应复用 currency module normalization rules，除非 future contract 明确收紧 protocol-visible syntax。
- Grant 和 spend requests 应接受 positive amount input；transaction deltas 仍由 repository 拥有。
- Metadata JSON 不是 log-safe，跨 boundary 时必须 copied 或 immutable。
- Missing expected version behavior 必须在 implementation tests 中明确。
- Invalid input 应尽可能在 repository mutation 前失败。
- Repository unavailable errors 必须保持 redacted。

## 9. Permission And Route Policy Posture

第一 route-policy posture 保守：

```yaml
route_policy_requirement: request_token_required
public_currency_wallet_routes_allowed: false
bound_connection_required_by_this_gate: false
session_validated_required_by_this_gate: false
bound_session_required_by_this_gate: false
```

Candidate permission families for later public contracts:

- read own wallet;
- list own wallet balances;
- list own wallet transactions;
- spend own wallet balance through a ratified service command;
- receive server-authoritative currency grants.

Rules:

- Currency wallet routes 必须是 protected routes。
- 第一姿态应使用现有 `request_token_required` route protection family，除非后续 route-policy ADR 改变 named routes。
- Public routes 不得读取或修改 currency wallets。
- Bound connection identity 与 durable session validation 可保留给 future route families，但本 gate 不要求它们，也不改变 ordinary protected route behavior。
- Metadata-only identity 必须 fail closed。

## 10. Conflict And Error Mapping

Future runtime behavior 必须把 currency repository errors 映射为 stable application classes：

```yaml
candidate_public_error_classes:
  - CURRENCY_WALLET_INVALID_REQUEST
  - CURRENCY_WALLET_NOT_FOUND
  - CURRENCY_WALLET_ALREADY_EXISTS
  - CURRENCY_WALLET_NOT_ACTIVE
  - CURRENCY_BALANCE_NOT_FOUND
  - CURRENCY_AMOUNT_INVALID
  - CURRENCY_INSUFFICIENT_BALANCE
  - CURRENCY_TRANSACTION_DUPLICATE
  - CURRENCY_TRANSACTION_CONFLICT
  - CURRENCY_WALLET_VERSION_MISMATCH
  - CURRENCY_BALANCE_VERSION_MISMATCH
  - CURRENCY_WALLET_FORBIDDEN
  - CURRENCY_WALLET_UNAVAILABLE
```

Rules:

- Not-found、owner mismatch、closed wallet、suspended wallet cases 必须避免 cross-player existence leaks。
- Insufficient balance 可以作为 stable conflict class 公开，但不得泄露 balance internals、wallet ids、SQL details、driver errors、DSNs、credentials、token material、verifier digests、route proof carriers。
- Duplicate idempotency 只能作为 stable duplicate/conflicting-duplicate class 公开。
- Repository `storage_unavailable` errors 必须映射到 retryable/unavailable class，且不暴露 platform internals。
- Request identity 未 validated 时，permission failure 必须在 repository access 前发生。
- Runtime behavior 不得在 existing application route-protection classes 之外添加 authentication/token/session failure detail。

## 11. Unit-Of-Work And Repository Handoff

Future runtime behavior 应使用现有 application transaction boundary：

```yaml
repository_handoff: unit_of_work_currency_wallet_repository_factory
unit_of_work_handoff_required: true
caller_owns_transaction_control: true
postgres_adapter_transaction_control_allowed: false
```

Rules:

- Application service 应通过 unit-of-work capability 或 implementation slice ratify 的等价 app-owned dependency 获取 `currency.Repository`。
- Wallet creation 与 transaction mutations 应在 caller-supplied unit of work 内执行。
- 只有 implementation slice 明确记录 composition 时，service 才可以在同一 unit of work 中组合 player account lookup 与 currency repository calls。
- PostgreSQL adapter 不得自己 start、commit 或 rollback transactions。
- Runtime behavior 不得通过 direct import PostgreSQL adapter 绕过 repository interface。

## 12. Test Expectations

Future runtime behavior implementation 应在实现前添加 focused fake-service tests：

```yaml
future_test_expectations:
  - rejects_missing_or_unvalidated_request_identity_before_unit_of_work
  - derives_wallet_owner_from_validated_player_identity
  - rejects_metadata_only_player_id_as_proof
  - ensures_player_wallet_through_repository_handoff
  - lists_balances_for_own_wallet_only
  - records_server_authoritative_grant_with_idempotency
  - records_player_spend_with_insufficient_balance_mapping
  - lists_transactions_for_own_wallet_only
  - maps_repository_conflicts_to_public_errors
  - redacts_wallet_owner_actor_idempotency_and_metadata_details
  - keeps_protocol_routes_and_generated_output_absent
```

本 gate 的默认 verification 仍是 static repository checks。Future implementation verification 应包含 `runtime/internal/app/currency` 下的 focused Go tests。

## 13. Stop Conditions

添加以下任何内容前必须停止并创建新的 bounded work item：

- runtime behavior implementation；
- runtime handlers 或 route dispatch wiring；
- startup composition；
- protocol routes；
- Protobuf source；
- generated output；
- repository interface changes；
- PostgreSQL adapter changes；
- migration changes；
- dependency additions；
- authentication/session behavior changes；
- reward integration；
- inventory integration；
- purchase behavior；
- catalog tables；
- event/audit tables；
- payment provider behavior；
- reservation、settlement、refund、transfer 或 paid-currency behavior；
- hosted deployment；
- SDK publication；
- release artifacts；
- public announcements 或 paid promotion；
- Pitaya-style distributed runtime；
- direct Nakama/Pitaya API compatibility。

## 14. Verification

本 gate 需要的 repository verification：

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.currency_wallet_runtime_behavior_gate
node tools/vibit check change define-currency-wallet-runtime-behavior-gate --json
node tools/vibit check module currency --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check protocol --json
node tools/vibit check schemas --json
node tools/vibit check memory --json
node tools/vibit check all --json
cd runtime && go test ./...
git diff --check
```

Accepted pre-existing warning 是 `runtime.identity_boundary`，位于 `runtime/internal/platform/persistence/postgres/authentication_repository.go`。
