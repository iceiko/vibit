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
- `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
- `../decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
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

Transport handlers 只负责把 requests 适配成 commands 或 queries。它们不得隐藏 business logic。

State-changing commands 应通过 `internal/app/` 进入，并在 application-owned unit of work 中运行。Repository mutations 和 domain event recording 应发生在同一个 unit of work 内。

Query handlers 不应改变状态，默认不需要 write transaction。

在 vibit 采纳明确的 event delivery 或 outbox standard 前，transaction 外的 event publication 继续 deferred。

## 6. Generated Files

Generated files 对 non-system agents 不可变。

如果 generated output 错了，应修改 source contract、schema、template 或 generator。除非 change spec 或 Agent Decision Record 明确授予 `generated_file_override`，不要手工编辑 generated files。

`internal/generated/proto/` 下的 Go Protobuf generated output 必须通过已接受的 Buf 和 `protoc-gen-go` 路径从 `../proto/` sources 生成。该目录下的文件必须是带有 `protoc-gen-go` marker 和 source trace 的 generated `*.pb.go` files，或者是在 generation 尚未运行时使用的临时 `.gitkeep` placeholders。

不要把 handwritten runtime code 放在 `internal/generated/proto/` 或 `internal/generated/contracts/` 下。

## 7. 当前状态

这个 runtime workspace 现在已经有第一批 generated Protobuf output 和第一段窄 runtime handoff slice。

但它仍然没有实现 WebSocket transport、PostgreSQL persistence、migrations 或完整的 application dispatcher。

## 8. 验证

从仓库根目录运行 repository verification：

```bash
node tools/vibit check runtime
node tools/vibit check generated
node tools/vibit check all
```

当 Go source files 存在且本地 Go toolchain 可用时，runtime verification 应包括：

```bash
go test ./...
go vet ./...
```

当 Go toolchain 不可用或测试没有运行时，不要声称已经完成 Go test verification。
