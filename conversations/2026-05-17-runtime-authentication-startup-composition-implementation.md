# Conversation: Runtime Authentication Startup Composition Implementation

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-implement-runtime-authentication-startup-composition/`

Related artifacts:

- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `.arch/work-items.yaml`

## Context

The startup composition gate authorized wiring the existing authentication service into the PostgreSQL runtime startup path.

## Maintainer Narrative

The maintainer selected startup composition as the next direction and asked for continued progress with strong Nakama and Pitaya reference discipline.

## Agent Response Summary

The agent implemented the bounded startup composition slice in `runtime/cmd/vibit-server`, wiring existing authentication service validation into request-level route protection for the PostgreSQL runtime path only. The memory runtime path remains metadata-only bootstrap behavior.

## Implementation Summary

The implementation:

- Loads verifier keys through the existing environment loader.
- Creates the existing `authentication.Service` with PostgreSQL unit-of-work, `crypto/rand.Reader`, a system clock, a token record id generator, token lifetime, and token audience.
- Creates a route access-token validator and application route protector.
- Injects the route protector into the Protobuf frame handler for the PostgreSQL runtime path.
- Preserves memory startup as metadata-only bootstrap behavior.

## Decisions

- Compose authentication startup only in `runtime/cmd/vibit-server`.
- Limit the composed runtime path to `VIBIT_RUNTIME_STORE=postgres`.
- Keep access-token proof validation request-level through the existing route protector.
- Use `VIBIT_AUTH_ACCESS_TOKEN_TTL` and `VIBIT_AUTH_TOKEN_AUDIENCE` as optional startup configuration, with safe defaults.
- Keep WebSocket transport credential-neutral and keep the existing envelope unchanged.

## Reference Notes

Nakama informed the authenticate-then-use-session/token capability expectation.

Pitaya informed the separation of connection handling, routing, session/context, and handler behavior.

The implementation adapts those patterns without adopting direct public API compatibility.

## Artifacts

- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `changes/2026-05-17-implement-runtime-authentication-startup-composition/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `runtime/AGENTS.md`
- `modules/authentication/AGENTS.md`

## Deferrals Preserved

- WebSocket handshake authentication.
- Session persistence.
- Authentication command routes.
- Repository interface changes.
- PostgreSQL adapter changes.
- Migrations.
- Generated files.
- Dependencies.
- Logout, refresh, cleanup, token rotation, and token validation audit mutation.

## Open Questions

- The next milestone direction remains blocked pending maintainer selection.
- Session persistence, handshake authentication, authentication command route exposure, logout, refresh, cleanup, token rotation, and operations posture remain future decisions.

## Follow-Up

- Run full runtime tests, repository checks, memory checks, change checks, and diff checks.
- Record any remaining warnings or unavailable live verification explicitly.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
