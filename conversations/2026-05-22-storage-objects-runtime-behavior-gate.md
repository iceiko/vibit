# Conversation: Storage Objects Runtime Behavior Gate

Date: 2026-05-22
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-runtime-behavior-gate/`

Related artifacts:

- `docs/storage-objects-runtime-behavior-gate.md`
- `docs/storage-objects-runtime-behavior-gate.zh-CN.md`
- `decisions/ADR-0116-storage-objects-runtime-behavior-gate.md`
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

`M-135/W-0207` completed the storage objects PostgreSQL adapter implementation under `runtime/internal/platform/persistence/postgres`. It preserved runtime behavior, runtime handlers, protocol routes, Protobuf sources, generated output, dependency changes, migration changes, authentication/session behavior changes, large object/blob storage, S3-compatible object storage, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0208 Define storage objects runtime behavior gate`.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。
```

The agent continued one bounded next-ready work item.

## Agent Response Summary

The agent advanced `W-0208` by defining the storage objects runtime behavior gate after the PostgreSQL adapter.

The gate records future application-owned runtime behavior, validated request identity requirements, owner derivation, metadata-only identity refusal, route-policy posture, permissions, validation, conflict mapping, unit-of-work repository handoff, test expectations, and stop conditions.

The work added:

- the English and Simplified Chinese gate standard;
- `ADR-0116`;
- a change spec package for `define-storage-objects-runtime-behavior-gate`;
- `runtime.storage_objects_runtime_behavior_gate` catalog and runtime check coverage;
- storage module manifest and guide updates;
- architecture, work item, reference, public next-work, and agent-guide continuation updates.

## Decisions

- Accept `docs/storage-objects-runtime-behavior-gate.md` as the gate standard.
- Record `ADR-0116`.
- Register `runtime.storage_objects_runtime_behavior_gate`.
- Complete `M-136/W-0208`.
- Open `M-137/W-0209 Implement storage objects runtime behavior` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable storage-object-like game state should become usable through runtime behavior after repository and adapter pieces exist.

Pitaya guided the layering pressure: handlers should receive normalized context and should not parse transport credentials or own persistence mechanics.

vibit adapted those lessons into its own model: an application-owned runtime behavior gate that requires validated request identity and keeps protocol/generated surfaces separate.

## Artifacts

- `docs/storage-objects-runtime-behavior-gate.md`
- `docs/storage-objects-runtime-behavior-gate.zh-CN.md`
- `decisions/ADR-0116-storage-objects-runtime-behavior-gate.md`
- `changes/2026-05-22-define-storage-objects-runtime-behavior-gate/`
- `conversations/2026-05-22-storage-objects-runtime-behavior-gate.md`
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

- Runtime storage object behavior implementation remains deferred to `W-0209`.
- Protocol routes and Protobuf messages remain deferred.
- Startup wiring remains deferred.
- Public command/query contracts and generated output remain deferred.
- Public ACLs, admin search, group/guild/party/room/match scopes, large object/blob storage, S3-compatible storage, and direct compatibility remain deferred.

## Follow-Up

- Implement storage objects runtime behavior in `W-0209`.
- Preserve protocol routes, generated output, dependency changes, migration changes, authentication/session behavior, hosted deployment, release artifact, announcement, promotion, blob/S3 storage, and direct compatibility deferrals unless a later bounded work item explicitly authorizes them.
- Keep metadata-only `player_id` and `session_id` insufficient as authenticated proof.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw storage object values from a real user are recorded in this conversation log.
