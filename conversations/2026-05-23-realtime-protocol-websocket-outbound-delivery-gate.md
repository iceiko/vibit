# Conversation: Realtime Protocol And WebSocket Outbound Delivery Gate

Date: 2026-05-23
Participants: Maintainer, Agent
Work item: `M-145/W-0217 Define realtime protocol and WebSocket outbound delivery gate`
Decision: `ADR-0125`

## Context

`M-143/W-0215` implemented the first application-owned realtime runtime slice. `M-144/W-0216` then selected `define_realtime_protocol_websocket_outbound_delivery_gate` and opened `M-145/W-0217` as next-ready.

The next step is gate-only. It must not implement:

- WebSocket outbound delivery;
- concrete socket writes;
- `proto/vibit/realtime/v1/realtime.proto`;
- generated realtime Protobuf output;
- realtime protocol bridge behavior;
- application bootstrap handlers;
- startup wiring;
- persistence, migrations, dependencies, or repository changes;
- authentication/session or route-protection behavior changes;
- stream subscriptions, chat, groups, broadcast fanout, delivery guarantees, distributed fanout, matchmaking, match runtime, or direct compatibility.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。
```

English summary: continue advancing the repository from the current next-ready work item.

## Agent Response Summary

The agent advanced one gate-only work item:

```text
W-0217 Define realtime protocol and WebSocket outbound delivery gate
```

The change defines the future handoff from application-owned realtime delivery intents to vibit-native realtime protocol payloads and WebSocket outbound delivery, without adding protocol source, generated output, socket writes, startup wiring, persistence, delivery guarantees, or direct compatibility.

## Decisions

The change accepts `docs/realtime-protocol-websocket-outbound-delivery-gate.md` and `ADR-0125`.

The gate records:

- application realtime service ownership for policy and recipient resolution;
- Protobuf protocol adapter ownership for future payload and envelope mapping;
- WebSocket transport ownership for future encoded binary frame writes;
- server-observed connection id and epoch authority;
- future file candidates for realtime Protobuf source, generated output, protocol bridge, bootstrap handler, and transport outbound delivery;
- `M-146/W-0218 Implement realtime protocol and WebSocket outbound delivery slice` as the next-ready implementation work item.

## Nakama/Pitaya Alignment

Nakama alignment:

- The gate moves vibit toward client-visible outbound realtime capability for notifications, streams, chat, and presence-adjacent behavior.
- It does not copy Nakama public APIs, route paths, runtime helper names, payload names, or compatibility promises.

Pitaya alignment:

- The gate preserves acceptor/transport, session/connection state, handler, protocol serialization, backend service intent, push/group/broadcast vocabulary, and later cluster/RPC separation.
- Groups, broadcast fanout, RPC, service discovery, frontend/backend split, and cluster behavior remain deferred.

Direct Nakama/Pitaya API compatibility remains deferred.

## Artifacts

- `docs/realtime-protocol-websocket-outbound-delivery-gate.md`
- `docs/realtime-protocol-websocket-outbound-delivery-gate.zh-CN.md`
- `decisions/ADR-0125-realtime-protocol-websocket-outbound-delivery-gate.md`
- `changes/2026-05-23-define-realtime-protocol-websocket-outbound-delivery-gate/`
- `conversations/2026-05-23-realtime-protocol-websocket-outbound-delivery-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- `W-0218` must choose the smallest implementation surface that satisfies ADR-0125 without broadening into stream subscriptions, chat, groups, broadcast fanout, persistence, or delivery guarantees.
- Future generated output must continue to trace to `.proto` source and must not be hand-edited.
- Transport delivery error reporting must stay redacted when socket writes are later added.

## Follow-Up

- Advance `W-0218` only within the boundaries recorded by `ADR-0125`.
- Keep stream subscription ownership, chat semantics, groups, broadcast fanout, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, distributed fanout, matchmaking, match runtime, and direct compatibility behind later bounded work items.

## Redaction Notes

No secrets, GitHub tokens, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, DSNs with credentials, raw storage object values from a real user, raw realtime payload from a real user, or concrete transport metadata from a real user are recorded in this conversation log.
