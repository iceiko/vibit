# Conversation: Storage Objects PostgreSQL Adapter Implementation

Date: 2026-05-22
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-storage-objects-postgresql-adapter/`

Related artifacts:

- `runtime/internal/platform/persistence/postgres/storage_object_repository.go`
- `runtime/internal/platform/persistence/postgres/storage_object_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `decisions/ADR-0115-storage-objects-postgresql-adapter-implementation.md`
- `changes/2026-05-22-implement-storage-objects-postgresql-adapter/`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Context

`M-134/W-0206` completed the storage objects PostgreSQL adapter gate. It authorized a later bounded implementation under `runtime/internal/platform/persistence/postgres` and preserved runtime behavior, protocol routes, generated output, dependency changes, migration changes, authentication/session behavior, blob/S3 storage, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0207 Implement storage objects PostgreSQL adapter`.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。
```

The agent continued using the locally installed Superpowers workflow discipline manually, including test-driven development and verification-before-completion.

## Agent Response Summary

The agent advanced one bounded work item and implemented the storage objects PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`.

The work added:

- `StorageObjectRepository`;
- `NewStorageObjectRepositoryForUnitOfWork`;
- `UnitOfWork.NewStorageObjectRepository`;
- create/read/list/update/delete SQL mapping for the existing `storage_objects` table;
- active-row filtering and soft-delete behavior;
- expected-version checks and server-side version increments;
- deterministic list ordering and one-row-overflow pagination;
- row scanning through storage module normalizers;
- redacted storage module repository errors;
- focused fake-executor adapter tests;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-135/W-0207`.
- Accept `ADR-0115`.
- Add `runtime.storage_objects_postgresql_adapter_implementation`.
- Keep storage object runtime behavior out of this slice.
- Select `M-136/W-0208 Define storage objects runtime behavior gate` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable storage-object-like state needs a concrete persistence path before a useful prototype runtime surface.

Pitaya guided the layering pressure: persistence concerns should remain below handlers, routes, RPC, and cluster behavior.

vibit adapted those lessons into its own model: a PostgreSQL adapter implementing a storage-neutral repository interface, with no direct public API compatibility and no runtime/protocol behavior in this slice.

## Artifacts

- `runtime/internal/platform/persistence/postgres/storage_object_repository.go`
- `runtime/internal/platform/persistence/postgres/storage_object_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `decisions/ADR-0115-storage-objects-postgresql-adapter-implementation.md`
- `changes/2026-05-22-implement-storage-objects-postgresql-adapter/`
- `conversations/2026-05-22-storage-objects-postgresql-adapter-implementation.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Open Questions

- Storage object runtime behavior remains deferred to `W-0208` and later implementation work.
- Protocol routes and Protobuf messages remain deferred.
- Permission model and route protection remain deferred.
- Admin search, public ACLs, cross-owner scopes, group/guild scopes, large object/blob storage, S3-compatible storage, and direct compatibility remain deferred.
- Live PostgreSQL adapter verification remains optional and unavailable in this default fake-executor slice.

## Follow-Up

- Define the storage objects runtime behavior gate.
- Only after that gate, implement runtime behavior in a separate bounded slice.
- Only after runtime behavior is ratified, define protocol routes and generated output.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw storage object values from a real user are recorded in this conversation log.
