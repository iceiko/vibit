# Storage Objects Persistence Schema Gate

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Gate for first PostgreSQL storage objects persistence schema before migration source
Depends on: `docs/storage-objects-behavior-gate.md`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0110`

The paired Simplified Chinese translation is `docs/storage-objects-persistence-schema-gate.zh-CN.md`. The English file is authoritative.

This document defines the first storage objects persistence schema gate. It is a gate artifact. It does not add SQL migration source, create the `storage_objects` table, implement storage objects runtime behavior, add protocol routes, add Protobuf source or generated output, add dependencies, add repository interfaces, add storage adapters, broaden operations/admin behavior, add hosted deployments, create release artifacts, run public announcements, run paid promotion, change authentication/session behavior, add large object/blob storage, add S3-compatible object storage, add a broad product module implementation, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The storage objects persistence schema gate record is:

```yaml
storage_objects_persistence_schema_gate: defined
completed_work_item: W-0202
decision: ADR-0110
check_rule: runtime.storage_objects_persistence_schema_gate
source_behavior_gate_decision: ADR-0109
source_behavior_gate_standard: docs/storage-objects-behavior-gate.md
gate_standard: docs/storage-objects-persistence-schema-gate.md
gate_standard_translation: docs/storage-objects-persistence-schema-gate.zh-CN.md
target_stage: prototype_ready_foundation
selected_first_storage_objects_store: postgres
future_storage_objects_logical_table: storage_objects
future_migration_source_candidate: runtime/migrations/postgres/000006_create_storage_objects.sql
future_repository_owner_candidate: runtime/internal/modules/storage
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
table_candidate_recorded: true
owner_identity_representation_recorded: true
collection_key_constraints_recorded: true
value_representation_recorded: true
version_representation_recorded: true
timestamp_posture_recorded: true
uniqueness_index_posture_recorded: true
redaction_posture_recorded: true
future_repository_adapter_boundaries_recorded: true
future_migration_source_candidate_recorded: true
schema_gate_only: true
migration_source_added: false
storage_objects_table_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
dependency_added: false
repository_interface_changed: false
storage_adapter_changed: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_migration_work_item: W-0203
future_migration_direction: storage_objects_migration_source
```

## 2. Product Intent

`ADR-0109` selected player-owned small JSON storage objects as the first general durable game-state behavior beyond inventory. The next step is to make the future persistence shape explicit before adding SQL.

This gate makes the future migration inspectable:

- the table candidate is known;
- ownership and identity columns are selected;
- collection/key constraints are bounded;
- JSON value storage is intentionally small-object game state;
- optimistic versioning has a concrete storage posture;
- indexes and uniqueness are planned before implementation;
- redaction and future repository/adapter boundaries are explicit.

The gate keeps the work conservative. It does not add a migration file. It prepares the next migration-source-only slice.

## 3. Selected Store And Table

The first storage objects persistence target is PostgreSQL:

```yaml
selected_first_storage_objects_store: postgres
future_storage_objects_logical_table: storage_objects
future_migration_source_candidate: runtime/migrations/postgres/000006_create_storage_objects.sql
future_repository_boundary: separate_future_work_item
future_postgresql_adapter: separate_future_work_item
```

Rationale:

- PostgreSQL is vibit's first accepted authoritative durable store.
- The current source alpha already uses PostgreSQL for inventory, player accounts, authentication verifier records, token verifier records, and runtime sessions.
- Storage objects need transactionally inspectable ownership and versioning before runtime behavior exists.
- Starting with PostgreSQL avoids adopting a document database or object/blob storage dependency before the need is concrete.

## 4. Future `storage_objects` Table Candidate

The future first migration may define one logical table:

```yaml
storage_objects:
  primary_key_candidate:
    - object_id
  required_columns:
    - object_id
    - owner_kind
    - owner_id
    - collection
    - object_key
    - value_json
    - version
    - created_at
    - updated_at
  nullable_columns:
    - deleted_at
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
    - blob_bytes
    - file_path
    - s3_bucket
    - s3_object_key
```

`object_id` is a server-generated opaque record id. It is not the public identity of the object. The logical identity remains:

```text
owner_kind + owner_id + collection + object_key
```

The first table candidate may use `object_key` instead of a bare `key` column to avoid SQL readability and tool ambiguity.

## 5. Owner Identity Representation

The first schema posture supports player-owned objects only:

```yaml
owner_kind_first_value: player
owner_id_source: validated_request_identity_player_id
owner_player_fk_candidate: player_accounts(player_id)
owner_kind_check_candidate: owner_kind = 'player'
```

Rules:

- `owner_kind` must be a closed vocabulary in the first migration source.
- `owner_id` must be non-blank.
- The first owner kind is `player`.
- A foreign key to `player_accounts(player_id)` is the first candidate.
- The table must not trust client-supplied owner id as proof; runtime behavior must derive it from validated identity in a later implementation slice.

Deferred owner identities:

- global;
- group/guild;
- party;
- room;
- match;
- server shard;
- public catalog;
- admin-owned objects.

## 6. Collection And Key Constraints

The future migration should constrain collection and key fields before runtime behavior depends on them.

First candidate:

```yaml
collection_column: collection
key_column: object_key
collection_type: TEXT
object_key_type: TEXT
collection_not_blank: true
object_key_not_blank: true
collection_max_length_candidate: 128
object_key_max_length_candidate: 256
path_semantics_allowed: false
case_sensitive: true
```

Recommended SQL-level checks:

- `length(btrim(collection)) > 0`;
- `length(collection) <= 128`;
- `length(btrim(object_key)) > 0`;
- `length(object_key) <= 256`.

Stricter ASCII-safe or pattern checks may be added in a later implementation gate or migration-source slice if the exact protocol identifier rules are ratified first.

## 7. Value Representation

The first value representation candidate is PostgreSQL `JSONB`:

```yaml
value_column: value_json
value_type_candidate: JSONB
value_top_level_shape: object
value_binary_storage: false
value_log_safe: false
```

Rules:

- The value must remain small durable game state, not large binary content.
- The value must not contain raw credentials, raw tokens, verifier keys, DSNs, digest bytes, transport metadata, or other secrets.
- The future migration source should include a SQL check that the top-level value is a JSON object when PostgreSQL support is available.
- Maximum byte size remains an implementation-gate requirement; the migration-source slice may choose a SQL check if the maximum is ratified there.

Deferred value behavior:

- arbitrary JSON scalar or array payloads;
- binary blobs;
- external file references;
- S3-compatible object references;
- encrypted-at-rest application payload envelope;
- document search indexes;
- JSON patch or merge behavior.

## 8. Version Representation

The first version representation candidate is a server-managed integer revision:

```yaml
version_column: version
version_type_candidate: BIGINT
initial_version_candidate: 1
version_increment_policy: server_managed_on_successful_mutation
client_authoritative_version: false
```

Rules:

- `version` must be positive.
- Successful create starts at version `1` unless the future migration-source decision selects another positive start.
- Successful update increments version.
- Delete version behavior must be defined by the future runtime behavior gate or implementation gate.
- The public protocol may later expose a string form or opaque token, but the first stored representation is a numeric revision candidate.

This gate does not add compare-and-swap runtime behavior. It only defines the persistence representation candidate that can support optimistic concurrency.

## 9. Timestamp And Deletion Posture

The first timestamp posture is:

```yaml
created_at: TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at: TIMESTAMPTZ NOT NULL DEFAULT now()
deleted_at: TIMESTAMPTZ NULL
soft_delete_candidate: true
hard_delete_candidate: deferred
```

Rules:

- `updated_at >= created_at` should be enforced.
- If `deleted_at` is present, `deleted_at >= created_at` should be enforced.
- Soft delete is the first schema candidate because it helps future conflict and audit behavior without immediately adding event history.
- Runtime behavior may still choose whether delete hides rows, tombstones rows, or later hard-deletes rows after a cleanup gate.

The schema does not add audit/event tables.

## 10. Uniqueness And Index Posture

The first uniqueness posture is:

```yaml
logical_identity_unique_candidate:
  - owner_kind
  - owner_id
  - collection
  - object_key
active_unique_index_candidate:
  - owner_kind
  - owner_id
  - collection
  - object_key
  - deleted_at IS NULL
```

Recommended indexes for the future migration-source slice:

- unique active logical identity index on `(owner_kind, owner_id, collection, object_key)` where `deleted_at IS NULL`;
- lookup index for `(owner_kind, owner_id, collection)`;
- updated-at index for future diagnostics or cleanup;
- optional deleted-at index only if cleanup or tombstone queries are authorized later.

This gate does not authorize global search, cross-owner search, JSONB GIN indexes, admin listing, analytics indexes, or operations dashboards.

## 11. Ownership Boundaries

Future storage object behavior should have its own module boundary:

```yaml
future_repository_owner_candidate: runtime/internal/modules/storage
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_contract_owner_candidate: contracts/storage
future_proto_source_candidate: proto/vibit/storage/v1
storage_module_owns_storage_objects: true
inventory_module_owns_storage_objects: false
player_module_owns_storage_objects: false
authentication_module_owns_storage_objects: false
websocket_transport_owns_storage_objects: false
```

Rules:

- The future storage module owns storage object domain behavior.
- The player module owns player account lifecycle, not storage object behavior.
- The inventory module owns inventory state, not general storage objects.
- Authentication owns credentials/tokens/sessions, not storage object data.
- WebSocket transport owns connection plumbing, not durable object records.
- PostgreSQL adapters may implement storage repositories only after repository boundaries are authorized.

## 12. Redaction Posture

Storage object values are not log-safe by default.

Not log-safe by default:

- `value_json`;
- raw object values in validation errors;
- owner ids when combined with collection/key in public diagnostics;
- collection and key until a later redaction decision makes them safe for specific logs;
- JSON path fragments from rejected payloads;
- database error details that include values.

Forbidden secret material:

- raw device credentials;
- raw access tokens;
- verifier keys;
- lookup or verifier digests;
- PostgreSQL DSNs with credentials;
- headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata.

Future adapters and handlers must return redacted errors.

## 13. Future Migration Source Expectations

The next bounded work item may add:

```text
runtime/migrations/postgres/000006_create_storage_objects.sql
```

That migration-source-only slice may add SQL DDL, comments, indexes, and migration checks for `storage_objects`.

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
- large object/blob storage;
- S3-compatible object storage;
- direct Nakama/Pitaya API compatibility.

## 14. Verification Expectations

The future migration-source slice should verify:

- goose up/down markers;
- table name and required columns;
- owner, collection, key, value, version, timestamp, and deletion checks;
- active logical identity uniqueness;
- no forbidden secret, digest, transport, blob, file, or S3 columns;
- no Go runtime behavior;
- repository checks for migration boundary.

Later repository/adapter/runtime work should add focused tests for create/read/update/delete/conflict/cross-owner behavior only after those implementation gates are accepted.

## 15. Stop Conditions

Stop and ask for maintainer authorization before doing any of the following:

- adding the SQL migration source in the same change as this gate;
- implementing storage objects runtime behavior;
- adding protocol routes;
- adding Protobuf source files or generated output;
- adding repository interfaces;
- adding PostgreSQL or other storage adapters;
- adding dependencies;
- changing authentication/session semantics;
- changing route protection semantics;
- adding cross-player, global, group, party, room, match, public, admin, or ACL object scopes;
- adding large object/blob storage or S3-compatible object storage;
- adding JSONB search indexes for product search behavior;
- adding server-side custom logic hooks;
- adding broad operations/admin behavior;
- adding hosted deployments or demos;
- creating release binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry publications, SDK packages, or additional release artifacts;
- executing public announcements beyond the GitHub release record;
- running paid promotion;
- adding direct Nakama/Pitaya API compatibility.

## 16. Next Work

The next bounded direction is:

```text
W-0203 Add storage objects migration source
```

That work may add the first SQL migration source for `storage_objects` and matching static checks, while keeping repository interfaces, adapters, protocol, and runtime behavior deferred.
