# Impact

## Runtime

- Selects a first-message connection binding implementation gate as the next work direction.
- Does not add runtime behavior.
- Does not add Protobuf source or generated output.
- Does not change the existing Protobuf envelope.
- Keeps WebSocket transport credential-neutral.

## Authentication

- Keeps request-level access-token validation as the current protected-route path.
- Keeps future connection binding tied to the existing access-token validation service boundary.
- Does not change authentication repositories, migrations, token lifecycle, logout, refresh, or cleanup behavior.

## Data

- Adds no session persistence schema.
- Adds no migration.
- Changes no repository interface.
- Adds no PostgreSQL adapter behavior.

## Reference Alignment

- Nakama guides the authenticated realtime socket lifecycle: a connection can become associated with authenticated state, but token/session and socket lifecycle remain distinct.
- Pitaya guides the implementation layering: transport acceptors stay separate from session-like binding state and route handlers.
