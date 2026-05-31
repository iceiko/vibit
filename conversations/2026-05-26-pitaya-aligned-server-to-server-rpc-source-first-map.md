# Conversation: Pitaya-Aligned Server To Server RPC Source-First Map

Date: 2026-05-31
Work item: `M-179/W-0251 Implement Pitaya-aligned server-to-server RPC source-first map`
Decision: `ADR-0159`
Check rule: `runtime.pitaya_aligned_server_to_server_rpc_source_first_map`

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Context

`M-178/W-0250` defined the Pitaya-aligned server-to-server RPC boundary gate, accepted `ADR-0158`, registered `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`, and opened `M-179/W-0251` as next-ready.

The active work item authorized only a source-first RPC map. It allowed making `server_to_server_rpc`, `remote_call`, `route_handler`, `module_handler`, `application_dispatch`, and `service_discovery` vocabulary inspectable, but it did not authorize server-to-server RPC implementation, remote call behavior, service discovery, frontend/backend role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.

## Maintainer Narrative

The maintainer asked to continue for "20 steps" and move closer to Pitaya. In repository terms, that means advancing the active work queue while preserving vibit's source-first, ADR-first, and checkable-boundary model.

For W-0251, the safe Pitaya-aligned step is an inspection surface that maps current single-process dispatch and module-handler responsibilities to deferred RPC and remote-call vocabulary. The request does not authorize changing runtime topology or adding RPC, remoting, service discovery, protocol carriers, persistence, dependencies, hosted deployment, SDK publication, or direct Pitaya compatibility.

## Agent Response Summary

The agent completed the source-first RPC map:

1. Confirmed W-0251 was the active next-ready work item.
2. Confirmed ADR-0158 authorized only a source-first RPC map.
3. Ran RED checks for the missing rule, missing command, and missing change artifacts.
4. Implemented `node tools/vibit inspect pitaya-rpc --json`.
5. Registered `runtime.pitaya_aligned_server_to_server_rpc_source_first_map`.
6. Added repository check coverage.
7. Accepted `ADR-0159`.
8. Added change artifacts and this conversation memory.
9. Updated manifests, guides, README, alpha, maturity, and roadmap memory.
10. Completed W-0251.
11. Opened W-0252 as the next Pitaya-aligned service discovery boundary gate.

## RED Check

Initial checks failed because the W-0251 rule, command, and change artifacts were absent:

```text
Unknown rule_id: runtime.pitaya_aligned_server_to_server_rpc_source_first_map
Unknown command
change directory does not exist
```

## Decisions

- Implement `node tools/vibit inspect pitaya-rpc --json` as a source-first repository inspection map.
- Accept `ADR-0159`.
- Register `runtime.pitaya_aligned_server_to_server_rpc_source_first_map`.
- Complete `M-179/W-0251`.
- Open `M-180/W-0252 Define Pitaya-aligned service discovery boundary gate`.

The map records:

- ADR-0158 RPC vocabulary;
- current single-process mappings for RPC and remote-call planning;
- source surfaces;
- runtime behavior, RPC implementation, remotes, service discovery, frontend/backend role implementation, distributed runtime, protocol, generated output, persistence, dependency, hosted, SDK, and direct compatibility deferrals;
- redaction flags for credentials, tokens, digests, verifier keys, DSNs, database payloads, and local secret file contents.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0159-pitaya-aligned-server-to-server-rpc-source-first-map.md`
- `conversations/2026-05-26-pitaya-aligned-server-to-server-rpc-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-server-to-server-rpc-source-first-map/`
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

`runtime.pitaya_aligned_server_to_server_rpc_source_first_map` is registered and checked.

`node tools/vibit inspect pitaya-rpc --json` reports the RPC vocabulary, current single-process mapping, source surfaces, deferrals, redaction posture, and W-0252 follow-up.

`M-179/W-0251` is completed.

`M-180/W-0252 Define Pitaya-aligned service discovery boundary gate` is next-ready.

## Boundaries Preserved

No server-to-server RPC implementation, remote call behavior, service discovery, frontend/backend server role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility were added.

## Open Questions

None for W-0251. W-0252 must define the service discovery boundary before any service registry, discovery dependency, node registry, distributed process topology, or discovery-backed routing implementation work starts.

## Follow-Up

Complete `W-0252` by defining a Pitaya-aligned service discovery boundary gate.

## Redaction Notes

No secrets, tokens, credentials, DSNs, database payloads, transport metadata values, or local ignored file contents are recorded here. The inspection output records redaction posture only and does not expose credential material, access-token material, lookup digests, verifier digests, verifier keys, database payloads, or local secret contents.
