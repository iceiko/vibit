# Conversation: First Message Binding Next Direction And Implementation Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-gate/`
- `changes/2026-05-17-define-first-message-connection-binding-implementation-gate/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/first-message-connection-binding-implementation-gate.md`
- `docs/first-message-connection-binding-implementation-gate.zh-CN.md`
- `decisions/ADR-0058-first-message-connection-binding-implementation-gate.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The first-message connection binding gate was completed in `M-052/W-0124`. It selected a future `runtime.authentication.BindConnection` system route but did not authorize implementation.

The work queue was blocked at `M-053/W-0125`, asking for the next major direction before crossing first-message connection binding implementation, PostgreSQL session persistence schema, logout/revocation active-connection behavior, reconnect/epoch behavior, operations posture, memory durable authentication behavior, direct Nakama/Pitaya API compatibility, or broader game backend scope.

## Maintainer Narrative

The maintainer asked:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。
```

The maintainer had also emphasized that Nakama and Pitaya should remain important references because existing game servers understand current game backend needs deeply.

## Agent Response Summary

The agent recommended selecting:

```text
define_first_message_connection_binding_implementation_gate
```

The reasoning was that `ADR-0057` had already selected the future `BindConnection` route, but adding Protobuf messages, generated output, protocol adapter behavior, application binding registries, startup composition, or route-policy bound identity still crosses explicit ask-first boundaries.

Nakama guides the authenticated realtime socket lifecycle and the distinction between token/session lifecycle and active socket lifecycle. Pitaya guides acceptor/session/handler separation and warns against placing identity logic in transport handlers.

The agent then defined a gate-only standard rather than implementing connection binding.

## Decisions

- Close `M-053/W-0125`.
- Select `define_first_message_connection_binding_implementation_gate`.
- Create and complete `M-054/W-0126` as a gate-only standard.
- Accept `ADR-0058`.
- Add `runtime.first_message_connection_binding_implementation_gate` as a repository check rule.
- Defer actual `BindConnection` implementation and all durable session or active connection behavior.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-gate/`
- `changes/2026-05-17-define-first-message-connection-binding-implementation-gate/`
- `docs/first-message-connection-binding-implementation-gate.md`
- `docs/first-message-connection-binding-implementation-gate.zh-CN.md`
- `decisions/ADR-0058-first-message-connection-binding-implementation-gate.md`
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

- Exact Protobuf field numbers are deferred to the future implementation slice.
- Whether route policy can use bound identity is deferred to the future implementation slice.
- Whether a process-local connection registry is sufficient before durable session persistence is deferred.
- Duplicate connection, kick, reconnect, epoch, logout-triggered invalidation, and presence/group attachment remain separate gates.

## Follow-Up

- Implement first-message connection binding only after the maintainer selects `implement_first_message_connection_binding`.
- Define PostgreSQL session persistence before any durable session table, repository, migration, or cleanup behavior.
- Define logout/revocation and reconnect/epoch behavior before bound connections can be invalidated, replaced, or resumed.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, or GitHub tokens are recorded in this conversation log.
