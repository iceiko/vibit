# ADR-0168: Pitaya-Aligned Route Handler Pipeline Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-route-handler-pipeline-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-route-handler-pipeline-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md`
- `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.zh-CN.md`
- `decisions/ADR-0167-pitaya-aligned-route-handler-pipeline-boundary-gate.md`
- `decisions/ADR-0166-select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map.md`
- `decisions/ADR-0165-pitaya-aligned-cluster-safe-session-routing-source-first-map.md`
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

`ADR-0167` defined a gate-only Pitaya-aligned route handler pipeline vocabulary boundary. It allowed `route_handler`, `route_key`, `handler_dispatch`, `handler_pipeline`, `pipeline_step`, `serializer_boundary`, `message_forwarding`, and `route_target` as future architecture vocabulary while preserving vibit's current single-process WebSocket Protobuf route flow.

The safe next step is source-first inspection. Agents need a concrete repository command that summarizes route handler pipeline vocabulary, current protocol/dispatch/bridge mapping, redaction posture, and deferrals before any route handler implementation, handler routing, handler pipeline, serializer, forwarding, backend route targeting, service discovery, RPC, distributed runtime, protocol carrier, persistence, or dependency is authorized.

## Decision

Implement `node tools/vibit inspect pitaya-routes --json` as the source-first Pitaya-aligned route handler pipeline map for `M-188/W-0260`.

The command reports:

- ADR-0167 as the source gate and ADR-0168 as the implementation decision.
- `runtime.pitaya_aligned_route_handler_pipeline_source_first_map` as the check rule.
- Allowed route handler pipeline vocabulary: `route_handler`, `route_key`, `handler_dispatch`, `handler_pipeline`, `pipeline_step`, `serializer_boundary`, `message_forwarding`, and `route_target`.
- Related vocabulary: `protocol_envelope`, `route_request`, `application_dispatch`, `command_handler`, `query_handler`, `protocol_bridge`, `target_scope`, `frontend_server`, `backend_server`, `server_to_server_rpc`, `remote_call`, `service_discovery`, and `cluster_safe_session_routing`.
- Current single-process mappings for protocol envelope, route request, application dispatch, transactional dispatch, protocol bridge, outbound message, and route target planning.
- Explicit false deferrals for route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, cluster-safe session routing, distributed runtime, service discovery, RPC, remote calls, frontend/backend roles, distributed groups, room broadcast fanout, protocol, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility.
- Redaction flags for credential, token, digest, key, DSN, database payload, local secret file, node credential, transport metadata, and route payload contents.
- `M-189/W-0261 Select next Pitaya-aligned direction after route handler pipeline map` as the next-ready follow-up.

## Alternatives Considered

- Implement route handlers or handler pipelines immediately.
- Add handler routing, middleware chains, serializer plugins, forwarding workers, or backend route targeting while adding the map.
- Keep route handler pipeline mapping only in ADR-0167 without a tool inspection surface.
- Fold route handler pipeline vocabulary into `node tools/vibit inspect pitaya-sessions --json` instead of adding a focused command.

## Rationale

Route handler pipeline vocabulary is useful for Pitaya-aligned planning, but it is easy to confuse with permission to add handler routing behavior, middleware chains, serializer plugins, forwarding workers, backend route targeting, service discovery, RPC, transport carriers, dependencies, or multi-node behavior. A dedicated source-first inspection command gives agents a precise place to inspect vocabulary, current mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all route handler implementation, handler routing, pipeline, middleware, serializer, forwarding, backend targeting, service discovery, RPC, distributed runtime, protocol, persistence, dependency, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  route_handler_pipeline_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned route handler pipeline vocabulary without reading every architecture document.

This decision does not add route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, runtime endpoint behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future Pitaya-aligned direction planning must start with W-0261 as a selection gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete route handler pipeline or distributed runtime model;
- the route handler pipeline inspection output creates confusion with public API compatibility;
- protocol adapter, application dispatch, transactional dispatch, protocol bridge, outbound delivery, or target-scope ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for serializer boundaries, forwarding, or backend route targeting.

## Follow-Up

- Complete `W-0261`: select the next Pitaya-aligned direction after the route handler pipeline map.
- Keep route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, cluster-safe session routing behavior, distributed runtime behavior, distributed groups, room broadcast fanout, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
