# ADR-0114: Storage Objects PostgreSQL Adapter Gate

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-postgresql-adapter-gate/`

Related conversations:

- `conversations/2026-05-22-storage-objects-postgresql-adapter-gate.md`

Related artifacts:

- `docs/storage-objects-postgresql-adapter-gate.md`
- `docs/storage-objects-postgresql-adapter-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/storage/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-133/W-0205` implemented the storage-neutral storage objects repository interface under `runtime/internal/modules/storage`. The next bounded step is to define the PostgreSQL adapter gate before any adapter code, SQL execution, runtime behavior, or protocol surface is added.

The existing storage object migration source already defines the `storage_objects` table, and the repository interface already defines create/read/list/update/delete vocabulary. A separate adapter gate keeps SQL mapping, transaction handoff, error mapping, and tests explicit before implementation.

Nakama keeps durable storage-object-like game state as a central backend capability. Pitaya reinforces keeping storage and persistence below handlers and runtime routing. vibit adapts those lessons through explicit adapter boundaries and checkable manifests, not direct public API compatibility.

## Decision

Accept `docs/storage-objects-postgresql-adapter-gate.md` as the gate for the future storage objects PostgreSQL adapter.

The gate records:

- future adapter owner `runtime/internal/platform/persistence/postgres`;
- future source candidate `runtime/internal/platform/persistence/postgres/storage_object_repository.go`;
- future test candidate `runtime/internal/platform/persistence/postgres/storage_object_repository_test.go`;
- repository interface source `runtime/internal/modules/storage.Repository`;
- SQL mapping posture for `storage_objects`;
- constructor and caller-supplied executor expectations;
- unit-of-work and transaction handoff expectations;
- redacted driver-error and conflict mapping;
- focused adapter implementation test expectations;
- stop conditions before implementation, runtime behavior, protocol routes, generated output, dependencies, migration changes, hosted deployment, release artifacts, public announcements, paid promotion, object/blob storage, S3-compatible storage, or direct compatibility.

This ADR does not add PostgreSQL storage adapters, SQL execution behavior, unit-of-work factory wiring, runtime handlers, protocol routes, Protobuf sources, generated output, dependencies, migration changes, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement the PostgreSQL adapter immediately after the repository interface.
- Reuse the inventory PostgreSQL repository implementation shape without a storage-specific gate.
- Put SQL execution under `runtime/internal/modules/storage`.
- Add protocol storage routes together with adapter implementation.
- Add direct Nakama-compatible storage object APIs.

## Rationale

The storage object repository has more conflict and redaction pressure than a trivial table adapter: active identity uniqueness, soft delete, owner-scope leakage, expected version conflicts, JSON value redaction, and transaction handoff all need to be explicit before implementation.

A gate-only ADR keeps the next implementation slice bounded and makes it possible for repository checks to reject accidental SQL or protocol behavior before that work item is authorized.

## Agent Reasoning Summary

The safest continuation from `W-0205` is a platform adapter gate. It gives future implementation a precise owner, constructor posture, SQL mapping checklist, and test list while preserving the separation between storage module vocabulary, PostgreSQL adapter behavior, runtime routing, protocol shape, and product-scope expansion.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  adapter_risk: contained_by_next_gate
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/storage-objects-postgresql-adapter-gate.md` and `.zh-CN.md` exist.
- `runtime.storage_objects_postgresql_adapter_gate` becomes the repository check rule for this slice.
- `M-134/W-0206` is completed.
- `M-135/W-0207 Implement storage objects PostgreSQL adapter` becomes the next-ready work item.
- Existing runtime behavior, protocol behavior, migrations, dependencies, and generated output remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- `runtime/internal/modules/storage.Repository` changes before adapter implementation;
- the `storage_objects` migration source changes materially;
- the project selects a different first storage objects persistence engine;
- transaction ownership moves away from caller-supplied unit-of-work boundaries;
- storage objects become large blob or S3-compatible object storage instead of small JSON game state;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Implement the storage objects PostgreSQL adapter only after this gate is accepted.
- Keep runtime behavior, permissions, protocol routes, and generated output behind later gates.
- Keep admin search, public ACLs, cross-owner scopes, group/guild scopes, large blob storage, S3-compatible storage, and direct compatibility behind future gates.
