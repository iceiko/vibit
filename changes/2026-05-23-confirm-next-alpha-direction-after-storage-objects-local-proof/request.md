# Request

## Original Request

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

English summary: continue the next ready work item, keep Nakama/Pitaya alignment explicit, then commit and push the result.

## Clarified Requirement

Advance `M-141/W-0213 Confirm next alpha direction after storage objects local proof`.

This is a direction-selection change only. It must select exactly one bounded alpha direction after the storage objects local proof and open one follow-up work item.

## Selected Direction

Select:

```text
define_first_server_push_realtime_messaging_gate
```

as the next prototype-ready alpha direction.

## User-Visible Outcome

Maintainers and agents can see that the next bounded work after storage objects is the first server push / realtime messaging gate, not matchmaking, match runtime, distributed runtime, direct compatibility, or an unbounded social module expansion.

## Nakama/Pitaya Alignment

Nakama guides the product pressure: after durable player-owned storage objects, prototype authors need a first outbound realtime vocabulary that can later support notifications, streams, chat, and presence-adjacent shared online services.

Pitaya guides the architecture pressure: push, broadcast, stream, group, and handler vocabulary must preserve transport/session/protocol-adapter/application/backend separation before any later cluster/RPC topology.

## Non-Goals

- Implementing server push, broadcast, streams, chat, notifications, or realtime messaging.
- Adding protocol messages or routes.
- Changing Protobuf source or generated output.
- Changing storage object behavior.
- Changing repository interfaces or PostgreSQL adapters.
- Adding migrations.
- Adding dependencies.
- Changing authentication/session behavior.
- Changing route protection.
- Adding hosted deployments or demos.
- Creating release binaries, packages, containers, checksums, provenance, signing artifacts, install scripts, registry publications, or SDK packages.
- Executing public announcements beyond the GitHub release record.
- Running paid promotion.
- Adding large object/blob storage or S3-compatible object storage.
- Adding direct Nakama/Pitaya API compatibility.
- Jumping to matchmaking, match runtime, distributed runtime, or broad social/competitive modules.

## Unknowns

- The exact first server-push/realtime vocabulary will be defined by the follow-up gate.
- Whether the first implementation slice should be server-to-one-player push, stream subscription, room/group broadcast, or notification-style delivery remains deferred.
- Persistence, delivery guarantees, offline inboxes, acknowledgements, ordering, backpressure, retries, and distributed fanout remain deferred.

## Acceptance Criteria

- [x] Exactly one next bounded direction is selected.
- [x] Nakama/Pitaya reference alignment is recorded.
- [x] One follow-up work item is opened.
- [x] Implementation, protocol, generated-output, storage, authentication/session, dependency, migration, hosted deployment, release artifact, public announcement, paid promotion, blob/S3, matchmaking, match runtime, distributed runtime, and direct compatibility deferrals are preserved.
