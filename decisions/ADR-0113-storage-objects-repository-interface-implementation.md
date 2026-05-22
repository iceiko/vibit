# ADR-0113: Storage Objects Repository Interface Implementation

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-storage-objects-repository-interface/`

Related conversations:

- `conversations/2026-05-22-storage-objects-repository-interface-implementation.md`

Related artifacts:

- `runtime/internal/modules/storage/repository.go`
- `runtime/internal/modules/storage/repository_test.go`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-132/W-0204` defined the storage objects repository boundary after the PostgreSQL `storage_objects` migration source. The next bounded step was to turn that boundary into storage-neutral Go repository vocabulary without adding adapter behavior, runtime behavior, protocol routes, or generated output.

Nakama keeps durable storage objects as a practical backend primitive for player game state. Pitaya reinforces that persistence and handler/runtime routing concerns should stay separated. vibit adapts those lessons through module-owned repository interfaces and explicit architecture checks instead of direct public API compatibility.

## Decision

Implement the storage-neutral repository interface under:

```text
runtime/internal/modules/storage
```

The package defines:

- `runtime/internal/modules/storage.Repository`;
- `StorageObject`, `StorageObjectOwner`, `StorageObjectIdentity`, `StorageObjectValue`, and `StorageObjectVersion`;
- closed first-posture vocabulary for `OwnerKindPlayer`, `StorageObjectStatusActive`, and `StorageObjectStatusDeleted`;
- create, get, list, update, and delete input/result types;
- optimistic conflict classes including `version_mismatch`;
- redacted repository error types;
- normalization helpers for records, owners, identities, values, and repository inputs;
- focused tests for neutrality, validation, copying, redaction, and forbidden material.

Add the first storage module manifest and module AGENTS guides so future agents can discover ownership before adding adapters or runtime behavior.

This ADR does not add PostgreSQL storage adapters, SQL execution behavior, unit-of-work factory wiring, runtime handlers, protocol routes, Protobuf sources, generated output, dependencies, migration changes, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Define the PostgreSQL adapter gate before writing the repository interface.
- Implement repository interface and PostgreSQL adapter in one slice.
- Place the interface under `runtime/internal/app/storage`.
- Reuse inventory repository abstractions for general storage objects.
- Add public storage object commands or Protobuf messages immediately.
- Add direct Nakama-compatible storage object APIs.

## Rationale

The repository boundary already selected the owner candidate and capability vocabulary. Implementing only the storage-neutral interface now reduces future adapter ambiguity while keeping SQL and runtime behavior behind later gates.

Putting the interface in `runtime/internal/modules/storage` makes storage objects a first-class domain module without making it own player accounts, authentication, sessions, protocol framing, transport behavior, or blob/S3 storage.

## Agent Reasoning Summary

The safest continuation from `W-0204` was an interface-only code slice. It gives later PostgreSQL adapter work a stable typed contract, adds tests for redaction and storage neutrality, and preserves the architectural stop conditions that keep protocol/runtime behavior from leaking into this package.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  adapter_risk: deferred
  protocol_risk: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/internal/modules/storage/repository.go` exists.
- `runtime/internal/modules/storage/repository_test.go` exists.
- `modules/storage/module.yaml` and paired module guides exist.
- `runtime.storage_objects_repository_interface_implementation` becomes the repository check rule for this slice.
- `M-133/W-0205` is completed.
- `M-134/W-0206 Define storage objects PostgreSQL adapter gate` becomes the next-ready work item.
- Existing runtime behavior, protocol behavior, migrations, and dependencies are unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- storage objects stop being a storage module-owned capability;
- the first adapter needs a different repository owner;
- the `storage_objects` table shape changes in a way that invalidates the value vocabulary;
- a later ADR selects a different conflict model;
- large object/blob or S3-compatible storage becomes a separate accepted capability that changes this repository boundary;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Define the storage objects PostgreSQL adapter gate.
- Implement the adapter only after the gate is accepted.
- Define runtime behavior, permissions, protocol routes, and generated output only after repository and adapter boundaries are accepted.
- Keep admin search, public ACLs, cross-owner scopes, group/guild scopes, large blob storage, S3-compatible storage, and direct compatibility behind future gates.
