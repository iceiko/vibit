# ADR-0198: Pitaya-Aligned Startup And Shutdown Hook Source-First Map

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-startup-shutdown-hook-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-startup-shutdown-hook-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md`
- `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.zh-CN.md`
- `decisions/ADR-0197-pitaya-aligned-startup-shutdown-hook-boundary-gate.md`
- `decisions/ADR-0196-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map.md`
- `decisions/ADR-0195-pitaya-aligned-component-discovery-module-loading-source-first-map.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0197` defined the Pitaya-aligned startup and shutdown hook boundary gate and opened `M-218/W-0290`.

The repository already has explicit process startup, bootstrap source files, application services, domain modules, protocol adapter composition, route registration, unit-of-work boundaries, repository interfaces, and the component discovery and module loading source-first map. It does not have startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency containers, component discovery or loading behavior, component module loading behavior, distributed lifecycle behavior, or direct Nakama/Pitaya API compatibility.

## Decision

Implement `node tools/vibit inspect pitaya-startup-shutdown --json` as the source-first Pitaya-aligned startup and shutdown hook map for `M-218/W-0290`.

The command reports:

- `ADR-0197` as the source gate and `ADR-0198` as the implementation decision.
- `runtime.pitaya_aligned_startup_shutdown_hook_source_first_map` as the check rule.
- Allowed startup and shutdown hook vocabulary from `ADR-0197`.
- Related component discovery, component loading, handler module, runtime component, application dispatch, unit-of-work, operations, redaction, and distributed runtime deferral vocabulary.
- Current mappings for explicit bootstrap composition, startup hooks, shutdown hooks, lifecycle hook execution, dependency handoff, module loading handoff, and distributed lifecycle deferral.
- Source surfaces for architecture manifests, startup/shutdown gate, component loading gate, handler module registration gate, runtime component lifecycle gate, protocol standards, bootstrap files, application and module packages, protocol bridge files, transaction boundaries, `tools/vibit`, and rule catalog entries.
- Explicit false deferrals for startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoints, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility.
- `M-219/W-0291 Select next Pitaya-aligned direction after startup and shutdown hook map` as the next-ready follow-up.

This decision does not add startup hook behavior. It also does not add shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add startup hook behavior immediately.
- Add shutdown hook behavior immediately.
- Add lifecycle hook execution interfaces immediately.
- Add a dependency container or lifecycle scheduler immediately.
- Extend the component discovery and module loading source-first map instead of adding a startup/shutdown-focused map.
- Return directly to broad Nakama product module expansion after W-0290.

## Rationale

Startup and shutdown hooks are useful Pitaya-class architecture vocabulary, but they can easily imply runtime lifecycle behavior. A focused source-first inspection command gives agents the current repository map without authorizing hook execution, startup rewiring, shutdown scheduling, dependency containers, component loading, protocol changes, persistence changes, distributed lifecycle behavior, or compatibility work.

## Agent Reasoning Summary

The active work item is a source-first inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all startup hook, shutdown hook, lifecycle hook execution, dependency container, component discovery, component loading, component module loading, handler module registration, handler registration, dynamic registration, runtime endpoint, dashboard/admin, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  startup_shutdown_mapping_clarity: high
  implementation_boundedness: high
  source_first_runtime_composition_reuse: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned startup and shutdown hook vocabulary without reading every architecture document.

This decision does not add startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete startup hook, shutdown hook, lifecycle hook execution, dependency container, startup wiring, shutdown scheduling, runtime endpoint, or distributed lifecycle implementation model;
- the startup and shutdown hook inspection output creates confusion with public API compatibility or live lifecycle behavior;
- current bootstrap composition, route registration, payload registry, application services, domain modules, repository wiring, protocol adapter composition, process startup, or transport posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate source-first maps for startup ordering, shutdown ordering, lifecycle execution, dependency handoff, module loading handoff, or distributed lifecycle.

## Follow-Up

- Complete `W-0291`: select the next Pitaya-aligned direction after the startup and shutdown hook source-first map.
- Keep startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery/loading behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
