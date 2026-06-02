# Pitaya-Aligned Handler Module Registration Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-02
Scope: Gate-only boundary for using Pitaya-aligned handler module registration vocabulary after the runtime component lifecycle source-first map
Depends on: `decisions/ADR-0190-select-next-pitaya-aligned-direction-after-runtime-component-lifecycle-map.md`, `decisions/ADR-0189-pitaya-aligned-runtime-component-lifecycle-source-first-map.md`, `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0191`

The paired Simplified Chinese translation is `docs/pitaya-aligned-handler-module-registration-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a handler module registration vocabulary gate only. It does not implement handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, transport behavior changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned handler module registration boundary gate record is:

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
current_source_first_handler_module_registration_mapping:
  module_bootstrap:
    current: runtime/internal/app/bootstrap source files and runtime/cmd/vibit-server composition
    future_vocabulary: module_bootstrap_boundary
    implementation_status: source_first_repository_inspection_only
  explicit_route_registration:
    current: application-owned route registration source files
    future_vocabulary: explicit_registration_source
    implementation_status: no_dynamic_handler_registration_behavior
  payload_registry:
    current: Protobuf bridge and payload registry source files
    future_vocabulary: payload_registry_boundary
    implementation_status: no_protocol_shape_change
  application_module_ownership:
    current: runtime/internal/app service packages and runtime/internal/modules domain packages
    future_vocabulary: handler_module_boundary
    implementation_status: no_component_module_loading_behavior
  dependency_handoff:
    current: unit-of-work and repository interface composition
    future_vocabulary: handler_dependency_boundary
    implementation_status: no_dependency_container_behavior
  execution_context:
    current: existing application request context and session handoff boundaries
    future_vocabulary: handler_execution_context_boundary
    implementation_status: no_runtime_endpoint_or_protocol_behavior
  distributed_registration:
    current: deferred
    future_vocabulary: distributed_handler_module_deferral
    implementation_status: no_distributed_runtime_behavior
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

`ADR-0190` selected handler module registration as the next Pitaya-aligned direction after the runtime component lifecycle source-first map.

The risk is that agents may treat handler module vocabulary as permission to add dynamic handler registries, component/module loading, startup or shutdown hooks, dependency containers, protocol routes, or distributed runtime behavior. This gate records vocabulary and source-first mapping only. It keeps the current local alpha runtime behavior unchanged and prepares a narrow source-first map follow-up.

## 3. Vocabulary

Allowed handler module registration vocabulary:

- `handler_module_boundary`: future planning vocabulary for application-owned handler module boundaries.
- `handler_registration_boundary`: future planning vocabulary for ownership of handler registration surfaces.
- `route_handler_ownership`: future planning vocabulary for the module that owns a route handler.
- `explicit_registration_source`: future planning vocabulary for committed source files that register routes or handlers explicitly.
- `payload_registry_boundary`: future planning vocabulary for the payload registry and bridge boundary that remains protocol-adapter owned.
- `module_bootstrap_boundary`: future planning vocabulary for explicit bootstrap composition of modules and routes.
- `handler_dependency_boundary`: future planning vocabulary for dependencies handed to handlers without introducing a dependency container.
- `handler_execution_context_boundary`: future planning vocabulary for request/session context handoff into handlers.
- `local_handler_inventory`: future planning vocabulary for a source-first local inventory of handler modules and registrations.
- `dynamic_registration_deferral`: future planning vocabulary for dynamic handler registration that remains deferred.
- `distributed_handler_module_deferral`: future planning vocabulary for distributed handler/module registration that remains deferred.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, module, lifecycle, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use handler module registration vocabulary as permission to add handler registries, dynamic handler registration, component discovery, component loading, component module loading, startup hooks, shutdown hooks, dependency containers, runtime endpoints, dashboards, admin console behavior, protocol messages, generated output, persistence, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, session data payloads, component lifecycle payloads, handler registration payloads, module registration payloads, or concrete transport metadata as log-safe in this gate.

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

Handler module registration vocabulary ownership:

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

Rules:

- Documentation and manifests may define handler module registration vocabulary and current source-first mapping.
- `tools/vibit` may later emit a source-first handler module registration map if a follow-up implementation work item authorizes it.
- Existing explicit startup, bootstrap, route registration, payload registry, repository, protocol adapter, and transport close behavior remain unchanged by this gate.
- Handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, protocol payloads, repository interfaces, migrations, generated output, dashboard behavior, admin console behavior, dependencies, hosted surfaces, SDKs, release artifacts, and distributed runtime behavior remain unchanged by this gate.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for runtime components, handlers, services, sessions, routes, RPC, service discovery, groups, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, component lifecycle parity, handler registration parity, handler module registration parity, startup/shutdown parity, or runtime behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- handler module registration behavior;
- handler registration behavior;
- dynamic handler registration;
- component discovery or loading;
- component module loading;
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

The repository check rule is `runtime.pitaya_aligned_handler_module_registration_boundary_gate`.

The check verifies the standard, translation, ADR, change artifacts, manifest references, next-ready state, vocabulary markers, and explicit implementation deferrals.
