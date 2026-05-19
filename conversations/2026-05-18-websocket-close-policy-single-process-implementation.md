# Conversation: WebSocket Close Policy Single Process Implementation

Date: 2026-05-18
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0077-websocket-close-policy-single-process-implementation.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-gate/`
- `changes/2026-05-18-implement-websocket-close-policy-single-process/`

Related artifacts:

- `runtime/internal/app/connection/close_policy.go`
- `runtime/internal/app/connection/close_policy_test.go`
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

The repository was blocked at `M-093/W-0165`, a next-direction confirmation gate after the WebSocket close policy gate.

Completed prerequisites included:

- Presented access-token logout behavior.
- The application-owned single-process active connection registry.
- Gate-only WebSocket close policy guidance in `docs/websocket-close-policy-gate.md`.

The remaining implementation gap was a server-owned close policy primitive that could resolve trusted registry targets and record invalidation intent without closing concrete sockets.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese and to keep Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended:

```text
implement_websocket_close_policy_single_process
```

The reason was that future logout routes, token revocation composition, runtime session revocation, duplicate replacement, reconnect/epoch behavior, admin disconnect, and operational drain all need one application-owned close policy vocabulary before transport can close real sockets.

Nakama guided the lifecycle pressure: realtime socket lifecycle and token/session lifecycle need explicit server semantics, and session logout/token invalidation is distinct from server-directed socket disconnect. Pitaya guided the layering: acceptors, sessions, route handlers, and connection management remain separate.

The agent used this ten-step plan:

1. Close `M-093/W-0165` by selecting `implement_websocket_close_policy_single_process`.
2. Create `M-094/W-0166` as the bounded implementation slice.
3. Create `M-095/W-0167` as the next confirmation gate.
4. Add `runtime/internal/app/connection/close_policy.go`.
5. Add `runtime/internal/app/connection/close_policy_test.go`.
6. Implement registry-backed target resolution by connection id/epoch, player id, runtime session id, and access-token record id.
7. Mark matched active bound records invalidated and emit redacted close intents with `mark_invalidated_only`.
8. Add `ADR-0077`, change specs, and conversation memory.
9. Update architecture manifests, module manifest, AGENTS guides, rules, and checks.
10. Run Go and repository verification.

## Decisions

- Close `M-093/W-0165`.
- Select `implement_websocket_close_policy_single_process`.
- Complete `M-094/W-0166` as a conservative close policy implementation slice.
- Accept `ADR-0077`.
- Add `runtime.websocket_close_policy_single_process_implementation` as a repository check rule.
- Preserve deferrals for concrete WebSocket close handoff, close codes, close reason text, protocol close messages, protocol logout routes, runtime session revocation close behavior, duplicate replacement, reconnect/epoch behavior, protocol session carriers, durable/distributed registry storage, operations/admin disconnect, dependencies, memory durable session behavior, direct Nakama/Pitaya API compatibility, and broader game backend behavior.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-gate/`
- `changes/2026-05-18-implement-websocket-close-policy-single-process/`
- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-implementation/`
- `decisions/ADR-0077-websocket-close-policy-single-process-implementation.md`
- `conversations/2026-05-18-websocket-close-policy-single-process-implementation.md`
- `runtime/internal/app/connection/close_policy.go`
- `runtime/internal/app/connection/close_policy_test.go`
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

- Protocol logout route behavior remains deferred.
- Concrete WebSocket transport close handoff remains deferred.
- Close codes, close reason text, kick/disconnect public behavior, and protocol close messages remain deferred.
- Reconnect, resume, duplicate replacement, and durable epoch policy remain deferred.
- Client-visible protocol session carriers remain deferred.
- Operations and observability posture remain deferred.
- Presence, rooms, parties, match runtime, groups, and broader game backend modules remain deferred.

## Follow-Up

Block at the next confirmation gate before choosing protocol logout routes, reconnect/epoch behavior, protocol session carriers, concrete transport close handoff, operations hardening, memory durable sessions, direct Nakama/Pitaya compatibility, or broader game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, close reason text, remote addresses, headers, or GitHub tokens are recorded in this conversation log.
