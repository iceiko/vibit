# Checklist

- [x] Add realtime Protobuf source.
- [x] Regenerate Go Protobuf output with Buf.
- [x] Register realtime generated payload types.
- [x] Add protocol bridge mapping for `server_notice`, `domain_event_push`, and `presence_signal`.
- [x] Keep `stream_message` future-only.
- [x] Add WebSocket outbound delivery by server-observed connection id and epoch.
- [x] Serialize writes per accepted socket.
- [x] Keep transport delivery credential-neutral and payload-policy-neutral.
- [x] Keep application realtime policy outside protocol and transport adapters.
- [x] Avoid startup wiring and realtime bootstrap handlers.
- [x] Preserve stream subscription, chat, group, broadcast fanout, persistence, delivery guarantee, distributed runtime, matchmaking, match runtime, release, and direct compatibility deferrals.
- [x] Add `ADR-0126`.
- [x] Add conversation and change records.
- [x] Update `.arch` manifests and next-ready state.
- [x] Register `runtime.realtime_protocol_websocket_outbound_delivery_implementation`.
- [x] Update docs and module guidance.
- [x] Run verification commands.
- [ ] Commit and push.
