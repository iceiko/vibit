# Conversation: Realtime Protocol And WebSocket Outbound Delivery Implementation

Date: 2026-05-23
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-23-implement-realtime-protocol-websocket-outbound-delivery-slice/`

Related artifacts:

- `proto/vibit/realtime/v1/realtime.proto`
- `runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go`
- `runtime/internal/platform/protocol/protobuf/realtime_bridge.go`
- `runtime/internal/platform/protocol/protobuf/realtime_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/payload_registry.go`
- `runtime/internal/platform/transport/ws/outbound.go`
- `runtime/internal/platform/transport/ws/outbound_test.go`
- `runtime/internal/platform/transport/ws/server.go`
- `decisions/ADR-0126-realtime-protocol-websocket-outbound-delivery-implementation.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-145/W-0217` defined the realtime protocol and WebSocket outbound delivery gate. The current next-ready item was `M-146/W-0218 Implement realtime protocol and WebSocket outbound delivery slice`.

The requested continuation required keeping Nakama/Pitaya alignment explicit and then committing and pushing the result.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

English summary: continue the work, keep Nakama/Pitaya alignment explicit, then commit and push the result.

## Agent Response Summary

The agent advanced one bounded runtime work item:

```text
W-0218 Implement realtime protocol and WebSocket outbound delivery slice
```

The change adds vibit-native realtime Protobuf payloads, generated Go Protobuf output, a protocol bridge from accepted `realtime.DeliveryIntent` values to encoded envelope bytes, and WebSocket transport delivery of already encoded binary frames to a server-observed connection id and epoch. It keeps application realtime policy outside protocol and transport adapters.

## Decisions

- Complete `M-146/W-0218`.
- Accept `ADR-0126`.
- Add `ServerNotice`, `DomainEventPush`, and `PresenceSignal` realtime payloads.
- Keep `stream_message` future-only until subscription ownership is separately defined.
- Add best-effort single-process WebSocket outbound binary frame delivery.
- Register `runtime.realtime_protocol_websocket_outbound_delivery_implementation`.
- Open `M-147/W-0219 Confirm next alpha direction after realtime outbound delivery slice` as next-ready.
- Keep startup wiring, stream subscriptions, chat, groups, broadcast fanout, persistence, delivery guarantees, matchmaking, match runtime, distributed runtime, and direct compatibility deferred.

## Nakama And Pitaya Reference Basis

Nakama guided the capability family: outbound client-visible notices, stream/chat-adjacent payload delivery, and presence-adjacent signals are common game-backend needs.

Pitaya guided the layering: application intent, protocol serialization, acceptor/transport mechanics, session and connection state, backend push, groups, broadcast, RPC, and cluster concerns must stay separated.

This slice adapts those references by adding only vibit-native protocol payloads and transport handoff mechanics. Direct Nakama/Pitaya API compatibility remains deferred.

## Artifacts

- `proto/vibit/realtime/v1/realtime.proto`
- `runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go`
- `runtime/internal/platform/protocol/protobuf/realtime_bridge.go`
- `runtime/internal/platform/protocol/protobuf/realtime_bridge_test.go`
- `runtime/internal/platform/transport/ws/outbound.go`
- `runtime/internal/platform/transport/ws/outbound_test.go`
- `decisions/ADR-0126-realtime-protocol-websocket-outbound-delivery-implementation.md`
- `changes/2026-05-23-implement-realtime-protocol-websocket-outbound-delivery-slice/`
- `conversations/2026-05-23-realtime-protocol-websocket-outbound-delivery-implementation.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The next direction after W-0218 must be confirmed before adding startup wiring, stream subscriptions, chat semantics, group membership, broadcast fanout, delivery guarantees, operations inspection, SDK surfaces, or distributed runtime work.
- Durable delivery semantics remain undefined; this slice is best-effort and single-process.
- Public client publish, subscribe, chat, stream, room, group, match, and broadcast routes remain future work.

## Follow-Up

- Complete `W-0219`: confirm the next alpha direction after realtime outbound delivery.
- Keep direct Nakama/Pitaya API compatibility, broad chat/social modules, matchmaking, match runtime, persistence, delivery guarantees, startup wiring, SDKs, hosted deployment, release artifacts, and distributed fanout behind later bounded work items unless explicitly authorized.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, raw storage object values, raw realtime payload from a real user, or concrete transport metadata from a real user are recorded in this conversation log.
