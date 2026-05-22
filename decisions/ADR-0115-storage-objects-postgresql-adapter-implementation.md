# ADR-0115: Storage Objects PostgreSQL Adapter Implementation

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-storage-objects-postgresql-adapter/`

Related conversations:

- `conversations/2026-05-22-storage-objects-postgresql-adapter-implementation.md`

Related artifacts:

- `runtime/internal/platform/persistence/postgres/storage_object_repository.go`
- `runtime/internal/platform/persistence/postgres/storage_object_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `modules/storage/module.yaml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-134/W-0206` defined the storage objects PostgreSQL adapter gate after the storage-neutral repository interface. The gate authorized a later bounded implementation under `runtime/internal/platform/persistence/postgres` and kept runtime handlers, protocol routes, generated output, dependencies, migration changes, and broader product behavior deferred.

The existing storage object migration source defines the `storage_objects` table. The storage module already owns `runtime/internal/modules/storage.Repository`, value types, normalizers, optimistic conflict vocabulary, and redacted repository errors.

Nakama keeps durable storage-object-like game state as a practical backend primitive. Pitaya reinforces keeping persistence below handlers and route/runtime behavior. vibit adapts those lessons through a platform adapter implementing a storage-neutral repository interface, not through direct public API compatibility.

## Decision

Implement the storage objects PostgreSQL adapter under:

```text
runtime/internal/platform/persistence/postgres
```

The implementation adds:

- `StorageObjectRepository`;
- `NewStorageObjectRepositoryForUnitOfWork`;
- `UnitOfWork.NewStorageObjectRepository`;
- create/read/list/update/delete mapping to the existing `storage_objects` table;
- active-row filtering through `deleted_at IS NULL`;
- deterministic list ordering by `object_key`;
- bounded pagination using one-row-overflow cursor detection;
- expected-version checks for update and delete;
- server-side version increment for update and delete;
- soft-delete behavior through `deleted_at = now()`;
- row scanning through storage module normalizers;
- redacted mapping to storage module repository errors;
- focused fake-executor tests.

This ADR does not add storage object runtime handlers, protocol routes, Protobuf sources, generated output, dependency changes, migration changes, automatic startup migration behavior, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Wait for a live PostgreSQL integration environment before implementing the adapter.
- Implement runtime storage behavior together with the adapter.
- Add protocol routes and Protobuf messages together with persistence.
- Put SQL execution under `runtime/internal/modules/storage`.
- Add a new storage dependency or object-storage SDK.
- Copy Nakama storage object API behavior directly.

## Rationale

The repository interface and SQL migration already exist, and the gate narrowed the adapter owner and SQL mapping expectations. Implementing the adapter now gives future runtime behavior a real persistence target while preserving the later gates for request identity, permission, route policy, protocol, and generated-output decisions.

Fake-executor tests follow the existing PostgreSQL adapter pattern in this repository. They verify query shape, argument normalization, row scanning, redaction, version-aware conflict handling, and absence of transaction-control SQL without making default checks depend on a live database.

## Agent Reasoning Summary

The safest continuation from `W-0206` was a platform adapter slice only. It provides durable storage-object persistence behind `storage.Repository`, keeps SQL out of the storage module, and leaves runtime and protocol behavior to later gates.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: medium
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/internal/platform/persistence/postgres/storage_object_repository.go` exists.
- `runtime/internal/platform/persistence/postgres/storage_object_repository_test.go` exists.
- `UnitOfWork.NewStorageObjectRepository` exists as an adapter handoff helper.
- `runtime.storage_objects_postgresql_adapter_implementation` becomes the repository check rule for this slice.
- `M-135/W-0207` is completed.
- `M-136/W-0208 Define storage objects runtime behavior gate` becomes the next-ready work item.
- Runtime handlers, protocol routes, generated output, migrations, dependencies, and authentication/session semantics remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- the storage repository interface changes materially;
- the `storage_objects` migration source changes materially;
- a later ADR selects a different pagination or conflict model;
- the project requires live PostgreSQL adapter tests as default verification;
- storage objects become large blob or S3-compatible object storage instead of small JSON game state;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Define the storage objects runtime behavior gate before adding runtime services or handlers.
- Define protocol routes and Protobuf messages only after runtime behavior is ratified.
- Keep admin search, public ACLs, cross-owner scopes, group/guild scopes, large blob storage, S3-compatible storage, and direct compatibility behind future gates.
