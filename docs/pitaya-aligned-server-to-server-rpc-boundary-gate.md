# Pitaya-Aligned Server To Server RPC Boundary Gate

Status: Accepted v0.1
Last updated: 2026-05-31
Scope: Gate-only boundary for using Pitaya-aligned server-to-server RPC and remote-call vocabulary after the frontend/backend role source-first map
Depends on: `decisions/ADR-0157-pitaya-aligned-frontend-backend-role-source-first-map.md`, `docs/pitaya-aligned-frontend-backend-role-boundary-gate.md`, `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0158`

The paired Simplified Chinese translation is `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a server-to-server RPC vocabulary gate only. It does not implement server-to-server RPC, remote calls, service discovery, frontend/backend server roles, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned server-to-server RPC boundary gate record is:

```yaml
pitaya_aligned_server_to_server_rpc_boundary_gate: defined
completed_work_item: W-0250
decision: ADR-0158
check_rule: runtime.pitaya_aligned_server_to_server_rpc_boundary_gate
role_source_first_map_decision: ADR-0157
role_source_first_map_check_rule: runtime.pitaya_aligned_frontend_backend_role_source_first_map
standard: docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md
translation: docs/pitaya-aligned-server-to-server-rpc-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: server_to_server_rpc_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_server_to_server_rpc_vocabulary
future_implementation_work_item: W-0251
future_implementation_direction: implement_pitaya_aligned_server_to_server_rpc_source_first_map
allowed_rpc_vocabulary:
  - server_to_server_rpc
  - remote_call
related_vocabulary:
  - route_handler
  - module_handler
  - application_dispatch
  - service_discovery
current_single_process_rpc_mapping:
  server_to_server_rpc:
    current: no_rpc_current_single_process_application_dispatch
    future_vocabulary: server_to_server_rpc
    implementation_status: deferred_future_architecture_reference
  remote_call:
    current: no_remote_call_current_in_process_module_invocation
    future_vocabulary: remote_call
    implementation_status: deferred_future_architecture_reference
  route_handler:
    current: current_application_dispatch_and_protocol_bridge
    future_vocabulary: backend_server_route_handler_boundary
    implementation_status: current_single_process_only
  module_handler:
    current: current_module_handler_in_process_function_call
    future_vocabulary: backend_server_module_handler_boundary
    implementation_status: current_single_process_only
  service_discovery:
    current: no_service_discovery_current_static_single_process_composition
    future_vocabulary: service_discovery
    implementation_status: deferred_future_architecture_reference
server_to_server_rpc_implementation_added: false
remote_call_behavior_added: false
service_discovery_added: false
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

`ADR-0157` made frontend/backend role vocabulary inspectable through `node tools/vibit inspect pitaya-roles --json`. The next Pitaya-aligned concept likely to cause implementation drift is server-to-server RPC. Pitaya uses distributed server calls as part of its architecture; vibit can adopt the vocabulary for planning only after the boundary is explicit.

This gate records the vocabulary and current mapping before any implementation. It does not add RPC transports, remote handlers, service registries, node identity, frontend/backend process topology, protocol carriers, or client-visible routes.

## 3. RPC Vocabulary

Allowed RPC vocabulary:

- `server_to_server_rpc`: a future server-internal call family for architecture planning. If implemented later, it must not bypass vibit module contracts, route contracts, authorization boundaries, identity/session validation, generated boundaries, or repository checks.
- `remote_call`: a future distributed invocation vocabulary item. It is distinct from client protocol commands, queries, events, and WebSocket routes.

Related vocabulary:

- `route_handler`: currently application dispatch plus protocol bridge; future planning may associate it with backend handler ownership.
- `module_handler`: currently in-process handwritten module behavior; future planning may associate it with internal service handling, but module ownership remains authoritative.
- `application_dispatch`: current vibit in-process command/query dispatch path.
- `service_discovery`: a future dependency-sensitive architecture vocabulary item, not a current implementation.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, or wire compatibility names from Pitaya.
- Do not use `server_to_server_rpc` or `remote_call` as permission to add RPC transports, remote calls, service discovery, node registry, distributed process topology, new endpoint behavior, protocol changes, generated output, persistence, or dependencies.
- Do not use RPC vocabulary to bypass module contracts, application dispatch boundaries, authentication/session validation gates, permission checks, generated output rules, or repository ownership.

## 4. Current Mapping

```yaml
current_single_process_rpc_mapping:
  server_to_server_rpc:
    current: no server-to-server RPC; application dispatch is in-process
    future_vocabulary: server_to_server_rpc
    status: deferred_future_architecture_reference
  remote_call:
    current: no remote call; module handlers are invoked in-process
    future_vocabulary: remote_call
    status: deferred_future_architecture_reference
  route_handler:
    current: Protobuf bridge plus application dispatch and module handlers
    future_vocabulary: backend_server route handler boundary
    status: current_application_dispatch_and_protocol_bridge
  module_handler:
    current: runtime/internal/modules handwritten behavior in one process
    future_vocabulary: backend_server module handler boundary
    status: current_module_handler_in_process_function_call
  service_discovery:
    current: none; composition is static and single-process
    future_vocabulary: service_discovery
    status: deferred_future_architecture_reference
```

## 5. Ownership

RPC vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md
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

- Documentation and manifests may define RPC vocabulary and mapping.
- `tools/vibit` may later emit a source-first RPC map if a follow-up implementation work item authorizes it.
- Runtime, protocol, repository, persistence, generated output, startup wiring, dependencies, and service discovery remain unchanged by this gate.
- Domain modules do not gain RPC ownership by default. Module contracts remain the source of module behavior and data ownership.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for frontend/backend roles, route handler placement, RPC/remotes, service discovery, groups, broadcast, and cluster routing.

Adopted as vocabulary:

- server-to-server RPC as future architecture-planning vocabulary;
- remote calls as future distributed invocation vocabulary;
- service discovery as future dependency-sensitive vocabulary;
- mapping current in-process dispatch and module handlers to deferred RPC concepts.

Adapted to vibit:

- Any future RPC must preserve vibit module ownership, application dispatch boundaries, server-authoritative validation, generated output rules, redaction, and repository checks.
- Current single-process runtime remains the concrete implementation.
- Any future RPC implementation must be separately gated and verified before behavior exists.

Rejected for now:

- direct Pitaya API compatibility;
- Pitaya package, method, or route naming compatibility;
- server-to-server RPC implementation;
- remote call behavior;
- service discovery;
- distributed runtime behavior;
- frontend/backend process split;
- protocol messages or routes for RPC.

## 7. Future Implementation Work

Open:

```text
M-179/W-0251 Implement Pitaya-aligned server-to-server RPC source-first map
```

The future work item may:

- add a source-first repository inspection map for server-to-server RPC and remote-call vocabulary;
- summarize current single-process dispatch and module handler mapping;
- update runbooks and acceptance docs to point to the RPC map;
- add repository checks that verify the RPC map remains gate-only and redacted.

The future work item must not:

- add server-to-server RPC implementation;
- add remote call behavior;
- add service discovery;
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
node tools/vibit inspect rule runtime.pitaya_aligned_server_to_server_rpc_boundary_gate
node tools/vibit check change define-pitaya-aligned-server-to-server-rpc-boundary-gate --json
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

- RPC transport behavior;
- remote call behavior;
- service discovery or node registry behavior;
- process topology changes;
- new goroutines, listeners, network protocols, or server roles;
- distributed group or broadcast behavior;
- cluster-safe session routing behavior;
- protocol or Protobuf changes;
- repository, adapter, migration, dependency, or generated-output changes;
- public API compatibility with Pitaya or Nakama;
- any sensitive runtime state exposure.
