# Impact Analysis

## Affected Modules

- `runtime`: gains the third SQL migration source under the accepted migration root.
- `player`: remains referenced by `player_id` only; player account lifecycle tables are not modified.

## Module Ownership Impact

The new table belongs to `runtime.authentication`:

- `authentication_device_credentials`

The player module continues to own:

- `player_accounts`
- `player_account_events`

The credential schema references `player_accounts(player_id)` but does not add credential, token, external identity, session, WebSocket, request-validation, inventory, or permission data to player lifecycle tables.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf schema, generated output, or WebSocket protocol changes.

The migration stores fields already ratified by `docs/credential-record-schema-boundary.md`.

## Data And Migration Impact

Adds the third PostgreSQL migration source file:

- `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`

The `Up` migration creates the credential verifier table, constraints, indexes, and a table comment. The `Down` migration drops indexes and the table in reverse dependency order.

## Test Impact

No authentication repository interface or adapter exists yet, so there are no repository integration tests for credential behavior in this step.

Repository source checks validate migration naming, goose markers, module trace, credential migration shape, selected login/token boundary status, and runtime boundary preservation.

## Documentation Impact

Update architecture manifests, agent guides, and selected login/token standards so future agents understand that the credential migration source exists while token migration, repository interfaces, PostgreSQL adapters, runtime authentication, Protobuf, WebSocket, and authentication dependencies remain deferred.

## Compatibility Risks

Low. No runtime code consumes this table yet.

The main future risk is editing the migration after it is treated as applied in a shared environment. Future data changes should add new migrations instead.
