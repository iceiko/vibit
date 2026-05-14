# Application Authentication Service Interface Boundary

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Future application-owned runtime authentication service interface boundary
Depends on: `docs/runtime-authentication-implementation-boundary.md`, `docs/authentication-generated-contract-shape-timing.md`, `docs/authentication-contract-error-permission-surfaces.md`
Canonical decision: `ADR-0039`

The paired Simplified Chinese translation is `docs/application-authentication-service-interface-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This document defines the service-interface boundary that future runtime authentication behavior must fit into.

It exists after metadata-only generated authentication contract shapes and before handwritten runtime authentication behavior. Its job is to make the future application service shape predictable for agents without adding login execution, token generation, verifier comparison, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, or migration schema changes.

This is a service-interface boundary only. It does not add Go service code.

## 2. Core Rule

Application authentication service interfaces are application-owned.

The future owner is:

```text
runtime/internal/app
```

The first implementation must expose authentication behavior through application-level request and result vocabulary before domain dispatch receives validated identity.

Target ownership:

```text
protocol-decoded request proof
-> application authentication service interface
-> application unit of work
-> authentication.Repository
-> persistence-only PostgreSQL adapter
-> application-owned RequestIdentity handoff
-> domain dispatch
```

No transport, Protobuf adapter, domain module, generated contract shape, repository, or PostgreSQL adapter may become the service-interface owner.

## 3. Planned Package Boundary

The future package may live directly under `runtime/internal/app` or under a child package owned by `runtime/internal/app`, such as:

```text
runtime/internal/app/authentication
```

A later implementation work item must choose the concrete Go package path before code is added.

Allowed interface-level dependencies after a later service-interface implementation gate:

- Go standard library types.
- Application-owned `RequestIdentity`, `RouteKey`, `Session`, `ApplicationResult`, and application error vocabulary.
- The application unit-of-work boundary represented by `runtime/internal/platform/tx`.
- The module-owned `authentication.Repository`, obtained through the application unit of work.

Forbidden interface-level dependencies:

- WebSocket transport packages.
- Generated Protobuf packages.
- PostgreSQL driver packages.
- Migration tooling.
- JWT, OAuth, OIDC, password-hashing, provider SDK, Redis-like token/session store, S3, or MinIO dependencies.
- Generated authentication contract shape packages as runtime registries or behavior owners.

Generated authentication contract shapes may inform names and mapping tables, but future service interfaces must not require domain modules or transport handlers to import generated authentication shape packages.

## 4. Service Vocabulary

The generated authentication contract shapes under:

```text
runtime/internal/generated/contracts/runtime/authentication/
```

inform the future service request and result vocabulary.

The first planned service surface is:

```yaml
service_boundary:
  owner: runtime/internal/app
  status: boundary_defined_no_code
  commands:
    AuthenticateWithDeviceCredential:
      request_vocabulary:
        - credential_proof
        - requested_player_id
        - client_instance_id
        - account_creation_intent
      result_vocabulary:
        - authentication_status
        - actor_kind
        - player_id
        - access_token
        - token_type
        - issued_at
        - expires_at
        - token_record_id
    ValidateAccessToken:
      request_vocabulary:
        - access_token
        - route_kind
        - route_module
        - route_name
        - connection_id
        - connection_epoch
      result_vocabulary:
        - validation_status
        - proof_status
        - actor_kind
        - actor_id
        - player_id
        - player_id_validated
        - session_validated
        - token_record_id
    LogoutAccessToken:
      request_vocabulary:
        - access_token
        - logout_reason
      result_vocabulary:
        - revoked
        - logout_scope
        - token_record_id
        - revoked_at
    RefreshAccessToken:
      first_posture: unsupported_reserved
      required_error: AUTHENTICATION_REFRESH_NOT_SUPPORTED
```

This vocabulary is not a public wire schema. Protobuf authentication messages remain deferred.

## 5. Redaction Expectations

Future service interfaces must classify secret fields before implementation.

Secret input fields:

- `credential_proof`
- `access_token`

Secret internal material:

- Raw credential material.
- Raw token material.
- Credential verifier digest inputs.
- Token verifier digest inputs.
- Verifier keys or key identifiers when a future standard marks them confidential.

Rules:

- Raw credential material and raw access token material must be redacted from logs, errors, conversation logs, change specs, test names, table rows, and audit records.
- Public errors must use registered authentication error codes, not raw proof details.
- Future tests must assert that redacted values do not appear in application errors or audit-safe facts.
- Token issuance may present the raw access token to the client once through a later ratified response carrier. That exception must be tested and documented by the token generation gate.

## 6. Unit-Of-Work And Repository Boundary

Future service behavior may use `authentication.Repository` only through the application unit-of-work boundary.

Allowed shape:

```text
application service method
-> runner.WithinUnitOfWork(...)
-> UnitOfWork.NewAuthenticationRepository(...)
-> authentication.Repository
```

Rules:

- The application service owns orchestration.
- The repository owns storage-neutral record operations only.
- The PostgreSQL adapter owns SQL persistence only.
- The application service must not bypass the unit-of-work boundary for state-changing authentication operations.
- The repository must not generate tokens, compare verifiers, validate proof, parse bearer tokens, publish events, or return public protocol responses.
- The PostgreSQL adapter must not make authentication decisions.

## 7. Request Identity Handoff

Future `ValidateAccessToken` service behavior must convert proven proof into `RequestIdentity` before production-sensitive domain dispatch.

Target handoff:

```yaml
request_identity_handoff:
  owner: runtime/internal/app
  input: ValidateAccessToken result
  output: RequestIdentity
  required_success_markers:
    status: validated
    actor_kind: player
    actor_id: player_id
    player_id_validated: true
    session_validated: false
  metadata_only_allowed_as_proof: false
```

`MetadataOnlySessionValidator` remains a bootstrap path and does not satisfy this boundary.

Domain modules must consume `RequestIdentity`; they must not parse, validate, or compare authentication proof.

## 8. Error, Permission, And Audit Handoff

Future service interfaces must map to existing semantic contracts.

Errors:

- `AuthenticateWithDeviceCredential` uses `AUTHENTICATION_PROOF_MISSING`, `AUTHENTICATION_PROOF_MALFORMED`, `AUTHENTICATION_CREDENTIAL_INVALID`, `AUTHENTICATION_ACCOUNT_DISABLED`, `AUTHENTICATION_RATE_LIMITED`, `AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE`, `AUTHENTICATION_TOKEN_STORE_UNAVAILABLE`, and `AUTHENTICATION_NOT_IMPLEMENTED`.
- `ValidateAccessToken` uses missing, malformed, invalid, expired, revoked, unavailable, disabled-account, and not-implemented token errors.
- `LogoutAccessToken` uses missing, malformed, invalid, token-store-unavailable, and not-implemented token errors.
- `RefreshAccessToken` maps to `AUTHENTICATION_REFRESH_NOT_SUPPORTED` in the first posture.

Permissions:

- `authentication_device_credential_login`
- `authentication_access_token_validate`
- `authentication_access_token_logout`
- `authentication_access_token_refresh`

Audit handoff:

- `AuthenticationSucceeded`
- `AuthenticationFailed`
- `TokenIssued`
- `TokenValidationFailed`
- `TokenRevoked`
- `LogoutRequested`

Audit publication and audit persistence remain separate gates. The service interface may define audit-safe facts later, but it must not store audit events until a separate storage path is ratified.

## 9. Nakama And Pitaya Alignment

Nakama remains a capability reference for account authentication, session token issuance, token expiration, logout, revocation, and realtime authenticated actor binding.

Pitaya remains a vocabulary reference for handler context, frontend acceptor separation, backend handler separation, and route identity context.

vibit adapts those lessons through an application-owned service boundary. It does not copy Nakama or Pitaya public APIs, and it does not make transport handlers or domain handlers responsible for token validation.

## 10. Verification Path

The repository check rule for this boundary is:

```text
runtime.application_authentication_service_interface_boundary
```

For changes that touch this boundary, run:

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

If a change spec exists, also run:

```bash
node tools/vibit check change <change-id> --json
```

Runtime Go tests are not required for this boundary-only standard unless a later work item adds or changes Go runtime code.

## 11. Non-Goals

This standard does not authorize:

- Go application authentication service code.
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
- Authentication dependencies.
- Authentication audit persistence.
- Runtime session persistence.
- Changes to `authentication.Repository`.
- Changes to ratified migration schemas.
