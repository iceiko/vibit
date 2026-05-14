# Credential Storage And External Identity Linking Boundaries

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Credential storage boundaries, external identity linking boundaries, player account lifecycle separation, login-method family deferral, provider subject deferral, and future implementation gates
Depends on: `docs/authentication-token-session-validation.md`

The paired Simplified Chinese translation is `docs/credential-storage-external-identity-linking-boundaries.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines what vibit means by credential storage and external identity linking before either is implemented.

The purpose is to keep future authentication work agent-readable and bounded. A future agent must be able to see that player account lifecycle storage, credential storage, provider identity linking, token/session behavior, and runtime request identity validation are separate responsibilities.

This standard does not choose:

- A supported login method.
- A password model.
- A credential schema.
- An external identity provider.
- OAuth, OIDC, social login, device login, guest login, or custom ID behavior.
- Password hashing, encryption, signing, key management, or provider dependencies.
- Credential tables, external identity tables, token tables, session tables, or migrations.
- Runtime credential lookup, account linking handlers, recovery flows, merge behavior, or WebSocket routes.

## 2. Required Reading

Read this standard together with:

- `docs/authentication-token-session-validation.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/player-identity-session-boundary.md`
- `docs/player-account-session-contracts.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `modules/player/module.yaml`
- `runtime/AGENTS.md`
- `ADR-0019`
- `ADR-0021`
- `ADR-0022`
- `ADR-0023`

Reference reading:

- Nakama authentication concepts: `https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama account/session capability surface: `https://heroiclabs.com/docs/nakama/`
- Pitaya session and route-handler vocabulary: `https://pitaya.readthedocs.io/`

Nakama and Pitaya are references. They do not govern vibit's public API shape, credential schema, linking behavior, generated file conventions, or agent workflow.

## 3. Core Vocabulary

### Credential

A credential is secret or proof material used by a login method or identity provider.

Examples include passwords, password hashes, device secrets, custom ID secrets, provider secrets, OAuth credentials, OIDC subjects, provider-issued identity material, or future service credentials.

Rules:

- Credentials are not player account lifecycle rows.
- Credentials are not runtime sessions.
- Credentials are not WebSocket connection metadata.
- Credentials are not Protobuf envelope session metadata.
- Credentials must not be stored in `player_accounts` or `player_account_events`.
- Credential storage remains unimplemented until a future standard ratifies schema, dependencies, redaction, tests, and ownership.

### External Identity Link

An external identity link maps a vibit-owned account identity to a provider namespace and provider subject.

Examples include future provider subjects for platform accounts, device identities, social identity providers, OIDC issuers, custom identity providers, or game-platform accounts.

Rules:

- External identity links are not player account lifecycle rows.
- External identity links are not credentials by themselves, though they may be related to a credential or login method.
- Provider subject semantics must be defined before storage exists.
- Link, unlink, conflict, recovery, and merge behavior remain unimplemented.
- External identity links must not be added to `player_accounts` or `player_account_events` by convenience.

### Player Account Lifecycle

Player account lifecycle storage owns stable account identity state.

Current lifecycle tables:

```text
player_accounts
player_account_events
```

Rules:

- These tables remain account lifecycle storage only.
- They may record lifecycle state such as created, disabled, or deleted account status as ratified by the player account PostgreSQL schema boundary.
- They must not store credential material, provider subjects, token state, runtime sessions, WebSocket connection state, or request validation results.

### Login Method Family

A login method family is a future category of ways to produce authentication proof.

Deferred families include:

- guest or anonymous-login family
- device-login family
- email/password family
- custom ID family
- social-login family
- OAuth family
- OIDC family
- external identity-provider family
- service-auth family

Rules:

- No family is selected by this standard.
- Listing a family is capability coverage only, not implementation permission.
- Every selected family must get a future contract, schema boundary, dependency review, error model, permission model, and verification path before implementation.

## 4. Ownership Boundaries

### Player Module

Owner:

```text
modules/player
runtime/internal/modules/player
```

Owns:

- stable `player_id`
- player account lifecycle contracts
- player account lifecycle repository interfaces
- ratified player account lifecycle persistence boundaries

Does not own:

- credential storage
- password hashing
- external identity provider subjects
- token issuance or validation
- runtime session persistence
- WebSocket connection binding
- request validation results

### Future Credential Boundary

Owner status:

```text
planned, not implemented
```

Future work must define:

- Credential owner module or runtime subsystem.
- Credential record lifecycle.
- Secret material and non-secret metadata.
- Hashing, encryption, signing, or provider dependency adoption.
- Redaction and logging rules.
- Access permissions.
- Failure classes and retryability.
- Migration ownership.
- Tests and repository checks.

### Future External Identity Boundary

Owner status:

```text
planned, not implemented
```

Future work must define:

- Provider namespace semantics.
- Provider subject semantics.
- Link and unlink authority.
- Duplicate provider-subject behavior.
- Conflict behavior when a provider subject maps to an existing account.
- Account recovery behavior.
- Account merge behavior, if any.
- Provider metadata retention and redaction rules.
- Audit events and client-visible events.

### Runtime Session Validation

Owner:

```text
runtime/internal/app
```

Rules:

- Runtime session validation may consume future authentication proof or session validation results only after those contracts are ratified.
- It must not query credential stores or external identity stores directly unless a future boundary explicitly grants that responsibility.
- Domain modules receive normalized request identity; they do not parse credentials or provider subjects.

## 5. Deferred Login-Method Coverage

Nakama demonstrates that a production game backend commonly supports multiple authentication methods and account/session concepts.

vibit adopts that as capability coverage, but defers concrete login-method selection.

| Capability family | Vibit position | Reason |
| --- | --- | --- |
| Device-style login | Deferred | Requires device identifier semantics, replay controls, secret treatment, account recovery behavior, and abuse controls. |
| Email/password login | Deferred | Requires password hashing, password reset, rate limiting, credential storage, secret redaction, and recovery flows. |
| Custom ID login | Deferred | Requires issuer semantics, collision behavior, account linking rules, and trusted caller boundaries. |
| Social/provider login | Deferred | Requires provider namespace, subject semantics, provider metadata, conflict behavior, and dependency adoption. |
| OAuth/OIDC-style login | Deferred | Requires provider trust, issuer/audience validation, key management, token validation, refresh behavior, and dependency adoption. |
| Guest/anonymous login | Deferred | Requires explicit anonymous actor rules, account upgrade behavior, data ownership, and anti-abuse posture. |
| Session token and refresh token concepts | Deferred | Requires token format, issuer, verifier, expiration, revocation, rotation, storage, and replay behavior. |
| Direct Nakama API compatibility | Rejected for now | vibit defines an agent-native contract surface; compatibility requires a future ADR. |

Pitaya demonstrates a useful separation between connection acceptors, sessions, route handlers, and server roles.

vibit adapts that vocabulary as follows:

| Pitaya pattern | Vibit position | Reason |
| --- | --- | --- |
| Session object separate from acceptor | Adopted as vocabulary | Transport connections must remain separate from application request identity and future runtime sessions. |
| Handler receives session context | Adapted | vibit handlers receive application-owned `RequestIdentity`, not a transport-owned session object. |
| Session binding | Adopted as vocabulary | Binding is useful, but must come from ratified validation results rather than raw metadata. |
| Frontend/backend split | Deferred | Distributed topology remains outside the current credential/linking boundary. |
| Direct Pitaya API compatibility | Rejected for now | vibit may learn from architecture vocabulary without copying public APIs. |

## 6. Forbidden Shortcuts

Agents must not:

- Add credential columns to `player_accounts`.
- Add provider subject columns to `player_accounts`.
- Add token or session columns to `player_accounts`.
- Add credential, provider subject, token, session, WebSocket state, or request-validation rows to `player_account_events`.
- Add credential tables or external identity tables without a future schema boundary.
- Add password hashing, OAuth, OIDC, JWT, provider SDK, cryptography, or key-management dependencies from this standard alone.
- Add runtime credential lookup or external identity lookup from this standard alone.
- Add account linking, unlinking, recovery, or merge behavior from this standard alone.
- Infer a login method from a reference project list.
- Treat direct Nakama or Pitaya API compatibility as an unstated goal.

## 7. Future Credential Storage Artifact Gates

Before implementing credential storage, future work must add or update:

- Change spec under `changes/`.
- ADR when the choice affects long-term architecture, provider dependency posture, generated file conventions, or security boundaries.
- Credential owner declaration in manifests.
- Login method contract.
- Credential schema boundary.
- Migration source, if data is stored.
- Secret-handling rules.
- Redaction and logging rules.
- Error catalog and failure-class mapping.
- Permission catalog.
- Tests for common success and failure modes.
- Negative tests proving credentials are not stored in player lifecycle tables.
- Repository checks when the rule can be statically enforced.
- English documentation and Simplified Chinese translation.

The future schema boundary must explicitly answer:

- Which records contain secret material.
- Which records contain non-secret lookup metadata.
- Which identifiers are safe to log.
- Which fields are unique.
- Which fields are mutable.
- Which account lifecycle states block login.
- Which operations require transactions with player account lifecycle changes.

## 8. Future External Identity Linking Artifact Gates

Before implementing external identity linking, future work must add or update:

- Change spec under `changes/`.
- ADR when provider semantics, merge behavior, recovery behavior, or direct API compatibility affect long-term architecture.
- External identity owner declaration in manifests.
- Provider namespace contract.
- Provider subject contract.
- Link command contract.
- Unlink command contract, if unlinking is supported.
- Conflict, duplicate, recovery, and merge semantics.
- External identity schema boundary.
- Migration source, if data is stored.
- Audit event catalog.
- Error catalog and retryability rules.
- Permission catalog.
- Tests for link, unlink, duplicate, conflict, disabled-account, deleted-account, and recovery behavior.
- Repository checks when the rule can be statically enforced.
- English documentation and Simplified Chinese translation.

The future linking boundary must explicitly answer:

- Whether provider subject IDs are globally unique or provider-scoped.
- Whether provider metadata is retained, normalized, redacted, or discarded.
- Whether one account can hold multiple provider links.
- Whether one provider subject can ever map to multiple vibit accounts.
- Whether linked accounts can be merged.
- Whether links can be moved between accounts.
- Which events are security/audit-only and which may be client-visible.

## 9. Relationship To Token And Session Work

Credential storage and external identity linking do not by themselves define runtime session behavior.

Rules:

- A credential can produce authentication proof only through a future ratified login method.
- An external identity link can identify an account only through a future ratified provider validation path.
- A token can be issued only after token format, issuance, expiration, revocation, storage, and validation behavior are ratified.
- A runtime session can be persisted only after the session persistence boundary is ratified.
- A WebSocket connection can be bound to identity only after request-level, first-message, handshake-level, every-request, or hybrid validation behavior is ratified.

## 10. Ask-First Boundaries

Ask the maintainer before:

- Choosing a supported login method.
- Choosing a password model.
- Choosing credential table shape.
- Choosing external identity table shape.
- Choosing provider namespace or provider subject semantics.
- Choosing OAuth, OIDC, social login, provider SDK, password hashing, encryption, signing, or key-management dependencies.
- Adding credential storage, external identity storage, token storage, or session storage.
- Adding runtime credential lookup, account linking handlers, account recovery flows, or account merge behavior.
- Changing the player account lifecycle table shape.
- Copying Nakama or Pitaya public API shape.

## 11. Verification

Default repository verification for this standard:

```bash
node tools/vibit check architecture --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check module player --json
node tools/vibit check change define-credential-storage-external-identity-linking-boundaries --json
node tools/vibit check all --json
```

Go tests are required only when runtime Go code changes. This boundary standard does not require Go runtime code changes.

## 12. Agent Rules

Agents must:

- Read this standard before adding credential storage, external identity linking, login methods, or provider-related behavior.
- Preserve player account lifecycle tables as account lifecycle storage only.
- Keep login-method family lists as deferred capability coverage until a future choice is ratified.
- Map reference patterns from Nakama and Pitaya as adopted, adapted, deferred, or rejected before using them for planning.
- Record verification honestly.

Agents must not:

- Store credentials, provider subjects, tokens, runtime sessions, WebSocket connection state, or request validation results in player account lifecycle tables.
- Choose a login method implicitly.
- Add credential or provider dependencies implicitly.
- Implement account linking, unlinking, recovery, or merge behavior from this boundary standard alone.
- Treat metadata-only `player_id` or `session_id` as authenticated proof.
