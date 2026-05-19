# Conversation: Logout Access Token Behavior Gate

Date: 2026-05-18
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0072-logout-access-token-behavior-gate.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-logout-revocation-active-connection-gate/`
- `changes/2026-05-18-define-logout-access-token-behavior-gate/`

Related artifacts:

- `docs/logout-access-token-behavior-gate.md`
- `docs/logout-access-token-behavior-gate.zh-CN.md`
- `decisions/ADR-0072-logout-access-token-behavior-gate.md`
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

The repository was blocked at `M-083/W-0155`, a next-direction confirmation gate after the logout/revocation active-connection gate.

Completed prerequisites included:

- Request-level access-token validation.
- Public device-credential login.
- Durable access-token records.
- Login-time durable runtime session creation.
- First-message connection binding.
- Explicit bound identity route policy.
- A gate-only logout/revocation active-connection boundary that separated token revocation, runtime session revocation, and active socket invalidation.

The remaining immediate gap was the first future behavior boundary for `LogoutAccessToken`.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese and to keep Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended:

```text
define_logout_access_token_behavior_gate
```

The reason was that the previous gate selected `presented_access_token_only` as the recommended first logout scope. Before implementing token revocation, vibit needs to define proof validation, verifier comparison, transaction ordering, public error collapse, and explicit deferrals for sessions, sockets, protocol routes, and reconnect behavior.

Nakama guided the session lifecycle pressure: revoked session or token material must not authorize future gameplay requests, while refresh, logout-all, session management, and realtime socket behavior remain distinct lifecycle surfaces. Pitaya guided the layering: session and connection lifecycle infrastructure stays separate from handler credential parsing.

The agent used this ten-step plan:

1. Fix the previous conversation memory check failure.
2. Inspect the current queue, manifests, service skeleton, repository vocabulary, and existing logout contract.
3. Select `define_logout_access_token_behavior_gate`.
4. Close `M-083/W-0155` with that selected direction.
5. Create `M-084/W-0156` as the bounded gate-only logout access-token behavior slice.
6. Create `M-085/W-0157` as the next confirmation gate.
7. Add the English gate standard.
8. Add the Simplified Chinese translation and `ADR-0072`.
9. Add change specs, update manifests, module manifests, AGENTS guides, rule catalog, and `tools/vibit`.
10. Run repository verification and record the result.

## Decisions

- Close `M-083/W-0155`.
- Select `define_logout_access_token_behavior_gate`.
- Complete `M-084/W-0156` as a gate-only logout access-token behavior boundary slice.
- Accept `ADR-0072`.
- Add `runtime.logout_access_token_behavior_gate` as a repository check rule.
- Preserve deferrals for `LogoutAccessToken` implementation, token revocation execution, runtime session revocation, active connection invalidation, connection registry, WebSocket close policy, Protobuf logout route, protocol session carriers, existing envelope changes, refresh, logout-all, admin revocation, cleanup jobs, dependencies, memory durable session behavior, and direct Nakama/Pitaya API compatibility.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-logout-revocation-active-connection-gate/`
- `changes/2026-05-18-define-logout-access-token-behavior-gate/`
- `docs/logout-access-token-behavior-gate.md`
- `docs/logout-access-token-behavior-gate.zh-CN.md`
- `decisions/ADR-0072-logout-access-token-behavior-gate.md`
- `conversations/2026-05-18-logout-access-token-behavior-gate.md`
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

- `LogoutAccessToken` implementation remains deferred.
- Whether already revoked tokens should eventually return idempotent success remains deferred.
- Runtime session revocation remains deferred.
- Active connection registry and socket invalidation remain deferred.
- Protocol logout route exposure remains deferred.
- Reconnect, resume, duplicate replacement, and connection epoch behavior remain deferred.
- Refresh, logout-all sessions, and admin revocation remain deferred.
- Broader Nakama/Pitaya-inspired game backend modules remain deferred.

## Follow-Up

- Block at the next confirmation gate before choosing logout implementation, active connection registry behavior, reconnect/epoch behavior, protocol logout route, protocol session carriers, operations hardening, memory durable sessions, direct Nakama/Pitaya compatibility, or broader game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, session ids, cookies, query tokens, WebSocket subprotocol token material, database credentials, or GitHub tokens are recorded in this conversation log.
