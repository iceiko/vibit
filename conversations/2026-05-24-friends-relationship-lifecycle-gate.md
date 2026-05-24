# Conversation: Friends Relationship Lifecycle Gate

Date: 2026-05-24
Related changes:

- `changes/2026-05-24-define-friends-relationship-lifecycle-gate/`

Related artifacts:

- `docs/friends-relationship-lifecycle-gate.md`
- `docs/friends-relationship-lifecycle-gate.zh-CN.md`
- `decisions/ADR-0139-friends-relationship-lifecycle-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The maintainer asked to continue toward the Nakama target and previously clarified that Pitaya should be deferred for now. The product purpose is AI-native development and AI-native testing: user requirement to spec, acceptance criteria, tests, implementation, verification, and durable memory.

`W-0230` completed a scaffolded Nakama feature request intake pilot, selected `friends_groups_and_parties`, and opened `W-0231 Define friends relationship lifecycle gate`.

## Maintainer Narrative

The maintainer wants vibit to move toward Nakama-class capability coverage while preserving the architecture's purpose: AI should help users turn requirements into implemented and tested backend behavior. The maintainer asked to keep the current target Nakama-first rather than balancing Nakama and Pitaya equally.

## Agent Response Summary

The agent completed `W-0231` as a semantic gate only. The gate defines future friend request, accept, reject, remove, block, unblock, list, and relationship-status semantics. It records future command, query, event, error, permission, actor-relative state, invariant, redaction, and test-plan vocabulary.

The agent opened `W-0232 Define friends relationship persistence schema gate` as the next bounded work item because social graph state should have a stable persistence posture before migration, repository, runtime behavior, or protocol work.

## Decisions

- Accept `ADR-0139`.
- Register `runtime.friends_relationship_lifecycle_gate`.
- Keep Nakama as the primary product capability reference.
- Keep Pitaya deferred as a future distributed architecture reference.
- Treat this slice as gate-only.
- Require validated player identity for every future friendship command/query.
- Treat metadata-only `player_id` and `session_id` as insufficient proof.
- Defer runtime behavior, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, broader social features, hosted surfaces, distributed runtime, and direct compatibility.

## Artifacts

- Added `docs/friends-relationship-lifecycle-gate.md`.
- Added `docs/friends-relationship-lifecycle-gate.zh-CN.md`.
- Added `decisions/ADR-0139-friends-relationship-lifecycle-gate.md`.
- Filled `changes/2026-05-24-define-friends-relationship-lifecycle-gate/`.
- Updated architecture manifests, roadmap docs, agent guides, `tools/vibit`, and `rules/check-rules.json`.

## Open Questions

- Exact persistence table names.
- Canonical pair identity encoding.
- Whether rejected and removed states are current state, tombstone state, audit-only facts, or derived history.
- Duplicate request idempotency behavior.
- Conflict resolution for simultaneous request, accept, reject, remove, block, and unblock operations.
- Retention and hard-delete policy for social graph records.

## Follow-Up

- Complete `W-0232 Define friends relationship persistence schema gate`.
- Keep implementation and protocol work behind later bounded work items.

## Redaction Notes

No raw device credentials, raw access tokens, verifier keys, credential or token digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private user data were recorded. Future relationship graph details are not log-safe by default.
