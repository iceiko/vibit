# ADR-0162: Pitaya-Aligned Distributed Group And Broadcast Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-distributed-group-broadcast-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-distributed-group-broadcast-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md`
- `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-service-discovery-boundary-gate.md`
- `decisions/ADR-0161-pitaya-aligned-service-discovery-source-first-map.md`
- `decisions/ADR-0160-pitaya-aligned-service-discovery-boundary-gate.md`
- `docs/first-server-push-realtime-messaging-gate.md`
- `docs/realtime-protocol-websocket-outbound-delivery-gate.md`
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

`ADR-0161` implemented `node tools/vibit inspect pitaya-discovery --json` as a source-first map for Pitaya-aligned service discovery vocabulary. It opened `M-182/W-0254` to define a distributed group and broadcast boundary gate before group, room, stream, fanout, or cluster routing vocabulary turns into runtime behavior.

Pitaya-style groups and broadcast are useful architecture vocabulary for future distributed runtimes, but vibit's current runtime remains a single-process server with target-scope metadata, narrow application-owned server-push intent, and single-process outbound delivery.

## Decision

Accept `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned distributed group and broadcast vocabulary.

Register `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate` as the repository check rule.

The gate defines:

- allowed group and broadcast vocabulary: `distributed_group`, `room_broadcast`, `broadcast_target`, `group_membership`, and `broadcast_fanout`;
- related vocabulary: `target_scope`, `server_push_intent`, `route_handler`, `module_handler`, `frontend_server`, `backend_server`, `service_discovery`, `server_to_server_rpc`, `remote_call`, and `cluster_safe_session_routing`;
- current single-process mapping for target-scope metadata, server-push intent, distributed group, group membership, room broadcast, broadcast target, and broadcast fanout concepts;
- ownership for distributed group and broadcast vocabulary and a future source-first map;
- stop conditions for any implementation behavior.

Complete `M-182/W-0254` and open `M-183/W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map` as next-ready.

This decision does not add distributed group implementation, room broadcast fanout, delivery guarantees, stream subscriptions, group membership registries, groups, parties, chat rooms, matchmaking, match runtime, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement distributed groups or room broadcast immediately.
- Add a source-first group/broadcast map without first defining a gate.
- Keep group and broadcast vocabulary only inside the general Pitaya vocabulary or service discovery map.
- Defer group and broadcast vocabulary entirely until distributed runtime implementation.

## Rationale

Group and broadcast vocabulary is high-risk because it can be mistaken for permission to add membership registries, room state, stream subscriptions, fanout workers, delivery guarantees, queues, retries, ordering, topology, dependencies, or distributed routing. Defining a gate first lets vibit preserve Pitaya-aligned planning vocabulary while keeping the concrete runtime single-process.

This preserves vibit's agent-native model: vocabulary, ownership, deferrals, stop conditions, checks, and memory before implementation.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first group/broadcast map, not distributed group behavior, room broadcast fanout, service discovery, RPC, remote calls, or distributed runtime behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  group_broadcast_boundary_clarity: high
  implementation_boundedness: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `distributed_group`, `room_broadcast`, `broadcast_target`, `group_membership`, and `broadcast_fanout` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate` becomes the check rule for W-0254.
- `M-182/W-0254` is completed.
- `M-183/W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map` becomes next-ready.
- Distributed groups, room broadcast fanout, group membership registries, stream subscriptions, delivery guarantees, service discovery implementation, RPC, remote calls, frontend/backend role implementation, distributed runtime behavior, cluster-safe session routing, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete distributed group, broadcast, or cluster routing model;
- the group/broadcast vocabulary creates confusion with public API compatibility;
- target-scope metadata or server-push intent ownership changes enough to require remapping;
- prototype feedback shows a stronger near-term need for a different product capability family.

## Follow-Up

- Complete `W-0255`: implement a source-first Pitaya-aligned distributed group and broadcast map.
- Keep distributed group implementation, room broadcast fanout, group membership registries, stream subscriptions, delivery guarantees, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, distributed runtime behavior, cluster-safe session routing, runtime endpoint expansion, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
