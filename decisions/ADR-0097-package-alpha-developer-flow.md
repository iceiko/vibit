# ADR-0097: Package Alpha Developer Flow

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-package-alpha-developer-flow/`

Related conversations:

- `conversations/2026-05-21-package-alpha-developer-flow.md`

Related artifacts:

- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The local alpha path has a README, runbook, redacted request-loop script, health/readiness/version/config endpoints, and alpha acceptance checklist. These artifacts were useful individually, but a contributor still needed to infer the intended sequence across several files.

`W-0189` was selected to package those existing entry points into one coherent local developer journey without publishing a release or adding runtime behavior.

## Decision

Add `docs/alpha-developer-flow.md` and `docs/alpha-developer-flow.zh-CN.md`.

The packaged flow connects:

- repository intake,
- prerequisites,
- static checks,
- Go tests,
- redacted request-loop proof,
- runtime status endpoints,
- PostgreSQL manual setup posture,
- redaction contract,
- acceptance checklist,
- and the next contribution entry point.

The repository check rule for this packaged alpha developer flow is `runtime.package_alpha_developer_flow`.

The packaged flow explicitly does not authorize publishing `v0.1 alpha`, release packaging, hosted deployment, runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, product module expansion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Publish `v0.1 alpha` immediately.
- Create release packaging before the local developer journey is coherent.
- Leave the flow scattered across README, runbook, request-loop script, and acceptance checklist.
- Add runtime code to create a more automated flow.

## Rationale

A packaged developer flow is the smallest useful step before any release-publishing decision. It improves contributor ergonomics while preserving the project's explicit release and runtime boundaries.

Nakama and Pitaya set expectations that serious server frameworks have coherent local entry points. vibit should meet that expectation with explicit documentation and checks before broadening product scope or publishing release artifacts.

## Agent Reasoning Summary

The maintainer asked to continue. `W-0189` was next ready. The implementation stayed documentation- and tooling-focused because the work item packages existing entry points rather than adding runtime behavior.

## Decision Weights

```yaml
decision_weights:
  contributor_path_clarity: high
  release_scope_restraint: high
  runtime_behavior_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- Contributors have one document for the current local alpha journey.
- Repository checks guard the packaging and its deferrals.
- The next work can be a release publishing decision gate without implying that publishing is already authorized.

## Reversal Conditions

Revisit this decision if a future release guide replaces the local developer flow, if publishing requirements require a different packaging artifact, or if the maintainer chooses to split contributor, operator, and release-publisher journeys.

## Follow-Up

- Define the release publishing decision gate without publishing a release unless the maintainer explicitly authorizes it.
