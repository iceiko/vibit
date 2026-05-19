# Impact

## Runtime

- Adds a standards gate for a future first-message connection binding implementation slice.
- Names future implementation artifacts for `BindConnection` protocol messages, generated output, Protobuf adapter behavior, application binding behavior, startup composition, and tests.
- Keeps the current runtime behavior unchanged.
- Keeps WebSocket transport credential-neutral.
- Keeps the existing Protobuf envelope unchanged.

## Authentication

- Future binding must validate opaque access-token proof through the existing application authentication service boundary.
- Future binding may normalize validated player identity for a server-observed connection id.
- This gate does not change authentication service behavior, repositories, token lifecycle, logout, refresh, or cleanup behavior.

## Data

- Adds no session table or persistent connection table.
- Adds no repository interface.
- Adds no PostgreSQL adapter.
- Adds no migration.

## Reference Alignment

- Nakama guides the need for authenticated realtime socket lifecycle and connection-associated identity.
- Pitaya guides the implementation split between acceptor/transport, session-like binding state, and route handlers.
- vibit keeps its own protocol/application boundaries and does not adopt direct Nakama or Pitaya public API compatibility.
