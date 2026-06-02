# Pitaya-Aligned Runtime Component Lifecycle Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-02
Scope: Gate-only boundary for using Pitaya-aligned runtime component lifecycle vocabulary after the dashboard/admin operations source-first map
Depends on: `decisions/ADR-0187-select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map.md`, `decisions/ADR-0186-pitaya-aligned-dashboard-admin-operations-source-first-map.md`, `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0188`

The paired Simplified Chinese translation is `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a runtime component lifecycle vocabulary gate only. It does not implement runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, transport behavior changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned runtime component lifecycle boundary gate record is:

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
current_source_first_runtime_component_lifecycle_mapping:
  bootstrap_composition:
    current: runtime cmd and app bootstrap source files
    future_vocabulary: bootstrap_composition_boundary
    implementation_status: source_first_repository_inspection_only
  handler_registration:
    current: explicit route registration and payload registry source files
    future_vocabulary: handler_registration_boundary
    implementation_status: no_dynamic_handler_registration_behavior
  application_services:
    current: application-owned service packages
    future_vocabulary: runtime_component_boundary
    implementation_status: no_component_lifecycle_interface
  repository_wiring:
    current: unit-of-work and repository interface composition
    future_vocabulary: component_dependency_boundary
    implementation_status: no_dependency_container_behavior
  lifecycle_state:
    current: process startup and transport close behavior only
    future_vocabulary: component_state_posture
    implementation_status: no_component_start_or_shutdown_hooks
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

`ADR-0187` selected runtime component lifecycle as the next Pitaya-aligned direction after the dashboard/admin operations source-first map.

The risk is that agents may treat component lifecycle vocabulary as permission to add dynamic components, handler registration behavior, startup or shutdown hooks, component discovery, dependency containers, or distributed runtime behavior. This gate records vocabulary and source-first mapping only. It keeps the current local alpha runtime behavior unchanged and prepares a narrow source-first map follow-up.

## 3. Vocabulary

Allowed runtime component lifecycle vocabulary:

- `runtime_component_boundary`: future planning vocabulary for a runtime-owned component boundary.
- `component_lifecycle_phase`: future planning vocabulary for lifecycle phases such as configured, started, stopping, or stopped.
- `component_start_boundary`: future planning vocabulary for start behavior that remains deferred.
- `component_shutdown_boundary`: future planning vocabulary for shutdown behavior that remains deferred.
- `handler_registration_boundary`: future planning vocabulary for route or handler registration ownership.
- `component_dependency_boundary`: future planning vocabulary for dependencies a component needs without introducing a dependency container.
- `bootstrap_composition_boundary`: future planning vocabulary for explicit process and application composition.
- `component_state_posture`: future planning vocabulary for state that may later be inspectable.
- `local_component_inventory`: future planning vocabulary for source-first component inventory.
- `distributed_component_lifecycle_deferral`: future planning vocabulary for distributed lifecycle behavior that remains deferred.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, component, lifecycle, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use component lifecycle vocabulary as permission to add lifecycle interfaces, dynamic handler registration, component discovery, component loading, startup hooks, shutdown hooks, dependency containers, runtime endpoints, dashboards, admin console behavior, protocol messages, generated output, persistence, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, session data payloads, component lifecycle payloads, handler registration payloads, or concrete transport metadata as log-safe in this gate.

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

Runtime component lifecycle vocabulary ownership:

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

Rules:

- Documentation and manifests may define runtime component lifecycle vocabulary and current source-first mapping.
- `tools/vibit` may later emit a source-first runtime component lifecycle map if a follow-up implementation work item authorizes it.
- Existing explicit startup, bootstrap, route registration, payload registry, repository, and transport close behavior remain unchanged by this gate.
- Runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, protocol payloads, repository interfaces, migrations, generated output, dashboard behavior, admin console behavior, dependencies, hosted surfaces, SDKs, release artifacts, and distributed runtime behavior remain unchanged by this gate.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for runtime components, handlers, services, sessions, routes, RPC, service discovery, groups, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, component lifecycle parity, handler registration parity, startup/shutdown parity, or runtime behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- runtime component lifecycle behavior;
- handler registration behavior;
- component discovery or loading;
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

The repository check rule is `runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate`.

The check verifies the standard, translation, ADR, change artifacts, manifest references, next-ready state, vocabulary markers, and explicit implementation deferrals.
