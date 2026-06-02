# ADR-0190: Select Next Pitaya-Aligned Direction After Runtime Component Lifecycle Map

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-runtime-component-lifecycle-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-runtime-component-lifecycle-map.md`

Related artifacts:

- `decisions/ADR-0189-pitaya-aligned-runtime-component-lifecycle-source-first-map.md`
- `decisions/ADR-0188-pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`
- `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`
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

`ADR-0189` completed the source-first Pitaya-aligned runtime component lifecycle map through `node tools/vibit inspect pitaya-component-lifecycle --json`.

The map consolidated runtime component lifecycle vocabulary, current source-first bootstrap composition, handler registration source files, application services, repository wiring, protocol adapter composition, lifecycle state, operations cross-reference, source surfaces, redaction posture, and implementation deferrals. The next bounded continuation should keep Pitaya-class architecture planning moving without turning the component lifecycle map into runtime behavior, dynamic handler registration, component discovery or loading, startup hooks, shutdown hooks, protocol, persistence, dependency, or distributed runtime authorization.

The safe next step is a boundary gate for handler module registration vocabulary. Pitaya-class architecture uses handler and module registration concepts to organize service ownership, route handling, component/module boundaries, and runtime extension points. vibit needs that vocabulary in an agent-native form before adding any handler registration behavior or component module loading.

## Decision

Select `define_pitaya_aligned_handler_module_registration_boundary_gate` as the next bounded Pitaya-aligned direction after the runtime component lifecycle source-first map.

Register `runtime.next_pitaya_aligned_direction_after_runtime_component_lifecycle_map` as the repository check rule.

Complete `M-210/W-0282` and open `M-211/W-0283 Define Pitaya-aligned handler module registration boundary gate` as next-ready.

This decision does not add runtime component lifecycle behavior. It also does not add handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement runtime component lifecycle behavior immediately.
- Implement handler registration behavior immediately.
- Add startup or shutdown hooks immediately.
- Add component discovery or component loading immediately.
- Define runtime component lifecycle again instead of moving to handler module registration.
- Implement dynamic handler registration or component module loading immediately.
- Return directly to Nakama product module expansion after the runtime component lifecycle map.
- Start distributed runtime implementation after the runtime component lifecycle vocabulary sequence.

## Rationale

Handler module registration is a high-value Pitaya alignment direction because it clarifies how future server components, route handlers, module ownership, extension points, and component lifecycle vocabulary relate to vibit's existing explicit route registration and application/runtime boundaries.

It is also broad enough to create accidental implementation scope. A selection-only step preserves the boundary-first sequence: the next work item can define vocabulary and deferrals before any handler registration behavior, dynamic module loading, component lifecycle code, startup hook, shutdown hook, protocol shape, persistence, dependency, dashboard, admin, hosted, SDK, release, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a selection-only continuation. The correct continuation is to record the selected follow-up direction, update repository memory, register the check rule, and open W-0283. Handler registration behavior and runtime component lifecycle behavior remain deferred until later gates and implementation slices explicitly authorize them.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  handler_module_registration_boundary_clarity: high
  implementation_boundedness: high
  agent_native_runtime_ownership_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

W-0283 becomes the next-ready work item and must define only a Pitaya-aligned handler module registration boundary gate.

This selection does not authorize runtime component lifecycle behavior, handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol changes, generated output, persistence changes, dependencies, hosted surfaces, SDK publication, distributed runtime behavior, or direct compatibility.

## Reversal Conditions

Revisit this decision if a later ADR chooses to return to Nakama product module expansion before any further Pitaya-aligned runtime architecture vocabulary, or if handler module registration needs to be split into separate gates for route handler ownership, component module ownership, dynamic registration, startup composition, or distributed runtime lifecycle.

## Follow-Up

- Complete `W-0283`: define the Pitaya-aligned handler module registration boundary gate.
- Keep runtime component lifecycle behavior, handler registration behavior, dynamic handler registration, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
