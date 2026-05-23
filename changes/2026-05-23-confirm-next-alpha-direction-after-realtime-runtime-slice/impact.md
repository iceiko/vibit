# Impact

## Runtime

- Completes `M-144/W-0216`.
- Selects `define_realtime_protocol_websocket_outbound_delivery_gate`.
- Creates `M-145/W-0217` as the next ready work item.
- Does not add Go runtime behavior.

## Protocol, Transport, And Generated Output

- No WebSocket outbound delivery changes.
- No concrete socket writes.
- No protocol bridge changes.
- No Protobuf source changes.
- No generated output changes.
- No protocol route changes.
- No envelope changes.
- No startup wiring changes.

## Storage, Data, And Dependencies

- No storage object behavior changes.
- No repository interface changes.
- No PostgreSQL adapter changes.
- No persistence, offline inbox, stream subscription storage, or delivery guarantee changes.
- No migrations.
- No dependencies.
- No blob/S3 object storage work.

## Authentication And Session

- No authentication behavior changes.
- No session behavior changes.
- No route-protection changes.
- No WebSocket handshake authentication changes.

## Product Scope

- Chooses the realtime protocol and WebSocket outbound delivery gate as the next prototype-ready shared online-service direction.
- Keeps chat, groups, broadcast fanout, matchmaking, match runtime, distributed runtime, SDKs, hosted deployment, release artifact expansion, public announcements, paid promotion, and direct Nakama/Pitaya compatibility deferred.

## Nakama/Pitaya Alignment

- Nakama alignment: the next selected family covers bounded planning for future client-visible outbound delivery that can support notifications, streams, chat, and presence-adjacent shared online services.
- Pitaya alignment: the future gate must preserve acceptor/session/handler/protocol/backend/push separation and avoid premature group, broadcast, cluster, RPC, service discovery, or direct route compatibility.
