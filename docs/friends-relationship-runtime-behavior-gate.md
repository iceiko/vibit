# Friends Relationship Runtime Behavior Gate

Status: Accepted v0.1
Last updated: 2026-05-26
Scope: Gate-only boundary for future application-owned friends relationship runtime behavior after the PostgreSQL adapter
Depends on: `docs/friends-relationship-lifecycle-gate.md`, `docs/friends-relationship-repository-boundary.md`, `docs/friends-relationship-postgresql-adapter-gate.md`, `runtime/internal/modules/friends/repository.go`, `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`, `docs/runtime-protocol-adapter.md`, `docs/bound-identity-route-policy-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0146`

The paired Simplified Chinese translation is `docs/friends-relationship-runtime-behavior-gate.zh-CN.md`. The English file is authoritative.

This document defines the friends relationship runtime behavior gate. It is a gate artifact. It does not add runtime behavior implementation, runtime handlers, startup wiring, protocol routes, Protobuf source, generated output, repository interface changes, PostgreSQL adapter changes, migration changes, dependencies, authentication/session behavior changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, event/audit tables, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The friends relationship runtime behavior gate record is:

```yaml
friends_relationship_runtime_behavior_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0238
decision: ADR-0146
check_rule: runtime.friends_relationship_runtime_behavior_gate
source_postgresql_adapter_implementation_decision: ADR-0145
source_postgresql_adapter: runtime/internal/platform/persistence/postgres/friend_relationship_repository.go
source_repository_interface_decision: ADR-0143
repository_interface: runtime/internal/modules/friends.Repository
repository_interface_source: runtime/internal/modules/friends/repository.go
future_runtime_owner_candidate: runtime/internal/app
future_friends_application_package_candidate: runtime/internal/app/friends
future_runtime_service_source_candidate: runtime/internal/app/friends/service.go
future_runtime_service_test_candidate: runtime/internal/app/friends/service_test.go
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
actor_identity_source: validated_request_identity_player_id
first_actor_kind: player
actor_relative_public_status_required: true
route_policy_requirement: request_token_required
service_application_owner: runtime/internal/app
repository_handoff: unit_of_work_friend_relationship_repository_factory
unit_of_work_handoff_required: true
runtime_behavior_gate_only: true
runtime_behavior_added: false
runtime_handlers_added: false
startup_wiring_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
event_audit_table_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_runtime_behavior_implementation_work_item: W-0239
future_runtime_behavior_implementation_direction: implement_friends_relationship_runtime_behavior
```

## 2. Purpose

`W-0237` implemented the PostgreSQL adapter for `runtime/internal/modules/friends.Repository`. The next useful boundary is not a route or protocol change. The next useful boundary is the runtime behavior gate that defines how application code may later turn a validated player request into friends repository operations.

This gate records the future behavior shape before implementation:

- application ownership for the service;
- actor identity derivation from validated request identity;
- actor-relative status and list behavior;
- permission and route-policy posture;
- validation and conflict mapping expectations;
- unit-of-work and repository handoff;
- redaction rules;
- test expectations;
- stop conditions that keep protocol, generated output, authentication/session changes, event/audit tables, and broader social features out of this slice.

Nakama motivates durable friends relationships as a core social graph capability. Pitaya motivates keeping handlers, sessions, route context, and persistence responsibilities separated. vibit adapts those references through explicit application-owned behavior and checks, not direct public API compatibility.

## 3. Ownership

Future runtime behavior is application-owned:

```yaml
future_runtime_owner_candidate: runtime/internal/app
future_friends_application_package_candidate: runtime/internal/app/friends
future_runtime_service_source_candidate: runtime/internal/app/friends/service.go
future_runtime_service_test_candidate: runtime/internal/app/friends/service_test.go
repository_interface_owner: runtime/internal/modules/friends
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
player_account_owner: runtime/internal/app/player
```

Rules:

- Future service behavior may live under `runtime/internal/app/friends` or an equivalent application-owned package ratified by the implementation slice.
- The service may call `runtime/internal/modules/friends.Repository` only through application or unit-of-work dependencies.
- The service may check target player account existence and account state only through existing application-owned player repository capabilities when the implementation slice authorizes that dependency handoff.
- The service must not import PostgreSQL adapter packages, SQL row types, migration packages, WebSocket transport packages, generated Protobuf packages, chat packages, group packages, party packages, matchmaking packages, match runtime packages, SDK packages, or distributed runtime packages.
- The friends module remains the owner of storage-neutral value types, normalizers, lifecycle/status vocabulary, and repository error vocabulary.
- The PostgreSQL adapter remains persistence-only and must not derive request identity, route policy, or public protocol errors.
- Transport and protocol adapters must not own friends relationship permission or business behavior.

## 4. Request Identity And Actor Derivation

The first runtime behavior posture is player-to-player:

```yaml
first_actor_kind: player
actor_identity_source: validated_request_identity_player_id
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_actor_id_allowed_as_proof: false
client_supplied_target_player_id_allowed: true
self_relationship_allowed: false
```

Rules:

- A future friends operation must derive the actor from a validated `app.RequestIdentity`.
- `RequestIdentity.Status` must be `validated`.
- `RequestIdentity.ActorKind` must be `player`.
- `RequestIdentity.PlayerIDValidated` must be true.
- `RequestIdentity.PlayerID` must be non-empty and must match the actor identity when both are present.
- Metadata-only `player_id` from envelope/session metadata must never satisfy this gate.
- A persisted `session_id` alone must never become proof.
- Client payloads may name the target player for player-to-player operations, but must not name the actor as proof.
- Self-targeting must fail before repository mutation when possible.

This gate does not change `RequestIdentity`, access-token validation, bound connection identity, durable runtime session validation, or WebSocket handshake behavior. It only records the identity requirements that future friends behavior must require before repository access.

## 5. Future Runtime Behavior Shape

The future first implementation may expose an application service with these candidate operations:

```yaml
candidate_operations:
  - send_friend_request
  - accept_friend_request
  - reject_friend_request
  - remove_friend
  - block_player
  - unblock_player
  - list_friend_relationships
  - get_friend_relationship_status
```

Recommended first posture:

- `send_friend_request` creates an outgoing pending relationship from the validated actor to a target player.
- `accept_friend_request` accepts only an incoming pending relationship.
- `reject_friend_request` rejects only an incoming pending relationship.
- `remove_friend` removes an existing friendship or pending relationship according to repository lifecycle rules, but must not remove an actor-specific block.
- `block_player` sets the actor-specific block state and prevents ordinary friend operations while blocked.
- `unblock_player` clears only the actor's block state and must not restore a prior friendship.
- `list_friend_relationships` lists actor-scoped relationships and computes actor-relative public statuses.
- `get_friend_relationship_status` reads one actor-target relationship and computes actor-relative public status.

Rules:

- Runtime behavior must use server-derived actor identity.
- Runtime behavior must validate target player id, self-targeting, expected version, list status filters, pagination cursors, and list limit before repository calls when possible.
- Runtime behavior must compute public status relative to the requesting actor.
- Runtime behavior must not expose another player's private social graph, arbitrary social graph search, admin inspection, group/guild relationships, party memberships, chat rooms, match rooms, matchmaking filters, server script hooks, or direct external API compatibility in the first implementation.
- Runtime behavior must not add public protocol routes or generated output unless a later protocol gate authorizes them.

## 6. Candidate Application Service Shape

The first implementation slice should define a small application-owned service boundary. Candidate inputs and outputs:

```yaml
candidate_request_fields:
  - request_identity
  - target_player_id
  - expected_relationship_version
  - status_filter
  - list_limit
  - after_relationship_cursor
candidate_result_fields:
  - relationship
  - relationships
  - actor_relative_public_status
  - next_relationship_cursor
  - relationship_version
  - public_error_code
```

Rules:

- The service should accept already-normalized application identity, not raw tokens, cookies, headers, WebSocket subprotocol values, or envelope proof carriers.
- The service should call friends module normalizers before repository handoff.
- The service should keep actor ids, target player ids, relationship ids, relationship state, block state, and relationship versions out of default errors and logs.
- The service should expose stable public error codes or classes for runtime handlers to map later.
- The service should not add route registration, Protobuf conversion, startup composition, or command/query dispatch wiring in the gate slice.

## 7. Validation Rules

Future runtime behavior must enforce validation before persistence:

```yaml
validation_required:
  - validated_player_identity
  - target_player_id_non_empty_and_normalized
  - self_target_forbidden
  - target_player_lookup_or_repository_conflict_handled
  - actor_must_be_pair_member_for_existing_relationship
  - expected_relationship_version_positive_when_present
  - list_limit_bounded
  - pagination_cursor_bounded
  - status_filter_known
```

Rules:

- Target player id validation should reuse existing player identity rules or friends module normalizers where available.
- Self-targeting must fail before repository mutation when possible.
- Missing expected version behavior must be explicit in implementation tests.
- Invalid input must fail before repository mutation when possible.
- Target player not found, relationship not found, blocked relationship, and privacy-sensitive cases must not leak hidden social graph state.
- Repository unavailable errors must remain redacted.

## 8. Permission And Route Policy Posture

The first route-policy posture is conservative:

```yaml
route_policy_requirement: request_token_required
public_friends_routes_allowed: false
bound_connection_required_by_this_gate: false
session_validated_required_by_this_gate: false
bound_session_required_by_this_gate: false
```

Candidate permission families for later public contracts:

- send friend request;
- accept incoming friend request;
- reject incoming friend request;
- remove friend;
- block player;
- unblock player;
- list own friend relationships;
- read own actor-relative relationship status.

Rules:

- Friends relationship routes must be protected routes.
- The first posture should use the existing `request_token_required` route protection family unless a later route-policy ADR changes named routes.
- Public routes must not read or mutate friends relationships.
- Bound connection identity and durable session validation may remain available for future route families, but this gate does not require them or change ordinary protected route behavior.
- Metadata-only identity must fail closed.

## 9. Conflict And Error Mapping

Future runtime behavior must map friends repository errors into stable application classes:

```yaml
candidate_public_error_classes:
  - FRIENDSHIP_INVALID_REQUEST
  - FRIENDSHIP_UNAUTHENTICATED
  - FRIENDSHIP_FORBIDDEN
  - FRIENDSHIP_TARGET_NOT_FOUND
  - FRIENDSHIP_RELATIONSHIP_NOT_FOUND
  - FRIENDSHIP_DUPLICATE_REQUEST
  - FRIENDSHIP_ALREADY_FRIENDS
  - FRIENDSHIP_BLOCKED_RELATIONSHIP
  - FRIENDSHIP_INVALID_TRANSITION
  - FRIENDSHIP_VERSION_MISMATCH
  - FRIENDSHIP_UNAVAILABLE
```

Rules:

- Missing or unvalidated request identity must map to the existing protected-route authentication posture before repository access.
- Self-targeting and malformed inputs may be public invalid-request classes.
- Duplicate pending request, already friends, blocked relationship, invalid transition, and version mismatch may be public conflict classes only when they do not reveal private graph state beyond the actor's own relationship.
- Target-player-not-found and relationship-not-found leakage rules must be conservative; public output may collapse privacy-sensitive cases to not-found or forbidden.
- Stored private relationship state, block details, actor ids, target ids, relationship ids, SQL details, driver errors, DSNs, credentials, token material, verifier digests, and route proof carriers must not leak.
- Repository `storage_unavailable` errors must map to an unavailable class without exposing platform internals.
- Runtime behavior must not add authentication/token/session failure detail beyond existing application route-protection classes.

## 10. Unit-Of-Work And Repository Handoff

Future runtime behavior should use the existing application transaction boundary:

```yaml
unit_of_work_handoff_required: true
repository_handoff: unit_of_work_friend_relationship_repository_factory
repository_interface: runtime/internal/modules/friends.Repository
postgresql_adapter: runtime/internal/platform/persistence/postgres.FriendRelationshipRepository
service_starts_transactions: false
repository_starts_transactions: false
```

Rules:

- State-changing operations must run through an application-owned unit-of-work or equivalent transaction boundary.
- The service should obtain `friends.Repository` from the unit-of-work rather than constructing the PostgreSQL adapter directly.
- The service must not issue `BEGIN`, `COMMIT`, or `ROLLBACK` SQL.
- The service must not import PostgreSQL-specific packages.
- Read-only operations may use repository dependencies directly or through a read-only unit-of-work if a later implementation gate ratifies that shape.
- Target player lookup and relationship mutation must be ordered deterministically inside the future unit-of-work when both are required.
- Repository errors must be mapped after rollback or failed unit-of-work outcomes without leaking platform internals.

## 11. Actor-Relative Status Rules

Public status is actor-relative:

```yaml
actor_relative_public_states:
  - none
  - outgoing_request_pending
  - incoming_request_pending
  - friends
  - blocked_by_actor
  - blocked_actor
  - mutual_blocked
  - removed
  - rejected
```

Rules:

- Persisted pair-oriented state must be translated relative to the requesting actor before public output.
- Pending requests must distinguish outgoing from incoming for the actor.
- Actor-specific block columns must distinguish actor blocked target, target blocked actor, and mutual block.
- `unblock_player` must not restore friendship or pending request state.
- Removed and rejected states must not expose another player's private relationship history unless a later public contract explicitly authorizes the output.
- List output must remain player-scoped and must not reveal relationships between other players.

## 12. Redaction And Logging

Friends relationship runtime data is private social graph data.

```yaml
private_social_graph_log_safe: false
actor_player_id_log_safe: false_by_default
target_player_id_log_safe: false_by_default
relationship_id_log_safe: conditional_after_validation
relationship_state_log_safe: conditional_after_validation
relationship_version_log_safe: conditional_after_validation
forbidden_runtime_log_material:
  - raw_access_token
  - raw_credential
  - credential_lookup_digest
  - credential_verifier_digest
  - token_lookup_digest
  - token_verifier_digest
  - verifier_key
  - verifier_key_id
  - authorization_header
  - cookie
  - query_string_token
  - websocket_subprotocol
  - websocket_connection_id
  - remote_address
  - sql_text
  - database_dsn
  - private_relationship_graph
  - chat_room_id
  - group_id
  - party_id
  - match_id
  - pitaya_server_id
```

Rules:

- Application errors must be redacted and typed.
- Raw player relationship details, target identifiers, block details, and storage driver errors must not be logged by default.
- Private social graph data must be treated as non-log-safe unless a later redaction policy narrows that.
- Runtime behavior must not store or return authentication material, token material, verifier digests, transport metadata, chat room state, group state, party state, match state, or distributed routing state.

## 13. Future Test Expectations

The later implementation slice should add focused application service tests:

```yaml
future_tests:
  - service_requires_validated_player_identity
  - service_rejects_metadata_only_identity
  - service_rejects_client_supplied_actor_as_proof
  - send_friend_request_validates_target_and_self_relationship
  - send_friend_request_maps_duplicate_already_friends_blocked_conflicts
  - accept_reject_require_incoming_pending_request
  - remove_friend_preserves_block_semantics
  - block_unblock_use_actor_specific_block_state
  - list_relationships_is_actor_scoped_and_status_filtered
  - get_status_computes_actor_relative_public_status
  - expected_version_mismatch_maps_to_public_conflict
  - repository_unavailable_is_redacted
  - state_changing_operations_use_unit_of_work
  - no_protocol_or_transport_dependency
```

Rules:

- Tests must use fakes or in-memory stubs for application dependencies unless a later implementation slice authorizes live PostgreSQL integration.
- Tests must verify that missing or metadata-only identity fails before repository mutation.
- Tests must verify actor-relative status conversion for pending, friends, actor block, target block, mutual block, removed, and rejected states.
- Tests must verify that target player not found, relationship not found, and privacy-sensitive failures do not leak hidden graph details.
- Tests must not require protocol routes, WebSocket transport, Protobuf generated code, or generated clients.
- Tests must not print raw credentials, tokens, verifier keys, digests, DSNs, query strings, authorization headers, cookies, player ids, or private relationship state.

## 14. Relationship To Runtime, Protocol, And Authentication

This gate does not change runtime or protocol behavior:

```yaml
runtime_friends_service_added: false
runtime_friends_handlers_added: false
friends_protocol_routes_added: false
protobuf_friends_messages_added: false
generated_friends_output_added: false
authentication_session_behavior_changed: false
request_identity_handoff_changed: false
```

Rules:

- Runtime friendship behavior implementation remains deferred.
- Protocol routes and generated friends contract shapes remain deferred.
- Request identity validation remains owned by authentication/session and route-policy boundaries.
- The future service must not parse bearer tokens, cookies, WebSocket subprotocols, envelope metadata, or transport connection identifiers.
- This gate does not alter access-token validation, runtime session validation, bound connection identity, route protection, WebSocket handshake behavior, logout behavior, or reconnect behavior.

## 15. Reference Alignment

Nakama provides the product capability pressure for durable friends relationship social graph behavior. Pitaya remains a deferred future architecture reference for distributed runtime concerns. vibit uses those references for capability planning only:

- no direct Nakama or Pitaya API compatibility is added;
- no public friends route is added by this gate;
- no server runtime hook, group, party, chat, matchmaking, match runtime, admin surface, or distributed routing behavior is added by this gate;
- the future runtime behavior is an application service above module repository vocabulary and below protocol route/public contract behavior.

## 16. Stop Conditions

Stop and open a later bounded work item before doing any of the following:

- implementing `runtime/internal/app/friends/service.go`;
- adding runtime friend request/list/status handlers;
- adding command/query dispatch registration or startup composition;
- adding protocol routes, Protobuf sources, generated output, or generated clients;
- changing authentication/session behavior or request identity validation;
- changing the friends repository interface or PostgreSQL adapter;
- adding or changing migrations;
- adding event/audit tables;
- adding chat, groups, parties, broadcast fanout, matchmaking, match runtime, SDK, hosted, release, or distributed runtime scope;
- adding direct Nakama or Pitaya public API compatibility.

