# Conversation: Pitaya-Aligned Frontend Backend Role Boundary Gate

Date: 2026-05-31
Work item: `M-176/W-0248 Define Pitaya-aligned frontend/backend role boundary gate`
Decision: `ADR-0156`
Check rule: `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Context

`M-175/W-0247` had implemented `node tools/vibit inspect pitaya-vocabulary --json`, accepted `ADR-0155`, registered `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map`, and opened `M-176/W-0248` as next-ready.

The active work item was explicitly gate-only. It allowed defining how `frontend_server` and `backend_server` vocabulary may be used for future planning, but did not authorize frontend/backend role implementation or distributed runtime behavior.

## Maintainer Narrative

The maintainer asked to continue for "20 steps" and move closer to Pitaya. In this repository, the safe interpretation is to keep advancing the work queue while turning Pitaya architecture pressure into explicit, checkable, source-first boundaries.

The request does not authorize frontend/backend process topology, server-to-server RPC, service discovery, distributed groups, cluster-safe session routing, direct Pitaya API compatibility, or bypassing existing vibit module contracts.

## Agent Response Summary

The agent completed the role boundary gate:

1. Confirmed W-0248 was the active next-ready work item.
2. Confirmed W-0248 was gate-only.
3. Ran RED checks for the missing rule and missing change artifacts.
4. Defined the frontend/backend role boundary standard and Simplified Chinese translation.
5. Accepted `ADR-0156`.
6. Added change artifacts.
7. Registered the repository check rule.
8. Updated manifests and continuation docs.
9. Completed W-0248.
10. Opened W-0249 as the next source-first frontend/backend role map slice.

## RED Check

Initial commands:

```bash
node tools/vibit inspect rule runtime.pitaya_aligned_frontend_backend_role_boundary_gate
node tools/vibit check change define-pitaya-aligned-frontend-backend-role-boundary-gate --json
```

Results before implementation:

```text
Unknown rule_id: runtime.pitaya_aligned_frontend_backend_role_boundary_gate
change directory does not exist: changes/define-pitaya-aligned-frontend-backend-role-boundary-gate
```

## Decisions

- Accept `ADR-0156`.
- Register `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`.
- Define role vocabulary for `frontend_server` and `backend_server`.
- Map current single-process acceptor, session binding, protocol bridge, application dispatch, and module handler responsibilities to deferred role vocabulary.
- Complete `M-176/W-0248`.
- Open `M-177/W-0249 Implement Pitaya-aligned frontend/backend role source-first map`.

## Artifacts

- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.md`
- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.zh-CN.md`
- `decisions/ADR-0156-pitaya-aligned-frontend-backend-role-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-frontend-backend-role-boundary-gate/`
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

The gate narrows Pitaya vocabulary to role boundaries:

- `frontend_server`
- `backend_server`
- `acceptor`
- `session_binding`
- `route_handler`

Current implementation remains single-process and modular-monolith.

## Boundary

This slice does not implement frontend/backend server roles, distributed runtime behavior, server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Open Questions

- Which command shape should W-0249 use for the role map: `inspect pitaya-roles`, `inspect frontend-backend-roles`, or extending `inspect pitaya-vocabulary`?

## Follow-Up

Complete `W-0249` by implementing a source-first Pitaya-aligned frontend/backend role map.

## Redaction Notes

No secrets, tokens, credentials, DSNs, database payloads, transport metadata values, or local ignored file contents are recorded here.
