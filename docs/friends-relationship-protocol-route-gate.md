# Friends Relationship Protocol Route Gate

Status: Accepted v0.1
Last updated: 2026-05-26
Scope: Gate-only boundary for future client-facing friends relationship protocol routes after application runtime behavior
Depends on: `docs/friends-relationship-runtime-behavior-gate.md`, `decisions/ADR-0147-friends-relationship-runtime-behavior-implementation.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/generated-output.md`, `docs/bound-identity-route-policy-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0148`

The paired Simplified Chinese translation is `docs/friends-relationship-protocol-route-gate.zh-CN.md`. The English file is authoritative.

This document defines the friends relationship protocol route gate. It is a gate artifact. It does not add protocol route implementation, Protobuf source, generated output, startup wiring, runtime handlers, repository interface changes, PostgreSQL adapter changes, migration changes, dependencies, authentication/session behavior changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, event/audit tables, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The friends relationship protocol route gate record is:

```yaml
friends_relationship_protocol_route_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0240
decision: ADR-0148
check_rule: runtime.friends_relationship_protocol_route_gate
source_runtime_behavior_implementation_decision: ADR-0147
source_runtime_behavior_implementation: runtime/internal/app/friends/service.go
source_runtime_behavior_tests: runtime/internal/app/friends/service_test.go
source_runtime_behavior_gate_decision: ADR-0146
source_repository_interface_decision: ADR-0143
repository_interface: runtime/internal/modules/friends.Repository
future_protocol_source_candidate: proto/vibit/friends/v1/friends.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/friends_bridge.go
future_protocol_bridge_test_candidate: runtime/internal/platform/protocol/protobuf/friends_bridge_test.go
future_application_handler_candidate: runtime/internal/app/bootstrap/friends.go
future_application_handler_test_candidate: runtime/internal/app/bootstrap/friends_test.go
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_actor_id_allowed_as_proof: false
client_supplied_target_player_id_allowed: true
first_actor_kind: player
first_payload_package: vibit.friends.v1
protobuf_envelope_change_status: unchanged
websocket_transport_credential_neutral: true
protocol_route_gate_only: true
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
runtime_handler_added: false
startup_wiring_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
delivery_guarantees_added: false
stream_subscription_added: false
chat_added: false
groups_added: false
parties_added: false
matchmaking_added: false
match_runtime_added: false
event_audit_table_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_protocol_route_implementation_work_item: W-0241
future_protocol_route_implementation_direction: implement_friends_relationship_protocol_route
```

## 2. Purpose

`W-0239` implemented application-owned friends relationship behavior under `runtime/internal/app/friends`. The next useful boundary is not route code or `.proto` generation. The next useful boundary is a protocol route gate that records how future WebSocket/Protobuf exposure should call that service without moving friends behavior into transport, generated files, or persistence adapters.

Nakama motivates the product surface: friends, friend requests, blocks, and actor-relative relationship status are common game backend social graph capabilities. vibit should cover that capability class.

Pitaya motivates the architecture posture: acceptors, sessions, route handlers, serializers, and backend services should remain separated. vibit adapts that by keeping WebSocket transport credential-neutral, keeping Protobuf payload bridging explicit, and invoking application-owned route handlers that call application-owned friends services.

This gate records the future protocol shape before implementation:

- candidate route names;
- candidate request and response message shapes;
- route protection and identity handoff posture;
- protocol adapter, application handler, and startup ownership;
- generated-output expectations;
- public error mapping and redaction expectations;
- Nakama/Pitaya reference mapping;
- stop conditions that keep implementation and generated artifacts out of this slice.

## 3. Future Route Surface

The first route family should expose player-to-player friends relationship operations for the validated actor only:

```yaml
candidate_routes:
  - kind: command
    module: friends
    name: SendFriendRequest
    route_id: friends.SendFriendRequest
    service_method: SendFriendRequest
  - kind: command
    module: friends
    name: AcceptFriendRequest
    route_id: friends.AcceptFriendRequest
    service_method: AcceptFriendRequest
  - kind: command
    module: friends
    name: RejectFriendRequest
    route_id: friends.RejectFriendRequest
    service_method: RejectFriendRequest
  - kind: command
    module: friends
    name: RemoveFriend
    route_id: friends.RemoveFriend
    service_method: RemoveFriend
  - kind: command
    module: friends
    name: BlockPlayer
    route_id: friends.BlockPlayer
    service_method: BlockPlayer
  - kind: command
    module: friends
    name: UnblockPlayer
    route_id: friends.UnblockPlayer
    service_method: UnblockPlayer
  - kind: query
    module: friends
    name: ListFriendRelationships
    route_id: friends.ListFriendRelationships
    service_method: ListFriendRelationships
  - kind: query
    module: friends
    name: GetFriendRelationshipStatus
    route_id: friends.GetFriendRelationshipStatus
    service_method: GetFriendRelationshipStatus
```

Rules:

- The route names must stay vibit-native and must not copy Nakama route paths or Pitaya route naming.
- Send, accept, reject, remove, block, and unblock routes are commands.
- List and status routes are queries.
- The first route family is only for the validated player actor. It must not expose arbitrary actor ids or another player's private social graph.
- Client payloads may provide a target player id, expected relationship version, status filter, list limit, and pagination cursor where the service already has vocabulary for them.
- Groups, parties, chat rooms, broadcast fanout, presence subscriptions, matchmaking filters, match social context, admin social graph inspection, account merge behavior, event/audit streams, and script hooks remain deferred.
- Future route implementation must register routes explicitly. No catch-all friends route or reflective handler is allowed.

## 4. Protocol Shape

The first friends relationship protocol source candidate is:

```text
proto/vibit/friends/v1/friends.proto
```

The first generated output candidate is:

```text
runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go
```

The first Protobuf package candidate is:

```text
vibit.friends.v1
```

Candidate messages:

```yaml
messages:
  FriendRelationship:
    fields:
      relationship_id: string
      player_low_id: string
      player_high_id: string
      requested_by_player_id: string
      lifecycle_state: string
      public_status: string
      version: int64
      created_at: string
      updated_at: string
  FriendRelationshipPage:
    fields:
      relationships: repeated FriendRelationship
      next_pair_token: string
  SendFriendRequestRequest:
    fields:
      target_player_id: string
  SendFriendRequestResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  AcceptFriendRequestRequest:
    fields:
      target_player_id: string
      expected_version: int64
  AcceptFriendRequestResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  RejectFriendRequestRequest:
    fields:
      target_player_id: string
      expected_version: int64
  RejectFriendRequestResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  RemoveFriendRequest:
    fields:
      target_player_id: string
      expected_version: int64
  RemoveFriendResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  BlockPlayerRequest:
    fields:
      target_player_id: string
      expected_version: int64
  BlockPlayerResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  UnblockPlayerRequest:
    fields:
      target_player_id: string
      expected_version: int64
  UnblockPlayerResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  ListFriendRelationshipsRequest:
    fields:
      status: string
      limit: int32
      after_pair_token: string
  ListFriendRelationshipsResponse:
    fields:
      page: FriendRelationshipPage
  GetFriendRelationshipStatusRequest:
    fields:
      target_player_id: string
  GetFriendRelationshipStatusResponse:
    fields:
      public_status: string
      relationship: FriendRelationship
      version: int64
```

Rules:

- The existing `proto/vibit/protocol/v1/envelope.proto` must remain unchanged unless a later protocol ADR explicitly changes envelope semantics.
- Time values should use RFC3339Nano UTC text when exposed.
- Optional `expected_version` mapping must preserve the service's optional expected-version vocabulary. Future implementation must make absence vs `0` semantics explicit in tests.
- Public status values must map to the service's actor-relative statuses: `none`, `outgoing_request_pending`, `incoming_request_pending`, `friends`, `blocked_by_actor`, `blocked_actor`, `mutual_blocked`, `removed`, and `rejected`.
- Response `status` values must map to service operation outcomes such as `sent`, `accepted`, `request_rejected`, `removed`, `blocked`, `unblocked`, `listed`, `found`, or `rejected`.
- The protocol shape must not include client-supplied actor id, raw access tokens, credential material, lookup digests, verifier digests, SQL details, private transport metadata, chat payloads, group ids, party ids, matchmaking fields, match runtime fields, or direct external API compatibility markers.
- Relationship ids, player ids, relationship versions, block state, and lifecycle state are not log-safe by default.

## 5. Route Protection And Identity Handoff

The first route-policy posture is:

```yaml
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_actor_id_allowed_as_proof: false
client_supplied_target_player_id_allowed: true
```

Rules:

- Future friends routes must be protected gameplay routes.
- Future handlers must receive a validated `app.RequestIdentity` from the existing protected-route flow.
- Metadata-only `player_id` or `session_id` from envelope/session metadata must never become friends actor proof.
- Client payloads must not choose actor ids.
- Client payloads may choose the target player id for player-to-player operations.
- The service remains responsible for rejecting invalid identity before id generation or repository access.
- This gate does not change authentication, token validation, session persistence, first-message binding, WebSocket handshake behavior, bound-identity policy, or route-protection semantics.

## 6. Future Route Flow

Future implementation must preserve this sequence:

```yaml
future_route_flow:
  - receive WebSocket/Protobuf envelope through existing request path
  - apply protected-route authenticated wrapper policy
  - obtain validated authenticated request identity
  - decode vibit.friends.v1 request payload
  - reject payload actor proof and derive actor only from request identity
  - map payload fields to runtime/internal/app/friends service request
  - call application-owned friends service
  - map service result to vibit.friends.v1 response payload
  - map service public errors to protocol error responses
  - keep transport, Protobuf bridge, application handler, service, repository, and PostgreSQL adapter ownership separated
```

Rules:

- WebSocket transport remains credential-neutral.
- The Protobuf bridge should map payloads and response shapes only; it must not own friends relationship behavior.
- Application bootstrap handlers should own route registration and service invocation.
- The application service remains the owner of identity checks, validation handoff, repository conflict mapping, and actor-relative public status.
- PostgreSQL adapters remain persistence-only.

## 7. Public Error Mapping

Future route implementation should map service public errors without leaking internal details:

```yaml
public_error_mapping:
  FRIENDSHIP_INVALID_REQUEST: invalid_request
  FRIENDSHIP_UNAUTHENTICATED: unauthenticated
  FRIENDSHIP_FORBIDDEN: forbidden
  FRIENDSHIP_TARGET_NOT_FOUND: target_not_found
  FRIENDSHIP_RELATIONSHIP_NOT_FOUND: relationship_not_found
  FRIENDSHIP_DUPLICATE_REQUEST: duplicate_request
  FRIENDSHIP_ALREADY_FRIENDS: already_friends
  FRIENDSHIP_BLOCKED_RELATIONSHIP: blocked_relationship
  FRIENDSHIP_INVALID_TRANSITION: invalid_transition
  FRIENDSHIP_VERSION_MISMATCH: version_mismatch
  FRIENDSHIP_UNAVAILABLE: unavailable
```

Rules:

- Public protocol errors may expose only stable public codes/classes and retryability posture if a later ADR authorizes that field.
- Internal repository errors, SQL details, target existence probes beyond the service public error, relationship ids, player ids, access-token material, credential material, lookup digests, verifier digests, and transport metadata must remain out of default logs and error messages.
- Authentication and route protection failures must use existing protected-route semantics. This gate does not invent a new proof carrier.

## 8. Generated Output Posture

Future generated output must follow `docs/generated-output.md`.

Rules:

- `proto/vibit/friends/v1/friends.proto` may be added only by a later implementation slice.
- `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go` may be added only as generated output from Buf/protoc.
- Generated Go output must contain the `protoc-gen-go` generated-code marker and trace to the source `.proto`.
- Agents must not hand-edit generated Go Protobuf files.
- This gate does not change `buf.yaml`, `buf.gen.yaml`, or generated output.

## 9. Nakama Reference Mapping

Nakama reference mapping:

```yaml
nakama_reference_mapping:
  capability_family: friends_groups_and_parties
  mapped_capabilities:
    - friend_request_send
    - friend_request_accept
    - friend_request_reject
    - friend_remove
    - player_block
    - player_unblock
    - list_relationships
    - get_relationship_status
  direct_api_compatibility: false
```

Nakama informs the useful capability class. vibit does not copy Nakama route paths, field names, permission semantics, runtime script APIs, or public API compatibility.

## 10. Pitaya Reference Mapping

Pitaya reference mapping:

```yaml
pitaya_reference_mapping:
  architecture_pressure:
    - acceptor_session_handler_separation
    - serializer_adapter_separation
    - backend_service_boundary
  distributed_architecture_status: deferred
  direct_api_compatibility: false
```

Pitaya informs layering pressure. This gate does not add Pitaya-style distributed topology, frontend/backend split, RPC, groups, service discovery, distributed social graph routing, or direct Pitaya API compatibility.

## 11. Required Future Tests

Future implementation tests should cover:

- route registration for all eight route ids;
- command/query kind mapping;
- Protobuf request and response bridge mapping;
- optional expected-version mapping;
- actor derivation from validated request identity;
- rejection of metadata-only identity through the existing protected-route wrapper;
- rejection of client-supplied actor proof;
- target player id validation and self-targeting behavior;
- public status mapping for none, pending incoming/outgoing, friends, blocked, removed, and rejected states;
- service public error to protocol error mapping;
- redaction of private relationship, player, token, credential, SQL, and transport details;
- no route behavior in WebSocket transport or PostgreSQL adapter packages;
- generated-output traceability if Protobuf source is added.

## 12. Stop Conditions

Stop and create a separate work item before adding:

- protocol route implementation;
- `proto/vibit/friends/v1/friends.proto`;
- generated Go Protobuf output;
- protocol bridge implementation;
- application bootstrap handlers;
- startup route registration;
- new dependencies;
- migration changes;
- repository interface changes;
- PostgreSQL adapter changes;
- authentication/session behavior changes;
- event/audit tables;
- delivery guarantees;
- stream subscriptions;
- chat rooms;
- groups;
- parties;
- broadcast fanout;
- matchmaking;
- match runtime;
- operations/admin behavior;
- SDK publication;
- generated client libraries;
- hosted deployments;
- release artifacts;
- public announcements;
- paid promotion;
- Pitaya-style distributed architecture;
- direct Nakama/Pitaya API compatibility.

## 13. Verification

The repository check rule for this gate is:

```text
runtime.friends_relationship_protocol_route_gate
```

Recommended verification after this gate:

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.friends_relationship_protocol_route_gate
node tools/vibit check change define-friends-relationship-protocol-route-gate --json
node tools/vibit check module friends --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Go tests are not required by this gate because no Go runtime behavior is added, but a full runtime test run remains useful before closing a development turn.
