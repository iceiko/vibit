# Request

## User Request

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

## Interpretation

Advance exactly one next-ready work item:

```text
W-0214 Define first server push and realtime messaging gate
```

The work must preserve Nakama/Pitaya alignment while staying inside the gate boundary. It must not implement server push, broadcast, streams, chat, notifications, realtime messaging, protocol messages or routes, Protobuf source, generated output, storage object behavior, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior, route protection changes, hosted deployments, release artifacts, public announcements beyond the GitHub release record, paid promotion, direct Nakama/Pitaya API compatibility, large object/blob storage, S3-compatible object storage, matchmaking, match runtime, distributed runtime, RPC, service discovery, or broad product expansion.

## Selected Scope

Define the first server push and realtime messaging gate by adding:

- an English authoritative gate standard;
- a Simplified Chinese translation;
- `ADR-0122`;
- change, verification, and conversation records;
- `.arch` manifest updates;
- `tools/vibit` rule coverage;
- next-ready `M-143/W-0215`.

## Nakama/Pitaya Fit

Nakama informs the capability family: outbound realtime behavior for notifications, streams, chat, and presence-adjacent messages.

Pitaya informs the layering: acceptors, sessions, handlers, push, groups, broadcast, backend services, and later cluster/RPC concerns remain separate.

The change adapts those references into vibit-native vocabulary and does not copy direct public APIs.
