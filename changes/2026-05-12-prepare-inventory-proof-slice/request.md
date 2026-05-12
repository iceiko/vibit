# Request

## Original Request

After runtime readiness decisions were committed and pushed, the maintainer asked to continue.

## Clarified Requirement

Prepare the `inventory` module as the first runtime proof slice by defining its module boundary, first public contracts, invariants, generated boundary, handwritten extension points, and required tests before runtime implementation begins.

## User-Visible Outcome

Future agents should be able to inspect `modules/inventory/module.yaml` and understand exactly what the first inventory runtime slice is supposed to prove.

## Non-Goals

- Do not implement TypeScript runtime code yet.
- Do not add package dependencies.
- Do not generate runtime files yet.
- Do not introduce persistence or migrations yet.
- Do not implement a full inventory system.

## Unknowns

- Exact schema source format for command, query, event, error, and permission contracts.
- Exact generated file paths once runtime packages exist.
- Exact test runner and runtime verification command.
- Whether initial persistence will be in-memory only or file-backed.

## Acceptance Criteria

- `modules/inventory/module.yaml` declares a concrete first proof-slice boundary.
- `modules/inventory/AGENTS.md` explains when to use the module, what is forbidden, and what must be tested.
- `modules/inventory/AGENTS.zh-CN.md` is updated with the same meaning.
- The change links to the runtime readiness decisions.
- Verification is run and recorded.
