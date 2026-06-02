# ADR-0194: Pitaya-Aligned Component Discovery And Module Loading Boundary Gate

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-component-discovery-module-loading-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-component-discovery-module-loading-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md`
- `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-handler-module-registration-boundary-gate.md`
- `decisions/ADR-0193-select-next-pitaya-aligned-direction-after-handler-module-registration-map.md`
- `decisions/ADR-0192-pitaya-aligned-handler-module-registration-source-first-map.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0193` selected `define_pitaya_aligned_component_discovery_module_loading_boundary_gate` as the next bounded Pitaya-aligned direction after the handler module registration source-first map.

The repository already has explicit process wiring, application bootstrap files, route registration source files, payload registry source files, application services, domain module packages, repository interfaces, unit-of-work boundaries, protocol adapter bridges, and a handler module registration source-first map. It does not yet have a ratified vocabulary boundary for component discovery boundaries, component loading boundaries, component module loading boundaries, explicit component inventory, bootstrap component sources, application module ownership, handler module registration handoff, dynamic loading deferral, or distributed component discovery deferrals.

## Decision

Define the Pitaya-aligned component discovery and module loading boundary gate for `M-214/W-0286`.

Register `runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate` as the repository check rule.

The gate allows future planning vocabulary for `component_discovery_boundary`, `component_loading_boundary`, `component_module_loading_boundary`, `explicit_component_inventory`, `bootstrap_component_source`, `application_module_ownership`, `handler_module_registration_handoff`, `dynamic_loading_deferral`, and `distributed_component_discovery_deferral`.

The gate maps current source-first bootstrap composition, application module ownership, handler module registration mapping, explicit source inventories, dynamic loading deferral, and distributed discovery deferral to that vocabulary.

Complete `M-214/W-0286` and open `M-215/W-0287 Implement Pitaya-aligned component discovery and module loading source-first map` as next-ready.

The selected follow-up direction is `implement_pitaya_aligned_component_discovery_module_loading_source_first_map`.

This decision does not add component discovery behavior. It also does not add component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement component discovery behavior immediately.
- Implement dynamic component loading or component module loading immediately.
- Add startup and shutdown hooks immediately.
- Add a runtime component registry immediately.
- Fold component discovery vocabulary into the handler module registration source-first map.
- Return directly to Nakama product module expansion after the handler module registration map.

## Rationale

Component discovery and module loading vocabulary is central to Pitaya-class architecture planning, but it can easily imply runtime scanning, dynamic module loaders, startup/shutdown hooks, dependency containers, protocol route changes, and distributed behavior.

A gate keeps this step vocabulary-only and source-first. It also lets a later source-first map report committed repository facts before any component discovery behavior, component loading behavior, component module loading behavior, startup hook, shutdown hook, protocol, persistence, dependency, hosted, SDK, release, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a gate-only boundary. The correct continuation is to add the standard, translation, ADR, change artifacts, rule catalog entry, `tools/vibit` check coverage, and repository memory updates while preserving all component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hook, shutdown hook, runtime endpoint, dashboard, admin, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  component_discovery_boundary_clarity: high
  component_loading_boundary_clarity: high
  implementation_boundedness: high
  agent_native_runtime_ownership_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now use component discovery and module loading vocabulary without treating it as permission to add runtime scanning, dynamic loaders, component registries, startup hooks, shutdown hooks, component module loading, dependency containers, endpoints, protocol changes, persistence, dependencies, hosted operations, or distributed runtime behavior.

This decision does not add component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later ADR selects a concrete component discovery or module loading implementation model;
- the boundary creates confusion with dynamic loading, process startup hooks, shutdown hooks, dependency containers, handler module registration, or distributed component discovery;
- existing bootstrap composition, route registration, payload registry, application service ownership, domain module ownership, repository wiring, request/session context handoff, process startup, or transport posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate gates for component inventory, module loading, bootstrap ownership, dependency wiring, or distributed component discovery.

## Follow-Up

- Complete `W-0287`: implement a source-first Pitaya-aligned component discovery and module loading map.
- Keep component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
