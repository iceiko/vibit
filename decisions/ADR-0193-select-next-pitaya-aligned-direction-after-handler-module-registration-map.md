# ADR-0193: Select Next Pitaya-Aligned Direction After Handler Module Registration Map

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-handler-module-registration-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-handler-module-registration-map.md`

Related artifacts:

- `decisions/ADR-0192-pitaya-aligned-handler-module-registration-source-first-map.md`
- `decisions/ADR-0191-pitaya-aligned-handler-module-registration-boundary-gate.md`
- `docs/pitaya-aligned-handler-module-registration-boundary-gate.md`
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

`ADR-0192` completed the source-first Pitaya-aligned handler module registration map through `node tools/vibit inspect pitaya-handler-modules --json`.

The map made current bootstrap composition, explicit route registration, payload registry ownership, application module ownership, dependency handoff, execution context handoff, local handler inventory, and distributed handler/module registration deferrals inspectable. The next continuation should keep Pitaya-class architecture planning moving without treating that map as permission to add handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery or module loading behavior, startup hooks, shutdown hooks, protocol changes, persistence changes, dependencies, distributed runtime behavior, or direct compatibility.

The safe next step is a boundary gate for component discovery and module loading vocabulary. This is the smallest follow-up because component discovery, component loading, and component module loading are repeated deferrals in the runtime component lifecycle and handler module registration sequence.

## Decision

Select `define_pitaya_aligned_component_discovery_module_loading_boundary_gate` as the next bounded Pitaya-aligned direction after the handler module registration source-first map.

Register `runtime.next_pitaya_aligned_direction_after_handler_module_registration_map` as the repository check rule.

Complete `M-213/W-0285` and open `M-214/W-0286 Define Pitaya-aligned component discovery and module loading boundary gate` as next-ready.

This decision does not add handler module registration behavior. It also does not add handler registration behavior, dynamic handler registration, component discovery or module loading behavior, component module loading behavior, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement handler module registration behavior immediately.
- Implement component discovery or component module loading immediately.
- Implement dynamic handler registration immediately.
- Add startup or shutdown hooks immediately.
- Define dependency container behavior before discovery and loading vocabulary.
- Return directly to Nakama product module expansion after the handler module registration map.
- Start distributed runtime implementation after the handler module registration vocabulary sequence.

## Rationale

Component discovery and module loading are useful Pitaya alignment concepts, but they can easily imply dynamic runtime loading, dependency containers, startup and shutdown hook behavior, protocol route changes, plugin behavior, service discovery, or distributed runtime implementation.

A selection-only step preserves the boundary-first sequence. The next work item can define vocabulary and deferrals before any discovery, loading, dynamic registration, startup hook, shutdown hook, protocol, persistence, dependency, dashboard, admin, hosted, SDK, release, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a selection-only continuation. The correct continuation is to record the selected follow-up direction, update repository memory, register the check rule, and open W-0286. Handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery behavior, component loading behavior, and component module loading behavior remain deferred until later gates and implementation slices explicitly authorize them.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  component_discovery_module_loading_boundary_clarity: high
  implementation_boundedness: high
  handler_module_registration_follow_up_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

W-0286 becomes the next-ready work item and must define only a Pitaya-aligned component discovery and module loading boundary gate.

This selection does not authorize handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery behavior, component loading behavior, component module loading behavior, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol changes, generated output, persistence changes, dependencies, hosted surfaces, SDK publication, distributed runtime behavior, or direct compatibility.

## Reversal Conditions

Revisit this decision if a later ADR chooses to return to Nakama product module expansion before any further Pitaya-aligned runtime architecture vocabulary, or if component discovery and module loading need to be split into separate gates for discovery inventory, static module declarations, dynamic loading, startup composition, shutdown ordering, dependency handoff, or distributed runtime lifecycle.

## Follow-Up

- Complete `W-0286`: define the Pitaya-aligned component discovery and module loading boundary gate.
- Keep handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery behavior, component loading behavior, component module loading behavior, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
