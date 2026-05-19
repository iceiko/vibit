# Conversation: Session PostgreSQL Adapter Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-session-repository-interface/`
- `changes/2026-05-17-define-session-postgresql-adapter-gate/`

Related artifacts:

- `docs/session-postgresql-adapter-gate.md`
- `docs/session-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0063-session-postgresql-adapter-gate.md`
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

`M-064/W-0136` implemented the storage-neutral session repository interface after the `runtime_sessions` migration source and session repository boundary. The repository still had no PostgreSQL adapter, no unit-of-work factory wiring, and no runtime session validation behavior.

The work queue was blocked at `M-065/W-0137`, a next-direction confirmation gate.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had already asked the agent to communicate in Chinese and to treat Nakama and Pitaya as important references for current game-server needs.

## Agent Response Summary

The agent recommended continuing with:

```text
define_session_postgresql_adapter_gate
```

The reason was that Nakama-style durable session lifecycle pressure now needs adapter-ready storage boundaries, while Pitaya-style session context separation means transport and routing must not own durable persistence.

The agent defined a gate-only standard for a future PostgreSQL adapter. The gate documents ownership, future files, repository method mapping, SQL shape, transaction handoff, error mapping, redaction, tests, Nakama/Pitaya reference mapping, and deferrals.

## Decisions

- Close `M-065/W-0137`.
- Select `define_session_postgresql_adapter_gate`.
- Create and complete `M-066/W-0138` as a gate-only milestone.
- Accept `ADR-0063`.
- Add `runtime.session_postgresql_adapter_gate` as a repository check rule.
- Preserve deferrals for adapter implementation, unit-of-work factory wiring, runtime session creation and validation, logout/revocation execution, route-policy session identity, WebSocket handshake authentication, transport credential carriers, Protobuf session messages, generated output, reconnect/epoch behavior, dependencies, memory durable session behavior, and direct Nakama/Pitaya API compatibility.

## Nakama And Pitaya Reference Basis

Nakama guided the adapter capability pressure: authenticated sessions need durable lifecycle lookup, expiration, revocation, logout, listing, and future management operations.

Pitaya guided the layering: session/context handoff should stay separate from acceptors, routing, and transport-owned state.

vibit adapted both into a PostgreSQL adapter gate that keeps persistence behind `runtime/internal/platform/persistence/postgres` and leaves validation and identity handoff in application-owned future gates.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-session-repository-interface/`
- `changes/2026-05-17-define-session-postgresql-adapter-gate/`
- `docs/session-postgresql-adapter-gate.md`
- `docs/session-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0063-session-postgresql-adapter-gate.md`
- `conversations/2026-05-17-session-postgresql-adapter-gate.md`
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

- The concrete PostgreSQL session adapter implementation remains deferred.
- Runtime session validation behavior remains deferred.
- Whether session creation belongs to login, BindConnection, or a later command remains deferred.
- Whether persisted session identity can satisfy ordinary protected route policy remains deferred.
- Whether logout/revocation should actively invalidate open WebSocket connections remains deferred.
- Reconnect/resume/duplicate replacement and durable connection epoch policy remain deferred.

## Follow-Up

- Block at `M-067/W-0139` before choosing the next direction.
- The likely next conservative direction is `implement_session_postgresql_adapter` if durable validation remains the target.
- Other valid candidates are runtime session validation gate, logout/revocation active-connection gate, reconnect/epoch gate, bound identity route-policy gate, operations hardening, or broader Nakama/Pitaya-inspired game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, or GitHub tokens are recorded in this conversation log.
