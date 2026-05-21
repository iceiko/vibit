# Release Execution Maintainer Decision 中文版

状态：Draft v0.1
最后更新：2026-05-21
范围：v0.1 alpha release execution path 的 maintainer go/no-go decision record
依赖：`docs/release-execution-authorization-gate.md`
Canonical decision: `ADR-0101`
说明：本文件是 `docs/release-execution-maintainer-decision.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文记录 `W-0193` 的 maintainer decision。它不发布 `v0.1 alpha`，不选择或创建最终 release identifier 或 tag，不创建 binaries、archives、containers、packages、checksums、provenance files、hosted deployments、GitHub release records、runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、broad product modules 或 direct Nakama/Pitaya API compatibility。

## 1. 目的

`W-0192` 有意停在 blocked maintainer decision gate。Maintainer 现在已经说明 decision outcome 是同意继续 release execution path。

这个 approval 是窄范围的。它授权仓库从 blocked go/no-go gate 进入一个有边界的 release identifier 和 artifact planning step。它不授权立即执行 release execution commands。

## 2. Decision Record

Maintainer decision 是：

```yaml
release_execution_maintainer_decision: recorded
completed_work_item: W-0193
decision: ADR-0101
check_rule: runtime.release_execution_maintainer_decision
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
next_direction: release_identifier_artifact_plan
```

这个 decision 只表示可以进入下一个 planning step。它不是任何会产出 release 的命令的 go。

## 3. Maintainer Input

Maintainer 在 conversation 中说明：

```text
我没看懂，你说需要我决策什么？我决策的结果为同意。
```

本记录把 `同意` 解释为：

- go，继续 release execution path；
- 尚未批准具体 release identifier；
- 尚未批准 tag；
- 尚未批准 artifact family；
- 尚未批准 publication surface；
- 尚未批准 release execution command。

这个保守解释保留 release authorization boundary，同时允许 queue 继续向前。

## 4. Required Next Step

下一个 bounded work item 是：

```text
W-0194 Define release identifier and artifact plan
```

该 work item 可以提议或定义 release identifier 和 artifact plan，但仍不得创建 tag、GitHub release、binary、archive、container、package、checksum、provenance file、hosted deployment 或 public release announcement，除非后续 work item 明确授权 execution。

## 5. Verification State

这个 decision record 应使用以下命令检查：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change confirm-release-execution-maintainer-decision --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

这个 documentation-only decision record 不要求运行 Go runtime tests；如果后续 release planning 或 execution step 改动 runtime behavior，再运行对应 Go tests。本记录不修改 Go runtime code。

## 6. Known Warning Handling

已知 repository warning 仍是：

```text
runtime.identity_boundary on runtime/internal/platform/persistence/postgres/authentication_repository.go
```

为了继续 release planning，本 decision 接受这个已知 warning。任何新 warning 都必须在 release execution 继续前 triage。

## 7. Stop Conditions

如果出现以下情况，在 release execution 前停止：

- required checks 失败；
- 出现新 warning 且未明确 triage；
- 没有 repository record 就选择最终 release identifier；
- 没有后续 execution authorization 就尝试 tag、GitHub release、artifact、checksum、provenance file、hosted deployment 或 publication command；
- 任何 tracked artifact 包含 raw credentials、raw access tokens、verifier keys、lookup digests、verifier digests、HMAC input/output bytes、带 credential 的 DSNs、transport proof carriers、GitHub tokens 或 private environment file content；
- 在本 decision slice 内修改 runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior、broad product modules 或 direct Nakama/Pitaya API compatibility。

## 8. Reference Alignment

Nakama 和 Pitaya 都把 release decisions 与 artifact creation、deployment mechanics 分开。本 decision 沿用这种纪律，同时保留 vibit 的 agent-native workflow，并且不采用 direct Nakama/Pitaya APIs、packaging models、deployment models 或 compatibility promises。

