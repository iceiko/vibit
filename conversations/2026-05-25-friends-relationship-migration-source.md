# Conversation: Friends Relationship Migration Source

Date: 2026-05-25
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-24-add-friends-relationship-migration-source/`

Related artifacts:

- `runtime/migrations/postgres/000007_create_friend_relationships.sql`
- `decisions/ADR-0141-friends-relationship-migration-source.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Context

`W-0232` defined the friends relationship persistence schema gate. That gate selected PostgreSQL, `friend_relationships`, canonical unordered player pair identity, lifecycle states, actor-specific block timestamps, relationship versioning, timestamp checks, uniqueness and index posture, event/audit deferral, and future repository/adapter ownership candidates.

The work queue identified `W-0233 Add friends relationship migration source` as next-ready. The maintainer continued the Nakama-first direction and the product goal that user requirements should flow through AI-written specs, tests, implementation, verification, and durable project memory.

## Maintainer Narrative

The maintainer said:

```text
继续
```

Earlier direction established that continuation means advancing the next-ready work item, staying aligned with Nakama, keeping Pitaya deferred, and committing and pushing after verification unless an ask-first boundary or failure blocks the work.

## Agent Response Summary

The agent added `runtime/migrations/postgres/000007_create_friend_relationships.sql`. The migration creates only the PostgreSQL `friend_relationships` current-state table and indexes described by the schema gate.

The agent also recorded `ADR-0141`, added change-spec artifacts, updated work-item and architecture manifests, extended `tools/vibit` static checks, and advanced the queue to `W-0234 Define friends relationship repository boundary`.

The agent did not add repository interfaces, PostgreSQL adapters, runtime friendship behavior, protocol routes, Protobuf sources, generated output, dependencies, automatic startup migration behavior, event/audit tables, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## Decisions

- Complete `M-161/W-0233` as a migration-source-only slice.
- Accept `ADR-0141`.
- Add `runtime.friends_relationship_migration_source` as the repository check rule.
- Add only `runtime/migrations/postgres/000007_create_friend_relationships.sql`.
- Advance the queue to `M-162/W-0234 Define friends relationship repository boundary`.
- Continue to defer friends relationship repositories, adapters, protocol, runtime behavior, event/audit tables, groups, parties, chat, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, and direct compatibility.

## Artifacts

- `changes/2026-05-24-add-friends-relationship-migration-source/`
- `runtime/migrations/postgres/000007_create_friend_relationships.sql`
- `decisions/ADR-0141-friends-relationship-migration-source.md`
- `conversations/2026-05-25-friends-relationship-migration-source.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- Storage-neutral repository interface shape remains deferred.
- PostgreSQL adapter mapping and unit-of-work exposure remain deferred.
- Runtime command/query behavior and conflict handling remain deferred.
- Protocol route and Protobuf payload shape remain deferred.
- Event/audit history remains deferred.
- Groups, parties, chat targeting, matchmaking filters, and match social context remain separate future capability work.

## Follow-Up

- Start `W-0234 Define friends relationship repository boundary`.
- Keep the next slice gate-only unless explicitly authorized to implement repository interfaces.
- Run repository checks and record any accepted warnings explicitly.

## Redaction Notes

No secrets, raw device credentials, raw access tokens, lookup digests, verifier digests, verifier keys, DSNs, cookies, query strings, WebSocket subprotocol values, remote addresses, private relationship rows, or GitHub tokens are recorded in this conversation log.
