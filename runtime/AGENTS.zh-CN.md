# Go Runtime Agent 指南

状态：Draft v0.1
最后更新：2026-05-13
范围：`runtime/` Go server runtime workspace
权威来源：`../CONSTITUTION.md`、`../AGENTS.md` 和 `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
说明：本文件是 `runtime/AGENTS.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本指南适用于第一版 Go server runtime implementation。

## 1. 目的

`runtime/` 是 vibit 第一版 Go server runtime 的 Go module。

Go module path 是：

```text
github.com/iceiko/vibit/runtime
```

Runtime 的目标是通过一个小而长期维护的 backend slice 证明 vibit 的核心命题：

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

不要把这个 workspace 当成一次性 demo。

## 2. 必读内容

修改 `runtime/` 下文件前，先阅读：

- `../CONSTITUTION.md`
- `../AGENTS.md`
- `../.arch/runtime.yaml`
- `../.arch/dependencies.yaml`
- `../.arch/contracts.yaml`
- `../docs/generated-output.md`
- `../docs/runtime-protocol-adapter.md`
- persistence work 前阅读 `../docs/postgresql-persistence-boundary.md`
- live PostgreSQL verification work 前阅读 `../docs/postgresql-verification-environment.md`
- authentication、token、credential、external identity、session persistence、request identity trust、WebSocket handshake、player handler 或 player route work 前阅读 `../docs/authentication-token-session-validation.md`
- `../docs/runtime-runbook.md`
- `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
- `../decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
- persistence work 前阅读 `../decisions/ADR-0020-postgresql-persistence-boundary.md`
- player account persistence work 前阅读 `../decisions/ADR-0022-player-account-postgresql-schema-boundary.md`
- authentication/session design 或 implementation work 前阅读 `../decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`
- 受影响的 module manifest，例如 `../modules/inventory/module.yaml`
- `../changes/` 下相关 change spec

## 3. Package Ownership

使用以下 package boundaries：

- `cmd/vibit-server/`：process startup、configuration wiring 和 lifecycle。
- `internal/app/`：command/query dispatch、application composition 和 transaction orchestration。
- `internal/platform/transport/ws/`：WebSocket transport adapter 和 `github.com/coder/websocket` ownership。
- `internal/platform/protocol/protobuf/`：Protobuf framing、envelope conversion 和 wire message adaptation。
- `internal/platform/persistence/postgres/`：PostgreSQL adapter implementation 和 `github.com/jackc/pgx/v5` ownership。
- `internal/platform/migrations/`：migration tooling invocation 和 validation。
- `internal/platform/events/`：event recording 和 publication mechanisms。
- `internal/platform/tx/`：transaction boundary 和 unit-of-work interfaces。
- `internal/modules/<module>/`：手写 domain module runtime logic。
- `internal/generated/contracts/`：生成的 Go contract shapes。
- `internal/generated/proto/`：生成的 Go Protobuf files。
- `migrations/postgres/`：SQL-first PostgreSQL migration sources。

## 4. 依赖规则

Domain modules 不得直接 import 第三方 transport、protocol、persistence、migration、object-storage 或 framework dependencies。

允许的 owner packages：

- `github.com/coder/websocket`：仅 `internal/platform/transport/ws/`。
- `google.golang.org/protobuf`：仅 generated protocol packages 和 protocol adapter packages。
- `github.com/jackc/pgx/v5`：仅 `internal/platform/persistence/postgres/`。
- `github.com/pressly/goose/v3`：仅 `internal/platform/migrations/`。

未经检查 `../.arch/dependencies.yaml` 并创建所需 adoption record，不要添加新的 foundational dependencies。

## 5. Runtime Boundary Rules

Runtime protocol handoff 必须遵循 `../docs/runtime-protocol-adapter.md`。

WebSocket transport 读写 frames。Protobuf protocol adaptation 解码和编码 envelopes。Application dispatch 路由 commands 和 queries。Domain modules 执行 invariants。Generated packages 只提供 shapes。

WebSocket transport handlers 把 opaque frame bytes 交给注入的 protocol/application composition。它们不直接把 requests 适配成 commands 或 queries，也不得隐藏 business logic。

State-changing commands 应通过 `internal/app/` 进入，并在 application-owned unit of work 中运行。Repository mutations 和 domain event recording 应发生在同一个 unit of work 内。

当前 transaction skeleton 是 `internal/platform/tx.Runner`、`internal/platform/tx.UnitOfWork` 和 `internal/app.TransactionalDispatcher`。Application code 可以 import 这个 transaction boundary package，但不得 import persistence、migration、protocol 或 transport platform adapters。Query routes 默认应不经过 write unit of work。

Query handlers 不应改变状态，默认不需要 write transaction。

在 vibit 采纳明确的 event delivery 或 outbox standard 前，transaction 外的 event publication 继续 deferred。

PostgreSQL persistence work 必须遵循 `../docs/postgresql-persistence-boundary.md`。Repository interfaces 保持 module-owned，`pgx` 保持在 `internal/platform/persistence/postgres/` 下，`goose` 保持在 `internal/platform/migrations/` 下，SQL migration sources 保持在 `migrations/postgres/` 下。

第一版 durable inventory implementation 中，`GrantItem` 必须使用 transaction-bound repository，并在读取当前 items、执行 capacity-sensitive mutation 前调用 `LockInventoryForMutation`。返回的 `MutationLock` 是 locked aggregate view，不是 transaction owner。Repositories 不得在 command flows 中偷偷开启独立 write transactions。

第一版 PostgreSQL inventory repository adapter 是 `internal/platform/persistence/postgres/inventory_repository.go`。使用 `NewInventoryRepositoryForUnitOfWork` 构造它，并传入由 application-owned unit of work 提供的 executor，例如 `pgx.Tx` 或兼容的 test executor。该 adapter 不得调用 `BEGIN`、`COMMIT` 或 `ROLLBACK`；transaction lifetime 属于 `internal/platform/tx` 和 `internal/app`。

PostgreSQL configuration 由 `internal/platform/persistence/postgres/config.go` 拥有。它读取 `VIBIT_POSTGRES_DSN`、`VIBIT_POSTGRES_MAX_CONNS` 和 `VIBIT_POSTGRES_MIN_CONNS`，构建 pgx pool configuration，并且普通 unit tests 不得要求 live PostgreSQL server。Connection strings 和 credentials 必须来自 environment 或显式 runtime input，不得 commit。

pgx-backed transaction runner 是 `internal/platform/persistence/postgres/runner.go`。它实现 `internal/platform/tx.Runner`，同时把 pgx transaction handles 保持在 PostgreSQL platform package 内部。它会 commit 成功的 command unit of work，rollback 失败的 callback unit of work，并提供 package-owned helpers，例如 `UnitOfWork.NewInventoryRepository`，供未来 persistent composition 使用。不要从 `internal/app/` 或 domain modules import PostgreSQL runner；persistent runtime wiring 必须发生在已批准的 composition boundary 中。

`GrantItemMutation` 携带 `event_id`、`occurred_at` 和 `reason`，这样 PostgreSQL adapter 可以在与 item quantity update 相同的 executor path 中持久化 `inventory_item_grants`。

第一版 inventory migration source 是 `migrations/postgres/000001_create_inventory_state.sql`。它创建 `inventory_accounts`、`inventory_items` 和 `inventory_item_grants`。当 migration sources 或 migration guidance 发生变化时，运行 `node ../tools/vibit check migrations`。设置 `VIBIT_POSTGRES_TEST_DSN` 后，opt-in live durable inventory request-loop verification 会覆盖 migration status 和 apply behavior。

Ratified player account PostgreSQL schema boundary 已记录在 `../docs/postgresql-persistence-boundary.md` 和 `../decisions/ADR-0022-player-account-postgresql-schema-boundary.md`。第一版 player account migration source 是 `migrations/postgres/000002_create_player_account_state.sql`。该 migration 只创建 `player_accounts` 和 `player_account_events` lifecycle state。它不得添加 credentials、password hashes、external identity links、access tokens、refresh tokens、runtime session rows、WebSocket connection state、request identity validation results、inventory state 或 permission grants。

Player account repository interface boundary 是 `internal/modules/player/repository.go`。它是 storage-neutral domain code，可以定义 account lifecycle structs、`Repository.CreatePlayerAccount`、`Repository.GetPlayerAccount`，以及 persistence adapters 所需的 durable mutation metadata。PostgreSQL adapter 是 `internal/platform/persistence/postgres/player_account_repository.go`，focused tests 位于 `internal/platform/persistence/postgres/player_account_repository_test.go`。它使用 `NewPlayerAccountRepositoryForUnitOfWork(executor)`，实现 `player.Repository`，从 application-owned unit of work 接收 executor，并且不得调用 `BEGIN`、`COMMIT` 或 `ROLLBACK`。`UnitOfWork.NewPlayerAccountRepository` 是 PostgreSQL package helper，不得向 application 或 domain packages 暴露 pgx。

Player account PostgreSQL adapter 不授权 runtime handlers、WebSocket routes、authentication、token behavior、credential storage、external identity linking 或 session persistence。除非后续 change ratify 更多行为，adapter 只能写入 `player_accounts`，为 `PlayerAccountCreated` 写入 `player_account_events`，并从 `player_accounts` 读取当前 lifecycle rows。

Authentication、token 和 session validation design boundary 记录在 `../docs/authentication-token-session-validation.md` 和 `../decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`。它分离 authentication proof、login methods、tokens、credentials、external identity links、runtime sessions、request identity、WebSocket handshake authentication、player account lifecycle、transport connection metadata 和 Protobuf envelope metadata。当前 `MetadataOnlySessionValidator` 是 non-authenticated bootstrap path。不要把 metadata-only `player_id`、`session_id` 或 `connection_id` 当作 production proof；未经单独 ratify，不要添加 authentication runtime code、token parsing、credential lookup、external identity linking、session persistence、Protobuf envelope authentication changes、WebSocket handshake authentication、runtime player handlers 或 WebSocket routes。`runtime.authentication_token_session_boundary` 是该边界的 repository check rule。

第一版显式 PostgreSQL migration runner 是 `internal/platform/migrations/postgres.go`。它拥有 `github.com/pressly/goose/v3`，接收调用方提供的 `*sql.DB` 和 migration source filesystem 或 directory，列出 SQL migration sources，报告结构化 status，并且只在被显式调用时应用 pending migrations。未经 change spec 授权，不要把它接入普通 `cmd/vibit-server` startup。

Live PostgreSQL verification 受 `../docs/postgresql-verification-environment.md` 约束。它通过 `VIBIT_POSTGRES_TEST_DSN` 选择性启用；普通 unit tests、`node ../tools/vibit check runtime` 和默认 repository checks 不得要求运行中的 PostgreSQL server。当 live PostgreSQL check 因为没有 disposable DSN 而跳过时，必须显式记录。

## 6. Generated Files

Generated files 对 non-system agents 不可变。

如果 generated output 错了，应修改 source contract、schema、template 或 generator。除非 change spec 或 Agent Decision Record 明确授予 `generated_file_override`，不要手工编辑 generated files。

`internal/generated/proto/` 下的 Go Protobuf generated output 必须通过已接受的 Buf 和 `protoc-gen-go` 路径从 `../proto/` sources 生成。该目录下的文件必须是带有 `protoc-gen-go` marker 和 source trace 的 generated `*.pb.go` files，或者是在 generation 尚未运行时使用的临时 `.gitkeep` placeholders。

不要把 handwritten runtime code 放在 `internal/generated/proto/` 或 `internal/generated/contracts/` 下。

## 7. 当前状态

这个 runtime workspace 现在已经有第一批 generated Protobuf output、第一段窄 runtime handoff slice、第一版 WebSocket transport adapter、一个用于 command 和 query routes 的小型 application dispatch skeleton、第一版 transaction boundary skeleton、带 command-safe mutation lock 的第一版 inventory repository/policy/handler runtime boundary、第一版 PostgreSQL configuration parser、第一版 pgx-backed transaction runner adapter、第一版 PostgreSQL inventory repository adapter、第一条 inventory Protobuf/domain payload bridge、第一条 application-error-to-Protobuf-error-envelope mapper、第一版 frame-to-Protobuf-to-application composition adapter、用于 Protobuf command/query tests 的 package-local request-loop test fixture、挂载 `/v1/ws` 的 minimal process wiring、显式 PostgreSQL inventory runtime composition path，以及 opt-in live PostgreSQL durable inventory request-loop verification test。

这个 workspace 已经有 documented PostgreSQL persistence boundary、transaction skeleton、PostgreSQL configuration parser、pgx-backed transaction runner、第一版 inventory migration source、第一版显式 migration apply/status runner、第一版 PostgreSQL repository adapter、显式 runtime store selection、已添加第一版 migration source、storage-neutral repository interface、focused PostgreSQL adapter implementation 和 PostgreSQL unit-of-work factory helper 的 ratified player account PostgreSQL lifecycle schema boundary，以及已经 ratify 但尚未实现 runtime authentication 的 authentication/token/session validation design boundary，并有只有设置 `VIBIT_POSTGRES_TEST_DSN` 才会运行 live branch 的 live verification test。`VIBIT_RUNTIME_STORE=memory` 仍然是默认值。提供 `VIBIT_POSTGRES_DSN` 时，`VIBIT_RUNTIME_STORE=postgres` 会启用 PostgreSQL-backed inventory composition。这个 workspace 仍然没有实现 generated route registration、generated protocol bridge creation、production authentication/session validation、runtime player account handlers、WebSocket player routes、automatic startup migrations 或 catalog-driven error retryability。

第一版手动 process run path 是：

```bash
cd runtime
go run ./cmd/vibit-server
```

第一版显式 persistent process run path 是：

```bash
cd runtime
VIBIT_RUNTIME_STORE=postgres VIBIT_POSTGRES_DSN='postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable' go run ./cmd/vibit-server
```

普通 server startup 不会自动 apply migrations。

第一条 opt-in live durable inventory verification command 是：

```bash
cd runtime
VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

如果未设置 `VIBIT_POSTGRES_TEST_DSN`，该 test 会 skip，并记录 live PostgreSQL verification 不可用。第一版 live execution 已经在 local Termux PostgreSQL 18.2 上通过。

## 8. 验证

从仓库根目录运行 repository verification：

```bash
node tools/vibit check runtime
node tools/vibit check generated
node tools/vibit check migrations
node tools/vibit check postgres-env
node tools/vibit check all
```

当 Go source files 存在且本地 Go toolchain 可用时，runtime verification 应包括：

```bash
go test ./...
go vet ./...
```

当 Go toolchain 不可用或测试没有运行时，不要声称已经完成 Go test verification。
