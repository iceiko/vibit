# Pitaya-Aligned Service Discovery Boundary Gate

Status: Accepted v0.1
Last updated: 2026-05-31
Scope: Gate-only boundary for using Pitaya-aligned service discovery vocabulary after the server-to-server RPC source-first map
Depends on: `decisions/ADR-0159-pitaya-aligned-server-to-server-rpc-source-first-map.md`, `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md`, `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0160`

The paired Simplified Chinese translation is `docs/pitaya-aligned-service-discovery-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a service discovery vocabulary gate only. It does not implement service discovery, service registries, service selectors, node registries, server identity, server-to-server RPC, remote calls, frontend/backend server roles, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned service discovery boundary gate record is:

```yaml
pitaya_aligned_service_discovery_boundary_gate: defined
completed_work_item: W-0252
decision: ADR-0160
check_rule: runtime.pitaya_aligned_service_discovery_boundary_gate
rpc_source_first_map_decision: ADR-0159
rpc_source_first_map_check_rule: runtime.pitaya_aligned_server_to_server_rpc_source_first_map
standard: docs/pitaya-aligned-service-discovery-boundary-gate.md
translation: docs/pitaya-aligned-service-discovery-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: service_discovery_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_service_discovery_vocabulary
future_implementation_work_item: W-0253
future_implementation_direction: implement_pitaya_aligned_service_discovery_source_first_map
allowed_service_discovery_vocabulary:
  - service_discovery
  - service_registry
  - service_instance
  - service_selector
related_vocabulary:
  - frontend_server
  - backend_server
  - server_to_server_rpc
  - remote_call
  - route_handler
  - module_handler
  - static_process_composition
current_single_process_service_discovery_mapping:
  service_discovery:
    current: no_service_discovery_current_static_single_process_composition
    future_vocabulary: service_discovery
    implementation_status: deferred_future_architecture_reference
  service_registry:
    current: no_registry_current_startup_composition
    future_vocabulary: service_registry
    implementation_status: deferred_future_architecture_reference
  service_instance:
    current: single_process_runtime_components_not_network_instances
    future_vocabulary: service_instance
    implementation_status: deferred_future_architecture_reference
  service_selector:
    current: no_selector_current_direct_in_process_dispatch
    future_vocabulary: service_selector
    implementation_status: deferred_future_architecture_reference
  route_handler:
    current: current_application_dispatch_and_protocol_bridge
    future_vocabulary: discoverable_backend_route_handler
    implementation_status: current_single_process_only
  module_handler:
    current: current_module_handler_in_process_function_call
    future_vocabulary: discoverable_backend_module_handler
    implementation_status: current_single_process_only
service_discovery_implementation_added: false
service_registry_added: false
service_selector_added: false
node_registry_added: false
server_identity_added: false
server_to_server_rpc_implementation_added: false
remote_call_behavior_added: false
frontend_server_implementation_added: false
backend_server_implementation_added: false
frontend_backend_server_roles_added: false
distributed_runtime_implementation_added: false
distributed_groups_added: false
room_broadcast_fanout_added: false
cluster_safe_session_routing_added: false
runtime_behavior_added: false
runtime_endpoint_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
hosted_deployment_added: false
sdk_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0159` made server-to-server RPC and remote-call vocabulary inspectable through `node tools/vibit inspect pitaya-rpc --json`. The next Pitaya-aligned concept that can accidentally imply implementation is service discovery. Distributed runtimes need a way to find services or server instances, but vibit currently uses one statically composed process.

This gate records service discovery vocabulary and current mapping before any implementation. It does not add registry storage, registry clients, selection algorithms, node identity, heartbeat behavior, routing tables, process topology, network listeners, RPC transports, or dependencies.

## 3. Service Discovery Vocabulary

Allowed service discovery vocabulary:

- `service_discovery`: a future architecture-planning family for locating server-side capabilities or instances.
- `service_registry`: a future record of available server-side capabilities or instances. It is not a database table, external registry, or in-memory runtime structure in this slice.
- `service_instance`: a future representation of a discoverable server-side instance. It is not a current runtime process identity.
- `service_selector`: a future selection concept for choosing a service instance. It is not a load balancer, routing algorithm, or retry policy in this slice.

Related vocabulary:

- `frontend_server` and `backend_server`: role vocabulary from the prior Pitaya-aligned role boundary. These remain future architecture vocabulary only.
- `server_to_server_rpc` and `remote_call`: RPC vocabulary from the prior boundary. Service discovery does not authorize RPC implementation.
- `route_handler` and `module_handler`: current in-process application/module handlers that may be mapped to future discoverable backend ownership.
- `static_process_composition`: the current concrete vibit implementation model.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, registry, or configuration compatibility names from Pitaya.
- Do not use service discovery vocabulary as permission to add service registries, selectors, heartbeats, node identity, runtime topology, registry storage, external discovery dependencies, RPC transports, remote calls, endpoint behavior, protocol changes, generated output, persistence, or dependencies.
- Do not use future service discovery vocabulary to bypass module contracts, application dispatch boundaries, authentication/session validation gates, permission checks, generated output rules, redaction rules, or repository ownership.

## 4. Current Mapping

```yaml
current_single_process_service_discovery_mapping:
  service_discovery:
    current: none; composition is static and single-process
    future_vocabulary: service_discovery
    status: deferred_future_architecture_reference
  service_registry:
    current: no registry; startup wires concrete handlers directly
    future_vocabulary: service_registry
    status: deferred_future_architecture_reference
  service_instance:
    current: runtime components are not discoverable network instances
    future_vocabulary: service_instance
    status: deferred_future_architecture_reference
  service_selector:
    current: no selector; dispatch is direct and in-process
    future_vocabulary: service_selector
    status: deferred_future_architecture_reference
  route_handler:
    current: Protobuf bridge plus application dispatch and module handlers
    future_vocabulary: discoverable backend route handler
    status: current_application_dispatch_and_protocol_bridge
  module_handler:
    current: runtime/internal/modules handwritten behavior in one process
    future_vocabulary: discoverable backend module handler
    status: current_module_handler_in_process_function_call
```

## 5. Ownership

Service discovery vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-service-discovery-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
dependency_owner: unchanged
```

Rules:

- Documentation and manifests may define service discovery vocabulary and mapping.
- `tools/vibit` may later emit a source-first service discovery map if a follow-up implementation work item authorizes it.
- Runtime, protocol, repository, persistence, generated output, startup wiring, dependencies, RPC, remote calls, and frontend/backend role behavior remain unchanged by this gate.
- Domain modules do not gain service discovery ownership by default. Module contracts remain the source of module behavior and data ownership.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for frontend/backend roles, route handler placement, RPC/remotes, service discovery, groups, broadcast, and cluster routing.

Adopted as vocabulary:

- service discovery as future architecture-planning vocabulary;
- service registry, service instance, and service selector as future discovery vocabulary;
- mapping current static startup composition to deferred service discovery concepts;
- mapping current route and module handlers to possible future discoverable backend ownership.

Adapted to vibit:

- Any future discovery model must preserve vibit module ownership, application dispatch boundaries, server-authoritative validation, generated output rules, redaction, and repository checks.
- Current single-process runtime remains the concrete implementation.
- Any future service discovery implementation must be separately gated and verified before behavior exists.

Rejected for now:

- direct Pitaya API compatibility;
- Pitaya package, method, registry, or route naming compatibility;
- service discovery implementation;
- service registry or selector behavior;
- node identity, heartbeats, or runtime topology;
- server-to-server RPC or remote call behavior;
- frontend/backend process split;
- protocol messages or routes for discovery.

## 7. Future Implementation Work

Open:

```text
M-181/W-0253 Implement Pitaya-aligned service discovery source-first map
```

The future work item may:

- add a source-first repository inspection map for service discovery vocabulary;
- summarize current static single-process composition and direct handler dispatch;
- update runbooks and acceptance docs to point to the service discovery map;
- add repository checks that verify the service discovery map remains gate-only and redacted.

The future work item must not:

- add service discovery implementation;
- add service registry or selector behavior;
- add node identity, heartbeat, membership, or topology behavior;
- add server-to-server RPC implementation;
- add remote call behavior;
- add frontend/backend server role implementation;
- add distributed runtime implementation;
- add distributed groups, room broadcast fanout, or delivery guarantees;
- add cluster-safe session routing behavior;
- add runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 8. Verification Expectations

This gate should verify:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_service_discovery_boundary_gate
node tools/vibit check change define-pitaya-aligned-service-discovery-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

This gate does not require Go tests by itself because it adds no Go runtime behavior.

## 9. Stop Conditions

Stop and create a separate gate if the work requires:

- service discovery implementation;
- service registry, selector, membership, heartbeat, or node identity behavior;
- server-to-server RPC or remote call behavior;
- process topology changes;
- new goroutines, listeners, network protocols, or server roles;
- distributed group or broadcast behavior;
- cluster-safe session routing behavior;
- protocol or Protobuf changes;
- repository, adapter, migration, dependency, or generated-output changes;
- public API compatibility with Pitaya or Nakama;
- any sensitive runtime state exposure.
