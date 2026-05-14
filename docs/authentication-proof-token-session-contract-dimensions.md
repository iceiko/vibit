# Authentication Proof And Token Session Contract Dimensions

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Semantic contract dimensions for authentication proof, token/session validation, request identity handoff, validation statuses, failure classes, retryability, errors, permissions, commands, queries, and events
Depends on: `docs/authentication-token-session-validation.md`

The paired Simplified Chinese translation is `docs/authentication-proof-token-session-contract-dimensions.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard ratifies the semantic contract dimensions that future authentication proof and token/session validation work must use.

It does not choose a login method, token format, token carrier, refresh behavior, signing behavior, credential store, external identity provider, session store, Protobuf envelope change, WebSocket handshake behavior, runtime player handler, or WebSocket route.

The purpose is narrower and more foundational: make the words and fields that future agents use for authentication and session contracts stable before implementation begins.

Future work may ratify concrete contracts such as login, token refresh, logout, session invalidation, connection binding, or service authentication. Those contracts must map back to the dimensions in this document unless a later ADR supersedes them.

## 2. Required Reading

Read this standard together with:

- `docs/authentication-token-session-validation.md`
- `docs/player-account-session-contracts.md`
- `docs/player-identity-session-boundary.md`
- `docs/game-protocol.md`
- `docs/runtime-protocol-adapter.md`
- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `contracts/runtime/session/commands/ValidateSession.yaml`
- `contracts/runtime/session/events/SessionValidated.yaml`
- `contracts/runtime/session/errors/session_errors.yaml`
- `contracts/runtime/session/permissions/session_permissions.yaml`
- `ADR-0019`
- `ADR-0021`
- `ADR-0023`

Nakama remains the reference for account, authentication, session token, refresh token, expiration, revocation, and realtime socket capability coverage.

Pitaya remains the reference for session binding, session context handoff, route handler vocabulary, and realtime connection/session separation.

Both references inform vocabulary. Neither governs vibit's API shape.

## 3. Ratified Dimensions

Every future authentication proof or token/session validation contract must say which of these dimensions it owns, consumes, or leaves untouched.

| Dimension | Meaning | Current status |
| --- | --- | --- |
| `actor_kind` | The kind of validated or candidate actor. | Ratified vocabulary. |
| `actor_id` | The stable identifier for the actor after validation. | Ratified handoff field. |
| `player_id` | Stable player identity owned by the player module. | Metadata-only until validated. |
| `session_id` | Logical session identifier or candidate session metadata. | Metadata-only until validated. |
| `connection_id` | Transport-local connection metadata. | Not proof. |
| `connection_epoch` | Reconnection or connection generation metadata. | Not proof. |
| `validation_status` | The semantic result category produced before domain dispatch. | Ratified vocabulary. |
| `proof_status` | Whether authentication proof was absent, present but unverified, proven, rejected, expired, revoked, unsupported, or unavailable. | Ratified vocabulary for future contracts. |
| `failure_class` | Machine-readable failure family for validation errors. | Ratified vocabulary. |
| `retryability` | Whether the caller or runtime may retry the same semantic operation. | Ratified expectation. |
| `request_identity_handoff` | The normalized identity context passed to application/domain handlers. | Owned by `runtime/internal/app`. |
| `permission_basis` | The trust state module permission policies may rely on. | Metadata-only is not sufficient for production permissions. |

These dimensions are semantic. They do not imply any token shape, storage schema, cryptographic primitive, transport header, envelope field, or database table.

## 4. Actor Kinds

Ratified actor kinds:

| Actor kind | Meaning | Production authority |
| --- | --- | --- |
| `unknown` | The actor kind is absent or cannot be trusted. | No player-owned or privileged authority. |
| `anonymous` | The request is intentionally unauthenticated. | Only explicitly anonymous capabilities. |
| `player` | A player actor, backed by a stable `player_id` after validation. | Player-owned permissions only after validation succeeds. |
| `service` | A trusted internal or external service actor. | Only service permissions explicitly modeled for that actor. |
| `admin` | An administrative actor. | Deferred; requires a future admin permission model. |

The current Go runtime has `unknown`, `player`, `service`, and `admin` in `runtime/internal/app`. The `anonymous` actor kind is ratified as contract vocabulary for future semantic contracts. Adding it to runtime code is a separate implementation change and is not required by this design-only step.

Rules:

- `player` must not be inferred from a raw client-supplied `player_id` for production permissions.
- `service` and `admin` require separately ratified proof and permission catalogs before use.
- `anonymous` must be explicit; it is not the same as missing metadata.

## 5. Validation Statuses

Ratified validation statuses:

| Status | Meaning | Current production authority |
| --- | --- | --- |
| `unknown` | No validation result is available. | None. |
| `anonymous` | The request was accepted as intentionally unauthenticated. | Anonymous-only permissions. |
| `metadata_only` | Identity metadata was normalized without authentication proof. | Not sufficient for production permissions. |
| `authentication_proven` | Authentication proof was verified, but logical session binding may still be separate. | Deferred until implementation. |
| `session_validated` | A logical session was validated and bound to request identity. | Future permission basis after ratification. |
| `service_validated` | Service authority was validated. | Future service permission basis after ratification. |
| `rejected` | Validation failed and the request must not dispatch to production-sensitive handlers. | None. |

The current `ValidateSession` semantic contract uses `metadata_only` and `validated`. This standard refines future vocabulary without requiring runtime code changes. During future implementation, `validated` should be replaced or mapped to a more specific status such as `authentication_proven`, `session_validated`, or `service_validated` before production permission policies rely on it.

Rules:

- `metadata_only` must remain non-authenticated.
- `validated` must not be used as a vague long-term production state when the validator can distinguish authentication proof, session binding, or service authority.
- Domain modules should consume normalized request identity, not token/session internals.

## 6. Proof Statuses

Ratified proof statuses for future contracts:

| Proof status | Meaning |
| --- | --- |
| `not_present` | No proof material was provided. |
| `present_unverified` | Proof material exists but has not been verified. |
| `proven` | Proof was verified by a ratified authenticator or validator. |
| `rejected` | Proof was checked and rejected. |
| `expired` | Proof is outside its accepted time window. |
| `revoked` | Proof was explicitly invalidated. |
| `unsupported` | The proof kind, issuer, method, or carrier is not supported by the current runtime. |
| `unavailable` | Validation could not run because a required validator or dependency was unavailable. |

This vocabulary intentionally does not choose:

- JWT or opaque tokens.
- Access token or refresh token structure.
- Signing, key management, issuer, audience, expiration duration, revocation store, or replay handling.
- Header, envelope, first-message, or payload carrier behavior.

## 7. Failure Classes

Ratified failure classes:

| Failure class | Meaning | Default retryability |
| --- | --- | --- |
| `missing_proof` | Required proof material was absent. | Not retryable without new proof. |
| `malformed_proof` | Proof material could not be parsed by a ratified validator. | Not retryable without changed input. |
| `unsupported_proof` | The proof kind, issuer, method, or carrier is not supported. | Not retryable until configuration or implementation changes. |
| `invalid_proof` | Proof was understood but rejected. | Not retryable with the same proof. |
| `expired_proof` | Proof was validly shaped but expired. | Retryable only through a ratified refresh or re-login flow. |
| `revoked_proof` | Proof was explicitly invalidated. | Not retryable with the same proof. |
| `session_not_found` | A referenced logical session does not exist. | Usually not retryable with the same session. |
| `session_expired` | A logical session expired. | Retryable only through a ratified refresh or re-login flow. |
| `session_revoked` | A logical session was invalidated. | Not retryable with the same session. |
| `actor_disabled` | The actor or player account is disabled. | Not retryable until state changes. |
| `permission_denied` | Validation succeeded but the actor lacks permission. | Not retryable without authority changes. |
| `validator_unavailable` | The validator or its dependency is unavailable. | Retryable. |
| `not_implemented` | The requested authentication or validation path is not implemented. | Not retryable until implementation changes. |

Future error catalogs must map each validation error to one of these classes or explicitly introduce a new class through a change spec and ADR when the class is long-lived.

## 8. Retryability Rules

Retryability is part of the contract and must not be guessed by clients, agents, or transport adapters.

Rules:

- Error catalogs must declare `retryable: true` or `retryable: false`.
- Retryability describes the same semantic request with the same proof material.
- Refresh, re-login, proof replacement, account recovery, or support intervention are separate flows, not retries of the same validation request.
- `validator_unavailable` may be retryable.
- `invalid_proof`, `revoked_proof`, `session_revoked`, `actor_disabled`, and `not_implemented` are not retryable with the same proof.
- `expired_proof` and `session_expired` are retryable only through a future ratified refresh or re-login flow.

## 9. Command Dimensions

Authentication proof and token/session validation commands must declare:

- Actor input semantics.
- Proof input semantics, when proof is in scope.
- Route and target context, when request-level validation is in scope.
- Transport metadata consumed as metadata only.
- Whether the command may read player accounts, credentials, token state, session state, or external identity links.
- Output validation status.
- Output request identity handoff fields.
- Failure classes.
- Retryability.
- Required permissions.
- Invariants that prevent transport, protocol, player repository, or domain handlers from owning validation behavior.

Current ratified command:

```text
ValidateSession
```

`ValidateSession` is application-owned and semantic-only. It describes the request identity handoff before domain dispatch. It does not implement token parsing, credential lookup, player account lookup, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes.

## 10. Query Dimensions

No runtime authentication/session validation query is ratified by this step.

Future queries may be ratified for read-only inspection, session status, public key metadata, or service validation metadata only after their ownership, exposure, permission model, and information-leakage behavior are specified.

Query rules:

- Queries must not leak credential material, token contents, secret keys, provider secrets, or raw proof material.
- Queries must distinguish operator/admin inspection from gameplay client behavior.
- Queries must declare whether they are safe for anonymous, player, service, or admin actors.
- Queries must not become a shortcut for validating tokens inside domain modules.

## 11. Event Dimensions

Authentication proof and token/session events must declare:

- Whether the event is a domain fact, security fact, audit fact, or runtime-observation fact.
- Whether it may be published to clients.
- Which identifiers are safe to expose.
- Whether raw proof material, token strings, credentials, provider subjects, or secrets are forbidden.
- Which command or validation path produced it.
- Compatibility and versioning rules.

Current ratified event:

```text
SessionValidated
```

`SessionValidated` is a semantic fact only. It may describe that a decoded request's session metadata was evaluated into request identity. It does not add an event bus, durable audit store, token/session store, or public client event stream.

Rules:

- A future `AuthenticationSucceeded`, `AuthenticationFailed`, `TokenRefreshed`, `SessionInvalidated`, or similar event requires separate ratification.
- Security-sensitive failure events should be designed for audit/operations before public exposure.
- Event payloads must not contain raw credentials, token strings, password hashes, provider secrets, or full third-party identity payloads.

## 12. Error Dimensions

Authentication and session errors must declare:

- Stable error code.
- Failure class.
- Category.
- Retryability.
- Public-safe message.
- Whether the error may be returned to clients.
- Whether additional detail is internal-only.
- Commands or queries that use it.

Current runtime session errors:

- `SESSION_INVALID`
- `SESSION_VALIDATOR_UNAVAILABLE`
- `SESSION_VALIDATION_NOT_IMPLEMENTED`

Required mapping:

| Error code | Failure class | Retryable |
| --- | --- | --- |
| `SESSION_INVALID` | `invalid_proof` or `session_not_found`, depending on the future validator path | `false` |
| `SESSION_VALIDATOR_UNAVAILABLE` | `validator_unavailable` | `true` |
| `SESSION_VALIDATION_NOT_IMPLEMENTED` | `not_implemented` | `false` |

The current metadata-only runtime should not emit `SESSION_INVALID` for missing production proof unless a future validator has been ratified to require proof.

## 13. Permission Dimensions

Authentication and session permissions must distinguish:

- Permission to run validation infrastructure.
- Permission granted by a validated actor to a module action.
- Permission to inspect authentication/session state.
- Permission to administer or revoke sessions.

Current runtime session permission:

```text
runtime_session_validate
```

`runtime_session_validate` allows application dispatch to evaluate decoded request metadata before domain handlers run. It is not a gameplay permission and does not grant player-owned domain authority.

Rules:

- Domain module permissions must not treat metadata-only identity as production proof.
- Service and admin permissions require future explicit catalogs.
- Token possession alone must not become a permission shortcut without a ratified validator and request identity handoff.

## 14. Request Identity Handoff

The request identity handoff is the boundary between validation and domain behavior.

Owner:

```text
runtime/internal/app
```

Required semantic fields:

- `validation_status`
- `actor_kind`
- `actor_id`
- `player_id`
- `session_id`
- `connection_id`
- `connection_epoch`
- `session_validated`
- `player_id_validated`
- `reason`

Rules:

- Domain modules receive request identity after application validation.
- Domain modules must not parse tokens, credentials, WebSocket headers, Protobuf envelope internals, or session stores directly.
- `player_id_validated: true` requires proof that the request actor is allowed to act as that `player_id`.
- `session_validated: true` requires proof that a logical session was validated and bound according to a future ratified session model.
- `metadata_only` identity may be useful for local proof slices and development behavior, but it is not a production permission basis.

## 15. Reference Pattern Map

### Nakama Patterns

| Pattern | Vibit position | Contract dimension impact |
| --- | --- | --- |
| Session token | Deferred | Token/session validation contracts must have proof status, failure class, retryability, and request identity handoff dimensions before implementation. |
| Refresh token | Deferred | Refresh is a separate flow; expired proof is not a retry of the same validation request. |
| Token expiration | Adopted as dimension | Expiration must map to `expired_proof` or `session_expired` before implementation. |
| Token revocation/logout | Adopted as dimension | Revocation must map to `revoked_proof`, `session_revoked`, and explicit events before implementation. |
| Realtime socket bound to authenticated session | Adapted | vibit keeps request identity handoff application-owned; handshake binding remains a later decision. |

### Pitaya Patterns

| Pattern | Vibit position | Contract dimension impact |
| --- | --- | --- |
| Session object separate from transport | Adopted | `connection_id` is metadata; request identity is application-owned. |
| Handler receives session context | Adapted | Domain handlers receive `RequestIdentity`, not a transport-owned session object. |
| Session binding | Adopted as vocabulary | Binding requires explicit `session_validated` and `player_id_validated` semantics. |
| Frontend/backend split | Deferred | Distributed session routing is not part of this semantic contract step. |

## 16. Non-Goals

This standard does not:

- Choose JWT, opaque tokens, refresh tokens, signing, expiration duration, revocation store, rotation, key management, or token storage.
- Choose guest, device, email/password, custom ID, social login, OAuth, OIDC, or external identity provider login.
- Add credential storage, password hashing, provider secrets, external identity tables, token tables, or session tables.
- Change the Protobuf envelope.
- Change WebSocket handshake behavior.
- Add runtime authentication code.
- Add token parsing.
- Add runtime player handlers or WebSocket routes.
- Add event bus publication or audit persistence.

## 17. Required Future Artifacts

When future work implements any concrete authentication or token/session behavior, it must add or update the relevant subset of:

- Change spec under `changes/`.
- ADR for long-lived architecture-sensitive choices.
- Contract sources under `contracts/`.
- `.arch/contracts.yaml`.
- `.arch/runtime.yaml`.
- Module manifests and guides.
- Error and permission catalogs.
- Persistence boundary and migrations, when state is stored.
- Dependency adoption records, when security or provider dependencies are introduced.
- Protocol sources and generated output, when wire shape changes.
- Runtime tests and architecture checks.
- English documentation and Simplified Chinese translation.

## 18. Verification

Default verification for this standard:

```bash
node tools/vibit check contracts --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check work --json
node tools/vibit check change ratify-authentication-proof-token-session-contract-dimensions --json
node tools/vibit check all --json
```

Go tests are not required for this design-only standard unless runtime Go code changes.

## 19. Agent Rules

Agents must:

- Use the ratified dimensions when designing authentication proof or token/session validation contracts.
- Keep metadata-only identity non-authenticated.
- Map future validation failures to explicit failure classes and retryability.
- Preserve request identity handoff through `runtime/internal/app`.
- Record when Nakama or Pitaya vocabulary is adopted, adapted, deferred, or rejected.

Agents must not:

- Treat token strings, player account rows, envelope metadata, or WebSocket connection IDs as authentication proof by themselves.
- Hide token validation in transport, protocol adapters, repositories, or domain handlers.
- Add concrete login methods, token behavior, credential storage, session persistence, Protobuf envelope changes, or WebSocket handshake behavior from this standard alone.
