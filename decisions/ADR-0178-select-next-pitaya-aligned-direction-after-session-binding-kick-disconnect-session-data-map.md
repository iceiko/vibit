# ADR-0178: Select Next Pitaya-Aligned Direction After Session Binding Kick Disconnect Session Data Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map.md`

Related artifacts:

- `decisions/ADR-0177-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map.md`
- `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md`
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

`ADR-0177` completed the source-first Pitaya-aligned session binding, kick/disconnect, and session data map through `node tools/vibit inspect pitaya-session-lifecycle --json`.

The remaining deferred Pitaya architecture pressure repeatedly recorded across the recent maps is runtime observability: runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, and admin/operations inspection. The repository already has a minimum source-first local operations inspection surface from `ADR-0152` and `ADR-0153`, so the safe next step is not implementation. The safe next step is to define a Pitaya-aligned runtime observability boundary gate that maps current local operations facts to future observability vocabulary.

## Decision

Select `define_pitaya_aligned_runtime_observability_boundary_gate` as the next bounded Pitaya-aligned direction after the session lifecycle source-first map.

Register `runtime.next_pitaya_aligned_direction_after_session_binding_kick_disconnect_session_data_map` as the repository check rule.

Complete `M-198/W-0270` and open `M-199/W-0271 Define Pitaya-aligned runtime observability boundary gate` as next-ready.

This decision does not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, session binding behavior, kick/disconnect behavior, session data behavior, session data persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, cluster-safe session routing behavior, distributed session routing, service discovery implementation, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, or admin console behavior immediately.
- Return directly to Nakama product module expansion after the session lifecycle map.
- Implement session binding, kick/disconnect, or session data behavior now that the source-first map exists.
- Start distributed runtime implementation after the Pitaya vocabulary sequence.

## Rationale

The Pitaya-aligned architecture pass has already covered frontend/backend roles, server-to-server RPC, service discovery, distributed groups and broadcast, cluster-safe session routing, route handler pipelines, serializer/message forwarding, acceptor/connection lifecycle, and session lifecycle vocabulary. Runtime observability is the next shared architecture surface that has been repeatedly deferred and is useful before any future production or distributed runtime work.

A gate keeps the next step bounded. It can define vocabulary and map current source-first operations surfaces without creating metrics, tracing, dashboard, runtime endpoint, dependency, or direct compatibility behavior.

## Agent Reasoning Summary

The active work item is a selection-only continuation. The correct continuation is to record the selected follow-up direction, update repository memory, register the check rule, and open W-0271. Runtime observability behavior remains deferred until a later gate and implementation slice explicitly authorize it.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  remaining_architecture_surface_value: high
  implementation_boundedness: high
  source_first_operations_reuse: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

W-0271 becomes the next-ready work item and must define only a runtime observability boundary gate.

This selection does not authorize runtime endpoints, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol changes, generated output, persistence changes, dependencies, hosted surfaces, SDK publication, distributed runtime behavior, or direct compatibility.

## Reversal Conditions

Revisit this decision if a later ADR chooses to return to Nakama product module expansion before any further Pitaya-aligned vocabulary, or if runtime observability is split into separate metrics, tracing, dashboard, and admin-operation gates.

## Follow-Up

- Complete `W-0271`: define the Pitaya-aligned runtime observability boundary gate.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, dependencies, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
