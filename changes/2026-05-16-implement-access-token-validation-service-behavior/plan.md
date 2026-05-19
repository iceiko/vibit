# Plan

1. Implement `ValidateAccessToken` using the `ADR-0052` service-owned sequence.
2. Decode access-token proof with the same opaque material constraints used by issued access tokens.
3. Use local unit-of-work capability interfaces for authentication and player repositories.
4. Add token posture checks for kind, status, verifier algorithm, verifier version, key id, audience, issue time, and expiration.
5. Compare token verifier digest before player account lookup and request identity construction.
6. Return validated player identity only after unit-of-work success, with `SessionValidated` false.
7. Add focused tests for proof rejection, repository handoff, public error collapse, lifecycle checks, verifier mismatch, player activity, redaction, and no identity return on dependency or commit failure.
8. Update architecture manifests, authentication module metadata, agent guides, and repository checks.
9. Record verification.

## Non-Goals

- Protocol carriers.
- WebSocket handshake authentication.
- Route protection.
- Session persistence.
- Logout or refresh execution.
- Cleanup jobs or token rotation.
- Repository interface changes.
- PostgreSQL adapter changes.
- Migrations.
- Generated files.
- External dependencies.
- Startup wiring.
