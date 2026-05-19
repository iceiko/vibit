# Conversation: Logout Revocation Active Connection Gate

Date: 2026-05-18
Participants: Maintainer, Agent
Related work items:

- `W-0153`
- `W-0154`
- `W-0155`

Related decision:

- `ADR-0071`

## Context

The previous completed slice implemented explicit bound identity route policy. Route protection can now distinguish public, request-token, bound-connection, session-validated, and bound-session route families while keeping ordinary protected routes on request-level access-token proof.

The queue was blocked at `M-081/W-0153`, asking for the next direction after bound identity route policy implementation.

The maintainer asked in Chinese for the agent to recommend the next ten steps and continue, with Nakama and Pitaya as key references.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese and to keep Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended selecting:

```text
define_logout_revocation_active_connection_gate
```

Reasoning:

- Route policy now has enough identity vocabulary; the next major lifecycle gap is what happens when token/session state is revoked while WebSocket connections remain active.
- Nakama provides strong reference pressure around authentication/session lifecycle, logout, expiration, refresh, and realtime socket implications.
- Pitaya provides strong reference pressure around separating acceptors, sessions, handlers, and connection management.
- The correct next move is a gate-only boundary, not immediate logout or socket-close code.

The agent used this ten-step plan:

1. Inspect the current work queue, git state, and validation output.
2. Re-read route policy, session validation, session creation, and connection binding boundaries.
3. Check Nakama and Pitaya reference alignment before selecting a lifecycle direction.
4. Select `define_logout_revocation_active_connection_gate`.
5. Add the English gate standard.
6. Add the Simplified Chinese translation.
7. Add `ADR-0071`.
8. Add change specs and conversation memory.
9. Update architecture manifests, module manifests, AGENTS guides, rules, and checks.
10. Run repository verification and record the result.

## Decisions

- Marked `M-081/W-0153` completed with selected direction `define_logout_revocation_active_connection_gate`.
- Added `M-082/W-0154` as the completed logout/revocation active-connection gate.
- Added `M-083/W-0155` as the next blocked confirmation gate.
- Accepted `ADR-0071`.
- Added `runtime.logout_revocation_active_connection_gate` as a repository check rule.
- Preserved deferrals for logout execution, token revocation execution, session revocation execution, active connection invalidation, connection registry, WebSocket close policy, Protobuf logout routes, session carriers, reconnect/epoch behavior, dependencies, memory durable session behavior, and direct Nakama/Pitaya compatibility.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-bound-identity-route-policy-implementation/`
- `changes/2026-05-18-define-logout-revocation-active-connection-gate/`
- `docs/logout-revocation-active-connection-gate.md`
- `docs/logout-revocation-active-connection-gate.zh-CN.md`
- `decisions/ADR-0071-logout-revocation-active-connection-gate.md`
- `conversations/2026-05-18-logout-revocation-active-connection-gate.md`
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

The gate establishes that future work must keep token revocation, runtime session revocation, and active-connection invalidation as separate state transitions. It recommends presented-token logout as the first future scope, requires an explicit connection registry before the server can target active sockets, and keeps WebSocket transport credential-neutral.

No runtime logout behavior, socket close behavior, Protobuf route, generated output, connection registry, or reconnect behavior was added.

## Follow-Up

- Define presented-token logout execution before implementing `LogoutAccessToken`.
- Define an active connection registry before targeting open sockets.
- Define transport close policy before using custom WebSocket close codes or reason text.
- Define reconnect and epoch behavior before duplicate replacement or resume behavior.
- Define protocol logout/session carriers before exposing logout commands, session ids, resume tokens, or connection epochs to clients.

Verification is recorded in `changes/2026-05-18-define-logout-revocation-active-connection-gate/verification.md`.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, session ids, cookies, query tokens, WebSocket subprotocol token material, database credentials, or GitHub tokens are recorded in this conversation log.
