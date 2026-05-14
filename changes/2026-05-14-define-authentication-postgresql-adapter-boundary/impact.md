# Impact

## Architecture Impact

This change defines the future PostgreSQL adapter boundary for `authentication.Repository`.

The storage-neutral repository remains owned by `runtime/internal/modules/authentication/`. The future PostgreSQL adapter is reserved under `runtime/internal/platform/persistence/postgres/`.

## Runtime Impact

No Go adapter implementation is added.

No runtime authentication, credential lookup behavior, token issuance, token validation, logout, refresh, cleanup, handlers, routes, WebSocket behavior, Protobuf behavior, generated output, or authentication dependencies are added.

## Data Impact

No migration source changes and no data shape changes.

## Test Impact

No new Go adapter tests are added because adapter behavior is intentionally deferred.

The future adapter implementation must add focused fake-executor tests for credential and token SQL shape, no transaction-control SQL, normalization, row mapping, error mapping, and default no-live-PostgreSQL behavior.

## Documentation Impact

Updates English and Simplified Chinese PostgreSQL persistence guidance, authentication module guides, runtime guides, architecture manifests, and checks.

## Compatibility Risks

Low. This change adds no runtime behavior and no public wire or database changes.
