# Impact Analysis

## Affected Modules

- `authentication`
- `runtime`

## Module Ownership Impact

The authentication module keeps ownership of the storage-neutral `authentication.Repository` interface, record structs, and normalization helpers.

The PostgreSQL platform package owns the new adapter implementation, SQL shape, pgx error mapping, and fake-executor tests. No domain package imports pgx.

## Public Contract Impact

No public commands, queries, events, errors, permissions, Protobuf messages, WebSocket envelope behavior, or WebSocket handshake behavior changed.

The adapter maps PostgreSQL missing-row, duplicate, foreign-key, and check-constraint failures into stable adapter sentinel errors. Runtime client-facing error mapping remains deferred until runtime authentication handlers are ratified.

## Data And Migration Impact

No migration source changed. The adapter uses the already-ratified `authentication_device_credentials` and `authentication_access_tokens` schemas from:

- `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`
- `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`

## Test Impact

Added focused fake-executor tests for:

- credential insert SQL
- credential lookup SQL
- token insert SQL
- token lookup SQL
- credential revocation SQL
- token revocation SQL
- cleanup eligibility query SQL
- no transaction-control SQL
- mutation normalization and UTC timestamps
- nullable terminal-state timestamp row mapping
- missing-row behavior
- duplicate, foreign-key, and check-constraint error mapping
- no live PostgreSQL dependency by default
- `UnitOfWork.NewAuthenticationRepository`

## Documentation Impact

Updated English and Simplified Chinese persistence and agent guidance to say the adapter is implemented while runtime authentication, login, token generation, token validation, logout execution, refresh, cleanup jobs, WebSocket behavior, Protobuf behavior, generated authentication shapes, and authentication dependencies remain deferred.

## Compatibility Risks

Runtime client behavior is unchanged because no handler, validator, route, generated shape, or Protobuf message consumes the adapter yet.

The main implementation risk is future agents mistaking adapter availability for runtime product availability. The manifests and `tools/vibit check runtime` now keep runtime authentication behavior markers false while allowing the bounded authentication PostgreSQL adapter.
