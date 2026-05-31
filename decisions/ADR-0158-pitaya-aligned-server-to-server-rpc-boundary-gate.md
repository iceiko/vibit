# ADR-0158: Pitaya-Aligned Server To Server RPC Boundary Gate

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-server-to-server-rpc-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-server-to-server-rpc-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.md`
- `decisions/ADR-0157-pitaya-aligned-frontend-backend-role-source-first-map.md`
- `decisions/ADR-0156-pitaya-aligned-frontend-backend-role-boundary-gate.md`
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

`ADR-0157` implemented `node tools/vibit inspect pitaya-roles --json` as a source-first map for Pitaya-aligned frontend/backend role vocabulary. It opened `M-178/W-0250` to define a server-to-server RPC boundary gate before any RPC planning turns into implementation.

Pitaya-style server-to-server RPC and remote calls are useful architecture vocabulary, but they can easily imply transports, service discovery, remote handler invocation, process topology, or protocol additions. vibit's current runtime remains a single-process modular monolith with in-process application dispatch and module handlers.

## Decision

Accept `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned server-to-server RPC and remote-call vocabulary.

Register `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate` as the repository check rule.

The gate defines:

- allowed RPC vocabulary: `server_to_server_rpc` and `remote_call`;
- related vocabulary: `route_handler`, `module_handler`, `application_dispatch`, and `service_discovery`;
- current single-process mapping for application dispatch, module handlers, route handling, and absent service discovery;
- ownership for RPC vocabulary and a future source-first RPC map;
- stop conditions for any implementation behavior.

Complete `M-178/W-0250` and open `M-179/W-0251 Implement Pitaya-aligned server-to-server RPC source-first map` as next-ready.

This decision does not add server-to-server RPC implementation, remote call behavior, service discovery, frontend/backend server role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement server-to-server RPC immediately.
- Add remote-call and service-discovery inspection without first defining a gate.
- Keep RPC vocabulary only in the broader Pitaya vocabulary map.
- Switch away from Pitaya architecture vocabulary after the role map.

## Rationale

RPC vocabulary is one of the highest-risk Pitaya-aligned concepts because it can be mistaken for permission to add distributed runtime behavior. Defining a gate first makes the planning vocabulary precise while preserving vibit's current single-process runtime and contract-first module boundaries.

This preserves vibit's agent-native model: vocabulary, ownership, deferrals, stop conditions, checks, and memory before implementation.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first RPC map, not RPC behavior, remote calls, service discovery, or distributed runtime behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  rpc_boundary_clarity: high
  implementation_boundedness: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `server_to_server_rpc` and `remote_call` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate` becomes the check rule for W-0250.
- `M-178/W-0250` is completed.
- `M-179/W-0251 Implement Pitaya-aligned server-to-server RPC source-first map` becomes next-ready.
- RPC implementation, remote calls, service discovery, frontend/backend role implementation, distributed runtime behavior, groups, broadcast fanout, cluster-safe session routing, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete distributed runtime or RPC model;
- the RPC vocabulary creates confusion with public API compatibility;
- single-process application dispatch or module handler ownership changes enough to require remapping;
- prototype feedback shows a stronger near-term need for a different product capability family.

## Follow-Up

- Complete `W-0251`: implement a source-first Pitaya-aligned server-to-server RPC map.
- Keep RPC implementation, remote calls, service discovery, frontend/backend role implementation, distributed runtime behavior, distributed groups, cluster-safe session routing, runtime endpoint expansion, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
