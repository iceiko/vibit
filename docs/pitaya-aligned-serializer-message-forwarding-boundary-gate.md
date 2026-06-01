# Pitaya-Aligned Serializer And Message Forwarding Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned serializer and message forwarding vocabulary after the route handler pipeline source-first map
Depends on: `decisions/ADR-0169-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map.md`, `decisions/ADR-0168-pitaya-aligned-route-handler-pipeline-source-first-map.md`, `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0170`

The paired Simplified Chinese translation is `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a serializer and message forwarding vocabulary gate only. It does not implement serializer behavior, message forwarding behavior, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, cluster-safe session routing behavior, distributed session routing, distributed runtime behavior, distributed groups, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC, remote calls, frontend/backend server roles, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned serializer and message forwarding boundary gate record is:

```yaml
pitaya_aligned_serializer_message_forwarding_boundary_gate: defined
completed_work_item: W-0262
decision: ADR-0170
check_rule: runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate
previous_direction_decision: ADR-0169
route_handler_pipeline_source_first_map_decision: ADR-0168
route_handler_pipeline_source_first_map_check_rule: runtime.pitaya_aligned_route_handler_pipeline_source_first_map
standard: docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md
translation: docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: serializer_message_forwarding_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_serializer_message_forwarding_vocabulary
future_implementation_work_item: W-0263
future_implementation_direction: implement_pitaya_aligned_serializer_message_forwarding_source_first_map
allowed_serializer_message_forwarding_vocabulary:
  - serializer_boundary
  - serializer_format
  - encode_boundary
  - decode_boundary
  - message_forwarding
  - forwarding_target
  - forwarding_envelope
  - delivery_handoff
related_vocabulary:
  - protocol_bridge
  - generated_payload_bridge
  - outbound_message
  - target_scope
  - route_target
  - route_handler
  - backend_server
  - server_to_server_rpc
  - remote_call
  - service_discovery
  - cluster_safe_session_routing
current_single_process_serializer_forwarding_mapping:
  protocol_bridge:
    current: explicit_generated_payload_bridge
    future_vocabulary: serializer_boundary
    implementation_status: no_pluggable_serializer_behavior
  envelope_encoding:
    current: protobuf_envelope_owned_by_protocol_adapter
    future_vocabulary: serializer_format
    implementation_status: no_serializer_registry
  payload_encoding:
    current: generated_payload_bridge_functions
    future_vocabulary: encode_boundary
    implementation_status: no_custom_encode_pipeline
  payload_decoding:
    current: generated_payload_bridge_functions
    future_vocabulary: decode_boundary
    implementation_status: no_custom_decode_pipeline
  outbound_message:
    current: server_push_intent_to_protocol_envelope
    future_vocabulary: message_forwarding
    implementation_status: no_cross_node_forwarding
  target_scope:
    current: metadata_only_target_scope
    future_vocabulary: forwarding_target
    implementation_status: no_backend_route_targeting
  forwarding_envelope:
    current: no_internal_forwarding_envelope
    future_vocabulary: forwarding_envelope
    implementation_status: not_implemented
  delivery_handoff:
    current: single_process_websocket_delivery
    future_vocabulary: delivery_handoff
    implementation_status: no_remote_delivery_handoff
serializer_behavior_added: false
message_forwarding_behavior_added: false
route_handler_implementation_added: false
handler_routing_behavior_added: false
handler_pipeline_behavior_added: false
pipeline_middleware_behavior_added: false
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

`ADR-0169` selected serializer and message forwarding vocabulary as the next Pitaya-aligned direction after the route handler pipeline source-first map.

The risk is that agents may treat serializer and forwarding terms as permission to add pluggable codecs, protocol shape changes, forwarding workers, backend route targeting, service discovery, RPC, remote delivery, or distributed runtime behavior. This gate records vocabulary and mapping only. It keeps vibit's concrete single-process WebSocket Protobuf flow unchanged.

## 3. Vocabulary

Allowed serializer and message forwarding vocabulary:

- `serializer_boundary`: future planning vocabulary for encode/decode ownership. Current Protobuf bridge functions remain the concrete boundary.
- `serializer_format`: future planning vocabulary for wire format selection. Current envelope format remains the existing Protobuf envelope.
- `encode_boundary`: future planning vocabulary for outbound payload encoding. Current generated bridge functions remain explicit.
- `decode_boundary`: future planning vocabulary for inbound payload decoding. Current generated bridge functions remain explicit.
- `message_forwarding`: future planning vocabulary for forwarding a message to another owner or node. Current runtime has no cross-node forwarding.
- `forwarding_target`: future planning vocabulary for the owner selected for forwarding. Current target scope metadata is not backend route targeting.
- `forwarding_envelope`: future planning vocabulary for an internal forwarding wrapper. No internal forwarding envelope exists in this slice.
- `delivery_handoff`: future planning vocabulary for handing delivery to another runtime owner. Current delivery remains single-process WebSocket delivery.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, serializer, forwarding, registry, selector, or configuration compatibility names from Pitaya or Nakama.
- Do not use serializer or forwarding vocabulary as permission to add codecs, serializer registries, forwarding workers, backend route targeting, service discovery, RPC, remote calls, protocol messages, generated output, persistence, dependencies, topology, or distributed runtime behavior.
- Do not move domain behavior into transport, Protobuf adapters, serializer boundaries, forwarding layers, or process startup.

## 4. Current Mapping

```yaml
current_single_process_serializer_forwarding_mapping:
  protocol_bridge:
    current: explicit generated Protobuf payload bridge functions
    future_vocabulary: serializer_boundary
    status: no_pluggable_serializer_behavior
  envelope_encoding:
    current: Protobuf envelope owned by the protocol adapter
    future_vocabulary: serializer_format
    status: no_serializer_registry
  payload_encoding:
    current: generated payload bridge functions
    future_vocabulary: encode_boundary
    status: no_custom_encode_pipeline
  payload_decoding:
    current: generated payload bridge functions
    future_vocabulary: decode_boundary
    status: no_custom_decode_pipeline
  outbound_message:
    current: server-push intent converted to protocol envelope
    future_vocabulary: message_forwarding
    status: no_cross_node_forwarding
  target_scope:
    current: metadata-only target scope
    future_vocabulary: forwarding_target
    status: no_backend_route_targeting
  forwarding_envelope:
    current: no internal forwarding envelope
    future_vocabulary: forwarding_envelope
    status: not_implemented
  delivery_handoff:
    current: single-process WebSocket delivery
    future_vocabulary: delivery_handoff
    status: no_remote_delivery_handoff
```

## 5. Ownership

Serializer and message forwarding vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md
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

- Documentation and manifests may define serializer and message forwarding vocabulary and current mapping.
- `tools/vibit` may later emit a source-first serializer and message forwarding map if a follow-up implementation work item authorizes it.
- Runtime, transport, protocol, repository, persistence, generated output, startup wiring, dependencies, service discovery, RPC, remote calls, frontend/backend role behavior, cluster-safe session routing, distributed group behavior, and room broadcast behavior remain unchanged by this gate.
- Domain modules do not gain serializer, forwarding, backend route targeting, service discovery, RPC, transport, or distributed runtime ownership by default.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for acceptors, sessions, route handlers, frontend/backend roles, RPC/remotes, service discovery, groups, broadcast, cluster routing, handler pipelines, serializers, and forwarding.

Adopted as vocabulary:

- serializer boundary and serializer format vocabulary for future architecture planning;
- encode boundary and decode boundary vocabulary for future payload handoff planning;
- message forwarding, forwarding target, forwarding envelope, and delivery handoff vocabulary for future distributed routing planning.

Adapted to vibit:

- Current serialization remains Protobuf adapter owned through explicit bridge functions.
- Current server push remains single-process and does not imply cross-node forwarding.
- Current target scope metadata does not imply backend route targeting.
- Any future serializer or forwarding implementation must be separately gated and verified before behavior exists.

Rejected for now:

- direct Pitaya or Nakama API compatibility;
- Pitaya or Nakama package, method, route, handler, serializer, forwarding, registry, selector, or configuration naming compatibility;
- serializer behavior, message forwarding behavior, forwarding workers, backend route targeting, service discovery, RPC, remote calls, protocol changes, generated output, persistence, migrations, dependencies, hosted deployment, SDK publication, or release artifacts.

## 7. Future Implementation Work

Open:

```text
M-191/W-0263 Implement Pitaya-aligned serializer and message forwarding source-first map
```

The future work item may:

- add a source-first repository inspection map for serializer and message forwarding vocabulary;
- summarize current protocol bridge, generated payload bridge, outbound message, target-scope metadata, and delivery handoff mappings;
- update runbooks and acceptance docs to point to the serializer and message forwarding map;
- add repository checks that verify the map remains gate-only and redacted.

The future work item must not:

- add serializer behavior or message forwarding behavior;
- add route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, or backend route targeting;
- add service discovery implementation, service registries, selectors, node identity, topology behavior, RPC, remote calls, or distributed runtime behavior;
- add protocol messages or routes, Protobuf source, generated output, persistence, migrations, dependencies, hosted deployment, SDK publication, or direct Nakama/Pitaya API compatibility.

## 8. Verification

Required checks:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate
node tools/vibit check change define-pitaya-aligned-serializer-message-forwarding-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Stop if any check fails.
