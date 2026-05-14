# Runtime Authentication Implementation Boundary

Status: Draft v0.1
Last updated: 2026-05-14
Scope: First runtime authentication implementation boundary planning after the authentication PostgreSQL adapter
Depends on: `docs/authentication-token-session-validation.md`, `docs/login-method-token-format-ratification.md`, `docs/selected-login-token-boundary-checks.md`, `docs/postgresql-persistence-boundary.md`
Canonical decision: `ADR-0036`

The paired Simplified Chinese translation is `docs/runtime-authentication-implementation-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This document defines the boundary for future runtime authentication implementation after vibit has:

- Ratified `device_credential_login` as the first login method.
- Ratified opaque high-entropy access tokens.
- Ratified explicit request proof payloads as the first request proof posture.
- Added credential and token verifier PostgreSQL migration sources.
- Added the storage-neutral `authentication.Repository` interface.
- Implemented the authentication PostgreSQL persistence adapter.

This document is a planning and boundary standard. It does not implement login, token generation, token validation, verifier comparison, logout execution, refresh, cleanup jobs, Protobuf messages, WebSocket proof carriers, generated authentication shapes, authentication dependencies, or production authentication behavior.

## 2. Core Rule

Runtime authentication is application-owned.

The first implementation must pass through the application boundary before any domain module receives authenticated identity:

```text
protocol decoded request
-> application authentication boundary
-> module-owned authentication repository interface
-> platform-owned PostgreSQL adapter through an application unit of work
-> application-owned request identity handoff
-> domain dispatch
```

No transport, protocol adapter, domain module, player repository, generated file, or PostgreSQL adapter may become the owner of authentication proof, token validation, verifier comparison, or permission decisions.

The current runtime still has only metadata-only identity. Metadata-only `player_id`, `session_id`, and `connection_id` remain unauthenticated.

## 3. Ownership Split

### Application Service Boundary

Owner:

```text
runtime/internal/app
```

Future responsibility after a bounded implementation work item authorizes code:

- Orchestrate device credential login after protocol decoding.
- Orchestrate access-token validation before production-sensitive domain dispatch.
- Orchestrate presented-token logout.
- Convert authentication outcomes into `RequestIdentity`.
- Map authentication failures to registered application errors.
- Request repository operations through the module-owned `authentication.Repository`.
- Run state-changing authentication operations inside the application unit of work.

Must not:

- Import WebSocket transport packages.
- Import generated Protobuf packages.
- Store raw token or credential material.
- Own PostgreSQL SQL text or driver handles.
- Hide verifier comparison inside generic dispatch.
- Treat metadata-only identity as proof.

### Authentication Module Boundary

Owner:

```text
runtime/internal/modules/authentication
```

Current responsibility:

- Own storage-neutral authentication repository interfaces and record shapes.
- Own closed status vocabulary for credential and token verifier records.
- Preserve copying, normalization, and UTC timestamp invariants for repository inputs and outputs.

Future responsibility:

- Continue owning storage-neutral repository interfaces when runtime authentication behavior starts.

Must not:

- Generate tokens.
- Compare credential or token verifiers.
- Parse bearer tokens.
- Execute login, logout, or cleanup jobs.
- Import PostgreSQL, WebSocket, Protobuf, JWT, OAuth, OIDC, password-hashing, provider SDK, or Redis-like dependencies.

### PostgreSQL Adapter Boundary

Owner:

```text
runtime/internal/platform/persistence/postgres
```

Current responsibility:

- Implement `authentication.Repository` using caller-supplied executors.
- Persist and read ratified credential and token verifier records.
- Preserve no-transaction-control behavior.
- Keep pgx details inside the platform package.

Must not:

- Generate raw credential or token material.
- Compare verifier digests.
- Interpret token proof.
- Execute authentication decisions.
- Emit domain or audit events directly.
- Know WebSocket, Protobuf, request identity, or permission behavior.

### Protocol Adapter Boundary

Owner:

```text
runtime/internal/platform/protocol/protobuf
```

Future responsibility after a protocol gate authorizes it:

- Convert ratified Protobuf authentication request and response messages into application route requests.
- Convert application authentication errors into public-safe Protobuf error envelopes.

Must not:

- Choose proof carrier semantics.
- Parse or validate tokens outside application-owned boundaries.
- Treat current `Session` metadata as proof.
- Generate tokens or compare verifiers.

### WebSocket Transport Boundary

Owner:

```text
runtime/internal/platform/transport/ws
```

Current responsibility:

- Accept connections.
- Read and write opaque binary frames.
- Delegate frame bytes to injected handlers.

Must not:

- Read `Authorization`, `Bearer`, `Cookie`, or `Sec-WebSocket-Protocol` as authentication proof until a later WebSocket proof-carrier decision authorizes it.
- Parse credentials or tokens.
- Bind player identity by connection metadata.
- Own domain permission decisions.

## 4. Token Material Lifecycle Placeholder

Opaque access-token generation remains a separate gate.

Future token generation must define:

- Raw token generation owner.
- Entropy source and minimum entropy verification.
- Text encoding.
- Redaction tests.
- One-time client presentation rules.
- Verifier digest derivation.
- Storage behavior through `authentication.Repository`.
- Rotation and revocation behavior.

Until that gate exists:

- No code may call `crypto/rand` or equivalent for authentication tokens.
- No raw token may appear in logs, tests, change specs, conversation logs, or database rows.
- No token string may be copied into Protobuf `Session` metadata.

## 5. Verifier Comparison Boundary

Verifier comparison remains a separate gate from repository persistence.

Future verifier comparison must define:

- Which verifier algorithm is used for device credential proof.
- Which verifier algorithm is used for access-token proof.
- Whether the Go standard library is sufficient or an external dependency adoption record is required.
- Constant-time comparison requirements.
- Error mapping for invalid, expired, revoked, unavailable, and malformed proof.
- Tests that prove raw secrets are not stored or echoed.

Repositories may store and retrieve verifier digests. They must not decide whether presented proof is valid.

## 6. Request Identity Handoff

Future token validation must produce application-owned request identity before domain dispatch.

The target handoff is:

```yaml
owner: runtime/internal/app
input: explicit_request_proof_payload_after_protocol_decode
output: RequestIdentity
required_success_markers:
  actor_kind: player
  validation_status: authentication_proven_or_session_validated
  player_id_validated: true
  session_validated: false_until_session_persistence_is_ratified
metadata_only_allowed_as_proof: false
```

Rules:

- Domain modules consume `RequestIdentity`; they do not validate tokens.
- Inventory permission policy may use validated identity only after the validation boundary exists.
- The current `MetadataOnlySessionValidator` remains a bootstrap path and does not grant production player-owned permissions.
- Runtime session persistence remains deferred.

## 7. Error, Permission, And Audit Mapping

Runtime authentication implementation must use the registered semantic surfaces under:

```text
contracts/runtime/authentication/
```

Required first mapping:

- `AuthenticateWithDeviceCredential` failures map to `authentication_errors`.
- `ValidateAccessToken` failures map to missing, malformed, invalid, expired, revoked, unsupported, unavailable, or not-implemented failure classes.
- `LogoutAccessToken` must distinguish presented-token revocation success, missing proof, invalid proof, expired proof, revoked proof, and store unavailability.
- `RefreshAccessToken` remains unsupported in the first implementation and must map to `AUTHENTICATION_REFRESH_NOT_SUPPORTED` if exposed.
- Raw credential and raw token material must never appear in public errors.

Permission surfaces remain semantic until implementation:

- `authentication_device_credential_authenticate`
- `authentication_access_token_validate`
- `authentication_access_token_logout`
- `authentication_access_token_refresh`

Audit event publication and audit persistence remain separate gates. Runtime authentication code may prepare event facts only after the event publication and storage path is ratified.

## 8. Implementation Queue

The first runtime authentication implementation queue must remain split into separately reviewable gates.

Recommended order:

1. Add or refine runtime authentication boundary checks.
2. Decide whether generated Go authentication contract shapes are needed before runtime code.
3. Ratify token and credential verifier algorithms, redaction tests, and dependency posture.
4. Add application-owned authentication service interfaces and tests without protocol or transport behavior.
5. Implement token material generation only inside the ratified application boundary.
6. Implement credential verifier comparison only inside the ratified application boundary.
7. Implement `AuthenticateWithDeviceCredential` execution.
8. Implement `ValidateAccessToken` execution and request identity handoff.
9. Implement `LogoutAccessToken` execution for the presented access token.
10. Define and implement cleanup job behavior.
11. Add Protobuf authentication messages and bridge code after protocol impact is ratified.
12. Add WebSocket proof-carrier behavior only if a later decision chooses a WebSocket carrier.

This order may be revised by a later ADR, but no gate may be collapsed into another gate silently.

## 9. Nakama And Pitaya Alignment

Nakama remains a capability reference for:

- Account authentication.
- Session token issuance.
- Token expiration.
- Token revocation and logout.
- Realtime socket binding to authenticated actors.

Pitaya remains a vocabulary reference for:

- Session context passed to handlers.
- Frontend acceptor versus backend handler separation.
- Route handler identity context.
- Session binding concepts.

vibit adapts these concepts into its own boundaries. It does not copy Nakama or Pitaya public API shapes, does not put authentication into WebSocket acceptors by default, and does not let route handlers validate tokens directly.

## 10. Verification Path

The repository check rule for this boundary is:

```text
runtime.authentication_implementation_boundary
```

For changes that touch this boundary, run:

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check architecture --json
node tools/vibit check all --json
```

If a change spec exists, also run:

```bash
node tools/vibit check change <change-id> --json
```

Implementation changes must also run the relevant Go tests. Live PostgreSQL is not required by this boundary standard unless the change explicitly adds or changes database behavior that cannot be verified statically.

## 11. Non-Goals

This standard does not authorize:

- Runtime login behavior.
- Access-token generation.
- Access-token validation.
- Credential verifier comparison.
- Token verifier comparison.
- Logout execution.
- Refresh-token behavior.
- Cleanup jobs.
- Protobuf authentication messages.
- WebSocket proof carriers.
- WebSocket handshake authentication.
- Generated authentication shapes.
- Authentication dependencies.
- Authentication audit persistence.
- Runtime session persistence.
- Changes to `authentication.Repository`.
- Changes to ratified migration schemas.
