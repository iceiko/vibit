# Conversation: Protocol Logout Route Gate

Date: 2026-05-18
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-18-define-protocol-logout-route-gate/`

Related artifacts:

- `docs/protocol-logout-route-gate.md`
- `docs/protocol-logout-route-gate.zh-CN.md`
- `decisions/ADR-0079-protocol-logout-route-gate.md`

## Context

The maintainer clarified that vibit should become a Nakama/Pitaya-class product and cover common game backend functionality. The previous roadmap slice ratified product parity and selected runtime lifecycle closure as the near-term path.

## Maintainer Narrative

The maintainer asked to replan the project development route and method around that product target, then continue.

## Agent Response Summary

The agent continued with the roadmap's first concrete lifecycle step: defining the protocol logout route gate. The gate prepares a future implementation of `runtime.authentication.LogoutAccessToken` as a WebSocket Protobuf command route while preserving token-service ownership, WebSocket transport neutrality, existing Protobuf envelope shape, and socket/session lifecycle deferrals.

## Reference Review

Nakama is used as the product reference for explicit session lifecycle behavior. Its server runtime exposes session logout separately from session disconnect, so vibit should not make token logout implicitly close sockets.

Pitaya is used as the architecture reference for separating acceptors, sessions, route handlers, and connection management. vibit adapts that by keeping logout proof validation in the application authentication service and keeping concrete socket close behind later policy and handoff gates.

## Decisions

- Add `docs/protocol-logout-route-gate.md`.
- Add `docs/protocol-logout-route-gate.zh-CN.md`.
- Record `ADR-0079`.
- Add repository check rule `runtime.protocol_logout_route_gate`.
- Complete `W-0169`.
- Create `M-098/W-0170` as the next ready implementation slice for the protocol logout route.

## Artifacts

- `docs/protocol-logout-route-gate.md`
- `docs/protocol-logout-route-gate.zh-CN.md`
- `decisions/ADR-0079-protocol-logout-route-gate.md`
- `changes/2026-05-18-define-protocol-logout-route-gate/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The implementation slice still needs to add the exact Protobuf fields, generated output, bridge mapping, route handler registration, transaction bypass, startup wiring, and tests.
- The future transport close handoff still needs its own gate.
- Bound connection/session behavior after logout remains deferred.

## Follow-Up

- Implement `W-0170`, the bounded protocol logout route implementation slice.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
