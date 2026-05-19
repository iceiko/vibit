# Conversation: Reconnect Connection Epoch Functional Slice

Date: 2026-05-20
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0083-reconnect-connection-epoch-functional-slice.md`
Related changes:

- `changes/2026-05-20-define-reconnect-connection-epoch-functional-slice/`

Related artifacts:

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Context

The repository was at `M-103/W-0175`, the first Tier 2 functional slice after `ADR-0082` adopted the gate density optimization strategy.

Completed prerequisites included:

- First-message connection binding.
- Single-process active connection registry.
- Application-owned close policy.
- Protocol logout route.
- WebSocket transport close handoff by server-observed connection id and epoch.

The missing smallest lifecycle primitive was explicit application-level connection epoch progression for a server-owned connection id.

## Maintainer Narrative

The maintainer asked in Chinese to continue advancing the project and emphasized that vibit needs to reach Nakama/Pitaya-class server capability faster:

```text
继续推进，我们要能够保证我们的服务器能较快地达到nakama pitaya的水平，否则我们永远做不完，没法上线。
```

The maintainer had also asked to review and reduce confirmation gates because the old process slowed development too much. This slice applies that direction by embedding the boundary in the change spec and implementing the bounded behavior in the same work item.

## Agent Response Summary

The agent advanced `W-0175` as a Tier 2 functional slice instead of creating another pure confirmation milestone.

The implementation updated the application-owned connection registry so a newer server-observed epoch for the same connection id supersedes earlier active records. Stale or repeated epochs after a newer epoch exists fail closed with `connection_epoch_stale`. Superseded records remain available through lifecycle inspection and are excluded from active targeting lists.

## Decisions

- Complete `M-103/W-0175`.
- Accept `ADR-0083`.
- Add `runtime.reconnect_connection_epoch_functional_slice` as the repository check rule for the implementation state.
- Keep reconnect tokens, resume behavior, protocol session carriers, logout-triggered close, runtime session revocation, close code mapping, close reason text, presence, operations/admin disconnect, dependencies, and direct Nakama/Pitaya API compatibility deferred.
- Move the work queue to `M-104/W-0176`, a Tier 2 protocol session carrier functional slice.

## Artifacts

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
- `changes/2026-05-20-define-reconnect-connection-epoch-functional-slice/`
- `decisions/ADR-0083-reconnect-connection-epoch-functional-slice.md`
- `conversations/2026-05-20-reconnect-connection-epoch-functional-slice.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- Protocol session carrier behavior remains next.
- Reconnect tokens and resume tokens remain deferred.
- Durable or distributed reconnect behavior remains deferred.
- Duplicate replacement socket close remains deferred.
- Close code mapping and close reason text remain deferred.
- Logout-triggered socket close remains deferred.
- Runtime session revocation remains deferred.
- Presence lifecycle remains deferred until lifecycle closure is more complete.

## Follow-Up

Continue with `M-104/W-0176 define_protocol_session_carrier_functional_slice` as the next Tier 2 lifecycle-closure step, unless the maintainer redirects the queue.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, close reason text, remote addresses, headers, or GitHub tokens are recorded in this conversation log.
