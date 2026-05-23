# ADR-0125: Realtime Protocol And WebSocket Outbound Delivery Gate

Status: Accepted
Date: 2026-05-23
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-23-define-realtime-protocol-websocket-outbound-delivery-gate/`

Related conversations:

- `conversations/2026-05-23-realtime-protocol-websocket-outbound-delivery-gate.md`

Related artifacts:

- `docs/realtime-protocol-websocket-outbound-delivery-gate.md`
- `docs/realtime-protocol-websocket-outbound-delivery-gate.zh-CN.md`
- `runtime/internal/app/realtime/service.go`
- `runtime/internal/app/realtime/service_test.go`
- `docs/first-server-push-realtime-messaging-gate.md`
- `decisions/ADR-0124-next-alpha-direction-realtime-protocol-websocket-outbound-delivery-gate.md`
- `decisions/ADR-0123-first-server-push-realtime-messaging-runtime-slice.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0123` added the application-owned realtime runtime slice. It validates server-authored outbound intents and resolves active delivery recipients, but it deliberately does not define client-visible protocol payloads or write WebSocket frames.

`ADR-0124` selected `define_realtime_protocol_websocket_outbound_delivery_gate` as the next prototype-ready direction. The next boundary must clarify where protocol payload mapping stops and WebSocket transport delivery starts before any Protobuf source, generated output, protocol bridge, startup wiring, or socket write code appears.

Nakama provides the product pressure: a useful game backend eventually needs outbound notifications, streams, chat, and presence-adjacent client-visible messages. Pitaya provides the layering pressure: acceptors, sessions, handlers, protocol serialization, backend service intent, push, groups, broadcast, and later cluster/RPC concerns must stay separated.

## Decision

Accept `docs/realtime-protocol-websocket-outbound-delivery-gate.md` as the gate for future realtime protocol payloads and WebSocket outbound delivery.

The gate defines:

- application-owned realtime service as the policy and recipient-resolution owner;
- Protobuf protocol adapter as the future payload and envelope mapping owner;
- WebSocket transport as the future encoded-frame write owner;
- server-observed connection id and epoch as the only connection targeting authority;
- future source candidates for `proto/vibit/realtime/v1/realtime.proto`, generated Go output, realtime protocol bridge, application bootstrap, and transport outbound delivery;
- a future implementation work item, `W-0218 Implement realtime protocol and WebSocket outbound delivery slice`;
- stop conditions preserving protocol source, generated output, concrete socket writes, startup wiring, persistence, delivery guarantees, stream subscriptions, chat, groups, broadcast fanout, distributed runtime, matchmaking, match runtime, and direct compatibility deferrals in this slice.

Create `M-146/W-0218 Implement realtime protocol and WebSocket outbound delivery slice` as the next ready work item.

This decision completes `M-145/W-0217`.

## Alternatives Considered

- Implement WebSocket outbound writes directly from `runtime/internal/app/realtime`.
- Add Protobuf realtime messages and generated output without a transport handoff gate.
- Put recipient authorization into the protocol adapter.
- Put protocol payload mapping into the WebSocket transport.
- Define stream subscriptions, chat, groups, or broadcast fanout before the single-process delivery boundary is explicit.
- Copy Nakama or Pitaya public APIs directly.

## Rationale

The previous runtime slice created a stable application handoff: validated server-authored delivery intents. The next useful step is to define how a future implementation can transform those intents into vibit-native protocol payloads and encoded WebSocket frames without changing ownership.

This gate keeps the application layer responsible for authorization, the protocol adapter responsible for serialization and envelope mapping, and the WebSocket transport responsible for byte delivery to server-observed connections.

That matches Nakama-class capability pressure while preserving Pitaya-style layering. It also keeps durable delivery, fanout, stream subscriptions, chat semantics, group/broadcast behavior, cluster/RPC, matchmaking, and match runtime out of the first outbound delivery slice.

## Agent Reasoning Summary

The maintainer asked to continue. The current `next_ready` item was `W-0217`, and its ask-first boundary forbids implementing protocol source, generated output, protocol bridges, socket writes, startup wiring, persistence, delivery guarantees, distributed runtime, or direct compatibility in this slice. A gate artifact is therefore the correct continuation: it records the future file candidates, ownership rules, verification expectations, and stop conditions before implementation.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  protocol_transport_boundary_clarity: high
  implementation_scope_control: high
  generated_output_risk: none_in_this_step
  socket_write_risk: none_in_this_step
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `docs/realtime-protocol-websocket-outbound-delivery-gate.md` and its Simplified Chinese translation are accepted.
- `runtime.realtime_protocol_websocket_outbound_delivery_gate` becomes the repository check rule for this gate.
- `M-145/W-0217` is completed.
- `M-146/W-0218` becomes next-ready.
- Protobuf source, generated output, realtime protocol bridge behavior, WebSocket outbound delivery, concrete socket writes, application bootstrap handler, startup wiring, persistence, migrations, dependencies, authentication/session behavior changes, route-protection changes, stream subscription persistence, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, distributed fanout, broad social modules, matchmaking, match runtime, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, and direct Nakama/Pitaya API compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- the application-owned realtime service handoff shape changes before implementation;
- a later protocol ADR changes envelope semantics before realtime payloads are added;
- connection registry semantics move away from server-observed connection id and epoch;
- alpha feedback shows operations visibility, examples, or failure-path verification is more urgent than outbound delivery implementation;
- the maintainer explicitly selects stream subscriptions, chat, groups, matchmaking, match runtime, distributed runtime, or direct compatibility instead.

## Follow-Up

- Complete `W-0218`: implement the smallest realtime protocol and WebSocket outbound delivery slice authorized by this gate and a matching implementation ADR.
- Keep stream subscription ownership, chat semantics, groups, broadcast fanout, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, distributed fanout, matchmaking, match runtime, and direct compatibility behind later bounded work items.
