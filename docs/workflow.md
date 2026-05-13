# Work Continuation Workflow

Status: Draft v0.1
Last updated: 2026-05-13
Scope: Maintainer-agent continuation, work item sequencing, and roadmap state

This standard defines how vibit agents decide what "continue" means.

The paired Simplified Chinese translation is `docs/workflow.zh-CN.md`. The English file is authoritative.

## 1. Purpose

The project needs deterministic progress.

When the maintainer asks an agent to continue, the agent should not choose an arbitrary next task from memory. It should read the project work queue, find the next ready bounded work item, execute that item, verify it, and update the queue.

This mechanism keeps the project self-bootstrapping without letting the process become heavy.

## 2. Artifact Roles

The project uses different identifiers for different purposes:

- `M-000`: milestone identifier.
- `W-0000`: work item identifier.
- `ADR-0000`: architecture decision identifier.
- `changes/YYYY-MM-DD-change-id`: concrete change spec and execution record.
- Git commit hash: immutable repository snapshot.
- Version tag: future release identifier.

Do not use ADR numbers as work-step numbers. ADRs record durable decisions. Work items record execution steps.

Do not use release versions as work-step numbers. Versions describe published capability sets after work is complete.

## 3. Work Item Definition

A work item is the unit of one continuation step.

A good work item:

- Has a stable `W-0000` id.
- Belongs to exactly one milestone.
- Has a clear status.
- Has dependencies when order matters.
- Has completion criteria.
- Names its expected change spec when it is not yet implemented.
- Links to its change spec and commits after completion.
- Lists ask-first boundaries when a step could branch into a larger decision.

## 4. Milestone Definition

A milestone groups related work items under a stage goal.

Milestones are not releases. A milestone says what capability the repository is trying to prove next. A release says what users can depend on.

## 5. Status Values

Milestone status values:

```text
planned
active
completed
paused
superseded
```

Work item status values:

```text
planned
next_ready
active
blocked
completed
paused
superseded
```

Rules:

- At least one milestone should be `active`.
- A `next_ready` work item must have all dependencies completed.
- A completed work item should link to a change spec or explain why no change spec was needed.
- There should normally be one `next_ready` item per active milestone. Multiple next-ready items are allowed only when independent work can proceed in parallel.

## 6. Continuation Semantics

Maintainer phrase:

```text
continue
继续
```

Meaning:

```text
advance one next_ready work item
```

Maintainer phrase:

```text
continue N steps
继续推进 N 步
```

Meaning:

```text
advance up to N next_ready work items in dependency order
```

The agent must stop early if:

- No next-ready work item exists.
- The next work item is blocked.
- An ask-first boundary is reached.
- Verification fails.
- The maintainer redirects the work.

## 7. Required Agent Intake

Before interpreting a continuation request, an agent should run or inspect:

```bash
node tools/vibit inspect work
node tools/vibit check work
```

The agent should then read the relevant work item, its dependencies, related decisions, and expected change spec.

## 8. Execution Rule

For one work item, the preferred execution flow is:

```text
work item
-> change spec
-> implementation or documentation
-> tests/checks
-> verification record
-> work item status update
-> commit
```

If the work item changes architecture, public contracts, generated-file conventions, module ownership, dependency adoption, or other ask-first boundaries, the agent must ask the maintainer before implementation.

## 9. Multi-Step Execution

When asked to continue multiple steps, the agent should still finish each step as a complete bounded change before starting the next.

Each step should have:

- A change spec when non-trivial.
- Verification.
- Updated work-item state.
- A commit when the repository state is ready to save.

The agent may stop before the requested number of steps if the next step is not ready or needs maintainer confirmation.

## 10. Verification

Current commands:

```bash
node tools/vibit check work
node tools/vibit inspect work
node tools/vibit check all
```

`check work` verifies the presence and basic consistency of `.arch/work-items.yaml` and this workflow standard.

Future checks may validate dependency ordering, one-step state transitions, and whether completed work items link to existing commits.
