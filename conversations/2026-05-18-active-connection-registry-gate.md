# Conversation: Active Connection Registry Gate

Date: 2026-05-18
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0074-active-connection-registry-gate.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-logout-access-token-behavior-implementation/`
- `changes/2026-05-18-define-active-connection-registry-gate/`
- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-gate/`

Related artifacts:

- `docs/active-connection-registry-gate.md`
- `docs/active-connection-registry-gate.zh-CN.md`
- `decisions/ADR-0074-active-connection-registry-gate.md`
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

The repository was blocked at `M-087/W-0159`, a next-direction confirmation gate after logout access-token behavior implementation.

Completed prerequisites included:

- Request-level access-token validation.
- Public device-credential login route.
- First-message connection binding.
- Runtime session persistence, validation, and login-time session creation.
- Bound identity route policy.
- Gate-only logout/revocation active-connection and logout access-token behavior standards.
- Presented-token-only `LogoutAccessToken` behavior implementation.

`LogoutAccessToken` now revokes the verified presented access-token record but deliberately does not close active WebSocket connections, revoke runtime sessions, add a connection registry, expose protocol logout routes, or change reconnect behavior.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese.

## Agent Response Summary

The agent recommended:

```text
define_active_connection_registry_gate
```

The reason was that service-level token revocation now exists, but future logout/revocation, kick/disconnect, duplicate replacement, reconnect, and multiplayer lifecycle behavior still need a safe server-owned model for active WebSocket connections.

Nakama guided the lifecycle pressure: realtime sockets and authenticated session material need coordinated policy. Pitaya guided the layering pressure: acceptors, sessions, handlers, and connection management remain separate surfaces.

The agent used this ten-step plan:

1. Confirm `M-087/W-0159` and current blocked state.
2. Select `define_active_connection_registry_gate`.
3. Add `M-088/W-0160` as a gate-only active connection registry slice.
4. Define registry owner, scope, first posture, and non-goals.
5. Record single-process in-memory non-durable first posture.
6. Preserve logout, socket close, reconnect, protocol carrier, distributed storage, and direct compatibility deferrals.
7. Add ADR, change specs, and conversation memory.
8. Update manifests, AGENTS guides, rule catalog, and `tools/vibit`.
9. Run Go and repository verification and fix issues.
10. Confirm `inspect next` stops at the next confirmation gate.

## Decisions

- Close `M-087/W-0159`.
- Select `define_active_connection_registry_gate`.
- Complete `M-088/W-0160` as the bounded gate-only active connection registry slice.
- Accept `ADR-0074`.
- Add `runtime.active_connection_registry_gate` as a repository check rule.
- Stop at `M-089/W-0161` for the next explicit maintainer decision.
- Preserve deferrals for active connection registry implementation, WebSocket close policy, kick/disconnect behavior, runtime session revocation, reconnect/epoch behavior, protocol logout routes, protocol session carriers, durable/distributed registry storage, dependencies, memory durable session behavior, broader game backend modules, and direct Nakama/Pitaya API compatibility.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-logout-access-token-behavior-implementation/`
- `changes/2026-05-18-define-active-connection-registry-gate/`
- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-gate/`
- `docs/active-connection-registry-gate.md`
- `docs/active-connection-registry-gate.zh-CN.md`
- `decisions/ADR-0074-active-connection-registry-gate.md`
- `conversations/2026-05-18-active-connection-registry-gate.md`
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

- Active connection registry implementation remains deferred.
- WebSocket close policy, close codes, close reasons, kick, and disconnect behavior remain deferred.
- Runtime session revocation remains deferred.
- Reconnect, resume, duplicate replacement, and durable epoch behavior remain deferred.
- Protocol logout routes and protocol session carriers remain deferred.
- Durable/distributed connection registry storage remains deferred.
- Broader Nakama/Pitaya-inspired game backend modules remain deferred.

## Follow-Up

- Block at the next confirmation gate before choosing active connection registry implementation, WebSocket close policy, reconnect/epoch behavior, protocol logout route, protocol session carriers, operations hardening, memory durable sessions, direct Nakama/Pitaya compatibility, or broader game backend planning.

## Redaction Notes

No raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, remote addresses, database credentials, or GitHub tokens are recorded in this conversation log.
