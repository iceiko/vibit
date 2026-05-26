# Friends Relationship PostgreSQL Adapter Gate

Status: Accepted v0.1
Last updated: 2026-05-26
Scope: Gate-only boundary for the future PostgreSQL adapter implementing `runtime/internal/modules/friends.Repository`
Depends on: `runtime/internal/modules/friends/repository.go`, `docs/friends-relationship-repository-boundary.md`, `runtime/migrations/postgres/000007_create_friend_relationships.sql`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0144`

The paired Simplified Chinese translation is `docs/friends-relationship-postgresql-adapter-gate.zh-CN.md`. The English file is authoritative.

This document defines the friends relationship PostgreSQL adapter gate. It is a gate artifact. It does not add PostgreSQL adapter implementation, SQL execution behavior, unit-of-work factory wiring, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, migrations, automatic startup migration behavior, event/audit tables, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The friends relationship PostgreSQL adapter gate record is:

```yaml
friends_relationship_postgresql_adapter_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0236
decision: ADR-0144
check_rule: runtime.friends_relationship_postgresql_adapter_gate
source_repository_interface_decision: ADR-0143
repository_interface: runtime/internal/modules/friends.Repository
repository_interface_source: runtime/internal/modules/friends/repository.go
repository_tests: runtime/internal/modules/friends/repository_test.go
source_migration_source_decision: ADR-0141
source_migration_source: runtime/migrations/postgres/000007_create_friend_relationships.sql
friend_relationships_logical_table: friend_relationships
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_adapter_source_candidate: runtime/internal/platform/persistence/postgres/friend_relationship_repository.go
future_adapter_tests_candidate: runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go
future_constructor_candidate: NewFriendRelationshipRepositoryForUnitOfWork
unit_of_work_handoff_required: true
transaction_owner: caller_supplied_unit_of_work
sql_mapping_owner: postgresql_adapter
adapter_gate_only: true
postgresql_adapter_added: false
sql_execution_added: false
unit_of_work_factory_added: false
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
direct_nakama_pitaya_api_compatibility_added: false
future_adapter_implementation_work_item: W-0237
future_adapter_implementation_direction: implement_friends_relationship_postgresql_adapter
```

## 2. Purpose

`W-0235` implemented the storage-neutral `runtime/internal/modules/friends.Repository` interface. The next useful boundary is the platform adapter gate that will later map that interface to the existing PostgreSQL `friend_relationships` table.

This gate records the future implementation shape before any adapter SQL is written:

- adapter ownership;
- constructor and executor handoff expectations;
- transaction and unit-of-work boundaries;
- SQL mapping posture for friend request, read, list, lifecycle transition, and block mutation methods;
- conflict, affected-row, and driver-error mapping;
- timestamp, relationship version, canonical pair, and block-column mapping;
- focused test expectations;
- stop conditions that keep runtime, protocol, generated output, broad social features, and direct compatibility out of the adapter slice.

This is not an implementation. The future adapter source path is named only so agents can verify later work against the accepted boundary.

## 3. Ownership

The future adapter owner is:

```yaml
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
repository_interface_owner: runtime/internal/modules/friends
sql_mapping_owner: runtime/internal/platform/persistence/postgres
transaction_owner: caller_supplied_unit_of_work
application_layer_owns_request_identity: true
friends_module_owns_repository_vocabulary: true
player_module_owns_player_lifecycle: true
websocket_transport_owns_friends_relationships: false
protocol_adapter_owns_friends_relationships: false
authentication_module_owns_friends_relationships: false
storage_module_owns_friends_relationships: false
```

Rules:

- The adapter may later implement `friends.Repository` under the PostgreSQL platform package.
- The adapter must not move SQL into `runtime/internal/modules/friends`.
- The adapter must not own request authentication, player identity validation, route policy, protocol parsing, WebSocket state, chat rooms, groups, parties, matchmaking, match runtime, or distributed topology.
- The adapter must receive already-normalized repository input or call friends module normalizers before SQL mapping.
- The adapter must return friends module value types and typed repository errors, not driver-specific errors.
- The adapter may check player-account foreign-key outcomes as storage conflicts, but it must not become the player account lifecycle owner.

## 4. Future Constructor And Executor Handoff

The first adapter implementation should follow existing PostgreSQL adapter patterns:

```yaml
future_constructor_candidate: NewFriendRelationshipRepositoryForUnitOfWork
future_repository_interface: runtime/internal/modules/friends.Repository
executor_source: caller_supplied
transaction_control_sql_allowed: false
unit_of_work_handoff_required: true
connection_pool_owned_by_adapter: false
context_required: true
```

Rules:

- The constructor should accept an existing executor or query interface supplied by a unit-of-work boundary rather than owning a pool directly.
- The adapter must not issue `BEGIN`, `COMMIT`, or `ROLLBACK`; transaction ownership remains with the unit-of-work runner.
- The adapter must not create automatic startup migrations.
- The adapter must not add a new dependency if existing PostgreSQL platform dependencies already cover the implementation.
- Any required dependency change must be a separate dependency-adoption decision.

## 5. SQL Mapping Posture

The future adapter may map repository methods to the `friend_relationships` table:

```yaml
logical_table: friend_relationships
primary_key_column: relationship_id
pair_columns:
  - player_low_id
  - player_high_id
lifecycle_column: lifecycle_state
actor_columns:
  - requested_by_player_id
  - responded_by_player_id
  - removed_by_player_id
block_columns:
  - blocked_by_low_at
  - blocked_by_high_at
version_column: relationship_version
created_at_column: created_at
updated_at_column: updated_at
state_changed_at_column: state_changed_at
rejected_at_column: rejected_at
removed_at_column: removed_at
pair_unique_index: friend_relationships_pair_uq
player_low_state_index: friend_relationships_player_low_state_idx
player_high_state_index: friend_relationships_player_high_state_idx
updated_at_index: friend_relationships_updated_at_idx
```

Method posture:

- `CreateOrUpdateFriendRequest` should insert a pending relationship for a canonical pair or update an ended relationship only according to the later implementation decision. It must preserve repository validation, reject self-relationships before SQL when possible, start at positive relationship version, and map active friendship, duplicate pending, blocked, and invalid transition outcomes to typed friends conflicts.
- `GetRelationshipByPair` should select by canonical pair and return a normalized `FriendRelationship`.
- `ListRelationshipsForPlayer` should be player-scoped, status-filtered, ordered deterministically, bounded by repository limits, and pagination-ready. It must not become arbitrary social graph search or admin inspection.
- `AcceptFriendRequest`, `RejectFriendRequest`, and `RemoveFriend` should update existing rows through expected lifecycle state and optional expected-version checks. Affected-row counts must distinguish not found, stale version, and invalid transition without leaking private social graph details.
- `SetPlayerBlock` and `ClearPlayerBlock` should update the actor-specific block timestamp column derived from the canonical pair member role. They must not implicitly restore friendship after unblock.

Rules:

- SQL text must remain inside the PostgreSQL adapter package.
- Private social graph data is not log-safe and must not be included in default error text.
- Driver-specific constraint names may be used internally for error mapping, but public module errors must remain friends-module neutral.
- Affected-row counts must be checked for update operations.
- The adapter must not hard-delete relationship history unless a later event/audit or retention decision explicitly authorizes it.

## 6. Transaction And Unit-Of-Work Boundary

The future adapter participates in existing transaction handoff:

```yaml
unit_of_work_handoff_required: true
adapter_starts_transactions: false
adapter_commits_transactions: false
adapter_rolls_back_transactions: false
adapter_safe_for_existing_runner: true
```

Rules:

- Application services or runtime composition may later obtain the adapter through an explicit unit-of-work boundary.
- This gate does not add that factory or composition.
- Adapter methods must use the caller's context.
- The adapter must not hide transaction failures by returning successful friends relationship results.
- The adapter must not perform route policy, session validation, access-token validation, or WebSocket close behavior.

## 7. Error Mapping

The future adapter must collapse PostgreSQL details into friends module errors:

```yaml
repository_error_owner: runtime/internal/modules/friends
driver_error_public_leakage_allowed: false
constraint_name_public_leakage_allowed: false
private_social_graph_public_leakage_allowed: false
conflict_classes:
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

Mapping expectations:

- Pair uniqueness conflicts should map to duplicate request, already friends, blocked relationship, invalid transition, or pair identity conflict according to the current row state and repository operation.
- Foreign-key failures for pair members or actor columns should map to target-player-not-found or invalid input without exposing database constraint names.
- No affected row with expected version should map to relationship-not-found, version-mismatch, stale-version, or invalid-transition without leaking hidden relationship history.
- Malformed input should be rejected before SQL execution when possible.
- Unknown driver or executor failures should map to storage-unavailable style errors with redacted reasons.
- Raw SQL, DSNs, credentials, token material, verifier digests, player ids, private relationship state, and driver stack details must not appear in public error strings.

## 8. Test Expectations

The later implementation slice should add focused PostgreSQL adapter tests. Fake-executor or query-capture tests are acceptable before live database verification is available.

Required test families for the implementation slice:

```yaml
future_tests:
  - constructor_requires_executor
  - create_or_update_friend_request_maps_insert_update_and_conflicts
  - get_relationship_selects_by_canonical_pair
  - list_relationships_is_player_scoped_status_filtered_and_ordered
  - accept_reject_remove_check_expected_version_and_transition_state
  - block_unblock_update_actor_specific_block_columns
  - rows_scan_through_friends_normalizers
  - driver_errors_are_redacted
  - transaction_control_sql_is_absent
  - default_tests_do_not_require_live_postgresql
```

Rules:

- Tests must not require protocol routes.
- Tests must not require WebSocket transport.
- Tests must not print raw credentials, tokens, verifier keys, digests, DSNs, query strings, authorization headers, cookies, player ids, or private relationship state.
- Live PostgreSQL verification may be added later when an implementation slice authorizes it; if unavailable, it must be explicitly recorded.

## 9. Relationship To Runtime, Protocol, And Authentication

This gate does not change runtime or protocol behavior:

```yaml
runtime_friends_handlers_added: false
friends_protocol_routes_added: false
protobuf_friends_messages_added: false
generated_friends_output_added: false
authentication_session_behavior_changed: false
request_identity_handoff_changed: false
```

Rules:

- Runtime friendship behavior remains deferred.
- Protocol routes and generated friends contract shapes remain deferred.
- Request identity validation remains owned by authentication/session boundaries.
- The adapter must not parse bearer tokens, cookies, WebSocket subprotocols, envelope metadata, or transport connection identifiers.

## 10. Reference Alignment

Nakama provides the product capability pressure for durable friends relationship social graph behavior. Pitaya remains a deferred future architecture reference for distributed runtime concerns. vibit uses those references for capability planning only:

- no direct Nakama or Pitaya API compatibility is added;
- no public friends route is added by this gate;
- no server runtime hook, group, party, chat, matchmaking, match runtime, or admin surface is added by this gate;
- the future adapter is a platform persistence detail below module/application behavior.

## 11. Stop Conditions

Stop and open a later bounded work item before doing any of the following:

- implementing `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`;
- adding SQL execution behavior for friends relationships;
- adding unit-of-work factory wiring or startup composition;
- adding runtime friend request/list/status behavior;
- adding protocol routes, Protobuf sources, generated output, or generated clients;
- changing authentication/session behavior or request identity validation;
- adding event/audit tables;
- adding chat, groups, parties, broadcast fanout, matchmaking, match runtime, SDK, hosted, release, or distributed runtime scope;
- adding direct Nakama or Pitaya public API compatibility.

