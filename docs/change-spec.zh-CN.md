# Change Spec Standard 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：`changes/<date>-<change-id>/`  
说明：本文件是 `docs/change-spec.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档定义 vibit 的标准 change-spec 工作流。

Change spec 是 humans 和 agents 的持久工作上下文。它的存在，是为了避免非平凡变更只依赖一小段 prompt 执行。

## 1. 目的

每个非平凡 feature、bug fix、migration、refactor 或 standard change 都应该有目录：

```text
changes/<date>-<change-id>/
```

示例：

```text
changes/2026-05-12-add-inventory-module/
```

该目录应保存：

- 原始请求
- 澄清后的 requirement
- Affected modules 和 contracts
- Impact analysis
- Implementation plan
- Acceptance checklist
- Verification results
- Open questions

## 2. 必需文件

推荐结构：

```text
changes/<date>-<change-id>/
  request.md
  spec.yaml
  impact.md
  plan.md
  checklist.md
  verification.md
```

对于小型纯文档变更，可以使用更轻量版本，但 agents 仍必须记录改了什么以及如何验证。

## 3. `request.md`

`request.md` 记录 human request 和澄清后的 requirement。

它应包含：

- Original request
- Clarified requirement
- User-visible outcome
- Non-goals
- Unknowns
- Acceptance criteria

## 4. `spec.yaml`

`spec.yaml` 是该变更的机器可读摘要。

示例：

```yaml
schema_version: 0.1
change_id: add-inventory-module
date: 2026-05-12
type: feature
status: draft

summary: Add the first inventory module prototype.

affected_modules:
  - inventory

contracts:
  commands:
    added:
      - AddItem
      - RemoveItem
    changed: []
    removed: []
  queries:
    added:
      - GetInventory
    changed: []
    removed: []
  events:
    added:
      - ItemAdded
      - ItemRemoved
    changed: []
    removed: []
  permissions:
    added:
      - inventory_read
      - inventory_write
    changed: []
    removed: []

data:
  migrations_required: false
  ownership_changes: []

compatibility:
  breaking_api: false
  breaking_events: false
  breaking_data: false

verification:
  required:
    - architecture_checks
    - command_tests
    - invariant_tests
  status: Not applicable
```

## 5. `impact.md`

`impact.md` 解释变更触及什么，以及为什么。

它应覆盖：

- Affected modules
- Module ownership impact
- Public contract impact
- Event impact
- Permission impact
- Data and migration impact
- Test impact
- Documentation impact
- Compatibility risks

## 6. `plan.md`

`plan.md` 是有边界的 implementation plan。

它应识别：

- Files to create
- Files to edit
- Generated artifacts
- Handwritten logic
- Tests to add or update
- Verification commands
- 相关时的 rollback 或 migration notes

## 7. `checklist.md`

`checklist.md` 跟踪完成状态。

使用简单 task states：

```text
- [ ] Pending task
- [x] Completed task
```

Checklist 应包含 contract、implementation、test、verification 和 documentation tasks。

## 8. `verification.md`

`verification.md` 记录检查过什么。

使用此格式：

```text
Verified:
- <command or check>

Not verified:
- <reason>

Not applicable:
- <reason>
```

没有运行验证时，绝不能声称变更已经验证。

## 9. Agent 规则

以下情况 agents 应创建或更新 change spec：

- 变更影响 public behavior。
- 变更影响 module ownership。
- 变更新增或修改 commands、queries、events、permissions、errors 或 data shape。
- 变更引入或修改 architectural standard。
- 变更足够大，另一个 agent 继续工作时需要持久上下文。

对于小 typo、纯 formatting edits 或窄范围文档更新，agents 可以跳过完整 change spec，但最终响应仍必须说明 verification status。

## 10. 命名规则

Change directory names 应使用：

```text
YYYY-MM-DD-short-kebab-case-id
```

规则：

- 使用 change spec 创建日期。
- 使用简短、稳定、描述性的 ID。
- 实现开始后不要重命名，除非必要。
- 让 ID 对未来 agents 有意义。
