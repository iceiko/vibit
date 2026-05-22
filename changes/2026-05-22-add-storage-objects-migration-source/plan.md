# Plan

1. Add the PostgreSQL `storage_objects` migration source under `runtime/migrations/postgres`.
2. Record the migration-source decision in `ADR-0111` and the conversation log.
3. Update work-item and architecture manifests so `W-0203` is completed and `W-0204` becomes next-ready.
4. Add repository checks for the storage objects migration source.
5. Update README, alpha, product maturity, and agent guide pointers.
6. Run static checks, Go tests, and diff hygiene checks.
