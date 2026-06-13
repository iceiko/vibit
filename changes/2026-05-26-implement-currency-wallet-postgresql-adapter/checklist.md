# Checklist: Implement Currency Wallet PostgreSQL Adapter

- [x] Write failing fake-executor adapter tests before implementation.
- [x] Implement `CurrencyWalletRepository`.
- [x] Add `NewCurrencyWalletRepositoryForUnitOfWork`.
- [x] Add `UnitOfWork.NewCurrencyWalletRepository`.
- [x] Map create/get/owner-get/balance-list/grant/spend/transaction-list SQL.
- [x] Scan rows through currency module normalizers.
- [x] Collapse PostgreSQL details into redacted currency repository errors.
- [x] Preserve caller-owned transaction control.
- [x] Keep tests independent of live PostgreSQL by default.
- [x] Add ADR and conversation artifacts.
- [x] Add change spec artifacts.
- [x] Update manifests, guides, rules, and `tools/vibit`.
- [x] Run required verification commands.
- [x] Record fresh verification evidence.
