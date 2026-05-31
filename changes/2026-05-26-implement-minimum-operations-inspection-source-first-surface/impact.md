# Impact Analysis

## Affected Modules

- `runtime`
- `reference`
- `repository_workflow`
- `operations`

## Module Ownership Impact

`tools/vibit` owns the source-first inspection command. Runtime endpoint ownership is unchanged. No domain module owns broad operations inspection.

## Public Contract Impact

No command, query, event, error, permission, Protobuf, or WebSocket contract changes.

## Data And Migration Impact

No data schema, migration, repository interface, PostgreSQL adapter, or persistence behavior changes.

## Test Impact

Adds repository check coverage through `runtime.minimum_operations_inspection_source_first_surface_implementation`. No Go test is required because no Go runtime behavior changes.

## Documentation Impact

Updates source-first runbook and acceptance checklist references. Adds ADR, change artifacts, and conversation memory.

## Pitaya Alignment Impact

The inspection command records Pitaya architecture vocabulary as a deferred map:

- current WebSocket acceptor remains single-process;
- current first-message connection binding remains single-process;
- current route handler model remains application dispatch plus protocol bridge;
- frontend/backend roles, RPC/remotes, distributed groups/broadcast, and service discovery remain deferred.

## Compatibility Risks

No direct Nakama/Pitaya API compatibility is added. The command output is repository tooling metadata, not a runtime API compatibility promise.
