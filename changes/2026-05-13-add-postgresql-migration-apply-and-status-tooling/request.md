# Request

## Original Request

The maintainer asked to continue development for up to ten work items. The current next-ready work item is `W-0018: Add migration apply and status tooling`.

## Clarified Requirement

Add explicit PostgreSQL migration apply and status tooling through the runtime migration platform package.

The implementation must keep `github.com/pressly/goose/v3` imports inside `runtime/internal/platform/migrations/`, use SQL-first migration sources under `runtime/migrations/postgres/`, and avoid applying migrations automatically during normal server startup.

## User-Visible Outcome

Maintainers and future agents can invoke a small Go migration API to validate migration options, set the PostgreSQL dialect, report migration status, and apply SQL-first migrations against a caller-supplied `*sql.DB`.

## Non-Goals

- Do not apply migrations automatically in `cmd/vibit-server`.
- Do not replace SQL-first migrations.
- Do not add another migration tool.
- Do not weaken existing migration source validation.
- Do not require a live PostgreSQL server for default unit tests.
- Do not define the disposable PostgreSQL verification environment yet; that is `W-0019`.

## Unknowns

- The exact live integration invocation will be defined in `W-0019`.
- Process startup wiring for persistent inventory remains deferred until `W-0020`.

## Acceptance Criteria

- Add migration status and apply helpers under `runtime/internal/platform/migrations/`.
- Keep `goose` imports only inside `runtime/internal/platform/migrations/`.
- Require explicit migration directory and database handles from callers.
- Cover option validation and live-database gating behavior with focused tests that do not require PostgreSQL by default.
- Update manifests and migration guidance.
- Mark `W-0018` complete and move the next work item to `next_ready`.
