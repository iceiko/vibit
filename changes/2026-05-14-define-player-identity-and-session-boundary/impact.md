# Impact

## Architecture

This change adds a durable boundary standard and ADR for player identity and runtime session work. It keeps transport, protocol, application dispatch, player module ownership, and inventory ownership separate.

The standard aligns with the existing Nakama/Pitaya reference baseline:

- Nakama informs the account, authentication, user, and session capability surface.
- Pitaya informs the separation between acceptors, sessions, route handlers, and server framework routing.
- Neither project is copied as an API surface.

## Runtime

No runtime code changes are required in this step.

The existing `app.Session` and Protobuf `Session` metadata remain unchanged. They are documented as metadata until session validation exists.

## Protocol

No Protobuf envelope or WebSocket handshake change is introduced.

The existing envelope session fields are interpreted more precisely:

- `connection_id` is transport-local metadata.
- `session_id` is reserved for authenticated logical sessions.
- `player_id` is reserved for authenticated player identity.
- `connection_epoch` is reserved for reconnect/lifecycle behavior.

## Modules

`player` becomes the planned owner of player identity and player account lifecycle.

`inventory` remains the owner of inventory records, item quantities, item grants, and inventory permissions. It may reference `player_id`, but it does not own player accounts, authentication, session validation, or token formats.

## Data

No migrations are added.

No player account schema is defined.

No inventory ownership is moved.

## Compatibility

This is a standards and planning change. It does not break existing runtime behavior, protocol schemas, database schemas, or public contracts.

## Risks

The main remaining risk is that a future work item crosses into authentication implementation before the authentication model is ratified. The work queue now records ask-first boundaries for those decisions.
