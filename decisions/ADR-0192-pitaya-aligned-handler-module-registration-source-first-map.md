# ADR-0192: Pitaya-Aligned Handler Module Registration Source-First Map

Status: Accepted
Date: 2026-06-02
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-handler-module-registration-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-handler-module-registration-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-handler-module-registration-boundary-gate.md`
- `docs/pitaya-aligned-handler-module-registration-boundary-gate.zh-CN.md`
- `decisions/ADR-0191-pitaya-aligned-handler-module-registration-boundary-gate.md`
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

`ADR-0191` defined the gate-only Pitaya-aligned handler module registration vocabulary boundary and opened `M-212/W-0284`.

The repository already has explicit bootstrap composition, route registration source files, payload registry source files, application service packages, domain module packages, unit-of-work and repository interface boundaries, protocol adapter bridges, and request/session context handoff boundaries. The safe next step is a source-first inspection command that reports those committed surfaces through the handler module registration vocabulary without adding handler module registration behavior, dynamic handler registration, component discovery/loading, component module loading, startup or shutdown hooks, protocol changes, dependencies, distributed runtime behavior, or direct compatibility.

## Decision

Implement `node tools/vibit inspect pitaya-handler-modules --json` as the source-first Pitaya-aligned handler module registration map for `M-212/W-0284`.

The command reports:

- ADR-0191 as the source gate and ADR-0192 as the implementation decision.
- `runtime.pitaya_aligned_handler_module_registration_source_first_map` as the check rule.
- Allowed handler module registration vocabulary from ADR-0191.
- Related vocabulary from runtime component lifecycle, route handler pipeline, serializer/message forwarding, application dispatch, unit-of-work boundaries, request/session context, source-first operations, redaction, and distributed runtime deferrals.
- Current mappings for module bootstrap, explicit route registration, payload registry ownership, application module ownership, dependency handoff, execution context handoff, and distributed registration deferral.
- Source surfaces for architecture manifests, the handler module registration gate, runtime protocol standards, bootstrap files, application and module packages, protocol bridge files, transaction boundaries, `tools/vibit`, and rule catalog entries.
- Explicit false deferrals for handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery/loading, component module loading, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility.
- `M-213/W-0285 Select next Pitaya-aligned direction after handler module registration map` as the next-ready follow-up.

This decision does not add handler module registration behavior. It also does not add handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add handler module registration behavior immediately.
- Add dynamic handler registration immediately.
- Add component discovery or component module loading immediately.
- Add startup or shutdown hooks immediately.
- Expand the runtime component lifecycle source-first map instead of adding a handler-module-focused map.
- Return directly to broad Nakama product module expansion after W-0283.

## Rationale

Handler module registration vocabulary is useful only if agents can inspect the current source tree without mistaking that visibility for permission to add runtime registration behavior. A focused command gives agents a bounded map before handler registries, dynamic registration, component discovery/loading, startup/shutdown hooks, dependency containers, protocol changes, distributed registration, or compatibility work is authorized.

## Agent Reasoning Summary

The active work item is a source-first inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all handler module registration, handler registration, dynamic registration, component discovery/loading, component module loading, startup/shutdown hook, runtime endpoint, dashboard/admin, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  handler_module_registration_mapping_clarity: high
  implementation_boundedness: high
  source_first_runtime_composition_reuse: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned handler module registration vocabulary without reading every architecture document.

This decision does not add handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery/loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete handler module registration, dynamic handler registration, component discovery/loading, startup/shutdown, dependency container, runtime endpoint, or distributed registration implementation model;
- the handler module registration inspection output creates confusion with public API compatibility or live registration behavior;
- current bootstrap composition, route registration, payload registry, application services, domain modules, repository wiring, protocol adapter composition, process startup, or transport close posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate source-first maps for handler modules, explicit route registration, payload registries, dependency handoff, execution context, local handler inventory, or distributed registration.

## Follow-Up

- Complete `W-0285`: select the next Pitaya-aligned direction after the handler module registration source-first map.
- Keep handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery/loading, component module loading, startup hooks, shutdown hooks, runtime endpoints, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
