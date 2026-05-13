# Verification

Verified:

- `node tools/vibit check work --json`
- `node tools/vibit inspect work --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check change plan-durable-inventory-runtime-work-queue --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go tests; this planning step does not change Go source.
- PostgreSQL live integration tests; this planning step does not run PostgreSQL.
