# Conversation: Storage Objects Protocol Route Gate

Date: 2026-05-22
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-protocol-route-gate/`

Related artifacts:

- `docs/storage-objects-protocol-route-gate.md`
- `docs/storage-objects-protocol-route-gate.zh-CN.md`
- `decisions/ADR-0118-storage-objects-protocol-route-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Context

The repository had completed `W-0209`, which implemented application-owned storage object runtime behavior under `runtime/internal/app/storage`. The next-ready item was `W-0210`, a gate-only protocol route planning step.

Relevant standards read:

- `docs/reference-game-server-alignment.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/storage-objects-runtime-behavior-gate.md`
- `decisions/ADR-0117-storage-objects-runtime-behavior-implementation.md`
- `docs/runtime-protocol-adapter.md`
- `docs/game-protocol.md`

Reference check:

- Nakama motivates storage objects as a client-facing game backend capability with collection/key/value/version and read/list/write/delete behavior.
- Pitaya motivates keeping acceptors, sessions, serializers, route handlers, and backend behavior separated.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。注意要贴合nakama pitaya
```

The agent continued one bounded next-ready work item and treated Nakama/Pitaya alignment as a required reference constraint.

## Agent Response Summary

The agent advanced `W-0210` by defining the storage objects protocol route gate after the application-owned runtime behavior implementation.

The gate records future own-player route ids, candidate message-shape posture, protected-route policy, authenticated request identity handoff, protocol adapter ownership, generated-output expectations, redaction rules, required future tests, and stop conditions.

The work added:

- the English and Simplified Chinese gate standard;
- `ADR-0118`;
- a change spec package for `define-storage-objects-protocol-route-gate`;
- `runtime.storage_objects_protocol_route_gate` catalog and runtime check coverage;
- storage module manifest and guide updates;
- architecture, work item, reference, public next-work, and agent-guide continuation updates.

## Decisions

- Accept `docs/storage-objects-protocol-route-gate.md` as the gate standard.
- Record `ADR-0118`.
- Register `runtime.storage_objects_protocol_route_gate`.
- Complete `M-138/W-0210`.
- Open `M-139/W-0211 Implement storage objects protocol route` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable player-owned storage objects should be exposed as a client-facing game backend capability with collection/key/value/version and read/list/write/delete behavior.

Pitaya guided the layering pressure: acceptors, sessions, serializers, route handlers, and backend behavior should remain separated.

vibit adapted those lessons into its own model: explicit `kind/module/name` route identity, protected-route authenticated wrapper requirements, validated request identity handoff, WebSocket transport neutrality, Protobuf bridge ownership, and application-owned handlers calling application-owned services.

## Artifacts

- `docs/storage-objects-protocol-route-gate.md`
- `docs/storage-objects-protocol-route-gate.zh-CN.md`
- `decisions/ADR-0118-storage-objects-protocol-route-gate.md`
- `changes/2026-05-22-define-storage-objects-protocol-route-gate/`
- `conversations/2026-05-22-storage-objects-protocol-route-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Open Questions

- Protocol route implementation remains deferred to `W-0211`.
- Protobuf source and generated output remain deferred to `W-0211`.
- Startup wiring remains deferred to `W-0211`.
- Public ACLs, admin search, group/guild/party/room/match scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, large object/blob storage, S3-compatible storage, and direct compatibility remain deferred.

## Scope Preserved

This work did not add protocol route implementation, Protobuf source, generated output, startup wiring, runtime handlers, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## Follow-Up

The next bounded work item is expected to be storage objects protocol route implementation, with Protobuf source/generated output and route registration only after this gate.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw storage object values from a real user are recorded in this conversation log.
