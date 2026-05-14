# Agent Tooling Standard 中文版

状态：Draft v0.1
最后更新：2026-05-14
范围：面向 Agent 的 inspection、generation 和 verification commands
说明：本文件是 `docs/agent-tooling.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本标准定义最小的机器可读 tooling surface，让 agents 可以在不依赖记忆或大范围读源码的情况下维护 vibit。

## 1. 目的

vibit 需要能用结构化输出回答窄问题的工具。

Agents 在改代码前，应能检查下一个 work item、已注册 contracts、generated outputs、reference planning context 和 verification rules。这样可以减少架构漂移，让“继续推进”变得可预测。

## 2. Required Commands

当前面向 agent 的 commands 是：

```bash
node tools/vibit inspect work --json
node tools/vibit inspect next --json
node tools/vibit inspect contracts --json
node tools/vibit inspect contracts --module inventory --json
node tools/vibit inspect contracts --type command --json
node tools/vibit inspect contracts --status draft --json
node tools/vibit inspect contract --module inventory --type command --id GrantItem
node tools/vibit inspect generated --json
node tools/vibit inspect generated --module inventory --json
node tools/vibit inspect generated --type command --json
node tools/vibit inspect reference --json
node tools/vibit check agent-tooling --json
node tools/vibit generate contract-shapes all
```

这些 commands 故意保持小而明确。它们不能替代在重大修改前阅读 governing standards。

## 3. Inspection Rules

Inspection commands 应该：

- 在提供 `--json` 时返回 JSON。
- 包含 `schema_version`、`kind` 和 `project`。
- 为每个持久事实报告 source artifact path。
- JSON fields 中的 repository-relative paths，例如 `artifact`、`path`、`source` 和 `output`，即使在 Windows 上也必须统一为 forward-slash form。
- 优先使用明确的 `exists`、`status` 和 `summary` fields，而不是只给 prose。
- 不隐藏 ask-first boundaries。

## 4. Generation Rules

Generation commands 应该：

- 读取 source contracts、manifests 或 schemas。
- 只写入已声明的 generated roots。
- 包含 generated markers、source traces 和 generator traces。
- 足够可复现，使 checks 能发现 drift。
- 除非由单独 work item 负责实现，否则不添加 runtime behavior。

`node tools/vibit generate contract-shapes all` 会从 semantic contract manifests 生成可检查的 Go contract shape files。这些 files 只是 generated artifacts，不实现 handlers、persistence、authentication、routes 或 protocol envelope behavior。

## 5. Verification

当前 verification：

```bash
node tools/vibit check agent-tooling --json
node tools/vibit check generated --json
node tools/vibit check all --json
```

`check agent-tooling` 会验证本标准和简体中文译本存在，验证 command surface 已被文档化，并验证 public docs 保持 translation pairs。

## 6. Agent Rules

Agents 必须：

- 当 work state 不清楚时，在解释 continuation request 前运行 `node tools/vibit inspect next --json`。
- 在广泛 contract 或 generator 工作前运行 `node tools/vibit inspect contracts --json`。
- 当只需要窄范围 contract slice 时，使用 `node tools/vibit inspect contracts --module <module> --json`、`node tools/vibit inspect contracts --type <contract-type> --json` 或 `node tools/vibit inspect contracts --status <status> --json`。
- 在修改 generated output 或 generator behavior 前运行 `node tools/vibit inspect generated --json`。
- 当只需要窄范围 generated-output slice 时，使用 `node tools/vibit inspect generated --module <module> --json` 或 `node tools/vibit inspect generated --type <contract-type> --json`。
- 在规划新的 game server capability family 前运行 `node tools/vibit inspect reference --json`。
- 通过 `node tools/vibit generate contract-shapes all` 重新生成 generated contract shapes，而不是手改这些文件。

Agents 不得：

- 把 generated contract shapes 当作 handwritten runtime implementation。
- 把 inspection output 当作越过 ask-first boundaries 的许可。
- 当 check 或 JSON inspection 能执行规则时，用宽泛的纯 prose standard 代替窄 tooling。
