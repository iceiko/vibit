# vibit

vibit 是一个开源 agent-native server framework，用于构建 AI coding agents 能从第一性原理理解、扩展、验证和维护的后端系统。

状态：宪法设计阶段

## 本项目所说的 Agent-Native 是什么

Agent-native 主要不是指服务器带有 AI 功能。

它指的是服务器架构被设计成让 AI coding agents 能可靠地在其中工作：

- 架构是显式的，而不是团队默契。
- Module ownership 是声明出来的，而不是靠猜。
- Public behavior 先有 contract，再有实现。
- 重复结构由 generator 生成。
- Business rules 作为 invariants 被测试。
- Cross-module communication 有明确边界。
- Change workflow 有文档且可验证。
- 文档同时服务 humans 和 agents。

NPC agents、memory、model routing、tool calling 和 simulations 等 AI gameplay features 未来可以成为扩展，但它们不是项目基础。

## 为什么需要这个项目

许多现有服务器代码库，是为带有本地上下文、长期记忆和隐性团队约定的人类维护者构建的。AI coding agents 在这些代码库里可以提供帮助，但当架构规则隐藏、模块边界薄弱、测试不完整或公共契约不清晰时，它们往往难以发挥。

vibit 从另一个前提出发：

> 下一代长期演进的服务器软件，应该从设计之初就让 agents 能安全地理解、修改、验证和扩展。

目标不是让 agents 魔法般变聪明，而是让代码库更清晰、更有边界、更多由 contract 驱动、更可生成、更可测试。

## 当前文档

- `CONSTITUTION.md`：权威项目宪法
- `CONSTITUTION.zh-CN.md`：简体中文译本
- `AGENTS.md`：仓库级 coding agent 操作指南
- `AGENTS.zh-CN.md`：简体中文译本
- `.arch/README.md`：机器可读 architecture manifest 入口
- `.arch/modules.yaml`：第一版 module registry manifest 草案
- `.arch/conventions.yaml`：第一版 repository convention manifest 草案
- `.arch/protocol.yaml`：第一版 WebSocket Protobuf envelope 的 game protocol framework manifest
- `.arch/runtime.yaml`：第一版 Go server runtime 方向的 runtime readiness manifest
- `.arch/contracts.yaml`：public command、query、event、error 和 permission source files 的 contract registry
- `.arch/dependencies.yaml`：foundational dependency decision slots 的 dependency adoption registry
- `.arch/reference.yaml`：game server capability planning 的 Nakama/Pitaya 主动参考基线 manifest
- `buf.yaml`：Protobuf 的 Buf source、lint 和 breaking-check configuration
- `buf.gen.yaml`：计划中的 Go Protobuf output 的 Buf generation configuration
- `proto/`：protocol envelope 和 module wire schemas 的 Protobuf source root
- `docs/module-manifest.md`：module manifest 标准
- `docs/module-manifest.zh-CN.md`：简体中文译本
- `docs/change-spec.md`：change spec 标准
- `docs/change-spec.zh-CN.md`：简体中文译本
- `changes/_template/`：可复用 change spec 模板
- `docs/conversation-log.md`：conversation log 标准
- `docs/conversation-log.zh-CN.md`：简体中文译本
- `conversations/`：maintainer-agent project memory
- `docs/agent-decision-record.md`：Agent Decision Record 标准
- `docs/agent-decision-record.zh-CN.md`：简体中文译本
- `decisions/`：持久 decision rationale
- `docs/schema-validation.md`：schema validation 标准
- `docs/schema-validation.zh-CN.md`：简体中文译本
- `docs/dependency-adoption.md`：dependency adoption 标准
- `docs/dependency-adoption.zh-CN.md`：简体中文译本
- `docs/game-protocol.md`：game protocol framework 标准
- `docs/game-protocol.zh-CN.md`：简体中文译本
- `docs/generated-output.md`：generated output 标准
- `docs/generated-output.zh-CN.md`：简体中文译本
- `docs/runtime-protocol-adapter.md`：runtime protocol adapter boundary 标准
- `docs/runtime-protocol-adapter.zh-CN.md`：简体中文译本
- `docs/runtime-runbook.md`：第一版 Go runtime process startup 和 manual verification runbook
- `docs/runtime-runbook.zh-CN.md`：简体中文译本
- `docs/reference-game-server-alignment.md`：Nakama 和 Pitaya 的 active game server reference alignment 标准
- `docs/reference-game-server-alignment.zh-CN.md`：简体中文译本
- `docs/nakama-pitaya-product-parity-roadmap.md`：面向 Nakama/Pitaya 同级常用能力覆盖的产品路线图
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`：简体中文译本
- `schema/`：用于机器可检查 standards 的 JSON Schema files
- `rules/`：面向机器可读 check metadata 的 rule catalogs

英文文档是权威版本。简体中文译本服务于人类阅读和早期项目讨论。

## 预期方向

vibit 应逐步演进出：

- `.arch/` 下的 architecture manifests
- 由 `.arch/runtime.yaml` 和 Agent Decision Records 约束的第一版 Go server runtime
- WebSocket 作为第一版 gameplay/client protocol
- Protobuf 作为第一版 client/server wire message format
- PostgreSQL 作为第一版 authoritative durable relational store
- `github.com/coder/websocket` 作为第一版 WebSocket platform adapter dependency
- `google.golang.org/protobuf`、`protoc-gen-go` 和 Buf CLI 作为第一版 Protobuf tooling stack
- `github.com/jackc/pgx/v5` 作为 platform persistence adapters 后面的第一版 PostgreSQL driver
- `github.com/pressly/goose/v3` 作为第一版 SQL-first migration tooling
- 第一版 Go module 位于 `runtime/go.mod`，module path 为 `github.com/iceiko/vibit/runtime`
- Go runtime package boundaries 位于 `runtime/cmd/vibit-server/`、`runtime/internal/app/`、`runtime/internal/platform/`、`runtime/internal/modules/` 和 `runtime/internal/generated/`
- Protobuf source files 位于 `proto/vibit/<module>/v1/`，生成的 Go Protobuf output 位于 `runtime/internal/generated/proto/`
- Generated output rules 位于 `docs/generated-output.md`，Go Protobuf output 在提交前应被检查
- Runtime protocol adapter boundary rules 位于 `docs/runtime-protocol-adapter.md`
- Protocol envelope source 位于 `proto/vibit/protocol/v1/envelope.proto`
- Buf generation configuration 位于 `buf.yaml` 和 `buf.gen.yaml`
- 由 `.arch/protocol.yaml`、`docs/game-protocol.md` 和 `ADR-0015` 约束的 game-aware WebSocket Protobuf envelope
- SQL-first PostgreSQL migration source files 位于 `runtime/migrations/postgres/`
- S3-compatible object storage 作为计划中的大对象存储抽象，MinIO 作为本地/自托管优先候选，但需要先完成 dependency adoption
- `modules/<module>/module.yaml` 中遵循 `docs/module-manifest.md` 的 module manifests
- `modules/<module>/AGENTS.md` 中的 module-level agent guides
- Contract-first 的 commands、queries、events、errors、permissions 和 migrations
- `contracts/` 下的 contract source files，并由 `.arch/contracts.yaml` 登记
- 由 `.arch/dependencies.yaml` 登记 foundational dependency decisions
- 为重复框架结构生成 scaffolds 的 generators
- 能验证 dependency、contract、event 和 generated-file rules 的 architecture checks
- `changes/<date>-<change-id>/` 下遵循 `docs/change-spec.md` 的 change specs
- `conversations/` 下遵循 `docs/conversation-log.md` 的 conversation logs
- `decisions/` 下遵循 `docs/agent-decision-record.md` 的 Agent Decision Records
- `schema/` 下遵循 `docs/schema-validation.md` 的 schema validation
- `rules/` 下的 rule catalogs，首先是 `rules/check-rules.json`

第一个严肃 prototype 应证明一个命题：

> 给定一个新的后端需求，AI coding agent 能识别 affected module，更新正确 contracts，生成正确结构，实现行为，添加测试，运行验证，并更新文档，同时不破坏无关架构。

## CLI Prototype

第一版可执行标准位于：

```bash
tools/vibit
```

初始命令：

```bash
node tools/vibit --help
node tools/vibit check all
node tools/vibit check all --json
node tools/vibit check schemas
node tools/vibit check schemas --json
node tools/vibit check memory
node tools/vibit check memory --json
node tools/vibit check contracts
node tools/vibit check contracts --json
node tools/vibit check protocol
node tools/vibit check protocol --json
node tools/vibit check generated
node tools/vibit check generated --json
node tools/vibit check runtime
node tools/vibit check runtime --json
node tools/vibit inspect module inventory
node tools/vibit inspect boundary --from inventory --to player
node tools/vibit inspect contract --module inventory --type command --id GrantItem
node tools/vibit inspect change bootstrap-vibit-cli
node tools/vibit inspect memory
node tools/vibit inspect rule check.subcheck
node tools/vibit inspect rules --category check
node tools/vibit check architecture
node tools/vibit check architecture --json
node tools/vibit check change bootstrap-vibit-cli
node tools/vibit check change bootstrap-vibit-cli --json
node tools/vibit check module inventory
node tools/vibit check module inventory --json
node tools/vibit generate module <module>
```

当前 CLI 只使用 Node.js standard-library APIs。它是 architecture checks、inspection 和 generators 的 prototype，不是 server runtime，也不决定 server runtime language。

当 agent 在 intake、verification 或 handoff 阶段需要机器可读检查结果时，使用 `--json`。面向人类的文本输出仍是默认行为。

每条 JSON check result item 都包含稳定的 `rule_id` 和 `artifact`，让 agent 不必解析自然语言就能定位失败原因和相关产物。`check all --json` 是紧凑总览；需要完整细节时，对具体失败检查单独运行 `--json`。

使用 `node tools/vibit check memory` 可以验证 conversation logs 和 Agent Decision Records 的必需结构。

使用 `node tools/vibit check contracts` 可以验证 `.arch/contracts.yaml` 与已登记 contract source files 的一致性。

使用 `node tools/vibit check protocol` 可以在添加或修改 `.proto` files 前验证 manifest-to-Protobuf alignment。当还没有 `.proto` files 时，它会报告计划中的 protocol sources 和 messages；一旦 `.proto` files 存在，它会检查 package names、source traces、expected messages 和 field names。

第一批 Protobuf source files 现在定义 protocol envelope 和 inventory wire messages：

```text
proto/vibit/protocol/v1/envelope.proto
proto/vibit/inventory/v1/inventory.proto
```

`buf.yaml` 和 `buf.gen.yaml` 定义计划中的 generation path，但在使用已接受 toolchain 实际运行 generation 前，不提交生成的 Go Protobuf output。不要手工创建或编辑生成的 Go Protobuf files。

使用 `node tools/vibit check generated` 可以验证 module 声明的 generated files 存在，并且包含 generated、source 和 generator trace markers。它也会检查 `runtime/internal/generated/proto/` 下计划中的 Go Protobuf output root；generated Protobuf Go files 必须使用 `*.pb.go` 后缀，包含 `protoc-gen-go` generated-code marker，并能追溯到现有 `.proto` sources。

使用 `node tools/vibit check runtime` 做 server runtime verification。在 Go runtime 尚不存在前，该检查会报告 runtime implementation 尚未开始；当 `runtime/go.mod` 已存在但 Go source files 尚不存在时，它会验证 ADR-0014 skeleton。一旦 Go source files 存在，它必须发现 Go test files 并运行 Go runtime test path。

使用 `node tools/vibit inspect contract --module <module> --type <type> --id <id>` 可以在 agent intake 阶段以 JSON 查询单个已登记 command、query、event、error catalog 或 permission catalog。

使用 `node tools/vibit inspect change <change-id>` 可以查询 change spec 目录及其 verification metadata，而不必手动打开每个文件。

使用 `node tools/vibit inspect memory` 可以把 change specs、conversation logs 和 Agent Decision Records 列成机器可读的 project memory index。

Check output 的 rule metadata 位于 `rules/check-rules.json`。

使用 `node tools/vibit inspect rule <rule-id>` 可以查询单条 rule，而不必解析整个 catalog。

使用 `node tools/vibit inspect rules` 或 `node tools/vibit inspect rules --category <category>` 可以发现可用 rules。

第一版 server runtime 方向是 Go，并采用 modular monolith single-process server model。WebSocket 是第一版 gameplay/client protocol，Protobuf 是第一版 client/server wire format。Semantic business contracts 仍然保留在 vibit manifests 和 contract source files 中；Protobuf 负责 wire schema shape。见 `.arch/runtime.yaml`、`decisions/ADR-0008-go-server-runtime-language.md` 和 `decisions/ADR-0009-websocket-protobuf-client-protocol.md`。

第一版 game protocol framework 是 WebSocket-framed Protobuf envelope，使用显式 `kind`、`module` 和 `name` routing fields，并包含 request correlation、session metadata、target scopes、server-authoritative message rules 和 error mapping。第一版 endpoint 是 `/v1/ws`，并已由 Go runtime process 挂载。第一版 inventory slice 使用 player-scoped command/query/event/error/system messages；room state sync、matchmaking、allocation、reconnect replay、presence、streams、realtime input 和 state patches 仍然 deferred，直到它们拥有独立 modules 和 standards。见 `.arch/protocol.yaml`、`docs/game-protocol.md` 和 `decisions/ADR-0015-game-protocol-framework.md`。

Runtime protocol adapter boundary 定义在 `docs/runtime-protocol-adapter.md` 和 `decisions/ADR-0018-runtime-protocol-adapter-boundary.md` 中。WebSocket transport 拥有 frames，Protobuf adapter 拥有 envelope conversion，application dispatch 拥有 command/query routing，domain modules 拥有 invariants 和 behavior，generated packages 只提供 shapes。

Nakama 和 Pitaya 是 capability planning 的主动参考基线。Nakama 应指导 broad game backend product surface：accounts、sessions、storage、social systems、chat、groups、parties、leaderboards、tournaments、matchmaking、realtime multiplayer、authoritative matches 和 operations。Pitaya 应指导 Go game server architecture vocabulary：acceptors、sessions、routes、handlers、remotes/RPC、frontend/backend roles、groups、broadcast、serializers 和 cluster service discovery。`ADR-0078` 已 ratify Nakama/Pitaya-class product parity roadmap：vibit 应覆盖同级常用 capability families，同时保留 vibit 的 agent-native constraints，并避免意外 direct API compatibility。近期优先级是 runtime lifecycle closure，然后再扩展 presence、chat、social modules、matchmaking 或 match runtime。见 `.arch/reference.yaml`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`decisions/ADR-0019-nakama-and-pitaya-reference-baseline.md` 和 `decisions/ADR-0078-nakama-pitaya-product-parity-roadmap.md`。

PostgreSQL 是 runtime state 的第一版 authoritative durable relational store。S3-compatible object storage 计划用于 replays、snapshots、exports、binary assets 和 diagnostic archives 等大对象 artifacts。MinIO 是这个 S3-compatible 角色的本地/自托管优先候选，但在具体 use case 和 dependency adoption record 证明它必要之前，它不是 mandatory runtime dependency。Domain modules 必须使用 vibit-owned storage interfaces，而不是直接依赖 database drivers 或 object-storage clients。见 `decisions/ADR-0011-postgresql-and-object-storage-persistence.md`。

第一批已接受的 foundational runtime dependencies 记录在 `decisions/ADR-0013-first-go-runtime-dependencies.md` 和 `.arch/dependencies.yaml` 中。它们只被接受用于 platform adapters 和 generation tooling，不允许 domain modules 直接使用。S3 client tooling、MinIO deployment、observability 和外部 Go test framework adoption 仍然 deferred，直到具体 runtime needs 证明它们必要。

第一版 Go runtime layout 记录在 `decisions/ADR-0014-go-runtime-layout-and-boundaries.md` 中。Runtime 现在已有 generated Go Protobuf output、纯 application handoff types、把 generated envelopes 转换为 application route requests 的 Protobuf protocol adapter、用于 command 和 query routes 的小型 application dispatch skeleton、inventory runtime handlers、WebSocket transport adapter，以及 `/v1/ws` 的 minimal process wiring。PostgreSQL persistence、migrations、transaction wiring、authentication/session validation 和 generated route registration 尚未开始。未来 Go files 应遵循这些边界：

- `runtime/cmd/vibit-server/`：process startup、configuration wiring 和 lifecycle。
- `runtime/internal/app/`：command/query dispatch、application composition 和 transaction orchestration。
- `runtime/internal/platform/`：WebSocket、Protobuf、PostgreSQL、migration、event 和 transaction platform adapters。
- `runtime/internal/modules/<module>/`：手写 domain module logic。
- `runtime/internal/generated/`：生成的 Go contract 和 Protobuf output。

State-changing commands 应在 application-owned unit of work 中运行，然后才进行 repository mutation 和 domain-event recording。在 vibit 采纳明确的 event delivery 或 outbox standard 前，transaction 外的 event publication 继续 deferred。

`node tools/vibit check runtime` 现在会在 Go source files 存在时验证 skeleton、import boundaries、app/domain layer boundaries、runtime test discovery 和 Go runtime test path。

用下面命令启动第一版 runtime process：

```bash
cd runtime
go run ./cmd/vibit-server
```

当前 endpoint 和 manual verification notes 见 `docs/runtime-runbook.md`。

## 早期参考领域

小型游戏后端适合作为第一版演示领域，因为它天然包含状态、权限、事件、一致性规则和长期存在的模块。

建议模块：

- Player accounts
- Inventory
- Currency
- Rewards
- Tasks or quests
- Match sessions

第一条 backend slice 应强调 maintainability 和 agent workflow，而不是功能数量。但它仍应被视为长期维护系统的起点，而不是一次性 demo code。

## 治理

项目决策由 `CONSTITUTION.md` 管理。

在修改宪法原则、正式确定项目名、引入重大架构模式或进行破坏性标准变更前，应阅读宪法，并记录 motivation、alternatives、compatibility impact 和 migration path。

## 名称

`vibit` 是当前产品名。

预期类别短语是：

```text
agent-native server framework
```

正式定名前，应检查主要公开平台注册和项目重名风险。
