# Release Publishing Decision Gate 中文版

状态：Draft v0.1
最后更新：2026-05-21
范围：v0.1 alpha path 的 gate-only release publishing decision boundary
依赖：`docs/v0.1-alpha-goal.md`、`docs/alpha-acceptance-checklist.md`、`docs/alpha-developer-flow.md`
权威决策：`ADR-0098`
说明：本文件是 `docs/release-publishing-decision-gate.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档定义 release publishing decision gate。它不发布 `v0.1 alpha`，不创建 release tags、binaries、archives、containers、packages、checksums、hosted deployments、runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、broad product modules 或 direct Nakama/Pitaya API compatibility。

## 1. 目的

Local alpha path 现在已经有 packaged developer journey、acceptance checklist、runtime runbook、request-loop script、status endpoints 和 focused authenticated gameplay proof。下一类风险是把“可以准备 release”误认为“已经执行 release publishing”。

本 gate 把这两个决定分开。

该 gate 记录后续 release execution step 可以开始准备前必须满足的条件。它也记录哪些 release artifacts 在后续明确 work item 授权前仍然 forbidden。

## 2. Core Rule

Release publishing decision gate 是：

```yaml
release_publishing_decision_gate: defined
completed_work_item: W-0190
decision: ADR-0098
check_rule: runtime.release_publishing_decision_gate
release_declared: false
release_publishing_authorized_by_this_gate: false
release_execution_authorized_by_this_gate: false
release_packaging_authorized_by_this_gate: false
release_artifacts_created_by_this_gate: false
hosted_deployment_authorized_by_this_gate: false
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
next_direction: release_execution_preparation_gate
```

该 gate 允许项目进入后续 release execution preparation gate。它不授权 release execution。

## 3. Publishing Prerequisites

只有在以下条件全部满足时，后续 release execution preparation step 才可以被考虑：

- `W-0189` 已经 packaged local alpha developer flow。
- `docs/alpha-acceptance-checklist.md` 记录 local alpha readiness、manual setup requirements、deferred work 和 release deferrals。
- `docs/alpha-developer-flow.md` 给出 coherent local developer journey。
- `docs/runtime-runbook.md` 描述 memory 和 PostgreSQL startup posture。
- `examples/local-alpha-request-loop.sh` 仍然是 authenticated gameplay E2E test 上的 redacted local proof。
- `/healthz`、`/readyz`、`/version` 和 `/configz` 仍是 local troubleshooting surfaces，且 `/configz` 保持 `secrets_redacted: true`。
- Required repository checks 和 Go tests 通过。
- 任何 known warning 都被明确记录，且不表示新的 release blocker。当前 known warning 仍是 `runtime.identity_boundary`，对应 `runtime/internal/platform/persistence/postgres/authentication_repository.go`。
- 面向公众的 release-path documents 都有英文权威文档和简体中文译文。
- Tracked artifacts 不记录 raw credential、raw access token、verifier key、digest、带 credentials 的 DSN、header、cookie、query string、WebSocket subprotocol value、remote address 或 GitHub token。
- Maintainer 明确授权后续 release execution preparation work item。

这些 prerequisites 对 release execution 是必要条件，但不是充分条件。

## 4. Decision Outcome

本 gate 的 decision outcome 是：

```yaml
release_decision_outcome:
  local_alpha_flow_packaged: true
  release_publishing_decision_gate_defined: true
  may_prepare_release_execution_gate_later: true
  may_publish_release_now: false
  may_create_release_artifacts_now: false
```

Repository 仍是 pre-alpha。Future release preparation 必须由 bounded work item 表达，并且必须继续把 release execution 和 release publication 分开，直到明确授权。

## 5. Release Artifact Boundary

只有在后续 work item 授权后，future release execution work 才可以讨论这些 artifact families：

- Git version tags。
- Release notes 或 changelog entries。
- Source archives。
- Checksums 或 provenance files。
- Optional binaries。
- Optional packages。
- Optional container images。
- Optional hosted deployments。

本 gate 不创建这些 artifacts。

本 gate 中 forbidden：

- 创建、签名或 push tags；
- 创建 binaries、archives、packages、containers、checksums 或 provenance files；
- 发布 GitHub releases 或等价 release records；
- 部署 hosted service；
- 为 release packaging 改变 runtime startup behavior；
- 改变 protocol routes 或 Protobuf sources；
- 为 release 原因重新生成 output；
- 添加 migrations 或 dependencies；
- 扩大 operations/admin behavior；
- 改变 authentication/session behavior；
- 添加 broad product modules；
- 选择 direct Nakama/Pitaya API compatibility。

## 6. Verification Requirements

把该 gate 视为 complete 前，使用这组命令：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-release-publishing-decision-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
```

Optional live PostgreSQL checks 仍通过 `VIBIT_POSTGRES_TEST_DSN` 和 disposable database opt-in。本 gate 不要求默认 repository checks 必须执行 live PostgreSQL verification。

## 7. Stop Conditions

如果发生以下任何情况，必须在后续 release execution preparation 前停止：

- required repository check 或 Go test 失败；
- 出现新的 warning 且没有明确 triage；
- 面向公众的英文文档缺少简体中文译文；
- 任何 tracked artifact 包含 raw credential material、raw access tokens、verifier key values、lookup digests、verifier digests、HMAC input/output bytes、带 credentials 的 DSNs、transport proof carriers 或 GitHub tokens；
- generated output 在没有 generation step 和 source trace 的情况下变化；
- Protobuf sources、generated output、migrations、dependencies 或 runtime behavior 在 release decision-only slice 中变化；
- release artifact boundary 被越过；
- 创建或暗示 hosted deployment；
- 添加 broad operations/admin behavior；
- authentication/session behavior 在没有独立授权 work item 时变化；
- 选择 direct Nakama/Pitaya API compatibility；
- maintainer 没有明确授权下一步 release preparation。

## 8. Redaction Expectations

Release decision records 必须保持可以安全提交到 repository。

不要包含：

- raw device credential text 或 bytes；
- raw access tokens；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- HMAC input 或 output bytes；
- verifier key values；
- concrete verifier key set ids；
- 带 credentials 的 PostgreSQL DSNs；
- 可能携带 secrets 的 headers、cookies、query strings、WebSocket subprotocol values 或 remote addresses；
- 包含 GitHub tokens 或其他 access credentials 的 local files。

## 9. Reference Alignment

Nakama 和 Pitaya 都体现了严肃 server framework 应具备 coherent local developer path 和明确 release support posture 的预期。本 gate 只用这个预期来塑造 release readiness discipline。

它不复制 Nakama 或 Pitaya APIs、release packaging、deployment model、cluster model、route names、data models、SDK surfaces 或 operations surfaces。

## 10. Next Work

下一步 bounded contribution 是：

```text
W-0191 Define release execution preparation gate
```

该 future step 可以准备 release execution plan，但仍不得发布 release 或创建 release artifacts，除非它自己的 scope 明确说明并且 maintainer 授权 ask-first boundary。
