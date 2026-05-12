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
- `schema/`：用于机器可检查 standards 的 JSON Schema files
- `rules/`：面向机器可读 check metadata 的 rule catalogs

英文文档是权威版本。简体中文译本服务于人类阅读和早期项目讨论。

## 预期方向

vibit 应逐步演进出：

- `.arch/` 下的 architecture manifests
- `modules/<module>/module.yaml` 中遵循 `docs/module-manifest.md` 的 module manifests
- `modules/<module>/AGENTS.md` 中的 module-level agent guides
- Contract-first 的 commands、queries、events、errors、permissions 和 migrations
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
node tools/vibit inspect module inventory
node tools/vibit inspect boundary --from inventory --to player
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

当前 CLI 只使用 Node.js standard-library APIs。它是 architecture checks 和 module skeleton generation 的 prototype，不是 server runtime。

当 agent 在 intake、verification 或 handoff 阶段需要机器可读检查结果时，使用 `--json`。面向人类的文本输出仍是默认行为。

每条 JSON check result item 都包含稳定的 `rule_id` 和 `artifact`，让 agent 不必解析自然语言就能定位失败原因和相关产物。`check all --json` 是紧凑总览；需要完整细节时，对具体失败检查单独运行 `--json`。

使用 `node tools/vibit inspect change <change-id>` 可以查询 change spec 目录及其 verification metadata，而不必手动打开每个文件。

使用 `node tools/vibit inspect memory` 可以把 change specs、conversation logs 和 Agent Decision Records 列成机器可读的 project memory index。

Check output 的 rule metadata 位于 `rules/check-rules.json`。

使用 `node tools/vibit inspect rule <rule-id>` 可以查询单条 rule，而不必解析整个 catalog。

使用 `node tools/vibit inspect rules` 或 `node tools/vibit inspect rules --category <category>` 可以发现可用 rules。

服务器实现语言和整体服务器实例架构目前有意保持开放。它们应在仓库标准和工具链足够强、可以评估 tradeoffs 后，通过 change specs 决定。

## 早期参考领域

小型游戏后端适合作为第一版演示领域，因为它天然包含状态、权限、事件、一致性规则和长期存在的模块。

建议模块：

- Player accounts
- Inventory
- Currency
- Rewards
- Tasks or quests
- Match sessions

Demo 应强调 maintainability 和 agent workflow，而不是功能数量。

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
