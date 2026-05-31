# Conversation: Pitaya-Aligned Distributed Runtime Vocabulary Reactivation Gate

Date: 2026-05-31
Work item: `M-174/W-0246 Define Pitaya-aligned distributed runtime vocabulary reactivation gate`
Decision: `ADR-0154`
Check rule: `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate`

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Context

`M-173/W-0245` had implemented `node tools/vibit inspect operations --json`, accepted `ADR-0153`, registered `runtime.minimum_operations_inspection_source_first_surface_implementation`, recorded a Pitaya deferred architecture map, and opened `M-174/W-0246` as next-ready.

The active work item was explicitly gate-only. It allowed Pitaya-aligned distributed runtime vocabulary reactivation, but did not authorize distributed runtime implementation.

## Maintainer Narrative

The maintainer asked to keep moving and to move closer to Pitaya. In this repository, "closer to Pitaya" currently means making the future distributed runtime vocabulary explicit and checkable while preserving the single-process implementation boundary.

The request does not authorize frontend/backend server roles, server-to-server RPC, service discovery, distributed groups, cluster-safe session routing, direct Pitaya API compatibility, or bypassing existing vibit module contracts.

## Agent Response Summary

The agent completed the vocabulary gate:

1. Confirmed the active next-ready work item was W-0246.
2. Read the reference alignment, roadmap, and previous operations inspection decision.
3. Ran a RED check for `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate`.
4. Defined the vocabulary standard and Simplified Chinese translation.
5. Accepted `ADR-0154`.
6. Added change artifacts.
7. Registered the repository check rule.
8. Updated manifests and continuation docs.
9. Completed W-0246.
10. Opened W-0247 as the next source-first vocabulary map slice.

## RED Check

Initial command:

```bash
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate
```

Result before implementation:

```text
Unknown rule_id: runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate
```

## Decisions

- Accept `ADR-0154`.
- Register `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate`.
- Define allowed vocabulary for acceptors, frontend/backend roles, route handlers, session binding, server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast, and cluster-safe session routing.
- Complete `M-174/W-0246`.
- Open `M-175/W-0247 Implement Pitaya-aligned distributed runtime vocabulary source-first map`.

## Artifacts

- `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
- `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.zh-CN.md`
- `decisions/ADR-0154-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate/`
- `tools/vibit`
- `rules/check-rules.json`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`

## Pitaya Alignment

The gate reactivates vocabulary only:

- `acceptor`
- `frontend_server`
- `backend_server`
- `route_handler`
- `session_binding`
- `server_to_server_rpc`
- `remote_call`
- `service_discovery`
- `distributed_group`
- `room_broadcast`
- `cluster_safe_session_routing`

Current implementation remains single-process and modular-monolith.

## Boundary

This slice does not implement distributed runtime behavior, frontend/backend server roles, server-to-server RPC, remote calls, service discovery, distributed groups, broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Open Questions

- Which source-first command shape should W-0247 use: `inspect pitaya-vocabulary`, `inspect distributed-vocabulary`, or a broader reference inspection surface?
- Which later product capability should follow after the vocabulary map is available: groups/parties, chat/broadcast vocabulary, matchmaking context, or operations hardening?

## Follow-Up

Complete `W-0247` by implementing a source-first Pitaya-aligned distributed runtime vocabulary map.

## Redaction Notes

No secrets, tokens, credentials, DSNs, database payloads, transport metadata values, or local ignored file contents are recorded here.
