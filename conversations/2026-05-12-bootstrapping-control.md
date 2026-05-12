# Conversation: Bootstrapping Control

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-bootstrapping-control/`

Related artifacts:

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `decisions/ADR-0002-productive-bootstrapping.md`

## Context

After several standards and tooling slices, the maintainer asked whether vibit might drift in the future and how to keep self-bootstrapping under control.

## Maintainer Narrative

The maintainer said:

> 比起现在是否已经跑偏，我更关注的是未来是否可能跑偏。我们要让整个项目处在一个自举可控的范围内，这样才能最终很好地达成我们的目标。

The maintainer also emphasized that the project should do what is actually needed, not add process for its own sake.

## Agent Response Summary

The agent pushed the completed local commits, then added a small governance rule and ADR. The rule treats self-bootstrapping as valid only when it directly supports concrete runtime progress, module boundaries, contracts, generated shapes, tests, verification, or agent context reduction.

## Decisions

- Add a constitutional principle for productive bootstrapping.
- Add an AGENTS operating rule for future agents.
- Record the durable decision in ADR-0002.
- Do not add a new CLI check yet; the rule is intentionally lightweight until runtime work validates where automation is needed.

## Artifacts

- Updated `CONSTITUTION.md`.
- Updated `CONSTITUTION.zh-CN.md`.
- Updated `AGENTS.md`.
- Updated `AGENTS.zh-CN.md`.
- Created `decisions/ADR-0002-productive-bootstrapping.md`.
- Created change spec `changes/2026-05-12-add-bootstrapping-control/`.

## Open Questions

- Should this governance rule become an automated check later?
- Which runtime language should be selected for the first vertical slice?
- Should the first vertical slice use the existing `inventory` module?

## Follow-Up

- Decide the first runtime language and minimal server instance architecture.
- Start the first small end-to-end backend capability.

## Redaction Notes

A GitHub token was provided in the conversation but intentionally not recorded here. No secrets, tokens, account details, or private data were committed.
