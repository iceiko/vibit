# ADR-0163: Pitaya-Aligned Distributed Group And Broadcast Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-distributed-group-broadcast-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-distributed-group-broadcast-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md`
- `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.zh-CN.md`
- `decisions/ADR-0162-pitaya-aligned-distributed-group-broadcast-boundary-gate.md`
- `decisions/ADR-0161-pitaya-aligned-service-discovery-source-first-map.md`
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

`ADR-0162` defined a gate-only Pitaya-aligned distributed group and broadcast boundary. It allowed `distributed_group`, `room_broadcast`, `broadcast_target`, `group_membership`, and `broadcast_fanout` as future architecture vocabulary and mapped current target-scope metadata and application-owned server-push intent to deferred group/broadcast concepts.

The safe next step is source-first inspection. Agents need a concrete repository command that summarizes group/broadcast vocabulary, current single-process mapping, and deferrals before any group membership registry, room state, fanout worker, delivery guarantee, stream subscription, cluster routing, protocol carrier, persistence, or dependency is authorized.

## Decision

Implement `node tools/vibit inspect pitaya-groups --json` as the source-first Pitaya-aligned distributed group and broadcast map for `M-183/W-0255`.

The command reports:

- ADR-0162 as the source gate and ADR-0163 as the implementation decision.
- `runtime.pitaya_aligned_distributed_group_broadcast_source_first_map` as the check rule.
- Allowed group/broadcast vocabulary: `distributed_group`, `room_broadcast`, `broadcast_target`, `group_membership`, and `broadcast_fanout`.
- Related vocabulary: `target_scope`, `server_push_intent`, `route_handler`, `module_handler`, `frontend_server`, `backend_server`, `service_discovery`, `server_to_server_rpc`, `remote_call`, and `cluster_safe_session_routing`.
- Current single-process mappings for target scope, server-push intent, distributed group, group membership, room broadcast, broadcast target, broadcast fanout, and cluster-safe session routing planning.
- Explicit false deferrals for distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, RPC implementation, remote calls, frontend/backend role implementation, distributed runtime implementation, cluster-safe routing, protocol, generated output, persistence, dependency, hosted, SDK, and direct compatibility surfaces.
- Redaction flags for credential, token, digest, key, DSN, database payload, local secret file, node credential, and transport metadata contents.
- `M-184/W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate` as the next-ready follow-up.

## Alternatives Considered

- Implement distributed groups or room broadcast immediately.
- Add group membership registries, fanout workers, or delivery guarantees while adding the map.
- Keep group and broadcast mapping only in ADR-0162 without a tool inspection surface.
- Fold group/broadcast vocabulary into `node tools/vibit inspect pitaya-discovery --json` instead of adding a focused command.

## Rationale

Distributed group and broadcast vocabulary is useful for Pitaya-aligned planning, but it is easy to confuse with permission to add group state, room routing, stream subscriptions, fanout topology, queues, retries, ordering, dependencies, or distributed runtime behavior. A dedicated source-first inspection command gives agents a precise place to inspect vocabulary, current mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all distributed group implementation, group membership registry, broadcast fanout, delivery guarantee, stream subscription, topology, protocol, persistence, dependency, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  group_broadcast_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned distributed group and broadcast vocabulary without reading every architecture document.

This decision does not add distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, groups, parties, chat rooms, matchmaking, match runtime, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, cluster-safe session routing behavior, runtime endpoint behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future cluster-safe session routing planning must start with W-0256 as a boundary gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete distributed group, room broadcast, or cluster-safe routing model;
- the group/broadcast inspection output creates confusion with public API compatibility;
- target-scope metadata, server-push intent, or connection/session binding ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for broadcast, groups, or cluster-safe routing.

## Follow-Up

- Complete `W-0256`: define a Pitaya-aligned cluster-safe session routing boundary gate.
- Keep distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, distributed runtime behavior, cluster-safe session routing behavior, runtime endpoint expansion, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
