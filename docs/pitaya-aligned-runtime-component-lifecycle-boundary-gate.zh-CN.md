# Pitaya-Aligned Runtime Component Lifecycle Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-02
Scope: dashboard/admin operations source-first map 之后使用 Pitaya-aligned runtime component lifecycle vocabulary 的 gate-only boundary
Depends on: `decisions/ADR-0187-select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map.md`, `decisions/ADR-0186-pitaya-aligned-dashboard-admin-operations-source-first-map.md`, `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0188`

英文原文 `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md` 是权威版本。本文件是配套简体中文翻译。

This document defines a runtime component lifecycle vocabulary gate only. 本 gate 不实现 runtime component lifecycle behavior、handler registration behavior、component discovery or loading、startup hooks、shutdown hooks、runtime endpoint behavior、dashboards、admin console behavior、transport behavior changes、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployment、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned runtime component lifecycle boundary gate 记录如下：

```yaml
pitaya_aligned_runtime_component_lifecycle_boundary_gate: defined
completed_work_item: W-0280
decision: ADR-0188
check_rule: runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate
selection_decision: ADR-0187
dashboard_admin_operations_source_first_map_decision: ADR-0186
dashboard_admin_operations_boundary_gate_decision: ADR-0185
standard: docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md
translation: docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: runtime_component_lifecycle_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_runtime_component_lifecycle_vocabulary
future_implementation_work_item: W-0281
future_implementation_direction: implement_pitaya_aligned_runtime_component_lifecycle_source_first_map
allowed_runtime_component_lifecycle_vocabulary:
  - runtime_component_boundary
  - component_lifecycle_phase
  - component_start_boundary
  - component_shutdown_boundary
  - handler_registration_boundary
  - component_dependency_boundary
  - bootstrap_composition_boundary
  - component_state_posture
  - local_component_inventory
  - distributed_component_lifecycle_deferral
runtime_component_lifecycle_behavior_added: false
component_lifecycle_behavior_added: false
handler_registration_behavior_added: false
component_discovery_added: false
component_loading_added: false
startup_hook_behavior_added: false
shutdown_hook_behavior_added: false
runtime_endpoint_behavior_added: false
dashboard_added: false
admin_console_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
hosted_deployment_added: false
sdk_added: false
release_artifact_added: false
distributed_runtime_implementation_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0187` 在 dashboard/admin operations source-first map 之后选择 runtime component lifecycle 作为下一项 Pitaya-aligned direction。

风险在于 agent 可能把 component lifecycle vocabulary 误解为可以添加 dynamic components、handler registration behavior、startup 或 shutdown hooks、component discovery、dependency containers 或 distributed runtime behavior。本 gate 只记录 vocabulary 和 source-first mapping，保持当前 local alpha runtime behavior 不变，并准备一个窄的 source-first map 后续项。

## 3. Vocabulary

允许的 runtime component lifecycle vocabulary：

- `runtime_component_boundary`: runtime-owned component boundary 的未来规划词汇。
- `component_lifecycle_phase`: configured、started、stopping、stopped 等 lifecycle phase 的未来规划词汇。
- `component_start_boundary`: 仍延后的 start behavior 的未来规划词汇。
- `component_shutdown_boundary`: 仍延后的 shutdown behavior 的未来规划词汇。
- `handler_registration_boundary`: route 或 handler registration ownership 的未来规划词汇。
- `component_dependency_boundary`: component 需要的 dependency 的未来规划词汇，但不引入 dependency container。
- `bootstrap_composition_boundary`: explicit process 和 application composition 的未来规划词汇。
- `component_state_posture`: 未来可能 inspectable 的 state posture 词汇。
- `local_component_inventory`: source-first component inventory 的未来规划词汇。
- `distributed_component_lifecycle_deferral`: distributed lifecycle behavior 仍延后的未来规划词汇。

禁止用法：

- 不要引入来自 Pitaya 或 Nakama 的 concrete public API、package、route、method、wire、handler、component、lifecycle、dashboard、metrics、trace、admin、console 或 inspector compatibility names。
- 不要把 component lifecycle vocabulary 当作添加 lifecycle interfaces、dynamic handler registration、component discovery、component loading、startup hooks、shutdown hooks、dependency containers、runtime endpoints、dashboards、admin console behavior、protocol messages、generated output、persistence、hosted surfaces、SDKs、release artifacts 或 distributed runtime behavior 的授权。
- 不要在本 gate 中把 raw tokens、credentials、lookup digests、verifier digests、verifier keys、DSNs、headers、cookies、query strings、subprotocol values、remote addresses、database payloads、local secret file contents、route payloads、session data payloads、component lifecycle payloads、handler registration payloads 或 concrete transport metadata 归类为 log-safe。

## 4. Current Mapping

```yaml
current_source_first_runtime_component_lifecycle_mapping:
  bootstrap_composition:
    current: runtime/cmd/vibit-server and runtime/internal/app/bootstrap source files
    future_vocabulary: bootstrap_composition_boundary
    status: source_first_repository_inspection_only
  handler_registration:
    current: route registration and payload registry source files
    future_vocabulary: handler_registration_boundary
    status: no_dynamic_handler_registration_behavior
  application_services:
    current: runtime/internal/app service packages
    future_vocabulary: runtime_component_boundary
    status: no_component_lifecycle_behavior
  persistence_composition:
    current: unit-of-work and repository interface boundaries
    future_vocabulary: component_dependency_boundary
    status: no_dependency_container_behavior
  lifecycle_state:
    current: process startup and existing transport close behavior
    future_vocabulary: component_state_posture
    status: no_component_start_or_shutdown_hooks
```

## 5. Ownership

Runtime component lifecycle vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_component_behavior_owner: deferred
handler_registration_owner: deferred
startup_hook_owner: deferred
shutdown_hook_owner: deferred
component_discovery_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
dependency_owner: unchanged
```

规则：

- Documentation 和 manifests 可以定义 runtime component lifecycle vocabulary 与当前 source-first mapping。
- 如果后续 implementation work item 授权，`tools/vibit` 可以输出 source-first runtime component lifecycle map。
- 现有 explicit startup、bootstrap、route registration、payload registry、repository 与 transport close behavior 不因本 gate 改变。
- Runtime component lifecycle behavior、handler registration behavior、component discovery or loading、startup hooks、shutdown hooks、protocol payloads、repository interfaces、migrations、generated output、dashboard behavior、admin console behavior、dependencies、hosted surfaces、SDKs、release artifacts 与 distributed runtime behavior 不因本 gate 改变。

## 6. Nakama And Pitaya Mapping

Nakama 仍是 broad game backend product capability pressure 的 primary product reference。Pitaya 仍是 runtime components、handlers、services、sessions、routes、RPC、service discovery、groups 和 operational concerns 的 architecture vocabulary reference。

本 gate 只把这些 reference 映射到 vibit-owned vocabulary。它不创建 direct compatibility、public API parity、component lifecycle parity、handler registration parity、startup/shutdown parity 或 runtime behavior。

## 7. Stop Conditions

添加以下内容前必须停止，并等待后续 bounded work item：

- runtime component lifecycle behavior；
- handler registration behavior；
- component discovery or loading；
- startup hooks；
- shutdown hooks；
- dependency containers；
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

Repository check rule 是 `runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate`。

该检查验证 standard、translation、ADR、change artifacts、manifest references、next-ready state、vocabulary markers 和 explicit implementation deferrals。
