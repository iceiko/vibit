# Request

## Original Request

The maintainer asked to continue development. The current next-ready work item is `W-0017: Add PostgreSQL configuration and transaction runner`.

## Clarified Requirement

Add explicit PostgreSQL runtime configuration and a pgx-backed transaction runner without wiring PostgreSQL into process startup yet.

The runner must implement the existing vibit-owned `tx.Runner` boundary, keep `pgx` imports inside the PostgreSQL platform owner package, and remain testable without a live PostgreSQL server.

## User-Visible Outcome

The repository has the platform plumbing needed for later persistent inventory composition:

- PostgreSQL connection input can be parsed from explicit config or environment values.
- A pgx-backed transaction runner can open, commit, and roll back an application-owned unit of work.
- Future composition can bind repositories to the transaction executor without exposing pgx to application or domain packages.

## Non-Goals

- Do not wire PostgreSQL into `cmd/vibit-server`.
- Do not make PostgreSQL mandatory for local startup.
- Do not add a new configuration framework dependency.
- Do not introduce transaction retry policy.
- Do not make live PostgreSQL integration tests mandatory.
- Do not change inventory command, query, event, permission, or Protobuf contracts.

## Unknowns

- The disposable PostgreSQL integration environment remains planned for a later work item.
- The persistent runtime composition path remains planned for a later work item.
- Transaction isolation and retry policy remain deferred until they have a concrete requirement.

## Acceptance Criteria

- Add a PostgreSQL config type and environment parser under `runtime/internal/platform/persistence/postgres/`.
- Add a pgx-backed transaction runner under `runtime/internal/platform/persistence/postgres/`.
- Keep `tx.Runner` and `tx.UnitOfWork` driver-neutral.
- Add focused tests for config parsing and transaction commit/rollback behavior without a live PostgreSQL server.
- Update architecture manifests and persistence guidance where needed.
- Run runtime and repository checks.
