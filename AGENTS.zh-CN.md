# Agent 操作指南

状态：Draft v0.1  
最后更新：2026-05-12  
范围：仓库级 coding agent 操作指令  
权威来源：`CONSTITUTION.md`  
说明：本文件是 `AGENTS.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

这份指南把宪法转化为 agent 的工作规则。它不取代宪法。如果本指南与 `CONSTITUTION.md` 冲突，应遵循 `CONSTITUTION.md`，并更新本指南。

配套英文源文档是 `AGENTS.md`。英文文件是权威版本。

## 1. 项目身份

工作名：

```text
vibit
```

类别：

```text
Agent-Native Server Framework
```

定位：

```text
vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.
```

在本仓库中，“AI-native” 首先指 agent-native maintainability。它主要不是指添加 AI gameplay features 或 AI product features。

## 2. 必读内容

在进行非平凡变更前，先阅读：

- `CONSTITUTION.md`
- 本文件
- `.arch/` 下相关 architecture manifests，如果它们已经存在
- `modules/<module>/module.yaml` 中的相关 module manifest，如果它已经存在
- `modules/<module>/AGENTS.md` 中的相关 module guide，如果它已经存在
- `changes/` 下相关 change spec，如果该变更已有 spec

如果预期产物尚不存在，不要发明隐藏假设。要么把缺失产物作为本次变更的一部分创建出来，要么记录它目前尚不可用。

## 3. 当前仓库状态

本仓库当前处于宪法和标准设计阶段。

现有基础：

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/README.md`
- `.arch/modules.yaml`
- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `docs/module-manifest.md`
- `docs/module-manifest.zh-CN.md`
- `docs/change-spec.md`
- `docs/change-spec.zh-CN.md`
- `changes/_template/`
- `docs/conversation-log.md`
- `docs/conversation-log.zh-CN.md`
- `conversations/`
- `docs/agent-decision-record.md`
- `docs/agent-decision-record.zh-CN.md`
- `decisions/`
- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `schema/`
- `rules/`

框架实现代码、generators、modules 和 verification commands 可能尚不存在。如果它们不存在，应记录 verification 当前不可用，而不是假装已经运行。

当前 runtime readiness decisions 指向 TypeScript on Node.js 作为第一参考实现、modular monolith single-process server model、contract-first commands/queries/events/errors/permissions，以及 `inventory` 作为优先的第一 proof slice。在创建 runtime implementation code 前，必须阅读 `.arch/runtime.yaml` 以及 `ADR-0003` 到 `ADR-0006`。

当前可执行工具：

```bash
npm run typecheck
npm test
npm run check
node tools/vibit --help
node tools/vibit check all
node tools/vibit check all --json
node tools/vibit check schemas
node tools/vibit check schemas --json
node tools/vibit check memory
node tools/vibit check memory --json
node tools/vibit check contracts
node tools/vibit check contracts --json
node tools/vibit check generated
node tools/vibit check generated --json
node tools/vibit check runtime
node tools/vibit check runtime --json
node tools/vibit inspect module <module>
node tools/vibit inspect boundary --from <module> --to <module>
node tools/vibit inspect contract --module <module> --type <type> --id <id>
node tools/vibit inspect change <change-id>
node tools/vibit inspect memory
node tools/vibit inspect rule <rule-id>
node tools/vibit inspect rules
node tools/vibit inspect rules --category <category>
node tools/vibit generate contract --module <module> --type <type> --id <id>
node tools/vibit check architecture
node tools/vibit check architecture --json
node tools/vibit check change <change-id>
node tools/vibit check change <change-id> --json
node tools/vibit check module <module>
node tools/vibit check module <module> --json
node tools/vibit generate module <module>
```

当 CLI tooling 可用时，默认使用 `node tools/vibit check all` 作为仓库验证命令。

当前 TypeScript package baseline 使用 npm scripts：

- `npm run typecheck` 对 runtime files 运行 no-emit TypeScript checks。
- `npm test` 运行 Node.js built-in module runtime tests。
- `npm run check` 运行 `node tools/vibit check all`。

根 TypeScript runtime files 使用 ESM。`tools/package.json` 让现有 `tools/vibit` CLI 保持在 CommonJS scope 内。

当 agent 在 intake、verification 或 handoff 阶段需要机器可读检查结果时，使用 `--json`。

每条 JSON check result item 都应包含 `rule_id` 和 `artifact`。把 `check all --json` 视为紧凑总览；需要完整细节时，对具体失败检查单独运行 `--json`。

当新增或修改 conversation logs 或 Agent Decision Records 时，使用 `node tools/vibit check memory`。

当新增或修改 contract source files 或 `.arch/contracts.yaml` 时，使用 `node tools/vibit check contracts`。

当新增或修改 generated files 或 module manifest 中的 `generated` declarations 时，使用 `node tools/vibit check generated`。

当新增或修改 runtime module behavior 或 tests 时，使用 `node tools/vibit check runtime`。当 package baseline 存在时，该检查会先运行 TypeScript typecheck，再运行 runtime tests。

当 agent 在 intake 阶段需要以 JSON 读取单个 contract 的 registry entry、source summary、module manifest declaration 和 consistency status 时，使用 `node tools/vibit inspect contract --module <module> --type <type> --id <id>`。

使用 `node tools/vibit generate contract --module <module> --type <type> --id <id>` 从 contract source files 重新生成已声明的 contract shapes，而不是手工编辑 generated output。

当 change spec 已存在，并且 agent 在 intake 或 handoff 阶段需要结构化了解它的文件、metadata、affected modules 和 verification state 时，使用 `node tools/vibit inspect change <change-id>`。

当 agent 在决定完整阅读哪些 artifacts 之前，需要结构化索引 change specs、conversation logs 和 Agent Decision Records 时，使用 `node tools/vibit inspect memory`。

使用 `rules/check-rules.json` 解读 check result 中的 `rule_id`。

当只需要单条 rule metadata 时，使用 `node tools/vibit inspect rule <rule-id>`。

使用 `node tools/vibit inspect rules --category <category>` 按 category 发现 rules。

使用 `.arch/runtime.yaml` 作为 runtime readiness 的机器可读 intake 入口。它链接了约束语言、服务器实例模型、contract 与 generation boundary，以及第一 proof slice 的 ADR。

## 4. 文档规则

英文是项目权威文档语言。

每个面向公众的文档都应该有：

- 英文源文档
- 简体中文可读译本

命名示例：

```text
CONSTITUTION.md
CONSTITUTION.zh-CN.md
AGENTS.md
AGENTS.zh-CN.md
docs/<name>.md
docs/<name>.zh-CN.md
.arch/README.md
.arch/README.zh-CN.md
```

规则：

- 英文源文档发生实质变更时，应在同一次变更中更新中文译本。
- 如果同一次变更无法更新译本，必须明确标记译本已过期。
- 机器可读标识符保持英文。
- 代码标识符、模块名、命令、事件、权限、错误码应使用英文，除非存在强领域理由。
- 翻译应保留意义。不要为了逐字翻译牺牲清晰度。

## 5. 标准变更工作流

每个非平凡 feature、bug fix、migration、refactor 或 standard change 都应遵循：

1. 澄清 requirement。
2. 识别 affected modules 和 contracts。
3. 当变更足够大、需要持久上下文时，编写或更新 change spec。
4. 当 public behavior 改变时，先更新 schemas、manifests 或 contracts，再实现。
5. 当 generators 存在时，用它们生成重复结构。
6. 只在声明过的边界内实现。
7. 新增或更新聚焦的 tests。
8. 运行相关 verification commands。
9. 更新文档和译本。
10. 记录已验证和未验证的内容。

对于早期纯设计变更，涉及代码、测试、生成器和验证的步骤可以不适用，但必须明确说明。

## 6. 架构规则

优先选择这样的设计：

- 给 agents 更少歧义上下文
- 创建更强 module boundaries
- 让 behavior 更容易验证
- 让 contracts 显式
- 减少 hidden coupling
- 支持 code generation
- 对人类开发者仍然实用

不要让架构规则只存在于维护者记忆里。如果一条规则重要，它最终应体现为文档、schema、manifest、test、generator 或 architecture check。

## 7. 模块规则

当 modules 存在时，每个 module 应声明：

- 它拥有什么
- 它不拥有什么
- Public commands
- Public queries
- Published events
- Subscribed events
- Allowed dependencies
- Forbidden dependencies
- Invariants
- Required tests
- Generated files
- Handwritten extension points

其他 modules 不能直接访问某个 module 的内部实现。跨模块通信应通过 commands、queries、events、public module APIs 或 generated clients 完成。

`modules/<module>/module.yaml` 应以 `docs/module-manifest.md` 作为源标准。

`changes/<date>-<change-id>/` 应以 `docs/change-spec.md` 作为源标准。

`conversations/` 应以 `docs/conversation-log.md` 作为源标准。

当 maintainer 引入 product intent、拒绝一种解释、命名一个概念或做出架构决策时，应在 conversation log 中保留该上下文。提交前必须脱敏 secrets。

`decisions/` 应以 `docs/agent-decision-record.md` 作为源标准。

当某个 decision 影响长期 architecture、generated file conventions、module ownership 或拒绝了一个合理替代方案时，应创建或更新 Agent Decision Record。Rationale 应简洁、公开；不要存储隐藏 chain-of-thought。

Generated files 对 non-system agents 不可变。如果 generated output 错了，应修改 source schema、template 或 generator，除非 change spec 或 decision record 显式授予 `generated_file_override`。

`schema/` 应以 `docs/schema-validation.md` 作为源标准。

当修改 module manifests、change specs、Agent Decision Records 或 tool JSON output 的结构时，必须更新对应 schema file，并运行 `node tools/vibit check schemas`。

## 8. 契约规则

公共行为应先定义，再实现。

带有 contract 属性的产物可以包括：

- API schemas
- Command schemas
- Query schemas
- Event schemas
- Error catalogs
- Permission catalogs
- Database migration schemas
- Generated clients

规则：

- Public contracts 必须先声明再使用。
- 对兼容性敏感的 contracts 必须 versioned。
- Breaking changes 必须显式说明。
- Generated output 必须能追踪到 source schema。
- 不要手工编辑 generated contract output，除非本次变更的目标就是 generator 本身。

## 9. 测试与验证

测试是架构的一部分，不是收尾步骤。

当实现代码存在时，相关验证可以包括：

- Unit tests
- Contract tests
- Invariant tests
- Integration tests
- Migration tests
- Replay tests
- Architecture checks
- Generator checks
- Documentation consistency checks

本仓库尚未定义最终 verification commands。在此之前，按以下方式记录验证：

```text
Verified: <commands or checks run>
Not verified: <reason>
Not applicable: <reason>
```

没有运行验证时，绝不能声称变更已经验证。

## 10. 先询问

以下情况应先询问人类维护者：

- 修改宪法原则
- 正式确定或替换项目名
- 重新定义 module ownership
- 引入新的架构模式
- 做 breaking API、command、query 或 event changes
- 修改 generated file conventions
- 删除 tests
- 削弱 validation 或 permission checks
- 在 modules 之间迁移 data ownership
- 添加重大的外部框架依赖

## 11. 禁止事项

禁止：

- 把 AI gameplay features 当作本项目基础
- 为方便绕过 module boundaries
- 把业务逻辑藏在 transport handlers
- 添加未登记 public events
- 添加未登记 permissions
- 添加无类型 cross-module payloads
- 在没有声明边界的情况下做大范围仓库编辑
- 无记录地手工编辑 generated files
- 英文公共文档发生实质变更后，让中文译本静默落后
- 没有运行验证却声称已验证

## 12. 新增标准时

新增标准应说明：

- 要解决的问题
- 引入的规则
- 该规则为什么帮助 agents
- 对人类的影响
- 预期产物
- 验证路径
- 从现有工作的迁移路径

优先选择可以被执行和检查的小标准，而不是无法检查的宏大表述。

## 13. 新增实现代码时

不要一开始就把框架代码分散到整个仓库。

从能证明核心命题的最小完整切片开始：

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

一个好的第一版实现目标，应包含一个小而完整的后端领域，例如 player accounts、inventory、currency、rewards、quests 或 match sessions。

## 14. 自举控制

Self-bootstrapping 只有在它能改善通向可工作服务器框架的路径时才有价值。

在新增 standard、inspect command、check command、schema、generator 或 workflow rule 前，先确认它直接支持以下至少一项：

- 下一个 runtime vertical slice
- 具体 module boundary
- Public contract 或 generated shape
- Test 或 verification path
- 针对预期实现任务减少 agent context

如果收益主要只是让 tooling 本身更完整，应推迟。

当仓库已经有足够 tooling 可以尝试一个小的端到端后端能力时，优先做 runtime readiness work，而不是继续增加 meta-tooling，然后再实现 runtime slice。

Runtime readiness 只应回答让第一个 slice 成立所必需的决策：

- Implementation language 和 package layout
- Minimal server instance model
- 第一个 module 和 capability boundary
- Contract format
- Generated files 与 handwritten files 的边界
- 最小 test 和 verification strategy
- Persistence 和 migration 假设

当这些选择仍然含糊时，不要仓促进入实现。但当 readiness work 不再改变第一个 slice 如何构建、验证或维护时，也不要继续扩张准备工作。

例外情况应记录在 change spec 或 Agent Decision Record 中。

## 15. 交接要求

每次变更结束时，都要给下一个 agent 或 human 留下足够上下文。

记录：

- 改了什么
- 为什么改
- 哪些文件变更
- 哪些 contracts 或 standards 变更
- 已验证什么
- 未验证什么
- 还有哪些开放问题

如果工作未完成，说明下一个具体动作。
