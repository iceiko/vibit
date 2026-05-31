# Pitaya-Aligned Frontend Backend Role Boundary Gate

Status: Accepted v0.1
Last updated: 2026-05-31
Scope: Gate-only boundary for using Pitaya-aligned frontend/backend role vocabulary after the source-first distributed runtime vocabulary map
Depends on: `decisions/ADR-0155-pitaya-aligned-distributed-runtime-vocabulary-source-first-map.md`, `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0156`

The paired Simplified Chinese translation is `docs/pitaya-aligned-frontend-backend-role-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a frontend/backend role vocabulary gate only. It does not implement frontend/backend server roles, distributed runtime behavior, server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned frontend/backend role boundary gate record is:

```yaml
pitaya_aligned_frontend_backend_role_boundary_gate: defined
completed_work_item: W-0248
decision: ADR-0156
check_rule: runtime.pitaya_aligned_frontend_backend_role_boundary_gate
source_vocabulary_map_decision: ADR-0155
source_vocabulary_map_check_rule: runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map
standard: docs/pitaya-aligned-frontend-backend-role-boundary-gate.md
translation: docs/pitaya-aligned-frontend-backend-role-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: frontend_backend_role_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_frontend_backend_role_vocabulary
future_implementation_work_item: W-0249
future_implementation_direction: implement_pitaya_aligned_frontend_backend_role_source_first_map
allowed_role_vocabulary:
  - frontend_server
  - backend_server
related_vocabulary:
  - acceptor
  - session_binding
  - route_handler
current_single_process_role_mapping:
  frontend_server:
    current: single_process_acceptor_and_dispatch
    future_vocabulary: frontend_server
    implementation_status: deferred_future_architecture_reference
  backend_server:
    current: application_dispatch_and_module_handlers_in_same_process
    future_vocabulary: backend_server
    implementation_status: deferred_future_architecture_reference
  acceptor:
    current: current_websocket_acceptor_single_process_only
    future_vocabulary: frontend_server_acceptor_boundary
    implementation_status: current_single_process_only
  session_binding:
    current: current_first_message_connection_binding_single_process_only
    future_vocabulary: frontend_server_session_ingress_boundary
    implementation_status: current_single_process_only
  route_handler:
    current: current_application_dispatch_and_protocol_bridge
    future_vocabulary: backend_server_handler_boundary
    implementation_status: current_single_process_only
frontend_server_implementation_added: false
backend_server_implementation_added: false
frontend_backend_server_roles_added: false
distributed_runtime_implementation_added: false
server_to_server_rpc_added: false
remote_call_behavior_added: false
service_discovery_added: false
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

`ADR-0155` made the broader Pitaya-aligned distributed runtime vocabulary inspectable through `node tools/vibit inspect pitaya-vocabulary --json`. The next safe step is to define a narrower boundary for the two role terms most likely to cause accidental implementation drift: `frontend_server` and `backend_server`.

This gate makes the role vocabulary usable for future architecture planning only. It does not split the process, add topology, add listeners, add handler remoting, add service discovery, or alter protocol routes.

## 3. Role Vocabulary

Allowed role vocabulary:

- `frontend_server`: a future client-facing role that may own acceptors, session ingress, connection lifecycle, and routing handoff. Current vibit state remains a single-process WebSocket acceptor and application dispatch path.
- `backend_server`: a future service-facing role that may own module handlers behind a frontend role. Current vibit state remains application dispatch and module handlers in the same process.

Related vocabulary:

- `acceptor`: currently the single-process WebSocket acceptor; future planning may associate it with a frontend role.
- `session_binding`: currently first-message connection binding; future planning may associate it with frontend session ingress.
- `route_handler`: currently application dispatch plus protocol bridge; future planning may associate it with backend handler ownership.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, or wire compatibility names from Pitaya.
- Do not use `frontend_server` or `backend_server` as permission to add process topology, runtime behavior, new listeners, RPC/remoting, service discovery, protocol changes, generated output, persistence, or dependencies.
- Do not move module ownership into role vocabulary. Module contracts and vibit ownership manifests remain authoritative.

## 4. Current Mapping

```yaml
current_single_process_role_mapping:
  frontend_server:
    current: single_process_acceptor_and_dispatch
    future_vocabulary: frontend_server
    status: deferred_future_architecture_reference
  backend_server:
    current: application_dispatch_and_module_handlers_in_same_process
    future_vocabulary: backend_server
    status: deferred_future_architecture_reference
  acceptor:
    current: runtime/cmd/vibit-server WebSocket acceptor
    future_role: frontend_server
    status: current_websocket_acceptor_single_process_only
  session_binding:
    current: first-message connection binding
    future_role: frontend_server
    status: current_first_message_connection_binding_single_process_only
  route_handler:
    current: Protobuf bridge plus application dispatch and module handlers
    future_role: backend_server
    status: current_application_dispatch_and_protocol_bridge
```

## 5. Ownership

Role vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-frontend-backend-role-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
```

Rules:

- Documentation and manifests may define role vocabulary and mapping.
- `tools/vibit` may later emit a source-first frontend/backend role map if a follow-up implementation work item authorizes it.
- Runtime, protocol, repository, persistence, generated output, startup wiring, and dependencies remain unchanged by this gate.
- No domain module owns frontend/backend role vocabulary by default.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for role topology, acceptor/session ingress, route-handler placement, RPC/remotes, service discovery, and cluster routing.

Adopted as vocabulary:

- frontend role as future client-facing ingress vocabulary;
- backend role as future service/handler ownership vocabulary;
- frontend/backend mapping to current acceptor, session binding, protocol bridge, application dispatch, and module handler responsibilities.

Adapted to vibit:

- Role vocabulary must stay behind contract-first, module-owned, repository-checkable boundaries.
- Current single-process runtime remains the concrete implementation.
- Any future role split must preserve vibit route contracts, module ownership, generated boundaries, and verification commands.

Rejected for now:

- direct Pitaya API compatibility;
- Pitaya package or route naming compatibility;
- runtime topology changes;
- frontend/backend process split;
- handler remoting, service discovery, distributed groups, broadcast fanout, or cluster-safe routing behavior.

## 7. Future Implementation Work

Open:

```text
M-177/W-0249 Implement Pitaya-aligned frontend/backend role source-first map
```

The future work item may:

- add a source-first repository inspection map for frontend/backend role vocabulary;
- summarize current single-process role mapping;
- update runbooks and acceptance docs to point to the role map;
- add repository checks that verify the role map remains gate-only and redacted.

The future work item must not:

- add frontend/backend server role implementation;
- add distributed runtime implementation;
- add server-to-server RPC or remote call behavior;
- add service discovery;
- add distributed groups, room broadcast fanout, or delivery guarantees;
- add cluster-safe session routing behavior;
- add runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 8. Verification Expectations

This gate should verify:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_frontend_backend_role_boundary_gate
node tools/vibit check change define-pitaya-aligned-frontend-backend-role-boundary-gate --json
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

- process topology changes;
- new goroutines, listeners, network protocols, or server roles;
- RPC/remoting behavior;
- service discovery dependencies;
- distributed group or broadcast behavior;
- cluster-safe session routing behavior;
- protocol or Protobuf changes;
- repository, adapter, migration, or generated-output changes;
- public API compatibility with Pitaya or Nakama;
- any sensitive runtime state exposure.
