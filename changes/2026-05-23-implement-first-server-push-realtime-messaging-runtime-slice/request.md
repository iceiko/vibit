# Request

## User Request

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

## Interpretation

Advance exactly one next-ready work item:

```text
W-0215 Implement first server push and realtime messaging runtime slice
```

The work must preserve Nakama/Pitaya alignment and stay inside `ADR-0122`. It may add application-owned runtime behavior under `runtime/internal/app/realtime`, but must not add WebSocket outbound delivery, Protobuf source, generated output, protocol routes, startup wiring, persistence, migrations, dependencies, authentication/session behavior changes, route-protection changes, public client publish routes, stream subscription persistence, offline inboxes, delivery guarantees, distributed fanout, broad chat/social behavior, matchmaking, match runtime, hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## Selected Scope

Implement the first application-owned runtime slice by adding:

- `runtime/internal/app/realtime/service.go`;
- `runtime/internal/app/realtime/service_test.go`;
- `ADR-0123`;
- change, verification, and conversation records;
- `.arch` manifest updates;
- `tools/vibit` rule coverage;
- next-ready `M-144/W-0216`.

## User-Visible Outcome

Maintainers and agents can inspect a small realtime application service that:

- accepts server-authored realtime message intents;
- supports `server_notice`, `domain_event_push`, `stream_message`, and `presence_signal` vocabulary;
- resolves connection and player-current-connection targets through server-observed registry state;
- rejects metadata-only and player-authored publish attempts;
- keeps stream subscription delivery future-only;
- returns redacted delivery intents for later protocol/transport slices.

## Non-Goals

- No concrete WebSocket writes.
- No `proto/vibit/realtime/v1/realtime.proto`.
- No generated Go Protobuf output.
- No protocol bridge or application bootstrap handler.
- No startup wiring.
- No persistence, migrations, dependencies, delivery guarantees, offline inboxes, acknowledgements, retries, ordering, backpressure, distributed fanout, chat routes, stream subscription state, matchmaking, match runtime, broad social modules, or direct Nakama/Pitaya compatibility.

## Nakama/Pitaya Fit

Nakama informs the capability family: notifications, streams, chat, and presence-adjacent outbound messages.

Pitaya informs the layering: acceptors, sessions, handlers, backend services, push, groups, broadcast, and later cluster/RPC concerns remain separate.

The implementation adapts those references into a vibit-native application service. It does not copy Nakama notification/channel/stream APIs or Pitaya push/group/broadcast APIs.

## Acceptance Criteria

- [x] Application-owned realtime runtime service exists under `runtime/internal/app/realtime`.
- [x] The service validates server-authored intent and target vocabulary.
- [x] Metadata-only and player-authored publish attempts are rejected before registry resolution.
- [x] Active bound connection and player-current-connection targets resolve through the connection registry.
- [x] Stream target remains future-only and unauthorized.
- [x] Delivery results are redacted and copied.
- [x] No protocol, generated output, transport delivery, persistence, dependency, authentication/session, broad social, matchmaking, match runtime, distributed runtime, or direct compatibility expansion is added.
