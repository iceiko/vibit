# Impact

## Runtime

- Completes `M-142/W-0214`.
- Defines the first server push and realtime messaging gate.
- Opens `M-143/W-0215` as next-ready.
- Does not add Go runtime behavior.
- Does not add transport delivery behavior.

## Protocol And Generated Output

- No Protobuf source changes.
- No generated output changes.
- No protocol route changes.
- No envelope changes.
- Future protocol candidates are recorded only as candidates.

## Storage, Data, And Dependencies

- No storage object behavior changes.
- No repository interface changes.
- No PostgreSQL adapter changes.
- No persistence or migration changes.
- No dependencies.
- No blob/S3 object storage work.

## Authentication And Session

- No authentication behavior changes.
- No session behavior changes.
- No route-protection changes.
- No WebSocket handshake authentication changes.
- Metadata-only identity remains rejected as proof.

## Product Scope

- Defines a vibit-native outbound realtime vocabulary: `server_notice`, `domain_event_push`, `stream_message`, and `presence_signal`.
- Keeps chat implementation, stream subscriptions, offline inboxes, delivery guarantees, distributed fanout, matchmaking, match runtime, social modules, SDKs, hosted deployment, release artifact expansion, public announcements, paid promotion, and direct Nakama/Pitaya compatibility deferred.

## Nakama/Pitaya Alignment

- Nakama alignment: the gate covers the first vocabulary needed before notifications, streams, chat, or presence-adjacent outbound messages.
- Pitaya alignment: the gate preserves acceptor/session/handler/protocol/backend separation and defers groups, broadcast fanout, cluster, RPC, and direct route compatibility.
