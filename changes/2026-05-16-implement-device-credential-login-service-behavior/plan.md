# Plan

1. Extend `ServiceDependencies` with verifier key set, access-token random reader, clock, token record id generator, token lifetime, and token audience.
2. Implement `AuthenticateWithDeviceCredential` using the `ADR-0051` service-owned sequence.
3. Use local unit-of-work capability interfaces for authentication and player repositories.
4. Add focused tests for proof pre-validation, repository handoff, public error collapse, verifier comparison before token generation, active player account requirement, digest-only token storage, redaction, and no token return on storage or commit failure.
5. Update architecture manifests, authentication module metadata, agent guides, and repository checks.
6. Record verification.

## Non-Goals

- Access-token validation execution.
- Logout or refresh execution.
- Cleanup jobs.
- Protocol carriers.
- WebSocket handshake authentication.
- Startup wiring.
- Repository interface changes.
- PostgreSQL adapter changes.
- Migrations.
- Generated files.
- External dependencies.
