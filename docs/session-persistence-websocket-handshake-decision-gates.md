# Session Persistence And WebSocket Handshake Decision Gates

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Session persistence decision gates, WebSocket handshake authentication decision gates, request-level validation options, reconnect gates, connection epoch gates, Protobuf envelope interaction gates, and future implementation artifacts
Depends on: `docs/authentication-token-session-validation.md`

The paired Simplified Chinese translation is `docs/session-persistence-websocket-handshake-decision-gates.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines the decision gates that must exist before vibit implements session persistence or WebSocket handshake authentication.

The goal is to keep future session work explicit enough for agents to extend without hiding authentication or binding behavior in transport, protocol, or domain modules.

This standard does not choose:

- Request-level validation as the production model.
- First-message validation as the production model.
- WebSocket handshake-level validation as the production model.
- Every-request validation as the production model.
- A hybrid validation model.
- A session store.
- Session tables or migrations.
- Token carrier behavior.
- Protobuf envelope changes.
- WebSocket handshake/system messages.
- WebSocket handshake authentication behavior.
- Route-level authentication implementation.

## 2. Required Reading

Read this standard together with:

- `docs/authentication-token-session-validation.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/player-identity-session-boundary.md`
- `docs/game-protocol.md`
- `docs/runtime-protocol-adapter.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `.arch/reference.yaml`
- `runtime/AGENTS.md`
- `ADR-0015`
- `ADR-0018`
- `ADR-0021`
- `ADR-0023`

Reference reading:

- Nakama authentication and sessions concepts: `https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama realtime socket concepts: `https://heroiclabs.com/docs/nakama/concepts/multiplayer/`
- Pitaya session, handler, and frontend/backend vocabulary: `https://pitaya.readthedocs.io/`

Nakama and Pitaya are references. They do not govern vibit's public API shape, validation model, session persistence model, envelope behavior, WebSocket handshake behavior, or agent workflow.

## 3. Core Vocabulary

### Session Persistence

Session persistence is server-side storage of logical session state beyond a single request or in-memory validation step.

Possible future state includes `session_id`, `actor_kind`, `actor_id`, `player_id`, claims, expiration, revocation status, refresh linkage, connection binding metadata, and audit metadata.

Rules:

- Session persistence is not implemented by this standard.
- Persisted sessions are not player account lifecycle rows.
- Persisted sessions are not WebSocket connections.
- Persisted sessions are not Protobuf envelope metadata.
- Session storage must not be added by convenience inside player, inventory, WebSocket transport, or Protobuf adapter code.

### WebSocket Handshake Authentication

WebSocket handshake authentication means validating proof before or during WebSocket connection establishment.

Rules:

- Handshake authentication is not implemented by this standard.
- Current WebSocket transport remains credential-neutral.
- WebSocket transport must not inspect `Authorization`, cookies, subprotocol authentication carriers, token query parameters, or credential headers unless a future handshake standard explicitly grants a narrow transport responsibility.
- Even if future handshake authentication is adopted, application dispatch must still receive normalized request identity.

### Request-Level Validation

Request-level validation means validating proof or session identity after an envelope is decoded and before domain dispatch for each request.

Rules:

- This is a future option, not a selected production model.
- It aligns with the current application-owned `SessionValidator` hook.
- It may be combined with other validation gates later.

### First-Message Validation

First-message validation means a connected client sends a protocol message after WebSocket connection establishment to authenticate or bind the connection before normal gameplay routes are accepted.

Rules:

- This is a future option, not a selected production model.
- It would require explicit protocol/system-message contracts before implementation.
- It must not be invented through ad hoc domain routes.

### Every-Request Validation

Every-request validation means every command or query carries enough proof or session metadata to validate independently before domain dispatch.

Rules:

- This is a future option, not a selected production model.
- It can reduce reliance on connection state, but it may increase validation cost and protocol verbosity.
- Token/session carrier behavior must be ratified before implementation.

### Hybrid Validation

Hybrid validation combines multiple validation gates.

Examples:

- Handshake validation plus request-level permission checks.
- First-message binding plus periodic session revalidation.
- Every-request validation for sensitive routes plus cached session context for low-risk routes.

Rules:

- This is a future option, not a selected production model.
- A hybrid model must declare which layer owns each step and how failures are handled.

## 4. Validation Model Decision Gate

Before selecting a production validation model, future work must compare at least these options:

| Option | Layer touched | Benefits | Risks | Required before implementation |
| --- | --- | --- | --- | --- |
| Request-level validation | Application dispatch after protocol decode | Clear domain handoff, works with multiple transports, fits current `SessionValidator` hook | May validate repeatedly unless cache/session rules exist | Validation contract, proof/session carrier, failure behavior, route policy mapping |
| First-message validation | Protocol/application after connection open | Keeps transport credential-neutral and supports explicit game protocol negotiation | Requires connection-bound state and system-message contracts | System-message contract, connection binding model, timeout/error behavior, reconnect rules |
| Handshake-level validation | WebSocket transport/process boundary | Rejects unauthenticated connections early and can reduce unauthenticated connection load | Risks putting auth into transport and complicates non-WebSocket transports | Transport auth boundary, header/subprotocol/cookie/query carrier decision, normalized identity handoff |
| Every-request validation | Protocol/application on each request | Stateless-friendly and clear per-request proof semantics | Higher overhead, larger envelopes or payloads, carrier leakage risk | Token/session carrier contract, replay handling, failure/retryability rules |
| Hybrid validation | Multiple layers | Can balance early rejection, connection binding, and route sensitivity | Easy to blur ownership unless manifests are precise | Layer ownership matrix, cache invalidation, failure precedence, verification plan |

Choosing one of these options is an ask-first decision.

Until a model is selected:

- Current metadata-only validation remains non-authenticated.
- `Session.connection_id`, `Session.session_id`, `Session.player_id`, and `Session.connection_epoch` remain metadata carriers only.
- Domain modules must not treat metadata-only identity as production proof.
- WebSocket transport remains credential-neutral.

## 5. Session Store Decision Gate

Before adding a session store, future work must define:

- Whether sessions are persisted server-side.
- Whether the store is PostgreSQL, memory, Redis-like, another ratified store, or no store.
- Session ID generation semantics.
- Lookup semantics.
- Expiration semantics.
- Refresh semantics.
- Revocation and logout semantics.
- Rotation and replay behavior, if tokens are involved.
- Connection binding and connection epoch behavior.
- Cleanup and migration behavior.
- Transaction boundaries with account lifecycle or token changes.
- Opt-in live verification requirements.

Rules:

- No session store is selected by this standard.
- No session table may be added by this standard.
- PostgreSQL is the first authoritative durable store for module state, but that does not automatically make PostgreSQL the session store.
- Memory may be useful for tests or local bootstrap, but that does not automatically make memory the production session store.
- Redis-like storage remains a possible future option only after dependency adoption and architecture ratification.

## 6. WebSocket Handshake Decision Gate

Before adding WebSocket handshake authentication, future work must define:

- Which proof carrier is used during or before the handshake.
- Whether the carrier is a header, cookie, query parameter, WebSocket subprotocol, mTLS-like service proof, or another ratified input.
- How browser and non-browser clients are supported.
- Whether failed proof rejects the handshake or accepts the connection into an anonymous/non-authenticated state.
- Which component creates normalized request identity from handshake validation.
- How identity is revalidated after handshake.
- How logout, revocation, expiration, refresh, reconnect, and connection migration affect active WebSocket connections.
- How errors are surfaced when the connection is rejected before a Protobuf envelope exists.
- What tests prove transport code stays narrow.

Rules:

- WebSocket transport must not own player account lookup.
- WebSocket transport must not own credential lookup.
- WebSocket transport must not own long-lived session persistence.
- Any future transport-level proof extraction must hand off to application-owned or auth-owned validation contracts.
- Route-level domain handlers must not assume handshake identity unless application dispatch passes normalized request identity.

## 7. Protobuf Envelope Interaction Gate

The current envelope includes:

```text
Session.connection_id
Session.session_id
Session.player_id
Session.connection_epoch
```

This standard does not change those fields.

Before changing envelope behavior, future work must define:

- Whether token/session proof is carried in envelope metadata, a system message, a module payload, a WebSocket handshake carrier, or another design.
- Whether `session_id` is a candidate identifier, a validated identifier, or both in different states.
- Whether `player_id` can appear before validation and how it is marked as untrusted.
- Whether `connection_epoch` is server-issued, client-presented, or both.
- Whether new fields require a Protobuf package version change.
- Generated output impact.
- Backward compatibility behavior.
- Error behavior for unsupported, missing, malformed, expired, or revoked proof.

Rules:

- Envelope fields remain metadata carriers until validation exists.
- Adding token fields, credential fields, authentication result fields, handshake fields, or new system messages requires a protocol change spec and ADR.
- Generated Protobuf output must not be hand-edited.

## 8. Reconnect And Connection Epoch Gate

Before implementing reconnect or connection epoch behavior, future work must define:

- Whether reconnect is tied to session persistence, token proof, connection IDs, or another proof.
- Who issues `connection_epoch`.
- Whether epoch increases on every connection, reconnect, migration, or only after successful validation.
- Whether stale epochs are rejected, ignored, or treated as hints.
- How duplicate active connections for the same session are handled.
- Whether previous connections are closed, replaced, or allowed to coexist.
- How room, match, party, presence, and stream membership will be restored later.
- What replay protection is required.

Rules:

- Current `connection_epoch` remains metadata only.
- Current WebSocket transport does not own reconnect semantics.
- Future reconnect behavior must be compatible with later multiplayer and presence modules.

## 9. Reference Pattern Map

### Nakama Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Session token after authentication | Deferred | Token format, issuance, validation, expiration, revocation, and carrier behavior remain separate choices. |
| Refresh token concept | Deferred | Refresh, rotation, replay, logout, and storage behavior require future ratification. |
| Realtime socket associated with an authenticated session | Adapted | vibit requires normalized request identity before production-sensitive dispatch, but does not yet choose handshake-level binding. |
| Session expiration and revocation vocabulary | Adopted as design dimensions | Future session work must handle these dimensions explicitly. |
| Direct Nakama API compatibility | Rejected for now | vibit defines agent-native contracts and may only adopt compatibility through a future ADR. |

### Pitaya Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Session object separate from acceptor | Adopted as vocabulary | Transport connection state must remain separate from application request identity. |
| Handler receives session context | Adapted | vibit handlers receive `RequestIdentity` through application dispatch instead of a transport-owned session object. |
| Session binding | Adopted as vocabulary | Binding is useful, but vibit requires ratified validation results and manifests. |
| Frontend/backend split | Deferred | Distributed routing and cluster topology remain future work. |
| Direct Pitaya API compatibility | Rejected for now | Architecture vocabulary may inform vibit without copying public APIs. |

## 10. Future Artifact Gates

Before implementing session persistence, future work must add or update:

- Change spec under `changes/`.
- ADR when selecting validation model, store, reconnect behavior, or transport responsibilities.
- Session persistence standard.
- Session contract sources.
- Error catalog and retryability rules.
- Permission catalog.
- Store ownership manifest.
- Schema boundary and migrations, when data is stored.
- Session cleanup and expiration plan.
- Tests for creation, lookup, expiration, revocation, logout, refresh, reconnect, and failure behavior.
- Repository checks when static enforcement is possible.
- English documentation and Simplified Chinese translation.

Before implementing WebSocket handshake authentication, future work must add or update:

- Change spec under `changes/`.
- ADR for handshake carrier and validation model.
- WebSocket handshake authentication standard.
- Transport/auth handoff contract.
- Protocol or non-protocol error surface definition.
- Browser and non-browser client compatibility notes.
- Tests proving transport responsibilities stay narrow.
- Tests proving application dispatch receives normalized request identity.
- Repository checks when static enforcement is possible.
- English documentation and Simplified Chinese translation.

## 11. Ask-First Boundaries

Ask the maintainer before:

- Choosing request-level, first-message, handshake-level, every-request, or hybrid validation as the production model.
- Choosing a session persistence store.
- Adding session tables or migrations.
- Choosing session expiration, refresh, revocation, reconnect, or connection epoch behavior.
- Choosing token/session carrier behavior.
- Changing Protobuf envelope fields or adding handshake/system messages.
- Adding WebSocket handshake authentication behavior.
- Adding route-level authentication implementation.
- Declaring metadata-only `player_id`, `session_id`, or `connection_id` sufficient for production permissions.

## 12. Verification

Default repository verification for this standard:

```bash
node tools/vibit check architecture --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change define-session-persistence-websocket-handshake-decision-gates --json
node tools/vibit check all --json
```

Go tests are required only when runtime Go code changes. This decision-gate standard does not require Go runtime code changes.

## 13. Agent Rules

Agents must:

- Read this standard before adding session persistence, WebSocket handshake authentication, reconnect behavior, connection epoch behavior, token/session carriers, or session-related protocol changes.
- Keep the current WebSocket transport credential-neutral until a future handshake standard is ratified.
- Keep current envelope session fields metadata-only until validation exists.
- Keep domain modules dependent on normalized request identity, not token, session, credential, or transport internals.
- Record reference patterns from Nakama and Pitaya as adopted, adapted, deferred, or rejected when used for planning.
- Record verification honestly.

Agents must not:

- Choose a validation model implicitly.
- Add session tables implicitly.
- Treat PostgreSQL, memory, Redis-like storage, or any other store as selected for sessions without ratification.
- Parse tokens, credentials, cookies, `Authorization` headers, WebSocket subprotocols, or query-parameter proof inside WebSocket transport without a future handshake standard.
- Add token, credential, authentication result, handshake, or system-message envelope fields from this standard alone.
- Treat metadata-only `player_id`, `session_id`, `connection_id`, or `connection_epoch` as authenticated proof.
