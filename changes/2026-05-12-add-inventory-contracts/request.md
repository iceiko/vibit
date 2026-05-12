# Request

## Original Request

The maintainer asked to continue after the inventory proof slice was prepared at the module-manifest level.

## Clarified Requirement

Declare the first inventory runtime contracts as source artifacts before runtime implementation or generation begins.

## User-Visible Outcome

Future agents should be able to inspect the contract registry and contract files to understand the exact payloads, permissions, errors, events, and generated outputs expected by the first inventory proof slice.

## Non-Goals

- Do not implement runtime TypeScript code yet.
- Do not generate files from the contracts yet.
- Do not add package dependencies.
- Do not add a full contract validator yet.
- Do not define contracts for modules other than `inventory`.

## Unknowns

- Exact TypeScript generator output paths.
- Exact runtime validation library.
- Exact package manager and test runner.
- Whether the first persistence adapter is in-memory only or file-backed.

## Acceptance Criteria

- Add a durable decision for the first contract source format.
- Add a machine-readable contract registry.
- Add source contract files for `GrantItem`, `GetInventory`, `ItemGranted`, inventory errors, and inventory permissions.
- Update inventory module metadata to link the contract files.
- Run verification and record the results.
