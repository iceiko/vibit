# ADR-0112: Storage Objects Repository Boundary

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-repository-boundary/`

Related conversations:

- `conversations/2026-05-22-storage-objects-repository-boundary.md`

Related artifacts:

- `docs/storage-objects-repository-boundary.md`
- `docs/storage-objects-repository-boundary.zh-CN.md`
- `runtime/migrations/postgres/000006_create_storage_objects.sql`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-131/W-0203` added the PostgreSQL `storage_objects` migration source without repository interfaces, adapters, runtime behavior, protocol routes, or generated output.

The next step toward product-usable storage objects is to define the storage-neutral repository boundary before writing Go repository code. Mature game backend references put pressure on this boundary: Nakama makes storage objects a common durable game-state primitive, while Pitaya reinforces a separation between handlers, routing, and persistence concerns.

## Decision

Accept `docs/storage-objects-repository-boundary.md` as the gate-only standard for the future storage objects repository.

The boundary records:

- `runtime/internal/modules/storage` as the future repository owner candidate;
- `runtime/internal/modules/storage.Repository` as the future interface candidate;
- `runtime/internal/platform/persistence/postgres` as the future PostgreSQL adapter owner;
- `storage_objects` as the logical table;
- candidate value types such as `StorageObject`, `StorageObjectOwner`, `StorageObjectIdentity`, `StorageObjectValue`, `StorageObjectVersion`, create/read/list/update/delete input types, conflict types, and repository errors;
- candidate repository capabilities `CreateStorageObject`, `GetStorageObject`, `ListStorageObjects`, `UpdateStorageObject`, and `DeleteStorageObject`;
- version and conflict handoff expectations;
- redaction expectations for values, owner identifiers, object identifiers, collection/key fields, and forbidden secret/transport/blob/S3 material;
- future PostgreSQL adapter expectations;
- `W-0205 Implement storage-neutral storage objects repository interface` as the next-ready work item.

This ADR does not add repository interface implementation, PostgreSQL adapter implementation, SQL execution behavior, runtime handlers, protocol routes, Protobuf sources, generated output, dependencies, migrations, automatic startup migration behavior, broad operations/admin behavior, authentication/session behavior changes, large object/blob storage, S3-compatible object storage, hosted deployment, release artifacts, public announcements, paid promotion, broad product module expansion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement the Go repository interface immediately after the migration source.
- Combine repository interface, PostgreSQL adapter, and runtime behavior in one slice.
- Put the repository under `runtime/internal/app/storage` instead of the future storage module.
- Reuse inventory repository types for general storage objects.
- Add direct Nakama-compatible storage object APIs.
- Defer repository design until protocol routes are selected.

## Rationale

The repository boundary is the smallest safe step after the SQL migration source. It gives future agents a precise ownership and vocabulary surface before adding code, while keeping SQL, protocol, runtime behavior, and route protection out of this slice.

Using `runtime/internal/modules/storage` as the owner candidate fits the project model: storage objects are durable game-state domain behavior, not authentication/session behavior, not transport behavior, and not a PostgreSQL-specific concern. A later repository interface implementation can define storage-neutral types and tests; a later PostgreSQL adapter gate can then map those types to the existing table.

## Agent Reasoning Summary

Storage objects are becoming vibit's first general durable game-state capability beyond inventory. Before adding code, the repository boundary needs to make ownership, value types, CRUD vocabulary, version conflicts, redaction, and adapter expectations explicit and machine-checkable.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  protocol_risk: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/storage-objects-repository-boundary.md` and its Simplified Chinese translation become the source standard for the future repository interface implementation.
- `runtime.storage_objects_repository_boundary` becomes the repository check rule for this slice.
- `M-132/W-0204` completes as a repository-boundary-only milestone.
- `M-133/W-0205 Implement storage-neutral storage objects repository interface` becomes the next-ready work item.
- No Go storage module or repository interface is created by this ADR.
- No PostgreSQL adapter, runtime behavior, protocol route, generated output, migration change, dependency, hosted deployment, release artifact, public announcement, paid promotion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility is added.

## Reversal Conditions

Revisit this decision if:

- storage objects are no longer planned as a storage module-owned capability;
- the first repository implementation needs a different owner than `runtime/internal/modules/storage`;
- the PostgreSQL table shape changes before repository implementation;
- a later ADR selects a direct external API compatibility target;
- large object/blob storage or S3-compatible object storage becomes a separate accepted capability family that changes the storage object boundary.

## Follow-Up

- Implement the storage-neutral storage objects repository interface in a later bounded work item.
- Define the storage objects PostgreSQL adapter gate after repository interface semantics are explicit.
- Define storage objects runtime behavior and protocol routes only after repository and adapter boundaries are accepted.
- Keep broad object/blob storage, S3-compatible storage, admin search, public ACLs, and direct compatibility behind future gates.
