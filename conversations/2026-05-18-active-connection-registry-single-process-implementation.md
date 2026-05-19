# Conversation: Active Connection Registry Single Process Implementation

Date: 2026-05-18
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0075-active-connection-registry-single-process-implementation.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-gate/`
- `changes/2026-05-18-implement-active-connection-registry-single-process/`

Related artifacts:

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
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

The repository was blocked at `M-089/W-0161`, a next-direction confirmation gate after the active connection registry gate.

Completed prerequisites included:

- Request-level access-token validation.
- First-message connection binding.
- Durable runtime session validation.
- Login-time runtime session creation.
- Explicit bound/session route policy families.
- Presented access-token logout behavior.
- Gate-only active connection registry guidance in `docs/active-connection-registry-gate.md`.

The remaining implementation gap was a server-owned active connection targeting primitive.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese and to keep Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended:

```text
implement_active_connection_registry_single_process
```

The reason was that future logout socket invalidation, kick/disconnect, duplicate replacement, reconnect/epoch behavior, presence, rooms, and match runtime all need a safe application-owned active connection model before they can target open sockets.

Nakama guided the lifecycle pressure: realtime sockets and authenticated session material need coordinated state. Pitaya guided the layering: acceptors, sessions, route handlers, and connection management remain separate.

The agent used this ten-step plan:

1. Close `M-089/W-0161` by selecting `implement_active_connection_registry_single_process`.
2. Create `M-090/W-0162` as the bounded implementation slice.
3. Create `M-091/W-0163` as the next confirmation gate.
4. Add `runtime/internal/app/connection/registry.go`.
5. Add `runtime/internal/app/connection/registry_test.go`.
6. Implement register, bind, close, invalidate, find, and list capabilities.
7. Add focused tests for validation, duplicate active records, metadata-only rejection, terminal-state exclusion, copy semantics, and record redaction.
8. Add `ADR-0075`, change specs, and conversation memory.
9. Update architecture manifests, module manifest, AGENTS guides, rules, and checks.
10. Run Go and repository verification.

## Decisions

- Close `M-089/W-0161`.
- Select `implement_active_connection_registry_single_process`.
- Complete `M-090/W-0162` as a conservative registry implementation slice.
- Accept `ADR-0075`.
- Add `runtime.active_connection_registry_single_process_implementation` as a repository check rule.
- Preserve deferrals for WebSocket close policy, concrete socket close handoff, kick/disconnect behavior, runtime session revocation, duplicate replacement, reconnect/epoch behavior, protocol logout routes, protocol session carriers, durable/distributed registry storage, dependencies, memory durable session behavior, direct Nakama/Pitaya API compatibility, and broader game backend behavior.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-gate/`
- `changes/2026-05-18-implement-active-connection-registry-single-process/`
- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-implementation/`
- `decisions/ADR-0075-active-connection-registry-single-process-implementation.md`
- `conversations/2026-05-18-active-connection-registry-single-process-implementation.md`
- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
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

- Whether logout or token revocation closes active WebSocket connections remains deferred.
- WebSocket close codes, close reasons, kick/disconnect behavior, and transport handoff remain deferred.
- Reconnect, resume, duplicate replacement, and durable epoch policy remain deferred.
- Client-visible protocol logout routes and session carriers remain deferred.
- Operations and observability posture remain deferred.
- Presence, rooms, parties, match runtime, groups, and broader game backend modules remain deferred.

## Follow-Up

Block at the next confirmation gate before choosing WebSocket close policy, protocol logout routes, reconnect/epoch behavior, protocol session carriers, operations hardening, memory durable sessions, direct Nakama/Pitaya compatibility, or broader game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, or GitHub tokens are recorded in this conversation log.
