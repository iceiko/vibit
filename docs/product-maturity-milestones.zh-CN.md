# Product Maturity Milestones 中文版

状态：Accepted v0.1
最后更新：2026-05-22
范围：vibit 从 source alpha 到 production-useful 产品阶段的成熟度路径
依赖：`docs/v0.1-alpha-goal.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/first-alpha-user-discovery-loop.md`
权威决策：`ADR-0105`
说明：本文件是 `docs/product-maturity-milestones.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文把维护者关于“必须推进到真正产品阶段、真正成为生产力”的意图沉淀成持久 milestones。当前 source-first alpha 已经证明 vibit 有真实后端闭环，但产品目标更大：从 first alpha 走到 prototype-ready foundation，再走到 single-node production-candidate foundation，最后走到 Nakama-first open-source server framework，并让 AI-native development 和 AI-native testing 成为默认用户体验。本文只提供 roadmap 和 feedback triage 指导。它不授权 runtime behavior changes、protocol route changes、Protobuf source 或 generated output changes、migrations、dependencies、hosted deployments、additional release artifacts、GitHub release record 之外的 public announcements、paid promotion、broad operations/admin behavior、authentication/session behavior changes、broad product module expansion、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Product maturity record 是：

```yaml
product_maturity_milestones: defined
completed_work_item: W-0197
decision: ADR-0105
check_rule: runtime.first_alpha_feedback_intake_surfaces
current_stage: source_first_alpha
current_release_identifier: v0.1.0-alpha.1
stage_1_source_first_alpha: reached
stage_2_prototype_ready_foundation: next_product_stage
stage_3_single_node_production_candidate_foundation: planned
stage_4_nakama_first_ai_native_product: long_term_target
stage_4_nakama_pitaya_class_product: long_term_target
reference_posture_update: ADR-0127
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
feedback_intake_surface: .github/ISSUE_TEMPLATE/alpha-feedback.yml
feedback_intake_standard: docs/first-alpha-feedback-intake-surfaces.md
next_direction: pilot_nakama_aligned_feature_request_workflow
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
direct_nakama_pitaya_api_compatibility_added: false
pitaya_distributed_architecture_added: false
```

## 2. Stage 1：Source-First Alpha

状态：

```text
stage_1_source_first_alpha: reached
release_identifier: v0.1.0-alpha.1
```

目的：证明 vibit 已经不只是设计练习。

当前 alpha 应被判断为 source-first developer alpha。它应该吸引能够在本地运行、检查架构并提供具体反馈的技术开发者。

Stage 1 存在的标准是开发者能够：

- clone repository；
- 运行 repository checks；
- 运行 Go tests；
- 运行 local alpha request loop；
- 理解 WebSocket + Protobuf + PostgreSQL posture；
- 理解 local onboarding、device credential login、access-token validation、runtime sessions、connection binding、protected inventory、protected presence query 和 logout；
- 在误把 alpha 当作 production server distribution 之前看到当前限制。

Stage 1 已经由 `v0.1.0-alpha.1` 达成。

Stage 1 不是：

- production deployment readiness；
- packaged distribution readiness；
- SDK readiness；
- hosted platform readiness；
- 与 Nakama、Pitaya、Colyseus、Pomelo、Agones 或 custom production backends 的 feature parity。
- 完整的用户可见 AI-native feature request、test、implementation 和 verification workflow。

## 3. Stage 2：Prototype-Ready Foundation

状态：

```text
stage_2_prototype_ready_foundation: next_product_stage
suggested_release_band: v0.2_or_v0.3
```

目的：让 vibit 可信地支持严肃 multiplayer 或 realtime backend prototype。

这个阶段意味着开发者可以用 vibit 作为小型 proof game 或 product prototype 的后端基础，而不是先自己补齐缺失的 common online-service layer。

Required capability groups：

- 超出 module-local inventory proof slice 的 storage objects 或类似 general durable game-state surface。
- 比当前最小 protected query 更清晰的 presence/status semantics。
- 第一版 server push、broadcast、stream 或 realtime messaging vocabulary。
- 一个 minimal chat 或 realtime messaging slice，或者明确记录另一个 shared online service 是更紧急的 prototype unlock。
- 更好的 local startup ergonomics，包括明确的 setup、migration 和 configuration flow。
- 一个 realistic example client 或 example app path，展示的不只是孤立请求。
- 能把 external friction 转换成 bounded work items 的 issue 和 feedback loops。
- AI-native requirement intake：把用户请求转换成 specs、acceptance criteria、test plans、tests、implementation boundaries、verification records 和 durable project memory。
- 针对 login、protected requests、logout、reconnect 和 database failure behavior 的基础 concurrency 和 failure-path verification。

Exit criteria：

- 技术能力足够的外部开发者可以基于 vibit 构建小型 prototype flow，而不是先改内部架构边界。
- 主要 setup friction 被记录为已接受或已降低。
- 下一个缺失能力是 product decision，而不是意外未知。
- 至少一个 non-maintainer feedback item 被 triaged 成 bounded work item，或被明确 deferred。
- 一个非平凡未来 feature 可以按已记录的 AI-native requirement 和 test workflow 推进，而不需要临时发明流程。

Stage 2 仍然不是 production-ready。它可以继续 single-node、local-first、source-first，但应该感觉有用，而不只是可检查。

## 4. Stage 3：Single-Node Production-Candidate Foundation

状态：

```text
stage_3_single_node_production_candidate_foundation: planned
suggested_release_band: v0.4_or_v0.5
```

目的：让 vibit 成为严肃 single-node foundation，使外部团队可以评估它是否适合真实项目开发。

这个阶段不意味着 broad distributed scalability 或完整 game-backend product parity。它意味着 single-node server path 已经有足够 security、operations、packaging 和 common backend capability，可以作为 production-candidate foundation。

Required capability groups：

- login、token validation、runtime sessions、connection binding、logout、close handoff、reconnect 和 session carriers 的稳定 lifecycle semantics。
- Hardened PostgreSQL path：migrations、indexes、transaction boundaries、connection handling、failure behavior 和 upgrade notes。
- Production configuration posture：secrets、redaction、validation、fail-closed behavior 和 environment separation。
- Observable operations baseline：structured logs、metrics 或 metrics boundary、health/readiness posture，以及 players、sessions、tokens、active connections 的 admin inspection path。
- 经明确授权后的 source archive 之外 release distribution：可能包括 container image、checksums、versioned release notes 和 upgrade documentation。
- Stable client ergonomics：SDK、client helper 或 documented protocol client example。
- Authentication、sessions、token redaction、configuration leakage 和 route permissions 的 security review。
- 适合 single-node usage 的 concurrency、soak 和 failure-mode verification。
- 可追溯到 requirement 和 acceptance artifacts 的 AI-generated 或 AI-maintained feature tests。

Exit criteria：

- 团队可以评估 vibit 用于真实 game backend project，而不是把 framework 当成研究 artifact。
- 已知 production risks 被记录为 accepted、fixed 或 blocking。
- Upgrade 和 operations expectations 是显式的。
- 距 product-class parity 的剩余差距主要是 breadth，而不是 core reliability posture。

## 5. Stage 4：Nakama-First AI-Native Product

状态：

```text
stage_4_nakama_first_ai_native_product: long_term_target
stage_4_nakama_pitaya_class_product: long_term_target
target: ai_era_nakama_pitaya_class_server_framework
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
```

目的：成为与 Nakama 同一广义产品级别的严肃开源 server framework，并围绕 vibit 的 agent-native maintainability model 重新适配。Pitaya-style distributed architecture 保持 future-only，直到后续 ADR 明确重新激活。

Product-class target 有两个不可分割的部分：

- Nakama-style backend capability coverage。
- AI-native development 和 testing，即用户需求通过 agent assistance 产出 specs、acceptance criteria、tests、implementation、verification 和 durable records。

这个阶段需要覆盖 `docs/nakama-pitaya-product-parity-roadmap.md` 中已经记录的 common capability families：

- identity、authentication、sessions；
- connection lifecycle、reconnect、logout；
- storage objects 和 durable game state；
- presence、status、notifications；
- chat、streams、realtime messaging；
- friends、groups、parties；
- leaderboards、tournaments、competitive systems；
- economy、inventory、rewards、currencies、progression；
- matchmaking、match listing、room lifecycle；
- realtime multiplayer 和 authoritative match runtime；
- server runtime hooks、RPC、custom logic；
- admin console、metrics、observability、operations；
- client SDKs、examples、developer experience；
- distributed runtime、frontend/backend roles、RPC、service discovery，仅在后续 architecture ADR 选择该路径之后；
- agent-native requirement、acceptance、test、implementation 和 verification workflow。

Stage 4 不意味着 direct Nakama/Pitaya API compatibility，除非后续 ADR 明确采纳 compatibility surface。

## 6. 按阶段分流反馈

早期用户反馈应映射到这些 maturity buckets 之一：

- `source_alpha_friction`：README、setup、checks、request-loop、runbook 或 concept clarity 阻塞。
- `prototype_ready_gap`：缺少 shared online service、example flow、client ergonomics 或 local development path，导致无法做 prototype。
- `production_candidate_gap`：security、operations、packaging、migration、observability、performance 或 failure behavior 阻塞真实项目评估。
- `product_class_gap`：social、competitive、matchmaking、match runtime、SDK、admin console、distributed runtime、extensibility breadth 或 AI-native feature workflow gaps 阻塞 Nakama-first product-class usefulness。
- `out_of_scope_for_now`：请求有效，但必须等待明确授权或后续阶段。

如果反馈要求 production claims、broad feature parity、hosted deployment、direct compatibility、binary/package/container publication、paid promotion 或 public announcement，不应默默接受，应转入 maintainer authorization。

## 7. 下一个产品方向

定义 agent-native feature request and test workflow 之后，下一个产品方向是：

```text
W-0221 Pilot Nakama-aligned feature request workflow
```

Agent-native feature request and test workflow 已记录在 `docs/agent-native-feature-request-test-workflow.md` 和 `ADR-0128`。下一项 work 应把该 workflow 试点应用到一个 bounded Nakama-aligned capability request，同时继续延后 runtime behavior、protocol changes、startup wiring、stream subscriptions、chat rooms、groups、broadcast fanout、delivery guarantees、persistence expansion、distributed runtime、matchmaking、match runtime、SDKs、hosted deployments、release artifacts、public announcements 和 direct compatibility，除非 pilot 打开后续明确 work item。

Stage 2 的追溯引用包括 `docs/prototype-ready-local-development-path-package.md`、`docs/storage-objects-behavior-gate.md`、`docs/storage-objects-persistence-schema-gate.md`、`runtime/migrations/postgres/000006_create_storage_objects.sql` 和 `docs/storage-objects-repository-boundary.md`。

上一项方向 `W-0198 Define prototype-ready foundation execution plan`、`W-0199 Define prototype-ready local development path gate`、`W-0200 Implement prototype-ready local development path package`、`W-0201 Define storage objects behavior gate`、`W-0202 Define storage objects persistence schema gate`、`W-0203 Add storage objects migration source`、`W-0204 Define storage objects repository boundary`、`W-0205 Implement storage-neutral storage objects repository interface`、`W-0206 Define storage objects PostgreSQL adapter gate`、`W-0207 Implement storage objects PostgreSQL adapter`、`W-0208 Define storage objects runtime behavior gate`、`W-0209 Implement storage objects runtime behavior`、`W-0210 Define storage objects protocol route gate`、`W-0211 Implement storage objects protocol route`、`W-0212 Prove storage objects protocol route in local alpha request flow`、`W-0213 Confirm next alpha direction after storage objects local proof`、`W-0214 Define first server push and realtime messaging gate` 和 `W-0215 Implement first server push and realtime messaging runtime slice` 已完成，并继续作为从 feedback intake 进入 Stage 2 execution plan 及其第一项 product capability path 的追溯记录。

已记录的候选重点现在是：

- 试点 Nakama-aligned feature request workflow；
- 加强现有 authenticated loop 的 concurrency 和 failure-path verification；
- 添加更清晰的 example client 或 example app path；
- 定义严肃 prototype 使用前需要的最小 operations inspection surface。

后续 work 必须保留 runtime behavior、protocol、generated output、dependencies、repository interfaces、storage adapters、operations/admin breadth、release artifacts、hosted deployments、public announcements、authentication/session behavior changes、broad product modules、large object/blob storage、S3-compatible object storage、Pitaya-style distributed architecture 和 direct Nakama/Pitaya compatibility 的 ask-first boundaries。

## 8. 非授权声明

本文记录产品成熟度目标。它本身不：

- 实现 runtime behavior；
- 添加 protocol routes；
- 添加 Protobuf source 或 generated output；
- 添加 migrations；
- 添加 dependencies；
- 发布 binaries、packages、containers、checksums、signing/provenance artifacts、install scripts、registry artifacts 或 hosted deployments；
- 授权 broad public announcements 或 paid promotion；
- 添加 operations/admin behavior；
- 改变 authentication/session behavior；
- 添加 broad product modules；
- 在后续 ADR 重新激活前添加 Pitaya-style distributed architecture；
- 添加 direct Nakama/Pitaya API compatibility；
- 声明当前 alpha production-ready。
