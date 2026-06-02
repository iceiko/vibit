# Pitaya-Aligned Startup And Shutdown Hook Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-02
Scope: Gate-only boundary for using Pitaya-aligned startup and shutdown hook vocabulary after the component discovery and module loading source-first map
Depends on: `decisions/ADR-0196-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map.md`, `decisions/ADR-0195-pitaya-aligned-component-discovery-module-loading-source-first-map.md`, `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0197`

The paired Simplified Chinese translation is `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a startup and shutdown hook vocabulary gate only. It does not implement startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoint behavior, dashboards, admin console behavior, transport behavior changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned startup and shutdown hook boundary gate record is:

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
allowed_startup_shutdown_hook_vocabulary:
  - startup_hook_boundary
  - shutdown_hook_boundary
  - lifecycle_hook_boundary
  - explicit_bootstrap_order_source
  - startup_ordering_deferral
  - shutdown_ordering_deferral
  - dependency_handoff_deferral
  - module_loading_handoff_deferral
  - distributed_lifecycle_deferral
current_source_first_startup_shutdown_hook_mapping:
  explicit_bootstrap_composition:
    current: runtime/cmd/vibit-server and runtime/internal/app/bootstrap source files
    future_vocabulary: explicit_bootstrap_order_source
    implementation_status: source_first_repository_inspection_only
  startup_hooks:
    current: deferred
    future_vocabulary: startup_hook_boundary
    implementation_status: no_startup_hook_behavior
  shutdown_hooks:
    current: deferred
    future_vocabulary: shutdown_hook_boundary
    implementation_status: no_shutdown_hook_behavior
  lifecycle_hook_execution:
    current: deferred
    future_vocabulary: lifecycle_hook_boundary
    implementation_status: no_lifecycle_hook_execution
  dependency_handoff:
    current: explicit constructor and unit-of-work handoff only
    future_vocabulary: dependency_handoff_deferral
    implementation_status: no_dependency_container_behavior
  distributed_lifecycle:
    current: deferred
    future_vocabulary: distributed_lifecycle_deferral
    implementation_status: no_distributed_runtime_behavior
startup_hook_behavior_added: false
shutdown_hook_behavior_added: false
lifecycle_hook_execution_added: false
dependency_container_behavior_added: false
component_discovery_behavior_added: false
component_loading_behavior_added: false
component_module_loading_behavior_added: false
handler_module_registration_behavior_added: false
handler_registration_behavior_added: false
dynamic_handler_registration_added: false
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

`ADR-0196` selected startup and shutdown hooks as the next Pitaya-aligned direction after the component discovery and module loading source-first map.

The risk is that agents may treat hook vocabulary as permission to wire process startup, add shutdown execution, introduce dependency containers, add plugin behavior, or change runtime lifecycle behavior. This gate records vocabulary and source-first mapping only. It keeps the local alpha runtime behavior unchanged and prepares a narrow source-first map follow-up.

## 3. Vocabulary

Allowed startup and shutdown hook vocabulary:

- `startup_hook_boundary`: future planning vocabulary for startup hook concepts without implementing startup hook behavior.
- `shutdown_hook_boundary`: future planning vocabulary for shutdown hook concepts without implementing shutdown hook behavior.
- `lifecycle_hook_boundary`: future planning vocabulary for hook execution contracts without implementing execution.
- `explicit_bootstrap_order_source`: future planning vocabulary for committed bootstrap composition sources.
- `startup_ordering_deferral`: future planning vocabulary for startup order that remains deferred.
- `shutdown_ordering_deferral`: future planning vocabulary for shutdown order that remains deferred.
- `dependency_handoff_deferral`: future planning vocabulary for dependency handoff without dependency container behavior.
- `module_loading_handoff_deferral`: future planning vocabulary for module loading handoff without module loading behavior.
- `distributed_lifecycle_deferral`: future planning vocabulary for distributed lifecycle behavior that remains deferred.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, hook, lifecycle, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use startup or shutdown hook vocabulary as permission to add process startup rewiring, shutdown execution, hook interfaces, dependency containers, dynamic loading, module loaders, component registries, runtime endpoints, dashboards, admin console behavior, protocol messages, generated output, persistence, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, session data payloads, component lifecycle payloads, handler registration payloads, component inventory payloads, module loading payloads, startup hook payloads, shutdown hook payloads, or concrete transport metadata as log-safe in this gate.

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

Startup and shutdown hook vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
startup_hook_behavior_owner: deferred
shutdown_hook_behavior_owner: deferred
lifecycle_hook_execution_owner: deferred
dependency_container_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
distributed_runtime_owner: deferred
```

Rules:

- Documentation and manifests may define startup and shutdown hook vocabulary and current source-first mapping.
- `tools/vibit` may later emit a source-first startup and shutdown hook map if a follow-up implementation work item authorizes it.
- Existing explicit startup, bootstrap, route registration, payload registry, repository, protocol adapter, and transport behavior remain unchanged by this gate.
- Startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency containers, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, protocol payloads, repository interfaces, migrations, generated output, dashboard behavior, admin console behavior, dependencies, hosted surfaces, SDKs, release artifacts, and distributed runtime behavior remain unchanged by this gate.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for components, handlers, services, sessions, routes, RPC, service discovery, groups, lifecycle hooks, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, startup/shutdown parity, lifecycle hook parity, dependency container parity, or runtime behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- startup hook behavior;
- shutdown hook behavior;
- lifecycle hook execution;
- dependency container behavior;
- component discovery behavior;
- component loading behavior;
- component module loading behavior;
- handler module registration behavior;
- handler registration behavior;
- dynamic handler registration;
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
node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate
node tools/vibit check change define-pitaya-aligned-startup-shutdown-hook-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
