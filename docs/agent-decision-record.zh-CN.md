# Agent Decision Record Standard 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：`decisions/`  
说明：本文件是 `docs/agent-decision-record.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档定义 vibit 的 Agent Decision Records。

Agent Decision Records 是面向 humans 和 agents 的 ADR-style 记录，用于保存影响 architecture、standards、modules、generated files 或 long-term maintainability 的决策。

## 1. 目的

Conversation logs 保存讨论如何展开。

Change specs 保存某个变更如何执行。

Agent Decision Records 保存未来 agents 必须尊重的持久决策依据。

它们存在的原因是：agent 写出的代码可能局部看起来正确，却违背长期设计意图。

## 2. 位置

Decision records 位于：

```text
decisions/
```

模板：

```text
decisions/_template/adr-agent.md
```

Decision files 应使用：

```text
decisions/ADR-0001-short-kebab-case-title.md
```

## 3. 记录什么

记录影响以下内容的决策：

- Constitutional principles
- Architecture standards
- Module ownership
- Public commands、queries、events、errors 或 permissions
- Generated file conventions
- Verification gates
- Implementation language
- Server instance architecture
- Long-term maintainability tradeoffs

不要为每个很小的实现细节创建 decision record。

## 4. 公开依据，不是隐藏 Chain-of-Thought

Decision records 必须包含 public rationale。

它们不应包含 private chain-of-thought、隐藏推理转储或不可验证的内部独白。

标准格式是：

- Context
- Decision
- Alternatives considered
- Rationale
- Agent reasoning summary
- Confidence
- Consequences
- Links

“Agent reasoning summary” 应简洁且可检查。它应解释为什么该决策合适，但不暴露隐藏推理。

## 5. 从其他产物链接

当某个 decision 影响 module 时，module 应引用 decision ID。

示例：

```yaml
decisions:
  - ADR-0001
```

当某个 decision 影响 change 时，change spec 应引用它。

示例：

```yaml
decisions:
  - ADR-0001
```

Decision record 才是持久 rationale 的位置。`module.yaml` 应包含链接和简洁 metadata，而不是长篇 reasoning text。

## 6. Confidence And Weight

Decision records 可以包含 confidence level：

```text
low
medium
high
```

也可以包含 decision weights：

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: medium
  implementation_cost: low
  reversibility: medium
  long_term_maintainability: high
```

这些 weights 不是数学真理。它们是紧凑的 public metadata，帮助后续 agents 理解决策时什么更重要。

## 7. Generated Immutability

Generated files 对 non-system agents 不可变。

规则：

- Generated files 必须在 manifests 中声明。
- Ordinary agents 不得手工编辑 generated files。
- 如果 generated output 错了，应修改 source schema、template 或 generator。
- 覆盖 generated file 需要显式 decision record 或 change spec。
- 任何 override 的原因都必须记录。

初始 permission concept：

```text
generated_file_override
```

在框架拥有真正 permission system 前，这是 standards-level permission。

## 8. 必需章节

每个 decision record 应包含：

```text
# ADR-0001: Title

Status:
Date:
Decision Makers:
Related changes:
Related conversations:
Related artifacts:

## Context

## Decision

## Alternatives Considered

## Rationale

## Agent Reasoning Summary

## Decision Weights

## Consequences

## Reversal Conditions

## Follow-Up
```

## 9. Agent 规则

以下情况 agents 应创建或更新 decision record：

- 决策会塑造未来 architecture。
- 决策拒绝了一个合理替代方案。
- 决策影响 generated files 或 module boundaries。
- 未来 agent 很难仅从代码推断该决策。

Agents 不得把 decision records 当作倾倒冗长 private reasoning 的地方。

Decision records 应保持简洁、公开、可链接。
