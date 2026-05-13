# Conversation: Work Continuation Semantics

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-add-work-item-continuation-system/`

Related artifacts:

- `.arch/work-items.yaml`
- `docs/workflow.md`
- `docs/workflow.zh-CN.md`
- `tools/vibit`

## Context

The maintainer asked how the project currently decides where work starts and stops, whether ADR numbers are work numbers, and whether a work-step system should exist.

The agent clarified that ADRs are decision identifiers, change specs are execution records, commits are repository snapshots, versions are release identifiers, and the project lacked a dedicated work-step queue.

## Maintainer Narrative

The maintainer then stated:

```text
建立这个机制，我希望我说继续就是继续推进一步，我说继续推进十步，你就会继续推进十步。
```

## Agent Response Summary

The agent interpreted this as a requirement to create a lightweight continuation mechanism where one step equals one work item, and multi-step continuation means advancing multiple ready work items in sequence.

The agent proposed making the mechanism machine-readable and checkable instead of relying on memory or informal judgment.

## Decisions

- A `Work Item` is the unit of one continuation step.
- A `Milestone` groups work items into a broader stage.
- `continue` / `继续` means advancing one `next_ready` work item unless blocked or confirmation is required.
- `continue N steps` / `继续推进 N 步` means advancing up to N next-ready work items in dependency order.
- Agents must stop early at blockers, verification failures, ask-first boundaries, or maintainer redirection.
- ADR IDs remain architecture decision IDs and must not become work-step IDs.

## Artifacts

- Added `.arch/work-items.yaml`.
- Added `docs/workflow.md`.
- Added `docs/workflow.zh-CN.md`.
- Added CLI support for `node tools/vibit check work`.
- Added CLI support for `node tools/vibit inspect work`.

## Open Questions

- Whether future work-item state transitions should be automated by dedicated CLI update commands remains deferred.
- Whether release version planning should be added to the roadmap remains deferred until the first runnable runtime proof is closer.

## Follow-Up

- Use `node tools/vibit inspect work` before interpreting future continuation requests.
- Keep work item state current after each completed continuation step.

## Redaction Notes

No secrets were included in this conversation log.
