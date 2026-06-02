# Pitaya-Aligned Handler Module Registration Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-02
Scope: runtime component lifecycle source-first map 之后使用 Pitaya-aligned handler module registration vocabulary 的 gate-only boundary
Depends on: `decisions/ADR-0190-select-next-pitaya-aligned-direction-after-runtime-component-lifecycle-map.md`, `decisions/ADR-0189-pitaya-aligned-runtime-component-lifecycle-source-first-map.md`, `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0191`

英文原文 `docs/pitaya-aligned-handler-module-registration-boundary-gate.md` 是权威版本。本文件是配套简体中文翻译。

This document defines a handler module registration vocabulary gate only. 本 gate 不实现 handler module registration behavior、handler registration behavior、dynamic handler registration、component discovery or loading、component module loading、startup hooks、shutdown hooks、runtime endpoint behavior、dashboards、admin console behavior、transport behavior changes、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployment、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned handler module registration boundary gate 记录如下：

```yaml
pitaya_aligned_handler_module_registration_boundary_gate: defined
completed_work_item: W-0283
decision: ADR-0191
check_rule: runtime.pitaya_aligned_handler_module_registration_boundary_gate
selection_decision: ADR-0190
runtime_component_lifecycle_source_first_map_decision: ADR-0189
runtime_component_lifecycle_boundary_gate_decision: ADR-0188
standard: docs/pitaya-aligned-handler-module-registration-boundary-gate.md
translation: docs/pitaya-aligned-handler-module-registration-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: handler_module_registration_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_handler_module_registration_vocabulary
future_implementation_work_item: W-0284
future_implementation_direction: implement_pitaya_aligned_handler_module_registration_source_first_map
allowed_handler_module_registration_vocabulary:
  - handler_module_boundary
  - handler_registration_boundary
  - route_handler_ownership
  - explicit_registration_source
  - payload_registry_boundary
  - module_bootstrap_boundary
  - handler_dependency_boundary
  - handler_execution_context_boundary
  - local_handler_inventory
  - dynamic_registration_deferral
  - distributed_handler_module_deferral
runtime_component_lifecycle_behavior_added: false
handler_module_registration_behavior_added: false
handler_registration_behavior_added: false
dynamic_handler_registration_added: false
component_discovery_added: false
component_loading_added: false
component_module_loading_added: false
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

`ADR-0190` 在 runtime component lifecycle source-first map 之后选择 handler module registration 作为下一项 Pitaya-aligned direction。

风险在于 agent 可能把 handler module vocabulary 误解为可以添加 dynamic handler registries、component/module loading、startup 或 shutdown hooks、dependency containers、protocol routes 或 distributed runtime behavior。本 gate 只记录 vocabulary 和 source-first mapping，保持当前 local alpha runtime behavior 不变，并准备一个窄的 source-first map 后续项。

## 3. Vocabulary

允许的 handler module registration vocabulary：

- `handler_module_boundary`: application-owned handler module boundary 的未来规划词汇。
- `handler_registration_boundary`: handler registration surface ownership 的未来规划词汇。
- `route_handler_ownership`: 拥有 route handler 的 module 的未来规划词汇。
- `explicit_registration_source`: 显式注册 route 或 handler 的已提交 source file 的未来规划词汇。
- `payload_registry_boundary`: payload registry 与 bridge boundary 的未来规划词汇，仍由 protocol adapter 拥有。
- `module_bootstrap_boundary`: modules 与 routes 的显式 bootstrap composition 未来规划词汇。
- `handler_dependency_boundary`: handler dependency handoff 的未来规划词汇，但不引入 dependency container。
- `handler_execution_context_boundary`: request/session context handoff into handlers 的未来规划词汇。
- `local_handler_inventory`: handler modules 和 registrations 的 source-first local inventory 未来规划词汇。
- `dynamic_registration_deferral`: dynamic handler registration 仍延后的未来规划词汇。
- `distributed_handler_module_deferral`: distributed handler/module registration 仍延后的未来规划词汇。

禁止用法：

- 不要引入来自 Pitaya 或 Nakama 的 concrete public API、package、route、method、wire、handler、module、lifecycle、dashboard、metrics、trace、admin、console 或 inspector compatibility names。
- 不要把 handler module registration vocabulary 当作添加 handler registries、dynamic handler registration、component discovery、component loading、component module loading、startup hooks、shutdown hooks、dependency containers、runtime endpoints、dashboards、admin console behavior、protocol messages、generated output、persistence、hosted surfaces、SDKs、release artifacts 或 distributed runtime behavior 的授权。
- 不要在本 gate 中把 raw tokens、credentials、lookup digests、verifier digests、verifier keys、DSNs、headers、cookies、query strings、subprotocol values、remote addresses、database payloads、local secret file contents、route payloads、session data payloads、component lifecycle payloads、handler registration payloads、module registration payloads 或 concrete transport metadata 归类为 log-safe。

## 4. Current Mapping

```yaml
current_source_first_handler_module_registration_mapping:
  module_bootstrap:
    current: runtime/cmd/vibit-server and runtime/internal/app/bootstrap source files
    future_vocabulary: module_bootstrap_boundary
    status: source_first_repository_inspection_only
  explicit_route_registration:
    current: runtime/internal/app route registration source files
    future_vocabulary: explicit_registration_source
    status: no_dynamic_handler_registration_behavior
  payload_registry:
    current: runtime/internal/platform/protocol/protobuf bridge and payload registry source files
    future_vocabulary: payload_registry_boundary
    status: no_protocol_shape_change
  application_module_ownership:
    current: runtime/internal/app service packages and runtime/internal/modules domain packages
    future_vocabulary: handler_module_boundary
    status: no_component_module_loading_behavior
  dependency_handoff:
    current: unit-of-work and repository interface boundaries
    future_vocabulary: handler_dependency_boundary
    status: no_dependency_container_behavior
  execution_context:
    current: existing application request context and session handoff boundaries
    future_vocabulary: handler_execution_context_boundary
    status: no_runtime_endpoint_or_protocol_behavior
  distributed_registration:
    current: deferred
    future_vocabulary: distributed_handler_module_deferral
    status: no_distributed_runtime_behavior
```

## 5. Ownership

Handler module registration vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-handler-module-registration-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
handler_module_registration_behavior_owner: deferred
handler_registration_owner: deferred
dynamic_registration_owner: deferred
component_module_loading_owner: deferred
startup_hook_owner: deferred
shutdown_hook_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
dependency_owner: unchanged
distributed_runtime_owner: deferred
```

规则：

- Documentation 和 manifests 可以定义 handler module registration vocabulary 与当前 source-first mapping。
- 如果后续 implementation work item 授权，`tools/vibit` 可以输出 source-first handler module registration map。
- 现有 explicit startup、bootstrap、route registration、payload registry、repository、protocol adapter 与 transport close behavior 不因本 gate 改变。
- Handler module registration behavior、handler registration behavior、dynamic handler registration、component discovery or loading、component module loading、startup hooks、shutdown hooks、protocol payloads、repository interfaces、migrations、generated output、dashboard behavior、admin console behavior、dependencies、hosted surfaces、SDKs、release artifacts 与 distributed runtime behavior 不因本 gate 改变。

## 6. Nakama And Pitaya Mapping

Nakama 仍是 broad game backend product capability pressure 的 primary product reference。Pitaya 仍是 runtime components、handlers、services、sessions、routes、RPC、service discovery、groups 和 operational concerns 的 architecture vocabulary reference。

本 gate 只把这些 reference 映射到 vibit-owned vocabulary。它不创建 direct compatibility、public API parity、component lifecycle parity、handler registration parity、handler module registration parity、startup/shutdown parity 或 runtime behavior。

## 7. Stop Conditions

添加以下内容前必须停止，并等待后续 bounded work item：

- handler module registration behavior；
- handler registration behavior；
- dynamic handler registration；
- component discovery or loading；
- component module loading；
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

Repository check rule 是 `runtime.pitaya_aligned_handler_module_registration_boundary_gate`。

该检查验证 standard、translation、ADR、change artifacts、manifest references、next-ready state、vocabulary markers 和 explicit implementation deferrals。
