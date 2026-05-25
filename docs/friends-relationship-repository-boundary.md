# Friends Relationship Repository Boundary

Status: Accepted v0.1
Last updated: 2026-05-25
Scope: Gate-only boundary for the future storage-neutral friends relationship repository after the PostgreSQL `friend_relationships` migration source
Depends on: `docs/friends-relationship-lifecycle-gate.md`, `docs/friends-relationship-persistence-schema-gate.md`, `decisions/ADR-0141-friends-relationship-migration-source.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0142`

The paired Simplified Chinese translation is `docs/friends-relationship-repository-boundary.zh-CN.md`. The English file is authoritative.

This document defines the friends relationship repository boundary. It is a gate artifact. It does not add Go repository interfaces, PostgreSQL adapter behavior, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, migrations, automatic startup migration behavior, event/audit tables, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The friends relationship repository boundary record is:

```yaml
friends_relationship_repository_boundary: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0234
decision: ADR-0142
check_rule: runtime.friends_relationship_repository_boundary
source_migration_source_decision: ADR-0141
source_migration_source: runtime/migrations/postgres/000007_create_friend_relationships.sql
source_schema_gate_decision: ADR-0140
source_lifecycle_gate_decision: ADR-0139
future_repository_owner_candidate: runtime/internal/modules/friends
future_repository_interface_candidate: runtime/internal/modules/friends.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
friend_relationships_logical_table: friend_relationships
repository_boundary_gate_only: true
repository_interface_added: false
postgresql_adapter_added: false
runtime_behavior_added: false
authentication_session_behavior_changed: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
dependency_added: false
migration_added: false
event_audit_table_added: false
chat_added: false
groups_added: false
parties_added: false
matchmaking_added: false
match_runtime_added: false
sdk_added: false
generated_client_library_added: false
hosted_deployment_added: false
release_artifact_added: false
distributed_runtime_added: false
pitaya_distributed_architecture_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_repository_interface_work_item: W-0235
future_repository_interface_direction: implement_friends_relationship_repository_interface
```

## 2. Purpose

`W-0233` added the PostgreSQL migration source for `friend_relationships`. The next useful boundary is the storage-neutral repository vocabulary that future implementation code can use without exposing SQL details, transport details, or protocol assumptions.

This boundary prepares the Nakama-class friends relationship path by recording:

- repository ownership;
- candidate value types;
- lifecycle command and query vocabulary;
- pair identity and actor handoff rules;
- version, conflict, and transaction handoff posture;
- redaction and error posture;
- PostgreSQL adapter expectations;
- stop conditions for future implementation work.

This is still not a runtime feature. No handler, route, adapter, repository interface, or protocol message can use friends relationships until later bounded work items explicitly authorize them.

## 3. Ownership

The future repository is friends module-owned:

```yaml
future_repository_owner_candidate: runtime/internal/modules/friends
future_repository_interface_candidate: runtime/internal/modules/friends.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
friend_relationships_table_owner: runtime.friends
application_layer_owns_request_identity: true
postgresql_adapter_owns_sql_mapping: true
websocket_transport_owns_friends_relationships: false
protocol_adapter_owns_friends_relationships: false
authentication_module_owns_friends_relationships: false
player_module_owns_friends_relationships: false
storage_module_owns_friends_relationships: false
```

Rules:

- The future repository interface must be storage-neutral and module-facing.
- The interface must not mention PostgreSQL, pgx, SQL rows, goose migrations, prepared statements, connection pools, transaction runners, or database driver errors.
- The PostgreSQL adapter may later implement the interface under `runtime/internal/platform/persistence/postgres`, but only after a separate adapter gate.
- Application or handler code may later call friends relationship behavior through module/application boundaries, not through SQL or transport state.
- Authentication and session code provide validated request identity; they do not own friends relationship records.
- Player account storage owns player lifecycle state, not social graph relationship state.
- Storage objects own player-owned small JSON game state, not friends relationships.
- WebSocket transport owns connection plumbing, not friendship state.
- Protocol adapters own wire conversion, not repository behavior.

## 4. Candidate Value Types

A later implementation gate may rename or reduce these shapes, but the first repository interface implementation should start from this vocabulary:

```yaml
candidate_value_types:
  - FriendRelationship
  - FriendRelationshipPair
  - FriendRelationshipID
  - FriendRelationshipStatus
  - FriendRelationshipLifecycleState
  - FriendRelationshipActor
  - FriendRelationshipVersion
  - FriendRelationshipBlockState
  - SendFriendRequestInput
  - AcceptFriendRequestInput
  - RejectFriendRequestInput
  - RemoveFriendInput
  - BlockPlayerInput
  - UnblockPlayerInput
  - GetFriendRelationshipInput
  - ListFriendRelationshipsInput
  - FriendRelationshipConflict
  - FriendRelationshipRepositoryError
```

First-posture record vocabulary:

```yaml
friend_relationship_record:
  relationship_id: server_generated_opaque_id
  pair: canonical_unordered_player_pair
  player_low_id: canonical_pair_member
  player_high_id: canonical_pair_member
  lifecycle_state: pending_or_friends_or_rejected_or_removed
  requested_by_player_id: nullable_pair_member_actor
  responded_by_player_id: nullable_pair_member_actor
  removed_by_player_id: nullable_pair_member_actor
  blocked_by_low_at: nullable_server_timestamp
  blocked_by_high_at: nullable_server_timestamp
  relationship_version: server_managed_bigint_revision
  created_at: server_timestamp
  updated_at: server_timestamp
  state_changed_at: server_timestamp
  rejected_at: nullable_server_timestamp
  removed_at: nullable_server_timestamp
```

Rules:

- Pair identity must be a canonical unordered player pair.
- `player_low_id` and `player_high_id` are persistence identity fields, not authentication proof.
- `requested_by_player_id`, `responded_by_player_id`, and `removed_by_player_id` are normalized pair-member actors, not proof.
- Future runtime behavior must derive the actor from validated request identity before calling the repository.
- Public actor-relative states such as outgoing pending, incoming pending, blocked by actor, and actor blocked target are computed later; they are not repository lifecycle states.
- `relationship_version` is server-managed and must not be client-authoritative state.
- Private social graph data is not log-safe by default.

## 5. Candidate Repository Capabilities

The first storage-neutral capability family is:

```yaml
candidate_repository_capabilities:
  - CreateOrUpdateFriendRequest
  - GetRelationshipByPair
  - ListRelationshipsForPlayer
  - AcceptFriendRequest
  - RejectFriendRequest
  - RemoveFriend
  - SetPlayerBlock
  - ClearPlayerBlock
```

Capability rules:

- `CreateOrUpdateFriendRequest` may create or update a relationship row only for an already validated actor and normalized target pair.
- `GetRelationshipByPair` is a storage lookup. It must not authenticate users, validate access tokens, or create request identity.
- `ListRelationshipsForPlayer` must be player-scoped and pagination-ready. It must not become arbitrary social graph search or admin inspection without a later gate.
- `AcceptFriendRequest` and `RejectFriendRequest` must be lifecycle transitions over an existing pending relationship.
- `RemoveFriend` must end an active friendship or pending relationship according to later behavior rules; it must not hard-delete audit-relevant history by default.
- `SetPlayerBlock` and `ClearPlayerBlock` must preserve actor-specific block semantics and must not restore friendship implicitly after unblock.
- All methods must return typed module-owned records and errors, not raw SQL rows or database driver errors.

The future repository interface may choose shorter names, but it must preserve the semantic split between request creation, relationship read/list, lifecycle transition, block mutation, and conflict handling.

## 6. Pair Identity And Request Identity Handoff

The repository boundary prepares canonical pair handling without implementing behavior:

```yaml
pair_identity: canonical_unordered_player_pair
self_relationship_allowed: false
actor_identity_source: validated_request_identity_before_repository_call
client_supplied_actor_id_as_proof_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
actor_relative_public_status_stored: false
```

Rules:

- The repository may receive normalized actor ids as data, but actor ids are not authentication proof.
- A player must never be allowed to form a relationship with itself.
- Pair canonicalization must be deterministic and independent of request direction.
- Public errors should not reveal hidden relationship history when future behavior gates require leakage collapse.
- The repository must not accept transport metadata, WebSocket connection identifiers, cookies, headers, tokens, or sessions as identity proof.

## 7. Version And Conflict Handoff

The repository boundary prepares optimistic concurrency without implementing behavior:

```yaml
version_storage: BIGINT
initial_create_version: 1
version_owner: server
client_authoritative_version_allowed: false
expected_version_handoff: future_behavior_or_interface_gate
conflict_public_shape: deferred_to_protocol_gate
```

Candidate conflict classes:

```yaml
candidate_conflict_classes:
  - relationship_not_found
  - target_player_not_found
  - self_relationship_forbidden
  - duplicate_pending_request
  - already_friends
  - blocked_relationship
  - invalid_transition
  - version_mismatch
  - stale_relationship_version
  - pair_identity_conflict
  - storage_unavailable
```

Rules:

- Repository methods may later distinguish internal typed conflicts, but public protocol error mapping remains deferred.
- Version equality is not authentication proof.
- A stale expected version must not be collapsed into a hidden successful write.
- Target-player-not-found and relationship-not-found leakage rules remain future behavior decisions.
- The PostgreSQL adapter must map unique-index, affected-row, and foreign-key outcomes into typed repository conflicts without exposing driver error text.

## 8. Redaction And Logging

Friends relationship state is private social graph data.

```yaml
private_social_graph_log_safe: false
relationship_id_log_safe: conditional_after_validation
player_id_log_safe: false_by_default
lifecycle_state_log_safe: conditional_after_validation
version_log_safe: conditional_after_validation
forbidden_repository_material:
  - raw_access_token
  - raw_credential
  - credential_lookup_digest
  - credential_verifier_digest
  - token_lookup_digest
  - token_verifier_digest
  - verifier_key
  - websocket_connection_id
  - websocket_subprotocol
  - remote_address
  - authorization_header
  - cookie
  - query_string_token
  - chat_room_id
  - group_id
  - party_id
  - match_id
  - pitaya_server_id
```

Rules:

- Repository errors must be redacted and typed.
- Raw player relationship details and storage driver errors must not be logged by default.
- Private social graph data must be treated as non-log-safe unless a later redaction policy narrows that.
- The repository must not store or return authentication material, token material, verifier digests, transport metadata, chat room state, group state, party state, match state, or distributed routing state.

## 9. PostgreSQL Adapter Expectations

The future PostgreSQL adapter may later map the repository to:

```yaml
logical_table: friend_relationships
pair_unique_index: friend_relationships_pair_uq
player_low_state_index: friend_relationships_player_low_state_idx
player_high_state_index: friend_relationships_player_high_state_idx
updated_at_index: friend_relationships_updated_at_idx
version_column: relationship_version
```

Adapter expectations:

- SQL execution belongs under `runtime/internal/platform/persistence/postgres`.
- Unit-of-work and transaction handoff must follow the existing platform transaction boundary.
- SQL must not leak into `runtime/internal/modules/friends`.
- The adapter must preserve canonical pair uniqueness over `player_low_id + player_high_id`.
- Updates must be affected-row checked and version-aware when expected version is supplied.
- Actor columns must be pair-member checked before storage or mapped to typed conflicts.
- Adapter tests should cover request creation, duplicate request, accept, reject, remove, block, unblock, pair lookup, player listing, stale version, foreign-key target absence, canonical pair normalization, timestamp mapping, and redacted errors after the adapter gate is accepted.

## 10. Relationship To Runtime, Protocol, And Authentication

This boundary does not change runtime or protocol behavior:

```yaml
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
authentication_session_behavior_changed: false
request_identity_validation_added: false
```

Rules:

- Future runtime behavior must derive actor identity from already validated request identity.
- This repository boundary does not authorize friendship request, accept, reject, remove, block, unblock, list, or status runtime handlers.
- This repository boundary does not authorize any protocol routes, Protobuf messages, generated clients, SDKs, or transport carriers.
- The repository does not authenticate users, parse tokens, validate sessions, bind WebSocket connections, or enforce route policy.

## 11. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - friends_relationships_are_core_social_graph_state
  - friend_request_accept_reject_remove_block_unblock_need_durable_state
  - list_and_status_queries_need_repository_ready_boundaries
adapted_concepts:
  - repository_is_vibit_storage_neutral_module_boundary
  - public_status_is_actor_relative_and_computed_later
  - schema_and_protocol_are_vibit_native_not_direct_api_compatibility
deferred_concepts:
  - groups
  - parties
  - chat_targeting
  - matchmaking_social_filters
  - match_runtime_social_context
rejected_for_now:
  - direct_nakama_session_or_friends_api_compatibility
```

Pitaya reference mapping:

```yaml
pitaya_reference_status: deferred_future_architecture_reference
deferred_concepts:
  - frontend_backend_cluster_social_graph_routing
  - distributed_group_membership
  - RPC_or_service_discovery_for_friendship_operations
  - distributed_session_social_context
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama is the primary product reference for near-term capability coverage. Pitaya remains deferred as a future distributed architecture reference. Neither reference overrides vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

## 12. Future Implementation Queue

After this boundary, future work should remain split:

```yaml
future_work_items:
  friends_relationship_repository_interface_implementation:
    work_item: W-0235
    may_add:
      - runtime/internal/modules/friends
      - storage-neutral repository interface and value types
      - focused repository vocabulary tests
    must_not_add:
      - PostgreSQL adapter behavior
      - runtime friendship behavior
      - protocol routes
      - Protobuf source
      - generated output
  friends_relationship_postgresql_adapter_gate:
    may_define:
      - adapter ownership
      - transaction handoff
      - SQL mapping
      - adapter tests
  friends_relationship_runtime_behavior_gate:
    may_define:
      - actor derivation from validated request identity
      - lifecycle command behavior
      - public conflict and leakage policy
  friends_relationship_protocol_route_gate:
    may_define:
      - route names
      - payloads
      - generated output posture
```

Do not combine these into one broad social subsystem slice without a new ADR.

## 13. Stop Conditions

Stop and ask for a separate bounded work item before:

- adding `runtime/internal/modules/friends`;
- adding a friends repository interface implementation;
- adding PostgreSQL friends adapter behavior;
- adding runtime friendship behavior;
- adding protocol routes, Protobuf source, or generated output;
- adding migrations or event/audit tables;
- adding chat, groups, parties, matchmaking, match runtime, SDKs, hosted surfaces, release artifacts, or distributed runtime;
- adding direct Nakama/Pitaya API compatibility.

## 14. Verification

Repository verification for this boundary is:

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-friends-relationship-repository-boundary --json
node tools/vibit check all --json
```

The repository check rule is:

```yaml
runtime.friends_relationship_repository_boundary
```
