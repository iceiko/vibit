# Request

## User Request

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

## Interpretation

Advance exactly one next-ready work item:

```text
W-0218 Implement realtime protocol and WebSocket outbound delivery slice
```

The work must preserve Nakama/Pitaya alignment and stay inside `ADR-0125`. It may add the smallest realtime Protobuf payload source, generated Go Protobuf output through Buf, protocol bridge mapping from already accepted realtime delivery intents, and WebSocket transport delivery of already encoded binary frames to a server-observed connection id and epoch.

## Selected Scope

Implement the bounded slice by adding:

- `proto/vibit/realtime/v1/realtime.proto`;
- generated `runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go`;
- `runtime/internal/platform/protocol/protobuf/realtime_bridge.go`;
- `runtime/internal/platform/protocol/protobuf/realtime_bridge_test.go`;
- `runtime/internal/platform/transport/ws/outbound.go`;
- `runtime/internal/platform/transport/ws/outbound_test.go`;
- realtime payload registry import;
- WebSocket outbound socket table registration and serialized binary writes;
- `ADR-0126`;
- change, verification, and conversation records;
- `.arch` manifest updates;
- `tools/vibit` rule coverage;
- next-ready `M-147/W-0219`.

## User-Visible Outcome

Maintainers and agents can inspect a narrow outbound realtime delivery path:

- application realtime service still owns policy and recipient resolution;
- Protobuf adapter maps accepted delivery intents to vibit-native realtime payloads inside the existing envelope;
- WebSocket transport can deliver already encoded binary frames to a server-observed connection id and epoch;
- redacted outcomes distinguish delivered, missing socket, stale epoch, closed socket, and write failure.

## Non-Goals

- No startup wiring for realtime outbound delivery.
- No `runtime/internal/app/bootstrap/realtime.go`.
- No public client publish route.
- No stream subscriptions, chat rooms, groups, parties, rooms, matches, broadcast fanout, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, or backpressure.
- No persistence, migrations, dependencies, authentication/session behavior changes, route-protection changes, matchmaking, match runtime, distributed runtime, SDKs, hosted deployments, release artifacts, public announcements, paid promotion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## Nakama/Pitaya Fit

Nakama informs the capability pressure: client-visible outbound notifications, streams, chat, and presence-adjacent features need a server-to-client delivery foundation.

Pitaya informs the layering: application policy, protocol serialization, acceptor/transport mechanics, connection/session state, backend push vocabulary, groups, broadcast, RPC, and cluster concerns must stay separated.

This implementation adapts those references into vibit-native payload and transport handoff shapes. It does not copy Nakama notification/channel/stream APIs or Pitaya push/group/broadcast APIs.

## Acceptance Criteria

- [x] Realtime Protobuf source and generated output exist and trace to ADR-0125/ADR-0126.
- [x] Protocol bridge maps accepted realtime delivery intents to envelope payload bytes.
- [x] WebSocket outbound delivery writes binary frames by server-observed connection id and epoch.
- [x] Transport delivery remains credential-neutral and payload-policy-neutral.
- [x] Application policy remains outside protocol and transport adapters.
- [x] Startup wiring remains deferred.
- [x] Stream subscriptions, chat, groups, broadcast fanout, delivery guarantees, persistence, distributed runtime, matchmaking, match runtime, and direct compatibility remain deferred.
