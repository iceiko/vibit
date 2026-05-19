# Conversation: Session Repository Boundary

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-runtime-sessions-migration-source/`
- `changes/2026-05-17-define-session-repository-boundary/`

Related artifacts:

- `docs/session-repository-boundary.md`
- `docs/session-repository-boundary.zh-CN.md`
- `decisions/ADR-0061-session-repository-boundary.md`
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

`M-060/W-0132` added `runtime/migrations/postgres/000005_create_runtime_sessions.sql` as a migration-source-only slice. It did not add repository interfaces, PostgreSQL adapters, runtime session validation, route policy, WebSocket handshake authentication, logout/revocation active-connection behavior, reconnect behavior, or direct Nakama/Pitaya API compatibility.

The work queue was blocked at `M-061/W-0133`, asking for the next direction after the migration source.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya
```

The maintainer had already asked the agent to communicate in Chinese and to treat Nakama and Pitaya as important references for current game-server needs.

## Agent Response Summary

The agent recommended continuing with:

```text
define_session_repository_boundary
```

The reason was that `runtime_sessions` now exists as a SQL source, but adding runtime behavior or a PostgreSQL adapter before the repository boundary would risk mixing session lifecycle storage, token validation, WebSocket connection state, and route policy.

The agent defined a gate-only standard for the future storage-neutral session repository boundary. The future repository owner candidate is `runtime/internal/app/session`; the future PostgreSQL adapter owner remains `runtime/internal/platform/persistence/postgres`.

## Decisions

- Close `M-061/W-0133`.
- Select `define_session_repository_boundary`.
- Create and complete `M-062/W-0134` as a gate-only milestone.
- Accept `ADR-0061`.
- Add `runtime.session_repository_boundary` as a repository check rule.
- Preserve deferrals for Go repository code, PostgreSQL adapters, runtime validation, session creation, logout/revocation execution, route-policy session identity, WebSocket handshake authentication, reconnect/epoch behavior, dependencies, memory durable session behavior, and direct Nakama/Pitaya API compatibility.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: session lifecycle lookup, expiration, refresh, logout, and operational management become central in production game backends.

Pitaya guided the layering: session-like context can be useful to handlers, but acceptors, routing, and transport should not own durable session persistence.

vibit adapted both into its own boundary: a storage-neutral application repository plan with no public API compatibility promise.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-runtime-sessions-migration-source/`
- `changes/2026-05-17-define-session-repository-boundary/`
- `docs/session-repository-boundary.md`
- `docs/session-repository-boundary.zh-CN.md`
- `decisions/ADR-0061-session-repository-boundary.md`
- `conversations/2026-05-17-session-repository-boundary.md`
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

- The concrete Go repository interface remains deferred.
- The PostgreSQL adapter query shape remains deferred.
- Whether session creation belongs to login, BindConnection, or a separate command remains deferred.
- Whether persisted session identity can satisfy ordinary protected route policy remains deferred.
- Whether logout/revocation should actively invalidate open WebSocket connections remains deferred.
- Reconnect/resume/duplicate replacement and durable connection epoch policy remain deferred.

## Follow-Up

- Block at `M-063/W-0135` before choosing the next direction.
- Likely next candidates are the session PostgreSQL adapter gate, repository implementation slice, runtime session validation gate, logout/revocation active-connection gate, reconnect/epoch gate, bound identity route-policy gate, operations hardening, or broader Nakama/Pitaya-inspired game backend planning.
- Run repository checks and record any remaining warning explicitly.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, or GitHub tokens are recorded in this conversation log.
