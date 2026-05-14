# Authentication, Token, And Session Validation Design Standard

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Authentication proof, token behavior, runtime session validation, credential boundaries, external identity boundaries, session persistence boundaries, request identity trust, Protobuf envelope interaction, and WebSocket handshake interaction
Canonical decision: `ADR-0023`

The paired Simplified Chinese translation is `docs/authentication-token-session-validation.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has player account lifecycle persistence and an application-owned session validation hook, but it still does not have production authentication.

This standard defines the design boundary before future agents implement login, token validation, credential storage, external identity linking, runtime session persistence, WebSocket handshake authentication, player account handlers, or player WebSocket routes.

The goal is to keep authentication and session behavior agent-readable and safe to extend. A future agent should be able to identify where proof is produced, where it is validated, where it is stored, where it becomes request identity, and where domain permissions consume it.

This standard does not choose:

- A concrete login method.
- JWT, opaque tokens, refresh tokens, signing, expiration, revocation, or token storage behavior.
- Credential storage, password hashing, OAuth, OIDC, social login, device login, guest login, or custom ID behavior.
- External identity linking tables or provider dependencies.
- Runtime session persistence.
- Protobuf envelope changes.
- WebSocket handshake authentication behavior.
- Runtime player account handlers or WebSocket routes.

## 2. Required Reading

Read this standard together with:

- `docs/player-identity-session-boundary.md`
- `docs/player-account-session-contracts.md`
- `docs/game-protocol.md`
- `docs/runtime-protocol-adapter.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `modules/player/module.yaml`
- `runtime/AGENTS.md`
- `ADR-0019`
- `ADR-0021`
- `ADR-0022`
- `ADR-0023`

Reference reading:

- Nakama documentation: `https://heroiclabs.com/docs/nakama/`
- Nakama authentication and sessions concepts: `https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama realtime multiplayer and socket concepts: `https://heroiclabs.com/docs/nakama/concepts/multiplayer/`
- Pitaya documentation: `https://pitaya.readthedocs.io/`

Nakama and Pitaya are references. They do not govern vibit's public API shape, runtime module layout, generated file conventions, or agent workflow.

## 3. Core Vocabulary

The following terms must remain distinct.

### Authentication Proof

Authentication proof is the result of verifying that an actor is allowed to bind to a player identity, service identity, or future administrative identity.

Rules:

- Authentication proof is not the same thing as a token string.
- Authentication proof is not the same thing as a player account row.
- Authentication proof must become a machine-readable validation result before domain handlers receive the request.
- Authentication proof must not be inferred from client-supplied `player_id`, `session_id`, or `connection_id` metadata.

### Login Method

A login method is the way a client obtains authentication proof or a session credential.

Examples include guest login, device login, email/password login, custom ID login, social login, OAuth, OIDC, and external identity-provider login.

Rules:

- No login method is ratified by this standard.
- Every login method must have a future contract, storage boundary, error model, and verification path before implementation.
- A login method must not be implemented inside WebSocket transport, Protobuf adapters, inventory handlers, or generic player account persistence.

### Token

A token is a credential-like artifact presented by a client or service for validation.

Examples include access tokens, session tokens, refresh tokens, opaque tokens, and signed tokens.

Rules:

- No token format is ratified by this standard.
- Token parsing, signing, verification, refresh, revocation, rotation, expiration, storage, and replay handling are all separate design dimensions.
- Token behavior must be declared before implementation.
- Token validation must not be hidden in domain modules, Protobuf payload bridges, or WebSocket frame handlers.

### Credential

A credential is secret or proof material used by a login method or identity provider.

Examples include passwords, password hashes, device secrets, provider secrets, OAuth credentials, OIDC subjects, and provider-issued identity material.

Rules:

- Credential storage is not ratified by this standard.
- Credentials must not be stored in `player_accounts` or `player_account_events`.
- Credential storage must have a separate schema boundary, secret-handling rule, dependency adoption record when needed, and verification path.

### External Identity Link

An external identity link maps a future vibit player account to an identity provider subject.

Rules:

- External identity linking is not ratified by this standard.
- Provider names, provider subject IDs, provider metadata, linking conflicts, unlinking, account recovery, and merge behavior require a future design.
- External identity links must not be added to player account lifecycle tables by convenience.

### Runtime Session

A runtime session is the server-side validation and binding state used after authentication/session validation is ratified.

Rules:

- Runtime sessions are not player accounts.
- Runtime sessions are not WebSocket connection IDs.
- Runtime sessions may eventually bind `session_id`, `player_id`, actor kind, claims, expiration, connection metadata, and revocation state, but no storage model is ratified here.
- Session validation is application-owned and runs before production-sensitive domain dispatch.

### Request Identity

Request identity is the application-facing identity context passed to command and query handlers.

Current owner:

```text
runtime/internal/app
```

Rules:

- `RequestIdentity` is the handoff type for domain permissions.
- Domain modules consume request identity; they do not validate tokens or credentials.
- Current metadata-only request identity is not authenticated proof.

### WebSocket Handshake Authentication

WebSocket handshake authentication means validating an actor before or during connection establishment.

Rules:

- WebSocket handshake authentication is not ratified by this standard.
- Until a future handshake standard exists, WebSocket transport must remain credential-neutral.
- Any future handshake design must keep transport mechanics separate from authentication proof and session binding semantics.

## 4. Trust States

vibit uses these trust states as design vocabulary.

| State | Meaning | Production authority |
| --- | --- | --- |
| `anonymous` | No actor proof is present. | Must not grant player-owned or privileged permissions. |
| `metadata_only` | Identity text came from envelope or transport metadata and was normalized by the current metadata-only path. | Must not grant production permissions by itself. |
| `authentication_proven` | A future ratified authenticator has verified a login method, token, or trusted service proof. | May be used only after the validating component and contracts are ratified. |
| `session_validated` | A future ratified session validator has bound request identity to a valid logical session. | May be used by module permission policies after validation succeeds. |
| `service_validated` | A future ratified service-auth path has verified non-player service authority. | May be used only by explicitly modeled service permissions. |

The current runtime has only metadata-only behavior. `MetadataOnlySessionValidator` preserves the existing proof-slice request flow; it does not authenticate clients.

## 5. Ownership Model

### Authentication Boundary

Owner status:

```text
planned, not implemented
```

Responsibilities when ratified:

- Verify login-specific proof or token/session credentials.
- Produce machine-readable authentication results.
- Map authentication failures to registered errors.
- Avoid leaking credential implementation details into domain modules.

Must not:

- Store account lifecycle rows as a side effect unless a contract requires it.
- Hide token parsing inside transport, protocol, player account repository, or inventory code.
- Treat metadata-only identity as proof.

### Application Session Validation

Owner:

```text
runtime/internal/app
```

Responsibilities:

- Invoke session validation after protocol decoding and before domain dispatch.
- Convert validation results into `RequestIdentity`.
- Preserve the current metadata-only behavior only as a non-authenticated bootstrap path.
- Keep authentication provider details behind explicit interfaces when implemented.

Must not:

- Import WebSocket transport libraries.
- Import generated Protobuf packages.
- Store credential records.
- Own player account lifecycle persistence.

### Player Module

Owner:

```text
modules/player
runtime/internal/modules/player
```

Responsibilities:

- Own stable `player_id` and player account lifecycle.
- Own player account contracts and lifecycle persistence.
- Remain the durable account owner used by future authentication and linking flows.

Must not:

- Store credentials, token state, refresh tokens, runtime sessions, WebSocket connection state, or request validation results in the account lifecycle tables.
- Validate tokens inside account repositories.
- Add runtime player handlers or WebSocket routes without separate ratification.

### WebSocket Transport

Owner:

```text
runtime/internal/platform/transport/ws
```

Responsibilities:

- Accept WebSocket connections.
- Read and write binary frames.
- Provide transport-local metadata such as connection IDs.
- Delegate opaque frame bytes to protocol/application composition.

Must not:

- Parse credentials or tokens.
- Authenticate players.
- Bind player accounts to connections.
- Own session persistence.
- Enforce domain permissions.

### Protobuf Protocol Adapter

Owner:

```text
runtime/internal/platform/protocol/protobuf
```

Responsibilities:

- Decode and encode the current envelope.
- Preserve existing session metadata fields.
- Convert envelope session metadata into application handoff types.
- Map validation or protocol errors into error envelopes after the relevant errors are registered.

Must not:

- Add credential or token fields to the envelope without a protocol change spec and ADR.
- Treat envelope `player_id` or `session_id` as proof.
- Own long-lived session state.
- Perform permission decisions.

### Domain Modules

Examples:

```text
inventory, currency, reward, quest, match
```

Responsibilities:

- Enforce module-owned permissions and invariants using `RequestIdentity`.
- Treat `metadata_only` as unauthenticated for production-sensitive operations.

Must not:

- Validate tokens.
- Parse credentials.
- Query credential or session stores directly.
- Create or link player accounts unless they are the player module and the behavior is ratified.

## 6. Request Flow

Current non-authenticated bootstrap flow:

```text
websocket frame
-> protobuf envelope
-> route request
-> metadata-only session validator
-> metadata-only request identity
-> application dispatch
-> domain handler
```

Future production-sensitive flow:

```text
websocket frame
-> protobuf envelope
-> route request
-> ratified authentication/session validation boundary
-> validated request identity
-> application dispatch
-> module permission policy
-> domain handler
```

Rules:

- Production-sensitive domain handlers must receive validated request identity before trusting player-owned or privileged actions.
- Request-level validation is the required application handoff before domain dispatch.
- Handshake-level validation, first-message validation, every-request validation, and hybrid validation are future design choices.
- A future handshake model must not remove the application-owned validation handoff.

## 7. Token Design Dimensions

Future token work must decide and record the following before implementation:

- Token format: opaque, signed, structured, or another form.
- Token issuer and verifier ownership.
- Subject semantics: player, service, admin, guest, or another actor kind.
- Audience and route scope.
- Session binding and connection binding.
- Expiration and clock-skew behavior.
- Refresh behavior.
- Revocation behavior.
- Rotation behavior.
- Replay detection behavior.
- Storage requirements.
- Secret and key management.
- Error codes and retryability.
- Logging and redaction.
- Test fixtures and negative tests.

This standard intentionally leaves those dimensions undecided.

## 8. Credential And External Identity Boundaries

Credential and external identity storage require a later standard before schema or code exists.

Future credential work must define:

- Which login methods are supported.
- Which secret material is stored.
- Which hashing, encryption, or provider dependencies are adopted.
- Which tables own credential records.
- How credential rows relate to player accounts.
- How disabled or deleted accounts affect login.
- How failures are logged without leaking secrets.
- How tests prove common failure modes.

Future external identity work must define:

- Provider namespace and subject semantics.
- Link and unlink behavior.
- Duplicate and conflict behavior.
- Account recovery and merge behavior.
- Provider metadata retention rules.
- Audit/event requirements.

Until those standards exist:

- `player_accounts` remains account lifecycle storage only.
- `player_account_events` remains lifecycle event storage only.
- No credential, provider subject, access token, refresh token, or session row may be added by convenience.

## 9. Session Persistence Boundaries

Session persistence is deferred.

Future session persistence work must decide:

- Whether sessions are persisted server-side.
- Whether a session store is PostgreSQL, memory, Redis-like, or another ratified store.
- Session ID generation and lookup semantics.
- Expiration and refresh semantics.
- Revocation and logout semantics.
- Reconnect and connection epoch semantics.
- Relationship to WebSocket connection lifecycle.
- Relationship to token validation.
- Cleanup and migration behavior.
- Opt-in live verification requirements.

Until ratified, the runtime must keep `session_id` and `player_id` metadata-only unless a future validator explicitly validates them.

## 10. Protobuf Envelope And WebSocket Handshake Interaction

The current Protobuf envelope already has:

```text
Session.connection_id
Session.session_id
Session.player_id
Session.connection_epoch
```

This standard does not change those fields.

Rules:

- Existing envelope session fields remain metadata carriers until validation exists.
- Adding credential fields, token fields, authentication result fields, or new handshake/system messages requires a protocol change spec and ADR.
- A future token carrier may be envelope metadata, a system message, a first request payload, a WebSocket subprotocol/header pattern, or another ratified design. This standard does not choose one.
- WebSocket transport must not parse or validate credentials unless a future handshake standard gives it a narrow transport-level responsibility.
- Even if handshake authentication is later adopted, application dispatch must still receive a normalized request identity that domain permissions can inspect.

## 11. Reference Pattern Map

### Nakama Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Account/user as first-class backend capability | Adopted | vibit already treats player identity and account lifecycle as first-class player module concerns. |
| Multiple authentication methods | Deferred | Login methods affect contracts, storage, dependencies, security posture, and public API shape. |
| Session token concept | Deferred | Token format, issuer, verifier, expiration, and validation behavior require separate ratification. |
| Refresh token concept | Deferred | Refresh, rotation, revocation, storage, and replay handling require separate ratification. |
| Realtime socket bound to authenticated session | Adapted | vibit requires request identity validation before production-sensitive domain dispatch, but WebSocket handshake binding remains undecided. |
| Session expiration and revocation vocabulary | Adopted as design dimensions | These must be considered in future token/session work, but no behavior is implemented here. |
| Direct Nakama public API compatibility | Rejected for now | vibit defines an agent-native contract surface; compatibility requires a future ADR. |

### Pitaya Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Session object separate from acceptor/transport | Adopted | vibit keeps transport connection metadata separate from application request identity and future runtime sessions. |
| Handler receives session context | Adapted | vibit domain handlers receive application-owned `RequestIdentity`, not Pitaya API session objects. |
| Session binding vocabulary | Adopted as vocabulary | Binding is useful, but vibit will bind through ratified validation results and manifests. |
| Frontend/backend server role split | Deferred | Distributed topology remains deferred until single-process boundaries and checks are stable. |
| Realtime session context | Adapted | vibit keeps realtime connection metadata, session validation, and domain permissions separate. |
| Direct Pitaya public API compatibility | Rejected for now | vibit may learn from Pitaya architecture vocabulary without copying its API surface. |

## 12. Required Future Artifacts

Before implementing production authentication, future work must add or update the relevant subset of:

- Change spec under `changes/`.
- ADR when the decision affects long-term architecture.
- Contract source files under `contracts/`.
- `.arch/contracts.yaml` when public commands, queries, events, errors, or permissions change.
- `.arch/runtime.yaml` for runtime ownership and implementation state.
- `.arch/dependencies.yaml` and dependency adoption records when new foundational dependencies are introduced.
- Module manifests and module guides.
- Database schema boundary and migrations when data is added.
- Protobuf source and generated output only after protocol impact is ratified.
- Runtime tests and repository checks.
- English documentation and Simplified Chinese translation.

## 13. Ask-First Boundaries

Ask the maintainer before:

- Choosing guest, device, email/password, custom ID, social login, OAuth, OIDC, or another external identity provider.
- Choosing JWT, opaque tokens, refresh tokens, signing, expiration, revocation, rotation, or token storage behavior.
- Adding credential storage, password hashing, cryptography, OAuth, OIDC, external identity, or session-store dependencies.
- Adding authentication runtime code, token parsing, credential lookup, external identity linking, or session persistence.
- Changing Protobuf envelope behavior.
- Changing WebSocket handshake authentication behavior.
- Adding runtime player account handlers or WebSocket routes.
- Declaring metadata-only `player_id` or `session_id` sufficient for production permissions.
- Copying Nakama or Pitaya public API shape.

## 14. Remaining M-011 Work Queue

The remaining authentication design milestone should proceed in bounded steps:

1. Add architecture checks that enforce the authentication/token/session design boundary.
2. Ratify semantic contract dimensions for authentication proof and token/session validation without choosing a concrete token format.
3. Define credential storage and external identity linking boundaries without adding schema or dependencies.
4. Define session persistence and WebSocket handshake decision gates without implementing either path.
5. Close the milestone or create a confirmation gate for the first implementation direction.

Each step must preserve metadata-only identity as non-authenticated until a real validator is separately ratified and implemented.

## 15. Verification

Default repository verification for this standard:

```bash
node tools/vibit check architecture --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

`node tools/vibit check runtime --json` includes `runtime.authentication_token_session_boundary`. That rule statically checks the standard/ADR references, implementation-status markers, metadata-only validator markers, Protobuf source and generated-output boundaries, WebSocket transport boundary, player repository boundary, and player account migration boundary without requiring live external services.

Go tests are required only when runtime Go code changes. This design standard does not require Go runtime code changes.

## 16. Agent Rules

Agents must:

- Read this standard before adding authentication, token, credential, external identity, session persistence, request identity, WebSocket handshake, or runtime player route behavior.
- Preserve separation between authentication proof, login methods, tokens, credentials, external identity links, runtime sessions, request identity, transport connections, and player account lifecycle.
- Use change specs, contracts, manifests, ADRs, and checks before implementing public behavior.
- Treat current metadata-only identity as unauthenticated.
- Record every Nakama or Pitaya pattern as adopted, adapted, deferred, or rejected when used for planning.

Agents must not:

- Add authentication shortcuts to WebSocket transport, Protobuf bridges, player repositories, or domain handlers.
- Treat envelope `player_id` or `session_id` as proof.
- Store credentials, tokens, provider subjects, or sessions in player account lifecycle tables.
- Add token, credential, OAuth, OIDC, password hashing, or session-store dependencies without ratification.
- Change the Protobuf envelope or WebSocket handshake behavior during design-only work.
