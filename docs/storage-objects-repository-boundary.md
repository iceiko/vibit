# Storage Objects Repository Boundary

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Gate-only boundary for the future storage-neutral storage objects repository after the PostgreSQL `storage_objects` migration source
Depends on: `docs/storage-objects-behavior-gate.md`, `docs/storage-objects-persistence-schema-gate.md`, `decisions/ADR-0111-storage-objects-migration-source.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0112`

The paired Simplified Chinese translation is `docs/storage-objects-repository-boundary.zh-CN.md`. The English file is authoritative.

This document defines the storage objects repository boundary. It is a gate artifact. It does not add Go repository interfaces, PostgreSQL storage adapters, runtime behavior, protocol routes, Protobuf source, generated output, dependencies, migrations, automatic startup migration behavior, broad operations/admin behavior, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The storage objects repository boundary record is:

```yaml
storage_objects_repository_boundary: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0204
decision: ADR-0112
check_rule: runtime.storage_objects_repository_boundary
source_migration_source_decision: ADR-0111
source_migration_source: runtime/migrations/postgres/000006_create_storage_objects.sql
future_repository_owner_candidate: runtime/internal/modules/storage
future_repository_interface_candidate: runtime/internal/modules/storage.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
storage_objects_logical_table: storage_objects
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
broad_operations_admin_behavior_added: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_repository_interface_work_item: W-0205
future_repository_interface_direction: storage_objects_repository_interface_implementation
```

## 2. Purpose

`W-0203` added the PostgreSQL migration source for `storage_objects`. The next useful boundary is the storage-neutral repository seam that later implementation code can use without exposing SQL details, transport details, or protocol assumptions.

This boundary prepares storage objects for prototype-ready use by recording:

- repository ownership;
- candidate value types;
- create/read/list/update/delete vocabulary;
- optimistic conflict and version handoff posture;
- redaction and error posture;
- PostgreSQL adapter expectations;
- stop conditions for future implementation work.

This is still not a runtime feature. No route, handler, adapter, or protocol message can use storage objects until later bounded work items explicitly authorize them.

## 3. Ownership

The future repository is storage module-owned:

```yaml
future_repository_owner_candidate: runtime/internal/modules/storage
future_repository_interface_candidate: runtime/internal/modules/storage.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
storage_objects_table_owner: runtime.storage
application_layer_owns_request_identity: true
postgresql_adapter_owns_sql_mapping: true
websocket_transport_owns_storage_objects: false
protocol_adapter_owns_storage_objects: false
authentication_module_owns_storage_objects: false
player_module_owns_storage_objects: false
```

Rules:

- The future repository interface must be storage-neutral and module-facing.
- The interface must not mention PostgreSQL, pgx, SQL rows, goose migrations, prepared statements, connection pools, or database transaction implementation details.
- The PostgreSQL adapter may later implement the interface under `runtime/internal/platform/persistence/postgres`, but only after a separate adapter gate.
- Application or handler code may later call storage object behavior through module/application boundaries, not through SQL or transport state.
- Authentication and session code provide validated request identity; they do not own storage object rows.
- Player account storage owns player lifecycle state, not general player-owned storage objects.
- WebSocket transport owns connection plumbing, not storage objects.
- Protocol adapters own wire conversion, not repository behavior.

## 4. Candidate Value Types

A later implementation gate may rename or reduce these shapes, but the first repository interface implementation should start from this vocabulary:

```yaml
candidate_value_types:
  - StorageObject
  - StorageObjectOwner
  - StorageObjectIdentity
  - StorageObjectValue
  - StorageObjectVersion
  - StorageObjectStatus
  - CreateStorageObjectInput
  - GetStorageObjectInput
  - ListStorageObjectsInput
  - UpdateStorageObjectInput
  - DeleteStorageObjectInput
  - StorageObjectConflict
  - StorageObjectRepositoryError
```

First-posture field vocabulary:

```yaml
storage_object_record:
  object_id: server_generated_opaque_id
  owner_kind: player
  owner_id: validated_player_id_handoff
  collection: bounded_text_identifier
  object_key: bounded_text_identifier
  value_json: small_json_object_payload
  version: server_managed_bigint_revision
  created_at: server_timestamp
  updated_at: server_timestamp
  deleted_at: nullable_server_timestamp
```

Rules:

- `owner_kind` must remain a closed vocabulary. The first allowed owner kind is `player`.
- `owner_id` is input identity handoff, not proof. Future runtime behavior must derive it from validated request identity before calling the repository.
- `collection` and `object_key` must remain bounded identifiers; they are not path names, file names, bucket names, or SQL fragments.
- `value_json` is not log-safe by default.
- `version` is server-managed and must not be accepted as client-authoritative state.
- `deleted_at` expresses soft-delete storage state only; public delete semantics remain a later behavior decision.

## 5. Candidate Repository Capabilities

The first storage-neutral capability family is:

```yaml
candidate_repository_capabilities:
  - CreateStorageObject
  - GetStorageObject
  - ListStorageObjects
  - UpdateStorageObject
  - DeleteStorageObject
```

Capability rules:

- `CreateStorageObject` may create a row for an already validated owner identity and server-validated collection/key/value.
- `GetStorageObject` is a storage lookup. It must not authenticate users, validate access tokens, or create request identity.
- `ListStorageObjects` must be collection-scoped and pagination-ready. It must not become arbitrary owner search or admin inspection without a later gate.
- `UpdateStorageObject` must support a future expected-version handoff and must return a typed conflict result rather than leaking storage driver errors.
- `DeleteStorageObject` must be version-aware when the caller supplies expected version, and it must preserve not-found and owner-mismatch leakage rules defined by future behavior gates.
- All methods must return typed module-owned records and errors, not raw SQL rows or database driver errors.

The future repository interface may choose shorter names, but it must preserve the semantic split between create, read, list, update, delete, and conflict handling.

## 6. Version And Conflict Handoff

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
  - object_already_exists
  - object_not_found
  - version_mismatch
  - invalid_expected_version
  - owner_scope_mismatch
  - deleted_object
  - storage_unavailable
```

Rules:

- Repository methods may later distinguish internal typed conflicts, but public protocol error mapping remains deferred.
- Version equality is not authentication proof.
- A stale expected version must not be collapsed into a hidden successful write.
- Owner mismatch must not reveal another player's object existence.
- The PostgreSQL adapter must map unique-index and affected-row outcomes into typed repository conflicts without exposing driver error text.

## 7. Redaction And Logging

Storage objects are game state, not authentication state, but they still require redaction discipline.

```yaml
value_log_safe: false
collection_log_safe: conditional_after_validation
object_key_log_safe: conditional_after_validation
object_id_log_safe: conditional_after_validation
owner_id_log_safe: false_by_default
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
  - blob_bytes
  - file_path
  - s3_bucket
  - s3_object_key
```

Rules:

- Repository errors must be redacted and typed.
- Raw values and storage driver errors must not be logged by default.
- `value_json` may contain player state and user-adjacent data; future behavior must treat it as non-log-safe unless a later redaction policy narrows that.
- The repository must not store or return authentication material, token material, verifier digests, transport metadata, file/blob payloads, or S3-compatible object references.

## 8. PostgreSQL Adapter Expectations

The future PostgreSQL adapter may later map the repository to:

```yaml
logical_table: storage_objects
active_identity_unique_index: storage_objects_active_identity_uq
owner_collection_index: storage_objects_owner_collection_idx
updated_at_index: storage_objects_updated_at_idx
soft_delete_column: deleted_at
version_column: version
```

Adapter expectations:

- SQL execution belongs under `runtime/internal/platform/persistence/postgres`.
- Unit-of-work and transaction handoff must follow the existing platform transaction boundary.
- SQL must not leak into `runtime/internal/modules/storage`.
- The adapter must preserve active-object uniqueness over `owner_kind + owner_id + collection + object_key`.
- Updates and deletes must be affected-row checked and version-aware when expected version is supplied.
- Soft-deleted rows must not be returned by ordinary active-object lookups unless a future admin/operations gate defines that behavior.
- Adapter tests should cover unique conflicts, not found, stale version, soft delete, owner/collection listing, value copying, timestamp mapping, and redacted errors after the adapter gate is accepted.

## 9. Relationship To Runtime, Protocol, And Authentication

This boundary does not change runtime or protocol behavior:

```yaml
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
authentication_session_behavior_changed: false
route_protection_semantics_changed: false
websocket_transport_credential_neutral: true
```

Rules:

- Future runtime behavior must derive owner identity from validated request identity before repository calls.
- The repository must not parse access tokens, validate sessions, inspect WebSocket metadata, or set request identity.
- The repository must not expose Protobuf request/response types.
- Future protocol routes remain planning candidates only until a protocol gate defines exact route names, message shapes, generated output, and error mapping.
- This boundary does not authorize public read/write permissions, ACLs, admin search, cross-player access, or direct Nakama/Pitaya API compatibility.

## 10. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - storage_objects_are_core_durable_game_state
  - collection_key_identity_supports_prototype_state
  - version_conflicts_need_first_class_semantics
adapted_concepts:
  - repository_is_vibit_storage_neutral_module_boundary
  - owner_identity_comes_from_validated_request_identity
  - public_api_compatibility_is_not_a_goal
deferred_concepts:
  - public_or_permission_bit_storage_objects
  - admin_storage_object_search
  - multi_owner_and_group_storage
  - direct_client_api_compatibility
rejected_for_now:
  - direct_nakama_storage_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - handlers_should_not_own_persistence_details
  - routing_and_storage_boundaries_should_remain_separate
adapted_concepts:
  - storage_repository_is_module_owned_not_transport_owned
  - durable_state_handoff_stays_below_protocol_adapter
deferred_concepts:
  - cluster_safe_storage_routing
  - distributed_cache_or_document_store
  - frontend_backend_rpc_storage_facade
rejected_for_now:
  - direct_pitaya_handler_or_rpc_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

## 11. Future Implementation Queue

After this boundary, future work should remain split:

```yaml
future_work_items:
  storage_objects_repository_interface_implementation:
    may_add:
      - runtime/internal/modules/storage
      - storage-neutral repository types
      - focused unit tests
    must_not_add:
      - PostgreSQL adapter
      - SQL execution
      - protocol routes
      - runtime behavior
  storage_objects_postgresql_adapter_gate:
    may_define:
      - adapter ownership
      - transaction handoff
      - SQL query shape
      - adapter tests
  storage_objects_postgresql_adapter_implementation:
    may_add:
      - platform persistence adapter
      - unit-of-work factory handoff if separately authorized
  storage_objects_runtime_behavior_gate:
    may_define:
      - owner identity derivation
      - permissions
      - validation
      - conflict semantics
      - route-policy expectations
  storage_objects_protocol_gate:
    may_define:
      - Protobuf source
      - generated output
      - route names
      - error mapping
```

Do not combine these into one broad storage subsystem slice without a new ADR.

## 12. Stop Conditions

Stop and require a later bounded work item before any of the following:

- creating `runtime/internal/modules/storage` source files;
- adding repository interface implementation;
- adding PostgreSQL adapter files or SQL execution;
- changing migrations;
- adding runtime handlers or application behavior;
- adding protocol routes, Protobuf sources, or generated output;
- adding dependencies;
- changing authentication/session behavior or route protection;
- adding broad operations/admin behavior;
- adding hosted deployments;
- creating release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, or SDK packages;
- executing public announcements beyond the GitHub release record;
- running paid promotion;
- adding large object/blob storage or S3-compatible object storage;
- adding direct Nakama/Pitaya API compatibility.

## 13. Verification

Repository verification for this boundary is:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit inspect rule runtime.storage_objects_repository_boundary
node tools/vibit check change define-storage-objects-repository-boundary --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./...
git diff --check
```

The repository check rule is:

```yaml
runtime.storage_objects_repository_boundary
```
