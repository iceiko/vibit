# Agent-Native Server Constitution 中文版

状态：Draft v0.2  
最后更新：2026-05-12  
范围：开源 agent-native 游戏/后端服务器架构  
产品名：vibit  
说明：本文件是 `CONSTITUTION.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

## 0. 目的

这份文档是项目宪法。

它不是为了描述某一个实现细节，而是为了定义所有实现细节都必须遵守的原则、边界、标准和决策规则。

这个项目的目标，是构建一个让 AI coding agents 能够低摩擦地理解、扩展、验证和维护的服务器框架。

在这份宪法里，“AI-native” 主要不是指服务器提供 AI 功能，而是指代码库、架构、工作流和扩展模型从第一天起就面向 agent-driven development 设计。

核心论断：

> 如果一个服务器架构足够显式、机器可读、边界清晰、可测试，并且由稳定契约生成结构，那么它就能显著降低 AI agents 正确工作的难度。

## 0.1 命名标准

项目需要一个对人类好记、同时能让已经隐约意识到这个问题的人一眼看懂的名字。

已确定的产品名：

> vibit

含义：

- `vi`：短小的前缀，可以联想到 visibility、vitality 和清晰接口。
- `bit`：数字系统里最熟悉的最小单位。这个名字表达小型、可组合、可验证的构建块。

推荐描述性副标题：

> vibit: an agent-native server framework for AI-maintainable backends.

长描述：

> vibit 是一个开源 agent-native server framework，用于构建 AI coding agents 能从第一性原理理解、扩展、验证和维护的后端系统。

命名规则：

- 公开名称必须简短、可读、易记。
- 副标题必须立刻说明项目类别。
- 名字不能让人误解为“服务器只是带了一些 AI 功能”。
- 名字应优先适配游戏服务器场景，但不能阻碍未来扩展到通用后端。
- 仓库描述中应包含 `agent-native server framework`。
- 最终定名前，应检查主要公开平台注册和项目重名风险。

考虑过的候选名：

- `AgentNative`：清楚，但太泛，且已有公开使用。
- `AgentFirst`：清楚，但范围过宽，且已有公开使用。
- `AgentReady`：短语有用，但已有代码库评估和 readiness 工具使用。
- `AgentFrame`：描述性强，但太泛，容易和 agent orchestration frameworks 混淆。
- `AgentForge`：好记，但已有大量使用。
- `Framewright`：有识别度，也和 framework-building 有关，但现在由更简洁的产品名 `vibit` 取代。

使用 `vibit` 作为产品名，使用 `Agent-Native Server Framework` 作为类别名。

## 0.2 文档语言标准

项目权威文档语言是英文。

每个面向公众的项目文档都应该有：

- 英文源文档
- 简体中文可读译本

英文文件是权威版本。中文文件是译本。

推荐命名：

```text
CONSTITUTION.md
CONSTITUTION.zh-CN.md
AGENTS.md
AGENTS.zh-CN.md
docs/<name>.md
docs/<name>.zh-CN.md
```

规则：

- 英文文档是架构、契约、标准和治理的事实来源。
- 中文文档服务于人类阅读，尤其是更习惯用中文思考的贡献者。
- 翻译应保留意义，不追求逐字翻译。
- 英文源文档发生实质变更时，对应中文译本必须在同一次变更中更新，或者明确标记为已过期。
- 机器可读 manifest 和 schema 应使用英文标识符。
- 代码标识符、模块名、命令、事件、权限、错误码应使用英文，除非有强领域理由。

这条规则的目的，是让项目对全球 agents 和工具保持可读，同时服务早期用中文讨论和塑造项目思想的人。

## 1. 项目定义

这个项目是一个 Agent-Native Server Framework。

如果聚焦游戏领域，它就是一个 Agent-Native Game Backend Framework。

它的目标，是让人类开发者和 AI agents 能通过标准化模块、spec、contract、测试、生成结构和自动化架构检查，安全地演进一个长期存在的服务器代码库。

项目未来可以提供 AI gameplay 功能，例如 NPC agents、memory、model routing 或 tool-using characters。这些都是合理的扩展，但它们不是基础。

基础是：

- Agent-friendly architecture
- Agent-readable project context
- Agent-safe change workflows
- Agent-verifiable behavior
- 面向 modules、APIs、events、data、tests 和 documentation 的 agent-compatible standards

## 2. 外部标准立场

我们应该参考已有标准和方法论，但不能盲目继承它们。

已知参考包括：

- Superpowers：可参考其结构化 agent workflows、skills、planning、test-first development 和 verification discipline。
- AGENTS.md：可参考其作为仓库级 agent instructions 的新兴约定。
- Spec-first development：可参考其“先需求和契约，后实现”的纪律。
- DDD、CQRS、Event Sourcing、Actor Model、Hexagonal Architecture：可参考其显式边界、状态迁移和依赖控制思想。
- OpenAPI、AsyncAPI、JSON Schema、Protobuf、GraphQL、database migration schemas：可参考其 contract-first implementation。
- Skynet、Nakama、Colyseus、Agones、Pomelo、ET、Zinx 等游戏/后端框架：可参考其工程经验，但不作为治理标准。

我们的规则：

> 外部系统为我们提供词汇、流程纪律和成熟模式。凡是 agent-native server maintainability 需要更严格或不同约束的地方，我们定义自己的标准。

项目应兼容 `AGENTS.md` 风格的指令，并可以生成 `AGENTS.md` 文件，但宪法层级高于 `AGENTS.md`。

项目可以借鉴 Superpowers 类系统里的 “skill” 概念，但我们的规范单元更广：module contracts、change specs、architecture manifests、generated scaffolds、verification gates 和 agent operating manuals。

## 3. 宪法原则

### 3.1 架构必须显式

关键架构规则不能只存在于维护者脑子里。

每个模块、边界、依赖、契约、事件、权限、错误码、测试要求和迁移要求，都应该以代码、schema、manifest 或生成文档的形式存在。

### 3.2 Agent 上下文是一等资源

框架必须降低 agent 正确修改代码前需要读取的上下文量。

一个好的模块应该告诉 agent：

- 它拥有什么
- 它不拥有什么
- 它暴露哪些 APIs
- 它发布哪些 events
- 它消费哪些 events
- 它允许哪些依赖
- 哪些 invariants 不能破坏
- 哪些 tests 验证该模块
- 哪些文件是 extension points
- 哪些文件是 generated files，不应手工编辑

### 3.3 禁止无边界修改

Agent 很少应该执行不受约束的全仓库修改。

每个有意义的变更都应表达为一个有边界的 change request，其中包含 affected modules、expected contracts、generated files、tests、migrations 和 verification steps。

### 3.4 Schema Before Code

公共行为必须先定义，再实现。

推荐顺序：

1. Requirement
2. Spec
3. Contract
4. Generated structure
5. Business logic
6. Tests
7. Verification
8. Documentation update

Agent 不应在没有更新相关 schema 或 manifest 的情况下，临时发明 API、event payload、database shape 或 error format。

### 3.5 生成结构，手写逻辑

框架应该生成可重复结构。

Agents 和人类应该专注于 behavior、invariants、edge cases 和 tests，而不是猜文件该放哪、样板代码该怎么写。

### 3.6 强模块边界

模块是 bounded contexts。

模块可以拥有 data、commands、queries、events、policies、repositories 和 tests。其他模块不能直接访问它的内部实现。

跨模块通信必须通过批准过的接口：

- Commands
- Queries
- Events
- Public module APIs
- Generated clients

### 3.7 Commands 表达意图，Events 表达事实

Commands 描述被请求的变更。

Events 描述已经发生的事实。

这个区分是强制性的，因为它帮助 agents 理解业务流，并让行为更容易测试、回放和验证。

### 3.8 服务器保持权威

Agent 写出的代码不能削弱 server authority。

所有会改变状态的动作都必须通过 validation、permission checks、invariants 和 transactional consistency rules。

这条原则同时适用于人类用户和 AI agents。

### 3.9 架构必须可测试

架构规则必须能自动检查。

如果一条规则重要，框架最终就应该提供命令验证它。

例子：

- 模块不得 import forbidden modules。
- Events 必须 versioned。
- Commands 必须有 schemas。
- Public APIs 必须有 contract tests。
- Migrations 必须有 rollback 或 compatibility tests。
- Permissioned operations 必须声明 permissions。
- Generated files 不应手工编辑。

### 3.10 文档必须可操作

文档应该帮助 agent 或人类采取下一个正确动作。

项目应避免装饰性文档，也就是看起来完整但不能指导 implementation、testing 或 verification 的文档。

每个模块最终都应该有 agent-readable operating guide。

## 4. 必需项目产物

项目应逐步演进出以下 artifact system。

### 4.1 Constitution

文件：

```text
CONSTITUTION.md
```

目的：

- 定义项目原则
- 定义什么是有效架构
- 定义标准如何采纳或修订
- 提供最高层决策框架

### 4.2 Repository Agent Guide

预期文件：

```text
AGENTS.md
```

目的：

- 给 coding agents 简明操作指令
- 列出 build、test、lint、generation、verification commands
- 定义 "always"、"ask first"、"never" 规则
- 指向 module-level guides 和 architecture manifests

`AGENTS.md` 应根据本宪法生成或审查。

### 4.3 Architecture Manifests

预期目录：

```text
.arch/
```

可能包含：

```text
.arch/modules.yaml
.arch/dependencies.yaml
.arch/conventions.yaml
.arch/commands.yaml
.arch/events.yaml
.arch/errors.yaml
.arch/permissions.yaml
.arch/test-matrix.yaml
.arch/generation.yaml
```

目的：

- 给 agents 一张机器可读的系统地图
- 支持 architecture checks
- 支持 code generation
- 支持 impact analysis

### 4.4 Module Manifests

每个模块预期文件：

```text
modules/<module>/module.yaml
```

最低字段示例：

```yaml
module: inventory
type: domain
owns:
  entities:
    - inventory
    - inventory_item
public_api:
  commands:
    - AddItem
    - RemoveItem
  queries:
    - GetInventory
events:
  publishes:
    - ItemAdded
    - ItemRemoved
  subscribes:
    - PlayerCreated
dependencies:
  allowed:
    - player
    - economy
  forbidden:
    - matchmaking
invariants:
  - item_count_must_not_be_negative
  - inventory_capacity_must_not_exceed_limit
tests:
  required:
    - command_tests
    - event_tests
    - invariant_tests
    - migration_tests
```

### 4.5 Module Agent Guides

每个模块预期文件：

```text
modules/<module>/AGENTS.md
```

目的：

- 说明什么时候使用该模块
- 说明什么时候不要使用该模块
- 列出 extension points
- 列出 forbidden shortcuts
- 列出 required tests
- 链接 schemas、commands、events 和 invariants

### 4.6 Change Specs

预期目录：

```text
changes/
```

每个非平凡变更应有独立目录：

```text
changes/<date>-<change-id>/
  request.md
  spec.yaml
  impact.md
  plan.md
  checklist.md
  verification.md
```

目的：

- 给 agents 持久上下文
- 让决策可追踪
- 支持 agents 和 humans 之间交接
- 让后续维护更容易

示例 `spec.yaml`：

```yaml
change_id: add-season-pass
type: feature
affected_modules:
  - player
  - economy
  - reward
new_commands:
  - PurchaseSeasonPass
  - ClaimSeasonPassReward
new_events:
  - SeasonPassPurchased
  - SeasonPassRewardClaimed
data_migrations:
  required: true
compatibility:
  breaking_api: false
  breaking_db: false
acceptance_tests:
  - player_can_purchase_pass
  - player_cannot_claim_locked_reward
  - duplicate_claim_is_rejected
```

## 5. 标准服务器形态

框架应逐步收敛到可预测结构。

示例结构：

```text
modules/
  player/
    module.yaml
    AGENTS.md
    commands/
    queries/
    events/
    models/
    repositories/
    policies/
    tests/
  inventory/
  economy/
  match/
  chat/
  notification/
schema/
  api/
  commands/
  events/
  errors/
  permissions/
  db/
.arch/
changes/
tools/
docs/
```

具体实现可以变化，但框架应保持可预测性。

## 6. 标准变更工作流

每个非平凡 feature、bug fix、migration 或 refactor 都应遵循此流程。

### 6.1 Intake

把用户请求转成清晰 requirement。

Requirement 应识别：

- Desired behavior
- User-visible outcome
- Affected modules
- Unknowns
- Risks
- Acceptance criteria

### 6.2 Impact Analysis

改代码前，识别：

- Modules affected
- Contracts affected
- Data affected
- Events affected
- Permissions affected
- Tests affected
- Migration requirements
- Compatibility risks

### 6.3 Plan

计划应短、边界清晰、可执行。

它应识别：

- 要创建或编辑的文件
- Generated artifacts
- Manual logic
- Required tests
- Verification commands

### 6.4 Contract Update

当行为改变公共或跨模块契约时，应先更新 schemas 和 manifests，再实现。

### 6.5 Code Generation

凡可用 generator 生成的重复结构，都应使用 generator。

手动创建 framework-shaped files 只能作为临时 fallback。

### 6.6 Implementation

只在定义好的 change boundary 内实现。

不要绕过 module public APIs、validators、repositories、migrations 或 permission systems。

### 6.7 Verification

变更只有在运行验证，或明确记录验证未运行时，才算完成。

验证应包括最小相关 tests 和 architecture checks。

### 6.8 Documentation Update

如果变更修改了 architecture、module ownership、public behavior、commands、events、schemas、permissions 或 test procedure，必须更新相关 docs 或 manifests。

## 7. 模块标准

每个模块应该有稳定内部形态。

推荐模块目录：

```text
commands/
queries/
events/
models/
repositories/
policies/
services/
tests/
```

定义：

- `commands/`：state-changing requests and handlers
- `queries/`：read-only access patterns
- `events/`：facts emitted by the module
- `models/`：domain models and value objects
- `repositories/`：persistence boundaries
- `policies/`：event reactions and cross-module orchestration
- `services/`：local domain services，不是 generic dumping grounds
- `tests/`：module-local tests

规则：

- 模块拥有自己的内部数据。
- 其他模块不能直接 import 内部文件。
- 公共行为必须通过声明过的 APIs 暴露。
- Events 一旦公开就必须 versioned。
- Commands 必须校验 input 和 permissions。
- Invariants 必须有测试。

## 8. 契约标准

框架应该把 contracts 作为一等对象。

Contract 类型：

- API contracts
- Command contracts
- Query contracts
- Event contracts
- Error contracts
- Permission contracts
- Database migration contracts
- Generated client contracts

规则：

- Public contracts 必须在 schema 中声明。
- 当兼容性重要时，contracts 必须 versioned。
- Breaking changes 必须显式声明。
- Generated code 必须能追踪到 source schema。
- Agents 不得手工编辑 generated contract output，除非目标就是 generator 本身。

## 9. 测试标准

测试不是流程装饰。测试是 agent-driven development 的架构护栏。

必需测试类别应包括：

- Unit tests for handlers, policies, validators, and domain logic
- Contract tests for public APIs, commands, queries, and events
- Invariant tests for business rules
- Integration tests for module collaboration
- Migration tests for database changes
- Replay tests where event history is used
- Architecture tests for dependency and manifest rules

每个模块应在 `module.yaml` 中声明 required test categories。

## 10. 架构验证标准

框架最终应提供类似命令：

```bash
server check architecture
server check contracts
server check module <module>
server check change <change-id>
server generate module <module>
server generate command <module> <command>
server generate event <module> <event>
```

这些命令应同时服务 humans 和 agents。

命令输出应简洁、确定、可行动。

## 11. Agent 操作规则

在本仓库工作的 agents 应遵循这些规则。

### 11.1 Always

- 修改模块前，读取相关 constitution、architecture manifest、module manifest 和 module `AGENTS.md`。
- 在实现前优先更新 schema 和 manifest。
- 把编辑限制在声明过的 change boundary 内。
- 行为变更必须新增或更新 tests。
- 可用时运行相关 verification commands。
- 记录 verification results。

### 11.2 Ask First

以下情况先询问：

- 修改宪法原则
- 重新定义 module ownership
- 引入新架构模式
- 做 breaking API 或 event changes
- 修改 generated file conventions
- 删除 tests 或削弱 validation
- 在模块间迁移 data ownership

### 11.3 Never

禁止：

- 为方便绕过 module boundaries
- 把业务逻辑藏在 transport handlers
- 添加未登记 public events
- 添加未登记 permissions
- 添加无类型 cross-module payloads
- 无理由手工编辑 generated files
- 没有运行验证却声称已验证

## 12. AI 功能边界

框架可以支持 AI gameplay 和 AI product features，但不能把这些和 agent-native maintainability 混为一谈。

AI gameplay features 可以包括：

- NPC agents
- Agent memory
- Model routing
- Tool calling
- Agent simulation
- Narrative systems

这些是 application-layer capabilities。

宪法基础仍然是：

- Explicit architecture
- Machine-readable standards
- Bounded changes
- Schema-first contracts
- Generated structure
- Automated verification
- Agent-readable context

## 13. 治理

宪法可以修改，但必须谨慎。

任何 constitutional amendment 都应包含：

- Motivation
- Problem being solved
- Alternatives considered
- Compatibility impact
- Required migration
- Updated examples

宪法修改应比实现修改少得多。

## 14. 决策规则

在两个设计之间选择时，优先选择：

1. 给 agents 更少歧义上下文的设计
2. 创建更强 module boundaries 的设计
3. 让 behavior 更容易验证的设计
4. 让 contracts 更显式的设计
5. 减少 hidden coupling 的设计
6. 支持 code generation 的设计
7. 支持 long-term maintenance 的设计
8. 同时对人类开发者保持实用性的设计

Agent-native 不等于 agent-only。

最好的架构应该让 humans 和 agents 都更有效。

## 15. 开放问题

以下问题有意保持开放：

- 参考实现使用哪种语言？
- 第一版实现应游戏专用还是后端通用？
- 框架是否应拥有 transport，还是集成现有 transports？
- Event sourcing 应强制还是可选？
- Module manifests 应使用 YAML、JSON、TOML 还是 code-native？
- Generators 是否应从第一天就内置在 framework CLI？
- 哪个数据库应作为第一优先参考目标？
- 多少规则靠 convention，多少规则靠 static analysis？
- 大型变更中如何协调 multi-agent collaboration？
- 框架如何衡量它是否真的让 agents 更容易维护？

这些问题应通过 design notes、prototypes 和具体实现经验回答。

## 16. 第一版实现目标

第一个严肃 prototype 应证明这个命题：

> 给定一个新的后端需求，AI coding agent 能识别 affected module，更新正确 contracts，生成正确结构，实现行为，添加测试，运行验证，并更新文档，同时不破坏无关架构。

Prototype 应小而完整。

建议示例领域：

- Player accounts
- Inventory
- Currency
- Rewards
- Tasks or quests
- Match sessions

Demo 应强调 maintainability 和 agent workflow，而不是功能数量。

## 17. 创始声明

这个项目不是试图让 agents 魔法般变聪明。

它试图让服务器架构更清晰、更有边界、更可验证、更可演进，从而让 agents 能稳定产出高质量工作。

标准很简单：

> 如果 agent 无法理解一个变更应该放在哪里、影响哪些 contracts、由哪些 tests 证明、受哪些 architecture rules 约束，那么框架还没有完成它的工作。
