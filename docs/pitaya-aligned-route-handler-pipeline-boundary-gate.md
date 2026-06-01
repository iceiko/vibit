# Pitaya-Aligned Route Handler Pipeline Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned route handler, handler pipeline, serializer, and message forwarding vocabulary after the cluster-safe session routing source-first map
Depends on: `decisions/ADR-0166-select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map.md`, `decisions/ADR-0165-pitaya-aligned-cluster-safe-session-routing-source-first-map.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0167`

The paired Simplified Chinese translation is `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a route handler pipeline vocabulary gate only. It does not implement route handlers, handler routing behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, distributed runtime behavior, distributed groups, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC, remote calls, frontend/backend server roles, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned route handler pipeline boundary gate record is:

```yaml
pitaya_aligned_route_handler_pipeline_boundary_gate: defined
completed_work_item: W-0259
decision: ADR-0167
check_rule: runtime.pitaya_aligned_route_handler_pipeline_boundary_gate
previous_direction_decision: ADR-0166
cluster_safe_session_routing_source_first_map_decision: ADR-0165
cluster_safe_session_routing_source_first_map_check_rule: runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map
standard: docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md
translation: docs/pitaya-aligned-route-handler-pipeline-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: route_handler_pipeline_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_route_handler_pipeline_vocabulary
future_implementation_work_item: W-0260
future_implementation_direction: implement_pitaya_aligned_route_handler_pipeline_source_first_map
allowed_route_handler_pipeline_vocabulary:
  - route_handler
  - route_key
  - handler_dispatch
  - handler_pipeline
  - pipeline_step
  - serializer_boundary
  - message_forwarding
  - route_target
related_vocabulary:
  - protocol_envelope
  - route_request
  - application_dispatch
  - command_handler
  - query_handler
  - protocol_bridge
  - target_scope
  - frontend_server
  - backend_server
  - server_to_server_rpc
  - remote_call
  - service_discovery
  - cluster_safe_session_routing
current_single_process_route_handler_mapping:
  protocol_envelope:
    current: kind_module_name_structured_routing
    future_vocabulary: route_key
    implementation_status: current_protocol_adapter_owned_shape
  route_request:
    current: app_route_request_handoff
    future_vocabulary: route_handler
    implementation_status: current_application_dispatch
  application_dispatch:
    current: explicit_command_query_dispatch
    future_vocabulary: handler_dispatch
    implementation_status: no_pitaya_handler_pipeline
  transactional_dispatch:
    current: application_unit_of_work_wrapper
    future_vocabulary: pipeline_step
    implementation_status: current_vibit_transaction_boundary_only
  protocol_bridge:
    current: explicit_generated_payload_bridge
    future_vocabulary: serializer_boundary
    implementation_status: no_pluggable_serializer_behavior
  outbound_message:
    current: server_push_intent_to_protocol_envelope
    future_vocabulary: message_forwarding
    implementation_status: no_cross_node_forwarding
route_handler_implementation_added: false
handler_routing_behavior_added: false
handler_pipeline_behavior_added: false
pipeline_middleware_behavior_added: false
serializer_behavior_added: false
message_forwarding_behavior_added: false
backend_route_targeting_added: false
cluster_safe_session_routing_added: false
session_location_registry_added: false
connection_owner_node_registry_added: false
routing_epoch_behavior_added: false
session_route_target_added: false
remote_connection_handoff_added: false
reconnect_route_added: false
distributed_session_routing_added: false
distributed_group_implementation_added: false
distributed_groups_added: false
group_membership_registry_added: false
room_broadcast_fanout_added: false
broadcast_delivery_guarantee_added: false
stream_subscription_added: false
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

`ADR-0166` selected route handler pipeline vocabulary as the next Pitaya-aligned direction after the cluster-safe session routing source-first map. This vocabulary is useful because Pitaya-style distributed systems separate client-facing route receipt, handler dispatch, handler pipelines, serialization, and forwarding decisions.

The risk is that agents may treat those words as permission to replace vibit's current protocol adapter, application dispatch, generated payload bridges, transaction wrapper, or server-push path. This gate records vocabulary and mapping only. It keeps vibit's concrete single-process WebSocket Protobuf route flow unchanged.

## 3. Route Handler Pipeline Vocabulary

Allowed route handler pipeline vocabulary:

- `route_handler`: future planning vocabulary for the application-facing handler selected by a route key. It is not new handler code in this slice.
- `route_key`: future planning vocabulary for the logical route identity. Current route identity remains structured `kind`, `module`, and `name` fields.
- `handler_dispatch`: future planning vocabulary for selecting a handler. Current dispatch remains the explicit application dispatcher.
- `handler_pipeline`: future planning vocabulary for ordered pre-handler or post-handler processing. It is not middleware behavior in this slice.
- `pipeline_step`: future planning vocabulary for a bounded pipeline unit. Current transactional dispatch is not generalized middleware.
- `serializer_boundary`: future planning vocabulary for encode/decode ownership. Current Protobuf bridge functions remain the only concrete serializer boundary.
- `message_forwarding`: future planning vocabulary for forwarding a message to another owner or node. Current runtime has no cross-node forwarding.
- `route_target`: future planning vocabulary for handler placement or target selection. Current target scope metadata is not backend route targeting.

Related vocabulary:

- `protocol_envelope`, `route_request`, `application_dispatch`, `command_handler`, `query_handler`, `protocol_bridge`, and `target_scope`: current vibit route-flow concepts.
- `frontend_server`, `backend_server`, `server_to_server_rpc`, `remote_call`, `service_discovery`, and `cluster_safe_session_routing`: prior Pitaya-aligned vocabulary families that remain implementation deferrals.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, pipeline, serializer, forwarding, registry, selector, or configuration compatibility names from Pitaya or Nakama.
- Do not use route handler pipeline vocabulary as permission to add handler routing behavior, middleware chains, serializer plugins, forwarding workers, backend route targeting, service discovery, RPC, remote calls, protocol messages, generated output, persistence, dependencies, topology, or distributed runtime behavior.
- Do not move domain behavior into transport, Protobuf adapters, serializer boundaries, or process startup.
- Do not bypass application dispatch, request/session validation, bound identity route policy, generated output rules, redaction rules, or module ownership.

## 4. Current Mapping

```yaml
current_single_process_route_handler_mapping:
  protocol_envelope:
    current: kind/module/name structured routing fields
    future_vocabulary: route_key
    status: current_protocol_adapter_owned_shape
  route_request:
    current: explicit application route request handoff
    future_vocabulary: route_handler
    status: current_application_dispatch
  application_dispatch:
    current: explicit command/query route registration and dispatch
    future_vocabulary: handler_dispatch
    status: no_pitaya_handler_pipeline
  transactional_dispatch:
    current: application-owned unit-of-work wrapper for commands
    future_vocabulary: pipeline_step
    status: current_vibit_transaction_boundary_only
  protocol_bridge:
    current: explicit generated Protobuf payload bridge functions
    future_vocabulary: serializer_boundary
    status: no_pluggable_serializer_behavior
  outbound_message:
    current: server-push intent converted to protocol envelope
    future_vocabulary: message_forwarding
    status: no_cross_node_forwarding
```

## 5. Ownership

Route handler pipeline vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
transport_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
dependency_owner: unchanged
```

Rules:

- Documentation and manifests may define route handler pipeline vocabulary and current mapping.
- `tools/vibit` may later emit a source-first route handler pipeline map if a follow-up implementation work item authorizes it.
- Runtime, transport, protocol, repository, persistence, generated output, startup wiring, dependencies, service discovery, RPC, remote calls, frontend/backend role behavior, cluster-safe session routing, distributed group behavior, and room broadcast behavior remain unchanged by this gate.
- Domain modules do not gain handler pipeline, serializer, forwarding, backend route targeting, service discovery, RPC, or transport ownership by default.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for acceptors, sessions, route handlers, frontend/backend roles, RPC/remotes, service discovery, groups, broadcast, cluster routing, handler pipelines, serializers, and forwarding.

Adopted as vocabulary:

- route handler pipeline vocabulary for future architecture planning;
- route key, handler dispatch, handler pipeline, pipeline step, serializer boundary, message forwarding, and route target as planning vocabulary;
- mapping vibit's existing protocol adapter, route request, application dispatcher, transactional dispatch, protocol bridges, and server-push intent to deferred Pitaya-aligned concepts.

Adapted to vibit:

- Current route identity remains structured `kind`, `module`, and `name` fields.
- Current dispatch remains application-owned and explicit.
- Current serialization remains Protobuf adapter owned through explicit bridge functions.
- Current server push remains single-process and does not imply cross-node forwarding.
- Any future route handler pipeline implementation must be separately gated and verified before behavior exists.

Rejected for now:

- direct Pitaya or Nakama API compatibility;
- Pitaya or Nakama package, method, route, handler, serializer, pipeline, forwarding, registry, selector, or configuration naming compatibility;
- route handler implementation, handler routing behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, or backend route targeting;
- service discovery, RPC, remote calls, frontend/backend process split, cluster-safe session routing, distributed groups, room broadcast fanout, or distributed runtime behavior;
- protocol messages or routes, generated output, persistence, migrations, dependencies, hosted deployment, SDK publication, or release artifacts for route handler pipelines.

## 7. Future Implementation Work

Open:

```text
M-188/W-0260 Implement Pitaya-aligned route handler pipeline source-first map
```

The future work item may:

- add a source-first repository inspection map for route handler pipeline vocabulary;
- summarize current protocol envelope, route request, application dispatch, transactional dispatch, protocol bridge, and outbound message mappings;
- update runbooks and acceptance docs to point to the route handler pipeline map;
- add repository checks that verify the map remains gate-only and redacted.

The future work item must not:

- add route handler implementation;
- add handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, or backend route targeting;
- add cluster-safe session routing behavior, session location registries, connection owner node registries, routing epochs, session route targets, remote connection handoff, reconnect routing, or distributed session routing;
- add service discovery implementation, service registries, selectors, node identity, or topology behavior;
- add server-to-server RPC implementation or remote call behavior;
- add frontend/backend server role implementation;
- add distributed runtime implementation;
- add distributed group implementation, group membership registries, stream subscriptions, room broadcast fanout, or delivery guarantees;
- add runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 8. Verification Expectations

This gate should verify:

- `runtime.pitaya_aligned_route_handler_pipeline_boundary_gate` is registered.
- `ADR-0167`, this standard, the Simplified Chinese translation, change artifacts, and conversation memory exist.
- W-0259 is completed and W-0260 is next-ready.
- Current route dispatch and protocol adapter mapping is recorded.
- Deferrals remain explicit for route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, cluster-safe session routing, distributed runtime, service discovery, RPC, remote calls, protocol, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility.

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_route_handler_pipeline_boundary_gate
node tools/vibit check change define-pitaya-aligned-route-handler-pipeline-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## 9. Stop Conditions

Stop and require a new bounded work item before:

- implementing route handlers, handler routing behavior, handler pipelines, pipeline middleware, serializer behavior, message forwarding, or backend route targeting;
- changing route identity, protocol envelope shape, Protobuf sources, generated output, application dispatch semantics, transaction behavior, protocol bridge behavior, or outbound delivery behavior;
- adding service discovery, RPC, remote calls, frontend/backend server role behavior, cluster-safe session routing, distributed groups, room broadcast fanout, or distributed runtime behavior;
- adding dependencies, migrations, repository interfaces, PostgreSQL adapters, hosted deployment surfaces, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.
