# Login Method And Token Format Ratification Standard

Status: Draft v0.1
Last updated: 2026-05-14
Scope: First production login-method selection, token model selection, proof carrier boundaries, lifecycle semantics, storage implications, and pre-implementation gates
Depends on: `docs/authentication-token-session-validation.md`
Canonical decision: `ADR-0024`

The paired Simplified Chinese translation is `docs/login-method-token-format-ratification.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines how vibit ratifies its first production login methods and token model before implementation.

The goal is not to add authentication code quickly. The goal is to make the first authentication choices explicit enough that future agents can implement them without smuggling security behavior into transport handlers, Protobuf adapters, player account persistence, or domain modules.

This standard may be used to ratify:

- The first login-method set.
- The first token model.
- Token and proof carrier boundaries.
- Token lifecycle semantics.
- Required contracts, schemas, dependencies, checks, and tests.
- A bounded implementation queue.

This standard does not by itself implement:

- Runtime authentication.
- Token parsing, signing, validation, refresh, revocation, rotation, replay handling, or storage.
- Credential storage.
- External identity linking.
- Session persistence.
- Protobuf envelope changes.
- WebSocket handshake authentication.
- Runtime player account handlers.
- WebSocket routes.

## 2. Required Reading

Read this standard together with:

- `docs/authentication-token-session-validation.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/session-persistence-websocket-handshake-decision-gates.md`
- `docs/player-identity-session-boundary.md`
- `docs/player-account-session-contracts.md`
- `docs/game-protocol.md`
- `docs/runtime-protocol-adapter.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/player/module.yaml`
- `runtime/AGENTS.md`
- `ADR-0019`
- `ADR-0021`
- `ADR-0022`
- `ADR-0023`
- `ADR-0024`

Reference reading:

- Nakama authentication concepts: `https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama session concepts: `https://heroiclabs.com/docs/nakama/concepts/session/`
- Pitaya API and handler vocabulary: `https://pitaya.readthedocs.io/en/latest/API.html`
- Pitaya session and server-role features: `https://pitaya.readthedocs.io/en/stable/features.html`

Nakama and Pitaya are references. They do not govern vibit's public API shape, credential schema, token format, generated file conventions, or agent workflow.

## 3. Ratification Posture

Authentication choices must be made as durable architecture decisions, not as incidental implementation details.

Rules:

- Login methods are public behavior and must be ratified before handlers exist.
- Token format is a security and operations decision and must be ratified before parser or validator code exists.
- Token carrier behavior must be ratified before changing Protobuf envelope fields, WebSocket handshake behavior, or per-request payload contracts.
- Credential storage and external identity linking remain separate boundaries.
- Session persistence remains a separate boundary.
- Current metadata-only `player_id`, `session_id`, and `connection_id` values remain non-authenticated.
- A selected login method does not automatically authorize all dependencies, tables, routes, or protocol changes required by that method.

## 4. Reference Alignment

### Nakama

Nakama demonstrates that a mature game backend commonly supports multiple authentication methods, account linking, session tokens, refresh tokens, session expiration, logout, and realtime socket behavior.

vibit adopts this as capability coverage, but not as an API shape.

Reference positions:

| Nakama concept | vibit position |
| --- | --- |
| Multiple authentication methods | Capability coverage to compare before selecting the first set. |
| Device-style authentication | Candidate login family, adapted only if credential and abuse semantics are explicit. |
| Email/password authentication | Candidate login family, deferred unless password hashing, reset, rate limit, and recovery gates are ratified. |
| Social/provider authentication | Candidate login family, deferred until provider subject and external identity linking are ratified. |
| Custom identifier authentication | Candidate login family, allowed only with issuer/trust boundaries. |
| Session token | Candidate token concept, not automatically Nakama-compatible. |
| Refresh token | Candidate lifecycle concept requiring storage, rotation, revocation, replay, and redaction decisions. |
| Session variables | Deferred; must not become hidden authority inside opaque token claims. |
| Logout invalidating tokens | Required comparison dimension for any token model. |
| Realtime socket bound to an authenticated session | Adapted through vibit's application-owned request identity and future session validation model. |

### Pitaya

Pitaya demonstrates Go game server vocabulary around handlers, request context, session binding, frontend/backend roles, message routing, push, and session lifecycle.

vibit adopts useful vocabulary, but not Pitaya's API shape or cluster assumptions.

Reference positions:

| Pitaya concept | vibit position |
| --- | --- |
| Handler receives request context containing session | Adapted as `RequestIdentity` and `SessionValidationResult` before domain dispatch. |
| Session binding to user ID | Candidate future session-binding concept, not a token format. |
| Session data storage between requests | Deferred until session persistence is ratified. |
| Frontend/backend session distinction | Deferred until distributed runtime planning. |
| Route-aware routing function | Reference for future route policy mapping; it must not bypass module contracts. |
| Push to bound users | Future realtime capability, not part of first authentication ratification unless required by session lifecycle tests. |

## 5. Login-Method Candidate Families

Future work must compare candidate login methods before selecting the first production set.

| Family | Description | Main benefits | Main risks | Required gates before implementation |
| --- | --- | --- | --- | --- |
| Device credential login | Client proves possession of a device-scoped high-entropy secret or identifier. | Low friction for games, useful first account entry. | Device identifiers can change; weak identifiers are replayable; account recovery and abuse controls matter. | Credential semantics, storage/hashing, account creation/linking policy, replay controls, rate limits, tests. |
| Guest or anonymous login | Server creates a temporary or low-assurance actor without durable external proof. | Fast onboarding and good for early game entry. | Abuse, data ownership, upgrade behavior, and production permission scope must be strict. | Anonymous actor rules, upgrade path, expiration, anti-abuse posture, permission limits. |
| Custom ID login | A trusted issuer maps a custom subject to a player. | Useful for integrating with an existing game platform or studio identity. | Unsafe if arbitrary clients can mint IDs; issuer trust must be explicit. | Issuer model, trusted caller boundary, subject collision rules, linking policy. |
| Email/password login | User authenticates with email or username plus password. | Familiar and recoverable. | Requires password hashing, reset flows, rate limits, breach posture, and secret handling. | Password dependency adoption, credential schema, reset/recovery flow, rate limiting, redaction. |
| External provider login | User authenticates through a platform, social, OAuth, OIDC, or game-provider account. | Strong platform fit and cross-device identity. | Provider dependencies, issuer/audience validation, account linking conflicts, token validation, availability. | Provider subject schema, dependency adoption, issuer/audience validation, linking/conflict/recovery behavior. |
| Service authentication | Non-player service proves authority. | Needed for future internal operations and server-to-server work. | Easy to over-privilege; separate from player login. | Service actor model, permissions, key management, rotation, audit events. |

Selection rules:

- The first set should be intentionally small.
- A selected login method must have an explicit confidence level and known gaps.
- A selected login method must state whether it creates player accounts, links to existing accounts, or only authenticates existing accounts.
- A selected login method must state whether it is production-capable, bootstrap-only, or local-development-only.
- No selected method may place credential parsing in WebSocket transport or Protobuf adapter code.

## 6. Token Model Dimensions

Future work must compare token models across these dimensions.

### Token Kinds

Candidate token kinds:

- Access token.
- Refresh token.
- Session token.
- One-time proof token.
- Service token.
- External provider token.

Rules:

- A token kind must state whether it is client-presented, server-stored, derived from an external provider, or only used inside a trusted server boundary.
- Access and refresh token behavior must be separate even if both are opaque strings.
- A refresh token must not be introduced without rotation, revocation, storage, replay, and redaction semantics.
- External provider tokens must not become vibit session tokens by convenience.

### Token Formats

Candidate formats:

| Format | Benefits | Risks | Required gates |
| --- | --- | --- | --- |
| Opaque high-entropy token | Simple to redact, easy to revoke server-side, implementation can remain storage-backed and explicit. | Requires lookup storage and careful hashing. | Token generation, hash storage, lookup index, expiration, rotation, cleanup. |
| Signed structured token | Can validate without storage lookup and can carry claims. | Key management, revocation, replay, claim drift, and secret rotation are harder. | Signing dependency, key management, issuer/audience/claims, revocation story. |
| External provider token | Reuses provider proof. | Provider-specific validation and availability; audience and issuer errors can be severe. | Provider validation, trust boundary, subject mapping, failure handling. |
| Plain session ID as token | Simple vocabulary. | Easy to confuse identifier with proof; unsafe unless high entropy and protected as a secret. | Must be treated as secret token, not as metadata-only `session_id`. |

Selection rules:

- A token format must state issuer, verifier, subject, audience, entropy or signing model, expiration, refresh, revocation, rotation, replay, redaction, and storage implications.
- Token strings must never be stored in plaintext unless a future decision explicitly grants an exception with rationale.
- Logs, errors, traces, conversation logs, and change specs must redact token values.
- Token validation errors must map to registered error dimensions before implementation.

## 7. Proof Carrier Boundaries

A proof carrier is where credentials, tokens, or session proof appear on the wire or inside a request.

Candidate carriers:

| Carrier | Layer touched | Position |
| --- | --- | --- |
| Login command payload | Application/protocol payload after decode | Candidate first carrier for credential exchange because it keeps WebSocket transport credential-neutral. |
| Per-request payload field | Application/protocol payload after decode | Candidate for explicit validation, but verbose and easy to duplicate. |
| Protobuf envelope extension | Wire envelope | Requires protocol decision; do not overload current metadata fields. |
| Current `Session.session_id` field | Existing envelope metadata | Metadata-only today; must not become proof without protocol ratification. |
| WebSocket handshake header | Transport/process boundary | Requires handshake-auth decision and normalized identity handoff. |
| WebSocket subprotocol, cookie, or query parameter | Transport/process boundary | Requires explicit risk analysis; not available by convenience. |
| First system message | Protocol/application after connection open | Requires system-message contract and connection binding rules. |

Rules:

- The current Protobuf `Session` message remains metadata-only until a future protocol decision says otherwise.
- Do not put access tokens or refresh tokens into `player_id`, `session_id`, `connection_id`, `connection_epoch`, route names, request IDs, or target IDs.
- WebSocket transport remains credential-neutral until a future handshake decision grants a narrow role.
- Carrier selection must explain how proof becomes `RequestIdentity` before domain handlers run.

## 8. Ratification Packet

Every future work item that selects a login method, token model, token carrier, or lifecycle rule must leave a ratification packet in its change spec or documentation.

The packet must include:

- Selected option.
- Rejected plausible options.
- Reference alignment with Nakama and Pitaya.
- Public rationale.
- Decision weights.
- Security and abuse notes.
- Public contract impact.
- Protobuf impact.
- WebSocket handshake impact.
- Persistence and migration impact.
- Dependency adoption impact.
- Error, permission, and audit-event impact.
- Test and repository-check impact.
- Redaction rules.
- Reversal conditions.
- Next implementation gate.

## 9. Required Artifacts Before Implementation

Before implementation code for a selected login method or token model exists, the repository must have:

- Semantic contracts for login commands, token refresh, logout, validation, and any account-linking behavior in scope.
- Error catalog entries for invalid proof, expired proof, revoked proof, replayed proof, malformed proof, credential mismatch, account disabled, rate limited, and dependency unavailable where applicable.
- Permission catalog entries for player login, token refresh, logout, account linking, service-auth validation, and administrative revocation where applicable.
- Schema boundaries for credential records, token/session records, or external identity links if storage is required.
- Migration sources only after schema boundaries are ratified.
- Dependency adoption records for hashing, signing, OAuth/OIDC, provider SDKs, Redis-like stores, or key-management libraries.
- Redaction and logging rules for credential and token material.
- Runtime ownership rules for validator interfaces and storage adapters.
- Repository checks guarding forbidden shortcut implementations.
- Focused tests for success, invalid proof, expired proof, revoked proof, replay/collision behavior, redaction, and boundary ownership.
- Bilingual documentation updates.

## 10. M-013 Work Queue

M-013 proceeds in bounded steps:

1. Define this ratification standard.
2. Compare first login-method candidates.
3. Ratify the first login-method set.
4. Compare token format and token carrier options.
5. Ratify the first token format and proof carrier posture.
6. Define token lifecycle and storage implications.
7. Define authentication contract, error, and permission surfaces.
8. Define credential, token, and session schema gates.
9. Add repository checks for selected login/token boundaries.
10. Close the milestone and create the next implementation gate.

This queue may be adjusted by a future work item if verification or reference analysis proves that a step is too broad or too narrow. It must not be collapsed into implementation work by convenience.

## 11. Repository Rules

Agents must not:

- Implement login handlers from this standard alone.
- Parse or validate tokens from this standard alone.
- Add token, credential, external identity, or session tables from this standard alone.
- Add password hashing, JWT, OAuth, OIDC, provider SDK, cryptography, key-management, Redis-like, or other major authentication dependencies from this standard alone.
- Change Protobuf envelope behavior from this standard alone.
- Change WebSocket handshake authentication behavior from this standard alone.
- Add runtime player account command/query handlers or WebSocket routes from this standard alone.
- Treat metadata-only identity as production proof.
- Copy Nakama or Pitaya public APIs without a compatibility ADR.

Agents may:

- Add comparison docs.
- Add decision records.
- Add change specs.
- Update manifests.
- Plan contracts, schemas, checks, and tests.
- Strengthen repository checks that prevent shortcut implementations.

## 12. Verification

For changes under this standard, default verification includes:

```bash
node tools/vibit inspect next --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check memory --json
node tools/vibit check change <change-id> --json
node tools/vibit check all --json
git diff --check
```

Runtime Go tests are required only when runtime code changes.

Live PostgreSQL verification is required only when persistence behavior or migrations change; otherwise record it as not applicable.

## 13. Exit Criteria

M-013 can close only when:

- The first login-method set is selected or explicitly deferred with a next gate.
- The first token model is selected or explicitly deferred with a next gate.
- Proof carrier posture is selected or explicitly deferred without implicit Protobuf or WebSocket behavior changes.
- Token lifecycle semantics are recorded.
- Required contracts, schemas, dependency records, checks, and tests are planned.
- No implementation boundary was crossed accidentally.
- `node tools/vibit check all --json` passes.
