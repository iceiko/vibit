# Conversation: Friends Relationship Repository Boundary

Date: 2026-05-25
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-25-define-friends-relationship-repository-boundary/`

Related artifacts:

- `docs/friends-relationship-repository-boundary.md`
- `docs/friends-relationship-repository-boundary.zh-CN.md`
- `decisions/ADR-0142-friends-relationship-repository-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0233` added the friends relationship migration source at `runtime/migrations/postgres/000007_create_friend_relationships.sql` and opened `W-0234 Define friends relationship repository boundary`.

The maintainer wants vibit to continue toward a Nakama-first server framework. Pitaya remains useful as a future distributed architecture reference, but it should not drive the near-term single-process prototype path.

The product purpose remains AI-native development and testing:

```text
user requirement -> AI-written bounded spec -> acceptance criteria -> test plan -> tests -> implementation -> verification -> durable project memory
```

## Maintainer Narrative

The maintainer said:

```text
继续
```

Earlier direction established that continuation means advancing the next-ready work item, staying aligned with Nakama, keeping Pitaya deferred, and committing and pushing after verification unless an ask-first boundary or failure blocks the work.

## Agent Response Summary

This slice defines the storage-neutral friends relationship repository boundary after `W-0233` added `runtime/migrations/postgres/000007_create_friend_relationships.sql`.

The boundary records the future owner `runtime/internal/modules/friends`, the future interface candidate `runtime/internal/modules/friends.Repository`, candidate value types, candidate repository capabilities, canonical pair identity, request identity handoff, conflict classes, redaction, PostgreSQL adapter expectations, and future implementation queue.

This slice does not add repository implementation, PostgreSQL adapter behavior, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, startup wiring, event/audit tables, groups, parties, chat, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Decisions

- Complete `M-162/W-0234` as a repository-boundary-only slice.
- Accept `ADR-0142`.
- Add `runtime.friends_relationship_repository_boundary` as the repository check rule.
- Define the future repository owner candidate as `runtime/internal/modules/friends`.
- Define the future interface candidate as `runtime/internal/modules/friends.Repository`.
- Keep PostgreSQL adapter behavior, runtime behavior, protocol routes, Protobuf source, generated output, dependencies, event/audit tables, SDKs, hosted surfaces, distributed runtime, and direct compatibility deferred.
- Advance the queue to `M-163/W-0235 Implement storage-neutral friends relationship repository interface`.

## Artifacts

- `docs/friends-relationship-repository-boundary.md`
- `docs/friends-relationship-repository-boundary.zh-CN.md`
- `decisions/ADR-0142-friends-relationship-repository-boundary.md`
- `conversations/2026-05-25-friends-relationship-repository-boundary.md`
- `changes/2026-05-25-define-friends-relationship-repository-boundary/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Exact Go interface method names remain deferred to `W-0235`.
- PostgreSQL adapter mapping and unit-of-work behavior remain deferred.
- Runtime command/query behavior and public conflict mapping remain deferred.
- Protocol route and Protobuf payload shape remain deferred.
- Event/audit history remains deferred.

## Follow-Up

The next ready work item is:

```text
W-0235 Implement storage-neutral friends relationship repository interface
```

The next slice may create `runtime/internal/modules/friends` and focused repository vocabulary tests only if it preserves the W-0234 stop conditions.

## Redaction Notes

No secrets, raw device credentials, raw access tokens, lookup digests, verifier digests, verifier keys, DSNs, cookies, query strings, WebSocket subprotocol values, remote addresses, private relationship rows, or GitHub tokens are recorded in this conversation log.
