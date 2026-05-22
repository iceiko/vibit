# Conversation: Storage Objects Behavior Gate

Date: 2026-05-22
Participants: Maintainer, Agent
Related work item: `W-0201`
Related decision: `ADR-0109`

## Context

The maintainer asked to continue advancing vibit after the prototype-ready local development path package was completed. The next-ready work item was `W-0201 Define storage objects behavior gate`.

The product objective is to keep moving toward a prototype-ready and later production-useful server framework. The next product gap is general durable game state beyond the current inventory proof slice.

## Maintainer Narrative

The maintainer wants vibit to become a real product foundation, not only a source alpha that can be inspected. Storage objects are a practical next step because a prototype author needs somewhere to store small durable player-owned game state before larger systems such as chat, friends, leaderboards, matchmaking, or match runtime are useful.

## Agent Response Summary

The agent defined the storage objects behavior gate as a planning-only, gate-only slice. The gate selects player-owned small JSON objects as the first posture, records ownership, scope/key rules, read/write behavior, permission posture, conflict semantics, protocol and data expectations, verification expectations, and stop conditions, then advances the queue to a persistence schema gate.

## Decisions

- The storage objects behavior gate is recorded in `docs/storage-objects-behavior-gate.md`.
- The paired Simplified Chinese translation is `docs/storage-objects-behavior-gate.zh-CN.md`.
- The decision record is `ADR-0109`.
- The repository check rule is `runtime.storage_objects_behavior_gate`.
- The first posture is player-owned small JSON storage objects addressed by `owner_kind + owner_id + collection + key`.
- The next bounded direction is `W-0202 Define storage objects persistence schema gate`.

## Artifacts

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/storage-objects-behavior-gate.md`
- `docs/storage-objects-behavior-gate.zh-CN.md`
- `decisions/ADR-0109-storage-objects-behavior-gate.md`
- `changes/2026-05-22-define-storage-objects-behavior-gate/`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- The exact PostgreSQL table shape, version representation, value constraints, indexes, and migration source remain deferred to `W-0202`.
- Exact Protobuf messages and route names remain planning candidates only until a later protocol gate.
- Runtime implementation, repository interfaces, and PostgreSQL adapters remain unauthorized until later bounded work.

## Follow-Up

The next work item should define the storage objects persistence schema gate. It should decide table and index posture before adding migration source or repository code.

## Redaction Notes

No private local environment file was read or printed. No secrets, raw credentials, raw tokens, verifier keys, DSNs, digests, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, object values, or concrete transport metadata were added to the record.
