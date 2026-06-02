# ADR-0188: Pitaya-Aligned Runtime Component Lifecycle Boundary Gate

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-runtime-component-lifecycle-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`
- `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md`
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

`ADR-0187` selected `define_pitaya_aligned_runtime_component_lifecycle_boundary_gate` as the next bounded Pitaya-aligned direction after the dashboard/admin operations source-first map.

The repository already has explicit process wiring, application bootstrap files, route registration source files, payload registry source files, application services, repository interfaces, unit-of-work boundaries, and transport close behavior. It does not yet have a ratified vocabulary boundary for describing runtime components, component lifecycle phases, handler registration ownership, start/shutdown semantics, component dependencies, local component inventory, or distributed lifecycle deferrals.

## Decision

Define the Pitaya-aligned runtime component lifecycle boundary gate for `M-208/W-0280`.

Register `runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate` as the repository check rule.

The gate allows future planning vocabulary for `runtime_component_boundary`, `component_lifecycle_phase`, `component_start_boundary`, `component_shutdown_boundary`, `handler_registration_boundary`, `component_dependency_boundary`, `bootstrap_composition_boundary`, `component_state_posture`, `local_component_inventory`, and `distributed_component_lifecycle_deferral`.

The gate maps current source-first bootstrap composition, route registration, payload registry, application service ownership, repository wiring, process startup, and transport close posture to that vocabulary.

Complete `M-208/W-0280` and open `M-209/W-0281 Implement Pitaya-aligned runtime component lifecycle source-first map` as next-ready.

The selected follow-up direction is `implement_pitaya_aligned_runtime_component_lifecycle_source_first_map`.

This decision does not add runtime component lifecycle behavior. It also does not add handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement runtime component lifecycle behavior immediately.
- Implement handler registration behavior immediately.
- Add startup and shutdown hooks immediately.
- Add component discovery or component loading immediately.
- Fold component lifecycle vocabulary into the dashboard/admin operations map.
- Return directly to Nakama product module expansion after the dashboard/admin operations map.

## Rationale

Runtime component lifecycle vocabulary is central to Pitaya-class architecture planning, but it can easily imply dynamic component loading, runtime hooks, handler registries, dependency containers, and distributed behavior.

A gate keeps this step vocabulary-only and source-first. It also lets a later source-first map report committed repository facts before any component lifecycle behavior, handler registration behavior, startup hook, shutdown hook, protocol, persistence, dependency, hosted, SDK, release, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a gate-only boundary. The correct continuation is to add the standard, translation, ADR, change artifacts, rule catalog entry, `tools/vibit` check coverage, and repository memory updates while preserving all runtime component lifecycle behavior, handler registration behavior, startup hook, shutdown hook, runtime endpoint, dashboard, admin, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

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

Agents can now use runtime component lifecycle vocabulary without treating it as permission to add component lifecycle behavior, handler registries, startup hooks, shutdown hooks, component discovery, dependency containers, endpoints, protocol changes, persistence, dependencies, hosted operations, or distributed runtime behavior.

This decision does not add runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later ADR selects a concrete runtime component lifecycle implementation model;
- the boundary creates confusion with dynamic handler registration, process startup hooks, shutdown hooks, dependency containers, or distributed component lifecycle;
- existing bootstrap composition, route registration, payload registry, application service ownership, repository wiring, process startup, or transport close posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate gates for component ownership, handler registration, startup, shutdown, dependency wiring, or distributed component lifecycle.

## Follow-Up

- Complete `W-0281`: implement a source-first Pitaya-aligned runtime component lifecycle map.
- Keep runtime component lifecycle behavior, handler registration behavior, startup hooks, shutdown hooks, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
