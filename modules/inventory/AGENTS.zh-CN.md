# inventory Module Agent Guide 中文版

状态：Draft v0.1
说明：本文件是 `modules/inventory/AGENTS.md` 的简体中文译本。英文版本是权威版本。

## 何时使用本模块

当需求需要修改或读取 player inventory state 时，使用本模块。

第一条 proof-slice capability 有意保持很小：

- 通过 `GrantItem` 给某个 player 的 inventory 发放 item。
- 通过 `GetInventory` 读取某个 player 的 inventory。
- 成功发放后发布 `ItemGranted`。
- 执行 inventory permissions、正数量、容量限制和 inventory item identity 相关规则。

本模块拥有 inventory records 和 inventory items。它可以引用 `player_id` 和 `item_id`，但不拥有 player accounts 或 item catalog。

## 何时不要使用本模块

不要把以下需求放入本模块：

- Player account lifecycle。
- Item catalog definitions 或 item balancing。
- Currency balances 或 purchases。
- Reward claim eligibility。
- Quest progress。
- Match 或 session lifecycle。

如果需求需要这些概念，应创建或更新对应 owning module 的 contract，而不是在本模块里隐藏 ownership。

## Extension Points

- Command handler：`GrantItem`
- Query handler：`GetInventory`
- Published event：`ItemGranted`
- Policies：inventory capacity 和 inventory permission checks
- Repository：位于 module boundary 后面的 inventory persistence
- Tests：command、query、event、contract、invariant 和 architecture tests

Runtime contract source files 现在位于 `contracts/inventory/` 下，并登记在 `.arch/contracts.yaml` 中。Generated files 还没有创建。

在实现第一条 runtime slice 前，应阅读：

- `contracts/inventory/commands/GrantItem.yaml`
- `contracts/inventory/queries/GetInventory.yaml`
- `contracts/inventory/events/ItemGranted.yaml`
- `contracts/inventory/errors/inventory_errors.yaml`
- `contracts/inventory/permissions/inventory_permissions.yaml`

## Forbidden Shortcuts

- 不要绕过 `module.yaml` 中声明的边界。
- 不要直接修改其他模块拥有的数据。
- 不要添加未登记的 public commands、queries、events 或 permissions。
- 不要把 inventory business rules 放进 HTTP 或 transport handlers。
- 不要手工编辑 generated files。如果 generated output 错了，应修改 source contract、template 或 generator。
- 不要在实现中临时发明 payload fields。必须先更新对应 contract source file。
- 未先更新 manifest 和 change spec，不要引入对 player、currency、reward、quest 或 match modules 的 dependency。
- Grant flows 中不要使用负数或零数量。

## Required Tests

参见 `module.yaml` 中的 `tests.required`。

第一条 runtime slice 的测试应覆盖：

- `GrantItem` 接受有效 grant 并记录 item。
- `GrantItem` 拒绝 invalid quantity。
- `GrantItem` 拒绝超过 capacity 的 grant。
- 成功 grant 时，`GrantItem` 只发布一次 `ItemGranted`。
- `GetInventory` 返回 inventory state 且不修改 state。
- Permission failure 使用 `INVENTORY_PERMISSION_DENIED`。
- Architecture checks 仍然通过。

如果 runtime test infrastructure 尚不存在，应把相关 tests 记录为 not available，而不是从 manifest 中移除。
