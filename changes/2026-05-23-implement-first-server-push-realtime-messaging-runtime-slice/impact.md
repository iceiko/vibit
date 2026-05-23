# Impact

## Affected Areas

This change adds the first application-owned realtime runtime package:

```text
runtime/internal/app/realtime
```

It also updates ADRs, change records, `.arch` manifests, repository checks, and continuation docs so `M-143/W-0215` is completed and `M-144/W-0216` becomes the next-ready direction-confirmation item.

## Runtime Impact

The runtime impact is intentionally narrow:

- adds `Service`, `NewService`, and `AcceptServerMessage`;
- validates server-authored outbound message intent vocabulary;
- requires validated service or admin sender authority;
- rejects metadata-only identity and validated player identity before registry access;
- resolves active bound recipients through the existing single-process connection registry;
- returns redacted delivery intents for later protocol and transport slices;
- keeps `stream_subscribers` as future vocabulary only.

The service does not write WebSocket frames, encode Protobuf messages, register routes, wire startup, persist messages, or implement delivery guarantees.

## Module Ownership

The service is application-owned under `runtime/internal/app/realtime`.

It uses the connection registry through a small local interface and does not move realtime policy into:

- WebSocket transport;
- Protobuf protocol adapters;
- storage modules;
- domain modules;
- generated output.

## Contract And Protocol Impact

No public contract, Protobuf source, generated output, or envelope shape changes are added.

Future protocol payload mapping remains deferred until a later bounded work item.

## Data And Migration Impact

No database schema, migration, repository interface, PostgreSQL adapter, storage object behavior, stream subscription table, offline inbox, or durable delivery state is added.

## Authentication And Session Impact

No authentication, token validation, session persistence, route protection, WebSocket handshake, or access-token carrier behavior changes are added.

The runtime service treats only validated service/admin identities as server-authorized senders. Metadata-only identity and validated player identity are not enough to publish server facts.

## Nakama/Pitaya Impact

Nakama informs the capability family: server push, notifications, streams, chat, and presence-adjacent outbound behavior are common game-backend needs.

Pitaya informs the layering: acceptors, sessions, handlers, backend services, push, groups, broadcast, cluster, and RPC concerns must remain separate.

This slice adapts those references into vibit-native application policy and recipient-resolution behavior. It does not copy Nakama notification/channel/stream APIs or Pitaya push/group/broadcast APIs.

## Compatibility Risk

Compatibility risk is low because the change adds no public wire protocol and no direct external API surface. The main risk is future misuse of `DeliveryIntent` as if it were transport delivery. The check rule and docs explicitly keep socket writes and protocol mapping deferred.
