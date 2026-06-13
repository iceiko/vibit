# 货币钱包协议路由门禁

状态：Accepted v0.1
最后更新：2026-06-08
范围：在 application runtime behavior 之后，为未来 client-facing currency wallet protocol routes 定义 gate-only boundary
依赖：`docs/currency-wallet-runtime-behavior-gate.md`、`decisions/ADR-0208-currency-wallet-runtime-behavior-implementation.md`、`docs/runtime-protocol-adapter.md`、`docs/game-protocol.md`、`docs/generated-output.md`、`docs/bound-identity-route-policy-gate.md`、`docs/runtime-session-validation-gate.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`
规范决策：`ADR-0209`

配对英文原文是 `docs/currency-wallet-protocol-route-gate.md`。英文文件为权威版本。

本文定义 currency wallet protocol route gate。它只是 gate artifact。本项不添加 protocol route implementation、Protobuf source、generated output、startup wiring、runtime handlers、repository interface changes、PostgreSQL adapter changes、migration changes、dependencies、authentication/session behavior changes、reward integration、inventory integration、purchase behavior、catalog tables、event/audit tables、payment behavior、reservation behavior、settlement behavior、refund behavior、transfer behavior、operations/admin behavior、SDK publication、generated client libraries、hosted deployments、release artifacts、public announcements、paid promotion、Pitaya-style distributed architecture，或 direct Nakama/Pitaya API compatibility。

## 1. 核心规则

currency wallet protocol route gate 记录为：

```yaml
currency_wallet_protocol_route_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0301
decision: ADR-0209
check_rule: runtime.currency_wallet_protocol_route_gate
source_runtime_behavior_implementation_decision: ADR-0208
source_runtime_behavior_implementation: runtime/internal/app/currency/service.go
source_runtime_behavior_tests: runtime/internal/app/currency/service_test.go
source_runtime_behavior_gate_decision: ADR-0207
source_repository_interface_decision: ADR-0204
repository_interface: runtime/internal/modules/currency.Repository
future_protocol_source_candidate: proto/vibit/currency/v1/currency.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/currency_bridge.go
future_protocol_bridge_test_candidate: runtime/internal/platform/protocol/protobuf/currency_bridge_test.go
future_application_handler_candidate: runtime/internal/app/bootstrap/currency.go
future_application_handler_test_candidate: runtime/internal/app/bootstrap/currency_test.go
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed_as_proof: false
client_supplied_wallet_id_allowed_for_lookup: false
client_supplied_actor_id_allowed_as_proof: false
server_authoritative_grant_policy_required: true
first_owner_kind: player
first_actor_kinds:
  - player
  - system
first_payload_package: vibit.currency.v1
protobuf_envelope_change_status: unchanged
websocket_transport_credential_neutral: true
protocol_route_gate_only: true
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
runtime_handler_added: false
startup_wiring_added: false
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
payment_behavior_added: false
reservation_behavior_added: false
settlement_behavior_added: false
refund_behavior_added: false
transfer_behavior_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_protocol_route_implementation_work_item: W-0302
future_protocol_route_implementation_direction: implement_currency_wallet_protocol_route
```

## 2. 目的

`W-0300` 已在 `runtime/internal/app/currency` 下实现 application-owned currency wallet behavior。下一步不应直接添加 route code 或 `.proto` generation，而应先定义 protocol route gate，记录未来 WebSocket/Protobuf 暴露如何调用该 service，同时不把 wallet behavior 移到 transport、generated files、persistence adapters、reward systems、purchase systems 或 payment integrations 中。

Nakama 给出产品面压力：durable wallets、balances、grants、spends 和 transaction history 是常见 game backend economy capability。vibit 应覆盖这个能力类别。

Pitaya 给出架构面压力：acceptors、sessions、route handlers、serializers 和 backend services 应保持分离。vibit 通过 credential-neutral WebSocket transport、显式 Protobuf payload bridge，以及调用 application-owned currency wallet services 的 application-owned route handlers 来适配这一点。

本 gate 在实现前记录：

- candidate route names；
- candidate request/response message shapes；
- route protection 与 identity handoff posture；
- protocol adapter、application handler 与 startup ownership；
- generated-output expectations；
- public error mapping 与 redaction expectations；
- local proof expectations；
- Nakama/Pitaya reference mapping；
- 将 implementation 与 generated artifacts 留到后续 slice 的 stop conditions。

## 3. 未来路由面

第一组 route family 应只暴露 own-player wallet operations：

```yaml
candidate_routes:
  - kind: command
    module: currency
    name: EnsurePlayerWallet
    route_id: currency.EnsurePlayerWallet
    service_method: EnsurePlayerWallet
  - kind: query
    module: currency
    name: GetOwnWallet
    route_id: currency.GetOwnWallet
    service_method: GetOwnWallet
  - kind: query
    module: currency
    name: ListOwnWalletBalances
    route_id: currency.ListOwnWalletBalances
    service_method: ListOwnWalletBalances
  - kind: command
    module: currency
    name: GrantCurrency
    route_id: currency.GrantCurrency
    service_method: GrantCurrency
  - kind: command
    module: currency
    name: SpendCurrency
    route_id: currency.SpendCurrency
    service_method: SpendCurrency
  - kind: query
    module: currency
    name: ListOwnWalletTransactions
    route_id: currency.ListOwnWalletTransactions
    service_method: ListOwnWalletTransactions
```

规则：

- Route names 必须保持 vibit-native，不复制 Nakama route paths 或 Pitaya route naming。
- `EnsurePlayerWallet`、`GrantCurrency`、`SpendCurrency` 是 commands。
- `GetOwnWallet`、`ListOwnWalletBalances`、`ListOwnWalletTransactions` 是 queries。
- 第一组 routes 仅面向 validated player wallet owner，不暴露 arbitrary owner ids 或 arbitrary wallet id lookup。
- `GrantCurrency` 在 service behavior 中仍是 server-authoritative。未来 route implementation 必须明确它是 local-proof/dev-only、admin-only、通过 application composition system-initiated，还是不注册为 public client route。它绝不能变成 unauthenticated 或 client-authoritative mint route。
- Client payloads 可在 service 已有 vocabulary 的范围内提供 currency code、positive amount、idempotency key/scope、reason code、external reference、metadata JSON、expected wallet version、expected balance version、list limit、currency pagination cursor 和 transaction pagination cursor。
- Client payloads 不得提供 owner ids、用于 lookup 的 wallet ids、作为 proof 的 actor ids、raw access tokens、credential material、lookup digests、verifier digests、SQL details、payment provider payloads、带 purchase semantics 的 catalog ids、带 reward execution semantics 的 reward ids、inventory mutation fields 或 direct external API compatibility markers。
- Reward execution、inventory integration、purchase behavior、currency catalogs、event/audit streams、payment providers、reserves、settlement、refunds、transfers、operations/admin inspection、script hooks、SDK/client libraries、hosted deployments、distributed runtime routing 和 direct compatibility 均保持 deferred。
- 未来 route implementation 必须显式注册 routes。不得添加 catch-all currency route 或 reflective handler。

## 4. 协议形状

第一候选 currency wallet protocol source 是：

```text
proto/vibit/currency/v1/currency.proto
```

第一候选 generated output 是：

```text
runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go
```

第一候选 Protobuf package 是：

```text
vibit.currency.v1
```

Candidate messages：

```yaml
messages:
  CurrencyWallet:
    fields:
      wallet_id: string
      owner_kind: string
      lifecycle_state: string
      wallet_version: int64
      created_at: string
      updated_at: string
      state_changed_at: string
  CurrencyWalletBalance:
    fields:
      currency_code: string
      balance_amount: int64
      balance_version: int64
      created_at: string
      updated_at: string
  CurrencyWalletTransaction:
    fields:
      transaction_id: string
      currency_code: string
      transaction_kind: string
      amount_delta: int64
      balance_after: int64
      actor_kind: string
      reason_code: string
      external_reference: string
      metadata_json: string
      created_at: string
  EnsurePlayerWalletRequest:
    fields: {}
  EnsurePlayerWalletResponse:
    fields:
      wallet: CurrencyWallet
      status: string
  GetOwnWalletRequest:
    fields: {}
  GetOwnWalletResponse:
    fields:
      wallet: CurrencyWallet
      status: string
  ListOwnWalletBalancesRequest:
    fields:
      limit: int32
      after_currency_code: string
  ListOwnWalletBalancesResponse:
    fields:
      balances: repeated CurrencyWalletBalance
      next_currency_code: string
      status: string
  GrantCurrencyRequest:
    fields:
      currency_code: string
      amount: int64
      idempotency_key: string
      idempotency_scope: string
      reason_code: string
      external_reference: string
      metadata_json: string
      expected_wallet_version: int64
      expected_balance_version: int64
  GrantCurrencyResponse:
    fields:
      transaction: CurrencyWalletTransaction
      status: string
  SpendCurrencyRequest:
    fields:
      currency_code: string
      amount: int64
      idempotency_key: string
      idempotency_scope: string
      reason_code: string
      external_reference: string
      metadata_json: string
      expected_wallet_version: int64
      expected_balance_version: int64
  SpendCurrencyResponse:
    fields:
      transaction: CurrencyWalletTransaction
      status: string
  ListOwnWalletTransactionsRequest:
    fields:
      currency_code: string
      limit: int32
      after_transaction_id: string
      after_transaction_time: string
  ListOwnWalletTransactionsResponse:
    fields:
      transactions: repeated CurrencyWalletTransaction
      next_transaction_id: string
      next_transaction_time: string
      status: string
```

规则：

- 除非后续 protocol ADR 明确改变 envelope semantics，现有 `proto/vibit/protocol/v1/envelope.proto` 必须保持 unchanged。
- 暴露时间值时应使用 RFC3339Nano UTC text。
- Optional expected-version mapping 必须保留 service 的 optional expected-version vocabulary。未来实现测试必须明确 field absence 与 `0` semantics。
- Grant 与 spend request 的 amount input 是 positive integer minor units。Transaction deltas 仍由 service/repository 所有，不得允许 client 提交 negative spend deltas。
- `metadata_json` 不是 log-safe。它不得出现在默认 error messages、logs、route policy diagnostics 或 test names 中。
- Wallet ids、owner ids、transaction ids、actor ids、idempotency keys、reason codes、external references、balance amounts、metadata JSON 和 payment-adjacent fields 默认都不是 log-safe。
- Protocol shape 不得包含 `owner_id`、raw `wallet_id` lookup、client-supplied actor id、raw access tokens、credential material、lookup digests、verifier digests、SQL details、payment provider payloads、catalog purchase fields、reward execution fields、inventory mutation fields、private transport metadata 或 direct external API compatibility markers。

## 5. 路由保护与身份交接

第一 route-policy posture 是：

```yaml
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed_as_proof: false
client_supplied_wallet_id_allowed_for_lookup: false
client_supplied_actor_id_allowed_as_proof: false
server_authoritative_grant_policy_required: true
```

规则：

- 未来 currency wallet routes 必须是 protected gameplay routes。
- 未来 handlers 必须从既有 protected-route flow 接收 validated `app.RequestIdentity`。
- Envelope/session metadata 中 metadata-only `player_id` 或 `session_id` 永远不能成为 wallet owner proof。
- Client payloads 不得选择 owner ids、用于 lookup 的 wallet ids 或作为 proof 的 actor ids。
- Service 仍负责在 unit-of-work access 或 repository access 前拒绝 invalid identity。
- WebSocket transport 保持 credential-neutral。
- 本 gate 不改变 authentication、token validation、session persistence、first-message binding、WebSocket handshake behavior、bound-identity policy、route-protection semantics 或 permission schemas。

## 6. 未来路由流

未来实现必须保留如下顺序：

```yaml
future_route_flow:
  - receive WebSocket/Protobuf envelope through existing request path
  - apply protected-route authenticated wrapper policy
  - obtain validated authenticated request identity
  - decode vibit.currency.v1 request payload
  - reject payload owner proof, wallet lookup proof, and actor proof
  - map payload fields to runtime/internal/app/currency service request
  - call application-owned currency wallet service
  - map service result to vibit.currency.v1 response payload
  - map service public errors to protocol error responses
  - keep transport, Protobuf bridge, application handler, service, repository, and PostgreSQL adapter ownership separated
```

规则：

- WebSocket transport 保持 credential-neutral 和 currency-neutral。
- Protobuf bridge 只应映射 payload 与 response shape；不得拥有 wallet behavior、minting policy、spend policy、permission decisions、repository calls 或 payment logic。
- Application bootstrap handlers 应拥有 route registration 和 service invocation。
- Application service 继续拥有 identity checks、validation handoff、repository conflict mapping、actor derivation、idempotency handoff 和 public service errors。
- PostgreSQL adapters 保持 persistence-only。

## 7. 公有错误映射

未来 route implementation 应映射 service public errors，且不泄漏内部细节：

```yaml
public_error_mapping:
  CURRENCY_WALLET_INVALID_REQUEST: invalid_request
  CURRENCY_WALLET_UNAUTHENTICATED: unauthenticated
  CURRENCY_WALLET_NOT_FOUND: not_found
  CURRENCY_WALLET_ALREADY_EXISTS: already_exists
  CURRENCY_WALLET_NOT_ACTIVE: wallet_not_active
  CURRENCY_WALLET_INSUFFICIENT_BALANCE: insufficient_balance
  CURRENCY_WALLET_DUPLICATE_TRANSACTION: duplicate_transaction
  CURRENCY_WALLET_VERSION_MISMATCH: version_mismatch
  CURRENCY_WALLET_UNAVAILABLE: unavailable
```

规则：

- Public protocol errors 只能暴露 stable public codes/classes，以及在后续 ADR 授权时暴露 retryability posture。
- Not-found、owner mismatch、inactive-wallet、insufficient-balance、duplicate/idempotency cases 不得泄漏 cross-player wallet existence 或 private transaction details。
- Internal repository errors、SQL details、wallet ids、owner ids、actor ids、transaction ids、idempotency keys、reason codes、external references、metadata JSON、access-token material、credential material、lookup digests、verifier digests、payment provider payloads 和 transport metadata 必须保持在默认 logs 与 error messages 之外。
- Authentication 与 route protection failures 必须使用既有 protected-route semantics。本 gate 不发明新的 proof carrier。

## 8. Generated Output Posture

未来 generated output 必须遵循 `docs/generated-output.md`。

规则：

- `proto/vibit/currency/v1/currency.proto` 只能由后续 implementation slice 添加。
- `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go` 只能作为 Buf/protoc 生成输出添加。
- Generated Go output 必须包含 `protoc-gen-go` generated-code marker，并可追溯到 source `.proto`。
- Agents 不得手工编辑 generated Go Protobuf files。
- 本 gate 不改变 `buf.yaml`、`buf.gen.yaml` 或 generated output。

## 9. Ownership

未来实现必须保留这些 owner：

```yaml
currency_service_owner: runtime/internal/app/currency
application_handler_owner: runtime/internal/app/bootstrap
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
currency_repository_interface_owner: runtime/internal/modules/currency
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
websocket_transport_owner: runtime/internal/platform/transport/ws
startup_owner: runtime/cmd/vibit-server
generated_output_owner: runtime/internal/generated/proto/vibit/currency/v1
protobuf_source_owner: proto/vibit/currency/v1
```

规则：

- Currency wallet runtime behavior 保持在 `runtime/internal/app/currency`。
- Protocol bridge code 只能映射 payload fields。
- Persistence code 仍只负责 currency adapter behavior。
- Generated output 必须从 `.proto` sources 生成，且不得手工编辑。
- Startup wiring、route registration 与 generated output 保持在后续 implementation work item 后面。

## 10. Nakama Reference Mapping

Nakama reference mapping：

```yaml
nakama_reference_mapping:
  capability_family: economy_wallets_and_ledgers
  mapped_capabilities:
    - wallet_create_or_get
    - wallet_balance_read
    - wallet_balance_list
    - currency_grant
    - currency_spend
    - wallet_transaction_history
  direct_api_compatibility: false
```

Nakama 只指导有用的 capability class。vibit 不复制 Nakama route paths、API names、runtime script APIs、storage model names、permission semantics、economy ledger semantics 或 public API compatibility。

## 11. Pitaya Reference Mapping

Pitaya reference mapping：

```yaml
pitaya_reference_mapping:
  architecture_pressure:
    - acceptor_session_handler_separation
    - serializer_adapter_separation
    - backend_service_boundary
  distributed_architecture_status: deferred
  direct_api_compatibility: false
```

Pitaya 只指导 layering pressure。本 gate 不添加 Pitaya-style distributed topology、frontend/backend split、RPC、groups、service discovery、distributed wallet routing 或 direct Pitaya API compatibility。

## 12. 未来必需测试

未来 implementation tests 应覆盖：

- all selected route ids 的 route registration；
- command/query kind mapping；
- Protobuf request/response bridge mapping；
- optional expected-version mapping；
- 从 validated request identity 派生 owner；
- 通过既有 protected-route wrapper 拒绝 metadata-only identity；
- 拒绝 client-supplied owner id、wallet id lookup proof 与 actor proof；
- grant route policy 不能变成 unauthenticated 或 client-authoritative minting；
- spend route 把 positive amount input 映射到 service-owned spend behavior；
- balance 与 transaction pagination mapping；
- metadata JSON redaction；
- public `CURRENCY_WALLET_*` error 到 protocol error 的映射；
- private wallet、balance、transaction、idempotency、token、credential、SQL、payment 与 transport details 的 redaction；
- WebSocket transport 或 PostgreSQL adapter packages 中没有 route behavior；
- 如果添加 Protobuf source，则验证 generated-output traceability；
- 后续 proof slice 授权时，对 ensure、get、list balances、spend 和 transaction listing 的 local proof expectations。

## 13. 停止条件

添加以下任一内容前必须停止并创建独立 work item：

- protocol route implementation；
- `proto/vibit/currency/v1/currency.proto`；
- generated Go Protobuf output；
- protocol bridge implementation；
- application bootstrap handlers；
- startup route registration；
- new dependencies；
- migration changes；
- repository interface changes；
- PostgreSQL adapter changes；
- authentication/session behavior changes；
- route-protection semantic changes；
- reward integration；
- inventory integration；
- purchase behavior；
- currency catalog tables；
- event/audit tables；
- payment behavior；
- reservation behavior；
- settlement behavior；
- refund behavior；
- transfer behavior；
- operations/admin behavior；
- SDK publication；
- generated client libraries；
- hosted deployments；
- release artifacts；
- public announcements；
- paid promotion；
- Pitaya-style distributed architecture；
- direct Nakama/Pitaya API compatibility。

## 14. Verification

本 gate 的 repository check rule 是：

```text
runtime.currency_wallet_protocol_route_gate
```

本 gate 后建议 verification：

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.currency_wallet_protocol_route_gate
node tools/vibit check change define-currency-wallet-protocol-route-gate --json
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

本 gate 不添加 Go runtime behavior，因此不要求新增 Go tests；但在关闭 development turn 前运行完整 runtime test 仍然有价值。
