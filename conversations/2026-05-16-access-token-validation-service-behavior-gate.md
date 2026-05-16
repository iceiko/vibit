# Conversation: Access Token Validation Service Behavior Gate

Date: 2026-05-16
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-16-define-access-token-validation-service-behavior-gate/`

Related artifacts:

- `docs/access-token-validation-service-behavior-gate.md`
- `docs/access-token-validation-service-behavior-gate.zh-CN.md`
- `decisions/ADR-0052-access-token-validation-service-behavior-gate.md`
- `.arch/work-items.yaml`

## Context

`W-0109` implemented device credential login and left `ValidateAccessToken` fail-closed. The next ready work item was `W-0110`, defining the bounded gate for future access-token validation behavior before any real validation execution is added.

## Maintainer Narrative

The maintainer asked the agent to continue.

The standing workflow preference remains: `继续` advances one `next_ready` work item unless a real maintainer decision is required. Routine technical implementation details inside an already ratified direction should be handled by the agent.

## Agent Response Summary

The agent advanced `W-0110` as a gate-only standard.

The change defines the future access-token validation sequence, proof input shape, dependency reuse, repository handoff, token lifecycle and audience checks, request identity handoff, public error collapse, redaction requirements, and required tests. It keeps `ValidateAccessToken` fail-closed and does not execute validation, expose protocol carriers, wire startup, change repositories or migrations, add dependencies, add session persistence, or add production authentication behavior.

## Decisions

- Treat first access-token proof as already-decoded, high-entropy, Base64URL unpadded opaque token material, not a bearer header, session id, route field, or transport metadata.
- Require proof pre-validation before any unit-of-work or repository call.
- Use existing authentication and player repositories through unit-of-work capabilities instead of changing repository interfaces in this gate.
- Require token lookup, lifecycle checks, audience check, verifier digest comparison, and active player account state before producing validated request identity.
- Keep `SessionValidated` false until a future session persistence gate.
- Collapse lookup miss, wrong token posture, expired token, revoked token, wrong audience, verifier mismatch, and inactive player to the same public invalid-token family for the first behavior.

## Artifacts

- `docs/access-token-validation-service-behavior-gate.md`
- `docs/access-token-validation-service-behavior-gate.zh-CN.md`
- `decisions/ADR-0052-access-token-validation-service-behavior-gate.md`
- `changes/2026-05-16-define-access-token-validation-service-behavior-gate/`

## Open Questions

- None for this gate-only change.

## Follow-Up

- The next work item may implement access-token validation inside `runtime/internal/app/authentication/service.go` only if it preserves the gate sequence and deferrals.
- Protocol carriers, WebSocket handshake authentication, route protection, session persistence, logout, refresh, cleanup, startup wiring, repository changes, migrations, dependencies, and production behavior remain separate work.

## Redaction Notes

No secrets, tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, account details, or private data are recorded in this conversation log.
