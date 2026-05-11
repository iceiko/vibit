# inventory Module Agent Guide 中文版

状态：Draft v0.1
说明：本文件是 `modules/inventory/AGENTS.md` 的简体中文译本。英文版本是权威版本。

## 何时使用本模块

填写本模块负责的需求范围。

## 何时不要使用本模块

填写不应放入本模块的需求。

## Extension Points

- Commands
- Queries
- Events
- Policies
- Tests

## Forbidden Shortcuts

- 不要绕过 `module.yaml` 中声明的边界。
- 不要直接修改其他模块拥有的数据。
- 不要添加未登记的 public commands、queries、events 或 permissions。

## Required Tests

参见 `module.yaml` 中的 `tests.required`。
