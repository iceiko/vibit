# Conversation: Pitaya-Aligned Distributed Runtime Vocabulary Source-First Map

Date: 2026-05-31
Work item: `M-175/W-0247 Implement Pitaya-aligned distributed runtime vocabulary source-first map`

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Context

`W-0246` completed the Pitaya-aligned distributed runtime vocabulary reactivation gate and opened `W-0247` as the next-ready implementation slice. The gate authorized only a source-first vocabulary map and explicitly forbade distributed runtime behavior, frontend/backend role implementation, server-to-server RPC, remote calls, service discovery, distributed groups, broadcast fanout, cluster-safe session routing, protocol shape changes, generated output, persistence changes, dependencies, hosted surfaces, SDKs, and direct compatibility.

## Maintainer Narrative

The maintainer asked to continue for "20 steps" and move closer to Pitaya. In repository terms, that means advancing the active work queue while preserving the source-first and gate-only boundaries already recorded in `.arch/work-items.yaml`.

The request authorizes making Pitaya-aligned distributed runtime vocabulary inspectable for future planning. It does not authorize implementing frontend/backend server roles, distributed runtime behavior, direct Pitaya API compatibility, or bypassing vibit module contracts.

## Agent Response Summary

The agent completed the source-first vocabulary map:

1. Confirmed W-0247 was the active next-ready work item.
2. Confirmed ADR-0154 authorized only a source-first map.
3. Ran RED checks for the missing rule, missing command, and missing change artifacts.
4. Implemented `node tools/vibit inspect pitaya-vocabulary --json`.
5. Registered `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map`.
6. Added repository check coverage.
7. Accepted `ADR-0155`.
8. Added change artifacts and this conversation memory.
9. Updated manifests, runbooks, guides, and product direction docs.
10. Completed W-0247.
11. Opened W-0248 as the next Pitaya-aligned frontend/backend role boundary gate.

## TDD RED

Initial checks failed because the W-0247 rule, command, and change artifacts were absent:

```text
Unknown rule_id: runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map
Unknown command
change directory does not exist: changes/implement-pitaya-aligned-distributed-runtime-vocabulary-source-first-map
```

## Decisions

- Implement `node tools/vibit inspect pitaya-vocabulary --json` as a source-first repository inspection map.
- Accept `ADR-0155`.
- Register `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map`.
- Complete `M-175/W-0247`.
- Open `M-176/W-0248 Define Pitaya-aligned frontend/backend role boundary gate`.

The map records:

- ADR-0154 allowed vocabulary;
- current single-process mapping;
- source surfaces;
- distributed runtime, protocol, persistence, dependency, hosted, SDK, and direct compatibility deferrals;
- redaction flags for credentials, tokens, digests, verifier keys, DSNs, database payloads, and local secret file contents.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0155-pitaya-aligned-distributed-runtime-vocabulary-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-distributed-runtime-vocabulary-source-first-map/`
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

`runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map` is registered and checked.

`M-175/W-0247` is completed.

`M-176/W-0248 Define Pitaya-aligned frontend/backend role boundary gate` is next-ready.

## Boundaries Preserved

No distributed runtime implementation, frontend/backend role implementation, server-to-server RPC, remote call behavior, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility were added.

## Open Questions

None for W-0247. W-0248 must define how `frontend_server` and `backend_server` vocabulary can be used for future planning before any implementation work starts.

## Follow-Up

Complete `W-0248` by defining a Pitaya-aligned frontend/backend role boundary gate.

## Redaction Notes

No secrets, tokens, credentials, DSNs, database payloads, transport metadata values, or local ignored file contents are recorded here. The inspection output records redaction posture only and does not expose credential material, access-token material, lookup digests, verifier digests, verifier keys, database payloads, or local secret contents.
