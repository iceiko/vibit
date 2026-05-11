# Module Manifest Standard 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：`modules/<module>/module.yaml`  
说明：本文件是 `docs/module-manifest.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档定义 vibit module manifest 的第一版草案。

Module manifest 是一个 module 与系统其他部分之间的本地契约。它告诉 agents 该 module 拥有什么、暴露什么、依赖什么、保护哪些 invariants，以及必须如何验证。

## 1. 目的

每个实现模块都必须有：

```text
modules/<module>/module.yaml
```

Manifest 的存在，是为了让 agents 不需要从分散代码中推断 module boundaries。

它应该回答：

- 这个 module 是什么？
- 它属于哪类 module？
- 它拥有哪些 data 和 concepts？
- 它暴露哪些 public commands 和 queries？
- 它发布和订阅哪些 events？
- 哪些 dependencies 被允许或禁止？
- 哪些 invariants 不能被破坏？
- 哪些文件是 generated？
- 哪些 tests 是必需的？

## 2. 最小示例

```yaml
schema_version: 0.1
module: inventory
category: domain
status: draft
summary: Owns player inventory state and item capacity rules.

owns:
  entities:
    - inventory
    - inventory_item
  data:
    - inventory_records
  permissions:
    - inventory_read
    - inventory_write

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
    - currency
  forbidden:
    - match

invariants:
  - item_count_must_not_be_negative
  - inventory_capacity_must_not_exceed_limit

generated:
  files: []

tests:
  required:
    - command_tests
    - query_tests
    - event_tests
    - invariant_tests
    - migration_tests
```

## 3. 必填字段

### `schema_version`

Manifest schema version。

初始值：

```yaml
schema_version: 0.1
```

### `module`

稳定的 module identifier。

规则：

- 使用 `snake_case`。
- 与 `modules/` 下的目录名一致。
- 未经显式 migration 不要重命名。

### `category`

Module category。

初始允许值：

```text
domain
platform
integration
application
```

定义维护在 `.arch/modules.yaml` 中。

### `status`

生命周期状态。

初始允许值：

```text
draft
active
deprecated
removed
```

### `summary`

一句话解释 module 的责任。

Summary 应帮助 agent 判断某个变更是否属于该 module。

## 4. Ownership

`owns` 声明该 module 控制的 concepts、data 和 permissions。

示例：

```yaml
owns:
  entities:
    - inventory
    - inventory_item
  data:
    - inventory_records
  permissions:
    - inventory_read
    - inventory_write
```

规则：

- Owned data 不得被其他 modules 直接修改。
- 在添加实现代码前，ownership 必须显式。
- 在 modules 之间移动 ownership 需要 change spec 和 maintainer approval。

## 5. Public API

`public_api` 声明 module 面向外部的预期表面。

示例：

```yaml
public_api:
  commands:
    - AddItem
  queries:
    - GetInventory
```

规则：

- Commands 表达改变状态的意图。
- Queries 表达只读访问。
- Public commands 和 queries 必须先有 schemas，再实现。
- Transport handlers 应调用 module APIs，而不是拥有业务逻辑。

## 6. Events

`events` 声明 module 发布的事实，以及从其他 modules 消费的事实。

示例：

```yaml
events:
  publishes:
    - ItemAdded
  subscribes:
    - PlayerCreated
```

规则：

- Events 表达已经发生的事实。
- 当兼容性重要时，public events 必须 versioned。
- Event payloads 必须由 schema 定义。
- Subscriptions 应通过 policies 或 explicit handlers 处理。

## 7. Dependencies

`dependencies` 声明允许和禁止的 module dependencies。

示例：

```yaml
dependencies:
  allowed:
    - player
  forbidden:
    - match
```

规则：

- 未列为 allowed 的 dependency，在声明前应视为 disallowed。
- Forbidden dependencies 记录特别危险或无效的耦合。
- Agents 不得添加违反本节的 imports 或 calls。

## 8. Invariants

`invariants` 声明必须持续为真的业务规则。

示例：

```yaml
invariants:
  - item_count_must_not_be_negative
```

规则：

- 使用 `snake_case`。
- 每个 invariant 最终都应该有 tests。
- Invariants 优先于实现便利。

## 9. Generated Files

`generated` 声明由 generators 产生的文件或目录。

示例：

```yaml
generated:
  files:
    - commands/add_item.generated.ts
  directories:
    - clients
```

规则：

- Agents 不得手工编辑 generated files，除非正在修改 generator 或记录原因。
- Generated files 必须能追踪到 source schema。

## 10. Tests

`tests.required` 声明必需 test categories。

示例：

```yaml
tests:
  required:
    - command_tests
    - query_tests
    - event_tests
    - invariant_tests
    - migration_tests
```

初始允许的 test categories：

```text
unit_tests
command_tests
query_tests
event_tests
contract_tests
invariant_tests
integration_tests
migration_tests
replay_tests
architecture_tests
```

规则：

- Behavior changes 应更新相关 tests。
- 如果必需 test category 还不能运行，应记录为 not available。
- 不要为了让变更更容易而移除 required test categories。

## 11. Agent Checklist

修改 module 前，agent 应检查：

- Module manifest 是否存在？
- 变更是否属于这个 module？
- 变更是否影响 ownership？
- 变更是否新增或修改 public command、query、event、error 或 permission？
- 是否需要先修改 schema？
- Dependencies 是否仍然被允许？
- 哪些 invariants 有风险？
- 哪些 tests 必须新增或更新？
- 哪些 generated files 不应手工编辑？

如果 manifest 不足，应先改进 manifest，再修改实现代码。
