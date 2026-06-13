# Plan

1. Read the repository constitution, module instructions, currency wallet lifecycle/persistence/migration/repository boundary artifacts, and current next work item.
2. Add `runtime/internal/modules/currency` tests first for storage-neutral repository value types, lifecycle/transaction/actor/conflict vocabulary, normalizers, idempotency inputs, redacted errors, and forbidden fields.
3. Run the focused Go package test and confirm it fails before implementation.
4. Add `runtime/internal/modules/currency/repository.go` with storage-neutral value types, normalizers, redacted errors, and the repository interface.
5. Add `modules/currency` manifest and module AGENTS guides.
6. Add change spec, ADR, conversation memory, check rule, and architecture manifest updates.
7. Advance `.arch/work-items.yaml` from `W-0296` to `W-0297 Define currency wallet PostgreSQL adapter gate`.
8. Run focused and full verification.
