# Agent Operating Guide

Status: Draft v0.1  
Last updated: 2026-05-13
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

This repository is currently in the constitutional and standards-design phase.

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

Framework implementation code, generators, modules, and verification commands may not exist yet. When they do not exist, document that verification is not available instead of pretending that it ran.

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

The runtime authentication implementation boundary standard is `docs/runtime-authentication-implementation-boundary.md`, with `docs/runtime-authentication-implementation-boundary.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0036` records the boundary decision. Future runtime authentication is application-owned under `runtime/internal/app`; it must use `authentication.Repository` through the application unit-of-work boundary, keep the PostgreSQL adapter persistence-only, and convert validated proof into `RequestIdentity` before domain dispatch. Token generation, verifier comparison, login execution, access-token validation, logout execution, cleanup jobs, Protobuf authentication messages, WebSocket proof carriers, generated authentication shapes, and authentication dependencies remain separate gates. `runtime.authentication_implementation_boundary` is the repository check rule for this boundary.

The authentication generated contract shape timing standard is `docs/authentication-generated-contract-shape-timing.md`, with `docs/authentication-generated-contract-shape-timing.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0038` records the timing decision. Generated Go authentication contract shapes now exist after the runtime authentication implementation boundary and before service interfaces, using `contracts/runtime/authentication/` as source and `runtime/internal/generated/contracts/runtime/authentication/` as the output root. Generated files remain immutable and metadata-only.

The application authentication service interface boundary standard is `docs/application-authentication-service-interface-boundary.md`, with `docs/application-authentication-service-interface-boundary.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0039` records the boundary decision. Future authentication service interfaces are application-owned under `runtime/internal/app`; generated authentication shapes inform service-level request/result vocabulary; service behavior may use `authentication.Repository` only through the application unit-of-work boundary; validated proof must become `RequestIdentity` before domain dispatch. This boundary does not authorize Go service code or runtime authentication behavior. `runtime.application_authentication_service_interface_boundary` is the repository check rule for this boundary.

`ADR-0037` closes the runtime authentication implementation boundary planning milestone and opens the generated authentication contract shape timing gate. `ADR-0038` completes the timing decision. `W-0089` completes generator/check support plus metadata-only generated authentication shape output. `ADR-0039` and `W-0090` complete the service-interface boundary step. `ADR-0040` and `W-0091` complete the verifier algorithm/redaction step. `ADR-0041` and `W-0092` complete the secret configuration/verifier key loading preparation step. `ADR-0042` and `W-0093` complete the material generation preparation step. `ADR-0043` and `W-0094` complete the verifier digest computation and comparison preparation step. `ADR-0044` and `W-0095` complete the implementation readiness step. `ADR-0045` and `W-0096` complete the local verifier key configuration loading gate. `W-0097` completes the explicit in-memory verifier key set validator implementation slice. `ADR-0046` and `W-0098` complete the environment verifier key loader gate. `W-0099` completes the environment verifier key loader implementation slice. `ADR-0047` and `W-0100` complete the token and credential material generation implementation gate. `W-0101` completes the token and credential material generation helper implementation slice. `ADR-0048` and `W-0102` complete the verifier digest helper implementation gate. `W-0103` completes the verifier digest computation helper implementation slice. `ADR-0049` and `W-0104` complete the verifier digest comparison helper gate.

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

- Update the Chinese translation in the same change when the English source changes materially.
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

## 11. Never

Never:

- Treat AI gameplay features as the foundation of this project
- Bypass module boundaries for convenience
- Hide business logic in transport handlers
- Add unregistered public events
- Add unregistered permissions
- Add untyped cross-module payloads
- Make broad repository edits without a declared boundary
- Hand-edit generated files without documenting why
- Leave an English public document materially changed while its Chinese translation silently falls behind
- Claim verification was run when it was not

## 12. When Adding New Standards

New standards should explain:

- The problem being solved
- The rule being introduced
- The reason the rule helps agents
- The impact on humans
- The expected artifacts
- The verification path
- The migration path from existing work

Prefer a small standard that can be enforced over a broad statement that cannot be checked.

## 13. When Adding Implementation Code

Do not start by scattering framework code across the repository.

Start from the smallest complete slice that proves the core claim:

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

A good first implementation target should include a small but complete backend domain, such as player accounts, inventory, currency, rewards, quests, or match sessions.

## 14. Bootstrapping Control

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

## 15. Handoff Requirements

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
