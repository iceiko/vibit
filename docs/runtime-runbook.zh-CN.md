# Runtime Runbook 中文版

状态：Draft v0.1
最后更新：2026-05-13
范围：第一版 Go runtime process startup 和 manual verification
说明：本文件是 `docs/runtime-runbook.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本 runbook 记录如何启动第一版 vibit Go runtime process。

## Current Runtime Surface

当前 runtime process 挂载一个 gameplay WebSocket endpoint：

```text
/v1/ws
```

该 endpoint 期望接收 binary WebSocket messages，其中 payload 是 `vibit.protocol.v1.Envelope` Protobuf bytes。

第一版已挂载的 request loop 是：

```text
WebSocket binary frame
-> Protobuf envelope
-> application dispatch
-> inventory command or query handler
-> Protobuf response envelope
-> WebSocket binary frame
```

## Start The Server

从 Go runtime module 启动：

```bash
cd runtime
go run ./cmd/vibit-server
```

默认 listen address 是：

```text
:8080
```

可以这样覆盖：

```bash
VIBIT_ADDR=:9090 go run ./cmd/vibit-server
```

默认 runtime store 是 in memory：

```text
VIBIT_RUNTIME_STORE=memory
```

如果要启动显式 PostgreSQL-backed inventory composition path，需要同时提供 store selector 和 PostgreSQL DSN：

```bash
VIBIT_RUNTIME_STORE=postgres VIBIT_POSTGRES_DSN='postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable' go run ./cmd/vibit-server
```

可选 PostgreSQL pool settings：

```text
VIBIT_POSTGRES_MAX_CONNS
VIBIT_POSTGRES_MIN_CONNS
```

普通 server startup 不会自动 apply migrations。在 fresh database 上使用 PostgreSQL store path 前，必须显式 apply 或 verify migrations。

## Manual Verification Path

1. 启动 server。
2. 连接 WebSocket client 到 `ws://127.0.0.1:8080/v1/ws`。
3. 发送 `inventory.GrantItem` 或 `inventory.GetInventory` 的 binary Protobuf `Envelope`。
4. 确认 response 是 binary Protobuf `Envelope`，并且带有相同的 `request_id`。

Text WebSocket messages 会被 transport adapter 拒绝。该 endpoint 不接受 JSON。

## Current Runtime Assumptions

- Runtime 默认使用 in-memory inventory repository。
- `VIBIT_RUNTIME_STORE=postgres` 会启用显式 PostgreSQL inventory composition path。
- Inventory bootstrap permissions 允许 grant 和 read operations。
- Authentication 和 session validation 尚未实现。
- PostgreSQL persistence 只在显式选择时接入 inventory runtime composition。Persistence boundary 已定义在 `docs/postgresql-persistence-boundary.md`。
- 普通 server startup 不会自动 apply PostgreSQL migrations。
- Optional live PostgreSQL verification 定义在 `docs/postgresql-verification-environment.md`；它要求 `VIBIT_POSTGRES_TEST_DSN`，且不属于默认 server startup。
- Generated route registration 尚未实现。

这些是第一版 request loop 的 bootstrap assumptions，不是长期 production policy。

## Verification Commands

除特别说明外，从仓库根目录运行：

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime
node tools/vibit check postgres-env
node tools/vibit check all
```

`node tools/vibit check postgres-env` 是静态标准检查。它不会连接 PostgreSQL。Live PostgreSQL verification 继续通过 `VIBIT_POSTGRES_TEST_DSN` 选择性启用。

使用下面的 command 针对 disposable PostgreSQL database 运行当前 live durable inventory verification：

```bash
cd runtime && VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

该 test 会显式 apply inventory migration，并通过 PostgreSQL-backed runtime composition 验证 WebSocket Protobuf `GrantItem` 后接 `GetInventory` 的 request loop。如果未设置 `VIBIT_POSTGRES_TEST_DSN`，test 会 skip，并记录 live PostgreSQL verification 不可用。

该 test 默认使用 `drop_schema` cleanup semantics。其他 cleanup modes 会被该 test 有意 skip，因为 migration apply 必须从 clean schema 验证。
