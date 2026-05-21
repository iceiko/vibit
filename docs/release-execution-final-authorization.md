# Release Execution Final Authorization

Status: Accepted v0.1
Last updated: 2026-05-21
Scope: Final maintainer authorization for the source-first `v0.1.0-alpha.1` release execution
Depends on: `docs/release-identifier-artifact-plan.md`, `docs/release-execution-authorization-gate.md`
Canonical decision: `ADR-0103`

The paired Simplified Chinese translation is `docs/release-execution-final-authorization.zh-CN.md`. The English file is authoritative.

This document records the final `W-0195` maintainer authorization before release-producing commands run. It authorizes exactly one source-first alpha release surface: the `v0.1.0-alpha.1` Git tag, GitHub release record, release notes, and the hosting-platform source archive generated from that tag. It does not authorize binaries, packages, containers, checksum files, provenance or signing artifacts, hosted deployments, install scripts, registry publication, public announcements beyond the GitHub release record, runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependency changes, broad operations/admin behavior, authentication/session behavior changes, broad product module expansion, or direct Nakama/Pitaya API compatibility.

## 1. Maintainer Decision

The maintainer explicitly authorized final release execution with this decision:

```text
授权执行 W-0195：go，release=v0.1.0-alpha.1，允许创建并推送 Git tag，允许创建 GitHub Release，仅发布 GitHub source archive，不发布二进制/包/容器/checksum/签名/部署/公告；若发现版本冲突或验证失败则立即停止
```

The maintainer also requested that the README become more attractive for developers to try because the project now needs to find users.

## 2. Core Rule

The final authorization record is:

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
next_direction: first_alpha_user_discovery
```

## 3. Authorized Commands

The following command families are authorized only after the conflict checks and verification commands in this document pass on the target commit:

```bash
git tag -a v0.1.0-alpha.1 -m "v0.1.0-alpha.1"
git push origin main
git push origin v0.1.0-alpha.1
gh release create v0.1.0-alpha.1 --title "v0.1.0-alpha.1" --notes-file <release-notes-file>
```

An equivalent GitHub API release creation command is allowed if `gh` is unavailable. The GitHub token may come from the ignored local environment file, but token values must not be printed, committed, or included in release notes.

## 4. Authorized Release Surface

The authorized public surface is intentionally small:

- Git tag `v0.1.0-alpha.1`;
- GitHub release record `v0.1.0-alpha.1`;
- release notes derived from repository-owned facts;
- GitHub's automatically generated source archive for the tag.

The README refresh is authorized because it is documentation and developer-experience work needed for the first alpha to be understandable by external users. It must stay factual about the alpha state and must not imply production readiness, package availability, hosted deployment availability, SDK availability, direct Nakama/Pitaya API compatibility, or unsupported binary/container distribution.

## 5. Deferred Surface

The following remain deferred:

- standalone binaries;
- packages;
- container images;
- checksum files;
- provenance files;
- signing artifacts;
- hosted deployments;
- install scripts;
- package or container registry publication;
- public announcements beyond the GitHub release record;
- runtime behavior changes;
- protocol route changes;
- Protobuf source or generated output changes;
- migrations;
- dependencies;
- broad operations/admin behavior;
- authentication/session behavior changes;
- broad product module expansion;
- direct Nakama/Pitaya API compatibility.

## 6. Verification Requirements

Before creating the tag or GitHub release, the repository must pass:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change confirm-release-execution-final-authorization --json
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

The GitHub release-record lookup for `v0.1.0-alpha.1` must also return not found before the release record is created.

## 7. Known Warning Handling

The known repository warning remains accepted for this release:

```text
runtime.identity_boundary on runtime/internal/platform/persistence/postgres/authentication_repository.go
```

Any new warning, failed check, version conflict, or target-commit change after verification stops release execution until reverified.

## 8. Stop Conditions

Stop before release execution if any of these occur:

- required checks fail;
- a new warning appears and is not explicitly triaged;
- `v0.1.0-alpha.1` already exists as a local tag, remote tag, GitHub release record, package, image, hosted deployment, or external artifact;
- the target commit changes after verification without rerunning verification;
- release notes or README content would include unredacted credentials, tokens, verifier keys, digests, DSNs, transport proof carriers, GitHub tokens, or private environment file content;
- a binary, package, container, checksum, provenance file, signing artifact, hosted deployment, install script, registry publication, or public announcement beyond the GitHub release record is attempted;
- runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior, broad product modules, or direct Nakama/Pitaya API compatibility changes inside this release execution slice.

## 9. Next Work

After release execution, the next bounded direction is:

```text
W-0196 Define first alpha user discovery loop
```

That future work should turn the README's improved invitation into a concrete user-discovery loop without changing runtime behavior or expanding the release artifact surface.

## 10. Reference Alignment

Nakama and Pitaya both make public release state visible while keeping runtime capability work separate from packaging and deployment posture. This authorization follows that discipline for vibit's first source-first alpha and keeps direct Nakama/Pitaya API compatibility out of scope.
