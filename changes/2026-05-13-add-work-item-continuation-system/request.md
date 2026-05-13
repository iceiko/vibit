# Request

## Original Request

```text
建立这个机制，我希望我说继续就是继续推进一步，我说继续推进十步，你就会继续推进十步。
```

## Clarified Requirement

Create a lightweight, agent-readable work continuation system.

The system must define:

- A `Work Item` as the unit of one continuation step.
- A `Milestone` as the larger stage that groups work items.
- A checked machine-readable queue that identifies current, next-ready, completed, blocked, and planned work.
- CLI inspection so agents can determine what `continue` means before starting work.
- Documentation that explains how maintainer continuation commands map to work item execution.

## User-Visible Outcome

When the maintainer says:

```text
继续
```

an agent should interpret it as:

```text
advance one next_ready work item
```

When the maintainer says:

```text
继续推进十步
```

an agent should interpret it as:

```text
advance up to ten next_ready work items in order, stopping early if blocked or if a confirmation boundary is reached
```

## Non-Goals

- Do not implement the next runtime feature slice in this change.
- Do not introduce a heavy project management system.
- Do not replace ADRs, change specs, commits, or release versions.
- Do not make work item numbering serve as architecture decision numbering.
- Do not add external tooling dependencies.

## Unknowns

- Whether future work-item state transitions should be fully automated by CLI commands remains deferred.
- Whether release version planning should be added to the roadmap remains deferred until the first runnable runtime proof is closer.

## Acceptance Criteria

- [ ] `.arch/work-items.yaml` defines milestones, work item statuses, continuation semantics, completed recent work, and next-ready work.
- [ ] English and Simplified Chinese workflow docs explain the continuation mechanism.
- [ ] `node tools/vibit check work` verifies the work-item queue.
- [ ] `node tools/vibit inspect work` returns machine-readable current/next work context.
- [ ] `node tools/vibit check all --json` includes the work check.
- [ ] Repository and agent guides point to the new mechanism.
- [ ] Verification is recorded.
