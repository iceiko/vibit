# Impact

## Affected Areas

This change adds the first realtime protocol payload source, generated output, protocol bridge, and WebSocket outbound delivery handoff:

```text
proto/vibit/realtime/v1/realtime.proto
runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
runtime/internal/platform/protocol/protobuf/realtime_bridge.go
runtime/internal/platform/transport/ws/outbound.go
```

It also updates ADRs, change records, `.arch` manifests, repository checks, and continuation docs so `M-146/W-0218` is completed and `M-147/W-0219` becomes the next-ready direction-confirmation item.

## Runtime Impact

The runtime impact is intentionally narrow:

- adds vibit-native `ServerNotice`, `DomainEventPush`, and `PresenceSignal` payloads;
- maps accepted `realtime.DeliveryIntent` values into existing `vibit.protocol.v1.Envelope` bytes;
- writes already encoded binary frames to a server-observed WebSocket connection id and epoch;
- serializes per-socket writes shared by synchronous handler responses and outbound delivery;
- returns redacted transport delivery outcomes.

The implementation does not wire realtime outbound delivery into process startup, add a public route, persist messages, implement delivery guarantees, or add stream/chat/group/broadcast semantics.

## Ownership Impact

Ownership remains separated:

- `runtime/internal/app/realtime` owns server-authored intent policy and recipient resolution.
- `runtime/internal/platform/protocol/protobuf` owns Protobuf payload and envelope mapping.
- `runtime/internal/platform/transport/ws` owns encoded-frame WebSocket delivery mechanics.
- `runtime/internal/app/bootstrap` and `runtime/cmd/vibit-server` remain unchanged for realtime outbound delivery.

## Contract And Protocol Impact

The existing protocol envelope is unchanged. This slice adds only the first realtime payload family under `vibit.realtime.v1` and registers generated payload types for decoding.

The implementation deliberately omits `StreamMessage` protocol payload behavior because stream subscription ownership is still future-only.

## Data And Migration Impact

No database schema, migration, repository interface, PostgreSQL adapter, stream subscription table, offline inbox, durable offset, or persisted delivery state is added.

## Authentication And Session Impact

No authentication, token validation, session persistence, route protection, WebSocket handshake, or access-token carrier behavior changes are added.

The transport delivery request remains credential-neutral and uses only server-observed connection id, connection epoch, encoded payload bytes, and delivery timestamps.

## Nakama/Pitaya Impact

Nakama informs the capability pressure toward useful client-visible notifications, stream/chat-adjacent delivery, and presence-adjacent outbound messages.

Pitaya informs the separation between application intent, protocol serialization, transport acceptors, connection/session state, backend push, group/broadcast, and later cluster/RPC concerns.

This slice adapts those references into vibit-native payload and transport handoff shapes. It does not copy public APIs, route paths, helper names, payload conventions, group behavior, or clustering semantics.

## Compatibility Risk

Compatibility risk is moderate but bounded because new Protobuf messages and generated output are added. The existing envelope and current request loop remain unchanged. The main future risk is treating this best-effort single-process delivery handoff as a durable guarantee; docs and checks explicitly keep acknowledgements, retries, ordering, durable offsets, backpressure, offline inboxes, and persistence deferred.
