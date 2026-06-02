# ADR-0196: Select Next Pitaya-Aligned Direction After Component Discovery And Module Loading Map

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map.md`

Related artifacts:

- `decisions/ADR-0195-pitaya-aligned-component-discovery-module-loading-source-first-map.md`
- `decisions/ADR-0194-pitaya-aligned-component-discovery-module-loading-boundary-gate.md`
- `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md`
- `docs/pitaya-aligned-handler-module-registration-boundary-gate.md`
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

`ADR-0195` completed the source-first Pitaya-aligned component discovery and module loading map through `node tools/vibit inspect pitaya-component-loading --json`.

The map made current explicit bootstrap composition, application module ownership, handler module registration mapping, explicit component inventory, dynamic loading deferrals, distributed discovery deferrals, source surfaces, and redaction posture inspectable. The next continuation should keep Pitaya-class architecture planning moving without treating that map as permission to add component discovery behavior, component loading behavior, component module loading behavior, dynamic runtime loading, startup hooks, shutdown hooks, dependency containers, protocol changes, persistence changes, dependencies, distributed runtime behavior, or direct compatibility.

The safe next step is a boundary gate for startup and shutdown hook vocabulary. This is the smallest follow-up because startup and shutdown hooks are repeated deferrals in the runtime component lifecycle, handler module registration, and component discovery/module loading sequence, but they must remain separate from any concrete runtime hook behavior.

## Decision

Select `define_pitaya_aligned_startup_shutdown_hook_boundary_gate` as the next bounded Pitaya-aligned direction after the component discovery and module loading source-first map.

Register `runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map` as the repository check rule.

Complete `M-216/W-0288` and open `M-217/W-0289 Define Pitaya-aligned startup and shutdown hook boundary gate` as next-ready.

This decision does not add component discovery behavior. It also does not add component loading behavior, component module loading behavior, startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement component discovery behavior immediately.
- Implement component loading or component module loading immediately.
- Implement startup hooks or shutdown hooks immediately.
- Define dependency container behavior before startup and shutdown hook vocabulary.
- Return directly to Nakama product module expansion after the component discovery and module loading map.
- Start distributed runtime implementation after the component discovery and module loading vocabulary sequence.

## Rationale

Startup and shutdown hooks are useful Pitaya alignment concepts, but they can easily imply concrete lifecycle execution, dynamic module loading, dependency containers, process startup rewiring, runtime endpoints, plugin behavior, service discovery, or distributed runtime implementation.

A selection-only step preserves the boundary-first sequence. The next work item can define vocabulary and deferrals before any hook execution, component discovery, component loading, dynamic registration, protocol, persistence, dependency, dashboard, admin, hosted, SDK, release, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a selection-only continuation. The correct continuation is to record the selected follow-up direction, update repository memory, register the check rule, and open W-0289. Component discovery behavior, component loading behavior, component module loading behavior, startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, handler module registration behavior, handler registration behavior, and dynamic handler registration remain deferred until later gates and implementation slices explicitly authorize them.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  startup_shutdown_hook_boundary_clarity: high
  implementation_boundedness: high
  component_discovery_module_loading_follow_up_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

W-0289 becomes the next-ready work item and must define only a Pitaya-aligned startup and shutdown hook boundary gate.

This selection does not authorize component discovery behavior, component loading behavior, component module loading behavior, startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency containers, runtime endpoints, dashboards, admin console behavior, protocol changes, generated output, persistence changes, dependencies, hosted surfaces, SDK publication, distributed runtime behavior, or direct compatibility.

## Reversal Conditions

Revisit this decision if a later ADR chooses to return to Nakama product module expansion before any further Pitaya-aligned runtime architecture vocabulary, or if startup and shutdown hook vocabulary needs to be split into separate gates for startup ordering, shutdown ordering, lifecycle hook contracts, dependency handoff, module loading integration, local component inventory, or distributed runtime lifecycle.

## Follow-Up

- Complete `W-0289`: define the Pitaya-aligned startup and shutdown hook boundary gate.
- Keep component discovery behavior, component loading behavior, component module loading behavior, startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, runtime endpoints, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
