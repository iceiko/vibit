# Conversation: Minimum Operations Inspection Source-First Surface Implementation

Date: 2026-05-31
Work item: `M-173/W-0245 Implement minimum operations inspection source-first surface`
Decision: `ADR-0153`
Check rule: `runtime.minimum_operations_inspection_source_first_surface_implementation`

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Context

The repository had completed `M-172/W-0244`, accepted `ADR-0152`, registered `runtime.minimum_operations_inspection_surface_gate`, and opened `M-173/W-0245` as next-ready.

The active gate allowed only source-first repository inspection and documentation updates. It explicitly deferred operations/admin endpoints, metrics endpoints, observability pipelines, dashboards, runtime behavior, protocol routes, Protobuf source, generated output, persistence, dependencies, hosted surfaces, SDKs, distributed runtime implementation, and direct Nakama/Pitaya API compatibility.

Pitaya remained a deferred future architecture reference for acceptors, session binding, route handler model, frontend/backend roles, RPC/remotes, groups/broadcast, and service discovery.

## Maintainer Narrative

The maintainer asked to continue for "20 steps" and move closer to Pitaya. In repository terms, that means advancing the active work queue while preserving the gate discipline already recorded in `.arch/work-items.yaml`.

The request does not authorize distributed runtime implementation, direct Pitaya API compatibility, or bypassing ask-first boundaries. It does authorize recording Pitaya-aligned vocabulary and opening the next bounded gate when the active source-first operations inspection slice is complete.

## Agent Response Summary

The agent interpreted "20 steps" as a bounded sequence inside the active work queue rather than permission to bypass gates:

1. Confirm current next-ready work.
2. Confirm Pitaya reference boundaries.
3. Run a RED check for `node tools/vibit inspect operations --json`.
4. Implement the source-first inspection command.
5. Add route-family inspection output.
6. Add local alpha flow inspection output.
7. Add persistence and migration posture output.
8. Add generated-output posture output.
9. Add verification posture output.
10. Add redaction flags.
11. Add Pitaya deferred architecture mapping.
12. Register the runtime check rule.
13. Add repository check coverage.
14. Add `ADR-0153`.
15. Add change artifacts.
16. Update runbook and checklist references.
17. Mark `W-0245` completed.
18. Open `W-0246` as a Pitaya-aligned gate.
19. Run verification.
20. Commit and push if verification passes.

## RED Check

Initial command:

```bash
node tools/vibit inspect operations --json
```

Result before implementation:

```text
Unknown command.
```

## Decisions

- Implement `node tools/vibit inspect operations --json`.
- Accept `ADR-0153 Minimum Operations Inspection Source-First Surface Implementation`.
- Register `runtime.minimum_operations_inspection_source_first_surface_implementation`.
- Complete `M-173/W-0245`.
- Open `M-174/W-0246 Define Pitaya-aligned distributed runtime vocabulary reactivation gate`.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0153-minimum-operations-inspection-source-first-surface-implementation.md`
- `changes/2026-05-26-implement-minimum-operations-inspection-source-first-surface/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`

## Pitaya Alignment

The implementation records Pitaya-aligned vocabulary as source-first planning context only:

- WebSocket/TCP acceptors map to the current single-process WebSocket acceptor posture.
- Session binding maps to current first-message connection binding.
- Route handler model maps to current application dispatch plus protocol bridge.
- Frontend/backend server roles remain deferred.
- RPC/remotes remain deferred.
- Groups, rooms, and broadcast remain deferred.
- Cluster service discovery remains deferred.

## Boundary

This slice does not implement operations/admin endpoints, metrics endpoints, observability pipelines, dashboards, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, event/audit tables, SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Open Questions

None for W-0245. W-0246 must define how far Pitaya vocabulary can be reactivated before any implementation work starts.

## Follow-Up

Complete `W-0246` by defining a Pitaya-aligned distributed runtime vocabulary reactivation gate.

## Redaction Notes

The inspection output must not expose raw device credentials, raw access tokens, lookup digests, verifier digests, verifier keys, concrete verifier key set ids, DSNs with credentials, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, sensitive identifiers, database payloads, or local secret file contents.
