# Impact

## Affected Modules

The `inventory` module remains the first proof slice.

This change creates the future runtime path for Inventory under `runtime/internal/modules/inventory/`, but it does not implement Inventory behavior.

## Module Ownership Impact

No module ownership moves.

The skeleton makes existing ownership rules concrete:

- Runtime startup belongs under `runtime/cmd/vibit-server/`.
- Application dispatch and composition belong under `runtime/internal/app/`.
- Platform adapters belong under `runtime/internal/platform/`.
- Handwritten domain runtime logic belongs under `runtime/internal/modules/<module>/`.
- Generated Go output belongs under `runtime/internal/generated/`.

## Public Contract Impact

No public contract changes.

The existing YAML contract sources remain the semantic source of truth.

## Event Impact

No event contract changes.

The skeleton does not implement event recording or publication.

## Permission Impact

No permission contract changes.

## Data And Migration Impact

No migrations are added.

The future SQL-first migration directory is created at `runtime/migrations/postgres/`.

## Test Impact

No Go tests are added because no Go source code exists yet.

`node tools/vibit check runtime` is updated so a Go module skeleton can pass without requiring `go test`. Once Go source files exist, runtime checks require Go test files and a working Go toolchain.

## Tooling Impact

`tools/vibit` now checks `runtime/go.mod` as the Go runtime module file. It no longer assumes a repository-root `go.mod`.

The rule catalog gains runtime-specific rule IDs for module-file validation, skeleton verification, and missing Go toolchain reporting.

## Documentation Impact

This change updates repository guides and README files and adds runtime/proto guide pairs.

## Compatibility Risks

The main risk is creating directories before implementation makes every package boundary fully necessary. The risk is acceptable because the directories directly follow ADR-0014 and carry no business logic or external dependencies.
