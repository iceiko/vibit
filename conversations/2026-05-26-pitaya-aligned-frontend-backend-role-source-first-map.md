# Conversation: Pitaya-Aligned Frontend Backend Role Source-First Map

Date: 2026-05-31
Work item: `M-177/W-0249 Implement Pitaya-aligned frontend/backend role source-first map`
Decision: `ADR-0157`
Check rule: `runtime.pitaya_aligned_frontend_backend_role_source_first_map`

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Context

`M-176/W-0248` defined the Pitaya-aligned frontend/backend role boundary gate, accepted `ADR-0156`, registered `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`, and opened `M-177/W-0249` as next-ready.

The active work item authorized only a source-first role map. It allowed making `frontend_server`, `backend_server`, `acceptor`, `session_binding`, and `route_handler` vocabulary inspectable, but it did not authorize frontend/backend role implementation, distributed runtime behavior, server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.

## Maintainer Narrative

The maintainer asked to continue for "20 steps" and move closer to Pitaya. In repository terms, that means advancing the active work queue while preserving vibit's source-first, ADR-first, and checkable-boundary model.

For W-0249, the safe Pitaya-aligned step is an inspection surface that maps current single-process responsibilities to deferred frontend/backend role vocabulary. The request does not authorize changing runtime topology or adding RPC, remoting, service discovery, protocol carriers, persistence, dependencies, hosted deployment, SDK publication, or direct Pitaya compatibility.

## Agent Response Summary

The agent completed the source-first role map:

1. Confirmed W-0249 was the active next-ready work item.
2. Confirmed ADR-0156 authorized only a source-first role map.
3. Ran RED checks for the missing rule, missing command, and missing change artifacts.
4. Implemented `node tools/vibit inspect pitaya-roles --json`.
5. Registered `runtime.pitaya_aligned_frontend_backend_role_source_first_map`.
6. Added repository check coverage.
7. Accepted `ADR-0157`.
8. Added change artifacts and this conversation memory.
9. Updated manifests, guides, README, alpha, maturity, and roadmap memory.
10. Completed W-0249.
11. Opened W-0250 as the next Pitaya-aligned server-to-server RPC boundary gate.

## RED Check

Initial checks failed because the W-0249 rule, command, and change artifacts were absent:

```text
Unknown rule_id: runtime.pitaya_aligned_frontend_backend_role_source_first_map
Unknown command: node tools/vibit inspect pitaya-roles --json
change directory does not exist: changes/implement-pitaya-aligned-frontend-backend-role-source-first-map
```

## Decisions

- Implement `node tools/vibit inspect pitaya-roles --json` as a source-first repository inspection map.
- Accept `ADR-0157`.
- Register `runtime.pitaya_aligned_frontend_backend_role_source_first_map`.
- Complete `M-177/W-0249`.
- Open `M-178/W-0250 Define Pitaya-aligned server-to-server RPC boundary gate`.

The map records:

- ADR-0156 role vocabulary;
- current single-process mappings for frontend/backend role planning;
- source surfaces;
- runtime behavior, role implementation, distributed runtime, RPC, remotes, service discovery, protocol, generated output, persistence, dependency, hosted, SDK, and direct compatibility deferrals;
- redaction flags for credentials, tokens, digests, verifier keys, DSNs, database payloads, and local secret file contents.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0157-pitaya-aligned-frontend-backend-role-source-first-map.md`
- `conversations/2026-05-26-pitaya-aligned-frontend-backend-role-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-frontend-backend-role-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`

## Outcome

`runtime.pitaya_aligned_frontend_backend_role_source_first_map` is registered and checked.

`node tools/vibit inspect pitaya-roles --json` reports the role vocabulary, current single-process mapping, source surfaces, deferrals, redaction posture, and W-0250 follow-up.

`M-177/W-0249` is completed.

`M-178/W-0250 Define Pitaya-aligned server-to-server RPC boundary gate` is next-ready.

## Boundaries Preserved

No frontend/backend server role implementation, distributed runtime behavior, server-to-server RPC, remote call behavior, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility were added.

## Open Questions

None for W-0249. W-0250 must define the server-to-server RPC boundary before any RPC, remote call, or service discovery implementation work starts.

## Follow-Up

Complete `W-0250` by defining a Pitaya-aligned server-to-server RPC boundary gate.

## Redaction Notes

No secrets, tokens, credentials, DSNs, database payloads, transport metadata values, or local ignored file contents are recorded here. The inspection output records redaction posture only and does not expose credential material, access-token material, lookup digests, verifier digests, verifier keys, database payloads, or local secret contents.
