# ADR-0171: Pitaya-Aligned Serializer And Message Forwarding Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-serializer-message-forwarding-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-serializer-message-forwarding-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md`
- `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.zh-CN.md`
- `decisions/ADR-0170-pitaya-aligned-serializer-message-forwarding-boundary-gate.md`
- `decisions/ADR-0169-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map.md`
- `decisions/ADR-0168-pitaya-aligned-route-handler-pipeline-source-first-map.md`
- `decisions/ADR-0167-pitaya-aligned-route-handler-pipeline-boundary-gate.md`
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

`ADR-0170` defined a gate-only Pitaya-aligned serializer and message forwarding vocabulary boundary. It allowed `serializer_boundary`, `serializer_format`, `encode_boundary`, `decode_boundary`, `message_forwarding`, `forwarding_target`, `forwarding_envelope`, and `delivery_handoff` as future architecture vocabulary while preserving vibit's current single-process WebSocket Protobuf route flow.

The safe next step is source-first inspection. Agents need a concrete repository command that summarizes serializer and message forwarding vocabulary, current protocol bridge, generated payload bridge, outbound message, target-scope metadata, delivery handoff mapping, redaction posture, and deferrals before any serializer behavior, message forwarding behavior, backend route targeting, service discovery, RPC, distributed runtime, protocol carrier, persistence, or dependency is authorized.

## Decision

Implement `node tools/vibit inspect pitaya-serializer-forwarding --json` as the source-first Pitaya-aligned serializer and message forwarding map for `M-191/W-0263`.

The command reports:

- ADR-0170 as the source gate and ADR-0171 as the implementation decision.
- `runtime.pitaya_aligned_serializer_message_forwarding_source_first_map` as the check rule.
- Allowed serializer and message forwarding vocabulary: `serializer_boundary`, `serializer_format`, `encode_boundary`, `decode_boundary`, `message_forwarding`, `forwarding_target`, `forwarding_envelope`, and `delivery_handoff`.
- Related vocabulary: `protocol_bridge`, `generated_payload_bridge`, `protobuf_envelope`, `payload_encoding`, `payload_decoding`, `outbound_message`, `target_scope`, `route_target`, `delivery_handoff`, `route_handler`, `frontend_server`, `backend_server`, `server_to_server_rpc`, `remote_call`, `service_discovery`, and `cluster_safe_session_routing`.
- Current single-process mappings for protocol bridge, envelope encoding, payload encoding, payload decoding, outbound message, target scope, forwarding envelope, and delivery handoff planning.
- Explicit false deferrals for route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, cluster-safe session routing, distributed runtime, service discovery, RPC, remote calls, frontend/backend roles, distributed groups, room broadcast fanout, protocol, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility.
- Redaction flags for credential, token, digest, key, DSN, database payload, local secret file, node credential, transport metadata, route payload contents, and forwarding payload contents.
- `M-192/W-0264 Select next Pitaya-aligned direction after serializer and message forwarding map` as the next-ready follow-up.

## Alternatives Considered

- Implement serializer behavior or message forwarding immediately.
- Add codec plugins, serializer registries, forwarding workers, remote delivery handoff, or backend route targeting while adding the map.
- Keep serializer and message forwarding mapping only in ADR-0170 without a tool inspection surface.
- Fold serializer and message forwarding vocabulary into `node tools/vibit inspect pitaya-sessions --json` instead of adding a focused command.

## Rationale

Serializer and message forwarding vocabulary is useful for Pitaya-aligned planning, but it is easy to confuse with permission to add codec plugins, serializer registries, forwarding workers, remote delivery handoff, backend route targeting, service discovery, RPC, transport carriers, dependencies, or multi-node behavior. A dedicated source-first inspection command gives agents a precise place to inspect vocabulary, current mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all serializer behavior, message forwarding behavior, route handler implementation, handler routing, pipeline, middleware, backend targeting, service discovery, RPC, distributed runtime, protocol, persistence, dependency, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  serializer_message_forwarding_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned serializer and message forwarding vocabulary without reading every architecture document.

This decision does not add route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, runtime endpoint behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future Pitaya-aligned direction planning must start with W-0264 as a selection gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete serializer and message forwarding or distributed runtime model;
- the serializer and message forwarding inspection output creates confusion with public API compatibility;
- protocol adapter, application dispatch, transactional dispatch, protocol bridge, outbound delivery, or target-scope ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for serializer boundaries, forwarding, or backend route targeting.

## Follow-Up

- Complete `W-0264`: select the next Pitaya-aligned direction after the serializer and message forwarding map.
- Keep route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, cluster-safe session routing behavior, distributed runtime behavior, distributed groups, room broadcast fanout, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
