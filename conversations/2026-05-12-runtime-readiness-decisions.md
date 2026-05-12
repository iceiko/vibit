# Conversation: Runtime Readiness Decisions

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-runtime-readiness-decisions/`

Related artifacts:

- `.arch/runtime.yaml`
- `decisions/ADR-0003-first-reference-runtime-language.md`
- `decisions/ADR-0004-minimal-server-instance-model.md`
- `decisions/ADR-0005-contract-and-generation-boundary.md`
- `decisions/ADR-0006-first-runtime-proof-slice.md`

## Context

The project had added a runtime readiness gate after the maintainer clarified that vibit should not rush into a minimal runtime instance. The next step was to convert that gate into concrete but reversible decisions that reduce ambiguity for future agents.

## Maintainer Narrative

The maintainer emphasized careful preparation before runtime implementation:

```text
不需要过于着急地跑通第一个最小实例，我们还是要做必要的准备，甚至应该做好万全的准备。因为我们要开发的东西本身针对性较强，所以提前想好并且提前做好准备是可以做到的。
```

The maintainer then asked to continue:

```text
继续
```

## Agent Response Summary

The agent interpreted the request as permission to continue runtime readiness work without starting server implementation code.

The agent chose a bounded decision package:

- TypeScript on Node.js as the first reference implementation.
- A modular monolith single-process server as the first server instance model.
- Commands, queries, events, errors, and permissions as the first public contract categories.
- `inventory` as the preferred first runtime proof slice.

The agent also kept the broader vibit standard language-neutral and recorded reversal conditions so future evidence can supersede individual choices.

## Decisions

- `ADR-0003`: Use TypeScript on Node.js for the first reference runtime.
- `ADR-0004`: Use a modular monolith single-process server model first.
- `ADR-0005`: Separate source contracts, generated structure, and handwritten extension points.
- `ADR-0006`: Prefer `inventory` as the first runtime proof slice.

## Artifacts

- Added `.arch/runtime.yaml`.
- Added four Agent Decision Records.
- Added `changes/2026-05-12-add-runtime-readiness-decisions/`.

## Open Questions

- Exact package manager.
- Exact TypeScript test runner.
- Exact schema source format.
- Initial persistence adapter.

## Follow-Up

- Refine `modules/inventory/module.yaml` before first runtime implementation.
- Create a dedicated change spec for the first inventory runtime slice.
- Choose package manager, test runner, schema source, and persistence only when implementation needs them.

## Redaction Notes

No secrets, tokens, credentials, private account details, or unrelated private data are recorded in this log.
