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

`player_id` 是由 `docs/player-identity-session-boundary.md` 和 `ADR-0021` 约束的 external reference。在 session validation 存在之前，请求或 envelope 中携带的 `player_id` 只是 metadata，不是 authenticated proof。Inventory 可以用 `player_id` 作为 inventory state 的 key，但不得创建、认证、验证、迁移或拥有 player accounts 或 runtime sessions。

## 何时不要使用本模块

不要把以下需求放入本模块：

- Player account lifecycle。
- Player identity lifecycle。
- Authentication、token formats、credential storage 或 session validation。
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

Runtime contract source files 位于 `contracts/inventory/` 下，并登记在 `.arch/contracts.yaml` 中。

Generated contract shapes：

- 计划中的 `GrantItem`、`GetInventory` 和 `ItemGranted` Go contract shapes
- 第一批 WebSocket client/server inventory messages 的 generated Protobuf wire schemas

在实现第一条 runtime slice 前，应阅读：

- `contracts/inventory/commands/GrantItem.yaml`
- `contracts/inventory/queries/GetInventory.yaml`
- `contracts/inventory/events/ItemGranted.yaml`
- `contracts/inventory/errors/inventory_errors.yaml`
- `contracts/inventory/permissions/inventory_permissions.yaml`

第一条 handwritten runtime 路径现在从 `runtime/internal/modules/inventory/` 开始：

- `GrantItemHandler`
- `GetInventoryHandler`
- 位于 module boundary 后面的 inventory repository interface
- Inventory capacity 和 permission policy interfaces
- `PermissionContext`，它携带 requested actor text、目标 `player_id` 和 `app.RequestIdentity`
- 面向 `runtime/internal/app.Dispatcher` 的 `RegisterRoutes`
- 覆盖 command、query、event、permission、capacity 和 dispatcher integration behavior 的 tests

Inventory permission policies 可以检查 application-owned request identity context。它们不得验证 tokens、查询 player accounts、拥有 sessions，或 import player module。`metadata_only` identity 不是 authenticated proof。`StaticPermissionPolicy` 只作为显式 bootstrap/local proof-slice policy 存在；identity-aware production policy 要等 authentication 和 session validation 被确认后再实现。

第一条 inventory Protobuf/domain bridge 位于 `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`。它把 generated inventory wire payloads 映射为 inventory runtime request structs，也把 inventory runtime results/events 映射回 generated Protobuf payloads。

不要在本模块中直接 import generated Protobuf types。Protocol adapters 或 generated bridges 应把 wire payloads 转换成 inventory runtime request structs。

PostgreSQL persistence work 必须遵循 `docs/postgresql-persistence-boundary.md`、`docs/postgresql-verification-environment.md` 和 `ADR-0020`。

第一版 durable implementation 中：

- Inventory repository interfaces 继续由本模块拥有。
- PostgreSQL adapter code 位于 `runtime/internal/platform/persistence/postgres/`。
- SQL migrations 位于 `runtime/migrations/postgres/`。
- 第一版 migration source 是 `runtime/migrations/postgres/000001_create_inventory_state.sql`。
- `GrantItem` 必须在 request validation 和 permission checks 之后调用 `LockInventoryForMutation`，然后用返回的 `MutationLock` 读取 current inventory 并执行 grant mutation。
- PostgreSQL adapters 必须在 application-owned unit of work 内，把这个 lock 实现为 `player_id` 对应的 inventory account row lock。
- `MutationLock.Release` 只释放 aggregate lock 或 adapter-local resource；它不得 commit 或 roll back transaction。
- 第一版 PostgreSQL adapter 是 `runtime/internal/platform/persistence/postgres/inventory_repository.go`，focused tests 位于 `runtime/internal/platform/persistence/postgres/inventory_repository_test.go`。
- Durable grant behavior 必须在同一个 application-owned unit of work 中记录 item quantity change 和 `ItemGranted` grant record。`GrantItemMutation` 为此携带 `event_id`、`occurred_at` 和 `reason`。
- Migration apply/status 和 repository integration verification 需要通过 `VIBIT_POSTGRES_TEST_DSN` 显式提供 disposable PostgreSQL environment。
- Live PostgreSQL repository integration tests 是 opt-in；当 `VIBIT_POSTGRES_TEST_DSN` 未设置时，必须 skip 并在 verification notes 中明确记录。

## Forbidden Shortcuts

- 不要绕过 `module.yaml` 中声明的边界。
- 不要直接修改其他模块拥有的数据。
- 不要添加未登记的 public commands、queries、events 或 permissions。
- 不要把 inventory business rules 放进 WebSocket、HTTP、Protobuf 或 transport handlers。
- 不要让本模块直接依赖第三方 WebSocket 或 Protobuf libraries。
- 不要让本模块直接依赖 PostgreSQL drivers、S3 SDKs 或 MinIO clients。开始 persistence implementation 时，应使用 vibit-owned repository 和 storage interfaces。
- 不要在 command flows 中把 transaction creation 隐藏到 inventory repositories 里。
- 不要在 mutation lock 之外，为 capacity-sensitive grants 读取 current inventory。
- 不要从 `GrantItemMutation` 中移除 event metadata；durable adapters 需要它把 `inventory_item_grants` 与 quantity update 原子化持久化。
- 不要手工编辑 generated files。如果 generated output 错了，应修改 source contract、template 或 generator。
- 不要在实现中临时发明 payload fields。必须先更新对应 contract source file。
- 未先更新 manifest 和 change spec，不要引入对 player、currency、reward、quest 或 match modules 的 dependency。
- 在 player/session boundary 存在真实 validation implementation 前，不要把 client-supplied `player_id` 或 `session_id` metadata 当作 authenticated proof。
- 不要把 `RequestIdentity.Status == metadata_only` 当作 privileged inventory grants 的充分条件。
- 不要在 inventory 中添加 player account migrations、token parsing、credential storage 或 authentication provider code。
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

修改 inventory runtime behavior 后，运行 `node tools/vibit check runtime`。当 Go 可用时，还应运行 `cd runtime && go test ./...`。

当 inventory persistence 开始时，PostgreSQL 是第一版 authoritative durable store。除非未来 contract 引入 inventory 自己拥有的大对象 artifacts，否则第一条 inventory slice 不需要 S3-compatible object storage。
