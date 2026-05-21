# Prototype-Ready Foundation Execution Plan

状态：Accepted v0.1
最后更新：2026-05-22
范围：从 source-first alpha 走向 prototype-ready vibit foundation 的执行计划
依赖：`docs/product-maturity-milestones.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/first-alpha-feedback-intake-surfaces.md`、`docs/alpha-developer-flow.md`
权威决策：`ADR-0106`

配套英文源文档是 `docs/prototype-ready-foundation-execution-plan.md`。英文文件是权威版本。

本文定义 vibit 从 `v0.1.0-alpha.1` source-first alpha 走向 prototype-ready game/backend foundation 的第一份执行计划。它是 planning artifact。它不实现 runtime behavior、不添加 protocol routes、不添加 Protobuf source 或 generated output、不添加 migrations、不添加 dependencies、不扩大 operations/admin behavior、不添加 hosted deployments、不创建 release artifacts、不执行 public announcements、不进行 paid promotion、不改变 authentication/session behavior、不添加 broad product modules，也不添加 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Prototype-ready execution record 是：

```yaml
prototype_ready_foundation_execution_plan: defined
completed_work_item: W-0198
decision: ADR-0106
check_rule: runtime.prototype_ready_foundation_execution_plan
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
product_maturity_milestones_standard: docs/product-maturity-milestones.md
execution_plan_standard: docs/prototype-ready-foundation-execution-plan.md
execution_plan_standard_translation: docs/prototype-ready-foundation-execution-plan.zh-CN.md
recommended_sequence_recorded: true
candidate_work_items_recorded: true
maturity_stage_mapping_recorded: true
nakama_pitaya_capability_mapping_recorded: true
success_criteria_recorded: true
stop_conditions_recorded: true
selected_first_execution_slice: prototype_ready_local_development_path_gate
next_work_item: W-0199
next_direction: prototype_ready_local_development_path_gate
planning_only: true
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Interpretation

`v0.1.0-alpha.1` 证明 vibit 已经是一个真实的 source-first alpha：它有本地 authenticated request loop、PostgreSQL-backed state、WebSocket + Protobuf transport、local onboarding、protected inventory 和 presence paths、runtime sessions、connection binding、logout、runbooks、checks、release notes 和 feedback intake。

Prototype-ready 的意思比 production-ready 更窄，但更有实际价值：

- developer 可以从 repository 出发，构建一个严肃的小型多人或 realtime backend prototype；
- local setup、migration、configuration 和 example path 足够清晰，可以反复执行；
- 至少有一个 shared online-service capability 超出当前 proof slices；
- runtime lifecycle behavior 有足够 verification，让 prototype author 能信任核心 loop；
- 未完成的 production concerns 是显式的，而不是偶然遗漏。

下一阶段仍然是 source-first，也可以继续保持 single-node。它应该变得有用，而不仅仅是可检查。

## 3. Recommended Sequence

第一条 prototype-ready sequence 是：

1. 完成 local development path gate。
2. 在 gate 授权精确文件和行为后，实现 local development path 和更完整的 example flow。
3. 定义第一个超出 inventory 的 general storage-object behavior。
4. 实现最小 storage-object functional slice。
5. 定义第一个 server push、stream、broadcast 或 realtime messaging vocabulary。
6. 实现最小 realtime messaging 或 server-push functional slice。
7. 强化 login、protected requests、logout、reconnect 和 database failures 的 concurrency 与 failure-path verification。
8. 定义并添加严肃 prototype 使用前需要的最小 operations inspection surface。
9. 基于 external feedback 和 maturity criteria 执行 prototype-ready exit review。

这个顺序围绕 developer usefulness 排列。只有先有可重复的 local development path 和能覆盖多个请求的 example，storage、push、chat 或 matchmaking 等共享能力才更容易被评估。

## 4. First Execution Slice

第一个选中的 slice 是：

```text
W-0199 Define prototype-ready local development path gate
```

目的：定义一个可重复 local development path 的 gate，后续才能让 setup、migrations、configuration 和 example flow 对 prototype authors 可信。

该 gate 应记录：

- supported local prerequisites；
- startup 和 migration expectations；
- local secret/configuration expectations；
- example client 或 example app shape；
- 后续哪些 runtime、script、docs 或 test files 可以被修改；
- verification expectations；
- 在任何 behavior、protocol、migration、dependency、release 或 hosted deployment 扩展之前的 stop conditions。

本计划先选择 gate，是因为 developer trust 从可复现路径开始。如果用户不能可靠启动 server、运行 migrations、安全配置 secrets，并运行一个有意义的 example，那么 storage objects、push、chat 或 matchmaking 等更高层能力都很难评估。

## 5. Candidate Work Families

Prototype-ready execution plan 后续可以打开这些 bounded work families：

- `prototype_ready_local_development_path`：local setup、migration、configuration 和 example ergonomics。
- `storage_objects_and_durable_game_state`：inventory 之外的 general durable object behavior。
- `server_push_streams_or_realtime_messaging`：第一个 outbound realtime vocabulary。
- `failure_and_concurrency_verification`：lifecycle 和 database edge cases 的测试。
- `minimal_operations_inspection`：players、sessions、tokens、active connections 和 runtime state 的本地 inspection。
- `feedback_triage_to_work_items`：把真实 alpha feedback 转成 bounded work。

每个 family 仍必须遵循 vibit 的正常顺序：requirement、gate 或 spec、contract、必要时 generated shape、logic、tests、checks、docs 和 memory。

## 6. Maturity Mapping

本计划映射到 Stage 2 `prototype_ready_foundation`：

- Setup friction 由 local development path gate 和后续 implementation 处理。
- Example/client ergonomics 也由 local path 处理，因为 example 应展示更真实的 multi-request flow。
- General durable game state 由 storage-object family 处理。
- Realtime usefulness 由 server push、streams、broadcast 或 messaging 处理。
- 对现有 core loop 的信任由 failure 和 concurrency verification 处理。
- Operational confidence 由 minimal inspection surfaces 处理。
- User discovery 由现有 feedback intake loop 和后续 triage 处理。

本计划不宣称 Stage 3 single-node production-candidate readiness。Stage 3 仍需要更强的 security review、packaging、observability、upgrade posture、operational runbooks 和 failure-mode hardening。

## 7. Nakama/Pitaya Mapping

Nakama pressure：

- storage objects 和 durable game state；
- presence、status、notifications、chat 和 streams；
- SDK/example ergonomics；
- operational inspection；
- 清晰的 account/session lifecycle。

Pitaya pressure：

- connection lifecycle 和 session binding；
- handler route clarity；
- push、groups、broadcast 和 stream vocabulary；
- Go-first local runtime ergonomics；
- later frontend/backend 和 RPC topology，这些要等 single-process semantics 稳定后再推进。

本计划只把 Nakama 和 Pitaya 作为 capability baselines。它不添加 direct API compatibility、不复制 routes、不复制 schemas、不复制 clustering behavior，也不承诺 compatibility。

## 8. Success Criteria

Prototype-ready foundation track 成功的条件：

- technically capable external developer 可以从 source 出发完成可重复的 local prototype path；
- example flow 展示的不只是孤立的一次性 request；
- 至少一个 inventory 之外的 shared online-service capability 已存在，或被明确选为下一项 implementation；
- setup、configuration、migrations 和 secret redaction 对重复本地使用足够清晰；
- authenticated request loop 的 failure 和 concurrency risks 被 focused tests 覆盖，或被记录为 blockers；
- 至少一个真实 feedback item 被 triage 成 bounded work item 或显式 defer；
- 剩余 production gaps 是显式的，并被映射到 maturity stage。

## 9. Stop Conditions

如果执行计划需要以下事项，必须停止并请求 maintainer authorization：

- 在具体 implementation work item 授权前改变 runtime behavior；
- 改变 protocol routes 或 Protobuf source/generated output；
- 添加 migrations 或 repository/storage adapter changes；
- 添加 dependencies；
- 扩大 broad operations/admin behavior；
- 改变 authentication/session semantics；
- 扩大 broad product module；
- 添加 direct Nakama/Pitaya API compatibility；
- 添加 hosted deployments 或 demos；
- 添加 release binaries、packages、containers、checksums、signing/provenance artifacts、install scripts、registry publications 或 SDK packages；
- 执行 GitHub release record 之外的 public announcements；
- paid promotion；
- 处理或披露 secrets。

## 10. Next Work

下一项 bounded direction 是：

```text
W-0199 Define prototype-ready local development path gate
```

该工作应定义一个 gate，让 vibit 更容易在本地启动、配置、迁移和运行示例，这是走向 Stage 2 的第一个 product-useful step。

## 11. Verification

Repository 应验证：

- 本执行计划及其翻译存在；
- `ADR-0106` 记录该 decision；
- `.arch` manifests 标记 `W-0198` completed，并打开 `W-0199`；
- README、alpha goal、developer flow、acceptance checklist、AGENTS guides 和 product roadmap 指向新的 next work；
- runtime、protocol、generated output、migration、dependency、operations/admin、authentication/session、product module、hosted deployment、release artifact、public announcement、paid promotion 和 direct compatibility deferrals 都保持不变。

