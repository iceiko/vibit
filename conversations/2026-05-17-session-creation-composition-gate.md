# Conversation: Session Creation Composition Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-runtime-session-validation-implementation/`
- `changes/2026-05-17-define-session-creation-composition-gate/`

Related artifacts:

- `docs/session-creation-composition-gate.md`
- `docs/session-creation-composition-gate.zh-CN.md`
- `decisions/ADR-0067-session-creation-composition-gate.md`
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

`M-072/W-0144` implemented the application-owned persistent runtime session validator. The repository already had durable runtime session storage, a storage-neutral session repository interface, a PostgreSQL adapter, request-level access-token validation, and first-message connection binding.

The remaining gap was that no production path created durable runtime session rows. Login can issue an access token, and validation can validate a persisted session, but creation composition was not defined.

The work queue was blocked at `M-073/W-0145`, a next-direction confirmation gate.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese and to use Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended:

```text
define_session_creation_composition_gate
```

The reason was that route policy should not start depending on session validation until vibit defines how durable sessions are created, linked to token records, and committed with login.

Nakama guided the lifecycle pressure: login/session issuance has expiration, refresh, logout, and management implications, and active socket disconnect behavior should stay separate. Pitaya guided the layering: acceptors, session context, handler routing, and durable session state changes should remain separate responsibilities.

## Decisions

- Close `M-073/W-0145`.
- Select `define_session_creation_composition_gate`.
- Complete `M-074/W-0146` as a gate-only standard slice.
- Accept `ADR-0067`.
- Add `runtime.session_creation_composition_gate` as a repository check rule.
- Preserve deferrals for session creation implementation, authentication service behavior changes, session id generation, route-policy session identity, WebSocket handshake authentication, transport credential carriers, Protobuf session messages, existing envelope changes, logout/revocation active-connection invalidation, reconnect/epoch behavior, cleanup jobs, dependencies, memory durable session behavior, direct Nakama/Pitaya API compatibility, and broader game backend behavior.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-runtime-session-validation-implementation/`
- `changes/2026-05-17-define-session-creation-composition-gate/`
- `docs/session-creation-composition-gate.md`
- `docs/session-creation-composition-gate.zh-CN.md`
- `decisions/ADR-0067-session-creation-composition-gate.md`
- `conversations/2026-05-17-session-creation-composition-gate.md`
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

- Session creation implementation remains deferred.
- Session id generation posture remains deferred.
- Client-visible session carrier and login response shape remain deferred.
- Route-policy use of session-validated identity remains deferred.
- Logout/revocation active-connection behavior remains deferred.
- Reconnect/resume/duplicate replacement and connection epoch behavior remain deferred.
- Operations and observability classification for session ids remains deferred.

## Follow-Up

- Block at the next confirmation gate before choosing session creation implementation, bound/session identity route policy, logout/revocation active-connection behavior, reconnect/epoch behavior, operations hardening, or broader Nakama/Pitaya-inspired game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, or GitHub tokens are recorded in this conversation log.
