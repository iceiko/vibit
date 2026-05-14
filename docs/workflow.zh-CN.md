# Work Continuation Workflow 中文版

状态：Draft v0.1
最后更新：2026-05-13
范围：maintainer-agent continuation、work item sequencing 和 roadmap state
说明：本文件是 `docs/workflow.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本标准定义 vibit agents 如何判断“继续”是什么意思。

## 1. 目的

项目需要确定性的推进方式。

当 maintainer 要求 agent 继续时，agent 不应该只凭记忆临场选择下一个任务。它应该读取项目工作队列，找到下一个 ready 的有边界 work item，执行该 item，验证它，并更新队列。

这个机制让项目可以自举推进，同时避免流程变重。

## 2. Artifact Roles

项目用不同编号表达不同含义：

- `M-000`：milestone identifier。
- `W-0000`：work item identifier。
- `ADR-0000`：architecture decision identifier。
- `changes/YYYY-MM-DD-change-id`：具体 change spec 和 execution record。
- Git commit hash：不可变 repository snapshot。
- Version tag：未来 release identifier。

不要把 ADR 编号当成工作步骤编号。ADR 记录长期决策。Work items 记录执行步骤。

不要把 release versions 当成工作步骤编号。Version 描述工作完成后用户可以依赖的能力集合。

## 3. Work Item Definition

一个 work item 是一次 continuation step 的单位。

好的 work item 应该：

- 有稳定的 `W-0000` id。
- 只属于一个 milestone。
- 有明确 status。
- 在顺序重要时声明 dependencies。
- 有 completion criteria。
- 未实现时声明 expected change spec。
- 完成后链接 change spec 和 commits。
- 在可能分支成更大决策时列出 ask-first boundaries。

## 4. Milestone Definition

Milestone 把相关 work items 组织到一个阶段目标下。

Milestone 不是 release。Milestone 说明 repository 接下来要证明什么能力。Release 说明用户可以依赖什么能力。

## 5. Status Values

Milestone status values：

```text
planned
active
completed
paused
superseded
```

Work item status values：

```text
planned
next_ready
active
blocked
completed
paused
superseded
```

规则：

- 至少应该有一个 milestone 是 `active`。
- `next_ready` work item 的所有 dependencies 必须已经 completed。
- completed work item 应链接 change spec，或解释为什么不需要 change spec。
- 通常每个 active milestone 应只有一个 `next_ready` item。只有在多个独立工作可以并行推进时，才允许多个 next-ready items。
- 如果当前唯一工作被明确标为 `blocked` 并写明 `block_reason`，active milestone 没有 `next_ready` item 是有效状态。这代表有意的 maintainer decision gate，不是空队列。

## 6. Continuation Semantics

Maintainer phrase：

```text
continue
继续
```

含义：

```text
advance one next_ready work item
```

Maintainer phrase：

```text
continue N steps
继续推进 N 步
```

含义：

```text
advance up to N next_ready work items in dependency order
```

Agent 必须提前停止，如果：

- 没有 next-ready work item。
- 下一个 work item 被 blocked。
- 触发 ask-first boundary。
- Verification 失败。
- Maintainer 改变方向。

## 7. Required Agent Intake

在解释 continuation request 前，agent 应运行或查看：

```bash
node tools/vibit inspect work
node tools/vibit check work
```

然后 agent 应阅读相关 work item、dependencies、related decisions 和 expected change spec。

## 8. Execution Rule

单个 work item 的推荐执行流程是：

```text
work item
-> change spec
-> implementation or documentation
-> tests/checks
-> verification record
-> work item status update
-> commit
```

如果 work item 会改变 architecture、public contracts、generated-file conventions、module ownership、dependency adoption 或其他 ask-first boundaries，agent 必须先询问 maintainer。

## 9. Multi-Step Execution

当 maintainer 要求继续多步时，agent 仍然应在开始下一步前，把每一步作为完整 bounded change 完成。

每一步应该有：

- 非平凡工作对应的 change spec。
- Verification。
- 更新后的 work-item state。
- 当 repository state 适合保存时创建 commit。

如果下一步尚未 ready 或需要 maintainer confirmation，agent 可以在达到请求步数前停止。

## 10. Verification

当前命令：

```bash
node tools/vibit check work
node tools/vibit inspect work
node tools/vibit check all
```

`check work` 验证 `.arch/work-items.yaml` 和本 workflow standard 是否存在并基本一致。

未来 checks 可以验证 dependency ordering、单步 state transitions，以及 completed work items 是否链接到存在的 commits。
