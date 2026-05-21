# First Alpha Feedback Intake Surfaces 中文版

状态：Accepted v0.1
最后更新：2026-05-22
范围：面向 early alpha user feedback 的 repository-owned intake surfaces
依赖：`docs/first-alpha-user-discovery-loop.md`、`docs/product-maturity-milestones.md`
权威决策：`ADR-0105`
说明：本文件是 `docs/first-alpha-feedback-intake-surfaces.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文为 `v0.1.0-alpha.1` 准备第一版 repository-owned feedback intake surface。它在任何 broader outreach 之前记录 intake format、feedback fields、redaction expectations、triage labels、product maturity mapping 和 stop conditions。它不执行 GitHub release record 之外的 public announcements，不运行 paid promotion，不添加 hosted deployments，不创建 release binaries、packages、containers、checksums、signing/provenance artifacts、install scripts、registry publications 或 SDK packages，不改变 runtime behavior，不添加 protocol routes 或 generated output，不添加 migrations 或 dependencies，不扩大 operations/admin behavior，不改变 authentication/session behavior，不添加 broad product modules，也不添加 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Feedback intake record 是：

```yaml
first_alpha_feedback_intake_surfaces: prepared
completed_work_item: W-0197
decision: ADR-0105
check_rule: runtime.first_alpha_feedback_intake_surfaces
release_identifier: v0.1.0-alpha.1
source_first_alpha_release_available: true
user_discovery_loop_defined: true
product_maturity_milestones_defined: true
feedback_intake_surface: .github/ISSUE_TEMPLATE/alpha-feedback.yml
feedback_intake_standard: docs/first-alpha-feedback-intake-surfaces.md
feedback_intake_standard_translation: docs/first-alpha-feedback-intake-surfaces.zh-CN.md
issue_template_added: true
feedback_fields_recorded: true
redaction_expectations_recorded: true
triage_labels_recorded: true
product_maturity_mapping_recorded: true
next_direction: prototype_ready_foundation_execution_plan
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
additional_release_artifacts_authorized: false
hosted_deployment_added: false
runtime_behavior_added: false
authentication_session_behavior_changed: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
product_module_expansion_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Intake Surface Choice

第一版 intake surface 是 GitHub issue form：

```text
.github/ISSUE_TEMPLATE/alpha-feedback.yml
```

理由：

- 它由 repository 拥有，对 contributors 可见。
- 不需要 hosted deployment、binary package、SDK、container 或 external service。
- 可以引导用户避免泄露 secrets。
- 可以把 feedback 直接映射到 product maturity stages。
- 给 maintainers 一个稳定入口，把 friction 转成 bounded work items。

GitHub Discussions 以后可以添加，但第一版 feedback surface 不强制需要。Discussion forum 对 open-ended comparison 或 design feedback 有用，但第一轮更需要 structured reports，而不是 conversation volume。

## 3. Intake Format

每个 alpha feedback issue 应捕捉：

- developer segment；
- attempted action；
- outcome；
- first blocker；
- product maturity bucket；
- 涉及的 commands 或 pages，必要时 redacted；
- expected next useful step；
- reporter 是否检查过 redaction guidance。

Issue template 默认不要求粘贴 logs。如果确实需要 logs，应要求 reporter 先移除 secrets。

## 4. Redaction Expectations

不要要求用户粘贴：

- raw device credentials；
- raw access tokens；
- verifier key values；
- 带 credentials 的 PostgreSQL DSNs；
- lookup digests 或 verifier digests；
- HMAC input 或 output；
- headers；
- cookies；
- query strings；
- WebSocket subprotocol values；
- remote addresses；
- 私有 `.env`、shell profile 或 local environment files；
- GitHub tokens 或任何其他 API tokens。

安全片段应缩短并 redacted。优先使用 command names、error class names、file paths 和 documentation page names，而不是完整 environment output。

## 5. Triage Labels

推荐 feedback issues 使用这些 labels：

```text
alpha-feedback
source-alpha-friction
prototype-ready-gap
production-candidate-gap
product-class-gap
redaction-review
needs-reproduction
candidate-work-item
deferred
```

这些 labels 已被记录，但本 work item 不要求通过 GitHub API 创建 repository labels。Label creation 可由 maintainer 或后续 repository-admin work item 完成。

## 6. Product Maturity Mapping

把 incoming feedback 映射到一个 maturity bucket：

- `source_alpha_friction`：reporter 不能理解、clone、check、test、运行 request loop，或不能跟随当前 local alpha path。
- `prototype_ready_gap`：reporter 能检查 alpha，但因为缺少 shared online-service、example、client、setup 或 local development capability，无法构建严肃 prototype。
- `production_candidate_gap`：reporter 被 packaging、security、operations、migration、observability、performance、reliability 或 failure-mode concerns 阻塞。
- `product_class_gap`：reporter 需要更广的产品能力，如 social systems、chat、leaderboards、matchmaking、match runtime、SDKs、admin console、extensibility 或 distributed runtime。
- `out_of_scope_for_now`：请求有效，但没有后续授权或更后 maturity stage 时不能处理。

这种映射避免 early feedback 被压平成无序 feature list。下一项计划方向是 `prototype_ready_foundation_execution_plan`，因此 prototype-blocking feedback 应获得特别关注，但不能假装当前 release 已 production-ready。

## 7. Triage Posture

默认 triage flow：

1. 确认 issue 不包含 secrets。
2. 分配 product maturity bucket。
3. 判断报告属于 setup friction、concept confusion、command failure、runtime gap、missing capability、trust gap、packaging gap 或 roadmap question。
4. 如果 report 包含具体 command failure，在本地复现。
5. 将 actionable findings 转成 bounded work item，或明确 defer。
6. 链接相关 doc、ADR、change spec 或 future work item。

Issue replies 中不要承诺 production readiness、direct Nakama/Pitaya compatibility、hosted demos、SDK availability、binary/container publication 或 broad feature parity，除非后续明确 work item 授权该 claim。

## 8. Stop Conditions

如果 feedback handling 需要以下事项，应停止并请求 maintainer authorization：

- GitHub release record 之外的 public announcements；
- paid promotion；
- hosted deployment 或 demo；
- release binaries、packages、containers、checksums、signing/provenance artifacts、install scripts、registry publications 或 SDK packages；
- runtime behavior changes；
- protocol route changes；
- Protobuf source 或 generated output changes；
- migrations；
- dependencies；
- broad operations/admin behavior；
- authentication/session behavior changes；
- broad product module expansion；
- direct Nakama/Pitaya API compatibility；
- secret disclosure 或 secret handling。

## 9. Next Work

下一个 bounded direction 是：

```text
W-0198 Define prototype-ready foundation execution plan
```

该 work 应认真使用 maturity milestones，并定义从 source-first alpha 到 prototype-ready foundation 的最小路径。若已有 early feedback，可以使用；但不应无限等待 feedback 才定义第一版 execution plan。

## 10. Verification

Repository 应验证：

- issue form 存在；
- 本 standard 和译本存在；
- product maturity milestones 已记录；
- `ADR-0105` 记录该 decision；
- `.arch` manifests 指向 `W-0198`；
- runtime、protocol、generated output、migration、dependency、release artifact、hosted deployment、announcement、paid promotion、authentication/session、broad product 和 direct compatibility deferrals 仍被保留。
