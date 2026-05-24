# Plan

1. Add the PostgreSQL `friend_relationships` migration source under `runtime/migrations/postgres`.
2. Record the migration-source decision in `ADR-0141` and the conversation log.
3. Update work-item and architecture manifests so `W-0233` is completed and `W-0234` becomes next-ready.
4. Add repository checks for the friends relationship migration source.
5. Update README, alpha, product maturity, roadmap, and agent guide pointers.
6. Run static checks, Go tests, and diff hygiene checks.
