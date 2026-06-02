# ADR-0191: Pitaya-Aligned Handler Module Registration Boundary Gate

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-handler-module-registration-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-handler-module-registration-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-handler-module-registration-boundary-gate.md`
- `docs/pitaya-aligned-handler-module-registration-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`
- `decisions/ADR-0190-select-next-pitaya-aligned-direction-after-runtime-component-lifecycle-map.md`
- `decisions/ADR-0189-pitaya-aligned-runtime-component-lifecycle-source-first-map.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0190` selected `define_pitaya_aligned_handler_module_registration_boundary_gate` as the next bounded Pitaya-aligned direction after the runtime component lifecycle source-first map.

The repository already has explicit process wiring, application bootstrap files, route registration source files, payload registry source files, application services, domain module packages, repository interfaces, unit-of-work boundaries, protocol adapter bridges, and transport close behavior. It does not yet have a ratified vocabulary boundary for describing handler module ownership, explicit handler registration sources, payload registry ownership, handler dependency handoff, handler execution context handoff, local handler inventory, dynamic registration deferral, or distributed handler/module registration deferrals.

## Decision

Define the Pitaya-aligned handler module registration boundary gate for `M-211/W-0283`.

Register `runtime.pitaya_aligned_handler_module_registration_boundary_gate` as the repository check rule.

The gate allows future planning vocabulary for `handler_module_boundary`, `handler_registration_boundary`, `route_handler_ownership`, `explicit_registration_source`, `payload_registry_boundary`, `module_bootstrap_boundary`, `handler_dependency_boundary`, `handler_execution_context_boundary`, `local_handler_inventory`, `dynamic_registration_deferral`, and `distributed_handler_module_deferral`.

The gate maps current source-first bootstrap composition, explicit route registration, payload registry and bridge ownership, application module ownership, domain module ownership, repository and unit-of-work handoff, request/session context handoff, and distributed registration deferral to that vocabulary.

Complete `M-211/W-0283` and open `M-212/W-0284 Implement Pitaya-aligned handler module registration source-first map` as next-ready.

The selected follow-up direction is `implement_pitaya_aligned_handler_module_registration_source_first_map`.

This decision does not add handler module registration behavior. It also does not add handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement handler module registration behavior immediately.
- Implement dynamic handler registration behavior immediately.
- Add startup and shutdown hooks immediately.
- Add component discovery or component module loading immediately.
- Fold handler module registration vocabulary into the runtime component lifecycle source-first map.
- Return directly to Nakama product module expansion after the runtime component lifecycle map.

## Rationale

Handler module registration vocabulary is central to Pitaya-class architecture planning, but it can easily imply dynamic handler registries, component/module loading, startup hooks, dependency containers, protocol route changes, and distributed behavior.

A gate keeps this step vocabulary-only and source-first. It also lets a later source-first map report committed repository facts before any handler module registration behavior, handler registration behavior, dynamic registration, startup hook, shutdown hook, protocol, persistence, dependency, hosted, SDK, release, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a gate-only boundary. The correct continuation is to add the standard, translation, ADR, change artifacts, rule catalog entry, `tools/vibit` check coverage, and repository memory updates while preserving all handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery/loading, component module loading, startup hook, shutdown hook, runtime endpoint, dashboard, admin, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

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

Agents can now use handler module registration vocabulary without treating it as permission to add dynamic handler registries, startup hooks, shutdown hooks, component discovery, component module loading, dependency containers, endpoints, protocol changes, persistence, dependencies, hosted operations, or distributed runtime behavior.

This decision does not add handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later ADR selects a concrete handler module registration implementation model;
- the boundary creates confusion with dynamic handler registration, process startup hooks, shutdown hooks, dependency containers, component module loading, or distributed handler/module registration;
- existing bootstrap composition, route registration, payload registry, application service ownership, domain module ownership, repository wiring, request/session context handoff, process startup, or transport close posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate gates for handler ownership, handler registration, module bootstrap, dependency wiring, execution context, or distributed handler/module registration.

## Follow-Up

- Complete `W-0284`: implement a source-first Pitaya-aligned handler module registration map.
- Keep handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
