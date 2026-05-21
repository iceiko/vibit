# Release Identifier And Artifact Plan

Status: Draft v0.1
Last updated: 2026-05-21
Scope: v0.1 alpha 路径的 planning-only release identifier 与 artifact/publication boundary
Depends on: `docs/release-execution-maintainer-decision.md`, `docs/release-execution-authorization-gate.md`
Canonical decision: `ADR-0102`

配套英文源文档是 `docs/release-identifier-artifact-plan.md`。英文文件是权威版本。

本文档定义 `W-0194` 的 release identifier and artifact plan。它不发布 `v0.1 alpha`，不创建或推送 release tags，不创建 binaries、archives、containers、packages、checksums、provenance files、hosted deployments、GitHub release records、runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、broad product modules 或 direct Nakama/Pitaya API compatibility。

## 1. Purpose

`W-0193` 已经记录 maintainer go decision，允许继续 release execution planning。剩余缺口是在考虑任何 release-producing command 之前，选择一个具体 proposed release identifier，并识别最小 artifact/publication surface。

本 plan 记录 proposal 和 stop conditions。它是 planning artifact，不是 release execution artifact。

## 2. Core Rule

release identifier and artifact plan 是：

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

proposal 使用 `v0.1.0-alpha.1`，因为它是明确的 SemVer-style prerelease 文本，保留 alpha posture，并且不暗示 production readiness。它不是已创建 tag，不是已发布 release，也不是 execution authorization。

## 3. Conflict Check

proposed identifier conflict check 已在 2026-05-21 执行，未创建或推送任何内容：

```bash
git tag --list 'v0.1*'
git ls-remote --tags origin 'refs/tags/v0.1*'
curl -fsS https://api.github.com/repos/iceiko/vibit/releases/tags/v0.1.0-alpha.1
```

观察结果：

```yaml
local_v0_1_tags_found: false
remote_origin_v0_1_tags_found: false
github_release_record_for_v0_1_0_alpha_1_found: false
github_release_record_lookup_status: 404_not_found
```

这表示 planning 时没有观察到 `v0.1.0-alpha.1` 冲突。后续 execution step 必须在任何 tag 或 release record 创建前立即重新运行 conflict checks。

## 4. Proposed Artifact Surface

最小 future public release surface 是：

- Git tag: `v0.1.0-alpha.1`。
- `v0.1.0-alpha.1` 的 GitHub release record。
- 从 repository-owned facts 派生的 release notes/changelog content。
- 由 hosting platform 在 release record 或 tag workflow 中产生的 source archive。

这些只作为后续 execution step 的 proposal。本 plan 不创建它们。

## 5. Deferred Artifact Families

以下保持 deferred：

- standalone release binaries；
- packages；
- container images；
- checksum files；
- provenance 或 signing artifacts；
- hosted deployments；
- generated SDK packages；
- package registry publication；
- binary install scripts；
- release record 之外的 public announcement。

这些 deferral 让 first alpha 保持 source-first、可 inspect，同时避免尚未 ratified 的 packaging 和 operations promises。

## 6. Future Command Boundary

后续 execution record 只有在 maintainer approval 后，才可以授权这一类命令：

```bash
git tag v0.1.0-alpha.1 <verified-commit>
git push origin v0.1.0-alpha.1
```

后续 execution record 也可以通过 maintainer-approved command 或 UI/API workflow 授权创建 GitHub release record。

本 plan 不授权运行这些命令。它不选择 exact final commit、tag annotation/signing posture、release-note text、GitHub release command 或 publication timestamp。

## 7. Verification Requirements

任何后续 execution step 之前，repository 应通过：

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

GitHub release-record conflict check 也应针对 exact final identifier 重复执行。

## 8. Known Warning Handling

已知 repository warning 仍然是：

```text
runtime.identity_boundary on runtime/internal/platform/persistence/postgres/authentication_repository.go
```

本 plan 只为 release planning 接受该已知 warning。任何 new warning 都必须在 release execution 之前 triage。

## 9. Stop Conditions

出现以下情况时，在 release execution 前停止：

- required checks 失败；
- 出现 new warning 且未明确 triage；
- `v0.1.0-alpha.1` 已经作为 local tag、remote tag、GitHub release record、package、image、hosted deployment 或 external artifact 存在；
- target commit 在 verification 后变化但没有重新运行 verification；
- release notes 会包含未 redacted credentials、tokens、verifier keys、digests、DSNs、transport proof carriers、GitHub tokens 或 private environment file content；
- 在没有后续 final authorization record 的情况下尝试 tag、GitHub release record、source archive、binary、package、container、checksum、provenance file、hosted deployment 或 announcement；
- runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior、broad product modules 或 direct Nakama/Pitaya API compatibility 在本 planning slice 中变化。

## 10. Next Work

本 plan 有意将队列阻塞在：

```text
W-0195 Confirm release execution final authorization
```

该 future work item 必须在任何 release-producing command 运行前记录 final maintainer go/no-go decision。

## 11. Reference Alignment

Nakama 和 Pitaya 都分离 release identity、packaging choices、deployment posture 与 runtime capability work。本 plan 采用这种纪律，同时保留 vibit 的 agent-native workflow，并且不采用 direct Nakama/Pitaya APIs、packaging models、deployment models 或 compatibility promises。
