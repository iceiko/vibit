# ADR-0126: Realtime Protocol And WebSocket Outbound Delivery Implementation

Status: Accepted
Date: 2026-05-23
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-23-implement-realtime-protocol-websocket-outbound-delivery-slice/`

Related conversations:

- `conversations/2026-05-23-realtime-protocol-websocket-outbound-delivery-implementation.md`

Related artifacts:

- `proto/vibit/realtime/v1/realtime.proto`
- `runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go`
- `runtime/internal/platform/protocol/protobuf/realtime_bridge.go`
- `runtime/internal/platform/protocol/protobuf/realtime_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/payload_registry.go`
- `runtime/internal/platform/transport/ws/outbound.go`
- `runtime/internal/platform/transport/ws/outbound_test.go`
- `runtime/internal/platform/transport/ws/server.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0125` defined the realtime protocol and WebSocket outbound delivery gate. It authorized the next bounded implementation slice to add a small realtime Protobuf payload family, generated output, protocol bridge mapping, and WebSocket encoded-frame delivery while keeping application realtime policy, protocol serialization, and transport mechanics separate.

The runtime already has an application-owned realtime service under `runtime/internal/app/realtime`. That service validates server-authored outbound intents and resolves active recipients into `DeliveryIntent` values, but it deliberately does not encode protocol payloads or write sockets.

The runtime also has a WebSocket transport that accepts binary frames and carries server-observed connection id and epoch metadata. The missing piece is a narrow handoff that can turn already accepted delivery intents into encoded frames and deliver those frames to the corresponding server-observed socket.

Nakama provides the product pressure: useful game backends need client-visible outbound notices, stream/chat-adjacent messages, and presence-adjacent signals. Pitaya provides the layering pressure: acceptor/session/handler/protocol/backend push concerns must remain separated. vibit adapts both references without public API compatibility.

## Decision

Implement the realtime protocol and WebSocket outbound delivery slice authorized by `ADR-0125`.

The implementation adds:

- `vibit.realtime.v1.ServerNotice`;
- `vibit.realtime.v1.DomainEventPush`;
- `vibit.realtime.v1.PresenceSignal`;
- generated Go Protobuf output through Buf;
- payload registry import for realtime generated messages;
- `BuildRealtimeDeliveryEnvelope`;
- `BuildRealtimeDeliveryFrame`;
- route mapping for `realtime.ServerNotice`, `realtime.DomainEventPush`, and `realtime.PresenceSignal`;
- `stream_message` rejection as future-only;
- WebSocket `DeliverBinaryFrame`;
- per-socket serialized writes;
- redacted outbound outcomes `delivered`, `socket_not_found`, `epoch_mismatch`, `already_closed`, and `write_failed`;
- reuse of the same per-socket writer path for synchronous handler responses;
- focused tests for protocol mapping, copied payloads, unsupported intents, credential-neutral frame shape, stale epochs, closed sockets, write failure redaction, write serialization, and accepted-socket delivery.

This ADR does not add startup wiring, `runtime/internal/app/bootstrap/realtime.go`, public client publish routes, protocol routes, stream subscriptions, chat rooms, groups, broadcast fanout, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, persistence, migrations, dependencies, authentication/session behavior changes, route-protection changes, matchmaking, match runtime, distributed runtime, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Wire realtime outbound delivery into process startup in the same slice.
- Add stream, chat, group, or broadcast payloads immediately.
- Add durable delivery guarantees before any transport proof.
- Expose client-originated publish or subscribe routes.
- Put recipient authorization in the WebSocket transport.
- Copy Nakama notification/channel/stream or Pitaya push/group/broadcast APIs directly.

## Rationale

The smallest useful implementation is a handoff, not a product feature. Application code already owns authorization and recipient resolution. The protocol adapter can safely map accepted intents into vibit-native Protobuf payloads and existing envelope metadata. The transport can safely write already encoded binary frames by server-observed connection id and epoch.

Keeping startup wiring out of this slice lets the protocol and transport mechanics be verified without hiding product behavior in process assembly. Keeping `stream_message` future-only avoids pretending subscription ownership exists. Avoiding delivery guarantees keeps this slice best-effort and single-process until persistence, acknowledgements, retries, ordering, durable offsets, and backpressure are separately designed.

This moves vibit toward Nakama-class outbound realtime usefulness while preserving Pitaya-style separation between application intent, protocol serialization, transport mechanics, connection state, group/broadcast vocabulary, and later cluster concerns.

## Agent Reasoning Summary

The maintainer asked to continue and keep Nakama/Pitaya alignment explicit. The current next-ready work item was W-0218, and ADR-0125 authorized a bounded implementation slice. The safe continuation was to add the narrow protocol and WebSocket handoff with focused tests and checks, then open a next-direction confirmation item before broader realtime, social, delivery guarantee, or distributed work.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  transport_protocol_app_separation: high
  implementation_scope_control: high
  generated_output_traceability: high
  delivery_guarantee_risk: deferred
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `proto/vibit/realtime/v1/realtime.proto` exists.
- `runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go` exists and is generated from the Protobuf source.
- `runtime/internal/platform/protocol/protobuf/realtime_bridge.go` exists.
- `runtime/internal/platform/transport/ws/outbound.go` exists.
- `runtime.realtime_protocol_websocket_outbound_delivery_implementation` becomes the repository check rule for this slice.
- `M-146/W-0218` is completed.
- The next bounded direction is `W-0219 Confirm next alpha direction after realtime outbound delivery slice`.
- Startup wiring, delivery guarantees, stream subscriptions, chat, groups, broadcast fanout, persistence, distributed runtime, matchmaking, match runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- the existing envelope can no longer carry server-to-client realtime payloads safely;
- a later subscription or delivery guarantee ADR requires different frame metadata;
- connection epoch semantics change away from server-observed targeting;
- outbound delivery needs durable or distributed routing before local prototype users can evaluate it;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal through a later ADR.

## Follow-Up

- Complete `W-0219`: confirm the next alpha direction after realtime outbound delivery.
- Keep stream subscriptions, chat semantics, groups, broadcast fanout, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, persistence, distributed fanout, matchmaking, match runtime, SDKs, hosted deployment, release artifacts, and direct compatibility behind later bounded work items.
