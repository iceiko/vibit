# Conversation: Bound Identity Route Policy Implementation

Date: 2026-05-18
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0070-bound-identity-route-policy-implementation.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-bound-identity-route-policy-gate/`
- `changes/2026-05-18-implement-bound-identity-route-policy/`

Related artifacts:

- `runtime/internal/app/route_authentication.go`
- `runtime/internal/app/route_authentication_test.go`
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

The repository was blocked at `M-079/W-0151`, a next-direction confirmation gate after the bound identity route policy gate.

Completed prerequisites included:

- Request-level access-token route protection through an explicit Protobuf wrapper.
- First-message connection binding.
- Durable `runtime_sessions` persistence and repository.
- Persistent runtime session validation.
- Login-time durable runtime session creation.
- Gate-only route policy guidance in `docs/bound-identity-route-policy-gate.md`.

The remaining implementation gap was route policy: vibit had the gate and the identity sources, but `RouteProtector` still had only public-versus-protected behavior.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese and to keep Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended:

```text
implement_bound_identity_route_policy
```

The reason was that the previous gate had already defined policy families, and the safest next step was to implement the application-owned vocabulary and fail-closed checks without changing ordinary protected route behavior.

Nakama guided the session-backed access model: authenticated session material may become central to gameplay access, but logout, refresh, session management, and active socket behavior remain separate lifecycle decisions. Pitaya guided the layering: acceptors, sessions, route policy, and handlers stay separated.

The agent used this ten-step plan:

1. Close `M-079/W-0151` by selecting `implement_bound_identity_route_policy`.
2. Create `M-080/W-0152` as the bounded implementation slice.
3. Create `M-081/W-0153` as the next confirmation gate.
4. Add explicit route policy requirement vocabulary.
5. Preserve default request-token behavior for ordinary protected routes.
6. Add explicit bound, session-validated, and bound-session fail-closed checks.
7. Add focused tests for classification, protected defaults, metadata-only rejection, source agreement, and redaction.
8. Add `ADR-0070`, change specs, and conversation memory.
9. Update architecture manifests, module manifest, AGENTS guides, rules, and checks.
10. Run Go and repository verification.

## Decisions

- Close `M-079/W-0151`.
- Select `implement_bound_identity_route_policy`.
- Complete `M-080/W-0152` as a conservative route policy implementation slice.
- Accept `ADR-0070`.
- Add `runtime.bound_identity_route_policy_implementation` as a repository check rule.
- Preserve deferrals for production route reclassification, WebSocket handshake authentication, transport credential carriers, Protobuf session carriers, existing envelope changes, logout/revocation active-connection invalidation, reconnect/epoch behavior, cleanup jobs, dependencies, memory durable session behavior, direct Nakama/Pitaya API compatibility, and broader game backend behavior.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-bound-identity-route-policy-gate/`
- `changes/2026-05-18-implement-bound-identity-route-policy/`
- `decisions/ADR-0070-bound-identity-route-policy-implementation.md`
- `conversations/2026-05-18-bound-identity-route-policy-implementation.md`
- `runtime/internal/app/route_authentication.go`
- `runtime/internal/app/route_authentication_test.go`
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
- Reconnect, resume, duplicate replacement, and connection epoch policy remain deferred.
- Client-visible protocol session carriers remain deferred.
- Production route reclassification away from request-token proof remains deferred.
- Operations and observability posture remain deferred.
- Presence, rooms, parties, match runtime, groups, and broader game backend modules remain deferred.

## Follow-Up

- Block at the next confirmation gate before choosing logout/revocation active-connection behavior, reconnect/epoch behavior, protocol/session carriers, route-policy production reclassification, operations hardening, memory durable sessions, direct Nakama/Pitaya compatibility, or broader game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, session ids, cookies, query tokens, WebSocket subprotocol token material, database credentials, or GitHub tokens are recorded in this conversation log.
