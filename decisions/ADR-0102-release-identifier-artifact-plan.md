# ADR-0102: Release Identifier And Artifact Plan

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-define-release-identifier-artifact-plan/`

Related conversations:

- `conversations/2026-05-21-release-identifier-artifact-plan.md`

Related artifacts:

- `docs/release-identifier-artifact-plan.md`
- `docs/release-identifier-artifact-plan.zh-CN.md`
- `docs/release-execution-maintainer-decision.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0101` recorded the maintainer's go decision as permission to continue into release identifier and artifact planning. It intentionally did not approve release execution commands, a final tag, artifacts, hosted deployment, or GitHub release records.

`W-0194` needed to define a concrete proposed identifier and the smallest release surface before any release-producing command can be considered.

## Decision

Record `v0.1.0-alpha.1` as the proposed release identifier.

Record the smallest future artifact/publication surface as:

- Git tag `v0.1.0-alpha.1`;
- GitHub release record for `v0.1.0-alpha.1`;
- release notes derived from repository-owned facts;
- source archive produced by the hosting platform as part of the tag or release workflow.

This decision does not authorize creating the tag, pushing the tag, creating the GitHub release record, creating source archives, publishing `v0.1 alpha`, or running release execution commands.

The repository check rule for this decision is `runtime.release_identifier_artifact_plan`.

## Conflict Check Record

The following checks were run on 2026-05-21 without creating or pushing anything:

```bash
git tag --list 'v0.1*'
git ls-remote --tags origin 'refs/tags/v0.1*'
curl -fsS https://api.github.com/repos/iceiko/vibit/releases/tags/v0.1.0-alpha.1
```

Observed result:

```yaml
local_v0_1_tags_found: false
remote_origin_v0_1_tags_found: false
github_release_record_for_v0_1_0_alpha_1_found: false
github_release_record_lookup_status: 404_not_found
```

The conflict check must be repeated before any later execution command.

## Decision Record

```yaml
release_identifier_artifact_plan: defined
completed_work_item: W-0194
decision: ADR-0102
check_rule: runtime.release_identifier_artifact_plan
proposed_release_identifier: v0.1.0-alpha.1
release_identifier_conflict_checked: true
release_identifier_conflict_found: false
release_identifier_selected_for_execution: false
release_declared: false
release_publishing_authorized_by_this_plan: false
release_execution_commands_authorized_by_this_plan: false
release_tag_authorized_by_this_plan: false
release_artifacts_authorized_by_this_plan: false
release_packaging_authorized_by_this_plan: false
hosted_deployment_authorized_by_this_plan: false
github_release_authorized_by_this_plan: false
release_tag_created: false
release_binary_created: false
release_archive_created: false
release_container_created: false
release_package_created: false
release_checksum_created: false
release_provenance_created: false
next_direction: release_execution_final_authorization
```

## Alternatives Considered

- Use `v0.1-alpha`.
- Use `v0.1.0-alpha`.
- Use `v0.1.0-alpha.1`.
- Publish immediately after recording the identifier plan.
- Include binaries, containers, checksums, provenance, packages, or hosted deployment in the first alpha surface.
- Keep the queue next-ready after the plan instead of blocking for final authorization.

## Rationale

`v0.1.0-alpha.1` is explicit, prerelease-scoped, and compatible with common SemVer-style ordering. It communicates first alpha posture without implying production readiness.

Keeping the release surface source-first reduces packaging and operations burden. Blocking the next work item at final authorization preserves the ask-first boundary before public release state is created.

Nakama and Pitaya remain reference baselines for keeping release identity, packaging, deployment, and runtime capability concerns separated. This decision follows that discipline without adopting their APIs, packaging models, deployment models, SDK surfaces, operations surfaces, or direct compatibility.

## Agent Reasoning Summary

The maintainer asked to continue while `W-0194` was next-ready. The agent treated this as permission to complete the planning item only, because the active ask-first boundary still forbids publishing, tag creation, artifact creation, GitHub release records, hosted deployment, and release execution commands. The agent selected a conservative SemVer-style alpha identifier, verified that no matching local tag, remote tag, or GitHub release record was visible, and blocked the next work item for final authorization.

## Decision Weights

```yaml
decision_weights:
  alpha_posture_clarity: high
  release_scope_restraint: high
  source_inspectability: high
  artifact_boundary_clarity: high
  public_release_risk: low
  runtime_behavior_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0194` is completed as a planning-only release identifier and artifact plan.
- `W-0195 Confirm release execution final authorization` becomes blocked.
- The repository remains pre-alpha until a later authorized execution step changes public release state.
- No tag, GitHub release record, release artifact, hosted deployment, publication command, runtime behavior, protocol route, generated output, migration, dependency, operations/admin behavior, authentication/session behavior change, direct compatibility, or broad product module is created by this decision.

## Reversal Conditions

Revisit this decision if `v0.1.0-alpha.1` is found to conflict with an existing tag, release record, package, image, hosted deployment, or external artifact; if the maintainer chooses a different identifier; or if the first release surface should include more than source-first release records.

## Follow-Up

- Complete `W-0195 Confirm release execution final authorization` before any release-producing command runs.
- Repeat tag and release-record conflict checks immediately before any execution step.
