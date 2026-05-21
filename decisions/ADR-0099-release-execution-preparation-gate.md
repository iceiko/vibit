# ADR-0099: Release Execution Preparation Gate

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-define-release-execution-preparation-gate/`

Related conversations:

- `conversations/2026-05-21-release-execution-preparation-gate.md`

Related artifacts:

- `docs/release-execution-preparation-gate.md`
- `docs/release-execution-preparation-gate.zh-CN.md`
- `docs/release-publishing-decision-gate.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/runtime-runbook.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0098` defined the release publishing decision gate without authorizing release execution or artifact creation. The next bounded step is to define what a future release execution plan must contain before any authorization review can happen.

`W-0191` was selected to define the release execution preparation gate without publishing a release or creating release artifacts.

## Decision

Add `docs/release-execution-preparation-gate.md` and `docs/release-execution-preparation-gate.zh-CN.md`.

The gate defines:

- release execution planning inputs,
- release-note input boundaries,
- artifact plan boundaries,
- maintainer approval points,
- verification requirements,
- rollback notes,
- stop conditions,
- redaction expectations,
- and the next bounded contribution direction.

The repository check rule for this gate is `runtime.release_execution_preparation_gate`.

The gate explicitly keeps:

```yaml
release_declared: false
release_publishing_authorized_by_this_gate: false
release_execution_authorized_by_this_gate: false
release_packaging_authorized_by_this_gate: false
release_artifacts_created_by_this_gate: false
hosted_deployment_authorized_by_this_gate: false
```

The next bounded work item is `W-0192 Define release execution authorization gate`.

## Alternatives Considered

- Execute release publication immediately.
- Create release tags, binaries, archives, packages, containers, checksums, or hosted deployments in the preparation slice.
- Write final release notes as a release artifact before authorization.
- Skip preparation and rely on maintainer memory.
- Move directly to broad product features.

## Rationale

Preparation is useful only if it separates planning from execution. The project can now define the inputs and approval points for a release run while keeping all artifact creation and publication behind later explicit authorization.

Nakama and Pitaya set expectations that server frameworks distinguish release readiness, artifacts, support posture, and deployment choices. vibit should keep that discipline without copying their APIs, package model, deployment model, or operations surfaces.

## Agent Reasoning Summary

The maintainer asked to continue. `W-0191` was next ready. The implementation stayed documentation- and tooling-focused because this work item prepares release execution boundaries rather than executing a release.

## Decision Weights

```yaml
decision_weights:
  release_scope_restraint: high
  execution_planning_clarity: high
  artifact_boundary_clarity: high
  maintainer_approval_clarity: high
  runtime_behavior_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- The repository has an explicit release execution preparation gate.
- Repository checks guard the preparation boundary and its deferrals.
- The project may proceed to a release execution authorization gate.
- The repository remains pre-alpha and no release artifact is created by this decision.

## Reversal Conditions

Revisit this decision if maintainers choose a different release process, if release notes or artifact planning need a separate document family, or if release execution needs to be split into smaller authorization, artifact, and publication gates.

## Follow-Up

- Advance `W-0192 Define release execution authorization gate`.
- Keep release publication, release tags, binaries, archives, containers, packages, checksums, provenance files, hosted deployments, runtime changes, protocol changes, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct compatibility, and broad product expansion behind later explicit work items.
