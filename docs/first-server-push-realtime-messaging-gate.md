# First Server Push And Realtime Messaging Gate

Status: Accepted v0.1
Last updated: 2026-05-23
Scope: Gate-only boundary for the first server push and realtime messaging vocabulary after storage objects local proof
Depends on: `decisions/ADR-0121-next-alpha-direction-first-server-push-realtime-messaging-gate.md`, `docs/game-protocol.md`, `docs/runtime-protocol-adapter.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/prototype-ready-foundation-execution-plan.md`
Canonical decision: `ADR-0122`

The paired Simplified Chinese translation is `docs/first-server-push-realtime-messaging-gate.zh-CN.md`. The English file is authoritative.

This document defines the first server push and realtime messaging gate. It is a gate artifact. It does not add runtime behavior, transport delivery, protocol routes, Protobuf source, generated output, startup wiring, persistence, migrations, dependencies, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, matchmaking, match runtime, distributed runtime, broad chat/social modules, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The first server push and realtime messaging gate record is:

```yaml
first_server_push_realtime_messaging_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0214
decision: ADR-0122
check_rule: runtime.first_server_push_realtime_messaging_gate
source_next_direction_decision: ADR-0121
source_storage_objects_local_proof_decision: ADR-0120
future_runtime_owner_candidate: runtime/internal/app/realtime
future_runtime_service_source_candidate: runtime/internal/app/realtime/service.go
future_runtime_service_test_candidate: runtime/internal/app/realtime/service_test.go
future_protocol_source_candidate: proto/vibit/realtime/v1/realtime.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/realtime_bridge.go
future_application_handler_candidate: runtime/internal/app/bootstrap/realtime.go
future_transport_delivery_candidate: runtime/internal/platform/transport/ws/outbound.go
future_runtime_slice_work_item: W-0215
future_runtime_slice_direction: first_server_push_realtime_messaging_runtime_slice
first_delivery_model_candidate: single_process_bound_connection_delivery
first_message_intent_vocabulary_recorded: true
first_target_scope_candidates:
  - connection
  - player
  - stream
first_envelope_kind_candidates:
  - event
  - system
websocket_transport_credential_neutral: true
protocol_adapter_payload_mapping_only: true
application_policy_owner_required: true
backend_intent_owner_required: true
realtime_gate_only: true
runtime_behavior_added: false
transport_delivery_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
startup_wiring_added: false
persistence_added: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
matchmaking_added: false
match_runtime_added: false
distributed_runtime_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`W-0212` proved own-player storage object routes in the local alpha WebSocket/Protobuf request flow. `W-0213` selected the first server push and realtime messaging gate as the next prototype-ready direction.

The next useful boundary is not implementation. The next useful boundary is a vocabulary and ownership gate for outbound realtime behavior. Without this gate, future agents could hide server push inside the WebSocket transport, mix message policy with serializers, or copy public route/API shapes from Nakama or Pitaya.

Nakama motivates the product pressure: useful game backends commonly need notifications, streams, chat, presence-adjacent signals, and realtime socket messages after durable storage exists.

Pitaya motivates the architecture pressure: acceptors, sessions, handlers, push, groups, broadcast, backend services, and later cluster/RPC topology must remain separate.

vibit adapts both references with an agent-native boundary:

- transport moves bytes and owns connection mechanics;
- protocol adapters encode, decode, and map payloads;
- application policy decides who may receive outbound messages;
- backend/domain services own message intent and invariants;
- persistence, delivery guarantees, retry behavior, and distributed fanout remain separately gated.

## 3. Ownership

Future first-posture runtime behavior should be application-owned:

```yaml
future_runtime_owner_candidate: runtime/internal/app/realtime
future_runtime_service_source_candidate: runtime/internal/app/realtime/service.go
future_runtime_service_test_candidate: runtime/internal/app/realtime/service_test.go
connection_registry_owner: runtime/internal/app/connection
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
future_application_handler_candidate: runtime/internal/app/bootstrap/realtime.go
```

Rules:

- WebSocket transport must remain credential-neutral and payload-policy-neutral.
- Protocol adapters may map realtime payload bytes and envelope metadata but must not decide recipients or delivery authorization.
- Application-owned realtime behavior must decide recipient targets from validated identity, connection registry state, explicit subscriptions, or future module policy.
- Backend/domain services may emit intent, not perform socket writes directly.
- Domain modules must not import WebSocket transport packages, generated Protobuf packages for delivery, or Pitaya/Nakama SDKs.
- Future implementation must not use direct Nakama/Pitaya API names, route names, or public compatibility shims.

## 4. First Vocabulary

This gate reserves a narrow vibit-native outbound realtime vocabulary:

```yaml
message_intent_kinds:
  - server_notice
  - domain_event_push
  - stream_message
  - presence_signal
recipient_targets:
  - connection_id_and_epoch
  - player_current_connections
  - stream_subscribers
delivery_outcomes:
  - accepted
  - no_active_recipient
  - recipient_not_authorized
  - payload_invalid
  - delivery_unavailable
```

Rules:

- `server_notice` is for server-authored operational or lifecycle notices visible to the client.
- `domain_event_push` is for future server facts emitted by modules, not for client-authored facts.
- `stream_message` is a reserved future vocabulary for stream-targeted delivery. It does not implement chat, channels, rooms, groups, or subscriptions in this gate.
- `presence_signal` is reserved for future presence-adjacent outbound facts. It does not change existing presence lifecycle behavior.
- `connection_id_and_epoch` targets one server-observed connection only.
- `player_current_connections` targets currently active bound connections for one validated player.
- `stream_subscribers` is future vocabulary only until subscription ownership is defined.
- No delivery outcome may reveal raw access tokens, verifier material, credentials, DSNs, storage object values, SQL details, private account data, or hidden recipient existence beyond the authorized public class.

## 5. Future Protocol Shape

The first future protocol source candidate is:

```text
proto/vibit/realtime/v1/realtime.proto
```

The first future generated output candidate is:

```text
runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
```

Candidate future message names:

```yaml
future_messages:
  ServerNotice:
    purpose: server-authored notice payload
  RealtimeEnvelope:
    purpose: stable public payload wrapper for outbound realtime messages
  StreamMessage:
    purpose: future stream-targeted message payload
```

Rules:

- This gate does not add the Protobuf source or generated output.
- The existing `vibit.protocol.v1.Envelope` remains unchanged in this gate.
- Future server-to-client messages should use existing envelope `kind` values such as `event` or `system` unless a later protocol ADR changes envelope semantics.
- Future target scopes may use existing envelope target vocabulary, but connection-specific targeting must remain server-observed and must not trust client-supplied connection ids as authority.
- Future payloads must be vibit-native. They must not copy Nakama notification/channel/stream payloads or Pitaya route payload conventions.

## 6. Future Runtime Flow

The first future runtime slice should preserve this sequence:

```yaml
future_outbound_realtime_flow:
  - backend_or_application_service_creates_server_authored_message_intent
  - application_realtime_service_validates_intent_and_policy
  - application_realtime_service_resolves_allowed_recipients
  - protocol_adapter_maps_intent_to_existing_envelope_and_payload_bytes
  - transport_delivery_adapter_writes_binary_frames_to_server_observed_connections
  - delivery_result_reports_redacted_outcomes
```

Rules:

- Server-authored message intent must be created by trusted runtime code, not accepted blindly from clients.
- Future client-originated chat, stream publish, or room/match messages require separate gates.
- Future implementation may start single-process only.
- Offline inboxes, persistence, acknowledgements, ordering guarantees, retries, backpressure policy, durable stream offsets, and distributed fanout are deferred.
- Route handlers must not write directly to sockets unless a later implementation ADR explicitly defines a transport delivery adapter boundary.

## 7. Identity And Authorization

The first posture is conservative:

```yaml
sender_authority: server_only
client_published_facts_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
connection_id_client_authority_allowed: false
request_token_required_for_client_initiated_realtime_requests: true
```

Rules:

- Server push must be server-authored.
- Client-supplied `player_id`, `session_id`, stream id, room id, match id, or connection id must not grant delivery authority.
- A future client request that asks to publish or subscribe must pass existing route-protection and application policy before any recipient changes.
- Existing request-token protected route behavior remains unchanged by this gate.
- This gate does not change authentication, access-token validation, runtime session persistence, WebSocket handshake authentication, first-message binding, or bound identity route policy.

## 8. Nakama/Pitaya Reference Mapping

Nakama reference mapping:

- This gate moves vibit toward a common game-backend realtime surface: notifications, streams, chat, and presence-adjacent outbound messages.
- It adopts the capability family pressure, not Nakama route paths, REST APIs, runtime API names, payload names, or compatibility promises.

Pitaya reference mapping:

- This gate adopts Pitaya's separation pressure around acceptors, sessions, handlers, push, groups, broadcast, backend services, and cluster/RPC vocabulary.
- It keeps group/broadcast/remote/RPC behavior deferred until single-process delivery and application policy are explicit.
- It does not copy Pitaya handler names, route conventions, package APIs, or cluster topology.

## 9. Verification Expectations

Future implementation slices should include focused tests for:

- recipient target validation;
- metadata-only identity refusal;
- single-process bound connection resolution;
- redacted delivery errors;
- protocol adapter mapping if a realtime payload is added;
- transport delivery behavior if socket writes are added;
- no direct Nakama/Pitaya compatibility markers;
- no storage object value, token, credential, digest, DSN, or transport metadata leakage.

This gate itself is verified by:

- document and translation presence;
- ADR and change spec presence;
- `.arch` manifest status markers;
- `tools/vibit` check coverage;
- no Go runtime behavior, Protobuf source, generated output, migration, dependency, startup wiring, or direct compatibility additions.

## 10. Stop Conditions

Stop and ask before adding any of the following:

- runtime service code;
- WebSocket outbound delivery code;
- Protobuf source or generated output;
- protocol routes or startup registration;
- chat, channels, groups, parties, rooms, matches, matchmaking, or match runtime;
- stream subscription persistence;
- offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, or backpressure mechanisms;
- distributed fanout, frontend/backend split, service discovery, RPC, or cluster groups;
- credential, token, authentication, session, handshake, or route-protection semantic changes;
- repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, or direct Nakama/Pitaya API compatibility.

## 11. Next Work

The next bounded work item is:

```text
W-0215 Implement first server push and realtime messaging runtime slice
```

The next work should implement only the smallest runtime slice authorized by this gate and its implementation ADR. It should keep broad chat/social behavior, protocol expansion, generated output, persistence, delivery guarantees, distributed fanout, matchmaking, match runtime, and direct compatibility behind later bounded work items unless explicitly authorized.
