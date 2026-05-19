# Conversation: Transport Close Handoff Next Direction

Date: 2026-05-18
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-transport-close-handoff-gate/`

Related artifacts:

- `docs/transport-close-handoff-gate.md`
- `docs/websocket-close-policy-gate.md`
- `docs/protocol-logout-route-gate.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`

## Context

`W-0172` defined the gate-only application-to-WebSocket transport close handoff boundary. The gate selected server-observed `connection_id + epoch` as the first future concrete target and kept WebSocket transport credential-neutral and policy-neutral.

The repository still has no concrete socket close handoff. Application close policy can mark registry records invalidated, but the WebSocket transport does not yet expose a narrow in-process close target.

## Maintainer Narrative

The maintainer asked to inspect the project directory and continue advancing, with timely commit and push.

## Agent Response Summary

The agent selected `implement_transport_close_handoff_single_process` as the next lifecycle-closure direction. This follows `ADR-0080` directly: before reconnect, protocol session carriers, presence, operations/admin behavior, or broader product modules, vibit should implement the narrow single-process handoff from application-selected close intent to WebSocket transport-owned concrete socket close mechanics.

## Reference Review

Nakama remains the product reference for explicit lifecycle behavior across sessions, logout, realtime sockets, and server-directed disconnects.

Pitaya remains the architecture reference for separating acceptors, sessions, handlers, and connection-management mechanics.

vibit adapts both by keeping close decisions application-owned while letting WebSocket transport own only concrete socket mechanics for a server-observed connection id and epoch.

## Decisions

- Complete `M-101/W-0173`.
- Select `implement_transport_close_handoff_single_process`.
- Create `M-102/W-0174` as the next ready implementation work item.
- Keep close code mapping, close reason text, logout-triggered close, runtime session revocation, reconnect/epoch, protocol session carriers, presence, operations/admin disconnect, dependencies, direct compatibility, and broad product modules deferred.

## Artifacts

- `changes/2026-05-18-confirm-next-direction-after-transport-close-handoff-gate/`
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

## Open Questions

- The exact transport handoff Go type names remain for `W-0174`.
- WebSocket close code mapping and close reason text remain deferred.
- Whether logout triggers socket close remains deferred.
- Runtime session revocation remains deferred.
- Reconnect/epoch and protocol session carrier ordering remains dependent on the close handoff outcome.

## Follow-Up

Implement a single-process transport close handoff that targets only server-observed `connection_id + epoch`, remains testable without live network dependencies, and does not parse credentials or own application lifecycle policy.

## Redaction Notes

No raw access tokens, device credentials, generated secrets, digest bytes, HMAC input bytes, verifier keys, database secrets, player private data, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens are recorded in this conversation log.
