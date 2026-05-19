# Impact

## Runtime

- Adds a standards gate for future first-message connection binding.
- Defines `runtime.authentication.BindConnection` as the future system route candidate.
- Defines `vibit.authentication.v1.BindConnectionRequest` and `vibit.authentication.v1.BindConnectionResponse` as future payload candidates.
- Keeps WebSocket transport credential-neutral.
- Keeps the existing Protobuf envelope unchanged.
- Adds no Go runtime behavior.

## Authentication

- Reuses opaque access-token validation as the future proof validation source.
- Keeps raw access-token proof in the future protocol/application handoff only.
- Keeps `RequestIdentity.SessionValidated` false until durable session persistence is separately implemented.

## Data

- Adds no session persistence schema.
- Adds no migration.
- Changes no repository interface.
- Adds no PostgreSQL adapter behavior.

## Reference Alignment

- Nakama guides the connection lifecycle posture: authentication precedes realtime socket features, while active sockets and token/session lifecycle are not treated as the same object.
- Pitaya guides the architecture posture: transport acceptors, sessions, binding, route handlers, groups, and cluster concerns remain separate.
