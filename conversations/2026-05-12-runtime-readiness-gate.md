# Conversation: Runtime Readiness Gate

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-runtime-readiness-gate/`

Related artifacts:

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `decisions/ADR-0002-productive-bootstrapping.md`

## Context

After adding bootstrapping governance, the maintainer clarified that preventing meta-tooling drift should not mean rushing into a minimal implementation before the project is ready.

## Maintainer Narrative

The maintainer said:

> 不需要过于着急地跑通第一个最小实例，我们还是要做必要的准备，甚至应该做好万全的准备。因为我们要开发的东西本身针对性较强，所以提前想好并且提前做好准备是可以做到的。

## Agent Response Summary

The agent updated the governance rule to distinguish runtime readiness from premature implementation. The updated standard says the project should settle the minimum architecture needed for the first slice before implementation, while still avoiding readiness work that no longer changes how the slice will be built, verified, or maintained.

## Decisions

- Runtime readiness is a valid stage before the first backend vertical slice.
- Necessary preparation includes language, package layout, instance model, first module boundary, contract format, generated/handwritten boundary, test strategy, verification commands, and persistence assumptions.
- Preparation becomes drift only when it stops reducing ambiguity for the first slice.

## Artifacts

- Updated `CONSTITUTION.md`.
- Updated `CONSTITUTION.zh-CN.md`.
- Updated `AGENTS.md`.
- Updated `AGENTS.zh-CN.md`.
- Updated `decisions/ADR-0002-productive-bootstrapping.md`.
- Created change spec `changes/2026-05-12-add-runtime-readiness-gate/`.

## Open Questions

- Which implementation language should be selected?
- Which contract format should be selected?
- Which persistence assumption is acceptable for the first runtime slice?

## Follow-Up

- Prepare a runtime readiness decision set before the first backend vertical slice.

## Redaction Notes

No secrets, tokens, account details, or private data were recorded in this log.
