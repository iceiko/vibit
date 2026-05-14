# Request

## Original Request

Continue the next work item in the queue.

## Clarified Requirement

Advance `W-0084 Implement authentication PostgreSQL adapter` by implementing only the authentication PostgreSQL repository adapter inside the PostgreSQL platform package, adding focused fake-executor tests, updating architecture metadata, and preserving all runtime authentication, login, token generation, token validation, logout execution, refresh, cleanup job, protocol, WebSocket, generated output, and dependency deferrals.

## User-Visible Outcome

Agents can now see a bounded authentication PostgreSQL adapter implementation at `runtime/internal/platform/persistence/postgres/authentication_repository.go`, focused tests at `runtime/internal/platform/persistence/postgres/authentication_repository_test.go`, and a PostgreSQL package helper `UnitOfWork.NewAuthenticationRepository`.

The adapter remains internal persistence infrastructure. It is not exposed through runtime authentication handlers, token validators, login routes, WebSocket routes, Protobuf messages, or generated authentication shapes.

## Non-Goals

- Do not add runtime login handlers.
- Do not generate credential material or access tokens.
- Do not compare credential or token verifiers.
- Do not validate access tokens.
- Do not implement logout execution, refresh behavior, or cleanup jobs.
- Do not add WebSocket proof carriers, WebSocket routes, or handshake authentication.
- Do not add Protobuf authentication messages or generated authentication shapes.
- Do not change the ratified authentication migration schemas.
- Do not change the module-owned `authentication.Repository` interface.
- Do not make live PostgreSQL verification mandatory for default checks.

## Unknowns

- Runtime authentication behavior remains deferred to a later bounded milestone.
- Exact credential and token verifier algorithms, secret configuration, and comparison behavior remain deferred.
- Runtime cleanup scheduling or explicit maintenance command behavior remains deferred.

## Acceptance Criteria

- [x] The adapter implements `authentication.Repository` in the PostgreSQL platform package.
- [x] Credential store and lookup use `authentication_device_credentials` through a caller-supplied executor.
- [x] Token store, lookup, revocation, and cleanup eligibility reads use `authentication_access_tokens` through a caller-supplied executor.
- [x] The adapter does not call `BEGIN`, `COMMIT`, or `ROLLBACK`.
- [x] Focused fake-executor tests cover credential storage/lookup, token storage/lookup, revocation, cleanup eligibility, error mapping, UTC timestamp mapping, and no-live-PostgreSQL default behavior.
- [x] Architecture metadata and bilingual guidance record that the adapter is implemented while runtime authentication behavior remains deferred.
