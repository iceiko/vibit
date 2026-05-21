# Release Execution Preparation Gate 中文版

状态：Draft v0.1
最后更新：2026-05-21
范围：v0.1 alpha path 的 gate-only release execution preparation boundary
依赖：`docs/release-publishing-decision-gate.md`、`docs/alpha-developer-flow.md`、`docs/alpha-acceptance-checklist.md`、`docs/runtime-runbook.md`
权威决策：`ADR-0099`
说明：本文件是 `docs/release-execution-preparation-gate.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档定义 release execution preparation gate。它不发布 `v0.1 alpha`，不创建 release tags、binaries、archives、containers、packages、checksums、provenance files、hosted deployments、runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、broad product modules 或 direct Nakama/Pitaya API compatibility。

## 1. 目的

Release publishing decision gate 已经确认 local alpha readiness 可以走向 release preparation，但它没有授权 release execution。

本 gate 定义 future release execution plan 在任何 artifact creation 或 publication 被考虑前必须包含什么。它是 preparation boundary，不是 release run。

该 gate 记录 planning inputs、release-note boundaries、artifact plan boundaries、maintainer approval points、verification requirements、rollback notes 和 stop conditions。

## 2. Core Rule

Release execution preparation gate 是：

```yaml
release_execution_preparation_gate: defined
completed_work_item: W-0191
decision: ADR-0099
check_rule: runtime.release_execution_preparation_gate
release_declared: false
release_publishing_authorized_by_this_gate: false
release_execution_authorized_by_this_gate: false
release_packaging_authorized_by_this_gate: false
release_artifacts_created_by_this_gate: false
hosted_deployment_authorized_by_this_gate: false
release_execution_plan_defined: true
release_note_inputs_defined: true
artifact_plan_boundaries_defined: true
maintainer_approval_points_defined: true
rollback_stop_conditions_defined: true
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
next_direction: release_execution_authorization_gate
```

该 gate 允许项目进入后续 release execution authorization gate。它不授权 release execution 或 artifact creation。

## 3. Preparation Inputs

只有在 preparation record 能指向这些 inputs 后，后续 release execution authorization step 才可以被考虑：

- `docs/v0.1-alpha-goal.md`，用于 release target 和 non-goals。
- `docs/alpha-acceptance-checklist.md`，用于 local readiness、manual setup 和 deferrals。
- `docs/alpha-developer-flow.md`，用于面向 contributor 的 local flow。
- `docs/release-publishing-decision-gate.md`，用于 publication decision boundary。
- `docs/runtime-runbook.md`，用于 runtime startup、PostgreSQL、verifier key 和 redaction posture。
- `examples/local-alpha-request-loop.sh`，用于 redacted local proof command。
- `node tools/vibit inspect next` 和 repository checks，用于 machine-readable state。
- Current Git branch 和 commit range，供 future release note 总结。
- Known warning state，当前是 `runtime.identity_boundary`，对应 `runtime/internal/platform/persistence/postgres/authentication_repository.go`。
- Maintainer 确认下一步是 authorization review，而不是 execution。

本 gate 不选择 final release identifier，不创建 tag，不签名任何内容，也不创建 artifact manifest。

## 4. Release-Note Inputs

Future release notes 应只从 repository-owned facts 准备：

- current project status：pre-alpha moving toward `v0.1 alpha`；
- local developer flow summary；
- runtime surfaces：`/v1/ws`、`/healthz`、`/readyz`、`/version` 和 `/configz`；
- authenticated gameplay proof path：onboarding -> login -> bind connection -> protected inventory -> presence query -> logout -> revoked-token rejection；
- PostgreSQL manual setup posture；
- redaction expectations 和 secret handling；
- verification command set 和 results；
- known warnings 及其 triage status；
- explicit non-goals 和 deferrals；
- 通过 `.arch/work-items.yaml` 的 contribution entry point。

Future release notes 不得包含 raw credentials、raw access tokens、verifier keys、digests、带 credentials 的 DSNs、GitHub tokens、private environment files 或 unredacted transport proof carriers。

本 gate 不写 final release notes，也不把 changelog content 作为 release artifact 发布。

## 5. Artifact Plan Boundary

后续 release execution plan 可以讨论这些 artifact families：

- Git version tag。
- Release notes 或 changelog entry。
- Source archive。
- Checksum 或 provenance file。
- Optional binary build。
- Optional package。
- Optional container image。
- Optional hosted deployment。

本 gate 不创建它们。

任何 future artifact plan 都必须回答：

- 哪些 artifact families 在 scope 内；
- 哪些 artifact families 仍然 deferred；
- 每个 artifact 会由哪个 command 创建；
- 每个 artifact 会存放在哪里；
- creation 前必须通过哪些 verification；
- creation 前必须经过哪个 maintainer approval point；
- creation 失败时如何安全停止。

第一版 prepared alpha path 应优先选择最小 public release surface，以保持 source inspectability 和 contributor onboarding。Optional binaries、packages、containers 和 hosted deployments 仍然 deferred，除非后续 work item 明确授权。

## 6. Maintainer Approval Points

以下 actions 需要后续明确 maintainer authorization：

- 选择 final release identifier；
- 创建 Git tag；
- push Git tag；
- 创建 source archive；
- 创建 checksums 或 provenance files；
- 构建 binaries、packages 或 containers；
- 创建 hosted deployments；
- 创建 GitHub release 或等价 release record；
- 宣告 public release availability。

Authorization 应由 bounded work item 和 repository record 表达。仅聊天中的 approval 不足以构成 durable project state。

## 7. Verification Requirements

把该 preparation gate 视为 complete 前，使用这组命令：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-release-execution-preparation-gate --json
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

## 8. Rollback And Stop Conditions

因为本 gate 不创建 release artifacts，rollback 是 documentation 和 tooling revert。

如果发生以下任何情况，必须在后续 release execution authorization 前停止：

- required repository check 或 Go test 失败；
- 出现新的 warning 且没有明确 triage；
- known warning state 变化但没有 documented decision；
- 面向公众的英文文档缺少简体中文译文；
- 任何 tracked artifact 包含 raw credential material、raw access tokens、verifier key values、lookup digests、verifier digests、HMAC input/output bytes、带 credentials 的 DSNs、transport proof carriers 或 GitHub tokens；
- generated output 在没有 generation step 和 source trace 的情况下变化；
- Protobuf sources、generated output、migrations、dependencies 或 runtime behavior 在 preparation-only slice 中变化；
- 创建任何 release artifact；
- 创建、签名或 push 任何 release tag；
- hosted deployment work 开始；
- 添加 broad operations/admin behavior；
- authentication/session behavior 在没有独立授权 work item 时变化；
- 选择 direct Nakama/Pitaya API compatibility；
- maintainer 没有明确授权下一步 release execution authorization gate。

## 9. Redaction Expectations

Preparation records 必须保持可以安全提交到 repository。

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

## 10. Reference Alignment

Nakama 和 Pitaya 都体现了严肃 server framework 应分离 release readiness、artifacts、support posture 和 deployment choices。本 gate 采用这种纪律，但不采用它们的 APIs、data models、route names、release packaging、deployment model、cluster model、SDK surfaces 或 operations surfaces。

## 11. Next Work

`W-0192` 现在已经在 `docs/release-execution-authorization-gate.md` 中定义 release execution authorization gate，`W-0193 Confirm release execution maintainer decision` 已经在 `docs/release-execution-maintainer-decision.md` 中记录 maintainer go decision。

下一步 bounded contribution 是：

```text
W-0194 Define release identifier and artifact plan
```

该 future step 可以规划 release identifier 和 artifact/publication surface，但不得发布 release、创建 release tags、创建 release artifacts、创建 hosted deployments 或创建 GitHub release records，除非后续 execution scope 明确说明并且 maintainer 授权 ask-first boundary。
