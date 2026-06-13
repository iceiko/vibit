# Plan: Implement Currency Wallet PostgreSQL Adapter

1. Write failing fake-executor tests for the adapter constructor, SQL mapping, row scanning, error mapping, redaction, transaction-control absence, default live PostgreSQL deferral, and unit-of-work helper.
2. Implement `CurrencyWalletRepository` under `runtime/internal/platform/persistence/postgres`.
3. Add `UnitOfWork.NewCurrencyWalletRepository`.
4. Keep SQL mapping inside the PostgreSQL adapter package and convert domain value aliases to database argument primitives at the SQL boundary.
5. Map PostgreSQL and row-shape failures into redacted currency repository errors.
6. Add ADR, conversation, change spec, manifest, module guide, and check-rule updates.
7. Run focused Go tests, full Go tests, repository checks, diff checks, and token scan.
8. Record fresh verification evidence in `verification.md`.
