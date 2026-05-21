# vibit

vibit 是一个开源 **agent-native server framework**，用于构建 AI coding agents 和人类开发者都能从第一性原理理解、扩展、验证和维护的后端系统。

最新 source alpha：`v0.1.0-alpha.1`

这是一个面向开发者的早期 alpha：适合检查架构、在本地跑通第一条 authenticated gameplay loop，并参与塑造一个 AI-maintainable backend framework。它还不是 production server distribution。

## 为什么值得试

大多数 backend frameworks 设计时，AI coding agents 还不是工程流程的一部分。它们可能很强，但 architecture rules、ownership boundaries、change history 和 verification paths 往往散落在代码、文档、issues 和 maintainer 记忆里。

vibit 正在验证一种更严格的模型：

- architecture rules 存在 machine-readable manifests 中；
- changes 由 work items、specs、ADRs 和 verification records 约束；
- public behavior 是 contract-first；
- generated output 可追溯；
- module ownership 是显式的；
- agents 可以用 `tools/vibit` 检查下一步安全任务；
- humans 不需要翻旧聊天记录也能审计一个变更为什么存在。

第一个领域重点是 game/backend servers。长期目标是成为 AI 时代的 Nakama/Pitaya-class open-source backend framework，但围绕 agent-native maintainability 重新组织，而不是追求 direct API compatibility。

## 试用 Alpha

最快 source checkout 路径：

```bash
git clone https://github.com/iceiko/vibit.git
cd vibit
node tools/vibit check all
cd runtime && go test ./...
cd .. && examples/local-alpha-request-loop.sh
```

今天这条路径能证明：

- repository architecture checks 可以通过 `tools/vibit` 运行；
- Go runtime tests 可以运行；
- local alpha request loop 会跑通 authenticated gameplay path；
- script 不会打印 raw credentials、raw access tokens、verifier keys、DSNs、digests 或 transport metadata。

Prerequisites：

- Go；
- Node.js；
- PostgreSQL，用于 persistent runtime path；
- Buf 和 Protobuf tooling，仅在重新生成 Protobuf output 时需要。

## 当前已有内容

本仓库已经不只是纯设计阶段。当前实现基础包括：

- `runtime/` 下的 Go runtime；
- `/v1/ws` WebSocket gameplay endpoint；
- Protobuf envelope 和 generated Go Protobuf output；
- PostgreSQL migration sources 和 platform persistence adapters；
- inventory proof slice；
- player account persistence boundaries 和 adapters；
- development 用 local onboarding/device credential issuance；
- device credential login service behavior 和 protocol route；
- protected routes 的 opaque access-token validation；
- runtime session persistence 和 response session metadata；
- first-message WebSocket connection binding；
- protected inventory 和 protected presence query path；
- logout service behavior 和 protocol route；
- single-process active connection lifecycle、close handoff、reconnect epoch handling 和 presence lifecycle snapshot；
- health、readiness、version 和 redacted config endpoints；
- 面向 agents 和 humans 的 `tools/vibit` checks 与 inspection commands。

Runtime 默认监听 `:8080`，并挂载：

```text
/v1/ws
/healthz
/readyz
/version
/configz
```

`/v1/ws` 期望 binary WebSocket messages，其中包含 `vibit.protocol.v1.Envelope` Protobuf bytes。这个 gameplay endpoint 不接受 JSON。

`/healthz`、`/readyz`、`/version` 和 `/configz` 是小型 JSON troubleshooting endpoints。`/configz` 只报告 redacted runtime posture，不暴露 verifier keys、raw credentials、raw tokens、DSNs、digests、headers、cookies、query strings、subprotocol values、remote addresses 或 concrete transport metadata。

## 适合谁

如果你符合下面任意一项，现在就值得试 vibit：

- 你构建或运营 game/backend servers；
- 你用过或评估过 Nakama、Pitaya、Colyseus、Pomelo、Agones 或 custom Go backends；
- 你希望 AI coding agents 在严肃 backend codebase 中更安全地改代码；
- 你关心 explicit architecture、contracts、generated structure、tests 和 durable decision records；
- 你想在形态固化前参与定义第一个有用的 agent-native server framework。

这个 alpha 还不适合 production deployment、plug-and-play SDK use、hosted operations，或需要 packaged binaries/containers 的团队。

## 当前限制

`v0.1.0-alpha.1` 是 source-first：

- 没有 release binaries；
- 没有 packages；
- 没有 container images；
- 没有 checksum files；
- 没有 provenance 或 signing artifacts；
- 没有 hosted deployment；
- 没有 install script；
- 没有 SDK package；
- 没有 direct Nakama/Pitaya API compatibility 承诺。

PostgreSQL runtime path 是目前最完整的本地路径，但仍然需要 development setup：migrations、`VIBIT_POSTGRES_DSN` 和 authentication verifier key environment variables。见 `docs/runtime-runbook.md`。

## Release Notes

Source alpha release notes 位于：

- `docs/releases/v0.1.0-alpha.1.md`
- `docs/releases/v0.1.0-alpha.1.zh-CN.md`

Release authorization record 是：

- `docs/release-execution-final-authorization.md`
- `decisions/ADR-0103-release-execution-final-authorization.md`

First alpha user discovery loop 是：

- `docs/first-alpha-user-discovery-loop.md`
- `decisions/ADR-0104-first-alpha-user-discovery-loop.md`

First alpha feedback intake surface 和 product maturity milestones 是：

- `.github/ISSUE_TEMPLATE/alpha-feedback.yml`
- `docs/first-alpha-feedback-intake-surfaces.md`
- `docs/product-maturity-milestones.md`
- `decisions/ADR-0105-first-alpha-feedback-intake-and-product-maturity-milestones.md`

Prototype-ready foundation execution plan 是：

- `docs/prototype-ready-foundation-execution-plan.md`
- `decisions/ADR-0106-prototype-ready-foundation-execution-plan.md`

Prototype-ready local development path gate 是：

- `docs/prototype-ready-local-development-path-gate.md`
- `decisions/ADR-0107-prototype-ready-local-development-path-gate.md`

## 继续开发

如果你是继续推进项目的 agent 或 contributor，先运行：

```bash
node tools/vibit inspect next
node tools/vibit check work --json
```

当前 next work item 是：

```text
W-0200 Implement prototype-ready local development path package
```

使用 `.arch/work-items.yaml` 作为 continuation source of truth。`continue` 或 `继续推进` 的意思是：推进一个 `next_ready` work item，除非遇到 ask-first boundary、verification failure 或 required maintainer confirmation。

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
- `docs/release-publishing-decision-gate.md`：release publishing decision boundary。
- `docs/release-execution-preparation-gate.md`：release execution preparation boundary。
- `docs/release-execution-authorization-gate.md`：release execution authorization criteria。
- `docs/release-execution-final-authorization.md`：final release authorization record。
- `docs/first-alpha-user-discovery-loop.md`：first alpha user discovery loop。
- `docs/first-alpha-feedback-intake-surfaces.md`：first alpha feedback intake surface。
- `docs/product-maturity-milestones.md`：source alpha、prototype-ready、production-candidate 和 product-class milestones。
- `docs/releases/v0.1.0-alpha.1.md`：alpha release notes。
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
