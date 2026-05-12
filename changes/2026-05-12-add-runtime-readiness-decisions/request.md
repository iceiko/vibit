# Request

## Original Request

The maintainer asked to continue, after clarifying that the project should not rush into the first minimal runtime instance and should make necessary preparations because vibit is a specialized architecture.

Relevant maintainer statement:

```text
不需要过于着急地跑通第一个最小实例，我们还是要做必要的准备，甚至应该做好万全的准备。因为我们要开发的东西本身针对性较强，所以提前想好并且提前做好准备是可以做到的。
```

## Clarified Requirement

Create a runtime readiness decision package that settles the minimum architecture needed before the first backend vertical slice, without starting runtime implementation code.

## User-Visible Outcome

Future agents can inspect explicit decisions for:

- The first reference implementation language.
- The first server instance model.
- The contract and generated-file boundary.
- The first runtime proof slice.

## Non-Goals

- Do not implement the server runtime yet.
- Do not add package dependencies.
- Do not choose every final framework technology.
- Do not build the first inventory capability in this change.

## Unknowns

- Exact package manager.
- Exact TypeScript test runner.
- Exact schema source format for the first command, query, and event.
- Initial persistence adapter.

## Acceptance Criteria

- Add durable Agent Decision Records for runtime readiness decisions.
- Add a machine-readable runtime architecture manifest.
- Link the decisions from the relevant change and conversation records.
- Update only necessary project guides and README references.
- Run repository verification and record results.
