# Impact

## Affected Modules

- `player`: becomes the focus of the next contract ratification milestone, but still receives no runtime implementation, migrations, credentials, or authentication provider in this change.
- `runtime`: records that the next work will define account/session contracts while preserving the existing metadata-only request identity and session validator boundary.

## Module Ownership Impact

No ownership changes are made.

The selected direction preserves the existing split:

- `modules/player` owns player identity and player account lifecycle vocabulary.
- `runtime/internal/app` owns request identity context and session validation handoff.
- `runtime/internal/platform/transport/ws` owns WebSocket connection metadata only.
- `runtime/internal/platform/protocol/protobuf` owns envelope metadata conversion only.
- Domain modules such as `inventory` must not own player accounts, authentication, tokens, credentials, or runtime sessions.

## Public Contract Impact

No public command, query, event, permission, error, Protobuf, WebSocket, or database contract is added by this confirmation step.

The next milestone will ratify those contracts before implementation.

## Reference Impact

The selected direction explicitly continues the Nakama/Pitaya reference baseline:

- Nakama informs the account, user, authentication method, session token, refresh, and realtime socket capability vocabulary.
- Pitaya informs session binding, handler routing, frontend/backend server role vocabulary, and session state visibility.

This change does not make Nakama or Pitaya governing API shapes for vibit.

## Data And Migration Impact

No migration is added or changed.

Player account schema, credential storage, token storage, and session persistence remain deferred until ratified in a separate work item.

## Test Impact

No Go tests are required for this planning and workflow step.

Repository work checks should verify that one `next_ready` item exists after the confirmation gate closes.

## Compatibility Risks

Low. The change selects a conservative contract-first direction and stops before production authentication, persistence, protocol, or transport handshake choices.

The main risk is that future agents might interpret "player account/session contracts" as permission to implement authentication. The new work item keeps those implementation choices behind ask-first boundaries.
