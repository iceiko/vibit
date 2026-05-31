# Conversation: Friends Relationship PostgreSQL Adapter Implementation

Date: 2026-05-26
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-friends-relationship-postgresql-adapter/`

Related artifacts:

- `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`
- `runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `decisions/ADR-0145-friends-relationship-postgresql-adapter-implementation.md`
- `changes/2026-05-26-implement-friends-relationship-postgresql-adapter/`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`

## Context

`M-164/W-0236` completed the friends relationship PostgreSQL adapter gate. It authorized a later bounded implementation under `runtime/internal/platform/persistence/postgres` and preserved runtime behavior, protocol routes, generated output, dependency changes, migration changes, authentication/session behavior, event/audit tables, broad social features, hosted surfaces, SDKs, distributed runtime, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0237 Implement friends relationship PostgreSQL adapter`.

## Maintainer Narrative

The maintainer asked:

```text
继续
```

The standing product direction was Nakama-first, with Pitaya deferred as a future distributed architecture reference, and with AI-native requirement-to-spec-to-tests-to-implementation-to-verification as the product design goal.

## Agent Response Summary

The agent advanced one bounded work item and implemented the friends relationship PostgreSQL adapter under `runtime/internal/platform/persistence/postgres` using the existing `friend_relationships` table.

The work added:

- `FriendRelationshipRepository`;
- `NewFriendRelationshipRepositoryForUnitOfWork`;
- `UnitOfWork.NewFriendRelationshipRepository`;
- SQL mapping for create/update request, pair lookup, player-scoped list, accept, reject, remove, block, and unblock;
- canonical pair handling;
- status filtering and one-row-overflow pagination;
- optional expected-version transition checks;
- actor-specific block column selection;
- row scanning through friends module normalizers;
- redacted friends module repository errors;
- focused fake-executor adapter tests;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-165/W-0237`.
- Accept `ADR-0145`.
- Add `runtime.friends_relationship_postgresql_adapter_implementation`.
- Keep friends relationship runtime behavior out of this slice.
- Select `M-166/W-0238 Define friends relationship runtime behavior gate` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable friends relationship state needs a concrete persistence path before a useful runtime friend request, accept, reject, remove, block, unblock, list, or status surface can be defined.

Pitaya guided the layering pressure: persistence concerns should remain below handlers, routes, RPC, and cluster behavior.

vibit adapted those lessons into its own model: a PostgreSQL adapter implementing a storage-neutral repository interface, with no direct public API compatibility and no runtime/protocol behavior in this slice.

## Artifacts

- `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`
- `runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `decisions/ADR-0145-friends-relationship-postgresql-adapter-implementation.md`
- `changes/2026-05-26-implement-friends-relationship-postgresql-adapter/`
- `conversations/2026-05-26-friends-relationship-postgresql-adapter-implementation.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`

## Open Questions

- Friends relationship runtime behavior remains deferred to `W-0238` and later implementation work.
- Protocol routes and Protobuf messages remain deferred.
- Permission model and route protection remain deferred.
- Event/audit tables, chat, groups, parties, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, and direct compatibility remain deferred.
- Live PostgreSQL adapter verification remains optional and unavailable in this default fake-executor slice.

## Follow-Up

- Define the friends relationship runtime behavior gate.
- Only after that gate, implement runtime behavior in a separate bounded slice.
- Only after runtime behavior is ratified, define protocol routes and generated output.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw private social graph data from a real user are recorded in this conversation log.
