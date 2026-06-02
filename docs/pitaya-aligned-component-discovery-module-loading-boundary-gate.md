# Pitaya-Aligned Component Discovery And Module Loading Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-02
Scope: Gate-only boundary for using Pitaya-aligned component discovery and module loading vocabulary after the handler module registration source-first map
Depends on: `decisions/ADR-0193-select-next-pitaya-aligned-direction-after-handler-module-registration-map.md`, `decisions/ADR-0192-pitaya-aligned-handler-module-registration-source-first-map.md`, `docs/pitaya-aligned-handler-module-registration-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0194`

The paired Simplified Chinese translation is `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a component discovery and module loading vocabulary gate only. It does not implement component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, transport behavior changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned component discovery and module loading boundary gate record is:

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
current_source_first_component_discovery_module_loading_mapping:
  explicit_bootstrap_composition:
    current: runtime/cmd/vibit-server and runtime/internal/app/bootstrap source files
    future_vocabulary: bootstrap_component_source
    implementation_status: source_first_repository_inspection_only
  application_module_ownership:
    current: runtime/internal/app service packages and runtime/internal/modules domain packages
    future_vocabulary: application_module_ownership
    implementation_status: no_component_module_loading_behavior
  handler_module_registration_mapping:
    current: docs/pitaya-aligned-handler-module-registration-boundary-gate.md and node tools/vibit inspect pitaya-handler-modules --json
    future_vocabulary: handler_module_registration_handoff
    implementation_status: no_handler_module_registration_behavior
  component_inventory:
    current: committed source surfaces only
    future_vocabulary: explicit_component_inventory
    implementation_status: no_runtime_discovery_or_loader
  dynamic_loading:
    current: deferred
    future_vocabulary: dynamic_loading_deferral
    implementation_status: no_dynamic_loading_behavior
  distributed_discovery:
    current: deferred
    future_vocabulary: distributed_component_discovery_deferral
    implementation_status: no_distributed_runtime_behavior
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

`ADR-0193` selected component discovery and module loading as the next Pitaya-aligned direction after the handler module registration source-first map.

The risk is that agents may treat component discovery vocabulary as permission to add runtime scanning, dynamic component loading, module loaders, startup/shutdown hooks, dependency containers, runtime endpoints, or distributed discovery. This gate records vocabulary and source-first mapping only. It keeps the local alpha runtime behavior unchanged and prepares a narrow source-first map follow-up.

## 3. Vocabulary

Allowed component discovery and module loading vocabulary:

- `component_discovery_boundary`: future planning vocabulary for discovering component definitions without implementing discovery.
- `component_loading_boundary`: future planning vocabulary for loading component definitions without implementing a loader.
- `component_module_loading_boundary`: future planning vocabulary for module loading ownership that remains deferred.
- `explicit_component_inventory`: future planning vocabulary for committed source surfaces that enumerate components.
- `bootstrap_component_source`: future planning vocabulary for bootstrap composition source files.
- `application_module_ownership`: future planning vocabulary for application service and domain module ownership.
- `handler_module_registration_handoff`: future planning vocabulary for handoff from handler module registration maps.
- `dynamic_loading_deferral`: future planning vocabulary for dynamic loading that remains deferred.
- `distributed_component_discovery_deferral`: future planning vocabulary for distributed discovery that remains deferred.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, module, lifecycle, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use component discovery or module loading vocabulary as permission to add runtime scanning, dynamic loading, module loaders, component registries, startup hooks, shutdown hooks, dependency containers, runtime endpoints, dashboards, admin console behavior, protocol messages, generated output, persistence, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, session data payloads, component lifecycle payloads, handler registration payloads, component inventory payloads, module loading payloads, or concrete transport metadata as log-safe in this gate.

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

Component discovery and module loading vocabulary ownership:

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

Rules:

- Documentation and manifests may define component discovery and module loading vocabulary and current source-first mapping.
- `tools/vibit` may later emit a source-first component discovery and module loading map if a follow-up implementation work item authorizes it.
- Existing explicit startup, bootstrap, route registration, payload registry, repository, protocol adapter, and transport behavior remain unchanged by this gate.
- Component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, protocol payloads, repository interfaces, migrations, generated output, dashboard behavior, admin console behavior, dependencies, hosted surfaces, SDKs, release artifacts, and distributed runtime behavior remain unchanged by this gate.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for components, handlers, services, sessions, routes, RPC, service discovery, groups, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, component discovery parity, component loading parity, module loading parity, startup/shutdown parity, or runtime behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- component discovery behavior;
- component loading behavior;
- component module loading behavior;
- handler module registration behavior;
- handler registration behavior;
- dynamic handler registration;
- startup hooks;
- shutdown hooks;
- dependency containers;
- runtime endpoint behavior;
- dashboards;
- admin console behavior;
- protocol messages or routes;
- Protobuf source;
- generated output;
- repository interfaces;
- PostgreSQL adapters;
- migrations;
- dependencies;
- hosted deployment;
- SDK publication;
- release artifacts;
- distributed runtime behavior;
- direct Nakama/Pitaya API compatibility.

## 8. Verification

Repository verification:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate
node tools/vibit check change define-pitaya-aligned-component-discovery-module-loading-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
