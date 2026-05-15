# Go Runtime Agent Guide

Status: Draft v0.1
Last updated: 2026-05-13
Scope: `runtime/` Go server runtime workspace
Canonical source: `../CONSTITUTION.md`, `../AGENTS.md`, and `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`

This guide applies to the first Go server runtime implementation.

The paired Simplified Chinese translation is `runtime/AGENTS.zh-CN.md`. The English file is authoritative.

## 1. Purpose

`runtime/` is the first Go module for vibit's server runtime.

The Go module path is:

```text
github.com/iceiko/vibit/runtime
```

The runtime exists to prove vibit's core claim through a small, long-lived backend slice:

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

Do not treat this workspace as a disposable demo.

## 2. Required Reading

Before changing files under `runtime/`, read:

- `../CONSTITUTION.md`
- `../AGENTS.md`
- `../.arch/runtime.yaml`
- `../.arch/dependencies.yaml`
- `../.arch/contracts.yaml`
- `../docs/generated-output.md`
- `../docs/runtime-protocol-adapter.md`
- `../docs/postgresql-persistence-boundary.md`, before persistence work
- `../docs/postgresql-verification-environment.md`, before live PostgreSQL verification work
- `../docs/authentication-token-session-validation.md`, before authentication, token, credential, external identity, session persistence, request identity trust, WebSocket handshake, player handler, or player route work
- `../docs/authentication-proof-token-session-contract-dimensions.md`, before authentication proof, token/session validation, session error, session permission, or validation event contract work
- `../docs/runtime-authentication-implementation-boundary.md`, before runtime authentication implementation planning or code
- `../docs/runtime-runbook.md`
- `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
- `../decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
- `../decisions/ADR-0020-postgresql-persistence-boundary.md`, before persistence work
- `../decisions/ADR-0022-player-account-postgresql-schema-boundary.md`, before player account persistence work
- `../decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`, before authentication/session design or implementation work
- `../decisions/ADR-0036-runtime-authentication-implementation-boundary.md`, before runtime authentication implementation planning or code
- The affected module manifest, such as `../modules/inventory/module.yaml`
- The relevant change spec under `../changes/`

## 3. Package Ownership

Use these package boundaries:

- `cmd/vibit-server/`: process startup, configuration wiring, and lifecycle.
- `internal/app/`: command/query dispatch, application composition, and transaction orchestration.
- `internal/platform/transport/ws/`: WebSocket transport adapter and `github.com/coder/websocket` ownership.
- `internal/platform/protocol/protobuf/`: Protobuf framing, envelope conversion, and wire message adaptation.
- `internal/platform/persistence/postgres/`: PostgreSQL adapter implementation and `github.com/jackc/pgx/v5` ownership.
- `internal/platform/migrations/`: migration tooling invocation and validation.
- `internal/platform/events/`: event recording and publication mechanisms.
- `internal/platform/tx/`: transaction boundary and unit-of-work interfaces.
- `internal/modules/<module>/`: handwritten domain module runtime logic.
- `internal/generated/contracts/`: generated Go contract shapes.
- `internal/generated/proto/`: generated Go Protobuf files.
- `migrations/postgres/`: SQL-first PostgreSQL migration sources.

## 4. Dependency Rules

Domain modules must not import third-party transport, protocol, persistence, migration, object-storage, or framework dependencies directly.

Allowed owner packages:

- `github.com/coder/websocket`: `internal/platform/transport/ws/` only.
- `google.golang.org/protobuf`: generated protocol packages and protocol adapter packages only.
- `github.com/jackc/pgx/v5`: `internal/platform/persistence/postgres/` only.
- `github.com/pressly/goose/v3`: `internal/platform/migrations/` only.

Do not add new foundational dependencies without checking `../.arch/dependencies.yaml` and creating the required adoption record.

## 5. Runtime Boundary Rules

Runtime protocol handoff must follow `../docs/runtime-protocol-adapter.md`.

WebSocket transport reads and writes frames. Protobuf protocol adaptation decodes and encodes envelopes. Application dispatch routes commands and queries. Domain modules enforce invariants. Generated packages provide shapes only.

WebSocket transport handlers pass opaque frame bytes to injected protocol/application composition. They do not adapt requests into commands or queries directly, and they must not hide business logic.

State-changing commands should enter through `internal/app/` and run inside an application-owned unit of work. Repository mutations and domain event recording should happen inside that same unit of work.

The current transaction skeleton is `internal/platform/tx.Runner`, `internal/platform/tx.UnitOfWork`, and `internal/app.TransactionalDispatcher`. Application code may import this transaction boundary package, but it must not import persistence, migration, protocol, or transport platform adapters. Query routes should pass through without a write unit of work by default.

Query handlers should not mutate state and do not require a write transaction by default.

Event publication outside the transaction remains deferred until vibit adopts an explicit event delivery or outbox standard.

PostgreSQL persistence work must follow `../docs/postgresql-persistence-boundary.md`. Repository interfaces stay module-owned, `pgx` stays under `internal/platform/persistence/postgres/`, `goose` stays under `internal/platform/migrations/`, and SQL migration sources stay under `migrations/postgres/`.

For the first durable inventory implementation, `GrantItem` must use a transaction-bound repository and call `LockInventoryForMutation` before reading current items and applying capacity-sensitive mutations. The returned `MutationLock` is a locked aggregate view, not a transaction owner. Repositories must not silently open independent write transactions for command flows.

The first PostgreSQL inventory repository adapter is `internal/platform/persistence/postgres/inventory_repository.go`. Construct it with `NewInventoryRepositoryForUnitOfWork` and pass an executor supplied by the application-owned unit of work, such as a `pgx.Tx` or compatible test executor. The adapter must not call `BEGIN`, `COMMIT`, or `ROLLBACK`; transaction lifetime belongs to `internal/platform/tx` and `internal/app`.

PostgreSQL configuration is owned by `internal/platform/persistence/postgres/config.go`. It reads `VIBIT_POSTGRES_DSN`, `VIBIT_POSTGRES_MAX_CONNS`, and `VIBIT_POSTGRES_MIN_CONNS`, builds pgx pool configuration, and must not require a live PostgreSQL server in normal unit tests. Connection strings and credentials must come from environment or explicit runtime input and must not be committed.

The pgx-backed transaction runner is `internal/platform/persistence/postgres/runner.go`. It implements `internal/platform/tx.Runner` while keeping pgx transaction handles inside the PostgreSQL platform package. It commits successful command units of work, rolls back failed callback units of work, and exposes package-owned helpers such as `UnitOfWork.NewInventoryRepository` for future persistent composition. Do not import the PostgreSQL runner from `internal/app/` or domain modules; persistent runtime wiring must happen in an approved composition boundary.

`GrantItemMutation` carries `event_id`, `occurred_at`, and `reason` so the PostgreSQL adapter can persist `inventory_item_grants` in the same executor path as the item quantity update.

The first inventory migration source is `migrations/postgres/000001_create_inventory_state.sql`. It creates `inventory_accounts`, `inventory_items`, and `inventory_item_grants`. Run `node ../tools/vibit check migrations` when migration sources or migration guidance change. Migration status and apply behavior are covered by the opt-in live durable inventory request-loop verification when `VIBIT_POSTGRES_TEST_DSN` is set.

The ratified player account PostgreSQL schema boundary is documented in `../docs/postgresql-persistence-boundary.md` and `../decisions/ADR-0022-player-account-postgresql-schema-boundary.md`. The first player account migration source is `migrations/postgres/000002_create_player_account_state.sql`. That migration creates only `player_accounts` and `player_account_events` lifecycle state. It must not add credentials, password hashes, external identity links, access tokens, refresh tokens, runtime session rows, WebSocket connection state, request identity validation results, inventory state, or permission grants.

The player account repository interface boundary is `internal/modules/player/repository.go`. It is storage-neutral domain code and may define account lifecycle structs, `Repository.CreatePlayerAccount`, `Repository.GetPlayerAccount`, and durable mutation metadata for persistence adapters. The PostgreSQL adapter is `internal/platform/persistence/postgres/player_account_repository.go`, with focused tests at `internal/platform/persistence/postgres/player_account_repository_test.go`. It uses `NewPlayerAccountRepositoryForUnitOfWork(executor)`, implements `player.Repository`, receives its executor from the application-owned unit of work, and must not call `BEGIN`, `COMMIT`, or `ROLLBACK`. `UnitOfWork.NewPlayerAccountRepository` is a PostgreSQL package helper and must not expose pgx to application or domain packages.

The player account PostgreSQL adapter does not authorize runtime handlers, WebSocket routes, authentication, token behavior, credential storage, external identity linking, or session persistence. The adapter may only write `player_accounts`, write `player_account_events` for `PlayerAccountCreated`, and read current lifecycle rows from `player_accounts` until a later change ratifies more behavior.

The authentication, token, and session validation design boundary is documented in `../docs/authentication-token-session-validation.md` and `../decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`. It separates authentication proof, login methods, tokens, credentials, external identity links, runtime sessions, request identity, WebSocket handshake authentication, player account lifecycle, transport connection metadata, and Protobuf envelope metadata. The current `MetadataOnlySessionValidator` is a non-authenticated bootstrap path. Do not treat metadata-only `player_id`, `session_id`, or `connection_id` as production proof, and do not add authentication runtime code, token parsing, credential lookup, external identity linking, session persistence, Protobuf envelope authentication changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes until separately ratified. `runtime.authentication_token_session_boundary` is the repository check rule for this boundary.

Authentication proof and token/session contract dimensions are documented in `../docs/authentication-proof-token-session-contract-dimensions.md`. Use that standard for actor kinds, validation statuses, proof statuses, failure classes, retryability, request identity handoff, session error metadata, session permission metadata, and validation event metadata. These dimensions are semantic vocabulary only and do not grant permission to implement login methods, token formats, credential lookup, session persistence, Protobuf envelope changes, WebSocket handshake changes, runtime player handlers, or WebSocket routes.

Credential storage and external identity linking boundaries are documented in `../docs/credential-storage-external-identity-linking-boundaries.md`. Use that standard before adding credential storage, external identity linking, login methods, provider subjects, password hashing, OAuth, OIDC, provider SDKs, account linking, recovery flows, merge behavior, or related schema. The boundary preserves `player_accounts` and `player_account_events` as lifecycle-only tables and does not authorize credential tables, external identity tables, provider dependencies, runtime lookup code, player lifecycle table changes, or direct Nakama/Pitaya API compatibility.

Session persistence and WebSocket handshake decision gates are documented in `../docs/session-persistence-websocket-handshake-decision-gates.md`. Use that standard before adding session persistence, WebSocket handshake authentication, reconnect behavior, connection epoch behavior, token/session carriers, session-related Protobuf envelope changes, handshake/system messages, or route-level authentication. It does not select request-level, first-message, handshake-level, every-request, or hybrid validation as the production model. It does not authorize session tables, a session store, envelope changes, or handshake authentication behavior.

Login method and token format ratification is documented in `../docs/login-method-token-format-ratification.md` and `../decisions/ADR-0024-login-method-token-format-ratification-boundary.md`. Use that standard before selecting first login methods, token model, token format, proof carrier posture, token lifecycle semantics, credential/token/session schema gates, or implementation queue. It does not authorize runtime authentication, token parsing, credential lookup, external identity linking, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes.

Selected login/token boundary checks are documented in `../docs/selected-login-token-boundary-checks.md` and `../decisions/ADR-0030-selected-login-token-boundary-checks.md`. The selected first posture is `device_credential_login` with opaque high-entropy access tokens, login-command token issuance, explicit request proof payloads, no refresh token in the first implementation, and PostgreSQL schema gates for credential and token verifier records. This posture remains runtime-behavior-deferred after the ratified credential and token verifier migration sources and metadata-only generated authentication contract shapes. Do not add token validators, token issuers, credential repositories, token repositories, authentication Protobuf sources, additional authentication migrations, WebSocket proof carriers, or authentication dependencies until a later bounded implementation milestone authorizes them. `runtime.selected_login_token_boundary` is the repository check rule for this selected posture.

Credential record schema boundary is documented in `../docs/credential-record-schema-boundary.md` and `../decisions/ADR-0032-credential-record-schema-boundary.md`. The credential migration source now exists at `migrations/postgres/000003_create_authentication_device_credentials.sql` for the ratified `authentication_device_credentials` semantics. Runtime code must remain unchanged until later bounded work explicitly authorizes repositories, adapters, runtime lookup, handlers, routes, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

Token verifier record schema boundary is documented in `../docs/token-verifier-record-schema-boundary.md` and `../decisions/ADR-0033-token-verifier-record-schema-boundary.md`. The token verifier migration source now exists at `migrations/postgres/000004_create_authentication_access_tokens.sql` for the ratified `authentication_access_tokens` semantics. Runtime code must remain unchanged until later bounded work explicitly authorizes repositories, adapters, runtime token issuance, validation, logout, refresh, cleanup, handlers, routes, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

Authentication schema migration queue planning is documented in `../docs/authentication-schema-migration-queue.md` and `../decisions/ADR-0034-authentication-schema-migration-queue.md`. The adapter implementation gate is documented in `../decisions/ADR-0035-authentication-postgresql-adapter-implementation-gate.md`. Credential and token verifier migration sources, static checks, the storage-neutral authentication repository interface boundary, and the authentication PostgreSQL adapter boundary are now present. Runtime authentication code must remain unchanged until a later bounded work item explicitly authorizes runtime credential lookup, token issuance, token validation, logout execution, refresh, cleanup jobs, handlers, routes, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

`M-019 Token And Credential Verifier Algorithm Redaction Boundary` is completed. `W-0091` defined the first planned high-entropy verifier posture in `../docs/token-credential-verifier-algorithm-redaction-boundary.md` and `../decisions/ADR-0040-token-credential-verifier-algorithm-redaction-boundary.md` without adding runtime authentication behavior. The first planned verifier algorithm family is `vibit_hmac_sha256_v1`; future first-posture code may use Go standard library packages `crypto/hmac`, `crypto/sha256`, `crypto/subtle`, `crypto/rand`, and `encoding/base64` after a later code gate. Lookup digests and verifier digests are not log-safe, verifier digest comparison must be constant-time, and raw credential/token material must have at least 256 bits of entropy. `runtime.token_credential_verifier_algorithm_redaction_boundary` is the repository check rule for the completed verifier boundary.

`M-020 Secret Configuration And Verifier Key Loading Boundary` is completed. `W-0092` defined future key loading posture in `../docs/secret-configuration-verifier-key-loading-boundary.md` and `../decisions/ADR-0041-secret-configuration-verifier-key-loading-boundary.md` without adding runtime authentication behavior. Future verifier key loading is application-owned under `internal/app`; the first local implementation may use process environment configuration or explicit runtime secret input after a later code gate; external KMS or secret-manager integration remains behind dependency and operations gates. Four separated logical verifier keys are required, `verifier_key_id` is not log-safe by default, production key configuration must fail closed when invalid, and committed production-like secret values remain forbidden. Do not implement secret loading, environment parsing, token generation, credential generation, verifier comparison, login, token validation, logout execution, refresh, cleanup jobs, Protobuf messages, WebSocket behavior, authentication dependencies, repository changes, migration schema changes, KMS integration, secret-manager integration, or production authentication behavior from this boundary. `runtime.secret_configuration_verifier_key_loading_boundary` is the repository check rule for this boundary.

`M-021 Token And Credential Material Generation Boundary` is completed. `W-0093` defined future raw device credential and opaque access-token material generation posture in `../docs/token-credential-material-generation-boundary.md` and `../decisions/ADR-0042-token-credential-material-generation-boundary.md` without adding runtime authentication behavior. Future material generation is application-owned under `internal/app`; the first device credential and access token are server-issued and application-generated; raw material must be 32 cryptographically random bytes with at least 256 bits of entropy; text presentation is URL-safe unpadded Base64 or equivalent; raw material is one-time client-visible and must not be stored. Repository handoff remains digest-only. Future first-posture generation helpers may use Go standard library `crypto/rand` and `encoding/base64` after a later code gate. Do not implement token generation, credential generation, secret loading, verifier digest computation, verifier comparison, login, token validation, logout execution, refresh, cleanup jobs, Protobuf messages, WebSocket behavior, authentication dependencies, repository changes, migration schema changes, or production authentication behavior from this boundary. `runtime.token_credential_material_generation_boundary` is the repository check rule for this boundary.

`M-022 Verifier Digest Computation And Comparison Boundary` is completed. `W-0094` defined future verifier digest computation and constant-time comparison posture in `../docs/verifier-digest-computation-comparison-boundary.md` and `../decisions/ADR-0043-verifier-digest-computation-comparison-boundary.md` without adding runtime authentication behavior. Future digest helpers are application-owned under `internal/app`; canonical HMAC input is versioned and length-prefixed; lookup digest equality is record selection only; validation may compute lookup candidates across active and accepted previous key sets; stored `verifier_key_id` selects the verifier key set for verifier digest computation; verifier digest comparison must be constant-time; invalid lookup, mismatch, unknown key id, unsupported algorithm, malformed proof, and expired or revoked proof collapse to the same public invalid-proof class. Future first-posture helpers may use Go standard library `crypto/hmac`, `crypto/sha256`, and `crypto/subtle` after a later code gate. Do not implement verifier digest computation, verifier comparison, token generation, credential generation, secret loading, login, token validation, logout execution, refresh, cleanup jobs, Protobuf messages, WebSocket behavior, authentication dependencies, repository changes, migration schema changes, or production authentication behavior from this boundary. `runtime.verifier_digest_computation_comparison_boundary` is the repository check rule for this boundary.

`M-023 Authentication Service Implementation Readiness Gate` is completed. `W-0095` defined the readiness gate in `../docs/authentication-service-implementation-readiness-gate.md` and `../decisions/ADR-0044-authentication-service-implementation-readiness-gate.md` without adding runtime authentication behavior. Future service implementation remains application-owned under `internal/app`, with `internal/app/authentication` as the package candidate. The first code slice must be separately authorized; the recommended next gate is local verifier key configuration loading. This gate defines prior boundaries, allowed and forbidden write areas, test classes, redaction expectations, Nakama/Pitaya capability mapping, and deferrals. Do not implement authentication service code, secret loading, token generation, credential generation, verifier digest computation, verifier comparison, login, token validation, logout execution, refresh, cleanup jobs, Protobuf messages, WebSocket behavior, authentication dependencies, repository changes, migration schema changes, or production authentication behavior from this gate. `runtime.authentication_service_implementation_readiness_gate` is the repository check rule for this gate.

`M-024 Local Verifier Key Configuration Loading Gate` is completed. `W-0096` defined the gate in `../docs/local-verifier-key-configuration-loading-gate.md` and `../decisions/ADR-0045-local-verifier-key-configuration-loading-gate.md` without adding runtime authentication behavior. `W-0097` implemented the explicit in-memory verifier key set validator under `internal/app/authentication`, with tests for accepted input, copying, immutable accessors, invalid key sets, and redacted errors. Environment parsing and Base64 text decoding are deferred to `W-0098`, the environment verifier key loader gate. Do not parse process environment variables, read local secret files, wire startup, integrate KMS or cloud secret managers, generate tokens or credentials, compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, add authentication dependencies, or add production authentication behavior from this gate. `runtime.local_verifier_key_configuration_loading_gate` is the repository check rule for this gate.

`M-026 Environment Verifier Key Loader Gate` is completed. `W-0098` defined the gate in `../docs/environment-verifier-key-loader-gate.md` and `../decisions/ADR-0046-environment-verifier-key-loader-gate.md` without adding Go code or runtime authentication behavior. Future process environment loader work belongs under `internal/app/authentication`, should use `verifier_key_env.go` and `verifier_key_env_test.go`, must decode required environment key text and call `NewVerifierKeySet`, and must keep environment variable values and full concrete key set ids redacted. Do not wire startup, read local secret files, parse `.env` files, accept CLI secret input, integrate KMS or cloud secret managers, generate tokens or credentials, compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, add authentication dependencies, or add production authentication behavior from this gate. `runtime.environment_verifier_key_loader_gate` is the repository check rule for this gate.

`M-027 Environment Verifier Key Loader Implementation` is completed. `W-0099` implemented `internal/app/authentication/verifier_key_env.go` and `internal/app/authentication/verifier_key_env_test.go`. The loader accepts an explicit lookup function, provides a tiny `os.LookupEnv` process adapter, decodes Base64URL unpadded and standard padded Base64 key text, calls `NewVerifierKeySet`, and returns redacted typed errors. Do not wire this loader into process startup, read local secret files, parse `.env` files, accept CLI secret input, integrate KMS or cloud secret managers, generate tokens or credentials, compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, add authentication dependencies, or add production authentication behavior without a later bounded work item.

`M-028 Token And Credential Material Generation Implementation Gate` is completed. `W-0100` defined the gate in `../docs/token-credential-material-generation-implementation-gate.md` and `../decisions/ADR-0047-token-credential-material-generation-implementation-gate.md` without adding Go code or runtime authentication behavior. Future helper work belongs under `internal/app/authentication`, should use `material_generation.go` and `material_generation_test.go`, must use 32 random bytes, URL-safe unpadded Base64 presentation, explicit `io.Reader` entropy-source handoff, copied raw bytes, redacted errors, and focused tests. Do not compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior from this gate. `runtime.token_credential_material_generation_implementation_gate` is the repository check rule for this gate.

`M-029 Token And Credential Material Generation Helper Implementation` is completed. `W-0101` implemented `internal/app/authentication/material_generation.go` and `internal/app/authentication/material_generation_test.go`. The helpers expose `GenerateDeviceCredentialMaterial` and `GenerateAccessTokenMaterial`, accept an explicit `io.Reader`, read 32 bytes with `io.ReadFull`, encode URL-safe unpadded Base64 text, preserve `MaterialKind`, copy raw bytes on return, reject nil readers, read errors, short reads, all-zero material, and repeated-single-byte material, and return redacted typed errors. Do not expand these helpers to compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior.

`M-031 Verifier Digest Computation Helper Implementation` is completed. `W-0103` implemented `internal/app/authentication/verifier_digest.go` and `internal/app/authentication/verifier_digest_test.go`. The helpers expose `ComputeCredentialLookupDigest`, `ComputeCredentialVerifierDigest`, `ComputeTokenLookupDigest`, and `ComputeTokenVerifierDigest`, build the canonical `vibit.auth.verifier.input.v1` HMAC input, use the matching logical key from `VerifierKeySet`, return copied 32-byte digest bytes through `ComputedDigest`, and use redacted typed errors. Do not expand these helpers to compare verifier digests, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior. The next ready work item is the verifier digest comparison helper gate.

`M-032 Verifier Digest Comparison Helper Gate` is completed. `W-0104` defined the gate in `../docs/verifier-digest-comparison-helper-gate.md` and `../decisions/ADR-0049-verifier-digest-comparison-helper-gate.md` without adding Go code or runtime authentication behavior. Future comparison helper work belongs under `internal/app/authentication`, should use `verifier_comparison.go` and `verifier_comparison_test.go`, must keep `verifier_digest.go` computation-only, compare only computed verifier digest bytes to stored verifier digest bytes, prefer `crypto/hmac.Equal`, allow `crypto/subtle.ConstantTimeCompare` only as an equivalent constant-time alternative, reject lookup digest classes and malformed input, and keep raw material, lookup digests, key ids, database-only equality, protocol metadata, and public failure details out of comparison. Do not implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior from this gate. `runtime.verifier_digest_comparison_helper_gate` is the repository check rule for this gate.

`M-033 Verifier Digest Comparison Helper Implementation` is completed. `W-0105` implemented `internal/app/authentication/verifier_comparison.go` and `internal/app/authentication/verifier_comparison_test.go`. The helpers expose `CompareCredentialVerifierDigest` and `CompareTokenVerifierDigest`, return `VerifierComparisonResult`, use `crypto/hmac.Equal`, compare only `ComputedDigest` verifier bytes to stored verifier digest bytes, reject lookup digest classes, wrong verifier classes, missing input, malformed computed input, and malformed stored input, and keep errors redacted. Do not expand these helpers to implement authentication service behavior, login execution, token validation, logout, refresh, cleanup, protocol carriers, repository calls, migration changes, startup wiring, authentication dependencies, or production authentication behavior without a later bounded work item. The next ready work item is the authentication service behavior implementation gate.

The implemented authentication PostgreSQL adapter boundary is `internal/platform/persistence/postgres/authentication_repository.go`, with focused tests at `internal/platform/persistence/postgres/authentication_repository_test.go`. Its constructor is `NewAuthenticationRepositoryForUnitOfWork(executor)`, and `UnitOfWork.NewAuthenticationRepository` creates an `authentication.Repository` from the caller-owned executor. The adapter remains persistence-adapter-only.

Runtime authentication implementation boundary planning is documented in `../docs/runtime-authentication-implementation-boundary.md` and `../decisions/ADR-0036-runtime-authentication-implementation-boundary.md`. Future runtime authentication is application-owned under `internal/app`; it must use `authentication.Repository` through the application unit-of-work boundary, keep the PostgreSQL adapter persistence-only, and convert validated proof into `RequestIdentity` before domain dispatch. Token generation, verifier comparison, login execution, access-token validation, logout execution, cleanup jobs, Protobuf authentication messages, WebSocket proof carriers, generated authentication shapes, and authentication dependencies remain separately gated. `runtime.authentication_implementation_boundary` is the repository check rule for this boundary.

Authentication generated contract shape timing is documented in `../docs/authentication-generated-contract-shape-timing.md` and `../decisions/ADR-0038-authentication-generated-contract-shape-timing.md`. The source is `contracts/runtime/authentication/`, and the output root is `runtime/internal/generated/contracts/runtime/authentication/`. Generated authentication shapes are metadata-only and immutable; service interfaces and runtime behavior remain separately gated.

Application authentication service interface boundary is documented in `../docs/application-authentication-service-interface-boundary.md` and `../decisions/ADR-0039-application-authentication-service-interface-boundary.md`. Future authentication service interfaces are application-owned under `internal/app`; generated authentication shapes inform service-level request/result vocabulary; service behavior may use `authentication.Repository` only through the application unit-of-work boundary; validated proof must become `RequestIdentity` before domain dispatch. This boundary does not authorize Go service code or runtime authentication behavior. `runtime.application_authentication_service_interface_boundary` is the repository check rule for this boundary.

`../decisions/ADR-0037-close-runtime-auth-boundary-and-open-generated-shape-gate.md` closes M-016 and opens M-017. `ADR-0038` completes the timing decision. `W-0089` completes generator/check support plus metadata-only generated authentication shape output. `ADR-0039` and `W-0090` complete the service-interface boundary step. `ADR-0040` and `W-0091` complete the verifier algorithm/redaction step. `ADR-0041` and `W-0092` complete the secret configuration/verifier key loading preparation step. `ADR-0042` and `W-0093` complete the material generation preparation step. `ADR-0043` and `W-0094` complete the verifier digest computation and comparison preparation step. `ADR-0044` and `W-0095` complete the implementation readiness step. `ADR-0045` and `W-0096` complete the local verifier key configuration loading gate. `W-0097` completes the explicit in-memory verifier key set validator implementation slice. `ADR-0046` and `W-0098` complete the environment verifier key loader gate. `W-0099` completes the environment verifier key loader implementation slice. `ADR-0047` and `W-0100` complete the token and credential material generation implementation gate. `W-0101` completes the token and credential material generation helper implementation slice. `ADR-0048` and `W-0102` complete the verifier digest helper implementation gate. `W-0103` completes the verifier digest computation helper implementation slice. `ADR-0049` and `W-0104` complete the verifier digest comparison helper gate. `W-0105` completes the verifier digest comparison helper implementation slice.

The first explicit PostgreSQL migration runner is `internal/platform/migrations/postgres.go`. It owns `github.com/pressly/goose/v3`, accepts a caller-supplied `*sql.DB` and migration source filesystem or directory, lists SQL migration sources, reports structured status, and applies pending migrations only when explicitly invoked. Do not wire it into normal `cmd/vibit-server` startup without a change spec.

Live PostgreSQL verification is governed by `../docs/postgresql-verification-environment.md`. It is opt-in through `VIBIT_POSTGRES_TEST_DSN`; normal unit tests, `node ../tools/vibit check runtime`, and default repository checks must not require a running PostgreSQL server. When a live PostgreSQL check is skipped because no disposable DSN is available, record that explicitly.

## 6. Generated Files

Generated files are immutable to non-system agents.

If generated output is wrong, change the source contract, schema, template, or generator. Do not hand-edit generated files unless a change spec or Agent Decision Record explicitly grants `generated_file_override`.

Go Protobuf generated output under `internal/generated/proto/` must be produced from `../proto/` sources through the accepted Buf and `protoc-gen-go` path. Files under that root must be generated `*.pb.go` files with the `protoc-gen-go` marker and source trace, or temporary `.gitkeep` placeholders while generation has not run.

Do not place handwritten runtime code under `internal/generated/proto/` or `internal/generated/contracts/`.

## 7. Current State

This runtime workspace now has the first generated Protobuf output, the first narrow runtime handoff slice, the first WebSocket transport adapter, a small application dispatch skeleton for command and query routes, the first transaction boundary skeleton, the first inventory repository/policy/handler runtime boundary with a command-safe mutation lock, the first PostgreSQL configuration parser, the first pgx-backed transaction runner adapter, the first PostgreSQL inventory repository adapter, the first inventory Protobuf/domain payload bridge, the first application-error-to-Protobuf-error-envelope mapper, the first frame-to-Protobuf-to-application composition adapter, a package-local request-loop test fixture for Protobuf command/query tests, minimal process wiring that mounts `/v1/ws`, an explicit PostgreSQL inventory runtime composition path, and an opt-in live PostgreSQL durable inventory request-loop verification test.

The workspace has a documented PostgreSQL persistence boundary, transaction skeleton, PostgreSQL configuration parser, pgx-backed transaction runner, first inventory migration source, first explicit migration apply/status runner, first PostgreSQL repository adapter, explicit runtime store selection, a ratified player account PostgreSQL lifecycle schema boundary with its first migration source, storage-neutral repository interface, focused PostgreSQL adapter implementation, PostgreSQL unit-of-work factory helper, ratified authentication credential and token verifier migration sources with static checks, a storage-neutral authentication repository interface boundary, an implemented authentication PostgreSQL adapter, a ratified authentication/token/session validation design boundary without runtime authentication implementation, a runtime authentication implementation boundary planning standard, and metadata-only generated authentication contract shapes, plus a live verification test that skips unless `VIBIT_POSTGRES_TEST_DSN` is set. `VIBIT_RUNTIME_STORE=memory` remains the default. `VIBIT_RUNTIME_STORE=postgres` enables PostgreSQL-backed inventory composition when `VIBIT_POSTGRES_DSN` is provided. The workspace still does not implement generated route registration, generated protocol bridge creation, production authentication/session validation, runtime player account handlers, WebSocket player routes, automatic startup migrations, or catalog-driven error retryability yet.

The first manual process run path is:

```bash
cd runtime
go run ./cmd/vibit-server
```

The first explicit persistent process run path is:

```bash
cd runtime
VIBIT_RUNTIME_STORE=postgres VIBIT_POSTGRES_DSN='postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable' go run ./cmd/vibit-server
```

Migrations are not applied automatically during normal server startup.

The first opt-in live durable inventory verification command is:

```bash
cd runtime
VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

If `VIBIT_POSTGRES_TEST_DSN` is unset, this test skips and records that live PostgreSQL verification was unavailable. The first live execution has passed on local Termux PostgreSQL 18.2.

## 8. Verification

Run repository verification from the repository root:

```bash
node tools/vibit check runtime
node tools/vibit check generated
node tools/vibit check migrations
node tools/vibit check postgres-env
node tools/vibit check all
```

When Go source files exist and the local Go toolchain is available, runtime verification should include:

```bash
go test ./...
go vet ./...
```

Do not claim Go test verification when the Go toolchain is unavailable or tests were not run.
