# ADR-0096: Alpha Acceptance Checklist

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-add-alpha-acceptance-checklist/`

Related conversations:

- `conversations/2026-05-21-alpha-acceptance-checklist.md`

Related artifacts:

- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `README.md`
- `README.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The local alpha path now has onboarding service behavior, protocol login, connection binding, protected inventory, protected presence, logout, a focused E2E proof, a runbook, a minimal request-loop script, and health/readiness/version/config troubleshooting endpoints.

The repository still needs a single acceptance checklist that makes the alpha readiness state reviewable without declaring or publishing a release.

## Decision

Add `docs/alpha-acceptance-checklist.md` and `docs/alpha-acceptance-checklist.zh-CN.md`.

The checklist classifies alpha readiness items as ready, manual, deferred, or blocked. It covers repository intake, prerequisites, local configuration, migration posture, runtime surfaces, authenticated gameplay flow, verification commands, redaction, and contribution entry points.

The repository check rule for this alpha acceptance checklist is `runtime.alpha_acceptance_checklist`.

The checklist explicitly does not authorize publishing `v0.1 alpha`, release packaging, runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependencies, broad operations/admin behavior, product module expansion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Publish `v0.1 alpha` immediately.
- Add release packaging before acceptance criteria exist.
- Keep alpha readiness scattered across README, runbook, and work queue text.
- Add a machine-only check without human-readable acceptance guidance.

## Rationale

A human-readable checklist is the smallest useful artifact before packaging the developer flow. It gives maintainers and contributors a stable review surface while keeping release publishing and packaging behind later explicit work.

Nakama and Pitaya set expectations that a serious server framework should have a clear developer entry point. vibit should meet that expectation through explicit, checkable acceptance criteria rather than copying external APIs or release mechanics prematurely.

## Agent Reasoning Summary

The maintainer asked to continue. `W-0188` was next ready. The implementation stayed documentation- and tooling-focused because the work item was an acceptance checklist, not a release or runtime behavior slice.

## Decision Weights

```yaml
decision_weights:
  alpha_readiness_clarity: high
  release_scope_restraint: high
  contributor_onboarding_value: high
  runtime_behavior_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- Maintainers and contributors can review alpha readiness from one checklist.
- Repository checks can guard the checklist and its deferrals.
- The next work can package the alpha developer flow without publishing a release.

## Reversal Conditions

Revisit this decision if alpha readiness becomes fully machine-checkable, if release packaging requires a different acceptance structure, or if the maintainer chooses to split acceptance criteria by audience.

## Follow-Up

- Package the alpha developer flow and document prerequisites without publishing a release.
