# Conversation: Friends Relationship Persistence Schema Gate

Date: 2026-05-24

## Context

The maintainer asked to continue toward the Nakama target and previously clarified that Nakama should be the primary reference while Pitaya remains deferred. The product purpose is AI-native development and AI-native testing: a user requirement should become a bounded spec, acceptance criteria, tests, implementation, verification, and durable project memory.

`W-0231` completed the friends relationship lifecycle gate and opened `W-0232 Define friends relationship persistence schema gate`.

## Maintainer Narrative

The maintainer wants sustained progress on Nakama-class capability coverage, while keeping the architecture optimized for AI agents to safely turn requirements into tested implementation slices.

## Agent Response Summary

The agent treated `W-0232` as a schema gate only. It defined the future PostgreSQL friends relationship persistence posture before SQL migration source, repositories, adapters, runtime behavior, protocol routes, generated output, dependencies, or broader social features.

The selected posture uses a future `friend_relationships` current-state table, canonical unordered player pair identity, lifecycle state values, actor-specific block timestamps, relationship versioning, timestamp checks, uniqueness/index expectations, event/audit deferral, redaction, and future repository/adapter ownership candidates.

## Decisions

- Accepted `ADR-0140`.
- Registered `runtime.friends_relationship_persistence_schema_gate`.
- Selected `runtime/migrations/postgres/000007_create_friend_relationships.sql` as the future migration source candidate.
- Opened `M-161/W-0233 Add friends relationship migration source` as the next-ready work item.

## Artifacts

- `docs/friends-relationship-persistence-schema-gate.md`
- `docs/friends-relationship-persistence-schema-gate.zh-CN.md`
- `decisions/ADR-0140-friends-relationship-persistence-schema-gate.md`
- `changes/2026-05-24-define-friends-relationship-persistence-schema-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Exact SQL checks and index names are deferred to `W-0233`.
- Repository interface shape is deferred.
- PostgreSQL adapter behavior is deferred.
- Runtime behavior and protocol routes are deferred.
- Event/audit tables and outbox behavior are deferred.

## Follow-Up

- Complete `W-0233 Add friends relationship migration source`.
- Keep repository interfaces, adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, event/audit tables, chat, groups, parties, matchmaking, match runtime, hosted surfaces, SDKs, distributed runtime, and direct compatibility behind later bounded work items.

## Redaction Notes

No raw credentials, raw access tokens, verifier keys, digests, DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private social graph data were recorded beyond explicit planning vocabulary.
