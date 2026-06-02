# ADR-0189: Pitaya-Aligned Runtime Component Lifecycle Source-First Map

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-runtime-component-lifecycle-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-runtime-component-lifecycle-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`
- `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.zh-CN.md`
- `decisions/ADR-0188-pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`
- `decisions/ADR-0187-select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map.md`
- `decisions/ADR-0186-pitaya-aligned-dashboard-admin-operations-source-first-map.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0188` defined the gate-only Pitaya-aligned runtime component lifecycle vocabulary boundary and opened `M-209/W-0281`.

The repository already exposes source-first Pitaya-aligned maps for dashboard/admin operations, metrics/tracing, runtime observability, session lifecycle, acceptor/connection lifecycle, serializer/message forwarding, route handler pipeline, cluster-safe session routing, distributed groups, service discovery, RPC, frontend/backend roles, and distributed runtime vocabulary. The safe next step is a source-first runtime component lifecycle inspection command that reports vocabulary, current source mappings, redaction posture, source surfaces, and implementation deferrals without adding component lifecycle behavior, dynamic handler registration, component discovery/loading, startup or shutdown hooks, runtime endpoints, dependencies, distributed runtime behavior, or direct compatibility.

## Decision

Implement `node tools/vibit inspect pitaya-component-lifecycle --json` as the source-first Pitaya-aligned runtime component lifecycle map for `M-209/W-0281`.

The command reports:

- ADR-0188 as the source gate and ADR-0189 as the implementation decision.
- `runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map` as the check rule.
- Allowed runtime component lifecycle vocabulary from ADR-0188.
- Related vocabulary from dashboard/admin operations, route handler pipeline, payload registry, runtime observability, source-first operations, verification, redaction, and node-local runtime surfaces.
- Current mappings for bootstrap composition, handler registration source files, application services, repository wiring, protocol adapter composition, lifecycle state, operations cross-reference, and distributed component lifecycle deferral.
- Explicit false deferrals for runtime component lifecycle behavior, handler registration behavior, component discovery/loading, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility.
- `M-210/W-0282 Select next Pitaya-aligned direction after runtime component lifecycle map` as the next-ready follow-up.

This decision does not add runtime component lifecycle behavior. It also does not add handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add runtime component lifecycle behavior immediately.
- Add dynamic handler registration immediately.
- Add component discovery or loading immediately.
- Add startup or shutdown hooks immediately.
- Expand the dashboard/admin source-first map instead of adding a component-lifecycle-focused map.
- Return directly to broad product module expansion after W-0280.

## Rationale

Runtime component lifecycle vocabulary is useful only if agents can see how the current source tree composes runtime pieces without mistaking that visibility for permission to add lifecycle behavior. A focused command gives agents a bounded map before lifecycle interfaces, dynamic registration, discovery/loading, hooks, dependency containers, endpoints, distributed lifecycle behavior, or compatibility work is authorized.

## Agent Reasoning Summary

The active work item is a source-first inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all runtime component lifecycle, handler registration, component discovery/loading, startup/shutdown hook, runtime endpoint, dashboard/admin, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  runtime_component_lifecycle_mapping_clarity: high
  implementation_boundedness: high
  source_first_runtime_composition_reuse: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned runtime component lifecycle vocabulary without reading every architecture document.

This decision does not add runtime component lifecycle behavior, handler registration behavior, component discovery/loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete component lifecycle, handler registration, component discovery/loading, startup/shutdown, dependency container, runtime endpoint, or distributed lifecycle implementation model;
- the component lifecycle inspection output creates confusion with public API compatibility or live lifecycle behavior;
- current bootstrap composition, route registration, payload registry, application services, repository wiring, protocol adapter composition, process startup, or transport close posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate source-first maps for runtime components, handler registration, component discovery, lifecycle hooks, dependency containers, or distributed lifecycle coordination.

## Follow-Up

- Complete `W-0282`: select the next Pitaya-aligned direction after the runtime component lifecycle source-first map.
- Keep runtime component lifecycle behavior, handler registration behavior, component discovery/loading, startup hooks, shutdown hooks, runtime endpoints, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
