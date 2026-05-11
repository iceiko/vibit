# Conversation Log Standard 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：`conversations/`  
说明：本文件是 `docs/conversation-log.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档定义 vibit 如何把 maintainer-agent conversations 记录为持久项目记忆。

Conversation logs 不取代 constitution、change specs、architecture manifests 或 issue trackers。它们解释这些产物背后的推理路径。

## 1. 目的

vibit 正在通过 maintainer-agent 的密切协作被设计出来。

Maintainer 的叙述是 product intent 的第一手来源。它经常包含：

- 问题在形成正式语言前的真实感受
- 对误导性解释的拒绝
- 命名意图
- 架构品味
- 产品边界
- 标准背后的历史推理

未来 agents 不应该只能从最终文档里反推这些历史。

## 2. 位置

Conversation logs 位于：

```text
conversations/
```

可复用模板：

```text
conversations/_template/session.md
```

Session logs 应使用：

```text
conversations/YYYY-MM-DD-short-session-id.md
```

## 3. 权威语言规则

项目文档语言是英文，但 conversation logs 有一条特殊规则：

- Maintainer statements 应尽可能保留原始语言。
- 应添加英文摘要，方便全球可读。
- Agent responses 可以用英文摘要记录。
- 中文 maintainer statements 不需要逐行翻译，但关键决策和意图应有英文摘要。

这个例外存在，是因为 maintainer 原话本身就是项目记忆的一部分。

## 4. 记录什么

记录：

- Date
- Participants
- Context
- Maintainer narrative，高保真保留
- Agent response summary
- Decisions made
- Artifacts created or changed
- Open questions
- Follow-up actions

当 maintainer 的原话包含 product intent、architectural judgment 或 naming rationale 时，优先保留原文。

## 5. 什么可以摘要

Agent responses 可以摘要，除非某段原话对决策有必要。

摘要应捕捉：

- Recommendations
- Alternatives considered
- Warnings or constraints
- Concrete actions taken
- Verification status

不要因为 agent 输出很长就原样保存。

## 6. 不记录什么

绝不要提交：

- Access tokens
- API keys
- Passwords
- Private keys
- Session cookies
- One-time codes
- 与项目历史无关的个人数据
- 项目历史不需要的 private account details

如果 conversation 中出现 secret，用以下内容替换：

```text
[REDACTED_SECRET]
```

如果 private account identifier 不必要，用以下内容替换：

```text
[REDACTED_ACCOUNT]
```

## 7. 与 Change Specs 的关系

Conversation logs 解释项目如何走到某处。

Change specs 解释某个具体变更如何执行。

当 conversation 导致非平凡变更时，应双向链接：

- Conversation log 应提到对应 change spec。
- Change spec 应提到触发它的 conversation。

## 8. Session 必需章节

每个 session log 应包含：

```text
# Conversation: <title>

Date:
Participants:
Related changes:
Related artifacts:

## Context

## Maintainer Narrative

## Agent Response Summary

## Decisions

## Artifacts

## Open Questions

## Follow-Up

## Redaction Notes
```

## 9. Agent 规则

以下情况 agents 应更新 conversation logs：

- Maintainer 引入或澄清 product direction。
- Maintainer 命名或重命名一个概念。
- Maintainer 拒绝一种解释。
- 创建 governance、architecture 或 workflow standard。
- 讨论了有意义的 tradeoff。
- 从 conversation 创建了 public artifact。

对于日常实现工作，logs 可以简短。

## 10. 验证

提交 conversation logs 前，应运行适合仓库的 secret scan。

初始最小检查：

```bash
rg -n "ghp_|github_pat_|TOKEN|Token:|api[_-]?key|password|secret" .
```

这个检查不是完整安全工具，但在专用工具出现前，可以防止最明显的错误。
