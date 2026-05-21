# ADR-0101: Release Execution Maintainer Decision

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-confirm-release-execution-maintainer-decision/`

Related conversations:

- `conversations/2026-05-21-release-execution-maintainer-decision.md`

Related artifacts:

- `docs/release-execution-maintainer-decision.md`
- `docs/release-execution-maintainer-decision.zh-CN.md`
- `docs/release-execution-authorization-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0100` defined the release execution authorization gate and intentionally blocked the queue at `W-0193 Confirm release execution maintainer decision`.

The maintainer asked what decision was required and then stated:

```text
我决策的结果为同意。
```

The maintainer did not provide a concrete release identifier, tag name, artifact family, publication surface, or release execution command set.

## Decision

Record the maintainer decision as `go_to_release_identifier_artifact_plan`.

This means:

- the release execution path may continue;
- the next bounded work item is `W-0194 Define release identifier and artifact plan`;
- the known `runtime.identity_boundary` warning is accepted for continuing release planning;
- no release execution command is authorized by this decision;
- no final release identifier is approved by this decision;
- no tag, artifact, GitHub release, hosted deployment, package, checksum, provenance file, binary, archive, container image, or announcement is authorized by this decision.

The repository check rule for this decision is `runtime.release_execution_maintainer_decision`.

## Decision Record

```yaml
maintainer_decision: go_to_release_identifier_artifact_plan
maintainer_approval_recorded: true
release_execution_path_may_continue: true
release_declared: false
release_publishing_authorized_by_this_decision: false
release_execution_commands_authorized_by_this_decision: false
release_identifier_approved_by_this_decision: false
release_tag_authorized_by_this_decision: false
release_artifacts_authorized_by_this_decision: false
release_packaging_authorized_by_this_decision: false
hosted_deployment_authorized_by_this_decision: false
github_release_authorized_by_this_decision: false
release_identifier_selected: false
release_tag_created: false
release_artifacts_created: false
```

## Alternatives Considered

- Treat `同意` as permission to publish `v0.1 alpha` immediately.
- Choose a release identifier on behalf of the maintainer.
- Create a tag, GitHub release, source archive, checksum, binary, package, container image, provenance file, hosted deployment, or announcement.
- Keep `W-0193` blocked despite the maintainer's explicit go/no-go answer.
- Interpret the decision as no-go.

## Rationale

The maintainer gave a clear go/no-go answer, but not the details required for safe release execution. Recording a narrow go decision preserves momentum without exceeding the stated authorization boundary.

This keeps release execution aligned with the repository's agent-native workflow: each step is durable, inspectable, and bounded before commands create public artifacts.

Nakama and Pitaya remain reference baselines for separating release posture, artifact choices, and deployment promises from runtime capability work. This decision follows that discipline without adopting their APIs, data models, release packaging, deployment models, SDK surfaces, operations surfaces, or direct compatibility.

## Agent Reasoning Summary

The agent asked for an explicit release execution decision because `W-0193` was blocked. The maintainer answered that the decision result is approval. Since no release identifier, artifact family, publication surface, or command boundary was specified, the agent recorded the approval as permission to continue to a release identifier and artifact plan, not as release execution permission.

## Decision Weights

```yaml
decision_weights:
  maintainer_intent_respected: high
  release_scope_restraint: high
  artifact_boundary_clarity: high
  public_release_risk: low
  runtime_behavior_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0193` is completed as a durable maintainer decision record.
- `W-0194 Define release identifier and artifact plan` becomes the next ready work item.
- The repository remains pre-alpha.
- No release identifier, release tag, release artifact, hosted deployment, GitHub release, package, checksum, provenance file, runtime behavior, protocol route, generated output, migration, dependency, broad operations/admin behavior, authentication/session behavior change, direct compatibility, or broad product module is created by this decision.

## Reversal Conditions

Revisit this decision if the maintainer intended immediate release publication, if a no-go decision should replace the go-to-plan decision, or if release identifier/artifact approval should be captured in this same work item rather than a separate bounded plan.

## Follow-Up

- Complete `W-0194 Define release identifier and artifact plan`.
- Keep release execution commands behind later explicit authorization.
