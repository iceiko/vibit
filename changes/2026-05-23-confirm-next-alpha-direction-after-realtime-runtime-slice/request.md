# Request

## Original Request

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。提交所用的Key在Git忽略的文件里有
```

English summary: continue the next ready work item, keep Nakama/Pitaya alignment explicit, then commit and push the result using the ignored local Git credential material.

## Clarified Requirement

Advance `M-144/W-0216 Confirm next alpha direction after realtime runtime slice`.

This is a direction-selection change only. It must select exactly one bounded alpha direction after the first server push and realtime messaging runtime slice and open one follow-up work item.

## Selected Direction

Select:

```text
define_realtime_protocol_websocket_outbound_delivery_gate
```

as the next prototype-ready alpha direction.

## User-Visible Outcome

Maintainers and agents can see that the next bounded work after the application-owned realtime runtime slice is a gate for realtime protocol and WebSocket outbound delivery planning, not immediate socket writes, generated output, chat, groups, broadcast fanout, matchmaking, match runtime, distributed runtime, or direct compatibility.

## Nakama/Pitaya Alignment

Nakama guides the product pressure: after server-authored delivery intents exist, prototype authors need a path toward client-visible outbound behavior that can later support notifications, streams, chat, and presence-adjacent shared online services.

Pitaya guides the architecture pressure: acceptor, session, handler, protocol serializer, backend intent, push, group, broadcast, and later cluster/RPC concerns must remain separated before concrete socket writes or distributed fanout.

## Non-Goals

- Implementing WebSocket outbound delivery or socket writes.
- Adding protocol bridge behavior.
- Adding Protobuf source or generated output.
- Adding application bootstrap handlers or startup wiring.
- Adding persistence, stream subscription storage, offline inboxes, acknowledgements, ordering guarantees, retries, backpressure, durable offsets, or delivery guarantees.
- Adding chat, groups, parties, rooms, broadcast fanout, matchmaking, or match runtime.
- Changing authentication/session behavior.
- Changing route protection.
- Adding repository interfaces, PostgreSQL adapters, migrations, or dependencies.
- Adding hosted deployments or demos.
- Creating release binaries, packages, containers, checksums, provenance, signing artifacts, install scripts, registry publications, or SDK packages.
- Executing public announcements beyond the GitHub release record.
- Running paid promotion.
- Adding large object/blob storage or S3-compatible object storage.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- The exact realtime protocol payload posture will be defined by the follow-up gate.
- Whether the first transport delivery implementation should target one connection, all current player connections, or another bounded target remains deferred.
- Persistence, delivery guarantees, offline inboxes, acknowledgements, ordering, backpressure, retries, distributed fanout, and client SDK behavior remain deferred.

## Acceptance Criteria

- [x] Exactly one next bounded direction is selected.
- [x] Nakama/Pitaya reference alignment is recorded.
- [x] One follow-up work item is opened.
- [x] WebSocket outbound delivery, socket writes, protocol bridge, Protobuf source, generated output, startup wiring, persistence, repository, adapter, migration, dependency, authentication/session, route-protection, hosted deployment, release artifact, public announcement, paid promotion, blob/S3, matchmaking, match runtime, distributed runtime, broad social module, and direct compatibility deferrals are preserved.
