# Conversation: Storage Objects Repository Interface Implementation

Date: 2026-05-22
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-storage-objects-repository-interface/`

Related artifacts:

- `runtime/internal/modules/storage/repository.go`
- `runtime/internal/modules/storage/repository_test.go`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `decisions/ADR-0113-storage-objects-repository-interface-implementation.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `rules/check-rules.json`
- `tools/vibit`

## Context

`M-132/W-0204` completed the storage objects repository boundary. It recorded `runtime/internal/modules/storage` as the future repository owner candidate and kept adapter, SQL, runtime, protocol, generated output, dependency, migration, authentication/session, blob/S3 storage, and direct compatibility concerns deferred.

The next-ready work item was `W-0205 Implement storage-neutral storage objects repository interface`.

## Maintainer Narrative

The maintainer asked:

```text
使用superpowers继续推进。
```

The agent used the Superpowers workflow discipline manually from the installed plugin files, including test-driven development and verification-before-completion.

## Agent Response Summary

The agent continued the selected bounded work item and implemented the storage-neutral repository interface under `runtime/internal/modules/storage`.

The work added:

- storage object value types and closed first-posture owner/status vocabulary;
- create/read/list/update/delete repository input and result types;
- optimistic conflict and redacted repository error vocabulary;
- normalization helpers for values, owner identity, object identity, records, and repository inputs;
- focused Go tests;
- storage module manifest and module AGENTS guides;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-133/W-0205`.
- Accept `ADR-0113`.
- Add `runtime.storage_objects_repository_interface_implementation`.
- Register the `storage` module in `.arch/modules.yaml`.
- Keep the first implementation storage-neutral and module-owned.
- Select `M-134/W-0206 Define storage objects PostgreSQL adapter gate` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable storage objects are a common practical game backend primitive.

Pitaya guided the layering pressure: persistence concerns should remain below handlers, routes, RPC, and cluster behavior.

vibit adapted those lessons into its own model: an explicit module-owned repository interface with no direct public API compatibility and no runtime/protocol behavior in this slice.

## Artifacts

- `runtime/internal/modules/storage/repository.go`
- `runtime/internal/modules/storage/repository_test.go`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `decisions/ADR-0113-storage-objects-repository-interface-implementation.md`
- `changes/2026-05-22-implement-storage-objects-repository-interface/`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- PostgreSQL adapter mapping remains deferred to `W-0206` and later implementation work.
- Runtime storage object handlers remain deferred.
- Protocol routes and Protobuf messages remain deferred.
- Permission model and route protection remain deferred.
- Admin search, public ACLs, cross-owner scopes, group/guild scopes, large object/blob storage, S3-compatible storage, and direct compatibility remain deferred.

## Follow-Up

- Define the storage objects PostgreSQL adapter gate.
- Only after that gate, implement the adapter in a separate bounded slice.
- Only after repository and adapter boundaries, define runtime behavior and protocol routes.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw storage object values from a real user are recorded in this conversation log.
