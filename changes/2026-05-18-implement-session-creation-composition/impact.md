# Impact Analysis

## Affected Modules

- `runtime`: The startup authentication composition provides a runtime session id generator and the PostgreSQL unit-of-work already exposes `NewSessionRepository`.
- `authentication`: The application service composes device-credential login, token storage, and runtime session creation.

## Module Ownership Impact

`runtime/internal/app/authentication` owns the orchestration because login is the first selected token issuance entry point. `runtime/internal/app/session` still owns session repository contracts, and `runtime/internal/platform/persistence/postgres` still owns SQL persistence.

## Public Contract Impact

No Protobuf source, generated Protobuf output, command schema, query schema, event schema, permission, or error contract changes. `AuthenticationResult` grows application-owned session metadata, but protocol response mapping remains unchanged.

## Data And Migration Impact

No migration is added. The implementation uses the existing `runtime_sessions` table through the existing repository interface and PostgreSQL adapter.

## Test Impact

Focused Go tests cover:

- Successful login creates an active runtime session linked to the stored token record.
- Session repository acquisition failures are redacted and do not return token/session material.
- Session id generation failures are redacted and do not create sessions.
- Session creation failures are redacted and do not return token/session material.
- Commit failures do not return token/session material.
- Access-token validation remains unchanged with `SessionValidated` false.
- Startup random session id shape is explicit.

## Documentation Impact

ADR, conversation memory, work queue, manifests, module guides, rule catalog, and checks are updated for the implementation slice.

## Compatibility Risks

The main compatibility risk is that persistent PostgreSQL login now requires a functioning `runtime_sessions` table and session repository capability. That is intentional for this bounded slice because durable session creation is now part of successful production login composition. Protocol clients do not see a new response field in this slice.
