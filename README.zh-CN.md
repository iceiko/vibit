# vibit

vibit 是一个开源 agent-native server framework，用于构建 AI coding agents 能从第一性原理理解、扩展、验证和维护的后端系统。

状态：pre-alpha，正在推进 `v0.1 alpha`

短期目标是做出第一个 developer-usable `v0.1 alpha`：一个 single-node、PostgreSQL-backed、WebSocket + Protobuf game backend runtime，让真实开发者能在本地运行、检查，并把它作为加入开发的起点。

长期目标是成为 AI 时代的 Nakama 或 Pitaya：同一产品级别的严肃 game/backend server capability，但围绕 vibit 的 agent-native maintainability 模型重构。这不表示 direct Nakama 或 Pitaya API compatibility。

## 当前已有内容

本仓库已经不只是纯设计阶段。当前实现基础包括：

- Agent-readable governance：`CONSTITUTION.md`、`AGENTS.md`、change specs、ADRs、conversation logs 和 machine-readable architecture manifests。
- `runtime/` 下的 Go runtime。
- `/v1/ws` WebSocket gameplay endpoint。
- Protobuf envelope 和生成的 Go Protobuf output。
- PostgreSQL migration sources 和 platform persistence adapters。
- Inventory proof slice。
- Player account persistence boundaries 和 adapters。
- Device credential login service behavior 和 protocol route。
- Protected routes 的 opaque access-token validation。
- Logout service behavior 和 protocol route。
- Runtime session persistence 和 response session metadata。
- First-message WebSocket connection binding。
- Single-process active connection lifecycle、close handoff、reconnect epoch handling 和 presence lifecycle snapshot。
- 面向 agents 和 humans 的 `tools/vibit` checks 与 inspection commands。

这还不是完成态 alpha。Authenticated gameplay end-to-end path 现在已经通过 focused Go test 证明，runbook 和 request-loop script 已存在，runtime 也暴露了一个很小的 health/readiness/version/config surface，alpha acceptance checklist 已经记录本地 readiness state，并且 `docs/alpha-developer-flow.md` 已把这些入口整理成一条 coherent local developer journey。当前最关键的剩余缺口是 release publishing decision gate。

## 本地试用

当前开发工作流 prerequisites：

- Go，用于 runtime。
- Node.js，用于 `tools/vibit`。
- PostgreSQL，用于 persistent runtime path。
- Buf 和 Protobuf tooling，用于重新生成 Protobuf output。

运行仓库检查：

```bash
node tools/vibit check all
```

运行 Go runtime tests：

```bash
cd runtime
go test ./...
```

启动 bootstrap in-memory runtime：

```bash
cd runtime
go run ./cmd/vibit-server
```

Runtime 默认监听 `:8080`，并挂载：

```text
/v1/ws
/healthz
/readyz
/version
/configz
```

该 endpoint 期望 binary WebSocket messages，其中包含 `vibit.protocol.v1.Envelope` Protobuf bytes。这个 gameplay endpoint 不接受 JSON。

`/healthz`、`/readyz`、`/version` 和 `/configz` 是小型 JSON troubleshooting endpoints。`/configz` 只报告 redacted runtime posture，不暴露 verifier keys、raw credentials、raw tokens、DSNs、digests、headers、cookies、query strings、subprotocol values、remote addresses 或 concrete transport metadata。

PostgreSQL runtime path 更完整，但它需要 migrations、`VIBIT_POSTGRES_DSN` 和 authentication verifier key environment variables。当前 operational notes 见 `docs/runtime-runbook.md`。Runbook 是 v0.1 alpha hardening path 的一部分，应视为开发文档，而不是已打磨的 release guide。

## 下一目标：v0.1 Alpha

持久目标记录在 `docs/v0.1-alpha-goal.md`。

当前 local acceptance checklist 记录在 `docs/alpha-acceptance-checklist.md`。

Packaged local developer journey 记录在 `docs/alpha-developer-flow.md`。

`v0.1 alpha` 应让具备技术能力的开发者能够：

- clone repo；
- 准备 local config 且不提交 secrets；
- apply 或 verify required PostgreSQL migrations；
- 创建或获得 first device credential；
- 通过 protocol route 登录；
- 获得 opaque access token 和 runtime session metadata；
- bind WebSocket connection；
- 调用 protected inventory route；
- 查询 presence；
- logout；
- 运行 checks；
- 找到明确的下一步 contribution。

推荐下一步顺序：

1. 完成 `W-0178`：protected presence protocol query。已完成。
2. 先定义并添加 first local onboarding/device credential issuance。已完成。
3. 选择并证明 onboarding -> login -> bind connection -> protected inventory -> presence query -> logout 的 end-to-end path。已完成。
4. 围绕真实 alpha path 刷新 runtime runbook。已完成。
5. 添加最小 example client 或 request-loop script。已完成。
6. 添加 health/readiness/version/config surfaces。已完成。
7. 添加 alpha acceptance checklist 或 check。已完成。
8. 整理 alpha developer flow，并记录 prerequisites。已完成。
9. 定义 release publishing decision gate。

运行 minimal local alpha request loop：

```bash
examples/local-alpha-request-loop.sh
```

它包装 focused authenticated gameplay E2E proof，并且不会打印 raw credentials、raw access tokens、verifier keys、DSNs、digests 或 transport metadata。

## 继续开发

如果你是继续推进项目的 agent 或 contributor，先运行：

```bash
node tools/vibit inspect next
node tools/vibit check work --json
```

在本 README 更新时，next ready item 是：

```text
W-0190 Define release publishing decision gate
```

使用 `.arch/work-items.yaml` 作为 continuation source of truth。`continue` 或 `继续推进` 的意思是：推进一个 `next_ready` work item，除非遇到 ask-first boundary 或 verification failure。

## Agent-Native 的含义

Agent-native 主要不是指服务器带有 AI 功能。

它指的是代码库被设计成让 AI coding agents 能可靠工作：

- architecture rules 是显式的；
- module ownership 是声明的；
- public behavior 是 contract-first；
- generated structure 可追溯；
- business rules 有测试；
- change workflow 有边界；
- repository state 可检查；
- project memory 存在持久 artifacts 中，而不是只留在聊天记录里。

NPC agents、memory、model routing、tool calling 和 simulations 等 AI gameplay features 未来可以成为扩展，但它们不是项目基础。

## Nakama 与 Pitaya 目标

Nakama 和 Pitaya 是 capability planning 的主动参考基线。

- Nakama 指导 broad game backend surface：accounts、sessions、storage、social systems、chat、groups、parties、leaderboards、tournaments、matchmaking、realtime multiplayer、authoritative matches、operations、SDKs 和 examples。
- Pitaya 指导 Go game server architecture vocabulary：acceptors、sessions、routes、handlers、remotes/RPC、frontend/backend roles、groups、broadcast、serializers 和 cluster service discovery。

vibit 应随时间覆盖同级常用能力，同时保留自己的架构：explicit manifests、contracts、generation、tests、ADRs、repository checks 和 bounded agent workflow。

参见：

- `docs/reference-game-server-alignment.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/v0.1-alpha-goal.md`

## 项目地图

重要入口：

- `CONSTITUTION.md`：权威项目宪法。
- `AGENTS.md`：仓库级 coding agent 操作指南。
- `.arch/README.md`：architecture manifest 入口。
- `.arch/work-items.yaml`：continuation queue。
- `.arch/runtime.yaml`：runtime readiness 和 implementation state。
- `.arch/reference.yaml`：Nakama/Pitaya reference 和 product parity planning。
- `docs/v0.1-alpha-goal.md`：短期 alpha 和长期产品目标。
- `docs/alpha-developer-flow.md`：packaged local alpha developer journey。
- `docs/alpha-acceptance-checklist.md`：本地 v0.1 alpha acceptance checklist。
- `docs/runtime-runbook.md`：当前 runtime startup 和 verification notes。
- `docs/nakama-pitaya-product-parity-roadmap.md`：长期 capability roadmap。
- `changes/`：具体 change specs 和 verification records。
- `conversations/`：持久 maintainer-agent project memory。
- `decisions/`：Agent Decision Records。
- `runtime/`：第一版 Go reference runtime。
- `proto/`：Protobuf protocol source files。
- `tools/vibit`：architecture check、inspection 和 generator CLI。

英文文档是权威版本。简体中文译本服务于人类阅读和早期项目讨论。

## CLI

当前 CLI 是：

```bash
node tools/vibit --help
node tools/vibit inspect next
node tools/vibit inspect work
node tools/vibit inspect reference
node tools/vibit check all
node tools/vibit check all --json
node tools/vibit check runtime
node tools/vibit check work
node tools/vibit check memory
node tools/vibit check schemas
```

当 agent 在 intake、verification 或 handoff 阶段需要机器可读结果时，使用 `--json`。

## 治理

项目决策由 `CONSTITUTION.md` 管理。

在修改 constitutional principles、引入重大 architecture pattern、改变 public protocol shape、添加 dependencies 或改变 release direction 前，应在 change spec 和必要的 ADR 中记录 motivation、alternatives、compatibility impact 和 migration path。

`vibit` 是产品名。预期类别短语是：

```text
agent-native server framework
```
