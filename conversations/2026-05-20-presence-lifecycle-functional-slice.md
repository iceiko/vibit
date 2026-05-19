# Conversation: Presence Lifecycle Functional Slice

Date: 2026-05-20
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0085-presence-lifecycle-functional-slice.md`
Related changes:

- `changes/2026-05-20-define-presence-lifecycle-functional-slice/`

Related artifacts:

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/platform/transport/ws/server.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/connection_lifecycle.go`
- `runtime/cmd/vibit-server/connection_binding_registry.go`
- `runtime/internal/app/connection/registry_test.go`
- `runtime/internal/platform/transport/ws/server_test.go`
- `runtime/cmd/vibit-server/connection_lifecycle_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `rules/check-rules.json`
- `tools/vibit`

## Context

The repository had completed `M-104/W-0176`, the protocol session carrier functional slice. Successful login responses could now carry the server-created runtime session id through existing envelope session metadata.

The next ready work item was `M-105/W-0177`, a Tier 2 presence lifecycle functional slice under `ADR-0082`.

## Maintainer Narrative

The maintainer asked to keep moving quickly toward Nakama/Pitaya-class capability and to reduce confirmation gates:

```text
继续推进，我们要能够保证我们的服务器能较快地达到nakama pitaya的水平，否则我们永远做不完，没法上线。
```

## Agent Response Summary

The agent advanced W-0177 as a Tier 2 functional slice. The implementation did not create another pure confirmation gate.

The first presence behavior is server-owned and registry-backed:

- WebSocket open/close can notify a credential-neutral lifecycle observer.
- PostgreSQL startup composition wires that observer to the active connection registry.
- Successful first-message connection binding records validated player identity in the registry.
- `PresenceForPlayer` derives online/offline status from active bound player connection records.

## Decisions

- Complete `M-105/W-0177`.
- Accept `ADR-0085`.
- Add `runtime.presence_lifecycle_functional_slice` as the repository check rule.
- Keep Protobuf presence routes, generated output, subscriptions, broadcasts, durable/distributed presence, chat, social modules, matchmaking, match runtime, operations/admin behavior, reconnect/resume tokens, logout-triggered close, runtime session revocation, dependencies, and direct Nakama/Pitaya API compatibility deferred.
- Move the work queue to `M-106/W-0178`, a protected presence protocol query functional slice.

## Artifacts

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
- `runtime/internal/platform/transport/ws/server.go`
- `runtime/internal/platform/transport/ws/server_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/connection_lifecycle.go`
- `runtime/cmd/vibit-server/connection_binding_registry.go`
- `runtime/cmd/vibit-server/connection_lifecycle_test.go`
- `changes/2026-05-20-define-presence-lifecycle-functional-slice/`
- `decisions/ADR-0085-presence-lifecycle-functional-slice.md`
- `conversations/2026-05-20-presence-lifecycle-functional-slice.md`
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

- Presence protocol query shape remains next.
- Runtime session and access-token-record presence linkage from access-token validation remains deferred because `RequestIdentity` does not currently carry those values.
- Presence subscriptions and broadcasts remain deferred.
- Durable/distributed presence remains deferred.
- Direct Nakama/Pitaya API compatibility remains deferred.

## Follow-Up

Continue with `M-106/W-0178 define_presence_protocol_query_functional_slice` as the next Tier 2 product-parity step, unless the maintainer redirects the queue.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, close reason text, remote addresses, headers, or GitHub tokens are recorded in this conversation log.
