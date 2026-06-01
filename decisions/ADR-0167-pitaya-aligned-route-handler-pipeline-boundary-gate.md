# ADR-0167: Pitaya-Aligned Route Handler Pipeline Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-route-handler-pipeline-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-route-handler-pipeline-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md`
- `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.zh-CN.md`
- `decisions/ADR-0166-select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map.md`
- `decisions/ADR-0165-pitaya-aligned-cluster-safe-session-routing-source-first-map.md`
- `docs/runtime-protocol-adapter.md`
- `docs/game-protocol.md`
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

`ADR-0166` selected `define_pitaya_aligned_route_handler_pipeline_boundary_gate` as the next bounded Pitaya-aligned direction after the cluster-safe session routing source-first map.

Pitaya-style route handler, handler pipeline, serializer, and forwarding vocabulary is useful architecture pressure for future distributed runtime planning. vibit's current runtime, however, already has an explicit single-process WebSocket Protobuf flow: protocol envelope, route request, application dispatch, command/query handlers, explicit Protobuf bridge functions, transactional dispatch, and outbound message handling.

## Decision

Accept `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned route handler pipeline vocabulary.

Register `runtime.pitaya_aligned_route_handler_pipeline_boundary_gate` as the repository check rule.

The gate defines:

- allowed route handler pipeline vocabulary: `route_handler`, `route_key`, `handler_dispatch`, `handler_pipeline`, `pipeline_step`, `serializer_boundary`, `message_forwarding`, and `route_target`;
- related vocabulary: `protocol_envelope`, `route_request`, `application_dispatch`, `command_handler`, `query_handler`, `protocol_bridge`, `target_scope`, `frontend_server`, `backend_server`, `server_to_server_rpc`, `remote_call`, `service_discovery`, and `cluster_safe_session_routing`;
- current single-process mapping for protocol envelope, route request, application dispatch, transactional dispatch, protocol bridge, and outbound message concepts;
- ownership for route handler pipeline vocabulary and a future source-first map;
- stop conditions for any implementation behavior.

Complete `M-187/W-0259` and open `M-188/W-0260 Implement Pitaya-aligned route handler pipeline source-first map` as next-ready.

This decision does not add route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, distributed runtime behavior, distributed groups, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement route handlers or handler pipelines immediately.
- Add a source-first route handler pipeline map without first defining a gate.
- Fold route handler vocabulary into the cluster-safe session routing map.
- Return directly to Nakama product module expansion before route handler pipeline vocabulary is bounded.

## Rationale

Route handler pipeline vocabulary is high-risk because it can be mistaken for permission to add handler routing behavior, middleware chains, serializer plugins, forwarding workers, backend route targeting, service discovery, RPC, protocol changes, dependencies, or distributed runtime behavior. Defining a gate first lets vibit preserve Pitaya-aligned planning vocabulary while keeping the concrete runtime single-process and source-first.

This preserves vibit's agent-native model: vocabulary, ownership, deferrals, stop conditions, checks, and memory before implementation.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first route handler pipeline map, not route handlers, handler routing behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, service discovery, RPC, remote calls, protocol changes, generated output, or distributed runtime behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  route_handler_pipeline_boundary_clarity: high
  implementation_boundedness: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `route_handler`, `route_key`, `handler_dispatch`, `handler_pipeline`, `pipeline_step`, `serializer_boundary`, `message_forwarding`, and `route_target` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_route_handler_pipeline_boundary_gate` becomes the check rule for W-0259.
- `M-187/W-0259` is completed.
- `M-188/W-0260 Implement Pitaya-aligned route handler pipeline source-first map` becomes next-ready.
- Route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, service discovery implementation, RPC, remote calls, frontend/backend role implementation, cluster-safe session routing, distributed runtime behavior, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete route handler pipeline or distributed runtime model;
- the route handler pipeline vocabulary creates confusion with public API compatibility;
- protocol adapter, application dispatch, transactional dispatch, protocol bridge, or outbound delivery ownership changes enough to require remapping;
- prototype feedback shows a stronger near-term need for a different product capability family.

## Follow-Up

- Complete `W-0260`: implement a source-first Pitaya-aligned route handler pipeline map.
- Keep route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, cluster-safe session routing behavior, distributed runtime behavior, distributed groups, room broadcast fanout, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
