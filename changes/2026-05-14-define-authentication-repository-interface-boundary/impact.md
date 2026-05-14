# Impact

## Architecture Impact

This change advances M-014 from migration static checks to a storage-neutral authentication repository interface boundary.

The new module manifest `modules/authentication/module.yaml` declares authentication as the owner of credential and token verifier repository interfaces while preserving player account lifecycle ownership in the player module.

## Runtime Impact

Go code is added only under `runtime/internal/modules/authentication/`.

The package defines records, mutation/query shapes, a repository interface, and normalization helpers. It does not add runtime authentication behavior, handlers, token generation, token validation, logout execution, cleanup jobs, WebSocket routes, Protobuf messages, generated authentication shapes, PostgreSQL adapters, or dependencies.

## Data Impact

No migration source changes and no data shape changes.

The existing migration sources remain:

- `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`
- `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`

## Test Impact

Focused Go tests cover repository interface shape, allowed status values, required-field normalization, digest copying, UTC timestamp normalization, and storage-neutral mutation/query validation.

## Documentation Impact

Adds English and Simplified Chinese module guides for authentication and updates architecture state so the next work item is the PostgreSQL adapter boundary.

## Compatibility Risks

Low. No public wire schema, public runtime handler, WebSocket behavior, or database migration changes.
