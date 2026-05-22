# Conversation: Storage Objects Repository Boundary

Date: 2026-05-22
Participants: Maintainer, Agent
Related work item: `W-0204 Define storage objects repository boundary`
Related decision: `ADR-0112`
Related change: `changes/2026-05-22-define-storage-objects-repository-boundary/`

## Context

`W-0203` completed the storage objects migration-source-only slice and pushed `runtime/migrations/postgres/000006_create_storage_objects.sql`.

The next-ready work item was `W-0204 Define storage objects repository boundary`. Its ask-first boundary permits a gate-only repository boundary definition but does not permit repository interface implementation, PostgreSQL adapter implementation, runtime behavior, protocol routes, Protobuf or generated output changes, dependencies, migration changes, hosted deployment, public announcements, paid promotion, release artifacts, large object/blob storage, S3-compatible storage, broad product module expansion, or direct Nakama/Pitaya API compatibility.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Earlier direction established that continuation means advancing one next-ready work item, committing, and pushing after verification unless an ask-first boundary or failure blocks the work.

## Agent Response Summary

Define the storage objects repository boundary as a gate-only artifact.

The boundary records future ownership, candidate value types, create/read/list/update/delete vocabulary, version conflict handoff, redaction posture, PostgreSQL adapter expectations, and stop conditions. It opens `W-0205 Implement storage-neutral storage objects repository interface` as the next-ready work item.

The agent did not create `runtime/internal/modules/storage`, add repository interfaces, add PostgreSQL adapters, add SQL execution behavior, change migrations, add runtime handlers, add protocol routes, add Protobuf source, add generated output, add dependencies, change authentication/session behavior, add hosted deployments, create release artifacts, execute public announcements, run paid promotion, add large object/blob storage, add S3-compatible object storage, or add direct Nakama/Pitaya API compatibility.

## Decisions

- Complete `M-132/W-0204` as a repository-boundary-only slice.
- Accept `ADR-0112`.
- Add `runtime.storage_objects_repository_boundary` as the repository check rule.
- Record future owner `runtime/internal/modules/storage`.
- Record future interface candidate `runtime/internal/modules/storage.Repository`.
- Record future PostgreSQL adapter owner `runtime/internal/platform/persistence/postgres`.
- Advance the queue to `M-133/W-0205 Implement storage-neutral storage objects repository interface`.

## Artifacts

- `docs/storage-objects-repository-boundary.md`
- `docs/storage-objects-repository-boundary.zh-CN.md`
- `decisions/ADR-0112-storage-objects-repository-boundary.md`
- `changes/2026-05-22-define-storage-objects-repository-boundary/`
- `conversations/2026-05-22-storage-objects-repository-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- Exact Go repository type names remain deferred to `W-0205`.
- PostgreSQL adapter mapping and unit-of-work exposure remain deferred.
- Runtime permission, validation, and owner identity derivation remain deferred.
- Protocol command/query shape remains deferred.
- Admin, public, cross-owner, group/guild, room, party, match, object/blob, and S3-compatible storage scopes remain deferred.

## Follow-Up

- Start `W-0205 Implement storage-neutral storage objects repository interface`.
- Keep `W-0205` interface-only unless explicitly authorized to add adapters, runtime behavior, protocol, generated output, dependencies, or migration changes.
- Run repository checks and record any accepted warnings explicitly.

## Redaction Notes

No private local environment file was read or printed. No secrets, raw credentials, raw tokens, verifier keys, DSNs, digests, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, object values, GitHub tokens, or concrete transport metadata were added to the record.

## Verification Plan

- `node -c tools/vibit`
- `node tools/vibit inspect next`
- `node tools/vibit inspect rule runtime.storage_objects_repository_boundary`
- `node tools/vibit check change define-storage-objects-repository-boundary --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `git diff --check`
- Secret scan for obvious GitHub token patterns excluding `.git/`, `.vibit.local.env`, and `node_modules/`
