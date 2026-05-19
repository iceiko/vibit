# Conversation: Runtime Session Validation Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-session-postgresql-adapter-implementation/`
- `changes/2026-05-17-define-runtime-session-validation-gate/`

Related artifacts:

- `docs/runtime-session-validation-gate.md`
- `docs/runtime-session-validation-gate.zh-CN.md`
- `decisions/ADR-0065-runtime-session-validation-gate.md`
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

`M-068/W-0140` implemented the PostgreSQL adapter for `runtime/internal/app/session.Repository`. The repository now has durable runtime session storage, but runtime behavior still does not validate persisted sessions or set `RequestIdentity.SessionValidated` true.

The work queue was blocked at `M-069/W-0141`, a next-direction confirmation gate.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had already asked the agent to communicate in Chinese and to use Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended:

```text
define_runtime_session_validation_gate
```

The reason was that durable session storage and the PostgreSQL adapter now exist, but the repository still needs a ratified validation boundary before any code treats a persisted session row as authenticated request identity.

Nakama guided the lifecycle pressure: sessions can expire, be revoked, and interact with logout and socket behavior. Pitaya guided the layering: session context may reach handlers through application context, while acceptors and routing remain separate from durable validation.

## Decisions

- Close `M-069/W-0141`.
- Select `define_runtime_session_validation_gate`.
- Create and complete `M-070/W-0142` as a gate-only milestone.
- Accept `ADR-0065`.
- Add `runtime.runtime_session_validation_gate` as a repository check rule.
- Preserve deferrals for runtime validation implementation, session creation composition, route-policy use of session or bound identity, WebSocket handshake authentication, transport credential carriers, Protobuf session messages, existing envelope changes, logout/revocation active-connection invalidation, reconnect/epoch behavior, cleanup jobs, dependencies, memory durable session behavior, direct Nakama/Pitaya API compatibility, and broader game backend behavior.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-session-postgresql-adapter-implementation/`
- `changes/2026-05-17-define-runtime-session-validation-gate/`
- `docs/runtime-session-validation-gate.md`
- `docs/runtime-session-validation-gate.zh-CN.md`
- `decisions/ADR-0065-runtime-session-validation-gate.md`
- `conversations/2026-05-17-runtime-session-validation-gate.md`
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

- Runtime session validation implementation remains deferred.
- Session creation composition remains deferred.
- Whether session-validated identity can satisfy ordinary protected route policy remains deferred.
- Logout/revocation active-connection behavior remains deferred.
- Reconnect/resume/duplicate replacement and connection epoch behavior remain deferred.
- Operations and observability classification for session ids remains deferred.

## Follow-Up

- Block at the next confirmation gate before choosing runtime session validation implementation or another major direction.
- The likely next conservative direction is a bounded runtime session validation implementation slice with fake repository tests.
- Other valid candidates are session creation composition, bound identity route policy, logout/revocation active-connection behavior, reconnect/epoch behavior, operations hardening, or broader Nakama/Pitaya-inspired game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, or GitHub tokens are recorded in this conversation log.
