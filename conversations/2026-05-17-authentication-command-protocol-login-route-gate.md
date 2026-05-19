# Conversation: Authentication Command Protocol Login Route Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-define-authentication-command-protocol-login-route-gate/`

Related artifacts:

- `docs/authentication-command-protocol-login-route-gate.md`
- `docs/authentication-command-protocol-login-route-gate.zh-CN.md`
- `decisions/ADR-0055-authentication-command-protocol-login-route-gate.md`

## Context

The maintainer authorized continuing with the agent's recommended next direction after startup composition.

## Maintainer Narrative

The maintainer asked the agent to recommend the next ten steps and continue according to that recommendation.

## Agent Response Summary

The agent defined a gate-only standard for exposing the existing `AuthenticateWithDeviceCredential` service method through a public Protobuf command route. The gate keeps the existing envelope unchanged, keeps WebSocket transport credential-neutral, limits the first composed runtime path to PostgreSQL, and preserves session persistence and WebSocket handshake authentication deferrals.

## Reference Review

Nakama is used as the capability reference for authenticating before normal gameplay and realtime features.

Pitaya is used as the architecture reference for keeping connection acceptors, session/context, route dispatch, and handlers separated.

The gate adapts these references into vibit's own boundaries rather than copying public APIs.

## Decisions

- Define `docs/authentication-command-protocol-login-route-gate.md`.
- Record `ADR-0055`.
- Authorize a later bounded implementation slice for login command protocol messages, generated output, bridge behavior, application route registration, startup composition, and tests.
- Preserve deferrals for session persistence, WebSocket handshake authentication, logout, refresh, cleanup, token rotation, repository changes, migrations, dependencies, and direct Nakama/Pitaya compatibility.

## Artifacts

- `docs/authentication-command-protocol-login-route-gate.md`
- `docs/authentication-command-protocol-login-route-gate.zh-CN.md`
- `decisions/ADR-0055-authentication-command-protocol-login-route-gate.md`
- `changes/2026-05-17-define-authentication-command-protocol-login-route-gate/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The implementation slice still needs to add the exact `.proto` source, generated output, protocol bridge, route handler registration, and startup tests.
- Live PostgreSQL login verification remains optional until a later verification work item selects it.

## Follow-Up

- Implement the bounded public login route slice after the gate is completed and `W-0120` becomes active.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
