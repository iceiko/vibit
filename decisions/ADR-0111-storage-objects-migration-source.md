# ADR-0111: Storage Objects Migration Source

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-add-storage-objects-migration-source/`

Related conversations:

- `conversations/2026-05-22-storage-objects-migration-source.md`

Related artifacts:

- `runtime/migrations/postgres/000006_create_storage_objects.sql`
- `docs/storage-objects-persistence-schema-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0202` defined the storage objects persistence schema gate. That gate selected PostgreSQL as the first durable store, `storage_objects` as the future logical table, JSONB as the first value representation, BIGINT as the first version representation, and `runtime/migrations/postgres/000006_create_storage_objects.sql` as the future migration source candidate.

The work queue reached `M-131/W-0203`, which authorized adding the SQL migration source only. The repository is moving from a source-first alpha toward a prototype-ready foundation, and storage objects are the first general durable game-state capability beyond the inventory proof slice.

## Decision

Add only the PostgreSQL migration source:

```text
runtime/migrations/postgres/000006_create_storage_objects.sql
```

The migration creates `storage_objects` with:

- `object_id` as the server-generated opaque record id;
- first owner posture `owner_kind = 'player'` and `owner_id` linked to `player_accounts(player_id)`;
- logical active identity over `owner_kind + owner_id + collection + object_key`;
- constrained `collection` and `object_key` text fields;
- `value_json JSONB` constrained to a top-level JSON object;
- positive server-managed `version BIGINT` with default `1`;
- created, updated, and optional soft-delete timestamps;
- active logical identity uniqueness for rows where `deleted_at IS NULL`;
- lookup and updated-at indexes.

This ADR does not add storage object repository interfaces, PostgreSQL storage adapters, runtime handlers, protocol routes, Protobuf sources, generated output, dependencies, automatic startup migration behavior, broad operations/admin behavior, authentication/session behavior changes, large object/blob storage, S3-compatible object storage, hosted deployment, additional release artifacts, public announcements, paid promotion, broad product module expansion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Define the storage object repository boundary before adding the migration source.
- Add repository interfaces and PostgreSQL adapter behavior in the same slice.
- Add protocol routes or Protobuf sources together with the table.
- Use arbitrary JSON payloads instead of requiring a top-level JSON object.
- Store large binary objects or S3-compatible object references in the same table.
- Allow global, group, party, room, match, public, or admin object scopes in the first migration.
- Use direct Nakama/Pitaya table or API compatibility as the storage target.

## Rationale

Nakama demonstrates that small durable per-user storage objects are a core game backend capability. Pitaya reinforces that durable state should remain behind explicit module and persistence boundaries rather than leaking into handlers or transport. The vibit adaptation is a small PostgreSQL migration first, with strict ownership and redaction posture, before repository, adapter, protocol, or runtime behavior.

Keeping this slice migration-only gives future agents a stable table shape while avoiding accidental product claims. A later repository boundary can define storage-neutral create/read/update/delete/conflict semantics. A later PostgreSQL adapter slice can map those semantics to this table. A later protocol/runtime slice can decide how player identity, permissions, conflict responses, and value redaction are exposed.

## Agent Reasoning Summary

The smallest product-useful step after the schema gate is the SQL source itself. It makes the planned storage object state inspectable and migratable while preserving the more sensitive behavior decisions for later bounded work.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  migration_safety: high
  agent_readability: high
  nakama_pitaya_alignment: high
  runtime_behavior_risk: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/migrations/postgres/000006_create_storage_objects.sql` exists.
- `runtime.storage_objects_migration_source` becomes the repository check rule for this slice.
- `M-131/W-0203` is completed as a migration-source-only milestone.
- The work queue advances to `M-132/W-0204 Define storage objects repository boundary`.
- Existing runtime behavior is not changed by this ADR.
- Storage objects are not yet exposed through repository interfaces, adapters, protocol routes, runtime handlers, or generated output.

## Reversal Conditions

Revisit this decision if:

- A later ADR changes the first durable storage object backing store away from PostgreSQL.
- The first repository or runtime behavior requires a different minimum identity, value, version, or soft-delete shape before external users depend on it.
- A later explicit compatibility ADR adopts a direct Nakama or Pitaya public API compatibility target.
- A product requirement for large object/blob storage or S3-compatible object storage is accepted as a separate capability family.

## Follow-Up

- Define the storage objects repository boundary.
- Define the storage objects PostgreSQL adapter gate after repository semantics are explicit.
- Define protocol and runtime behavior only after repository and permission semantics are accepted.
- Keep large object/blob storage and S3-compatible object storage behind separate future gates.
