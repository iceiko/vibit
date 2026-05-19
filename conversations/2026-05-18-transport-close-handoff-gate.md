# Conversation: Transport Close Handoff Gate

Date: 2026-05-18
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-18-define-transport-close-handoff-gate/`

Related artifacts:

- `docs/transport-close-handoff-gate.md`
- `docs/transport-close-handoff-gate.zh-CN.md`
- `decisions/ADR-0080-transport-close-handoff-gate.md`

## Context

The maintainer clarified that vibit should be a Nakama/Pitaya-class game backend product, and the active roadmap keeps lifecycle closure before broad feature expansion.

`W-0170` implemented the protocol logout route. `W-0171` selected `define_transport_close_handoff_gate` as the next direction.

## Maintainer Narrative

The maintainer asked the agent to replan the development route and method around the Nakama/Pitaya-class target, then continue development.

## Agent Response Summary

The agent defined a gate-only standard for future transport close handoff. The gate keeps close policy decisions application-owned, keeps concrete socket mechanics in WebSocket transport, and selects server-observed `connection_id + epoch` as the first future handoff target.

## Reference Review

Nakama informs the product expectation that logout, session lifecycle, realtime disconnect, and server-directed disconnect are related but distinct lifecycle concerns.

Pitaya informs the architecture expectation that acceptors, sessions, handlers, groups/RPC, and kick/disconnect style connection management remain separate surfaces.

vibit adapts these lessons by defining a narrow application-to-transport handoff instead of letting logout, protocol handlers, authentication service, registry state, or domain modules close sockets directly.

## Decisions

- Add `docs/transport-close-handoff-gate.md`.
- Add `docs/transport-close-handoff-gate.zh-CN.md`.
- Add `ADR-0080`.
- Add repository check rule `runtime.transport_close_handoff_gate`.
- Complete `W-0172`.
- Preserve concrete socket close implementation, close code mapping, close reason text, logout-triggered close, runtime session revocation, reconnect/epoch behavior, protocol session carriers, operations/admin disconnect, dependencies, direct Nakama/Pitaya API compatibility, and broad product modules for later gates.

## Artifacts

- `docs/transport-close-handoff-gate.md`
- `docs/transport-close-handoff-gate.zh-CN.md`
- `decisions/ADR-0080-transport-close-handoff-gate.md`
- `changes/2026-05-18-define-transport-close-handoff-gate/`
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

- Whether to implement the single-process transport close handoff next remains a W-0173 confirmation decision.
- Exact WebSocket close codes and close reason text remain unselected.
- Whether logout, runtime session revocation, admin disconnect, or duplicate connection policy should call the future handoff remains deferred.

## Follow-Up

The next likely implementation slice is a single-process transport close handoff that closes by server-observed connection id and epoch only, after explicit confirmation.

## Redaction Notes

No raw access tokens, device credentials, generated secrets, digest bytes, HMAC input bytes, verifier keys, database secrets, player private data, remote addresses, headers, cookies, query strings, subprotocol values, or GitHub tokens are recorded in this conversation log.
