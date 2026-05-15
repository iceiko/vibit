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
- Protobuf messages or handwritten authentication behavior under generated authentication shape paths.
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

Runtime authentication implementation boundary planning is documented in `docs/runtime-authentication-implementation-boundary.md` and `decisions/ADR-0036-runtime-authentication-implementation-boundary.md`. Future runtime authentication must be application-owned under `runtime/internal/app`, must use this module's `authentication.Repository` through the application unit-of-work boundary, and must convert validated proof into `RequestIdentity` before domain dispatch. This module must not absorb token generation, verifier comparison, login execution, access-token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, handwritten logic under generated authentication shape paths, or authentication dependencies.

`M-019 Token And Credential Verifier Algorithm Redaction Boundary` is completed. `W-0091` defined the first planned high-entropy verifier posture in `docs/token-credential-verifier-algorithm-redaction-boundary.md` and `ADR-0040` without adding service code or runtime authentication behavior. The earlier application authentication service boundary remains documented in `docs/application-authentication-service-interface-boundary.md` and `ADR-0039`. The first planned verifier algorithm family is `vibit_hmac_sha256_v1`; future first-posture code may use Go standard library packages `crypto/hmac`, `crypto/sha256`, `crypto/subtle`, `crypto/rand`, and `encoding/base64` after a later code gate. This module's `authentication.Repository` may store and retrieve already-computed digest material only; it must not compute HMACs, generate token or credential material, compare verifiers, load secret keys, or decide authentication outcomes.

`M-020 Secret Configuration And Verifier Key Loading Boundary` is completed. `W-0092` defined future key loading posture in `docs/secret-configuration-verifier-key-loading-boundary.md` and `ADR-0041` without adding service code or runtime authentication behavior. Future verifier key loading is application-owned under `runtime/internal/app`; the first local implementation may use process environment configuration or explicit runtime secret input after a later code gate; external KMS or secret-manager integration remains behind dependency and operations gates. Four separated logical verifier keys are required, `verifier_key_id` is not log-safe by default, production key configuration must fail closed when invalid, and committed production-like secret values remain forbidden. This module must not load verifier keys, parse environment variables, choose secret managers, rotate keys, compute HMACs, generate token or credential material, compare verifiers, or decide authentication outcomes.

`M-021 Token And Credential Material Generation Boundary` is completed. `W-0093` defined future raw device credential and opaque access-token material generation posture in `docs/token-credential-material-generation-boundary.md` and `ADR-0042` without adding service code or runtime authentication behavior. Future material generation is application-owned under `runtime/internal/app`; the first device credential and access token are server-issued and application-generated; raw material must be 32 cryptographically random bytes with at least 256 bits of entropy; text presentation is URL-safe unpadded Base64 or equivalent; raw material is one-time client-visible and must not be stored. This module's `authentication.Repository` may store and retrieve already-computed digest material only; it must not generate raw token or credential material, encode raw material, accept raw material for storage, compute digests, compare verifiers, or decide authentication outcomes. `runtime.token_credential_material_generation_boundary` is the repository check rule for this boundary.

`M-022 Verifier Digest Computation And Comparison Boundary` is completed. `W-0094` defined future verifier digest computation and constant-time comparison posture in `docs/verifier-digest-computation-comparison-boundary.md` and `ADR-0043` without adding service code or runtime authentication behavior. Future digest computation and comparison is application-owned under `runtime/internal/app`; lookup digest equality may select a candidate record only; verifier digest comparison must be constant-time; invalid lookup, mismatch, unknown key id, unsupported algorithm, malformed proof, and expired or revoked proof must collapse to the same public invalid-proof class. This module's `authentication.Repository` may store and retrieve already-computed digest material only; it must not compute HMACs, choose verifier key sets, compare verifier digests, disclose lookup misses, or decide authentication outcomes. `runtime.verifier_digest_computation_comparison_boundary` is the repository check rule for this boundary.

`M-023 Authentication Service Implementation Readiness Gate` is completed. `W-0095` defined the readiness gate in `docs/authentication-service-implementation-readiness-gate.md` and `ADR-0044` without adding service code or runtime authentication behavior. Future service implementation remains application-owned under `runtime/internal/app`, with `runtime/internal/app/authentication` as the package candidate. This module remains a storage-neutral repository boundary. It must not absorb service orchestration, secret loading, material generation, digest computation, verifier comparison, login execution, token validation, protocol behavior, WebSocket behavior, or production authentication decisions. `runtime.authentication_service_implementation_readiness_gate` is the repository check rule for this gate.

`M-024 Local Verifier Key Configuration Loading Gate` is completed. `W-0096` defined the gate in `docs/local-verifier-key-configuration-loading-gate.md` and `ADR-0045` without adding service code or runtime authentication behavior. `W-0097` implemented explicit in-memory verifier key set validation under `runtime/internal/app/authentication`, not this module. This module remains storage-neutral and may store already-computed verifier metadata only. It must not load verifier keys, parse environment variables, decode key text, hold key material, validate key sets, rotate keys, compute digests, compare verifiers, or decide authentication outcomes. `runtime.local_verifier_key_configuration_loading_gate` is the repository check rule for this gate.

`M-026 Environment Verifier Key Loader Gate` is completed. `W-0098` defined the gate in `docs/environment-verifier-key-loader-gate.md` and `ADR-0046` without adding service code or runtime authentication behavior. Future process environment verifier key loading belongs under `runtime/internal/app/authentication`, not this module, and must call `NewVerifierKeySet` there. This module remains storage-neutral and may store already-computed verifier metadata only. It must not parse environment variables, decode key text, hold key material, validate key sets, read local secret files, parse `.env` files, rotate keys, compute digests, compare verifiers, or decide authentication outcomes. `runtime.environment_verifier_key_loader_gate` is the repository check rule for this gate.

## Forbidden Shortcuts

- Do not store raw credential or token material.
- Do not compute HMACs, generate verifier digests, load verifier keys, or compare verifier material inside this module.
- Do not add `AuthService`, `Authenticator`, `TokenValidator`, `TokenIssuer`, `TokenVerifier`, `TokenRepository`, or `CredentialRepository` implementation types.
- Do not add PostgreSQL adapter files under this module.
- Do not import `pgx`, `goose`, WebSocket, Protobuf, JWT, OAuth, OIDC, bcrypt, argon2, provider SDKs, Redis-like clients, S3 SDKs, or MinIO clients.
- Do not add runtime authentication behavior while editing the repository interface boundary.
- Do not change the ratified credential or token verifier migration schemas without a separate change spec and decision.
- Do not make player account lifecycle tables own authentication state.

## Required Tests

See `tests.required` in `module.yaml`.

For the current module state, focused tests are limited to repository interface shape, closed status sets, required-field normalization, digest copying, timestamp UTC normalization, and storage-neutral mutation/query validation. PostgreSQL adapter tests belong under `runtime/internal/platform/persistence/postgres/` during `M-015`; runtime authentication, WebSocket, Protobuf, and generated-shape tests become required only after their separate boundaries are ratified.
