# Request

Advance `W-0178`: define and implement the first bounded protected presence protocol query functional slice.

The maintainer asked to continue. Per `docs/workflow.md`, continuation advances one `next_ready` work item unless blocked. `.arch/work-items.yaml` identifies `W-0178` as the current next-ready work item.

## Scope

- Embed the presence protocol query boundary in this change spec.
- Add the smallest protected query over the existing server-owned registry-backed player presence snapshot.
- Keep presence subscriptions, broadcasts, chat, friends, groups, parties, matchmaking, match runtime, cluster, SDK, operations/admin behavior, dependencies, direct Nakama/Pitaya API compatibility, reconnect/resume tokens, logout-triggered socket close, runtime session revocation, durable/distributed presence, and broad product behavior deferred.
