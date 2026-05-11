# Schema Validation Standard 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：`schema/`  
说明：本文件是 `docs/schema-validation.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档定义 vibit 的第一版 schema validation layer。

Schema validation 把 standards 从可读文档推进为机器可检查契约。

## 1. 目的

vibit 使用 manifests、change specs、decision records 和 tool JSON output 作为 agent-readable architecture context。

Agents 不应只依赖 prose 或肉眼检查。重要结构应由工具验证。

## 2. 位置

Schema files 位于：

```text
schema/
```

初始文件：

```text
schema/module-manifest.schema.json
schema/change-spec.schema.json
schema/agent-decision-record.schema.json
schema/inspect-output.schema.json
```

## 3. 验证策略

第一版实现刻意保持轻量。

当前规则：

- Schemas 使用 JSON Schema draft 2020-12。
- 验证 schema files 存在且能作为 JSON 解析。
- 在无外部依赖的情况下验证 manifests 和 specs 的 critical fields。
- 在 standards 稳定前保持小范围严格验证。

未来规则：

- 当 dependency strategy 确定后，使用完整 JSON Schema validator。
- 从 schemas 生成 examples。
- 在可行时，把 schemas 作为 generator inputs。

## 4. CLI Commands

初始命令：

```bash
node tools/vibit check schemas
```

聚合命令：

```bash
node tools/vibit check all
```

`check all` 应包含 schema checks。

## 5. 当前检查什么

初始检查应覆盖：

- Required schema files 存在。
- Schema files 可作为 JSON 解析。
- `modules/<module>/module.yaml` 声明 critical fields。
- Change specs 声明 critical fields 和允许的 verification status。
- Agent Decision Records 包含 required sections。
- Tool JSON output schemas 覆盖 inspect output 和 check result output。
- Architecture conventions 声明 schema artifacts。

这还不是完整 YAML schema validation。

## 6. Agent 规则

当改变以下结构时，agents 必须更新 schemas：

- Module manifests
- Change specs
- Agent Decision Records
- Tool JSON output，包括 inspect output 和 check result output

如果 CLI 尚不能完整验证某个结构，agents 必须在相关 change spec 中记录缺口。

## 7. 开放问题

- 未来应使用哪个完整 JSON Schema validator？
- YAML manifests 是否应先转换为 JSON 再验证？
- Schemas 是否应直接驱动 code generation？
- Tool JSON output 是否应从一开始就拥有稳定 versioned schemas？
