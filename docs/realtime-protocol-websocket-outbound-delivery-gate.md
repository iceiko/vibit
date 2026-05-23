# Realtime Protocol And WebSocket Outbound Delivery Gate

Status: Accepted v0.1
Last updated: 2026-05-23
Scope: Gate-only boundary for realtime protocol payloads and WebSocket outbound delivery after the application-owned realtime runtime slice
Depends on: `decisions/ADR-0124-next-alpha-direction-realtime-protocol-websocket-outbound-delivery-gate.md`, `decisions/ADR-0123-first-server-push-realtime-messaging-runtime-slice.md`, `docs/first-server-push-realtime-messaging-gate.md`, `docs/game-protocol.md`, `docs/runtime-protocol-adapter.md`, `docs/generated-output.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0125`

The paired Simplified Chinese translation is `docs/realtime-protocol-websocket-outbound-delivery-gate.zh-CN.md`. The English file is authoritative.

This document defines the realtime protocol and WebSocket outbound delivery gate. It is a gate artifact. It does not add runtime behavior, WebSocket outbound delivery, concrete socket writes, Protobuf source, generated output, protocol bridge behavior, protocol routes, application bootstrap handlers, startup wiring, persistence, migrations, dependencies, authentication/session behavior changes, route-protection changes, hosted deployments, release artifacts, public announcements, paid promotion, stream subscription persistence, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, distributed fanout, matchmaking, match runtime, broad social modules, blob/S3 storage, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The realtime protocol and WebSocket outbound delivery gate record is:

```yaml
realtime_protocol_websocket_outbound_delivery_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0217
decision: ADR-0125
check_rule: runtime.realtime_protocol_websocket_outbound_delivery_gate
source_next_direction_decision: ADR-0124
source_realtime_runtime_slice_decision: ADR-0123
source_realtime_gate_decision: ADR-0122
runtime_intent_owner: runtime/internal/app/realtime
connection_registry_owner: runtime/internal/app/connection
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
future_protocol_source_candidate: proto/vibit/realtime/v1/realtime.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/realtime_bridge.go
future_application_handler_candidate: runtime/internal/app/bootstrap/realtime.go
future_transport_delivery_candidate: runtime/internal/platform/transport/ws/outbound.go
future_implementation_slice_work_item: W-0218
future_implementation_slice_direction: realtime_protocol_websocket_outbound_delivery_implementation_slice
first_protocol_delivery_model_candidate: single_process_server_observed_connection_delivery
first_envelope_kind_candidates:
  - event
  - system
first_payload_family_candidates:
  - server_notice
  - domain_event_push
  - stream_message
  - presence_signal
application_policy_owner_required: true
protocol_adapter_payload_mapping_only: true
websocket_transport_payload_policy_neutral: true
server_observed_connection_id_required: true
server_observed_connection_epoch_required: true
client_connection_id_authority_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
direct_nakama_pitaya_api_compatibility_added: false
runtime_behavior_added: false
websocket_outbound_delivery_added: false
socket_write_added: false
protocol_bridge_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
startup_wiring_added: false
persistence_added: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
route_protection_changed: false
delivery_guarantees_added: false
distributed_runtime_added: false
matchmaking_added: false
match_runtime_added: false
broad_social_module_added: false
```

## 2. Purpose

`W-0215` added an application-owned realtime service under `runtime/internal/app/realtime`. That service validates server-authored outbound intents and resolves allowed active recipients, but it deliberately returns delivery intents instead of writing WebSocket frames.

`W-0216` selected this gate as the next bounded direction. The missing boundary is now the planned handoff from application-owned delivery intent to future client-visible protocol payloads and WebSocket outbound frame delivery.

Without this gate, future work could mix several different responsibilities:

- message authorization and recipient resolution;
- Protobuf payload selection and generated-output ownership;
- protocol envelope mapping;
- WebSocket connection mechanics and socket writes;
- fanout, delivery guarantees, offline storage, and distributed runtime concerns.

This gate separates those concerns before any wire or socket behavior is implemented.

## 3. Ownership

The future outbound path must keep these owners separate:

```yaml
runtime_intent_owner: runtime/internal/app/realtime
connection_registry_owner: runtime/internal/app/connection
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
process_wiring_owner: runtime/cmd/vibit-server
application_bootstrap_owner: runtime/internal/app/bootstrap
```

Rules:

- `runtime/internal/app/realtime` owns outbound intent validation, recipient target validation, and policy-facing delivery outcomes.
- `runtime/internal/app/connection` owns server-observed connection id and epoch state.
- `runtime/internal/platform/protocol/protobuf` may map already-authorized delivery intents to Protobuf payload bytes and existing envelope metadata.
- `runtime/internal/platform/transport/ws` may write already-encoded binary frames to server-observed connections.
- WebSocket transport must not decide recipient authorization, parse domain payloads, or construct domain-specific payloads.
- Protocol adapters must not decide who may receive a message.
- Domain modules must not import the WebSocket transport package or write socket frames directly.
- Future startup wiring must be explicit and must not hide business behavior in process assembly.

## 4. Future Protocol Shape

The first future protocol source candidate remains:

```text
proto/vibit/realtime/v1/realtime.proto
```

The first future generated output candidate remains:

```text
runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
```

The first future protocol bridge candidate is:

```text
runtime/internal/platform/protocol/protobuf/realtime_bridge.go
```

Candidate future payload family:

```yaml
future_payloads:
  ServerNotice:
    intent_kind: server_notice
    purpose: server-authored lifecycle or operational notice visible to an authorized client
  DomainEventPush:
    intent_kind: domain_event_push
    purpose: server-authored module fact prepared for outbound delivery
  StreamMessage:
    intent_kind: stream_message
    purpose: future stream-targeted payload after subscription ownership is defined
  PresenceSignal:
    intent_kind: presence_signal
    purpose: future presence-adjacent outbound signal without changing current presence lifecycle semantics
```

Rules:

- This gate does not add the Protobuf source or generated output.
- The existing `vibit.protocol.v1.Envelope` remains unchanged.
- Future realtime server-to-client payloads should use existing envelope `kind` values such as `event` or `system` unless a later ADR changes envelope semantics.
- Future payloads must be vibit-native. They must not copy Nakama notification, channel, stream, chat, or Pitaya route payload conventions.
- Generated Go Protobuf output must still follow `docs/generated-output.md`; ordinary agents must not hand-edit generated files.
- Connection-specific targeting must rely on server-observed connection id and epoch, not client-supplied authority.

## 5. Future Outbound Flow

The first future implementation slice should preserve this planned flow:

```yaml
future_outbound_delivery_flow:
  - trusted_runtime_code_creates_server_authored_realtime_intent
  - application_realtime_service_validates_intent_and_recipient_policy
  - application_realtime_service_resolves_delivery_intents
  - protocol_adapter_maps_delivery_intents_to_envelope_and_payload_bytes
  - transport_outbound_adapter_writes_binary_frames_to_server_observed_connections
  - delivery_result_reports_redacted_outcomes
```

Rules:

- The application realtime service remains the policy owner.
- Protocol mapping starts only after the application service returns an accepted delivery intent.
- Transport delivery starts only after the protocol adapter returns encoded binary frame bytes.
- Delivery results must be redacted and must not reveal raw credentials, raw tokens, verifier digests, DSNs, SQL details, private account data, or hidden recipient existence beyond the authorized public class.
- A failed write must not mutate domain state unless a later delivery guarantee gate explicitly defines durable delivery semantics.

## 6. WebSocket Transport Delivery Boundary

The future transport delivery candidate is:

```text
runtime/internal/platform/transport/ws/outbound.go
```

Future transport delivery may:

- accept an encoded binary frame;
- target a server-observed connection id and epoch;
- write through the existing WebSocket connection owner;
- report redacted write outcomes to the caller.

Future transport delivery must not:

- choose recipients;
- inspect domain payload semantics;
- validate authentication credentials;
- create Protobuf payloads;
- implement stream subscriptions, chat, groups, broadcast fanout, offline inboxes, delivery guarantees, retries, ordering, backpressure, durable offsets, cluster routing, RPC, service discovery, or direct compatibility shims.

## 7. Identity And Authorization

The first posture remains server-authored and conservative:

```yaml
sender_authority: server_or_admin_validated_identity_only
client_published_facts_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
connection_id_client_authority_allowed: false
route_protection_changed_by_this_gate: false
websocket_handshake_authentication_changed_by_this_gate: false
```

Rules:

- A client-supplied player id, session id, stream id, room id, match id, or connection id must not grant outbound delivery authority.
- Future client-originated publish, subscribe, chat, stream, room, or match messages require separate gates.
- Existing route-protection behavior, first-message connection binding behavior, runtime session behavior, logout behavior, and access-token validation behavior remain unchanged.
- Metadata-only identity remains insufficient proof.

## 8. Nakama/Pitaya Reference Mapping

Nakama reference mapping:

- This gate moves vibit toward client-visible outbound realtime capability needed by notifications, streams, chat, and presence-adjacent features.
- It adopts the capability pressure only. It does not copy Nakama public APIs, route paths, runtime helper names, payload names, or compatibility promises.

Pitaya reference mapping:

- This gate preserves Pitaya-style separation among acceptor/transport mechanics, sessions and connection state, handlers, protocol serialization, backend service intent, push/group/broadcast vocabulary, and later cluster concerns.
- It keeps group membership, broadcast fanout, remote calls, RPC, service discovery, and frontend/backend role separation deferred.
- It does not copy Pitaya route conventions, package APIs, handler naming, or cluster topology.

## 9. Verification Expectations

Future implementation slices should include focused tests for:

- protocol payload mapping from accepted realtime delivery intents;
- rejection of wrong intent kinds, target kinds, and malformed payloads;
- binary frame write handoff to server-observed connection id and epoch;
- stale epoch or missing connection handling;
- transport write error redaction;
- application policy staying outside protocol and transport adapters;
- no domain module imports of WebSocket transport packages;
- no generated output hand edits;
- no direct Nakama/Pitaya compatibility markers.

This gate itself is verified by:

- English and Simplified Chinese standard presence;
- ADR and change spec presence;
- `.arch` manifest status markers;
- `runtime.realtime_protocol_websocket_outbound_delivery_gate` check coverage;
- no Go runtime behavior, Protobuf source, generated output, protocol bridge, WebSocket outbound writer, startup wiring, migration, dependency, delivery guarantee, distributed runtime, or direct compatibility additions.

## 10. Stop Conditions

Stop and ask before adding any of the following outside a later explicitly authorized implementation work item:

- `proto/vibit/realtime/v1/realtime.proto`;
- `runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go`;
- `runtime/internal/platform/protocol/protobuf/realtime_bridge.go`;
- `runtime/internal/platform/transport/ws/outbound.go`;
- `runtime/internal/app/bootstrap/realtime.go`;
- startup wiring for outbound delivery;
- concrete socket writes;
- protocol routes or client publish routes;
- stream subscription persistence;
- chat, channels, groups, parties, rooms, matches, matchmaking, or match runtime;
- offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, or backpressure mechanisms;
- distributed fanout, frontend/backend split, service discovery, RPC, or cluster groups;
- credential, token, authentication, session, handshake, or route-protection semantic changes;
- repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, or direct Nakama/Pitaya API compatibility.

## 11. Next Work

The next bounded work item is:

```text
W-0218 Implement realtime protocol and WebSocket outbound delivery slice
```

The next work should implement only the smallest slice authorized by this gate and a matching implementation ADR. It should keep stream subscription ownership, chat semantics, groups, broadcast fanout, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, distributed fanout, matchmaking, match runtime, and direct compatibility behind later bounded work items unless explicitly authorized.
