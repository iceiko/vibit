# Conversation: PostgreSQL Session Persistence Schema Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-implementation/`
- `changes/2026-05-17-define-postgres-session-persistence-schema-gate/`

Related artifacts:

- `docs/postgres-session-persistence-schema-gate.md`
- `docs/postgres-session-persistence-schema-gate.zh-CN.md`
- `decisions/ADR-0059-postgres-session-persistence-schema-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-056/W-0128` implemented first-message `runtime.authentication.BindConnection`. That implementation intentionally preserved deferrals for durable session persistence, route-policy use of bound identity, logout/revocation active-connection behavior, reconnect/epoch policy, migrations, repositories, dependencies, memory durable authentication behavior, and direct Nakama/Pitaya API compatibility.

The work queue was blocked at `M-057/W-0129`, asking for the next direction after first-message connection binding implementation.

## Maintainer Narrative

The maintainer said:

```text
按照你的建议继续推进。
```

The maintainer had already asked the agent to continue in Chinese and to keep Nakama and Pitaya as important references for game-server needs.

## Agent Response Summary

The agent selected the previously recommended direction:

```text
define_postgres_session_persistence_schema_gate
```

The agent defined a gate-only PostgreSQL session persistence schema boundary. The gate records PostgreSQL as the first durable target, `runtime_sessions` as the future logical table candidate, and `runtime/migrations/postgres/000005_create_runtime_sessions.sql` as the future migration source candidate.

The agent did not add SQL migration source, repository interfaces, PostgreSQL adapters, runtime session validation behavior, route-policy bound identity, WebSocket handshake authentication, logout/revocation active-connection behavior, reconnect/epoch behavior, dependencies, or direct Nakama/Pitaya API compatibility.

## Decisions

- Close `M-057/W-0129`.
- Select `define_postgres_session_persistence_schema_gate`.
- Create and complete `M-058/W-0130` as a gate-only milestone.
- Accept `ADR-0059`.
- Add `runtime.postgres_session_persistence_schema_gate` as a repository check rule.
- Select PostgreSQL as the first durable session store target.
- Record `runtime_sessions` as the future logical table candidate.
- Defer actual migration source and all runtime session behavior.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-implementation/`
- `changes/2026-05-17-define-postgres-session-persistence-schema-gate/`
- `docs/postgres-session-persistence-schema-gate.md`
- `docs/postgres-session-persistence-schema-gate.zh-CN.md`
- `decisions/ADR-0059-postgres-session-persistence-schema-gate.md`
- `conversations/2026-05-17-postgres-session-persistence-schema-gate.md`
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

- Exact SQL constraints and indexes for `runtime_sessions` remain deferred.
- Whether session creation belongs to login, BindConnection, or a separate route remains deferred.
- Whether session validation can satisfy ordinary route policy remains deferred.
- Whether logout/revocation should actively invalidate open WebSocket connections remains deferred.
- Reconnect/resume/duplicate replacement and durable connection epoch policy remain deferred.

## Follow-Up

- Block at the next direction confirmation gate before adding `runtime_sessions` migration source or choosing another major direction.
- Future candidates include implementing the session migration source, defining the session repository boundary, defining logout/revocation active-connection behavior, defining reconnect/epoch behavior, or continuing broader Nakama/Pitaya-inspired backend planning.
- Run repository checks and record any remaining warning explicitly.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, or GitHub tokens are recorded in this conversation log.
