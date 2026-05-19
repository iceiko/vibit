# Conversation: Logout Access Token Behavior Implementation

Date: 2026-05-18
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0073-logout-access-token-behavior-implementation.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-logout-access-token-behavior-gate/`
- `changes/2026-05-18-implement-logout-access-token-behavior/`

Related artifacts:

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `decisions/ADR-0073-logout-access-token-behavior-implementation.md`
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

The repository was blocked at `M-085/W-0157`, a next-direction confirmation gate after the logout access-token behavior gate.

Completed prerequisites included:

- Application authentication service skeleton.
- Device credential login behavior.
- Access-token validation behavior.
- Durable access-token records and revocation mutation vocabulary.
- Runtime session persistence and validation behavior.
- Session creation composition.
- Bound identity route policy.
- Gate-only logout/revocation active-connection and logout access-token behavior standards.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese.

## Agent Response Summary

The agent recommended:

```text
implement_logout_access_token_behavior
```

The reason was that `ADR-0072` already defined the first safe logout posture and the implementation could stay narrow: application service only, presented-token-only, verifier comparison before revocation, commit before success, and no session/socket/protocol side effects.

Nakama guided the token/session lifecycle pressure: revoked material should not authorize future protected gameplay requests. Pitaya guided the layering pressure: connection and session lifecycle infrastructure stays separate from handler credential parsing and application credential proof decisions.

The agent used this ten-step plan:

1. Confirm current queue, required boundaries, and existing implementation state.
2. Recheck Nakama/Pitaya reference points for this slice.
3. Implement `LogoutAccessToken` application service behavior.
4. Add focused logout behavior tests and fake repository recording.
5. Add implementation ADR, change specs, and conversation memory.
6. Update work queue, runtime/contracts/reference/conventions/module manifests.
7. Update AGENTS memory and bilingual guidance.
8. Update rule catalog and `tools/vibit` checks.
9. Run Go and repository verification and fix issues.
10. Confirm `inspect next` and report the next recommended direction.

## Decisions

- Close `M-085/W-0157`.
- Select `implement_logout_access_token_behavior`.
- Complete `M-086/W-0158` as the bounded logout access-token behavior implementation slice.
- Accept `ADR-0073`.
- Add `runtime.logout_access_token_behavior_implementation` as a repository check rule.
- Preserve deferrals for runtime session revocation, active connection invalidation, connection registry, WebSocket close policy, Protobuf logout route, protocol session carriers, existing envelope changes, refresh, logout-all, admin revocation, cleanup jobs, dependencies, memory durable session behavior, broader game backend modules, and direct Nakama/Pitaya API compatibility.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-logout-access-token-behavior-gate/`
- `changes/2026-05-18-implement-logout-access-token-behavior/`
- `decisions/ADR-0073-logout-access-token-behavior-implementation.md`
- `conversations/2026-05-18-logout-access-token-behavior-implementation.md`
- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
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

- Runtime session revocation remains deferred.
- Active connection registry and socket invalidation remain deferred.
- Protocol logout route exposure remains deferred.
- Reconnect, resume, duplicate replacement, and connection epoch behavior remain deferred.
- Refresh, logout-all sessions, admin revocation, and cleanup jobs remain deferred.
- Broader Nakama/Pitaya-inspired game backend modules remain deferred.

## Follow-Up

- Block at the next confirmation gate before choosing active connection registry behavior, reconnect/epoch behavior, protocol logout route, protocol session carriers, operations hardening, memory durable sessions, direct Nakama/Pitaya compatibility, or broader game backend planning.

## Redaction Notes

No raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, or GitHub tokens are recorded in this conversation log.
