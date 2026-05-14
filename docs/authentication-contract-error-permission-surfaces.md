# Authentication Contract, Error, And Permission Surfaces

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Semantic contract, error, permission, and audit surfaces for the selected first login and token posture
Depends on: `docs/token-lifecycle-storage-implications.md`
Canonical decision: `ADR-0028`

The paired Simplified Chinese translation is `docs/authentication-contract-error-permission-surfaces.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This document defines the semantic surfaces required before vibit implements the selected first authentication posture:

```yaml
login_method: device_credential_login
access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
refresh_token: not_in_first_implementation
renewal_method: reauthenticate_with_device_credential_login
```

The goal is to make future authentication work readable and verifiable for agents before runtime code exists.

This document does not implement authentication, token generation, token validation, logout, refresh, credential lookup, token storage, audit persistence, runtime handlers, WebSocket routes, Protobuf messages, generated output, migrations, or schema changes.

## 2. Problem

Authentication is a common place for framework code to become hard for agents to change safely.

Without explicit semantic surfaces, a future agent may:

- Put credential parsing into WebSocket transport.
- Put token validation into Protobuf adapters or domain modules.
- Use metadata-only `player_id` or `session_id` as proof.
- Add ad hoc error codes.
- Add route handlers before permission and failure models exist.
- Treat refresh token behavior as implied by access-token issuance.
- Store raw token or credential material in logs, tests, events, or database records.

The solution is to define contract, error, permission, and audit surfaces now, while keeping implementation blocked behind schema, repository, test, and runtime milestone gates.

## 3. Rule

Future authentication implementation must start from the registered semantic contract sources under:

```text
contracts/runtime/authentication/
```

The runtime authentication family is registered in:

```text
.arch/contracts.yaml
```

Agents must use `node tools/vibit inspect contracts --module runtime --json` and `node tools/vibit check contracts --json` before adding authentication behavior.

Runtime authentication remains application-owned:

```text
runtime/internal/app
```

No other layer may own credential parsing, access-token validation, or first logout behavior unless a later ADR changes that ownership.

## 4. Selected Surfaces

The first selected semantic surfaces are:

```yaml
commands:
  - AuthenticateWithDeviceCredential
  - ValidateAccessToken
  - LogoutAccessToken
  - RefreshAccessToken
events:
  - AuthenticationSucceeded
  - AuthenticationFailed
  - TokenIssued
  - TokenValidationFailed
  - TokenRevoked
  - LogoutRequested
errors:
  - authentication_errors
permissions:
  - authentication_permissions
queries: []
```

These are semantic source contracts only. They do not create generated Go contract shapes, Protobuf messages, runtime handlers, WebSocket routes, database tables, or migrations.

## 5. Command Surfaces

### `AuthenticateWithDeviceCredential`

Purpose:

```text
Authenticate a player using the selected high-entropy device credential login method and issue an opaque access token.
```

Contract source:

```text
contracts/runtime/authentication/commands/AuthenticateWithDeviceCredential.yaml
```

This command may only be implemented after future gates define credential schema, token verifier schema, repository boundaries, redaction tests, and a runtime authentication implementation milestone.

It must not treat a raw operating-system device ID, player ID, session ID, connection ID, or other metadata as credential proof.

### `ValidateAccessToken`

Purpose:

```text
Validate explicit opaque access-token proof into application-owned request identity before domain dispatch.
```

Contract source:

```text
contracts/runtime/authentication/commands/ValidateAccessToken.yaml
```

The selected request proof carrier is explicit request payload proof. Current Protobuf `Session` metadata remains metadata-only and must not become proof by reinterpretation.

### `LogoutAccessToken`

Purpose:

```text
Revoke the presented opaque access token.
```

Contract source:

```text
contracts/runtime/authentication/commands/LogoutAccessToken.yaml
```

The first logout scope is:

```yaml
logout_scope: presented_access_token
```

This does not authorize logout-all-sessions, account-wide token revocation, credential-wide token revocation, WebSocket close behavior, or admin revocation.

### `RefreshAccessToken`

Purpose:

```text
Reserve the token refresh semantic surface while keeping refresh tokens out of the first implementation.
```

Contract source:

```text
contracts/runtime/authentication/commands/RefreshAccessToken.yaml
```

Refresh tokens are not part of the first posture. The first renewal method remains:

```yaml
renewal_method: reauthenticate_with_device_credential_login
```

The refresh contract exists so agents can see that refresh is intentionally considered and intentionally unsupported in the first implementation, not forgotten.

## 6. Error Surface

The error catalog is:

```text
contracts/runtime/authentication/errors/authentication_errors.yaml
```

Required first error families include:

```yaml
missing_proof:
  - AUTHENTICATION_PROOF_MISSING
  - AUTHENTICATION_TOKEN_MISSING
malformed_proof:
  - AUTHENTICATION_PROOF_MALFORMED
  - AUTHENTICATION_TOKEN_MALFORMED
invalid_proof:
  - AUTHENTICATION_CREDENTIAL_INVALID
  - AUTHENTICATION_TOKEN_INVALID
expired_proof:
  - AUTHENTICATION_TOKEN_EXPIRED
revoked_proof:
  - AUTHENTICATION_TOKEN_REVOKED
unsupported_proof:
  - AUTHENTICATION_REFRESH_NOT_SUPPORTED
actor_disabled:
  - AUTHENTICATION_ACCOUNT_DISABLED
validator_unavailable:
  - AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
  - AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
not_implemented:
  - AUTHENTICATION_NOT_IMPLEMENTED
```

Errors must be public-safe and must not include raw credential material, raw token values, token prefixes, verifier hashes, password hashes, provider secrets, or hidden validation details.

## 7. Permission Surface

The permission catalog is:

```text
contracts/runtime/authentication/permissions/authentication_permissions.yaml
```

First permissions:

```yaml
authentication_device_credential_login:
  dimension: unauthenticated_login_entrypoint
authentication_access_token_validate:
  dimension: validation_infrastructure_permission
authentication_access_token_logout:
  dimension: player_token_lifecycle_permission
authentication_access_token_refresh:
  dimension: deferred_token_lifecycle_permission
authentication_admin_revoke_token:
  dimension: deferred_admin_permission
```

These permissions do not grant domain module authority by themselves. Domain modules consume normalized request identity after validation and their own permission policies.

Metadata-only identity remains insufficient for production permissions.

## 8. Audit Event Surface

The first authentication audit-oriented event surfaces are:

```yaml
AuthenticationSucceeded:
  reason: successful selected login proof
AuthenticationFailed:
  reason: rejected or unavailable selected login proof
TokenIssued:
  reason: opaque access-token verifier record issued
TokenValidationFailed:
  reason: access-token proof rejected or unavailable
TokenRevoked:
  reason: token verifier record revoked
LogoutRequested:
  reason: presented-token logout requested
```

These events are semantic only. They do not add an event bus, client event stream, audit table, durable audit persistence, or runtime publication behavior.

Raw credentials, raw tokens, token verifier hashes, password hashes, provider secrets, and full provider payloads are forbidden in these events.

## 9. Account Linking Surface

Account linking is not in scope for the first selected posture.

Status:

```yaml
account_linking: deferred
external_identity_storage_required_for_first_posture: false
```

Future account linking requires provider subject semantics, link/unlink permissions, conflict behavior, recovery behavior, merge behavior, schema gates, audit events, and tests.

## 10. Protocol And Transport Impact

This standard does not change:

- Protobuf envelope fields.
- Current Protobuf `Session` metadata semantics.
- WebSocket handshake authentication.
- WebSocket transport behavior.
- First system-message authentication.
- Runtime player handlers.
- WebSocket routes.

The selected request proof carrier remains semantic explicit request payload proof until a later protocol decision ratifies a wire shape.

## 11. Storage Impact

This standard does not add storage.

Future implementation remains blocked on W-0071 schema gates for:

- Credential records.
- Token verifier records.
- Optional session records, if a later posture needs them.
- External identity records, if a later login/linking posture needs them.

Player account lifecycle tables remain credential-free, token-free, external-identity-free, and session-free.

## 12. Agent Impact

This standard helps agents by making the next implementation boundary inspectable:

```bash
node tools/vibit inspect contracts --module runtime --json
node tools/vibit inspect contract --module runtime --type command --id AuthenticateWithDeviceCredential --json
node tools/vibit inspect contract --module runtime --type command --id ValidateAccessToken --json
node tools/vibit check contracts --json
```

Agents should read the relevant command, event, error, and permission manifests before editing implementation code.

## 13. Human Impact

Humans get a stable vocabulary for discussing the first authentication slice without needing to inspect future Go code.

The tradeoff is that the repository gains more manifest files before runtime behavior exists. That is intentional: the project values explicit, inspectable boundaries over implicit security behavior.

## 14. Verification Path

Default verification for this standard:

```bash
node tools/vibit inspect next --json
node tools/vibit check contracts --json
node tools/vibit inspect contracts --module runtime --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check agent-tooling --json
node tools/vibit check memory --json
node tools/vibit check work --json
node tools/vibit check change define-authentication-contract-error-permission-surfaces --json
node tools/vibit check all --json
git diff --check
```

Go tests are not required for this design and contract-source step because no Go runtime behavior changes.

## 15. Migration Path

Future implementation should proceed in this order:

1. Define credential, token, and session schema gates.
2. Add repository checks for selected login/token boundaries.
3. Close M-013.
4. Select an implementation milestone.
5. Add schema, migrations, and repository boundaries before runtime behavior.
6. Add runtime authentication interfaces behind `runtime/internal/app`.
7. Add protocol/wire behavior only after separate Protobuf or WebSocket decisions, if needed.

## 16. Non-Authorization

This standard does not authorize:

- Runtime authentication code.
- Login handlers.
- Token generation, parsing, validation, refresh, revocation, rotation, replay handling, or storage.
- Credential tables.
- External identity tables.
- Token tables.
- Session tables.
- Migrations.
- Generated contract shapes.
- Protobuf messages or generated Protobuf output.
- Password hashing, JWT, OAuth, OIDC, provider SDK, Redis-like, cryptography, key-management, or major authentication dependencies.
- Protobuf envelope changes.
- WebSocket handshake authentication.
- First system-message authentication.
- Runtime player handlers.
- WebSocket routes.
- Treating metadata-only `player_id`, `session_id`, `connection_id`, or `connection_epoch` as proof.

## 17. Follow-Up

Next work:

```text
W-0071 Define credential token session schema gates
```

W-0071 must define future schema gates without adding migrations, repository implementations, runtime lookup code, handlers, or routes.
