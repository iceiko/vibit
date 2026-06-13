# Change Request: Implement Currency Wallet PostgreSQL Adapter

## Request

Implement the currency wallet PostgreSQL adapter after the accepted adapter gate.

## Scope

Allowed:

- Add a PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`.
- Implement `runtime/internal/modules/currency.Repository`.
- Use the existing caller-supplied executor and unit-of-work pattern.
- Add `UnitOfWork.NewCurrencyWalletRepository`.
- Add focused fake-executor/query-capture tests.
- Add ADR, change, conversation, manifest, and check-rule updates for the implementation slice.

Forbidden in this change:

- Runtime wallet behavior, services, handlers, or route dispatch.
- Grant/spend execution above the adapter.
- Protocol routes.
- Protobuf source.
- Generated output.
- Dependencies.
- Migrations or schema changes.
- Startup wiring.
- Authentication/session behavior changes.
- Reward, inventory, purchase, catalog, event/audit, paid-currency, reservation, settlement, refund, or transfer behavior.
- SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- `CurrencyWalletRepository` implements `currency.Repository`.
- The adapter maps wallet create/get/owner-get, balance list, grant, spend, and transaction list to existing currency wallet tables.
- The adapter receives a caller-owned executor and does not issue transaction-control SQL.
- Row scanning returns normalized currency module records.
- PostgreSQL details are collapsed into redacted currency repository errors.
- Tests do not require live PostgreSQL by default.
- `runtime.currency_wallet_postgresql_adapter_implementation` verifies the implementation scope and deferrals.
- The next-ready work item becomes `W-0299 Define currency wallet runtime behavior gate`.
