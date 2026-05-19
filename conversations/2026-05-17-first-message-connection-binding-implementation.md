# Conversation: First Message Connection Binding Implementation

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-implementation-gate/`
- `changes/2026-05-17-implement-first-message-connection-binding/`

Related artifacts:

- `proto/vibit/authentication/v1/authentication.proto`
- `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go`
- `runtime/internal/app/connection_binding.go`
- `runtime/internal/platform/protocol/protobuf/connection_binding.go`
- `runtime/internal/platform/transport/ws/server.go`
- `runtime/cmd/vibit-server/main.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-054/W-0126` had completed the first-message connection binding implementation gate. `ADR-0057` had selected a future `runtime.authentication.BindConnection` system route, and `ADR-0058` had defined the bounded implementation slice. The work queue was blocked at `M-055/W-0127` pending maintainer direction before implementing BindConnection or moving into session persistence, logout/revocation active-connection behavior, reconnect/epoch behavior, operations posture, memory durable authentication behavior, direct Nakama/Pitaya API compatibility, or broader game backend scope.

## Maintainer Narrative

The maintainer asked in Chinese:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya
```

The maintainer had also asked the agent to speak Chinese and repeatedly emphasized that Nakama and Pitaya should be treated as important references because existing game servers understand current game backend needs deeply.

## Agent Response Summary

The agent recommended selecting:

```text
implement_first_message_connection_binding
```

The reasoning was that the route and implementation boundary were already selected by `ADR-0057` and `ADR-0058`, while durable session storage, route-policy use of bound identity, logout/revocation, reconnect, operations posture, and broader realtime features should remain separate gates.

The agent implemented the bounded W-0128 slice: authentication Protobuf messages, Buf-generated output, application-owned binding, Protobuf adapter handling, WebSocket connection metadata handoff, PostgreSQL startup injection, focused tests, architecture manifest updates, AGENTS guide updates, change specs, conversation memory, and repository checks.

Nakama informed the authenticate-before-realtime-socket lifecycle and the distinction between token/session lifetime and socket connection lifetime. Pitaya informed transport acceptor, session-like binding, and handler-context separation. The implementation adapts those ideas without direct public API compatibility.

## Decisions

- Close `M-055/W-0127`.
- Select `implement_first_message_connection_binding`.
- Complete `M-056/W-0128`.
- Add `runtime.first_message_connection_binding_implementation` as a repository check rule.
- Add `runtime.authentication.BindConnection` as a Protobuf `system` route.
- Keep WebSocket transport credential-neutral.
- Keep the existing Protobuf envelope unchanged.
- Keep ordinary protected routes on the request-level `AuthenticatedRequest` wrapper.
- Keep `SessionValidated` false until session persistence is separately ratified and implemented.
- Keep memory startup binding unavailable.
- Defer durable sessions, connection registries, route-policy bound identity, logout/revocation active-connection behavior, reconnect/epoch policy, and broader realtime/game backend behavior.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-implementation-gate/`
- `changes/2026-05-17-implement-first-message-connection-binding/`
- `proto/vibit/authentication/v1/authentication.proto`
- `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go`
- `runtime/internal/app/route_authentication.go`
- `runtime/internal/app/connection_binding.go`
- `runtime/internal/app/connection_binding_test.go`
- `runtime/internal/platform/protocol/protobuf/connection_binding.go`
- `runtime/internal/platform/protocol/protobuf/connection_binding_test.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/transport/ws/server.go`
- `runtime/internal/platform/transport/ws/server_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `conversations/2026-05-17-first-message-connection-binding-implementation.md`
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
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Whether durable connection/session binding should be modeled first through PostgreSQL session persistence schema.
- Whether logout/revocation should actively invalidate open WebSocket connections.
- Whether reconnect/resume/duplicate connection replacement should use a strict epoch policy.
- Whether connection-bound identity can later satisfy ordinary protected route policy.
- Which broader Nakama/Pitaya-inspired game backend family should follow after the authentication/session transport foundation is stable.

## Follow-Up

- Block at `M-057/W-0129` until the maintainer selects the next direction.
- Candidate next directions are PostgreSQL session persistence schema, logout/revocation active-connection behavior, reconnect/epoch behavior, bound identity route policy, operations/observability/admin tooling, or broader game backend modules after Nakama/Pitaya review.
- Run full runtime tests, Buf checks, repository checks, memory checks, change checks, and diff checks.
- Record any remaining warnings or unavailable live PostgreSQL verification explicitly.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, or GitHub tokens are recorded in this conversation log.
