# Storage Objects PostgreSQL Adapter Gate

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Gate-only boundary for the future PostgreSQL adapter implementing `runtime/internal/modules/storage.Repository`
Depends on: `runtime/internal/modules/storage/repository.go`, `docs/storage-objects-repository-boundary.md`, `runtime/migrations/postgres/000006_create_storage_objects.sql`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0114`

The paired Simplified Chinese translation is `docs/storage-objects-postgresql-adapter-gate.zh-CN.md`. The English file is authoritative.

This document defines the storage objects PostgreSQL adapter gate. It is a gate artifact. It does not add PostgreSQL adapter implementation, SQL execution behavior, unit-of-work factory wiring, runtime behavior, protocol routes, Protobuf source, generated output, dependencies, migrations, automatic startup migration behavior, broad operations/admin behavior, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The storage objects PostgreSQL adapter gate record is:

```yaml
storage_objects_postgresql_adapter_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0206
decision: ADR-0114
check_rule: runtime.storage_objects_postgresql_adapter_gate
source_repository_interface_decision: ADR-0113
repository_interface: runtime/internal/modules/storage.Repository
repository_interface_source: runtime/internal/modules/storage/repository.go
source_migration_source_decision: ADR-0111
source_migration_source: runtime/migrations/postgres/000006_create_storage_objects.sql
storage_objects_logical_table: storage_objects
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_adapter_source_candidate: runtime/internal/platform/persistence/postgres/storage_object_repository.go
future_adapter_tests_candidate: runtime/internal/platform/persistence/postgres/storage_object_repository_test.go
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
broad_operations_admin_behavior_added: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_adapter_implementation_work_item: W-0207
future_adapter_implementation_direction: storage_objects_postgresql_adapter_implementation
```

## 2. Purpose

`W-0205` implemented the storage-neutral `runtime/internal/modules/storage.Repository` interface. The next useful boundary is the platform adapter gate that will later map that interface to the existing PostgreSQL `storage_objects` table.

This gate records the future implementation shape before any SQL is written:

- adapter ownership;
- constructor and executor handoff expectations;
- transaction and unit-of-work boundaries;
- SQL mapping posture for create/read/list/update/delete;
- conflict and driver-error mapping;
- timestamp, JSON, version, and soft-delete mapping;
- focused test expectations;
- stop conditions that keep runtime, protocol, generated output, and product expansion out of the adapter slice.

This is not an implementation. The future adapter source path is named only so agents can verify later work against the accepted boundary.

## 3. Ownership

The future adapter owner is:

```yaml
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
repository_interface_owner: runtime/internal/modules/storage
sql_mapping_owner: runtime/internal/platform/persistence/postgres
transaction_owner: caller_supplied_unit_of_work
application_layer_owns_request_identity: true
storage_module_owns_repository_vocabulary: true
websocket_transport_owns_storage_objects: false
protocol_adapter_owns_storage_objects: false
authentication_module_owns_storage_objects: false
player_module_owns_storage_objects: false
```

Rules:

- The adapter may later implement `storage.Repository` under the PostgreSQL platform package.
- The adapter must not move SQL into `runtime/internal/modules/storage`.
- The adapter must not own request authentication, player identity validation, route policy, protocol parsing, or WebSocket state.
- The adapter must receive already-normalized repository input or call storage module normalizers before SQL mapping.
- The adapter must return storage module value types and typed repository errors, not driver-specific errors.

## 4. Future Constructor And Executor Handoff

The first adapter implementation should follow existing PostgreSQL adapter patterns:

```yaml
future_constructor_candidate: NewStorageObjectRepository
future_repository_interface: runtime/internal/modules/storage.Repository
executor_source: caller_supplied
transaction_control_sql_allowed: false
unit_of_work_handoff_required: true
connection_pool_owned_by_adapter: false
context_required: true
```

Rules:

- The constructor should accept an existing executor or query interface rather than owning a pool directly.
- The adapter must not issue `BEGIN`, `COMMIT`, or `ROLLBACK`; transaction ownership remains with the unit-of-work runner.
- The adapter must not create automatic startup migrations.
- The adapter must not add a new dependency if existing PostgreSQL platform dependencies already cover the implementation.
- Any required dependency change must be a separate dependency-adoption decision.

## 5. SQL Mapping Posture

The future adapter may map the repository methods to the `storage_objects` table:

```yaml
logical_table: storage_objects
primary_key_column: object_id
owner_columns:
  - owner_kind
  - owner_id
identity_columns:
  - collection
  - object_key
value_column: value_json
version_column: version
created_at_column: created_at
updated_at_column: updated_at
soft_delete_column: deleted_at
active_identity_unique_index: storage_objects_active_identity_uq
owner_collection_index: storage_objects_owner_collection_idx
updated_at_index: storage_objects_updated_at_idx
```

Method posture:

- `CreateStorageObject` should insert active rows only after validation, preserve `object_id`, start at positive version, and map active identity conflicts to `object_already_exists`.
- `GetStorageObject` should return only active rows for the requested owner, collection, and key.
- `ListStorageObjects` should be owner and collection scoped, apply deterministic key ordering, honor bounded limits, and use `AfterObjectKey` as an opaque-enough pagination cursor for the first implementation.
- `UpdateStorageObject` should update active rows, increment version server-side, preserve ownership, and enforce expected-version checks when supplied.
- `DeleteStorageObject` should soft-delete active rows, be version-aware when expected version is supplied, and preserve owner-mismatch leakage rules.

Rules:

- SQL text must remain inside the PostgreSQL adapter package.
- `value_json` is not log-safe and must not be included in default error text.
- Driver-specific constraint names may be used internally for error mapping, but public module errors must remain storage-neutral.
- Affected-row counts must be checked for update and delete operations.

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
- The adapter must not hide transaction failures by returning successful storage results.

## 7. Error Mapping

The future adapter must collapse PostgreSQL details into storage module errors:

```yaml
repository_error_owner: runtime/internal/modules/storage
driver_error_public_leakage_allowed: false
constraint_name_public_leakage_allowed: false
value_json_public_leakage_allowed: false
conflict_classes:
  - object_already_exists
  - object_not_found
  - version_mismatch
  - invalid_expected_version
  - owner_scope_mismatch
  - deleted_object
  - storage_unavailable
```

Mapping expectations:

- Unique active identity conflicts should map to `object_already_exists`.
- No affected row with expected version should map to a typed not-found or version conflict without leaking another owner's object existence.
- Malformed input should be rejected before SQL execution when possible.
- Unknown driver or executor failures should map to storage-unavailable style errors with redacted reasons.
- Raw SQL, DSNs, credentials, token material, verifier digests, `value_json`, and driver stack details must not appear in public error strings.

## 8. Test Expectations

The later implementation slice should add focused PostgreSQL adapter tests. Fake-executor or query-capture tests are acceptable before live database verification is available.

Required test families for the implementation slice:

```yaml
future_tests:
  - constructor_requires_executor
  - create_maps_insert_and_unique_conflict
  - get_maps_active_lookup_and_not_found
  - list_is_owner_collection_scoped_and_ordered
  - update_checks_expected_version_and_increments_version
  - delete_soft_deletes_and_checks_expected_version
  - json_value_is_copied_and_redacted
  - timestamps_are_utc_normalized
  - driver_errors_are_redacted
  - transaction_control_sql_is_absent
```

Rules:

- Tests must not require protocol routes.
- Tests must not require WebSocket transport.
- Tests must not print raw values, DSNs, credentials, tokens, or digests.
- Live PostgreSQL verification may be added later when an implementation slice authorizes it; if unavailable, it must be explicitly recorded.

## 9. Relationship To Runtime, Protocol, And Authentication

This gate does not change runtime or protocol behavior:

```yaml
runtime_storage_handlers_added: false
storage_protocol_routes_added: false
protobuf_storage_messages_added: false
generated_storage_output_added: false
authentication_session_behavior_changed: false
request_identity_handoff_changed: false
```

Rules:

- Runtime handlers remain deferred.
- Protocol routes and generated storage contract shapes remain deferred.
- Request identity validation remains owned by authentication/session boundaries.
- The adapter must not parse bearer tokens, cookies, WebSocket subprotocols, or envelope metadata.

## 10. Reference Alignment

Nakama provides durable storage-object-like game state capability. Pitaya reinforces keeping persistence below handlers and runtime routing. vibit uses those references for capability planning only:

- no direct Nakama or Pitaya API compatibility is added;
- no public storage route is added by this gate;
- no server runtime hook or admin surface is added by this gate.

## 11. Verification

Required verification after this gate:

```text
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.storage_objects_postgresql_adapter_gate
node tools/vibit check change define-storage-objects-postgresql-adapter-gate --json
node tools/vibit check module storage --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./...
git diff --check
```

Live PostgreSQL adapter verification is not required for this gate because no adapter implementation or SQL execution is added.

## 12. Stop Conditions

Stop and create a later bounded work item before adding:

- `runtime/internal/platform/persistence/postgres/storage_object_repository.go`;
- PostgreSQL storage adapter implementation;
- SQL execution behavior;
- unit-of-work factory wiring;
- runtime storage object handlers;
- protocol routes;
- Protobuf source or generated output;
- dependency changes;
- migration changes;
- automatic startup migration behavior;
- authentication/session behavior changes;
- route-protection changes;
- broad operations/admin behavior;
- hosted deployments;
- release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publication, or SDK packages;
- public announcements beyond the GitHub release record;
- paid promotion;
- large object/blob storage;
- S3-compatible object storage;
- direct Nakama/Pitaya API compatibility.
