# Plan: Implement Currency Wallet Runtime Behavior

1. Read the currency wallet runtime behavior gate, repository interface, PostgreSQL adapter, and comparable storage/friends application services.
2. Add focused failing tests under `runtime/internal/app/currency/service_test.go`.
3. Implement `runtime/internal/app/currency/service.go` with validated identity owner derivation, unit-of-work repository handoff, wallet ensure/get/list behavior, grant/spend behavior, transaction listing, and redacted public error mapping.
4. Add ADR-0208 and change/conversation artifacts.
5. Update manifests, guides, rule catalog, and `tools/vibit` to mark `W-0300` complete and open `W-0301`.
6. Run focused Go tests and repository checks.
7. Record verification evidence.
