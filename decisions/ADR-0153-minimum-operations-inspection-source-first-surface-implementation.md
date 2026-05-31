# ADR-0153: Minimum Operations Inspection Source-First Surface Implementation

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-minimum-operations-inspection-source-first-surface/`

Related conversations:

- `conversations/2026-05-26-minimum-operations-inspection-source-first-surface-implementation.md`

Related artifacts:

- `tools/vibit`
- `rules/check-rules.json`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`

## Context

`ADR-0152` accepted the minimum operations inspection surface gate and selected `source_first_local_operations_inspection` as the first operations posture. The gate authorized a repository-owned inspection surface, preferably `tools/vibit inspect operations`, and explicitly rejected adding runtime operations/admin endpoints, metrics endpoints, observability pipelines, dashboards, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, distributed runtime implementation, or direct Nakama/Pitaya API compatibility.

The maintainer asked to continue multiple steps and move closer to Pitaya. The safe interpretation is to complete the current source-first operations inspection slice and make Pitaya architecture vocabulary more visible as deferred, checkable planning context. Pitaya remains a future architecture reference, not an implementation authorization.

## Decision

Implement `node tools/vibit inspect operations --json` as the minimum source-first operations inspection surface.

The command emits a JSON `operations_inspection` record that summarizes:

- existing local troubleshooting endpoints: `/healthz`, `/readyz`, `/version`, and `/configz`;
- source surfaces such as `.arch/work-items.yaml`, `.arch/runtime.yaml`, `.arch/reference.yaml`, runbooks, alpha docs, examples, `tools/vibit`, and rule catalog metadata;
- the ten minimum inspectable categories from `ADR-0152`;
- current local alpha route families for authentication, inventory, presence, storage, and friends;
- local alpha flow steps;
- persistence and migration posture without database contents;
- generated-output and Protobuf posture without hand-editing generated files;
- verification posture and known warning class;
- redaction flags for raw credentials, raw tokens, digests, verifier keys, DSNs, transport metadata, sensitive identifiers, database payloads, and local secret file contents;
- Pitaya deferred architecture mapping for acceptors, session binding, route handler model, frontend/backend roles, RPC/remotes, groups/broadcast, and service discovery.

Register `runtime.minimum_operations_inspection_source_first_surface_implementation` as the check rule for this slice.

Complete `M-173/W-0245` and open `M-174/W-0246 Define Pitaya-aligned distributed runtime vocabulary reactivation gate` as next-ready.

This decision does not add operations/admin endpoints, metrics endpoints, observability pipelines, dashboards, runtime endpoint behavior, gameplay protocol routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, event/audit tables, hosted deployments, SDKs, distributed runtime implementation, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add a runtime HTTP operations endpoint immediately.
- Add metrics output or Prometheus-style endpoint names immediately.
- Add a dashboard or local admin console.
- Add live player/session/token/friends/storage inspectors.
- Jump directly to Pitaya frontend/backend, RPC, service discovery, or distributed group implementation.
- Keep Pitaya hidden until a much later distributed-runtime milestone.

## Rationale

The source-first command gives prototype authors and agents a compact way to inspect the current local alpha posture without widening runtime behavior. It also creates a concrete bridge toward Pitaya by making future distributed architecture vocabulary visible and redacted now, while keeping implementation deferred.

This fits vibit's agent-native model: state is inspectable through committed manifests and tools, requirements are bounded by ADRs and work items, and future architecture pressure is recorded before code expands.

## Agent Reasoning Summary

The active work item is implementation, not a new gate. The correct implementation is a repository inspection command plus checks and documentation. To honor the maintainer's Pitaya direction without violating the current boundary, the output records Pitaya-style acceptor/session/handler and deferred frontend/backend/RPC/group/service-discovery vocabulary, then opens a dedicated gate for that vocabulary.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  operations_visibility_value: high
  pitaya_alignment_value: medium
  redaction_safety: high
  runtime_behavior_risk: none
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `node tools/vibit inspect operations --json` is available.
- `runtime.minimum_operations_inspection_source_first_surface_implementation` becomes the check rule for W-0245.
- Operations inspection remains source-first and local.
- Pitaya architecture vocabulary is visible in a deferred map.
- `M-173/W-0245` is completed.
- `M-174/W-0246 Define Pitaya-aligned distributed runtime vocabulary reactivation gate` becomes next-ready.
- Runtime operations/admin endpoints, metrics, observability, dashboards, protocol changes, generated output, persistence, dependencies, distributed runtime implementation, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- prototype feedback shows that source-first operations inspection is insufficient;
- a later operations gate authorizes a runtime endpoint or admin console;
- redaction policy changes and allows a currently forbidden identifier class;
- Pitaya distributed architecture becomes active through a later ADR;
- the command output becomes stale relative to committed runtime route families.

## Follow-Up

- Complete `W-0246`: define a Pitaya-aligned distributed runtime vocabulary reactivation gate.
- Keep cluster/RPC/service-discovery implementation, distributed groups, cluster-safe session routing, runtime endpoint expansion, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
