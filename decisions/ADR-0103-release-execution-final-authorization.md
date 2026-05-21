# ADR-0103: Release Execution Final Authorization

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-confirm-release-execution-final-authorization/`

Related conversations:

- `conversations/2026-05-21-release-execution-final-authorization.md`

Related artifacts:

- `docs/release-execution-final-authorization.md`
- `docs/release-execution-final-authorization.zh-CN.md`
- `docs/release-identifier-artifact-plan.md`
- `README.md`
- `README.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0102` recorded `v0.1.0-alpha.1` as the proposed release identifier and defined a source-first future release surface. The queue then blocked at `W-0195 Confirm release execution final authorization` because tag creation, tag push, GitHub release creation, and source archive publication create public release state.

The maintainer explicitly authorized release execution and asked that the README be made more attractive to developers so the project can begin finding users.

## Decision

Record the maintainer decision as `go`.

Authorize exactly:

- release identifier `v0.1.0-alpha.1`;
- creating and pushing the Git tag `v0.1.0-alpha.1`;
- creating the GitHub release record `v0.1.0-alpha.1`;
- release notes derived from repository-owned facts;
- GitHub's source archive generated from the tag;
- README messaging changes that make the first alpha easier and more compelling for developers to try.

Do not authorize binaries, packages, containers, checksum files, provenance files, signing artifacts, hosted deployments, install scripts, registry publication, public announcements beyond the GitHub release record, runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependency changes, broad operations/admin behavior, authentication/session behavior changes, broad product module expansion, or direct Nakama/Pitaya API compatibility.

The repository check rule for this decision is `runtime.release_execution_final_authorization`.

## Decision Record

```yaml
release_execution_final_authorization: recorded
completed_work_item: W-0195
decision: ADR-0103
check_rule: runtime.release_execution_final_authorization
maintainer_decision: go
authorized_release_identifier: v0.1.0-alpha.1
release_identifier_selected_for_execution: true
release_declared: true
release_publishing_authorized_by_this_decision: true
release_execution_commands_authorized_by_this_decision: true
release_tag_authorized_by_this_decision: true
github_release_authorized_by_this_decision: true
release_notes_authorized_by_this_decision: true
hosting_platform_source_archive_authorized_by_this_decision: true
release_artifacts_authorized_by_this_decision: source_archive_only
release_binary_authorized_by_this_decision: false
release_package_authorized_by_this_decision: false
release_container_authorized_by_this_decision: false
release_checksum_authorized_by_this_decision: false
release_provenance_authorized_by_this_decision: false
release_signing_authorized_by_this_decision: false
hosted_deployment_authorized_by_this_decision: false
install_script_authorized_by_this_decision: false
registry_publication_authorized_by_this_decision: false
public_announcement_authorized_by_this_decision: github_release_record_only
readme_user_acquisition_refresh_authorized: true
next_direction: first_alpha_user_discovery
```

## Alternatives Considered

- Treat the release as no-go.
- Keep the queue blocked and ask the maintainer for another authorization.
- Publish binaries, packages, containers, checksums, signatures, provenance, hosted deployments, or install scripts in the first alpha.
- Create a release announcement beyond the GitHub release record.
- Change runtime behavior or protocol surface while publishing the release.
- Keep the README as an internal status report.

## Rationale

The maintainer provided the exact release identifier, tag action, GitHub release action, artifact family, exclusions, and stop conditions required by the authorization gate. The authorized surface is small enough to execute safely while making the project externally discoverable.

Refreshing the README is part of release execution because the first source alpha needs a clear value proposition and a short try path. The README still preserves alpha honesty: vibit is not production-ready, does not ship binaries or containers in this release, and does not claim direct Nakama/Pitaya API compatibility.

Nakama and Pitaya remain reference baselines for capability class, not API compatibility. This release makes vibit's agent-native foundation visible without expanding runtime capability scope.

## Agent Reasoning Summary

The agent previously asked for final `go` or `no-go` because `W-0195` was blocked. The maintainer answered with a precise `go` authorization for `v0.1.0-alpha.1`, tag creation and push, GitHub release creation, source archive only, and stop-on-conflict or stop-on-verification-failure conditions. The agent recorded that decision, refreshed README entry points for users, and preserved all runtime and artifact deferrals outside the authorized source-first release surface.

## Decision Weights

```yaml
decision_weights:
  maintainer_intent_respected: high
  release_scope_clarity: high
  user_trial_path_clarity: high
  artifact_surface_restraint: high
  public_release_risk: medium
  runtime_behavior_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0195` is completed as the durable final release authorization record.
- `v0.1.0-alpha.1` is selected for release execution.
- The authorized release surface is tag plus GitHub release record/source archive only.
- The README is updated as the first external user entry point.
- The next bounded direction is first alpha user discovery.
- No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, hosted deployment, operations/admin expansion, authentication/session behavior change, broad product module, binary/package/container/checksum/provenance/signing artifact, or direct Nakama/Pitaya API compatibility is added by this decision.

## Reversal Conditions

Revisit this decision if `v0.1.0-alpha.1` conflicts with an existing public artifact, if verification fails, if the maintainer withdraws authorization, or if the first alpha must include an artifact family that this decision explicitly defers.

## Follow-Up

- Run the required conflict checks and verification commands.
- Create and push tag `v0.1.0-alpha.1` only after verification passes.
- Create the GitHub release record only after tag push succeeds.
- Continue with `W-0196 Define first alpha user discovery loop`.
