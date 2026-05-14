# player Module Agent Guide 中文版

状态：Draft v0.1
说明：本文件是 `modules/player/AGENTS.md` 的简体中文译本。英文版本是权威版本。

## 何时使用本模块

当需求定义稳定 player identity 或 player account lifecycle 的 ownership 时，使用本模块。

本模块当前拥有 semantic 和 wire contract artifacts、已 ratify 的 PostgreSQL account lifecycle schema boundary、storage-neutral runtime repository interface，以及已实现的 PostgreSQL adapter boundary。它拥有以下概念的 vocabulary 和 ownership boundary：

- `player_id` 作为稳定的 domain identity。
- Player account lifecycle：creation、lookup、linking、disabling、deletion、recovery。
- 已 ratify 的 player account semantic contracts。
- 已 ratify 的 player account Protobuf wire messages。
- 已 ratify 的 `player_accounts` 和 `player_account_events` PostgreSQL lifecycle schema boundaries。
- 位于 `runtime/internal/modules/player/repository.go` 的 module-owned repository interface。
- 位于 `runtime/internal/platform/persistence/postgres/player_account_repository.go` 的 PostgreSQL adapter。

当前状态已经包含第一版 player account migration source、module-owned repository interface boundary，以及 focused PostgreSQL adapter implementation。它不实现 authentication、credentials、runtime sessions、WebSocket routes、WebSocket handshake behavior，也不实现 runtime account/session handlers。

## 何时不要使用本模块

不要把本模块用于：

- WebSocket connection mechanics。
- Protobuf envelope routing 或 payload conversion。
- Runtime session validation 或 session storage。
- Token formats、credential parsing、password storage、JWT、OAuth、OIDC、guest login、device login、social login 或 external identity providers。
- Inventory state、item quantities、item grants 或 inventory persistence。
- Currency、reward、quest、match 或 realtime room behavior。
- Authentication、credential、token、external identity 或 session tables。
- 添加超出已 ratify PostgreSQL adapter scope 的额外 player account repository adapter behavior。

如果需求涉及上述概念，应更新对应 owner 的 boundary，或先创建单独 ratified change，再添加代码。

## Extension Points

- Public commands 被 ratify 之后的 future command handlers。
- Public queries 被 ratify 之后的 future query handlers。
- Public events 被 ratify 之后的 future lifecycle events。
- Player identity 和 account lifecycle policies。
- Player account repository interface：`runtime/internal/modules/player/repository.go`。
- 已存在的第一版 player account migration source：`runtime/migrations/postgres/000002_create_player_account_state.sql`。
- Boundary and invariant tests。

`module.yaml` 中的 vocabulary placeholders 不是 public contracts。不要把它们当成实现 APIs 的许可。

在添加任何 player public contract 之前，创建或更新：

- `contracts/player/...`
- `.arch/contracts.yaml`
- `modules/player/module.yaml`
- `changes/` 下的 change spec
- 只有在 protocol impact 明确之后，才添加所需 Protobuf source

第一版已 ratify 的 player account Protobuf source 是 `proto/vibit/player/v1/player.proto`。生成的 Go Protobuf output 是 `runtime/internal/generated/proto/vibit/player/v1/player.pb.go`。不要把这些 wire shapes 当作添加 runtime handlers、WebSocket routes、authentication、token behavior、credential storage、player persistence implementation 或 session persistence 的许可。

已 ratify 的 player account PostgreSQL schema boundary 位于 `docs/postgresql-persistence-boundary.md` 和 `ADR-0022`。第一版允许的 player account lifecycle tables 是：

- `player_accounts`
- `player_account_events`

第一版 player account schema 不得包含 credentials、password hashes、authentication provider subjects、external identity links、access tokens、refresh tokens、runtime session rows、WebSocket connection state、request identity validation results、inventory state 或 permission grants。

已 ratify 的 runtime repository interface 是 `runtime/internal/modules/player/repository.go`。它可以定义 storage-neutral account lifecycle structs、`Repository.CreatePlayerAccount`、`Repository.GetPlayerAccount`，以及 persistence adapters 所需的 mutation/query shapes。它不得导入 PostgreSQL、WebSocket、Protobuf、authentication、token、credential、session、S3 或 MinIO dependencies。

PostgreSQL adapter 已实现于 `runtime/internal/platform/persistence/postgres/player_account_repository.go`，tests 位于 `runtime/internal/platform/persistence/postgres/player_account_repository_test.go`。它用 `NewPlayerAccountRepositoryForUnitOfWork(executor)` 构造，实现 `player.Repository`，使用 caller-supplied unit-of-work executor，并避免 `BEGIN`、`COMMIT` 和 `ROLLBACK`。它的第一版 SQL scope 仅限于插入 `player_accounts`，为 `PlayerAccountCreated` 插入 `player_account_events`，并从 `player_accounts` 读取当前 account lifecycle rows。它不授权 runtime player handlers、WebSocket routes、authentication、tokens、credentials、external identity links 或 session persistence。

## Forbidden Shortcuts

- 不要绕过 `module.yaml` 中声明的边界。
- 不要直接修改其他模块拥有的数据。
- 不要添加未登记的 public commands、queries、events 或 permissions。
- 不要让 transport connection metadata 充当 player identity。
- 不要把客户端提供的 `player_id` 或 `session_id` 当作 authenticated proof。
- 未经单独 decision，不要添加 authentication provider、token、credential、password、JWT、OAuth、OIDC、guest login、device login、social login 或 identity-provider code。
- 除非 future work item 明确 ratify schema change，否则不要再添加新的 player account migration source。
- 不要在 `runtime/internal/platform/persistence/postgres/` 之外添加 player account repository adapters。
- 未经单独 work item，不要扩展 player account PostgreSQL adapter 超出已 ratify 的 SQL scope。
- 不要让 player account PostgreSQL adapter 开启隐藏 transactions、解码 Protobuf、执行 WebSocket behavior、解析 credentials、validate tokens、bind sessions 或执行 permissions。
- 不要在 persistence work 中添加 runtime player account handlers 或 WebSocket routes。
- 不要让本模块拥有 inventory state，也不要依赖 inventory internals。
- 不要从本模块改变 Protobuf envelope 或 WebSocket handshake contract。
- 不要把 WebSocket、Protobuf、PostgreSQL、S3 或 MinIO dependencies 直接导入 domain module。
- 不要手工编辑 generated files。如果 generated output 有问题，应修改 source schema、template 或 generator。

## Required Tests

参见 `module.yaml` 中的 `tests.required`。

在当前 semantic、wire、migration、repository-interface 和 PostgreSQL-adapter-implemented 状态下，verification 是 repository-level architecture、module、contracts、protocol、generated-output、runtime、work queue、change-spec checks，以及 focused PostgreSQL adapter behavior tests。只有在添加 Go player handlers 或 runtime protocol adapters 时，player handler runtime behavior tests 才会成为必需。
