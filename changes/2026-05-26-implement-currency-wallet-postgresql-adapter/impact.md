# Impact: Implement Currency Wallet PostgreSQL Adapter

## Runtime

This change adds platform persistence code under `runtime/internal/platform/persistence/postgres` and one unit-of-work helper in `runner.go`.

It does not add application services, handlers, routes, startup composition, transport behavior, protocol behavior, generated output, or runtime wallet behavior.

## Currency Module

The storage-neutral currency module remains the owner of repository vocabulary, value types, normalizers, and typed repository errors.

The module does not import PostgreSQL or execute SQL.

## Data

No migration is added or changed.

The adapter maps to the existing migration source:

- `runtime/migrations/postgres/000008_create_currency_wallets.sql`
- `currency_wallets`
- `currency_wallet_balances`
- `currency_wallet_transactions`

## Compatibility

No public API, Protobuf envelope, Protobuf message, event, permission, generated client, migration, dependency, hosted deployment, SDK, release artifact, or direct compatibility surface changes in this slice.

## Verification

Default verification uses fake-executor tests and repository checks. Live PostgreSQL verification is not required by default for this implementation slice.
