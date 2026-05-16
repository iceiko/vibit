# Request

Implement `W-0109`, the device credential login service behavior slice.

The change must execute `AuthenticateWithDeviceCredential` inside `runtime/internal/app/authentication/service.go` using the `ADR-0051` sequence.

It must not implement access-token validation, logout, refresh, cleanup jobs, protocol carriers, WebSocket handshake authentication, startup wiring, repository interface changes, migrations, generated files, external dependencies, or broader production authentication behavior.

## Clarified Requirement

- Reject missing or malformed device credential proof before unit-of-work.
- Treat proof as server-issued Base64URL unpadded 32-byte material.
- Use verifier helpers to compute lookup and verifier digests.
- Obtain authentication and player repositories through local unit-of-work capabilities.
- Require active credential and active player account state.
- Generate opaque access-token material only after proof and player acceptance.
- Store token lookup and verifier digests only.
- Return raw access-token text only after token storage and unit-of-work success.

## Acceptance Criteria

- Focused service tests cover proof rejection, helper order, repository handoff, public error collapse, player account checks, digest-only token storage, redaction, and no token return on store or commit failure.
- Runtime checks recognize `W-0109` as the authorized service implementation slice while preserving remaining deferrals.
