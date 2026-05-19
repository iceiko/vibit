# Request

Implement `W-0111`, the access-token validation service behavior slice.

The change must execute `ValidateAccessToken` inside `runtime/internal/app/authentication/service.go` using the `ADR-0052` sequence.

It must not add protocol carriers, WebSocket handshake authentication, route protection, session persistence, logout, refresh, cleanup jobs, repository interface changes, PostgreSQL adapter changes, migrations, generated files, external dependencies, startup wiring, or broader production authentication behavior.

## Clarified Requirement

- Reject missing or malformed access-token proof before unit-of-work.
- Treat proof as opaque Base64URL unpadded 32-byte access-token material.
- Use verifier helpers to compute token lookup and verifier digests.
- Obtain authentication and player repositories through local unit-of-work capabilities.
- Find the token by lookup digest before verifier comparison.
- Require access-token kind, active token state, supported verifier posture, configured audience, non-future issue time, and non-expired time window.
- Compare token verifier digest before producing request identity.
- Require active player account state.
- Return application-owned validated player identity only after unit-of-work success.
- Keep `SessionValidated` false until session persistence is ratified.
- Collapse public invalid-token failures while preserving redacted internal failure classes.

## Acceptance Criteria

- Focused service tests cover proof rejection, helper order, repository handoff, public error collapse, token lifecycle checks, audience checks, verifier mismatch, active player account checks, request identity handoff, redaction, and no validated identity return on dependency or commit failure.
- Runtime checks recognize `W-0111` as the authorized validation service implementation slice while preserving remaining deferrals.
