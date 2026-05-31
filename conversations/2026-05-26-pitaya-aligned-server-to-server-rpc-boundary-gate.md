# Conversation: Pitaya-Aligned Server To Server RPC Boundary Gate

Date: 2026-05-31
Work item: `M-178/W-0250 Define Pitaya-aligned server-to-server RPC boundary gate`
Decision: `ADR-0158`
Check rule: `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Context

`M-177/W-0249` implemented `node tools/vibit inspect pitaya-roles --json`, accepted `ADR-0157`, registered `runtime.pitaya_aligned_frontend_backend_role_source_first_map`, and opened `M-178/W-0250` as next-ready.

The active work item was explicitly gate-only. It allowed defining how `server_to_server_rpc` and `remote_call` vocabulary may be used for future planning, but did not authorize RPC implementation, remote calls, service discovery, distributed runtime behavior, frontend/backend role implementation, protocol changes, generated output, persistence changes, dependencies, hosted deployment, SDK publication, or direct Nakama/Pitaya API compatibility.

## Maintainer Narrative

The maintainer asked to continue for "20 steps" and move closer to Pitaya. In this repository, the safe interpretation is to keep advancing the work queue while turning Pitaya architecture pressure into explicit, checkable, source-first boundaries.

The request does not authorize server-to-server RPC behavior, remote handler invocation, service discovery, distributed process topology, public compatibility with Pitaya APIs, or bypassing existing vibit module contracts.

## Agent Response Summary

The agent completed the server-to-server RPC boundary gate:

1. Confirmed W-0250 was the active next-ready work item.
2. Confirmed W-0250 was gate-only.
3. Ran RED checks for the missing rule and missing change artifacts.
4. Defined the server-to-server RPC boundary standard and Simplified Chinese translation.
5. Accepted `ADR-0158`.
6. Added change artifacts.
7. Registered the repository check rule.
8. Updated manifests and continuation docs.
9. Completed W-0250.
10. Opened W-0251 as the next source-first server-to-server RPC map slice.

## RED Check

Initial commands:

```bash
node tools/vibit inspect rule runtime.pitaya_aligned_server_to_server_rpc_boundary_gate
node tools/vibit check change define-pitaya-aligned-server-to-server-rpc-boundary-gate --json
```

Results before implementation:

```text
Unknown rule_id: runtime.pitaya_aligned_server_to_server_rpc_boundary_gate
change directory does not exist: changes/define-pitaya-aligned-server-to-server-rpc-boundary-gate
```

## Decisions

- Accept `ADR-0158`.
- Register `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`.
- Define RPC vocabulary for `server_to_server_rpc` and `remote_call`.
- Map current single-process application dispatch, protocol bridge, route handler, and module handler responsibilities to deferred RPC vocabulary.
- Complete `M-178/W-0250`.
- Open `M-179/W-0251 Implement Pitaya-aligned server-to-server RPC source-first map`.

## Artifacts

- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.zh-CN.md`
- `decisions/ADR-0158-pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-server-to-server-rpc-boundary-gate/`
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

The gate narrows Pitaya vocabulary to future internal call boundaries:

- `server_to_server_rpc`
- `remote_call`
- `route_handler`
- `module_handler`
- `application_dispatch`
- `service_discovery`

Current implementation remains single-process and modular-monolith.

## Boundary

This slice does not implement server-to-server RPC, remote call behavior, service discovery, frontend/backend server roles, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Open Questions

- Which command shape should W-0251 use for the RPC map: `inspect pitaya-rpc`, `inspect pitaya-remotes`, or extending `inspect pitaya-vocabulary`?

## Follow-Up

Complete `W-0251` by implementing a source-first Pitaya-aligned server-to-server RPC map.

## Redaction Notes

No secrets, tokens, credentials, DSNs, database payloads, transport metadata values, or local ignored file contents are recorded here.
