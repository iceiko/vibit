# ADR-0197: Pitaya-Aligned Startup And Shutdown Hook Boundary Gate

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-startup-shutdown-hook-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-startup-shutdown-hook-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md`
- `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.zh-CN.md`
- `decisions/ADR-0196-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map.md`
- `decisions/ADR-0195-pitaya-aligned-component-discovery-module-loading-source-first-map.md`
- `decisions/ADR-0194-pitaya-aligned-component-discovery-module-loading-boundary-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0196` selected `define_pitaya_aligned_startup_shutdown_hook_boundary_gate` after the component discovery and module loading source-first map.

Startup and shutdown hook vocabulary is useful for Pitaya alignment, but the current runtime remains explicitly composed through committed Go startup and bootstrap source files. There is no hook execution contract, dependency container, lifecycle hook scheduler, dynamic module loading, distributed lifecycle model, or direct Nakama/Pitaya API compatibility.

## Decision

Define `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md` and its paired Simplified Chinese translation as the startup and shutdown hook vocabulary gate.

Register `runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate` as the repository check rule.

Complete `M-217/W-0289` and open `M-218/W-0290 Implement Pitaya-aligned startup and shutdown hook source-first map` as next-ready.

The selected follow-up direction is `implement_pitaya_aligned_startup_shutdown_hook_source_first_map`.

This decision does not add startup hook behavior. It also does not add shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement startup hook behavior immediately.
- Implement shutdown hook behavior immediately.
- Add lifecycle hook execution interfaces immediately.
- Add a dependency container before hook vocabulary.
- Add dynamic module loading before hook vocabulary.
- Return directly to Nakama product module expansion after W-0288.

## Rationale

Startup and shutdown hooks are often coupled to component loading, lifecycle ordering, dependency containers, process startup wiring, and distributed runtime behavior. In vibit, those concerns need a source-first map before any runtime implementation slice.

A gate-only decision lets agents use the vocabulary for future planning while preserving current local alpha behavior and all implementation deferrals.

## Agent Reasoning Summary

The active work item is a boundary gate. The correct continuation is to add the standard, ADR, rule, change artifacts, and repository memory updates while opening a future source-first inspection map. Startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency containers, component discovery/loading behavior, protocol changes, persistence changes, dependencies, distributed runtime behavior, and direct compatibility remain deferred.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  startup_shutdown_hook_boundary_clarity: high
  implementation_boundedness: high
  source_first_follow_up_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now use startup and shutdown hook vocabulary without treating it as permission to add hook behavior or startup/shutdown wiring.

This decision does not add startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoints, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if a later architecture ADR selects a concrete startup hook, shutdown hook, lifecycle hook execution, dependency container, component loading, dynamic loading, process startup wiring, runtime endpoint, or distributed lifecycle implementation model.

## Follow-Up

- Complete `W-0290`: implement a source-first startup and shutdown hook inspection map.
- Keep startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery/loading behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
