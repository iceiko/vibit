# Impact

## Affected Modules

The `inventory` module remains the first proof slice. Its contracts and ownership do not change.

This change defines where future inventory runtime code will live.

## Module Ownership Impact

No ownership moves between modules.

The runtime layout strengthens the existing ownership model:

- Domain module code belongs under `runtime/internal/modules/<module>/`.
- Platform dependencies belong under `runtime/internal/platform/`.
- Generated protocol code belongs under `runtime/internal/generated/`.
- `.proto` source files remain source artifacts under repository-root `proto/`.

## Public Contract Impact

No public contract changes.

The layout declares how contracts will later flow into generated Go and Protobuf shapes.

## Event Impact

No event contract changes.

The transaction boundary decision affects future event publication: command handlers should record state changes and emitted events inside an application transaction boundary.

## Permission Impact

No permission contract changes.

## Data And Migration Impact

No migrations are added.

The future migration source directory is declared as `runtime/migrations/postgres/`, with SQL-first migration files.

## Test Impact

No runtime tests are added because runtime code has not started.

Future runtime tests should live near the Go packages they verify and use Go standard-library `testing` first.

## Documentation Impact

This change updates runtime manifests, architecture README files, ADRs, repository README files, AGENTS guides, and conversation memory.

## Compatibility Risks

The main risk is choosing a package layout that becomes too restrictive. The risk is bounded by:

- Keeping the first layout inside one Go module under `runtime/`.
- Using `internal` to protect package boundaries.
- Recording reversal conditions in the ADR.
- Avoiding implementation code in this change.
