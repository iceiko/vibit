# Conversation: Storage Objects PostgreSQL Adapter Gate

Date: 2026-05-22
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-postgresql-adapter-gate/`

Related artifacts:

- `docs/storage-objects-postgresql-adapter-gate.md`
- `docs/storage-objects-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0114-storage-objects-postgresql-adapter-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Context

`M-133/W-0205` completed the storage-neutral storage objects repository interface under `runtime/internal/modules/storage`. It preserved PostgreSQL adapter implementation, SQL execution, unit-of-work factory wiring, runtime behavior, protocol routes, generated output, dependency changes, migration changes, authentication/session behavior, large object/blob storage, S3-compatible object storage, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0206 Define storage objects PostgreSQL adapter gate`.

## Maintainer Narrative

The maintainer asked:

```text
使用superpowers继续推进。
```

The agent used the locally installed Superpowers workflow discipline manually, including gate-focused execution and verification-before-completion.

## Agent Response Summary

The agent advanced one bounded work item by defining the PostgreSQL adapter gate after the storage-neutral storage objects repository interface.

The gate records future adapter ownership, constructor and unit-of-work expectations, SQL mapping posture, transaction handoff, error mapping, tests, and stop conditions.

The work added:

- the English and Simplified Chinese gate standard;
- `ADR-0114`;
- a change spec package for `define-storage-objects-postgresql-adapter-gate`;
- `runtime.storage_objects_postgresql_adapter_gate` catalog and runtime check coverage;
- storage module manifest and guide updates;
- architecture, work item, reference, public next-work, and agent-guide continuation updates.

## Decisions

- Accept `docs/storage-objects-postgresql-adapter-gate.md` as the gate standard.
- Record `ADR-0114`.
- Register `runtime.storage_objects_postgresql_adapter_gate`.
- Complete `M-134/W-0206`.
- Open `M-135/W-0207 Implement storage objects PostgreSQL adapter` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable storage objects should have a concrete PostgreSQL persistence path before runtime-facing behavior.

Pitaya guided the layering pressure: storage persistence should remain below handlers, routes, RPC, and cluster behavior.

vibit adapted those lessons into its own model: a PostgreSQL adapter gate for the existing storage-neutral repository interface, without direct public API compatibility or protocol/runtime behavior in this slice.

## Artifacts

- `docs/storage-objects-postgresql-adapter-gate.md`
- `docs/storage-objects-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0114-storage-objects-postgresql-adapter-gate.md`
- `changes/2026-05-22-define-storage-objects-postgresql-adapter-gate/`
- `conversations/2026-05-22-storage-objects-postgresql-adapter-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Open Questions

- PostgreSQL adapter implementation remains deferred to `W-0207`.
- Unit-of-work factory wiring remains deferred.
- Runtime storage object handlers remain deferred.
- Protocol routes and Protobuf messages remain deferred.
- Permission model and route protection remain deferred.
- Admin search, public ACLs, cross-owner scopes, group/guild scopes, large object/blob storage, S3-compatible storage, and direct compatibility remain deferred.

## Follow-Up

- Implement the storage objects PostgreSQL adapter in `W-0207`.
- Preserve runtime behavior, protocol routes, generated output, dependency changes, migration changes, authentication/session behavior, hosted deployment, release artifact, announcement, promotion, blob/S3 storage, and direct compatibility deferrals unless a later bounded work item explicitly authorizes them.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw storage object values from a real user are recorded in this conversation log.
