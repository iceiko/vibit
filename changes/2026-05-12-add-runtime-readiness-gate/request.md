# Request

## Original Request

The maintainer clarified that the project should not rush into a first minimal implementation:

> 不需要过于着急地跑通第一个最小实例，我们还是要做必要的准备，甚至应该做好万全的准备。因为我们要开发的东西本身针对性较强，所以提前想好并且提前做好准备是可以做到的。

## Clarified Requirement

Refine the bootstrapping governance rule so it does not push agents into premature runtime implementation. Add a runtime readiness gate that defines the minimum decisions needed before the first backend vertical slice.

## User-Visible Outcome

Future agents understand that the project should avoid both infinite meta-tooling and premature runtime coding. Necessary preparation is valid when it makes the first runtime slice more coherent and reduces foreseeable churn.

## Non-Goals

- Do not choose the runtime language in this change.
- Do not start the first runtime slice in this change.
- Do not add a new CLI check yet.
- Do not expand readiness into a large new process.

## Unknowns

- Which runtime language will be chosen.
- Which contract format will be chosen.
- Which persistence assumptions will be accepted for the first slice.

## Acceptance Criteria

- [x] Constitution distinguishes runtime readiness from rushing implementation.
- [x] AGENTS gives future agents an operational runtime readiness checklist.
- [x] ADR-0002 is updated to reflect the maintainer's clarification.
- [x] Simplified Chinese translations are updated.
- [x] Conversation log records the maintainer's clarification.
- [x] Verification is recorded.
