# Plan

1. Add vibit-native realtime Protobuf payload source.
2. Regenerate Go Protobuf output through Buf.
3. Add protocol bridge mapping from accepted realtime delivery intents to envelopes and encoded frame payloads.
4. Add WebSocket outbound socket table and binary frame delivery by server-observed connection id and epoch.
5. Add focused protocol bridge and WebSocket outbound tests.
6. Record `ADR-0126`, change files, and conversation memory.
7. Mark `M-146/W-0218` completed and open `M-147/W-0219` as next-ready.
8. Register and implement `runtime.realtime_protocol_websocket_outbound_delivery_implementation` check coverage.
9. Run Buf, Go, vibit, schema, and diff verification.

This is an implementation slice only for ADR-0125. Do not wire startup, add realtime bootstrap handlers, implement stream subscriptions, chat rooms, groups, broadcast fanout, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, persistence, migrations, dependencies, authentication/session changes, route-protection changes, matchmaking, match runtime, distributed runtime, hosted/release surfaces, or direct Nakama/Pitaya API compatibility.
