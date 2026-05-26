# Conversation: Friends Relationship PostgreSQL Adapter Gate

Date: 2026-05-26  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-26-define-friends-relationship-postgresql-adapter-gate/`

Related artifacts:

- `docs/friends-relationship-postgresql-adapter-gate.md`
- `docs/friends-relationship-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0144-friends-relationship-postgresql-adapter-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Context

`M-163/W-0235` had implemented the storage-neutral friends relationship repository interface. The next-ready work item was `W-0236 Define friends relationship PostgreSQL adapter gate`, which needed to define the PostgreSQL adapter boundary before any adapter implementation or SQL execution behavior could be added.

## Maintainer Narrative

The maintainer repeatedly asked to continue:

```text
继续
```

The standing product direction was Nakama-first, with Pitaya deferred as a future distributed architecture reference, and with AI-native requirement-to-spec-to-tests-to-implementation-to-verification as the product design goal.

## Agent Response Summary

The agent defined the friends relationship PostgreSQL adapter gate without adding adapter implementation, SQL execution behavior, runtime friendship behavior, protocol routes, Protobuf source, generated output, migration changes, event/audit tables, broader social features, hosted or SDK scope, distributed runtime, or direct Nakama/Pitaya API compatibility.

The work added the gate standard, Simplified Chinese translation, ADR, change record, rule registration, manifest updates, module guide updates, and next-ready queue entry for `M-165/W-0237 Implement friends relationship PostgreSQL adapter`. The future adapter gate maps to the existing `friend_relationships` logical table from migration source `runtime/migrations/postgres/000007_create_friend_relationships.sql`.

## Decisions

- Complete `M-164/W-0236`.
- Accept `ADR-0144`.
- Add `runtime.friends_relationship_postgresql_adapter_gate`.
- Keep Nakama as the primary product capability reference for friends relationship persistence.
- Keep Pitaya deferred as a future distributed architecture reference.
- Select `M-165/W-0237 Implement friends relationship PostgreSQL adapter` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable friends relationship storage is required before a Nakama-class social graph can expose friend request, accept, reject, remove, block, unblock, list, or status behavior.

Pitaya remained deferred; no distributed topology, RPC routing, frontend/backend split, group broadcast, or service discovery behavior was added.

vibit adapted those references into its own model: an explicit platform-owned PostgreSQL adapter gate below the module-owned repository interface, with no direct public API compatibility.

## Artifacts

- `docs/friends-relationship-postgresql-adapter-gate.md`
- `docs/friends-relationship-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0144-friends-relationship-postgresql-adapter-gate.md`
- `changes/2026-05-26-define-friends-relationship-postgresql-adapter-gate/`
- `conversations/2026-05-26-friends-relationship-postgresql-adapter-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- PostgreSQL adapter implementation remains deferred to `W-0237`.
- Runtime friend request, accept, reject, remove, block, unblock, list, and status behavior remain deferred.
- Protocol routes and Protobuf messages remain deferred.
- Permission model and route protection remain deferred.
- Event/audit tables, chat, groups, parties, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, and direct compatibility remain deferred.

## Follow-Up

The next-ready work item is:

```text
M-165/W-0237 Implement friends relationship PostgreSQL adapter
```

That next slice may implement the PostgreSQL adapter only within the gate boundaries accepted by `ADR-0144`.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw private social graph data from a real user are recorded in this conversation log.
