# Conversation: Session Handshake Next Direction And First Message Binding Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-session-handshake-ratification/`
- `changes/2026-05-17-define-first-message-connection-binding-gate/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/first-message-connection-binding-gate.md`
- `docs/first-message-connection-binding-gate.zh-CN.md`
- `decisions/ADR-0057-first-message-connection-binding-gate.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The session persistence and WebSocket handshake posture was ratified in `M-050/W-0122`. The work queue was blocked at `M-051/W-0123`, asking for the next major direction before crossing first-message connection binding, session persistence schema, WebSocket handshake authentication, logout/revocation active-connection behavior, reconnect/epoch behavior, operations posture, direct Nakama/Pitaya API compatibility, or broader game backend scope.

## Maintainer Narrative

The maintainer asked:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。
```

The maintainer had also emphasized that Nakama and Pitaya should remain important references because existing game servers understand current game backend needs deeply.

## Agent Response Summary

The agent recommended selecting:

```text
define_first_message_connection_binding_gate
```

The reasoning was that `ADR-0056` already prefers protocol/application first-message binding for connection-level identity. Nakama supports the capability sequence of authentication followed by realtime socket usage with authenticated state. Pitaya supports the architecture separation of transport acceptors, session binding, handlers, groups, and cluster concerns.

The agent then defined a gate-only standard rather than implementing connection binding. The gate selects a future `runtime.authentication.BindConnection` system route, keeps WebSocket transport credential-neutral, keeps the existing Protobuf envelope unchanged, preserves request-level access-token validation as the current path, and defers session persistence, route policy use of bound identity, reconnect/epoch, logout/revocation, migrations, dependencies, and direct Nakama/Pitaya API compatibility.

## Decisions

- Close `M-051/W-0123`.
- Select `define_first_message_connection_binding_gate`.
- Create and complete `M-052/W-0124` as a gate-only standard.
- Accept `ADR-0057`.
- Add `runtime.first_message_connection_binding_gate` as a repository check rule.
- Defer implementation of connection binding and all durable session or active connection behavior.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-session-handshake-ratification/`
- `changes/2026-05-17-define-first-message-connection-binding-gate/`
- `docs/first-message-connection-binding-gate.md`
- `docs/first-message-connection-binding-gate.zh-CN.md`
- `decisions/ADR-0057-first-message-connection-binding-gate.md`
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

- The exact `BindConnectionRequest` and `BindConnectionResponse` field numbers are deferred.
- Connection binding registry shape and startup composition are deferred.
- Whether bound identity can satisfy protected routes without per-request wrapper is deferred to an implementation gate and route-policy update.
- Durable session persistence remains separate from connection binding.
- Duplicate connection, kick, reconnect, epoch, logout-triggered invalidation, and presence/group attachment remain separate gates.

## Follow-Up

- Define a first-message connection binding implementation gate before adding protocol messages, generated output, connection registries, route-policy bound identity, or startup wiring.
- Define PostgreSQL session persistence before any durable session table, repository, migration, or cleanup behavior.
- Define logout/revocation and reconnect/epoch behavior before bound connections can be invalidated, replaced, or resumed.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, or GitHub tokens are recorded in this conversation log.
