# Pitaya-Aligned Component Discovery And Module Loading Boundary Gate

状态：Accepted v0.1
最后更新：2026-06-02
范围：在 handler module registration source-first map 之后使用 Pitaya-aligned component discovery and module loading vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0193-select-next-pitaya-aligned-direction-after-handler-module-registration-map.md`、`decisions/ADR-0192-pitaya-aligned-handler-module-registration-source-first-map.md`、`docs/pitaya-aligned-handler-module-registration-boundary-gate.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
Canonical decision：`ADR-0194`

英文版 `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md` 是权威版本。

本文档只定义 component discovery and module loading vocabulary gate。它不实现 component discovery behavior、component loading behavior、component module loading behavior、handler module registration behavior、handler registration behavior、dynamic handler registration、startup hooks、shutdown hooks、runtime endpoint behavior、dashboards、admin console behavior、transport behavior changes、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployment、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned component discovery and module loading boundary gate record 是：

```yaml
pitaya_aligned_component_discovery_module_loading_boundary_gate: defined
completed_work_item: W-0286
decision: ADR-0194
check_rule: runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate
selection_decision: ADR-0193
handler_module_registration_source_first_map_decision: ADR-0192
handler_module_registration_boundary_gate_decision: ADR-0191
standard: docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md
translation: docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: component_discovery_module_loading_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_component_discovery_module_loading_vocabulary
future_implementation_work_item: W-0287
future_implementation_direction: implement_pitaya_aligned_component_discovery_module_loading_source_first_map
allowed_component_discovery_module_loading_vocabulary:
  - component_discovery_boundary
  - component_loading_boundary
  - component_module_loading_boundary
  - explicit_component_inventory
  - bootstrap_component_source
  - application_module_ownership
  - handler_module_registration_handoff
  - dynamic_loading_deferral
  - distributed_component_discovery_deferral
component_discovery_added: false
component_loading_added: false
component_module_loading_added: false
handler_module_registration_behavior_added: false
handler_registration_behavior_added: false
dynamic_handler_registration_added: false
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

`ADR-0193` 在 handler module registration source-first map 之后选择了 component discovery and module loading 方向。

风险在于 agent 可能把 component discovery vocabulary 理解为可以添加 runtime scanning、dynamic component loading、module loaders、startup/shutdown hooks、dependency containers、runtime endpoints 或 distributed discovery。此 gate 只记录 vocabulary 和 source-first mapping，保持 local alpha runtime behavior 不变，并准备后续 source-first map。

## 3. Vocabulary

允许的 vocabulary：

- `component_discovery_boundary`：用于未来规划 component definition discovery 的 vocabulary，不实现 discovery。
- `component_loading_boundary`：用于未来规划 component definition loading 的 vocabulary，不实现 loader。
- `component_module_loading_boundary`：用于未来规划 module loading ownership 的 vocabulary，仍保持 deferred。
- `explicit_component_inventory`：用于 committed source surfaces 的 source-first inventory vocabulary。
- `bootstrap_component_source`：用于 bootstrap composition source files 的 vocabulary。
- `application_module_ownership`：用于 application service 和 domain module ownership 的 vocabulary。
- `handler_module_registration_handoff`：用于从 handler module registration maps handoff 的 vocabulary。
- `dynamic_loading_deferral`：用于仍然 deferred 的 dynamic loading。
- `distributed_component_discovery_deferral`：用于仍然 deferred 的 distributed discovery。

Forbidden vocabulary use：

- 不要引入来自 Pitaya 或 Nakama 的 concrete public API、package、route、method、wire、handler、module、lifecycle、dashboard、metrics、trace、admin、console 或 inspector compatibility names。
- 不要把 component discovery 或 module loading vocabulary 当作添加 runtime scanning、dynamic loading、module loaders、component registries、startup hooks、shutdown hooks、dependency containers、runtime endpoints、dashboards、admin console behavior、protocol messages、generated output、persistence、hosted surfaces、SDKs、release artifacts 或 distributed runtime behavior 的授权。
- 不要把 raw tokens、credentials、lookup digests、verifier digests、verifier keys、DSNs、headers、cookies、query strings、subprotocol values、remote addresses、database payloads、local secret file contents、route payloads、session data payloads、component lifecycle payloads、handler registration payloads、component inventory payloads、module loading payloads 或 concrete transport metadata 归类为 log-safe。

## 4. Current Mapping

```yaml
current_source_first_component_discovery_module_loading_mapping:
  explicit_bootstrap_composition:
    current: runtime/cmd/vibit-server and runtime/internal/app/bootstrap source files
    future_vocabulary: bootstrap_component_source
    status: source_first_repository_inspection_only
  application_module_ownership:
    current: runtime/internal/app service packages and runtime/internal/modules domain packages
    future_vocabulary: application_module_ownership
    status: no_component_module_loading_behavior
  handler_module_registration_mapping:
    current: docs/pitaya-aligned-handler-module-registration-boundary-gate.md and tools/vibit handler module inspection
    future_vocabulary: handler_module_registration_handoff
    status: no_handler_module_registration_behavior
  component_inventory:
    current: committed source surfaces only
    future_vocabulary: explicit_component_inventory
    status: no_runtime_discovery_or_loader
  dynamic_loading:
    current: deferred
    future_vocabulary: dynamic_loading_deferral
    status: no_dynamic_loading_behavior
  distributed_discovery:
    current: deferred
    future_vocabulary: distributed_component_discovery_deferral
    status: no_distributed_runtime_behavior
```

## 5. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
component_discovery_behavior_owner: deferred
component_loading_behavior_owner: deferred
component_module_loading_owner: deferred
dynamic_loading_owner: deferred
startup_hook_owner: deferred
shutdown_hook_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
dependency_owner: unchanged
distributed_runtime_owner: deferred
```

规则：

- 文档和 manifests 可以定义 component discovery and module loading vocabulary 及当前 source-first mapping。
- 如果后续 work item 授权，`tools/vibit` 可以发出 source-first component discovery and module loading map。
- 现有 explicit startup、bootstrap、route registration、payload registry、repository、protocol adapter 和 transport behavior 不因本 gate 改变。
- Component discovery behavior、component loading behavior、component module loading behavior、handler module registration behavior、handler registration behavior、dynamic handler registration、startup hooks、shutdown hooks、protocol payloads、repository interfaces、migrations、generated output、dashboard behavior、admin console behavior、dependencies、hosted surfaces、SDKs、release artifacts 和 distributed runtime behavior 均保持不变。

## 6. Stop Conditions

添加以下内容前必须停止并等待后续 bounded work item：

- component discovery behavior；
- component loading behavior；
- component module loading behavior；
- handler module registration behavior；
- handler registration behavior；
- dynamic handler registration；
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

## 7. Verification

```text
node tools/vibit inspect rule runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate
node tools/vibit check change define-pitaya-aligned-component-discovery-module-loading-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
