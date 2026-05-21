# Agent Operating Guide

Status: Draft v0.2
Last updated: 2026-05-20
Scope: Repository-level operating instructions for coding agents  
Canonical source: `CONSTITUTION.md`

This guide turns the constitution into working rules for agents. It does not replace the constitution. When this guide and `CONSTITUTION.md` conflict, follow `CONSTITUTION.md` and update this guide.

The paired Simplified Chinese translation is `AGENTS.zh-CN.md`. The English file is authoritative.

## 1. Project Identity

Working name:

```text
vibit
```

Category:

```text
Agent-Native Server Framework
```

Positioning:

```text
vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.
```

In this repository, "AI-native" means agent-native maintainability first. It does not primarily mean adding AI gameplay features or AI product features.

## 2. Required Reading

Before making a non-trivial change, read:

- `CONSTITUTION.md`
- This file
- The relevant architecture manifests under `.arch/`, when they exist
- The relevant module manifest at `modules/<module>/module.yaml`, when it exists
- The relevant module guide at `modules/<module>/AGENTS.md`, when it exists
- The relevant change spec under `changes/`, when the change has one

If an expected artifact does not exist yet, do not invent hidden assumptions. Either create the missing artifact as part of the change or record that it is not yet available.

## 3. Current Repository State

This repository is currently pre-alpha and building toward `v0.1 alpha`.

The short-term target is a first developer-usable alpha that external developers can download, run locally, inspect, and use as a contribution entry point. The durable goal document is `docs/v0.1-alpha-goal.md`, with `docs/v0.1-alpha-goal.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0086` records the decision.

The long-term product target is an AI-era Nakama/Pitaya-class open-source game/backend server framework. This means same-class common capability coverage adapted to vibit's agent-native maintainability model. It does not mean direct Nakama/Pitaya API compatibility.

Existing foundation:

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/README.md`
- `.arch/modules.yaml`
- `.arch/conventions.yaml`
- `.arch/protocol.yaml`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/dependencies.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `docs/module-manifest.md`
- `docs/module-manifest.zh-CN.md`
- `docs/change-spec.md`
- `docs/change-spec.zh-CN.md`
- `changes/_template/`
- `docs/conversation-log.md`
- `docs/conversation-log.zh-CN.md`
- `conversations/`
- `docs/agent-decision-record.md`
- `docs/agent-decision-record.zh-CN.md`
- `decisions/`
- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `docs/dependency-adoption.md`
- `docs/dependency-adoption.zh-CN.md`
- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `docs/reference-game-server-alignment.md`
- `docs/reference-game-server-alignment.zh-CN.md`
- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/authentication-proof-token-session-contract-dimensions.zh-CN.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/credential-storage-external-identity-linking-boundaries.zh-CN.md`
- `docs/session-persistence-websocket-handshake-decision-gates.md`
- `docs/session-persistence-websocket-handshake-decision-gates.zh-CN.md`
- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `docs/workflow.md`
- `docs/workflow.zh-CN.md`
- `schema/`
- `rules/`

Framework implementation code now exists under `runtime/`, generated output exists under `runtime/internal/generated/`, and verification commands exist through `tools/vibit` plus Go tests. When a capability or verification path does not exist yet, document that it is not available instead of pretending that it ran.

Runtime readiness decisions currently point to Go as the first server runtime implementation language, WebSocket as the first gameplay/client protocol, Protobuf as the first wire message format, a modular monolith single-process server model, contract-first commands/queries/events/errors/permissions, and `inventory` as the preferred first proof slice. Read `.arch/runtime.yaml`, `ADR-0004` through `ADR-0010`, and note that `ADR-0003` is superseded before creating runtime implementation code.

The first game protocol framework is recorded in `.arch/protocol.yaml`, `docs/game-protocol.md`, `ADR-0015`, and `ADR-0016`. It defines a WebSocket-framed Protobuf envelope with explicit `kind`, `module`, and `name` routing fields, session metadata, game target scopes, server-authoritative message rules, error mapping, and compatibility expectations. Read it before adding `.proto` files, WebSocket protocol handlers, generated protocol output, or client/server protocol rules.

The first protocol source files are `proto/vibit/protocol/v1/envelope.proto` and `proto/vibit/inventory/v1/inventory.proto`. Buf configuration lives at `buf.yaml` and `buf.gen.yaml`. `ADR-0016` records the envelope and generation configuration decision. Generated Go Protobuf output remains planned under `runtime/internal/generated/proto/`; do not create or edit generated Go Protobuf files by hand.

The generated output standard is `docs/generated-output.md`, with `docs/generated-output.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0017` records the generated output decision. Read these before adding generated files, generated output checks, or generator behavior. Go Protobuf output under `runtime/internal/generated/proto/` must be `*.pb.go`, must contain the `protoc-gen-go` generated-code marker, and must trace to an existing `.proto` source.

The runtime protocol adapter boundary standard is `docs/runtime-protocol-adapter.md`, with `docs/runtime-protocol-adapter.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0018` records the boundary decision. Read these before adding WebSocket transport code, Protobuf runtime adapter code, application dispatch code, or domain runtime handlers.

The active game server reference alignment standard is `docs/reference-game-server-alignment.md`, with `docs/reference-game-server-alignment.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0019` records Nakama and Pitaya as active reference baselines. Read `.arch/reference.yaml` and the standard before adding new game server capability families, runtime subsystems, social/realtime features, matchmaking, match runtime, cluster/RPC work, or operational surfaces. Nakama and Pitaya guide capability planning; they do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

The authentication, token, and session validation design standard is `docs/authentication-token-session-validation.md`, with `docs/authentication-token-session-validation.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0023` records the design boundary. Read it before adding authentication, token behavior, credential storage, external identity linking, session persistence, request identity trust changes, Protobuf envelope authentication changes, WebSocket handshake authentication, runtime player account handlers, or WebSocket player routes. Metadata-only `player_id` and `session_id` are not authenticated proof. `runtime.authentication_token_session_boundary` is the repository check rule for this boundary.

The authentication proof and token/session contract dimensions standard is `docs/authentication-proof-token-session-contract-dimensions.md`, with `docs/authentication-proof-token-session-contract-dimensions.zh-CN.md` as the paired Simplified Chinese translation. Read it before defining or changing authentication proof, token/session validation, actor kinds, validation statuses, proof statuses, failure classes, retryability, request identity handoff, session errors, session permissions, or validation events. It ratifies semantic vocabulary only; it does not choose login methods, token formats, credential storage, session persistence, Protobuf envelope behavior, WebSocket handshake behavior, runtime player handlers, or WebSocket routes.

The credential storage and external identity linking boundary standard is `docs/credential-storage-external-identity-linking-boundaries.md`, with `docs/credential-storage-external-identity-linking-boundaries.zh-CN.md` as the paired Simplified Chinese translation. Read it before adding credential storage, external identity linking, login methods, provider subjects, password hashing, OAuth, OIDC, provider SDKs, account linking, account recovery, account merge behavior, or related schema. The standard defines deferred capability families only. It preserves `player_accounts` and `player_account_events` as account lifecycle tables and does not authorize credential tables, external identity tables, provider dependencies, runtime lookup code, player lifecycle table changes, or direct Nakama/Pitaya API compatibility.

The session persistence and WebSocket handshake decision-gates standard is `docs/session-persistence-websocket-handshake-decision-gates.md`, with `docs/session-persistence-websocket-handshake-decision-gates.zh-CN.md` as the paired Simplified Chinese translation. Read it before adding session persistence, WebSocket handshake authentication, reconnect behavior, connection epoch behavior, token/session carriers, session-related Protobuf envelope changes, handshake/system messages, or route-level authentication. It defines request-level, first-message, handshake-level, every-request, and hybrid validation as future choices only; it does not select a production model, session store, session tables, envelope behavior, or handshake behavior.

The login method and token format ratification standard is `docs/login-method-token-format-ratification.md`, with `docs/login-method-token-format-ratification.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0024` records the ratification boundary. Read it before selecting first login methods, token model, token format, proof carrier posture, token lifecycle semantics, credential/token/session schema gates, or implementation queue. It may guide comparison and ratification only; it does not authorize runtime authentication, token parsing, credential storage, external identity linking, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes.

The selected login/token boundary check standard is `docs/selected-login-token-boundary-checks.md`, with `docs/selected-login-token-boundary-checks.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0030` records the repository check decision. Read it before adding runtime authentication, token validation, token issuance, logout, refresh behavior, credential or token repositories, authentication Protobuf sources, generated authentication contract shapes, WebSocket proof carriers, authentication migrations, authentication dependencies, or changes to the selected `device_credential_login` plus opaque access-token posture. `runtime.selected_login_token_boundary` is the repository check rule for this selected posture.

The credential record schema boundary standard is `docs/credential-record-schema-boundary.md`, with `docs/credential-record-schema-boundary.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0032` records the boundary decision. The credential migration source now exists at `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`; it does not authorize repositories, PostgreSQL adapters, runtime lookup, handlers, routes, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

The token verifier record schema boundary standard is `docs/token-verifier-record-schema-boundary.md`, with `docs/token-verifier-record-schema-boundary.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0033` records the boundary decision. The token verifier migration source now exists at `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`; it does not authorize repositories, PostgreSQL adapters, runtime token issuance, validation, logout, refresh, cleanup, handlers, routes, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

The authentication schema migration queue standard is `docs/authentication-schema-migration-queue.md`, with `docs/authentication-schema-migration-queue.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0034` records the queue decision, and `ADR-0035` records the adapter implementation gate. Credential and token verifier migration sources, static checks, the storage-neutral authentication repository interface boundary, and the authentication PostgreSQL adapter boundary are now present. This queue does not authorize runtime credential lookup, token issuance, token validation, logout execution, refresh, cleanup jobs, handlers, routes, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

`M-019 Token And Credential Verifier Algorithm Redaction Boundary` is completed. `W-0091` defined the first planned high-entropy verifier posture in `docs/token-credential-verifier-algorithm-redaction-boundary.md` and `ADR-0040` without adding runtime authentication behavior. The first planned verifier algorithm family is `vibit_hmac_sha256_v1`; future first-posture code may use Go standard library packages `crypto/hmac`, `crypto/sha256`, `crypto/subtle`, `crypto/rand`, and `encoding/base64` after a later code gate. Lookup digests and verifier digests are not log-safe, verifier digest comparison must be constant-time, and raw credential/token material must have at least 256 bits of entropy. `runtime.token_credential_verifier_algorithm_redaction_boundary` is the repository check rule for the completed verifier boundary.

`M-020 Secret Configuration And Verifier Key Loading Boundary` is completed. `W-0092` defined the future key loading posture in `docs/secret-configuration-verifier-key-loading-boundary.md` and `ADR-0041` without adding runtime authentication behavior. Future verifier key loading is application-owned under `runtime/internal/app`; the first local implementation may use process environment configuration or explicit runtime secret input after a later code gate; external KMS or secret-manager integration remains behind dependency and operations gates. Four separated logical verifier keys are required, `verifier_key_id` is not log-safe by default, production key configuration must fail closed when invalid, and committed production-like secret values remain forbidden. This boundary must not implement secret loading, environment parsing, token generation, credential generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, migration schema changes, KMS integration, secret-manager integration, or production authentication behavior. `runtime.secret_configuration_verifier_key_loading_boundary` is the repository check rule for this boundary.

`M-021 Token And Credential Material Generation Boundary` is completed. `W-0093` defined future raw device credential and opaque access-token material generation posture in `docs/token-credential-material-generation-boundary.md` and `ADR-0042` without adding runtime authentication behavior. Future material generation is application-owned under `runtime/internal/app`; the first device credential and access token are server-issued and application-generated; raw material must be 32 cryptographically random bytes with at least 256 bits of entropy; text presentation is URL-safe unpadded Base64 or equivalent; raw material is one-time client-visible and must not be stored. Repository handoff remains digest-only. Future first-posture generation helpers may use Go standard library `crypto/rand` and `encoding/base64` after a later code gate. This boundary must not implement token generation, credential generation, secret loading, verifier digest computation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, migration schema changes, or production authentication behavior. `runtime.token_credential_material_generation_boundary` is the repository check rule for this boundary.

`M-022 Verifier Digest Computation And Comparison Boundary` is completed. `W-0094` defined future verifier digest computation and constant-time comparison posture in `docs/verifier-digest-computation-comparison-boundary.md` and `ADR-0043` without adding runtime authentication behavior. Future digest helpers are application-owned under `runtime/internal/app`; canonical HMAC input is versioned and length-prefixed; lookup digest equality is record selection only; validation may compute lookup candidates across active and accepted previous key sets; stored `verifier_key_id` selects the verifier key set for verifier digest computation; verifier digest comparison must be constant-time; invalid lookup, mismatch, unknown key id, unsupported algorithm, malformed proof, and expired or revoked proof collapse to the same public invalid-proof class. Future first-posture helpers may use Go standard library `crypto/hmac`, `crypto/sha256`, and `crypto/subtle` after a later code gate. This boundary must not implement verifier digest computation, verifier comparison, token generation, credential generation, secret loading, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, migration schema changes, or production authentication behavior. `runtime.verifier_digest_computation_comparison_boundary` is the repository check rule for this boundary.

`M-023 Authentication Service Implementation Readiness Gate` is completed. `W-0095` defined the readiness gate in `docs/authentication-service-implementation-readiness-gate.md` and `ADR-0044` without adding runtime authentication behavior. Future service implementation remains application-owned under `runtime/internal/app`, with `runtime/internal/app/authentication` as the package candidate. The first code slice must be separately authorized; the recommended next gate is local verifier key configuration loading. This gate defines prior boundaries, allowed and forbidden write areas, test classes, redaction expectations, Nakama/Pitaya capability mapping, and deferrals. It must not implement authentication service code, secret loading, token generation, credential generation, verifier digest computation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, migration schema changes, or production authentication behavior. `runtime.authentication_service_implementation_readiness_gate` is the repository check rule for this gate.

`M-024 Local Verifier Key Configuration Loading Gate` is completed. `W-0096` defined the gate in `docs/local-verifier-key-configuration-loading-gate.md` and `ADR-0045` without adding Go code or runtime authentication behavior. `W-0097` implemented the explicit in-memory verifier key set validator under `runtime/internal/app/authentication`, with tests for accepted input, copying, immutable accessors, invalid key sets, and redacted errors. Environment parsing and Base64 text decoding are deferred to `W-0098`, the environment verifier key loader gate. It must not parse process environment variables, read local secret files, wire startup, integrate KMS or cloud secret managers, generate tokens or credentials, compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, add authentication dependencies, or add production authentication behavior. `runtime.local_verifier_key_configuration_loading_gate` is the repository check rule for this gate.

`M-026 Environment Verifier Key Loader Gate` is completed. `W-0098` defined the gate in `docs/environment-verifier-key-loader-gate.md` and `ADR-0046` without adding Go code or runtime authentication behavior. Future process environment loader work belongs under `runtime/internal/app/authentication`, should use `verifier_key_env.go` and `verifier_key_env_test.go`, must decode required environment key text and call `NewVerifierKeySet`, and must keep environment variable values and full concrete key set ids redacted. It must not wire startup, read local secret files, parse `.env` files, accept CLI secret input, integrate KMS or cloud secret managers, generate tokens or credentials, compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, add authentication dependencies, or add production authentication behavior. `runtime.environment_verifier_key_loader_gate` is the repository check rule for this gate.

`M-027 Environment Verifier Key Loader Implementation` is completed. `W-0099` implemented `runtime/internal/app/authentication/verifier_key_env.go` and `runtime/internal/app/authentication/verifier_key_env_test.go`. The loader accepts an explicit lookup function, provides a tiny `os.LookupEnv` process adapter, decodes Base64URL unpadded and standard padded Base64 key text, calls `NewVerifierKeySet`, and returns redacted typed errors. Do not wire this loader into process startup, read local secret files, parse `.env` files, accept CLI secret input, integrate KMS or cloud secret managers, generate tokens or credentials, compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, add authentication dependencies, or add production authentication behavior without a later bounded work item.

`M-028 Token And Credential Material Generation Implementation Gate` is completed. `W-0100` defined the gate in `docs/token-credential-material-generation-implementation-gate.md` and `ADR-0047` without adding Go code or runtime authentication behavior. Future material generation helper work belongs under `runtime/internal/app/authentication`, should use `material_generation.go` and `material_generation_test.go`, must use 32 random bytes, URL-safe unpadded Base64 presentation, explicit `io.Reader` entropy-source handoff, copied raw bytes, redacted errors, and focused tests. It must not compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior. `runtime.token_credential_material_generation_implementation_gate` is the repository check rule for this gate.

`M-029 Token And Credential Material Generation Helper Implementation` is completed. `W-0101` implemented `runtime/internal/app/authentication/material_generation.go` and `runtime/internal/app/authentication/material_generation_test.go`. The helpers expose `GenerateDeviceCredentialMaterial` and `GenerateAccessTokenMaterial`, accept an explicit `io.Reader`, read 32 bytes with `io.ReadFull`, encode URL-safe unpadded Base64 text, preserve `MaterialKind`, copy raw bytes on return, reject nil readers, read errors, short reads, all-zero material, and repeated-single-byte material, and return redacted typed errors. Do not expand these helpers to compute digests, compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior without a later bounded work item.

`M-030 Verifier Digest Helper Implementation Gate` is completed. `W-0102` defined the verifier digest helper implementation gate in `docs/verifier-digest-helper-implementation-gate.md` and `ADR-0048` without adding Go code or runtime authentication behavior. Future digest helper work belongs under `runtime/internal/app/authentication`, should use `verifier_digest.go` and `verifier_digest_test.go`, must build deterministic canonical input with versioned ASCII header, null separator, length-prefixed purpose label and raw material, compute HMAC-SHA-256 with the matching logical key from an already-validated `VerifierKeySet`, return copied 32-byte digest bytes, expose redacted errors, and add focused tests. It must not compare verifiers, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior. `runtime.verifier_digest_helper_implementation_gate` is the repository check rule for this gate.

`M-031 Verifier Digest Computation Helper Implementation` is completed. `W-0103` implemented `runtime/internal/app/authentication/verifier_digest.go` and `runtime/internal/app/authentication/verifier_digest_test.go`. The helpers expose `ComputeCredentialLookupDigest`, `ComputeCredentialVerifierDigest`, `ComputeTokenLookupDigest`, and `ComputeTokenVerifierDigest`, build the canonical `vibit.auth.verifier.input.v1` HMAC input, use the matching logical key from `VerifierKeySet`, return copied 32-byte digest bytes through `ComputedDigest`, and use redacted typed errors. Do not expand these helpers to compare verifier digests, implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior without a later bounded work item. The next ready work item is the verifier digest comparison helper gate.

`M-032 Verifier Digest Comparison Helper Gate` is completed. `W-0104` defined the verifier digest comparison helper gate in `docs/verifier-digest-comparison-helper-gate.md` and `ADR-0049` without adding Go code or runtime authentication behavior. Future comparison helper work belongs under `runtime/internal/app/authentication`, should use `verifier_comparison.go` and `verifier_comparison_test.go`, must keep `verifier_digest.go` computation-only, compare only computed verifier digest bytes to stored verifier digest bytes, prefer `crypto/hmac.Equal`, allow `crypto/subtle.ConstantTimeCompare` only as an equivalent constant-time alternative, reject lookup digest classes and malformed input, and keep raw material, lookup digests, key ids, database-only equality, protocol metadata, and public failure details out of comparison. It must not implement authentication service behavior, expose protocol carriers, change repositories, change migrations, wire startup, add authentication dependencies, or add production authentication behavior. `runtime.verifier_digest_comparison_helper_gate` is the repository check rule for this gate. The next ready work item is the verifier digest comparison helper implementation slice.

`M-033 Verifier Digest Comparison Helper Implementation` is completed. `W-0105` implemented `runtime/internal/app/authentication/verifier_comparison.go` and `runtime/internal/app/authentication/verifier_comparison_test.go`. The helpers expose `CompareCredentialVerifierDigest` and `CompareTokenVerifierDigest`, return `VerifierComparisonResult`, use `crypto/hmac.Equal`, compare only `ComputedDigest` verifier bytes to stored verifier digest bytes, reject lookup digest classes, wrong verifier classes, missing input, malformed computed input, and malformed stored input, and keep errors redacted. Do not expand these helpers to implement authentication service behavior, login execution, token validation, logout, refresh, cleanup, protocol carriers, repository calls, migration changes, startup wiring, authentication dependencies, or production authentication behavior without a later bounded work item. The next ready work item is the authentication service behavior implementation gate.

`M-034 Authentication Service Behavior Implementation Gate` is completed. `W-0106` defined the gate in `docs/authentication-service-behavior-implementation-gate.md` and `ADR-0050` without adding service code or runtime authentication behavior. It requires future service behavior to remain application-owned under `runtime/internal/app/authentication`, to use `authentication.Repository` only through the application unit-of-work boundary, to compose verifier key, material generation, digest computation, and comparison helpers in a fixed order, and to collapse public proof failures. `runtime.authentication_service_behavior_implementation_gate` is the repository check rule for this gate and skeleton boundary.

`M-035 Authentication Service Behavior Skeleton` is completed. `W-0107` added `runtime/internal/app/authentication/service.go` and `runtime/internal/app/authentication/service_test.go` as a skeleton-only application service shape. The skeleton defines typed dependencies, request/result vocabulary, redacted internal failure classes, public error codes, and fail-closed `AUTHENTICATION_NOT_IMPLEMENTED` or `AUTHENTICATION_REFRESH_NOT_SUPPORTED` behavior. It does not execute login, validate access tokens, issue tokens, call repositories, revoke tokens, refresh tokens, expose protocol carriers, wire startup, add dependencies, change migrations, or add production authentication behavior. The next ready work item is the device credential login service behavior gate.

`M-036 Device Credential Login Service Behavior Gate` is completed. `W-0108` defined the gate in `docs/device-credential-login-service-behavior-gate.md` and `ADR-0051` without adding login execution code or runtime authentication behavior. Future device credential login must stay application-owned inside `runtime/internal/app/authentication/service.go`, treat `CredentialProof` as server-issued Base64URL unpadded 32-byte high-entropy material, reject missing or malformed proof before unit-of-work, use existing authentication and player repositories through unit-of-work capabilities, compare credential verifier digest before token generation, require active player account state, store token digests only, and return raw access-token text only after unit-of-work success. It must preserve deferrals for access-token validation, logout, refresh, cleanup, protocol carriers, startup wiring, repository changes, migrations, generated files, dependencies, and broader production authentication behavior. `runtime.device_credential_login_service_behavior_gate` is the repository check rule for this gate. The next ready work item is the device credential login service behavior implementation slice.

`M-037 Device Credential Login Service Behavior Implementation` is completed. `W-0109` implemented `AuthenticateWithDeviceCredential` inside `runtime/internal/app/authentication/service.go` and focused tests in `runtime/internal/app/authentication/service_test.go`. The service validates Base64URL unpadded 32-byte proof before unit-of-work, uses verifier helpers, obtains authentication and player repositories through unit-of-work capabilities, requires active credential and active player account state, stores opaque access-token digests only, and returns raw access-token text only after successful token storage and unit-of-work commit. It does not implement access-token validation, logout, refresh, cleanup, protocol carriers, startup wiring, repository interface changes, migrations, generated files, dependencies, or broader production authentication behavior. `runtime.device_credential_login_service_behavior_implementation` is the repository check rule for this slice. The next ready work item is the access-token validation service behavior gate.

`M-038 Access Token Validation Service Behavior Gate` is completed. `W-0110` defined the future access-token validation behavior gate in `docs/access-token-validation-service-behavior-gate.md` and `ADR-0052` without adding validation execution code. Future validation behavior must stay application-owned inside `runtime/internal/app/authentication`, treat `AccessToken` as opaque Base64URL unpadded 32-byte proof, reject missing or malformed proof before unit-of-work, use authentication and player repositories through unit-of-work capabilities, compare token verifier digest before request identity, require active player account state, keep `SessionValidated` false until session persistence is ratified, and collapse public invalid-token failures. It must preserve deferrals for WebSocket proof carriers, handshake authentication, route protection, session persistence, logout, refresh, cleanup, protocol carriers, startup wiring, repository changes, migrations, generated files, dependencies, and broader production authentication behavior. `runtime.access_token_validation_service_behavior_gate` is the repository check rule for this gate. The next ready work item is the access-token validation service behavior implementation slice.

`M-039 Access Token Validation Service Behavior Implementation` is completed. `W-0111` implemented `ValidateAccessToken` inside `runtime/internal/app/authentication/service.go` with focused tests in `runtime/internal/app/authentication/service_test.go`. The service rejects missing or malformed opaque Base64URL unpadded 32-byte proof before unit-of-work, computes token lookup and verifier digests with existing helpers, obtains authentication and player repositories through unit-of-work capabilities, checks token kind, status, verifier posture, audience, issue time, and expiration, compares token verifier digest before request identity, requires active player account state, returns validated player identity only after unit-of-work success, and keeps `SessionValidated` false. It does not add protocol carriers, WebSocket handshake authentication, route protection, session persistence, logout, refresh, cleanup, startup wiring, repository interface changes, PostgreSQL adapter changes, migrations, generated files, dependencies, or broader production authentication behavior. `runtime.access_token_validation_service_behavior_implementation` is the repository check rule for this slice.

`M-040 Next Direction Confirmation Gate` is completed. `W-0112` selected `expose_access_token_protocol_carrier_and_route_protection_gate` after the maintainer authorized continued progress. `M-041 Access Token Protocol Carrier And Route Protection Gate` is completed. `W-0113` defined `docs/access-token-protocol-carrier-route-protection-gate.md`, `docs/access-token-protocol-carrier-route-protection-gate.zh-CN.md`, and `ADR-0053` without adding implementation. The gate selected request-level validation with the future Protobuf payload wrapper candidate `vibit.authentication.v1.AuthenticatedRequest`, kept the existing Protobuf envelope unchanged, kept WebSocket transport credential-neutral, and required application-owned route policy before protected domain dispatch. `runtime.access_token_protocol_carrier_route_protection_gate` is the repository check rule for this gate.

`M-042 Access Token Protocol Carrier And Route Protection Implementation` is completed. `W-0114` added `proto/vibit/authentication/v1/authenticated_request.proto` and generated `runtime/internal/generated/proto/vibit/authentication/v1/authenticated_request.pb.go` through Buf, added `runtime/internal/app/route_authentication.go`, `runtime/internal/app/authentication/route_validator.go`, and focused app/protocol/WebSocket tests. The Protobuf adapter unwraps `vibit.authentication.v1.AuthenticatedRequest`, validates the access-token proof through application route protection before protected domain dispatch, rejects metadata-only identity for protected routes, keeps `SessionValidated` false, leaves the existing envelope route fields unchanged, and keeps WebSocket transport credential-neutral. It does not add WebSocket handshake authentication, session persistence, startup wiring, repository changes, migrations, dependencies, logout, refresh, cleanup, token rotation, or broader production authentication behavior. The work queue is now blocked at `M-043/W-0115`, a next-direction confirmation gate.

`M-044 Runtime Authentication Startup Composition Gate` is completed. `W-0116` defined `docs/runtime-authentication-startup-composition-gate.md`, `docs/runtime-authentication-startup-composition-gate.zh-CN.md`, and `ADR-0054` after the maintainer selected `wire_runtime_authentication_startup_composition` and asked to keep Nakama and Pitaya as strong references. The gate allows only process startup composition in `runtime/cmd/vibit-server`, with the first composed path limited to `VIBIT_RUNTIME_STORE=postgres`. It keeps the memory store as metadata-only bootstrap behavior, keeps WebSocket transport credential-neutral, keeps the existing Protobuf envelope unchanged, and preserves deferrals for session persistence, WebSocket handshake authentication, authentication command routes, repository changes, migrations, dependencies, logout, refresh, cleanup, token rotation, and broader production behavior. `runtime.authentication_startup_composition_gate` is the repository check rule for this gate.

`M-045 Runtime Authentication Startup Composition Implementation` is completed. `W-0117` implemented startup composition in `runtime/cmd/vibit-server/main.go` with focused tests in `runtime/cmd/vibit-server/main_test.go`. The PostgreSQL runtime path now loads verifier keys through the existing environment loader, builds the existing application authentication service with PostgreSQL unit-of-work, `crypto/rand.Reader`, a startup-owned clock, a startup-owned token record id generator, 1h/default or configured access-token lifetime, and `vibit_gameplay_runtime_requests`/default or configured audience, then injects `app.NewRouteProtector(authentication.NewRouteAccessTokenValidator(service))` into the Protobuf frame handler. The memory runtime path remains metadata-only bootstrap behavior. Do not expand this slice into session persistence, WebSocket handshake authentication, authentication command routes, repository changes, migrations, generated files, dependencies, logout, refresh, cleanup, token rotation, token validation audit mutation, or direct Nakama/Pitaya API compatibility without a later bounded work item. The work queue is now blocked at `M-046/W-0118`, a next-direction confirmation gate.

`M-046 Next Direction Confirmation Gate` is completed. `W-0118` selected `add_authentication_command_protocol_messages_and_login_route_registration` after startup composition, with Nakama and Pitaya still used as reference baselines. `M-047 Authentication Command Protocol And Login Route Gate` is completed. `W-0119` defined `docs/authentication-command-protocol-login-route-gate.md`, `docs/authentication-command-protocol-login-route-gate.zh-CN.md`, and `ADR-0055`. The gate authorizes only the public `runtime.authentication.AuthenticateWithDeviceCredential` Protobuf command route around the existing authentication service. `runtime.authentication_command_protocol_login_route_gate` is the repository check rule for this gate and implementation slice.

`M-048 Authentication Command Protocol And Login Route Implementation` is completed. `W-0120` added `proto/vibit/authentication/v1/authentication.proto` and generated `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go` through Buf, added `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`, added `runtime/internal/app/bootstrap/authentication.go`, registered the public login route only in PostgreSQL startup composition, and added a transaction-wrapper bypass for `runtime.authentication.AuthenticateWithDeviceCredential` because the authentication service owns its own unit-of-work. The existing Protobuf envelope is unchanged, WebSocket transport remains credential-neutral, memory durable login remains unavailable, and repository interfaces, PostgreSQL adapters, migrations, dependencies, session persistence, WebSocket handshake authentication, logout, refresh, cleanup, token rotation, token validation audit mutation, and direct Nakama/Pitaya API compatibility remain deferred. The work queue is now blocked at `M-049/W-0121`, a next-direction confirmation gate.

`M-049 Next Direction Confirmation Gate` is completed. `W-0121` selected `ratify_session_persistence_and_websocket_handshake_authentication` after the public login route, with Nakama and Pitaya still used as reference baselines. `M-050 Session Persistence And WebSocket Handshake Ratification` is completed. `W-0122` defined `docs/session-persistence-websocket-handshake-ratification.md`, `docs/session-persistence-websocket-handshake-ratification.zh-CN.md`, and `ADR-0056`. The ratified current path remains request-level opaque access-token validation through `vibit.authentication.v1.AuthenticatedRequest`; WebSocket transport remains credential-neutral; the existing Protobuf envelope remains unchanged; future connection-level identity is deferred to a first-message protocol/application binding gate; and future durable session persistence is deferred to a PostgreSQL-first schema/repository/migration gate. No session tables, migrations, repositories, dependencies, WebSocket handshake proof carriers, logout, refresh, cleanup, token rotation, reconnect/epoch behavior, memory durable authentication behavior, or direct Nakama/Pitaya API compatibility are authorized by this ratification. `runtime.session_persistence_websocket_handshake_ratification` is the repository check rule. The work queue is now blocked at `M-051/W-0123`, a next-direction confirmation gate.

`M-051 Next Direction Confirmation Gate` is completed. `W-0123` selected `define_first_message_connection_binding_gate` after session and handshake ratification, with Nakama and Pitaya still used as reference baselines. `M-052 First Message Connection Binding Gate` is completed. `W-0124` defined `docs/first-message-connection-binding-gate.md`, `docs/first-message-connection-binding-gate.zh-CN.md`, and `ADR-0057`. The gate selects the future `runtime.authentication.BindConnection` system route with `vibit.authentication.v1.BindConnectionRequest` and `vibit.authentication.v1.BindConnectionResponse` payload candidates. WebSocket transport remains credential-neutral; the existing Protobuf envelope remains unchanged; request-level access-token validation remains the current protected-route path; and future connection-bound identity remains behind an implementation gate. This gate does not authorize Protobuf source changes, generated output, connection binding registries, route-policy use of bound identity, session persistence, repositories, migrations, dependencies, logout/revocation, reconnect/epoch behavior, memory durable authentication behavior, or direct Nakama/Pitaya API compatibility. `runtime.first_message_connection_binding_gate` is the repository check rule. The work queue is now blocked at `M-053/W-0125`, a next-direction confirmation gate.

`M-053 Next Direction Confirmation Gate` is completed. `W-0125` selected `define_first_message_connection_binding_implementation_gate` after the first-message connection binding gate, with Nakama and Pitaya still used as reference baselines. `M-054 First Message Connection Binding Implementation Gate` is completed. `W-0126` defined `docs/first-message-connection-binding-implementation-gate.md`, `docs/first-message-connection-binding-implementation-gate.zh-CN.md`, and `ADR-0058`. The gate defines the future bounded implementation slice for `runtime.authentication.BindConnection`, including planned `BindConnectionRequest`, `BindConnectionResponse`, `ConnectionBindingStatus`, Buf-generated Go output, Protobuf adapter bridge, application binding boundary, PostgreSQL startup composition, public error mapping, and required tests. This gate still does not authorize Go runtime behavior, Protobuf source changes, generated output, connection binding registries, route-policy use of bound identity, session persistence, repositories, migrations, dependencies, logout/revocation, reconnect/epoch behavior, memory durable authentication behavior, or direct Nakama/Pitaya API compatibility. WebSocket transport remains credential-neutral, the existing Protobuf envelope remains unchanged, and request-level access-token validation remains the current protected-route path. `runtime.first_message_connection_binding_implementation_gate` is the repository check rule. The work queue is now blocked at `M-055/W-0127`, a next-direction confirmation gate.

`M-055 Next Direction Confirmation Gate` is completed. `W-0127` selected `implement_first_message_connection_binding`, with Nakama and Pitaya still used as reference baselines. `M-056 First Message Connection Binding Implementation` is completed. `W-0128` implemented the bounded `runtime.authentication.BindConnection` system route slice authorized by `ADR-0058`. The authentication Protobuf source now includes `BindConnectionRequest`, `BindConnectionResponse`, and `ConnectionBindingStatus`, and generated Go output was updated through Buf. The application owns `ConnectionBinder`; the Protobuf adapter handles the system route before ordinary route protection and dispatch; PostgreSQL startup composition injects the binder when authentication is composed; and WebSocket transport carries server-observed connection id and epoch metadata without parsing credentials. This implementation follows Nakama's pattern that authenticated session material precedes authenticated realtime socket use and Pitaya's acceptor/session/handler separation, but does not add direct Nakama/Pitaya API compatibility. It also does not add durable session persistence, connection registries, route-policy use of bound identity for ordinary protected routes, WebSocket handshake authentication, transport credential carriers, repositories, migrations, dependencies, logout/revocation active-connection invalidation, reconnect/resume/duplicate replacement policy, memory durable authentication behavior, presence, rooms, parties, match runtime, or broader game backend behavior. `runtime.first_message_connection_binding_implementation` is the repository check rule.

`M-057 Next Direction Confirmation Gate` is completed. `W-0129` selected `define_postgres_session_persistence_schema_gate` after first-message connection binding implementation, with Nakama and Pitaya still used as reference baselines. `M-058 PostgreSQL Session Persistence Schema Gate` is completed. `W-0130` defined `docs/postgres-session-persistence-schema-gate.md`, `docs/postgres-session-persistence-schema-gate.zh-CN.md`, and `ADR-0059`. The gate selects PostgreSQL as the first durable runtime session target, `runtime_sessions` as the future logical table candidate, and `runtime/migrations/postgres/000005_create_runtime_sessions.sql` as the future migration source candidate. It adapts Nakama's first-class session lifecycle and Pitaya's session-context separation, but it does not add SQL migration source, create session tables, add session repositories, add PostgreSQL adapters, create runtime session behavior, set `RequestIdentity.SessionValidated` true, change route policy, change the existing Protobuf envelope, add WebSocket handshake authentication, parse transport credential carriers, add logout/revocation active-connection invalidation, add reconnect/epoch behavior, add dependencies, or provide direct Nakama/Pitaya API compatibility. `runtime.postgres_session_persistence_schema_gate` is the repository check rule. The work queue is now blocked at `M-059/W-0131`, a next-direction confirmation gate.

`M-059 Next Direction Confirmation Gate` is completed. `W-0131` selected `implement_runtime_sessions_migration_source` after the PostgreSQL session persistence schema gate, with Nakama and Pitaya still used as reference baselines. `M-060 Runtime Sessions Migration Source` is completed. `W-0132` added `runtime/migrations/postgres/000005_create_runtime_sessions.sql` and `ADR-0060`. The migration creates the PostgreSQL `runtime_sessions` lifecycle table with actor/player identity, session status, issued/expires/last_seen timestamps, optional revocation fields, and optional `authentication_access_tokens(token_record_id)` linkage. It does not store raw access-token text, raw credentials, token digests, credential digests, WebSocket connection state, or connection registry rows. It also does not add session repositories, PostgreSQL session adapters, runtime session validation, session creation at login or BindConnection, route-policy use of session or bound identity, WebSocket handshake authentication, logout/revocation active-connection invalidation, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility. `runtime.runtime_sessions_migration_source` is the repository check rule. The work queue is now blocked at `M-061/W-0133`, a next-direction confirmation gate.

`M-061 Next Direction Confirmation Gate` is completed. `W-0133` selected `define_session_repository_boundary` after the runtime sessions migration source, with Nakama and Pitaya still used as reference baselines. `M-062 Session Repository Boundary Gate` is completed. `W-0134` defined `docs/session-repository-boundary.md`, `docs/session-repository-boundary.zh-CN.md`, and `ADR-0061`. The gate records `runtime/internal/app/session` as the future storage-neutral repository owner candidate and `runtime/internal/platform/persistence/postgres` as the future PostgreSQL adapter owner. It defines candidate lifecycle capabilities such as `CreateRuntimeSession`, `FindActiveSessionByID`, `UpdateRuntimeSessionLastSeen`, `MarkRuntimeSessionExpired`, and `RevokeRuntimeSession`, but it does not add Go repository code, PostgreSQL adapter behavior, runtime session validation, session creation at login or BindConnection, route-policy use of persisted session or bound identity, WebSocket handshake authentication, transport credential carriers, Protobuf session messages, logout/revocation active-connection invalidation, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility. `runtime.session_repository_boundary` is the repository check rule. The work queue is now blocked at `M-063/W-0135`, a next-direction confirmation gate.

`M-063 Next Direction Confirmation Gate` is completed. `W-0135` selected `implement_session_repository_interface` after the session repository boundary, with Nakama and Pitaya still used as reference baselines. `M-064 Session Repository Interface Implementation` is completed. `W-0136` added `runtime/internal/app/session/repository.go`, `runtime/internal/app/session/repository_test.go`, and `ADR-0062`. The package defines storage-neutral runtime session lifecycle value types, active/expired/revoked status vocabulary, first-posture player actor sessions, `runtime/internal/app/session.Repository`, lifecycle query and mutation types, and normalization helpers for creation, lookup, active lookup, last-seen update, expiration, revocation, and bounded active-session listing. It does not add PostgreSQL session adapters, SQL query execution, unit-of-work factory wiring, runtime session creation or validation, `RequestIdentity.SessionValidated = true`, route-policy use of persisted session or bound identity, WebSocket handshake authentication, transport credential carriers, Protobuf session messages, generated output, logout/revocation active-connection invalidation, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility. `runtime.session_repository_interface_implementation` is the repository check rule. The work queue is now blocked at `M-065/W-0137`, a next-direction confirmation gate.

`M-065 Next Direction Confirmation Gate` is completed. `W-0137` selected `define_session_postgresql_adapter_gate` after the session repository interface implementation, with Nakama and Pitaya still used as reference baselines. `M-066 Session PostgreSQL Adapter Gate` is completed. `W-0138` defined `docs/session-postgresql-adapter-gate.md`, `docs/session-postgresql-adapter-gate.zh-CN.md`, and `ADR-0063`. The gate records `runtime/internal/platform/persistence/postgres` as the future adapter owner for `runtime/internal/app/session.Repository`, defines future SQL shape, transaction handoff, error mapping, redaction, and adapter test requirements, but it does not add PostgreSQL session adapter files, SQL execution, unit-of-work factory wiring, runtime session creation or validation, `RequestIdentity.SessionValidated = true`, route-policy use of persisted session or bound identity, WebSocket handshake authentication, transport credential carriers, Protobuf session messages, generated output, logout/revocation active-connection invalidation, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility. `runtime.session_postgresql_adapter_gate` is the repository check rule. The work queue is now blocked at `M-067/W-0139`, a next-direction confirmation gate.

`M-067 Next Direction Confirmation Gate` is completed. `W-0139` selected `implement_session_postgresql_adapter` after the session PostgreSQL adapter gate, with Nakama and Pitaya still used as reference baselines. `M-068 Session PostgreSQL Adapter Implementation` is completed. `W-0140` added `runtime/internal/platform/persistence/postgres/session_repository.go`, `runtime/internal/platform/persistence/postgres/session_repository_test.go`, `UnitOfWork.NewSessionRepository()`, and `ADR-0064`. The adapter implements `runtime/internal/app/session.Repository` against `runtime_sessions`, normalizes inputs and returned records through the session package, maps PostgreSQL errors to redacted typed errors, and remains transaction-neutral. It does not create sessions at login or BindConnection, validate runtime sessions, set `RequestIdentity.SessionValidated = true`, change WebSocket handshake authentication, add transport credential carriers, add Protobuf session messages, change the existing envelope, use persisted session or bound identity for route policy, invalidate active connections on logout/revocation, add reconnect/epoch behavior, add cleanup jobs, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility. `runtime.session_postgresql_adapter_implementation` is the repository check rule. The work queue is now blocked at `M-069/W-0141`, a next-direction confirmation gate.

`M-069 Next Direction Confirmation Gate` is completed. `W-0141` selected `define_runtime_session_validation_gate` after the session PostgreSQL adapter implementation, with Nakama and Pitaya still used as reference baselines. `M-070 Runtime Session Validation Gate` is completed. `W-0142` defined `docs/runtime-session-validation-gate.md`, `docs/runtime-session-validation-gate.zh-CN.md`, and `ADR-0065`. The gate makes future runtime session validation application-owned under `runtime/internal/app`, requires already-validated actor identity before trusting persisted `runtime_sessions` rows, defines active/expired/revoked and actor-mismatch public failure collapse, and keeps request identity handoff explicit. It does not implement runtime session validation, set `RequestIdentity.SessionValidated = true`, create sessions at login or BindConnection, change route policy, change WebSocket handshake authentication, add transport credential carriers, add Protobuf session messages, change the existing envelope, add logout/revocation active-connection invalidation, add reconnect/epoch behavior, add cleanup jobs, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility. `runtime.runtime_session_validation_gate` is the repository check rule. The work queue is now blocked at `M-071/W-0143`, a next-direction confirmation gate.

`M-071 Next Direction Confirmation Gate` is completed. `W-0143` selected `implement_runtime_session_validation` after the runtime session validation gate, with Nakama and Pitaya still used as reference baselines. `M-072 Runtime Session Validation Implementation` is completed. `W-0144` added `runtime/internal/app/runtime_session_validator.go`, `runtime/internal/app/runtime_session_validator_test.go`, and `ADR-0066`. The application-owned `PersistentSessionValidator` uses `runtime/internal/app/session.Repository.FindActiveSessionByID`, requires already validated player identity before trusting a persisted runtime session row, validates active/unexpired actor-player match, collapses public invalid-session failures to a stable redacted reason, and sets `RequestIdentity.SessionValidated = true` only after durable validation succeeds. It is not wired into startup or route policy, does not create sessions, does not update `last_seen_at`, does not change WebSocket or Protobuf behavior, and does not add logout/revocation active-connection behavior, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility. `runtime.runtime_session_validation_implementation` is the repository check rule. The work queue is now blocked at `M-073/W-0145`, a next-direction confirmation gate.

`M-073 Next Direction Confirmation Gate` is completed. `W-0145` selected `define_session_creation_composition_gate` after runtime session validation implementation, with Nakama and Pitaya still used as reference baselines. `M-074 Session Creation Composition Gate` is completed. `W-0146` defined `docs/session-creation-composition-gate.md`, `docs/session-creation-composition-gate.zh-CN.md`, and `ADR-0067`. The gate makes future durable runtime session creation application-owned under `runtime/internal/app`, identifies `AuthenticateWithDeviceCredential` as the first future login-time composition candidate, records future `session.Repository.CreateRuntimeSession` use through unit-of-work capabilities, keeps `access_token_record_id` as private server metadata, and defines session id, lifetime, redaction, and future test expectations. It does not implement session creation, modify authentication service behavior, generate session ids, create sessions at login or BindConnection, change runtime session validation or route policy, change WebSocket or Protobuf behavior, add logout/revocation active-connection behavior, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility. `runtime.session_creation_composition_gate` is the repository check rule. The work queue is now blocked at `M-075/W-0147`, a next-direction confirmation gate.

`M-075 Next Direction Confirmation Gate` is completed. `W-0147` selected `implement_session_creation_composition` after the session creation composition gate, with Nakama and Pitaya still used as reference baselines. `M-076 Session Creation Composition Implementation` is completed. `W-0148` updated `runtime/internal/app/authentication/service.go`, `runtime/internal/app/authentication/service_test.go`, `runtime/cmd/vibit-server/main.go`, `runtime/cmd/vibit-server/main_test.go`, and `ADR-0068`. Successful device-credential login now stores the access-token record and creates one active durable runtime session linked to `access_token_record_id` in the same unit of work, with server-owned session id generation and token-aligned first lifetime. It does not expose session ids through Protobuf, change the existing envelope, change route policy, set `SessionValidated` true during token validation, change WebSocket handshake authentication, add transport credential carriers, add logout/revocation active-connection behavior, add reconnect/epoch behavior, add dependencies, add memory durable session behavior, or add direct Nakama/Pitaya API compatibility. `runtime.session_creation_composition_implementation` is the repository check rule. The work queue is now blocked at `M-077/W-0149`, a next-direction confirmation gate.

`M-077 Next Direction Confirmation Gate` is completed. `W-0149` selected `define_bound_identity_route_policy_gate` after session creation composition implementation, with Nakama and Pitaya still used as reference baselines. `M-078 Bound Identity Route Policy Gate` is completed. `W-0150` defined `docs/bound-identity-route-policy-gate.md`, `docs/bound-identity-route-policy-gate.zh-CN.md`, and `ADR-0069`. The gate makes future route-policy use of request-token, bound-connection, session-validated, and bound-session identity application-owned under `runtime/internal/app`, route-scoped, fail-closed, and redacted. The recommended first implementation posture keeps ordinary protected domain routes on request-level access-token proof; bound identity and session-validated identity are explicit future policy families only. It does not implement route-policy use of bound or session identity, remove per-request token proof, change WebSocket handshake authentication, add transport credential carriers, expose session ids through Protobuf, change the existing envelope, add logout/revocation active-connection behavior, add reconnect/epoch behavior, add dependencies, add memory durable session behavior, or add direct Nakama/Pitaya API compatibility. `runtime.bound_identity_route_policy_gate` is the repository check rule. The work queue is now blocked at `M-079/W-0151`, a next-direction confirmation gate.

`M-079 Next Direction Confirmation Gate` is completed. `W-0151` selected `implement_bound_identity_route_policy` after the bound identity route policy gate, with Nakama and Pitaya still used as reference baselines. `M-080 Bound Identity Route Policy Implementation` is completed. `W-0152` updated `runtime/internal/app/route_authentication.go`, `runtime/internal/app/route_authentication_test.go`, and `ADR-0070`. The application route protector now has explicit route policy families for `public`, `request_token_required`, `bound_connection_required`, `session_validated_required`, and `bound_session_required`. `runtime.authentication.AuthenticateWithDeviceCredential` remains the explicit public route; ordinary protected domain routes still default to request-level access-token proof and still clear `SessionValidated` after token validation. Bound/session identities can satisfy only explicitly classified routes, metadata-only identity is rejected for every protected policy family, and bound-session routes require identity-source agreement. This does not wire WebSocket handshake authentication, transport credential carriers, Protobuf session carriers, existing envelope changes, connection registries, persistent session validation into frame handling, logout/revocation active-connection behavior, reconnect/epoch behavior, dependencies, memory durable session behavior, broader game backend behavior, or direct Nakama/Pitaya API compatibility. `runtime.bound_identity_route_policy_implementation` is the repository check rule. The work queue is now blocked at `M-081/W-0153`, a next-direction confirmation gate.

`M-081 Next Direction Confirmation Gate` is completed. `W-0153` selected `define_logout_revocation_active_connection_gate` after bound identity route policy implementation, with Nakama and Pitaya still used as reference baselines. `M-082 Logout Revocation Active Connection Gate` is completed. `W-0154` defined `docs/logout-revocation-active-connection-gate.md`, `docs/logout-revocation-active-connection-gate.zh-CN.md`, and `ADR-0071`. The gate makes future logout/revocation active-connection policy application-owned under `runtime/internal/app`, keeps presented-token logout, runtime session revocation, and active socket invalidation as separate future decisions, recommends presented-token logout as the first future scope, and requires an explicit connection registry before targeting open sockets. It does not implement `LogoutAccessToken`, revoke tokens, revoke runtime sessions, close WebSocket connections, add active connection registries, add WebSocket close policy, add Protobuf logout routes, add protocol session carriers, change the existing envelope, add reconnect/epoch behavior, add cleanup jobs, add dependencies, add memory durable session behavior, or add direct Nakama/Pitaya API compatibility. `runtime.logout_revocation_active_connection_gate` is the repository check rule. The work queue is now blocked at `M-083/W-0155`, a next-direction confirmation gate.

`M-083 Next Direction Confirmation Gate` is completed. `W-0155` selected `define_logout_access_token_behavior_gate` after the logout/revocation active-connection gate, with Nakama and Pitaya still used as reference baselines. `M-084 Logout Access Token Behavior Gate` is completed. `W-0156` defined `docs/logout-access-token-behavior-gate.md`, `docs/logout-access-token-behavior-gate.zh-CN.md`, and `ADR-0072`. The gate made future `LogoutAccessToken` behavior application-owned under `runtime/internal/app/authentication`, kept the first logout scope to `presented_access_token_only`, required lookup digest plus verifier digest comparison before revocation, and required success only after unit-of-work commit. `runtime.logout_access_token_behavior_gate` is the repository check rule.

`M-085 Next Direction Confirmation Gate` is completed. `W-0157` selected `implement_logout_access_token_behavior` after the logout access-token behavior gate. `M-086 Logout Access Token Behavior Implementation` is completed. `W-0158` implemented `LogoutAccessToken` in `runtime/internal/app/authentication/service.go` and focused tests in `runtime/internal/app/authentication/service_test.go`. The implementation rejects missing or malformed opaque access-token proof before unit of work, computes token lookup digest, requires active access-token posture, compares verifier digest before revocation, calls `RevokeToken` once with `logout_presented_access_token`, and returns `LogoutStatusRevoked` only after commit. It does not revoke runtime sessions, close WebSocket connections, add connection registries, add WebSocket close policy, add Protobuf logout routes, add protocol session carriers, change the existing envelope, add refresh, logout-all, admin revocation, reconnect/epoch behavior, cleanup jobs, dependencies, memory durable session behavior, broader game backend behavior, or direct Nakama/Pitaya API compatibility. `runtime.logout_access_token_behavior_implementation` is the repository check rule. The work queue then blocked at `M-087/W-0159`, a next-direction confirmation gate.

`M-087 Next Direction Confirmation Gate` is completed. `W-0159` selected `define_active_connection_registry_gate` after logout access-token behavior implementation, with Nakama and Pitaya still used as reference baselines. `M-088 Active Connection Registry Gate` is completed. `W-0160` defined `docs/active-connection-registry-gate.md`, `docs/active-connection-registry-gate.zh-CN.md`, and `ADR-0074`. The gate makes future active connection registry behavior application-owned under `runtime/internal/app/connection`, selects a single-process in-memory non-durable first posture, treats registry records as server-observed connection state with validated identity linkage rather than client proof, and keeps WebSocket transport credential-neutral. It does not implement a registry, close WebSocket connections, add kick/disconnect behavior, revoke runtime sessions, add WebSocket close policy, add Protobuf logout routes, add protocol session carriers, change the existing envelope, add reconnect/epoch behavior, add durable/distributed registry storage, add dependencies, add memory durable session behavior, broaden game backend behavior, or add direct Nakama/Pitaya API compatibility. `runtime.active_connection_registry_gate` is the repository check rule. The work queue is now blocked at `M-089/W-0161`, a next-direction confirmation gate.

`M-089 Next Direction Confirmation Gate` is completed. `W-0161` selected `implement_active_connection_registry_single_process` after the active connection registry gate, with Nakama and Pitaya still used as reference baselines. `M-090 Active Connection Registry Single Process Implementation` is completed. `W-0162` added `runtime/internal/app/connection/registry.go` and `runtime/internal/app/connection/registry_test.go`, plus `ADR-0075`. The registry is application-owned, single-process, in-memory, and non-durable; it registers server-observed connection id and epoch, binds validated player identity with optional runtime session id and access-token record id, marks records closed or invalidated, finds records by connection id/epoch, and lists active bound records by player/session/token record. It does not wire startup or transport handoff, close WebSocket connections, add kick/disconnect behavior, revoke runtime sessions, replace duplicates, add reconnect/epoch behavior, add Protobuf logout routes, add protocol session carriers, change the existing envelope, add durable/distributed registry storage, add dependencies, add memory durable session behavior, broaden game backend behavior, or add direct Nakama/Pitaya API compatibility. `runtime.active_connection_registry_single_process_implementation` is the repository check rule. The work queue is now blocked at `M-091/W-0163`, a next-direction confirmation gate.

The runtime authentication implementation boundary standard is `docs/runtime-authentication-implementation-boundary.md`, with `docs/runtime-authentication-implementation-boundary.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0036` records the boundary decision. Future runtime authentication is application-owned under `runtime/internal/app`; it must use `authentication.Repository` through the application unit-of-work boundary, keep the PostgreSQL adapter persistence-only, and convert validated proof into `RequestIdentity` before domain dispatch. Token generation, verifier comparison, login execution, access-token validation, logout execution, cleanup jobs, Protobuf authentication messages, WebSocket proof carriers, generated authentication shapes, and authentication dependencies remain separate gates. `runtime.authentication_implementation_boundary` is the repository check rule for this boundary.

`M-092 WebSocket Close Policy Gate` is completed. `W-0164` defined `docs/websocket-close-policy-gate.md`, `docs/websocket-close-policy-gate.zh-CN.md`, and `ADR-0076`. Future WebSocket close policy is application-owned under `runtime/internal/app`; the active connection registry remains target state, not policy; and WebSocket transport may only own a future narrow concrete close handoff after application policy emits a redacted close intent. Do not expand this gate into transport close handoff code, close codes, close reason text, kick/disconnect behavior, logout-triggered socket close, runtime session revocation, duplicate replacement, reconnect/epoch behavior, Protobuf logout routes, protocol session carriers, existing envelope changes, WebSocket handshake authentication, transport credential carriers, durable/distributed registry storage, dependencies, broader game backend behavior, or direct Nakama/Pitaya API compatibility. `runtime.websocket_close_policy_gate` is the repository check rule.

`M-094 WebSocket Close Policy Single Process Implementation` is completed. `W-0166` added `runtime/internal/app/connection/close_policy.go`, `runtime/internal/app/connection/close_policy_test.go`, and `ADR-0077`. The policy is application-owned, single-process, and registry-backed; it targets only active bound registry records by connection id/epoch, player id, runtime session id, or access-token record id; it marks matched records invalidated and emits redacted `CloseIntent` values with `mark_invalidated_only`. It does not close concrete WebSocket sockets, add transport close handoff, choose close codes or reason text, add protocol close messages, change logout behavior, revoke runtime sessions, replace duplicate connections, add reconnect/epoch behavior, add Protobuf logout routes, add protocol session carriers, change generated output, add durable/distributed registry storage, add dependencies, broaden game backend behavior, or add direct Nakama/Pitaya API compatibility. `runtime.websocket_close_policy_single_process_implementation` is the repository check rule. The work queue is now blocked at `M-095/W-0167`, a next-direction confirmation gate.

The authentication generated contract shape timing standard is `docs/authentication-generated-contract-shape-timing.md`, with `docs/authentication-generated-contract-shape-timing.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0038` records the timing decision. Generated Go authentication contract shapes now exist after the runtime authentication implementation boundary and before service interfaces, using `contracts/runtime/authentication/` as source and `runtime/internal/generated/contracts/runtime/authentication/` as the output root. Generated files remain immutable and metadata-only.

The application authentication service interface boundary standard is `docs/application-authentication-service-interface-boundary.md`, with `docs/application-authentication-service-interface-boundary.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0039` records the boundary decision. Future authentication service interfaces are application-owned under `runtime/internal/app`; generated authentication shapes inform service-level request/result vocabulary; service behavior may use `authentication.Repository` only through the application unit-of-work boundary; validated proof must become `RequestIdentity` before domain dispatch. This boundary does not authorize Go service code or runtime authentication behavior. `runtime.application_authentication_service_interface_boundary` is the repository check rule for this boundary.

`ADR-0037` closes the runtime authentication implementation boundary planning milestone and opens the generated authentication contract shape timing gate. `ADR-0038` completes the timing decision. `W-0089` completes generator/check support plus metadata-only generated authentication shape output. `ADR-0039` and `W-0090` complete the service-interface boundary step. `ADR-0040` and `W-0091` complete the verifier algorithm/redaction step. `ADR-0041` and `W-0092` complete the secret configuration/verifier key loading preparation step. `ADR-0042` and `W-0093` complete the material generation preparation step. `ADR-0043` and `W-0094` complete the verifier digest computation and comparison preparation step. `ADR-0044` and `W-0095` complete the implementation readiness step. `ADR-0045` and `W-0096` complete the local verifier key configuration loading gate. `W-0097` completes the explicit in-memory verifier key set validator implementation slice. `ADR-0046` and `W-0098` complete the environment verifier key loader gate. `W-0099` completes the environment verifier key loader implementation slice. `ADR-0047` and `W-0100` complete the token and credential material generation implementation gate. `W-0101` completes the token and credential material generation helper implementation slice. `ADR-0048` and `W-0102` complete the verifier digest helper implementation gate. `W-0103` completes the verifier digest computation helper implementation slice. `ADR-0049` and `W-0104` complete the verifier digest comparison helper gate. `W-0105` completes the verifier digest comparison helper implementation slice. `ADR-0050` and `W-0106` complete the authentication service behavior implementation gate.

The work continuation standard is `docs/workflow.md`, with `docs/workflow.zh-CN.md` as the paired Simplified Chinese translation. The machine-readable work queue is `.arch/work-items.yaml`. When the maintainer says "continue" or "继续", interpret that as advancing one `next_ready` work item unless blocked or confirmation is required. When the maintainer asks to continue multiple steps, advance up to that many work items in order, stopping at blockers, verification failures, ask-first boundaries, or maintainer redirection.

Current executable tooling:

```bash
node tools/vibit --help
node tools/vibit check all
node tools/vibit check all --json
node tools/vibit check schemas
node tools/vibit check schemas --json
node tools/vibit check memory
node tools/vibit check memory --json
node tools/vibit check contracts
node tools/vibit check contracts --json
node tools/vibit check protocol
node tools/vibit check protocol --json
node tools/vibit check generated
node tools/vibit check generated --json
node tools/vibit check agent-tooling
node tools/vibit check agent-tooling --json
node tools/vibit check migrations
node tools/vibit check migrations --json
node tools/vibit check postgres-env
node tools/vibit check postgres-env --json
node tools/vibit check runtime
node tools/vibit check runtime --json
node tools/vibit check work
node tools/vibit check work --json
node tools/vibit inspect module <module>
node tools/vibit inspect boundary --from <module> --to <module>
node tools/vibit inspect contract --module <module> --type <type> --id <id>
node tools/vibit inspect contracts
node tools/vibit inspect contracts --module <module>
node tools/vibit inspect generated
node tools/vibit inspect next
node tools/vibit inspect reference
node tools/vibit inspect change <change-id>
node tools/vibit inspect work
node tools/vibit inspect work --json
node tools/vibit inspect memory
node tools/vibit inspect rule <rule-id>
node tools/vibit inspect rules
node tools/vibit inspect rules --category <category>
node tools/vibit check architecture
node tools/vibit check architecture --json
node tools/vibit check change <change-id>
node tools/vibit check change <change-id> --json
node tools/vibit check module <module>
node tools/vibit check module <module> --json
node tools/vibit generate module <module>
node tools/vibit generate contract-shapes <module|all>
```

Use `node tools/vibit check all` as the default repository verification command when CLI tooling is available.

The current CLI is Node.js standard-library tooling only. Do not treat CLI implementation language as the server runtime language.

Use `--json` when an agent needs machine-readable check results during intake, verification, or handoff.

Every JSON check result item should include `rule_id` and `artifact`. Repository-relative JSON paths in fields such as `artifact`, `path`, `source`, and `output` must use forward slashes on every platform, including Windows. Treat `check all --json` as a compact overview, then run the specific failing check with `--json` for full detail.

Use `node tools/vibit check memory` when conversation logs or Agent Decision Records are added or changed.

Use `node tools/vibit check contracts` when contract source files or `.arch/contracts.yaml` are added or changed.

Use `node tools/vibit check protocol` before creating or changing `.proto` files, generated Protobuf output, or protocol generation rules. Missing `.proto` files may pass while protocol sources are still planned, but once a `.proto` file exists it must align with registered command, query, and event contracts.

Use `node tools/vibit check generated` when generated files, module manifest `generated` declarations, generated output standards, or Go Protobuf generated output are added or changed.

Use `docs/agent-tooling.md` and `node tools/vibit check agent-tooling` when agent-facing inspection, generation, or verification commands are added or changed.

Use `node tools/vibit check migrations` when SQL migration sources, migration ownership manifests, migration guidance, or persistence migration standards are added or changed. This check validates PostgreSQL migration naming, goose markers, SQL-first boundaries, owning-module traces, and first inventory migration table references.

Use `node tools/vibit check postgres-env` when disposable PostgreSQL verification environment standards, live PostgreSQL verification guidance, or persistence verification environment manifests are added or changed. This is a static standards check; it must not connect to PostgreSQL or require Docker, Podman, cloud PostgreSQL, or another service manager.

Use `node tools/vibit check runtime` when runtime module behavior, runtime adapter boundaries, runtime guidance, or tests are added or changed. Before the Go runtime exists, this check should pass as not applicable because runtime implementation has not started. After `runtime/go.mod` exists but before Go source files exist, this check should verify the ADR-0014 skeleton and the ADR-0018 runtime protocol adapter boundary, and pass without running `go test`. Once Go source files exist, runtime checks require Go test files and a local Go toolchain.

Use `node tools/vibit inspect rule runtime.authentication_token_session_boundary --json` when a runtime check fails on authentication, token, credential, external identity, session persistence, Protobuf envelope authentication, WebSocket handshake authentication, runtime player handler, or WebSocket route boundaries.

Use `node tools/vibit inspect rule runtime.selected_login_token_boundary --json` when a runtime check fails on the selected `device_credential_login`, opaque access-token, explicit request proof payload, generated authentication shape metadata boundary, authentication Protobuf deferral, WebSocket carrier deferral, schema-gate, migration, repository, adapter, or dependency boundary.

Use `node tools/vibit inspect work` before interpreting a continuation request. Use `node tools/vibit check work` when `.arch/work-items.yaml`, workflow docs, or work item state changes. The default continuation unit is one work item.

Use `node tools/vibit inspect next --json` when the immediate continuation step is unclear.

Use `node tools/vibit inspect contracts --json` when an agent needs the full registered contract index before contract, generator, or runtime planning work.

Use `node tools/vibit inspect contract --module <module> --type <type> --id <id>` during intake when an agent needs one contract's registry entry, source summary, module manifest declaration, and consistency status as JSON.

Use `node tools/vibit inspect generated --json` before editing generated output standards, generated output checks, or generator behavior.

Use `node tools/vibit inspect reference --json` before planning new game server capability families from Nakama or Pitaya reference context.

Use `node tools/vibit generate contract-shapes all` to regenerate Go contract shape files from semantic contract manifests. Do not hand-edit files under `runtime/internal/generated/contracts/`.

Use `node tools/vibit inspect change <change-id>` during intake or handoff when a change spec exists and an agent needs a structured summary of its files, metadata, affected modules, and verification state.

Use `node tools/vibit inspect memory` when an agent needs a structured index of change specs, conversation logs, and Agent Decision Records before choosing which artifacts to read in full.

Use `rules/check-rules.json` to interpret check result `rule_id` values.

Use `node tools/vibit inspect rule <rule-id>` when only one rule's metadata is needed.

Use `node tools/vibit inspect rules --category <category>` to discover rules by category.

Use `.arch/runtime.yaml` as the machine-readable intake point for runtime readiness. It links the ADRs that govern language, server instance model, contract and generation boundary, client protocol, dependency adoption, and first proof slice.

Use `.arch/protocol.yaml` as the machine-readable intake point for game protocol framework decisions. It links `ADR-0015` and defines the first WebSocket Protobuf envelope, route fields, session model, target scopes, authority rules, error model, and first inventory slice protocol scope.

Use `ADR-0016`, `ADR-0017`, `buf.yaml`, `buf.gen.yaml`, `proto/README.md`, and `docs/generated-output.md` before changing the protocol envelope, inventory Protobuf source, Buf generation configuration, generated output checks, or generated Go Protobuf output path.

Use `ADR-0018` and `docs/runtime-protocol-adapter.md` before changing runtime code that sits between WebSocket transport, Protobuf protocol adaptation, application dispatch, generated code, and domain modules.

Use `.arch/dependencies.yaml` as the machine-readable intake point before adding foundational dependencies. Use `docs/dependency-adoption.md` and `docs/_templates/dependency-adoption.md` for adoption records.

Use `docs/postgresql-verification-environment.md` before adding live PostgreSQL migration checks, repository integration tests, transaction-runner integration tests, or persistent-runtime end-to-end checks. Live PostgreSQL verification is opt-in through `VIBIT_POSTGRES_TEST_DSN`; default repository checks must not require a running database.

Use `.arch/reference.yaml` as the machine-readable intake point for Nakama/Pitaya reference alignment. Nakama is the primary reference for broad game backend product capability surface. Pitaya is the primary reference for Go game server framework architecture vocabulary. Preserve vibit's Agent-Native constraints when adapting reference patterns, and record why a reference pattern is adopted, adapted, or rejected.

Use `docs/authentication-token-session-validation.md`, `docs/authentication-proof-token-session-contract-dimensions.md`, `docs/credential-storage-external-identity-linking-boundaries.md`, `docs/session-persistence-websocket-handshake-decision-gates.md`, `docs/login-method-token-format-ratification.md`, `ADR-0023`, `ADR-0024`, and `runtime.authentication_token_session_boundary` before changing authentication proof, login methods, token behavior, credential storage, external identity linking, runtime session persistence, request identity trust, Protobuf envelope authentication behavior, WebSocket handshake authentication, runtime player handlers, or WebSocket routes. The design standard separates authentication proof, login methods, tokens, credentials, external identity links, runtime sessions, request identity, transport metadata, envelope metadata, and player account lifecycle. The dimensions standard ratifies actor kinds, validation statuses, proof statuses, failure classes, retryability, request identity handoff, session error dimensions, session permission dimensions, and validation event dimensions. The credential/external identity boundary standard keeps player lifecycle tables credential-free and provider-subject-free while deferring login-method families, credential schema, provider subject semantics, account linking, recovery, merge behavior, and dependencies. The session/handshake gates standard keeps request-level, first-message, handshake-level, every-request, and hybrid validation as future choices until separately selected. The login/token ratification standard defines how to compare and select the first login-method set, token model, proof carrier posture, lifecycle semantics, schema gates, checks, and implementation queue without granting implementation permission.

Use `.arch/work-items.yaml` as the machine-readable intake point for continuation. Work item IDs such as `W-0007` are execution steps; ADR IDs remain architectural decisions; change spec IDs remain concrete execution records; Git hashes remain repository snapshots; versions remain release identifiers.

Use `docs/v0.1-alpha-goal.md` as the short-term release-target intake. The current release state is source-first `v0.1.0-alpha.1` with the first alpha user discovery loop, feedback intake surface, product maturity milestones, prototype-ready execution plan, and local development path gate defined. The target is `v0.1 alpha`, product maturity milestones live at `docs/product-maturity-milestones.md`, the prototype-ready execution plan lives at `docs/prototype-ready-foundation-execution-plan.md`, the local development path gate lives at `docs/prototype-ready-local-development-path-gate.md`, and the next work item is `W-0200 prototype_ready_local_development_path_package`. The alpha acceptance checklist lives at `docs/alpha-acceptance-checklist.md`, the packaged local alpha developer flow lives at `docs/alpha-developer-flow.md`, the release publishing decision gate lives at `docs/release-publishing-decision-gate.md`, the release execution preparation gate lives at `docs/release-execution-preparation-gate.md`, the release execution authorization gate lives at `docs/release-execution-authorization-gate.md`, the release execution maintainer decision lives at `docs/release-execution-maintainer-decision.md`, the release identifier plan lives at `docs/release-identifier-artifact-plan.md`, final authorization lives at `docs/release-execution-final-authorization.md`, the first alpha user discovery loop lives at `docs/first-alpha-user-discovery-loop.md`, and first alpha feedback intake lives at `docs/first-alpha-feedback-intake-surfaces.md`. Do not treat source-first alpha authorization, the user discovery loop, feedback intake, prototype-ready execution plan, or local development path gate as permission to add release binaries, packages, containers, checksums, signing/provenance artifacts, hosted deployments, direct Nakama/Pitaya API compatibility, runtime behavior, protocol routes, generated output, migrations, dependencies, public announcements beyond the GitHub release record, paid promotion, or to skip `.arch/work-items.yaml`.

Use `ADR-0014` before changing Go runtime files. The first Go module lives at `runtime/go.mod` with module path `github.com/iceiko/vibit/runtime`. Keep process startup under `runtime/cmd/vibit-server/`, application dispatch and composition under `runtime/internal/app/`, platform adapters under `runtime/internal/platform/`, handwritten domain module logic under `runtime/internal/modules/<module>/`, generated Go contract shapes under `runtime/internal/generated/contracts/`, generated Go Protobuf output under `runtime/internal/generated/proto/`, SQL-first PostgreSQL migrations under `runtime/migrations/postgres/`, and Protobuf source files under repository-root `proto/vibit/<module>/v1/`.

## 4. Documentation Rules

English is the canonical documentation language.

Every public-facing document should have:

- An English source document
- A Simplified Chinese human-readable translation

Naming examples:

```text
CONSTITUTION.md
CONSTITUTION.zh-CN.md
AGENTS.md
AGENTS.zh-CN.md
docs/<name>.md
docs/<name>.zh-CN.md
.arch/README.md
.arch/README.zh-CN.md
```

Rules:

- **Core documents** (`CONSTITUTION.md`, `AGENTS.md`, `.arch/README.md`): Update Chinese translation synchronously in the same change.
- **Feature-level and milestone-specific documents** (`docs/*.md`, `decisions/ADR-*.md`, change specs): Chinese translation may be **deferred/asynchronous**. Mark the English source with `Translation: pending` if the translation is not updated in the same change.
- If the translation cannot be updated in the same change, mark it clearly as out of date.
- Keep machine-readable identifiers in English.
- Use English for code identifiers, module names, commands, events, permissions, and errors unless a strong domain reason exists.
- Preserve meaning in translation. Do not force literal word-by-word translation when it reduces clarity.

## 5. Standard Change Workflow

For every non-trivial feature, bug fix, migration, refactor, or standard change:

1. Clarify the requirement.
2. Identify affected modules and contracts.
3. Write or update the change spec when the change is large enough to need durable context.
4. Update schemas, manifests, or contracts before implementation when public behavior changes.
5. Generate repeatable structure when generators exist.
6. Implement only inside the declared boundary.
7. Add or update focused tests.
8. Run relevant verification commands.
9. Update documentation and translations.
10. Record what was verified and what was not verified.

For early design-only changes, steps involving code, tests, generators, and verification may be not applicable. Say that explicitly.

## 6. Architecture Rules

Prefer designs that:

- Give agents less ambiguous context
- Create stronger module boundaries
- Make behavior easier to verify
- Make contracts explicit
- Reduce hidden coupling
- Support code generation
- Remain practical for human developers

Do not rely on maintainer memory for architecture rules. If a rule matters, it should eventually be represented in a document, schema, manifest, test, generator, or architecture check.

## 7. Module Rules

When modules exist, each module should declare:

- What it owns
- What it does not own
- Public commands
- Public queries
- Published events
- Subscribed events
- Allowed dependencies
- Forbidden dependencies
- Invariants
- Required tests
- Generated files
- Handwritten extension points

Other modules must not reach into a module's internals directly. Cross-module communication should happen through commands, queries, events, public module APIs, or generated clients.

Use `docs/module-manifest.md` as the source standard for `modules/<module>/module.yaml`.

Use `docs/change-spec.md` as the source standard for `changes/<date>-<change-id>/`.

Use `docs/conversation-log.md` as the source standard for `conversations/`.

When the maintainer introduces product intent, rejects an interpretation, names a concept, or makes an architectural decision, preserve that context in a conversation log. Redact secrets before committing.

Use `docs/agent-decision-record.md` as the source standard for `decisions/`.

When a decision affects long-term architecture, generated file conventions, module ownership, or a rejected plausible alternative, create or update an Agent Decision Record. Keep rationale concise and public; do not store hidden chain-of-thought.

Generated files are immutable to non-system agents. If generated output is wrong, change the source schema, template, or generator unless a change spec or decision record explicitly grants a `generated_file_override`.

For the server runtime, Go is the first implementation language. WebSocket is the first gameplay/client protocol. Protobuf is the first wire message format. PostgreSQL is the first authoritative durable relational store. S3-compatible object storage is a planned object-storage abstraction, with MinIO as the preferred local/self-hosted candidate pending a dependency adoption record. Domain modules must not depend directly on third-party transport, protocol, persistence, object-storage, or framework libraries; platform adapters own those dependencies behind vibit-owned interfaces.

Accepted first Go runtime dependencies are recorded in `ADR-0013` and `.arch/dependencies.yaml`:

- `github.com/coder/websocket` for the platform WebSocket transport adapter.
- `google.golang.org/protobuf` and `google.golang.org/protobuf/cmd/protoc-gen-go` for Go Protobuf runtime and generation.
- Buf CLI for Protobuf linting, breaking checks, formatting, and generation orchestration.
- `github.com/jackc/pgx/v5` for PostgreSQL platform persistence adapters.
- `github.com/pressly/goose/v3` for SQL-first migration tooling.
- Go standard-library `testing` first; no external test framework is adopted yet.

Direct imports or invocations of accepted dependencies are allowed only in their declared owner layers. Domain runtime logic and domain modules must use vibit-owned interfaces, generated contracts, repositories, and adapters.

Goose migrations should be SQL-first. Go migrations require a change spec explaining why SQL is insufficient and must not hide domain business logic.

When adding Go runtime code, follow the ADR-0014 package boundary:

- `runtime/cmd/vibit-server/` owns startup, configuration wiring, and process lifecycle.
- `runtime/internal/app/` owns command/query dispatch, application service composition, and transaction orchestration.
- `runtime/internal/platform/transport/ws/` owns `github.com/coder/websocket`.
- `runtime/internal/platform/protocol/protobuf/` owns Protobuf framing and envelope conversion.
- `runtime/internal/platform/persistence/postgres/` owns `github.com/jackc/pgx/v5`.
- `runtime/internal/platform/migrations/` owns `github.com/pressly/goose/v3` invocation and migration validation.
- `runtime/internal/platform/events/` owns event recording and publication mechanisms.
- `runtime/internal/platform/tx/` owns unit-of-work and transaction boundary interfaces.
- `runtime/internal/modules/<module>/` owns handwritten domain behavior only.

Runtime protocol handoff must follow `docs/runtime-protocol-adapter.md`: WebSocket transport reads and writes frames, the Protobuf adapter converts envelopes and payloads, application dispatch routes commands and queries, domain modules enforce invariants, and generated packages provide shapes only.

State-changing commands should enter through application dispatch and run inside an application-owned unit of work. Domain events produced by a command should be recorded in the same unit of work. Query handlers should not mutate state and do not require a write transaction by default. Event publication outside the transaction remains deferred until an explicit event delivery or outbox decision exists.

Before adding persistence implementation, agents must declare or update the relevant repository interfaces, migration expectations, transaction boundaries, and storage verification path. Do not add PostgreSQL drivers, migration tools, S3 SDKs, or MinIO clients without a change spec or adoption record that follows `ADR-0010` and `ADR-0011`.

Before adding foundational dependencies, agents must check `.arch/dependencies.yaml`. Do not change a dependency slot to `accepted` until an adoption record documents the problem solved, license, maintenance activity, abstraction boundary, allowed owners, forbidden owners, replacement path, and verification path.

Before adding new game server capability families or runtime subsystems, agents must check `.arch/reference.yaml` and `docs/reference-game-server-alignment.md`. Map the proposal to the relevant Nakama/Pitaya capability family, then keep the implementation sequence aligned with vibit's contract-first, manifest-first, generated, and checkable architecture. Do not copy external APIs without an explicit compatibility ADR. Do not add Pitaya-style cluster/RPC/service-discovery work before the modular monolith proof slice is stable.

Use `ADR-0012` for decision authority. After explicit maintainer authorization, agents may professionally evaluate and decide technical sub-decisions inside an already ratified direction. Still ask the maintainer before changing constitutional principles, product direction, runtime language, primary protocol direction, persistence direction, major architecture patterns, module ownership, breaking contracts, validation or permission strength, licensing-risk acceptance, hosting, cost, operations, or vendor-lock-in commitments.

Use `docs/schema-validation.md` as the source standard for `schema/`.

When changing the shape of module manifests, change specs, Agent Decision Records, or tool JSON output, update the paired schema file and run `node tools/vibit check schemas`.

## 8. Contract Rules

Public behavior should be specified before implementation.

Contract-bearing artifacts may include:

- API schemas
- Command schemas
- Query schemas
- Event schemas
- Error catalogs
- Permission catalogs
- Database migration schemas
- Generated clients

Rules:

- Public contracts must be declared before use.
- Compatibility-sensitive contracts must be versioned.
- Breaking changes must be explicit.
- Generated output must be traceable to source schema.
- Do not hand-edit generated contract output unless the generator itself is being changed.

## 9. Testing And Verification

Testing is part of the architecture, not a finishing step.

When implementation code exists, relevant verification may include:

- Unit tests
- Contract tests
- Invariant tests
- Integration tests
- Migration tests
- Replay tests
- Architecture checks
- Generator checks
- Documentation consistency checks

This repository does not yet define final verification commands. Until it does, record verification as one of:

```text
Verified: <commands or checks run>
Not verified: <reason>
Not applicable: <reason>
```

Never claim that a change is verified when verification was not run.

## 10. Ask First

Ask the human maintainer before:

- Changing constitutional principles
- Ratifying or replacing the project name
- Redefining module ownership
- Introducing a new architectural pattern
- Making breaking API, command, query, or event changes
- Changing generated file conventions
- Removing tests
- Weakening validation or permission checks
- Moving data ownership between modules
- Accepting meaningful licensing-risk, hosting, cost, operations, or vendor-lock-in commitments
- Changing server runtime language, primary protocol direction, persistence direction, or core project thesis
- Adding a major external framework dependency

## 11. Tiered Gate Strategy

Not all work items require the same level of gating overhead. Use the following tiers to balance safety with velocity:

### Tier 1 — Security-Critical (Two-Step: Gate + Implementation)

Applies to: cryptography, verifier algorithms, Protobuf wire format, credential schema, token lifecycle, secret configuration.

- Requires a separate gate milestone defining the boundary before implementation.
- Gate must include: threat model, redaction rules, fail-closed behavior, dependency posture.
- Implementation follows as a separate bounded work item.

### Tier 2 — Functional Implementation (Single Step)

Applies to: transport features, application policy, registry behavior, route registration, protocol bridge, session lifecycle, connection lifecycle.

- Single implementation work item with boundary definition embedded in the change spec.
- No separate gate milestone required.
- Must still include: focused tests, `vibit check all`, verification record.

### Tier 3 — Lightweight (Direct Implementation)

Applies to: documentation, translation, simple check rules, small tooling edits, and already-ratified mechanical migration-source updates.

- Direct implementation without change spec unless non-trivial.
- Do not classify new data ownership, new schema semantics, new dependencies, or new runtime behavior as Tier 3.
- Verification through `vibit check all` is sufficient.

### Direction Confirmation

- Direction confirmation milestones are **not required** as separate work items.
- Direction is managed through `ask_first` fields on work items and `recommended_direction` on continuation semantics.
- If a direction choice is significant enough, record it as an ADR.
- When the recommended direction is already explicit and the work is Tier 2, continue directly into a bounded functional slice instead of creating a pure confirmation milestone.

## 12. Universal Prohibitions

The following prohibitions apply to **all** milestones and work items. Do not repeat them in individual milestone `non_goals` or work item `ask_first` lists:

Never:

- Treat AI gameplay features as the foundation of this project
- Bypass module boundaries for convenience
- Hide business logic in transport handlers
- Add unregistered public events
- Add unregistered permissions
- Add untyped cross-module payloads
- Make broad repository edits without a declared boundary
- Hand-edit generated files without documenting why
- Leave an English core document materially changed while its Chinese translation silently falls behind
- Claim verification was run when it was not
- Add direct Nakama/Pitaya API compatibility without an explicit ADR
- Add dependencies without a dependency adoption record

## 13. When Adding New Standards

New standards should explain:

- The problem being solved
- The rule being introduced
- The reason the rule helps agents
- The impact on humans
- The expected artifacts
- The verification path
- The migration path from existing work

Prefer a small standard that can be enforced over a broad statement that cannot be checked.

## 14. When Adding Implementation Code

Do not start by scattering framework code across the repository.

Start from the smallest complete slice that proves the core claim:

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

A good first implementation target should include a small but complete backend domain, such as player accounts, inventory, currency, rewards, quests, or match sessions.

## 15. Bootstrapping Control

Self-bootstrapping is useful only while it improves the path to a working server framework.

Before adding a new standard, inspect command, check command, schema, generator, or workflow rule, confirm that it directly supports at least one of:

- The next runtime vertical slice
- A concrete module boundary
- A public contract or generated shape
- A test or verification path
- Agent context reduction for an expected implementation task

If the benefit is mainly that the tooling becomes more complete, defer it.

When the repository already has enough tooling to attempt a small end-to-end backend capability, prefer runtime readiness work over additional meta-tooling, then implement the runtime slice.

Runtime readiness should answer only the decisions needed to make the first slice coherent:

- Implementation language and package layout
- Minimal server instance model
- First module and capability boundary
- Contract format
- Generated versus handwritten file boundary
- Minimum test and verification strategy
- Persistence and migration assumptions

Do not rush into implementation when these choices are still ambiguous. Also do not extend readiness work after it stops changing how the first slice will be built, verified, or maintained.

Record exceptions in a change spec or Agent Decision Record.

## 16. Handoff Requirements

At the end of a change, leave enough context for the next agent or human to continue.

Record:

- What changed
- Why it changed
- Which files changed
- Which contracts or standards changed
- What was verified
- What was not verified
- Which open questions remain

If the work is incomplete, state the next concrete action.

## 16. Current Product Parity Roadmap

`M-096 Nakama Pitaya Product Parity Roadmap` is completed. `W-0168` added `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`, and `ADR-0078`. The maintainer clarified that vibit should become a Nakama/Pitaya-class game backend product and cover common capability families, not merely use those projects as loose inspiration. Product parity means common capability coverage and comparable usefulness; it does not mean direct API compatibility. `runtime.reference_product_parity_roadmap` is the repository check rule.

`M-097 Protocol Logout Route Gate` is completed. `W-0169` added `docs/protocol-logout-route-gate.md`, `docs/protocol-logout-route-gate.zh-CN.md`, and `ADR-0079`. The gate defines future `runtime.authentication.LogoutAccessToken` route semantics with `access_token_in_logout_request_payload` and an `explicit_service_validated_token_lifecycle_route` posture. It authorized the bounded W-0170 implementation slice but did not add socket close, runtime session revocation, active connection invalidation, reconnect behavior, protocol session carriers, dependencies, or direct Nakama/Pitaya API compatibility. `runtime.protocol_logout_route_gate` is the repository check rule.

`M-098 Protocol Logout Route Implementation` is completed. `W-0170` exposed the existing `LogoutAccessToken` service behavior as the explicit `runtime.authentication.LogoutAccessToken` command route. The slice added `LogoutAccessTokenRequest` and `LogoutAccessTokenResponse` to `proto/vibit/authentication/v1/authentication.proto`, regenerated Go Protobuf output through Buf, added protocol bridge mapping, application bootstrap route registration and handler behavior, PostgreSQL startup registration, transaction bypass, and focused tests. It rejects `AuthenticatedRequest` wrapping for logout so the proof comes from `LogoutAccessTokenRequest.access_token`. It does not close sockets, revoke runtime sessions, invalidate active connection records, add reconnect/epoch behavior, add protocol session carriers, add dependencies, or add direct Nakama/Pitaya API compatibility. `runtime.protocol_logout_route_implementation` is the repository check rule.

`M-099 Next Direction Confirmation After Protocol Logout Route Implementation` is completed. `W-0171` selected `define_transport_close_handoff_gate` because protocol logout is now visible and application close policy can invalidate registry records, but no narrow handoff exists for concrete WebSocket socket close. This selection follows Nakama's lifecycle separation and Pitaya's acceptor/session/handler/connection-management layering. The current work queue is active at `M-100/W-0172`, and the next ready work item is `define_transport_close_handoff_gate`. Do not implement concrete socket close, close codes, close reason text, logout-triggered close, runtime session revocation, reconnect/epoch behavior, protocol session carriers, presence, chat, social modules, matchmaking, match runtime, dependencies, or direct Nakama/Pitaya API compatibility in the direction-confirmation step.

`M-100 Transport Close Handoff Gate` is completed. `W-0172` added `docs/transport-close-handoff-gate.md`, `docs/transport-close-handoff-gate.zh-CN.md`, and `ADR-0080`. The gate defines future narrow application-to-WebSocket concrete close handoff, keeps close decisions application-owned, keeps WebSocket transport credential-neutral and policy-neutral, and selects server-observed `connection_id + epoch` as the first handoff target. It does not implement concrete socket close, close codes, close reason text, logout-triggered close, runtime session revocation, reconnect/epoch behavior, protocol session carriers, operations/admin disconnect, dependencies, direct Nakama/Pitaya API compatibility, or broad product modules. `runtime.transport_close_handoff_gate` is the repository check rule. The current work queue is active at `M-101/W-0173`, and the next ready work item is `confirm_next_direction_after_transport_close_handoff_gate`.

`M-101 Next Direction Confirmation After Transport Close Handoff Gate` is completed. `W-0173` selected `implement_transport_close_handoff_single_process` as the next lifecycle-closure direction. The current work queue is active at `M-102/W-0174`, and the next ready work item is `implement_transport_close_handoff_single_process`. This implementation may target only server-observed `connection_id + epoch` through a WebSocket transport-owned handoff while preserving application-owned close policy and transport credential neutrality. Do not select close codes, close reason text, logout-triggered socket close, runtime session revocation, reconnect/epoch behavior, protocol session carriers, operations/admin disconnect, dependencies, direct Nakama/Pitaya API compatibility, or broad product modules in this slice.

`M-102 Transport Close Handoff Single Process Implementation` is completed. `W-0174` added `runtime/internal/platform/transport/ws/close_handoff.go`, `runtime/internal/platform/transport/ws/close_handoff_test.go`, and `ADR-0081`. WebSocket transport now owns a single-process in-memory accepted socket table and exposes `RequestClose` by server-observed `connection_id + epoch`, returning redacted outcomes for close requested, socket not found, epoch mismatch, already closed, and close failed. The implementation does not parse credentials, does not change Protobuf envelope/logout/session behavior, does not select close codes or reason text, and does not add logout-triggered socket close, runtime session revocation, reconnect behavior, protocol session carriers, operations/admin disconnect, dependencies, direct Nakama/Pitaya API compatibility, or broad product modules. The current work queue is active at `M-103/W-0175`, and the next ready work item is `define_reconnect_connection_epoch_functional_slice`.

`M-103 Reconnect Connection Epoch Functional Slice` is completed. `W-0175` added `ADR-0083` and implemented the first application-owned server-observed connection epoch progression behavior in `runtime/internal/app/connection`. The active connection registry now marks earlier active epochs for the same connection id as `superseded`, records `superseded_at` and `superseded_by_epoch`, rejects stale or repeated epochs after a newer epoch exists with `connection_epoch_stale`, keeps superseded records inspectable, and excludes them from active target lists. This slice does not add reconnect tokens, resume routes, Protobuf changes, protocol session carriers, logout-triggered socket close, runtime session revocation, presence, operations/admin disconnect, dependencies, direct Nakama/Pitaya API compatibility, or broad product modules. The current work queue is active at `M-104/W-0176`, and the next ready work item is `define_protocol_session_carrier_functional_slice`.

`M-104 Protocol Session Carrier Functional Slice` is completed. `W-0176` added `ADR-0084` and reused existing `Envelope.Session` metadata as the first protocol-visible runtime session carrier. Successful `runtime.authentication.AuthenticateWithDeviceCredential` responses now carry the server-created runtime session id and authenticated player id in the response envelope session metadata without changing Protobuf sources or generated output. Response envelopes may derive session metadata from already validated application identity, but metadata-only identity remains metadata-only and `session_id` is not proof. This slice does not add reconnect tokens, resume routes, WebSocket handshake authentication, logout-triggered socket close, runtime session revocation, presence behavior, operations/admin disconnect, dependencies, direct Nakama/Pitaya API compatibility, or broad product modules. The current work queue is active at `M-105/W-0177`, and the next ready work item is `define_presence_lifecycle_functional_slice`.

`M-105 Presence Lifecycle Functional Slice` is completed. `W-0177` added `ADR-0085`, registry-backed player presence snapshots in `runtime/internal/app/connection`, a credential-neutral WebSocket lifecycle observer in `runtime/internal/platform/transport/ws`, and PostgreSQL startup composition adapters under `runtime/cmd/vibit-server`. The first presence behavior derives online/offline state from active bound server-owned connection registry records; successful first-message connection binding can feed validated player identity into the registry. This slice does not add Protobuf presence messages, generated output, protocol presence queries, subscriptions, broadcasts, chat, friends, groups, parties, matchmaking, match runtime, operations/admin behavior, durable/distributed presence, reconnect/resume tokens, logout-triggered close, runtime session revocation, dependencies, direct Nakama/Pitaya API compatibility, or broad product modules. The current work queue is active at `M-106/W-0178`, and the next ready work item is `define_presence_protocol_query_functional_slice`.

`M-107 v0.1 Alpha Goal And Long Term Product Target` is completed as a documentation and roadmap slice. `W-0179` added `docs/v0.1-alpha-goal.md`, `docs/v0.1-alpha-goal.zh-CN.md`, `ADR-0086`, and the conversation/change records that define `v0.1 alpha` as the short-term developer-usable release target and AI-era Nakama/Pitaya-class server framework as the long-term product target. It did not publish a release, add runtime behavior, add protocol messages, change generated output, alter migrations, add dependencies, or choose direct Nakama/Pitaya API compatibility. The current work queue remains active at `M-106/W-0178`, and the next ready work item remains `define_presence_protocol_query_functional_slice`.

`M-108 Next Alpha Direction Selection` is completed. `W-0180` added `ADR-0088` and selected `define_local_onboarding_device_credential_issuance_gate` as the next alpha-enabling direction after the protected presence query. This direction selection did not add runtime behavior, protocol messages, generated output, migrations, dependencies, release artifacts, or direct Nakama/Pitaya API compatibility. Because onboarding/device credential issuance touches credential material, verifier digests, one-time raw secret presentation, player account creation, and repository mutation ordering, the next step is the gate-only `M-109/W-0181` work item, not immediate implementation.

`M-109 Local Onboarding Device Credential Issuance Gate` is completed. `W-0181` added `docs/local-onboarding-device-credential-issuance-gate.md`, `docs/local-onboarding-device-credential-issuance-gate.zh-CN.md`, and `ADR-0089`. The gate defines a future local-only application service under `runtime/internal/app/authentication` that may create an active player account and active device credential record in one unit of work, store only credential digests, and return raw device credential text only once after commit. It does not authorize public signup, protocol routes, generated output, migrations, dependencies, access-token issuance from onboarding, runtime session creation from onboarding, external identity providers, password login, account recovery, multi-device linking, release publishing, direct Nakama/Pitaya API compatibility, or broad product modules. `runtime.local_onboarding_device_credential_issuance_gate` is the repository check rule.

`M-110 Local Onboarding Device Credential Issuance Implementation` is completed. `W-0182` added `OnboardLocalPlayerWithDeviceCredential`, local onboarding request/result vocabulary, explicit device credential entropy and id generator dependencies, startup dependency composition, focused tests, and `ADR-0090`. The service generates server-issued device credential material with an explicit entropy reader, computes credential lookup and verifier digests with existing helpers, creates a player account and digest-only credential record in the same unit of work, returns raw credential text only after commit, does not issue access tokens or runtime sessions from onboarding, and leaves the existing login route non-creating. It did not add public protocol routes, Protobuf sources or generated output, migrations, repository interface changes, dependencies, production signup, external identity, password login, recovery, multi-device linking, release artifacts, direct Nakama/Pitaya API compatibility, or broad product modules. `runtime.local_onboarding_device_credential_issuance_implementation` is the repository check rule.

`M-111 Next Alpha Direction Selection After Local Onboarding` is completed. `W-0183` added `ADR-0091` and selected `define_authenticated_gameplay_e2e_slice` as the next alpha-enabling direction after local onboarding. This direction selection did not implement authenticated gameplay E2E, add protocol routes, add generated output, change migrations, add dependencies, publish a release, select direct Nakama/Pitaya API compatibility, add production signup, external identity providers, password login, account recovery, multi-device linking, or broad product modules.

`M-112 Authenticated Gameplay E2E Functional Slice` is completed. `W-0184` added `ADR-0092` and `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go` as the focused proof that local onboarding, protocol login, first-message connection binding, protected inventory grant/read, protected presence query, logout, and post-logout protected-route rejection compose as one local alpha path. It did not add production runtime behavior, protocol routes, Protobuf sources, generated output, migrations, repository interface changes, dependencies, release artifacts, production signup, broad product modules, or direct Nakama/Pitaya API compatibility. `runtime.authenticated_gameplay_e2e_functional_slice` is the repository check rule.

`M-113 Runtime Runbook Alpha Path Refresh` is completed. `W-0185` refreshed `docs/runtime-runbook.md` and `docs/runtime-runbook.zh-CN.md` around the now-proven local alpha path, including memory vs PostgreSQL runtime posture, verifier key handling, the focused authenticated gameplay E2E proof, application-service-only local onboarding, and redaction expectations. It did not add runtime behavior, change startup configuration semantics, add protocol routes, change Protobuf sources or generated output, add migrations, add dependencies, publish a release, add production signup, broad product modules, or direct Nakama/Pitaya API compatibility. `runtime.runtime_runbook_alpha_path_refresh` is the repository check rule.

`M-114 Minimal Example Client Or Request Loop` is completed. `W-0186` added `examples/local-alpha-request-loop.sh` and `ADR-0094` as the minimal local alpha request-loop script over the focused authenticated gameplay E2E proof. The script prints a redacted path summary and Go test status, and it does not add runtime behavior, change startup configuration semantics, add public protocol onboarding, change Protobuf sources or generated output, add migrations, add dependencies, publish a release, choose production signup, add broad product modules, or direct Nakama/Pitaya API compatibility. `runtime.minimal_example_client_or_request_loop` is the repository check rule.

`M-115 Health Readiness Version Config Surface` is completed. `W-0187` added minimal JSON `/healthz`, `/readyz`, `/version`, and `/configz` endpoints in `runtime/cmd/vibit-server`, with focused tests and `ADR-0095`. The surface is for local alpha troubleshooting only; it reports redacted runtime posture and must not expose verifier keys, raw credentials, raw tokens, DSNs, digests, headers, cookies, query strings, subprotocol values, remote addresses, or concrete transport metadata. It did not add broad operations/admin behavior, observability dependencies, authentication/session behavior changes, startup configuration semantic changes, Protobuf sources, generated output, migrations, release artifacts, broad product modules, or direct Nakama/Pitaya API compatibility. `runtime.health_readiness_version_config_surface` is the repository check rule.

`M-116 Alpha Acceptance Checklist` is completed. `W-0188` added `docs/alpha-acceptance-checklist.md`, `docs/alpha-acceptance-checklist.zh-CN.md`, and `ADR-0096` as the local v0.1 alpha acceptance checklist. It covers repository intake, prerequisites, migration posture, local configuration, local onboarding posture, login, connection binding, protected inventory, presence query, logout, status endpoints, checks, redaction, contribution entry points, and release deferrals. It did not publish `v0.1 alpha`, add release packaging, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules. `runtime.alpha_acceptance_checklist` is the repository check rule.

`ADR-0082` adopts the tiered gate density strategy. Security-critical work remains two-step gate plus implementation. Functional work such as transport features, application policy, registry behavior, route registration, protocol bridges, session lifecycle, and connection lifecycle should usually be a single bounded work item with the boundary embedded in its change spec. Lightweight docs/tooling/translation/check-rule work may be direct when trivial. Direction confirmation milestones are no longer mandatory for future Tier 2 work; use `ask_first`, `recommended_direction`, and ADRs for significant direction changes. `M-103/W-0175` is the first prospective application of this strategy and should proceed as the bounded reconnect/connection epoch functional slice rather than as a pure confirmation step.

`M-117 Package Alpha Developer Flow` is completed. `W-0189` added `docs/alpha-developer-flow.md`, `docs/alpha-developer-flow.zh-CN.md`, and `ADR-0097` as the packaged local alpha developer journey. It connects README intake, v0.1 alpha goal, alpha acceptance checklist, runtime runbook, redacted request-loop script, status endpoints, PostgreSQL manual setup posture, redaction rules, verification commands, and the next contribution path. It did not publish `v0.1 alpha`, create release tags or binaries, add hosted deployments, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules. `runtime.package_alpha_developer_flow` is the repository check rule.

`M-118 Release Publishing Decision Gate` is completed. `W-0190` added `docs/release-publishing-decision-gate.md`, `docs/release-publishing-decision-gate.zh-CN.md`, and `ADR-0098` as the release publishing decision gate. It defines release-publishing prerequisites, release artifact boundaries, verification requirements, stop conditions, redaction expectations, and the next release execution preparation direction. It did not publish `v0.1 alpha`, create release tags, binaries, archives, containers, packages, checksums, hosted deployments, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules. `runtime.release_publishing_decision_gate` is the repository check rule.

`M-119 Release Execution Preparation Gate` is completed. `W-0191` added `docs/release-execution-preparation-gate.md`, `docs/release-execution-preparation-gate.zh-CN.md`, and `ADR-0099` as the release execution preparation gate. It defines release execution planning inputs, release-note input boundaries, artifact plan boundaries, maintainer approval points, verification requirements, rollback notes, stop conditions, redaction expectations, and the next release execution authorization direction. It did not publish `v0.1 alpha`, create release tags, binaries, archives, containers, packages, checksums, provenance files, hosted deployments, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules. `runtime.release_execution_preparation_gate` is the repository check rule.

`M-120 Release Execution Authorization Gate` is completed. `W-0192` added `docs/release-execution-authorization-gate.md`, `docs/release-execution-authorization-gate.zh-CN.md`, and `ADR-0100` as the release execution authorization gate. It defines final go/no-go criteria, required verification state, release identifier review, artifact authorization boundaries, maintainer approval requirements, authorization outcome, stop conditions, redaction expectations, and a blocked next release execution maintainer decision. It did not publish `v0.1 alpha`, choose or create release identifiers or tags, create release binaries, archives, containers, packages, checksums, provenance files, hosted deployments, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules. `runtime.release_execution_authorization_gate` is the repository check rule.

`M-121 Release Execution Maintainer Decision` is completed. `W-0193` recorded the maintainer decision as `go_to_release_identifier_artifact_plan`, added `docs/release-execution-maintainer-decision.md`, `docs/release-execution-maintainer-decision.zh-CN.md`, and `ADR-0101`, and added `runtime.release_execution_maintainer_decision` check coverage. The decision allows the release execution path to continue to planning only. It did not approve a final release identifier, create or authorize release tags, create or authorize release binaries, archives, containers, packages, checksums, provenance files, hosted deployments, GitHub release records, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules.

`M-122 Release Identifier And Artifact Plan` is completed. `W-0194` added `docs/release-identifier-artifact-plan.md`, `docs/release-identifier-artifact-plan.zh-CN.md`, and `ADR-0102`, and added `runtime.release_identifier_artifact_plan` check coverage. The plan proposes `v0.1.0-alpha.1`, records that no local tag, remote origin tag, or GitHub release record conflict was observed on 2026-05-21, and defines a source-first future surface of Git tag, GitHub release record, release notes, and hosting-platform source archive. It did not publish `v0.1 alpha`, select the identifier for execution, create or push release tags, create release binaries, archives, containers, packages, checksums, provenance files, hosted deployments, GitHub release records, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules.

`M-123 Release Execution Final Authorization` is completed. `W-0195` recorded final maintainer `go` authorization for `v0.1.0-alpha.1`, added `docs/release-execution-final-authorization.md`, `docs/release-execution-final-authorization.zh-CN.md`, `ADR-0103`, release notes, and `runtime.release_execution_final_authorization` check coverage. The authorization allows creating and pushing Git tag `v0.1.0-alpha.1`, creating GitHub Release `v0.1.0-alpha.1`, and publishing only the GitHub source archive. It does not authorize release binaries, packages, containers, checksums, signing/provenance artifacts, hosted deployments, install scripts, registry publication, public announcements beyond the GitHub release record, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules.

`M-124 First Alpha User Discovery` is completed. `W-0196` added `docs/first-alpha-user-discovery-loop.md`, `docs/first-alpha-user-discovery-loop.zh-CN.md`, `ADR-0104`, and `runtime.first_alpha_user_discovery_loop` check coverage. The loop records target developer segments, outreach surfaces, feedback capture fields, review questions, success signals, and stop conditions. It does not authorize public announcements beyond the GitHub release record, paid promotion, hosted deployments, additional release artifacts, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules.

`M-125 First Alpha Feedback Intake Surfaces` is completed. `W-0197` added `.github/ISSUE_TEMPLATE/alpha-feedback.yml`, `docs/first-alpha-feedback-intake-surfaces.md`, `docs/first-alpha-feedback-intake-surfaces.zh-CN.md`, `docs/product-maturity-milestones.md`, `docs/product-maturity-milestones.zh-CN.md`, `ADR-0105`, and `runtime.first_alpha_feedback_intake_surfaces` check coverage. It records source-first alpha as reached, prototype-ready foundation as the next product stage, single-node production-candidate foundation as planned, and Nakama/Pitaya-class product as the long-term target. It does not authorize broad announcements, paid promotion, hosted deployments, additional release artifacts, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct compatibility, or broad product modules.

`M-126 Prototype Ready Foundation Execution Plan` is completed. `W-0198` added `docs/prototype-ready-foundation-execution-plan.md`, `docs/prototype-ready-foundation-execution-plan.zh-CN.md`, `ADR-0106`, and `runtime.prototype_ready_foundation_execution_plan` check coverage. It records the Stage 2 execution sequence, candidate work families, maturity-stage mapping, Nakama/Pitaya capability mapping, success criteria, and stop conditions, and selects `prototype_ready_local_development_path_gate` as the first execution slice. It does not authorize runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, hosted deployments, release artifacts, public announcements, paid promotion, broad operations/admin behavior, authentication/session behavior changes, broad product modules, or direct compatibility.

`M-127 Prototype Ready Local Development Path Gate` is completed. `W-0199` added `docs/prototype-ready-local-development-path-gate.md`, `docs/prototype-ready-local-development-path-gate.zh-CN.md`, `ADR-0107`, and `runtime.prototype_ready_local_development_path_gate` check coverage. It records supported prerequisites, startup expectations, migration expectations, configuration and secret posture, example-flow shape, allowed future write areas, verification expectations, and stop conditions. It does not authorize runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, hosted deployments, release artifacts, public announcements, paid promotion, broad operations/admin behavior, authentication/session behavior changes, broad product modules, or direct compatibility.

`M-128 Prototype Ready Local Development Path Package` is active. `W-0200 Implement prototype-ready local development path package` is next-ready. It should package setup, migration, configuration/secret redaction, example-flow, and verification ergonomics inside the W-0199 gate before broader product changes.

Future major work must map to a roadmap family: identity/auth/session, connection lifecycle, storage, presence/status/notifications, chat/realtime messaging, friends/groups/parties, leaderboards/tournaments, economy/progression, matchmaking, match runtime, server runtime hooks/RPC, operations, SDK/developer experience, or distributed runtime. The near-term priority is now the prototype-ready local development path package after the source-first release, feedback intake, execution plan, and local path gate. Do not jump directly to chat, social modules, matchmaking, match runtime, SDKs, distributed runtime, public announcements, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility before that package is recorded.
