# Impact Analysis

## Affected Modules

- `player`, because it owns player identity and account lifecycle.
- `runtime`, because PostgreSQL migration source, repository adapters, and runtime checks live under runtime boundaries.

## Module Ownership Impact

The player module now owns the planned PostgreSQL lifecycle tables:

- `player_accounts`
- `player_account_events`

This does not make the player module own authentication, credentials, tokens, runtime sessions, WebSocket connections, or inventory state.

## Public Contract Impact

No public command, query, event, error, or permission contract changes.

The schema boundary traces to existing player contracts:

- `CreatePlayerAccount`
- `GetPlayerAccount`
- `PlayerAccountCreated`

## Data And Migration Impact

No SQL migration source is added in this change.

The next migration source must use the ratified schema boundary.

## Test Impact

No Go runtime tests are added because no runtime behavior changes.

Static repository checks are updated so future migration files can be validated against the ratified schema boundary.

## Documentation Impact

The PostgreSQL persistence standard, module guide, runtime guide, manifests, ADRs, and change records are updated. The Simplified Chinese translation is updated with the same material meaning.

## Compatibility Risks

Low. This change defines schema expectations before data exists.

The main risk is over-constraining future authentication work. The schema intentionally excludes credential and session concerns so those can receive their own designs later.
