# authentication Module Agent Guide

Status: Draft v0.1

## When To Use This Module

Use this module for requirements that define storage-neutral authentication persistence boundaries for the selected first posture:

- `device_credential_login`
- opaque high-entropy access tokens
- PostgreSQL-backed credential and token verifier records

The current module-owned runtime boundary is:

```text
runtime/internal/modules/authentication/repository.go
```

It may define credential record structs, token verifier record structs, mutation/query shapes, repository interfaces, and validation helpers needed by future persistence adapters.

## When Not To Use This Module

Do not use this module to implement:

- Runtime login handlers.
- Token generation, parsing, issuance, validation, refresh, or bearer-token acceptance.
- Logout behavior beyond storage-neutral revocation mutation shapes.
- Cleanup jobs beyond storage-neutral cleanup query shapes.
- PostgreSQL adapters outside the separately authorized `M-015` platform boundary.
- WebSocket routes, proof carriers, or handshake authentication.
- Protobuf messages or generated authentication shapes.
- Password hashing, JWT, OAuth, OIDC, provider SDKs, key-management, Redis-like token stores, S3, or MinIO dependencies.
- Player account lifecycle storage.

If a requirement needs one of those surfaces, create or update a separate ratified change before adding code.

## Extension Points

- Authentication repository interface: `runtime/internal/modules/authentication/repository.go`.
- Authentication repository boundary tests: `runtime/internal/modules/authentication/repository_test.go`.
- Authentication PostgreSQL adapter source: `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- Authentication PostgreSQL adapter tests: `runtime/internal/platform/persistence/postgres/authentication_repository_test.go`.
- Existing credential migration source: `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`.
- Existing token verifier migration source: `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`.

The repository interface is storage-neutral domain code. It must not import PostgreSQL, WebSocket, Protobuf, OAuth, OIDC, JWT, password-hashing, Redis-like, S3, or MinIO packages.

The repository interface may normalize identifiers, digest byte slices, statuses, timestamps, and storage mutation/query shapes. It must not create credential material, generate tokens, compare verifiers, parse bearer tokens, validate access tokens, open transactions, execute SQL, publish events, or call transport/protocol code.

The implemented PostgreSQL adapter boundary is `runtime/internal/platform/persistence/postgres/authentication_repository.go`, with focused tests at `runtime/internal/platform/persistence/postgres/authentication_repository_test.go`. It uses `NewAuthenticationRepositoryForUnitOfWork(executor)` and implements `authentication.Repository`; `UnitOfWork.NewAuthenticationRepository` creates it from the caller-owned executor. `M-015` authorizes only platform-owned persistence adapter work; it still does not authorize runtime authentication behavior.

## Forbidden Shortcuts

- Do not store raw credential or token material.
- Do not add `AuthService`, `Authenticator`, `TokenValidator`, `TokenIssuer`, `TokenVerifier`, `TokenRepository`, or `CredentialRepository` implementation types.
- Do not add PostgreSQL adapter files under this module.
- Do not import `pgx`, `goose`, WebSocket, Protobuf, JWT, OAuth, OIDC, bcrypt, argon2, provider SDKs, Redis-like clients, S3 SDKs, or MinIO clients.
- Do not add runtime authentication behavior while editing the repository interface boundary.
- Do not change the ratified credential or token verifier migration schemas without a separate change spec and decision.
- Do not make player account lifecycle tables own authentication state.

## Required Tests

See `tests.required` in `module.yaml`.

For the current module state, focused tests are limited to repository interface shape, closed status sets, required-field normalization, digest copying, timestamp UTC normalization, and storage-neutral mutation/query validation. PostgreSQL adapter tests belong under `runtime/internal/platform/persistence/postgres/` during `M-015`; runtime authentication, WebSocket, Protobuf, and generated-shape tests become required only after their separate boundaries are ratified.
