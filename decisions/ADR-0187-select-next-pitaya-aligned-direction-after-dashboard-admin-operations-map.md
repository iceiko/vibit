# ADR-0187: Select Next Pitaya-Aligned Direction After Dashboard/Admin Operations Map

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map.md`

Related artifacts:

- `decisions/ADR-0186-pitaya-aligned-dashboard-admin-operations-source-first-map.md`
- `decisions/ADR-0185-pitaya-aligned-dashboard-admin-operations-boundary-gate.md`
- `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md`
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

`ADR-0186` completed the source-first Pitaya-aligned dashboard and admin operations map through `node tools/vibit inspect pitaya-dashboard-admin --json`.

The map consolidated dashboard/admin vocabulary, current source-first operations mapping, runtime observability, metrics/tracing, health/readiness/version/config, route inventory, verification posture, redaction posture, audit/event deferrals, admin authorization deferrals, source surfaces, and implementation deferrals. The next bounded continuation should keep Pitaya-class architecture planning moving without turning the dashboard/admin map into runtime behavior, dashboard, admin console, endpoint, protocol, persistence, dependency, or distributed runtime authorization.

The safe next step is a boundary gate for runtime component lifecycle vocabulary. Pitaya-class architecture uses component lifecycle concepts to organize services, handlers, startup, shutdown, and runtime ownership. vibit needs that vocabulary in an agent-native form before adding any component lifecycle behavior.

## Decision

Select `define_pitaya_aligned_runtime_component_lifecycle_boundary_gate` as the next bounded Pitaya-aligned direction after the dashboard/admin operations source-first map.

Register `runtime.next_pitaya_aligned_direction_after_dashboard_admin_operations_map` as the repository check rule.

Complete `M-207/W-0279` and open `M-208/W-0280 Define Pitaya-aligned runtime component lifecycle boundary gate` as next-ready.

This decision does not add runtime component lifecycle behavior. It also does not add handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement runtime component lifecycle behavior immediately.
- Implement handler registration behavior immediately.
- Add startup or shutdown hooks immediately.
- Add component discovery or component loading immediately.
- Add dashboard/admin behavior after the dashboard/admin source-first map.
- Return directly to Nakama product module expansion after the dashboard/admin map.
- Start distributed runtime implementation after the dashboard/admin operations vocabulary sequence.

## Rationale

Runtime component lifecycle is a high-value Pitaya alignment direction because it can clarify how future server components, handlers, bootstrapping, shutdown, dependency ownership, and distributed architecture vocabulary relate to vibit's existing application/runtime boundaries.

It is also broad enough to create accidental implementation scope. A selection-only step preserves the boundary-first sequence: the next work item can define vocabulary and deferrals before any component lifecycle code, handler registration behavior, startup hook, shutdown hook, protocol shape, persistence, dependency, dashboard, admin, hosted, SDK, release, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a selection-only continuation. The correct continuation is to record the selected follow-up direction, update repository memory, register the check rule, and open W-0280. Runtime component lifecycle behavior remains deferred until a later gate and implementation slice explicitly authorize it.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  runtime_component_lifecycle_boundary_clarity: high
  implementation_boundedness: high
  agent_native_runtime_ownership_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

W-0280 becomes the next-ready work item and must define only a Pitaya-aligned runtime component lifecycle boundary gate.

This selection does not authorize runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol changes, generated output, persistence changes, dependencies, hosted surfaces, SDK publication, distributed runtime behavior, or direct compatibility.

## Reversal Conditions

Revisit this decision if a later ADR chooses to return to Nakama product module expansion before any further Pitaya-aligned runtime architecture vocabulary, or if runtime component lifecycle needs to be split into separate gates for component ownership, handler registration, startup, shutdown, or distributed runtime lifecycle.

## Follow-Up

- Complete `W-0280`: define the Pitaya-aligned runtime component lifecycle boundary gate.
- Keep runtime component lifecycle behavior, handler registration behavior, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
