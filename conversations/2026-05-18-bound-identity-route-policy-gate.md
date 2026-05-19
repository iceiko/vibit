# Conversation: Bound Identity Route Policy Gate

Date: 2026-05-18
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0069-bound-identity-route-policy-gate.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-session-creation-composition-implementation/`
- `changes/2026-05-18-define-bound-identity-route-policy-gate/`

Related artifacts:

- `docs/bound-identity-route-policy-gate.md`
- `docs/bound-identity-route-policy-gate.zh-CN.md`
- `decisions/ADR-0069-bound-identity-route-policy-gate.md`
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

The repository was blocked at `M-077/W-0149`, a next-direction confirmation gate after session creation composition implementation.

Completed prerequisites included:

- Request-level access-token route protection through an explicit Protobuf wrapper.
- First-message connection binding.
- Durable `runtime_sessions` migration and repository.
- PostgreSQL session adapter.
- Persistent runtime session validation.
- Login-time durable runtime session creation with token storage in one unit of work.

The remaining design gap was route policy: vibit had several identity sources, but no boundary deciding when request-token identity, bound connection identity, or session-validated identity may authorize routes.

## Maintainer Narrative

The maintainer said:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

The maintainer had also asked the agent to communicate in Chinese and to use Nakama and Pitaya as important game-server references.

## Agent Response Summary

The agent recommended:

```text
define_bound_identity_route_policy_gate
```

The reason was that login now creates durable sessions, and validation/binding paths exist, but route policy has not defined when to trust request-token identity, bound connection identity, or session-validated identity.

Nakama guided the gameplay-access pressure: authenticated session material eventually becomes central to gameplay access while logout, refresh, session management, and active socket behavior stay separate lifecycle concerns. Pitaya guided the layering: acceptors, sessions, route handlers, and context handoff should stay separated.

The agent used this ten-step plan:

1. Inspect the current work queue, git state, and validation output.
2. Re-read route protection, connection binding, session validation, and session creation boundaries.
3. Check Nakama and Pitaya reference alignment before selecting a capability direction.
4. Select `define_bound_identity_route_policy_gate`.
5. Add the English gate standard.
6. Add the Simplified Chinese translation.
7. Add `ADR-0069`.
8. Add change spec and conversation memory.
9. Update architecture manifests, module manifests, AGENTS guides, and check rules.
10. Run repository verification and record the result.

## Decisions

- Close `M-077/W-0149`.
- Select `define_bound_identity_route_policy_gate`.
- Complete `M-078/W-0150` as a gate-only route-policy boundary slice.
- Accept `ADR-0069`.
- Add `runtime.bound_identity_route_policy_gate` as a repository check rule.
- Preserve deferrals for route-policy implementation, removing request-level proof from ordinary protected routes, WebSocket handshake authentication, transport credential carriers, Protobuf session carriers, existing envelope changes, generated output, logout/revocation active-connection invalidation, reconnect/epoch behavior, cleanup jobs, dependencies, memory durable session behavior, direct Nakama/Pitaya API compatibility, and broader game backend behavior.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-session-creation-composition-implementation/`
- `changes/2026-05-18-define-bound-identity-route-policy-gate/`
- `docs/bound-identity-route-policy-gate.md`
- `docs/bound-identity-route-policy-gate.zh-CN.md`
- `decisions/ADR-0069-bound-identity-route-policy-gate.md`
- `conversations/2026-05-18-bound-identity-route-policy-gate.md`
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

- Implementation of route-policy families remains deferred.
- Whether ordinary protected routes may ever stop requiring per-request proof remains deferred.
- Client-visible session carrier and login response session shape remain deferred.
- Logout/revocation active-connection behavior remains deferred.
- Reconnect/resume/duplicate replacement and connection epoch behavior remain deferred.
- Operations and observability classification for session ids remains deferred.
- Presence, rooms, parties, match runtime, groups, and broader game backend modules remain deferred.

## Follow-Up

- Block at the next confirmation gate before choosing route-policy implementation, logout/revocation active-connection behavior, reconnect/epoch behavior, protocol/session carriers, operations hardening, memory durable sessions, direct Nakama/Pitaya compatibility, or broader game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, session ids, cookies, query tokens, WebSocket subprotocol token material, database credentials, or GitHub tokens are recorded in this conversation log.
