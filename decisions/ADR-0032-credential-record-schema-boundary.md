# ADR-0032: Credential Record Schema Boundary

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-credential-record-schema-boundary/`

Related conversations:

- `conversations/2026-05-14-credential-record-schema-boundary.md`

Related artifacts:

- `docs/credential-record-schema-boundary.md`
- `docs/credential-record-schema-boundary.zh-CN.md`
- `docs/credential-token-session-schema-gates.md`
- `docs/first-login-method-set.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/selected-login-token-boundary-checks.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Context

M-014 exists to ratify credential and token verifier record schema boundaries before migrations, repositories, adapters, runtime lookup, handlers, routes, generated authentication shapes, Protobuf messages, WebSocket proof carriers, or authentication behavior are added.

The selected first login method is `device_credential_login`. Earlier gates established that credential storage is required, raw credentials and raw OS device IDs must not be stored, and player account lifecycle tables must remain credential-free.

W-0074 turns the credential record gate into a concrete schema boundary while preserving implementation deferral.

## Decision

Ratify the credential record schema boundary for the first `device_credential_login` posture.

The future logical table is:

```text
authentication_device_credentials
```

The owner is:

```text
runtime.authentication
```

The boundary ratifies:

- One credential record belongs to one player account.
- `player_id` on a credential record is immutable.
- A credential must not move between players.
- The first posture allows at most one active device credential per player.
- Historical records for rotation, revocation, replacement, and audit are allowed.
- Multi-device linking, account recovery, and account merge remain deferred.
- Raw credential material and raw operating-system device IDs are forbidden.
- Credential lookup and verifier material must be non-plaintext, versioned, and redacted.
- `credential_record_id` is the log-safe credential identifier.
- `credential_lookup_digest`, `credential_verifier_digest`, and raw client instance identifiers are not log-safe.
- Credential lifecycle states are `active`, `disabled`, `revoked`, and `replaced`.
- Disabled, deleted, revoked, or replaced credentials cannot authenticate.
- Disabled or deleted player account state blocks login even if the credential is active.
- Rotation creates a new record and marks the old active record as `replaced`; verifier material must not be overwritten in place for rotation.

This decision does not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime lookup, login handlers, token behavior, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## Alternatives Considered

- Add credential migration source immediately.
- Store credential verifier material in `player_accounts`.
- Store raw device identifiers as credentials.
- Use one generic credential table for every future login family immediately.
- Allow many active device credentials per player in the first posture.
- Permit credential records to move between player accounts.
- Ratify token verifier schema together with credential schema in one large work item.
- Copy Nakama device-authentication API and storage assumptions.
- Treat Pitaya session binding as credential persistence.

## Rationale

The first credential boundary should be explicit and narrow. A device credential is a proof verifier record, not a player account lifecycle row, not a provider identity link, not an access token, and not a runtime session.

The future logical table name is specific to device credentials because the first selected method is specific. A generic credential table would be more flexible, but it would also invite agents to mix unrelated password, provider, custom ID, service, token, and session semantics before those families are ratified.

At most one active device credential per player is the conservative first posture because multi-device linking and account recovery are not yet ratified. Historical records remain necessary for rotation and revocation auditability.

Nakama remains useful as proof that device-style login is a common game-backend capability, but vibit rejects raw public device identifiers as credential proof. Pitaya remains useful for session and handler vocabulary, but credential records must not become transport sessions or handler context state.

## Agent Reasoning Summary

The correct next step is to ratify credential record semantics without creating storage. This gives later agents enough structure to write a migration safely, while keeping runtime authentication blocked until token verifier schema, migration planning, repository boundaries, adapters, tests, and implementation gates exist.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_boundary_clarity: high
  schema_first_discipline: high
  player_account_separation: high
  reversibility: medium
  implementation_deferral: high
  game_backend_reference_alignment: medium
  future_multi_device_flexibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- The credential record schema boundary is ratified with no schema added.
- Future migration work has a stable logical table target.
- Player account lifecycle tables remain credential-free.
- Token verifier schema remains the next active schema ratification work.
- Credential migration sources, repository interfaces, PostgreSQL adapters, runtime lookup, handlers, generated output, Protobuf changes, WebSocket changes, authentication dependencies, and authentication behavior remain deferred.

## Reversal Conditions

Revisit this decision if:

- A security review requires token verifier schema to be ratified before credential schema details become final.
- The maintainer ratifies multi-device account linking in the first authentication implementation.
- A future recovery or merge design requires multiple active device credentials per player.
- A future compatibility ADR explicitly adopts a Nakama-like authentication API surface.
- A future distributed runtime decision changes the durable credential target away from PostgreSQL.

## Follow-Up

- Define token verifier record schema boundary in W-0075.
- Plan the authentication schema migration queue after credential and token verifier boundaries are both ratified.
- Keep runtime authentication implementation blocked until later bounded work explicitly authorizes it.
