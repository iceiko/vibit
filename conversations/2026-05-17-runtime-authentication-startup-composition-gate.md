# Conversation: Runtime Authentication Startup Composition Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-define-runtime-authentication-startup-composition-gate/`

Related artifacts:

- `docs/runtime-authentication-startup-composition-gate.md`
- `docs/runtime-authentication-startup-composition-gate.zh-CN.md`
- `decisions/ADR-0054-runtime-authentication-startup-composition-gate.md`

## Context

The maintainer selected `wire_runtime_authentication_startup_composition` after request-level route protection was implemented.

## Maintainer Narrative

The maintainer asked to select `wire_runtime_authentication_startup_composition`, continue for ten steps, and strongly reference Nakama and Pitaya as mature game-server baselines.

## Agent Response Summary

The agent defined a gate-only standard before startup composition implementation. The gate authorizes only process startup composition in `runtime/cmd/vibit-server`, limited the first composed path to `VIBIT_RUNTIME_STORE=postgres`, and preserved the memory runtime path as metadata-only bootstrap behavior.

## Reference Review

Nakama is used as a capability reference for authenticate-then-use-session/token behavior before gameplay and realtime features.

Pitaya is used as an architecture vocabulary reference for keeping acceptors, sessions, routes, and handlers separated.

The gate adapts these references into vibit's own boundaries rather than copying public APIs.

## Decisions

The gate allows a later implementation to compose authentication startup only in `runtime/cmd/vibit-server` and only for the explicit PostgreSQL runtime path.

The gate preserves:

- WebSocket transport credential neutrality.
- Existing Protobuf envelope stability.
- Session persistence deferral.
- WebSocket handshake authentication deferral.
- Repository and migration stability.
- Dependency stability.
- Logout, refresh, cleanup, token rotation, and audit mutation deferrals.

## Artifacts

- `docs/runtime-authentication-startup-composition-gate.md`
- `docs/runtime-authentication-startup-composition-gate.zh-CN.md`
- `decisions/ADR-0054-runtime-authentication-startup-composition-gate.md`
- `changes/2026-05-17-define-runtime-authentication-startup-composition-gate/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The implementation slice still needs to wire existing authentication pieces into process startup after the gate.
- Later production choices for session persistence, WebSocket handshake authentication, logout, refresh, cleanup, token rotation, operations, and direct compatibility posture remain separately gated.

## Follow-Up

- Implement the bounded startup composition slice in `runtime/cmd/vibit-server`.
- Add focused tests for verifier key configuration, token lifetime/audience defaults, route-protector injection, and memory bootstrap preservation.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
