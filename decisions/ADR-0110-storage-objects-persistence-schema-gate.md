# ADR-0110: Storage Objects Persistence Schema Gate

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-persistence-schema-gate/`

Related conversations:

- `conversations/2026-05-22-storage-objects-persistence-schema-gate.md`

Related artifacts:

- `docs/storage-objects-persistence-schema-gate.md`
- `docs/storage-objects-persistence-schema-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0109` defined the first storage objects behavior gate. It selected player-owned small JSON objects as the first durable game-state posture and deferred schema, migration, repository, adapter, protocol, and runtime behavior work.

The next narrow step is a schema gate that defines the future persistence posture before adding SQL. This keeps the migration-source slice inspectable and prevents storage objects from silently expanding into runtime behavior, document database adoption, blob storage, S3-compatible storage, broad admin search, or direct Nakama/Pitaya API compatibility.

## Decision

Define the storage objects persistence schema gate.

The first persistence posture is:

```text
selected store: PostgreSQL
future logical table: storage_objects
future migration source candidate: runtime/migrations/postgres/000006_create_storage_objects.sql
future repository owner candidate: runtime/internal/modules/storage
future PostgreSQL adapter owner: runtime/internal/platform/persistence/postgres
```

The gate records the table candidate, owner identity representation, collection/key constraints, JSONB value representation, BIGINT version representation, timestamp and soft-delete posture, uniqueness/index posture, redaction posture, future repository/adapter boundaries, and stop conditions.

The repository check rule for this decision is:

```text
runtime.storage_objects_persistence_schema_gate
```

The next bounded work item is:

```text
W-0203 Add storage objects migration source
```

## Decision Record

```yaml
storage_objects_persistence_schema_gate: defined
completed_work_item: W-0202
decision: ADR-0110
check_rule: runtime.storage_objects_persistence_schema_gate
source_behavior_gate_decision: ADR-0109
source_behavior_gate_standard: docs/storage-objects-behavior-gate.md
gate_standard: docs/storage-objects-persistence-schema-gate.md
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

## Alternatives Considered

- Add `runtime/migrations/postgres/000006_create_storage_objects.sql` immediately in this gate.
- Use a document database as the first storage object backend.
- Store object values as TEXT instead of JSONB.
- Store public version as an opaque string instead of BIGINT revision.
- Add JSONB search indexes or admin listing in the first schema.
- Treat storage objects as S3-compatible object storage or file/blob storage.
- Put storage objects under inventory or player module ownership.

## Rationale

PostgreSQL is already vibit's first authoritative durable store, and the alpha has established migration tooling and persistence boundaries. A single `storage_objects` table candidate is the smallest useful next step after the behavior gate.

The chosen schema posture keeps the first storage objects capability focused on small durable game-state records. JSONB supports object payloads without adding a new database dependency. BIGINT revisions are straightforward for optimistic concurrency. A partial unique active-object identity index supports player-owned collection/key behavior while leaving cleanup, tombstones, and hard delete policy for later gates.

## Agent Reasoning Summary

The agent treated `W-0202` as a schema gate only. It selected concrete future SQL shape and ownership candidates while explicitly preventing migration-source creation, repository interfaces, adapters, protocol, runtime behavior, blob storage, object-storage dependencies, and compatibility copying.

## Decision Weights

```yaml
decision_weights:
  migration_readiness: high
  prototype_usefulness: high
  schema_inspectability: high
  runtime_scope_change: none
  protocol_scope_change: none
  migration_source_added: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0202` completes as the storage objects persistence schema gate.
- The future migration source candidate is `runtime/migrations/postgres/000006_create_storage_objects.sql`.
- The next work item becomes `W-0203 Add storage objects migration source`.
- No SQL migration source, runtime behavior, protocol route, Protobuf source, generated output, dependency, repository interface, storage adapter, hosted deployment, broad operations/admin surface, authentication/session behavior, broad module implementation, release artifact, public announcement, paid promotion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility is added.

## Reversal Conditions

Revisit this decision if JSONB object storage is too loose for agent-native verification, if optimistic concurrency needs a non-numeric version before migration, if player account ownership cannot be represented safely through `player_accounts(player_id)`, or if a future ADR selects a different first storage backend.

## Follow-Up

- Add the storage objects migration source in `W-0203`.
- Keep repository interface, PostgreSQL adapter, protocol, and runtime behavior behind later bounded work items.
- Preserve ask-first boundaries for blob storage, S3-compatible object storage, direct compatibility, hosted deployment, release artifacts, public announcement, and paid promotion.
