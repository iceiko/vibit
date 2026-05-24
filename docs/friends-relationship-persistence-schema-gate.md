# Friends Relationship Persistence Schema Gate

Status: Accepted v0.1
Last updated: 2026-05-24
Scope: Gate for future PostgreSQL friends relationship persistence schema before migration source, repository interfaces, adapters, runtime behavior, protocol routes, generated output, or broader social features
Depends on: `docs/friends-relationship-lifecycle-gate.md`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0140`

The paired Simplified Chinese translation is `docs/friends-relationship-persistence-schema-gate.zh-CN.md`. The English file is authoritative.

This document defines the friends relationship persistence schema gate. It is a gate artifact. It does not add SQL migration source, create the `friend_relationships` table, implement friendship runtime behavior, add protocol routes, add Protobuf source or generated output, add dependencies, add repository interfaces, add PostgreSQL adapters, wire startup, change authentication/session behavior, add delivery guarantees, add stream subscriptions, add chat rooms, add groups, add parties, add broadcast fanout, add matchmaking, add match runtime, add operations/admin behavior, publish SDKs or generated client libraries, create hosted deployments or release artifacts, add Pitaya-style distributed architecture, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The friends relationship persistence schema gate record is:

```yaml
friends_relationship_persistence_schema_gate: defined
completed_work_item: W-0232
decision: ADR-0140
check_rule: runtime.friends_relationship_persistence_schema_gate
source_lifecycle_gate_decision: ADR-0139
source_lifecycle_gate_standard: docs/friends-relationship-lifecycle-gate.md
gate_standard: docs/friends-relationship-persistence-schema-gate.md
gate_standard_translation: docs/friends-relationship-persistence-schema-gate.zh-CN.md
selected_nakama_capability_family: friends_groups_and_parties
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
selected_first_friends_relationship_store: postgres
future_friends_relationships_logical_table: friend_relationships
future_friend_relationship_events_logical_table: deferred
future_migration_source_candidate: runtime/migrations/postgres/000007_create_friend_relationships.sql
future_repository_owner_candidate: runtime/internal/modules/friends
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
pair_identity_posture_recorded: true
relationship_state_representation_recorded: true
block_representation_recorded: true
index_uniqueness_posture_recorded: true
timestamp_posture_recorded: true
event_audit_posture_recorded: true
redaction_posture_recorded: true
future_repository_adapter_boundaries_recorded: true
future_migration_source_candidate_recorded: true
schema_gate_only: true
migration_source_added: false
friend_relationships_table_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_migration_work_item: W-0233
future_migration_direction: add_friends_relationship_migration_source
```

## 2. Product Intent

`ADR-0139` defined the future friends relationship lifecycle for a Nakama-class social graph. The next step is to make the future durable state explicit before adding SQL.

This gate makes the future migration inspectable:

- the table candidate is known;
- pair identity and canonical ordering are selected before SQL;
- relationship state is represented separately from actor-relative public status;
- block state is represented without restoring prior friendship on unblock;
- indexes and uniqueness are planned before implementation;
- event/audit posture is explicit before adding event storage;
- redaction and future repository/adapter boundaries are explicit.

The gate keeps the work conservative. It prepares the next migration-source-only slice but does not add the migration file.

## 3. Selected Store And Table

The first friends relationship persistence target is PostgreSQL:

```yaml
selected_first_friends_relationship_store: postgres
future_friends_relationships_logical_table: friend_relationships
future_migration_source_candidate: runtime/migrations/postgres/000007_create_friend_relationships.sql
future_repository_boundary: separate_future_work_item
future_postgresql_adapter: separate_future_work_item
```

Rationale:

- PostgreSQL is vibit's first accepted authoritative durable store.
- Friends relationship state is durable social graph state and must be transactionally inspectable.
- The first implementation should not introduce a graph database, document database, cache dependency, or distributed social graph subsystem before the local schema is proven.
- Nakama is the product capability reference; vibit uses its own schema and contract posture.

## 4. Future `friend_relationships` Table Candidate

The future first migration may define one logical current-state table:

```yaml
friend_relationships:
  primary_key_candidate:
    - relationship_id
  required_columns:
    - relationship_id
    - player_low_id
    - player_high_id
    - lifecycle_state
    - relationship_version
    - created_at
    - updated_at
    - state_changed_at
  nullable_columns:
    - requested_by_player_id
    - responded_by_player_id
    - removed_by_player_id
    - rejected_at
    - removed_at
    - blocked_by_low_at
    - blocked_by_high_at
  forbidden_columns:
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
    - channel_id
    - chat_room_id
    - group_id
    - party_id
    - match_id
    - pitaya_server_id
    - nakama_api_path
```

`relationship_id` is a server-generated opaque record id. It is not the public identity of the relationship. The logical identity is the unordered player pair:

```text
player_low_id + player_high_id
```

`player_low_id` and `player_high_id` are canonical pair members. The future migration source should enforce self-target prevention and canonical ordering, for example with a check equivalent to:

```text
player_low_id < player_high_id
```

The schema is pair-oriented. Public query output remains actor-relative and is computed in future application behavior.

## 5. Pair Identity And Player References

The first pair identity posture is:

```yaml
pair_identity: canonical_unordered_player_pair
player_low_id_source: canonicalized_player_pair_member
player_high_id_source: canonicalized_player_pair_member
player_fk_candidate: player_accounts(player_id)
self_relationship_allowed: false
client_supplied_actor_id_as_proof_allowed: false
metadata_only_player_id_allowed_as_proof: false
```

Rules:

- A relationship row must never represent a player targeting itself.
- Both pair members should reference existing player account records.
- Future runtime behavior must derive the actor from validated request identity, not from client-supplied actor ids.
- The pair columns are persistence identity only; they are not authentication proof.
- The table must not store session ids or transport connection identifiers.

## 6. Relationship State Representation

The first lifecycle state candidate is:

```yaml
lifecycle_state_column: lifecycle_state
lifecycle_state_type: TEXT
allowed_lifecycle_states:
  - pending
  - friends
  - rejected
  - removed
actor_relative_public_status_stored: false
```

Rules:

- `pending` records a request awaiting response.
- `friends` records an accepted relationship.
- `rejected` records an explicitly rejected request as current pair state for the first schema candidate.
- `removed` records an ended relationship or neutral row retained for block posture.
- Actor-relative public states such as `outgoing_request_pending`, `incoming_request_pending`, `blocked_by_actor`, and `blocked_actor` are computed by future query behavior and must not be stored as canonical database states.
- Duplicate request idempotency remains a future runtime behavior decision; the schema only ensures there is at most one current pair row.

The first schema candidate stores `rejected` and `removed` as current row states rather than audit-only facts. Retention, cleanup, and hard-delete behavior remain deferred to a later gate.

## 7. Request And Response Actor Columns

The first request/response posture is:

```yaml
requested_by_player_id: nullable_pair_member
responded_by_player_id: nullable_pair_member
removed_by_player_id: nullable_pair_member
```

Rules:

- `requested_by_player_id` should be present for `pending`, `friends`, and `rejected` states when the state originated from a request.
- `responded_by_player_id` may be present for accepted or rejected states.
- `removed_by_player_id` may be present for removed states.
- These columns are pair-member references for state history and conflict handling; they are not authentication proof.
- Public errors and logs must not expose hidden relationship history when privacy requires collapse.

## 8. Block Representation

Block state is actor-specific and independent from the lifecycle state:

```yaml
block_columns:
  - blocked_by_low_at
  - blocked_by_high_at
block_representation: per_pair_member_timestamp
mutual_block_representation: both_block_columns_present
unblock_restores_prior_friendship: false
```

Rules:

- `blocked_by_low_at` means `player_low_id` blocked `player_high_id`.
- `blocked_by_high_at` means `player_high_id` blocked `player_low_id`.
- If both are present, the public actor-relative status is mutual block.
- Block must override pending or friends state in future behavior.
- Unblock clears the actor's block timestamp and must not automatically restore a prior friendship.
- A block-only row may use `lifecycle_state: removed` until a future behavior gate chooses a different representation.

## 9. Version, Timestamp, And Retention Posture

The first version and timestamp posture is:

```yaml
relationship_version_column: relationship_version
relationship_version_type_candidate: BIGINT
initial_relationship_version_candidate: 1
created_at: TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at: TIMESTAMPTZ NOT NULL DEFAULT now()
state_changed_at: TIMESTAMPTZ NOT NULL DEFAULT now()
soft_delete_column: deferred
hard_delete_policy: deferred
```

Rules:

- `relationship_version` must be positive and server-managed.
- Future command behavior should increment the version on successful state mutation.
- `updated_at >= created_at` should be enforced.
- `state_changed_at >= created_at` should be enforced.
- `rejected_at`, `removed_at`, `blocked_by_low_at`, and `blocked_by_high_at`, when present, should not precede `created_at`.
- Cleanup, retention windows, hard delete, and tombstone pruning are deferred.

## 10. Uniqueness And Index Posture

The first uniqueness posture is:

```yaml
logical_pair_unique_candidate:
  - player_low_id
  - player_high_id
```

Recommended indexes for the future migration-source slice:

- unique pair identity index on `(player_low_id, player_high_id)`;
- lookup index for `(player_low_id, lifecycle_state)`;
- lookup index for `(player_high_id, lifecycle_state)`;
- updated-at index for future diagnostics or cleanup;
- optional block indexes only if the migration-source slice can keep them narrow and justified.

This gate does not authorize global player search, relationship recommendations, social graph traversal, analytics indexes, admin dashboards, chat targeting indexes, group or party indexes, matchmaking indexes, or distributed graph routing.

## 11. Event And Audit Posture

The first migration source candidate should add only the current-state table:

```yaml
future_friend_relationship_events_logical_table: deferred
outbox_table_added_by_schema_gate: false
audit_table_added_by_schema_gate: false
domain_events_defined_by_lifecycle_gate: true
```

Rationale:

- The lifecycle gate already defined future domain events.
- A current-state table is enough for the first migration source.
- Event history, audit retention, outbox delivery, and analytics should be separated into a later bounded gate.

Future runtime behavior must still keep emitted events and state changes consistent within the unit-of-work boundary once behavior is authorized.

## 12. Ownership Boundaries

Future friends relationship behavior should have its own module boundary:

```yaml
future_repository_owner_candidate: runtime/internal/modules/friends
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_contract_owner_candidate: contracts/friends
future_proto_source_candidate: proto/vibit/friends/v1
friends_module_owns_friend_relationships: true
storage_module_owns_friend_relationships: false
player_module_owns_friend_relationships: false
authentication_module_owns_friend_relationships: false
websocket_transport_owns_friend_relationships: false
```

Rules:

- The future friends module owns relationship lifecycle domain behavior.
- The player module owns player account lifecycle, not social graph transitions.
- Authentication owns credentials, tokens, and sessions, not friendship state.
- Storage objects own generic player-owned JSON objects, not social graph relationships.
- WebSocket transport owns connection plumbing, not durable relationship records.
- PostgreSQL adapters may implement friends repositories only after repository boundaries are authorized.

## 13. Redaction Posture

Friends relationship records are not log-safe by default.

Not log-safe by default:

- relationship ids;
- pair member ids when combined into social graph records;
- lifecycle state;
- request, response, removal, rejection, or block actor ids;
- block timestamps;
- conflict details;
- database errors that expose pair identity or private relationship history.

Forbidden secret and transport material:

- raw device credentials;
- raw access tokens;
- verifier keys;
- lookup or verifier digests;
- PostgreSQL DSNs with credentials;
- headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata.

Future adapters and handlers must return redacted errors, and public failures must not leak hidden relationship details where privacy requires collapse.

## 14. Future Migration Source Expectations

The next bounded work item may add:

```text
runtime/migrations/postgres/000007_create_friend_relationships.sql
```

That migration-source-only slice may add SQL DDL, comments, indexes, and migration checks for `friend_relationships`.

It must not add:

- Go repository interfaces;
- PostgreSQL adapter behavior;
- runtime handlers;
- protocol routes;
- Protobuf source files;
- generated output;
- startup wiring;
- automatic migration apply behavior;
- dependencies;
- chat, groups, parties, matchmaking, match runtime, or operations/admin behavior;
- Pitaya-style distributed architecture;
- direct Nakama/Pitaya API compatibility.

## 15. Verification Expectations

The future migration-source slice should verify:

- goose up/down markers;
- table name and required columns;
- pair ordering and self-target prevention checks;
- lifecycle state vocabulary checks;
- relationship version and timestamp checks;
- pair uniqueness and list-query indexes;
- no forbidden secret, digest, transport, chat, group, party, match, Pitaya, or Nakama compatibility columns;
- no Go runtime behavior;
- repository checks for migration boundary.

Later repository/adapter/runtime work should add focused tests for request, accept, reject, remove, block, unblock, list, status, privacy, redaction, and concurrency behavior only after those implementation gates are accepted.

## 16. Stop Conditions

Stop and ask for maintainer authorization before doing any of the following:

- adding the SQL migration source in the same change as this gate;
- creating `friend_relationships`;
- implementing friends relationship runtime behavior;
- adding protocol routes;
- adding Protobuf source files or generated output;
- adding repository interfaces;
- adding PostgreSQL or other storage adapters;
- adding dependencies;
- changing authentication/session semantics;
- changing route protection semantics;
- adding chat rooms, groups, parties, matchmaking, match runtime, operations/admin behavior, social graph search, recommendation, analytics, or distributed graph routing;
- adding hosted deployments or demos;
- creating release binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry publications, SDK packages, or additional release artifacts;
- executing public announcements beyond the GitHub release record;
- running paid promotion;
- adding direct Nakama/Pitaya API compatibility.

## 17. Next Work

The next bounded direction is:

```text
W-0233 Add friends relationship migration source
```

That work may add the first SQL migration source for `friend_relationships` and matching static checks, while keeping repository interfaces, adapters, protocol, generated output, and runtime behavior deferred.
