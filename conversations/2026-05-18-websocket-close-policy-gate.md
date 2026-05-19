# Conversation: WebSocket Close Policy Gate

Date: 2026-05-18
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0076-websocket-close-policy-gate.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-implementation/`
- `changes/2026-05-18-define-websocket-close-policy-gate/`
- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-gate/`

Related artifacts:

- `docs/websocket-close-policy-gate.md`
- `docs/websocket-close-policy-gate.zh-CN.md`
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

The repository was blocked at `M-091/W-0163`, a next-direction confirmation gate after the active connection registry single-process implementation.

Completed prerequisites included:

- Request-level access-token validation.
- Public device-credential login route.
- First-message connection binding.
- Runtime session persistence, validation, and login-time session creation.
- Bound identity route policy.
- Presented access-token logout behavior.
- Gate-only logout/revocation active-connection guidance.
- Application-owned in-memory active connection registry behavior.

The active connection registry can now find active bound records by player id, runtime session id, access-token record id, or connection id/epoch, but it deliberately cannot close sockets or choose lifecycle policy.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese and to keep Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended:

```text
define_websocket_close_policy_gate
```

The reason was that the registry can now target active connections, but no policy yet defines when invalidation closes sockets, which close reason classes are allowed, whether logout close failure affects logout success, or which layer performs the concrete socket close.

Nakama guided the lifecycle pressure: session logout/token invalidation and realtime socket disconnection are distinct lifecycle actions. Pitaya guided the layering pressure: acceptors, sessions, handlers, and connection management remain separate surfaces.

The agent used this ten-step plan:

1. Complete `M-091/W-0163` by selecting `define_websocket_close_policy_gate`.
2. Create `M-092/W-0164` as the bounded gate-only WebSocket close policy slice.
3. Create `M-093/W-0165` as the next confirmation gate.
4. Add the English WebSocket close policy gate standard.
5. Add the Simplified Chinese translation.
6. Add `ADR-0076`.
7. Add change specs and conversation memory.
8. Update architecture manifests, module manifest, and AGENTS guides.
9. Add `runtime.websocket_close_policy_gate` to the rule catalog and runtime checks.
10. Run repository verification and confirm `inspect next` stops at `M-093/W-0165`.

## Decisions

- Close `M-091/W-0163`.
- Select `define_websocket_close_policy_gate`.
- Complete `M-092/W-0164` as a conservative gate-only close policy boundary.
- Accept `ADR-0076`.
- Add `runtime.websocket_close_policy_gate` as a repository check rule.
- Stop at `M-093/W-0165` for the next explicit maintainer decision.
- Preserve deferrals for WebSocket close implementation, transport close handoff code, close codes, close reasons, kick/disconnect behavior, logout-triggered socket close, runtime session revocation, reconnect/epoch behavior, protocol logout routes, protocol session carriers, durable/distributed registry storage, dependencies, memory durable session behavior, direct Nakama/Pitaya API compatibility, and broader game backend behavior.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-implementation/`
- `changes/2026-05-18-define-websocket-close-policy-gate/`
- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-gate/`
- `docs/websocket-close-policy-gate.md`
- `docs/websocket-close-policy-gate.zh-CN.md`
- `decisions/ADR-0076-websocket-close-policy-gate.md`
- `conversations/2026-05-18-websocket-close-policy-gate.md`
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

- Whether logout closes the presenting socket remains deferred.
- Whether token revocation targets by token record, runtime session, player, or connection id remains deferred.
- WebSocket close codes, close reasons, kick/disconnect messages, and protocol close messages remain deferred.
- Reconnect, resume, duplicate replacement, and durable epoch policy remain deferred.
- Client-visible protocol logout routes and session carriers remain deferred.
- Operations and admin disconnect surfaces remain deferred.
- Presence, rooms, parties, match runtime, groups, and broader game backend modules remain deferred.

## Follow-Up

Block at the next confirmation gate before choosing WebSocket close policy implementation, protocol logout routes, reconnect/epoch behavior, protocol session carriers, operations hardening, memory durable sessions, direct Nakama/Pitaya compatibility, or broader game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, remote addresses, database credentials, or GitHub tokens are recorded in this conversation log.
