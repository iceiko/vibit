# Agent Decision Records 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：vibit 的持久决策依据  
说明：本文件是 `decisions/README.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

这个目录保存 Agent Decision Records。

权威标准见 `docs/agent-decision-record.md`。

规则：

- 记录持久决策，而不是每个小实现细节。
- Rationale 应公开、简洁、可检查。
- 不存储 private chain-of-thought。
- 相关时，从 modules、change specs 和 conversation logs 链接 decisions。
- 用 decision records 解释 generated-file overrides 和长期架构选择。

模板：

```text
decisions/_template/adr-agent.md
```
