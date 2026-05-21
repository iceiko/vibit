# ADR-0098: Release Publishing Decision Gate

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-define-release-publishing-decision-gate/`

Related conversations:

- `conversations/2026-05-21-release-publishing-decision-gate.md`

Related artifacts:

- `docs/release-publishing-decision-gate.md`
- `docs/release-publishing-decision-gate.zh-CN.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-acceptance-checklist.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The v0.1 alpha path now has a packaged local developer flow, acceptance checklist, runtime runbook, request-loop script, local status endpoints, and a focused authenticated gameplay E2E proof.

The next risk is procedural: a contributor or agent could treat local readiness as permission to publish a release or create release artifacts.

`W-0190` was selected to define a release publishing decision gate without executing release publishing.

## Decision

Add `docs/release-publishing-decision-gate.md` and `docs/release-publishing-decision-gate.zh-CN.md`.

The gate defines:

- release-publishing prerequisites,
- release artifact boundaries,
- verification requirements,
- stop conditions,
- redaction expectations,
- and the next bounded contribution direction.

The repository check rule for this gate is `runtime.release_publishing_decision_gate`.

The gate explicitly keeps:

```yaml
release_declared: false
release_publishing_authorized_by_this_gate: false
release_execution_authorized_by_this_gate: false
release_packaging_authorized_by_this_gate: false
release_artifacts_created_by_this_gate: false
hosted_deployment_authorized_by_this_gate: false
```

The next bounded work item is `W-0191 Define release execution preparation gate`.

## Alternatives Considered

- Publish `v0.1 alpha` immediately.
- Create release tags, binaries, archives, packages, containers, or hosted deployments in the decision slice.
- Treat the alpha acceptance checklist as implicit release authorization.
- Skip a gate and move directly to broad product features.
- Leave release artifact boundaries undocumented.

## Rationale

A release decision gate is the smallest safe step after the packaged alpha developer flow. It allows maintainers and agents to reason about publishing prerequisites without creating artifacts or changing runtime behavior.

Nakama and Pitaya set expectations that server frameworks have clear release and support posture. vibit should meet that expectation through explicit project workflow before expanding product scope or executing publication steps.

## Agent Reasoning Summary

The maintainer asked to continue. `W-0190` was next ready. The implementation stayed documentation- and tooling-focused because this work item defines a release decision boundary rather than release execution.

## Decision Weights

```yaml
decision_weights:
  release_scope_restraint: high
  contributor_readiness_clarity: high
  artifact_boundary_clarity: high
  runtime_behavior_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- The repository has an explicit release publishing decision gate.
- Repository checks guard the gate and its deferrals.
- The project may proceed to a release execution preparation gate.
- The repository remains pre-alpha and no release artifact is created by this decision.

## Reversal Conditions

Revisit this decision if maintainers choose a different release process, if a hosting or package distribution target becomes a first-class alpha requirement, or if release execution needs to be split into smaller preparation, artifact, and publication gates.

## Follow-Up

- Advance `W-0191 Define release execution preparation gate`.
- Keep release publication, release tags, binaries, archives, containers, packages, hosted deployments, runtime changes, protocol changes, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct compatibility, and broad product expansion behind later explicit work items.
