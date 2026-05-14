# ADR-0022: Player Account PostgreSQL Schema Boundary

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-confirm-next-direction-after-tooling-hardening/`
- `changes/2026-05-14-define-player-account-postgresql-persistence-schema-boundary/`

Related conversations:

- `conversations/2026-05-14-player-account-postgresql-persistence-direction.md`

Related artifacts:

- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/player/module.yaml`
- `modules/player/AGENTS.md`
- `modules/player/AGENTS.zh-CN.md`
- `decisions/ADR-0019-nakama-and-pitaya-reference-baseline.md`
- `decisions/ADR-0020-postgresql-persistence-boundary.md`
- `decisions/ADR-0021-player-identity-and-session-boundary.md`

## Context

The player account semantic contracts and Protobuf wire messages are ratified, but runtime handlers, authentication, token behavior, credential storage, session persistence, and player account persistence are still deferred.

After the agent tooling milestone, the maintainer authorized the professional recommendation to proceed with player account PostgreSQL persistence schema ratification. This is the smallest next step that advances durable player account work while preserving the existing identity/session boundary.

The schema decision must happen before migration source or repository code. Otherwise a later agent could hide login, credentials, tokens, session state, or external identity linking inside a table that looks like ordinary account persistence.

## Decision

Ratify the first player account PostgreSQL schema boundary as account lifecycle persistence only.

The player module owns these future PostgreSQL tables:

```text
player_accounts
player_account_events
```

`player_accounts` stores stable player account lifecycle state:

- `player_id`
- `display_name`
- `account_state`
- `created_at`
- `updated_at`
- `disabled_at`
- `deleted_at`

`player_account_events` stores durable player account lifecycle event records for the first account creation flow and future lifecycle transitions:

- `event_id`
- `event_type`
- `occurred_at`
- `player_id`
- `requested_by`
- `account_state`
- `display_name`
- `metadata`
- `recorded_at`

The first migration source must be SQL-first under `runtime/migrations/postgres/`, use goose markers, declare `-- Module: player`, and create only the ratified lifecycle tables and indexes.

The schema must not store credentials, password hashes, authentication providers, external identity links, tokens, refresh tokens, runtime sessions, WebSocket connection state, inventory state, or permissions grants.

## Alternatives Considered

- Add player account migration source immediately without a schema boundary standard.
- Store authentication credentials in `player_accounts`.
- Store token or session state next to player account rows.
- Add a generalized user identity/linking table before choosing login methods.
- Reuse inventory `player_id` rows as the account table.
- Copy Nakama account/auth tables directly.

## Rationale

The first player account schema should only persist the account lifecycle that has already been ratified semantically. This keeps the data model aligned with current contracts and avoids choosing authentication by accident.

Credentials, tokens, external identities, and sessions have different security, expiration, rotation, revocation, and compatibility requirements. They need their own design and verification path.

A durable event/audit table is useful even before broad event delivery exists because it gives agents a concrete place to preserve lifecycle facts such as `PlayerAccountCreated` without hiding event recording in handler logic.

Keeping the first schema small also makes the first migration easy to review, easy to roll back in development, and easy for future agents to validate with static checks.

## Agent Reasoning Summary

The next safe persistence step is to ratify tables for the already-known player account lifecycle, not to implement login. The schema must make exclusions explicit so future agents do not smuggle security or session behavior into account persistence.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: low
  reversibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- Player account migration source is allowed only after this schema boundary is present.
- Runtime repository implementation remains a later work item.
- Runtime player account handlers and WebSocket routes remain deferred.
- Authentication, token behavior, credential storage, external identity linking, session persistence, Protobuf envelope changes, and WebSocket handshake authentication remain deferred.
- Repository checks may allow player account migration files only when they match this ratified schema boundary.

## Reversal Conditions

Revisit this decision if the first authentication design requires external identity linking to be present in the same initial migration, if account deletion semantics require a different lifecycle model, or if event/outbox standards require lifecycle events to use a shared platform event table instead of a player-owned event table.

## Follow-Up

- Add the first SQL player account migration source.
- Add player module repository interfaces after migration source is ratified.
- Add a PostgreSQL player account repository adapter after the interface exists.
- Add runtime player account handlers and protocol bridge only after persistence and runtime implementation boundaries are ready.
