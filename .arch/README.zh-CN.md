# Architecture Manifests 中文版

状态：Draft v0.1  
最后更新：2026-05-13
范围：vibit 的机器可读架构入口  
说明：本文件是 `.arch/README.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

这个目录包含面向 agents、humans、generators 和未来 verification commands 的 architecture manifests。

这些 manifests 不是装饰性文档。它们的目标是成为可执行的架构上下文。

## 目的

`.arch/` 目录应该回答 agent 在修改代码前必须解决的问题：

- 存在哪些 modules？
- 每个 module 拥有什么？
- 哪些 dependencies 被允许？
- 哪些 contracts 定义 public behavior？
- 哪些 events、commands、queries、errors 和 permissions 已登记？
- 哪些文件是 generated？
- 哪些 tests 或 checks 证明 architecture rules？

## 当前文件

```text
.arch/
  README.md
  README.zh-CN.md
  modules.yaml
  conventions.yaml
  runtime.yaml
  contracts.yaml
  dependencies.yaml
```

这是第一版草案。这些文件在实现代码存在前先描述预期形态。

`runtime.yaml` 记录第一版 Go server runtime 方向的 runtime readiness decisions。它指向约束第一语言、服务器实例模型、contract boundary、client protocol、wire format、persistence direction、dependency adoption 和 proof slice 的 Agent Decision Records。

`ADR-0014` 记录第一版 Go runtime package layout 和 boundary rules。计划中的 Go module root 是 `runtime/`，process startup 放在 `runtime/cmd/vibit-server/`，application orchestration 放在 `runtime/internal/app/`，platform adapters 放在 `runtime/internal/platform/`，手写 domain runtime logic 放在 `runtime/internal/modules/<module>/`，生成的 Go outputs 放在 `runtime/internal/generated/`，SQL-first PostgreSQL migrations 放在 `runtime/migrations/postgres/`，Protobuf source files 放在仓库根目录 `proto/`。

`contracts.yaml` 登记 public command、query、event、error 和 permission contract source files。Contract files 位于 `contracts/` 下，是 semantic source artifacts，不是 generated output。Protobuf wire schemas 计划位于 `proto/` 下，并且必须与这些 semantic contracts 对齐。

`dependencies.yaml` 记录 foundational dependency decision slots。它标识哪些 dependency categories 在 implementation import 或 require 具体 packages 前需要 adoption records。

第一批已接受的 Go runtime dependencies 由 `ADR-0013` 记录：

- `github.com/coder/websocket` 用于 platform WebSocket transport adapter。
- `google.golang.org/protobuf`、`protoc-gen-go` 和 Buf CLI 用于 Protobuf runtime、generation、linting、breaking checks、formatting 和 orchestration。
- `github.com/jackc/pgx/v5` 用于 PostgreSQL platform persistence adapters。
- `github.com/pressly/goose/v3` 用于 SQL-first migration tooling。

S3 client tooling、MinIO deployment、observability 和外部 Go test framework adoption 仍然 deferred，直到具体 runtime needs 证明它们必要。

Go runtime implementation code 尚未开始。当它开始时，agents 应先更新相关 manifests，再实现，并把第三方 transport、protocol、persistence 和 migration dependencies 保持在它们声明过的 owner packages 中。

## 未来预期文件

```text
.arch/test-matrix.yaml
.arch/generation.yaml
```

当第一版 prototype 需要它们时，项目应逐步添加这些文件。

## Agent 规则

修改实现代码前，agents 应：

1. 阅读 `CONSTITUTION.md`。
2. 阅读 `AGENTS.md`。
3. 阅读 `.arch/modules.yaml`。
4. 阅读 `.arch/conventions.yaml`。
5. 在修改或创建 runtime implementation code 前，阅读 `.arch/runtime.yaml`。
6. 在新增或修改 public contracts 前，阅读 `.arch/contracts.yaml`。
7. 在添加 foundational dependencies 前，阅读 `.arch/dependencies.yaml`。
8. 在相关 module 存在时，阅读其 `module.yaml`。
9. 当 public architecture 变化时，先更新 manifests，再实现。

如果 manifest 缺少安全变更所需的信息，应更新 manifest 或记录这个缺口。

Decision authority boundary 以 `ADR-0012` 为准。在 maintainer 授权后，agents 可以在已确认方向内按专业评估决定 technical sub-decisions，但修改 product direction、constitutional principles、runtime language、primary protocol direction、persistence direction、major architecture patterns、module ownership、breaking contracts、validation 或 permission 强度，以及接受 licensing-risk、hosting、cost、operations 或 vendor-lock-in commitments 前，仍必须询问。

## 验证方向

这些 manifests 最终应支持类似命令：

```bash
vibit check architecture
vibit check module <module>
vibit check contracts
vibit check change <change-id>
```

在这些命令存在前，agents 必须记录 architecture verification 当前不可用。
