# ADR-0160: Pitaya-Aligned Service Discovery Boundary Gate

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-service-discovery-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-service-discovery-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-service-discovery-boundary-gate.md`
- `docs/pitaya-aligned-service-discovery-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `decisions/ADR-0159-pitaya-aligned-server-to-server-rpc-source-first-map.md`
- `decisions/ADR-0158-pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0159` implemented `node tools/vibit inspect pitaya-rpc --json` as a source-first map for Pitaya-aligned server-to-server RPC and remote-call vocabulary. It opened `M-180/W-0252` to define a service discovery boundary gate before discovery vocabulary turns into runtime topology, service registries, selectors, node identity, or dependencies.

Pitaya-style service discovery is useful architecture vocabulary for future distributed runtimes, but vibit's current runtime remains a statically composed single-process modular monolith with in-process dispatch and module handlers.

## Decision

Accept `docs/pitaya-aligned-service-discovery-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned service discovery vocabulary.

Register `runtime.pitaya_aligned_service_discovery_boundary_gate` as the repository check rule.

The gate defines:

- allowed service discovery vocabulary: `service_discovery`, `service_registry`, `service_instance`, and `service_selector`;
- related vocabulary: `frontend_server`, `backend_server`, `server_to_server_rpc`, `remote_call`, `route_handler`, `module_handler`, and `static_process_composition`;
- current single-process mapping for static composition, startup wiring, direct dispatch, route handlers, and module handlers;
- ownership for service discovery vocabulary and a future source-first service discovery map;
- stop conditions for any implementation behavior.

Complete `M-180/W-0252` and open `M-181/W-0253 Implement Pitaya-aligned service discovery source-first map` as next-ready.

This decision does not add service discovery implementation, service registries, service selectors, node registries, server identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement service discovery immediately.
- Add a service discovery source-first map without first defining a gate.
- Keep service discovery vocabulary only inside the RPC source-first map.
- Defer Pitaya service discovery vocabulary entirely until distributed runtime implementation.

## Rationale

Service discovery vocabulary is high-risk because it can be mistaken for permission to add node identity, registries, membership, routing, dependencies, or distributed topology. Defining a gate first lets vibit preserve Pitaya-aligned planning vocabulary while keeping the concrete runtime single-process.

This preserves vibit's agent-native model: vocabulary, ownership, deferrals, stop conditions, checks, and memory before implementation.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first service discovery map, not service discovery behavior, service registries, selectors, node identity, RPC, remote calls, or distributed runtime behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  service_discovery_boundary_clarity: high
  implementation_boundedness: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `service_discovery`, `service_registry`, `service_instance`, and `service_selector` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_service_discovery_boundary_gate` becomes the check rule for W-0252.
- `M-180/W-0252` is completed.
- `M-181/W-0253 Implement Pitaya-aligned service discovery source-first map` becomes next-ready.
- Service discovery implementation, service registries, selectors, node identity, server-to-server RPC, remote calls, frontend/backend role implementation, distributed runtime behavior, groups, broadcast fanout, cluster-safe session routing, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete distributed runtime or discovery model;
- the service discovery vocabulary creates confusion with public API compatibility;
- single-process startup composition or module handler ownership changes enough to require remapping;
- prototype feedback shows a stronger near-term need for a different product capability family.

## Follow-Up

- Complete `W-0253`: implement a source-first Pitaya-aligned service discovery map.
- Keep service discovery implementation, service registries, selectors, node identity, server-to-server RPC, remote calls, frontend/backend role implementation, distributed runtime behavior, distributed groups, cluster-safe session routing, runtime endpoint expansion, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
