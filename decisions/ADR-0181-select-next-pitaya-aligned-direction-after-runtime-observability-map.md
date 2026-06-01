# ADR-0181: Select Next Pitaya-Aligned Direction After Runtime Observability Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-runtime-observability-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-runtime-observability-map.md`

Related artifacts:

- `decisions/ADR-0180-pitaya-aligned-runtime-observability-source-first-map.md`
- `decisions/ADR-0179-pitaya-aligned-runtime-observability-boundary-gate.md`
- `docs/pitaya-aligned-runtime-observability-boundary-gate.md`
- `docs/minimum-operations-inspection-surface-gate.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0180` completed the source-first Pitaya-aligned runtime observability map through `node tools/vibit inspect pitaya-observability --json`.

The map consolidated current local operations inspection, health/readiness/version/config posture, route inventory, verification posture, redaction posture, source surfaces, and deferred operations surfaces. The repeatedly deferred operational surface now has a shared vocabulary, but concrete metrics and tracing semantics remain intentionally undefined.

The safe next step is a boundary gate for metrics and tracing vocabulary. That gate can define how future metrics and trace signals relate to current source-first operations posture without adding runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, dependencies, or direct compatibility.

## Decision

Select `define_pitaya_aligned_metrics_tracing_boundary_gate` as the next bounded Pitaya-aligned direction after the runtime observability source-first map.

Register `runtime.next_pitaya_aligned_direction_after_runtime_observability_map` as the repository check rule.

Complete `M-201/W-0273` and open `M-202/W-0274 Define Pitaya-aligned metrics and tracing boundary gate` as next-ready.

This decision does not add metrics endpoints. It also does not add runtime endpoint behavior, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, session binding behavior, kick/disconnect behavior, session data behavior, session data persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, cluster-safe session routing behavior, distributed session routing, service discovery implementation, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement metrics endpoints immediately.
- Implement tracing pipelines immediately.
- Implement dashboards or admin console behavior immediately.
- Return directly to Nakama product module expansion after the runtime observability map.
- Start distributed runtime implementation after the runtime observability vocabulary sequence.

## Rationale

Runtime observability has now been mapped at the source-first level, but metrics and tracing are still high-risk because they can quickly imply new endpoints, dependencies, live runtime inspection, deployment surfaces, dashboards, and production operations behavior.

A gate keeps the next step bounded. It can define metrics and tracing vocabulary, allowed future source surfaces, redaction expectations, and implementation deferrals before any endpoint, pipeline, dependency, or dashboard work is authorized.

## Agent Reasoning Summary

The active work item is a selection-only continuation. The correct continuation is to record the selected follow-up direction, update repository memory, register the check rule, and open W-0274. Metrics and tracing behavior remain deferred until a later gate and implementation slice explicitly authorize it.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  operations_observability_continuity: high
  implementation_boundedness: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

W-0274 becomes the next-ready work item and must define only a metrics and tracing boundary gate.

This selection does not authorize runtime endpoints, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol changes, generated output, persistence changes, dependencies, hosted surfaces, SDK publication, distributed runtime behavior, or direct compatibility.

## Reversal Conditions

Revisit this decision if a later ADR chooses to return to Nakama product module expansion before any further Pitaya-aligned operations vocabulary, or if metrics and tracing must be split into separate gates before any shared boundary is useful.

## Follow-Up

- Complete `W-0274`: define the Pitaya-aligned metrics and tracing boundary gate.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, dependencies, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
