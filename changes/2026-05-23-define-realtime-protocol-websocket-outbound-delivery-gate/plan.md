# Plan

1. Define the realtime protocol and WebSocket outbound delivery gate standard.
2. Add the paired Simplified Chinese translation.
3. Accept `ADR-0125`.
4. Record the conversation memory.
5. Complete `M-145/W-0217` and open `M-146/W-0218`.
6. Update `.arch` manifests, README/alpha pointers, AGENTS guides, and storage module references.
7. Register and implement `runtime.realtime_protocol_websocket_outbound_delivery_gate` check coverage.
8. Verify repository checks.

## Stop Conditions

Stop before adding Protobuf source, generated output, protocol bridge behavior, WebSocket outbound writers, concrete socket writes, startup wiring, persistence, migrations, dependencies, authentication/session behavior changes, route-protection changes, stream subscriptions, chat, groups, broadcast fanout, delivery guarantees, distributed runtime, matchmaking, match runtime, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, or direct Nakama/Pitaya API compatibility.
