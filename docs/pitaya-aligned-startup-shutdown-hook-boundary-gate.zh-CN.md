# Pitaya-Aligned Startup And Shutdown Hook Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-02
Scope: 在 component discovery and module loading source-first map 之后，仅用于 Pitaya-aligned startup and shutdown hook vocabulary 的 gate-only boundary
Depends on: `decisions/ADR-0196-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map.md`, `decisions/ADR-0195-pitaya-aligned-component-discovery-module-loading-source-first-map.md`, `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0197`

英文版 `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md` 是权威版本。本文件是配套简体中文说明。

本文只定义 startup and shutdown hook vocabulary gate。它不实现 startup hook behavior、shutdown hook behavior、lifecycle hook execution、dependency container behavior、component discovery behavior、component loading behavior、component module loading behavior、handler module registration behavior、handler registration behavior、dynamic handler registration、runtime endpoint behavior、dashboards、admin console behavior、transport behavior changes、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployment、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned startup and shutdown hook boundary gate 记录是：

```yaml
pitaya_aligned_startup_shutdown_hook_boundary_gate: defined
completed_work_item: W-0289
decision: ADR-0197
check_rule: runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate
selection_decision: ADR-0196
source_map_decision: ADR-0195
component_discovery_module_loading_gate_decision: ADR-0194
standard: docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md
translation: docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: startup_shutdown_hook_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_startup_shutdown_hook_vocabulary
future_implementation_work_item: W-0290
future_implementation_direction: implement_pitaya_aligned_startup_shutdown_hook_source_first_map
```

该 gate 只允许记录 startup/shutdown hook 相关 vocabulary 与当前 source-first mapping，不授权任何运行时 hook 执行或启动/关闭 wiring。

## 2. Purpose

`ADR-0196` 在 component discovery and module loading source-first map 之后选择了 startup and shutdown hooks 作为下一项 Pitaya-aligned direction。

风险在于 agent 可能把 hook vocabulary 当作添加 process startup rewiring、shutdown execution、dependency containers、plugin behavior 或 runtime lifecycle behavior 的许可。本 gate 只记录 vocabulary 和 source-first mapping，保持 local alpha runtime behavior 不变，并准备后续 narrow source-first map。

## 3. Vocabulary

允许的 vocabulary：

- `startup_hook_boundary`
- `shutdown_hook_boundary`
- `lifecycle_hook_boundary`
- `explicit_bootstrap_order_source`
- `startup_ordering_deferral`
- `shutdown_ordering_deferral`
- `dependency_handoff_deferral`
- `module_loading_handoff_deferral`
- `distributed_lifecycle_deferral`

禁止使用 vocabulary 去添加 process startup rewiring、shutdown execution、hook interfaces、dependency containers、dynamic loading、module loaders、component registries、runtime endpoints、dashboards、admin console behavior、protocol messages、generated output、persistence、hosted surfaces、SDKs、release artifacts 或 distributed runtime behavior。

## 4. Current Mapping

```yaml
current_source_first_startup_shutdown_hook_mapping:
  explicit_bootstrap_composition:
    current: runtime/cmd/vibit-server and runtime/internal/app/bootstrap source files
    future_vocabulary: explicit_bootstrap_order_source
    status: source_first_repository_inspection_only
  startup_hooks:
    current: deferred
    future_vocabulary: startup_hook_boundary
    status: no_startup_hook_behavior
  shutdown_hooks:
    current: deferred
    future_vocabulary: shutdown_hook_boundary
    status: no_shutdown_hook_behavior
  lifecycle_hook_execution:
    current: deferred
    future_vocabulary: lifecycle_hook_boundary
    status: no_lifecycle_hook_execution
  dependency_handoff:
    current: explicit constructor and unit-of-work handoff only
    future_vocabulary: dependency_handoff_deferral
    status: no_dependency_container_behavior
  distributed_lifecycle:
    current: deferred
    future_vocabulary: distributed_lifecycle_deferral
    status: no_distributed_runtime_behavior
```

## 5. Ownership

Architecture vocabulary 由 `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md`、`.arch/reference.yaml`、`.arch/runtime.yaml` 拥有。后续 source-first map 可以由 `tools/vibit` 输出，但必须等待后续 bounded work item 授权。

现有 startup、bootstrap、route registration、payload registry、repository、protocol adapter 和 transport behavior 都不因本 gate 变化。

## 6. Nakama And Pitaya Mapping

Nakama 仍是广义 game backend product capability pressure 的主要产品参考。Pitaya 仍作为 components、handlers、services、sessions、routes、RPC、service discovery、groups、lifecycle hooks 和 operations 的 architecture vocabulary reference。

本 gate 只映射到 vibit-owned vocabulary，不创建 direct compatibility、public API parity、startup/shutdown parity、lifecycle hook parity、dependency container parity 或 runtime behavior。

## 7. Stop Conditions

添加以下内容前必须停止并等待后续 bounded work item：

- startup hook behavior；
- shutdown hook behavior；
- lifecycle hook execution；
- dependency container behavior；
- component discovery behavior；
- component loading behavior；
- component module loading behavior；
- handler module registration behavior；
- handler registration behavior；
- dynamic handler registration；
- runtime endpoint behavior；
- dashboards；
- admin console behavior；
- protocol messages or routes；
- Protobuf source；
- generated output；
- repository interfaces；
- PostgreSQL adapters；
- migrations；
- dependencies；
- hosted deployment；
- SDK publication；
- release artifacts；
- distributed runtime behavior；
- direct Nakama/Pitaya API compatibility。

## 8. Verification

Repository verification:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate
node tools/vibit check change define-pitaya-aligned-startup-shutdown-hook-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
