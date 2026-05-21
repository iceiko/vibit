# Release Identifier And Artifact Plan

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Planning-only release identifier and artifact/publication boundary for the v0.1 alpha path
Depends on: `docs/release-execution-maintainer-decision.md`, `docs/release-execution-authorization-gate.md`
Canonical decision: `ADR-0102`

The paired Simplified Chinese translation is `docs/release-identifier-artifact-plan.zh-CN.md`. The English file is authoritative.

This document defines the release identifier and artifact plan for `W-0194`. It does not publish `v0.1 alpha`, create or push release tags, create binaries, archives, containers, packages, checksums, provenance files, hosted deployments, GitHub release records, runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, broad product modules, or direct Nakama/Pitaya API compatibility.

## 1. Purpose

`W-0193` recorded a maintainer go decision to continue release execution planning. The remaining gap is to choose a concrete proposed release identifier and identify the smallest artifact/publication surface before any release-producing command is considered.

This plan records the proposal and stop conditions. It is a planning artifact, not a release execution artifact.

## 2. Core Rule

The release identifier and artifact plan is:

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
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
direct_nakama_pitaya_api_compatibility_added: false
next_direction: release_execution_final_authorization
```

The proposal uses `v0.1.0-alpha.1` because it is explicit SemVer-style prerelease text, carries alpha posture, and avoids implying production readiness. It is not a created tag, not a published release, and not final execution authorization by itself.

## 3. Conflict Check

The proposed identifier conflict check was run on 2026-05-21 without creating or pushing anything:

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

This means no conflict was observed for `v0.1.0-alpha.1` at planning time. A later execution step must re-run conflict checks immediately before any tag or release record creation.

## 4. Proposed Artifact Surface

The smallest future public release surface is:

- Git tag: `v0.1.0-alpha.1`.
- GitHub release record for `v0.1.0-alpha.1`.
- Release notes/changelog content derived from repository-owned facts.
- Source archive produced by the hosting platform as part of the release record or tag workflow.

These are proposed for a later execution step only. This plan creates none of them.

## 5. Deferred Artifact Families

The following remain deferred:

- standalone release binaries,
- packages,
- container images,
- checksum files,
- provenance or signing artifacts,
- hosted deployments,
- generated SDK packages,
- package registry publication,
- binary install scripts,
- public announcement beyond a release record.

Deferring these keeps the first alpha inspectable and source-first while avoiding packaging and operations promises the repository has not ratified.

## 6. Future Command Boundary

A later execution record may authorize commands in this family only after maintainer approval:

```bash
git tag v0.1.0-alpha.1 <verified-commit>
git push origin v0.1.0-alpha.1
```

A later execution record may also authorize creating a GitHub release record through a maintainer-approved command or UI/API workflow.

This plan does not authorize running those commands. It does not choose the exact final commit, tag annotation/signing posture, release-note text, GitHub release command, or publication timestamp.

## 7. Verification Requirements

Before any later execution step, the repository should pass:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
git status --short --branch
git tag --list 'v0.1*'
git ls-remote --tags origin 'refs/tags/v0.1*'
```

The GitHub release-record conflict check should also be repeated for the exact final identifier.

## 8. Known Warning Handling

The known repository warning remains:

```text
runtime.identity_boundary on runtime/internal/platform/persistence/postgres/authentication_repository.go
```

This plan accepts that known warning for release planning only. Any new warning must be triaged before release execution can proceed.

## 9. Stop Conditions

Stop before release execution if any of these occur:

- required checks fail;
- a new warning appears and is not explicitly triaged;
- `v0.1.0-alpha.1` already exists as a local tag, remote tag, GitHub release record, package, image, hosted deployment, or external artifact;
- the target commit changes after verification without rerunning verification;
- release notes would include unredacted credentials, tokens, verifier keys, digests, DSNs, transport proof carriers, GitHub tokens, or private environment file content;
- a tag, GitHub release record, source archive, binary, package, container, checksum, provenance file, hosted deployment, or announcement is attempted without a later final authorization record;
- runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior, broad product modules, or direct Nakama/Pitaya API compatibility changes inside this planning slice.

## 10. Next Work

This plan intentionally blocks the queue at:

```text
W-0195 Confirm release execution final authorization
```

That future work item must record a final maintainer go/no-go decision before any release-producing command runs.

## 11. Reference Alignment

Nakama and Pitaya both separate release identity, packaging choices, deployment posture, and runtime capability work. This plan follows that discipline while preserving vibit's agent-native workflow and without adopting direct Nakama/Pitaya APIs, packaging models, deployment models, or compatibility promises.
