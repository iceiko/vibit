# ADR-0195: Pitaya-Aligned Component Discovery And Module Loading Source-First Map

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-component-discovery-module-loading-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-component-discovery-module-loading-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md`
- `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.zh-CN.md`
- `decisions/ADR-0194-pitaya-aligned-component-discovery-module-loading-boundary-gate.md`
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

`ADR-0194` defined the Pitaya-aligned component discovery and module loading boundary gate and opened `M-215/W-0287`.

The repository already has explicit process startup, bootstrap source files, application service packages, domain module packages, route registration source files, payload registry source files, protocol bridge files, unit-of-work and repository interface boundaries, and a handler module registration source-first map. It does not have component discovery behavior, component loading behavior, component module loading behavior, runtime scanning, dynamic module loading, startup hooks, shutdown hooks, dependency containers, distributed component discovery, or direct Nakama/Pitaya API compatibility.

## Decision

Implement `node tools/vibit inspect pitaya-component-loading --json` as the source-first Pitaya-aligned component discovery and module loading map for `M-215/W-0287`.

The command reports:

- `ADR-0194` as the source gate and `ADR-0195` as the implementation decision.
- `runtime.pitaya_aligned_component_discovery_module_loading_source_first_map` as the check rule.
- Allowed component discovery and module loading vocabulary from `ADR-0194`.
- Related handler module registration, runtime component lifecycle, application dispatch, unit-of-work, operations, redaction, and distributed runtime deferral vocabulary.
- Current mappings for explicit bootstrap composition, application module ownership, handler module registration mapping, explicit component inventory, dynamic loading deferral, and distributed discovery deferral.
- Source surfaces for architecture manifests, component discovery and module loading gate, handler module registration gate, runtime protocol standards, bootstrap files, application and module packages, protocol bridge files, transaction boundaries, `tools/vibit`, and rule catalog entries.
- Explicit false deferrals for component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility.
- `M-216/W-0288 Select next Pitaya-aligned direction after component discovery and module loading map` as the next-ready follow-up.

This decision does not add component discovery behavior. It also does not add component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add component discovery behavior immediately.
- Add dynamic component loading or a component module loader immediately.
- Add startup and shutdown hooks immediately.
- Add a dependency container or runtime component registry immediately.
- Extend the handler module registration source-first map instead of adding a component-loading-focused map.
- Return directly to broad Nakama product module expansion after W-0287.

## Rationale

Component discovery and module loading are useful Pitaya-class architecture vocabulary, but they can easily imply dynamic runtime behavior. A focused source-first inspection command gives agents the current repository map without authorizing scanners, loaders, registries, dependency containers, lifecycle hooks, protocol changes, persistence changes, distributed discovery, or compatibility work.

## Agent Reasoning Summary

The active work item is a source-first inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all component discovery, component loading, component module loading, handler module registration, handler registration, dynamic registration, startup/shutdown hook, runtime endpoint, dashboard/admin, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  component_discovery_mapping_clarity: high
  component_loading_mapping_clarity: high
  implementation_boundedness: high
  source_first_runtime_composition_reuse: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned component discovery and module loading vocabulary without reading every architecture document.

This decision does not add component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete component discovery, component loading, component module loading, dynamic loading, dependency container, startup/shutdown, runtime endpoint, or distributed discovery implementation model;
- the component discovery and module loading inspection output creates confusion with public API compatibility or live loading behavior;
- current bootstrap composition, route registration, payload registry, application services, domain modules, repository wiring, protocol adapter composition, process startup, or transport posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate source-first maps for component inventory, bootstrap ownership, module loading, dependency wiring, local component inventory, or distributed discovery.

## Follow-Up

- Complete `W-0288`: select the next Pitaya-aligned direction after the component discovery and module loading source-first map.
- Keep component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoints, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
