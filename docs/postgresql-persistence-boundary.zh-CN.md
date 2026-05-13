# PostgreSQL Persistence Boundary Standard 中文版

状态：Draft v0.1  
最后更新：2026-05-13  
范围：第一版 durable Go runtime 的 PostgreSQL repository、transaction、migration 和 event-recording boundaries  
说明：本文件是 `docs/postgresql-persistence-boundary.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本标准定义 vibit 在添加 PostgreSQL-backed module state 前必须遵守的 implementation boundary。

本标准应与 `.arch/runtime.yaml`、`.arch/dependencies.yaml`、`modules/inventory/module.yaml`、`ADR-0011`、`ADR-0013`、`ADR-0014` 和 `ADR-0020` 一起使用。

## 1. 目的

Durable persistence 是长期 server project 最容易失去 architecture clarity 的位置。

主要风险不是写不出 PostgreSQL repository。主要风险是 agents 把 transactions、SQL ownership、migration semantics、permission checks、event recording 或 cross-module data access 隐藏在当时最方便的 package 里。

本标准在 persistent inventory behavior 实现前定义 ownership，防止这种漂移。

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
- 在 tooling 实现后提供 migration status、apply、rollback 和 validation helpers。
- 让 migration execution 保持 operationally explicit。

不得：

- 把 domain business behavior 隐藏在 Go migrations 中。
- 在没有 `runtime/migrations/postgres/` 下 source files 的情况下生成或应用 migrations。

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

不得：

- 让 repositories 在 command flows 中偷偷开启独立 write transactions。
- 让 domain modules 直接访问 transaction handles。
- 把 WebSocket connection 当成 transaction 或 session authority。

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

## 4. Repository Rules

Repository interfaces 属于 module。PostgreSQL adapters 实现这些 interfaces。

规则：

- Repository method 不得执行 permission checks。
- Repository method 不得解码 Protobuf payloads。
- Repository method 不得知道 WebSocket sessions。
- Repository method 必须保留那些也可表达为 database constraints 的 module-owned invariants。
- Persistent command flow 必须使用 transaction-bound repository。
- In-memory repository 可继续用于 tests 和 pre-persistence bootstrap，但它不是 authoritative durable store。

对于 inventory，在 `GrantItem` 变成 PostgreSQL-backed 前，persistent repository boundary 必须支持 command-safe mutation lock。

## 5. Migration Rules

SQL migrations 是 source artifacts。

规则：

- Migrations 必须作为 contract-bearing persistence changes 审查。
- Migration 一旦被视为已在 shared environment 应用，就不得继续编辑；应添加新的 migration。
- 每个创建 module-owned state 的 migration 都必须命名 owning module。
- 每个执行 invariant 的 migration 都应能映射回 module invariant 或 persistence boundary rule。
- Destructive migrations 需要带 rollback 和 data compatibility notes 的 change spec。
- 在 persistent repository 被视为 verified 前，必须记录 migration validation。

## 6. Test And Verification Rules

Persistence work 应在拥有行为的层添加 tests：

- Domain tests 不依赖 PostgreSQL，覆盖 invariants 和 permission behavior。
- PostgreSQL adapter tests 覆盖 row mapping、atomic mutations、constraint handling 和 concurrency-sensitive behavior。
- Migration tests 或 checks 覆盖 SQL file naming、`goose` markers 和 apply/rollback validation。
- Runtime wiring tests 覆盖 configuration 和 composition，不覆盖 SQL behavior。

在 disposable PostgreSQL test environment 存在前，PostgreSQL integration tests 可以由显式 environment variable 控制。Agents 必须记录这些 tests 是否被 skip 以及原因。

当前 repository verification 仍是：

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime
node tools/vibit check all
```

未来 persistence verification 应增加：

```bash
goose status
goose up
goose down
PostgreSQL repository integration tests
```

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
- 把 database constraints 当成唯一可见的 business rules 来源。
- 在没有 concrete large-object use case 和 dependency adoption record 的情况下，向 inventory persistence 添加 object storage。
- 让 MinIO 成为 durable inventory runtime 的 mandatory dependency。
