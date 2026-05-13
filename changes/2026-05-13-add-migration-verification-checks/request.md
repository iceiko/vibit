# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0014` by adding a repository check for SQL-first PostgreSQL migration source files.

## User-Visible Outcome

Maintainers and agents can run:

```bash
node tools/vibit check migrations
```

The aggregate `node tools/vibit check all` command also runs the migration check.

## Non-Goals

- Do not add PostgreSQL migration apply or rollback execution.
- Do not add a PostgreSQL test environment standard.
- Do not add Go migrations.
- Do not implement the PostgreSQL inventory repository adapter.
- Do not change migration ownership or durable data ownership.

## Unknowns

- Disposable PostgreSQL test environment setup remains undecided.
- Future migration apply/rollback tooling remains deferred.

## Acceptance Criteria

- [x] Add a `check migrations` CLI command with JSON output support.
- [x] Include migration checks in `check all`.
- [x] Validate migration root presence, SQL naming, goose markers, absence of unapproved Go migrations, owning module traces, and first inventory table references.
- [x] Update rule catalog metadata for migration check rule IDs.
- [x] Update English and Simplified Chinese guidance.
