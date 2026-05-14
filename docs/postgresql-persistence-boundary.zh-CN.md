# PostgreSQL Persistence Boundary Standard 中文版

状态：Draft v0.2
最后更新：2026-05-14
范围：第一版 durable Go runtime 的 PostgreSQL repository、transaction、migration 和 event-recording boundaries  
说明：本文件是 `docs/postgresql-persistence-boundary.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本标准定义 vibit 在添加 PostgreSQL-backed module state 前必须遵守的 implementation boundary。

本标准应与 `.arch/runtime.yaml`、`.arch/dependencies.yaml`、`modules/inventory/module.yaml`、`modules/player/module.yaml`、`ADR-0011`、`ADR-0013`、`ADR-0014`、`ADR-0020`、`ADR-0021` 和 `ADR-0022` 一起使用。

## 1. 目的

Durable persistence 是长期 server project 最容易失去 architecture clarity 的位置。

主要风险不是写不出 PostgreSQL repository。主要风险是 agents 把 transactions、SQL ownership、migration semantics、permission checks、event recording 或 cross-module data access 隐藏在当时最方便的 package 里。

本标准在 persistent inventory behavior 实现前，以及后续 player account persistence adapters 添加前定义 ownership，防止这种漂移。

## 2. Layer Ownership

### Domain Module

Owner：

```text
runtime/internal/modules/<module>/
```

职责：

- 拥有 module repository interfaces。
- 拥有 module invariants 和 domain validation。
- 拥有 domain request、response、event、permission 和 error semantics。
- 只依赖 vibit-owned interfaces 和 standard-library packages。

不得：

- Import `github.com/jackc/pgx/v5`。
- Import `github.com/pressly/goose/v3`。
- Import S3 SDKs 或 MinIO clients。
- 构造 SQL strings。
- 直接管理 database transactions。
- 访问其他 module 的 tables。

### PostgreSQL Adapter

Owner：

```text
runtime/internal/platform/persistence/postgres/
```

职责：

- 拥有 `github.com/jackc/pgx/v5`。
- 使用 PostgreSQL 实现 module repository interfaces。
- 把 PostgreSQL rows 转换成 module-owned runtime structs。
- 在 adapter boundaries 把 PostgreSQL errors 转换为稳定的 module 或 application errors。
- 使用 application transaction boundary 提供的 transaction handles。
- 拥有 PostgreSQL runtime configuration parsing 和 pgx connection-pool construction。
- 拥有 vibit-owned unit-of-work boundary 的 pgx-backed transaction runner adapter。
- 保持 SQL behavior 显式，并由 focused tests 覆盖。

不得：

- 执行 module permissions。
- 用隐藏的 database-only logic 取代 module invariant checks。
- Import WebSocket 或 Protobuf platform packages。
- 向 domain modules 暴露 `pgx` types。

### Migration Tooling

Owner：

```text
runtime/internal/platform/migrations/
```

职责：

- 拥有 `github.com/pressly/goose/v3` invocation 和 validation。
- 提供 migration status、apply 和 validation helpers。Rollback helpers 在有变更需要前，仍属于未来的显式工作。
- 让 migration execution 保持 operationally explicit。

不得：

- 把 domain business behavior 隐藏在 Go migrations 中。
- 在没有 `runtime/migrations/postgres/` 下 source files 的情况下生成或应用 migrations。

第一版显式 migration runner 位于：

```text
runtime/internal/platform/migrations/postgres.go
```

它暴露 vibit-owned `PostgresRunner`，接收调用方提供的 `*sql.DB` 和 migration source filesystem 或 directory，配置 PostgreSQL dialect 的 goose，禁用全局 Go migration registry，列出已知 migration sources，报告结构化 migration status，并应用 pending 的 SQL-first migrations。它有意不负责构造 database connection，不关闭调用方拥有的 database handles，也不会在普通 `cmd/vibit-server` startup 中运行。

PostgreSQL platform owner 可以提供窄 helper，把 package-owned pool 适配为 migration tooling 需要的 `*sql.DB`，例如：

```text
runtime/internal/platform/persistence/postgres/sql_db.go
```

这些 helpers 必须把 pgx stdlib imports 留在 PostgreSQL platform owner package 内。它们不会把 migration execution ownership 从 `runtime/internal/platform/migrations/` 移走。

### Migration Sources

Owner：

```text
runtime/migrations/postgres/
```

职责：

- 保存 SQL-first PostgreSQL migration source files。
- 使用 deterministic sequence numbers。
- 同时包含 `-- +goose Up` 和 `-- +goose Down` sections。
- 显式声明 module-owned tables、indexes、constraints 和 foreign-key-like ownership references。

命名模式：

```text
runtime/migrations/postgres/000001_create_inventory_state.sql
```

只有在 change spec 解释 SQL 为什么不足时，才允许 Go migrations。它们不得包含 domain business logic。

### Transaction Boundary

Owner：

```text
runtime/internal/platform/tx/
```

Orchestrator：

```text
runtime/internal/app/
```

职责：

- State-changing commands 在 application-owned unit of work 中运行。
- Repository mutations 和 durable event recording 发生在同一个 unit of work 中。
- Query handlers 不修改 state，默认不需要 write transaction。
- Transaction retry policy 属于 platform，使用前必须显式定义。
- 在 event delivery 或 outbox standard 存在前，transaction 外的 event publication 继续 deferred。

这个 boundary 的第一版 Go skeleton 是：

```text
runtime/internal/platform/tx.Runner
runtime/internal/platform/tx.UnitOfWork
runtime/internal/app.TransactionalDispatcher
```

`TransactionalDispatcher` 为 command routes 开启 unit of work，并默认让 query routes 不经过 write unit of work。这个 skeleton 有意保持 driver-neutral；它不得向 application handlers 或 domain modules 暴露 PostgreSQL driver handles。

这个 boundary 的第一版 pgx-backed implementation 位于：

```text
runtime/internal/platform/persistence/postgres/runner.go
```

它实现 `runtime/internal/platform/tx.Runner`，在 PostgreSQL platform package 内部开启 `pgx` transaction；当 application callback 成功时 commit，当 callback 失败时 rollback，并在 commit failure 后尝试 rollback。传给 callback 的 unit of work 仍然满足 vibit-owned `tx.UnitOfWork` interface。PostgreSQL-specific repository factories，例如 `UnitOfWork.NewInventoryRepository`，必须保持为 package-owned helpers，不得迫使 application 或 domain packages import `pgx`。

不得：

- 让 repositories 在 command flows 中偷偷开启独立 write transactions。
- 让 domain modules 直接访问 transaction handles。
- 把 WebSocket connection 当成 transaction 或 session authority。

### PostgreSQL Runtime Configuration

Owner：

```text
runtime/internal/platform/persistence/postgres/
```

第一版 PostgreSQL configuration source 是：

```text
runtime/internal/platform/persistence/postgres/config.go
```

它通过这些 environment variables 读取显式 process input：

```text
VIBIT_POSTGRES_DSN
VIBIT_POSTGRES_MAX_CONNS
VIBIT_POSTGRES_MIN_CONNS
```

规则：

- 打开 PostgreSQL pool 时，`VIBIT_POSTGRES_DSN` 是必填项。
- Pool size settings 是可选项，且必须是非负整数。
- Configuration parsing 可以构建 `pgxpool.Config`，但普通 unit tests 不得要求 live PostgreSQL server。
- Connection strings 和 credentials 必须来自 environment 或显式 runtime input，不得存入 tracked files。
- Process startup 必须保持 PostgreSQL optional。`VIBIT_RUNTIME_STORE=postgres` 会显式选择 PostgreSQL-backed inventory composition path；默认仍然是 in-memory startup。
- Migration tooling 需要 `*sql.DB` adaptation 时，该 adaptation 必须留在 PostgreSQL platform owner package 内。

## 3. Inventory Persistence Boundary

第一版 durable module state 是 inventory。

Inventory module 拥有：

```text
inventory account row
inventory item quantity row
inventory item grant event/audit row
```

第一版 persistent schema 应建模：

- 每个 `player_id` 一个 inventory account row。
- 每个 `(player_id, item_id)` 一个 inventory item row。
- 每个 emitted `ItemGranted` event 一个 durable grant record。

Persistent grant flow 必须是 atomic：

```text
application command dispatch
-> open unit of work
-> lock inventory account for player_id
-> read current inventory
-> enforce permission and capacity policies
-> upsert item quantity
-> record ItemGranted grant row
-> commit unit of work
-> return application result
```

这个 lock 是必须的，因为 inventory capacity 依赖玩家当前 item 集合。第一版 PostgreSQL implementation 应使用显式 inventory account row lock，而不是隐式 best-effort read。

这个 flow 对应的 module-owned Go boundary 是：

```text
Repository.LockInventoryForMutation(ctx, player_id) -> MutationLock
MutationLock.GetInventory(ctx, player_id)
MutationLock.GrantItem(ctx, GrantItemMutation)
MutationLock.Release()
```

`MutationLock` 是单个 inventory aggregate 的 locked repository view。它不是 transaction owner。`Release` 只释放 aggregate lock 或 adapter-local resource；它不得 commit 或 roll back application-owned unit of work。

`GrantItemMutation` 必须携带 PostgreSQL adapters 记录 durable grant record 所需的 metadata：

```text
event_id
occurred_at
player_id
item_id
quantity
reason
```

Domain handler 必须在调用 `MutationLock.GrantItem` 前创建这些 metadata，这样 adapter 才能用同一个 application-owned executor 记录 item quantity change 和 `inventory_item_grants` row。

## 3.1 Player Account Persistence Schema Boundary

Player account persistence 被 ratify 为 account lifecycle state only。

Player module 拥有计划中的 PostgreSQL state：

```text
player account lifecycle row
player account lifecycle event/audit row
```

第一版 player account persistent schema 必须建模：

- 每个 stable `player_id` 一个 `player_accounts` row。
- 每个 durable player account lifecycle fact 一个 `player_account_events` row。
- `PlayerAccountCreated` 是第一种必须能记录的 lifecycle event type。

计划中的 `player_accounts` table 拥有这些 columns：

```text
player_id
display_name
account_state
created_at
updated_at
disabled_at
deleted_at
```

Column rules：

- `player_id` 是 stable primary key，必须是 non-blank text。
- `display_name` 必须是 non-blank text。Uniqueness 继续 deferred。
- `account_state` 必须被 constraint 限制为 `active`、`disabled` 或 `deleted`。
- `created_at` 和 `updated_at` 是 required timestamps。
- `disabled_at` 可为空，只能用于 disabled 或 deleted accounts。
- `deleted_at` 可为空，只能用于 deleted accounts。
- 第一版 migration 应使用显式 check constraints 表达 non-blank text 和 lifecycle state values。

计划中的 `player_account_events` table 拥有这些 columns：

```text
event_id
event_type
occurred_at
player_id
requested_by
account_state
display_name
metadata
recorded_at
```

Column rules：

- `event_id` 是 primary key，必须是 non-blank text。
- `event_type` 必须是 non-blank text。第一种 required value 是 `PlayerAccountCreated`。
- `occurred_at` 记录 domain event time。
- `player_id` 引用 `player_accounts(player_id)`，并且必须是 non-blank text。
- `requested_by` 必须是 non-blank text，但它本身不是 authenticated proof。
- `account_state` 记录 event 发生后的 lifecycle state。
- `display_name` 在相关事件中记录 event 发生后的 display name。
- `metadata` 可以使用 `JSONB NOT NULL DEFAULT '{}'::jsonb` 保存 structured non-secret event metadata。
- `recorded_at` 记录 PostgreSQL 存储该 row 的时间。

第一版 player account migration 应包含这些 indexes：

```text
player_accounts(account_state)
player_account_events(player_id, occurred_at)
player_account_events(event_type, occurred_at)
```

第一版 player account schema 不得存储：

- authentication credentials
- password hashes
- authentication provider names 或 provider subject IDs
- external identity links
- access tokens
- refresh tokens
- token signing metadata
- runtime session rows
- WebSocket connection state
- request identity validation results
- inventory state
- permission grants

这些 concerns 都需要在 implementation 前分别拥有 standards、contracts、migrations 和 verification。

当 runtime implementation 被单独 ratify 后，第一版 persistent player account creation flow 应该是 atomic：

```text
application command dispatch
-> open unit of work
-> enforce request validation and account creation permission
-> insert player_accounts row
-> insert player_account_events row for PlayerAccountCreated
-> commit unit of work
-> return application result
```

Player account repository boundary 由 module 拥有，并且 storage-neutral：

```text
Repository.CreatePlayerAccount(ctx, CreatePlayerAccountMutation)
Repository.GetPlayerAccount(ctx, player_id)
```

该 boundary 的第一版 source 是：

```text
runtime/internal/modules/player/repository.go
```

它定义 account lifecycle view、`Repository` interface，以及 durable event record 所需的 `CreatePlayerAccountMutation` fields。它必须继续不依赖 PostgreSQL、WebSocket、Protobuf、authentication、token、credential 或 session。

`CreatePlayerAccountMutation` 必须携带 PostgreSQL adapters 记录 durable event metadata 所需的信息：

```text
event_id
occurred_at
player_id
display_name
account_state
requested_by
```

PostgreSQL adapter 必须使用 application-owned unit of work 提供的 executor。它不得开启隐藏 write transactions，不得解析 credentials，不得 validate tokens，不得 bind sessions，不得解码 Protobuf payloads，也不得执行 WebSocket behavior。

第一版 player account migration source 是：

```text
runtime/migrations/postgres/000002_create_player_account_state.sql
```

它使用 existing migration files 之后的下一个 deterministic SQL migration number，并包含：

```text
-- +goose Up
-- Module: player
-- +goose Down
```

Runtime player account handlers、WebSocket route wiring、authentication、token behavior、credential storage、external identity linking 和 session persistence 继续 deferred，直到分别被 ratify。

## 3.2 Player Account PostgreSQL Adapter Boundary

第一版 player account PostgreSQL adapter 已经实现，但它仍然只是 persistence adapter。它不授权 runtime handlers、WebSocket routes、authentication、token behavior、credential storage、external identity linking 或 session persistence。

Adapter source path 是：

```text
runtime/internal/platform/persistence/postgres/player_account_repository.go
```

Focused test path 是：

```text
runtime/internal/platform/persistence/postgres/player_account_repository_test.go
```

Adapter 用以下方式构造：

```text
NewPlayerAccountRepositoryForUnitOfWork(executor)
```

该 constructor 返回以下 interface 的 implementation：

```text
player.Repository
```

Executor 必须由 application-owned unit of work 提供，通常来自 transaction-bound handle，例如 `pgx.Tx`。Adapter 为了 testability 使用与 inventory adapter 相同的小型 pgx-shaped executor interface。它不得调用 `BEGIN`、`COMMIT` 或 `ROLLBACK`；transaction lifetime 属于 application-owned unit of work。

Adapter 不得：

- 调用 `BEGIN`、`COMMIT` 或 `ROLLBACK`。
- 开启自己的隐藏 write transaction。
- 打开 PostgreSQL pool 或读取 PostgreSQL configuration。
- Apply migrations 或 inspect migration status。
- 解析 authentication credentials。
- Validate tokens。
- Bind 或 persist runtime sessions。
- 解码 Protobuf payloads。
- 知道 WebSocket handshake 或 connection behavior。
- 执行 permissions。
- Import transport、protocol、application bootstrap、authentication、credential、token、session、inventory、S3 或 MinIO packages。

第一版允许的 SQL operation scope 有意保持很窄：

- `CreatePlayerAccount` normalize `player.CreatePlayerAccountMutation`。
- `CreatePlayerAccount` insert 一条 `player_accounts` row。
- `CreatePlayerAccount` 在同一个 caller-supplied executor path 中，为 `PlayerAccountCreated` insert 一条 `player_account_events` row。
- `GetPlayerAccount` normalize `player_id`。
- `GetPlayerAccount` 从 `player_accounts` 读取当前 account lifecycle row。
- 两个 method 都不得读取或写入 credentials、tokens、external identity links、runtime sessions、WebSocket state、inventory state 或 permission grants。

Error mapping expectations：

- Missing rows 在 runtime handlers 把 error 暴露给 client 前，必须映射到稳定的 player account not-found error path。
- Duplicate `player_id` 或 duplicate `event_id` constraint violations 必须映射到稳定的 duplicate/conflict error paths。
- Check constraint violations 必须映射到稳定的 invariant 或 validation error paths。
- Unexpected PostgreSQL errors 可以带 adapter context wrap，但 `pgx` 或 `pgconn` types 不得泄露到 module-owned repository interface。

第一版 adapter implementation 包含以下 focused tests：

- Account creation 和 account lookup 的 fake-executor SQL shape tests。
- 证明 adapter 不发出 `BEGIN`、`COMMIT` 或 `ROLLBACK` 的 no-transaction-control test。
- Mutation normalization 和 UTC timestamp tests。
- Nullable lifecycle timestamps 的 row mapping tests。
- Duplicate account、duplicate event、check constraint 和 missing-row error mapping tests。
- 通过 `node tools/vibit check runtime` 执行 import-boundary tests。
- Default repository checks 不依赖 live PostgreSQL。
- Optional live integration coverage 只能通过 `VIBIT_POSTGRES_TEST_DSN` 启用。

PostgreSQL unit-of-work helper 通过以下入口从 transaction executor 构造该 repository：

```text
runtime/internal/platform/persistence/postgres.UnitOfWork.NewPlayerAccountRepository
```

这个 helper 必须留在 PostgreSQL platform package 中。它不得改变 module-owned `player.Repository` interface，不得迫使 application 或 domain packages import `pgx`，也不意味着 runtime player account handlers、WebSocket routes、authentication、credentials、tokens、external identity links 或 session persistence 已经被允许。

## 4. Repository Rules

Repository interfaces 属于 module。PostgreSQL adapters 实现这些 interfaces。

规则：

- Repository method 不得执行 permission checks。
- Repository method 不得解码 Protobuf payloads。
- Repository method 不得知道 WebSocket sessions。
- Repository method 必须保留那些也可表达为 database constraints 的 module-owned invariants。
- Persistent command flow 必须使用 transaction-bound repository。
- In-memory repository 可继续用于 tests 和 pre-persistence bootstrap，但它不是 authoritative durable store。
- Repositories 从 application composition 获得 transaction binding。它们不得为 command flow 创建隐藏的独立 write transaction。

对于 inventory，`GrantItem` 必须在 request validation 和 permission checks 之后、为 capacity enforcement 读取 current inventory 之前，获取 `LockInventoryForMutation`。Capacity-sensitive reads 和 grant mutation 必须通过返回的 `MutationLock` 执行。

PostgreSQL adapter 必须在 application-owned unit of work 内，用 `player_id` 上的显式 inventory account row lock 实现该 lock。除非有 superseding decision，否则不得用 advisory lock 或隐藏的 repository-owned transaction 替代。

第一版 adapter source 是：

```text
runtime/internal/platform/persistence/postgres/inventory_repository.go
```

构造方式是：

```text
NewInventoryRepositoryForUnitOfWork(executor)
```

Executor 必须由 application composition 提供，通常来自 transaction-bound handle，例如 `pgx.Tx`。Adapter 可以为了 testability 依赖一个小型 pgx-shaped executor interface，但它不得自行调用 `BEGIN`、`COMMIT` 或 `ROLLBACK`。

第一版 PostgreSQL unit-of-work helper 可以通过以下入口，从 transaction executor 构造该 repository：

```text
runtime/internal/platform/persistence/postgres.UnitOfWork.NewInventoryRepository
```

这个 helper 有意由 PostgreSQL platform package 拥有。它不改变 module-owned repository interface，也不向 inventory domain code 暴露 `pgx`。

第一版 persistent runtime composition 是显式启用的：

```text
VIBIT_RUNTIME_STORE=postgres
```

该路径会从 `VIBIT_POSTGRES_DSN` 打开 PostgreSQL pool，用 `runtime/internal/platform/persistence/postgres.Runner` 拥有 command unit-of-work，用 `postgres.UnitOfWork.NewInventoryRepository` 创建 command repository，并用 PostgreSQL inventory repository 服务 query routes。默认 server path 仍然是 `VIBIT_RUNTIME_STORE=memory`。

普通 server startup 不会自动 apply migrations。除非未来 change spec 授权 startup 自动迁移行为，migration execution 仍然必须是显式 operator 或 tooling action。

## 5. Migration Rules

SQL migrations 是 source artifacts。

规则：

- Migrations 必须作为 contract-bearing persistence changes 审查。
- Migration 一旦被视为已在 shared environment 应用，就不得继续编辑；应添加新的 migration。
- 每个创建 module-owned state 的 migration 都必须命名 owning module。
- 每个执行 invariant 的 migration 都应能映射回 module invariant 或 persistence boundary rule。
- Player account migrations 在 authentication、credential、token、external identity 或 session persistence standards 被分别 ratify 前，只能创建 3.1 节中 ratified 的 lifecycle tables。
- Destructive migrations 需要带 rollback 和 data compatibility notes 的 change spec。
- 在 persistent repository 被视为 verified 前，必须记录 migration validation。
- Migration execution helpers 必须由 operator、tool 或已批准的 runtime composition path 显式调用。除非后续 change spec 授权，普通 server startup 不得自动 apply migrations。

## 6. Test And Verification Rules

Persistence work 应在拥有行为的层添加 tests：

- Domain tests 不依赖 PostgreSQL，覆盖 invariants 和 permission behavior。
- PostgreSQL adapter tests 覆盖 row mapping、atomic mutations、constraint handling 和 concurrency-sensitive behavior。
- PostgreSQL config tests 覆盖 environment parsing 和 pool config construction，且不打开 live connection。
- PostgreSQL transaction runner tests 在 disposable PostgreSQL environment 尚未定义时，使用 fake pgx transactions 覆盖 begin、commit、rollback 和 dependency validation。
- Migration tests 或 checks 覆盖 SQL file naming、`goose` markers 和 apply/rollback validation。
- Runtime wiring tests 覆盖 configuration 和 composition，不覆盖 SQL behavior。

在 disposable PostgreSQL test environment 存在前，PostgreSQL integration tests 可以由显式 environment variable 控制。Agents 必须记录这些 tests 是否被 skip 以及原因。

当前 PostgreSQL adapter 已经有 focused fake-executor tests，覆盖 SQL shape 和 transaction-bound behavior。当前 runtime wiring tests 覆盖 store selection 和 application composition，默认不会打开 live PostgreSQL connection。Live repository integration tests 通过 disposable PostgreSQL verification environment standard 选择性启用。

当前 PostgreSQL configuration 和 transaction runner 的 focused tests 位于：

```text
runtime/internal/platform/persistence/postgres/config_test.go
runtime/internal/platform/persistence/postgres/runner_test.go
runtime/internal/app/bootstrap/inventory_test.go
runtime/cmd/vibit-server/main_test.go
```

当前 PostgreSQL migration runner 的 focused tests 位于：

```text
runtime/internal/platform/migrations/postgres_test.go
```

这些 tests 在不要求 live PostgreSQL server 的情况下验证 option handling、SQL source discovery 和 error propagation。

当前 repository verification 仍是：

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check migrations
node tools/vibit check runtime
node tools/vibit check all
```

当前 migration source check 是：

```bash
node tools/vibit check migrations
```

它验证 SQL migration naming、`goose` Up/Down markers、没有未批准的 Go migrations、owning-module traces、architecture manifest references，以及第一版 inventory table references。它还不会针对 PostgreSQL 执行 apply 或 rollback。

当前 migration apply/status API 是：

```text
runtime/internal/platform/migrations.NewPostgresRunner
runtime/internal/platform/migrations.PostgresRunner.Status
runtime/internal/platform/migrations.PostgresRunner.Apply
```

Live migration apply/status verification 继续 deferred，直到 disposable PostgreSQL verification environment 被定义。

未来 persistence verification 应在该 disposable environment 上增加 PostgreSQL-backed migration status、apply、rollback 和 repository integration checks。

## 7. Agent Rules

Agents 必须：

- 在添加 PostgreSQL repositories、migrations、transaction code 或 persistent runtime wiring 前阅读本标准。
- 修改 durable data ownership 前更新 module manifests。
- 把 `pgx` 保持在 `runtime/internal/platform/persistence/postgres/` 内。
- 把 `goose` 保持在 `runtime/internal/platform/migrations/` 内。
- 把 SQL source files 保持在 `runtime/migrations/postgres/` 下。
- 显式记录 skipped PostgreSQL integration verification。

Agents 不得：

- 向 domain modules 添加 PostgreSQL imports。
- 向 command handlers 添加 migration side effects。
- 未经 change spec 授权，向普通 server startup 添加 migration side effects。
- 把 database constraints 当成唯一可见的 business rules 来源。
- 在没有 concrete large-object use case 和 dependency adoption record 的情况下，向 inventory persistence 添加 object storage。
- 让 MinIO 成为 durable inventory runtime 的 mandatory dependency。
