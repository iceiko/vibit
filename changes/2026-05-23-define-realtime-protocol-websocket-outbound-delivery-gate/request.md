# Request

## Summary

Advance `M-145/W-0217 Define realtime protocol and WebSocket outbound delivery gate`.

## Maintainer Intent

Continue the current next-ready work item while preserving Nakama/Pitaya alignment.

## Scope

Define the gate for future realtime protocol payloads and WebSocket outbound delivery after the application-owned realtime runtime slice.

The change should:

- accept a gate standard for realtime protocol and WebSocket outbound delivery ownership;
- record `ADR-0125`;
- update `.arch` manifests, current alpha pointers, and agent guides;
- register `runtime.realtime_protocol_websocket_outbound_delivery_gate`;
- complete `M-145/W-0217`;
- open `M-146/W-0218 Implement realtime protocol and WebSocket outbound delivery slice` as next-ready.

## Non-Goals

This change must not add:

- WebSocket outbound delivery;
- concrete socket writes;
- `proto/vibit/realtime/v1/realtime.proto`;
- generated realtime Protobuf output;
- realtime protocol bridge behavior;
- application bootstrap handlers;
- startup wiring;
- persistence, migrations, dependencies, or repository changes;
- authentication/session or route-protection behavior changes;
- stream subscription persistence;
- offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, or backpressure;
- chat, groups, broadcast fanout, matchmaking, match runtime, distributed runtime, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, or direct Nakama/Pitaya API compatibility.

## Reference Alignment

Nakama reference pressure:

- Preserve the path toward notifications, streams, chat, and presence-adjacent outbound realtime capability.
- Do not copy Nakama APIs or payloads.

Pitaya reference pressure:

- Preserve acceptor/session/handler/protocol/backend push separation.
- Keep groups, broadcast fanout, cluster, RPC, and service discovery deferred.
