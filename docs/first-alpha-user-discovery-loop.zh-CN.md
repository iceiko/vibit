# First Alpha User Discovery Loop 中文版

状态：Accepted v0.1
最后更新：2026-05-22
范围：`v0.1.0-alpha.1` 之后寻找早期开发者的有边界 discovery loop
依赖：`docs/release-execution-final-authorization.md`、`docs/alpha-developer-flow.md`
权威决策：`ADR-0104`
说明：本文件是 `docs/first-alpha-user-discovery-loop.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文定义 source-first `v0.1.0-alpha.1` release 之后的第一轮 user-discovery loop。它记录 vibit 首先应向哪些开发者学习、可以准备哪些 surface、应捕获什么反馈、什么算有效信号，以及 loop 必须在哪里停止。它不授权 GitHub release record 之外的 public announcements、paid promotion、hosted deployments、release binaries、packages、containers、checksums、provenance files、signing artifacts、install scripts、registry publication、runtime behavior changes、protocol route changes、Protobuf source 或 generated output changes、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、broad product module expansion 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

User discovery record 是：

```yaml
first_alpha_user_discovery_loop: defined
completed_work_item: W-0196
decision: ADR-0104
check_rule: runtime.first_alpha_user_discovery_loop
release_identifier: v0.1.0-alpha.1
source_first_alpha_release_available: true
user_discovery_loop_defined: true
target_developer_segments_recorded: true
outreach_surfaces_recorded: true
feedback_capture_recorded: true
success_signals_recorded: true
stop_conditions_recorded: true
public_announcements_beyond_github_release_authorized_by_this_decision: false
paid_promotion_authorized_by_this_decision: false
hosted_deployment_authorized_by_this_decision: false
additional_release_artifacts_authorized_by_this_decision: false
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
next_direction: first_alpha_feedback_intake_surfaces
```

## 2. Discovery Goal

第一轮 loop 要回答一个问题：

```text
技术能力足够的开发者能否理解 vibit 为什么存在，在本地跑起 source alpha，并告诉我们哪些 friction 或缺失能力阻塞了下一步？
```

这不是增长 campaign。这个 loop 是一个学习系统，用少量真实开发者反应转化出具体 work items。

## 3. Target Developer Segments

从最可能感受到 vibit 所解决问题的开发者开始：

- 正在评估 server frameworks 的 Go game/backend developers。
- 用过或评估过 Nakama、Pitaya、Colyseus、Pomelo、Agones 或 custom backend stacks 的开发者。
- 在 backend codebases 中使用 AI coding agents，并遇到 unsafe 或 context-heavy changes 的工程师。
- 关心 explicit architecture、contracts、generated output、ADRs 和 repository checks 的 open-source contributors。
- 正在构建 multiplayer 或 realtime backend prototypes、且能接受 source-first alpha setup 的 technical founders 或 solo developers。

第一轮不要面向 production buyers、hosted-platform users、SDK-first developers，或需要 packaged binaries、container images、commercial support、direct Nakama/Pitaya API compatibility 或 production hardening 的团队。

## 4. Outreach Surfaces

本 work item 只定义 surfaces 和 copy posture，不执行 broad public announcements。

可以准备或指向的 surfaces：

- `README.md` 和 `README.zh-CN.md`。
- GitHub repository description 和 topics，如果后续明确 work item 授权编辑。
- `v0.1.0-alpha.1` 的 GitHub release page。
- `docs/releases/v0.1.0-alpha.1.md`。
- `docs/alpha-developer-flow.md`。
- `docs/alpha-acceptance-checklist.md`。
- 后续明确授权后的 GitHub issue 或 discussion intake surface。
- 后续明确授权后的小范围 maintainer-to-maintainer review requests，且不能 spam。

继续 deferred 的 surfaces：

- 在 social media、news aggregators、大型 community forums、newsletters 或 paid promotion 上 broad posts。
- Hosted demos。
- Binary、package、container、registry、checksum、signing 或 provenance publication。
- Production readiness、feature parity、SDK availability 或 direct compatibility 声明。

## 5. Feedback Capture

每一次 early-user conversation 或 issue 都应尽量捕获这些字段：

```yaml
feedback_record:
  source: github_issue_or_direct_review_or_discussion
  developer_segment: go_game_backend | nakama_pitaya_evaluator | ai_agent_backend_engineer | oss_contributor | prototype_builder | other
  attempted_action: read_readme | clone_repo | run_checks | run_go_tests | run_request_loop | inspect_architecture | inspect_runtime | evaluate_contribution
  outcome: completed | blocked | abandoned | unclear
  first_blocker: prerequisite | command_failure | docs_confusion | concept_confusion | runtime_gap | missing_feature | trust_gap | packaging_gap | other
  exact_command_or_page: redacted_if_needed
  redaction_checked: true
  next_action_candidate: docs_fix | check_fix | runbook_fix | issue_template | runtime_work_item | protocol_work_item | roadmap_clarification | defer
```

不要要求用户分享 raw device credentials、raw access tokens、verifier keys、带 credentials 的 DSNs、digests、headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 private environment files。

## 6. Review Questions

请求反馈时使用这些问题：

- 是什么让 vibit 值得你看一眼？
- 你能否在五分钟内理解 README？
- 你最先运行了哪条命令？结果是什么？
- `node tools/vibit check all`、`go test ./...` 或 `examples/local-alpha-request-loop.sh` 是否失败？
- Agent-native workflow 是否可理解，还是感觉像内部流程噪音？
- 缺少哪项能力会阻止你第二次尝试 vibit？
- 你愿意开 issue 或贡献一个小 slice 吗？如果不愿意，阻塞点是什么？

## 7. Success Signals

出现以下至少一项时，第一轮 loop 就算成功：

- 外部开发者 clone repository 并反馈结果。
- 非 maintainer 提出有用 issue 或 discussion。
- Maintainer 收到与 README、runbook、request-loop、setup、architecture clarity 或 next capability 相关的具体反馈。
- Try path 中的失败可复现，并被转化为有边界 work item。
- 开发者读完 README 后能准确解释 vibit 的 agent-native server framework premise。

次级信号：

- GitHub stars、forks 或 release views。
- 关于 game backend roadmap scope 的 inbound questions。
- 针对 Nakama、Pitaya、Colyseus、Pomelo、Agones 或 custom Go backends 的具体 comparison feedback。

不要只把 vanity metrics 当作成功。

## 8. Stop Conditions

如果出现以下情况，停止 loop 或请求 maintainer authorization：

- 将要发布 GitHub release record 之外的 public announcement。
- 考虑 paid promotion。
- 提出 hosted deployment 或 demo。
- 本 loop 中请求 binary、package、container、checksum、signing/provenance artifact、install script、registry publication 或 SDK package。
- Feedback 需要 runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、operations/admin expansion、authentication/session behavior、broad product module expansion 或 direct Nakama/Pitaya API compatibility。
- 用户报告 security 或 secret-redaction 风险。
- README、release notes、runbook 或 request-loop instructions 存在实质错误。
- Loop 开始产生重复低信号 outreach，而不是具体学习。

## 9. Next Work

下一项有边界方向是：

```text
W-0197 Prepare first alpha feedback intake surfaces
```

该 work 应准备一个最小的、repository-owned 的 early-user feedback 接收位置，例如 issue/discussion guidance 或 labels；但不执行 broad announcements、不添加 hosted deployment，也不改变 runtime behavior。

## 10. Reference Alignment

Nakama 和 Pitaya 都依靠公开 repository 和 documentation 周围清晰的 contributor/user intake surfaces。vibit 应学习这种姿态，但第一轮要更小：source-first、诚实说明 alpha limitations，并聚焦 agent-native maintainability feedback，而不是 broad feature parity claims。
