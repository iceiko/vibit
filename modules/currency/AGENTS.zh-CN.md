# currency 模块 Agent 指南

状态：Draft v0.1

## 什么时候使用本模块

本模块用于 currency wallet repository vocabulary、storage-neutral value types、application-owned wallet behavior，以及当前已批准的 protected currency wallet protocol route family。

当前已实现的 slice 刻意保持狭窄：

- `runtime/internal/modules/currency.Repository`
- `CurrencyWallet`、`CurrencyWalletOwner`、`CurrencyWalletBalance`、`CurrencyWalletTransaction`、actor、lifecycle、idempotency、amount 和 version value types
- create、wallet lookup、owner lookup、balance list、grant record、spend record 和 transaction list input/result types
- conflict classes 和 redacted repository errors
- normalization helpers 和 focused Go tests

`M-224 Currency Wallet Repository Interface Implementation` 已由 `W-0296` 完成。检查规则是 `runtime.currency_wallet_repository_interface_implementation`。

`W-0297 Define currency wallet PostgreSQL adapter gate` 已完成。它接受 `ADR-0205`，注册 `runtime.currency_wallet_postgresql_adapter_gate`。

`W-0298 Implement currency wallet PostgreSQL adapter` 已完成。它接受 `ADR-0206`，注册 `runtime.currency_wallet_postgresql_adapter_implementation`，添加 `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`、focused fake-executor tests 和 `UnitOfWork.NewCurrencyWalletRepository`。

`W-0299 Define currency wallet runtime behavior gate` 已完成。它接受 `ADR-0207`，注册 `runtime.currency_wallet_runtime_behavior_gate`，并定义 `docs/currency-wallet-runtime-behavior-gate.md`。

`W-0300 Implement currency wallet runtime behavior` 已完成。它接受 `ADR-0208`，注册 `runtime.currency_wallet_runtime_behavior_implementation`，并添加 `runtime/internal/app/currency/service.go` 和聚焦 service tests。

`W-0301 Define currency wallet protocol route gate` 已完成。它接受 `ADR-0209`，注册 `runtime.currency_wallet_protocol_route_gate`，并定义 `docs/currency-wallet-protocol-route-gate.md`。

`W-0302 Implement currency wallet protocol route` 已完成。它接受 `ADR-0210`，注册 `runtime.currency_wallet_protocol_route_implementation`，添加 `proto/vibit/currency/v1/currency.proto`、generated Go output、route keys、bootstrap handlers、protocol bridge mapping、payload registry dispatch、server-side grant policy enforcement、focused tests 和 authenticated local route proof。

No repository next work item is currently ready。不要发明 `W-0303`；reward integration、inventory integration、purchase behavior、catalog/event tables、payment behavior、migrations、authentication/session behavior changes、SDK publication、hosted deployments、distributed runtime 和 direct Nakama/Pitaya API compatibility 都必须留在后续 bounded work items 后面。

## 什么时候不要使用本模块

不要用本模块处理：

- WebSocket 或 HTTP transport behavior。
- 已批准的 `vibit.currency.v1` route artifacts 之外的 Protobuf 或 generated wire behavior。
- 本模块下的 PostgreSQL adapter implementation 或 SQL execution。
- `runtime/internal/app/currency` 之外的 runtime wallet behavior。
- reward、inventory、purchase 或 payment behavior。
- player account lifecycle。
- authentication、token formats、credential storage 或 session validation。
- currency catalog management、event/audit tables、reservations、settlement、refunds、transfers、paid currency、matchmaking 或 match runtime。
- direct Nakama 或 Pitaya public API compatibility。

如果需求涉及这些概念，应创建或更新对应 owner boundary，而不是在本模块隐藏 ownership。

## 扩展点

- Repository interface：`runtime/internal/modules/currency.Repository`
- Repository value types：`CurrencyWallet`、`CurrencyWalletBalance`、`CurrencyWalletTransaction`、`CurrencyWalletOwner`、`CurrencyWalletActor`、`CurrencyWalletVersion`、`CurrencyBalanceVersion`
- Lifecycle vocabulary：`active`、`suspended`、`closed`
- Transaction vocabulary：`grant`、`spend`
- Actor vocabulary：`system`、`player`、`operation`
- Normalizers：wallet records、balance records、transaction records、list results、owner identity、actor identity、idempotency fields、metadata JSON 和 repository inputs
- Tests：`runtime/internal/modules/currency/repository_test.go`
- PostgreSQL adapter owner candidate：`runtime/internal/platform/persistence/postgres`
- PostgreSQL adapter implementation：`runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
- Runtime behavior gate：`runtime.currency_wallet_runtime_behavior_gate`
- Runtime behavior implementation：`runtime.currency_wallet_runtime_behavior_implementation`
- Runtime behavior implementation source：`runtime/internal/app/currency/service.go`
- Protocol route gate：`runtime.currency_wallet_protocol_route_gate`
- Protocol route implementation：`runtime.currency_wallet_protocol_route_implementation`
- Protocol route source：`proto/vibit/currency/v1/currency.proto`
- Protocol route bridge：`runtime/internal/platform/protocol/protobuf/currency_bridge.go`
- Protocol route handler：`runtime/internal/app/bootstrap/currency.go`

未来 runtime behavior 必须先从 validated request identity 或 service-authoritative context 派生 owner 和 actor identity，再调用 repository interface；client-supplied player ids 不是 authentication proof。

## 禁止的捷径

- 不要绕过 `module.yaml` 声明的边界。
- 不要添加未注册的 public commands、queries、events、errors 或 permissions。
- 不要在本模块下添加 PostgreSQL adapter code。
- 不要在本模块导入 `pgx`、`database/sql`、WebSocket packages、generated Protobuf packages、SDK packages 或 distributed runtime packages。
- 不要在 currency module source 中执行 SQL 或写入具体 SQL statements。
- 不要从本模块修改 migrations。
- 不要从本模块接线 transport behavior。
- 不要添加已批准的 `W-0302` currency wallet route family 之外的新 protocol routes、Protobuf sources 或 generated output。
- 不要在 currency value types 中存放 raw credentials、raw tokens、verifier material、lookup digests、verifier digests、cookies、headers、transport subprotocols、connection metadata、payment provider payloads、Nakama API paths、Pitaya server ids 或 direct compatibility fields。
- 不要把 wallet ids、owner ids、actor ids、`player_id`、`session_id`、idempotency keys 或 transport metadata 当成 authenticated proof。

## 必需测试

见 `module.yaml` 中的 `tests.required`。

当前 repository interface slice 的测试必须覆盖：

- Repository interface storage neutrality。
- owner kind、lifecycle state、transaction kind、actor kind 和 conflict vocabulary 的 closed set。
- Wallet record normalization。
- Balance record normalization。
- Transaction record normalization、idempotency、actor、metadata 和 delta validation。
- Create/get/owner/list input normalization。
- Grant/spend input normalization 和 expected-version pointer copying。
- List result copying。
- Redacted conflict and repository errors。
- 不存在 secret、transport、protocol、PostgreSQL、session、distributed、payment 和 direct compatibility fields。

修改 currency runtime source 后运行 `node tools/vibit check runtime`。Go 可用时也运行 `cd runtime && go test ./...`。
