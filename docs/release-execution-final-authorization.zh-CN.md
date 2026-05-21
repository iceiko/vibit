# Release Execution Final Authorization

状态：Accepted v0.1
最后更新：2026-05-21
范围：`v0.1.0-alpha.1` source-first release execution 的最终 maintainer authorization
依赖：`docs/release-identifier-artifact-plan.md`、`docs/release-execution-authorization-gate.md`
权威决策：`ADR-0103`

英文原文是 `docs/release-execution-final-authorization.md`。英文文件是权威版本。

本文记录 `W-0195` 的最终 maintainer authorization，必须在任何 release-producing command 运行前存在。本授权只允许一个 source-first alpha release surface：`v0.1.0-alpha.1` Git tag、GitHub release record、release notes，以及 GitHub 基于 tag 自动生成的 source archive。它不授权 binaries、packages、containers、checksum files、provenance 或 signing artifacts、hosted deployments、install scripts、registry publication、GitHub release record 之外的 public announcements、runtime behavior changes、protocol route changes、Protobuf source 或 generated output changes、migrations、dependency changes、broad operations/admin behavior、authentication/session behavior changes、broad product module expansion，或 direct Nakama/Pitaya API compatibility。

## 1. Maintainer Decision

Maintainer 明确授权 final release execution：

```text
授权执行 W-0195：go，release=v0.1.0-alpha.1，允许创建并推送 Git tag，允许创建 GitHub Release，仅发布 GitHub source archive，不发布二进制/包/容器/checksum/签名/部署/公告；若发现版本冲突或验证失败则立即停止
```

Maintainer 同时要求 README 变得更吸引开发者尝试，因为项目现在需要找到用户。

## 2. Core Rule

最终授权记录是：

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

以下 command families 只有在本文冲突检查和验证命令通过后才被授权：

```bash
git tag -a v0.1.0-alpha.1 -m "v0.1.0-alpha.1"
git push origin main
git push origin v0.1.0-alpha.1
gh release create v0.1.0-alpha.1 --title "v0.1.0-alpha.1" --notes-file <release-notes-file>
```

如果 `gh` 不可用，可以使用等价 GitHub API release creation command。GitHub token 可以来自 ignored local environment file，但 token value 不得打印、提交或写入 release notes。

## 4. Authorized Release Surface

授权的 public surface 很小：

- Git tag `v0.1.0-alpha.1`；
- GitHub release record `v0.1.0-alpha.1`；
- 基于 repository-owned facts 写成的 release notes；
- GitHub 基于 tag 自动生成的 source archive。

README refresh 被授权，因为它是第一次 alpha 被外部用户理解和尝试所需的 documentation 与 developer-experience work。README 必须如实说明 alpha 状态，不得暗示 production readiness、package availability、hosted deployment availability、SDK availability、direct Nakama/Pitaya API compatibility，或尚未支持的 binary/container distribution。

## 5. Deferred Surface

以下内容继续 deferred：

- standalone binaries；
- packages；
- container images；
- checksum files；
- provenance files；
- signing artifacts；
- hosted deployments；
- install scripts；
- package 或 container registry publication；
- GitHub release record 之外的 public announcements；
- runtime behavior changes；
- protocol route changes；
- Protobuf source 或 generated output changes；
- migrations；
- dependencies；
- broad operations/admin behavior；
- authentication/session behavior changes；
- broad product module expansion；
- direct Nakama/Pitaya API compatibility。

## 6. Verification Requirements

创建 tag 或 GitHub release 前，仓库必须通过：

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

`v0.1.0-alpha.1` 的 GitHub release-record lookup 也必须在创建 release record 前返回 not found。

## 7. Known Warning Handling

本次 release 继续接受既有 repository warning：

```text
runtime.identity_boundary on runtime/internal/platform/persistence/postgres/authentication_repository.go
```

任何新 warning、failed check、version conflict，或 verification 后 target commit 改变，都会停止 release execution，直到重新验证。

## 8. Stop Conditions

如果发生以下情况，必须在 release execution 前停止：

- required checks 失败；
- 出现新 warning 且未被明确 triage；
- `v0.1.0-alpha.1` 已作为 local tag、remote tag、GitHub release record、package、image、hosted deployment 或 external artifact 存在；
- target commit 在 verification 后改变且未重新验证；
- release notes 或 README content 会包含未 redacted credentials、tokens、verifier keys、digests、DSNs、transport proof carriers、GitHub tokens 或 private environment file content；
- 尝试 binary、package、container、checksum、provenance file、signing artifact、hosted deployment、install script、registry publication，或 GitHub release record 之外的 public announcement；
- 本 release execution slice 内改变 runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior、broad product modules，或 direct Nakama/Pitaya API compatibility。

## 9. Next Work

Release execution 之后，下一个有边界的方向是：

```text
W-0196 Define first alpha user discovery loop
```

该后续工作应把 README 中更明确的尝试入口转化为具体 user-discovery loop，同时不改变 runtime behavior，也不扩大 release artifact surface。

## 10. Reference Alignment

Nakama 和 Pitaya 都会让 public release state 可见，同时把 runtime capability work 与 packaging/deployment posture 分开。本授权遵循这一纪律，用于 vibit 第一次 source-first alpha，并继续排除 direct Nakama/Pitaya API compatibility。
