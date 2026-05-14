# Impact

## Affected Modules

- `player`: receives its first semantic public contract manifests for account lifecycle.

## Module Ownership Impact

No ownership changes. The player module already owns player identity and player account lifecycle. This change registers semantic contracts for that ownership.

## Public Contract Impact

Adds:

- Command: `CreatePlayerAccount`
- Query: `GetPlayerAccount`
- Event: `PlayerAccountCreated`
- Error catalog: `player_account_errors`
- Permission catalog: `player_account_permissions`

These are semantic contracts only. They are not implemented Go handlers, Protobuf messages, database schema, or WebSocket routes.

## Data And Migration Impact

No migration is added or changed. Account persistence schema remains deferred.

## Protocol Impact

No Protobuf source is added and the envelope is unchanged. The player module stops being boundary-only, so future protocol checks may require player Protobuf messages once generation is ratified. That follow-up remains separate.

## Test Impact

No Go tests are added because no Go runtime code changes.

Repository checks must verify contracts and protocol state.

## Compatibility Risks

Moderate. Registering semantic public contracts means future agents may try to implement them. The contract manifests and work queue explicitly block runtime implementation, persistence, authentication, token behavior, credential storage, Protobuf, and handshake changes until later ratified steps.
