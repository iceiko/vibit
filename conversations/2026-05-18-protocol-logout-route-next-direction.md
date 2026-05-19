# Conversation: Protocol Logout Route Next Direction

Date: 2026-05-18
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-protocol-logout-route-implementation/`

Related artifacts:

- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/protocol-logout-route-gate.md`
- `docs/websocket-close-policy-gate.md`
- `.arch/work-items.yaml`

## Context

`W-0170` exposed `runtime.authentication.LogoutAccessToken` through the WebSocket Protobuf protocol route while keeping socket close, runtime session revocation, reconnect, and protocol session carriers deferred.

The product roadmap requires vibit to become a Nakama/Pitaya-class game backend product. The near-term phase remains lifecycle closure before presence, chat, social, matchmaking, match runtime, SDK, or distributed runtime expansion.

## Maintainer Narrative

The maintainer asked to replan the project development route and method around the Nakama/Pitaya-class target, then continue.

## Agent Response Summary

The agent selected `define_transport_close_handoff_gate` as the next direction. This is the next lifecycle step after protocol logout route exposure because application close policy can currently invalidate registry records and produce close intents, but WebSocket transport cannot yet receive a narrow concrete close handoff.

## Reference Review

Nakama is used as the product reference for explicit lifecycle behavior across authentication, sessions, realtime socket disconnect, and operator actions.

Pitaya is used as the architecture reference for separating acceptor/transport mechanics, sessions, handlers, and kick/disconnect style connection management.

vibit adapts both by keeping close decisions application-owned and defining a future narrow transport handoff before any concrete socket close implementation.

## Decisions

- Complete `M-099/W-0171`.
- Select `define_transport_close_handoff_gate`.
- Create `M-100/W-0172` as the next ready gate-only work item.
- Keep implementation, close codes, close reason text, logout-triggered close, runtime session revocation, reconnect/epoch, protocol session carriers, presence, operations/admin disconnect, dependencies, direct compatibility, and broad product modules deferred.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-protocol-logout-route-implementation/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`

## Open Questions

- Whether to implement the transport close handoff immediately after the gate remains the next confirmation choice.
- Close code mapping and close reason text remain deferred.
- Reconnect/epoch and protocol session carrier ordering remains dependent on the close handoff outcome.

## Follow-Up

Define `docs/transport-close-handoff-gate.md` and its paired Simplified Chinese translation, then record the corresponding ADR and repository check rule before implementing any WebSocket socket close handoff.

## Redaction Notes

No raw access tokens, device credentials, generated secrets, digest bytes, HMAC input bytes, verifier keys, database secrets, player private data, or GitHub tokens are recorded in this conversation log.
